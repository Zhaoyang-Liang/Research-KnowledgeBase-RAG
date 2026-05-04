# Related Work 定位表

> paper_id: paper_2027_eurocrypt_fhe | 2026-05-04

| # | 论文 | 做了什么 | 与本文相同点 | 与本文不同点 | 本文不能claim | 本文可以claim | 精读? | citation? |
|---|------|---------|------------|------------|-------------|-------------|------|----------|
| 1 | CKKS 2017 (src_000041) | CKKS bootstrapping框架 | 使用C2S/EvalMod/S2C core | 不做分解 | core decomposition | — | 否 | 是 |
| 2 | PaCo (src_000040, src_000029) | 新bootstrapping结构 | bootstrapping加速 | 用 blind rotations替代传统core | PaCo的加速方法 | 保留传统core,只改粒度 | 轻读 | 是 |
| 3 | HERMES (src_000043) | Ring packing + transciphering | 使用HERMES-style ring switching | 输入不同: 多LWE→CKKS | Ring switching本身 | subring-ciphertext decomposition | 是 | 是 |
| 4 | GHPS Field Switching (src_000028) | Ring switching foundational theory | ring switching数学基础 | 聚焦correctness/noise proof | Field switching theory | — | 是 | 是 |
| 5 | Subring Secret Encapsulation | Subring-secret加速bootstrapping | 都使用subring | 只加速secret,不分解ciphertext | 密钥加速 | 真正ciphertext分解 | 轻读 | 是 |
| 6 | CKKS bootstrapping 优化 (multiple) | 优化individual linear transforms | 都在bootstrapping pipeline内 | 优化单个步骤 | 那些优化方法 | core factorization | 否 | 是 |
| 7 | Geelen SlotToCoeff (src_000030) | Slot-to-coeff 算法 | FFT-like分解 | power-of-2 BGV/BFV slot-to-coeff | Slot-to-coeff具体算法 | — | 轻读 | 可 |
| 8 | Ring switching old (src_000044) | 早期 ring switching | 概念相关 | 需要slot恢复 | 传统方法 | 不做slot恢复 | 轻读 | 可 |
| 9 | Bootstrapping综述 Wang (src_000042) | 综述 | bootstrapping概览 | — | — | — | 轻读 | 可 |

## 本文可 write 的区别陈述

1. 与 HERMES: 输入是单个CKKS密文（非多LWE），目标刷新自身（非packing）
2. 与 subring-secret: 分解ciphertext（非仅加速secret）
3. 与 PaCo: 保留传统 C2S/EvalMod/S2C 框架
4. 与传统优化: 改变core执行粒度（非优化单个transform）
5. 与 ring switching 传统: 不在split后恢复big-ring slots
