# RAG ready / defer / exclude 清单

> 2026-05-04T18:00 | 知识库工程完成阶段

## Ready：45 张卡片 + 165 条 Q&A

条件：high/medium confidence + source_id 完整 + 非 needs_human_review + 非 candidate + 非 translation_note

### 卡片

| card_id | title | domain | confidence |
|---------|-------|--------|------------|
| K001 | 分式域与分裂域的区别 | abstract_algebra | high |
| K002 | 共轭元素的多项式角度理解 | abstract_algebra | high |
| K003 | 复嵌入的两步分解：σ_i = ι ∘ τ_i | fhe | high |
| K004 | Trace落入基域的证明与field switching中的应用 | abstract_algebra | high |
| K005 | 自同构个数与实/复嵌入的对应 | abstract_algebra | high |
| K006 | 分圆域中的素理想分裂：i≡1(mod m₀)条件的来源 | abstract_algebra | high |
| K007 | 代数扩张的第一同构定理与代值同态 | fhe | high |
| K008 | SIMD代数总结：Frobenius轨道→不可约因子→槽的三线推导 | fhe | high |
| K009 | 分圆多项式在F_p上的分解与槽结构 | abstract_algebra | high |
| K010 | 选择分圆多项式的十大优势 | abstract_algebra | high |
| K011 | NTT与Cooley-Tukey快速算法 | fhe | high |
| K012 | Gentry09/BGV/BFV/CKKS四方案全对比(v2增强) | rlwe_pqc | high |
| K013 | BGV/CKKS完整同态计算流程 | fhe | high |
| K014 | 密钥切换通用公式与重线性化特例 | fhe | high |
| K015 | 模数切换(ModSwitch)与重缩放(Rescale)的本质区别 | fhe | high |
| K017 | BGV与BFV的核心区别：消息低位vs高位 | fhe | high |
| K018 | CKKS典范嵌入与U矩阵表示推导 | fhe | high |
| K021 | 同态Trace与Product操作的旋转分解 | abstract_algebra | high |
| K022 | Y=X^{N/(2n)}代换的环嵌入与提升 | fhe | high |
| K027 | Babai最近平面算法 | rlwe_pqc | high |
| K037 | BGV/BFV Bootstrapping讲义总览：229页系统教程 | fhe | ? |
| K038 | BGV Bootstrapping完整流程：6阶段与关键步骤 | fhe | ? |
| K039 | CKKS编码管线：从复数向量到多项式（σ^{-1}全流程） | fhe | high |
| K040 | π^{-1}共轭补全与σ^{-1}插值：CKKS编码的两个核心子步骤 | abstract_algebra | high |
| K042 | CKKS coefficient view与slot view：Vandermo | fhe | high |
| K043 | CKKS编码toy example：用Z[X]/(X^4+1)编码复数2+i和1 | fhe | medium |
| K044 | SlotToCoeff与CoeffToSlot：BGV/BFV槽-系数双表示转换 | fhe | medium |
| K045 | C2S/S2C作为线性变换T的数学结构与矩阵表示 | fhe | medium |
| K046 | M→U_ℓ→M^{-1}三步分解管线：C2S/S2C的层次化实现 | fhe | medium |
| K047 | M矩阵：槽内局部基统一化与Frobenius算子展开 | fhe | medium |
| K048 | U_ℓ矩阵：块Fourier/Vandermonde结构与FFT-like分解 | fhe | medium |
| K050 | SIMD编码：CRT分解与槽并行计算（BGV/BFV context） | fhe | medium |
| K051 | Slot编号的群论基础：Frobenius轨道、Z_m^*/⟨p⟩与Slot P | fhe | medium |
| K052 | τ_j自同构：多项式变量替换a(X)↦a(X^j)的数学结构与槽视角 | abstract_algebra | medium |
| K055 | Ψ_{n,i}:K_n→L_i嵌入与第一同构定理：为什么槽里能放小域元素 | abstract_algebra | medium |
| K056 | FHE方案总览：Gentry09 / BGV / BFV / CKKS 定义、公 | rlwe_pqc | medium |
| K058 | FHE同态计算完整流程：乘法→重线性化→模数切换/重缩放→解密 | fhe | medium |
| K059 | 密钥切换定义与三种典型场景 | fhe | medium |
| K060 | 模数切换(ModSwitch)与密钥切换(KeySwitch)的本质区别 | fhe | medium |
| K061 | 模数切换(ModSwitch)与重缩放(Rescaling)的本质区别 | fhe | medium |
| K065 | LWE与RLWE的基本关系与困难性概述 | rlwe_pqc | medium |
| K066 | CVP、SVP与Babai最近平面算法的几何解释 | rlwe_pqc | medium |
| K067 | LWE to SVP归约主线：从格高斯到Fourier bias再到LWE | rlwe_pqc | medium |
| K068 | 重线性化与位分解（bit decomposition）的完整机制 | fhe | medium |
| K069 | Ideal Lattice与RLWE：理想格的定义、性质与安全性 | rlwe_pqc | medium |

