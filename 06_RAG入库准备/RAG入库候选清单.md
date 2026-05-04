# RAG 入库候选清单

> 2026-05-04T18:15 | 知识库工程完成阶段（权威来源：canonical_card_manifest.jsonl）
>
> ⚠️ 仅做候选清单，不执行 embedding，不切 chunk。

## 可以优先入 RAG（45 张）

条件：high/medium confidence + 非 needs_human_review + 非 candidate + source_id 完整

| card_id | title | confidence | batch | notes |
|---------|-------|------------|-------|-------|
| K001 | 分式域与分裂域的区别 | high | ? |  |
| K002 | 共轭元素的多项式角度理解 | high | ? |  |
| K003 | 复嵌入的两步分解：σ_i = ι ∘ τ_i | high | ? |  |
| K004 | Trace落入基域的证明与field switching中的应用 | high | ? |  |
| K005 | 自同构个数与实/复嵌入的对应 | high | ? |  |
| K006 | 分圆域中的素理想分裂：i≡1(mod m₀)条件的来源 | high | ? |  |
| K007 | 代数扩张的第一同构定理与代值同态 | high | ? |  |
| K008 | SIMD代数总结：Frobenius轨道→不可约因子→槽的三线推导 | high | ? |  |
| K009 | 分圆多项式在F_p上的分解与槽结构 | high | ? |  |
| K010 | 选择分圆多项式的十大优势 | high | ? |  |
| K011 | NTT与Cooley-Tukey快速算法 | high | ? |  |
| K012 | Gentry09/BGV/BFV/CKKS四方案全对比(v2增强) | high | ? |  |
| K013 | BGV/CKKS完整同态计算流程 | high | ? |  |
| K014 | 密钥切换通用公式与重线性化特例 | high | ? |  |
| K015 | 模数切换(ModSwitch)与重缩放(Rescale)的本质区别 | high | ? |  |
| K017 | BGV与BFV的核心区别：消息低位vs高位 | high | ? |  |
| K018 | CKKS典范嵌入与U矩阵表示推导 | high | ? |  |
| K021 | 同态Trace与Product操作的旋转分解 | medium | ? | review_when_used_for_paper; medium |
| K022 | Y=X^{N/(2n)}代换的环嵌入与提升 | high | ? |  |
| K027 | Babai最近平面算法 | high | ? |  |
| K037 | BGV/BFV Bootstrapping讲义总览：229页系统教程 | ? | Batch-04 |  |
| K038 | BGV Bootstrapping完整流程：6阶段与关键步骤 | ? | Batch-04 |  |
| K039 | CKKS编码管线：从复数向量到多项式（σ^{-1}全流程） | high | Batch-05a1 |  |
| K040 | π^{-1}共轭补全与σ^{-1}插值：CKKS编码的两个核心子步骤 | high | Batch-05a1 |  |
| K042 | CKKS coefficient view与slot view：Vandermonde变换 | high | Batch-05a1 |  |
| K043 | CKKS编码toy example：用Z[X]/(X^4+1)编码复数2+i和1.60+0 | medium | Batch-05a1 | medium |
| K044 | SlotToCoeff与CoeffToSlot：BGV/BFV槽-系数双表示转换 | medium | Batch-05a2 | medium |
| K045 | C2S/S2C作为线性变换T的数学结构与矩阵表示 | medium | Batch-05a2 | medium |
| K046 | M→U_ℓ→M^{-1}三步分解管线：C2S/S2C的层次化实现 | medium | Batch-05a2 | medium |
| K047 | M矩阵：槽内局部基统一化与Frobenius算子展开 | medium | Batch-05a2 | medium |
| K048 | U_ℓ矩阵：块Fourier/Vandermonde结构与FFT-like分解 | medium | Batch-05a2 | medium |
| K050 | SIMD编码：CRT分解与槽并行计算（BGV/BFV context） | medium | Batch-05a3 | medium |
| K051 | Slot编号的群论基础：Frobenius轨道、Z_m^*/⟨p⟩与Slot Permut | medium | Batch-05a3 | medium |
| K052 | τ_j自同构：多项式变量替换a(X)↦a(X^j)的数学结构与槽视角 | medium | Batch-05a3 | medium |
| K055 | Ψ_{n,i}:K_n→L_i嵌入与第一同构定理：为什么槽里能放小域元素 | medium | Batch-05a3 | medium |
| K056 | FHE方案总览：Gentry09 / BGV / BFV / CKKS 定义、公式与对比 | medium | Batch-05b1 | medium |
| K058 | FHE同态计算完整流程：乘法→重线性化→模数切换/重缩放→解密 | medium | Batch-05b1 | medium |
| K059 | 密钥切换定义与三种典型场景 | medium | Batch-05b1 | medium |
| K060 | 模数切换(ModSwitch)与密钥切换(KeySwitch)的本质区别 | medium | Batch-05b1 | review_when_used_for_paper; medium |
| K061 | 模数切换(ModSwitch)与重缩放(Rescaling)的本质区别 | medium | Batch-05b1 | review_when_used_for_paper; medium |
| K065 | LWE与RLWE的基本关系与困难性概述 | medium | Batch-05c | review_when_used_for_paper; medium |
| K066 | CVP、SVP与Babai最近平面算法的几何解释 | medium | Batch-05c | review_when_used_for_paper; medium |
| K067 | LWE to SVP归约主线：从格高斯到Fourier bias再到LWE | medium | Batch-05c | review_when_used_for_paper; medium |
| K068 | 重线性化与位分解（bit decomposition）的完整机制 | medium | Batch-05c | review_when_used_for_paper; medium |
| K069 | Ideal Lattice与RLWE：理想格的定义、性质与安全性 | medium | Batch-05c | review_when_used_for_paper; medium |

