# K029: R 与 R^v 在 BGV/GHPS Field Switching 中的符号差异

## 摘要

R = Z[X]/(X^N+1) 是实现层记号（Geelen/SEAL/HElib）；
R^v = {a in K: Tr(aR) subset of Z} = t^{-1}R 是 GHPS 的 dual/codifferent 记号。

## 关键区别

| 维度 | GHPS (R^v) | Geelen/实现 (R) |
|------|-----------|----------------|
| 定义 | R^v = {a: Tr(aR) subset of Z} | R = Z[X]/(X^N+1) |
| 密文空间 | (R_q^v)^2 | R_q^2 |
| 噪声空间 | e in R^v | e in R |
| 为什么这样写 | trace 保 duality (Tr(R^v)=R0^v) | 实现友好 |
| 涉及 dual? | 是 | 否 |

## 当前论文建议

- 主文算法: R_q^2
- GHPS 引用: (R_q^v)^2
- Preliminaries: 增加 "Dual/codifferent convention" 小节

## 状态

- needs_human_review: True
- confidence: high（原文定义清晰）
- claim_source_type: Qclaw推断（跨源合成）
