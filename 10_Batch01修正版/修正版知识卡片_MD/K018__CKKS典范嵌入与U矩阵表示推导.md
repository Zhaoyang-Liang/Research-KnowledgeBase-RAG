# CKKS典范嵌入与U矩阵表示推导

- **卡片ID**: K018 | **版本**: v2
- **分类**: FHE核心技术模块 / CKKS/编码与解码
- **source_id**: `src_000018`
- **来源**: `01_源文件库/04_FHE核心技术模块/典范嵌入极简与U矩阵表示推导.md`
- **置信度**: high | **需审核**: False
- **声明**: 我的笔记

## 摘要
CKKS编码使用典范嵌入σ:R→C^{N/2}将多项式映射为复数向量：取多项式在N/2个本原单位根ζ^{s_j}处的求值。完整典范嵌入τ_n使用全部N个根求值(U_n矩阵)。U_n是N×N的Vandermonde矩阵，CKKS使用U_n的一半（取共轭不重复的根）。U_n矩阵在slot-to-coeff中关键——将系数表示转换为slot表示。

## 关键知识点

### K018-1: 典范嵌入τ_n
- `τ_n(a)=(a(ζ^{s₀}),…,a(ζ^{s_{N-1}})), S⊂Z_m^×`
- τ_n将多项式a(X)映射到在N个本原单位根处的求值。s_j∈Z_m^×/⟨p⟩代表元。

### K018-2: U_n矩阵
- `U_n=(ζ^{s_i·s_j})_{i,j=0}^{N-1}, N×N Vandermonde`
- U_n的(i,j)元=ζ^{s_i·s_j}。U_n·coeff = slot_values。行列式非零→可逆。

### K018-3: CKKS使用半槽
- `σ(a)=(a(ζ^{s₀}),…,a(ζ^{s_{N/2-1}})), 共轭约束`
- 因a∈R的系数是实数，a(ζ^−s_j)=conj(a(ζ^{s_j}))→后半槽由共轭自动确定→只需N/2个槽。

### K018-4: U_n逆矩阵与slot-to-coeff
- `coeff = U_n^{-1}·slot; SlotToCoeff≈同态矩阵乘法`
- 从slot值恢复多项式系数需U_n^{-1}。同态实现时需将U_n^{-1}乘法转化为可同态计算的操作。

### K018-5: U_n矩阵的Cooley-Tukey分解
- `U_n因子化→logN层sparse矩阵乘积→高效slot-to-coeff`
- PaCo论文利用U_n的CT分解将O(N²)的矩阵乘法降为O(N log N)，是PaCo加速的核心。

## 重要公式
- `τ_n(a)=(a(ζ^{s₀}),…,a(ζ^{s_{N-1}}))` — 典范嵌入定义
- `U_n=(ζ^{s_i·s_j})_{i,j=0}^{N-1}` — 典范嵌入矩阵
- `σ:R→C^{N/2}, a↦(a(ζ^{s₀}),…,a(ζ^{s_{N/2-1}}))` — CKKS编码使用的半槽映射

## 术语
- **典范嵌入** (canonical embedding): 将代数数域元素映射为复数向量的标准方法
- **U_n矩阵** (U_n matrix): 典范嵌入的Vandermonde矩阵表示
- **SlotToCoeff** (SlotToCoeff): 从slot表示恢复系数表示的同态操作

## 关系
- abstract_algebra: 复嵌入理论
- RLWE_lattice: -
- FHE_technique: CKKS编码核心、PaCo slot-to-coeff关键矩阵
- current_paper: 当前论文PaCo bootstrapping的slot-to-coeff组件

## RAG建议: 极高优先级，CKKS编码与PaCo核心