## 暂缓入 RAG（24 张）

条件：needs_human_review / low confidence / translation_note / candidate / paper_project

| card_id | title | 原因 | batch |
|---------|-------|------|-------|
| K016 | 噪声控制、重线性化、密钥切换的层次关系 | paper_project | ? |
| K019 | CoeffToSlot (C2S) 方向与 U_n 矩阵推导 | needs review | Batch-01 |
| K020 | SlotToCoeff (S2C) 方向与 U_n^{-1} 逆变换推导 | needs review | Batch-01 |
| K023 | Coefficient split vs Slot-preserving spl | needs review | ? |
| K024 | PaCo论文全文翻译摘要 | paper_project | ? |
| K025 | 多项式分解方法总览 | needs review | ? |
| K026 | 多项式分解正向实现细节 | needs review | ? |
| K028 | ⚠ CKKS slot与BGV/BFV CRT slot的区别——避免论文混用 | needs review, paper_project | ? |
| K029 | R 与 R^v 在 BGV/GHPS Field Switching 中的符号差 | needs review | ? |
| K030 | CKKS编码的逆变换τ_n^{-1}: 复向量到多项式的完整推导 | needs review, translation_note | Batch-03 |
| K031 | SlotToCoeff在密文态下的语义：表示切换的完整过程 | needs review, translation_note | Batch-03 |
| K032 | C^B到C^n的嵌入与编码保持性质 | needs review, translation_note | Batch-03 |
| K033 | PaCo算法结构：参数依赖图与索引系统 | needs review, translation_note | Batch-03 |
| K034 | CKKS中的Trace和Product操作及同态内积计算 | needs review, translation_note | Batch-03 |
| K035 | CKKS中5作为单位根生成元的数论性质 | needs review, translation_note | Batch-03 |
| K036 | 典范嵌入正向τ_n(p(Y)): 多项式到复向量的完整推导 | needs review, translation_note | Batch-03 |
| K041 | CKKS complex slot与BGV/BFV CRT slot的本质区别 | needs review | Batch-05a1 |
| K049 | C2S/S2C在CKKS bootstrapping与当前论文中的位置 | needs review, paper_project | Batch-05a2 |
| K053 | DFT/FFT/NTT定义与区别：复数域vs有限域vs快速算法 | needs review, low confidence | Batch-05a3 |
| K054 | Vandermonde矩阵、DFT矩阵、典范嵌入矩阵的关系与区别 | needs review, low confidence | Batch-05a3 |
| K057 | BGV vs BFV本质区别：消息低位高位、模数切换重缩放 | needs review | Batch-05b1 |
| K062 | BGV/BFV/CKKS/FHE方案层与当前论文CKKS Bootstrappi | paper_project | Batch-05b1 |
| K063 | CKKS噪声来源总览：加密/乘法/重缩放/KeySwitch/Bootstrap | paper_project | Batch-05b2-pre |
| K064 | KeySwitch噪声的来源、分解与上界 | paper_project | Batch-05b2-pre |

## 不建议入 RAG

- 临时草稿 / raw data（实验日志 .txt）
- side_material（Falcon、NTRU、Klein-GPV）
- 论文专属 candidate_knowledge（C-NOISE-001~010）
- 未导入或 source 不完整内容
- 重复副本（多目录中的副本保留 canonical 一份即可）

## 特别标注

### 本次修正（2026-05-04T18:15）
- K024 (PaCo 翻译摘要) → 从 ready 移到 defer（translation_note，需核原文）
- K016 (噪声/KeySwitch 层次关系) → defer（paper_project domain）
- K062/K063/K064 → defer（paper_project domain）
- K037/K038 → 补入 ready（讲义资源卡片）
- K021 → confidence high→medium（trace formula 跨源综合）
- K033 → defer（translation_note，非 PaCo 原文卡片）

### 值得注意
- K037/K038 的 confidence 为 ?（来自讲义资源，非用户笔记）
- K011 与 K053 的区别：K011 是标准 Cooley-Tukey 算法（high confidence），K053 是 DFT/FFT/NTT 三者区别（low confidence，已 defer）
- K022 是环嵌入/变量代换（标准代数操作，high confidence）

## RAG 准备建议

1. **优先批次**：45 张 ready → 先建 metadata index → 再切 chunk
2. **暂缓批次**：24 张 → 等 review/low_conf 解决后再入
3. **不做**：不在此阶段执行 embedding，仅做候选清单和 metadata 准备
