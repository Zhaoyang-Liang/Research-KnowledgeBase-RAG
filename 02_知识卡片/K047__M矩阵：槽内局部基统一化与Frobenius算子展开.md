# M矩阵：槽内局部基统一化与Frobenius算子展开

**card_id**: K047 | **batch**: Batch-05a2
**source_id**: src_000072, src_000073, src_000074, src_000075, src_000076 | **claim_source_type**: 用户笔记

## Knowledge Points

1. M的任务：将第i个槽中按局部基{ζ_{m,i}^j}表示的元素改写为统一公共基{ζ_m^j}表示。系数m_{i,j}不变，仅改变基的解释。
2. 不同槽的局部基不同是CRT分解的固有性质——各槽对应域的不同嵌入，扮演'ψ'角色的具体元素可能各异。
3. M是slot-wise线性变换：不跨槽混合，在每个槽内独立执行相同的d×d基变换矩阵。
4. 同态实现：将槽内线性变换展开为Frobenius自同构σ的线性组合L(m)=Σ_{v=0}^{d-1}κ(v)·σ^v(m)。Frobenius在每个槽内天然独立作用→适合同态原语。
5. κ(v)求解：在单槽内将目标矩阵M和Frobenius幂矩阵S_v拉直，解线性方程组M=Σκ(v)S_v。可用CRT将局部κ_i(v)拼成全局κ(v)。
6. BSGS/hoisting优化：利用σ^v可组合的特性，用baby-step giant-step加速多Frobenius求值。
