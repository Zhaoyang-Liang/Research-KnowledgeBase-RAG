# 模数切换(ModSwitch)与密钥切换(KeySwitch)的本质区别

**card_id**: K060 | **batch**: Batch-05b1 | **confidence**: medium | **needs_human_review**: false | **review_policy**: review_when_used_for_paper
**source_id**: src_000105 | **claim_source_type**: 用户笔记

## Knowledge Points

1. ModSwitch：改变模数q→q'，密钥s不变。c'=⌊(q'/q)c⌉。目的=降低噪声绝对值。不引入新噪声（仅缩放舍入）。BGV中大量使用。
2. KeySwitch：改变密钥s→s'，模数q不变。c'=SwitchKey(c,KSK)。目的=换密钥/重线性化。引入少量额外噪声（KSK解密+分解）。所有FHE方案通用。
3. 判定表格：(1)改q→modswitch；(2)改s→keyswitch；(3)降噪→modswitch；(4)降维→keyswitch(rlinearization)；(5)授权→keyswitch。
4. 组合使用：BGV乘法后先KeySwitch（重线性化s²→s）再ModSwitch（降噪）；CKKS乘法后先KeySwitch再Rescale。实践中两者总是成对出现，但目的不同。
5. 比喻：ModSwitch=换尺子（量程变小误差等比例缩小），KeySwitch=换锁+钥匙（锁体不变钥匙变了）。

**交叉引用**: K058(Y-X): 计算流程中ModSwitch时序；K059(Y-X): KeySwitch三场景；K061(Y-X): ModSwitch vs Rescale（ModSwitch的不同版本）。

**备注**: 全库术语：C2S=coeff→slot, S2C=slot→coeff。
