package main

import (
	"flag"
	"fmt"
	"math"
	"math/cmplx"
	"runtime"
	"sync"
	"time"

	"github.com/tuneinsight/lattigo/v6/circuits/ckks/bootstrapping"
	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/ring"
	"github.com/tuneinsight/lattigo/v6/schemes/ckks"
	"github.com/tuneinsight/lattigo/v6/utils"
)

// Simple no-B_k CKKS bootstrapping decomposition prototype.
//
// This program compares:
//
//   1. Direct full bootstrapping core:
//        ScaleDown -> ModUp -> top C2S -> top EvalMod -> top S2C
//
//   2. Factored no-B_k full bootstrapping core:
//        ScaleDown -> ModUp -> Split
//        -> parallel leaf C2S
//        -> parallel leaf EvalMod
//        -> parallel leaf S2C
//        -> Merge
//
// There is deliberately no B_k, no B_k^{-1}, no twiddle logic, and no
// B-related level reservation.
//
// The program prints:
//   - parameters and ring-basis compatibility;
//   - Level/Scale/Degree at every major ciphertext stage;
//   - split/merge roundtrip error;
//   - direct-vs-factored C2S alignment under residue and blocked layouts;
//   - direct-vs-factored EvalMod alignment under residue and blocked layouts;
//   - direct final error vs original plaintext;
//   - factored raw final error vs original plaintext;
//   - factored/k final error vs original plaintext;
//   - direct-vs-factored raw and div-k differences;
//   - timing and speedup.
//
// Suggested runs:
//
//   k=2:
//     go run . -logN=13 -layers=1 -leafOrder=bitrev -reps=1 -warmup=0
//
//   k=4:
//     go run . -logN=16 -layers=2 -leafOrder=natural -parallelLeaves=true -leafWorkers=4 -reps=3 -warmup=1
//

type benchResult struct {
	Name  string
	Avg   time.Duration
	Total time.Duration
	Reps  int
}

type complexLeafPair struct {
	Real []*rlwe.Ciphertext
	Imag []*rlwe.Ciphertext
}

type directTrace struct {
	Scaled      *rlwe.Ciphertext
	Boot        *rlwe.Ciphertext
	C2SReal     *rlwe.Ciphertext
	C2SImag     *rlwe.Ciphertext
	EvalReal    *rlwe.Ciphertext
	EvalImag    *rlwe.Ciphertext
	Output      *rlwe.Ciphertext
	C2SRealVals []complex128
	C2SImagVals []complex128
	EvalRealVal []complex128
	EvalImagVal []complex128
	OutputVals  []complex128
}

type factoredTrace struct {
	Scaled         *rlwe.Ciphertext
	Boot           *rlwe.Ciphertext
	SplitLeaves    []*rlwe.Ciphertext
	C2S            complexLeafPair
	Eval           complexLeafPair
	S2CLeaves      []*rlwe.Ciphertext
	Output         *rlwe.Ciphertext
	C2SRealVals    [][]complex128
	C2SImagVals    [][]complex128
	EvalRealVals   [][]complex128
	EvalImagVals   [][]complex128
	S2CLeafOutVals [][]complex128
	OutputVals     []complex128
}

func bench(name string, warmup, reps int, f func() error) benchResult {
	if reps <= 0 {
		panic("reps must be positive")
	}

	for i := 0; i < warmup; i++ {
		if err := f(); err != nil {
			panic(fmt.Errorf("%s warmup failed: %w", name, err))
		}
	}

	runtime.GC()
	start := time.Now()

	for i := 0; i < reps; i++ {
		if err := f(); err != nil {
			panic(fmt.Errorf("%s benchmark failed: %w", name, err))
		}
	}

	total := time.Since(start)
	avg := total / time.Duration(reps)
	fmt.Printf("%-72s total=%10.4fs  avg=%10.6fs  reps=%d\n", name, total.Seconds(), avg.Seconds(), reps)

	return benchResult{Name: name, Avg: avg, Total: total, Reps: reps}
}

func makeResidualParams(logN int) ckks.Parameters {
	params, err := ckks.NewParametersFromLiteral(ckks.ParametersLiteral{
		LogN: logN,
		LogQ: []int{
			55,
			40, 40, 40, 40, 40,
			40, 40, 40, 40, 40,
		},
		LogP:            []int{61, 61, 61},
		LogDefaultScale: 40,
		Xs:              ring.Ternary{H: 192},
	})
	if err != nil {
		panic(err)
	}
	return params
}

