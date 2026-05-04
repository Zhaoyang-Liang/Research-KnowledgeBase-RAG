# LWE与RLWE的基本关系与困难性概述

**card_id**: K065 | **batch**: Batch-05c | **confidence**: medium | **needs_human_review**: false | **review_policy**: review_when_used_for_paper
**source_id**: src_000110, src_000119, src_000058 | **claim_source_type**: 用户笔记 + 标准文献总结（Regev 2005, LPR 2010）
**备注**: 全库术语：LWE=Learning With Errors, RLWE=Ring-LWE, SVP=Shortest Vector Problem, SIVP=Shortest Independent Vectors Problem。标准归约细节见 Regev 2005 和 LPR 2010。

## Knowledge Points

1. 1. LWE (Learning With Errors) 定义：给定 A ∈ Z_q^{n×m}（均匀随机）和 b = A·s + e mod q（s 为短密钥，e 为小噪声），区分 (A,b) 与均匀随机。搜索版：从 (A,b) 恢复 s。判定版：区分分布。
2. 2. RLWE (Ring-LWE) 定义：将 LWE 的 Z_q^n→Z_q^m 映射替换为多项式环 R_q = Z_q[X]/(Φ_M(X)) 上的乘法。样本：(a, a·s + e) ∈ R_q^2，其中 s, e 为小多项式。Ring 结构使每个 RLWE sample 携带 n 个 LWE sample 的信息。
3. 3. RLWE 效率优势：RLWE 单次加密 = R_q 上一次乘法和加法（可 NTT 加速到 O(n log n)），而 LWE 需要矩阵乘法 O(n·m)。密钥和密文大小：RLWE = O(n)，LWE = O(n²)。此优势是 FHE 所有方案（BGV/BFV/CKKS）使用 RLWE 而非 plain LWE 的根本原因。
4. 4. LWE 困难性来源（Regev 2005）：存在从最坏情况格问题（GapSVP, SIVP）到平均情况 LWE 的量子归约。这意味着：如果有人能破解（平均）LWE，就有人能用量子算法破解（最坏）SVP。因此 LWE 的安全性不低于格问题的难度。
5. 5. RLWE 困难性来源（Lyubashevsky-Peikert-Regev 2010）：存在从 Ideal-SVP（理想格上的最坏情况 SVP）到 RLWE 的量子归约。Ideal lattice = 多项式理想对应的格，具有 ring 结构的额外代数对称性。归约假设分圆多项式 Φ_M(X) 不可约 mod 某些参数，且 noise distribution 满足一定条件。
6. 6. 当前论文中的位置：CKKS bootstrapping 的安全性依赖 RLWE。论文不需要重新证明 RLWE 的安全性——引用标准归约即可。但论文中 Split/Merge 引入的小环操作必须确保不会降低安全性（小环 RLWE 安全性需要额外参数说明，参见 P0-7 安全参数讨论）。
7. ⚠️ Caveat：LWE/RLWE 的归约涉及量子算法（非经典归约的最优形式），且归约中的噪声界和参数关系非常复杂。本卡仅是概览层次的理解性总结，不替代标准文献。

**交叉引用**: K056(Y-X): FHE方案总览（各方案均使用RLWE）；K062(Y-X): FHE方案与当前论文关系（安全参数讨论）。
