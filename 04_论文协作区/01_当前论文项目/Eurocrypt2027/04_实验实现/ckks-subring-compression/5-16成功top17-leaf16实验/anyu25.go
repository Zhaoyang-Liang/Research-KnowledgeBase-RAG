// anyu25.go embeds the AnyuWang 2025 subring-secret bootstrapping experiment in
// this package, avoiding the previous subprocess call to subkey25/main.
//
// This is still a standalone AnyuWang benchmark path: it generates its own
// parameters, keys, plaintext, and ciphertext. The next step for true stacking is
// to wrap this parameter/evaluator construction as a leafBootstrapper so it can
// consume split leaf ciphertexts from main.go.
package main

import (
	"fmt"
	"math"
	"time"

	"github.com/tuneinsight/lattigo/v6/circuits/ckks/bootstrapping"
	"github.com/tuneinsight/lattigo/v6/circuits/ckks/dft"
	"github.com/tuneinsight/lattigo/v6/circuits/ckks/mod1"
	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/failure"
	"github.com/tuneinsight/lattigo/v6/ring"
	"github.com/tuneinsight/lattigo/v6/schemes/ckks"
	"github.com/tuneinsight/lattigo/v6/utils/sampling"
)

type anyuBtpParamSet struct {
	LogDefaultScale                                   int
	q0, qiSlotsToCoeffs, qiCoeffsToSlots, piBootstrap []int
	qEvalMod                                          int
	StCLevels                                         []int
	LogMessageRatio                                   int

	K, mod1Degree, DoubleAngle int
	mod1Type                   mod1.Type

	StCSubringSecretExtDegree, bottomLevelSubringSecretExtDegree int
	Xs                                                           ring.Ternary
	isSlim                                                       bool
	EphemeralSecretWeight                                        int
}

type anyuLeafConfig struct {
	Profile       string
	SubringLogN   int
	MaxLogQP      int
	MaxLevels     int
	Q0Bits        int
	LogScale      int
	StCBits       int
	CtSBits       int
	EvalModBits   int
	NumP          int
	Mod1Degree    int
	DoubleAngle   int
	CtSFirstDepth int
	LogFailure    float64
	PrintParams   bool
	ExactTopBasis bool
}

func defaultAnyuLeafConfig(profile string, topLogN, leafLogN int) anyuLeafConfig {
	if profile == "auto" {
		if leafLogN <= 15 {
			profile = "lowq"
		} else {
			profile = "i1"
		}
	}

	cfg := anyuLeafConfig{
		Profile:       profile,
		SubringLogN:   12,
		MaxLogQP:      1747,
		MaxLevels:     100,
		Q0Bits:        45,
		LogScale:      35,
		StCBits:       33,
		CtSBits:       52,
		EvalModBits:   60,
		NumP:          5,
		Mod1Degree:    127,
		DoubleAngle:   3,
		CtSFirstDepth: 0,
		LogFailure:    -16,
		PrintParams:   true,
		ExactTopBasis: true,
	}

	switch profile {
	case "i1":
		// AnyuWang Table 4, I1-style parameters. For our split pipeline this is
		// appropriate when the leaf ring is logN=16, i.e. top logN=17 with one split.
		cfg.MaxLogQP = 1747
		cfg.NumP = 5
		cfg.StCBits = 33
		cfg.CtSBits = 52
		cfg.EvalModBits = 60
		cfg.Mod1Degree = 127
		cfg.DoubleAngle = 3
	case "lowq":
		// Experimental low-Q adaptation for top logN=16, one split, leaf logN=15.
		// It keeps the subkey target close to AnyuWang's N'=2^12, but compresses
		// the bootstrapping modulus budget to the 900-bit window used by our
		// previous leaf15 experiments.
		cfg.MaxLogQP = 946
		cfg.NumP = 1
		cfg.Q0Bits = 43
		cfg.StCBits = 30
		cfg.CtSBits = 50
		cfg.EvalModBits = 55
		cfg.Mod1Degree = 127
		cfg.DoubleAngle = 3
	default:
		panic(fmt.Errorf("unknown Anyu profile %q; valid options: auto, i1, lowq", profile))
	}

	if cfg.SubringLogN > leafLogN {
		cfg.SubringLogN = leafLogN
	}
	_ = topLogN
	return cfg
}

func anyuSumSlice(s []int) (sum int) {
	for _, ele := range s {
		sum += ele
	}
	return sum
}

