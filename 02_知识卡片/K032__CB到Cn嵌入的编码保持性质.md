# C^B到C^n的嵌入与编码保持性质

**card_id**: K032 | **publish_target**: 03_FHE知识库 (translation_note, needs_verification — 含原文笔误修正)  
**source_id**: src_000045 | **source_type**: translation_note  
**needs_human_review**: True | **needs_verification**: True  
**cross_reference**: K022 (Y=X^{N/(2n)}代换)  

## Knowledge Points

1. 定义嵌入ι:C^B→C^n为ι(z_1)=(z_1,z_1,…,z_1)重复M=n/B次
2. 关键等式：Ecd_n(ι(z_1))=Ecd_B(z_1)——编码在嵌入下保持一致
3. 代数解释：重复向量构成C^n的子环（在Hadamard积下封闭），ι是环同态
4. 原因是周期为B的向量的IDFT只有指数为M=n/B的倍数的项非零，经过变量替换Y=X^{N/(2n)}后与低维编码完全一致
5. 原文中的"2M-1"求和上限是笔误，应为"2B-1"（即k=0,...,2B-1）
6. 应用：在PaCo中用于将多个多项式的系数打包进一个密文的不同槽

**notes**: 翻译笔记。指出了原文求和范围错误(2M vs 2B)。需核原文验证。 ⚠ 已纠正原文求和范围错误(2M→2B)，需作者人工确认。
