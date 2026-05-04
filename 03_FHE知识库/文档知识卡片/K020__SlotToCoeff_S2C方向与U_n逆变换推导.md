# SlotToCoeff (S2C) 方向与 U_n^{-1} 逆变换推导

**card_id**: K020 | **batch**: Batch-01
**normalized_direction**: slot_to_coeff | **normalized_name**: S2C / SlotToCoeff
**original_note_convention**: CoeffToSlot反向（原笔记命名；与全库统一术语方向相同但命名相反）
**confidence**: medium | **needs_human_review**: true
**source_id**: src_000020

## Knowledge Points

1. 从slot值恢复多项式系数。U_n^{-1}乘以slot复数向量得N个复数系数。 coeff=U_n^{-1}·slot
2. 因多项式系数是实数，复系数向量满足共轭对称→可压缩为N个实数。 coeff[−j]=conj(coeff[j])→只需N/2个独立复数
3. N个独立实数→拆分为N/2个实部+N/2个虚部→两个不同的多项式(两个RLWE密文)。 c₀=Re(coeff), c₁=Im(coeff)→两个实系数密文
4. 因U_n=(ζ^{s_i s_j})的行向量近似正交，U_n^{-1}≈(1/N)U_n^H，可用类似正向的CT分解。 U_n^{-1}≈(1/N)·U_n^H(共轭转置)
5. 系数表示需要N个实数(一个多项式密文)，但slot提供N/2个复数→分解为两个多项式——每个存一半实数的加密。 N个实数→N/2个复数→N/2实部+N/2虚部→两个多项式

⚠️ **术语审计修正** (2026-05-04T15:20): 原笔记标题'CoeffToSlot反向'实际讨论U_n^-1·slot=coeff方向→按全库统一术语归入S2C/SlotToCoeff。
