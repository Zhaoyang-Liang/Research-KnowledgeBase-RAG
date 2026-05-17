# 5-16: low-output-level split BTS experiment

本实验对应新的论文定位：

\[
\BTS_N^{\mathrm{lowQ}}:
\mathrm{Ct}_{q_0,N}\to \mathrm{Ct}_{Q_s,N}
\]

对比

\[
\Merge_2\circ(\BTS_{N/2}^{\mathrm{safe}})^{\oplus 2}\circ\Split_2:
\mathrm{Ct}_{q_0,N}\to \mathrm{Ct}_{Q_s,N}.
\]

两条路线使用相同的低输出 modulus budget。实验目标不是替代 full-output-level top BTS，而是在后续只需要较少 levels 时比较 direct top lowQ refresh 和 split lowQ refresh 的时间。

## 默认命令

```bash
./run_lowq_lattigo.sh
```

默认参数：

```text
logN=16, layers=1, leafLogN=15
q0Bits=34, circuitLevels=1, circuitPrimeBits=20
numP=1, defaultScale=25
```

这组 Lattigo 默认 bootstrapping 参数给出约：

```text
top/leaf logQ ≈ 875 bits
top/leaf logQP ≈ 936 bits
```

因此它是“约 868 bits”的低输出预算实验。当前这份 Lattigo fork 的 bootstrapping keygen 需要至少 1 个 P prime，所以实际 `logQP` 会多一个约 61-bit 的 P prime；公平比较时 direct top lowQ 和 split lowQ 使用同一组 Q/P。这里先用 Lattigo 自带 BTS 和当前 split/merge 方法，不接入 subkey/Anyu 引擎。

## 输出文件

- `results/lowq_lattigo.tsv`: 每次运行的一行汇总结果。
- `results/lowq_lattigo_p35.tsv`: 20-bit 级精度诊断实验结果。
- `results/lowq_lattigo_p35_sec128.tsv`: leaf 128-bit 安全调整后的实验结果。
- `logs/lowq_lattigo_*.log`: 完整终端日志。
- `SECURITY.md`: 使用 `/Users/mac/opt/lattice-estimator/fhe_estimate.py` 得到的安全性估算。

TSV 中记录：

- top/leaf `logQ`, `logQP`, Q/P prime 数量；
- direct top lowQ BTS 的分项时间和总时间；
- split lowQ BTS 的 split、leaf ScaleDown、leaf ModUp、leaf C2S、leaf EvalMod、leaf S2C、merge 分项时间；
- 两条路线的精度；
- speedup。

## 当前要验证的结论

如果两条路线输出相同的 \(Q_s\)，那么比较是公平的：

\[
\mathrm{Ct}_{q_0,N}\to\mathrm{Ct}_{Q_s,N}.
\]

期望看到：

\[
T_{\mathrm{split-lowQ}}
<
T_{\mathrm{top-lowQ}}.
\]

这支持论文中的 time-level tradeoff 表述：牺牲 full top BTS 的输出 level，换取更快的安全 refresh。

## 2026-05-16 单次结果

命令：

```bash
./run_lowq_lattigo.sh
```

日志：

```text
logs/lowq_lattigo_20260516_151144.log
```

结果汇总：

| method | output level | output circuit moduli | precision | total time |
|---|---:|---|---:|---:|
| direct top lowQ BTS | 1 | `[20]` | 12.08 bits | 50.63 s |
| split lowQ BTS, k=2 | 1 | `[20]` | 11.50 bits | 24.62 s |

Speedup:

\[
50.63/24.62 \approx 2.06\times.
\]

分项时间：

| split lowQ stage | time |
|---|---:|
| Split | 0.025 s |
| leaf ScaleDown | 0.001 s |
| leaf ModUp | 0.645 s |
| leaf C2S | 20.723 s |
| leaf EvalMod | 1.598 s |
| leaf S2C | 1.567 s |
| Merge | 0.065 s |

这个结果初步支持新的实验定位：在相同低输出 level budget 下，split no-slot-recovery BTS 可以比 direct top lowQ BTS 更快；本次单次实验约为 \(2.06\times\)。

## 2026-05-16 更高精度 lowQ 结果

上面的默认脚本使用：

```text
q0Bits=34, circuitPrimeBits=20, defaultScale=25
```

