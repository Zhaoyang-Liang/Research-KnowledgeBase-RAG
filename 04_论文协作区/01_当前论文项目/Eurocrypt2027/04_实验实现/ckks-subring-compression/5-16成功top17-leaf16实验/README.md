# 5-16 成功 top17-leaf16 实验

本目录保存一个已经跑通的 **Lattigo-only** split bootstrapping 实验点：

```text
top logN = 17
leaf logN = 16
split layers = 1
leaves = 2
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

这次 top17->leaf16 不需要新增专门的 Lattigo preset，直接使用当前程序的 `lattigoPreset=default`，但把 residual/output budget 设为：

```text
q0Bits = 45
circuitLevels = 15
circuitPrimeBits = 35
numP = 5
leafNumP = 5
defaultScale = 35
```

参数结构可以粗略理解为：

```text
Residual/output Q = 45 + 15 * 35 = 570 bits
Bootstrapping internal Q = 821 bits
P = 5 * 61 = 305 bits
-----------------------------------------------
total logQP = 1696 bits
```

实际 dry-run 打印：

```text
top:  logN=17 N=131072 qCount=31 pCount=5 logQ=1391.0 logQP=1696.0
leaf: logN=16 N=65536  qCount=31 pCount=5 logQ=1391.0 logQP=1696.0
leaf logQP <= 1747, margin=51.0 bits [PASS]
```

输出 level：

```text
output level = 15
output circuit moduli = [35 35 35 35 35 35 35 35 35 35 35 35 35 35 35]
```

## 安全性

必须按 PQ 口径估算，不看 Q-only。

Estimator 命令：

```bash
PATH="/opt/anaconda3/envs/sage/bin:$PATH" \
/opt/anaconda3/envs/sage/bin/python /Users/mac/opt/lattice-estimator/fhe_estimate.py \
  --n 65535 \
  --logQP 1696 \
  --sigma 3.2 \
  --tag leaf16_top17_split_lattigo_Lout15_QP1696 \
  --attacks primal_usvp,primal_bdd,dual_hybrid
```

已保存输出：

```text
logs/security_leaf16_QP1696.log
```

关键结果：

```text
primal_usvp  ≈ 132.09 bits
primal_bdd   ≈ 132.09 bits
dual_hybrid  ≈ 133.11 bits
minimum estimated bits = 132.09138343060562
```

备注：本地 `/Users/mac/opt/lattice-estimator/fhe_estimate.py` 已修复大 `logQP` 溢出问题，把 `q=2**logQP` 改为 Sage 整数幂 `q=ZZ(2) ** int(round(logQP))`。

## Dry Run

```bash
go run . \
  -dryRunParams=true \
  -logN=17 \
  -layers=1 \
  -run=ours \
  -leafEngine=lattigo \
  -q0Bits=45 \
  -circuitLevels=15 \
  -circuitPrimeBits=35 \
  -numP=5 \
  -leafNumP=5 \
  -defaultScale=35 \
  -splitFirstLeafPreprocess=true
```

应看到：

```text
top:  logN=17 N=131072 qCount=31 pCount=5 logQ=1391.0 logQP=1696.0
leaf: logN=16 N=65536  qCount=31 pCount=5 logQ=1391.0 logQP=1696.0
128-bit reference: leaf logQP <= 1747, margin=51.0 bits [PASS]
```

## 正式运行命令

只跑 split route，不跑 baseline：

```bash
go run . \
  -logN=17 \
  -layers=1 \
  -run=ours \
  -leafEngine=lattigo \
  -q0Bits=45 \
  -circuitLevels=15 \
  -circuitPrimeBits=35 \
  -numP=5 \
  -leafNumP=5 \
  -defaultScale=35 \
  -splitFirstLeafPreprocess=true \
  -parallelLeaves=true \
  -leafWorkers=2 \
  -reps=1 \
  -outputTSV=results/split_top17_leaf16_lattigo_Lout15.tsv \
  2>&1 | tee logs/split_top17_leaf16_lattigo_Lout15.log
```

## 已得实验结果

日志：

```text
logs/split_top17_leaf16_lattigo_Lout15.log
```

TSV：

```text
results/split_top17_leaf16_lattigo_Lout15.tsv
```

关键结果：

```text
top logN = 17
leaf logN = 16
leaves = 2
top logQP = 1696.000268
leaf logQP = 1696.000268
output level = 15
```

时间：

```text
Total         = 51.651618667 s
Split         = 0.100103292 s
leafScaleDown = 0.000995208 s
leafModUp     = 1.137921875 s
leafC2S       = 27.825931500 s
leafEvalMod   = 6.483571875 s
leafS2C       = 14.964930292 s
Merge         = 1.138164625 s
```

精度：

```text
input ciphertext vs plaintext L1 precision      ≈ 21.09 bits
split BTS output vs input ciphertext L1 precision ≈ 20.47 bits
```

## 结论

top17->leaf16 的 Lattigo-only split bootstrapping 已经跑通，并且在 PQ/logQP 口径下满足 128-bit 安全估计。和 top16->leaf15 不同，leaf16 的安全窗口足够大，可以保留 15 个输出 level，不需要像 leaf15 那样牺牲到 output level 3。
