# BGV/CKKS完整同态计算流程

- **卡片ID**: K013 | **版本**: v2
- **分类**: FHE总论与基础方案 / 同态计算流程
- **source_id**: `src_000012`
- **来源**: `01_源文件库/03_FHE总论与基础方案/BGV-CKKS完整流程举例.md`
- **置信度**: high | **需审核**: False
- **声明**: 我的笔记

## 摘要
一次同态乘法流程：(1)张量积→三维密文(c₀d₀,c₀d₁+c₁d₀,c₁d₁)，解密需s²；(2)重线性化→用RLK消除s²项降为二维；(3)噪声评估→检查是否超过阈值B_noise≈q/(2t)；(4)噪声控制→BGV:ModSwitch缩小模数，CKKS:Rescale消耗模数层。每层乘法消耗一层模数。乘法深度由模数链q₀>q₁>…>q_L决定。

## 关键知识点

### K013-1: 张量积→维度膨胀
- `(c₀,c₁)⊗(d₀,d₁)=(c₀d₀,c₀d₁+c₁d₀,c₁d₁)`
- 乘法使密文2维→3维，解密多项式deg 1→2(含s²)。维度膨胀是FHE根本挑战。

### K013-2: 重线性化降维
- `KeySwitch_{s²→s}(c₀,c₁,c₂)→(c'₀,c'₁)`
- RLK=(b_i=a_i·s+t·e_i+B^i·s²,a_i), 按基B分解c₂后消除s²。

### K013-3: 噪声评估与阈值
- `η≤B_noise? 若η>B需降噪，否则解密失败`
- B_noise≈q/(2t)(BGV)或q/(2Δ)(CKKS)。超出则需ModSwitch/Rescale/Bootstrap。

### K013-4: 模数链与乘法深度
- `L层乘法→log₂(q₀/q_L) bits模数预算`
- 初始q₀越大→支持L层越多。模数链是FHE的'燃料量表'。

### K013-5: BGV vs CKKS流程差异
- `BGV: ModSwitch任意时机; CKKS: Rescale每乘一次`
- BGV灵活但需跟踪，CKKS结构化但刚性。

## 重要公式
- `c⊗d=(c₀d₀,c₀d₁+c₁d₀,c₁d₁)` — 密文张量积（乘法）
- `KeySwitch_{s²→s}(c₀,c₁,c₂)=(c₀+Σc₂ⁱ·b_i,c₁−Σc₂ⁱ·a_i)` — 重线性化：s²密钥切换为s

## 术语
- **重线性化** (relinearization): 乘法后高维密文经密钥切换降为二维
- **模数链** (modulus chain): 从q₀到q_L的有序模数序列
- **噪声预算** (noise budget): 当前密文剩余噪声容忍度(单位bits)

## 关系
- abstract_algebra: -
- RLWE_lattice: 噪声增长由RLWE误差分布决定
- FHE_technique: FHE计算核心流程
- current_paper: bootstrapping需重置模数链

## RAG建议: 高优先级，FHE计算流程核心
