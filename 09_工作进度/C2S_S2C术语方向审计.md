# C2S / S2C 术语方向审计报告

> 2026-05-04T15:22 | 审计人：AI | 下次复查：Batch-05a3 完成后

---

## 审计范围

检查了 **20 张卡片**：K018~K049（含 K018, K019, K020, K021, K022, K024, K029, K030, K031, K036, K039, K040, K041, K042, K044, K045, K046, K047, K048, K049）

---

## 全库统一术语约定

全库默认采用以下方向：

| 缩写 | 全称 | 方向 | 矩阵公式 |
|------|------|------|---------|
| **C2S** | CoeffToSlot | coefficient → slot | slot = U_n · coeff |
| **S2C** | SlotToCoeff | slot → coefficient | coeff = U_n^-1 · slot |

> ⚠️ 原始笔记中若出现相反命名，必须标注「原笔记convention，与全库统一术语不同」。

---

## 三层概念区分

### 1. CKKS 基础编码层
- encoding: complex vector → σ^-1 → polynomial
- decoding: polynomial → σ → complex vector
- **与 bootstrapping 的 C2S/S2C 不等同**（见 K049, Q143, Q144）

### 2. Bootstrapping 线性变换层
- C2S: coeff → slot（供 EvalMod 使用）
- EvalMod: slot-wise modular reduction
- S2C: slot → coeff（返回结果）

### 3. BGV/BFV vs CKKS slot 严格区分
- BGV/BFV: CRT slot
- CKKS: complex slot
- **不可混用**（见 K041）

---

## 逐卡审计结果

### K018 — ✅ 无问题
- title: CKKS典范嵌入与U矩阵表示推导
- 方向: τ_n: poly → slot values (正向)，slot = U_n·coeff
- 不涉及 C2S/S2C 术语命名

### K019 — ⚠️ 已修正（方向冲突）
- **修正前标题**: SlotToCoeff正向完整推导
- **修正后标题**: CoeffToSlot (C2S) 方向与 U_n 矩阵推导
- **修正前方向**: 标题说SlotToCoeff(slot→coeff)，但公式 slot=U_n·coeff
- **修正后方向**: normalized_direction = coeff_to_slot, normalized_name = C2S/CoeffToSlot
- **修正后 confidence**: medium | **needs_human_review**: true
- **原因**: 原笔记标题与实际公式方向相反。公式讨论的是 C2S (coeff→slot)，非 S2C
- **影响**: 4 个文件已重命名（JSON+MD across 2 locations）

### K020 — ⚠️ 已修正（方向冲突）
- **修正前标题**: CoeffToSlot反向推导：一个密文变两个
- **修正后标题**: SlotToCoeff (S2C) 方向与 U_n^-1 逆变换推导
- **修正前方向**: 标题说CoeffToSlot(coeff→slot)，但公式 coeff=U_n^-1·slot
- **修正后方向**: normalized_direction = slot_to_coeff, normalized_name = S2C/SlotToCoeff
- **修正后 confidence**: medium | **needs_human_review**: true
- **影响**: 4 个文件已重命名

### K021 — ⚠️ 已修正（Product oversimplification）
- **问题**: Product 操作描述为与 Trace 等同
- **修正**: Product 相关知识点标注 confidence medium + caveat
- **caveat**: Product 涉及乘法深度、重线性化、level 消耗、scale 管理
- **needs_human_review**: true

### K022 — ⚠️ 已修正（Trace formula + ring embedding）
- **问题**: Y=X^{N/(2n)} 嵌入 + Trace 公式高置信
- **修正**: confidence→medium, needs_human_review→true
- **notation_mismatch_possible**: true（与当前论文 N=k·n, Y=X^k 可能不一致）

### K024 — ⚠️ 已修正（translation_note 降级）
- **修正前**: claim_source_type=翻译材料, confidence=high, needs_human_review=false
- **修正后**: claim_source_type=translation_note, confidence=medium, needs_human_review=true
- **needs_original_paper_check**: true（p vs Δ, Packing, Blind Rotation, B=N/(4h), 1.5-2.6× 加速）
- **状态**: 从 high-confidence primary source 降为 support card

