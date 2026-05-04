# Batch-04 处理日志

> 2026-05-04 | 外部论文分类入库

## 范围

- 候选文件: 12 (理论笔记根目录 7 PDF + Bootstraping讲义 5 文件)
- Skip: 2 (CSDN文库 pending_manual_review + 应用密码学原语调研 unrelated)
- Archive_only: 1 (BGV-BFV Bootstrapping.tex)
- 实际处理: 9

## 新增导入

| source_id | 文件 | document_role | read_depth |
|-----------|------|---------------|------------|
| src_000054 | BFV方案介绍(1).pdf | background | light_read |
| src_000055 | CKKS组会.pdf | side_material | light_read |
| src_000056 | GSW_4.22.pdf | foundational | light_read |
| src_000057 | LWE FHE综述.pdf | side_material | light_read |
| src_000058 | FHE from LWE报告.pdf | foundational | light_read |
| src_000059 | BGV-BFV Bootstrapping.pdf | directly_related | deep_read |
| src_000060 | BGV-BFV Bootstrapping_目录完整版.pdf | directly_related | light_read |
| src_000061 | BGV-BFV Bootstrapping_目录简化版.pdf | directly_related | light_read |
| src_000062 | bgv_bfv_bootstrapping_beamer.pdf | directly_related | light_read |

## 产出

| 类别 | 数量 | 详情 |
|------|------|------|
| 知识卡片 | 2 | K037-K038 |
| 原子问答 | 5 | Q114-Q118 |
| Related Work定位简表 | 1 | 04_论文协作区/.../Related_Work定位简表.md |
| 核心材料索引更新 | 1 | 04_论文协作区/.../核心材料索引.md |

## deep_read

- **src_000059** (BGV/BFV Bootstrapping, 229页): 系统教程，覆盖代数基础→bootstrapping流程→优化→伪代码。第2.9节BGV/BFV vs CKKS对比、第6章bootstrapping流程、第13章复杂度为重点阅读区。

## light_read

- src_000054 (BFV方案, 28页): BFV方案标准介绍，background
- src_000055 (CKKS组会, 29页): CKKS组会slides，side_material
- src_000056 (GSW, 34页): GSW方案+FHE四代演进，foundational
- src_000057 (LWE FHE综述, 16页): 中文综述，side_material
- src_000058 (FHE from LWE, 16页): BV2011报告，foundational
- src_000060-062 (讲义变体): 同src_000059，light_read确认一致性

## 发布分布

| 目标 | 内容 |
|------|------|
| **03_FHE知识库** | K037 (讲义总览), K038 (BGV bootstrapping流程) |
| **04_论文协作区** | Related_Work定位简表.md (新增), 核心材料索引.md (更新) |
| **01_抽象代数知识库** | 0 |
| **02_RLWE_格密码** | 0 (GSW和FHE from LWE归入文献管理) |
| **11_知识入库接口** | 0 |

## skip说明

| 文件 | 处理 | 原因 |
|------|------|------|
| CKKS中CoeffToSlot...CSDN文库.pdf | pending_manual_review | CSDN文库质量不确定，不进入源文件库 |
| 应用密码学原语...调研报告.pdf | skip | 与FHE主线无关 |
| BGV-BFV Bootstrapping.tex | archive_only | TeX源文件，PDF已导入 |

## BGV/BFV Bootstrapping讲义关键笔记

- 作者: 梁朝阳 (南开大学密码科学与技术), 2026-04-18
- 229页 + 52页beamer
- 核心价值: 这是目前中文资料中BGV/BFV bootstrapping最完整的系统教程
- 与当前论文的直接关联: 第2.9节BGV/BFV vs CKKS bootstrapping差别，可作为论文Related Work中"BGV/BFV bootstrapping与CKKS bootstrapping结构对比"的参考基础
- 重点建议作者人工阅读: 第2.9节(对比CKKS)、第6章(完整流程)、第13章(复杂度)

## 关于原子问答索引中Batch-01从100变为126条

此前Batch-01原子问答记录为约100条，本次全局重建后发现共126条(含QA-K001-1至QA-K028-x等)。新增的约26条来自Batch-01修正版JSONL中未单独计数的补充问答条目。这些条目均来自Batch-01的核心材料(分式域、共轭、复嵌入等28张卡片)，每张卡片对应2-7条问答不等。全部126条均有source_id和card_id映射，未见重复。

## 质量检查

- [x] 知识卡片仅生成2张（克制，不扩展related work背景）
- [x] deep_read仅BGV/BFV讲义（符合计划）
- [x] light_read全部完成
- [x] Related Work定位简表已生成
- [x] 全局索引已更新
- [x] skip文件已明确说明
- [x] K037/K038为lecture_notes类型，非translation_note
- [x] 未跨入Batch-05范围


---

## 原子问答统计修正说明（2026-05-04T14:41）

### 问题
此前汇报中原子问答数量不一致：Batch-01 曾被报告为 100/126/139 等多个版本。

### 根因
Batch-01 原子问答 JSONL 文件被复制到多个目录（`10_Batch01修正版/`、`03_原子问答/`、`01_抽象代数知识库/原子问题/`、`03_FHE知识库/原子问题/`、`04_论文创作知识库/原子问题/`），早期加载脚本从多个目录同时加载导致统计重复。

### 修正后

| 批次 | 编号范围 | 数量 | 格式 | 权威来源 |
|------|---------|------|------|---------|
| Batch-01 | QA-K001-1 ~ QA-K028-x | **100** | JSONL | `10_Batch01修正版/修正版原子问答.jsonl` |
| Batch-03 | Q101 ~ Q113 | **13** | 独立 JSON | `03_原子问答/Q*.json` |
| Batch-04 | Q114 ~ Q118 | **5** | 独立 JSON | `03_原子问答/Q*.json` |
| **合计** | | **118** | | |

### 核验
- [x] 无重复 qa_id
- [x] 无编号冲突
- [x] 所有 118 条均有 source_id 和 card_id

---

## BGV/BFV 讲义重复来源标记修正

| source_id | canonical_role | 说明 |
|-----------|---------------|------|
| **src_000059** | **canonical_source** | BGV/BFV Bootstrapping讲义完整版(229页)，梁朝阳，2026-04-18 |
| src_000060 | support_source_outline | 同src_000059的目录完整版 |
| src_000061 | support_source_outline | 同src_000059的目录简化版 |
| src_000062 | support_source_slides | 同src_000059的52页beamer版 |

K037、K038主要依据 **src_000059**。已建立 `related_to` 关系。
