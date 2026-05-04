# K063: C2S/S2C方向统一 — PaCo语义权威定义

> **权威依据**: `05_知识关联/C2S_S2C_PaCo_术语方向说明.md` & PaCo论文
> **级别**: 全库最高优先级术语规则卡片

---

## 统一规范

```
C2S = CoeffToSlot = coefficient representation → slot representation
S2C = SlotToCoeff = slot representation → coefficient representation
```

## PaCo 证据

传统CKKS bootstrapping流程:

```
ModRaise → CoeffToSlot → EvalMod → SlotToCoeff
```

1. ModRaise: 密文 m 扩模得 m'=m+qJ
2. **CoeffToSlot**: 将多项式系数 m_i+qJ_i 搬入加密 slots
3. EvalMod: slotwise 近似计算去除 qJ_i
4. **SlotToCoeff**: 将 slots 中近似系数 m_i 移回系数多项式

PaCo 的 partial CoeffToSlot: 将 blind rotation 后 b'_v(Z) 的 coefficients 嵌入 slots — 仍属 CoeffToSlot 家族。

## 矩阵层 vs 操作语义层

U_n 矩阵: coefficient vector → evaluation/slot vector（裸矩阵方向）
但操作名以 bootstrapping 语义输入输出为准，不以矩阵方向为准。

PaCo conventional SlotToCoeff 中虽出现 U_n（coeff→eval方向），但操作语义是 slots→coefficients。

## 禁止表达式（canonical字段一律错误）

- ❌ SlotToCoeff(C2S)
- ❌ CoeffToSlot(S2C)
- ❌ C2S = slot → coeff
- ❌ S2C = coeff → slot

这些仅允许出现在 `original_note_convention` 字段。

## 与 CKKS 编码/解码的区别

| 层次 | 操作 | 方向 |
|------|------|------|
| CKKS 编码 | Ecd(z) | external complex vector → plaintext polynomial |
| CKKS 解码 | Dcd(m) | plaintext polynomial → complex vector |
| Bootstrapping | C2S (CoeffToSlot) | coefficient data → slots |
| Bootstrapping | S2C (SlotToCoeff) | slots → coefficient polynomial |

## 交叉引用

- K019: CKKS C2S = U_n·coeff（矩阵层方向）
- K020: CKKS S2C = U_n⁻¹·slot（矩阵层方向）
- K044: BGV-BFV 双表示转换（已按本卡修正）
- 权威依据: `05_知识关联/C2S_S2C_PaCo_术语方向说明.md`
