# Batch-03 处理日志

> 2026-05-04 | 当前论文相关补充材料与 Geelen/SlotToCoeff 翻译笔记整理

## 范围

- 候选文件: 21 (翻译 13 + 代数基础 6 + 规划文件 2)
- 已导入: 12 (src_000001~000006, src_000018~000024) — 不重复处理
- 新增导入: 9 (src_000045~000053)

## 新增导入文件

| source_id | 文件 | 类型 |
|-----------|------|------|
| src_000045 | 翻译/C^B到C^N的嵌入.md | translation_note |
| src_000046 | 翻译/SlotToCeff-密文下例子.md | translation_note |
| src_000047 | 翻译/algo1.md | translation_note |
| src_000048 | 翻译/同态内积操作.md | translation_note |
| src_000049 | 翻译/数字5mod2的幂的性质.md | translation_note |
| src_000050 | 翻译/SlotToCeff-CKKS/CKKS编码正向.md | translation_note |
| src_000051 | 翻译/SlotToCeff-CKKS/tao(p(Y)).md | translation_note |
| src_000052 | 安全性_复杂度_噪声分析/汇总代办路线.md | paper_planning |
| src_000053 | 安全性_复杂度_噪声分析/可能的尝试写法.md | paper_idea |

## 产出

| 类别 | 数量 | 详情 |
|------|------|------|
| 知识卡片 | 7 | K030-K036 |
| 原子问答 | 13 | Q101-Q113（以 QA-Kxxx 格式计入 Batch-01 的 126 条不包含在内） |
| 术语 | 8 | SlotToCoeff语义/τ_n/C^B嵌入/PaCo密钥/Trace/Product/split-compatible/5的阶 |
| 公式 | 10 | τ_n^{-1}/SlotToCoeff密文/嵌入编码保持/PaCo压缩/Trace/τ_n正向/5的阶/criterion/γ_impl/主定理 |
| 候选知识 | 4 | cand_001~004 (split-compatible criterion/CKKS instance/γ_impl/leaf security) |

## 发布分布（修正后）

| 目标 | 内容 |
|------|------|
| **03_FHE知识库 (support_ingest, translation_note)** | K030-K036 全部7张卡片 |
| **01_抽象代数知识库** | 0（Batch-03无新增代数卡片，代数基础已在Batch-01处理） |
| **02_RLWE_格密码_后量子知识库** | 0 |
| **04_论文协作区** | 汇总代办路线.md (完整版) + 论文写作方向_通用框架版.md |
| **11_知识入库接口** | cand_001~004 (candidate_knowledge, needs_validation) |

## 卡片状态明细

| card_id | 与Batch-01重复？ | needs_review | 发布目标 | 核原文需求 |
|---------|-----------------|-------------|----------|-----------|
| K030 | ❌ 互补(K018/K019) | ⚠ | 03_FHE (translation_note) | 验证U_n矩阵形式 |
| K031 | ❌ 互补(K019) | ⚠ | 03_FHE (translation_note) | 验证密文态语义 |
| K032 | ❌ 唯一主题 | ⚠ | 03_FHE (needs_verification) | 确认2M→2B修正 |
| K033 | ❌ 互补(K024) | ⚠ | 03_FHE (needs_verification) | 对照PaCo原文 |
| K034 | ❌ 互补(K021) | ⚠ | 03_FHE (translation_note) | 基础内容 |
| K035 | ❌ 唯一主题 | ⚠ | 03_FHE (needs_verification) | 标题表述调整 |
| K036 | ❌ 互补(K018/K030) | ⚠ | 03_FHE (translation_note) | 基础内容 |

## 质量检查

- [x] 翻译文件全部标记 translation_note
- [x] 翻译内容未宣称等同于原文结论
- [x] 与 Batch-01 已有卡片做标题/主题去重 — 全部为互补关系
- [x] 规划文件进入论文协作区，未发布为 validated knowledge
- [x] 论文想法标记为 candidate_knowledge
- [x] 所有新增 source_id 已更新索引
- [x] K030-K036 来源一致性检查通过
- [x] K032 原文笔误保留 needs_human_review
- [x] K033 标注来源为翻译笔记非PaCo原文
- [x] K035 标题表述保守标注
- [x] 发布目标按内容类型重新标注

## 下一步

建议 Batch-04：外部论文分类入库（~10个文件，BGV/BFV Bootstrapping精读，其余轻读）
