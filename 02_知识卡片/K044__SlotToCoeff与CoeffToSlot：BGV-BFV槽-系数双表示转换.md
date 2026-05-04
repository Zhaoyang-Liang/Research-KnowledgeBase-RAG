# K044: CoeffToSlot与SlotToCoeff — 基于PaCo统一术语的系数-槽双表示转换

> **术语权威依据**: `05_知识关联/C2S_S2C_PaCo_术语方向说明.md` (2026-05-04)
> **统一规则**: C2S = CoeffToSlot = coefficient → slot; S2C = SlotToCoeff = slot → coefficient

## 核心定义

1. **CoeffToSlot (C2S)**: 将系数表示(coefficient representation)同态转换为槽表示(slot representation)。
   输入：多项式系数 m_i（密文态下隐藏）；输出：加密槽中的 m_i 值。
   PaCo 证据：传统 CKKS bootstrapping 中，ModRaise 后 CoeffToSlot 将多项式系数 m_i+qJ_i 搬入加密 slots。

2. **SlotToCoeff (S2C)**: C2S 的逆操作。将槽表示中已得到的系数信息同态转回系数多项式。
   PaCo 证据：EvalMod 后 SlotToCoeff 将 slots 中近似系数移回 plaintext polynomial coefficients。

3. C2S/S2C ≠ CKKS 编码/解码。CKKS 编码是外部复数向量→多项式(σ⁻¹)；C2S/S2C 是同一环元素两种表示间的线性变换。

4. 在 BGV/BFV 中: 槽表示来自 CRT 分解（每槽=R mod 不可约因子），系数表示来自原始多项式环。两者是同一数学对象的不同 view。

## 原笔记旧命名

⚠️ **原笔记使用了相反命名**，已降级为 `original_note_convention`：
- 旧: SlotToCoeff(C2S) = slot → coeff
- 旧: CoeffToSlot(S2C) = coeff → slot

PaCo 论文和知识库统一术语采用相反方向：
- 新: CoeffToSlot(C2S) = coeff → slot
- 新: SlotToCoeff(S2C) = slot → coeff

## 交叉引用

- K019: CKKS C2S = U_n·coeff 方向（矩阵层）
- K062: CKKS context 中 C2S/S2C
- 权威依据: `05_知识关联/C2S_S2C_PaCo_术语方向说明.md`
