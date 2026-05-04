# SIMD编码：CRT分解与槽并行计算（BGV/BFV context）

**card_id**: K050 | **batch**: Batch-05a3 | **confidence**: medium
**source_id**: src_000087, src_000088 | **claim_source_type**: 用户笔记

## Knowledge Points

1. BGV/BFV明文空间R_t=Z_t[X]/(Φ_M(X))经CRT分解为直积∏_{i=0}^{L-1}Z_t[X]/(F_i(X))。每个F_i是不可约因子，对应一个槽。
2. SIMD(Single Instruction Multiple Data)：将长度L的明文向量编码为一个多项式a(X)，使得一次环加法/乘法等价于对所有槽的分量并行加/乘。Encode(v₀,...,v_{L-1})=a(X)，Decode(a)=(a mod F₀,...,a mod F_{L-1})。
3. 解码必须基于Φ_M(X)的根(不可约因子的零点)——这保证求值映射在商环上良定义且乘法分量独立。不能随意取其他点，否则破坏环同态结构。
4. 当d=1时每个槽=Z_t(整数模t)；当d>1时每个槽=扩域F_{t^d}(或Galois环)。d=ord_M(t)是t在(Z/MZ)^×中的乘法阶。
5. SIMD编码的核心=CRT插值：给定槽值向量→构造总环元素M(X)，使得M(X)≡消息_i (mod F_i(X))。编码公式：M=Σκ_i(σ_{n,i}(X))H_i(X)G_i(X)。
6. ⚠️ 这是BGV/BFV CRT-based SIMD，≠CKKS complex slot packing。CKKS的packing基于典范嵌入(σ)而非CRT分解。不同FHE方案层的packing机制本质不同，不可混用。

**交叉引用**: K041(Y-X): CKKS complex slot vs BGV/BFV CRT slot区别；K002(Y-X): 分圆多项式CRT分解；K044(Y-X): C2S/S2C通用定义(同为BGV/BFV context)。K050专门讲BGV/BFV的SIMD编码机制。

**备注**: 全库术语约定：C2S=coeff→slot, S2C=slot→coeff。本卡讨论BGV/BFV CRT slot，非CKKS complex slot。