func anyuMakeSliceInit(length, val int) []int {
	s := make([]int, length)
	for idx := range s {
		s[idx] = val
	}
	return s
}

func firstOrZero(s []int) int {
	if len(s) == 0 {
		return 0
	}
	return s[0]
}

func anyuCtSManualDepthSplit(logSlots, matrixCount, subringExtDegree, manualFirstDepth int) []int {
	if matrixCount <= 0 || subringExtDegree <= 1 {
		return nil
	}

	freeDepth := int(math.Round(math.Log2(float64(subringExtDegree))))
	if manualFirstDepth > 0 {
		freeDepth = manualFirstDepth
	}
	if freeDepth <= 0 || freeDepth >= logSlots {
		return nil
	}

	out := make([]int, matrixCount)
	out[0] = freeDepth
	remaining := logSlots - freeDepth
	for i := 1; i < matrixCount; i++ {
		partsLeft := matrixCount - i
		depth := int(math.Ceil(float64(remaining) / float64(partsLeft)))
		out[i] = depth
		remaining -= depth
	}
	return out
}

func (param anyuBtpParamSet) getXs() *ring.Ternary {
	xs := new(ring.Ternary)
	xs.ExtDegree = param.bottomLevelSubringSecretExtDegree
	if param.EphemeralSecretWeight > 0 {
		xs.H = param.EphemeralSecretWeight
	} else {
		xs.P = 2.0 / 3
	}
	return xs
}

func (param anyuBtpParamSet) EvalModLevels() int {
	return int(math.Ceil(math.Log2(float64(param.mod1Degree)))) + param.DoubleAngle
}

func (param anyuBtpParamSet) MinimumK(logfailure float64) int {
	logN := 16 - int(math.Round(math.Log2(float64(max(param.bottomLevelSubringSecretExtDegree, 1)))))
	return failure.FindSuitableK(param.getXs(), logN, 15, logfailure)
}

func (param anyuBtpParamSet) MinimumKSimple(logfailure float64) int {
	logN := 16 - int(math.Round(math.Log2(float64(max(param.bottomLevelSubringSecretExtDegree, 1)))))
	return failure.FindSuitableKSimple(param.getXs(), logN, 15, logfailure, 13.1)
}

func (param anyuBtpParamSet) FailureProbabilitySimple() float64 {
	logN := 16 - int(math.Round(math.Log2(float64(max(param.bottomLevelSubringSecretExtDegree, 1)))))
	return failure.ProbabilitySimple(param.getXs(), param.K, logN, 15, 13.1)
}

func (param anyuBtpParamSet) minimumKSimpleForLogN(fullLogN int, logfailure float64) int {
	logN := fullLogN - int(math.Round(math.Log2(float64(max(param.bottomLevelSubringSecretExtDegree, 1)))))
	return failure.FindSuitableKSimple(param.getXs(), logN, 15, logfailure, 13.1)
}

func anyuDefaultI1ParamSet() anyuBtpParamSet {
	return anyuBtpParamSet{
		LogDefaultScale: 35,
		q0:              []int{45},
		qiSlotsToCoeffs: anyuMakeSliceInit(3, 35),
		qEvalMod:        60,
		qiCoeffsToSlots: anyuMakeSliceInit(4, 60),
		piBootstrap:     anyuMakeSliceInit(5, 61),
		StCLevels:       []int{1, 1, 1},
		LogMessageRatio: 8,
		K:               512,
		mod1Degree:      127,
		DoubleAngle:     3,
		mod1Type:        mod1.CosContinuous,

		StCSubringSecretExtDegree:         1 << 3,
		bottomLevelSubringSecretExtDegree: 1 << 4,
		Xs:                                ring.Ternary{P: 2.0 / 3.0},
		isSlim:                            false,
	}
}

