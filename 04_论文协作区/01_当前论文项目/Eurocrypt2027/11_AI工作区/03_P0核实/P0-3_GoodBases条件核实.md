# P0-3: Good Bases 条件核实 (r ∈ {1,2} 是否满足)

> 基于: P0-2 的计算 + GHPS 2013 Lemma 2.6 (src_000028)

---

## 一、计算结论

由于当前论文使用 **power-of-two cyclotomic** (N, n 均为 2 的幂):

```
m = 2N,  m₀ = 2n
rad(m) = 2,  rad(m₀) = 2
r = rad(m) / rad(m₀) = 1
```

**r = 1 ∈ {1,2} → 条件满足。**

## 二、GHPS Good Bases Lemma 2.6 的含义

GHPS Lemma 2.6 (Section 2.1.5):
- 当 r = 1 或 r = 2 时，可以构造 good basis B，其 Gram matrix 是 **scaled unitary**
- 这意味着: 对偶基 B^∨ 的 embedding norms 与 B 的关系清晰
- s₁(B) = √(m̂/m₀) 可以直接使用（噪声各向同性）
- canonical embedding 范数 ∥a∥^{can}₂ 和 ∥a∥∞ 之间有良好关系

**对当前论文的影响**:
- GHPS field switching 的噪声分析可以直接引用
- 不需要额外处理 anisotropic noise
- 但当前论文是否需要引用 GHPS Good Bases **取决于 paper 是否做 field-switching noise proof**

## 三、当前论文是否需要引用 GHPS Good Bases？

| 如果 | 则需要 |
|------|--------|
| 论文做 ring switching noise analysis | 可能需要引用 |
| 论文只是把 ring switching 作为工具 | 可以只引用不展开 |
| 论文不包含 field switching proof | Good Bases 可能不需要显式出现 |

## 四、如果允许非 2-power 参数

如果后续论文考虑 non-power-of-2 m/m₀:
- r 可能 > 2 → scaled unitary 性质不成立
- 噪声分析需做 adaptation
- 建议: **论文中明确声明 power-of-two 假设**

## 五、结论

| 问题 | 答案 |
|------|------|
| r = ? | **r = 1**（N, n 都是 powers of 2） |
| r ∈ {1,2} 是否满足？ | ✅ 满足 |
| Good Bases scaled unitary? | ✅ 是（r=1） |
| 噪声是否各向同性？ | ✅ 是 |
| 是否可直接引用？ | ✅ 可以 |
| 当前论文是否需要？ | ⚠ 取决于论文是否做 ring switching noise proof |
| 状态 | ✅ 已核实 (high confidence) |
