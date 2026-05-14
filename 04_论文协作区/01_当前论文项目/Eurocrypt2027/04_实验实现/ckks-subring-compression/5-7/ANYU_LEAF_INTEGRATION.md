# AnyuWang leaf-engine integration notes

目标是把当前 no-\(B_k\) split pipeline 的 leaf bootstrapping engine

\[
\mathsf{ScaleDown}_n \to \mathsf{ModUp}_n \to \mathsf{C2S}_n
\to \mathsf{EvalMod}_n \to \mathsf{S2C}_n
\]

替换为 AnyuWang 2025 的 subring-secret optimized bootstrapping engine，从而叠加两层加速：

1. 外层：本实验的 \(N \to n=N/k\) coefficient split。
2. 内层：每个 leaf ring 上使用 AnyuWang 的 subring-secret bootstrapping 优化。

## Current state

`main.go` 现在已经把 leaf engine 抽象成 `leafBootstrapper` 接口：

```go
type leafBootstrapper interface {
    ScaleDown(*rlwe.Ciphertext) (*rlwe.Ciphertext, *rlwe.Scale, error)
    ModUp(*rlwe.Ciphertext) (*rlwe.Ciphertext, error)
    CoeffsToSlots(*rlwe.Ciphertext) (*rlwe.Ciphertext, *rlwe.Ciphertext, error)
    EvalMod(*rlwe.Ciphertext) (*rlwe.Ciphertext, error)
    SlotsToCoeffs(*rlwe.Ciphertext, *rlwe.Ciphertext) (*rlwe.Ciphertext, error)
    MaxBootLevel() int
}
```

The current split implementation uses `lattigoLeafBootstrapper`, a thin wrapper
around standard `*bootstrapping.Evaluator`.

`anyu25.go` now embeds the AnyuWang I1 parameter construction directly in this
package. It no longer calls the external `subkey25/main` binary.

There are now two Anyu paths:

1. `go run ./5-7/. -run=subKey25 -reps=1`
   runs the embedded standalone AnyuWang benchmark.
2. `go run ./5-7/. -run=ours -leafEngine=anyu -logN=17 -layers=1 -splitFirstLeafPreprocess -reps=1`
   uses AnyuWang's subkey bootstrapping evaluator as the actual leaf engine in
   the split pipeline.

## Why direct replacement is not yet one-line

The AnyuWang experiment under

`../../下一步实验/别人的实验对比main-25AnyuWang/subkey25/main.go`

is a standalone program. It constructs its own parameters, keys, plaintext,
ciphertext, bootstrapping evaluator, and timing loop. It does not expose a Go
library function that accepts our leaf ciphertexts.

It also depends on the local Lattigo fork:

```go
replace github.com/tuneinsight/lattigo/v6 => /Users/mac/Desktop/做实验/subring_bts-main
```

That fork adds subring-secret fields such as:

```go
dft.MatrixLiteral.SubringExtDegree
dft.MatrixLiteral.LevelsInSubringSecret
bootstrapping.Parameters.SubringSecretExtDegree
```

These fields are not present in the published Lattigo v6.2.0 dependency used by
the current `ckks-subring-compression/go.mod`.

## Required integration path

1. The experiment module now uses the local AnyuWang fork through:

```go
replace github.com/tuneinsight/lattigo/v6 => /Users/mac/Desktop/做实验/subring_bts-main
```

2. The standalone Anyu parameter construction has been copied into `anyu25.go`
   and refactored into reusable helper functions, parameterized by the leaf ring
   dimension:

```go
func makeAnyuLeafBootstrappingParams(
    leafLogN int,
    topOrLeafBasis ckks.Parameters,
    cfg AnyuLeafConfig,
) (ckks.Parameters, bootstrapping.Parameters, error)
```

3. `makeAnyuLeafEvaluatorPool` builds a leaf evaluator pool returning
   `[]leafBootstrapper`, analogous to `makeLeafEvaluatorPoolNoBK`, but using
   AnyuWang's subring-secret `bootstrapping.Parameters`.

4. Verify exact Q/P basis compatibility with the split/merge path. The current
   no-\(B_k\) implementation expects leaf ciphertexts to live in a basis
   compatible with the top ring. If Anyu parameters use a different modulus
   layout, then split/merge cannot be composed directly without an explicit
   modulus/layout conversion.

5. Run three correctness tests before timing:

```text
standard leaf engine:
Split -> leaf ScaleDown/ModUp/C2S -> EvalMod -> S2C -> Merge

Anyu leaf engine, same input values:
Split -> Anyu leaf ScaleDown/ModUp/C2S -> EvalMod -> S2C -> Merge

embedded Anyu standalone:
go run ./5-7/. -run=subKey25 -reps=1
```

The second row is now implemented and is the real integrated path. The third row
is only a benchmark reference.

## First integrated result

Command:

```bash
env GOCACHE=/private/tmp/go-build go run ./5-7/. \
  -run=ours \
  -leafEngine=anyu \
  -logN=17 \
  -layers=1 \
  -splitFirstLeafPreprocess \
  -reps=1
```

Observed output:

```text
top:  LogN=17, leaf: LogN=16, k=2
logQP ≈ 1730
leaf 128-bit reference bound: 1747
Output L1 precision: 19.38 bits
Total time: 51.72s
```

This confirms that split leaf ciphertexts can be consumed by the embedded Anyu
leaf evaluator when the top and leaf Q/P bases are synchronized.

## Expected risk

The main risk is modulus-chain mismatch. AnyuWang's optimized bootstrapping
chooses a specialized `LogQ` layout:

\[
q_0 \parallel Q_{\mathrm{circuit}} \parallel Q_{\mathrm{StC}}
\parallel Q_{\mathrm{EvalMod}} \parallel Q_{\mathrm{CtS}},
\]

whereas the current split experiment often forces the leaf ring to use the exact
top-ring Q/P basis. If those bases differ, the optimization must either:

1. choose Anyu parameters using the same top-compatible Q/P basis, or
2. add an explicit modulus-switching/basis-conversion stage between split and
   leaf bootstrapping.

The first option is the cleaner experiment to try first.
