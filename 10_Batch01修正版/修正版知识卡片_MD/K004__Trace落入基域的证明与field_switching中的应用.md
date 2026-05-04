# Trace落入基域的证明与field switching中的应用

- **卡片ID**: K004 | **版本**: v2
- **分类**: 数学基础 / 域论/Trace与Norm
- **source_id**: `src_000004`
- **来源文件**: `01_源文件库/01_抽象代数与代数数论/Trace与固定.md`
- **置信度**: high | **需要人工审核**: False
- **声明来源**: 我的笔记

## 摘要
扩张K/F的迹Tr_{K/F}(α)=Σ_{σ∈Gal(K/F)}σ(α)一定落入基域F，因为迹和被所有自同构固定（重排不变）。在field switching中，Tr_{n→B}通过旋转-加法链实现同态降维，将大环密文投影到子环，复杂度从O(n²)降到O(n log n)。这是GHPS field switching算法的核心操作。

## 关键知识点

### K004-1: 迹的代数定义
- 公式: `Tr_{K/F}(α)=Σ_{σ∈Gal(K/F)}σ(α)∈F`
- 含义: 扩张K/F中所有F-自同构下像的和，结果必在基域。

### K004-2: 落入基域的证明
- 公式: `τ(Tr(α))=Στ∘σ(α)=Σσ'(α)=Tr(α), ∀τ∈Gal`
- 含义: 自同构作用于迹和只是重排求和，故迹不变→被所有自同构固定→必在基域。

### K004-3: 同态迹的实现
- 公式: `Tr(c)=c+Rot₁(c)+Rot₂(c)+…+Rot_{d-1}(c)`
- 含义: 通过d-1次自同构旋转再加和实现同态迹。

### K004-4: Field switching中的迹
- 公式: `Tr_{n→B}: R_E→R_{E_B}, 将大环多项式投影到子环`
- 含义: 迹降维是field switching核心：减小bootstrapping密文尺寸。

### K004-5: 迹与范数的区别
- 公式: `Norm_{K/F}(α)=Πσ(α); Tr=Σ, Norm=Π`
- 含义: 迹是求和(加法性)，范数是乘积(乘法性)。FHE中主要用迹的加法结构。

## 重要公式
- `Tr_{K/F}(α)=Σ_{i=0}^{d-1}τⁱ(α), d=[K:F]` — Galois扩张中的迹公式
- `Tr(c)=c+Rot₁(c)+…+Rot_{d-1}(c)` — 同态迹的旋转-加法实现

## 术语
- **迹** (trace): 扩张中元素在所有F-自同构下像的和
- **基域** (base field): 扩张的起始域F
- **Field switching** (field switching): GHPS技术：通过迹将密文从大分圆环转换到小子环

## 关系
- abstract_algebra: 域扩张的迹理论
- RLWE_lattice: 迹满足Tr(a+b)=Tr(a)+Tr(b)，线性操作
- FHE_technique: field switching / bootstrapping降维核心
- current_paper: 当前论文field switching核心依赖迹

## RAG入库建议: 极高优先级，field switching数学基础