func makeBootstrappingParams(logN int) (ckks.Parameters, bootstrapping.Parameters) {
	params := makeResidualParams(logN)

	btpParametersLit := bootstrapping.ParametersLiteral{
		LogN: utils.Pointy(logN),
		LogP: []int{61, 61, 61, 61},
		Xs:   params.Xs(),
	}

	btpParams, err := bootstrapping.NewParametersFromLiteral(params, btpParametersLit)
	if err != nil {
		panic(err)
	}

	// This matches the previous prototype's convention for small rings.
	if logN < 16 {
		btpParams.Mod1ParametersLiteral.LogMessageRatio += 16 - params.LogN()
	}

	return params, btpParams
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

func sameRingBasisSilent(a, b ckks.Parameters) bool {
	aQ := a.RingQ()
	bQ := b.RingQ()
	if len(aQ.SubRings) != len(bQ.SubRings) {
		return false
	}
	for i := range aQ.SubRings {
		if aQ.SubRings[i].Modulus != bQ.SubRings[i].Modulus {
			return false
		}
	}

	aP := a.RingP()
	bP := b.RingP()
	if len(aP.SubRings) != len(bP.SubRings) {
		return false
	}
	for i := range aP.SubRings {
		if aP.SubRings[i].Modulus != bP.SubRings[i].Modulus {
			return false
		}
	}
	return true
}

func reportSameRingBasis(name string, top, leaf ckks.Parameters) bool {
	ok := true
	topQ := top.RingQ()
	leafQ := leaf.RingQ()

	fmt.Println()
	fmt.Println("=== Ring-basis compatibility check ===")
	fmt.Printf("%s\n", name)
	fmt.Printf("top LogN=%d, leaf LogN=%d\n", top.LogN(), leaf.LogN())
	fmt.Printf("Q count: top=%d, leaf=%d\n", len(topQ.SubRings), len(leafQ.SubRings))

	if len(topQ.SubRings) != len(leafQ.SubRings) {
		ok = false
		fmt.Println("Q count mismatch.")
	} else {
		for i := range topQ.SubRings {
			qa := topQ.SubRings[i].Modulus
			qb := leafQ.SubRings[i].Modulus
			same := qa == qb
			if !same {
				ok = false
			}
			if i < 8 || i >= len(topQ.SubRings)-8 || !same {
				fmt.Printf("Q[%02d]: top=%d leaf=%d same=%v\n", i, qa, qb, same)
			}
		}
	}

	topP := top.RingP()
	leafP := leaf.RingP()
	fmt.Printf("P count: top=%d, leaf=%d\n", len(topP.SubRings), len(leafP.SubRings))

	if len(topP.SubRings) != len(leafP.SubRings) {
		ok = false
		fmt.Println("P count mismatch.")
	} else {
		for i := range topP.SubRings {
			pa := topP.SubRings[i].Modulus
			pb := leafP.SubRings[i].Modulus
			same := pa == pb
			if !same {
				ok = false
			}
			fmt.Printf("P[%02d]: top=%d leaf=%d same=%v\n", i, pa, pb, same)
		}
	}

	if ok {
		fmt.Println("Ring-basis check: PASS. Top and leaf Q/P moduli match exactly.")
	} else {
		fmt.Println("Ring-basis check: FAIL. Top and leaf Q/P moduli do not match.")
	}
	return ok
}

func makeDirectEvaluator(logN int, sk *rlwe.SecretKey) (ckks.Parameters, bootstrapping.Parameters, *bootstrapping.Evaluator) {
	_, btpParams := makeBootstrappingParams(logN)
	paramsBoot := btpParams.BootstrappingParameters

	fmt.Printf("Generating direct top bootstrapping keys for LogN=%d, bootLogQP=%.1f, maxLevel=%d...\n",
		logN, paramsBoot.LogQP(), paramsBoot.MaxLevel())

	btpKeys, _, err := btpParams.GenEvaluationKeys(sk)
	if err != nil {
		panic(fmt.Errorf("GenEvaluationKeys LogN=%d failed: %w", logN, err))
	}

	btpEval, err := bootstrapping.NewEvaluator(btpParams, btpKeys)
	if err != nil {
		panic(fmt.Errorf("NewEvaluator LogN=%d failed: %w", logN, err))
	}

	return paramsBoot, btpParams, btpEval
}

func makeLeafEvaluatorPoolNoBK(topBootParams ckks.Parameters, leafLogN int, skLeaf *rlwe.SecretKey, workers int) (ckks.Parameters, bootstrapping.Parameters, []*bootstrapping.Evaluator) {
	if workers <= 0 {
		panic("workers must be positive")
	}

	exactLeafParams := makeParamsFromExactBasis(topBootParams, leafLogN)
	_, leafBtpParams := makeBootstrappingParams(leafLogN)

	// Important: no B_k and no B_k^{-1}; do not reserve any B-related levels.
	leafBtpParams.ResidualParameters = exactLeafParams
	leafBtpParams.BootstrappingParameters = exactLeafParams

	paramsBoot := leafBtpParams.BootstrappingParameters

	fmt.Printf("Generating no-B_k EXACT-BASIS leaf bootstrapping keys for LogN=%d, logQP=%.1f, maxLevel=%d, workers=%d...\n",
		leafLogN, paramsBoot.LogQP(), paramsBoot.MaxLevel(), workers)

	if !sameRingBasisSilent(topBootParams, paramsBoot) {
		reportSameRingBasis("top boot params vs exact-basis leaf params", topBootParams, paramsBoot)
		panic("exact-basis construction failed: leaf params still do not match top Q/P")
	}

	btpKeys, _, err := leafBtpParams.GenEvaluationKeys(skLeaf)
	if err != nil {
		panic(fmt.Errorf("GenEvaluationKeys exact-basis LogN=%d failed: %w", leafLogN, err))
	}

	evals := make([]*bootstrapping.Evaluator, workers)
	for i := 0; i < workers; i++ {
		ev, err := bootstrapping.NewEvaluator(leafBtpParams, btpKeys)
		if err != nil {
			panic(fmt.Errorf("NewEvaluator exact-basis LogN=%d worker=%d failed: %w", leafLogN, i, err))
		}
		evals[i] = ev
	}

	return paramsBoot, leafBtpParams, evals
}

func makeValues(slots int) []complex128 {
	values := make([]complex128, slots)
	for i := range values {
		realPart := float64((i%17)-8) / 8.0
		imagPart := float64((i%11)-5) / 16.0
		values[i] = complex(realPart, imagPart)
	}
	return values
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

func decodeValues(params ckks.Parameters, encoder *ckks.Encoder, decryptor *rlwe.Decryptor, ct *rlwe.Ciphertext) []complex128 {
	pt := decryptor.DecryptNew(ct)
	values := make([]complex128, params.MaxSlots())
	if err := encoder.Decode(pt, values); err != nil {
		panic(err)
	}
	return values
}

func badComplex(z complex128) bool {
	return math.IsNaN(real(z)) || math.IsNaN(imag(z)) || math.IsInf(real(z), 0) || math.IsInf(imag(z), 0)
}

func maxAbs(a []complex128) (float64, int) {
	var max float64
	idx := -1
	for i, x := range a {
		if badComplex(x) {
			return math.Inf(1), i
		}
		v := cmplx.Abs(x)
		if v > max {
			max = v
			idx = i
		}
	}
	return max, idx
}

func diffStats(want, got []complex128, gotDiv float64) (maxErr float64, maxIdx int, rmse float64, meanErr float64) {
	if len(want) != len(got) {
		panic(fmt.Errorf("diffStats length mismatch: want=%d got=%d", len(want), len(got)))
	}
	if gotDiv == 0 {
		panic("gotDiv must be non-zero")
	}

	maxIdx = -1
	var sumSq float64
	var sum float64
	div := complex(gotDiv, 0)

	for i := range want {
		if badComplex(want[i]) || badComplex(got[i]) {
			return math.Inf(1), i, math.Inf(1), math.Inf(1)
		}
		d := cmplx.Abs(want[i] - got[i]/div)
		sum += d
		sumSq += d * d
		if d > maxErr {
			maxErr = d
			maxIdx = i
		}
	}

	if len(want) > 0 {
		meanErr = sum / float64(len(want))
		rmse = math.Sqrt(sumSq / float64(len(want)))
	}

	return
}

func bestFitScalar(want, got []complex128) complex128 {
	if len(want) != len(got) {
		panic(fmt.Errorf("bestFitScalar length mismatch: want=%d got=%d", len(want), len(got)))
	}

	var num complex128
	var den float64
	for i := range want {
		num += want[i] * cmplx.Conj(got[i])
		den += cmplx.Abs(got[i]) * cmplx.Abs(got[i])
	}
	if den == 0 {
		return 0
	}
	return num / complex(den, 0)
}

func compareVectors(name string, want, got []complex128, gotDiv float64, printSlots int) (float64, float64) {
	maxErr, maxIdx, rmse, meanErr := diffStats(want, got, gotDiv)
	wantMax, wantMaxIdx := maxAbs(want)
	gotMax, gotMaxIdx := maxAbs(got)

	fmt.Println()
	fmt.Printf("--- %s ---\n", name)
	if gotDiv != 1.0 {
		fmt.Printf("got normalization divisor = %.0f\n", gotDiv)
	}
	fmt.Printf("len=%d\n", len(want))
	fmt.Printf("want max abs = %.17e at slot %d\n", wantMax, wantMaxIdx)
	fmt.Printf("got  max abs = %.17e at slot %d\n", gotMax, gotMaxIdx)
	fmt.Printf("max abs error = %.17e at slot %d\n", maxErr, maxIdx)
	fmt.Printf("rmse          = %.17e\n", rmse)
	fmt.Printf("mean abs err  = %.17e\n", meanErr)

	if maxIdx >= 0 && maxIdx < len(want) && !math.IsInf(maxErr, 0) {
		gotVal := got[maxIdx] / complex(gotDiv, 0)
		fmt.Printf("worst slot[%d]: want=%+.17e%+.17ei got=%+.17e%+.17ei diff=%.17e\n",
			maxIdx,
			real(want[maxIdx]), imag(want[maxIdx]),
			real(gotVal), imag(gotVal),
			cmplx.Abs(want[maxIdx]-gotVal),
		)
	}

	for i := 0; i < printSlots && i < len(want); i++ {
		gotVal := got[i] / complex(gotDiv, 0)
		fmt.Printf("slot[%d]: want=%+.12e%+.12ei got=%+.12e%+.12ei diff=%.3e\n",
			i,
			real(want[i]), imag(want[i]),
			real(gotVal), imag(gotVal),
			cmplx.Abs(want[i]-gotVal),
		)
	}

	return maxErr, rmse
}

func compareVectorsWithAlpha(name string, want, got []complex128) (rawErr float64, alphaErr float64, alpha complex128) {
	rawErr, rawIdx, rawRMSE, rawMean := diffStats(want, got, 1.0)
	alpha = bestFitScalar(want, got)

	scaled := make([]complex128, len(got))
	for i := range got {
		scaled[i] = alpha * got[i]
	}
	alphaErr, alphaIdx, alphaRMSE, alphaMean := diffStats(want, scaled, 1.0)

	wantMax, wantIdx := maxAbs(want)
	gotMax, gotIdx := maxAbs(got)

	fmt.Printf("%s\n", name)
	fmt.Printf("  want max abs = %.17e at %d\n", wantMax, wantIdx)
	fmt.Printf("  got  max abs = %.17e at %d\n", gotMax, gotIdx)
	fmt.Printf("  raw max err  = %.17e at %d, rmse=%.17e, mean=%.17e\n", rawErr, rawIdx, rawRMSE, rawMean)
	fmt.Printf("  best alpha   = %.17e%+.17ei\n", real(alpha), imag(alpha))
	fmt.Printf("  alpha err    = %.17e at %d, rmse=%.17e, mean=%.17e\n", alphaErr, alphaIdx, alphaRMSE, alphaMean)
	if rawIdx >= 0 && rawIdx < len(want) {
		fmt.Printf("  raw worst[%d]: want=%+.12e%+.12ei got=%+.12e%+.12ei\n",
			rawIdx,
			real(want[rawIdx]), imag(want[rawIdx]),
			real(got[rawIdx]), imag(got[rawIdx]),
		)
	}

	return rawErr, alphaErr, alpha
}

func printCTInfo(name string, ct *rlwe.Ciphertext) {
	if ct == nil {
		fmt.Printf("%s: nil ciphertext\n", name)
		return
	}
	fmt.Printf("%s: LogN=%d Level=%d Degree=%d Scale=%v MetaData=%+v\n",
		name, ct.LogN(), ct.Level(), ct.Degree(), ct.Scale, ct.MetaData)
}

func printLeafCTInfos(stage string, cts []*rlwe.Ciphertext) {
	fmt.Println()
	fmt.Printf("=== %s ciphertext metadata ===\n", stage)
	for i, ct := range cts {
		printCTInfo(fmt.Sprintf("%s leaf[%d]", stage, i), ct)
	}
}

func bitReverse(x, bits int) int {
	var y int
	for i := 0; i < bits; i++ {
		y <<= 1
		y |= (x >> i) & 1
	}
	return y
}

func residueToLeafIndex(residue int, layers int, leafOrder string) int {
	switch leafOrder {
	case "natural":
		return residue
	case "bitrev":
		return bitReverse(residue, layers)
	default:
		panic(fmt.Errorf("unknown leafOrder=%s", leafOrder))
	}
}

func splitTree(eval *rlwe.RingPackingEvaluator, ct *rlwe.Ciphertext, layers int) ([]*rlwe.Ciphertext, error) {
	nodes := []*rlwe.Ciphertext{ct}
	for layer := 0; layer < layers; layer++ {
		next := make([]*rlwe.Ciphertext, 0, 2*len(nodes))
		for _, node := range nodes {
			even, odd, err := eval.SplitNew(node)
			if err != nil {
				return nil, fmt.Errorf("SplitNew at layer %d, logN=%d: %w", layer, node.LogN(), err)
			}
			next = append(next, even, odd)
		}
		nodes = next
	}
	return nodes, nil
}

func mergeTree(eval *rlwe.RingPackingEvaluator, leaves []*rlwe.Ciphertext) (*rlwe.Ciphertext, error) {
	nodes := leaves
	for len(nodes) > 1 {
		if len(nodes)%2 != 0 {
			return nil, fmt.Errorf("mergeTree expects even number of nodes, got %d", len(nodes))
		}
		next := make([]*rlwe.Ciphertext, 0, len(nodes)/2)
		for i := 0; i < len(nodes); i += 2 {
			parent, err := eval.MergeNew(nodes[i], nodes[i+1])
			if err != nil {
				return nil, fmt.Errorf("MergeNew children logN=%d: %w", nodes[i].LogN(), err)
			}
			next = append(next, parent)
		}
		nodes = next
	}
	return nodes[0], nil
}

func cloneCiphertexts(cts []*rlwe.Ciphertext) []*rlwe.Ciphertext {
	out := make([]*rlwe.Ciphertext, len(cts))
	for i := range cts {
		out[i] = cts[i].CopyNew()
	}
	return out
}

func c2SOnly(eval *bootstrapping.Evaluator, ct *rlwe.Ciphertext) (*rlwe.Ciphertext, *rlwe.Ciphertext, error) {
	ctReal, ctImag, err := eval.CoeffsToSlots(ct.CopyNew())
	if err != nil {
		return nil, nil, fmt.Errorf("CoeffsToSlots failed: %w", err)
	}
	if ctReal == nil || ctImag == nil {
		return nil, nil, fmt.Errorf("CoeffsToSlots returned nil output: real nil=%v imag nil=%v", ctReal == nil, ctImag == nil)
	}
	return ctReal, ctImag, nil
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

func directFullTrace(
	eval *bootstrapping.Evaluator,
	paramsTop ckks.Parameters,
	encoderTop *ckks.Encoder,
	decryptorTop *rlwe.Decryptor,
	ct *rlwe.Ciphertext,
) (directTrace, error) {

	var tr directTrace
	var err error

	tr.Scaled, tr.Boot, err = prepareBootstrapInput(eval, ct.CopyNew())
	if err != nil {
		return tr, err
	}

	fmt.Println()
	fmt.Println("=== Direct full bootstrap stages ===")
	printCTInfo("direct after ScaleDown", tr.Scaled)
	printCTInfo("direct after ModUp", tr.Boot)

	tr.C2SReal, tr.C2SImag, err = c2SOnly(eval, tr.Boot.CopyNew())
	if err != nil {
		return tr, fmt.Errorf("direct C2S failed: %w", err)
	}
	printCTInfo("direct C2S real", tr.C2SReal)
	printCTInfo("direct C2S imag", tr.C2SImag)

	tr.EvalReal, err = eval.EvalMod(tr.C2SReal.CopyNew())
	if err != nil {
		return tr, fmt.Errorf("direct EvalMod real failed: %w", err)
	}
	tr.EvalImag, err = eval.EvalMod(tr.C2SImag.CopyNew())
	if err != nil {
		return tr, fmt.Errorf("direct EvalMod imag failed: %w", err)
	}
	printCTInfo("direct EvalMod real", tr.EvalReal)
	printCTInfo("direct EvalMod imag", tr.EvalImag)

	tr.Output, err = eval.SlotsToCoeffs(tr.EvalReal.CopyNew(), tr.EvalImag.CopyNew())
	if err != nil {
		return tr, fmt.Errorf("direct S2C failed: %w", err)
	}
	printCTInfo("direct final S2C output", tr.Output)

	tr.C2SRealVals = decodeValues(paramsTop, encoderTop, decryptorTop, tr.C2SReal)
	tr.C2SImagVals = decodeValues(paramsTop, encoderTop, decryptorTop, tr.C2SImag)
	tr.EvalRealVal = decodeValues(paramsTop, encoderTop, decryptorTop, tr.EvalReal)
	tr.EvalImagVal = decodeValues(paramsTop, encoderTop, decryptorTop, tr.EvalImag)
	tr.OutputVals = decodeValues(paramsTop, encoderTop, decryptorTop, tr.Output)

	return tr, nil
}

func leafC2SOnlyParallel(rpEval *rlwe.RingPackingEvaluator, leafEvals []*bootstrapping.Evaluator, ct *rlwe.Ciphertext, layers int) ([]*rlwe.Ciphertext, complexLeafPair, error) {
	if len(leafEvals) == 0 {
		return nil, complexLeafPair{}, fmt.Errorf("empty leaf evaluator pool")
	}

	leaves, err := splitTree(rpEval, ct.CopyNew(), layers)
	if err != nil {
		return nil, complexLeafPair{}, err
	}

	outRealLeaves := make([]*rlwe.Ciphertext, len(leaves))
	outImagLeaves := make([]*rlwe.Ciphertext, len(leaves))

	type job struct{ idx int }
	type result struct {
		idx  int
		real *rlwe.Ciphertext
		imag *rlwe.Ciphertext
		err  error
	}

	jobs := make(chan job)
	results := make(chan result, len(leaves))

	workerCount := len(leafEvals)
	if workerCount > len(leaves) {
		workerCount = len(leaves)
	}

	var wg sync.WaitGroup
	for w := 0; w < workerCount; w++ {
		eval := leafEvals[w]
		wg.Add(1)
		go func(eval *bootstrapping.Evaluator) {
			defer wg.Done()
			for jb := range jobs {
				ctReal, ctImag, err := c2SOnly(eval, leaves[jb.idx].CopyNew())
				if err != nil {
					results <- result{idx: jb.idx, err: fmt.Errorf("leaf %d C2S failed: %w", jb.idx, err)}
					continue
				}
				results <- result{idx: jb.idx, real: ctReal, imag: ctImag}
			}
		}(eval)
	}

	for i := range leaves {
		jobs <- job{idx: i}
	}
	close(jobs)
	wg.Wait()
	close(results)

	for res := range results {
		if res.err != nil {
			return nil, complexLeafPair{}, res.err
		}
		outRealLeaves[res.idx] = res.real
		outImagLeaves[res.idx] = res.imag
	}

	for i := range leaves {
		if outRealLeaves[i] == nil || outImagLeaves[i] == nil {
			return nil, complexLeafPair{}, fmt.Errorf("missing C2S output at leaf %d", i)
		}
	}

	return leaves, complexLeafPair{Real: outRealLeaves, Imag: outImagLeaves}, nil
}

func leafEvalModParallel(leafEvals []*bootstrapping.Evaluator, in complexLeafPair) (complexLeafPair, error) {
	if len(in.Real) != len(in.Imag) {
		return complexLeafPair{}, fmt.Errorf("real/imag leaf length mismatch")
	}
	if len(leafEvals) == 0 {
		return complexLeafPair{}, fmt.Errorf("empty leaf evaluator pool")
	}

	k := len(in.Real)
	outReal := make([]*rlwe.Ciphertext, k)
	outImag := make([]*rlwe.Ciphertext, k)

	type job struct{ idx int }
	type result struct {
		idx  int
		real *rlwe.Ciphertext
		imag *rlwe.Ciphertext
		err  error
	}

	jobs := make(chan job)
	results := make(chan result, k)

	workerCount := len(leafEvals)
	if workerCount > k {
		workerCount = k
	}

	var wg sync.WaitGroup
	for w := 0; w < workerCount; w++ {
		eval := leafEvals[w]
		wg.Add(1)
		go func(eval *bootstrapping.Evaluator) {
			defer wg.Done()
			for jb := range jobs {
				r, err := eval.EvalMod(in.Real[jb.idx].CopyNew())
				if err != nil {
					results <- result{idx: jb.idx, err: fmt.Errorf("leaf %d EvalMod real failed: %w", jb.idx, err)}
					continue
				}
				im, err := eval.EvalMod(in.Imag[jb.idx].CopyNew())
				if err != nil {
					results <- result{idx: jb.idx, err: fmt.Errorf("leaf %d EvalMod imag failed: %w", jb.idx, err)}
					continue
				}
				results <- result{idx: jb.idx, real: r, imag: im}
			}
		}(eval)
	}

	for i := 0; i < k; i++ {
		jobs <- job{idx: i}
	}
	close(jobs)
	wg.Wait()
	close(results)

	for res := range results {
		if res.err != nil {
			return complexLeafPair{}, res.err
		}
		outReal[res.idx] = res.real
		outImag[res.idx] = res.imag
	}

	return complexLeafPair{Real: outReal, Imag: outImag}, nil
}

func leafS2CAndMergeParallel(rpEval *rlwe.RingPackingEvaluator, leafEvals []*bootstrapping.Evaluator, in complexLeafPair) ([]*rlwe.Ciphertext, *rlwe.Ciphertext, error) {
	if len(in.Real) != len(in.Imag) {
		return nil, nil, fmt.Errorf("real/imag leaf length mismatch")
	}
	if len(leafEvals) == 0 {
		return nil, nil, fmt.Errorf("empty leaf evaluator pool")
	}

	k := len(in.Real)
	coeffLeaves := make([]*rlwe.Ciphertext, k)

	type job struct{ idx int }
	type result struct {
		idx int
		ct  *rlwe.Ciphertext
		err error
	}

	jobs := make(chan job)
	results := make(chan result, k)

	workerCount := len(leafEvals)
	if workerCount > k {
		workerCount = k
	}

	var wg sync.WaitGroup
	for w := 0; w < workerCount; w++ {
		eval := leafEvals[w]
		wg.Add(1)
		go func(eval *bootstrapping.Evaluator) {
			defer wg.Done()
			for jb := range jobs {
				ct, err := eval.SlotsToCoeffs(in.Real[jb.idx].CopyNew(), in.Imag[jb.idx].CopyNew())
				if err != nil {
					results <- result{idx: jb.idx, err: fmt.Errorf("leaf %d S2C failed: %w", jb.idx, err)}
					continue
				}
				results <- result{idx: jb.idx, ct: ct}
			}
		}(eval)
	}

	for i := 0; i < k; i++ {
		jobs <- job{idx: i}
	}
	close(jobs)
	wg.Wait()
	close(results)

	for res := range results {
		if res.err != nil {
			return nil, nil, res.err
		}
		coeffLeaves[res.idx] = res.ct
	}

	out, err := mergeTree(rpEval, coeffLeaves)
	if err != nil {
		return nil, nil, err
	}

	return coeffLeaves, out, nil
}

func decodeLeafOutputs(
	paramsLeaf ckks.Parameters,
	encoderLeaf *ckks.Encoder,
	decryptorLeaf *rlwe.Decryptor,
	leaves []*rlwe.Ciphertext,
) [][]complex128 {

	out := make([][]complex128, len(leaves))
	for i := range leaves {
		out[i] = decodeValues(paramsLeaf, encoderLeaf, decryptorLeaf, leaves[i])
	}
	return out
}

func factoredFullTraceNoBK(
	topEval *bootstrapping.Evaluator,
	rpEval *rlwe.RingPackingEvaluator,
	leafEvals []*bootstrapping.Evaluator,
	paramsTop ckks.Parameters,
	encoderTop *ckks.Encoder,
	decryptorTop *rlwe.Decryptor,
	paramsLeaf ckks.Parameters,
	encoderLeaf *ckks.Encoder,
	decryptorLeaf *rlwe.Decryptor,
	ct *rlwe.Ciphertext,
	layers int,
) (factoredTrace, error) {

	var tr factoredTrace
	var err error

	tr.Scaled, tr.Boot, err = prepareBootstrapInput(topEval, ct.CopyNew())
	if err != nil {
		return tr, err
	}

	fmt.Println()
	fmt.Println("=== Factored no-B_k full bootstrap stages ===")
	printCTInfo("factored after ScaleDown", tr.Scaled)
	printCTInfo("factored after ModUp", tr.Boot)

	tr.SplitLeaves, tr.C2S, err = leafC2SOnlyParallel(rpEval, leafEvals, tr.Boot.CopyNew(), layers)
	if err != nil {
		return tr, fmt.Errorf("factored Split + leaf C2S failed: %w", err)
	}
	printLeafCTInfos("factored Split output", tr.SplitLeaves)
	printLeafCTInfos("factored leaf C2S real", tr.C2S.Real)
	printLeafCTInfos("factored leaf C2S imag", tr.C2S.Imag)

	tr.Eval, err = leafEvalModParallel(leafEvals, tr.C2S)
	if err != nil {
		return tr, fmt.Errorf("factored leaf EvalMod failed: %w", err)
	}
	printLeafCTInfos("factored leaf EvalMod real", tr.Eval.Real)
	printLeafCTInfos("factored leaf EvalMod imag", tr.Eval.Imag)

	tr.S2CLeaves, tr.Output, err = leafS2CAndMergeParallel(rpEval, leafEvals, tr.Eval)
	if err != nil {
		return tr, fmt.Errorf("factored leaf S2C + Merge failed: %w", err)
	}
	printLeafCTInfos("factored leaf S2C outputs before Merge", tr.S2CLeaves)
	printCTInfo("factored final Merge output", tr.Output)

	tr.C2SRealVals = decodeLeafOutputs(paramsLeaf, encoderLeaf, decryptorLeaf, tr.C2S.Real)
	tr.C2SImagVals = decodeLeafOutputs(paramsLeaf, encoderLeaf, decryptorLeaf, tr.C2S.Imag)
	tr.EvalRealVals = decodeLeafOutputs(paramsLeaf, encoderLeaf, decryptorLeaf, tr.Eval.Real)
	tr.EvalImagVals = decodeLeafOutputs(paramsLeaf, encoderLeaf, decryptorLeaf, tr.Eval.Imag)
	tr.S2CLeafOutVals = decodeLeafOutputs(paramsLeaf, encoderLeaf, decryptorLeaf, tr.S2CLeaves)
	tr.OutputVals = decodeValues(paramsTop, encoderTop, decryptorTop, tr.Output)

	return tr, nil
}

func extractByResidue(values []complex128, residue, k, leafSlots int) []complex128 {
	out := make([]complex128, leafSlots)
	for s := 0; s < leafSlots; s++ {
		idx := residue + s*k
		if idx >= len(values) {
			panic(fmt.Errorf("extractByResidue out of range: idx=%d len=%d residue=%d k=%d leafSlots=%d", idx, len(values), residue, k, leafSlots))
		}
		out[s] = values[idx]
	}
	return out
}

func extractByBlock(values []complex128, block, leafSlots int) []complex128 {
	out := make([]complex128, leafSlots)
	start := block * leafSlots
	for s := 0; s < leafSlots; s++ {
		idx := start + s
		if idx >= len(values) {
			panic(fmt.Errorf("extractByBlock out of range: idx=%d len=%d block=%d leafSlots=%d", idx, len(values), block, leafSlots))
		}
		out[s] = values[idx]
	}
	return out
}

func compareTopAgainstLeaves(
	stage string,
	topRealValues []complex128,
	topImagValues []complex128,
	leafRealValues [][]complex128,
	leafImagValues [][]complex128,
	layers int,
	leafOrder string,
) float64 {

	k := 1 << uint(layers)
	if len(leafRealValues) != k || len(leafImagValues) != k {
		panic(fmt.Errorf("%s expected %d leaves, got real=%d imag=%d", stage, k, len(leafRealValues), len(leafImagValues)))
	}
	leafSlots := len(leafRealValues[0])

	fmt.Println()
	fmt.Printf("=== %s: top vs Split+leaf layout comparison ===\n", stage)
	fmt.Println("residue layout: top index = residue + s*k")
	fmt.Println("blocked layout: top index = residue*leafSlots + s")
	fmt.Println("For the no-B_k experiment, the correct implementation convention is usually blocked layout.")

	bestErr := math.Inf(1)
	bestName := ""

	for residue := 0; residue < k; residue++ {
		leafIdx := residueToLeafIndex(residue, layers, leafOrder)
		leafReal := leafRealValues[leafIdx]
		leafImag := leafImagValues[leafIdx]

		topRealResidue := extractByResidue(topRealValues, residue, k, leafSlots)
		topImagResidue := extractByResidue(topImagValues, residue, k, leafSlots)
		raw, alphaErr, _ := compareVectorsWithAlpha(
			fmt.Sprintf("%s REAL residue=%d leafIdx=%d natural-residue", stage, residue, leafIdx),
			topRealResidue,
			leafReal,
		)
		if raw < bestErr {
			bestErr = raw
			bestName = fmt.Sprintf("%s REAL residue=%d natural-residue raw", stage, residue)
		}
		if alphaErr < bestErr {
			bestErr = alphaErr
			bestName = fmt.Sprintf("%s REAL residue=%d natural-residue alpha", stage, residue)
		}

		raw, alphaErr, _ = compareVectorsWithAlpha(
			fmt.Sprintf("%s IMAG residue=%d leafIdx=%d natural-residue", stage, residue, leafIdx),
			topImagResidue,
			leafImag,
		)
		if raw < bestErr {
			bestErr = raw
			bestName = fmt.Sprintf("%s IMAG residue=%d natural-residue raw", stage, residue)
		}
		if alphaErr < bestErr {
			bestErr = alphaErr
			bestName = fmt.Sprintf("%s IMAG residue=%d natural-residue alpha", stage, residue)
		}

		topRealBlock := extractByBlock(topRealValues, residue, leafSlots)
		topImagBlock := extractByBlock(topImagValues, residue, leafSlots)
		raw, alphaErr, _ = compareVectorsWithAlpha(
			fmt.Sprintf("%s REAL residue=%d leafIdx=%d blocked", stage, residue, leafIdx),
			topRealBlock,
			leafReal,
		)
		if raw < bestErr {
			bestErr = raw
			bestName = fmt.Sprintf("%s REAL residue=%d blocked raw", stage, residue)
		}
		if alphaErr < bestErr {
			bestErr = alphaErr
			bestName = fmt.Sprintf("%s REAL residue=%d blocked alpha", stage, residue)
		}

		raw, alphaErr, _ = compareVectorsWithAlpha(
			fmt.Sprintf("%s IMAG residue=%d leafIdx=%d blocked", stage, residue, leafIdx),
			topImagBlock,
			leafImag,
		)
		if raw < bestErr {
			bestErr = raw
			bestName = fmt.Sprintf("%s IMAG residue=%d blocked raw", stage, residue)
		}
		if alphaErr < bestErr {
			bestErr = alphaErr
			bestName = fmt.Sprintf("%s IMAG residue=%d blocked alpha", stage, residue)
		}
	}

	fmt.Printf("%s best candidate: %s, err=%.17e\n", stage, bestName, bestErr)
	return bestErr
}

func runSplitMergeRoundtrip(
	label string,
	rpEval *rlwe.RingPackingEvaluator,
	paramsTop ckks.Parameters,
	encoderTop *ckks.Encoder,
	decryptorTop *rlwe.Decryptor,
	ct *rlwe.Ciphertext,
	layers int,
) float64 {

	fmt.Println()
	fmt.Printf("=== Split/Merge roundtrip: %s ===\n", label)
	leaves, err := splitTree(rpEval, ct.CopyNew(), layers)
	if err != nil {
		panic(err)
	}
	printLeafCTInfos("roundtrip split", leaves)

	merged, err := mergeTree(rpEval, cloneCiphertexts(leaves))
	if err != nil {
		panic(err)
	}
	printCTInfo("roundtrip merged", merged)

	want := decodeValues(paramsTop, encoderTop, decryptorTop, ct)
	got := decodeValues(paramsTop, encoderTop, decryptorTop, merged)
	maxErr, _ := compareVectors(fmt.Sprintf("Split/Merge roundtrip %s", label), want, got, 1.0, 4)
	return maxErr
}

func main() {
	logN := flag.Int("logN", 13, "top ring degree log2(N)")
	layers := flag.Int("layers", 1, "recursive split layers; split factor k=2^layers")
	reps := flag.Int("reps", 1, "benchmark repetitions")
	warmup := flag.Int("warmup", 0, "warmup repetitions")
	parallelLeaves := flag.Bool("parallelLeaves", true, "run leaf operations in parallel")
	leafWorkersFlag := flag.Int("leafWorkers", 0, "number of leaf workers; 0 means min(numLeaves, NumCPU)")
	leafOrder := flag.String("leafOrder", "bitrev", "leaf order convention: natural or bitrev")
	printInputSlots := flag.Int("printSlots", 4, "number of slots to print in vector comparisons")
	runTiming := flag.Bool("timing", true, "run timing benchmark after correctness trace")
	flag.Parse()

	if *layers <= 0 {
		panic("layers must be positive")
	}
	leafLogN := *logN - *layers
	if leafLogN < 1 {
		panic("layers too large")
	}

	numLeaves := 1 << uint(*layers)

	leafWorkers := *leafWorkersFlag
	if !*parallelLeaves {
		leafWorkers = 1
	} else {
		if leafWorkers <= 0 {
			leafWorkers = numLeaves
			if leafWorkers > runtime.NumCPU() {
				leafWorkers = runtime.NumCPU()
			}
		}
		if leafWorkers > numLeaves {
			leafWorkers = numLeaves
		}
	}

	fmt.Println("=== SIMPLE NO-B_k CKKS BOOTSTRAPPING DECOMPOSITION ===")
	fmt.Printf("top LogN=%d, layers=%d, k=%d, leaf LogN=%d\n", *logN, *layers, numLeaves, leafLogN)
	fmt.Printf("parallelLeaves=%v, leafWorkers=%d, leafOrder=%s, reps=%d, warmup=%d\n",
		*parallelLeaves, leafWorkers, *leafOrder, *reps, *warmup)
	fmt.Println("No B_k, no B_k^{-1}, no twiddles, no B-related level reservation.")

	_, btpTop := makeBootstrappingParams(*logN)
	paramsTop := btpTop.BootstrappingParameters

	fmt.Printf("Top boot params: logN=%d, N=%d, logSlots=%d, maxLevel=%d, logQP=%.1f\n",
		paramsTop.LogN(), paramsTop.N(), paramsTop.LogMaxSlots(), paramsTop.MaxLevel(), paramsTop.LogQP())

	kgenTop := rlwe.NewKeyGenerator(paramsTop)
	skTop, pkTop := kgenTop.GenKeyPairNew()

	encoderTop := ckks.NewEncoder(paramsTop)
	encryptorTop := rlwe.NewEncryptor(paramsTop, pkTop)
	decryptorTop := rlwe.NewDecryptor(paramsTop, skTop)

	values := makeValues(paramsTop.MaxSlots())

	ctHigh := encryptAtLevel(paramsTop, encoderTop, encryptorTop, values, paramsTop.MaxLevel())
	ctLow := encryptAtLevel(paramsTop, encoderTop, encryptorTop, values, 0)

	fmt.Println()
	fmt.Println("=== Input ciphertexts ===")
	printCTInfo("top high-level input", ctHigh)
	printCTInfo("top low-level bootstrap input", ctLow)

	highVals := decodeValues(paramsTop, encoderTop, decryptorTop, ctHigh)
	lowVals := decodeValues(paramsTop, encoderTop, decryptorTop, ctLow)
	compareVectors("High-level input decrypt/decode vs plaintext", values, highVals, 1.0, *printInputSlots)
	compareVectors("Low-level input decrypt/decode vs plaintext", values, lowVals, 1.0, *printInputSlots)

	fmt.Println()
	fmt.Println("Generating ring-switching keys for Split/Merge...")
	rpk := &rlwe.RingPackingEvaluationKey{}
	skMap, err := rpk.GenRingSwitchingKeys(paramsTop, skTop, leafLogN, rlwe.EvaluationKeyParameters{})
	if err != nil {
		panic(fmt.Errorf("GenRingSwitchingKeys failed: %w", err))
	}
	rpEval := rlwe.NewRingPackingEvaluator(rpk)
	fmt.Printf("RingPackingEvaluator: minLogN=%d, maxLogN=%d\n", rpk.MinLogN(), rpk.MaxLogN())

	skLeaf, ok := skMap[leafLogN]
	if !ok || skLeaf == nil {
		panic(fmt.Errorf("missing leaf secret key for LogN=%d", leafLogN))
	}

	fmt.Println()
	fmt.Println("Generating direct top evaluator...")
	_, _, topEval := makeDirectEvaluator(*logN, skTop)

	fmt.Println()
	fmt.Println("Generating no-B_k leaf evaluator pool...")
	paramsLeaf, _, leafEvals := makeLeafEvaluatorPoolNoBK(paramsTop, leafLogN, skLeaf, leafWorkers)
	encoderLeaf := ckks.NewEncoder(paramsLeaf)
	decryptorLeaf := rlwe.NewDecryptor(paramsLeaf, skLeaf)

	fmt.Printf("Leaf boot params: logN=%d, N=%d, logSlots=%d, maxLevel=%d, logQP=%.1f\n",
		paramsLeaf.LogN(), paramsLeaf.N(), paramsLeaf.LogMaxSlots(), paramsLeaf.MaxLevel(), paramsLeaf.LogQP())

	reportSameRingBasis("top boot params vs shared-basis leaf boot params", paramsTop, paramsLeaf)

	roundtripHighErr := runSplitMergeRoundtrip("high-level input", rpEval, paramsTop, encoderTop, decryptorTop, ctHigh, *layers)
	roundtripLowErr := runSplitMergeRoundtrip("low-level input", rpEval, paramsTop, encoderTop, decryptorTop, ctLow, *layers)

	direct, err := directFullTrace(topEval, paramsTop, encoderTop, decryptorTop, ctLow.CopyNew())
	if err != nil {
		panic(fmt.Errorf("direct full trace failed: %w", err))
	}

	factored, err := factoredFullTraceNoBK(
		topEval,
		rpEval,
		leafEvals,
		paramsTop,
		encoderTop,
		decryptorTop,
		paramsLeaf,
		encoderLeaf,
		decryptorLeaf,
		ctLow.CopyNew(),
		*layers,
	)
	if err != nil {
		panic(fmt.Errorf("factored no-B_k full trace failed: %w", err))
	}

	c2sBest := compareTopAgainstLeaves(
		"C2S input to EvalMod",
		direct.C2SRealVals,
		direct.C2SImagVals,
		factored.C2SRealVals,
		factored.C2SImagVals,
		*layers,
		*leafOrder,
	)

	evalBest := compareTopAgainstLeaves(
		"After EvalMod",
		direct.EvalRealVal,
		direct.EvalImagVal,
		factored.EvalRealVals,
		factored.EvalImagVals,
		*layers,
		*leafOrder,
	)

	fmt.Println()
	fmt.Println("=== Final output quality ===")
	directMaxErr, directRMSE := compareVectors("Direct full bootstrap vs plaintext", values, direct.OutputVals, 1.0, *printInputSlots)
	factoredRawMaxErr, factoredRawRMSE := compareVectors("Factored no-B_k full bootstrap RAW vs plaintext", values, factored.OutputVals, 1.0, *printInputSlots)
	factoredDivMaxErr, factoredDivRMSE := compareVectors("Factored no-B_k full bootstrap DIV-K vs plaintext", values, factored.OutputVals, float64(numLeaves), *printInputSlots)

	fmt.Println()
	fmt.Println("=== Direct vs factored final comparison ===")
	rawDFMaxErr, rawDFRMSE := compareVectors("Direct output vs factored RAW output", direct.OutputVals, factored.OutputVals, 1.0, *printInputSlots)
	divDFMaxErr, divDFRMSE := compareVectors("Direct output vs factored DIV-K output", direct.OutputVals, factored.OutputVals, float64(numLeaves), *printInputSlots)

	var directTiming benchResult
	var factoredTiming benchResult

	if *runTiming {
		fmt.Println()
		fmt.Println("=== Timing benchmark ===")
		directTiming = bench("Direct full: ScaleDown+ModUp+C2S+EvalMod+S2C", *warmup, *reps, func() error {
			_, err := directFullTraceQuiet(topEval, ctLow.CopyNew())
			return err
		})

		factoredTiming = bench("Factored no-B_k full: ScaleDown+ModUp+Split+C2S+EvalMod+S2C+Merge", *warmup, *reps, func() error {
			_, err := factoredFullTraceNoBKQuiet(topEval, rpEval, leafEvals, ctLow.CopyNew(), *layers)
			return err
		})
	}

	fmt.Println()
	fmt.Println("=== Summary ===")
	fmt.Printf("top LogN=%d, layers=%d, k=%d, leaf LogN=%d, leafOrder=%s\n", *logN, *layers, numLeaves, leafLogN, *leafOrder)
	fmt.Printf("Ring basis compatibility             = %v\n", sameRingBasisSilent(paramsTop, paramsLeaf))
	fmt.Printf("Split/Merge high-level roundtrip max = %.6e\n", roundtripHighErr)
	fmt.Printf("Split/Merge low-level roundtrip max  = %.6e\n", roundtripLowErr)
	fmt.Printf("Best C2S alignment error             = %.6e\n", c2sBest)
	fmt.Printf("Best after-EvalMod alignment error   = %.6e\n", evalBest)
	fmt.Printf("Direct final max/RMSE vs plaintext   = %.6e / %.6e\n", directMaxErr, directRMSE)
	fmt.Printf("Factored RAW max/RMSE vs plaintext   = %.6e / %.6e\n", factoredRawMaxErr, factoredRawRMSE)
	fmt.Printf("Factored DIV-K max/RMSE vs plaintext = %.6e / %.6e\n", factoredDivMaxErr, factoredDivRMSE)
	fmt.Printf("Direct vs factored RAW max/RMSE      = %.6e / %.6e\n", rawDFMaxErr, rawDFRMSE)
	fmt.Printf("Direct vs factored DIV-K max/RMSE    = %.6e / %.6e\n", divDFMaxErr, divDFRMSE)
	fmt.Printf("Observed gamma_k                     = %d if DIV-K comparison is small\n", numLeaves)
	fmt.Printf("Direct output Level                  = %d\n", direct.Output.Level())
	fmt.Printf("Factored output Level                = %d\n", factored.Output.Level())

	if *runTiming {
		fmt.Printf("Direct full avg                      = %.6fs\n", directTiming.Avg.Seconds())
		fmt.Printf("Factored no-B_k full avg             = %.6fs\n", factoredTiming.Avg.Seconds())
		fmt.Printf("Speedup direct/factored              = %.4fx\n", directTiming.Avg.Seconds()/factoredTiming.Avg.Seconds())
	}

	fmt.Println()
	fmt.Println("Decision guide:")
	fmt.Println("1. C2S alignment should be tiny under blocked layout; this checks Split+leafC2S = blocked coefficient grouping after top C2S.")
	fmt.Println("2. After-EvalMod alignment should remain small; this checks that top and leaf EvalMod implement the same coordinate-wise scalar map.")
	fmt.Println("3. Direct-vs-factored RAW is expected to show a global factor k.")
	fmt.Println("4. Direct-vs-factored DIV-K should be small; this is the main no-B_k bootstrapping equivalence check.")
	fmt.Println("5. The final Direct and Factored-DIV-K errors vs plaintext judge the actual bootstrapping quality.")
}

func directFullTraceQuiet(eval *bootstrapping.Evaluator, ct *rlwe.Ciphertext) (*rlwe.Ciphertext, error) {
	_, ctBoot, err := prepareBootstrapInput(eval, ct.CopyNew())
	if err != nil {
		return nil, err
	}

	ctReal, ctImag, err := c2SOnly(eval, ctBoot.CopyNew())
	if err != nil {
		return nil, err
	}

	ctReal, err = eval.EvalMod(ctReal)
	if err != nil {
		return nil, err
	}

	ctImag, err = eval.EvalMod(ctImag)
	if err != nil {
		return nil, err
	}

	out, err := eval.SlotsToCoeffs(ctReal, ctImag)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func factoredFullTraceNoBKQuiet(
	topEval *bootstrapping.Evaluator,
	rpEval *rlwe.RingPackingEvaluator,
	leafEvals []*bootstrapping.Evaluator,
	ct *rlwe.Ciphertext,
	layers int,
) (*rlwe.Ciphertext, error) {

	_, ctBoot, err := prepareBootstrapInput(topEval, ct.CopyNew())
	if err != nil {
		return nil, err
	}

	_, pair, err := leafC2SOnlyParallel(rpEval, leafEvals, ctBoot.CopyNew(), layers)
	if err != nil {
		return nil, err
	}

	pair, err = leafEvalModParallel(leafEvals, pair)
	if err != nil {
		return nil, err
	}

	_, out, err := leafS2CAndMergeParallel(rpEval, leafEvals, pair)
	if err != nil {
		return nil, err
	}

	return out, nil
}
