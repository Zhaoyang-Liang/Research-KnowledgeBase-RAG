package main

import (
	"flag"
	"fmt"
	"math"
	"math/big"
	"math/cmplx"

	"github.com/tuneinsight/lattigo/v6/circuits/ckks/bootstrapping"
	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/ring"
	"github.com/tuneinsight/lattigo/v6/schemes/ckks"
	"github.com/tuneinsight/lattigo/v6/utils"
)

// Custom coefficient/RNS ModRaise probe.
//
// Motivation:
//   The previous probes showed:
//     (1) Pure coefficient-wise RNS operations commute with coefficient/module split.
//     (2) Lattigo bootstrapping.Evaluator.ModUp and ScaleDown+ModUp cannot simply be
//         moved below RingSwitchDown/Split.
//
// This probe tests a lower-level custom ModRaise that only does centered coefficient
// lifting from level 0 to a chosen target level, optionally followed by a coefficient
// multiplier 2^mulPow2. It deliberately does NOT call bootstrapping.Evaluator.ModUp
// for the experimental path.
//
// It compares:
//   A  = Split( LattigoModUp_N(ScaleDown_N(ct)) )              [reference]
//   T  = Split( CustomModRaise_N(ScaleDown_N(ct)) )            [custom top then split]
//   L  = CustomModRaise_n( Split(ScaleDown_N(ct)) )            [split then custom leaf]
//
// Main readout:
//   - T vs L small  => custom coefficient/RNS ModRaise commutes with RingSwitchDown/Split
//                      up to low-vs-high key-switching noise.
//   - A vs T small  => custom ModRaise matches Lattigo ModUp convention.
//   - A vs L small  => the desired scheduling primitive works using custom ModRaise.
//
// Suggested runs:
//   go run lattigo_custom_modraise_probe.go -logN=14 -layers=2 -trials=3 -mulPow2=0
//   go run lattigo_custom_modraise_probe.go -logN=14 -layers=2 -trials=3 -mulPow2=5
//   go run lattigo_custom_modraise_probe.go -logN=14 -layers=2 -trials=3 -scanPow2=true -powMin=0 -powMax=12

type diag struct {
	maxErr  float64
	rmse    float64
	mean    float64
	gamma   complex128
	fitMax  float64
	fitRMSE float64
}

func makeResidualParams(logN int) ckks.Parameters {
	params, err := ckks.NewParametersFromLiteral(ckks.ParametersLiteral{
		LogN:            logN,
		LogQ:            []int{55, 40, 40, 40, 40, 40, 40, 40},
		LogP:            []int{61, 61, 61},
		LogDefaultScale: 40,
		Xs:              ring.Ternary{H: 192},
	})
	if err != nil {
		panic(err)
	}
	return params
}

func makeBootstrappingParams(logN int, smallRingAdjust bool) (ckks.Parameters, bootstrapping.Parameters) {
	params := makeResidualParams(logN)
	lit := bootstrapping.ParametersLiteral{
		LogN: utils.Pointy(logN),
		LogP: []int{61, 61, 61, 61},
		Xs:   params.Xs(),
	}
	btp, err := bootstrapping.NewParametersFromLiteral(params, lit)
	if err != nil {
		panic(err)
	}
	if smallRingAdjust && logN < 16 {
		btp.Mod1ParametersLiteral.LogMessageRatio += 16 - params.LogN()
	}
	return params, btp
}

func copyQ(params ckks.Parameters) []uint64 {
	rq := params.RingQ()
	q := make([]uint64, len(rq.SubRings))
	for i := range rq.SubRings {
		q[i] = rq.SubRings[i].Modulus
	}
	return q
}

func copyP(params ckks.Parameters) []uint64 {
	rp := params.RingP()
	p := make([]uint64, len(rp.SubRings))
	for i := range rp.SubRings {
		p[i] = rp.SubRings[i].Modulus
	}
	return p
}

