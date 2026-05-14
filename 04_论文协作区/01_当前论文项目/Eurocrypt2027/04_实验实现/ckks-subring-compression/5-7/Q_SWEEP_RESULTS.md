# Lattigo Split 后降 Q 的初步实验

日期：2026-05-08

目的：比较“Lattigo 自带 split 原语后再降低 ciphertext Q”的实际表现，并分析它与构造式低模数 split 的区别。

## 实验命令

### 1. 5-7 默认风格参数，`logN=16`

```bash
go run ./5-7 -run=qsweep \
  -logN=16 -layers=1 -reps=1 \
  -q0Bits=34 -circuitPrimeBits=20 -defaultScale=25 \
  -circuitLevels=1 -numP=1
```

### 2. Anyu lowq / 之前成功实验相关参数，`logN=16`

```bash
go run ./5-7 -run=qsweep \
  -leafEngine=anyu -anyuProfile=lowq -anyuSubringLogN=12 \
  -logN=16 -layers=1 -reps=1
```

`-run=qsweep` 执行的流程是：

\[
\text{high-level top ciphertext}
\to
\text{Lattigo SplitNew}
\to
\text{leaf Resize/drop Q}
\to
\text{Lattigo MergeNew}
\to
\text{decode and compare}.
\]

注意：这个实验只测试 split/drop/merge 的语义和时间，不包含 leaf bootstrapping。

## 实验 A：5-7 默认风格参数

参数：

| 项 | 值 |
|---|---:|
| top `logN` | 16 |
| leaf `logN` | 15 |
| input level | 16 |
| input `logQ` | 876 |
| input `logQP` | 936.2 |
| input precision | 12.08 bits |

代表性结果：

| target level | target logQ | saved Q bits | precision vs plaintext | precision vs input ct | split | merge | total |
|---:|---:|---:|---:|---:|---:|---:|---:|
| 16 | 876 | 0 | 8.91 | 8.93 | 256 ms | 223 ms | 555 ms |
| 14 | 764 | 112 | 8.91 | 8.93 | 256 ms | 169 ms | 486 ms |
| 10 | 532 | 344 | 9.01 | 9.02 | 263 ms | 98 ms | 416 ms |
| 5 | 232 | 644 | 9.32 | 9.34 | 243 ms | 28 ms | 303 ms |
| 0 | 35 | 841 | 9.40 | 9.42 | 265 ms | 4 ms | 270 ms |

观察：

- Lattigo `SplitNew` 后直接 `Resize` 降 Q，split/merge roundtrip 没有崩。
- 从 `logQ=876` 一直降到 `logQ=35` 仍然可解。
- 总时间下降主要来自 `MergeNew` 在更低 level 上更快；`SplitNew` 时间几乎不变，因为 split 发生在降 Q 之前。
- 精度没有随降 Q 变坏，反而略升，说明这里的主要误差不是 ciphertext Q 不够，而是 split/merge key switching 噪声和初始加密噪声。

## 实验 B：Anyu lowq / 之前成功实验相关参数

参数：

| 项 | 值 |
|---|---:|
| top `logN` | 16 |
| leaf `logN` | 15 |
| input level | 17 |
| input `logQ` | 869 |
| input `logQP` | 929.0 |
| input precision | 22.09 bits |
| top Q chain | `[43,35,30,30,30,55x10,50x3]` |
| P | `[61]` |

代表性结果：

| target level | target logQ | saved Q bits | precision vs plaintext | precision vs input ct | split | merge | total |
|---:|---:|---:|---:|---:|---:|---:|---:|
| 17 | 869 | 0 | 21.34 | 21.67 | 261 ms | 224 ms | 560 ms |
| 16 | 819 | 50 | 21.34 | 21.67 | 240 ms | 205 ms | 514 ms |
| 15 | 769 | 100 | 21.34 | 21.67 | 248 ms | 200 ms | 512 ms |
| 12 | 609 | 260 | 21.34 | 21.68 | 255 ms | 130 ms | 441 ms |
| 8 | 389 | 480 | 21.34 | 21.68 | 248 ms | 57 ms | 345 ms |
| 4 | 169 | 700 | 21.35 | 21.68 | 247 ms | 22 ms | 295 ms |
| 0 | 43 | 826 | 21.34 | 21.68 | 256 ms | 8 ms | 265 ms |

观察：

- 在更接近 5-7 成功实验的 Anyu lowq 参数下，Lattigo split 后降 Q 也非常稳。
- 从 `logQ=869` 降到 `logQ=43`，split/merge roundtrip 精度仍有约 21 bits。
- 这说明“split 后 ciphertext 降 Q”本身并不难，至少对于 split/merge 语义是可行的。

## 安全性解释：为什么这还不等于解决问题

Lattigo 路线做的是：

\[
\text{full-Q ciphertext}
\xrightarrow{\text{SplitNew using full-Q keys}}
\text{full-Q leaf ciphertexts}
\xrightarrow{\text{Resize}}
\text{low-q leaf ciphertexts}.
\]

这里降低的是 split 之后的 ciphertext modulus。

但是 `SplitNew` 的 ring-switching evaluation keys 仍然是在原始 full \(Q/P\) basis 下生成和使用的。因此安全性仍然要面对：

\[
\log(QP)\approx 929
\]

而 leaf 维度 `logN=15` 的 128-bit 参考上限大约是：

\[
868.
\]

所以：

\[
\boxed{
\text{Lattigo split 后再 Resize，可以降低输出 ciphertext Q，但不能降低 split key material 的 Q。}
}
\]

这就是它和构造式低模数 split 的核心区别。

## 对构造式方法的意义

构造式方法想要的是：

\[
\text{top ciphertext}
\xrightarrow{\text{low-q functional/module key switching}}
\text{low-q leaf ciphertexts}.
\]

也就是说，辅助 key 本身就在低模数 \(qP\) 下生成。例如入口低模数取：

\[
\log q_0=43,\qquad \log P=61,
\]

则：

\[
\log(q_0P)\approx 104,
\]

远低于 leaf `logN=15` 的 128-bit 上限 868。

因此，构造式方法真正要解决的是：

\[
\boxed{
\text{把“降 Q”移动到 split/key-switching 发生之前，而不是 split 之后。}
}
\]

## 对 leaf bootstrapping 参数的提示

对于 Anyu lowq 参数，top/leaf full basis 有：

\[
\log Q=869,\qquad \log P=61,\qquad \log(QP)=929.
\]

如果 leaf `logN=15` 要满足 128-bit 参考上限 868，则至少需要降低到：

| target level | target logQ | target logQ + 61 | 是否低于 868 |
|---:|---:|---:|---|
| 17 | 869 | 930 | 否 |
| 16 | 819 | 880 | 否 |
| 15 | 769 | 830 | 是 |

也就是说，从安全表角度，leaf bootstrapping 至少要去掉约 100 bits 的 Q，才能让：

\[
\log(Q_{\mathrm{leaf}}P_{\mathrm{leaf}})\le 868.
\]

本次 Lattigo q-sweep 表明，单就 split/drop/merge 语义而言，降到 `targetLevel=15` 完全没问题。但 leaf bootstrapping 在这个 reduced chain 上是否还能保持 EvalMod 精度，需要单独实验。