func anyuParamSetForLogN(logN int, cfg anyuLeafConfig) anyuBtpParamSet {
	curParam := anyuDefaultI1ParamSet()
	curParam.LogDefaultScale = cfg.LogScale
	curParam.q0 = []int{cfg.Q0Bits}
	curParam.qiSlotsToCoeffs = anyuMakeSliceInit(3, cfg.StCBits)
	curParam.qiCoeffsToSlots = anyuMakeSliceInit(4, cfg.CtSBits)
	curParam.qEvalMod = cfg.EvalModBits
	curParam.piBootstrap = anyuMakeSliceInit(cfg.NumP, 61)
	curParam.mod1Degree = cfg.Mod1Degree
	curParam.DoubleAngle = cfg.DoubleAngle
	curParam.bottomLevelSubringSecretExtDegree = 1 << max(logN-cfg.SubringLogN, 0)
	curParam.StCSubringSecretExtDegree = max(curParam.bottomLevelSubringSecretExtDegree/2, 1)
	curParam.K = curParam.minimumKSimpleForLogN(logN, cfg.LogFailure)
	if curParam.mod1Degree >= 2*(curParam.K-1) {
		curParam.mod1Type = mod1.CosDiscrete
	} else {
		curParam.mod1Type = mod1.CosContinuous
	}
	return curParam
}

func makeAnyuBootstrappingParams(logN int, exactBasis *ckks.Parameters) (ckks.Parameters, bootstrapping.Parameters, error) {
	return makeAnyuBootstrappingParamsWithConfig(logN, exactBasis, defaultAnyuLeafConfig("i1", logN, logN))
}

