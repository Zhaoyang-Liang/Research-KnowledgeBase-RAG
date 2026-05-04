# Batch-05a 综合小结

> 2026-05-04T15:50 | Batch-05a1 + 05a2 + 05a3 整合（已修正 K043 confidence + K042 标题截断）

---

## 一、三批内容的层次关系

Batch-05a 将「编码/表示/槽/变换」拆成三层 + 一个数学工具背景层：

```text
Layer 1: CKKS 编码 / 典范嵌入 (Batch-05a1)
   complex vector -> canonical embedding -> polynomial
   基础编码层：消息与环元素之间的转换

Layer 2: C2S / S2C 线性变换 (Batch-05a2)
   coeff representation <-> slot representation
   Bootstrapping 核心层：同一环元素的两种表示之间转换

Layer 3: SIMD / Slot Permutation / Rotation (Batch-05a3 前半)
   Slot 编号的群论基础 + 槽间置换操作
   槽操作层：多槽的组织、编号、旋转

Layer 0 (数学背景): DFT / FFT / NTT (Batch-05a3 后半)
   复数域 / 有限域上的快速求值-插值算法
   线性变换和快速算法的通用数学工具
```

### 各层区分要点

| 层 | 核心问题 | 域 | 典型操作 |
|----|---------|-----|---------|
| CKKS encoding | 复数向量如何变成多项式 | C (复数) | sigma^{-1} |
| CKKS decoding | 多项式如何变回复数向量 | C | sigma |
| C2S | 系数表示 -> 槽表示 | 环 R_q | U_n · coeff |
| S2C | 槽表示 -> 系数表示 | 环 R_q | U_n^{-1} · slot |
| SIMD (BGV/BFV) | 多个独立消息打包进一个多项式 | 有限域 F_{p^d} | CRT 编码 |
| Slot Permutation | 槽之间的循环位移 | 环 R | tau_j 自同构 |
| DFT / FFT | 复数域快速求值-插值 | C | Cooley-Tukey |
| NTT | 有限域快速求值-插值 | F_q | Cooley-Tukey in F_q |

### 关键不等同关系

1. **CKKS encoding != C2S/S2C**：CKKS encoding 是外部消息 -> 多项式（sigma^{-1}）；C2S/S2C 是同一已加密环元素的不同表示转换
2. **NTT != CKKS canonical embedding**：NTT 在有限域 F_q 中，服务于多项式乘法加速；CKKS 典范嵌入在复数域 C 中，服务于消息编码
3. **BGV/BFV CRT slot != CKKS complex slot**：前者基于有限域的 CRT 分解；后者基于复数域上的典范嵌入坐标
4. **tau_j rotation != C2S/S2C**：前者是槽间置换（使用 automorphism key）；后者是表示转换（使用线性变换/矩阵分解）

---

## 二、长期术语规则（固定）

全库统一采用以下约定：

| 缩写 | 全称 | 方向 | 矩阵公式 |
|------|------|------|---------|
| **C2S** | CoeffToSlot | coefficient -> slot | slot = U_n · coeff |
| **S2C** | SlotToCoeff | slot -> coefficient | coeff = U_n^{-1} · slot |

### 使用规则

- 原始笔记中如果使用不同 convention，必须标注为「原笔记 convention，与全库统一术语不同，需人工确认」
- 后续论文协作中优先使用 C2S / S2C 缩写，少用容易混淆的 SlotToCoeff / CoeffToSlot 中文直译
- K019 和 K020 已按此规则修正并重命名（详见 C2S_S2C术语方向审计.md）

---

## 三、Batch-05a 核心卡片

总计 17 张（K039 ~ K055）

### Batch-05a1: CKKS 编码 / 典范嵌入 (5 张)

| card_id | title | confidence | review |
|---------|-------|------------|--------|
| K039 | CKKS编码管线：从复数向量到多项式（σ^{-1}全流程） | high | ✓ |
| K040 | π^{-1}共轭补全与σ^{-1}插值：CKKS编码的两个核心子步骤 | high | ✓ |
| K041 | CKKS complex slot与BGV/BFV CRT slot的本质区别 | medium | ⚠️ |
| K042 | CKKS coefficient view与slot view：Vandermonde变换与两种 | high | ✓ |
| K043 | CKKS编码toy example：用Z[X]/(X^4+1)编码复数2+i和1.60+0.90 | medium | ✓ |

### Batch-05a2: C2S / S2C 线性变换 (6 张)

| card_id | title | confidence | review |
|---------|-------|------------|--------|
| K044 | SlotToCoeff与CoeffToSlot：BGV/BFV槽-系数双表示转换 | medium | ✓ |
| K045 | C2S/S2C作为线性变换T的数学结构与矩阵表示 | medium | ✓ |
| K046 | M→U_ℓ→M^{-1}三步分解管线：C2S/S2C的层次化实现 | medium | ✓ |
| K047 | M矩阵：槽内局部基统一化与Frobenius算子展开 | medium | ✓ |
| K048 | U_ℓ矩阵：块Fourier/Vandermonde结构与FFT-like分解 | medium | ✓ |
| K049 | C2S/S2C在CKKS bootstrapping与当前论文中的位置 | medium | ⚠️ |

### Batch-05a3: SIMD / Slot Permutation / NTT (6 张)

