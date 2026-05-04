# BGV Bootstrapping完整流程：6阶段与关键步骤

**card_id**: K038 | **publish_target**: 03_FHE知识库  
**source_id**: src_000059 | **source_type**: lecture_notes  
**document_role**: directly_related | **read_depth**: deep_read  
**needs_human_review**: False  

## Knowledge Points

1. 阶段1：模数提升(Modulus Raising/Lifting) — 将旧密文从低模数q提升到大模数Q，使得同态解密有足够的噪声空间。这是bootstrapping的起点。
2. 阶段2：同态线性解密(Homomorphic Linear Decryption) — 利用加密的sk在同态下计算解密内积⟨ct,sk⟩mod Q。这是最昂贵的同态计算步骤，通常使用BSGS优化。
3. 阶段3：CoeffToSlot — 将解密结果（coefficient表示→noisy plaintext）转换为slot表示。本质是线性变换，用Vandermonde矩阵实现。
4. 阶段4：unpack（fully packed时必需）— 如果原始密文是fully packed，需要将E-slot拆分为标量slot才能进行逐slot的digit extraction。本质是递归的even/odd splitting。
5. 阶段5：digit extraction/modular reduction — 在每个slot上独立执行digit extraction，将noisy plaintext的每个slot中的消息加噪声的高位digit移除。BGV使用digit extraction（逐位提升+减法），BFV使用rounding。这是bootstrapping的正确性核心。
6. 阶段6：repack+SlotToCoeff+输出 — 如果做了unpack则先repack逆向重组；然后SlotToCoeff将slot表示转回coefficient表示；最后输出为新密文（模数已恢复为q）。

**notes**: BGV/BFV bootstrapping与CKKS bootstrapping在结构上有相似性（C2S→中间处理→S2C），但中间处理不同（digit extraction vs EvalMod），环参数也不同（p^e vs plain modulus for CKKS）。