### K029 — ⚠️ 已修正（source_id + 记号 + 拆分）
- **修正前**: source_id 为分号拼接字符串 "src_000028; src_000030"
- **修正后**: source_ids = ["src_000028", "src_000030"] (数组)
- **R^v → R^∨**: 标题及 KP 中统一记号
- **拆分**: factual_claims（confidence high）+ recommendation_claims（confidence medium, needs_human_review）

### K030 / K031 — ✅ 已有正确审计标记
- confidence: medium, needs_human_review: true
- 翻译笔记已标注

### K036 — ⚠️ 已修正（翻译笔记降 confidence）
- 修正前: confidence=high（与翻译笔记身份矛盾）
- 修正后: confidence→medium, claim_source_type→translation_note

### K039-K049 — ✅ 已补 confidence 字段
- K039: high (CKKS编码管线, 用户笔记)
- K040: high (π^-1/σ^-1, 用户笔记)
- K041: medium (slot对比, 跨论文综合) + needs_human_review
- K042: high (coeff/slot view, 用户笔记)
- K044: medium (C2S/S2C通用定义, BGV/BFV context)
- K045: medium (T矩阵, BGV/BFV context)
- K046: medium (M→U_ℓ→M^-1, BGV/BFV context)
- K047: medium (M矩阵, 跨论文)
- K048: medium (U_ℓ Fourier, 跨论文)
- K049: medium (论文位置, 综合推导) + needs_human_review

---

## 全局统计

| 指标 | 审计后 |
|------|--------|
| 总卡片数 | 49 |
| confidence high | 30 |
| confidence medium | 16 |
| needs_human_review | 16 |
| 方向冲突已修正 | 2 (K019, K020) |
| translation_note 降级 | 2 (K024, K036) |
| confidence 补填 | 10 (K039-K049) |

---

## K019 修正前/后对照

| 字段 | 修正前 | 修正后 |
|------|--------|--------|
| title | SlotToCoeff正向完整推导 | CoeffToSlot (C2S) 方向与 U_n 矩阵推导 |
| normalized_direction | (无) | coeff_to_slot |
| normalized_name | (无) | C2S / CoeffToSlot |
| confidence | high | medium |
| needs_human_review | false | true |
| 文件名 | K019__SlotToCoeff正向... | K019__CoeffToSlot_C2S方向... |

## K020 修正前/后对照

| 字段 | 修正前 | 修正后 |
|------|--------|--------|
| title | CoeffToSlot反向推导... | SlotToCoeff (S2C) 方向与 U_n^-1 逆变换推导 |
| normalized_direction | (无) | slot_to_coeff |
| normalized_name | (无) | S2C / SlotToCoeff |
| confidence | high | medium |
| needs_human_review | false | true |
| 文件名 | K020__CoeffToSlot反向... | K020__SlotToCoeff_S2C方向... |

---

## 未解决的 C2S/S2C 术语风险

以下风险在本次审计中识别但需要作者手动判断：

1. **K019/K020 original_note_convention**: 原笔记的 "SlotToCoeff正向" 和 "CoeffToSlot反向" 命名与全库统一术语的 "正/反" 可能来自不同参照系。已标注但需作者确认
2. **K048 U_ℓ vs U_n**: BGV/BFV 的 U_ℓ 与 CKKS 的 U_n 虽然都是 Fourier 型矩阵，但参数/维度/上下文不同
3. **K044-K046 C2S/S2C in BGV/BFV vs CKKS**: 用户在 BGV/BFV context 中理解 C2S，但当前论文使用 CKKS context。两者概念对称但细节不同
4. **K022 notation mismatch**: Y=X^{N/(2n)} vs N=k·n, Y=X^k 是否一致需确认

---

## 是否可以继续 Batch-05a3

✅ **可以继续，前提是：**

1. 作者已审阅本审计报告并确认 K019/K020 方向修正
2. Batch-05a3 的 SIMD/slot permutation/NTT 卡片沿用审计后的统一术语
3. 所有新卡片默认遵守 C2S=coeff→slot, S2C=slot→coeff 约定

**如果作者暂时无法审阅**：可以先继续 Batch-05a3，但所有新涉及 C2S/S2C 的新卡片默认标 medium + needs_human_review。
