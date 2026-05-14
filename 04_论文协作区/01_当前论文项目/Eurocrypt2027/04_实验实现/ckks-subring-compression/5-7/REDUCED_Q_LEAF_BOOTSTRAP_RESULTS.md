# Reduced-Q Leaf Bootstrapping 实验记录

日期：2026-05-08

目标：验证降低 \(Q\) 后是否还能服务完整 leaf bootstrapping：

\[
\text{SplitNew}
\to
\text{reduced-Q leaf bootstrapping}
\to
\text{Merge}.
\]

## 1. Lattigo 默认 leaf bootstrapping：full-Q 对照

命令：

```bash
go run ./5-7 -run=ours \
  -logN=16 -layers=1 -reps=1 \
  -q0Bits=34 -circuitPrimeBits=20 -defaultScale=25 \
  -circuitLevels=1 -numP=1 \
  -leafQCount=17 \
  -splitFirstLeafPreprocess
```

结果：

| 项 | 结果 |
|---|---:|
| leaf `logN` | 15 |
| leaf `logQP` | 936.2 |
| leaf max level | 16 |
| 是否满足 leaf15 128-bit 参考 | 否，超过 868 |
| 是否完整跑通 | 是 |
| 输出 level | 1 |
| Avg L1 precision | 11.49 bits |
| 总时间 | 3.88 s |
| Split | 6.79 ms |
| leafScaleDown | 0.29 ms |
| leafModUp | 116.79 ms |
| leafC2S | 2.42 s |
| leafEvalMod | 1.09 s |
| leafS2C | 241.24 ms |
| Merge | 9.48 ms |

结论：Lattigo 默认 leaf bootstrapping 在 full-Q 下能跑，但安全性不满足 leaf15 的参考上限。

## 2. Lattigo 默认 leaf bootstrapping：reduced-Q 尝试

命令：

```bash
go run ./5-7 -run=ours \
  -logN=16 -layers=1 -reps=1 \
  -q0Bits=34 -circuitPrimeBits=20 -defaultScale=25 \
  -circuitLevels=1 -numP=1 \
  -leafQCount=16 \
  -splitFirstLeafPreprocess
```

结果：

| 项 | 结果 |
|---|---:|
| leaf `logQP` | 880.2 |
| leaf max level | 15 |
| 是否满足 leaf15 128-bit 参考 | 否，仍略高于 868 |
| 是否完整跑通 | 否 |
| 失败阶段 | `bootstrapping.NewEvaluator` 初始化 DFT 矩阵 |
| 失败原因 | DFT matrix 仍访问 level index 16，但 reduced basis 只有 16 个 Q prime，合法索引为 0--15 |

错误摘要：

```text
panic: runtime error: index out of range [16] with length 16
github.com/tuneinsight/lattigo/v6/circuits/ckks/dft.NewMatrixFromLiteral
```

进一步降到 `leafQCount=15` 也同样失败：

| leafQCount | leaf logQP | 状态 |
|---:|---:|---|
| 17 | 936.2 | 能跑 |
| 16 | 880.2 | evaluator 构造失败 |
| 15 | 824.2 | evaluator 构造失败 |

结论：

\[
\boxed{
\text{Lattigo 默认 leaf bootstrapping 不能直接服务 reduced-Q prefix。}
}
\]

原因不是 split/drop 语义不行，而是默认 CtS/EvalMod/StC 参数需要完整 level 深度。

## 3. Anyu/subKey leaf bootstrapping：原始 lowq 对照

命令：

```bash
go run ./5-7 -run=ours \
  -leafEngine=anyu -anyuProfile=lowq -anyuSubringLogN=12 \
  -logN=16 -layers=1 \
  -splitFirstLeafPreprocess -reps=1
```

参数：

| 项 | 值 |
|---|---:|
| top `logN` | 16 |
| leaf `logN` | 15 |
| subKey target `logN` | 12 |
| `logQ` | 868 |
| `logP` | 61 |
| `logQP` | 929 |
| EvalModBits | 55 |
| 是否满足 leaf15 128-bit 参考 | 否，929 > 868 |

结果：

| 项 | 结果 |
|---|---:|
| 是否完整跑通 | 是 |
| Avg L1 precision | 16.10 bits |
| 总时间 | 6.28 s |
| Split | 7.49 ms |
| leafScaleDown | 0.22 ms |
| leafModUp | 39.20 ms |
| leafC2S | 3.11 s |
| leafEvalMod | 2.82 s |
| leafS2C | 290.30 ms |
| Merge | 15.04 ms |

结论：原始 Anyu lowq 作为 leaf bootstrapping 能跑，精度约 16 bits，但 leaf15 安全性仍不够。

