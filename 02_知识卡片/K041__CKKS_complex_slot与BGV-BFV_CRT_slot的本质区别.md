# CKKS complex slot与BGV/BFV CRT slot的本质区别

**card_id**: K041 | **batch**: Batch-05a1 | **publish_target**: 03_FHE知识库
**source_id**: src_000066, src_000059 | **claim_source_type**: 用户笔记+讲义综合

## Knowledge Points

1. CKKS slot：复数槽。值空间=C，每个槽存储一个复数。槽个数=N/2（因共轭对称）。基础结构：典范嵌入σ:R→C^{N/2}。
2. BGV/BFV slot：CRT槽。值空间=F_{p^d}（有限域），每个槽存储一个mod p^e的元素。槽个数=N/d(d=ord_m(p))。基础结构：CRT分解R_{p^e}≅ΠF_{p^d}。
3. 代数基础不同：CKKS基于典范嵌入（分析/几何结构），BGV/BFV基于中国剩余定理（代数结构）。
4. 精确性不同：CKKS为近似编码（rounding引入误差，Δ控制精度）；BGV/BFV为精确编码（CRT分解无信息损失）。
5. CKKS slot值≡'多项式在单位根处的取值'；BGV/BFV slot值≡'多项式模各不可约因子的余数'。
6. 两者都是将大环分解为较小分量，但分解机制和分量含义完全不同。CRT slot≠complex slot，术语不可互换。

**交叉引用**: K002(Y-X): 分圆多项式CRT分解；K037(Y-X): BGV/BFV bootstrapping讲义（BGV/BFV侧slot结构）。
