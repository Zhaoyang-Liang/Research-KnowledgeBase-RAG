# Batch 质量检查清单

> 版本: v1 | 使用: 每个 Batch 完成后逐项检查

## 路径完整性

- [ ] 所有 source_id 指向的文件物理存在
- [ ] 所有生成文件中无嵌套路径（如 `01_源文件库/.../01_源文件库/...`）
- [ ] 所有生成文件中无绝对路径 `~/Desktop/Eurocrypt 2027`
- [ ] 所有生成文件中无临时池简写（如 `同态加密复习整理/`）

## 引用完整性

- [ ] 每个知识卡片有 card_id
- [ ] 每个知识卡片有 source_id
- [ ] 每个知识卡片有 source_relative_path
- [ ] 每个原子问答有 source_id
- [ ] 每个公式有 source_id
- [ ] 来源映射中所有 source_id 有对应文件

## 内容质量

- [ ] 知识卡片有 chinese_summary
- [ ] 知识卡片有 knowledge_points（3-5 个）
- [ ] 不确定内容标注 needs_human_review
- [ ] 推断内容标注 claim_source_type = "Qclaw推断"
- [ ] 术语使用时区分 CKKS slot vs CRT slot
- [ ] BFV Rescale vs CKKS Rescale 不混用
- [ ] 术语表中每种术语有中英文对照

## 发布前检查

- [ ] P0 问题已核实或已标注
- [ ] 发布日志已生成
- [ ] 未覆盖旧文件（或已有版本号）
- [ ] 未修改 01_源文件库
- [ ] 未修改临时资料池
- [ ] 卡片索引和 manifest 已更新

## RAG 前检查

- [ ] 所有 P0 问题已核实
- [ ] needs_human_review 条目已被用户确认
- [ ] 无笔记来源的未确认推导
