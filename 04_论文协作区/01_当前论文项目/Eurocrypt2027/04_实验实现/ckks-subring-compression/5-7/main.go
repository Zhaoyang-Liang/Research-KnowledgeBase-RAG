package main

import (
	"flag"
	"fmt"
	"math"
	"math/cmplx"
	"os"
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

type leafBootstrapper interface {
	ScaleDown(ctIn *rlwe.Ciphertext) (*rlwe.Ciphertext, *rlwe.Scale, error)
	ModUp(ctIn *rlwe.Ciphertext) (*rlwe.Ciphertext, error)
	CoeffsToSlots(ctIn *rlwe.Ciphertext) (*rlwe.Ciphertext, *rlwe.Ciphertext, error)
	EvalMod(ctIn *rlwe.Ciphertext) (*rlwe.Ciphertext, error)
	SlotsToCoeffs(ctReal, ctImag *rlwe.Ciphertext) (*rlwe.Ciphertext, error)
	MaxBootLevel() int
}

type lattigoLeafBootstrapper struct {
	*bootstrapping.Evaluator
}

func (l lattigoLeafBootstrapper) MaxBootLevel() int {
	return l.BootstrappingParameters.MaxLevel()
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

// ExpConfig holds all tunable experiment parameters.
// Pass it through the helper functions instead of individual arguments.
type ExpConfig struct {
	// Secret key Hamming weight. 0 = uniform ternary (P=2/3), i.e. AnyuWang-style dense secret.
	H int
	// Bit-size of the first residual Q prime.
	Q0Bits int
	// Number of circuit primes in LogQ.
	CircuitLevels int
	// Bit-size of each circuit prime in LogQ.
	CircuitPrimeBits int
	// Number of 61-bit P primes for key switching.
	NumP int
	// Default CKKS scale bit-size.
	DefaultScaleBits int
}

// Xs returns the ring.Ternary matching this config.
func (c ExpConfig) Xs() ring.Ternary {
	if c.H <= 0 {
		return ring.Ternary{P: 2.0 / 3}
	}
	return ring.Ternary{H: c.H}
}

// logQI builds the residual LogQ slice: [q0] + CircuitLevels circuit primes.
func (c ExpConfig) logQI() []int {
	q := make([]int, 1+c.CircuitLevels)
	q[0] = c.Q0Bits
	for i := 1; i <= c.CircuitLevels; i++ {
		q[i] = c.CircuitPrimeBits
	}
	return q
}

// logPI builds the LogP slice: NumP×[61].
func (c ExpConfig) logPI() []int {
	p := make([]int, c.NumP)
	for i := range p {
		p[i] = 61
	}
	return p
}

func makeResidualParams(logN int, cfg ExpConfig) ckks.Parameters {
	params, err := ckks.NewParametersFromLiteral(ckks.ParametersLiteral{
		LogN:            logN,
		LogQ:            cfg.logQI(),
		LogP:            cfg.logPI(),
		LogDefaultScale: cfg.DefaultScaleBits,
		Xs:              cfg.Xs(),
	})
	if err != nil {
		panic(err)
	}
	return params
}

func makeBootstrappingParams(logN int, cfg ExpConfig) (ckks.Parameters, bootstrapping.Parameters) {
	params := makeResidualParams(logN, cfg)

	btpParametersLit := bootstrapping.ParametersLiteral{
		LogN: utils.Pointy(logN),
		LogP: cfg.logPI(),
		Xs:   params.Xs(),
	}

	btpParams, err := bootstrapping.NewParametersFromLiteral(params, btpParametersLit)
	if err != nil {
		panic(err)
	}

	if logN < 16 {
		btpParams.Mod1ParametersLiteral.LogMessageRatio += 16 - params.LogN()
	}

	return params, btpParams
}

func copyQPrefix(params ckks.Parameters, count int) []uint64 {
	rq := params.RingQ()
	if count <= 0 || count > len(rq.SubRings) {
		panic(fmt.Errorf("invalid Q prefix count=%d for source count=%d", count, len(rq.SubRings)))
	}
	q := make([]uint64, count)
	for i := range q {
		q[i] = rq.SubRings[i].Modulus
	}
	return q
}

func copyPPrefix(params ckks.Parameters, count int) []uint64 {
	rp := params.RingP()
	if count < 0 || count > len(rp.SubRings) {
		panic(fmt.Errorf("invalid P prefix count=%d for source count=%d", count, len(rp.SubRings)))
	}
	p := make([]uint64, count)
	for i := range p {
		p[i] = rp.SubRings[i].Modulus
	}
	return p
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
	if rp == nil {
		return nil
	}
	p := make([]uint64, len(rp.SubRings))
	for i := range rp.SubRings {
		p[i] = rp.SubRings[i].Modulus
	}
	return p
}

func makeParamsFromPrefixBasis(template ckks.Parameters, logN, qCount, pCount int) ckks.Parameters {
	params, err := ckks.NewParametersFromLiteral(ckks.ParametersLiteral{
		LogN:            logN,
		Q:               copyQPrefix(template, qCount),
		P:               copyPPrefix(template, pCount),
		LogDefaultScale: template.LogDefaultScale(),
		Xs:              template.Xs(),
	})
	if err != nil {
		panic(fmt.Errorf("makeParamsFromPrefixBasis failed: %w", err))
	}
	return params
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

func reportRingBasisPrefix(name string, top, leaf ckks.Parameters) bool {
	ok := true
	topQ := top.RingQ()
	leafQ := leaf.RingQ()

	fmt.Println("=== Ring-basis prefix check:", name, "===")
	fmt.Printf("top LogN=%d, leaf LogN=%d\n", top.LogN(), leaf.LogN())
	fmt.Printf("Q count: top=%d, leaf=%d\n", len(topQ.SubRings), len(leafQ.SubRings))

	if len(leafQ.SubRings) > len(topQ.SubRings) {
		ok = false
		fmt.Println("Q prefix check: FAIL. Leaf has more Q primes than top.")
	} else {
		for i := range leafQ.SubRings {
			qa := topQ.SubRings[i].Modulus
			qb := leafQ.SubRings[i].Modulus
			same := qa == qb
			ok = ok && same
			if i < 8 || i >= len(leafQ.SubRings)-8 || !same {
				fmt.Printf("Q[%02d]: top=%d leaf=%d same=%v\n", i, qa, qb, same)
			}
		}
	}

	topP := top.RingP()
	leafP := leaf.RingP()
	topPCount := 0
	leafPCount := 0
	if topP != nil {
		topPCount = len(topP.SubRings)
	}
	if leafP != nil {
		leafPCount = len(leafP.SubRings)
	}
	fmt.Printf("P count: top=%d, leaf=%d\n", topPCount, leafPCount)

	if leafPCount > topPCount {
		ok = false
		fmt.Println("P prefix check: FAIL. Leaf has more P primes than top.")
	} else if leafP != nil {
		for i := range leafP.SubRings {
			pa := topP.SubRings[i].Modulus
			pb := leafP.SubRings[i].Modulus
			same := pa == pb
			ok = ok && same
			fmt.Printf("P[%02d]: top=%d leaf=%d same=%v\n", i, pa, pb, same)
		}
	}

	if ok {
		fmt.Println("Ring-basis prefix check: PASS. Leaf Q/P are prefixes of top Q/P.")
	} else {
		fmt.Println("Ring-basis prefix check: FAIL. Leaf Q/P are not prefixes of top Q/P.")
	}
	return ok
}

func reportRingBasisExact(name string, top, leaf ckks.Parameters) bool {
	ok := true
	topQ := top.RingQ()
	leafQ := leaf.RingQ()

	fmt.Println("=== Ring-basis exact check:", name, "===")
	fmt.Printf("top LogN=%d, leaf LogN=%d\n", top.LogN(), leaf.LogN())
	fmt.Printf("Q count: top=%d, leaf=%d\n", len(topQ.SubRings), len(leafQ.SubRings))

	if len(topQ.SubRings) != len(leafQ.SubRings) {
		ok = false
		fmt.Println("Q exact check: FAIL. Top and leaf Q counts differ.")
	} else {
		for i := range topQ.SubRings {
			qa := topQ.SubRings[i].Modulus
			qb := leafQ.SubRings[i].Modulus
			same := qa == qb
			ok = ok && same
			if i < 8 || i >= len(leafQ.SubRings)-8 || !same {
				fmt.Printf("Q[%02d]: top=%d leaf=%d same=%v\n", i, qa, qb, same)
			}
		}
	}

	topP := top.RingP()
	leafP := leaf.RingP()
	topPCount := 0
	leafPCount := 0
	if topP != nil {
		topPCount = len(topP.SubRings)
	}
	if leafP != nil {
		leafPCount = len(leafP.SubRings)
	}
	fmt.Printf("P count: top=%d, leaf=%d\n", topPCount, leafPCount)

	if topPCount != leafPCount {
		ok = false
		fmt.Println("P exact check: FAIL. Top and leaf P counts differ.")
	} else if topP != nil {
		for i := range topP.SubRings {
			pa := topP.SubRings[i].Modulus
			pb := leafP.SubRings[i].Modulus
			same := pa == pb
			ok = ok && same
			fmt.Printf("P[%02d]: top=%d leaf=%d same=%v\n", i, pa, pb, same)
		}
	}

	if ok {
		fmt.Println("Ring-basis exact check: PASS. Top and leaf Q/P moduli match exactly.")
	} else {
		fmt.Println("Ring-basis exact check: FAIL. Top and leaf Q/P moduli do not match exactly.")
	}
	return ok
}

func makeDirectEvaluator(logN int, sk *rlwe.SecretKey, cfg ExpConfig) (ckks.Parameters, bootstrapping.Parameters, *bootstrapping.Evaluator) {
	_, btpParams := makeBootstrappingParams(logN, cfg)
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

func makeEvaluatorFromBootstrappingParams(name string, btpParams bootstrapping.Parameters, sk *rlwe.SecretKey) *bootstrapping.Evaluator {
	paramsBoot := btpParams.BootstrappingParameters
	fmt.Printf("Generating %s bootstrapping keys for LogN=%d, bootLogQP=%.1f, maxLevel=%d...\n",
		name, paramsBoot.LogN(), paramsBoot.LogQP(), paramsBoot.MaxLevel())

	btpKeys, _, err := btpParams.GenEvaluationKeys(sk)
	if err != nil {
		panic(fmt.Errorf("%s GenEvaluationKeys LogN=%d failed: %w", name, paramsBoot.LogN(), err))
	}

	btpEval, err := bootstrapping.NewEvaluator(btpParams, btpKeys)
	if err != nil {
		panic(fmt.Errorf("%s NewEvaluator LogN=%d failed: %w", name, paramsBoot.LogN(), err))
	}

	return btpEval
}

func makeLeafEvaluatorPoolNoBK(topBootParams ckks.Parameters, leafLogN int, skLeaf *rlwe.SecretKey, workers int, cfg ExpConfig, leafQCountOverride int) (ckks.Parameters, bootstrapping.Parameters, []leafBootstrapper) {
	if workers <= 0 {
		panic("workers must be positive")
	}

	_, leafBtpParams := makeBootstrappingParams(leafLogN, cfg)
	leafParams := makeParamsFromExactBasis(topBootParams, leafLogN)
	if leafQCountOverride > 0 {
		leafParams = makeParamsFromPrefixBasis(topBootParams, leafLogN, leafQCountOverride, cfg.NumP)
	}

	// Important: no B_k and no B_k^{-1}; do not reserve any B-related levels.
	leafBtpParams.ResidualParameters = leafParams
	leafBtpParams.BootstrappingParameters = leafParams

	paramsBoot := leafBtpParams.BootstrappingParameters

	if leafQCountOverride > 0 {
		fmt.Printf("Generating no-B_k PREFIX-BASIS leaf bootstrapping keys for LogN=%d, qCount=%d, logQP=%.1f, maxLevel=%d, workers=%d...\n",
			leafLogN, leafQCountOverride, paramsBoot.LogQP(), paramsBoot.MaxLevel(), workers)
		if !reportRingBasisPrefix("top boot params vs reduced-prefix leaf params", topBootParams, paramsBoot) {
			panic("prefix-basis construction failed: leaf params are not top Q/P prefixes")
		}
	} else {
		fmt.Printf("Generating no-B_k EXACT-BASIS leaf bootstrapping keys for LogN=%d, logQP=%.1f, maxLevel=%d, workers=%d...\n",
			leafLogN, paramsBoot.LogQP(), paramsBoot.MaxLevel(), workers)
		if !reportRingBasisExact("top boot params vs exact-basis leaf params", topBootParams, paramsBoot) {
			panic("exact-basis construction failed: leaf params do not match top Q/P")
		}
	}

	btpKeys, _, err := leafBtpParams.GenEvaluationKeys(skLeaf)
	if err != nil {
		panic(fmt.Errorf("GenEvaluationKeys exact-basis LogN=%d failed: %w", leafLogN, err))
	}

	evals := make([]leafBootstrapper, workers)
	for i := 0; i < workers; i++ {
		ev, err := bootstrapping.NewEvaluator(leafBtpParams, btpKeys)
		if err != nil {
			panic(fmt.Errorf("NewEvaluator exact-basis LogN=%d worker=%d failed: %w", leafLogN, i, err))
		}
		evals[i] = lattigoLeafBootstrapper{Evaluator: ev}
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

func printLeafDecodedStats(stage string, values [][]complex128) {
	fmt.Println()
	fmt.Printf("=== %s decoded magnitude stats ===\n", stage)
	for i, vals := range values {
		maxVal, maxIdx := maxAbs(vals)
		var sum float64
		for _, v := range vals {
			sum += cmplx.Abs(v)
		}
		mean := 0.0
		if len(vals) > 0 {
			mean = sum / float64(len(vals))
		}
		first := complex(0, 0)
		if len(vals) > 0 {
			first = vals[0]
		}
		fmt.Printf("%s leaf[%d]: meanAbs=%.6e maxAbs=%.6e maxIdx=%d first=%+.8e%+.8ei\n",
			stage, i, mean, maxVal, maxIdx, real(first), imag(first))
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

func dropCiphertextsToLevel(cts []*rlwe.Ciphertext, targetLevel int) error {
	if targetLevel < 0 {
		return fmt.Errorf("targetLevel must be non-negative")
	}
	for i, ct := range cts {
		if ct == nil {
			return fmt.Errorf("nil ciphertext at index %d", i)
		}
		if ct.Level() < targetLevel {
			return fmt.Errorf("ciphertext[%d] level=%d is below target leaf level=%d", i, ct.Level(), targetLevel)
		}
		if ct.Level() > targetLevel {
			ct.Resize(ct.Degree(), targetLevel)
		}
	}
	return nil
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

func c2SOnly(eval leafBootstrapper, ct *rlwe.Ciphertext) (*rlwe.Ciphertext, *rlwe.Ciphertext, error) {
	ctReal, ctImag, err := eval.CoeffsToSlots(ct.CopyNew())
	if err != nil {
		return nil, nil, fmt.Errorf("CoeffsToSlots failed: %w", err)
	}
	if ctReal == nil || ctImag == nil {
		return nil, nil, fmt.Errorf("CoeffsToSlots returned nil output: real nil=%v imag nil=%v", ctReal == nil, ctImag == nil)
	}
	return ctReal, ctImag, nil
}

func c2SOnlyLattigo(eval *bootstrapping.Evaluator, ct *rlwe.Ciphertext) (*rlwe.Ciphertext, *rlwe.Ciphertext, error) {
	return c2SOnly(lattigoLeafBootstrapper{Evaluator: eval}, ct)
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

	tr.C2SReal, tr.C2SImag, err = c2SOnlyLattigo(eval, tr.Boot.CopyNew())
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

func leafC2SOnlyParallel(rpEval *rlwe.RingPackingEvaluator, leafEvals []leafBootstrapper, ct *rlwe.Ciphertext, layers int) ([]*rlwe.Ciphertext, complexLeafPair, error) {
	if len(leafEvals) == 0 {
		return nil, complexLeafPair{}, fmt.Errorf("empty leaf evaluator pool")
	}

	leaves, err := splitTree(rpEval, ct.CopyNew(), layers)
	if err != nil {
		return nil, complexLeafPair{}, err
	}
	leafMaxLevel := leafEvals[0].MaxBootLevel()
	if err := dropCiphertextsToLevel(leaves, leafMaxLevel); err != nil {
		return nil, complexLeafPair{}, fmt.Errorf("drop split leaves to leaf level: %w", err)
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
		go func(eval leafBootstrapper) {
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

func leafModUpAndC2SAfterSplitParallel(rpEval *rlwe.RingPackingEvaluator, leafEvals []leafBootstrapper, ctScaled *rlwe.Ciphertext, layers int) ([]*rlwe.Ciphertext, complexLeafPair, error) {
	if len(leafEvals) == 0 {
		return nil, complexLeafPair{}, fmt.Errorf("empty leaf evaluator pool")
	}

	leaves, err := splitTree(rpEval, ctScaled.CopyNew(), layers)
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

	var wg sync.WaitGroup
	workerCount := len(leafEvals)
	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		eval := leafEvals[w]
		go func() {
			defer wg.Done()
			for jb := range jobs {
				bootLeaf, err := eval.ModUp(leaves[jb.idx].CopyNew())
				if err != nil {
					results <- result{idx: jb.idx, err: fmt.Errorf("leaf %d ModUp failed: %w", jb.idx, err)}
					continue
				}

				real, imag, err := c2SOnly(eval, bootLeaf)
				if err != nil {
					results <- result{idx: jb.idx, err: fmt.Errorf("leaf %d C2S failed: %w", jb.idx, err)}
					continue
				}

				results <- result{idx: jb.idx, real: real, imag: imag}
			}
		}()
	}

	go func() {
		for i := range leaves {
			jobs <- job{idx: i}
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()

	var firstErr error
	for res := range results {
		if res.err != nil && firstErr == nil {
			firstErr = res.err
			continue
		}
		outRealLeaves[res.idx] = res.real
		outImagLeaves[res.idx] = res.imag
	}
	if firstErr != nil {
		return leaves, complexLeafPair{}, firstErr
	}
	for i := range leaves {
		if outRealLeaves[i] == nil || outImagLeaves[i] == nil {
			return leaves, complexLeafPair{}, fmt.Errorf("missing leaf ModUp+C2S output at leaf %d", i)
		}
	}

	return leaves, complexLeafPair{Real: outRealLeaves, Imag: outImagLeaves}, nil
}

func leafScaleDownModUpAndC2SAfterSplitParallel(rpEval *rlwe.RingPackingEvaluator, leafEvals []leafBootstrapper, ct *rlwe.Ciphertext, layers int) ([]*rlwe.Ciphertext, []*rlwe.Ciphertext, complexLeafPair, error) {
	if len(leafEvals) == 0 {
		return nil, nil, complexLeafPair{}, fmt.Errorf("empty leaf evaluator pool")
	}

	leaves, err := splitTree(rpEval, ct.CopyNew(), layers)
	if err != nil {
		return nil, nil, complexLeafPair{}, err
	}

	scaledLeaves := make([]*rlwe.Ciphertext, len(leaves))
	outRealLeaves := make([]*rlwe.Ciphertext, len(leaves))
	outImagLeaves := make([]*rlwe.Ciphertext, len(leaves))

	type job struct{ idx int }
	type result struct {
		idx    int
		scaled *rlwe.Ciphertext
		real   *rlwe.Ciphertext
		imag   *rlwe.Ciphertext
		err    error
	}

	jobs := make(chan job)
	results := make(chan result, len(leaves))

	var wg sync.WaitGroup
	workerCount := len(leafEvals)
	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		eval := leafEvals[w]
		go func() {
			defer wg.Done()
			for jb := range jobs {
				scaled, _, err := eval.ScaleDown(leaves[jb.idx].CopyNew())
				if err != nil {
					results <- result{idx: jb.idx, err: fmt.Errorf("leaf %d ScaleDown failed: %w", jb.idx, err)}
					continue
				}

				bootLeaf, err := eval.ModUp(scaled.CopyNew())
				if err != nil {
					results <- result{idx: jb.idx, err: fmt.Errorf("leaf %d ModUp failed: %w", jb.idx, err)}
					continue
				}

				real, imag, err := c2SOnly(eval, bootLeaf)
				if err != nil {
					results <- result{idx: jb.idx, err: fmt.Errorf("leaf %d C2S failed: %w", jb.idx, err)}
					continue
				}

				results <- result{idx: jb.idx, scaled: scaled, real: real, imag: imag}
			}
		}()
	}

	go func() {
		for i := range leaves {
			jobs <- job{idx: i}
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()

	var firstErr error
	for res := range results {
		if res.err != nil && firstErr == nil {
			firstErr = res.err
			continue
		}
		scaledLeaves[res.idx] = res.scaled
		outRealLeaves[res.idx] = res.real
		outImagLeaves[res.idx] = res.imag
	}
	if firstErr != nil {
		return leaves, nil, complexLeafPair{}, firstErr
	}
	for i := range leaves {
		if scaledLeaves[i] == nil || outRealLeaves[i] == nil || outImagLeaves[i] == nil {
			return leaves, nil, complexLeafPair{}, fmt.Errorf("missing leaf ScaleDown+ModUp+C2S output at leaf %d", i)
		}
	}

	return leaves, scaledLeaves, complexLeafPair{Real: outRealLeaves, Imag: outImagLeaves}, nil
}

func leafScaleDownParallel(leafEvals []leafBootstrapper, leaves []*rlwe.Ciphertext) ([]*rlwe.Ciphertext, error) {
	if len(leafEvals) == 0 {
		return nil, fmt.Errorf("empty leaf evaluator pool")
	}

	scaledLeaves := make([]*rlwe.Ciphertext, len(leaves))

	type job struct{ idx int }
	type result struct {
		idx    int
		scaled *rlwe.Ciphertext
		err    error
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
		go func(eval leafBootstrapper) {
			defer wg.Done()
			for jb := range jobs {
				scaled, _, err := eval.ScaleDown(leaves[jb.idx].CopyNew())
				if err != nil {
					results <- result{idx: jb.idx, err: fmt.Errorf("leaf %d ScaleDown failed: %w", jb.idx, err)}
					continue
				}
				results <- result{idx: jb.idx, scaled: scaled}
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
			return nil, res.err
		}
		scaledLeaves[res.idx] = res.scaled
	}

	for i := range scaledLeaves {
		if scaledLeaves[i] == nil {
			return nil, fmt.Errorf("missing leaf ScaleDown output at leaf %d", i)
		}
	}

	return scaledLeaves, nil
}

func leafModUpParallel(leafEvals []leafBootstrapper, leaves []*rlwe.Ciphertext) ([]*rlwe.Ciphertext, error) {
	if len(leafEvals) == 0 {
		return nil, fmt.Errorf("empty leaf evaluator pool")
	}

	bootLeaves := make([]*rlwe.Ciphertext, len(leaves))

	type job struct{ idx int }
	type result struct {
		idx int
		ct  *rlwe.Ciphertext
		err error
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
		go func(eval leafBootstrapper) {
			defer wg.Done()
			for jb := range jobs {
				ct, err := eval.ModUp(leaves[jb.idx].CopyNew())
				if err != nil {
					results <- result{idx: jb.idx, err: fmt.Errorf("leaf %d ModUp failed: %w", jb.idx, err)}
					continue
				}
				results <- result{idx: jb.idx, ct: ct}
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
			return nil, res.err
		}
		bootLeaves[res.idx] = res.ct
	}

	for i := range bootLeaves {
		if bootLeaves[i] == nil {
			return nil, fmt.Errorf("missing leaf ModUp output at leaf %d", i)
		}
	}

	return bootLeaves, nil
}

func leafC2SFromLeavesParallel(leafEvals []leafBootstrapper, leaves []*rlwe.Ciphertext) (complexLeafPair, error) {
	if len(leafEvals) == 0 {
		return complexLeafPair{}, fmt.Errorf("empty leaf evaluator pool")
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
		go func(eval leafBootstrapper) {
			defer wg.Done()
			for jb := range jobs {
				real, imag, err := c2SOnly(eval, leaves[jb.idx].CopyNew())
				if err != nil {
					results <- result{idx: jb.idx, err: fmt.Errorf("leaf %d C2S failed: %w", jb.idx, err)}
					continue
				}
				results <- result{idx: jb.idx, real: real, imag: imag}
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
			return complexLeafPair{}, res.err
		}
		outRealLeaves[res.idx] = res.real
		outImagLeaves[res.idx] = res.imag
	}

	for i := range leaves {
		if outRealLeaves[i] == nil || outImagLeaves[i] == nil {
			return complexLeafPair{}, fmt.Errorf("missing leaf C2S output at leaf %d", i)
		}
	}

	return complexLeafPair{Real: outRealLeaves, Imag: outImagLeaves}, nil
}

func leafEvalModParallel(leafEvals []leafBootstrapper, in complexLeafPair) (complexLeafPair, error) {
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
		go func(eval leafBootstrapper) {
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

func leafS2CAndMergeParallel(rpEval *rlwe.RingPackingEvaluator, leafEvals []leafBootstrapper, in complexLeafPair) ([]*rlwe.Ciphertext, *rlwe.Ciphertext, error) {
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
		go func(eval leafBootstrapper) {
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

func leafS2CParallel(leafEvals []leafBootstrapper, in complexLeafPair) ([]*rlwe.Ciphertext, error) {
	if len(in.Real) != len(in.Imag) {
		return nil, fmt.Errorf("real/imag leaf length mismatch")
	}
	if len(leafEvals) == 0 {
		return nil, fmt.Errorf("empty leaf evaluator pool")
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
		go func(eval leafBootstrapper) {
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
			return nil, res.err
		}
		coeffLeaves[res.idx] = res.ct
	}

	for i := range coeffLeaves {
		if coeffLeaves[i] == nil {
			return nil, fmt.Errorf("missing leaf S2C output at leaf %d", i)
		}
	}

	return coeffLeaves, nil
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
	leafEvals []leafBootstrapper,
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

// printBootstrapDebug prints precision statistics in AnyuWang 2025 output format.
// divK: divide decoded test values by this factor before comparison
//
//	(use 1.0 for no scaling; use float64(k) for the factored DIV-K comparison).
func printBootstrapDebug(
	params ckks.Parameters,
	ct *rlwe.Ciphertext,
	valuesWant []complex128,
	decryptor *rlwe.Decryptor,
	encoder *ckks.Encoder,
	divK float64,
) (valuesTest []complex128, logL1Err float64) {
	slots := ct.Slots()
	if !ct.IsBatched {
		slots *= 2
	}
	valuesTest = make([]complex128, slots)
	if err := encoder.Decode(decryptor.DecryptNew(ct), valuesTest); err != nil {
		panic(err)
	}
	if divK != 1 {
		inv := complex(1/divK, 0)
		for i := range valuesTest {
			valuesTest[i] *= inv
		}
	}

	fmt.Println()
	fmt.Printf("Level: %d (logQ = %d)\n", ct.Level(), params.LogQLvl(ct.Level()))
	fmt.Printf("Scale: 2^%f\n", math.Log2(ct.Scale.Float64()))
	fmt.Printf("ValuesTest: %10.14f %10.14f %10.14f %10.14f...\n",
		valuesTest[0], valuesTest[1], valuesTest[2], valuesTest[3])
	fmt.Printf("ValuesWant: %10.14f %10.14f %10.14f %10.14f...\n",
		valuesWant[0], valuesWant[1], valuesWant[2], valuesWant[3])

	precStats := ckks.GetPrecisionStats(params, encoder, nil, valuesWant, valuesTest, 0, false)

	// Log2L1Err not available in Lattigo v6.2.0; compute L1 precision manually.
	var l1Sum float64
	for i := range valuesTest {
		d := valuesTest[i] - valuesWant[i]
		l1Sum += math.Abs(real(d)) + math.Abs(imag(d))
	}
	logL1Err = -math.Log2(l1Sum / (2 * float64(len(valuesTest))))
	fmt.Printf("L1 Err = %f\n", logL1Err)
	fmt.Println(precStats.String())
	fmt.Println()
	return
}

func l1PrecisionBits(want, got []complex128) float64 {
	if len(want) != len(got) {
		panic(fmt.Errorf("l1PrecisionBits length mismatch: want=%d got=%d", len(want), len(got)))
	}
	var l1Sum float64
	for i := range want {
		d := got[i] - want[i]
		l1Sum += math.Abs(real(d)) + math.Abs(imag(d))
	}
	avg := l1Sum / (2 * float64(len(want)))
	if avg == 0 {
		return math.Inf(1)
	}
	return -math.Log2(avg)
}

func runLattigoSplitDropQSweep(
	paramsTop ckks.Parameters,
	encoderTop *ckks.Encoder,
	decryptorTop *rlwe.Decryptor,
	encryptorTop *rlwe.Encryptor,
	rpEval *rlwe.RingPackingEvaluator,
	values []complex128,
	layers int,
	reps int,
) {
	if reps <= 0 {
		reps = 1
	}

	maxLevel := paramsTop.MaxLevel()
	ctHigh := encryptAtLevel(paramsTop, encoderTop, encryptorTop, values, maxLevel)
	valuesHigh := decodeValues(paramsTop, encoderTop, decryptorTop, ctHigh)
	inputPrec := l1PrecisionBits(values, valuesHigh)

	fmt.Println()
	fmt.Println("=== Q-sweep: Lattigo SplitNew -> leaf Resize(drop Q) -> MergeNew ===")
	fmt.Printf("input: LogN=%d, splitLayers=%d, leaves=%d, inputLevel=%d, inputLogQ=%d, inputLogQP=%.1f, inputPrec=%.2f bits\n",
		paramsTop.LogN(), layers, 1<<uint(layers), maxLevel, paramsTop.LogQLvl(maxLevel), paramsTop.LogQP(), inputPrec)
	fmt.Println("note: this tests lowering ciphertext Q after Lattigo SplitNew. The ring-switching evaluation keys are still generated at the original full Q/P basis.")
	fmt.Println()
	fmt.Printf("| targetLevel | target logQ | saved Q bits | avg precision vs plaintext | avg precision vs input ct | split | drop | merge | total | status |\n")
	fmt.Printf("|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|\n")

	for targetLevel := maxLevel; targetLevel >= 0; targetLevel-- {
		var sumPrecPlain, sumPrecInput float64
		var sumSplit, sumDrop, sumMerge, sumTotal time.Duration
		status := "PASS"

		for rep := 0; rep < reps; rep++ {
			start := time.Now()

			t0 := time.Now()
			leaves, err := splitTree(rpEval, ctHigh.CopyNew(), layers)
			t1 := time.Now()
			if err != nil {
				status = fmt.Sprintf("Split failed: %v", err)
				break
			}

			err = dropCiphertextsToLevel(leaves, targetLevel)
			t2 := time.Now()
			if err != nil {
				status = fmt.Sprintf("Drop failed: %v", err)
				break
			}

			merged, err := mergeTree(rpEval, leaves)
			t3 := time.Now()
			if err != nil {
				status = fmt.Sprintf("Merge failed: %v", err)
				break
			}

			got := decodeValues(paramsTop, encoderTop, decryptorTop, merged)
			precPlain := l1PrecisionBits(values, got)
			precInput := l1PrecisionBits(valuesHigh, got)
			if math.IsNaN(precPlain) || math.IsInf(precPlain, -1) || math.IsNaN(precInput) || math.IsInf(precInput, -1) {
				status = "Decode failed"
				break
			}

			sumPrecPlain += precPlain
			sumPrecInput += precInput
			sumSplit += t1.Sub(t0)
			sumDrop += t2.Sub(t1)
			sumMerge += t3.Sub(t2)
			sumTotal += time.Since(start)
		}

		if status == "PASS" {
			den := float64(reps)
			saved := paramsTop.LogQLvl(maxLevel) - paramsTop.LogQLvl(targetLevel)
			fmt.Printf("| %d | %d | %d | %.2f | %.2f | %s | %s | %s | %s | %s |\n",
				targetLevel,
				paramsTop.LogQLvl(targetLevel),
				saved,
				sumPrecPlain/den,
				sumPrecInput/den,
				(sumSplit / time.Duration(reps)).String(),
				(sumDrop / time.Duration(reps)).String(),
				(sumMerge / time.Duration(reps)).String(),
				(sumTotal / time.Duration(reps)).String(),
				status)
		} else {
			saved := paramsTop.LogQLvl(maxLevel) - paramsTop.LogQLvl(targetLevel)
			fmt.Printf("| %d | %d | %d | - | - | - | - | - | - | %s |\n",
				targetLevel, paramsTop.LogQLvl(targetLevel), saved, status)
		}
	}
}

func securityLogQPLimit128(logN int) (float64, bool) {
	switch logN {
	case 15:
		return 868, true
	case 16:
		return 1747, true
	case 17:
		return 3523, true
	default:
		return 0, false
	}
}

func printParamDryRun(logN, layers int, cfg, leafCfg ExpConfig) {
	_, btpTop := makeBootstrappingParams(logN, cfg)
	paramsTop := btpTop.BootstrappingParameters

	leafLogN := logN - layers
	leafParams := makeParamsFromExactBasis(paramsTop, leafLogN)

	fmt.Println("=== Parameter dry run ===")
	fmt.Printf("top:  logN=%d N=%d qCount=%d pCount=%d logQ=%.1f logQP=%.1f\n",
		paramsTop.LogN(),
		1<<uint(paramsTop.LogN()),
		paramsTop.QCount(),
		paramsTop.PCount(),
		paramsTop.LogQ(),
		paramsTop.LogQP())
	fmt.Printf("leaf: logN=%d N=%d qCount=%d pCount=%d logQ=%.1f logQP=%.1f\n",
		leafParams.LogN(),
		1<<uint(leafParams.LogN()),
		leafParams.QCount(),
		leafParams.PCount(),
		leafParams.LogQ(),
		leafParams.LogQP())
	fmt.Printf("leaf config: q0Bits=%d circuitLevels=%d circuitPrimeBits=%d numP=%d defaultScale=%d\n",
		leafCfg.Q0Bits,
		leafCfg.CircuitLevels,
		leafCfg.CircuitPrimeBits,
		leafCfg.NumP,
		leafCfg.DefaultScaleBits)

	if limit, ok := securityLogQPLimit128(leafParams.LogN()); ok {
		logQMargin := limit - leafParams.LogQ()
		logQStatus := "PASS"
		if logQMargin < 0 {
			logQStatus = "FAIL"
		}
		logQPMargin := limit - leafParams.LogQP()
		logQPStatus := "PASS"
		if logQPMargin < 0 {
			logQPStatus = "FAIL"
		}
		fmt.Printf("128-bit reference: leaf logQ  <= %.0f, margin=%.1f bits [%s]\n", limit, logQMargin, logQStatus)
		fmt.Printf("128-bit reference: leaf logQP <= %.0f, margin=%.1f bits [%s]\n", limit, logQPMargin, logQPStatus)
	} else {
		fmt.Printf("128-bit reference: no local table entry for leaf logN=%d\n", leafParams.LogN())
	}
}

func main() {
	logN := flag.Int("logN", 16, "top ring degree log2(N)")
	layers := flag.Int("layers", 1, "recursive split layers; k=2^layers")
	reps := flag.Int("reps", 1, "number of repeated tests")
	_ = flag.Int("warmup", 0, "warmup repetitions (unused)")
	parallelLeaves := flag.Bool("parallelLeaves", true, "run leaf operations in parallel")
	leafWorkersFlag := flag.Int("leafWorkers", 0, "leaf workers; 0 = min(k, NumCPU)")
	_ = flag.String("leafOrder", "bitrev", "leaf order (unused)")
	_ = flag.Int("printSlots", 4, "print slots (unused)")
	_ = flag.Bool("timing", true, "timing (always on)")
	runMode := flag.String("run", "all",
		"which experiment(s) to run: all | baseline | ours | subKey25 | qsweep")
	leafEngine := flag.String("leafEngine", "lattigo",
		"leaf bootstrapping engine for proposed path: lattigo | anyu")
	anyuProfileFlag := flag.String("anyuProfile", "auto",
		"Anyu leaf parameter profile: auto | i1 | lowq")
	anyuSubringLogNFlag := flag.Int("anyuSubringLogN", 12,
		"target subring secret logN for Anyu leaf engine")
	anyuMaxLogQPFlag := flag.Int("anyuMaxLogQP", -1,
		"override Anyu MaxLogQP; -1 = profile default")
	anyuMaxLevelsFlag := flag.Int("anyuMaxLevels", -1,
		"override Anyu max ordinary circuit levels; -1 = profile default")
	anyuNumPFlag := flag.Int("anyuNumP", -1,
		"override Anyu number of 61-bit P primes; -1 = profile default")
	anyuQ0BitsFlag := flag.Int("anyuQ0Bits", -1,
		"override Anyu q0 bit-size; -1 = profile default")
	anyuLogScaleFlag := flag.Int("anyuLogScale", -1,
		"override Anyu default scale bit-size; -1 = profile default")
	anyuStCBitsFlag := flag.Int("anyuStCBits", -1,
		"override Anyu StC prime bit-size; -1 = profile default")
	anyuCtSBitsFlag := flag.Int("anyuCtSBits", -1,
		"override Anyu CtS prime bit-size; -1 = profile default")
	anyuEvalModBitsFlag := flag.Int("anyuEvalModBits", -1,
		"override Anyu EvalMod prime bit-size; -1 = profile default")
	anyuMod1DegreeFlag := flag.Int("anyuMod1Degree", -1,
		"override Anyu Mod1 polynomial degree; -1 = profile default")
	anyuDoubleAngleFlag := flag.Int("anyuDoubleAngle", -1,
		"override Anyu double-angle rounds; -1 = profile default")
	anyuCtSFirstDepthFlag := flag.Int("anyuCtSFirstDepth", -1,
		"override first collapsed CtS FFT depth; -1 = auto log2(N/N')")
	anyuLogFailureFlag := flag.Float64("anyuLogFailure", 0,
		"override Anyu log2 failure target; 0 = profile default")
	flagH := flag.Int("H", 0,
		"secret key Hamming weight; 0 = uniform ternary P=2/3 (AnyuWang-style dense)")
	flagQ0Bits := flag.Int("q0Bits", 34,
		"bit-size of the first residual Q prime")
	flagCircuit := flag.Int("circuitLevels", 1,
		"number of circuit primes in residual LogQ")
	flagCircuitPrimeBits := flag.Int("circuitPrimeBits", 20,
		"bit-size of each residual circuit prime in LogQ")
	flagNumP := flag.Int("numP", 1,
		"number of 61-bit P primes for key switching")
	flagDefaultScale := flag.Int("defaultScale", 25,
		"default CKKS scale bit-size")
	flagLeafCircuit := flag.Int("leafCircuitLevels", -1,
		"leaf number of circuit primes in LogQ; -1 = max(1, circuitLevels-layers)")
	flagLeafNumP := flag.Int("leafNumP", -1,
		"leaf number of 61-bit P primes for key switching; -1 = same as numP")
	flagLeafQCount := flag.Int("leafQCount", -1,
		"leaf Q prime count override for reduced-Q leaf bootstrapping; -1 = exact full top Q basis")
	dryRunParams := flag.Bool("dryRunParams", false,
		"print top/leaf bootstrapping parameters and exit before key generation")
	splitFirstModUp := flag.Bool("splitFirstModUp", false,
		"proposed path uses ScaleDown -> Split -> leaf ModUp -> leaf C2S instead of ScaleDown -> top ModUp -> Split -> leaf C2S")
	splitFirstLeafPreprocess := flag.Bool("splitFirstLeafPreprocess", false,
		"proposed path uses Split -> leaf ScaleDown -> leaf ModUp -> leaf C2S")
	stageTrace := flag.Bool("stageTrace", false,
		"print leaf metadata and decoded magnitude stats after each proposed-path stage")
	flag.Parse()

	validModes := map[string]bool{"all": true, "baseline": true, "ours": true, "subKey25": true, "qsweep": true}
	if !validModes[*runMode] {
		fmt.Fprintf(os.Stderr, "unknown -run value %q; valid options: all, baseline, ours, subKey25, qsweep\n", *runMode)
		os.Exit(1)
	}
	validLeafEngines := map[string]bool{"lattigo": true, "anyu": true}
	if !validLeafEngines[*leafEngine] {
		fmt.Fprintf(os.Stderr, "unknown -leafEngine value %q; valid options: lattigo, anyu\n", *leafEngine)
		os.Exit(1)
	}
	runBaseline := *runMode == "all" || *runMode == "baseline"
	runOurs := *runMode == "all" || *runMode == "ours"
	runSubKey25 := *runMode == "all" || *runMode == "subKey25"
	runQSweep := *runMode == "qsweep"

	nRepeat := *reps
	if nRepeat <= 0 {
		nRepeat = 1
	}

	if runSubKey25 && !runBaseline && !runOurs {
		fmt.Println()
		fmt.Println("=== subKey25: embedded AnyuWang 2025 subring-secret bootstrapping (I1 params) ===")
		a25AvgL1Err, a25AvgTotal := runAnyu25Capturing(nRepeat)
		fmt.Println()
		fmt.Println("=== Summary ===")
		fmt.Printf("subKey25  (AnyuWang embedded) Total: %s, Avg L1 prec: %.2f bits\n", a25AvgTotal, a25AvgL1Err)
		return
	}

	if *layers <= 0 {
		panic("layers must be positive")
	}
	leafLogN := *logN - *layers
	if leafLogN < 1 {
		panic("layers too large")
	}
	numLeaves := 1 << uint(*layers)
	anyuCfg := defaultAnyuLeafConfig(*anyuProfileFlag, *logN, leafLogN)
	anyuCfg.SubringLogN = *anyuSubringLogNFlag
	if anyuCfg.SubringLogN > leafLogN {
		anyuCfg.SubringLogN = leafLogN
	}
	if *anyuMaxLogQPFlag >= 0 {
		anyuCfg.MaxLogQP = *anyuMaxLogQPFlag
	}
	if *anyuMaxLevelsFlag >= 0 {
		anyuCfg.MaxLevels = *anyuMaxLevelsFlag
	}
	if *anyuNumPFlag >= 0 {
		anyuCfg.NumP = *anyuNumPFlag
	}
	if *anyuQ0BitsFlag >= 0 {
		anyuCfg.Q0Bits = *anyuQ0BitsFlag
	}
	if *anyuLogScaleFlag >= 0 {
		anyuCfg.LogScale = *anyuLogScaleFlag
	}
	if *anyuStCBitsFlag >= 0 {
		anyuCfg.StCBits = *anyuStCBitsFlag
	}
	if *anyuCtSBitsFlag >= 0 {
		anyuCfg.CtSBits = *anyuCtSBitsFlag
	}
	if *anyuEvalModBitsFlag >= 0 {
		anyuCfg.EvalModBits = *anyuEvalModBitsFlag
	}
	if *anyuMod1DegreeFlag >= 0 {
		anyuCfg.Mod1Degree = *anyuMod1DegreeFlag
	}
	if *anyuDoubleAngleFlag >= 0 {
		anyuCfg.DoubleAngle = *anyuDoubleAngleFlag
	}
	if *anyuCtSFirstDepthFlag >= 0 {
		anyuCfg.CtSFirstDepth = *anyuCtSFirstDepthFlag
	}
	if *anyuLogFailureFlag != 0 {
		anyuCfg.LogFailure = *anyuLogFailureFlag
	}

	cfg := ExpConfig{
		H:                *flagH,
		Q0Bits:           *flagQ0Bits,
		CircuitLevels:    *flagCircuit,
		CircuitPrimeBits: *flagCircuitPrimeBits,
		NumP:             *flagNumP,
		DefaultScaleBits: *flagDefaultScale,
	}
	leafCfg := cfg
	if *flagLeafCircuit >= 0 {
		leafCfg.CircuitLevels = *flagLeafCircuit
	}
	if *flagLeafNumP >= 0 {
		leafCfg.NumP = *flagLeafNumP
	}
	if leafCfg.CircuitLevels < 1 {
		panic("leafCircuitLevels must be positive")
	}
	if leafCfg.NumP < 0 {
		panic("leafNumP must be non-negative")
	}

	if cfg.H <= 0 {
		fmt.Printf("[config] leafEngine=%s, anyuProfile=%s, anyuSubringLogN=%d, Xs = uniform ternary P=2/3 (H≈%d for N=2^%d), q0Bits=%d, circuitPrimeBits=%d, defaultScale=%d, topCircuitLevels=%d, topNumP=%d, leafCircuitLevels=%d, leafNumP=%d, leafQCount=%d, splitFirstModUp=%v, splitFirstLeafPreprocess=%v, stageTrace=%v\n",
			*leafEngine, anyuCfg.Profile, anyuCfg.SubringLogN, int(float64(int(1)<<uint(*logN))*2.0/3.0), *logN, cfg.Q0Bits, cfg.CircuitPrimeBits, cfg.DefaultScaleBits, cfg.CircuitLevels, cfg.NumP, leafCfg.CircuitLevels, leafCfg.NumP, *flagLeafQCount, *splitFirstModUp, *splitFirstLeafPreprocess, *stageTrace)
	} else {
		fmt.Printf("[config] leafEngine=%s, anyuProfile=%s, anyuSubringLogN=%d, Xs = sparse H=%d, q0Bits=%d, circuitPrimeBits=%d, defaultScale=%d, topCircuitLevels=%d, topNumP=%d, leafCircuitLevels=%d, leafNumP=%d, leafQCount=%d, splitFirstModUp=%v, splitFirstLeafPreprocess=%v, stageTrace=%v\n",
			*leafEngine, anyuCfg.Profile, anyuCfg.SubringLogN, cfg.H, cfg.Q0Bits, cfg.CircuitPrimeBits, cfg.DefaultScaleBits, cfg.CircuitLevels, cfg.NumP, leafCfg.CircuitLevels, leafCfg.NumP, *flagLeafQCount, *splitFirstModUp, *splitFirstLeafPreprocess, *stageTrace)
	}
	if *leafEngine == "anyu" {
		fmt.Printf("[AnyuConfig] MaxLogQP=%d MaxLevels=%d NumP=%d Q0Bits=%d LogScale=%d StCBits=%d CtSBits=%d EvalModBits=%d Mod1Degree=%d DoubleAngle=%d CtSFirstDepth=%d LogFailure=%.1f\n",
			anyuCfg.MaxLogQP, anyuCfg.MaxLevels, anyuCfg.NumP, anyuCfg.Q0Bits, anyuCfg.LogScale, anyuCfg.StCBits, anyuCfg.CtSBits, anyuCfg.EvalModBits, anyuCfg.Mod1Degree, anyuCfg.DoubleAngle, anyuCfg.CtSFirstDepth, anyuCfg.LogFailure)
	}

	if *dryRunParams {
		if *leafEngine == "anyu" {
			printAnyuParamDryRun(*logN, *layers, anyuCfg)
		} else {
			printParamDryRun(*logN, *layers, cfg, leafCfg)
		}
		return
	}

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

	//============================
	//=== 1) SCHEME PARAMETERS ===
	//============================

	var btpTop bootstrapping.Parameters
	if *leafEngine == "anyu" {
		var err error
		_, btpTop, err = makeAnyuBootstrappingParamsWithConfig(*logN, nil, anyuCfg)
		if err != nil {
			panic(fmt.Errorf("makeAnyuBootstrappingParams top LogN=%d failed: %w", *logN, err))
		}
		if !*splitFirstLeafPreprocess {
			fmt.Println("[warning] -leafEngine=anyu is intended for -splitFirstLeafPreprocess; other schedules may not match Anyu normalization.")
		}
	} else {
		_, btpTop = makeBootstrappingParams(*logN, cfg)
	}
	paramsTop := btpTop.BootstrappingParameters

	fmt.Printf("Bootstrapping parameters: logN=%d, logSlots=%d, H(%d; %d), ExtDegree=%d, sigma=%v, logQP=%f, levels=%d, scale=2^%d\n",
		paramsTop.LogN(),
		paramsTop.LogMaxSlots(),
		paramsTop.XsHammingWeight(),
		btpTop.EphemeralSecretWeight,
		btpTop.SubringSecretExtDegree,
		paramsTop.Xe(),
		paramsTop.LogQP(),
		paramsTop.QCount(),
		paramsTop.LogDefaultScale())

	//===========================
	//=== 2) KEYGEN & ENCRYPT ===
	//===========================

	kgenTop := rlwe.NewKeyGenerator(paramsTop)
	skTop, pkTop := kgenTop.GenKeyPairNew()
	encoderTop := ckks.NewEncoder(paramsTop)
	encryptorTop := rlwe.NewEncryptor(paramsTop, pkTop)
	decryptorTop := rlwe.NewDecryptor(paramsTop, skTop)

	rpk := &rlwe.RingPackingEvaluationKey{}
	skMap, err := rpk.GenRingSwitchingKeys(paramsTop, skTop, leafLogN, rlwe.EvaluationKeyParameters{})
	if err != nil {
		panic(fmt.Errorf("GenRingSwitchingKeys: %w", err))
	}
	rpEval := rlwe.NewRingPackingEvaluator(rpk)
	skLeaf, ok := skMap[leafLogN]
	if !ok {
		panic(fmt.Errorf("missing leaf sk LogN=%d", leafLogN))
	}

	fmt.Println()
	fmt.Println("Generating bootstrapping evaluation keys...")
	var topEval *bootstrapping.Evaluator
	needTopEval := runBaseline || (runOurs && !*splitFirstLeafPreprocess)
	if needTopEval {
		topEval = makeEvaluatorFromBootstrappingParams("direct top", btpTop, skTop)
	}
	var leafEvals []leafBootstrapper
	var paramsLeaf ckks.Parameters
	var encoderLeaf *ckks.Encoder
	var decryptorLeaf *rlwe.Decryptor
	if runOurs {
		if *leafEngine == "anyu" {
			paramsLeaf, _, leafEvals = makeAnyuLeafEvaluatorPool(paramsTop, leafLogN, skLeaf, leafWorkers, anyuCfg)
		} else {
			paramsLeaf, _, leafEvals = makeLeafEvaluatorPoolNoBK(paramsTop, leafLogN, skLeaf, leafWorkers, leafCfg, *flagLeafQCount)
		}
		encoderLeaf = ckks.NewEncoder(paramsLeaf)
		decryptorLeaf = rlwe.NewDecryptor(paramsLeaf, skLeaf)
	}
	fmt.Println("Done")

	//========================
	//=== 3) BOOTSTRAPPING ===
	//========================

	values := makeValues(paramsTop.MaxSlots())
	ctLow := encryptAtLevel(paramsTop, encoderTop, encryptorTop, values, 0)

	fmt.Println()
	fmt.Println("Precision of values vs. ciphertext")
	valuesTest, _ := printBootstrapDebug(paramsTop, ctLow, values, decryptorTop, encoderTop, 1)

	if runQSweep {
		runLattigoSplitDropQSweep(paramsTop, encoderTop, decryptorTop, encryptorTop, rpEval, values, *layers, nRepeat)
		return
	}

	//===========================================
	//=== BASELINE: direct single-ring bootstrap
	//===========================================
	var (
		bAvgL1Err                                           float64
		bAvgModUp, bAvgCtS, bAvgEvalMod, bAvgStC, bAvgTotal time.Duration
		bRan                                                bool
	)
	if runBaseline {
		bRan = true
		fmt.Println()
		fmt.Printf("=== Baseline: direct single-ring bootstrapping (k=1) ===\n")
		fmt.Println("Bootstrapping...")

		for repeat := 0; repeat < nRepeat; repeat++ {
			_ = repeat
			tStart := time.Now()

			_, ctBoot, err := prepareBootstrapInput(topEval, ctLow.CopyNew())
			if err != nil {
				panic(err)
			}
			t1 := time.Now()

			ctReal, ctImag, err := c2SOnlyLattigo(topEval, ctBoot)
			if err != nil {
				panic(err)
			}
			t2 := time.Now()

			ctReal, err = topEval.EvalMod(ctReal)
			if err != nil {
				panic(err)
			}
			ctImag, err = topEval.EvalMod(ctImag)
			if err != nil {
				panic(err)
			}
			t3 := time.Now()

			outBaseline, err := topEval.SlotsToCoeffs(ctReal, ctImag)
			if err != nil {
				panic(err)
			}
			t4 := time.Now()

			modUp := t1.Sub(tStart)
			cts := t2.Sub(t1)
			evalMod := t3.Sub(t2)
			stc := t4.Sub(t3)
			total := t4.Sub(tStart)

			bAvgModUp += modUp
			bAvgCtS += cts
			bAvgEvalMod += evalMod
			bAvgStC += stc
			bAvgTotal += total

			fmt.Printf("Bootstrapping finished in %s\nStC: %s, ModUp: %s, CtS: %s, EvalMod: %s\n",
				total, stc, modUp, cts, evalMod)
			fmt.Println("Done")
			fmt.Printf("Output ciphertext level = %d, circuit moduli = %v\n",
				outBaseline.Level(), paramsTop.LogQi()[1:outBaseline.Level()+1])

			fmt.Println()
			fmt.Println("Precision of ciphertext vs. Bootstrap(ciphertext)")
			_, logL1 := printBootstrapDebug(paramsTop, outBaseline, valuesTest, decryptorTop, encoderTop, 1)
			bAvgL1Err += math.Pow(2, -logL1)
		}

		bAvgL1Err = -math.Log2(bAvgL1Err / float64(nRepeat))
		bAvgModUp /= time.Duration(nRepeat)
		bAvgCtS /= time.Duration(nRepeat)
		bAvgEvalMod /= time.Duration(nRepeat)
		bAvgStC /= time.Duration(nRepeat)
		bAvgTotal /= time.Duration(nRepeat)

		fmt.Printf("Avg L1 error is %f\n", bAvgL1Err)
		fmt.Printf("Avg time... ModUp: %s, CtS: %s, EvalMod: %s, StC: %s, Total: %s\n",
			bAvgModUp, bAvgCtS, bAvgEvalMod, bAvgStC, bAvgTotal)
	} // end runBaseline

	//==============================================
	//=== PROPOSED: factored no-B_k bootstrapping
	//==============================================
	var (
		avgL1Err float64
		avgSplit, avgScaleDown, avgModUp, avgCtS,
		avgEvalMod, avgStC, avgMerge, avgTotal time.Duration
		oursRan bool
	)
	if runOurs {
		oursRan = true
		fmt.Println()
		fmt.Printf("=== Proposed: factored no-B_k bootstrapping (splitLayers=%d, leaves=%d) ===\n", *layers, numLeaves)
		fmt.Println("Bootstrapping...")

		for repeat := 0; repeat < nRepeat; repeat++ {
			_ = repeat
			tStart := time.Now()

			var pair complexLeafPair
			var err error
			if *splitFirstLeafPreprocess {
				leaves, err := splitTree(rpEval, ctLow.CopyNew(), *layers)
				if err != nil {
					panic(err)
				}
				tSplit := time.Now()
				if *stageTrace {
					printLeafCTInfos("trace split leaves", leaves)
					printLeafDecodedStats("trace split leaves", decodeLeafOutputs(paramsLeaf, encoderLeaf, decryptorLeaf, leaves))
				}

				scaledLeaves, err := leafScaleDownParallel(leafEvals, leaves)
				if err != nil {
					panic(err)
				}
				tScaleDown := time.Now()
				if *stageTrace {
					printLeafCTInfos("trace leaf ScaleDown", scaledLeaves)
					printLeafDecodedStats("trace leaf ScaleDown", decodeLeafOutputs(paramsLeaf, encoderLeaf, decryptorLeaf, scaledLeaves))
				}

				bootLeaves, err := leafModUpParallel(leafEvals, scaledLeaves)
				if err != nil {
					panic(err)
				}
				tModUp := time.Now()
				if *stageTrace {
					printLeafCTInfos("trace leaf ModUp", bootLeaves)
					printLeafDecodedStats("trace leaf ModUp", decodeLeafOutputs(paramsLeaf, encoderLeaf, decryptorLeaf, bootLeaves))
				}

				pair, err = leafC2SFromLeavesParallel(leafEvals, bootLeaves)
				if err != nil {
					panic(err)
				}
				tC2S := time.Now()
				if *stageTrace {
					printLeafCTInfos("trace leaf C2S real", pair.Real)
					printLeafDecodedStats("trace leaf C2S real", decodeLeafOutputs(paramsLeaf, encoderLeaf, decryptorLeaf, pair.Real))
					printLeafCTInfos("trace leaf C2S imag", pair.Imag)
					printLeafDecodedStats("trace leaf C2S imag", decodeLeafOutputs(paramsLeaf, encoderLeaf, decryptorLeaf, pair.Imag))
				}

				pair, err = leafEvalModParallel(leafEvals, pair)
				if err != nil {
					panic(err)
				}
				tEvalMod := time.Now()
				if *stageTrace {
					printLeafCTInfos("trace leaf EvalMod real", pair.Real)
					printLeafDecodedStats("trace leaf EvalMod real", decodeLeafOutputs(paramsLeaf, encoderLeaf, decryptorLeaf, pair.Real))
					printLeafCTInfos("trace leaf EvalMod imag", pair.Imag)
					printLeafDecodedStats("trace leaf EvalMod imag", decodeLeafOutputs(paramsLeaf, encoderLeaf, decryptorLeaf, pair.Imag))
				}

				coeffLeaves, err := leafS2CParallel(leafEvals, pair)
				if err != nil {
					panic(err)
				}
				tS2C := time.Now()
				if *stageTrace {
					printLeafCTInfos("trace leaf S2C", coeffLeaves)
					printLeafDecodedStats("trace leaf S2C", decodeLeafOutputs(paramsLeaf, encoderLeaf, decryptorLeaf, coeffLeaves))
				}

				out, err := mergeTree(rpEval, coeffLeaves)
				if err != nil {
					panic(err)
				}
				tMerge := time.Now()

				split := tSplit.Sub(tStart)
				scaleDown := tScaleDown.Sub(tSplit)
				modUp := tModUp.Sub(tScaleDown)
				cts := tC2S.Sub(tModUp)
				evalMod := tEvalMod.Sub(tC2S)
				stc := tS2C.Sub(tEvalMod)
				merge := tMerge.Sub(tS2C)
				total := tMerge.Sub(tStart)

				avgSplit += split
				avgScaleDown += scaleDown
				avgModUp += modUp
				avgCtS += cts
				avgEvalMod += evalMod
				avgStC += stc
				avgMerge += merge
				avgTotal += total

				fmt.Printf("Bootstrapping finished in %s\nSplit: %s, leafScaleDown: %s, leafModUp: %s, leafC2S: %s, leafEvalMod: %s, leafS2C: %s, Merge: %s\n",
					total, split, scaleDown, modUp, cts, evalMod, stc, merge)
				fmt.Println("Done")
				fmt.Printf("Output ciphertext level = %d, circuit moduli = %v\n",
					out.Level(), paramsTop.LogQi()[1:out.Level()+1])

				fmt.Println()
				fmt.Println("Precision of ciphertext vs. Bootstrap(ciphertext)")
				_, logL1 := printBootstrapDebug(paramsTop, out, valuesTest, decryptorTop, encoderTop, 1)
				avgL1Err += math.Pow(2, -logL1)
				continue
			}

			if *splitFirstModUp {
				ctScaled, _, err := topEval.ScaleDown(ctLow.CopyNew())
				if err != nil {
					panic(fmt.Errorf("ScaleDown failed: %w", err))
				}
				t1 := time.Now()

				_, pair, err = leafModUpAndC2SAfterSplitParallel(rpEval, leafEvals, ctScaled, *layers)
				if err != nil {
					panic(err)
				}
				t2 := time.Now()

				pair, err = leafEvalModParallel(leafEvals, pair)
				if err != nil {
					panic(err)
				}
				t3 := time.Now()

				_, out, err := leafS2CAndMergeParallel(rpEval, leafEvals, pair)
				if err != nil {
					panic(err)
				}
				t4 := time.Now()

				modUp := t2.Sub(tStart)
				cts := t2.Sub(t1)
				evalMod := t3.Sub(t2)
				stc := t4.Sub(t3)
				total := t4.Sub(tStart)

				avgModUp += modUp
				avgCtS += cts
				avgEvalMod += evalMod
				avgStC += stc
				avgTotal += total

				fmt.Printf("Bootstrapping finished in %s\nStC: %s, ScaleDown+Split+leafModUp+C2S: %s, EvalMod: %s\n",
					total, stc, modUp, evalMod)
				fmt.Println("Done")
				fmt.Printf("Output ciphertext level = %d, circuit moduli = %v\n",
					out.Level(), paramsTop.LogQi()[1:out.Level()+1])

				fmt.Println()
				fmt.Println("Precision of ciphertext vs. Bootstrap(ciphertext)")
				_, logL1 := printBootstrapDebug(paramsTop, out, valuesTest, decryptorTop, encoderTop, 1)
				avgL1Err += math.Pow(2, -logL1)
				continue
			}

			_, ctBoot, err := prepareBootstrapInput(topEval, ctLow.CopyNew())
			if err != nil {
				panic(err)
			}
			t1 := time.Now()

			_, pair, err = leafC2SOnlyParallel(rpEval, leafEvals, ctBoot, *layers)
			if err != nil {
				panic(err)
			}
			t2 := time.Now()

			pair, err = leafEvalModParallel(leafEvals, pair)
			if err != nil {
				panic(err)
			}
			t3 := time.Now()

			_, out, err := leafS2CAndMergeParallel(rpEval, leafEvals, pair)
			if err != nil {
				panic(err)
			}
			t4 := time.Now()

			modUp := t1.Sub(tStart)
			cts := t2.Sub(t1)
			evalMod := t3.Sub(t2)
			stc := t4.Sub(t3)
			total := t4.Sub(tStart)

			avgModUp += modUp
			avgCtS += cts
			avgEvalMod += evalMod
			avgStC += stc
			avgTotal += total

			fmt.Printf("Bootstrapping finished in %s\nStC: %s, ModUp: %s, CtS: %s, EvalMod: %s\n",
				total, stc, modUp, cts, evalMod)
			fmt.Println("Done")
			fmt.Printf("Output ciphertext level = %d, circuit moduli = %v\n",
				out.Level(), paramsTop.LogQi()[1:out.Level()+1])

			fmt.Println()
			fmt.Println("Precision of ciphertext vs. Bootstrap(ciphertext)")
			_, logL1 := printBootstrapDebug(paramsTop, out, valuesTest, decryptorTop, encoderTop, 1)
			avgL1Err += math.Pow(2, -logL1)
		}

		avgL1Err = -math.Log2(avgL1Err / float64(nRepeat))
		avgSplit /= time.Duration(nRepeat)
		avgScaleDown /= time.Duration(nRepeat)
		avgModUp /= time.Duration(nRepeat)
		avgCtS /= time.Duration(nRepeat)
		avgEvalMod /= time.Duration(nRepeat)
		avgStC /= time.Duration(nRepeat)
		avgMerge /= time.Duration(nRepeat)
		avgTotal /= time.Duration(nRepeat)

		fmt.Printf("Avg L1 error is %f\n", avgL1Err)
		if *splitFirstLeafPreprocess {
			fmt.Printf("Avg time... Split: %s, leafScaleDown: %s, leafModUp: %s, leafC2S: %s, leafEvalMod: %s, leafS2C: %s, Merge: %s, Total: %s\n",
				avgSplit, avgScaleDown, avgModUp, avgCtS, avgEvalMod, avgStC, avgMerge, avgTotal)
		} else {
			fmt.Printf("Avg time... ModUp: %s, CtS: %s, EvalMod: %s, StC: %s, Total: %s\n",
				avgModUp, avgCtS, avgEvalMod, avgStC, avgTotal)
		}
	} // end runOurs

	//==============================================
	//=== subKey25: AnyuWang 2025 subring-secret
	//==============================================
	var (
		a25AvgL1Err float64
		a25AvgTotal time.Duration
		a25Ran      bool
	)
	if runSubKey25 {
		a25Ran = true
		fmt.Println()
		fmt.Println("=== subKey25: AnyuWang 2025 subring-secret bootstrapping (I1 params) ===")
		// RunAnyu25 prints its own per-repeat lines and avg summary;
		// we capture total/L1 for the cross-method Summary below.
		a25AvgL1Err, a25AvgTotal = runAnyu25Capturing(nRepeat)
	}

	//=================================
	//=== SUMMARY
	//=================================
	fmt.Println()
	fmt.Println("=== Summary ===")
	if bRan {
		fmt.Printf("Baseline  (k=1)      Total: %s, Avg L1 prec: %.2f bits\n", bAvgTotal, bAvgL1Err)
	}
	if oursRan {
		fmt.Printf("Proposed  (splitLayers=%d, leaves=%d) Total: %s, Avg L1 prec: %.2f bits\n", *layers, numLeaves, avgTotal, avgL1Err)
	}
	if a25Ran {
		fmt.Printf("subKey25  (AnyuWang) Total: %s, Avg L1 prec: %.2f bits\n", a25AvgTotal, a25AvgL1Err)
	}
	if bRan && oursRan {
		fmt.Printf("Speedup baseline→proposed: %.2fx\n", float64(bAvgTotal)/float64(avgTotal))
	}
	if a25Ran && oursRan {
		fmt.Printf("Speedup subKey25→proposed: %.2fx\n", float64(a25AvgTotal)/float64(avgTotal))
	}
}

func directFullTraceQuiet(eval *bootstrapping.Evaluator, ct *rlwe.Ciphertext) (*rlwe.Ciphertext, error) {
	_, ctBoot, err := prepareBootstrapInput(eval, ct.CopyNew())
	if err != nil {
		return nil, err
	}

	ctReal, ctImag, err := c2SOnlyLattigo(eval, ctBoot.CopyNew())
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
	leafEvals []leafBootstrapper,
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
