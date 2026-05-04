# CVP、SVP与Babai最近平面算法的几何解释

**card_id**: K066 | **batch**: Batch-05c | **confidence**: medium | **needs_human_review**: false | **review_policy**: review_when_used_for_paper
**source_id**: src_000025, src_000111, src_000117, src_000113 | **claim_source_type**: 用户笔记 + 标准格密码教材
**备注**: Babai算法的正确性依赖优质基（Gram-Schmidt向量衰减不剧烈）。窃听者用公钥劣质基（长、歪斜），Babai 不会恢复最近格点——这就是陷门的核心。

## Knowledge Points

1. 1. SVP (Shortest Vector Problem)：给定格 L = {Bx : x ∈ Z^n}，找非零格向量 v ∈ L 使得 ‖v‖ = λ₁(L)。SVP 是格中最基本的问题，NP-hard in the worst case。γ-近似 SVP (GapSVP)：区分 λ₁ ≤ d 还是 λ₁ > γd。
2. 2. CVP (Closest Vector Problem)：给定格 L 和目标向量 t ∈ R^n，找格点 v ∈ L 使得 ‖v−t‖ 最小（= dist(t,L)）。CVP 直接对应解密问题：密文 c 接近某个格点 ℓ，私钥找最近的 ℓ'。BDD (Bounded Distance Decoding) 是 CVP 的特化：已知 t 离格 ≤ r << λ₁/2，此时存在唯一最近格点。
3. 3. Babai最近平面算法（Babai Nearest Plane）：输入格基 B = {b₁,...,b_n} 和目标 t。算法从最后一维开始：计算 t_n 在 b_n 方向上的投影，找到最近的整数倍如 c_n = ⌈⟨t,b̃_n⟩/⟨b̃_n,b̃_n⟩⌋。更新 t_{n-1} = t − c_n·b_n，递归到下一维。输出 v = Σ_j c_j·b_j。结果是格中的点（保证正确），且当 B 是优质基（Gram-Schmidt 向量不衰减太快）时 v 接近最近格点。
4. 4. Babai 在 FHE 中的应用：Gentry09 解密使用 Babai 算法：用私钥优质基 B_sk 找离密文 c 最近的格点 ℓ'。当 a = m+2e 足够短时（噪声远小于 λ₁），Babai 恢复正确的 ℓ'=ℓ，解密成功。这对应「强路线」a' = a。
5. 5. Gentry09 两种正确性路线：(a) 强路线 a'=a：Babai 恢复加密时的格点 ℓ'=ℓ，需 a 足够短。(b) 弱路线 a'≡a mod 2：不要求 ℓ'=ℓ，只要求 a'−a ∈ 2Z^m（模2不变）。弱路线对 bit 解密已足够，且更宽松。
6. 6. 与签名的关系：GGH/NTRUSign 用 Babai 做签名（找接近 m 的格点）。失败原因：签名 s = v−m 落在基的平行六面体 P(B) 内，分布形状暴露私钥基。GPV 框架通过高斯采样替代 Babai 舍入，使签名分布变成球面高斯，与基无关，解决该问题。
7. ⚠️ 用你的话总结：CVP = 几何问题，Babai = 几何算法实现，Gentry = 证明 Babai 能正确解密的充分条件（a 足够短）。

**交叉引用**: K065(Y-X): LWE/RLWE困难性（归约到SVP/SIVP）。
