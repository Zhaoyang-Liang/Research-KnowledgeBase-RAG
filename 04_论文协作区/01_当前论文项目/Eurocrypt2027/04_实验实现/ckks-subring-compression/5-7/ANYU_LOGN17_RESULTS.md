# Anyu leaf engine: logN=17 -> leaf logN=16 notes

Date: 2026-05-08

This note records the first parameter sweep for the factored no-\(B_k\) path
using AnyuWang-style subring-secret bootstrapping on the leaves.

## Setting

- Top ring: `logN=17`, \(N=131072\)
- Split layers: `1`, so there are two leaves
- Leaf ring: `logN=16`, \(n=65536\)
- Anyu subring secret target: `anyuSubringLogN=12`
- Top Anyu extension degree: \(2^{17}/2^{12}=32\)
- Leaf Anyu extension degree: \(2^{16}/2^{12}=16\)
- Auto CtS first depth:
  - top: `ctsFirstDepth=5`
  - leaf: `ctsFirstDepth=4`

This is closer to AnyuWang's original CKKS setting than the previous
`logN=16 -> leaf logN=15` experiment, because the leaf bootstrapper is now
exactly in the \(2^{16}\to 2^{12}\) downsampling regime.

## Security Budget Reading

Using the working 128-bit reference table:

| Ring | Reference max log(PQ) |
|---|---:|
| `logN=15`, \(N=32768\) | about 868 |
| `logN=16`, \(N=65536\) | about 1747 |
| `logN=17`, \(N=131072\) | about 3523 |

For the split method the relevant bottleneck is the leaf ciphertext after
splitting. Thus for `logN=17, layers=1`, the leaf ring is `logN=16`, and
`logQP` should be compared against about 1747.

## Results

| Label | Command overrides | logQ | logQP | output levels | Avg L1 precision | Total time | Notes |
|---|---|---:|---:|---:|---:|---:|---|
| lowq tuned | `-anyuProfile=lowq -anyuEvalModBits=56 -anyuStCBits=26 -anyuCtSBits=50` | 866 | 927 | 1 | 15.47 bits | 45.66s | Safe but too little modulus budget for good precision |
| lowq default | `-anyuProfile=lowq` | 868 | 929 | 1 | 14.64 bits | 42.97s | Worse than tuned lowq |
| i1 default | `-anyuProfile=i1` | 1425 | 1730 | 15 | 19.40 bits | 52.03s | Best balanced setting so far |
| i1 CtS55 | `-anyuProfile=i1 -anyuCtSBits=55` | 1434 | 1739 | 15 | 19.42 bits | 50.51s | Only marginally better; logQP margin is only about 8 bits |
| i1 d255 | `-anyuProfile=i1 -anyuMod1Degree=255` | 1415 | 1720 | 13 | 19.39 bits | 45.33s | Higher EvalMod degree did not improve precision |
| direct baseline | `-run=baseline -anyuProfile=i1` | 1425 | 1730 | 15 | 18.61 bits | 135.08s | Single-ring top bootstrap baseline |

## Current Takeaway

The `logN=17 -> leaf logN=16` setting is much more convincing than the low-Q
`logN=16 -> leaf logN=15` setting.

The current best practical command is:

```bash
go run ./5-7/. \
  -run=ours \
  -leafEngine=anyu \
  -anyuProfile=i1 \
  -anyuSubringLogN=12 \
  -logN=17 \
  -layers=1 \
  -splitFirstLeafPreprocess \
  -reps=1
```

It gives about `19.40` bits of average L1 precision, keeps 15 output circuit
levels, and is about \(135.08/52.03 \approx 2.6\times\) faster than the direct
top-ring bootstrap in this local run.

The `-anyuCtSBits=55` variant gives only a tiny precision gain and consumes
almost all remaining `logQP` security margin, so the default `i1` profile is the
cleaner candidate unless later repeated trials show otherwise.

