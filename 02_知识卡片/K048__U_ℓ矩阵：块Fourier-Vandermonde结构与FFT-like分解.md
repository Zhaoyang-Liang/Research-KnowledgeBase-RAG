# U_ℓ矩阵：块Fourier/Vandermonde结构与FFT-like分解

**card_id**: K048 | **batch**: Batch-05a2
**source_id**: src_000077, src_000078, src_000079 | **claim_source_type**: 用户笔记

## Knowledge Points

1. U_ℓ是ℓ×ℓ块矩阵，第(r,c)块=ω^{rc}·I_d=ζ_m^{drc}·I_d(d×d单位块)。ω=ζ_m^d是外层mixing根。
2. U_ℓ天然像Fourier矩阵的原因：它执行roots-of-unity上的值↔系数重组——这类变换的矩阵形式本来就是Vandermonde/DFT/block-FFT结构。
3. U_ℓ不是任意dense矩阵——它有Fourier/Vandermonde结构，可做Cooley-Tukey式递归分解为logℓ层稀疏矩阵乘积，复杂度从O(ℓ²)降为O(ℓlogℓ)。
4. 与K018/K019中U_n的关系：U_n是CKKS中N×N Vandermonde（τ_n矩阵），U_ℓ是BGV/BFV中C2S管线的跨槽混合矩阵。两者均属Fourier型矩阵，但上下文不同。
5. 论文第4节将U_ℓ在不同情形下递归分解：p≡1(mod4)用块矩阵分解；p≡3(mod4)更接近CKKS的FFT矩阵。
6. U_ℓ的核心洞察：经过M统一坐标系后，跨槽重排不再是任意线性变换，而是与roots-of-unity的Fourier结构强相关——这正是可加速的关键。

**交叉引用**: K018(Y-X): CKKS的U_n Vandermonde矩阵；K019(Y-X): CKKS的C2S(U_n·coeff)。K048覆盖BGV/BFV中U_ℓ的块FFT结构。
