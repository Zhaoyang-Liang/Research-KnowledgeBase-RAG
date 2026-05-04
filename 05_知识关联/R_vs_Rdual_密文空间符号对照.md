# R 与 R^∨：BGV / GHPS Field Switching 密文空间符号对照

> 版本: v1 | 生成: 2026-05-04 | 基于: GHPS 2013 (src_000028) + Geelen 2023/2024 (src_000030)

---

## 一、总结判断

### 核心结论

`R` 和 `R^∨` **不是同一个数学对象**。

- `R` 是整环（ring of integers），`R^∨` 是其 codifferent（一个分式理想，通常是 R 的真上集）
- GHPS Field Switching **必须使用** `R^∨`，因为它的 correctness proof 和 security proof 建立在对偶结构和 trace pairing 之上
- Geelen BGV/BFV SlotToCoeff 使用 `R` 是因为它聚焦算法/实现层面，不涉及 RLWE 安全证明中的 dual/codifferent 讨论
- 在实现中（固定 basis 后），二者都可以表示为系数向量，**但这只是表示层面的巧合，不是数学等同**

### 当前论文的建议

**主文算法层**：采用 `R_q^2`（BGV/BFV 实现记号），与 Geelen、HElib、SEAL 等保持一致。

**涉及 GHPS Field Switching 之处**：切换到 GHPS 双标记号，明确写在 `(R_q^∨)^2` 中，
并在 Preliminaries 中增加 "Dual/co-different convention" 小节说明两套记号如何联系。

**关键转换**：乘以（或除以）scaling element `t` 在 `R` 和 `R^∨` 之间转换。
GHPS 明确指出 `R^∨ = t^{-1}R`，因此 `a ∈ R^∨ ⟹ t·a ∈ R`。

---

## 二、两篇论文符号对照表

| 对象 | GHPS Field Switching (2013) | Geelen SlotToCoeff BGV/BFV (2023) | 当前论文建议 |
|------|---------------------------|----------------------------------|-------------|
| cyclotomic ring | `R = Z[X]/(Φ_m(X))` | `R = Z[X]/(X^N+1)` (只含 2-power) | 看参数：2-power 用 X^N+1，一般用 Φ_m(X) |
| 整数环 (ring of integers) | `R` = ring of integers of `K` | `R` = `Z[X]/(X^N+1)` | `R` (含义取决于上下文) |
| 分圆数域 | `K = Q(ζ_m)` | **未显式定义** | `K = Q[X]/(Φ_m(X))` |
| 子域/子环 | `K₀ = Q(ζ_{m₀})`, `R₀ = Z[ζ_{m₀}]` | `R' = Z[X^{N/N'}]/(X^N+1)` | `R₀` (子环) |
| **codifferent (dual)** | `R^∨ = {a ∈ K: Tr_{K/Q}(aR) ⊆ Z}` | **未定义，完全不用** | `R^∨` 仅在引用 GHPS 时使用 |
| **密文空间** | `(R_q^∨)^2` | `R_q^2` ("a ciphertext is a pair of ring elements, i.e. it lives in R_q^2") | 主文 `R_q^2`，GHPS 部分 `(R_q^∨)^2` |
| **明文空间** | `R_p` = `R/pR` | `R_{p^e}` | `R_{p^e}`（如上） |
| **噪声空间** | `e ∈ R^∨` | `e ∈ R` (满足 `‖e‖_∞ < (q/p^e − 1)/2`) | 主文 `e ∈ R`，GHPS 部分 `e ∈ R^∨` |
| **解密关系** | `c₀ + c₁·s = e mod qR^∨` | BGV: `c₀ + c₁·s = m + p^e·e (mod q)` | 看上下文 |
| **trace target** | `Tr_{K/K₀}(R^∨) = R₀^∨` | **不涉及** (Geelen 不做 field switching) | 仅在 field switching 中使用 |
| **subring / embedding** | `K₀ → K`（数域扩张） | `R' = Z[Y]/(Y^{N'}+1)` 其中 `Y = X^{N/N'}` | `R₀ → R` |
| **slot algebra** | 不显式讨论（使用 BGV CRT slot） | `E = Z_{p^e}[ζ_m]`, slots = CRT factors | `CRT slot`（BGV/BFV） |
| **automorphism** | `τ_i: ζ_m ↦ ζ_m^i` | `τ_j: a(X) ↦ a(X^j)` | `τ_j` 或 `κ_j` |
| **decoding basis** | 显式定义 decoding basis for `R_q^∨` | **不定义**（系数向量默认表示 R 中的元素） | 实现层省略，证明层提及 |
| **实现表示** | `R^∨` 元素通过 fixed basis 表示为 Z-系数向量 | `R` 元素天然是系数不超过 N-1 的多项式 | 主文用多项式系数表示 |
| **RLWE 安全性** | RLWE over `K` with secret in `R₀^∨` | **不讨论**（假设安全性已知） | 引用 GHPS |

