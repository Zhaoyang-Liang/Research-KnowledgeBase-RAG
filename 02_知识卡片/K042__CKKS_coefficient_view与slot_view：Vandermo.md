# CKKS coefficient view与slot view：Vandermonde变换与两种表示

**card_id**: K042 | **batch**: Batch-05a1 | **publish_target**: 03_FHE知识库
**source_id**: src_000066 | **claim_source_type**: 用户笔记

## Knowledge Points

1. coefficient view：多项式m(X)=Σa_iX^i，用系数向量(a_0,...,a_{N-1})表示。系数为整数（最终可加密形式）。
2. slot view：m在N个单位根处的取值(m(ζ^{s_0}),...,m(ζ^{s_{N-1}}))。值为复数。这是解码方向看到的表示。
3. 两种表示之间的变换：slot = U_n·coeff（Vandermonde矩阵乘法），coeff = U_n^{-1}·slot（逆Vandermonde）。
4. U_n是N×N Vandermonde矩阵：(U_n)_{i,j}=ζ^{s_i·s_j}。已知在互不相同单位根处求值→U_n可逆→变换是双射。
5. CKKS编码从slot view出发：用户指定想要的slot值→通过σ^{-1}(=U_n^{-1})反推coefficient→得到多项式。
6. CKKS解码回到slot view：m(X)→σ→(m(ζ^{s_j}))→取前N/2个独立槽→≈原始消息。

**交叉引用**: K018(Y-X): τ_n和U_n矩阵；K036(Y-X): τ_n正向推导(翻译笔记)。K042解释coefficient/slot两种表示之间的Vandermonde变换关系。
