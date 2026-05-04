# CKKS中的Trace和Product操作及同态内积计算

**card_id**: K034 | **publish_target**: 03_FHE知识库 (translation_note, support_ingest)  
**source_id**: src_000048 | **source_type**: translation_note  
**needs_human_review**: True | **needs_verification**: True  
**cross_reference**: K021 (同态Trace与Product操作的旋转分解)  

## Knowledge Points

1. Trace操作Tr_{n→B}：将长向量按分组"水平求和"成短向量，每组n/B个元素求和
2. Product操作Pr_{n→B}: 将组内元素"水平求积"
3. Trace是"免费"操作（仅加法，不消耗乘法级别），Product消耗对数级别
4. 同态内积标准做法：(1)逐分量密文乘法，(2)对结果反复Trace至1维得到总内积
5. 这两种操作是数域扩张中Trace和Norm在槽表示下的类比
6. 自举场景：将多个多项式系数"折叠"成更少的槽

**notes**: 翻译笔记。Trace和Product是CKKS标准操作，confidence设为high因为这是教科书级内容。
