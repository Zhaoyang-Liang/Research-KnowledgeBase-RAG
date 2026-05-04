# FHE同态计算完整流程：乘法→重线性化→模数切换/重缩放→解密

**card_id**: K058 | **batch**: Batch-05b1 | **confidence**: medium
**source_id**: src_000103 | **claim_source_type**: 用户笔记

## Knowledge Points

1. 通用流程：加密→同态运算→噪声控制→解密。每次乘法使噪声≈增长为平方级别，必须通过噪声控制操作才能继续计算。加法噪声线性增长，一般不需要立即控制。
2. Step 1 乘法（张量积）：c×=(c₀d₀, c₀d₁+c₁d₀, c₁d₁)——维度3，有效密钥s²，噪声大幅增加。
3. Step 2 重线性化（密钥切换s²→s）：使用RLK将3维密文降回2维。重线性化=密钥切换的特例。引入少量额外噪声（远<乘法噪声）。每次乘法后必须立即执行。
4. Step 3 噪声控制：BGV=模数切换（任意缩比q→q'，降噪声绝对值）；BFV/CKKS=重缩放（固定缩比1/t或1/p，保持缩放因子一致）。
5. 加法流程：a+b直接对应分量相加，噪声线性叠加，维度不变，不需要重线性化，一般不需要立即噪声控制。
6. 自举时机：当模数链耗尽（无法继续降模数）且噪声接近阈值时，触发Bootstrapping重置噪声和模数。当前论文聚焦CKKS bootstrapping中的coeff→slot→evalMod→slot→coeff管线。

**交叉引用**: K008(Y-X): FHE方案对比中的各方案噪声控制；K056(Y-X): 方案总览（各方案加密/解密公式）；K059(Y-X): 密钥切换三场景（重线性化详解）；K061(Y-X): ModSwitch vs Rescale。

**备注**: 全库术语：C2S=coeff→slot, S2C=slot→coeff。current_paper_relation提到CKKS bootstrapping管线但不展开。
