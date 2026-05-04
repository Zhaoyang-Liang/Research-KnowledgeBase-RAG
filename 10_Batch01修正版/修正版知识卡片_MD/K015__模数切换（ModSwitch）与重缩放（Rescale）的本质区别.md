# 模数切换(ModSwitch)与重缩放(Rescale)的本质区别

- **卡片ID**: K015 | **版本**: v2
- **分类**: FHE总论与基础方案 / 噪声控制机制
- **source_id**: `src_000014`
- **来源**: `01_源文件库/03_FHE总论与基础方案/模数切换和重缩放的区别.md`
- **置信度**: high | **需审核**: False
- **声明**: 我的笔记

## 摘要
ModSwitch(BGV)：c'=⌊(q'/q)·c⌉，任意q'<q，噪声等比例缩小，消息m不变（因消息在低位模t）。Rescale(BFV/CKKS)：c'=⌊(p^{-1})·c⌉，p固定=t(BFV)或编码因子(CKKS)。BFV Rescale同步恢复缩放因子保持消息尺度一致；CKKS Rescale除以Δ≈大素数，降低模数层级，消息尺度变化。核心区别：ModSwitch消息不变，Rescale消息尺度变化。

## 关键知识点

### K015-1: ModSwitch定义
- `c'=⌊(q'/q)·c⌉, 任意q'<q, m不变`
- BGV专属：缩小模数→噪声缩水→消息在低位(m mod t)不变。可任意时刻执行。

### K015-2: BFV Rescale定义
- `c'=⌊(1/t)·c⌉, t固定, Δ恢复为原始值`
- 重缩放固定比例÷t→恢复消息尺度。scale-invariant：消息表示始终保持一致。

### K015-3: CKKS Rescale定义
- `c'=⌊(1/Δ)·c⌉, Δ≈大素数, 模数降级`
- CKKS Rescale消费模数链层级。消息编码尺度变化但约等关系保持。

### K015-4: 消息尺度变化
- `BGV:消息不变; BFV:消息恢复; CKKS:消息尺度改`
- ModSwitch保持m不变；BFV Rescale恢复Δ倍后的消息；CKKS Rescale改变消息尺度。

### K015-5: ⚠ 术语陷阱
- `BFV Rescale≠CKKS Rescale≠BGV ModSwitch`
- 三种操作机制不同但文献中经常混用rescale/scale-down等术语。论文写作必须明确定义。

## 重要公式
- `ModSwitch: c'=⌊(q'/q)·c⌉, m不变` — BGV模数切换
- `BFV Rescale: c'=⌊(1/t)·c⌉, Δ恢复` — BFV重缩放
- `CKKS Rescale: c'=⌊(1/Δ)·c⌉, 模数降级` — CKKS重缩放

## 术语
- **模数切换** (Modulus Switching): BGV中缩放密文到更小模数降噪，消息不变
- **重缩放(BFV)** (Rescale(BFV)): BFV中÷t恢复消息尺度
- **重缩放(CKKS)** (Rescale(CKKS)): CKKS中÷Δ降低模数层级

## 关系
- abstract_algebra: -
- RLWE_lattice: 噪声缩放
- FHE_technique: 噪声控制的核心机制差异
- current_paper: 需明确当前论文使用的噪音控制术语

## RAG建议: 极高优先级，避免术语混淆
