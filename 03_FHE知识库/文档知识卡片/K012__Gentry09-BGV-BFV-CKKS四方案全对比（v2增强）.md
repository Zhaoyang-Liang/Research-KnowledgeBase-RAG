# Gentry09/BGV/BFV/CKKS四方案全对比(v2增强)

- **卡片ID**: K012 | **版本**: v2
- **分类**: FHE总论与基础方案 / FHE方案对比
- **source_id**: `src_000011`
- **来源**: `01_源文件库/03_FHE总论与基础方案/gentry-BGFV-CKKS.md`
- **置信度**: high | **需审核**: False
- **声明**: 我的笔记+原文精读

## 摘要
四大FHE方案核心区别：BGV明文在低位+ModSwitch任意比例降噪；BFV明文放大Δ=⌊q/t⌋在高位+Rescale固定比例÷t；CKKS复数近似编码+Rescale。BGV与BFV均为精确整数但噪声控制方式不同；CKKS牺牲精度换取效率（bootstrapping更高效）。注意BFV和CKKS的Rescale名称相同但机制不同：BFV除以t恢复尺度，CKKS除以编码因子Δ≈大素数降低模数层。

## 关键知识点

### K012-1: LWE/RLWE加密统一形式
- `LWE:(⟨a,s⟩+e+m,−a); RLWE:(a·s+e+m,−a) in R_q`
- 加解密结构同，差异在m位置和噪声控制。

### K012-2: BGV噪声控制
- `c←⌊(q'/q)·c⌉, 任意q'<q, 消息m不变`
- 模数切换缩小模数→噪声比例降低→消息在低位(m mod t)不变。

### K012-3: BFV尺不变性
- `Δ=⌊q/t⌋, c'←⌊(t/q)·c⌉(重缩放)`
- 消息放大Δ倍写高位，重缩放÷t恢复原始尺度。尺不变：消息表示一致。

### K012-4: CKKS近似复数编码
- `m=⌊Δ·σ(z)⌉, z∈C^{N/2}, Decode≈z/Δ`
- 典范嵌入编码复数向量为多项式，放大Δ取整。解码≈z/Δ。

### K012-5: Bootstrapping比较
- `CKKS用sin多项式近似取模→更低深度`
- CKKS允许近似→bootstrapping更高效。BGV/BFV需精确整数电路→更昂贵。

### K012-6: ⚠ BFV vs CKKS Rescale差异
- `BFV: ÷t恢复尺度; CKKS: ÷Δ(素数)降模数层`
- 同名不同义！BFV Rescale保持消息尺度不变，CKKS Rescale消费模数链层级。

## 重要公式
- `BGV: c=(a·s+t·e+m,−a), Dec: m=[c₀+c₁·s]_q mod t` — BGV加解密，明文在低位
- `BFV: c=(a·s+Δ·m+e,−a), Δ=⌊q/t⌋` — BFV加密，明文放大Δ在高位
- `CKKS: Encode(z,Δ)=⌊Δ·σ^{-1}(z)⌉` — CKKS编码：逆典范嵌入后放大取整

## 术语
- **模数切换** (Modulus Switching (ModSwitch)): BGV中缩放密文至更小模数降噪，消息不变
- **重缩放** (Rescale/Scale-down): BFV: ÷t恢复尺度; CKKS: ÷Δ降模数层
- **尺不变** (scale-invariant): BFV特征：消息表示在计算中保持一致尺度（由Δ固定）

## 关系
- abstract_algebra: RLWE环选择
- RLWE_lattice: 所有方案基于RLWE
- FHE_technique: 三方案噪声控制是FHE核心技术
- current_paper: 当前论文基于CKKS bootstrapping改进

## RAG建议: 最高优先级，FHE总览卡片