| card_id | title | confidence | review |
|---------|-------|------------|--------|
| K050 | SIMD编码：CRT分解与槽并行计算（BGV/BFV context） | medium | ✓ |
| K051 | Slot编号的群论基础：Frobenius轨道、Z_m^*/⟨p⟩与Slot Permutati | medium | ✓ |
| K052 | τ_j自同构：多项式变量替换a(X)↦a(X^j)的数学结构与槽视角 | medium | ✓ |
| K053 | DFT/FFT/NTT定义与区别：复数域vs有限域vs快速算法 | low | ⚠️ |
| K054 | Vandermonde矩阵、DFT矩阵、典范嵌入矩阵的关系与区别 | low | ⚠️ |
| K055 | Ψ_{n,i}:K_n→L_i嵌入与第一同构定理：为什么槽里能放小域元素 | medium | ✓ |

### A. needs_human_review 卡片 (4 张)

这些卡片在元数据中标记了 needs_human_review: true，知识卡片层面需要作者审阅：

| card_id | title | confidence | 原因 |
|---------|------|------------|------|
| K041 | CKKS complex slot与BGV/BFV CRT slot的本质区别 | medium | ⚠️ needs_human_review：slot术语在CKKS和BGV/BFV中含义不同。当前论文使用CKKS co |
| K049 | C2S/S2C在CKKS bootstrapping与当前论文中的位置 | medium | ⚠️ candidate_knowledge + needs_human_review。当前论文的C2S/S2C位置判断 |
| K053 | DFT/FFT/NTT定义与区别：复数域vs有限域vs快速算法 | low | ⚠️ rough_note，confidence=low。Needs validation: Cooley-Tukey分 |
| K054 | Vandermonde矩阵、DFT矩阵、典范嵌入矩阵的关系与区别 | low | ⚠️ cross_source_synthesis, confidence=low。跨CKKS/BGV/BFV/信号处理 |

### B. 低置信度卡片 (2 张)

- **K053**: DFT/FFT/NTT定义与区别：复数域vs有限域vs快速算法
- **K054**: Vandermonde矩阵、DFT矩阵、典范嵌入矩阵的关系与区别

---

## 四、作者后续确认问题（独立于卡片 needs_human_review）

以下问题不是卡片 metadata 层面的 needs_human_review，而是在论文写作和 RAG 入库前需要作者确认的决策性问题。
这些问题**不阻塞后续 Batch-05b/05c**，但需要记录追踪。

### 术语与记号问题

**Q1: K019 / K020 原笔记 convention 是否接受**

K019 原笔记标题「SlotToCoeff正向」实际讨论 U_n·coeff=slot 方向，已按全库术语归入 C2S。
K020 原笔记标题「CoeffToSlot反向」实际讨论 U_n^{-1}·slot=coeff 方向，已按全库术语归入 S2C。

需确认：原笔记「正向/反向」的参照系是什么？全库统一术语是否与你的理解一致？

**Q2: K022 Y=X^{N/(2n)} vs Y=X^k 记号统一**

K022 使用 Y=X^{N/(2n)} 描述小环到大环的嵌入。当前论文上下文为 N=k·n, Y=X^k。

需确认：N/(2n) 和 k 是否等价（若 N=2kn 则 N/(2n)=k）？论文最终记号是哪个？

### 类比与综合问题

**Q3: K048 U_ell (BGV/BFV) vs U_n (CKKS) 类比范围**

需确认：两者是否只作为类比/结构启发使用，还是存在精确的数学对应关系（如 Cooley-Tukey 分解在两者中实现相同）？

**Q4: K054 Vandermonde/DFT/典范嵌入统一视角**

K054 试图统一 Vandermonde/DFT 矩阵/CKKS 典范嵌入矩阵/BGV/BFV U_ell 矩阵为「Vandermonde at roots-of-unity」。

需确认：这个视角是否只作为 intuition / 内部参考？是否与具体实现有冲突？

### 技术补充问题

**Q5: K053 Cooley-Tukey/bit-reversal 是否需要补充**

K053 (rough_note) 中缺少完整蝶形图和 bit-reversal permutation。

需确认：这些细节在论文中是否需要精确使用（DIT vs DIF、bit-reversal 顺序）？

### 入库策略问题

**Q6: FFT/DFT/NTT 笔记入库策略**

K053 和 K054 标记 low confidence + rough_note，当前不入正式知识库。

需确认：这些笔记在作者审阅后是否升级为 validated knowledge，还是永久作为 internal reference？

---

## 五、数据汇总

| 指标 | 数值 |
|------|------|
| Batch-05a 总卡片 | 17 |
| 其中 low confidence | 2 |
| 其中 needs_human_review | 4 |
| Batch-05a 导入文件 | 32 (05a1:5 + 05a2:13 + 05a3:14) |
| 作者后续确认问题 | 6 (Q1-Q6) |

---

## 六、是否可以进入 Batch-05b

✅ **可以进入 Batch-05b：FHE 方案总论 + 噪声/安全/参数**

以上 6 项问题已记录在案，不阻塞 Batch-05b 的执行。但在论文写作和 RAG 入库前需要作者确认。

后续 Batch-05b 继续遵守：
- C2S/S2C 术语规则
- CRT slot vs complex slot 区分
- 变换类型区分（DFT/FFT/NTT/canonical embedding）
- 跨方案比较标 medium + needs_human_review