func makeAnyuBootstrappingParamsWithConfig(logN int, exactBasis *ckks.Parameters, cfg anyuLeafConfig) (ckks.Parameters, bootstrapping.Parameters, error) {
	curParam := anyuParamSetForLogN(logN, cfg)
	stcSubringSecretExtDegree := curParam.StCSubringSecretExtDegree
	bottomLevelSubringSecretExtDegree := curParam.bottomLevelSubringSecretExtDegree
	isSlim := curParam.isSlim
	if !isSlim {
		stcSubringSecretExtDegree = 1
	}

	logDefaultScale := curParam.LogDefaultScale
	q0 := curParam.q0
	qiSlotsToCoeffs := append([]int(nil), curParam.qiSlotsToCoeffs...)
	qiEvalMod := anyuMakeSliceInit(curParam.EvalModLevels(), curParam.qEvalMod)
	qiCoeffsToSlots := append([]int(nil), curParam.qiCoeffsToSlots...)
	piBootstrap := curParam.piBootstrap

	extraCtSLevel := 0
	if bottomLevelSubringSecretExtDegree > 1 {
		qiCoeffsToSlots = qiCoeffsToSlots[:len(qiCoeffsToSlots)-1]
		extraCtSLevel = 1
	}

	nCircuitLevels := (cfg.MaxLogQP - anyuSumSlice(q0) - anyuSumSlice(qiSlotsToCoeffs) - anyuSumSlice(qiEvalMod) -
		anyuSumSlice(qiCoeffsToSlots) - anyuSumSlice(piBootstrap)) / logDefaultScale
	if nCircuitLevels < 0 {
		return ckks.Parameters{}, bootstrapping.Parameters{}, fmt.Errorf("Anyu profile %s fixed bootstrapping levels exceed MaxLogQP: fixed=%d MaxLogQP=%d",
			cfg.Profile,
			anyuSumSlice(q0)+anyuSumSlice(qiSlotsToCoeffs)+anyuSumSlice(qiEvalMod)+anyuSumSlice(qiCoeffsToSlots)+anyuSumSlice(piBootstrap),
			cfg.MaxLogQP)
	}
	qiCircuitSlots := anyuMakeSliceInit(min(nCircuitLevels, cfg.MaxLevels), logDefaultScale)

	logQ := append([]int(nil), q0...)
	if isSlim {
		logQ = append(logQ, qiSlotsToCoeffs...)
		logQ = append(logQ, qiCircuitSlots...)
	} else {
		logQ = append(logQ, qiCircuitSlots...)
		logQ = append(logQ, qiSlotsToCoeffs...)
	}
	logQ = append(logQ, qiEvalMod...)
	logQ = append(logQ, qiCoeffsToSlots...)

	var params ckks.Parameters
	var err error
	if exactBasis != nil {
		params, err = ckks.NewParametersFromLiteral(ckks.ParametersLiteral{
			LogN:            logN,
			Q:               copyQ(*exactBasis),
			P:               copyP(*exactBasis),
			LogDefaultScale: logDefaultScale,
			Xs:              exactBasis.Xs(),
		})
	} else {
		params, err = ckks.NewParametersFromLiteral(ckks.ParametersLiteral{
			LogN:            logN,
			LogQ:            logQ,
			LogP:            piBootstrap,
			LogDefaultScale: logDefaultScale,
			Xs:              curParam.Xs,
		})
	}
	if err != nil {
		return ckks.Parameters{}, bootstrapping.Parameters{}, err
	}

	if cfg.PrintParams {
		fmt.Printf("[AnyuProfile:%s] logN=%d subringLogN=%d extDegree=%d MaxLogQP=%d LogQ=%v LogP=%v circuitLevels=%d K=%d d=%d r=%d ctsFirstDepth=%d\n",
			cfg.Profile, logN, cfg.SubringLogN, bottomLevelSubringSecretExtDegree, cfg.MaxLogQP, logQ, piBootstrap, len(qiCircuitSlots), curParam.K, curParam.mod1Degree, curParam.DoubleAngle,
			firstOrZero(anyuCtSManualDepthSplit(logN-1, 4, bottomLevelSubringSecretExtDegree, cfg.CtSFirstDepth)))
	}

	coeffsToSlotsParameters := dft.MatrixLiteral{
		Type:         dft.HomomorphicEncode,
		Format:       dft.RepackImagAsReal,
		LogSlots:     params.LogMaxSlots(),
		LevelQ:       params.MaxLevelQ(),
		LevelP:       params.MaxLevelP(),
		LogBSGSRatio: 1,
		Levels:       []int{1, 1, 1, 1},

		SubringExtDegree:      bottomLevelSubringSecretExtDegree,
		LevelsInSubringSecret: 1,
		ManualDepthSplit: anyuCtSManualDepthSplit(
			params.LogMaxSlots(),
			4,
			bottomLevelSubringSecretExtDegree,
			cfg.CtSFirstDepth,
		),
	}

	mod1ParametersLiteral := mod1.ParametersLiteral{
		LevelQ:          params.MaxLevel() - coeffsToSlotsParameters.Depth(true) + extraCtSLevel,
		LogScale:        curParam.qEvalMod,
		Mod1Type:        curParam.mod1Type,
		Mod1Degree:      curParam.mod1Degree,
		DoubleAngle:     curParam.DoubleAngle,
		K:               curParam.K,
		LogMessageRatio: curParam.LogMessageRatio,
		Mod1InvDegree:   0,
	}

	slotsToCoeffsParameters := dft.MatrixLiteral{
		Type:         dft.HomomorphicDecode,
		LogSlots:     params.LogMaxSlots(),
		LogBSGSRatio: 1,
		LevelP:       params.MaxLevelP(),
		Levels:       curParam.StCLevels,

		LevelsInSubringSecret: 1,
		SubringExtDegree:      stcSubringSecretExtDegree,
		Base2Decomp:           0,
		ManualDepthSplit:      []int{},
	}

	if isSlim {
		slotsToCoeffsParameters.LevelQ = len(slotsToCoeffsParameters.Levels)
	} else {
		slotsToCoeffsParameters.LevelQ = mod1ParametersLiteral.LevelQ - len(qiEvalMod)
	}

	btpParams := bootstrapping.Parameters{
		ResidualParameters:      params,
		BootstrappingParameters: params,
		SlotsToCoeffsParameters: slotsToCoeffsParameters,
		Mod1ParametersLiteral:   mod1ParametersLiteral,
		CoeffsToSlotsParameters: coeffsToSlotsParameters,
		EphemeralSecretWeight:   curParam.EphemeralSecretWeight,
		SubringSecretExtDegree:  bottomLevelSubringSecretExtDegree,
		CircuitOrder:            bootstrapping.DecodeThenModUp,
	}
	if !isSlim {
		btpParams.CircuitOrder = bootstrapping.ModUpThenEncode
	}

	return params, btpParams, nil
}

type anyuLeafBootstrapper struct {
	*bootstrapping.Evaluator
}

func (a anyuLeafBootstrapper) MaxBootLevel() int {
	return a.BootstrappingParameters.MaxLevel()
}