## Defer：24 张卡片 + 21 条 Q&A

条件：needs_human_review / low confidence / translation_note / candidate / paper_project

### 卡片

| card_id | title | 原因 |
|---------|-------|------|
| K016 | 噪声控制、重线性化、密钥切换的层次关系 | paper_project |
| K019 | CoeffToSlot (C2S) 方向与 U_n 矩阵推导 | needs_review |
| K020 | SlotToCoeff (S2C) 方向与 U_n^{-1} 逆变换推 | needs_review |
| K023 | Coefficient split vs Slot-preservin | needs_review |
| K024 | PaCo论文全文翻译摘要 | paper_project |
| K025 | 多项式分解方法总览 | needs_review |
| K026 | 多项式分解正向实现细节 | needs_review |
| K028 | ⚠ CKKS slot与BGV/BFV CRT slot的区别——避免 | needs_review+paper_project |
| K029 | R 与 R^v 在 BGV/GHPS Field Switching  | needs_review |
| K030 | CKKS编码的逆变换τ_n^{-1}: 复向量到多项式的完整推导 | needs_review |
| K031 | SlotToCoeff在密文态下的语义：表示切换的完整过程 | needs_review |
| K032 | C^B到C^n的嵌入与编码保持性质 | needs_review |
| K033 | PaCo算法结构：参数依赖图与索引系统 | needs_review |
| K034 | CKKS中的Trace和Product操作及同态内积计算 | needs_review |
| K035 | CKKS中5作为单位根生成元的数论性质 | needs_review |
| K036 | 典范嵌入正向τ_n(p(Y)): 多项式到复向量的完整推导 | needs_review |
| K041 | CKKS complex slot与BGV/BFV CRT slot的 | needs_review |
| K049 | C2S/S2C在CKKS bootstrapping与当前论文中的位置 | needs_review+paper_project |
| K053 | DFT/FFT/NTT定义与区别：复数域vs有限域vs快速算法 | needs_review+low_conf |
| K054 | Vandermonde矩阵、DFT矩阵、典范嵌入矩阵的关系与区别 | needs_review+low_conf |
| K057 | BGV vs BFV本质区别：消息低位高位、模数切换重缩放 | needs_review |
| K062 | BGV/BFV/CKKS/FHE方案层与当前论文CKKS Bootst | paper_project |
| K063 | CKKS噪声来源总览：加密/乘法/重缩放/KeySwitch/Boot | paper_project |
| K064 | KeySwitch噪声的来源、分解与上界 | paper_project |

## Exclude：0 张卡片 + 0 条 Q&A

条件：side_material / archive_only / 源文件不完整 / 纯实验日志

| 类别 | 说明 |
|------|------|
| side_material | Falcon (src_000116), NTRU (src_000118), Klein-GPV (src_000114) |
| 纯实验日志 | Batch-06 暂停项（2 个 .txt） |
| duplicate_paths | 28 张卡片的副本路径（RAG 只读 canonical） |
| skipped source_ids | src_000095~000100 |
