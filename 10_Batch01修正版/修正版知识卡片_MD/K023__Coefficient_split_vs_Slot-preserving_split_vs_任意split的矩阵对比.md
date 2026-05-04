# Coefficient split vs Slot-preserving split vs 任意split的矩阵对比

- **卡片ID**: K023 | **版本**: v2
- **分类**: FHE核心技术模块 / 多项式分解
- **source_id**: `src_000023`
- **来源**: `01_源文件库/04_FHE核心技术模块/补充视角-slot情况下的拆分.md`
- **置信度**: medium | **需审核**: True
- **声明**: 我的笔记

## 摘要
多项式拆分有三种类型：(1)Coefficient split——按系数奇偶拆分p↦(E_p,O_p)，简单直接但破坏slot结构；(2)Slot-preserving split——p↦(E_p+H·O_p, E_p−H·O_p)，保留slot信息但需额外矩阵H；(3)任意split——p↦F_n^{-1}·S·F_N·p，通用框架，S为任意对角矩阵。三种拆分对应不同的分解策略，当前论文需要明确使用的是哪种拆分方式。

## 关键知识点

### K023-1: Coefficient split
- `p↦(E_p,O_p), E_p=p(X²), O_p=p(X)·(1−(−1)ⁿ)/2 + …`
- 按系数奇偶拆分：偶系数多项式+奇系数多项式(×X^{-1})。简单但破坏slot含义。

### K023-2: Slot-preserving split
- `p↦(E_p+H·O_p, E_p−H·O_p), H=某个对角矩阵`
- 通过额外矩阵H保持slot语义→拆分后每个分量仍对应原槽的子集。

### K023-3: 任意split框架
- `p↦F_n^{-1}·S·F_N·p, S任意对角矩阵`
- F_n,F_N=NTT/DFT矩阵，S对角→最一般形式的拆分。coeff split和slot-preserving是S取特定值。

### K023-4: 对Bootstrapping的影响
- `PaCo使用哪种split? → 盲旋转中需要多项式打包→拆分为打包做准备`
- 不同的split策略影响后续的操作效率和噪声。需确认PaCo具体使用的拆分类型。

### K023-5: ⚠ 当前论文需确认
- `construction中使用coefficient split还是slot-preserving split?`
- 关键决策：coeff split简单但破坏slot→需要额外操作恢复；slot-preserving保持slot但复杂。需要人工确认。

## 重要公式
- `Coeff split: p↦(E_p,O_p)` — 系数奇偶拆分
- `Slot-preserving: p↦(E_p+H·O_p, E_p−H·O_p)` — 保持slot的拆分
- `General: p↦F_n^{-1}·S·F_N·p` — 任意矩阵拆分框架

## 术语
- **系数拆分** (coefficient split): 按多项式系数奇偶位置拆分为两个多项式
- **槽保持拆分** (slot-preserving split): 拆分后每个分量仍保留原槽语义
- **F_n** (F_n matrix): DFT/NTT变换矩阵，将多项式转换到频域

## 关系
- abstract_algebra: 多项式分解
- RLWE_lattice: -
- FHE_technique: Bootstrapping前处理
- current_paper: 当前论文construction的关键技术选择

## RAG建议: 高优先级，当前论文关键决策点，标注需人工确认
