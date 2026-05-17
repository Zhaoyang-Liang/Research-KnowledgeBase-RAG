# Low-output-level parameter selection

This note records the current parameter-selection logic for the `5-16`
low-output-level CKKS split bootstrapping experiment.

## 1. Decide what security modulus means

There are two different bit counts in the experiment:

```text
params.LogQ()  = total Q-chain allocated to the bootstrapping evaluator
params.LogQP() = params.LogQ() plus special P-primes used by key switching
out.LogQLvl()  = active Q still carried by the final output ciphertext
```

For output-level discussion, the relevant numbers are:

```text
out.Level()
params.LogQLvl(out.Level())
```

In the current low-output experiments, the final output ciphertext has:

```text
out.Level() = 1
active output logQ ~= 50--54 bits
```

The much larger `params.LogQ() ~= 870--936` is the total bootstrapping working
chain, not the final active output modulus.

For conservative RLWE/LWE security estimation, the safer convention is to
estimate with `logQP`, because evaluation keys and key switching use the
extended basis. Under this convention, `leaf logN=15` reaches about 128 bits
only around:

```text
n = 32767, logQP ~= 868
```

Using only `logQ` is a weaker accounting convention. It is useful as a
diagnostic, but the paper should explicitly state which convention is used.

## 2. The strict `logQP ~= 868` point does not run in native Lattigo

The closest strict candidate is:

```text
logN=16, leafLogN=15
q0Bits=26, circuitPrimeBits=18, defaultScale=25
numP=0
```

Dry run:

```text
top/leaf logQ  ~= 866.6
top/leaf logQP ~= 866.6
```

This is inside the estimated 128-bit window, but native Lattigo crashes during
bootstrapping-key generation with `numP=0`. In this implementation, the
bootstrapping key-generation path expects a P-basis.

With `numP=1`, the same output-Q point becomes:

```text
top/leaf logQ  ~= 866.6
top/leaf logQP ~= 927.6
```

and it is no longer 128-bit under the conservative `QP` convention.

## 3. ScaleDown correctness imposes another window

Even before security, CKKS bootstrapping needs enough room between the initial
modulus and the scale. A useful empirical rule in this code path is:

```text
q0Bits - defaultScale >= 8.
```

For example:

```text
q0Bits=26, defaultScale=25
```

fails with:

```text
ScaleDown failed: initial Q/Scale = 2.003906
```

The closest runnable point to total-chain `params.LogQ() ~= 868` that satisfies
this ScaleDown condition is:

```text
q0Bits=30, circuitPrimeBits=18, defaultScale=22, numP=1
top/leaf logQ  ~= 870.6
top/leaf logQP ~= 931.6
```

It runs, but split precision collapses:

| point | method | output level | active output logQ | precision | total time |
|---|---|---:|---:|---:|---:|
| total `params.LogQ() ~= 870.6` | direct top BTS | 1 | 50 | 9.08 bits | 51.64 s |
| total `params.LogQ() ~= 870.6` | split BTS | 1 | 50 | 0.97 bits | 25.67 s |

This is not a valid correctness point for the split method.

Estimator values for the same point:

| accounting | estimator input | minimum estimated bits |
|---|---:|---:|
| Q only | 870.585 | 127.51 |
| actual QP | 931.585 | 118.92 |

## 4. First working near-boundary point

The next nearby point that works is:

```text
q0Bits=34, circuitPrimeBits=18, defaultScale=25, numP=1
top/leaf logQ  ~= 874.6
top/leaf logQP ~= 935.6
```

Result:

| point | method | output level | output logQ at level | precision | total time |
|---|---|---:|---:|---:|---:|
| `logQ ~= 874.6` | direct top BTS | 1 | 54 | 12.08 bits | 53.64 s |
| `logQ ~= 874.6` | split BTS | 1 | 54 | 11.51 bits | 27.43 s |

Speedup:

```text
53.64 / 27.43 ~= 1.96x
```

This is the closest currently working total-chain `params.LogQ() ~= 868`-level
experiment. It is not 128-bit secure under the conservative `QP` convention,
but it answers the low-output-level timing question: both methods output a
level-1 ciphertext with active output `logQ ~= 54`.

Estimator values:

| accounting | estimator input | minimum estimated bits |
|---|---:|---:|
| Q only | 874.585 | 126.94 |
| actual QP | 935.585 | 118.37 |

## 5. Practical selection flow

Use the following flow for future experiments.

1. Choose the target output precision.

   In these L1 precision measurements, a rough empirical rule is:

   ```text
   output precision ~= defaultScale - 13 bits.
   ```

   So 20-bit output needs `defaultScale` around 33--35, while 12-bit output
   can use `defaultScale` around 25.

2. Choose `q0Bits` so ScaleDown has room.

   Use:

   ```text
   q0Bits >= defaultScale + 8.
   ```

3. Run a dry scan to see actual Lattigo `logQ/logQP`.

   Do not infer total bootstrapping `logQ` by hand. Lattigo adds internal
   primes for CtS/EvalMod/StC.

4. Check security at the leaf dimension.

   For one split layer:

   ```text
   leafLogN = topLogN - 1.
   n = 2^leafLogN - 1.
   ```

   Estimate with `logQP` for conservative claims; optionally record `logQ`
   as a diagnostic.

5. Run the direct top lowQ and split lowQ experiments with the same output
   `Q` and the same output level.

   The fair comparison is:

   ```text
   Ct_{q0,N} -> Ct_{Q_s,N}
   ```

   for both direct and split paths.

6. Accept the point only if all of the following hold:

   ```text
   leaf security >= 128 bits under the chosen convention
   direct precision is acceptable
   split precision is close to direct precision
   output levels match
   split time < direct time
   ```

## 6. Current interpretation

There are now three distinct categories:

| category | example | interpretation |
|---|---|---|
| strict secure leaf `QP` | `leaf logQP <= 868` at `logN=15` | Native Lattigo cannot run this with `P=0`; with `P=1`, `QP` exceeds the bound. |
| low output `Q` diagnostic | `logQ ~= 874.6`, `logQP ~= 935.6` | Runs and shows about `1.96x` speedup, but not 128-bit under conservative `QP`. |
| secure p35 high-precision | `top logN=17`, `leaf logN=16`, `logQP ~= 962` | Conservative security passes, but ring dimension is larger. |

This is exactly the feasible-window tension:

```text
correctness/precision lower bound <= chosen Q budget <= security upper bound.
```

For `leaf logN=15`, the native Lattigo implementation makes the conservative
`QP` window too tight. For a paper-quality 128-bit experiment without changing
Lattigo's key-switching internals, the clean route is to raise the top ring to
`logN=17` so that split leaves are at `logN=16`.
