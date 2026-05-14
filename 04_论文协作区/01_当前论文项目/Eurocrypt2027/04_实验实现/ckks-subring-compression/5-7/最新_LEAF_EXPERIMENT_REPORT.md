# subKey Bootstrapping 叶子替换实验报告

日期：2026-05-08

本报告记录这次把 no-\(B_k\) 拆分方案中的叶子 Bootstrapping 替换为
**subKey Bootstrapping** 的实验历程、失败尝试、参数变化和当前最好结果。

说明：代码和命令行参数里还保留了 `anyu` 这个名字，这是因为代码最早从
Anyu Wang 那篇 subKey Bootstrapping 实验代码整理而来。论文和报告里建议统一称为
**subKey Bootstrapping**，不要称为 “Anyu Bootstrapping”。

## 1. 实验目标

原来的 no-\(B_k\) 路线是：

\[
\mathsf{Split}
\to
\bigoplus \mathsf{C2S}_{\mathrm{leaf}}
\to
\bigoplus \mathsf{EvalMod}_{\mathrm{leaf}}
\to
\bigoplus \mathsf{S2C}_{\mathrm{leaf}}
\to
\mathsf{Merge}.
\]

这次实验想做的是：把每个叶子上的普通 Bootstrapping 换成
subKey Bootstrapping，从而尝试叠加两类加速：

1. **no-\(B_k\) 拆分加速**：把一个大环 Bootstrapping 拆成多个小环 Bootstrapping。
2. **subKey Bootstrapping 加速**：在每个叶子 Bootstrapping 内部继续利用子环密钥结构。

目标路线可以写成：

\[
\mathsf{Split}
\to
\bigoplus \mathsf{BTS}^{\mathrm{subKey}}_{\mathrm{leaf}}
\to
\mathsf{Merge}.
\]

其中 \(\mathsf{BTS}^{\mathrm{subKey}}_{\mathrm{leaf}}\) 表示叶子上的
subKey Bootstrapping。

## 2. 安全性判断方式

本实验暂时使用下面的 128-bit 安全参考值：

| 环参数 | 环维度 | 128-bit 安全参考上限 |
|---:|---:|---:|
| `logN=15` | \(32768\) | \(\log(PQ)\approx 868\) |
| `logN=16` | \(65536\) | \(\log(PQ)\approx 1747\) |
| `logN=17` | \(131072\) | \(\log(PQ)\approx 3523\) |

拆分方案的安全瓶颈不能只看大环，而是要看拆分后的叶子环。

如果大环是：

\[
\log N_{\mathrm{top}}=16,
\]

并且只拆分一层，那么叶子环是：

\[
\log N_{\mathrm{leaf}}=15.
\]

这时叶子的 \(\log(PQ)\) 应该跟 868 比。

如果大环是：

\[
\log N_{\mathrm{top}}=17,
\]

并且只拆分一层，那么叶子环是：

\[
\log N_{\mathrm{leaf}}=16.
\]

这时叶子的 \(\log(PQ)\) 应该跟 1747 比。

这就是为什么后来发现 `top=17 -> leaf=16` 比 `top=16 -> leaf=15`
更适合：叶子环能安全承受的模数预算大很多。

## 3. 术语说明

| 记号或参数 | 中文含义 |
|---|---|
| `top logN` | 大环的 \(\log_2 N\) |
| `leaf logN` | 拆分后叶子环的 \(\log_2 n\) |
| `layers` | 向下拆分层数；`layers=1` 表示拆成 2 个叶子 |
| `subringLogN` | subKey Bootstrapping 内部使用的子环密钥维度 |
| `logQ` | 密文主模数 \(Q\) 的比特长度 |
| `logQP` | \(Q\) 加 key switching 辅助模数 \(P\) 后的总比特长度 |
| `StC` | 槽到系数变换（SlotToCoeff） |
| `CtS` | 系数到槽变换（CoeffToSlot） |
| `EvalMod` | Bootstrapping 中的模约化近似步骤 |
| `Mod1Degree` | EvalMod 使用的近似多项式次数 |
| `CtSFirstDepth` | CtS 第一段在子环密钥下折叠的深度 |
| 输出计算层 | Bootstrapping 后还能继续做多少普通乘法层级 |
| 平均 L1 精度 | 平均 L1 precision，单位是 bit，越高越好 |

## 4. 代码实现状态

当前代码位置：

- `5-7/main.go`
- `5-7/anyu25.go`

虽然命令行里还有 `-leafEngine=anyu`，但这只是代码内部名字。报告和论文中建议称为
subKey Bootstrapping。

