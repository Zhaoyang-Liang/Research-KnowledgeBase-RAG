# ⚠ CKKS slot与BGV/BFV CRT slot的区别——避免论文混用

- **卡片ID**: K028 | **版本**: v2
- **分类**: 跨源合成 / 术语辨析
- **source_id**: `src_000009+src_000008+src_000011`
- **来源**: `跨源合成卡片（见relations）`
- **置信度**: high | **需审核**: True
- **声明**: Qclaw推断（跨源合成）

## 摘要
CKKS的slot与BGV/BFV的（CRT）slot虽然都叫slot，但数学来源不同：BGV/BFV的slot来自CRT分解——R_p≅∏F_{p^d}，每个F_{p^d}因子=一个slot，可存一个F_{p^d}元素；CKKS的slot来自复嵌入——典范嵌入将多项式映射到C^{N/2}，每个复数坐标=一个slot。关键区别：(1)BGV/BFV槽在有限域上操作（精确整数），CKKS槽在复数上操作（近似）；(2)BGV/BFV所有L=φ(m)/d个槽均可用，CKKS仅N/2个槽因共轭减半；(3)编码方式完全不同——CRT编码 vs 复嵌入编码。论文写作中必须明确区分，不可混用。

## 关键知识点

### K028-1: BGV/BFV槽的数学来源
- `R_p=F_p[X]/(Φ_m)≅∏_{j=1}^{L}F_{p^d}, 每个因子=一个CRT槽`
- CRT同构给出L=φ(m)/d个槽。每个槽存F_{p^d}元素。这是代数结构层面的槽。

### K028-2: CKKS槽的数学来源
- `σ:R→C^{N/2}, a(X)↦(a(ζ^{s₀}),…,a(ζ^{s_{N/2-1}}))`
- 复嵌入将多项式映射为N/2个复数值。每个求值点=一个CKKS槽。这是分析（嵌入）层面的槽。

### K028-3: 槽数的差异
- `BGV/BFV: L_slots=φ(m)/d; CKKS: L_slots=φ(m)/2`
- 在p≡1(mod m)的特殊情况下，BGV/BFV有φ(m)个槽，CKKS有φ(m)/2个。BGV/BFV槽数是CKKS的两倍。

### K028-4: 编码方式的根本不同
- `BGV/BFV: encode via CRT; CKKS: encode via canonical embedding`
- CRT编码逐个槽分配F_{p^d}元素→精确；典范嵌入整体映射→近似。两种编码不可互换。

### K028-5: ⚠ 论文混用风险
- `CKKS slot≠BGV/BFV slot≠CRT slot`
- 三种slot术语不能混用！当前论文如基于CKKS应使用'复数slot'或'CKKS slot'，称'CRT slot'需要定义。BGV/BFV才严格说CRT slot。

## 重要公式
- `BGV/BFV: R_p≅∏_{j=1}^{L}F_{p^d} (CRT分解)` — BGV/BFV槽的CRT来源
- `CKKS: σ(a)=(a(ζ^{sⱼ}))_{j=0}^{N/2-1}∈C^{N/2}` — CKKS槽的复嵌入来源
- `BGV/BFV槽数=φ(m)/d vs CKKS槽数=φ(m)/2` — 槽数差异

## 术语
- **CRT槽** (CRT slot): BGV/BFV中通过中国剩余定理分解得到的槽
- **CKKS槽** (CKKS slot / complex slot): CKKS中通过复嵌入得到的复数槽
- **槽** (slot): 单指令多数据(SIMD)中的独立数据单元

## 关系
- abstract_algebra: CRT vs 复嵌入
- RLWE_lattice: 不同方案的环使用方式
- FHE_technique: BGV/BFV vs CKKS的根本差异
- current_paper: 当前论文必须明确定义使用的槽类型

## RAG建议: 极高优先级，论文写作避错关键卡片
