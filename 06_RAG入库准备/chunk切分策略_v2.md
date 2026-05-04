# Chunk 切分策略 v2

> 版本: v2 | 生成: 2026-05-04 | 状态: 建议阶段

## 切分目标

- 每个 chunk 是一个独立的检索单元
- chunk 大小: 300-500 字（中文）/ 200-400 tokens
- 重叠: 50-100 字

## 切分来源

### 知识卡片 (推荐优先)
- 每个卡片 as-is: 250-400 字，天然 chunk
- chunk_id = card_id + "_full"
- source_id 直接从卡片获取

### 源文件 (md/tex/txt)
- 按 ## 或 ### 标题切分
- 超长段落按 500 字截断 + 50 字重叠
- chunk_id = source_id + "_chunk_NNN"

### 原子问答
- 每个问答 as-is: 100-200 字
- chunk_id = qa_id
- source_id 从问答获取

## 不切分的内容

- 数学公式集中区（保留完整性）
- 论文原文（保留原文，不切分）
- 关系图（保留结构化文本）

## 字段规范

每个 chunk 必须有:
```json
{
  "chunk_id": "唯一标识",
  "text": "切分后的文本",
  "source_id": "src_NNNNNN",
  "source_relative_path": "01_源文件库/...",
  "card_id": "KNNN (如果有)",
  "category": "分类",
  "needs_human_review": false
}
```

## 暂不执行
