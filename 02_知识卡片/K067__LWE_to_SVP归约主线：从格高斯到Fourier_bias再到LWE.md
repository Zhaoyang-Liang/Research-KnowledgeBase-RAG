# LWE to SVP归约主线：从格高斯到Fourier bias再到LWE

**card_id**: K067 | **batch**: Batch-05c | **confidence**: medium | **needs_human_review**: false | **review_policy**: review_when_used_for_paper
**source_id**: src_000110 | **claim_source_type**: 用户笔记（LWE_SVP_reduc.tex）总结 Regev 2005
**备注**: 全库术语：torus T = R/Z 是连续模1的圆环。D_{L,s} = 格上的离散高斯。f_s(x) 在 torus 上定义为离散分布的线性组合。具体的 Fourier 展开公式和量化界见 LWE_SVP_reduc.tex 原文。

## Knowledge Points

1. 1. 归约目标（Regev 2005）：Worst-case lattice problem (GapSVP/SIVP) ≤_quantum Average-case LWE。即存在量子归约将最坏情况格问题归约到平均情况 LWE。这给 LWE 安全性提供了坚实的理论基础。归约路线不是 trivial 的，经过多个中间步骤。
2. 2. 归约主线五步结构：(a) 格 L → 连续模格高斯 D_{R^n/L, s}（在 torus 上采样）。(b) D → Poisson 求和 / Fourier 展开 → 揭示格的对偶格 L^∨ 的结构。(c) Fourier 分析 → f_s(x) = Uniform + Fourier bias（在一定参数下，LWE 样本不是纯均匀的——它带有来自格的 Fourier 偏向，表现为格对偶点在频率域上的贡献）。(d) 离散化 → a = ⌊qx⌋ mod q → 离散 Fourier 结构。(e) Fourier 偏向 + 噪声 → LWE 样本。
3. 3. 核心 insight：f_s(x) 在 torus 上近似均匀，但其 Fourier 系数由对偶格 L^∨ 上的高斯质量决定。当 s 足够大时，Fourier 偏向（非常数 terms）exponentially small。LWE 求解器需要能从微小的非均匀信号中提取信息——这等价于解决对偶格上的近距解码问题（BDD）。
4. 4. 量子归约环节（为什么需要量子）：经典归约可以实现 BDD-to-LWE，但需要额外的技术假设。Regev 2005 用量子算法在第一步做高斯采样（从连续模格高斯采样需要量子），后续步骤是经典的。这是方案需要量子归约的核心原因。近年有部分经典归约工作，但噪声率和近似因子不如量子版紧。
5. 5. 与 FHE 的关系：(a) FHE 方案（BGV/BFV/CKKS）的安全性基于 RLWE（LWE 的 ring variant）。(b) 归约给了我们理由相信：只要最坏情况格问题难，随机 (A, b) 就不可区分。(c) 参数选择必须匹配归约条件（noise rate α 与模数 q 的关系），否则归约的保证失效。
6. ⚠️ Caveat：归约中的参数关系非常技术化。量子归约是否可完全经典化是一个活跃研究课题。本卡只给主线，具体推导和常数见原始文献 Regev 2005。

**交叉引用**: K065(Y-X): LWE/RLWE基本关系（归约是其基础）；K066(Y-X): SVP/CVP问题定义（归约的源问题）。