目前代码已经做到：

1. 可以用 `-leafEngine=anyu` 调用叶子 subKey Bootstrapping。
2. 不再外部调用别人的 main 文件，而是把 subKey Bootstrapping 相关逻辑复制进
   `5-7`。
3. 顶层和叶子的 \(Q/P\) 模数基完全对齐。
4. 已经给 subKey CtS 阶段生成需要的 Galois key。
5. 已经加入 `ManualDepthSplit` 修正。
6. 支持命令行直接调主要参数。

编译检查通过：

```bash
env GOCACHE=/private/tmp/go-build go test ./5-7
```

输出：

```text
?    ckks-subring-compression/5-7    [no test files]
```

## 5. 实验历程

### 5.1 第一次接入：出现 Galois key 错误

最早尝试命令类似：

```bash
go run . \
  -run=ours \
  -leafEngine=anyu \
  -anyuProfile=lowq \
  -anyuSubringLogN=12 \
  -logN=16 \
  -layers=1 \
  -splitFirstLeafPreprocess \
  -reps=1
```

当时程序在叶子系数到槽变换（CoeffToSlot, CtS）阶段报错：

```text
cannot MultiplyByDiagMatrix:
Automorphism:
CheckAndGetGaloisKey:
evaluation key interface is nil
```

含义是：subKey Bootstrapping 的子环 CtS 阶段需要一套特殊的旋转密钥
（Galois keys），但是原来的 evaluation key 没有把这部分密钥传进去。

修正方法：

1. 给 Bootstrapping evaluation keys 增加 subring CtS 专用的旋转密钥。
2. 在 `NewEvaluator` 时把这套密钥挂到 subring CtS evaluator 上。

修正后，程序不再 panic。

### 5.2 第二次问题：程序能跑，但精度崩溃

修完 Galois key 后，程序可以跑完，但数值非常差。

典型现象是输出值变成几千级别，平均精度甚至是负数。这说明不是普通噪声变大，
而是语义缩放或线性变换结构不对。

后来定位到：subKey CtS 的第一段折叠深度不对。

对于叶子 `logN=15`、subKey 内部子环 `logN=12`：

\[
2^{15}/2^{12}=2^3.
\]

所以叶子 CtS 第一段只能折叠 3 层。

对于叶子 `logN=16`、subKey 内部子环 `logN=12`：

\[
2^{16}/2^{12}=2^4.
\]

所以叶子 CtS 第一段应该折叠 4 层。

修正后加入了 `ManualDepthSplit`：

| 情况 | 叶子环 | subKey 子环 | 叶子 CtS 第一段深度 |
|---|---:|---:|---:|
| `top=16, layers=1` | `logN=15` | `logN=12` | 3 |
| `top=17, layers=1` | `logN=16` | `logN=12` | 4 |

这个修正非常关键。没有这个修正，后面的实验没有意义。

### 5.3 第三步：先在 `top=16 -> leaf=15` 上压低 Q

此时大环是：

\[
\log N_{\mathrm{top}}=16.
\]

拆分一层后叶子是：

\[
\log N_{\mathrm{leaf}}=15.
\]

按照安全参考表，叶子 \(\log(PQ)\) 应该接近或低于 868。

所以我们尝试压缩 subKey Bootstrapping 的模数预算，也就是 `lowq`
配置。这一阶段的主要目的不是直接找到最终最优参数，而是先验证融合路线能不能稳定跑通。

### 5.4 第四步：发现 `top=16 -> leaf=15` 精度和安全性都不够漂亮

`top=16 -> leaf=15` 最好的实验结果约为：

```text
平均 L1 精度 = 16.74 bits
总时间 = 6.59s
```

这个结果说明融合路线能工作。

但是它有两个问题：

1. 只保留 1 个输出计算层。
2. 如果严格把 \(QP\) 都算进叶子安全预算，`logQP≈927` 超过 leaf15 的
   868 参考值。

所以它适合作为工程验证，不适合作为主打参数。

### 5.5 第五步：切换到 `top=17 -> leaf=16`

后来把大环提高到：

\[
\log N_{\mathrm{top}}=17.
\]

拆分一层后叶子是：

\[
\log N_{\mathrm{leaf}}=16.
\]

此时叶子的安全参考上限是：

\[
\log(PQ)\approx 1747.
\]

这就能承受 subKey Bootstrapping 原始 I1 风格的大模数预算。

更重要的是，叶子上的 subKey Bootstrapping 变成：

