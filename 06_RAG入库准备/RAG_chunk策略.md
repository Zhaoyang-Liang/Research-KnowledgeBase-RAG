# RAG Chunk 策略

> 2026-05-04T18:00 | 知识库工程完成阶段 | 仅策略，不执行

## 1. 知识卡片 Chunk（优先入 RAG）

每张卡片按结构化字段切分为多个 chunk：

### summary 层（1 chunk/card）
包含：card_id, title, domain, topic, confidence, source_id。
适合：概览检索，"有哪些关于 X 的知识？"

### knowledge_points 层（N chunks/card）
每个 knowledge_point 独立为一个 chunk。
包含：card_id, kp_index, kp_text, 关联的公式和术语。
适合：精确检索，"C2S 的方向是什么？"

### formulas 层（可选）
如果卡片包含独立公式，单独切分。

### relations 层（可选）
卡片间的 cross_reference 作为关系边。

### 附加 metadata
每个 chunk 必须携带：
- source_id（可回溯到源文件）
- canonical_card_path（RAG 只读 canonical 版本）
- rag_status

## 2. 原子问答 Chunk

一条 Q&A 一个 chunk。
包含：qa_id, question, answer, card_id, source_id。
适合：问答式检索。

## 3. 源文件 Chunk（暂缓）

只对 source_status 稳定（非 translation_note、非草稿、非实验日志）的源文件切分。
策略：按段落切分，保留 section 层级。
当前暂不执行——先等卡片和 Q&A 的 RAG 验证后再启动。

## 4. 关系图 Chunk（暂缓）

适合做概念检索的关系图：
- 抽象代数到 FHE 关联
- 概念依赖图
- 参数关系图
- C2S/S2C 术语规则
- R vs R^∨ 符号对照

当前关系图内容以 Markdown 表格为主，可逐行切分为 relation chunk。

## 5. Candidate Chunk（默认 defer）

candidate_knowledge（C-NOISE-001~010）不进 ready。
等作者确认后升级。

## 6. 切分粒度建议

| 内容类型 | 目标大小 | 切分方式 |
|----------|---------|----------|
| card summary | 200-400 chars | 单一 chunk |
| knowledge_point | 300-800 chars | 每 KB 一个 chunk |
| Q&A | 完整 pair | 单一 chunk |
| source paragraph | 500-1000 chars | 按段落 |

## 7. 执行顺序

1. 先对 ready 卡片做 sample chunk（选 5 张）
2. 手动检查 chunk 质量
3. 如合格，批量切分全部 ready 卡片和 Q&A
4. 不做 embedding
5. metadata index 先行，chunk 内容准备跟随
