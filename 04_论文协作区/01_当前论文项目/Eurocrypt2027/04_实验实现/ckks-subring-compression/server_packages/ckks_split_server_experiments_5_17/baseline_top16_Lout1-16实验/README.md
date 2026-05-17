# top16 direct baseline Lout 1-16 sweep

这个目录用于跑一个没有 split、没有 leaf、没有 ring-switch 优化的 direct baseline：

```text
top logN = 16
run = baseline
leafEngine = lattigo
target output level = 1..16
```

每个点固定使用：

```text
q0Bits = 45
circuitPrimeBits = 35
numP = 5
defaultScale = 35
```

只改变：

```text
circuitLevels = target output level
```

## 文件说明

```text
main.go
anyu25.go
verify_logN15_direct_bootstrap.go
go.mod
go.sum
deps/subring_bts-main/
vendor/
run_top16_baseline_lout1_16.sh
run_split_success_points.sh
run_split_level_sweeps.sh
scripts/summarize_baseline_log.py
logs/
results/
```

`deps/subring_bts-main/` 是当前实验使用的本地 Lattigo fork，`go.mod` 使用相对路径：

```text
replace github.com/tuneinsight/lattigo/v6 => ./deps/subring_bts-main
```

因此上传服务器时，把整个目录上传即可，不要只上传 `main.go`。

`vendor/` 已经包含 Go 间接依赖；运行脚本会设置：

```text
GOFLAGS=-mod=vendor
```

因此服务器不需要连接 GitHub 下载 Go module。

## 运行命令

在服务器进入本目录后运行：

```bash
chmod +x run_top16_baseline_lout1_16.sh
./run_top16_baseline_lout1_16.sh
```

脚本会同时：

```text
1. 在终端正常输出所有运行日志；
2. 写入总日志 master.log；
3. 为每个 level 写单独日志；
4. 生成一个汇总 TSV。
```

默认跑 `1..16`，每个点 `reps=1`。

如果服务器中途断了，可以从某个 level 继续，例如：

```bash
START_LEVEL=8 END_LEVEL=16 ./run_top16_baseline_lout1_16.sh
```

如果想重复多次：

```bash
REPS=3 ./run_top16_baseline_lout1_16.sh
```

## 输出位置

每次运行会创建带时间戳的目录，例如：

```text
logs/top16_baseline_lout1_16_YYYYMMDD_HHMMSS/
results/top16_baseline_lout1_16_YYYYMMDD_HHMMSS/
```

总日志：

```text
logs/top16_baseline_lout1_16_YYYYMMDD_HHMMSS/master.log
```

每个 level 的单独日志：

```text
logs/top16_baseline_lout1_16_YYYYMMDD_HHMMSS/baseline_top16_Lout1.log
...
logs/top16_baseline_lout1_16_YYYYMMDD_HHMMSS/baseline_top16_Lout16.log
```

最重要的汇总结果：

```text
results/top16_baseline_lout1_16_YYYYMMDD_HHMMSS/summary.tsv
```

Go 程序原生 TSV：

```text
results/top16_baseline_lout1_16_YYYYMMDD_HHMMSS/raw_go_summary.tsv
```

## 怎么看结果

优先看 `summary.tsv`。字段含义：

```text
target_lout            命令中设置的 circuitLevels
actual_output_level    程序日志中实际输出的 ciphertext level
logQP                  实际 bootstrapping 参数 logQP
qCount                 Q prime 数量
total_sec              direct bootstrapping 总时间
stc_sec                SlotsToCoeffs 时间
modup_sec              ModUp / bootstrap preprocess 时间
cts_sec                CoeffsToSlots 时间
evalmod_sec            EvalMod 时间
avg_l1_prec_bits       output vs input 的平均 L1 precision
log_file               该点的完整日志
```

正常情况下应看到：

```text
target_lout = actual_output_level
```

如果不相等，优先打开 `log_file` 查 `Output ciphertext level = ...` 和 panic/error 信息。

## 跑 split 对比点

本目录也包含 split 实现。服务器上如果要跑之前三个成功 split 点，用：

```bash
chmod +x run_split_success_points.sh
./run_split_success_points.sh
```

它会跑：

```text
top16 -> leaf15, output level 3, n15qp849
top17 -> leaf15, output level 3, n15qp849
top17 -> leaf16, output level 15, default Lout=15
```

输出位置：

```text
logs/split_success_points_YYYYMMDD_HHMMSS/
results/split_success_points_YYYYMMDD_HHMMSS/split_success_points.tsv
```

这几个 split 结果应该和 direct baseline sweep 在同一台服务器上比较，避免本机/服务器 CPU 差异污染 speedup。

## 跑 split level sweep

如果要扫 split 输出 level，用：

```bash
chmod +x run_split_level_sweeps.sh
./run_split_level_sweeps.sh
```

它会跑：

```text
top16 -> leaf15: output level 1..3
top17 -> leaf15: output level 1..3
top17 -> leaf16: output level 1..15
```

其中 leaf15 的 level 1/2/3 使用 `n15lout1/n15lout2/n15lout3` preset；leaf16 的 level 1..15 使用 default schedule 并改变 `circuitLevels`。

输出位置：

```text
logs/split_level_sweeps_YYYYMMDD_HHMMSS/
results/split_level_sweeps_YYYYMMDD_HHMMSS/split_level_sweeps.tsv
```

## 上传服务器

推荐从上级目录打包：

```bash
tar -czf baseline_top16_Lout1-16实验.tar.gz baseline_top16_Lout1-16实验
```

服务器解压后：

```bash
tar -xzf baseline_top16_Lout1-16实验.tar.gz
cd baseline_top16_Lout1-16实验
chmod +x run_top16_baseline_lout1_16.sh
./run_top16_baseline_lout1_16.sh
```

服务器需要有：

```text
Go 1.25.x 或兼容版本
python3
```

当前包已经包含 `vendor/`，服务器不需要联网下载 Go 依赖。