\[
2^{16}\to 2^{12},
\]

这正好接近 subKey Bootstrapping 论文中的主要实验设置。

### 5.6 第六步：`top=17 -> leaf=16` 的大预算配置明显更好

使用 I1 风格参数后，结果变成：

```text
平均 L1 精度 = 19.40 bits
总时间 = 52.03s
输出计算层 = 15
```

同时，直接在 `top logN=17` 大环上做单环 Bootstrapping 的 baseline 是：

```text
平均 L1 精度 = 18.61 bits
总时间 = 135.08s
输出计算层 = 15
```

所以当前拆分 + 叶子 subKey Bootstrapping 的路线大约有：

\[
135.08/52.03\approx 2.6\times
\]

的加速，并且这次实验里精度还略高。

## 6. 所有主要实验结果

### 6.1 `top=16 -> leaf=15` 的实验

共同设置：

```text
大环 top logN = 16
拆分层数 layers = 1
叶子 leaf logN = 15
subKey 子环 subringLogN = 12
叶子安全参考 log(PQ) ≈ 868
```

| 编号 | 尝试内容 | 主要参数 | logQ | logQP | 输出计算层 | 平均 L1 精度 | 时间 | 结果判断 |
|---:|---|---|---:|---:|---:|---:|---:|---|
| 1 | 默认低 Q | `lowq` | 868 | 929 | 1 | 16.16 bits | 6.73s | 能跑通，但安全叙事不干净 |
| 2 | 增加 EvalMod，降低 CtS | `EvalMod=56, CtS=46` | 866 | 927 | 1 | 16.42 bits | 6.55s | 略好 |
| 3 | 继续增加 EvalMod，CtS 更低 | `EvalMod=57, CtS=43` | 867 | 928 | 1 | 13.77 bits | 6.69s | CtS 太低，失败 |
| 4 | 降低 EvalMod，提高 CtS | `EvalMod=54, CtS=52` | 864 | 925 | 1 | 15.12 bits | 6.79s | EvalMod 不够，变差 |
| 5 | EvalMod 拉高，CtS 极低 | `EvalMod=60, CtS=33` | 867 | 928 | 1 | 3.77 bits | 6.62s | 明显失败 |
| 6 | 降低 StC，保留 CtS | `EvalMod=56, StC=26, CtS=50` | 866 | 927 | 1 | 16.74 bits | 6.59s | top16 目前最好 |
| 7 | StC 继续降低 | `EvalMod=57, StC=23, CtS=50` | 867 | 928 | 1 | 14.02 bits | 6.02s | StC 太低，失败 |
| 8 | 改成更大子环密钥维度 | `subringLogN=13` | 868 | 929 | 1 | -5.52 bits | 8.03s | 失败 |

这一组实验的主要教训：

1. 系数到槽变换（CoeffToSlot, CtS）不能压得太低。
2. 槽到系数变换（SlotToCoeff, StC）可以从 30 降到 26，但不能降到 23。
3. EvalMod 的 bit 数不是越高越好，因为会挤压 CtS/StC 的预算。
4. `subringLogN=13` 在当前拆分叶子设置下失败。
5. `top=16 -> leaf=15` 的路线可以证明融合可行，但不是最适合主推的参数。

### 6.2 `top=17 -> leaf=16` 的实验

共同设置：

```text
大环 top logN = 17
拆分层数 layers = 1
叶子 leaf logN = 16
subKey 子环 subringLogN = 12
叶子安全参考 log(PQ) ≈ 1747
```

| 编号 | 尝试内容 | 主要参数 | logQ | logQP | 输出计算层 | 平均 L1 精度 | 时间 | 结果判断 |
|---:|---|---|---:|---:|---:|---:|---:|---|
| 1 | 沿用 top16 的低 Q 调参 | `lowq, EvalMod=56, StC=26, CtS=50` | 866 | 927 | 1 | 15.47 bits | 45.66s | 安全裕量大，但精度和层数不够 |
| 2 | 默认低 Q | `lowq` | 868 | 929 | 1 | 14.64 bits | 42.97s | 不如低 Q 调参 |
| 3 | I1 风格大预算 | `i1` | 1425 | 1730 | 15 | 19.40 bits | 52.03s | 当前综合最好 |
| 4 | I1，加大 CtS bit | `i1, CtS=55` | 1434 | 1739 | 15 | 19.42 bits | 50.51s | 精度最高，但安全裕量小 |
| 5 | I1，提高 EvalMod 多项式次数 | `i1, Mod1Degree=255` | 1415 | 1720 | 13 | 19.39 bits | 45.33s | 没提升精度，还少 2 层 |
| 6 | 直接大环 baseline | `direct baseline, i1` | 1425 | 1730 | 15 | 18.61 bits | 135.08s | 对照组 |

