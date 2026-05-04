# Babai最近平面算法

- **卡片ID**: K027 | **版本**: v2
- **分类**: 格密码/后量子 / CVP算法
- **source_id**: `src_000025`
- **来源**: `01_源文件库/02_RLWE_格密码_后量子/Babai.md`
- **置信度**: high | **需审核**: False
- **声明**: 我的笔记

## 摘要
Babai最近平面算法是解决格上最近向量问题(CVP)的近似算法。利用Gram-Schmidt正交化将格基分层：目标向量s投影到各层平面，逐层选择最近的偏移c·b̃_n（c=⌊⟨s,b̃_n⟩/|b̃_n|²⌉取整），递归降维。在FHE中Babai算法用于解密时的噪声估计和格归约攻击分析。尽管非核心FHE技术，理解它有助于理解RLWE的安全性。

## 关键知识点

### K027-1: Gram-Schmidt分层
- `b̃_i=b_i−Σ_{j<i}μ_{ij}·b̃_j, μ_{ij}=⟨b_i,b̃_j⟩/|b̃_j|²`
- 将格基正交化，得到正交基b̃_i→定义各层平面P_n=span(b₁,…,b_{n-1})。

### K027-2: α=⟨s,b̃_n⟩/|b̃_n|²的含义
- `α是将s投影到b̃_n方向的坐标`
- s可写为α·b̃_n+(span(b₁,…,b_{n-1})中的分量)。取整c=⌊α⌉最小化距离。

### K027-3: c=⌊α⌉的原因
- `|α−c|·|b̃_n|最小化c为整数时`
- 各层平面在b̃_n方向间隔|b̃_n|，选c=⌊α⌉使s距第c层平面最近。

### K027-4: 递归降维
- `s'=s−c·b_n, 继续在n−1维格上`
- 选定c后投影到低维格→递归→复杂度O(n³)（含Gram-Schmidt）。

### K027-5: FHE中的关联
- `RLWE解密→CVP→Babai可近似求解但非多项式(精确CVP是NP-hard)`
- Babai给的是近似解。RLWE安全性基于CVP的难度。在FHE中Babai用于格攻击分析而非构造。

## 重要公式
- `c=⌊⟨s,b̃_n⟩/|b̃_n|²⌉` — Babai最近平面选择规则
- `s'=s−c·b_n, 递归至dim=1` — 降维递归

## 术语
- **最近向量问题** (Closest Vector Problem (CVP)): 给定格和目标向量，求格中最接近的向量
- **Gram-Schmidt正交化** (Gram-Schmidt orthogonalization): 将线性无关向量组转化为正交向量组
- **Babai算法** (Babai's nearest plane algorithm): CVP的近似算法，通过分层投影求解

## 关系
- abstract_algebra: 线性代数/格
- RLWE_lattice: CVP→RLWE安全性
- FHE_technique: 噪声估计与攻击分析
- current_paper: 间接相关（安全性分析）

## RAG建议: 中优先级，格密码安全性