## 4. Anyu/subKey leaf bootstrapping：reduced-Q 到安全线内

命令：

```bash
go run ./5-7 -run=ours \
  -leafEngine=anyu -anyuProfile=lowq -anyuSubringLogN=12 \
  -anyuMaxLogQP=868 -anyuEvalModBits=52 \
  -logN=16 -layers=1 \
  -splitFirstLeafPreprocess -reps=1
```

参数：

| 项 | 值 |
|---|---:|
| `logQ` | 803 |
| `logP` | 61 |
| `logQP` | 864 |
| EvalModBits | 52 |
| 是否满足 leaf15 128-bit 参考 | 是，864 < 868 |

结果：

| 项 | 结果 |
|---|---:|
| 是否完整跑通 | 是 |
| Avg L1 precision | 13.16 bits |
| 总时间 | 4.90 s |
| Split | 5.76 ms |
| leafScaleDown | 0.12 ms |
| leafModUp | 13.32 ms |
| leafC2S | 2.30 s |
| leafEvalMod | 2.37 s |
| leafS2C | 207.30 ms |
| Merge | 4.53 ms |

结论：

\[
\boxed{
\text{reduced-Q leaf bootstrapping 可以完整跑通，并且能落入 leaf15 的 128-bit 参考线内。}
}
\]

代价是精度从约 16.10 bits 降到约 13.16 bits。

## 5. 更激进 reduced-Q

命令：

```bash
go run ./5-7 -run=ours \
  -leafEngine=anyu -anyuProfile=lowq -anyuSubringLogN=12 \
  -anyuMaxLogQP=830 -anyuEvalModBits=48 \
  -logN=16 -layers=1 \
  -splitFirstLeafPreprocess -reps=1
```

参数：

| 项 | 值 |
|---|---:|
| `logQ` | 763 |
| `logP` | 61 |
| `logQP` | 824 |
| EvalModBits | 48 |
| 是否满足 leaf15 128-bit 参考 | 是，824 < 868 |

结果：

| 项 | 结果 |
|---|---:|
| 是否完整跑通 | 是 |
| Avg L1 precision | 9.14 bits |
| 总时间 | 4.80 s |
| Split | 5.92 ms |
| leafScaleDown | 0.12 ms |
| leafModUp | 13.56 ms |
| leafC2S | 2.31 s |
| leafEvalMod | 2.27 s |
| leafS2C | 199.91 ms |
| Merge | 5.00 ms |

结论：继续降低 \(Q\) 仍能跑通，但 EvalMod precision 明显下降。

## 6. 当前总表

| 方法 | leaf logQP | 安全性 | 是否跑通 | 精度 | 总时间 | 主要结论 |
|---|---:|---|---|---:|---:|---|
| Lattigo 默认 full-Q | 936.2 | 不安全 | 是 | 11.49 bits | 3.88 s | 能跑但不安全 |
| Lattigo 默认 reduced-Q | 880.2 | 仍略不安全 | 否 | - | - | evaluator 构造失败 |
| Lattigo 默认 reduced-Q | 824.2 | 安全 | 否 | - | - | evaluator 构造失败 |
| Anyu/subKey lowq 原始 | 929 | 不安全 | 是 | 16.10 bits | 6.28 s | 精度最好但不安全 |
| Anyu/subKey reduced | 864 | 安全 | 是 | 13.16 bits | 4.90 s | 当前最重要正结果 |
| Anyu/subKey aggressive reduced | 824 | 安全 | 是 | 9.14 bits | 4.80 s | 能跑但精度偏低 |

## 7. 当前判断

1. 单纯 split/drop/merge 不是瓶颈。之前 q-sweep 已经说明 SplitNew 后 leaf ciphertext 可以降到很低 Q 而不破坏语义。

2. Lattigo 默认 bootstrapping 不能直接使用 reduced-Q prefix，因为默认 CtS/EvalMod/StC 的 level 布局需要完整 Q 深度。

3. Anyu/subKey 风格的可调参数能够服务 reduced-Q leaf bootstrapping。

4. 当前最有价值的参数点是：

\[
\log(QP)=864,\qquad \text{precision}\approx 13.16\text{ bits}.
\]

它满足 leaf15 的 128-bit 参考线：

\[
864 < 868.
\]

5. 精度损失主要来自 EvalMod prime 从 55 bits 降到 52 bits，以及没有额外 circuit level。

下一步应围绕 `logQP≈864--868` 微调：

- `EvalModBits=52/53`
- `CtSBits=48/49/50`
- `StCBits=28/29/30`
- `Mod1Degree=127` 是否可提高
- `DoubleAngle=3/4`