这一组实验的主要教训：

1. `top=17 -> leaf=16` 明显比 `top=16 -> leaf=15` 更适合。
2. I1 风格大预算配置可以把精度拉回到约 19.4 bits。
3. 加大 CtS 到 55 只提升约 0.02 bits，但安全裕量从约 17 bits 降到约 8 bits。
4. 提高 EvalMod 多项式次数到 255 没有提升精度。
5. 当前误差很可能不主要来自 EvalMod 多项式次数，也不主要来自 CtS bit 数不足。

## 7. 补充对照实验：只有拆分、只有 subKey、直接 baseline

为了更清楚地判断“组合方案”的收益，本节补充三类对照：

1. **只有拆分**：使用 no-\(B_k\) 拆分路线，但叶子上用普通 Bootstrapping，
   不使用 subKey Bootstrapping。
2. **只有 subKey**：不做 no-\(B_k\) 拆分，直接在大环上使用 subKey
   Bootstrapping。
3. **直接普通 baseline**：不拆分，也不使用 subKey，直接在大环上做普通
   Bootstrapping。

注意：这里的时间都是本地单次运行结果，只能作为初步对照。

### 7.1 `top=17 -> leaf=16`，I1 风格主参数的对照

组合方案，也就是“拆分 + 叶子 subKey Bootstrapping”的命令是：

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

这里 `-leafEngine=anyu` 在代码里表示使用 subKey Bootstrapping。

对照实验结果：

| 方法 | 是否拆分 | 是否使用 subKey | 安全参考对象 | logQP | 安全裕量 | 输出计算层 | 平均 L1 精度 | 时间 | 备注 |
|---|---|---|---|---:|---:|---:|---:|---:|---|
| 组合方案 | 是 | 是 | leaf16，参考 1747 | 1730 | 约 17 bits | 15 | 19.40 bits | 52.03s | 当前主推 |
| 只有拆分 | 是 | 否 | leaf16，参考 1747 | 1696 | 约 51 bits | 15 | 20.47 bits | 63.19s | 精度更高，但慢于组合方案 |
| 只有 subKey | 否 | 是 | top17，参考 3523 | 1730 | 约 1793 bits | 15 | 18.61 bits | 135.08s | 不拆分时很慢 |
| 直接普通 baseline | 否 | 否 | top17，参考 3523 | 1696 | 约 1827 bits | 15 | 20.87 bits | 143.29s | 精度最高，但最慢 |

对应命令如下。

只有拆分：

```bash
go run ./5-7/. \
  -run=ours \
  -leafEngine=lattigo \
  -logN=17 \
  -layers=1 \
  -q0Bits=45 \
  -defaultScale=35 \
  -circuitLevels=15 \
  -circuitPrimeBits=35 \
  -numP=5 \
  -leafCircuitLevels=15 \
  -leafNumP=5 \
  -splitFirstLeafPreprocess \
  -reps=1
```

只有 subKey：

```bash
go run ./5-7/. \
  -run=baseline \
  -leafEngine=anyu \
  -anyuProfile=i1 \
  -anyuSubringLogN=12 \
  -logN=17 \
  -layers=1 \
  -splitFirstLeafPreprocess \
  -reps=1
```

直接普通 baseline：

```bash
go run ./5-7/. \
  -run=baseline \
  -leafEngine=lattigo \
  -logN=17 \
  -layers=1 \
  -q0Bits=45 \
  -defaultScale=35 \
  -circuitLevels=15 \
  -circuitPrimeBits=35 \
  -numP=5 \
  -reps=1
```

这一组的主要结论是：

1. 组合方案相比“只有 subKey”快约

   \[
   135.08/52.03\approx 2.6\times.
   \]

2. 组合方案相比“直接普通 baseline”快约

   \[
   143.29/52.03\approx 2.75\times.
   \]

3. “只有拆分”比组合方案精度高一些，约 20.47 bits vs 19.40 bits，但时间更慢：

   \[
   63.19/52.03\approx 1.21\times.
   \]

   也就是说，subKey 叶子替换带来了额外加速，但牺牲了约 1 bit 精度。

### 7.2 `top=17 -> leaf=16`，I1 + CtS55 参数的对照

组合方案命令：

