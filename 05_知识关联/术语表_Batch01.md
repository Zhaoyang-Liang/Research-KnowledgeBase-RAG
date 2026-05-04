# Batch-01 修正版术语表

> 生成时间: 2026-05-04 | 版本: v2 | 术语数: 35

## 基础代数

| 术语 | 英文 | 定义 |
|------|------|------|
| 分式域 | field of fractions | 整环中形式分数等价类构成的域 |
| 分裂域 | splitting field | 多项式完全分解为一次因子所需的最小扩域 |
| 整环 | integral domain | 无零因子的交换环 |
| 共轭 | conjugate | 同一最小多项式的不同根 |
| 自同构 | automorphism | 域到自身的保持代数结构的双射 |
| Galois群 | Galois group | 扩张K/F的全体F-自同构组成的群 |
| Galois扩张 | Galois extension | 可分且正规的代数扩张 |
| 复嵌入 | complex embedding | 数域到C的域同态 |
| 典范嵌入 | canonical embedding | 使用全部复嵌入将域元素映射到C^N |
| 迹 | trace | 扩张中元素在所有F-自同构下像的和 |
| 范数 | norm | 扩张中元素在所有F-自同构下像的积 |
| 固定域 | fixed field | 被自同构子群逐点固定的最大子域 |
| 代值同态 | evaluation homomorphism | 将多项式代入特定元素得到的环同态 |
| 第一同构定理 | first isomorphism theorem | 若φ:A→B环同态则A/ker(φ)≅im(φ) |

## 分圆域与有限域

| 术语 | 英文 | 定义 |
|------|------|------|
| 分圆多项式 | cyclotomic polynomial | 以所有m次本原单位根为根的首一多项式 |
| ord_m(p) | multiplicative order | 使p^k≡1(mod m)的最小正整数k |
| 完全分裂 | split completely | 多项式在域上分解为一次因子的乘积 |
| Frobenius自同构 | Frobenius automorphism | 特征p域上的映射x→x^p |
| Frobenius轨道 | Frobenius orbit | Frobenius映射下元素的轨道 |
| NTT | Number Theoretic Transform | 有限域F_q上的DFT |
| 蝴蝶操作 | butterfly operation | Cooley-Tukey算法基本单元 |
| 位反转 | bit-reversal | 按二进制位反转索引的置换 |
| Eisenstein判别法 | Eisenstein criterion | Q上多项式不可约的充分条件 |
| SIMD | Single Instruction Multiple Data | 将多项式在各槽求值解释为独立数据 |

## FHE方案与操作

| 术语 | 英文 | 定义 |
|------|------|------|
| 模数切换 (BGV) | Modulus Switching | BGV中缩放密文至更小模数降噪，消息不变 |
| 重缩放 (BFV) | Rescale (BFV) | BFV中÷t恢复消息尺度 |
| 重缩放 (CKKS) | Rescale (CKKS) | CKKS中÷Δ降低模数层级 |
| 尺不变 | scale-invariant | BFV特征：消息表示在计算中保持尺度一致 |
| 密钥切换 | Key Switching | 密文从一个密钥同态转换到另一个密钥 |
| 数字分解 | digit decomposition | 大整数表示为基B的多位数字 |
| 重线性化 | relinearization | 密钥切换特例(s²→s)消除乘法维度膨胀 |
| 模数链 | modulus chain | 从q₀到q_L的有序模数序列 |
| 噪声预算 | noise budget | 当前密文剩余噪声容忍度(单位bits) |

## 高级FHE技术

| 术语 | 英文 | 定义 |
|------|------|------|
| 盲旋转 | blind rotation | 用同态旋转和test vector实现的不暴露密钥的解密 |
| 多项式打包 | polynomial packing | 将多个解密块结果合并为一个多项式 |
| 结构化密钥 | structured secret key | 将密钥拆分为稀疏块的技巧 |
| 同态迹 | homomorphic trace | 通过旋转和加法同态计算的迹操作 |
| 环嵌入 | ring embedding | 一个多项式环到另一个多项式环的单同态 |
| 系数拆分 | coefficient split | 按多项式系数奇偶位置拆分为两个多项式 |
| 槽保持拆分 | slot-preserving split | 拆分后每个分量仍保留原槽语义 |
| SlotToCoeff | SlotToCoeff | 将系数表示转换为slot表示 |
| CoeffToSlot | CoeffToSlot | SlotToCoeff的逆：slot表示→系数表示 |
| U_n矩阵 | U_n matrix | 典范嵌入的Vandermonde矩阵表示 |
| Field switching | field switching | GHPS技术：通过迹将密文从大分圆环转换到小子环 |

## 槽类型（⚠ 重要区分）

| 术语 | 英文 | 定义 |
|------|------|------|
| CRT槽 | CRT slot | BGV/BFV中通过CRT分解得到的槽，值为F_{p^d}元素 |
| CKKS槽 | CKKS slot / complex slot | CKKS中通过复嵌入得到的槽，值为复数 |
| CKKS slot vs BGV/BFV CRT slot | - | 数学来源不同：复嵌入 vs CRT分解。论文写作必须区分！ |

## 密钥材料

| 术语 | 英文 | 定义 |
|------|------|------|
| 切换密钥 | switching key (KSK) | 预计算密钥材料{(b_i,a_i)}用于密钥切换 |
| 重线性化密钥 | relinearization key (RLK) | KSK_{s²→s}特例 |
| 旋转密钥 | rotation key / automorphism key | KSK_{s(X)→s(X^k)}用于同态旋转 |
| Galois密钥 | Galois key | 旋转密钥的另一种称呼 |

## 来源

- 所有术语定义基于 Batch-01 已读 27 个源文件
- 标注 ⚠ 的区分项需要人工确认（详见 `需要人工确认清单.md`）
