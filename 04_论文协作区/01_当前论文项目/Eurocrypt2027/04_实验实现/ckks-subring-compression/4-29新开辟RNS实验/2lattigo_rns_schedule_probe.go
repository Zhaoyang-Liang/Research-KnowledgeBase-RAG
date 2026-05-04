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

// Lattigo ciphertext-level RNS scheduling probe.
//
// Goal: test whether the following two schedules agree at ciphertext level:
//
//   Path A: ScaleDown_N -> ModUp_N -> RingSwitchDown/Split_{N->n}
//   Path B: ScaleDown_N -> RingSwitchDown/Split_{N->n} -> ModUp_n
//
// This is the first ciphertext-level test needed before trying reduced-Q leaf
// bootstrapping. It isolates the scheduling question and does NOT run C2S,
// EvalMod, or S2C.
//
// Suggested run:
//   go run lattigo_rns_schedule_probe.go -logN=13 -layers=1 -trials=3
//
// Larger split:
//   go run lattigo_rns_schedule_probe.go -logN=14 -layers=2 -trials=3
//
// Interpretation:
//   - leaf A vs leaf B small error: Split and ModUp commute well even with
//     Lattigo RingPacking key switching in the loop.
//   - merged A/B vs top ModUp small error: Merge preserves the same high-level
//     ModUp plaintext after the schedule change.
//   - large fitted gamma but tiny fit error: algebraic equality holds up to an
//     implementation scaling convention.

type diag struct {
    maxErr float64
    rmse   float64
    mean   float64
    gamma  complex128 // got ~= gamma * want
    fitMax float64
    fitRMSE float64
}

func makeResidualParams(logN int) ckks.Parameters {
    params, err := ckks.NewParametersFromLiteral(ckks.ParametersLiteral{
        LogN: logN,
        LogQ: []int{55, 40, 40, 40, 40, 40, 40, 40},
        LogP: []int{61, 61, 61},
        LogDefaultScale: 40,
        Xs: ring.Ternary{H: 192},
    })
    if err != nil {
        panic(err)
    }
    return params
}

