# PaCo算法结构：参数依赖图与索引系统

**card_id**: K033 | **publish_target**: 03_FHE知识库 (translation_note, needs_verification — 源为翻译笔记非PaCo原文)  
**source_id**: src_000047 | **source_type**: translation_note  
**needs_human_review**: True | **needs_verification**: True  
**cross_reference**: K024 (PaCo论文全文翻译摘要)  

## Knowledge Points

1. PaCo使用结构化密钥：h组，每组4条候选支路，每组恰好选择1条(d_α=1)
2. 参数依赖链：N,h→B=N/(4h)，j∈[0,h)→t_j→λ_j=t_j·h+j→d_{λ_j}=1
3. u_α独立采样于[0,B)，与d_α共同决定密钥s=Σ d_α·X^{u_α·4h+α}
4. a_α由密文公开构造：a_0=ct_0+ct_1，a_α=X^α·ct_1 (α≥1)
5. 解密式压缩：m=Σ d_α·X^{u_α·4h}·a_α→Σ X^{u_{λ_j}·4h}·a_{λ_j}（每组只留激活支路）
6. d_α,u_α,a_α三套量彼此独立：d选路、u移位、a载公开内容

**notes**: 翻译笔记。PaCo算法结构的依赖图梳理。需核PaCo原文验证索引定义。 ⚠ 来源为翻译笔记(algo1.md)，非PaCo原文。需对照PaCo原文验证所有索引定义和依赖关系。
