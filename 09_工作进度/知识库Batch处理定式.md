# 知识库 Batch 处理定式

> v1 | 2026-05-04

## 标准流程 (12步)

```text
1. 确定 Batch 范围
   → 从待处理池选文件，确定类型(A/B/C)，控制数量(10-30/批)

2. 导入源文件到 01_源文件库
   → 复制，不移动不删除原文件

3. 生成 source_id
   → src_XXXXXX，记录到 原始路径到库内路径映射.jsonl

4. 轻读 / 精读决策
   → 相关论文=精读，背景资料=轻读

5. 生成知识卡片
   → JSON+MD 双格式，card_id, knowledge_points, source_id

6. 生成原子问答
   → question/answer 对，标注 card_id, source_id, confidence

7. 抽取术语和公式
   → 术语表、公式表，关联 source_id

8. 更新知识关联
   → 概念依赖图、参数关系图、交叉引用

9. 标记 needs_human_review
   → 不确定处标注，不自动标记 verified

10. 发布到正式知识区
    → 03_FHE知识库/ 按主题分类

11. 更新全局索引
    → 卡片索引、问答索引、术语公式索引

12. 准备 RAG chunk，但不自动 embedding
```

## 三种 Batch 类型

### A. 外部论文 Batch

**场景**: GHPS、PaCo、CKKS、HERMES、BGV/BFV 等 PDF

**输出重点**:
- 文献摘要和角色标注 (foundational / related_work / background)
- 核心技术卡片
- 原子问答
- related work 定位表
- 与当前论文关系的说明

**精读决策**: foundational 论文精读，background/side 论文轻读

### B. 理论笔记 Batch

**场景**: 抽象代数、RLWE、FHE 基础、证明笔记

**输出重点**:
- 概念卡片 (定义、定理、证明技巧)
- 公式表
- 术语表
- 抽象代数 → FHE 的关联链
- 可复用 proof trick

**精读决策**: 与当前论文直接相关的精读，基础概念轻读

### C. 实验 / 实现 Batch

**场景**: OpenFHE/Lattigo 参数、实验日志、benchmark

**输出重点**:
- 实验结论摘要
- 参数表
- 失败原因记录
- implementation caveat
- 是否能支撑论文 claim

**精读决策**: 支撑当前论文实验的优先，其余轻读

## 黄金规则

```text
Batch 的目标是整理和入库，不是替作者做最终判断。
一切 confidence: high 的标记都不代表作者确认。
```
