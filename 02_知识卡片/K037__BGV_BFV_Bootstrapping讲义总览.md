# BGV/BFV Bootstrapping讲义总览：229页系统教程

**card_id**: K037 | **publish_target**: 03_FHE知识库  
**source_id**: src_000059 | **source_type**: lecture_notes  
**document_role**: directly_related | **read_depth**: deep_read  
**needs_human_review**: False  

## Knowledge Points

1. 作者：梁朝阳，2026年4月18日，229页 + 52页beamer。覆盖BGV/BFV bootstrapping所需全部代数基础与算法流程。
2. 三部分结构：Part I 代数与数论基础（分圆环、Frobenius轨道、CRT、slot几何、automorphism/rotation）；Part II 主流程与表示变换（解密关系、sparse/fully packed、CoeffToSlot/SlotToCoeff、unpack/repack）；Part III 优化与复杂度（digit extraction、BSGS、hoisting、full/thin bootstrap）。
3. BGV/BFV bootstrapping核心流程：旧密文→同态解密→CoeffToSlot→digit extraction(或BFV rounding)→SlotToCoeff→新密文。关键区分：BGV用digit extraction移除噪声的高位digit，BFV用rounding。
4. BGV/BFV与CKKS bootstrapping的本质区别：BGV/BFV工作在finite ring arithmetic（目标是消除一个精确的模t消息上的噪声）；CKKS工作在approximate arithmetic（目标是把近似复数上的缩放因子恢复）。因此两者的C2S/S2C/中间步骤的数学本质不同。
5. 讲义包括完整伪代码附录（参数生成、编码解码、加解密、key switching、weighted rotations、BSGS、hoisting、端到端bootstrap）。

**notes**: 作者标注为南开大学密码科学与技术。讲义质量极高，是BGV/BFV领域系统性的中文资料。建议人工阅读第2.9节（与CKKS对比）和第13章（复杂度）。
