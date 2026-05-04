# GHPS Duality / Good Bases 精读报告

> 生成时间: 2026-05-04 | 来源: GHPS 2013 论文 2.1.4-2.1.5 | source_id: src_000028

## 1. 什么是 R^∨ (R dual / codifferent)

GHPS 定义：

> R^∨ = {a ∈ K : Tr_{K/Q}(aR) ⊆ Z} ⊇ R

即 R^∨ 是 K 中所有满足"与 R 的迹落在 Z 中"的元素集合。
它是 R 的**上集**（codifferent，分数理想）。

关键性质：
- R^∨ = t^{-1}R，是**主分数理想**（principal fractional ideal）
- 除以 t: R → R^∨ 是一个 R-模同构
- Tr_{K/K₀}(R^∨) = R₀^∨（公式 2.3）

等价定义：R^∨ = {a ∈ K : Tr_{K/K₀}(aR) ⊆ R₀^∨}

### 为什么不用 R 而用 R^∨？

GHPS 指出：Tr_{K/K₀}(R) **不等于** R₀（而是 R₀ 的某个真理想），
所以直接用 R 做 trace 不方便。用 R^∨ 代替 R：
- Tr 映射 R^∨ → R₀^∨（干净的满射）
- 在 field switching 中用 t 的除法作为 R ↔ R^∨ 之间的转换

## 2. Duality 在 GHPS 中服务于什么？

**服务于 field switching 的正确性**。

Field switching 需要：
1. 将 R_p ≅ R/p 上的密文映射到 R_{p₀} 上
2. 核心操作：同态 trace Tr_{K/K₀}(·) 将大环密文投影到子环
3. 但 trace 直接作用在 R 上时，像不一定是 R₀——这是一个代数障碍

**Duality 的解决方案**：
- 先用 t 将密文从 R 转换到 R^∨
- 在 R^∨ 上做 trace（干净）
- 再用 t₀ 从 R₀^∨ 转换回 R₀
- 这就是 GHPS 的 "decode basis" 机制

### trace pairing 形式

对任意 R₀-线性函数 L: R → R₀，存在 r^∨ ∈ t₀R^∨ 使得：
L(a) = Tr_{K/K₀}(r^∨·a)

在 duality 框架下，GHPS 转换为在 R^∨ 上工作：
L^∨(a^∨) = Tr_{K/K₀}(r·a^∨), 其中 r = (t/t₀)·r^∨ ∈ R

**关键**：这让 trace pairing 在不含分母的整环 R 上就能表示。

## 3. Good Bases 的定义

> 我们寻求 R 的一个 R₀-基 ~b，使得 B（典范嵌入矩阵）几乎保持欧几里得范数，
> 即所有奇异值几乎相等。

形式化：存在 ~b = (b_j) ∈ R^{n/n₀} 满足：
- 任何 K₀-系数的短向量（在 σ₀ 下短）对应 K 中的短元素（在 σ 下短）
- 典范嵌入矩阵 B: σ(~b) = B·σ₀(~a)，其中 B 的奇异值接近 √(m̂/m₀)

### Lemma 2.6

GHPS 构造了一个 R₀-基 ~b，满足：
- s₁(B) = √(m̂/m₀)
- sₙ(B) = √(m/(r·m₀))
- 当 r ∈ {1, 2} 时，B 是 unitary 矩阵乘以 √(m̂/m₀) 因子

其中：m̂ = m/2（如果 m 偶且 m₀ 奇）/ m（其他），r = rad(m)/rad(m₀)。

### 对偶基

对偶基 ~b^∨ 定义为 Tr_{K/K₀}(b_j·b_{j'}^∨) = δ_{jj'}。
对应矩阵 B^∨ = B^{-T}，奇异值是 B 的逆。

## 4. Good Bases 为什么影响 field switching？

**直接影响 field switching 的正确性和安全性**：

1. **正确性**：field switching 中的密文是从 R 的元素构造的。
   Good bases 保证了操作不引入过多噪声——因为短系数 ↔ 短元素，
   不会因基变换导致噪声爆炸。

2. **安全性**：Good bases 用于证明 "RLWE over K with secret in R₀ 的安全性
   归约到 RLWE over K₀ with secret in R₀"（Section 3.1）。
   这是 field switching 安全性证明的基础。

3. **噪声控制**：公式 (2.5) 量化了 K 和 K₀ 之间的范数关系：
   ‖σ(a)‖ ≤ √(m̂/m₀)·‖σ₀(~a)‖
   这是 field switching 噪声分析的关键不等式。

## 5. Trace pairing 与 Good Bases 的关系

Trace pairing 是 Good Bases 定义的**基础**：
- 对偶基通过 trace 关系定义
- B^∨ = B^{-T} 直接来源于 trace pairing
- 对偶基的奇异值是 B 的逆，用于分析 R^∨ 上的范数

简言之：trace pairing 定义了对偶结构，Good Bases 让这个结构在几何上"良好"。

## 6. 是否影响当前论文的 correctness proof？

**很有可能影响**：
- 如果当前论文使用了 field switching（大概率是），则 GHPS 的 correctness
  建立在 duality + good bases 之上
- 噪声估计公式直接依赖 good bases 的几何性质
- 如果当前论文使用了不同的基或不同的环结构，需重新证明

## 7. 是否影响噪声分析？

**直接影响**：
- 公式 (2.5) 是噪声分析的核心：K 元素的范数 vs K₀-系数向量的范数
- 每次 trace 操作（Tr_{K→K₀}）的噪声放大因子由 good bases 的奇异值决定
- s₁(B) = √(m̂/m₀) 是主要的噪声放大因子
- 在 r ∉ {1,2} 时，sₙ(B) ≠ s₁(B) 会导致各向异性噪声

## 8. 仍需人工确认的内容

| # | 问题 | 原因 |
|----|------|------|
| 1 | 当前论文使用的分圆环参数（m, m₀, r 值）是什么？ | 决定 good bases 的具体几何量 |
| 2 | r ∈ {1,2} 条件在论文参数中是否满足？ | 若满足则噪声各向同性（简化分析） |
| 3 | 论文的 noise analysis 是否引用了 good bases 的范数不等式？ | 若无引用需补充 |
| 4 | decode basis 具体是哪个 Z-基？ | 影响 field switching 的实现细节 |
| 5 | 当前论文是否在 R 还是 R^∨ 上定义密文？ | GHPS 的 duality 转换是这个选择的关键 |

## 建议后续

1. 🔴 确认当前论文的环参数 (m, N, K, K₀)，计算 r 值
2. 🟡 精读 GHPS Section 3.1 (security proof) 和 3.2-3.3 (correctness)
3. 🟡 对比论文 noise analysis 与 GHPS 的 bound

## 来源

- GHPS 2013, 2.1.4 Duality, 2.1.5 Good Bases (src_000028)
- K004 (Trace), K006 (Prime Splitting), K021 (Trace操作)
