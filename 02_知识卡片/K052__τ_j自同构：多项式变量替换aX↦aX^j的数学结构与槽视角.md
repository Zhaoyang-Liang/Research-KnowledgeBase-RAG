# τ_j自同构：多项式变量替换a(X)↦a(X^j)的数学结构与槽视角

**card_id**: K052 | **batch**: Batch-05a3 | **confidence**: medium
**source_id**: src_000081, src_000083, src_000084 | **claim_source_type**: 用户笔记

## Knowledge Points

1. τ_j定义：τ_j(a(X))=a(X^j)，其中j∈Z_m^*。即把多项式中所有X替换为X^j：若a(X)=Σa_kX^k，则τ_j(a)=Σa_kX^{kj}（再模X^N+1化简）。系数不变，幂次乘j。
2. τ_j≠乘X^j：τ_j(a)=a(X^j)是变量替换，而非X^j·a(X)（乘单项式）。两者完全不同。例如a=1+2X，τ₃(a)=1+2X³而非X³(1+2X)=X³+2X⁴。
3. τ_j是环自同构：保持加法τ_j(a+b)=τ_j(a)+τ_j(b)和乘法τ_j(ab)=τ_j(a)·τ_j(b)。j∈Z_m^*保证X↦X^j对应单位根重排，不破坏cyclotomic环结构。
4. 槽视角：τ_j把根标签h变成hj，在slot层面诱导陪集置换[h]↦[hj]。这是"合法"的，因为(hp^k)j=hj·p^k。
5. Rotation=τ_j的特殊情况：在BGV/BFV中，rotation是在加密数据上执行特定自同构(使用自同构密钥)来循环移动槽的值。τ_{g_i}对应沿第i维的one-dimensional rotation。
6. Rotation≠C2S/S2C：Rotation是槽之间的置换(同态操作)，C2S/S2C是槽↔系数之间的表示转换(线性变换)。两者使用同态原语不同(Rotation用automorphism key，C2S用矩阵/CT分解)。

**交叉引用**: K021(Y-X): 同态Trace与Product的旋转分解（τ_j在同态实现中的应用）；K051(Y-X): Slot编号与slot permutation的群论基础。K052讲τ_j本身的数学结构。

**备注**: 全库术语：C2S=coeff→slot, S2C=slot→coeff。τ_j的rotation≠C2S/S2C的方向转换。