func makeAnyuLeafEvaluatorPool(topBootParams ckks.Parameters, leafLogN int, skLeaf *rlwe.SecretKey, workers int, cfg anyuLeafConfig) (ckks.Parameters, bootstrapping.Parameters, []leafBootstrapper) {
	if workers <= 0 {
		panic("workers must be positive")
	}

	leafCfg := cfg
	leafCfg.PrintParams = true
	leafParams, leafBtpParams, err := makeAnyuBootstrappingParamsWithConfig(leafLogN, &topBootParams, leafCfg)
	if err != nil {
		panic(fmt.Errorf("makeAnyuBootstrappingParams leaf LogN=%d failed: %w", leafLogN, err))
	}

	paramsBoot := leafBtpParams.BootstrappingParameters
	fmt.Printf("Generating AnyuWang subkey leaf bootstrapping keys for LogN=%d, logQP=%.1f, maxLevel=%d, workers=%d...\n",
		leafLogN, paramsBoot.LogQP(), paramsBoot.MaxLevel(), workers)

	if !reportRingBasisExact("top boot params vs Anyu exact-basis leaf params", topBootParams, paramsBoot) {
		panic("Anyu exact-basis construction failed: leaf params do not match top Q/P")
	}

	btpKeys, _, err := leafBtpParams.GenEvaluationKeys(skLeaf)
	if err != nil {
		panic(fmt.Errorf("Anyu GenEvaluationKeys LogN=%d failed: %w", leafLogN, err))
	}

	evals := make([]leafBootstrapper, workers)
	for i := 0; i < workers; i++ {
		ev, err := bootstrapping.NewEvaluator(leafBtpParams, btpKeys)
		if err != nil {
			panic(fmt.Errorf("Anyu NewEvaluator LogN=%d worker=%d failed: %w", leafLogN, i, err))
		}
		evals[i] = anyuLeafBootstrapper{Evaluator: ev}
	}

	return leafParams, leafBtpParams, evals
}

