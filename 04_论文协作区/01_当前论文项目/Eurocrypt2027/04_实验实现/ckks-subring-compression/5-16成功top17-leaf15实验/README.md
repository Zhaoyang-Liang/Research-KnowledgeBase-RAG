# 5-16 成功 top17-leaf15 实验

本目录保存一个已经跑通的 **Lattigo-only** split bootstrapping 实验点：

```text
top logN = 17
leaf logN = 15
split layers = 2
leaves = 4
leaf engine = Lattigo 自带 CKKS bootstrapping
slot recovery = 不使用 B_k / B_k^{-1}
security accounting = PQ / logQP
baseline = 不跑
```

本实验不使用 AnyuWang/subKey bootstrapping。`anyu25.go` 只是保留原程序编译依赖，正式命令固定使用 `-leafEngine=lattigo`。

## 目录内容

```text
main.go
anyu25.go
verify_logN15_direct_bootstrap.go
go.mod
go.sum
logs/
results/
参数选取办法.md
```

本目录单独带 `go.mod`，因为目录名含中文，不能作为上级 Go module 的普通 package path 直接 `go run .`。

## 成功参数

本实验使用 top16->leaf15 成功点中的同一套 leaf15 低输出参数：

```text
lattigoPreset = n15qp849
```

参数结构：

```text
Residual Q = [40, 31, 31, 31]
P          = [56, 56]
StC        = [30, 30]
EvalMod    = 8 * 55
CtS        = [52, 52]
scale      = 2^31
```

因此：

```text
bootstrapping logQ  ≈ 737
bootstrapping logQP ≈ 849
output level = 3
output circuit moduli = [31, 31, 31]
```

注意：使用 `-lattigoPreset=n15qp849` 时，命令行里默认打印的 `q0Bits/circuitLevels/circuitPrimeBits/numP/defaultScale` 不再决定实际 bootstrapping modulus；实际以 preset 中的 `LogQ/LogP` 为准。

## 安全性

必须按 leaf15 的 PQ/logQP 口径估算，不看 Q-only。因为 top17 拆到 leaf15 后，leaf15 是安全瓶颈。

Estimator 命令：

```bash
PATH="/opt/anaconda3/envs/sage/bin:$PATH" \
/opt/anaconda3/envs/sage/bin/python /Users/mac/opt/lattice-estimator/fhe_estimate.py \
  --n 32767 \
  --logQP 849 \
  --sigma 3.2 \
  --tag leaf15_top17_split_lattigo_n15qp849 \
  --attacks primal_usvp,primal_bdd,dual_hybrid
```

已保存输出：

```text
logs/security_leaf15_QP849.log
```

关键结果：

```text
primal_usvp  ≈ 131.08 bits
primal_bdd   ≈ 130.90 bits
dual_hybrid  ≈ 131.88 bits
minimum estimated bits = 130.89845721434037
```

## Dry Run

```bash
go run . \
  -dryRunParams=true \
  -logN=17 \
  -layers=2 \
  -run=ours \
  -leafEngine=lattigo \
  -lattigoPreset=n15qp849 \
  -splitFirstLeafPreprocess=true
```

应看到：

```text
top:  logN=17 N=131072 qCount=15 pCount=2 logQ=737.0 logQP=849.0
leaf: logN=15 N=32768  qCount=15 pCount=2 logQ=737.0 logQP=849.0
leaf logQP <= 868, margin=19.0 bits [PASS]
```

## 正式运行命令

只跑 split route，不跑 baseline：

```bash
go run . \
  -logN=17 \
  -layers=2 \
  -run=ours \
  -leafEngine=lattigo \
  -lattigoPreset=n15qp849 \
  -splitFirstLeafPreprocess=true \
  -parallelLeaves=true \
  -leafWorkers=4 \
  -reps=1 \
  -outputTSV=results/split_top17_leaf15_lattigo_n15qp849.tsv \
  2>&1 | tee logs/split_top17_leaf15_lattigo_n15qp849.log
```

## 已得实验结果

日志：

```text
logs/split_top17_leaf15_lattigo_n15qp849.log
```

TSV：

```text
results/split_top17_leaf15_lattigo_n15qp849.tsv
```

关键结果：

```text
top logN = 17
leaf logN = 15
layers = 2
leaves = 4
top logQP = 849.002817
leaf logQP = 849.002817
output level = 3
output circuit moduli = [31 31 31]
```

时间：

```text
Total         = 12.318991750 s
Split         = 0.057511208 s
leafScaleDown = 0.000347583 s
leafModUp     = 0.156895875 s
leafC2S       = 9.691315875 s
leafEvalMod   = 1.073625292 s
leafS2C       = 1.079163333 s
Merge         = 0.260132584 s
```

精度：

```text
input ciphertext vs plaintext L1 precision        ≈ 17.08 bits
split BTS output vs input ciphertext L1 precision ≈ 16.19 bits
```

## 结论

top17->leaf15 的 Lattigo-only split bootstrapping 已经跑通。由于参数完全按 leaf15 的 `n15qp849` 选择，输出 level 与 top16->leaf15 一样是 3；同时由于 leaf 数从 2 增加到 4，总时间从 top16->leaf15 的约 5.12s 增加到约 12.32s，但仍显著小于 top17->leaf16 的约 51.65s。
