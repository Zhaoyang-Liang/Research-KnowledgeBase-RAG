# BGV与BFV的核心区别：消息低位vs高位

- **卡片ID**: K017 | **版本**: v2
- **分类**: FHE总论与基础方案 / BGV vs BFV对比
- **source_id**: `src_000016`
- **来源**: `01_源文件库/03_FHE总论与基础方案/BFV与BGV的高低位明文.md`
- **置信度**: high | **需审核**: False
- **声明**: 我的笔记

## 摘要
BGV与BFV的根本区别在于明文消息在密文中的位置。BGV：消息在低位(c₀=a·s+t·e+m)，模t后噪声e被消掉只剩m。模数切换缩小模数q→q'，噪声绝对值等比例缩小，消息不变。BFV：消息放大Δ=⌊q/t⌋倍放在高位(c₀=a·s+Δ·m+e)，重缩放÷t恢复尺度。用比喻：BGV像在尺子低刻度写数字→噪声大了换短尺子；BFV像放大后在尺子高刻度写→噪声大了把尺子和物体同时缩小。

## 关键知识点

### K017-1: BGV消息在低位
- `c₀=a·s+t·e+m, Dec: [c₀+c₁·s]_q mod t=m`
- 解密时先模q消去a·s，再模t消去t·e（若|t·e|<q/2）。消息z在低位不受模数大小影响。

### K017-2: BFV消息在高位
- `c₀=a·s+Δ·m+e, Δ=⌊q/t⌋, Dec: ⌊(t/q)·[c₀+c₁·s]_q⌉`
- 消息被放大Δ倍占据密文系数的高位。解密需除以Δ恢复消息，精度依赖q足够大。

### K017-3: ModSwitch(BGV) vs Rescale(BFV)
- `BGV:⌊(q'/q)·c⌉(消息不变); BFV:⌊(1/t)·c⌉(尺度恢复)`
- ModSwitch任意比例缩小模数→灵活；Rescale固定比例÷t→结构化。BGV更灵活但参数跟踪更复杂。

### K017-4: 比喻理解
- `BGV=低刻度换尺子; BFV=放大后同比例缩小`
- BGV:改变测量工具(模数); BFV:改变被测量对象(消息放大因子)和测量工具同步缩小。

### K017-5: 对Bootstrapping的影响
- `BGV:需精确整数模约化; BFV:可用模数切换降噪换模数`
- BGV bootstrapping需精确整数电路→昂贵。BFV可利用尺度不变性简化某些操作。

## 重要公式
- `BGV Dec: m=[c₀+c₁·s]_q mod t` — BGV解密：先模q后模t
- `BFV Dec: m=⌊(t/q)·[c₀+c₁·s]_q⌉` — BFV解密：缩放恢复

## 术语
- **低位明文** (low-bit plaintext): BGV：消息放在系数最低位，模t后即消息
- **高位明文** (high-bit plaintext): BFV：消息放大Δ倍占据系数高位
- **Δ** (Delta): BFV的缩放因子=⌊q/t⌋

## 关系
- abstract_algebra: -
- RLWE_lattice: 加密结构差异
- FHE_technique: ModSwitch vs Rescale的技术根源
- current_paper: 当前基于CKKS，消息尺度变化需明确定义

## RAG建议: 高优先级，避免BGV/BFV混淆