```bash
go run ./5-7/. \
  -run=ours \
  -leafEngine=anyu \
  -anyuProfile=i1 \
  -anyuSubringLogN=12 \
  -anyuCtSBits=55 \
  -logN=17 \
  -layers=1 \
  -splitFirstLeafPreprocess \
  -reps=1
```

对照实验结果：

| 方法 | 是否拆分 | 是否使用 subKey | 安全参考对象 | logQP | 安全裕量 | 输出计算层 | 平均 L1 精度 | 时间 | 备注 |
|---|---|---|---|---:|---:|---:|---:|---:|---|
| 组合方案 | 是 | 是 | leaf16，参考 1747 | 1739 | 约 8 bits | 15 | 19.42 bits | 50.51s | 精度最高，但安全裕量小 |
| 只有拆分 | 是 | 否 | leaf16，参考 1747 | 1696 | 约 51 bits | 15 | 20.47 bits | 63.19s | 与 7.1 共用普通拆分对照 |
| 只有 subKey | 否 | 是 | top17，参考 3523 | 1739 | 约 1784 bits | 15 | 18.60 bits | 135.15s | 不拆分时仍很慢 |
| 直接普通 baseline | 否 | 否 | top17，参考 3523 | 1696 | 约 1827 bits | 15 | 20.87 bits | 143.29s | 与 7.1 共用普通 baseline |

只有 subKey 的命令：

```bash
go run ./5-7/. \
  -run=baseline \
  -leafEngine=anyu \
  -anyuProfile=i1 \
  -anyuSubringLogN=12 \
  -anyuCtSBits=55 \
  -logN=17 \
  -layers=1 \
  -splitFirstLeafPreprocess \
  -reps=1
```

这一组的主要结论是：

1. `CtS=55` 的组合方案比默认 I1 只高约 0.02 bits。
2. 但是 leaf16 安全裕量从约 17 bits 降到约 8 bits。
3. 因此它可以作为“最高精度观察点”，但不是最稳妥的主推参数。

### 7.3 `top=16 -> leaf=15`，lowq tuned 参数的对照

组合方案命令：

```bash
go run ./5-7/. \
  -run=ours \
  -leafEngine=anyu \
  -anyuProfile=lowq \
  -anyuSubringLogN=12 \
  -anyuEvalModBits=56 \
  -anyuStCBits=26 \
  -anyuCtSBits=50 \
  -logN=16 \
  -layers=1 \
  -splitFirstLeafPreprocess \
  -reps=1
```

对照实验结果：

| 方法 | 是否拆分 | 是否使用 subKey | 安全参考对象 | logQP | 安全裕量 | 输出计算层 | 平均 L1 精度 | 时间 | 备注 |
|---|---|---|---|---:|---:|---:|---:|---:|---|
| 组合方案 | 是 | 是 | leaf15，参考 868 | 927 | 约 -59 bits | 1 | 16.74 bits | 6.59s | 工程上能跑，但严格 QP 安全不干净 |
| 只有拆分，split-first | 是 | 否 | leaf15，参考 868 | 960 | 约 -92 bits | - | 失败 | - | leaf ScaleDown 条件不满足 |
| 只有拆分，旧路线 | 是 | 否 | leaf15，参考 868 | 960 | 约 -92 bits | 1 | 1.51 bits | 9.39s | 输出约放大 2 倍，语义归一化失败 |
| 只有 subKey | 否 | 是 | top16，参考 1747 | 927 | 约 820 bits | 1 | 15.92 bits | 32.88s | 不拆分时慢很多 |
| 直接普通 baseline | 否 | 否 | top16，参考 1747 | 960 | 约 787 bits | 1 | 21.97 bits | 60.85s | 精度高，但不是 leaf 安全对照 |

只有拆分，split-first 路线命令：

```bash
go run ./5-7/. \
  -run=ours \
  -leafEngine=lattigo \
  -logN=16 \
  -layers=1 \
  -q0Bits=43 \
  -defaultScale=35 \
  -circuitLevels=1 \
  -circuitPrimeBits=35 \
  -numP=1 \
  -leafCircuitLevels=1 \
  -leafNumP=1 \
  -splitFirstLeafPreprocess \
  -reps=1
```

失败原因：

```text
leaf 0 ScaleDown failed:
initial Q/Scale = 255.999912
< 0.5*Q[0]/MessageRatio = 512.000000
```

只有拆分，旧路线命令：

