# RAG Metadata Schema

> 2026-05-04T18:00 | 知识库工程完成阶段

## Chunk Metadata 字段

```json
{
  "chunk_id": "card_K001_kp_3",
  "chunk_type": "card_knowledge_point",
  "source_id": "src_000001",
  "source_path": "01_源文件库/02_RLWE_格密码_后量子知识库/...",
  "card_id": "K001",
  "card_path": "01_抽象代数知识库/文档知识卡片/K001__...json",
  "qa_id": null,
  "domain": "abstract_algebra",
  "topic": "分裂域/分式域",
  "confidence": "medium",
  "needs_human_review": false,
  "knowledge_status": "source_preserved",
  "claim_source_type": "用户笔记",
  "paper_id": null,
  "idea_id": null,
  "candidate_id": null,
  "paper_project_status": null,
  "rag_status": "ready",
  "canonical_path": "01_抽象代数知识库/文档知识卡片/K001__...json",
  "citation_path": "01_源文件库/.../对应源文件"
}
```

## chunk_type 枚举

| 值 | 说明 | 入 RAG |
|----|------|--------|
| `card_knowledge_point` | 知识卡片的单个知识点 | ✅ ready |
| `card_summary` | 卡片整体摘要 | ✅ ready |
| `qa_pair` | 原子问答 | ✅ ready |
| `source_paragraph` | 源文件段落（仅 stable） | ⏸ defer |
| `relation_edge` | 知识关联关系 | ⏸ defer |
| `candidate_knowledge` | 候选知识（论文专属 idea） | ⏸ defer |
| `paper_idea` | 论文新 idea（未经确认） | ⏸ defer |
| `paper_draft` | 论文草稿 | ❌ exclude |
| `experiment_log` | 实验日志 | ❌ exclude |
| `translation_raw` | 翻译笔记原文 | ❌ exclude |
| `negative_result` | 已确认的失败探索 | ✅ ready（标记） |

## 关键约束

1. **回溯能力**：每个 chunk 必须可通过 `source_id` + `card_id` 或 `qa_id` 追溯到原始来源
2. **状态感知**：`candidate_knowledge` 和 `needs_human_review: true` 的内容默认不进 ready
3. **论文隔离**：`paper_id` 不为空的 chunk 默认 defer，仅作者确认后升级
4. **副本排除**：RAG 只使用 `canonical_card_path`，不读取 `duplicate_paths` 中的副本
5. **低置信排除**：`confidence: low` 的 chunk 不进 ready
6. **论文 idea 保护**：`idea_id` 不为空的 chunk，论文专属内容不进 ready（除 `negative_result` 已确认外）
7. **跨论文复用**：`knowledge_status: validated` + `paper_id: null` 的内容可跨论文直接使用
8. **论文项目状态**：`paper_project_status: archived` 的论文专属 chunk 可选择性降级为 exclude

## 统计

| chunk_type | ready | defer | exclude |
|------------|-------|-------|---------|
| card_knowledge_point | 45 | 24 | 0 |
| qa_pair | 165 | 21 | 0 |
| candidate_knowledge | 0 | 10 | 0 |
| paper_idea | 0 | 0 | 0 |
| negative_result | 0 | 0 | 0 |
| source_paragraph | 0 | 待定 | 0 |

## 未来论文项目 RAG 对接

### rag_status 决策矩阵

| knowledge_status | paper_id | idea_id | rag_status |
|-----------------|----------|---------|-----------|
| validated | null | null | ✅ ready |
| validated | 有值 | null | ✅ ready（论文产出沉淀为通用知识后） |
| candidate | 有值 | 有值 | ⏸ defer |
| needs_validation | 有值 | 有值 | ⏸ defer |
| negative_result | 有值 | 有值 | ✅ ready（标记为 negative_result） |
| deprecated | 有值 | 有值 | ❌ exclude |

### 论文 idea 保护墙

```
论文项目的 idea（paper_id != null, status != validated）
  → chunk_type = "paper_idea" 或 "candidate_knowledge"
  → rag_status = defer
  → 检索时不会作为确定答案返回
  → 仅在明确查询 "{paper_id} 有哪些未解决 idea" 时返回
```
