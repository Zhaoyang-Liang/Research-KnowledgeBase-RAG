# SIMD代数总结：Frobenius轨道→不可约因子→槽的三线推导

- **卡片ID**: K008 | **版本**: v2
- **分类**: 数学基础 / 代数/SIMD与Frobenius轨道
- **source_id**: `src_000008`
- **来源文件**: `01_源文件库/01_抽象代数与代数数论/SIMD代数版本总结.md`
- **置信度**: high | **需要人工审核**: False
- **声明来源**: 我的笔记

## 摘要
SIMD槽结构可从三条线理解：(1)群论线——S=Z_m^×/<p>商群代表元，每个代表元k对应一个槽；(2)有限域线——d=ord_m(p)，不可约因子f_j次数=d，槽=F_{p^d}；(3)Frobenius轨道线——Frobenius映射x→x^p在ζ_m根上形成轨道，每轨道对应一个不可约因子和槽。三条线统一描述同一结构。CKKS槽=BGV/BFV槽数的一半（因共轭减半）。

## 关键知识点

### K008-1: 群论线：商群S
- 公式: `|S|=φ(m)/d=L, S={r₁,…,r_L}`
- 含义: L个代表元定义L个槽。CKKS中(r_j,-r_j)合并为一个复槽。

### K008-2: 有限域线：不可约因子
- 公式: `Φ_m(X)≡∏_{j=1}^{L}f_j(X)(mod p), deg(f_j)=d`
- 含义: d=ord_m(p), L=φ(m)/d。每个因子定义一个槽=F_{p^d}。

### K008-3: Frobenius轨道线
- 公式: `Orb(ζ_m^{r_j})={ζ_m^{r_j},ζ_m^{pr_j},…,ζ_m^{p^{d-1}r_j}}`
- 含义: Frobenius映射x→x^p作用于根，轨道长d，轨数L=φ(m)/d。

### K008-4: CKKS槽vs BGV/BFV槽
- 公式: `CKKS: slot=paired real/complex(半槽); BGV/BFV: slot=F_{p^d}(全槽)`
- 含义: 来源相同但CKKS利用复嵌入取半(N/2个槽)，BGV/BFV用CRT取全部L个槽。

### K008-5: 三条线统一
- 公式: `Z_m^×/<p>≅{Frobenius轨道}≅{不可约因子f_j}`
- 含义: 三种视角描述同一槽结构，选视角取决于当前问题。

## 重要公式
- `d=ord_m(p)=min{k>0:p^k≡1(mod m)}` — p模m的乘法阶=每个不可约因子次数
- `L=φ(m)/d` — 不可约因子数=BGV/BFV全部槽数，CKKS槽数的一半
- `Orb(ζ_m^r)={ζ_m^{rp^i}:i=0,…,d-1}` — Frobenius轨道定义

## 术语
- **Frobenius自同构** (Frobenius automorphism): 特征p域上的映射x→x^p
- **Frobenius轨道** (Frobenius orbit): Frobenius映射下元素的轨道
- **ord_m(p)** (multiplicative order): 使p^k≡1(mod m)的最小正整数k

## 关系
- abstract_algebra: 有限域Galois理论、群论
- RLWE_lattice: RLWE环CRT分解→SIMD
- FHE_technique: 所有方案SIMD编码代数基础
- current_paper: PaCo多槽编码、field switching槽重映射

## RAG入库建议: 极高优先级，SIMD理论基础
