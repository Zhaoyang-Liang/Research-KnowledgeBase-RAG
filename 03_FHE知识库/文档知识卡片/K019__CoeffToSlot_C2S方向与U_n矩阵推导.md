# CoeffToSlot (C2S) 方向与 U_n 矩阵推导

**card_id**: K019 | **batch**: Batch-01
**normalized_direction**: coeff_to_slot | **normalized_name**: C2S / CoeffToSlot
**original_note_convention**: SlotToCoeff正向（原笔记命名；与全库统一术语方向相同但命名相反）
**confidence**: medium | **needs_human_review**: true
**source_id**: src_000019

## Knowledge Points

1. 给定系数→求slot值。标准线性代数：槽值=U_n乘系数向量。同态实现需避免明文的矩阵乘法。 slot_values = U_n·coeff
2. 系数多项式是实数→需要第二个密文提供虚部→合成复数向量→乘U_n得slot复数值。 coeff=Re(c_0)+i·Im(c_1), 复向量经U_n→复数slot值
3. 每层W_k仅O(N)非零元→logN层→总复杂度O(N log N)。这是PaCo加速的关键。 U_n=W_{logN}·…·W_1, 每层W_k稀疏
4. 因子分解中自然出现位反转，PaCo通过预排列索引映射吸收位反转。 Cooley-Tukey的奇偶分离→位反转置换
5. 互逆操作但对FHE复杂度不同。正向(U_n)和反向(U_n^{-1})的CT分解结构不同。 SlotToCoeff: coeff→slot(U_n); CoeffToSlot: slot→coeff(U_n^{-1})

⚠️ **术语审计修正** (2026-05-04T15:20): 原笔记标题'SlotToCoeff正向'实际讨论U_n·coeff=slot方向→按全库统一术语归入C2S/CoeffToSlot。