func makeParamsFromExactBasis(template ckks.Parameters, logN int) ckks.Parameters {
	params, err := ckks.NewParametersFromLiteral(ckks.ParametersLiteral{
		LogN:            logN,
		Q:               copyQ(template),
		P:               copyP(template),
		LogDefaultScale: template.LogDefaultScale(),
		Xs:              template.Xs(),
	})
	if err != nil {
		panic(fmt.Errorf("makeParamsFromExactBasis failed: %w", err))
	}
	return params
}

func sameRingBasis(a, b ckks.Parameters) bool {
	if len(a.RingQ().SubRings) != len(b.RingQ().SubRings) || len(a.RingP().SubRings) != len(b.RingP().SubRings) {
		return false
	}
	for i := range a.RingQ().SubRings {
		if a.RingQ().SubRings[i].Modulus != b.RingQ().SubRings[i].Modulus {
			return false
		}
	}
	for i := range a.RingP().SubRings {
		if a.RingP().SubRings[i].Modulus != b.RingP().SubRings[i].Modulus {
			return false
		}
	}
	return true
}

func makeTopEvaluator(logN int, sk *rlwe.SecretKey, smallAdjust bool) (ckks.Parameters, bootstrapping.Parameters, *bootstrapping.Evaluator) {
	_, btp := makeBootstrappingParams(logN, smallAdjust)
	keys, _, err := btp.GenEvaluationKeys(sk)
	if err != nil {
		panic(fmt.Errorf("top GenEvaluationKeys failed: %w", err))
	}
	ev, err := bootstrapping.NewEvaluator(btp, keys)
	if err != nil {
		panic(fmt.Errorf("top NewEvaluator failed: %w", err))
	}
	return btp.BootstrappingParameters, btp, ev
}

func makeValues(slots int, trial int) []complex128 {
	vals := make([]complex128, slots)
	off := trial + 1
	for i := range vals {
		re := float64(((i+3*off)%17)-8) / 16.0
		im := float64(((2*i+5*off)%19)-9) / 32.0
		vals[i] = complex(re, im)
	}
	return vals
}

func encryptAtLevel(params ckks.Parameters, encoder *ckks.Encoder, encryptor *rlwe.Encryptor, values []complex128, level int) *rlwe.Ciphertext {
	pt := ckks.NewPlaintext(params, level)
	if err := encoder.Encode(values, pt); err != nil {
		panic(err)
	}
	ct, err := encryptor.EncryptNew(pt)
	if err != nil {
		panic(err)
	}
	return ct
}

func decode(params ckks.Parameters, encoder *ckks.Encoder, decryptor *rlwe.Decryptor, ct *rlwe.Ciphertext) []complex128 {
	pt := decryptor.DecryptNew(ct)
	vals := make([]complex128, params.MaxSlots())
	if err := encoder.Decode(pt, vals); err != nil {
		panic(err)
	}
	return vals
}

func printCT(name string, ct *rlwe.Ciphertext) {
	fmt.Printf("%-42s LogN=%d Level=%d Degree=%d IsNTT=%v IsMont=%v Scale=%v\n",
		name, ct.LogN(), ct.Level(), ct.Degree(), ct.MetaData.IsNTT, ct.MetaData.IsMontgomery, ct.Scale)
}

func splitTree(eval *rlwe.RingPackingEvaluator, ct *rlwe.Ciphertext, layers int) ([]*rlwe.Ciphertext, error) {
	nodes := []*rlwe.Ciphertext{ct}
	for l := 0; l < layers; l++ {
		next := make([]*rlwe.Ciphertext, 0, 2*len(nodes))
		for _, node := range nodes {
			even, odd, err := eval.SplitNew(node)
			if err != nil {
				return nil, fmt.Errorf("SplitNew layer=%d logN=%d: %w", l, node.LogN(), err)
			}
			next = append(next, even, odd)
		}
		nodes = next
	}
	return nodes, nil
}