func printAnyuParamDryRun(logN, layers int, cfg anyuLeafConfig) {
	leafLogN := logN - layers
	if leafLogN < 1 {
		panic("layers too large")
	}

	topCfg := cfg
	topCfg.PrintParams = true
	topParams, _, err := makeAnyuBootstrappingParamsWithConfig(logN, nil, topCfg)
	if err != nil {
		panic(fmt.Errorf("makeAnyuBootstrappingParams top LogN=%d failed: %w", logN, err))
	}
	leafCfg := cfg
	leafCfg.PrintParams = true
	leafParams, _, err := makeAnyuBootstrappingParamsWithConfig(leafLogN, &topParams, leafCfg)
	if err != nil {
		panic(fmt.Errorf("makeAnyuBootstrappingParams leaf LogN=%d failed: %w", leafLogN, err))
	}

	fmt.Println("=== Anyu leaf-engine parameter dry run ===")
	fmt.Printf("top:  logN=%d N=%d qCount=%d pCount=%d logQ=%.1f logQP=%.1f scale=2^%d\n",
		topParams.LogN(), topParams.N(), topParams.QCount(), topParams.PCount(), topParams.LogQ(), topParams.LogQP(), topParams.LogDefaultScale())
	fmt.Printf("leaf: logN=%d N=%d qCount=%d pCount=%d logQ=%.1f logQP=%.1f scale=2^%d exactTopBasis=%v\n",
		leafParams.LogN(), leafParams.N(), leafParams.QCount(), leafParams.PCount(), leafParams.LogQ(), leafParams.LogQP(), leafParams.LogDefaultScale(),
		reportRingBasisExact("Anyu top params vs exact-basis Anyu leaf params", topParams, leafParams))

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

func runAnyu25Capturing(nRepeat int) (float64, time.Duration) {
	return RunAnyu25(nRepeat)
}

// RunAnyu25 executes the embedded AnyuWang 2025 I1-parameter experiment and
// returns average L1 precision in bits and average total bootstrapping time.
func RunAnyu25(nRepeat int) (avgL1Err float64, avgTotal time.Duration) {
	if nRepeat <= 0 {
		nRepeat = 1
	}

	const (
		logN         = 16
		logFailure   = -16.0
		maxLogQP     = 1747
		maxLevels    = 100
		scalestcBits = 33
		scalectsBits = 52
	)

	runCfg := defaultAnyuLeafConfig("i1", logN, logN)
	runCfg.StCBits = scalestcBits
	runCfg.CtSBits = scalectsBits
	runCfg.PrintParams = false

	curParam := anyuParamSetForLogN(logN, runCfg)
	curParam.K = curParam.MinimumKSimple(logFailure)
	if curParam.mod1Degree >= 2*(curParam.K-1) {
		curParam.mod1Type = mod1.CosDiscrete
	} else {
		curParam.mod1Type = mod1.CosContinuous
	}

	fmt.Printf("old K is %d\n", curParam.MinimumK(logFailure))
	fmt.Printf("Failure prob is 2^%f for K = %d\n", curParam.FailureProbabilitySimple(), curParam.K)
	fmt.Println(curParam)

	stcSubringSecretExtDegree := curParam.StCSubringSecretExtDegree
	bottomLevelSubringSecretExtDegree := curParam.bottomLevelSubringSecretExtDegree
	isSlim := curParam.isSlim
	if !isSlim {
		stcSubringSecretExtDegree = 1
	}

	logDefaultScale := curParam.LogDefaultScale
	q0 := curParam.q0
	qiSlotsToCoeffs := append([]int(nil), curParam.qiSlotsToCoeffs...)
	qiEvalMod := anyuMakeSliceInit(curParam.EvalModLevels(), curParam.qEvalMod)
	qiCoeffsToSlots := append([]int(nil), curParam.qiCoeffsToSlots...)
	piBootstrap := curParam.piBootstrap

	extraCtSLevel := 0
	if bottomLevelSubringSecretExtDegree > 1 {
		qiCoeffsToSlots = qiCoeffsToSlots[:len(qiCoeffsToSlots)-1]
		extraCtSLevel = 1
	}

	nCircuitLevels := (maxLogQP - anyuSumSlice(q0) - anyuSumSlice(qiSlotsToCoeffs) - anyuSumSlice(qiEvalMod) -
		anyuSumSlice(qiCoeffsToSlots) - anyuSumSlice(piBootstrap)) / logDefaultScale
	qiCircuitSlots := anyuMakeSliceInit(min(nCircuitLevels, maxLevels), logDefaultScale)
	fmt.Printf("Number of circuit primes = %d\n", len(qiCircuitSlots))

	logQ := append([]int(nil), q0...)
	if isSlim {
		logQ = append(logQ, qiSlotsToCoeffs...)
		logQ = append(logQ, qiCircuitSlots...)
	} else {
		logQ = append(logQ, qiCircuitSlots...)
		logQ = append(logQ, qiSlotsToCoeffs...)
	}
	logQ = append(logQ, qiEvalMod...)
	logQ = append(logQ, qiCoeffsToSlots...)

	fmt.Println("Primes length:", logQ)

	params, err := ckks.NewParametersFromLiteral(ckks.ParametersLiteral{
		LogN:            logN,
		LogQ:            logQ,
		LogP:            piBootstrap,
		LogDefaultScale: logDefaultScale,
		Xs:              curParam.Xs,
	})
	if err != nil {
		panic(err)
	}

	coeffsToSlotsParameters := dft.MatrixLiteral{
		Type:         dft.HomomorphicEncode,
		Format:       dft.RepackImagAsReal,
		LogSlots:     params.LogMaxSlots(),
		LevelQ:       params.MaxLevelQ(),
		LevelP:       params.MaxLevelP(),
		LogBSGSRatio: 1,
		Levels:       []int{1, 1, 1, 1},

		SubringExtDegree:      bottomLevelSubringSecretExtDegree,
		LevelsInSubringSecret: 1,
		ManualDepthSplit: anyuCtSManualDepthSplit(
			params.LogMaxSlots(),
			4,
			bottomLevelSubringSecretExtDegree,
			0,
		),
	}

	mod1ParametersLiteral := mod1.ParametersLiteral{
		LevelQ:          params.MaxLevel() - coeffsToSlotsParameters.Depth(true) + extraCtSLevel,
		LogScale:        curParam.qEvalMod,
		Mod1Type:        curParam.mod1Type,
		Mod1Degree:      curParam.mod1Degree,
		DoubleAngle:     curParam.DoubleAngle,
		K:               curParam.K,
		LogMessageRatio: curParam.LogMessageRatio,
		Mod1InvDegree:   0,
	}

	slotsToCoeffsParameters := dft.MatrixLiteral{
		Type:         dft.HomomorphicDecode,
		LogSlots:     params.LogMaxSlots(),
		LogBSGSRatio: 1,
		LevelP:       params.MaxLevelP(),
		Levels:       curParam.StCLevels,

		LevelsInSubringSecret: 1,
		SubringExtDegree:      stcSubringSecretExtDegree,
		Base2Decomp:           0,
		ManualDepthSplit:      []int{},
	}

	if isSlim {
		slotsToCoeffsParameters.LevelQ = len(slotsToCoeffsParameters.Levels)
	} else {
		slotsToCoeffsParameters.LevelQ = mod1ParametersLiteral.LevelQ - len(qiEvalMod)
	}

	btpParams := bootstrapping.Parameters{
		ResidualParameters:      params,
		BootstrappingParameters: params,
		SlotsToCoeffsParameters: slotsToCoeffsParameters,
		Mod1ParametersLiteral:   mod1ParametersLiteral,
		CoeffsToSlotsParameters: coeffsToSlotsParameters,
		EphemeralSecretWeight:   curParam.EphemeralSecretWeight,
		SubringSecretExtDegree:  bottomLevelSubringSecretExtDegree,
		CircuitOrder:            bootstrapping.DecodeThenModUp,
	}
	if !isSlim {
		btpParams.CircuitOrder = bootstrapping.ModUpThenEncode
	}

	fmt.Printf("Bootstrapping parameters: logN=%d, logSlots=%d, H(%d; %d), ExtDegree=%d, sigma=%v, logQP=%f, levels=%d, scale=2^%d\n",
		btpParams.BootstrappingParameters.LogN(),
		btpParams.BootstrappingParameters.LogMaxSlots(),
		btpParams.BootstrappingParameters.XsHammingWeight(),
		btpParams.EphemeralSecretWeight,
		btpParams.SubringSecretExtDegree,
		btpParams.BootstrappingParameters.Xe(),
		btpParams.BootstrappingParameters.LogQP(),
		btpParams.BootstrappingParameters.QCount(),
		btpParams.BootstrappingParameters.LogDefaultScale())

	var avgL1Linear float64
	var avgStC, avgEvalMod, avgCtS, avgModUp time.Duration

	for repeat := 0; repeat < nRepeat; repeat++ {
		kgen := rlwe.NewKeyGenerator(params)
		sk, pk := kgen.GenKeyPairNew()

		encoder := ckks.NewEncoder(params)
		decryptor := rlwe.NewDecryptor(params, sk)
		encryptor := rlwe.NewEncryptor(params, pk)

		fmt.Println()
		fmt.Println("Generating bootstrapping evaluation keys...")
		evk, _, err := btpParams.GenEvaluationKeys(sk)
		if err != nil {
			panic(err)
		}
		fmt.Println("Done")

		eval, err := bootstrapping.NewEvaluator(btpParams, evk)
		if err != nil {
			panic(err)
		}

		valuesWant := make([]complex128, params.MaxSlots())
		for i := range valuesWant {
			valuesWant[i] = sampling.RandComplex128(-1, 1)
		}

		startLevel := 0
		if isSlim {
			startLevel = slotsToCoeffsParameters.LevelQ
		}
		plaintext := ckks.NewPlaintext(params, startLevel)
		if err := encoder.Encode(valuesWant, plaintext); err != nil {
			panic(err)
		}

		ciphertext, err := encryptor.EncryptNew(plaintext)
		if err != nil {
			panic(err)
		}

		fmt.Println()
		fmt.Println("Precision of values vs. ciphertext")
		valuesTest, _ := printAnyuDebug(params, ciphertext, valuesWant, decryptor, encoder)

		fmt.Println("Bootstrapping...")

		var afterStCTime, afterModUpTime, afterCtSTime, afterEvalTime time.Time
		startTime := time.Now()

		if isSlim {
			if ciphertext, err = eval.SlotsToCoeffs(ciphertext, nil); err != nil {
				panic(err)
			}
			afterStCTime = time.Now()
		}

		if ciphertext, _, err = eval.ScaleDown(ciphertext); err != nil {
			panic(err)
		}

		if ciphertext, err = eval.ModUp(ciphertext); err != nil {
			panic(err)
		}
		afterModUpTime = time.Now()

		var realCt, imagCt *rlwe.Ciphertext
		if realCt, imagCt, err = eval.CoeffsToSlots(ciphertext); err != nil {
			panic(err)
		}
		afterCtSTime = time.Now()

		if realCt, err = eval.EvalMod(realCt); err != nil {
			panic(err)
		}
		if imagCt, err = eval.EvalMod(imagCt); err != nil {
			panic(err)
		}

		if isSlim {
			if err = eval.Evaluator.Mul(imagCt, 1i, imagCt); err != nil {
				panic(err)
			}
			if err = eval.Evaluator.Add(realCt, imagCt, ciphertext); err != nil {
				panic(err)
			}
		}
		afterEvalTime = time.Now()

		if !isSlim {
			if ciphertext, err = eval.SlotsToCoeffs(realCt, imagCt); err != nil {
				panic(err)
			}
			afterStCTime = time.Now()
		}

		var totalTime, stcTime, modUpTime, ctsTime, evalModTime time.Duration
		if isSlim {
			totalTime, stcTime, modUpTime, ctsTime, evalModTime =
				afterEvalTime.Sub(startTime), afterStCTime.Sub(startTime), afterModUpTime.Sub(afterStCTime), afterCtSTime.Sub(afterModUpTime), afterEvalTime.Sub(afterCtSTime)
		} else {
			totalTime, stcTime, modUpTime, ctsTime, evalModTime =
				afterStCTime.Sub(startTime), afterStCTime.Sub(afterEvalTime), afterModUpTime.Sub(startTime), afterCtSTime.Sub(afterModUpTime), afterEvalTime.Sub(afterCtSTime)
		}

		avgCtS += ctsTime
		avgStC += stcTime
		avgModUp += modUpTime
		avgEvalMod += evalModTime
		avgTotal += totalTime

		fmt.Printf("Bootstrapping finished in %s\nStC: %s, ModUp: %s, CtS: %s, EvalMod: %s\n", totalTime, stcTime, modUpTime, ctsTime, evalModTime)
		fmt.Println("Done")

		fmt.Printf("Output ciphertext level = %d, circuit moduli = ", ciphertext.Level())
		if isSlim {
			fmt.Println(btpParams.BootstrappingParameters.LogQi()[slotsToCoeffsParameters.LevelQ+1 : ciphertext.Level()+1])
		} else {
			fmt.Println(btpParams.BootstrappingParameters.LogQi()[1 : ciphertext.Level()+1])
		}

		fmt.Println()
		fmt.Println("Precision of ciphertext vs. Bootstrap(ciphertext)")
		_, logL1Err := printAnyuDebug(params, ciphertext, valuesTest, decryptor, encoder)
		avgL1Linear += math.Pow(2, -logL1Err)
	}

	avgL1Err = -math.Log2(avgL1Linear / float64(nRepeat))
	avgCtS /= time.Duration(nRepeat)
	avgStC /= time.Duration(nRepeat)
	avgModUp /= time.Duration(nRepeat)
	avgEvalMod /= time.Duration(nRepeat)
	avgTotal /= time.Duration(nRepeat)

	fmt.Printf("Avg L1 error is %f\n", avgL1Err)
	fmt.Printf("Avg time... ModUp: %s, CtS: %s, EvalMod: %s, StC: %s, Total: %s\n",
		avgModUp, avgCtS, avgEvalMod, avgStC, avgTotal)

	return avgL1Err, avgTotal
}

func printAnyuDebug(params ckks.Parameters, ciphertext *rlwe.Ciphertext, valuesWant []complex128, decryptor *rlwe.Decryptor, encoder *ckks.Encoder) (valuesTest []complex128, logL1Err float64) {
	slots := ciphertext.Slots()
	if !ciphertext.IsBatched {
		slots *= 2
	}

	valuesTest = make([]complex128, slots)
	if err := encoder.Decode(decryptor.DecryptNew(ciphertext), valuesTest); err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Printf("Level: %d (logQ = %d)\n", ciphertext.Level(), params.LogQLvl(ciphertext.Level()))
	fmt.Printf("Scale: 2^%f\n", math.Log2(ciphertext.Scale.Float64()))
	fmt.Printf("ValuesTest: %10.14f %10.14f %10.14f %10.14f...\n", valuesTest[0], valuesTest[1], valuesTest[2], valuesTest[3])
	fmt.Printf("ValuesWant: %10.14f %10.14f %10.14f %10.14f...\n", valuesWant[0], valuesWant[1], valuesWant[2], valuesWant[3])

	precStats := ckks.GetPrecisionStats(params, encoder, nil, valuesWant, valuesTest, 0, false)
	logL1Err = precStats.Log2L1Err
	fmt.Printf("L1 Err = %f\n", logL1Err)
	fmt.Println(precStats.String())
	fmt.Println()

	return valuesTest, logL1Err
}
