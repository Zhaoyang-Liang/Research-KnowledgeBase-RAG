# RAG 检索测试问题集

> 2026-05-04T18:00 | 知识库工程完成阶段 | 仅测试问题，不执行检索

## 抽象代数到 FHE

| # | 问题 | 期望来源 | RAG status |
|---|------|----------|------------|
| T01 | 分式域和分裂域有什么不同？ | K001 | ready |
| T02 | trace 为什么和 field switching 有关？ | K007 | ready |
| T03 | R 和 R^∨ 的区别是什么？ | K029 | ready |
| T04 | 分圆多项式和 slot 数量有什么关系？ | K041 | ready |
| T05 | 复嵌入的分解步骤是什么？ | K039 | ready |
| T06 | i≡1(mod m₀) 条件是什么意思？ | K006 | ready |

## CKKS 编码 / C2S / S2C

| # | 问题 | 期望来源 | RAG status |
|---|------|----------|------------|
| T07 | CKKS encoding 和 C2S/S2C 有什么区别？ | K039, K044 | ready |
| T08 | C2S 和 S2C 的方向是什么？ | K019, K020 | ready |
| T09 | CKKS complex slot 和 BGV/BFV CRT slot 有什么区别？ | K041 | ready |
| T10 | 典范嵌入 τ_n 的作用是什么？ | K039 | ready |
| T11 | SIMD 编码的 Frobenius 轨道如何工作？ | K050 | ready |
| T12 | SlotToCoeff 在 bootstrapping 中起什么作用？ | K046, K049 | ready |
| T13 | Vandermonde 矩阵和典范嵌入矩阵是什么关系？ | K054 | defer |

## FHE 方案与操作

| # | 问题 | 期望来源 | RAG status |
|---|------|----------|------------|
| T14 | BGV 和 BFV 的主要区别是什么？ | K057 | ready |
| T15 | ModSwitch 和 Rescale 有什么区别？ | K058 | ready |
| T16 | KeySwitch 和 Relinearization 的关系是什么？ | K059, K068 | ready |
| T17 | CKKS 噪声控制的完整流程是什么？ | K063 | ready |
| T18 | 重线性化中为什么需要位分解？ | K068, Q179 | ready |
| T19 | RNS representation 是什么？ | K055 | ready |

## RLWE / 格密码

| # | 问题 | 期望来源 | RAG status |
|---|------|----------|------------|
| T20 | LWE 和 RLWE 有什么关系？ | K065 | ready |
| T21 | CVP / SVP / Babai 分别是什么？ | K066 | ready |
| T22 | Ideal lattice 和 RLWE 有什么关系？ | K069 | ready |
| T23 | LWE to SVP 归约说明什么？ | K067 | ready |
| T24 | Shor 算法为什么对格密码无效？ | Q183 | ready |
| T25 | GGH 签名为什么被破解？ | Q182 | ready |

## Bootstrapping

| # | 问题 | 期望来源 | RAG status |
|---|------|----------|------------|
| T26 | CKKS bootstrapping 的大致流程是什么？ | K042, K043 | ready |
| T27 | PaCo 算法的核心创新是什么？ | K033 (defer: translation_note) | defer / needs_original_check |

## 当前论文相关（应 defer）

| # | 问题 | 期望 RAG 行为 | RAG status |
|---|------|---------------|------------|
| T28 | reduced-Q 是否成立？ | 不应返回 ready chunk | defer |
| T29 | leaf bootstrapping error 如何合并？ | 不应返回 ready chunk | defer |
| T30 | γ_impl 是否进入数学定理？ | 不应返回 ready chunk | defer |
| T31 | 当前论文 C2S/S2C 的位置？ | 不应返回 ready chunk | defer |
| T32 | 多项式拆分是 slot-preserving 还是 coefficient split？ | 应返回 K025/K026 但标记 defer | defer |

## 测试设计

**ready 问题**（T01-T27）：期望 RAG 返回 card/Q&A chunk，答案可在知识卡片或原子问答中找到。

**defer 问题**（T28-T32）：期望 RAG 不返回确定答案，或返回标记了 `rag_status: defer` 的 chunk 并附带 `needs_human_review` 警告。

## 检索质量指标

| 指标 | 目标 |
|------|------|
| ready 问题 top-3 包含正确答案 | ≥ 80% |
| defer 问题不返回 false certainty | 100% |
| 所有返回 chunk 可回溯 source_id | 100% |
| 重复 chunk（副本源） | 0 |
