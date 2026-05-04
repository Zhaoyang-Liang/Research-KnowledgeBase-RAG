# BGV vs BFV本质区别：消息低位高位、模数切换重缩放

**card_id**: K057 | ⚠️ **needs_human_review**: true

**batch**: Batch-05b1 | **confidence**: medium
**source_id**: src_000102 | **claim_source_type**: 用户笔记

## Knowledge Points

1. 加密公式区别：BGV c₀=a·s+t·e+m（消息在低位，t乘噪声），BFV c₀=a·s+Δ·m+e（消息在高位，Δ=⌊q/t⌋乘消息）。BGV解密m=[c₀+c₁·s]_q mod t；BFV解密m=⌊(t/q)·(c₀+c₁·s)_q⌉。
2. 消息位置：BGV消息在系数低位（模t提取），BFV消息在系数高位（乘以Δ抬高到q的高比特）。这导致噪声与消息的相对位置不同，影响噪声管理和参数选择。
3. 噪声控制：BGV用模数切换（ModSwitch）——任意比例缩小q→q'，噪声绝对值等比例缩小；BFV用重缩放（Rescale）——固定比例缩小t倍（q→q/t），同时恢复缩放因子Δ。
4. 模数链：BGV可自由选择q链（需要精心设计），BFV的q最好是t的幂次（更自然）。实现上BFV更简单（尺度不变），BGV更灵活（适合细粒度噪声控制）。
5. 缩放因子：BGV无Δ（消息直接嵌在低位），BFV有Δ=⌊q/t⌋（消息被放大后嵌入）。CKKS也有缩放因子但含义不同：CKKS的p是编码精度控制，BFV的Δ是环格结构需要。
6. ⚠️ 本卡重点：BGV和BFV都是精确整数FHE方案，但内部表示策略不同。当前论文以CKKS为context——CKKS的消息编码方式和噪声控制策略与两者都不同。

**交叉引用**: K008(Y-X): FHE方案对比总结；K056(Y-X): FHE方案总览（含CKKS和Gentry09）；K060(Y-X): ModSwitch vs KeySwitch区别；K061: ModSwitch vs Rescale。

⚠️ **Caveat**: 本卡中的「模数切换」和「重缩放」在 BGV/BFV context 中定义。CKKS 的重缩放（rescale by p）虽然数学形式上与 BFV 类似，但 p 是编码精度控制因子，非明文模数 t。当前论文在 CKKS context，引用本卡时需注意：BGV 的 ModSwitch ≠ CKKS 的 Rescale，BFV 的 Rescale=t 除法 ≠ CKKS 的 Rescale=p 除法。

**备注**: 全库术语：C2S=coeff→slot, S2C=slot→coeff。本卡BGV/BFV context，CRT slot。
