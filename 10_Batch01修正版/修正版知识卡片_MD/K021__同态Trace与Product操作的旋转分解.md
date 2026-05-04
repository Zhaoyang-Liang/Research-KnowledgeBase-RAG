# 同态Trace与Product操作的旋转分解

- **卡片ID**: K021 | **版本**: v2
- **分类**: FHE核心技术模块 / CKKS/Trace与旋转操作
- **source_id**: `src_000021`
- **来源**: `01_源文件库/04_FHE核心技术模块/trace操作.md`
- **置信度**: high | **需审核**: False
- **声明**: 我的笔记

## 摘要
FHE中的同态Trace=旋转-加法链。Tr(c)=c+Rot₁(c)+Rot₂(c)+…+Rot_{d-1}(c)，将密文在d个旋转位置上的值求和。实现通过重复加倍技巧：Rot_{2^k}(c)用之前的旋转累积，复杂度O(log d)而非O(d)。Product操作Prod(c)=c⊙Rot₁(c)⊙…⊙Rot_{d-1}(c)类似但用乘法。两者都在bootstrapping和field switching中关键。

## 关键知识点

### K021-1: Trace的同态实现
- `Tr(c)=c+Rot₁(c)+…+Rot_{d-1}(c)`
- 密文在d个自同构下的像的加法链。每个Rot_k用自同构密钥。

### K021-2: 重复加倍技巧
- `Rot_{2^k}=Rot_{2^{k-1}}∘Rot_{2^{k-1}}, 累积加速`
- 只用log d次旋转累积而非d次。例：d=8时Rot₄=Rot₂∘Rot₂。

### K021-3: Product操作
- `Prod(c)=c⊙Rot₁(c)⊙…⊙Rot_{d-1}(c)`
- 密文乘积链，应用于范数计算和某些同态聚合。复杂度同Trace。

### K021-4: Trace在field switching中的角色
- `Tr_{n→B}: 将N环密文投影到B子环`
- Trace降维是field switching降维的关键：大环→子环。

### K021-5: 旋转密钥的存储与复用
- `预计算{KSK_{s(X)→s(X^k)}: k in generator set of Z_m^×/<p>}`
- 只需预计算生成元集合的旋转密钥，其他旋转通过组合获得。

## 重要公式
- `Tr(c)=c+Σ_{k=1}^{d-1}Rot_k(c)` — 同态迹的旋转-加法链
- `Rot_{2^k}=Rot_{2^{k-1}}∘Rot_{2^{k-1}}` — 重复加倍技巧

## 术语
- **同态迹** (homomorphic trace): 通过旋转和加法同态计算的迹操作
- **旋转** (rotation/automorphism): FHE中X→X^k的自同构操作
- **重复加倍** (repeated doubling): 通过自复合累积旋转加速的技巧

## 关系
- abstract_algebra: 迹理论
- RLWE_lattice: 自同构密钥
- FHE_technique: bootstrapping/field switching核心
- current_paper: field switching的Tr_{n→B}实现

## RAG建议: 高优先级，field switching关键技术