func mergeTree(eval *rlwe.RingPackingEvaluator, leaves []*rlwe.Ciphertext) (*rlwe.Ciphertext, error) {
	nodes := make([]*rlwe.Ciphertext, len(leaves))
	for i := range leaves {
		nodes[i] = leaves[i].CopyNew()
	}
	for len(nodes) > 1 {
		if len(nodes)%2 != 0 {
			return nil, fmt.Errorf("odd number of nodes: %d", len(nodes))
		}
		next := make([]*rlwe.Ciphertext, 0, len(nodes)/2)
		for i := 0; i < len(nodes); i += 2 {
			parent, err := eval.MergeNew(nodes[i], nodes[i+1])
			if err != nil {
				return nil, fmt.Errorf("MergeNew logN=%d: %w", nodes[i].LogN(), err)
			}
			next = append(next, parent)
		}
		nodes = next
	}
	return nodes[0], nil
}

func prepareBootstrapInput(eval *bootstrapping.Evaluator, ct *rlwe.Ciphertext) (*rlwe.Ciphertext, *rlwe.Ciphertext, error) {
	ctScaled, _, err := eval.ScaleDown(ct.CopyNew())
	if err != nil {
		return nil, nil, fmt.Errorf("ScaleDown failed: %w", err)
	}
	ctBoot, err := eval.ModUp(ctScaled.CopyNew())
	if err != nil {
		return nil, nil, fmt.Errorf("ModUp failed: %w", err)
	}
	return ctScaled, ctBoot, nil
}

func modBigSignedPow2FromResidue(x, q0 uint64, pow2 int, qi uint64) uint64 {
	var z big.Int
	if x <= q0/2 {
		z.SetUint64(x)
	} else {
		var tmp big.Int
		tmp.SetUint64(q0 - x)
		z.Neg(&tmp)
	}
	if pow2 > 0 {
		z.Lsh(&z, uint(pow2))
	}
	var mod big.Int
	mod.SetUint64(qi)
	z.Mod(&z, &mod)
	if z.Sign() < 0 {
		z.Add(&z, &mod)
	}
	return z.Uint64()
}

// customCenteredModRaise lifts a level-0 ciphertext to targetLevel by centered
// coefficient lifting from q0 to every qi. If the input is in NTT domain, it first
// applies INTT at level 0 and transforms the output back to NTT at targetLevel.
//
// It is intentionally simple and slow; this is a correctness probe, not production code.
func customCenteredModRaise(params ckks.Parameters, ctIn *rlwe.Ciphertext, targetLevel int, scaleSource *rlwe.Ciphertext, mulPow2 int) (*rlwe.Ciphertext, error) {
	if ctIn.Level() != 0 {
		return nil, fmt.Errorf("customCenteredModRaise expects level-0 input, got level=%d", ctIn.Level())
	}
	if targetLevel < 0 || targetLevel > params.MaxLevel() {
		return nil, fmt.Errorf("invalid targetLevel=%d max=%d", targetLevel, params.MaxLevel())
	}
	if ctIn.MetaData.IsMontgomery {
		return nil, fmt.Errorf("input is in Montgomery domain; this probe only handles standard residues")
	}

	rq := params.RingQ()
	rq0 := rq.AtLevel(0)
	rqT := rq.AtLevel(targetLevel)
	q0 := rq.SubRings[0].Modulus

	ctOut := rlwe.NewCiphertext(params, ctIn.Degree(), targetLevel)
	ctOut.MetaData = ctIn.MetaData.CopyNew()
	ctOut.MetaData.IsNTT = ctIn.MetaData.IsNTT
	ctOut.MetaData.IsMontgomery = false
	ctOut.Scale = ctIn.Scale
	if scaleSource != nil {
		ctOut.Scale = scaleSource.Scale
	}

	for d := 0; d <= ctIn.Degree(); d++ {
		coeffPoly := rq0.NewPoly()
		if ctIn.MetaData.IsNTT {
			rq0.INTT(ctIn.Value[d], coeffPoly)
		} else {
			coeffPoly.CopyLvl(0, ctIn.Value[d])
		}

		for j := 0; j <= targetLevel; j++ {
			qj := rq.SubRings[j].Modulus
			in0 := coeffPoly.Coeffs[0]
			outj := ctOut.Value[d].Coeffs[j]
			for i := range in0 {
				outj[i] = modBigSignedPow2FromResidue(in0[i]%q0, q0, mulPow2, qj)
			}
		}

		if ctIn.MetaData.IsNTT {
			rqT.NTT(ctOut.Value[d], ctOut.Value[d])
		}
	}

	return ctOut, nil
}