```bash
go run ./5-7/. \
  -run=ours \
  -leafEngine=lattigo \
  -logN=16 \
  -layers=1 \
  -q0Bits=43 \
  -defaultScale=35 \
  -circuitLevels=1 \
  -circuitPrimeBits=35 \
  -numP=1 \
  -leafCircuitLevels=1 \
  -leafNumP=1 \
  -reps=1
```

只有 subKey 命令：

```bash
go run ./5-7/. \
  -run=baseline \
  -leafEngine=anyu \
  -anyuProfile=lowq \
  -anyuSubringLogN=12 \
  -anyuEvalModBits=56 \
  -anyuStCBits=26 \
  -anyuCtSBits=50 \
  -logN=16 \
  -layers=1 \
  -splitFirstLeafPreprocess \
  -reps=1
```

直接普通 baseline 命令：

```bash
go run ./5-7/. \
  -run=baseline \
  -leafEngine=lattigo \
  -logN=16 \
  -layers=1 \
  -q0Bits=43 \
  -defaultScale=35 \
  -circuitLevels=1 \
  -circuitPrimeBits=35 \
  -numP=1 \
  -reps=1
```

这一组的主要结论是：

1. `top=16 -> leaf=15` 的组合方案确实快，约 6.59s。
2. 但如果完整计入 leaf15 的 \(QP\)，安全裕量是负的。
3. 普通 only-split 路线在这一组参数下不稳定：split-first 直接失败，旧路线输出归一化不对。
4. 因此 `top=16` 更适合作为工程探索记录，不适合作为论文主推安全参数。

### 7.4 代码路径

本节所有实验均在下面路径运行：

```text
/Users/mac/Desktop/Research-KB-RAG/04_论文协作区/01_当前论文项目/Eurocrypt2027/04_实验实现/ckks-subring-compression
```

主要代码文件：

```text
5-7/main.go
5-7/anyu25.go
```

报告文件：

```text
5-7/最新_LEAF_EXPERIMENT_REPORT.md
```

## 8. 当前综合效果最好的三个情况

排序规则：

1. 安全性第一。
2. 精度第二。
3. 加速比第三。

### 第一名：`top=17 -> leaf=16`，I1 风格 subKey Bootstrapping

推荐命令：

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

| 指标 | 数值 |
|---|---:|
| 大环 | `logN=17` |
| 叶子环 | `logN=16` |
| 叶子安全参考 | \(\log(PQ)\approx 1747\) |
| 实际 logQ | 1425 |
| 实际 logQP | 1730 |
| 安全裕量 | 约 17 bits |
| 输出计算层 | 15 |
| 平均 L1 精度 | 19.40 bits |
| 本方案时间 | 52.03s |
| 只有拆分时间 | 63.19s |
| 只有 subKey 时间 | 135.08s |
| 直接普通 baseline 时间 | 143.29s |
| 相对只有拆分加速 | 约 1.21 倍 |
| 相对只有 subKey 加速 | 约 2.6 倍 |
| 相对直接普通 baseline 加速 | 约 2.75 倍 |

这是目前最适合写进论文主实验表的参数。它的优点是：安全叙事干净、精度较高、
输出层数多、相对“只有拆分”和“只有 subKey”都能看到额外收益。

### 第二名：`top=17 -> leaf=16`，I1 风格但加大 CtS

命令：

```bash
go run ./5-7/. \
  -run=ours \
  -leafEngine=anyu \
  -anyuProfile=i1 \
  -anyuSubringLogN=12 \
  -anyuCtSBits=55 \
  -logN=17 \
  -layers=1 \
  -splitFirstLeafPreprocess \
  -reps=1
```

| 指标 | 数值 |
|---|---:|
| 大环 | `logN=17` |
| 叶子环 | `logN=16` |
| 叶子安全参考 | \(\log(PQ)\approx 1747\) |
| 实际 logQ | 1434 |
| 实际 logQP | 1739 |
| 安全裕量 | 约 8 bits |
| 输出计算层 | 15 |
| 平均 L1 精度 | 19.42 bits |
| 本方案时间 | 50.51s |
| 只有拆分时间 | 63.19s |
| 只有 subKey 时间 | 135.15s |
| 直接普通 baseline 时间 | 143.29s |
| 相对只有拆分加速 | 约 1.25 倍 |
| 相对只有 subKey 加速 | 约 2.68 倍 |
| 相对直接普通 baseline 加速 | 约 2.84 倍 |

这是目前精度最高的一组。但是它只比第一名高约 0.02 bits，安全裕量却明显变小。
因此它适合作为对照组，但不如第一名适合作为默认主推参数。