---

## 三、关键差异的数学来源

### 为什么 GHPS 必须用 `R^∨`？

1. **Trace 的代数性质**：`Tr_{K/K₀}(R)` **不等于** `R₀`（而是 R₀ 的某个真理想）。
   如果密文写在 `R` 上，trace 后的结果不在 `R₀` 中 → correctness proof 断裂。

2. **对偶解决**：`Tr_{K/K₀}(R^∨) = R₀^∨`（干净）。Field switching 的核心操作是
   对小域元素 `r ∈ R₀` 构造 trace pairing 函数:
   `L_r(a^∨) = Tr_{K/K₀}(r · a^∨)`，其中 `a^∨ ∈ R^∨`，结果在 `R₀^∨`。

3. **Good Bases**：范数估计依赖 `R` 和 `R^∨` 作为对偶模块的几何性质。
   Lemma 2.6 构造的 "good" 基同时适用于 `R` 和 `R^∨`。

4. **Security Proof**：GHPS 需要证明 "RLWE over K with secret in R₀ 的安全性
   归约到 RLWE over K₀"。这个证明依赖对偶结构和 good bases。

### 为什么 Geelen（实现论文）不需要 `R^∨`？

1. **不涉及 RLWE 安全证明**：Geelen 论文假设 BGV/BFV 的安全性已知，不重新证明。

2. **不涉及 field switching 的代数处理**：算法在固定的大环 `R` 上做 slot-to-coeff
   变换，不需要 trace 到子环。

3. **实现友好**：`R = Z[X]/(X^N+1)` 是 power-of-two 环，
   元素天然表示为长度 N 的多项式，Haskell/SEAL/HElib 都这样处理。

4. **BGV 解密公式中 `e ∈ R` 是简化写法**：
   严格 RLWE 表述中噪声在 `R^∨`，但在实现层面可以通过固定基
   （decoding basis）将二者等同。Geelen 选择了算法层面的简化。

---

## 四、风险点（如果混用）

| 风险 | 说明 |
|------|------|
| correctness proof 不严谨 | trace 目标空间写错，decryption relation 写错 |
| noise bound 写错 | GHPS Good Bases 的范数不等式适用于 `R^∨`/`R₀^∨`，直接套用在 `R` 上可能差一个 scaling factor |
| trace 的目标空间写错 | `Tr(R^∨) = R₀^∨` ≠ `Tr(R) = (不是 R₀)` |
| Good Bases 作用解释不清 | 读者不理解为什么需要 "good bases"，以为可以直接用多项式系数 |
| plaintext/noise/ciphertext 空间混淆 | reader 无法判断你采用的是 GHPS 记号还是实现记号 |
| RLWE security embedding 不准确 | Security proof 中的对偶结构被省略 |

---

## 五、论文写作建议措辞

### 建议在 Preliminaries 中新增 "Dual/codifferent convention"

```text
本文采用以下符号约定:
- 主文算法描述使用 R_q^2 记号，与 BGV/BFV 实现文献 (Geelen, HElib, SEAL) 一致。
- 在引用 GHPS field switching 的 correctness 和 security proof 时，
  切换到 GHPS 的 dual notation: 密文空间写作 (R_q^∨)^2，噪声空间写作 R^∨。
- 二者的关系: R^∨ = t^{-1}R 是主分式理想，乘以 t 可在 R^∨ ↔ R 之间转换。
  在实现中，通过固定 decoding basis，R^∨ 的元素可表示为系数向量（与 R 相同）。
```

### 转换示意图

```
GHPS notation (proof layer):
  密文 ∈ (R_q^∨)^2    噪声 ∈ R^∨    trace: R^∨ → R₀^∨

             ↓ 固定 decoding basis ↓

Implementation notation (algorithm layer):
  密文 ∈ R_q^2        噪声 ∈ R      操作同 Geelen/SEAL
```

---

## 六、来源引用

| 引用 | 来源 | source_id |
|------|------|-----------|
| `R^∨` 定义 | GHPS 2013, 2.1.4 | src_000028 |
| `Tr(R^∨)=R₀^∨` | GHPS 2013, Eq (2.3) | src_000028 |
| `R^∨ = t^{-1}R` | GHPS 2013, 2.1.4 | src_000028 |
| Good Bases Lemma 2.6 | GHPS 2013, 2.1.5 | src_000028 |
| `R = Z[X]/(X^N+1)`, `R_k = R/kR` | Geelen 2023, 2.1 | src_000030 |
| ciphertext in `R_q^2` | Geelen 2023, 2.2 | src_000030 |
| BGV decryption: `c₀+c₁·s = m+p^e·e` where `e ∈ R` | Geelen 2023, 2.2 | src_000030 |