func customModRaiseLeaves(paramsLeaf ckks.Parameters, leaves []*rlwe.Ciphertext, targetLevel int, scaleSource *rlwe.Ciphertext, mulPow2 int) ([]*rlwe.Ciphertext, error) {
	out := make([]*rlwe.Ciphertext, len(leaves))
	for i := range leaves {
		ct, err := customCenteredModRaise(paramsLeaf, leaves[i], targetLevel, scaleSource, mulPow2)
		if err != nil {
			return nil, fmt.Errorf("leaf[%d] custom ModRaise: %w", i, err)
		}
		out[i] = ct
	}
	return out, nil
}

func bad(z complex128) bool {
	return math.IsNaN(real(z)) || math.IsNaN(imag(z)) || math.IsInf(real(z), 0) || math.IsInf(imag(z), 0)
}

func diffStats(want, got []complex128) (maxErr float64, rmse float64, mean float64) {
	if len(want) != len(got) {
		panic(fmt.Errorf("length mismatch: want=%d got=%d", len(want), len(got)))
	}
	var sum, sumSq float64
	for i := range want {
		if bad(want[i]) || bad(got[i]) {
			return math.Inf(1), math.Inf(1), math.Inf(1)
		}
		d := cmplx.Abs(want[i] - got[i])
		if d > maxErr {
			maxErr = d
		}
		sum += d
		sumSq += d * d
	}
	if len(want) > 0 {
		mean = sum / float64(len(want))
		rmse = math.Sqrt(sumSq / float64(len(want)))
	}
	return
}

func estimateGamma(want, got []complex128) complex128 {
	var num complex128
	var den float64
	for i := range want {
		if bad(want[i]) || bad(got[i]) {
			return complex(math.NaN(), math.NaN())
		}
		num += got[i] * cmplx.Conj(want[i])
		den += cmplx.Abs(want[i]) * cmplx.Abs(want[i])
	}
	if den == 0 {
		return 0
	}
	return num / complex(den, 0)
}

func divideVector(in []complex128, gamma complex128) []complex128 {
	out := make([]complex128, len(in))
	if gamma == 0 || bad(gamma) {
		for i := range out {
			out[i] = complex(math.NaN(), math.NaN())
		}
		return out
	}
	for i := range in {
		out[i] = in[i] / gamma
	}
	return out
}

func compare(name string, want, got []complex128) diag {
	maxErr, rmse, mean := diffStats(want, got)
	gamma := estimateGamma(want, got)
	fitGot := divideVector(got, gamma)
	fitMax, fitRMSE, _ := diffStats(want, fitGot)
	fmt.Printf("%-58s rawMax=%10.3e rmse=%10.3e gamma=%+.6e%+.6ei fitMax=%10.3e fitRMSE=%10.3e\n",
		name, maxErr, rmse, real(gamma), imag(gamma), fitMax, fitRMSE)
	return diag{maxErr: maxErr, rmse: rmse, mean: mean, gamma: gamma, fitMax: fitMax, fitRMSE: fitRMSE}
}

func updateWorst(cur, x float64) float64 {
	if x > cur {
		return x
	}
	return cur
}

