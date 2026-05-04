# CKKS编码的逆变换τ_n^{-1}: 复向量到多项式的完整推导

**card_id**: K030 | **publish_target**: 03_FHE知识库 (translation_note, support_ingest)  
**source_id**: src_000050 | **source_type**: translation_note  
**needs_human_review**: True | **needs_verification**: True  
**cross_reference**: K018 (典范嵌入), K019 (SlotToCoeff正向), K036 (τ_n正向)  

## Knowledge Points

1. 给定z∈C^n，计算p(Y)=τ_n^{-1}(z)的完整步骤：(1)解U_n·w=z得w∈C^n，(2)将w拆为p_0+i·p_1，(3)组装p(Y)=Σ p_0_j·Y^j + Σ p_1_j·Y^{j+n}
2. U_n是Vandermonde矩阵，U_n[i,j]=ζ_{n,i}^j，其中ζ_{n,i}=exp(2πi·5^i/(4n))
3. 由于ζ_{n,i}^n=i，后一半系数的贡献被转化为纯虚数倍，从而实系数多项式p(Y)在复点上求值得到复向量
4. 因此τ_n^{-1}本质上是：IDFT → 分离虚实部 → 分别放入多项式的高低半部分
5. CKKS编码定义为Ecd(z)=⌊Δ·τ_n^{-1}(z)⌋∈Z[X]/(X^N+1)，其中令Y=X^{N/(2n)}

**notes**: 翻译笔记，非原始论文。推导正确但需对照Geelen原文验证Vandermonde矩阵具体形式。 需对照Geelen原文验证U_n矩阵具体形式。
