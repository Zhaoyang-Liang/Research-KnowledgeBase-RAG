# 模数切换(ModSwitch)与重缩放(Rescaling)的本质区别

**card_id**: K061 | ⚠️ **needs_human_review**: true

**batch**: Batch-05b1 | **confidence**: medium
**source_id**: src_000106 | **claim_source_type**: 用户笔记

## Knowledge Points

1. ModSwitch（BGV风格）：c'=⌊(q'/q)c⌉，缩比q'/q任意，消息不变（因为m在低位mod t不变），主要用于降噪声绝对值。
2. Rescaling（BFV风格）：c'=⌊(1/t)c⌉，缩比固定=1/t，消息Δm的尺度从Δ²恢复为Δ（因为Δ=⌊q/t⌋），新模数=q/t。保持缩放因子一致。
3. Rescaling（CKKS风格）：c'=⌊(1/p)c⌉，缩比固定=1/p，p=编码缩放因子。乘法后p²→p，保持编码精度。
4. 本质区别：(1)目的——ModSwitch降噪不改变消息尺度，Rescale保持缩放因子同步修复。(2)缩比——ModSwitch任意，Rescale固定。(3)方案——ModSwitch for BGV，Rescale for BFV/CKKS。(4)消息含义——ModSwitch后消息不变，Rescale后消息缩放因子归一。
5. 为什么容易混淆：两者在代码实现上几乎一样（都是对密文系数做整数缩放舍入），但数学语义不同。尤其是BFV中，Rescale本质也是一种ModSwitch，但带有缩放因子同步的含义。
6. ⚠️ CKKS中Rescale还有精度管理含义（缩比越大丢失精度越多）。当前论文CKKS bootstrapping中频繁使用Rescale——需区分是精度控制还是噪声控制。

**交叉引用**: K057(Y-X): BGV vs BFV高低位；K058(Y-X): 计算流程中Rescale/ModSwitch的触发时机；K060(Y-X): ModSwitch vs KeySwitch。

⚠️ **Caveat**: 本卡 ModSwitch vs Rescale 跨 BGV/BFV/CKKS context。CKKS 的 Rescale 含精度管理含义，≠BGV 的 ModSwitch。引用时必须在明确方案 context 下使用。

**备注**: 全库术语：C2S=coeff→slot, S2C=slot→coeff。CKKS context：Rescale含义需结合当前论文判断。