### 第三名：`top=16 -> leaf=15`，低 Q 调参版本

命令：

```bash
go run ./5-7/. \
  -run=ours \
  -leafEngine=anyu \
  -anyuProfile=lowq \
  -anyuSubringLogN=12 \
  -anyuEvalModBits=56 \
  -anyuStCBits=26 \
  -anyuCtSBits=50 \
  -logN=16 \
  -layers=1 \
  -splitFirstLeafPreprocess \
  -reps=1
```

| 指标 | 数值 |
|---|---:|
| 大环 | `logN=16` |
| 叶子环 | `logN=15` |
| 叶子安全参考 | \(\log(PQ)\approx 868\) |
| 实际 logQ | 866 |
| 实际 logQP | 927 |
| 安全裕量 | 约 -59 bits，如果完整计入 leaf15 的 \(QP\) |
| 输出计算层 | 1 |
| 平均 L1 精度 | 16.74 bits |
| 本方案时间 | 6.59s |
| 只有拆分时间 | split-first 失败；旧路线 9.39s 但精度 1.51 bits |
| 只有 subKey 时间 | 32.88s |
| 直接普通 baseline 时间 | 60.85s |
| 相对只有 subKey 加速 | 约 4.99 倍 |
| 相对直接普通 baseline 加速 | 约 9.23 倍 |

这是 `top=16` 目前最好的结果。

但是它的问题也比较明显：

1. 如果完整计入 \(QP\)，`logQP≈927` 超过 leaf15 的 868 参考。
2. 输出计算层只有 1 层。
3. 精度低于 `top=17 -> leaf=16` 的 I1 配置。

所以它适合作为说明“融合路线可行”的工程实验，不适合作为最终主推安全参数。

## 9. 总结

目前最有价值的结论是：

\[
\boxed{
\texttt{top logN=17, layers=1, leaf logN=16, subKey Bootstrapping I1}
}
\]

这组参数的综合效果最好：

1. 叶子安全参数基本合理：\(\log(PQ)\approx 1730 < 1747\)。
2. 叶子正好是 \(2^{16}\to2^{12}\) 的 subKey Bootstrapping 设置。
3. 平均 L1 精度约 19.40 bits。
4. 输出后保留 15 个计算层。
5. 相比只有 subKey 的大环 Bootstrapping，本地实验约 2.6 倍加速。
6. 相比直接普通大环 Bootstrapping，本地实验约 2.75 倍加速。
7. 相比只有拆分的普通叶子 Bootstrapping，本地实验约 1.21 倍加速。

下一步应该做：

1. 对第一名参数重复多次实验，统计平均值和方差。
2. 跑 `-stageTrace`，确认剩余约 2 bits 精度损失主要来自哪个阶段。
3. 增加内存和 evaluation key size 统计。
4. 用同一套机器设置重新跑 direct baseline 和 proposed path，避免单次运行波动。

## 10. 总览大表

下面这张表把目前主要成功和失败实验放在一起，便于横向比较。

说明：

- “组合方案”表示 no-\(B_k\) 拆分 + 叶子 subKey Bootstrapping。
- “只有拆分”表示 no-\(B_k\) 拆分 + 普通叶子 Bootstrapping。
- “只有 subKey”表示不拆分，直接在大环上使用 subKey Bootstrapping。
- “直接普通”表示不拆分，也不使用 subKey。
- “安全裕量”按本报告使用的 128-bit 参考表粗略计算。
- `失败` 表示该路线没有得到有意义的正确结果。

