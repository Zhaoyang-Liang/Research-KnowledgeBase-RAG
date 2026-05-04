# Batch 依赖关系图

> 版本: v1 | 生成: 2026-05-04

```text
Batch-01 (基础链条) ✅
  |
  v
Batch-02 (P0 核实)
  |    \
  v     \
Batch-03 (论文 core)   Batch-04 (FHE 文献) [可并行]
  |       \              |
  v        \             v
Batch-05 (Bootstrapping)  Batch-07 (分析模块)
  |                        |
  v                        v
Batch-06 (格密码 P2)    Batch-08 (实现 P2)
  |                        |
  +----- Batch-09 (文献) --+
              |
              v
         Batch-10 (RAG 入库)
```

## 依赖说明

| Batch | 直接依赖 | 原因 |
|-------|---------|------|
| Batch-02 | Batch-01 | 需要 v2 卡片中的 P0 标注 |
| Batch-03 | Batch-02 | P0 核实后才能处理论文 core |
| Batch-04 | 无 | 可独立启动 |
| Batch-05 | Batch-02 + Batch-03 | 需要 PaCo P0 核实 + 论文 core |
| Batch-06 | 无 | 可独立启动 |
| Batch-07 | Batch-03 | 分析材料来自论文 |
| Batch-08 | Batch-03 | 实验数据来自论文 |
| Batch-09 | Batch-01-08 | 需要全局文献视图 |
| Batch-10 | 全部 P0 完成 | RAG 需要可信数据 |

## 论文知识沉淀流 (穿插于各 Batch)

```text
论文项目产生想法
  |
  v
11_知识入库接口 (候选箱)
  |
  +-- draft_idea
  +-- candidate_knowledge
  +-- needs_validation
  +-- validated
  |
  v
正式知识库 (published_to_kb)
```
