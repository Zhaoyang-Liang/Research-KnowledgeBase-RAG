# Security estimates for 5-16 lowQ experiments

Estimator command source:

```text
/Users/mac/opt/lattice-estimator/fhe_estimate.py
```

I added command-line arguments to that script so the same file can estimate different \(n,\log QP\) pairs.

Run all checks:

```bash
./estimate_security_5_16.sh
```

## Current p35 experiment is not 128-bit after split

The current high-precision lowQ experiment uses:

```text
top logN=16, split leaf logN=15
top/leaf logQ  ~= 901 bits
top/leaf logQP ~= 962 bits
```

The security bottleneck is the leaf, because after split each bootstrapping is performed in the smaller ring.

Using `n=32767` for the `logN=15` leaf:

| case | n | logQP used by estimator | minimum estimated bits |
|---|---:|---:|---:|
| p35 leaf, Q only diagnostic | 32767 | 901 | 123.07 |
| p20 leaf, actual QP | 32767 | 936 | 118.36 |
| p35 leaf, actual QP | 32767 | 962 | 115.05 |

Thus the current `run_lowq_lattigo_p35.sh` is useful as a correctness and speed diagnostic, but it should not be reported as a 128-bit secure split experiment if the estimator is run on the full \(QP\).

## Leaf threshold at logN=15

For the same leaf dimension:

| n | logQP | minimum estimated bits |
|---:|---:|---:|
| 32767 | 880 | 126.11 |
| 32767 | 868 | 128.03 |
| 32767 | 850 | 130.82 |

So a `logN=15` leaf needs roughly

\[
\log(QP) \le 868
\]

under this estimator setup. The p35 experiment needs about \(962\) bits, which is about \(94\) bits outside this window.

## 128-bit adjustment

The direct adjustment is to raise the top ring by one level:

```text
top logN=17, split leaf logN=16
q0Bits=45, circuitLevels=1, circuitPrimeBits=35
numP=1, defaultScale=35
```

Then the leaf estimate uses `n=65535` and the same `logQP ~= 962`:

| case | n | logQP | minimum estimated bits |
|---|---:|---:|---:|
| p35 leaf after top logN=17 adjustment | 65535 | 962 | 246.05 |

This is comfortably above 128-bit. The corresponding experiment script is:

```bash
./run_lowq_lattigo_p35_sec128.sh
```

This changes the experiment from a fast diagnostic at `top logN=16` into a leaf-secure split experiment at `top logN=17`.