func makeBootstrappingParams(logN int) (ckks.Parameters, bootstrapping.Parameters) {
    params := makeResidualParams(logN)
    lit := bootstrapping.ParametersLiteral{
        LogN: utils.Pointy(logN),
        LogP: []int{61, 61, 61, 61},
        Xs: params.Xs(),
    }
    btp, err := bootstrapping.NewParametersFromLiteral(params, lit)
    if err != nil {
        panic(err)
    }
    // Same small-ring convention as the previous prototype.
    if logN < 16 {
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
        LogN: logN,
        Q: copyQ(template),
        P: copyP(template),
        LogDefaultScale: template.LogDefaultScale(),
        Xs: template.Xs(),
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

func makeTopEvaluator(logN int, sk *rlwe.SecretKey) (ckks.Parameters, bootstrapping.Parameters, *bootstrapping.Evaluator) {
    _, btp := makeBootstrappingParams(logN)
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

func makeExactBasisLeafEvaluator(topBootParams ckks.Parameters, leafLogN int, skLeaf *rlwe.SecretKey) (ckks.Parameters, bootstrapping.Parameters, *bootstrapping.Evaluator) {
    exactLeaf := makeParamsFromExactBasis(topBootParams, leafLogN)
    _, leafBtp := makeBootstrappingParams(leafLogN)
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
    fmt.Printf("%-34s LogN=%d Level=%d Degree=%d Scale=%v\n", name, ct.LogN(), ct.Level(), ct.Degree(), ct.Scale)
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
    gamma := estimateGamma(want, got) // got ~= gamma * want
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
    return diag{maxErr: maxErr, rmse: math.Sqrt(sumSq/n), mean: sum/n, gamma: gamma, fitMax: fitMax, fitRMSE: math.Sqrt(fitSumSq/n)}
}

func printDiag(name string, d diag) {
    fmt.Printf("%-48s rawMax=%10.3e rmse=%10.3e gamma=%+.6e%+.6ei fitMax=%10.3e fitRMSE=%10.3e\n",
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

func main() {
    logN := flag.Int("logN", 13, "top log2 ring degree")
    layers := flag.Int("layers", 1, "recursive split layers; k=2^layers")
    trials := flag.Int("trials", 1, "number of randomized plaintext trials")
    flag.Parse()

    if *layers <= 0 || *layers >= *logN {
        panic("invalid layers")
    }
    leafLogN := *logN - *layers
    k := 1 << uint(*layers)

    fmt.Println("=== Lattigo ciphertext-level RNS scheduling probe ===")
    fmt.Printf("top LogN=%d, layers=%d, k=%d, leaf LogN=%d, trials=%d\n", *logN, *layers, k, leafLogN, *trials)
    fmt.Println("Testing: Split(ModUp_N(ScaleDown(ct))) vs ModUp_n(Split(ScaleDown(ct))).")

    _, topBtp := makeBootstrappingParams(*logN)
    paramsTop := topBtp.BootstrappingParameters
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
    _, _, topEval := makeTopEvaluator(*logN, skTop)
    paramsLeaf, _, leafEval := makeExactBasisLeafEvaluator(paramsTop, leafLogN, skLeaf)

    fmt.Printf("top  params: LogN=%d MaxLevel=%d LogQP=%.1f Slots=%d\n", paramsTop.LogN(), paramsTop.MaxLevel(), paramsTop.LogQP(), paramsTop.MaxSlots())
    fmt.Printf("leaf params: LogN=%d MaxLevel=%d LogQP=%.1f Slots=%d exactBasis=%v\n", paramsLeaf.LogN(), paramsLeaf.MaxLevel(), paramsLeaf.LogQP(), paramsLeaf.MaxSlots(), sameRingBasis(paramsTop, paramsLeaf))

    encTop := ckks.NewEncoder(paramsTop)
    encLeaf := ckks.NewEncoder(paramsLeaf)
    encryptorTop := rlwe.NewEncryptor(paramsTop, pkTop)
    decTop := rlwe.NewDecryptor(paramsTop, skTop)
    decLeaf := rlwe.NewDecryptor(paramsLeaf, skLeaf)

    var worstLeafRaw, worstLeafFit, worstMergedRaw, worstMergedFit float64

    for t := 0; t < *trials; t++ {
        values := makeValues(paramsTop.MaxSlots(), t)
        ctLow := encryptAtLevel(paramsTop, encTop, encryptorTop, values, 0)

        scaled, _, err := topEval.ScaleDown(ctLow.CopyNew())
        if err != nil {
            panic(fmt.Errorf("ScaleDown failed: %w", err))
        }
        topBoot, err := topEval.ModUp(scaled.CopyNew())
        if err != nil {
            panic(fmt.Errorf("top ModUp failed: %w", err))
        }

        // Path A: ScaleDown -> top ModUp -> Split.
        leavesA, err := splitTree(rpEval, topBoot.CopyNew(), *layers)
        if err != nil {
            panic(fmt.Errorf("path A split failed: %w", err))
        }

        // Path B: ScaleDown -> Split -> leaf ModUp.
        leavesScaled, err := splitTree(rpEval, scaled.CopyNew(), *layers)
        if err != nil {
            panic(fmt.Errorf("path B split scaled failed: %w", err))
        }
        leavesB, err := modUpLeaf(leafEval, leavesScaled)
        if err != nil {
            panic(fmt.Errorf("path B leaf ModUp failed: %w", err))
        }

        if t == 0 {
            fmt.Println("\nMetadata on trial 0:")
            printCT("input level-0", ctLow)
            printCT("after ScaleDown_N", scaled)
            printCT("after ModUp_N", topBoot)
            for i := range leavesA {
                printCT(fmt.Sprintf("A leaf[%d]=Split(ModUp)", i), leavesA[i])
                printCT(fmt.Sprintf("B leaf[%d]=ModUp(Split)", i), leavesB[i])
            }
        }

        fmt.Printf("\nTrial %d leaf comparisons:\n", t)
        for i := range leavesA {
            valsA := decode(paramsLeaf, encLeaf, decLeaf, leavesA[i])
            valsB := decode(paramsLeaf, encLeaf, decLeaf, leavesB[i])
            d := compare(valsA, valsB)
            printDiag(fmt.Sprintf("leaf[%d] A vs B", i), d)
            if d.maxErr > worstLeafRaw { worstLeafRaw = d.maxErr }
            if d.fitMax > worstLeafFit { worstLeafFit = d.fitMax }
        }

        mergedA, err := mergeTree(rpEval, leavesA)
        if err != nil { panic(fmt.Errorf("merge A failed: %w", err)) }
        mergedB, err := mergeTree(rpEval, leavesB)
        if err != nil { panic(fmt.Errorf("merge B failed: %w", err)) }

        valsTop := decode(paramsTop, encTop, decTop, topBoot)
        valsA := decode(paramsTop, encTop, decTop, mergedA)
        valsB := decode(paramsTop, encTop, decTop, mergedB)
        dA := compare(valsTop, valsA)
        dB := compare(valsTop, valsB)
        dAB := compare(valsA, valsB)
        fmt.Printf("\nTrial %d merged/top comparisons:\n", t)
        printDiag("Merge(Split(ModUp)) vs top ModUp", dA)
        printDiag("Merge(ModUp(Split)) vs top ModUp", dB)
        printDiag("Merged A vs merged B", dAB)
        if dB.maxErr > worstMergedRaw { worstMergedRaw = dB.maxErr }
        if dB.fitMax > worstMergedFit { worstMergedFit = dB.fitMax }
    }

    fmt.Println("\n=== Summary ===")
    fmt.Printf("worst leaf raw max error          = %.6e\n", worstLeafRaw)
    fmt.Printf("worst leaf fitted-gamma max error = %.6e\n", worstLeafFit)
    fmt.Printf("worst merged raw max error        = %.6e\n", worstMergedRaw)
    fmt.Printf("worst merged fitted max error     = %.6e\n", worstMergedFit)
    fmt.Println("\nInterpretation:")
    fmt.Println("  This tests the ciphertext-level schedule change needed before reduced-Q experiments.")
    fmt.Println("  If errors are small, ModUp can plausibly be moved after RingSwitchDown/Split.")
    fmt.Println("  This still does not prove reduced-Q bootstrapping correctness; it only validates the scheduling primitive with the same Q/P basis.")
}
