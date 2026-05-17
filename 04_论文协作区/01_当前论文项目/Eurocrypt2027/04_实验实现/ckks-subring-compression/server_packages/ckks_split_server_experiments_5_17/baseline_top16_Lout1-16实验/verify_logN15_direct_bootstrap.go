//go:build directverify

package main

import (
	"flag"
	"fmt"
	"math"
	"time"

	"github.com/tuneinsight/lattigo/v6/circuits/ckks/bootstrapping"
	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/ring"
	"github.com/tuneinsight/lattigo/v6/schemes/ckks"
	"github.com/tuneinsight/lattigo/v6/utils"
)

func main() {
	logN := flag.Int("logN", 15, "ring degree log2(N)")
	q0Bits := flag.Int("q0Bits", 40, "first residual Q prime bit-size")
	circuitLevels := flag.Int("circuitLevels", 4, "number of post-bootstrap circuit primes")
	circuitPrimeBits := flag.Int("circuitPrimeBits", 31, "post-bootstrap circuit prime bit-size")
	scaleBits := flag.Int("scaleBits", 31, "default CKKS scale bit-size")
	flag.Parse()

	logQ := make([]int, 1+*circuitLevels)
	logQ[0] = *q0Bits
	for i := 1; i < len(logQ); i++ {
		logQ[i] = *circuitPrimeBits
	}

	params, err := ckks.NewParametersFromLiteral(ckks.ParametersLiteral{
		LogN:            *logN,
		LogQ:            logQ,
		LogP:            []int{56, 56},
		Xs:              ring.Ternary{H: 16384},
		LogDefaultScale: *scaleBits,
	})
	if err != nil {
		panic(err)
	}

	btpLit := bootstrapping.ParametersLiteral{
		LogN: utils.Pointy(*logN),
		LogP: []int{56, 56},
		Xs:   params.Xs(),
		SlotsToCoeffsFactorizationDepthAndLogScales: [][]int{
			{30, 30},
		},
		CoeffsToSlotsFactorizationDepthAndLogScales: [][]int{
			{52},
			{52},
		},
		EvalModLogScale: utils.Pointy(55),
	}

	btpParams, err := bootstrapping.NewParametersFromLiteral(params, btpLit)
	if err != nil {
		panic(err)
	}

	// Same correction used by Lattigo's small-LogN bootstrapping examples.
	btpParams.Mod1ParametersLiteral.LogMessageRatio += 16 - params.LogN()

	fmt.Printf("Residual: logN=%d slots=%d logQ=%.1f logQP=%.1f levels=%d scale=2^%d H=%d\n",
		params.LogN(), params.MaxSlots(), params.LogQ(), params.LogQP(), params.QCount(), params.LogDefaultScale(), params.XsHammingWeight())
	fmt.Printf("Bootstrap: logN=%d slots=%d logQ=%.1f logQP=%.1f levels=%d scale=2^%d H=%d ephH=%d\n",
		btpParams.BootstrappingParameters.LogN(),
		btpParams.BootstrappingParameters.MaxSlots(),
		btpParams.BootstrappingParameters.LogQ(),
		btpParams.BootstrappingParameters.LogQP(),
		btpParams.BootstrappingParameters.QCount(),
		btpParams.BootstrappingParameters.LogDefaultScale(),
		btpParams.BootstrappingParameters.XsHammingWeight(),
		btpParams.EphemeralSecretWeight)

	kgen := rlwe.NewKeyGenerator(params)
	sk, pk := kgen.GenKeyPairNew()
	encoder := ckks.NewEncoder(params)
	encryptor := rlwe.NewEncryptor(params, pk)
	decryptor := rlwe.NewDecryptor(params, sk)

	fmt.Println("Generating bootstrapping keys...")
	t0 := time.Now()
	evk, _, err := btpParams.GenEvaluationKeys(sk)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Keygen done in %.3fs\n", time.Since(t0).Seconds())

	eval, err := bootstrapping.NewEvaluator(btpParams, evk)
	if err != nil {
		panic(err)
	}

	valuesWant := make([]complex128, params.MaxSlots())
	for i := range valuesWant {
		re := -1.0 + 2.0*float64(i%257)/256.0
		im := -0.25 + 0.5*float64(i%131)/130.0
		valuesWant[i] = complex(re, im)
	}

	pt := ckks.NewPlaintext(params, 0)
	if err := encoder.Encode(valuesWant, pt); err != nil {
		panic(err)
	}
	ct, err := encryptor.EncryptNew(pt)
	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println("Precision of values vs. input ciphertext")
	valuesInput := printDebug(params, ct, valuesWant, decryptor, encoder)

	fmt.Println("Bootstrapping...")
	t1 := time.Now()
	out, err := eval.Bootstrap(ct)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Bootstrap done in %.3fs\n", time.Since(t1).Seconds())
	fmt.Printf("Output level=%d activeLogQ=%d scale=2^%.6f\n",
		out.Level(), params.LogQLvl(out.Level()), math.Log2(out.Scale.Float64()))

	fmt.Println()
	fmt.Println("Precision of input ciphertext vs. bootstrapped ciphertext")
	_ = printDebug(params, out, valuesInput, decryptor, encoder)
}

func printDebug(params ckks.Parameters, ct *rlwe.Ciphertext, valuesWant []complex128, decryptor *rlwe.Decryptor, encoder *ckks.Encoder) []complex128 {
	valuesTest := make([]complex128, ct.Slots())
	if err := encoder.Decode(decryptor.DecryptNew(ct), valuesTest); err != nil {
		panic(err)
	}

	fmt.Printf("Level: %d activeLogQ=%d scale=2^%.6f\n", ct.Level(), params.LogQLvl(ct.Level()), math.Log2(ct.Scale.Float64()))
	fmt.Printf("ValuesTest: %.10f %.10f %.10f %.10f...\n", valuesTest[0], valuesTest[1], valuesTest[2], valuesTest[3])
	fmt.Printf("ValuesWant: %.10f %.10f %.10f %.10f...\n", valuesWant[0], valuesWant[1], valuesWant[2], valuesWant[3])
	fmt.Println(ckks.GetPrecisionStats(params, encoder, nil, valuesWant, valuesTest, 0, false).String())
	return valuesTest
}
