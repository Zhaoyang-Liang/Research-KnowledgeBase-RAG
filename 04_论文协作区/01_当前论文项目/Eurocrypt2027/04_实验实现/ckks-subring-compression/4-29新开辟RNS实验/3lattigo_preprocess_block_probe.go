package main

import (
	"flag"
	"fmt"
	"math"
	"math/cmplx"

	"github.com/tuneinsight/lattigo/v6/circuits/ckks/bootstrapping"
	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/ring"
	"github.com/tuneinsight/lattigo/v6/schemes/ckks"
	"github.com/tuneinsight/lattigo/v6/utils"
)

// Lattigo preprocessing-block scheduling probe.
//
// This is the follow-up to lattigo_rns_schedule_probe.go.
// The previous probe tested the too-aggressive schedule:
//
//   A: ScaleDown_N -> ModUp_N -> Split
//   B: ScaleDown_N -> Split -> ModUp_n
//
// and it failed. This file tests whether the whole preprocessing pair can be
// moved below RingSwitchDown/Split:
//
//   A: ScaleDown_N -> ModUp_N -> Split
//   C: Split -> ScaleDown_n -> ModUp_n
//
// It also keeps the previous B path as a negative/control path:
//
//   B: ScaleDown_N -> Split -> ModUp_n
//
// Suggested runs:
//
//   go run lattigo_preprocess_block_probe.go -logN=14 -layers=2 -trials=3 -leafMode=leaf
//   go run lattigo_preprocess_block_probe.go -logN=14 -layers=2 -trials=3 -leafMode=top
//   go run lattigo_preprocess_block_probe.go -logN=14 -layers=2 -trials=3 -leafMode=none
//
// Interpretation:
//   - A vs C small: the entire ScaleDown+ModUp preprocessing block can be moved
//     after Split/RingSwitchDown, at least with the same Q/P basis.
//   - A vs B large but A vs C small: ScaleDown and ModUp must move together.
//   - A vs C large for leafMode=leaf but small for leafMode=top: convention issue,
//     likely LogMessageRatio / Mod1 scaling mismatch between top and leaf.
//   - A vs C large in every mode: Lattigo preprocessing is not schedule-compatible
//     in this form; reduced-Q needs a lower-level custom preprocessing/RNS lift.

