# KeySwitch噪声的来源、分解与上界

**card_id**: K064 | **batch**: Batch-05b2-pre | **confidence**: medium | **needs_human_review**: false | **review_policy**: review_when_used_for_paper
**source_id**: src_000015, src_000013, src_000104 | **claim_source_type**: 用户笔记总结标准RLWE/KS噪声知识
**备注**: 全库术语：C2S=coeff→slot, S2C=slot→coeff。本卡为通用KeySwitch噪声知识。β, B_KS数值见参数文件。

## Knowledge Points

1. 1. KeySwitch 的数学定义：Input = ct (in s_old, modulus q), KSK = {Enc_{s_new}(s_old · g^i)}_i。Output = ct' (in s_new, modulus q)。过程：(a) 将 ct.a 按基 g 分解为 digits a^{(i)}，(b) 计算 ct' = ct.b · (1,0) + Σ_i a^{(i)} ⊙ KSK_i，(c) KeySwitch noise = Σ_i a^{(i)} · e_KSK,i，其中 e_KSK,i 是 KSK 中的加密噪声。
2. 2. 噪声的三大来源：(a) KSK加密噪声——每个 KSK_i 本身是 RLWE 密文 Enc_{s_new}(s_old·g^i) + e_enc,i，解密噪声 e_enc,i ~ D_σ，(b) Digit decomposition 近似——将 s_old 分解为 s_old ≈ Σ_i s^{(i)}g^i 时引入逼近误差（如 g=2^L, β=⌈log_g(qp)⌉ 段），(c) ModDown——如果 KS 过程含模数下降（QP→Q），ModDown 引入额外的舍入噪声。
3. 3. 噪声上界（保守形式）：设每次 KS 使用 β 个 KSK 分量，每个 KSK 的加密噪声界 B_KS。则单次 KS 噪声保守上界：B_KS,total ≤ β · B_KS。其中 β = digit decomposition 段数，B_KS = |e_enc|_∞ 的上界。概率界：Pr[|e_KS|_∞ > β·B_KS] ≤ β · δ_enc（union bound）。
4. 4. 噪声上界（tight subgaussian 形式）：假设 e_enc 和 a^{(i)} 独立且 subgaussian，则 KS 噪声的标准差 ≈ √(β) · σ_enc · |a^{(i)}|_2。这在参数大时比保守上界紧得多。实践中常用 √β·B_KS 而非 β·B_KS。
5. 5. 重线性化 = KeySwitch 的特殊形式：Relinearize(ct=(c₀,c₁,c₂)) 本质是 KS(s²→s)。需提供 RLK = {Enc_s(s² · g^i)}_i。其他 KS 场景（密钥授权 s₁→s₂, bootstrapping s_old→s_fresh）使用相同的噪声模型，仅 KSK 不同。
6. 6. 当前论文中的 KS 出现点：(a) Relinearization（EvalMod 内的乘法后）,(b) C2S/S2C 内的 rotation key switching,(c) 如果 Split/Merge 带 key switching。KS 是 bootstrapping 中除 EvalMod 外最主要的噪声来源。
7. ⚠️ Caveat：本卡给出 KS 噪声的标准形式（从 RLWE 文献和 CKKS 2017 中抽象），但具体的 β, B_KS, δ_enc 数值依赖参数。当前论文中 KS 噪声的叠加方式参见 src_000038 §4.5。不混用 BGV ModSwitch 的噪声概念。

**交叉引用**: K059(Y-X): KeySwitch三场景（重线性化/授权/适配）——互补（K059讲「是什么」，K064讲「噪声怎么来」）；K060(Y-X): ModSwitch vs KeySwitch——ModSwitch不引入噪声，KeySwitch引入噪声；K058(Y-X): 计算流程中KeySwitch出现的时机；K063(Y-X): CKKS噪声来源总览（KeySwitch作为来源之一）。
