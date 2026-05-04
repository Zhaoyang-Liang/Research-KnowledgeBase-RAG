# NTT与Cooley-Tukey快速算法

- **卡片ID**: K011 | **版本**: v2
- **分类**: 数学基础 / NTT/FFT
- **source_id**: `src_000017`
- **来源**: `01_源文件库/01_抽象代数与代数数论/变换与NTT-Cooley-Tukey.md`
- **置信度**: high | **需审核**: False
- **声明**: 我的笔记

## 摘要
NTT(数论变换)是有限域F_q上的DFT，使用本原n次单位根ω(满足ω^n=1, ω^k≠1)，运算精确无浮点误差。Cooley-Tukey算法将NTT分解为log₂n层蝴蝶操作x'=x+ω^k·y, y'=x−ω^k·y。位反转因递归分奇偶打乱输出顺序。在BGV同态NTT中位反转被吸收到索引映射中避免显式旋转。

## 关键知识点

### K011-1: DFT vs NTT
- `DFT: A_k=Σa_j·e^{2πijk/n}; NTT: A_k=Σa_j·ω^{jk} mod q`
- 形式同但运算域不同：DFT复数浮点误差，NTT有限域精确。要求n|(q-1)。

### K011-2: Cooley-Tukey蝴蝶
- `x'=x+ω^k·y, y'=x−ω^k·y (mod q)`
- 每层n/2蝴蝶，log₂n层→O(n log n)。加/减来自偶奇子问题合并。

### K011-3: 位反转
- `rev(k)=Σb_i·2^{log₂n−1−i}, 索引按bits反转`
- 递归分奇偶→输入按位反转输出自然序。实践中在输入或输出端做一次。

### K011-4: NTT与多项式乘法
- `c=a∗b⇔NTT(c)=NTT(a)⊙NTT(b)`
- 卷积定理：乘法的NTT=系数的NTT逐点乘。O(n log n) vs O(n²)。

### K011-5: BGV同态NTT的位反转处理
- `bit-reversal absorbed via pre-permuted index maps`
- 位反转不通过显式旋转而通过预排列索引映射实现，避免昂贵同态旋转。

## 重要公式
- `A_k=Σ_{j=0}^{n-1}a_j·ω^{jk} (mod q), ω^n≡1, ω^k≠1` — NTT正变换
- `x'=x+ω^k·y, y'=x−ω^k·y` — Cooley-Tukey蝴蝶操作

## 术语
- **NTT** (Number Theoretic Transform): 有限域F_q上DFT
- **蝴蝶操作** (butterfly operation): Cooley-Tukey基本单元，同时计算加权和与差
- **位反转** (bit-reversal): 按二进制位反转索引的置换

## 关系
- abstract_algebra: 有限域傅里叶分析
- RLWE_lattice: RLWE多项式乘法快速实现
- FHE_technique: PaCo slot-to-coeff、BGV同态NTT
- current_paper: PaCo中Vandermonde矩阵Cooley-Tukey分解

## RAG建议: 高优先级，Bootstrapping/NTT核心技术