func main() {
	logN := flag.Int("logN", 14, "top ring degree log2(N)")
	layers := flag.Int("layers", 2, "recursive split layers; k=2^layers")
	trials := flag.Int("trials", 3, "number of trials")
	smallAdjust := flag.Bool("topSmallAdjust", true, "apply small-logN LogMessageRatio adjustment to the top evaluator")
	mulPow2 := flag.Int("mulPow2", 0, "multiply lifted signed coefficients by 2^mulPow2 before reducing to target moduli")
	scanPow2 := flag.Bool("scanPow2", false, "scan mulPow2 range and summarize best A-vs-L/T-vs-L errors")
	powMin := flag.Int("powMin", 0, "minimum mulPow2 when scanPow2=true")
	powMax := flag.Int("powMax", 12, "maximum mulPow2 when scanPow2=true")
	flag.Parse()

	if *layers <= 0 {
		panic("layers must be positive")
	}
	leafLogN := *logN - *layers
	if leafLogN < 1 {
		panic("layers too large")
	}
	k := 1 << uint(*layers)

	fmt.Println("=== Custom coefficient/RNS ModRaise scheduling probe ===")
	fmt.Printf("top LogN=%d, layers=%d, k=%d, leaf LogN=%d, trials=%d\n", *logN, *layers, k, leafLogN, *trials)
	fmt.Printf("topSmallAdjust=%v, mulPow2=%d, scanPow2=%v\n", *smallAdjust, *mulPow2, *scanPow2)
	fmt.Println("A = Split(Lattigo ModUp_N(ScaleDown_N(ct)))")
	fmt.Println("T = Split(Custom ModRaise_N(ScaleDown_N(ct)))")
	fmt.Println("L = Custom ModRaise_n(Split(ScaleDown_N(ct)))")

	_, btpTopForParams := makeBootstrappingParams(*logN, *smallAdjust)
	paramsTop := btpTopForParams.BootstrappingParameters
	paramsLeaf := makeParamsFromExactBasis(paramsTop, leafLogN)

	fmt.Println("Generating top keys and evaluators...")
	kgenTop := rlwe.NewKeyGenerator(paramsTop)
	skTop, pkTop := kgenTop.GenKeyPairNew()
	encoderTop := ckks.NewEncoder(paramsTop)
	encryptorTop := rlwe.NewEncryptor(paramsTop, pkTop)
	decryptorTop := rlwe.NewDecryptor(paramsTop, skTop)

	paramsTop, btpTop, topEval := makeTopEvaluator(*logN, skTop, *smallAdjust)
	_ = btpTop

	fmt.Println("Generating ring switching keys...")
	rpk := &rlwe.RingPackingEvaluationKey{}
	skMap, err := rpk.GenRingSwitchingKeys(paramsTop, skTop, leafLogN, rlwe.EvaluationKeyParameters{})
	if err != nil {
		panic(fmt.Errorf("GenRingSwitchingKeys failed: %w", err))
	}
	rpEval := rlwe.NewRingPackingEvaluator(rpk)
	skLeaf := skMap[leafLogN]
	if skLeaf == nil {
		panic("missing leaf secret key")
	}
	encoderLeaf := ckks.NewEncoder(paramsLeaf)
	decryptorLeaf := rlwe.NewDecryptor(paramsLeaf, skLeaf)

	fmt.Printf("top params : LogN=%d MaxLevel=%d LogQP=%.1f Slots=%d exactBasis=true\n", paramsTop.LogN(), paramsTop.MaxLevel(), paramsTop.LogQP(), paramsTop.MaxSlots())
	fmt.Printf("leaf params: LogN=%d MaxLevel=%d LogQP=%.1f Slots=%d exactBasis=%v\n", paramsLeaf.LogN(), paramsLeaf.MaxLevel(), paramsLeaf.LogQP(), paramsLeaf.MaxSlots(), sameRingBasis(paramsTop, paramsLeaf))

	powList := []int{*mulPow2}
	if *scanPow2 {
		powList = nil
		for p := *powMin; p <= *powMax; p++ {
			powList = append(powList, p)
		}
	}

	type summary struct {
		pow        int
		worstAvsT  float64
		worstAvsL  float64
		worstTvsL  float64
		worstMAvsT float64
		worstMAvsL float64
		worstMTvsL float64
	}
	summaries := make([]summary, 0, len(powList))

	for _, pow := range powList {
		fmt.Println()
		fmt.Printf("=== Testing custom ModRaise with mulPow2=%d ===\n", pow)

		var worstAvsT, worstAvsL, worstTvsL float64
		var worstMAvsT, worstMAvsL, worstMTvsL float64

		for trial := 0; trial < *trials; trial++ {
			values := makeValues(paramsTop.MaxSlots(), trial)
			ctLow := encryptAtLevel(paramsTop, encoderTop, encryptorTop, values, 0)

			topScaled, topModUp, err := prepareBootstrapInput(topEval, ctLow.CopyNew())
			if err != nil {
				panic(err)
			}

			A, err := splitTree(rpEval, topModUp.CopyNew(), *layers)
			if err != nil {
				panic(fmt.Errorf("A split failed: %w", err))
			}

			customTop, err := customCenteredModRaise(paramsTop, topScaled.CopyNew(), paramsTop.MaxLevel(), topModUp, pow)
			if err != nil {
				panic(fmt.Errorf("custom top ModRaise failed: %w", err))
			}
			T, err := splitTree(rpEval, customTop.CopyNew(), *layers)
			if err != nil {
				panic(fmt.Errorf("T split failed: %w", err))
			}

			scaledLeaves, err := splitTree(rpEval, topScaled.CopyNew(), *layers)
			if err != nil {
				panic(fmt.Errorf("split topScaled failed: %w", err))
			}
			L, err := customModRaiseLeaves(paramsLeaf, scaledLeaves, paramsLeaf.MaxLevel(), topModUp, pow)
			if err != nil {
				panic(fmt.Errorf("custom leaf ModRaise failed: %w", err))
			}

			if trial == 0 && !*scanPow2 {
				fmt.Println("\nMetadata on trial 0:")
				printCT("input level-0", ctLow)
				printCT("top after ScaleDown_N", topScaled)
				printCT("top after Lattigo ModUp_N", topModUp)
				printCT("custom top ModRaise", customTop)
				for i := range A {
					printCT(fmt.Sprintf("A leaf[%d]", i), A[i])
					printCT(fmt.Sprintf("T leaf[%d]", i), T[i])
					printCT(fmt.Sprintf("L leaf[%d]", i), L[i])
				}
			}

			if !*scanPow2 {
				fmt.Printf("\nTrial %d leaf comparisons:\n", trial)
			}
			for i := range A {
				aVals := decode(paramsLeaf, encoderLeaf, decryptorLeaf, A[i])
				tVals := decode(paramsLeaf, encoderLeaf, decryptorLeaf, T[i])
				lVals := decode(paramsLeaf, encoderLeaf, decryptorLeaf, L[i])

				var d diag
				if !*scanPow2 {
					d = compare(fmt.Sprintf("leaf[%d] A vs T", i), aVals, tVals)
				} else {
					d.maxErr, d.rmse, _ = diffStats(aVals, tVals)
				}
				worstAvsT = updateWorst(worstAvsT, d.maxErr)
				if !*scanPow2 {
					d = compare(fmt.Sprintf("leaf[%d] A vs L", i), aVals, lVals)
				} else {
					d.maxErr, d.rmse, _ = diffStats(aVals, lVals)
				}
				worstAvsL = updateWorst(worstAvsL, d.maxErr)
				if !*scanPow2 {
					d = compare(fmt.Sprintf("leaf[%d] T vs L", i), tVals, lVals)
				} else {
					d.maxErr, d.rmse, _ = diffStats(tVals, lVals)
				}
				worstTvsL = updateWorst(worstTvsL, d.maxErr)
			}

			mA, err := mergeTree(rpEval, A)
			if err != nil {
				panic(err)
			}
			mT, err := mergeTree(rpEval, T)
			if err != nil {
				panic(err)
			}
			mL, err := mergeTree(rpEval, L)
			if err != nil {
				panic(err)
			}

			aTop := decode(paramsTop, encoderTop, decryptorTop, mA)
			tTop := decode(paramsTop, encoderTop, decryptorTop, mT)
			lTop := decode(paramsTop, encoderTop, decryptorTop, mL)
			bootTop := decode(paramsTop, encoderTop, decryptorTop, topModUp)

			if !*scanPow2 {
				fmt.Printf("\nTrial %d merged/top comparisons:\n", trial)
				compare("Merge(A) vs Lattigo top ModUp", bootTop, aTop)
				d := compare("Merge(A) vs Merge(T)", aTop, tTop)
				worstMAvsT = updateWorst(worstMAvsT, d.maxErr)
				d = compare("Merge(A) vs Merge(L)", aTop, lTop)
				worstMAvsL = updateWorst(worstMAvsL, d.maxErr)
				d = compare("Merge(T) vs Merge(L)", tTop, lTop)
				worstMTvsL = updateWorst(worstMTvsL, d.maxErr)
			} else {
				m, _, _ := diffStats(aTop, tTop)
				worstMAvsT = updateWorst(worstMAvsT, m)
				m, _, _ = diffStats(aTop, lTop)
				worstMAvsL = updateWorst(worstMAvsL, m)
				m, _, _ = diffStats(tTop, lTop)
				worstMTvsL = updateWorst(worstMTvsL, m)
			}
		}

		s := summary{pow: pow, worstAvsT: worstAvsT, worstAvsL: worstAvsL, worstTvsL: worstTvsL, worstMAvsT: worstMAvsT, worstMAvsL: worstMAvsL, worstMTvsL: worstMTvsL}
		summaries = append(summaries, s)

		fmt.Println()
		fmt.Printf("Summary for mulPow2=%d:\n", pow)
		fmt.Printf("  worst leaf A vs T raw = %.6e\n", worstAvsT)
		fmt.Printf("  worst leaf A vs L raw = %.6e\n", worstAvsL)
		fmt.Printf("  worst leaf T vs L raw = %.6e\n", worstTvsL)
		fmt.Printf("  worst merged A vs T raw = %.6e\n", worstMAvsT)
		fmt.Printf("  worst merged A vs L raw = %.6e\n", worstMAvsL)
		fmt.Printf("  worst merged T vs L raw = %.6e\n", worstMTvsL)
	}

	if len(summaries) > 1 {
		fmt.Println()
		fmt.Println("=== mulPow2 scan summary ===")
		fmt.Println("pow | leaf A-T | leaf A-L | leaf T-L | merged A-T | merged A-L | merged T-L")
		for _, s := range summaries {
			fmt.Printf("%3d | %.3e | %.3e | %.3e | %.3e | %.3e | %.3e\n", s.pow, s.worstAvsT, s.worstAvsL, s.worstTvsL, s.worstMAvsT, s.worstMAvsL, s.worstMTvsL)
		}
	}

	fmt.Println()
	fmt.Println("Interpretation:")
	fmt.Println("  A vs T small: custom centered ModRaise matches Lattigo top ModUp convention.")
	fmt.Println("  T vs L small: custom ModRaise is schedule-compatible with RingSwitchDown/Split.")
	fmt.Println("  A vs L small: desired low-modulus RingSwitchDown followed by custom leaf ModRaise matches the reference path.")
	fmt.Println("  If T vs L is small but A vs T/A vs L are large, the algebraic custom lift is compatible but not Lattigo-compatible; next step is to feed the custom path into C2S/EvalMod or adjust the scaling/multiplier convention.")
}
