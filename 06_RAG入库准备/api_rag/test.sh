python3 query_rag.py "C2S 和 S2C 的方向是什么？" -k 5 -v
python3 query_rag.py "CoeffToSlot 和 SlotToCoeff 分别是什么意思？" -k 5 -v
python3 query_rag.py "C2S/S2C 和 CKKS encoding/decoding 有什么本质区别？" -k 5 -v
python3 query_rag.py "CKKS 编码管线是什么？从复数向量到多项式经历哪些步骤？" -k 5 -v
python3 query_rag.py "CKKS 编码中的 σ^{-1}、π^{-1}、共轭补全分别是什么？" -k 5 -v
python3 query_rag.py "CKKS 的典范嵌入和 Vandermonde 矩阵有什么关系？" -k 5 -v
python3 query_rag.py "CKKS 的 slot 和 BGV BFV 的 CRT slot 有什么不同？" -k 5 -v
python3 query_rag.py "BGV、BFV、CKKS 三种方案的核心区别是什么？" -k 5 -v
python3 query_rag.py "BGV 的 ModSwitch 和 CKKS 的 Rescale 有什么本质区别？" -k 5 -v
python3 query_rag.py "ModSwitch、Rescale、KeySwitch 三者有什么区别？" -k 5 -v
python3 query_rag.py "BFV 的 Rescale 和 CKKS 的 Rescale 是一回事吗？" -k 5 -v
python3 query_rag.py "KeySwitching 的通用公式是什么？" -k 5 -v
python3 query_rag.py "重线性化为什么是 KeySwitch 的特例？" -k 5 -v
python3 query_rag.py "同态乘法后为什么要做重线性化？不做会怎样？" -k 5 -v
python3 query_rag.py "Bootstrapping 在 FHE 计算流程中什么时候触发？" -k 5 -v
python3 query_rag.py "CKKS bootstrapping 为什么比 BGV BFV bootstrapping 更高效？" -k 5 -v
python3 query_rag.py "BGV BFV bootstrapping 的完整流程是什么？" -k 5 -v
python3 query_rag.py "EvalMod 在 CKKS bootstrapping 中起什么作用？" -k 5 -v
python3 query_rag.py "SlotToCoeff 和 EvalMod 在 bootstrapping 中是什么关系？" -k 5 -v
python3 query_rag.py "PaCo 的 partial CoeffToSlot 为什么叫 partial？" -k 5 -v
python3 query_rag.py "PaCo 和传统 CKKS bootstrapping 的区别是什么？" -k 5 -v
python3 query_rag.py "Trace 操作为什么会落入基域？" -k 5 -v
python3 query_rag.py "Trace 在 field switching 中为什么重要？" -k 5 -v
python3 query_rag.py "FHE 中如何同态计算 Trace？" -k 5 -v
python3 query_rag.py "Rotation 和 C2S/S2C 有什么区别？" -k 5 -v
python3 query_rag.py "τ_j 自同构 a(X) 到 a(X^j) 在 slot 视角下是什么意思？" -k 5 -v
python3 query_rag.py "Frobenius 轨道和 slot permutation 有什么关系？" -k 5 -v
python3 query_rag.py "SIMD 编码在 BGV BFV 中是怎么通过 CRT 实现的？" -k 5 -v
python3 query_rag.py "NTT、FFT、DFT 的区别是什么？" -k 5 -v
python3 query_rag.py "NTT 为什么可以加速多项式乘法？" -k 5 -v
python3 query_rag.py "Vandermonde、DFT、典范嵌入三者有什么关系？" -k 5 -v
python3 query_rag.py "LWE 和 RLWE 有什么关系？为什么 FHE 更常用 RLWE？" -k 5 -v
python3 query_rag.py "RLWE 的安全性和 ideal lattice 有什么关系？" -k 5 -v
python3 query_rag.py "CVP、SVP、Babai 最近平面算法分别是什么？" -k 5 -v
python3 query_rag.py "LWE 到 SVP 的归约主线是什么？" -k 5 -v
python3 query_rag.py "重线性化和 bit decomposition 有什么关系？" -k 5 -v
python3 query_rag.py "Gentry09、BGV、BFV、CKKS 的历史关系是什么？" -k 5 -v
python3 query_rag.py "R 和 R^∨ 在 GHPS 和当前论文中有什么区别？" -k 5 -v
python3 query_rag.py "GHPS 中为什么密文空间写成 (R_q^∨)^2？" -k 5 -v
python3 query_rag.py "当前论文能否直接把 R^∨ 省略成 R？" -k 5 -v
python3 query_rag.py "coefficient split 和 slot-preserving split 有什么区别？" -k 5 -v
python3 query_rag.py "当前论文中 Split、C2S、EvalMod、S2C、Merge 的关系是什么？" -k 5 -v
python3 query_rag.py "Ring switching noise 和 KeySwitch noise 有什么关系？" -k 5 -v
python3 query_rag.py "CKKS 噪声来源有哪些？" -k 5 -v
python3 query_rag.py "KeySwitch 噪声来源、分解基和上界是什么？" -k 5 -v
python3 query_rag.py "leaf bootstrapping error 如何合并？" -k 5 -v
python3 query_rag.py "reduced-Q leaf bootstrapping 是否成立？" -k 5 -v
python3 query_rag.py "当前论文的完整噪声定理是什么？" -k 5 -v
python3 query_rag.py "当前论文最终主定理应该怎么写？" -k 5 -v
python3 query_rag.py "K004" -k 5 -v
python3 query_rag.py "K044" -k 5 -v
python3 query_rag.py "K063" -k 5 -v
python3 query_rag.py "K068" -k 5 -v
python3 query_rag.py "Q131" -k 5 -v
python3 query_rag.py "Q177" -k 5 -v
python3 query_rag.py "src_000059" -k 5 -v
python3 query_rag.py "src_000068" -k 5 -v
python3 query_rag.py "src_000119" -k 5 -v