# P0-4: R vs R^∨ 符号核实

> 基于: 初稿.pdf (src_000031) + GHPS 2013 (src_000028) + R_vs_Rdual_密文空间符号对照.md

---

## 一、问题重新表述

原问题: "论文密文空间到底 R 还是 R^∨？" → **不准确的二选一**

正确表述:
```
当前论文需要统一实现记号 R_{q,N}^2 与 GHPS dual/codifferent 记号 (R_q^∨)^2。
```

## 二、初稿采用的记号

初稿 Section 2.1:
```
R_{q,N} = Z_q[X]/(X^N + 1)
ct = (c_0, c_1) ∈ R_{q,N}^2
```

**这是纯实现层记号，完全使用 R，不涉及 R^∨。**

## 三、与 GHPS 的差异

| 维度 | 初稿 | GHPS |
|------|------|------|
| 环记号 | R_{q,N} = Z_q[X]/(X^N+1) | R_q = Z_q[ζ_m] |
| 密文空间 | ct ∈ R_{q,N}^2 | ct ∈ (R_q^∨)^2 |
| 对偶记号 | 不定义 | R^∨ = {a: Tr(aR) ⊆ Z} |
| 说明 | 实现友好 | RLWE 对偶结构 |

## 四、判断

| 场景 | 推荐记号 |
|------|----------|
| 主文算法/实现章节 | R_{q,N}^2 ← **初稿已经这样做** ✅ |
| Section 5.5 安全证明 | 如果涉及 GHPS 归约 → 需要说明 R^∨ |
| 引用 GHPS ring switching 的 trace 性质 | 需要切换到 GHPS 记号并说明 |
| Preliminaries | 增加 "Notation for dual" 小节 |

## 五、当前论文是否需要补充 codifferent？

取决于:
1. **如果不做 GHPS security proof** → 不需要（只需 notation caveat）
2. **如果做 GHPS security proof** → 必需

建议: 即使不完整重做 GHPS 安全性，也应在 Preliminaries 中加一句:
```text
Remark (Dual/codifferent convention).
GHPS field switching formulates ciphertexts in (R_q^∨)^2,
where R^∨ = t⁻¹R is the codifferent. In our algorithmic layers,
we use R_{q,N}^2 for implementation clarity.
The conversion is via a fixed decoding basis.
```

## 六、结论

| 问题 | 答案 |
|------|------|
| 初稿使用什么记号？ | R_{q,N}^2（实现层） |
| 是否正确？ | ✅ 正确 |
| 是否需要补充 R^∨？ | ⚠ 是，在 Preliminaries 中加 notation caveat |
| 是否需要完整展开 dual？ | 取决于是否做 GHPS security proof |
| 状态 | ⚠ 符号体系冲突，非方案错误；推荐 Notation caveat |
