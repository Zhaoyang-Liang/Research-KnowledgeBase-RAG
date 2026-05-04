# Metadata 字段规范 v2

> 版本: v2 | 生成: 2026-05-04

## 必需字段

| 字段 | 类型 | 说明 | 来源 |
|------|------|------|------|
| chunk_id | string | 唯一标识 | 自动生成 |
| text | string | chunk 文本内容 | 切分 |
| source_id | string | 源文件 ID | 映射表 |
| source_relative_path | string | 库内相对路径 | 映射表 |

## 推荐字段

| 字段 | 类型 | 说明 |
|------|------|------|
| card_id | string | 关联知识卡片 |
| qa_id | string | 关联原子问答 |
| category | string | 学科分类 |
| subtopic | string | 子主题 |
| needs_human_review | bool | 需确认标记 |
| confidence | string | 置信度 |
| claim_source_type | string | 声明来源 |
| language | string | zh/en/mixed |
| doc_type | string | card/qa/source/formula/term |

## 检索元数据

| 字段 | 类型 | 说明 |
|------|------|------|
| embedding_model | string | 使用模型 |
| embedding_dim | int | 向量维度 |
| chunk_index | int | chunk 序号 |
| parent_doc | string | 父文档路径 |
