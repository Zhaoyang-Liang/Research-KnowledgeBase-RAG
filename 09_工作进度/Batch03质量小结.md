# Batch-03 质量小结

> 2026-05-04

---

## 1. Batch-03 是否通过

✅ **通过**。9个文件全部正确导入、阅读、处理。7张卡片+13条Q&A+8术语+10公式+4候选知识全部产出。

---

## 2. 哪些卡片可以进入正式知识库

| card_id | 发布目标 | 状态 |
|---------|----------|------|
| K030 | 03_FHE知识库 | support_ingest (translation_note) |
| K031 | 03_FHE知识库 | support_ingest (translation_note) |
| K034 | 03_FHE知识库 | support_ingest (translation_note) — Trace/内积是标准知识 |
| K035 | 03_FHE知识库 | support_ingest (needs_verification for title) |
| K036 | 03_FHE知识库 | support_ingest (translation_note) |

> 注：所有卡片均为 translation_note，不可直接作为 validated knowledge。需核原文后降级 needs_human_review 标记。

---

## 3. 哪些卡片只作为 translation_note support

全部 7 张（K030-K036）均为 translation_note。
其中：
- K034 (Trace/内积)、K036 (τ_n正向) — 教科书级内容，confidence 较高，review 需求较低
- K030 (τ_n^{-1})、K031 (密文态SlotToCoeff) — 中等复杂度，需对照原文
- K032 (C^B嵌入)、K033 (PaCo结构)、K035 (5的幂) — 高优先级核原文

---

## 4. 哪些卡片需要核原文

| 优先级 | card_id | 核原文事项 |
|--------|---------|-----------|
| 🔴 高 | K032 | 原文求和范围笔误(2M vs 2B) — 必须作者确认 |
| 🔴 高 | K033 | 全部索引定义需对照PaCo原文验证 |
| 🟡 中 | K030 | Vandermonde矩阵U_n的具体形式 |
| 🟡 中 | K031 | SlotToCoeff密文态语义描述 |
| 🟡 中 | K035 | 标题"5作为单位根生成元"不精确，应改为"半本原根生成元" |
| 🟢 低 | K034 | Trace/内积是标准知识 |
| 🟢 低 | K036 | τ_n正向是标准推导 |

---

## 5. 哪些内容进入论文协作区

- `04_论文协作区/.../00_项目管理/作者代办路线_完整版.md` — 汇总代办路线
- `04_论文协作区/.../00_项目管理/论文写作方向_通用框架版.md` — 可能的尝试写法

---

## 6. 哪些内容进入知识入库接口

4 条候选知识 (cand_001~004):

| id | 内容 | 类型 |
|----|------|------|
| cand_001 | split-compatible sandwich-circuit criterion | construction_idea |
| cand_002 | CKKS bootstrapping core is non-trivial instance | proof_trick |
| cand_003 | γ_impl as implementation normalization factor | notation_convention |
| cand_004 | leaf security 不继承 top security | parameter_lesson |

全部标记 needs_validation + needs_human_review。

---

## 7. 与 Batch-01 的重复/交叉引用情况

| Batch-03 卡片 | Batch-01 相关卡片 | 关系 |
|--------------|------------------|------|
| K030 (τ_n^{-1}) | K018 (典范嵌入), K019 (SlotToCoeff正向) | 互补 — τ_n^{-1} vs τ_n正向 |
| K031 (密文态) | K019 (SlotToCoeff正向) | 互补 — 密文语义 vs 数学推导 |
| K032 (C^B嵌入) | K022 (Y=X^{N/(2n)}代换) | 互补 — 嵌入 vs 代换 |
| K033 (PaCo结构) | K024 (PaCo摘要) | 互补 — 算法结构 vs 全文摘要 |
| K034 (内积) | K021 (Trace旋转分解) | 互补 — 应用 vs 理论 |
| K035 (5的幂) | K009 (F_p分解) | 互补 — 数论性质 vs 代数分解 |
| K036 (τ_n正向) | K018 (典范嵌入), K030 (τ_n^{-1}) | 互补 — 正向 vs 逆变换 |

**结论：无重复。** 所有 7 张新卡与已有卡片均为互补关系。

---

## 8. 是否可以继续 Batch-04

✅ **可以继续。** Batch-03 收尾修正已完成：
- 全局索引已重建（36 卡片 + 139 Q&A）
- K030-K036 来源一致性检查通过
- 发布目标已按内容类型重新标注
- 核原文清单已明确

Batch-04（~10 外部论文）不涉及 translation_note，处理复杂度较低。