type diag struct {
	maxErr  float64
	rmse    float64
	mean    float64
	gamma   complex128 // got ~= gamma * want
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

func makeExactBasisLeafEvaluator(topBootParams ckks.Parameters, topBtp bootstrapping.Parameters, leafLogN int, skLeaf *rlwe.SecretKey, mode string) (ckks.Parameters, bootstrapping.Parameters, *bootstrapping.Evaluator) {
	exactLeaf := makeParamsFromExactBasis(topBootParams, leafLogN)

	smallAdjust := true
	if mode == "none" {
		smallAdjust = false
	}
	_, leafBtp := makeBootstrappingParams(leafLogN, smallAdjust)

	switch mode {
	case "leaf":
		// Keep the leaf's own small-ring convention.
	case "top":
		// Force leaf Mod1/MessageRatio convention to match the top evaluator.
		leafBtp.Mod1ParametersLiteral.LogMessageRatio = topBtp.Mod1ParametersLiteral.LogMessageRatio
	case "none":
		// No small-ring LogMessageRatio bump.
	default:
		panic("leafMode must be one of: leaf, top, none")
	}

	leafBtp.ResidualParameters = exactLeaf
	leafBtp.BootstrappingParameters = exactLeaf
	if !sameRingBasis(topBootParams, leafBtp.BootstrappingParameters) {
		panic("exact-basis leaf params do not share top Q/P basis")
	}
	keys, _, err := leafBtp.GenEvaluationKeys(skLeaf)
	if err != nil {
		panic(fmt.Errorf("leaf GenEvaluationKeys failed: %w", err))
	}
	ev, err := bootstrapping.NewEvaluator(leafBtp, keys)
	if err != nil {
		panic(fmt.Errorf("leaf NewEvaluator failed: %w", err))
	}
	return leafBtp.BootstrappingParameters, leafBtp, ev
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
	fmt.Printf("%-38s LogN=%d Level=%d Degree=%d Scale=%v\n", name, ct.LogN(), ct.Level(), ct.Degree(), ct.Scale)
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

func bad(z complex128) bool {
	return math.IsNaN(real(z)) || math.IsNaN(imag(z)) || math.IsInf(real(z), 0) || math.IsInf(imag(z), 0)
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

func compare(want, got []complex128) diag {
	if len(want) != len(got) {
		panic(fmt.Errorf("length mismatch want=%d got=%d", len(want), len(got)))
	}
	var maxErr, sum, sumSq float64
	for i := range want {
		d := cmplx.Abs(want[i] - got[i])
		if d > maxErr {
			maxErr = d
		}
		sum += d
		sumSq += d * d
	}
	gamma := estimateGamma(want, got)
	var fitMax, fitSumSq float64
	if gamma == 0 || bad(gamma) {
		fitMax = math.Inf(1)
		fitSumSq = math.Inf(1)
	} else {
		for i := range want {
			d := cmplx.Abs(want[i] - got[i]/gamma)
			if d > fitMax {
				fitMax = d
			}
			fitSumSq += d * d
		}
	}
	n := float64(len(want))
	return diag{maxErr: maxErr, rmse: math.Sqrt(sumSq / n), mean: sum / n, gamma: gamma, fitMax: fitMax, fitRMSE: math.Sqrt(fitSumSq / n)}
}

func printDiag(name string, d diag) {
	fmt.Printf("%-54s rawMax=%10.3e rmse=%10.3e gamma=%+.6e%+.6ei fitMax=%10.3e fitRMSE=%10.3e\n",
		name, d.maxErr, d.rmse, real(d.gamma), imag(d.gamma), d.fitMax, d.fitRMSE)
}

func modUpLeaf(eval *bootstrapping.Evaluator, leaves []*rlwe.Ciphertext) ([]*rlwe.Ciphertext, error) {
	out := make([]*rlwe.Ciphertext, len(leaves))
	for i := range leaves {
		ct, err := eval.ModUp(leaves[i].CopyNew())
		if err != nil {
			return nil, fmt.Errorf("leaf[%d] ModUp failed: %w", i, err)
		}
		out[i] = ct
	}
	return out, nil
}

func scaleDownModUpLeaf(eval *bootstrapping.Evaluator, leaves []*rlwe.Ciphertext) ([]*rlwe.Ciphertext, []*rlwe.Ciphertext, error) {
	scaled := make([]*rlwe.Ciphertext, len(leaves))
	out := make([]*rlwe.Ciphertext, len(leaves))
	for i := range leaves {
		s, _, err := eval.ScaleDown(leaves[i].CopyNew())
		if err != nil {
			return nil, nil, fmt.Errorf("leaf[%d] ScaleDown failed: %w", i, err)
		}
		b, err := eval.ModUp(s.CopyNew())
		if err != nil {
			return nil, nil, fmt.Errorf("leaf[%d] ModUp after ScaleDown failed: %w", i, err)
		}
		scaled[i] = s
		out[i] = b
	}
	return scaled, out, nil
}

func updateWorst(d diag, raw *float64, fit *float64) {
	if d.maxErr > *raw {
		*raw = d.maxErr
	}
	if d.fitMax > *fit {
		*fit = d.fitMax
	}
}

func main() {
	logN := flag.Int("logN", 13, "top log2 ring degree")
	layers := flag.Int("layers", 1, "recursive split layers; k=2^layers")
	trials := flag.Int("trials", 1, "number of randomized plaintext trials")
	topSmallAdjust := flag.Bool("topSmallAdjust", true, "apply the same small-logN LogMessageRatio bump to top params")
	leafMode := flag.String("leafMode", "leaf", "leaf Mod1 convention: leaf, top, or none")
	flag.Parse()

	if *layers <= 0 || *layers >= *logN {
		panic("invalid layers")
	}
	leafLogN := *logN - *layers
	k := 1 << uint(*layers)

	fmt.Println("=== Lattigo preprocessing-block scheduling probe ===")
	fmt.Printf("top LogN=%d, layers=%d, k=%d, leaf LogN=%d, trials=%d\n", *logN, *layers, k, leafLogN, *trials)
	fmt.Printf("topSmallAdjust=%v, leafMode=%s\n", *topSmallAdjust, *leafMode)
	fmt.Println("A: ScaleDown_N -> ModUp_N -> Split")
	fmt.Println("B: ScaleDown_N -> Split -> ModUp_n       (previous control path)")
	fmt.Println("C: Split -> ScaleDown_n -> ModUp_n       (new preprocessing-block path)")

	_, topBtp0 := makeBootstrappingParams(*logN, *topSmallAdjust)
	paramsTop := topBtp0.BootstrappingParameters
	kgenTop := rlwe.NewKeyGenerator(paramsTop)
	skTop, pkTop := kgenTop.GenKeyPairNew()

	fmt.Println("Generating ring switching keys...")
	rpk := &rlwe.RingPackingEvaluationKey{}
	skMap, err := rpk.GenRingSwitchingKeys(paramsTop, skTop, leafLogN, rlwe.EvaluationKeyParameters{})
	if err != nil {
		panic(fmt.Errorf("GenRingSwitchingKeys failed: %w", err))
	}
	skLeaf := skMap[leafLogN]
	if skLeaf == nil {
		panic("missing leaf secret")
	}
	rpEval := rlwe.NewRingPackingEvaluator(rpk)

	fmt.Println("Generating bootstrapping evaluators...")
	_, topBtp, topEval := makeTopEvaluator(*logN, skTop, *topSmallAdjust)
	paramsLeaf, leafBtp, leafEval := makeExactBasisLeafEvaluator(paramsTop, topBtp, leafLogN, skLeaf, *leafMode)

	fmt.Printf("top  params: LogN=%d MaxLevel=%d LogQP=%.1f Slots=%d LogMessageRatio=%d\n", paramsTop.LogN(), paramsTop.MaxLevel(), paramsTop.LogQP(), paramsTop.MaxSlots(), topBtp.Mod1ParametersLiteral.LogMessageRatio)
	fmt.Printf("leaf params: LogN=%d MaxLevel=%d LogQP=%.1f Slots=%d LogMessageRatio=%d exactBasis=%v\n", paramsLeaf.LogN(), paramsLeaf.MaxLevel(), paramsLeaf.LogQP(), paramsLeaf.MaxSlots(), leafBtp.Mod1ParametersLiteral.LogMessageRatio, sameRingBasis(paramsTop, paramsLeaf))

	encTop := ckks.NewEncoder(paramsTop)
	encLeaf := ckks.NewEncoder(paramsLeaf)
	encryptorTop := rlwe.NewEncryptor(paramsTop, pkTop)
	decTop := rlwe.NewDecryptor(paramsTop, skTop)
	decLeaf := rlwe.NewDecryptor(paramsLeaf, skLeaf)

	var worstBLeafRaw, worstBLeafFit, worstCLeafRaw, worstCLeafFit float64
	var worstBMergedRaw, worstBMergedFit, worstCMergedRaw, worstCMergedFit float64

	for t := 0; t < *trials; t++ {
		values := makeValues(paramsTop.MaxSlots(), t)
		ctLow := encryptAtLevel(paramsTop, encTop, encryptorTop, values, 0)

		topScaled, _, err := topEval.ScaleDown(ctLow.CopyNew())
		if err != nil {
			panic(fmt.Errorf("top ScaleDown failed: %w", err))
		}
		topBoot, err := topEval.ModUp(topScaled.CopyNew())
		if err != nil {
			panic(fmt.Errorf("top ModUp failed: %w", err))
		}

		// A: top ScaleDown + top ModUp, then Split.
		leavesA, err := splitTree(rpEval, topBoot.CopyNew(), *layers)
		if err != nil {
			panic(fmt.Errorf("path A split failed: %w", err))
		}

		// B: top ScaleDown, Split, then leaf ModUp. Kept as a control path.
		leavesTopScaled, err := splitTree(rpEval, topScaled.CopyNew(), *layers)
		if err != nil {
			panic(fmt.Errorf("path B split topScaled failed: %w", err))
		}
		leavesB, err := modUpLeaf(leafEval, leavesTopScaled)
		if err != nil {
			panic(fmt.Errorf("path B leaf ModUp failed: %w", err))
		}

		// C: Split raw low-level ciphertext, then leaf ScaleDown + leaf ModUp.
		leavesRaw, err := splitTree(rpEval, ctLow.CopyNew(), *layers)
		if err != nil {
			panic(fmt.Errorf("path C split raw failed: %w", err))
		}
		leavesCScaled, leavesC, err := scaleDownModUpLeaf(leafEval, leavesRaw)
		if err != nil {
			panic(fmt.Errorf("path C leaf ScaleDown+ModUp failed: %w", err))
		}

		if t == 0 {
			fmt.Println("\nMetadata on trial 0:")
			printCT("input level-0", ctLow)
			printCT("top after ScaleDown_N", topScaled)
			printCT("top after ModUp_N", topBoot)
			for i := range leavesA {
				printCT(fmt.Sprintf("A leaf[%d]=Split(top preproc)", i), leavesA[i])
				printCT(fmt.Sprintf("B leaf[%d]=ModUp(Split topScaled)", i), leavesB[i])
				printCT(fmt.Sprintf("C leaf[%d]=leaf ScaleDown", i), leavesCScaled[i])
				printCT(fmt.Sprintf("C leaf[%d]=leaf preproc", i), leavesC[i])
			}
		}

		fmt.Printf("\nTrial %d leaf comparisons, WANT=A=Split(top ScaleDown+ModUp):\n", t)
		for i := range leavesA {
			valsA := decode(paramsLeaf, encLeaf, decLeaf, leavesA[i])
			valsB := decode(paramsLeaf, encLeaf, decLeaf, leavesB[i])
			valsC := decode(paramsLeaf, encLeaf, decLeaf, leavesC[i])
			dB := compare(valsA, valsB)
			dC := compare(valsA, valsC)
			printDiag(fmt.Sprintf("leaf[%d] A vs B", i), dB)
			printDiag(fmt.Sprintf("leaf[%d] A vs C", i), dC)
			updateWorst(dB, &worstBLeafRaw, &worstBLeafFit)
			updateWorst(dC, &worstCLeafRaw, &worstCLeafFit)
		}

		mergedA, err := mergeTree(rpEval, leavesA)
		if err != nil {
			panic(fmt.Errorf("merge A failed: %w", err))
		}
		mergedB, err := mergeTree(rpEval, leavesB)
		if err != nil {
			panic(fmt.Errorf("merge B failed: %w", err))
		}
		mergedC, err := mergeTree(rpEval, leavesC)
		if err != nil {
			panic(fmt.Errorf("merge C failed: %w", err))
		}

		valsTop := decode(paramsTop, encTop, decTop, topBoot)
		valsA := decode(paramsTop, encTop, decTop, mergedA)
		valsB := decode(paramsTop, encTop, decTop, mergedB)
		valsC := decode(paramsTop, encTop, decTop, mergedC)
		dA := compare(valsTop, valsA)
		dB := compare(valsTop, valsB)
		dC := compare(valsTop, valsC)
		dAB := compare(valsA, valsB)
		dAC := compare(valsA, valsC)
		fmt.Printf("\nTrial %d merged/top comparisons:\n", t)
		printDiag("Merge(A) vs top ModUp", dA)
		printDiag("Merge(B) vs top ModUp", dB)
		printDiag("Merge(C) vs top ModUp", dC)
		printDiag("Merged A vs merged B", dAB)
		printDiag("Merged A vs merged C", dAC)
		updateWorst(dB, &worstBMergedRaw, &worstBMergedFit)
		updateWorst(dC, &worstCMergedRaw, &worstCMergedFit)
	}

	fmt.Println("\n=== Summary ===")
	fmt.Printf("B control path worst leaf raw / fit   = %.6e / %.6e\n", worstBLeafRaw, worstBLeafFit)
	fmt.Printf("C block path   worst leaf raw / fit   = %.6e / %.6e\n", worstCLeafRaw, worstCLeafFit)
	fmt.Printf("B control path worst merged raw / fit = %.6e / %.6e\n", worstBMergedRaw, worstBMergedFit)
	fmt.Printf("C block path   worst merged raw / fit = %.6e / %.6e\n", worstCMergedRaw, worstCMergedFit)

	fmt.Println("\nInterpretation:")
	fmt.Println("  This tests whether ScaleDown+ModUp can move as a pair below RingSwitchDown/Split.")
	fmt.Println("  If C is small while B is large, the preprocessing pair is schedule-compatible but ModUp alone is not.")
	fmt.Println("  If C is large in leafMode=leaf but small in leafMode=top, the issue is likely implementation scaling convention.")
	fmt.Println("  If C is large in all modes, reduced-Q experiments likely require a lower-level custom ModUp/RNS lift rather than Eval.ModUp as-is.")
}