因此输入密文自身的 L1 precision 只有约 12 bits。这个配置更像是一个极低 scale 的 stress test，不适合作为论文主结果。

更合理的 low-output-level 配置是：

```bash
./run_lowq_lattigo_p35.sh
```

参数：

```text
q0Bits=45, circuitLevels=1, circuitPrimeBits=35
numP=1, defaultScale=35
top/leaf logQ ≈ 901 bits
top/leaf logQP ≈ 962 bits
```

单次结果：

| method | output level | output circuit moduli | precision | total time |
|---|---:|---|---:|---:|
| direct top lowQ BTS | 1 | `[35]` | 21.96 bits | 58.21 s |
| split lowQ BTS, k=2 | 1 | `[35]` | 21.46 bits | 26.36 s |

Speedup:

\[
58.21/26.36 \approx 2.21\times.
\]

这组结果更接近预期的 20-bit 级精度，同时仍然保持低输出 level，并且 split 方法有约 \(2.21\times\) 加速。

## 安全性检查

使用 `/Users/mac/opt/lattice-estimator/fhe_estimate.py` 估算后，当前 `logN=16` 的 p35 诊断实验不是 128-bit split 参数：

```text
top logN=16, leaf logN=15, leaf logQP≈962 -> minimum ≈115.05 bits
```

即使只按 `leaf logQ≈901` 做诊断估算，也只有约 `123.07` bits。`leaf logN=15` 下要达到 128-bit，估算上需要约：

```text
logQP <= 868
```

这里的 `logQ/logQP` 指 bootstrapping evaluator 的总 modulus chain，不是最终输出密文仍保留这么多 bits。最终输出还要看 `out.Level()` 和 `LogQLvl(out.Level())`。

因此保留 `run_lowq_lattigo_p35.sh` 作为正确性和速度诊断脚本，但正式 128-bit 版本应使用：

```bash
./run_lowq_lattigo_p35_sec128.sh
```

该脚本使用：

```text
top logN=17, leaf logN=16, leaf logQP≈962
```

估算结果为：

```text
leaf minimum security ≈246.05 bits
```

也就是说，若当前 p35 低输出 level 配置不想牺牲精度，最直接的 128-bit 修正是把 top 环提升到 `logN=17`，让 split 后的 leaf 环是 `logN=16`。

## 接近 868-bit 总 Q-chain 的补充实验

严格的 `leaf logQP≈868` 候选：

```text
q0Bits=26, circuitPrimeBits=18, defaultScale=25, numP=0
top/leaf total logQ=logQP≈866.6
```

不能用原生 Lattigo 跑通，bootstrapping key generation 在 `P=0` 时崩溃。因此当前实现需要区分：

```text
输出 ciphertext Q 的 bits
和 key-switch/basis-extension 使用的 QP bits
```

把 `numP=1` 打开后，`q0Bits=26, defaultScale=25` 又会因为 ScaleDown 窗口太窄而失败。最近的可运行边界点是：

```text
q0Bits=30, circuitPrimeBits=18, defaultScale=22, numP=1
top/leaf total logQ≈870.6, logQP≈931.6
```

但这个点 split 精度崩掉：

| method | output level | output logQ | precision | total time |
|---|---:|---:|---:|---:|
| direct top lowQ BTS | 1 | 50 | 9.08 bits | 51.64 s |
| split lowQ BTS, k=2 | 1 | 50 | 0.97 bits | 25.67 s |

Estimator:

```text
Q-only:     127.51 bits
actual QP: 118.92 bits
```

第一个接近边界且 split 正常工作的点是：

```text
q0Bits=34, circuitPrimeBits=18, defaultScale=25, numP=1
top/leaf total logQ≈874.6, logQP≈935.6
```

结果：

| method | output level | output logQ | precision | total time |
|---|---:|---:|---:|---:|
| direct top lowQ BTS | 1 | 54 | 12.08 bits | 53.64 s |
| split lowQ BTS, k=2 | 1 | 54 | 11.51 bits | 27.43 s |

Estimator:

```text
Q-only:     126.94 bits
actual QP: 118.37 bits
```

Speedup:

\[
53.64/27.43\approx 1.96\times.
\]

参数选择流程见 `PARAM_SELECTION.md`。干跑扫描脚本见：

```bash
./select_lowq_params.py --security-modulus Q --max-log-sec 868 --target-logq 868 --top 20
```