| 编号 | 大类 | top logN | leaf logN | 方法 | 主要参数 | 安全参考 | logQ | logQP | 安全裕量 | 输出计算层 | 平均 L1 精度 | 时间 | 状态与结论 |
|---:|---|---:|---:|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | top17 主推 | 17 | 16 | 组合方案 | `i1` | 1747 | 1425 | 1730 | +17 | 15 | 19.40 bits | 52.03s | 当前综合最好 |
| 2 | top17 高 CtS | 17 | 16 | 组合方案 | `i1, CtS=55` | 1747 | 1434 | 1739 | +8 | 15 | 19.42 bits | 50.51s | 精度最高，但安全裕量小 |
| 3 | top17 高 degree | 17 | 16 | 组合方案 | `i1, Mod1Degree=255` | 1747 | 1415 | 1720 | +27 | 13 | 19.39 bits | 45.33s | 没提升精度，还少 2 层 |
| 4 | top17 低 Q 调参 | 17 | 16 | 组合方案 | `lowq, EvalMod=56, StC=26, CtS=50` | 1747 | 866 | 927 | +820 | 1 | 15.47 bits | 45.66s | 安全裕量大，但精度和层数不够 |
| 5 | top17 低 Q 默认 | 17 | 16 | 组合方案 | `lowq` | 1747 | 868 | 929 | +818 | 1 | 14.64 bits | 42.97s | 不如低 Q 调参 |
| 6 | top17 only split | 17 | 16 | 只有拆分 | 普通叶子 BTS，15 层 | 1747 | 1391 | 1696 | +51 | 15 | 20.47 bits | 63.19s | 精度高于组合方案，但更慢 |
| 7 | top17 only subKey | 17 | - | 只有 subKey | `i1` | 3523 | 1425 | 1730 | +1793 | 15 | 18.61 bits | 135.08s | 不拆分时很慢 |
| 8 | top17 only subKey 高 CtS | 17 | - | 只有 subKey | `i1, CtS=55` | 3523 | 1434 | 1739 | +1784 | 15 | 18.60 bits | 135.15s | CtS=55 没改善 only subKey |
| 9 | top17 direct | 17 | - | 直接普通 | 普通大环 BTS，15 层 | 3523 | 1391 | 1696 | +1827 | 15 | 20.87 bits | 143.29s | 精度最高，但最慢 |
| 10 | top16 最好 | 16 | 15 | 组合方案 | `lowq, EvalMod=56, StC=26, CtS=50` | 868 | 866 | 927 | -59 | 1 | 16.74 bits | 6.59s | 工程上能跑，但严格 QP 安全不干净 |
| 11 | top16 低 Q 默认 | 16 | 15 | 组合方案 | `lowq` | 868 | 868 | 929 | -61 | 1 | 16.16 bits | 6.73s | 能跑通，但不如调参版 |
| 12 | top16 Eval56/CtS46 | 16 | 15 | 组合方案 | `EvalMod=56, CtS=46` | 868 | 866 | 927 | -59 | 1 | 16.42 bits | 6.55s | 略好于默认，但不如 StC26/CtS50 |
| 13 | top16 Eval57/CtS43 | 16 | 15 | 组合方案 | `EvalMod=57, CtS=43` | 868 | 867 | 928 | -60 | 1 | 13.77 bits | 6.69s | CtS 太低，失败倾向 |
| 14 | top16 Eval54/CtS52 | 16 | 15 | 组合方案 | `EvalMod=54, CtS=52` | 868 | 864 | 925 | -57 | 1 | 15.12 bits | 6.79s | EvalMod 不够 |
| 15 | top16 Eval60/CtS33 | 16 | 15 | 组合方案 | `EvalMod=60, CtS=33` | 868 | 867 | 928 | -60 | 1 | 3.77 bits | 6.62s | CtS 过低，明显失败 |
| 16 | top16 StC23 | 16 | 15 | 组合方案 | `EvalMod=57, StC=23, CtS=50` | 868 | 867 | 928 | -60 | 1 | 14.02 bits | 6.02s | StC 太低 |
| 17 | top16 subringLogN13 | 16 | 15 | 组合方案 | `subringLogN=13` | 868 | 868 | 929 | -61 | 1 | -5.52 bits | 8.03s | 失败 |
| 18 | top16 only split split-first | 16 | 15 | 只有拆分 | 普通叶子 BTS，split-first | 868 | 899 | 960 | -92 | - | 失败 | - | leaf ScaleDown 条件不满足 |
| 19 | top16 only split 旧路线 | 16 | 15 | 只有拆分 | 普通叶子 BTS，旧路线 | 868 | 899 | 960 | -92 | 1 | 1.51 bits | 9.39s | 输出约放大 2 倍，语义归一化失败 |
| 20 | top16 only subKey | 16 | - | 只有 subKey | `lowq, EvalMod=56, StC=26, CtS=50` | 1747 | 866 | 927 | +820 | 1 | 15.92 bits | 32.88s | 比组合方案慢很多 |
| 21 | top16 direct | 16 | - | 直接普通 | 普通大环 BTS，1 层 | 1747 | 899 | 960 | +787 | 1 | 21.97 bits | 60.85s | 精度高，但慢且不是 leaf 安全对照 |

从这张总表看，当前最适合主推的仍然是编号 1：

\[
\boxed{
\texttt{top logN=17, leaf logN=16, i1, split + subKey}
}
\]

它不是精度最高的单点，但在安全性、精度、输出层数和加速比之间最均衡。
