# 典范嵌入正向τ_n(p(Y)): 多项式到复向量的完整推导

**card_id**: K036 | **publish_target**: 03_FHE知识库 (translation_note, support_ingest)  
**source_id**: src_000051 | **source_type**: translation_note  
**needs_human_review**: True | **needs_verification**: True  
**cross_reference**: K018 (CKKS典范嵌入), K030 (τ_n^{-1}逆变换)  

## Knowledge Points

1. τ_n:R[Y]/(Y^{2n}+1)→C^n定义为τ_n(p(Y))=(p(ζ_{n,i}))_{i∈[0,n)}
2. 正向推导：(1)拆p(Y)系数为p_0(前n)和p_1(后n)，(2)利用ζ_{n,i}^n=i简化，(3)定义复向量w=p_0+i·p_1，(4)τ_n(p)=U_n·w
3. 这就是"多项式→复槽向量"方向的典范嵌入
4. 与τ_n^{-1}(逆变换)构成对：τ_n^{-1}是IDFT+组装，τ_n是DFT+求值
5. τ_n和τ_n^{-1}互逆：τ_n∘τ_n^{-1}(z)=z，且τ_n^{-1}∘τ_n(p)=p

**notes**: 翻译笔记。τ_n正向是CKKS编码理论基础，confidence=high。
