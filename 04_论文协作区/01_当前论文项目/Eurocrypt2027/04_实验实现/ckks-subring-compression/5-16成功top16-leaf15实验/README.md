# 5-16 成功 top16-leaf15 实验

本目录保存一个已经跑通的 **Lattigo-only** split bootstrapping 实验点：

```text
top logN = 16
leaf logN = 15
split layers = 1
leaf engine = Lattigo 自带 CKKS bootstrapping
slot recovery = 不使用 B_k / B_k^{-1}
security accounting = PQ / logQP
```

不使用 AnyuWang/subKey bootstrapping，也不需要跑 top17/leaf16。

## 目录内容

```text
main.go
anyu25.go
verify_logN15_direct_bootstrap.go
go.mod
go.sum
logs/
results/
```

说明：

- `main.go` 是 split top16-leaf15 实验代码。
- `anyu25.go` 只是为了保留原程序的编译依赖；本实验命令固定使用 `-leafEngine=lattigo`，不走 Anyu 路线。
- `verify_logN15_direct_bootstrap.go` 带 `directverify` build tag，用于单独验证 leaf logN=15 direct bootstrapping。
- 本目录单独带 `go.mod`，因为目录名含中文，不能作为上级 Go module 的普通 package path 直接 `go run .`。

## 成功参数

新增的 Lattigo preset 名为：

```text
n15qp849
```

它来自 Lattigo 的 N15 bootstrapping preset，并减少一个输出 computation level。

参数结构：

```text
Residual Q = [40, 31, 31, 31]
P          = [56, 56]
StC        = [30, 30]
EvalMod    = 8 × 55
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

## 参数选择方法

先验证 leaf ring `logN=15` 上 Lattigo 自带 direct bootstrapping 是否可行。

原始 N15 preset：

```text
Residual Q = [40, 31, 31, 31, 31]
logQP = 880
output level = 4
```

用本地 estimator 按 PQ 口径估计：

```bash
PATH="/opt/anaconda3/envs/sage/bin:$PATH" \
/opt/anaconda3/envs/sage/bin/python /Users/mac/opt/lattice-estimator/fhe_estimate.py \
  --n 32767 \
  --logQP 880 \
  --sigma 3.2 \
  --tag logN15_direct_bts_QP880 \
  --attacks primal_usvp,primal_bdd,dual_hybrid
```

结果约为：

```text
minimum estimated bits = 126.11
```

所以原始 `logQP=880` 略低于 128-bit。减少一个 31-bit 输出 level 后：

```text
Residual Q = [40, 31, 31, 31]
logQP = 849
output level = 3
```

Estimator：

```bash
PATH="/opt/anaconda3/envs/sage/bin:$PATH" \
/opt/anaconda3/envs/sage/bin/python /Users/mac/opt/lattice-estimator/fhe_estimate.py \
  --n 32767 \
  --logQP 849 \
  --sigma 3.2 \
  --tag logN15_direct_bts_QP849_Lout3 \
  --attacks primal_usvp,primal_bdd,dual_hybrid
```

结果：

```text
minimum estimated bits = 130.90
```

因此该 leaf15 参数在 PQ 口径下满足 128-bit。

## Dry Run

先检查 top/leaf 参数：

```bash
go run . \
  -dryRunParams=true \
  -logN=16 \
  -layers=1 \
  -run=ours \
  -leafEngine=lattigo \
  -lattigoPreset=n15qp849 \
  -splitFirstLeafPreprocess=true
```

应看到：

```text
top:  logN=16 qCount=15 pCount=2 logQ=737.0 logQP=849.0
leaf: logN=15 qCount=15 pCount=2 logQ=737.0 logQP=849.0
leaf logQP <= 868, margin=19.0 bits [PASS]
```

## 正式运行命令

只跑 split route，不跑 baseline：

```bash
LOG="logs/split_top16_leaf15_lattigo_n15qp849_$(date +%Y%m%d_%H%M%S).log"

go run . \
  -logN=16 \
  -layers=1 \
  -run=ours \
  -leafEngine=lattigo \
  -lattigoPreset=n15qp849 \
  -splitFirstLeafPreprocess=true \
  -parallelLeaves=true \
  -leafWorkers=2 \
  -reps=1 \
  -outputTSV=results/split_top16_leaf15_lattigo_n15qp849.tsv \
  2>&1 | tee "$LOG"
```

## 已得实验结果

日志：

```text
logs/split_top16_leaf15_lattigo_n15qp849_20260516_215137.log
```

TSV：

```text
results/split_top16_leaf15_lattigo_n15qp849.tsv
```

关键结果：

```text
top logN  = 16
leaf logN = 15
leaves    = 2
top logQP = 849.000442
leaf logQP = 849.000442
output level = 3
output circuit moduli = [31 31 31]
```

时间：

```text
Total         = 5.119727667 s
Split         = 0.010600583 s
leafScaleDown = 0.000215459 s
leafModUp     = 0.077790000 s
leafC2S       = 3.073341708 s
leafEvalMod   = 0.780989667 s
leafS2C       = 1.152307416 s
Merge         = 0.024482834 s
```

精度：

```text
input ciphertext L1 precision      ≈ 18.08 bits
split BTS output vs input precision ≈ 17.13 bits
```

安全性：

```text
leaf logQP = 849
estimator minimum ≈ 130.90 bits
```

## 结论

该实验点证明：

```text
top logN=16 -> split leaf logN=15
Lattigo 自带 CKKS bootstrapping
不使用 slot recovery
PQ 口径 128-bit 安全
输出 level = 3
```

可以成功运行，并保持约 17-bit bootstrapping precision。

这正好支撑论文的新定位：split bootstrapping 不应被表述为 full-output-level top BTS 的无损替代，而应表述为一个 **low-output-level secure fast refresh primitive**。
