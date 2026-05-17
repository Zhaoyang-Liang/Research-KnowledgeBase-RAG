# 5-16 内核 subkey 实验

本目录整理 split 外层框架 + subkey/AnyuWang-style leaf bootstrapping 内核的兼容性实验。

定位：

```text
目的 = 验证 split 外层框架可以替换 leaf bootstrapping kernel
不是 = 主性能实验
```

主性能实验仍以 Lattigo leaf BTS 为主。这里的 subkey/Anyu leaf kernel 用来支撑论文中的 backend-agnostic / kernel-compatible 叙述。

## 目录内容

```text
main.go
anyu25.go
verify_logN15_direct_bootstrap.go
go.mod
go.sum
vendor/
deps/subring_bts-main/
run_subkey_dryrun.sh
run_subkey_current_params.sh
run_subkey_legacy_lowq_p1.sh
logs/
results/
```

本目录是离线 Go module：

```text
replace github.com/tuneinsight/lattigo/v6 => ./deps/subring_bts-main
GOFLAGS=-mod=vendor
```

因此服务器不需要连接 GitHub 下载 Go 依赖。

## 当前参数整理

### top17 -> leaf16, subkey i1

推荐作为主要 subkey 兼容性点。

Dry-run 参数：

```text
leaf logN = 16
profile = i1
logQ  ~= 1425
logQP ~= 1730
P = 5 * 61
local leaf16 PQ reference: PASS, margin ~= 17 bits
```

命令由 `run_subkey_current_params.sh` 自动运行。

### top16/top17 -> leaf15, subkey lowq, exploratory current-PQ attempt

为了贴近当前 leaf15 严格 PQ 口径，脚本使用：

```text
anyuProfile = lowq
anyuMaxLogQP = 868
anyuNumP = 0
```

Dry-run 参数：

```text
leaf logN = 15
logQ  ~= 868
logQP ~= 868
P = 0, i.e. LogP = []
local leaf15 PQ reference: at boundary
```

这里的 `P=0` 不是 AnyuWang 论文参数。它表示主 bootstrapping 参数不加入
key-switching auxiliary primes，
因此 `logQP = logQ`。这样做的唯一目的，是把 leaf15 的 subkey 参数压到当前严格
`logQP <= 868` 口径内。

注意这不等于 subring encapsulation key 的 `P' = 0`。Lattigo fork 在生成
subring key-switch 参数时会单独构造 `SubringParamsModUp`；当 `N'=2^12` 时，
`P'` 仍按约 61 bit 选取，与 AnyuWang Table 3 的 `log P' = 61` 对齐。

这个点只能作为探索性实验：它 dry-run 可以构造参数和精确继承 top ring basis，但是否能
完整通过 keygen/runtime 需要服务器实跑确认。论文中若使用该点，必须明确标注为
`leaf15 strict-PQ attempt`，不能称为 paper-faithful AnyuWang parameter。

### legacy lowq P=1 fallback

老 lowq profile 默认：

```text
logQ  ~= 868
logQP ~= 929
P = 1 * 61
```

这不满足当前 strict leaf15 PQ<=868 口径，因此只作为 compatibility fallback。如果 P=0 跑不通，可以跑这个证明外层框架能接 subkey 内核，但不能作为严格 128-bit leaf15 性能点。

## subring-secret 参数处理

`anyu25.go` 中 subring-secret 相关参数由 `anyuSubringLogN` 控制，默认值为 12，对应
AnyuWang CKKS 表格中的 `N' = 4096`。当前实现会按实际 leaf ring 自动计算：

```text
bottomLevelSubringSecretExtDegree = 2^(leafLogN - anyuSubringLogN)
StCSubringSecretExtDegree         = bottomLevelSubringSecretExtDegree / 2
Xs                                = uniform ternary over the subring
K                                 = FindSuitableKSimple(..., LogFailure=-16)
CtS ManualDepthSplit              = derived from extDegree
```

因此：

```text
leaf15: N=2^15, N'=2^12, extDegree=8
leaf16: N=2^16, N'=2^12, extDegree=16
```

这些参数已经在代码层面接入 `CoeffsToSlotsParameters.SubringExtDegree`、
`bootstrapping.Parameters.SubringSecretExtDegree` 和 `mod1.K`。不过，`lowq/P=0`
是我们的压缩适配，不是论文原始表格；最接近 AnyuWang Table 3/4 的兼容性实验仍是
`top17 -> leaf16, profile=i1, P=5*61, logQP ~= 1730`。

## 先做 dry-run

```bash
chmod +x run_subkey_dryrun.sh
./run_subkey_dryrun.sh
```

确认：

```text
Ring-basis exact check: PASS
leaf logQP 与预期一致
```

## 跑当前参数 subkey 兼容性实验

```bash
chmod +x run_subkey_current_params.sh
./run_subkey_current_params.sh
```

它会跑：

```text
top16 -> leaf15, Anyu lowq, P=0, logQP~=868
top17 -> leaf15, Anyu lowq, P=0, logQP~=868
top17 -> leaf16, Anyu i1,  logQP~=1730
```

输出：

```text
logs/subkey_current_params_YYYYMMDD_HHMMSS/
results/subkey_current_params_YYYYMMDD_HHMMSS/subkey_current_params.tsv
```

## fallback

如果 P=0 的 leaf15 点在 keygen/runtime 失败，可以跑：

```bash
chmod +x run_subkey_legacy_lowq_p1.sh
./run_subkey_legacy_lowq_p1.sh
```

注意这个 fallback 不满足当前 leaf15 strict PQ 口径，只用于证明 backend 兼容性。

## 论文使用建议

表格建议命名为 compatibility experiment：

```text
outer framework | inner kernel | top logN | leaf logN | profile | logQP | output level | precision | time | status
```

结论应写成：

```text
The split framework is compatible with different leaf bootstrapping kernels.
```

不要把 subkey 结果直接混入 Lattigo 主性能表，除非它的参数安全口径和输出 level 都已完全对齐。
