# 重线性化与位分解（bit decomposition）的完整机制

**card_id**: K068 | **batch**: Batch-05c | **confidence**: medium | **needs_human_review**: false | **review_policy**: review_when_used_for_paper
**source_id**: src_000112, src_000013, src_000015 | **claim_source_type**: 用户笔记（bit_decom.md + 密钥切换笔记）总结标准KeySwitch机制
**备注**: 全库术语：BitDecomp=位分解函数。Powersof2=2的幂。FHE文献中用BitDecomp+Powersof2的技术已成为KeySwitch的标准写法。BGV论文中有详细的公式。

## Knowledge Points

1. 1. 为什么需要重线性化：乘法 c×c' 产生三项密文 (α₀, α₁, α₂) 对应密钥 (1, s, s²)。继续乘 → s³, s⁴... 密文维数爆炸。重线性化将 (α₀, α₁, α₂) 压回 (β₀, β₁)，恢复到 (1, s) 下解密。核心等价于 KeySwitch(s²→s)。
2. 2. 朴素直觉：预发布 evk ≈ 's² 的加密'（γ₀+γ₁s ≈ s²）。用 α₂(γ₀+γ₁s) 替代 α₂s² → 噪声被 α₂ 放大，不可行。α₂ 是模 q 大数。
3. 3. 位分解的动机：避免直接乘 α₂（大数→噪声爆炸）。把 α₂ 按基 T（通常 T=2 或 2^w）分解为 digits：α₂ = Σ_{i=0}^{ℓ-1} α_{2,i} T^i（每个 α_{2,i} 小）。ℓ = ⌈log_T q⌉。然后发布一组密钥 rlk_i ≈ T^i·s² 的加密（而不是单一 s² 加密），用小 digit α_{2,i} 分别去乘。
4. 4. 完整过程：(a) BitDecomp(α₂)：α₂ → (α_{2,0}, α_{2,1}, ..., α_{2,ℓ-1})。(b) Σ_i α_{2,i}·rlk_i ≈ Σ_i α_{2,i}·T^i·s² + noise = α₂s² + noise。(c) β₀ = α₀ + Σ_i α_{2,i}·b_i，β₁ = α₁ + Σ_i α_{2,i}·a_i → (β₀, β₁) 是线性化后密文。
5. 5. 噪声分析：噪声增加 δ·ℓ·B_rlk（保守）或 √ℓ·B_rlk（subgaussian）。ℓ = 分解段数，B_rlk = rlk 中加密噪声界。位基 T 越大→ℓ 越小（分解段少）→噪声 Σ 项少→但每个 α_{2,i} 分布范围大→权衡点 T=2^w (w≈8-16)。
6. 6. KeySwitch 通用性：上述过程适用于任意 s_old→s_new 切换（不限于 s²→s）。KSK_i ≈ T^i·s_old 在 s_new 下的加密。重线性化 = KeySwitch 的 s²→s 特例。Rotation = KeySwitch 的 s(X^k)→s(X) 特例。C2S/S2C 中的旋转也通过此机制实现。
7. ⚠️ 用你的话总结：位分解 = 把大系数拆成小 digit × T^i，用多个小乘替代一个大乘，把噪声增长从 O(q·B) 降到 O(ℓ·B)。这是所有 KeySwitch/Relinearization 的核心技术。

**交叉引用**: K059(Y-X): 密钥切换三种场景（Relinearization是s²→s的特例）；K064(Y-X): KeySwitch噪声来源（digit decomposition是噪声源之一）；K058(Y-X): 同态计算流程（乘后Relin在flow中位置）。
