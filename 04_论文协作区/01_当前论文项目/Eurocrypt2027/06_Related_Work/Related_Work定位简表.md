# Related Work 定位简表

> 2026-05-04 | Batch-04 产物

## 一、外部论文定位

| source_id | 论文/材料 | document_role | read_depth | 与当前论文关系 | 是否进入长期知识库 | 是否进入Related Work | 建议人工阅读 |
|-----------|----------|---------------|------------|---------------|-------------------|---------------------|-------------|
| src_000054 | BFV方案介绍 | background | light_read | BFV方案基础，与当前CKKS论文间接相关 | 是(03_FHE) | 否(基础背景) | 否 |
| src_000055 | CKKS组会 | side_material | light_read | CKKS介绍，含bootstrapping概述 | 是(03_FHE) | 否(组会材料) | 否 |
| src_000056 | GSW_4.22 | foundational | light_read | GSW方案介绍，FHE四代演进综述 | 是(02_RLWE) | 可引用(历史背景) | 否 |
| src_000057 | LWE FHE综述 | side_material | light_read | 中文综述，含硬件优化，与当前论文理论方向间接相关 | 是(08_文献) | 可引用(综述引用) | 否 |
| src_000058 | FHE from LWE | foundational | light_read | BV2011，第二代FHE基础 | 是(02_RLWE) | 可引用(历史引用) | 否 |
| src_000059 | BGV/BFV Bootstrapping讲义 | directly_related | deep_read | BGV/BFV bootstrapping完整教程，与当前CKKS bootstrapping因子分解论文直接对比参考 | 是(03_FHE) | 是(对比讨论) | **是**(第2.9节+第13章) |
| src_000060 | 同上(目录完整版) | directly_related | light_read | 与src_000059内容一致，带完整目录 | 是(同src_000059) | 否(重复) | 否 |
| src_000061 | 同上(目录简化版) | directly_related | light_read | 与src_000059内容一致，带简化目录 | 是(同src_000059) | 否(重复) | 否 |
| src_000062 | 同上(beamer) | directly_related | light_read | 52页beamer版，与讲义同内容 | 是(同src_000059) | 否(重复) | 否 |

## 二、与当前论文的关联分析

### directly_related (4)
- src_000059-062: BGV/BFV Bootstrapping讲义 — **最重要的外部材料**
  - 与当前CKKS bootstrapping因子分解论文对比参考
  - 第2.9节明确讨论了BGV/BFV vs CKKS bootstrapping的差别
  - 第12-13章讨论的sparse/fully packed优化与论文的split/factorization技术有结构相似性
  - 讲义中unpack/repack的even/odd splitting与HERMES Split/Merge在思想上平行

### foundational (2)
- src_000056 (GSW): FHE四代演进，提供历史背景
- src_000058 (FHE from LWE): BV2011，第二代FHE基础

### 当前论文Related Work应讨论的外部工作（优先级排序）
1. **BGV/BFV Bootstrapping** (src_000059) — 对比CKKS vs BGV/BFV bootstrapping结构差异
2. **HERMES Ring Packing** (已导入，src_0000xx) — directly_related
3. **PaCo** (已导入，src_000029) — directly_related
4. **Geelen SlotToCoeff** (已导入，src_000030) — directly_related
5. **GSW** (src_000056) — 作为FHE历史引用
6. **FHE from LWE (BV2011)** (src_000058) — 作为FHE历史引用

## 三、skip/pending文件

| 文件 | 处理 | 原因 |
|------|------|------|
| CKKS中CoeffToSlot...CSDN文库.pdf | pending_manual_review | CSDN文库质量不确定 |
| 应用密码学原语...调研报告.pdf | skip | 与FHE主线无关 |
| BGV-BFV Bootstrapping.tex | archive_only | TeX源文件，PDF已导入 |
