# Batch-01 修正版公式表

> 生成时间: 2026-05-04 | 版本: v2

## 代数基础

| # | 公式 | 中文解释 | 概念 | source_id |
|---|------|----------|------|-----------|
| F001 | `Frac(R) = {a/b : a,b∈R, b≠0}/~` | 分式域的集合定义 | 分式域 | src_000001 |
| F002 | `K_splitting = F(α₁,…,αₙ), f=∏(X−αᵢ)` | 分裂域是包含f全部根的最小扩域 | 分裂域 | src_000001 |
| F003 | `τ_k(ζ_m) = ζ_m^k, k∈(Z/mZ)^×` | 分圆域中自同构的定义 | 共轭/自同构 | src_000002 |
| F004 | `|Gal(Q(ζ_m)/Q)| = φ(m)` | 分圆域Galois群阶 | Galois群 | src_000002 |
| F005 | `σᵢ = ι ∘ τᵢ` | 复嵌入的两步分解 | 典范嵌入 | src_000003 |
| F006 | `Tr_{K/F}(α) = Σ_{σ∈Gal(K/F)} σ(α) ∈ F` | 迹的定义与落入基域 | 迹 | src_000004 |
| F007 | `τ(Tr(α)) = Tr(α), ∀τ∈Gal(K/F)` | 迹被所有自同构固定的不变性 | 迹 | src_000004 |
| F008 | `[K:Q] = r₁ + 2r₂` | 扩张次数的实-复分解 | 自同构 | src_000005 |
| F009 | `τ_i(ζ_{m₀}) = ζ_{m₀} ⇔ i ≡ 1 (mod m₀)` | 自同构固定子域生成元的充要条件 | Prime Splitting | src_000006 |
| F010 | `F[X]/(f) ≅ F(ψ)` | 第一同构定理在多项式环上的应用 | 代值同态 | src_000007 |

## 有限域与分圆多项式

| # | 公式 | 中文解释 | 概念 | source_id |
|---|------|----------|------|-----------|
| F011 | `d = ord_m(p) = min{k>0: p^k ≡ 1 (mod m)}` | p模m的乘法阶 | Frobenius | src_000009 |
| F012 | `L = φ(m)/d` | 不可约因子个数=槽数 | SIMD | src_000009 |
| F013 | `Φ_m(X) ≡ ∏_{j=1}^{L} f_j(X) (mod p), deg(f_j)=d` | 分圆多项式在F_p上分解 | 分圆分解 | src_000009 |
| F014 | `Orb(ζ_m^r) = {ζ_m^{rp^i} : i=0,…,d-1}` | Frobenius轨道定义 | Frobenius轨道 | src_000008 |
| F015 | `X^N ≡ −1 (mod X^N+1)` | 稀疏分圆多项式的快速模约化 | 分圆优势 | src_000010 |
| F016 | `NTT: A_k = Σ_{j=0}^{n-1} a_j·ω^{jk} (mod q)` | NTT正变换 | NTT | src_000017 |
| F017 | `x' = x+ω^k·y, y' = x−ω^k·y` | Cooley-Tukey蝴蝶操作 | NTT/FFT | src_000017 |

## FHE加密/解密

| # | 公式 | 中文解释 | 概念 | source_id |
|---|------|----------|------|-----------|
| F018 | `BGV Enc: c = (a·s + t·e + m, −a)` | BGV加密（消息在低位） | BGV | src_000011 |
| F019 | `BGV Dec: m = [c₀ + c₁·s]_q mod t` | BGV解密 | BGV | src_000011 |
| F020 | `BFV Enc: c = (a·s + Δ·m + e, −a), Δ = ⌊q/t⌋` | BFV加密（消息放大Δ在高位） | BFV | src_000011 |
| F021 | `BFV Dec: m = ⌊(t/q)·[c₀ + c₁·s]_q⌉` | BFV解密 | BFV | src_000011 |
| F022 | `CKKS Encode: m = ⌊Δ·σ^{-1}(z)⌉` | CKKS编码（逆典范嵌入+放大取整） | CKKS | src_000011 |

## FHE噪声控制

| # | 公式 | 中文解释 | 概念 | source_id |
|---|------|----------|------|-----------|
| F023 | `ModSwitch: c' = ⌊(q'/q)·c⌉, m不变` | BGV模数切换 | ModSwitch | src_000014 |
| F024 | `BFV Rescale: c' = ⌊(1/t)·c⌉` | BFV重缩放（恢复消息尺度） | Rescale(BFV) | src_000014 |
| F025 | `CKKS Rescale: c' = ⌊(1/Δ)·c⌉` | CKKS重缩放（降低模数层级） | Rescale(CKKS) | src_000014 |
| F026 | `c⊗d = (c₀d₀, c₀d₁+c₁d₀, c₁d₁)` | 密文张量积（乘法） | 同态乘法 | src_000012 |
| F027 | `KeySwitch_{s²→s}: c'₀ = c₀ + Σc₂^{(i)}·b_i` | 重线性化公式 | 密钥切换 | src_000013 |

## CKKS编码与Bootstrapping

| # | 公式 | 中文解释 | 概念 | source_id |
|---|------|----------|------|-----------|
| F028 | `τ_n(a) = (a(ζ^{s₀}),…,a(ζ^{s_{N-1}}))` | 典范嵌入定义 | 典范嵌入 | src_000018 |
| F029 | `U_n = (ζ^{s_i·s_j})_{i,j=0}^{N-1}` | 典范嵌入矩阵 | U_n矩阵 | src_000018 |
| F030 | `slot = U_n·coeff, coeff = U_n^{-1}·slot` | SlotToCoeff和逆变换 | SlotToCoeff | src_000019 |
| F031 | `Tr(c) = c + Σ_{k=1}^{d-1} Rot_k(c)` | 同态迹的旋转-加法链 | 同态迹 | src_000021 |
| F032 | `s = Σ_{j=1}^{h} s_j, wt(s_j) = B = N/(4h)` | PaCo结构化密钥 | PaCo | src_000024 |
| F033 | `Y = X^{N/(2n)}: Z[Y]/(Y^{2n}+1) ↪ Z[X]/(X^N+1)` | 环嵌入（幂次代换） | 环嵌入 | src_000022 |

## 多项式分解

| # | 公式 | 中文解释 | 概念 | source_id |
|---|------|----------|------|-----------|
| F034 | `p ↦ (E_p, O_p)` | 系数拆分 | Coeff split | src_000023 |
| F035 | `p ↦ (E_p+H·O_p, E_p−H·O_p)` | 槽保持拆分 | Slot-preserving split | src_000023 |
| F036 | `p ↦ F_n^{-1}·S·F_N·p` | 任意矩阵拆分框架 | General split | src_000023 |

## 格密码

| # | 公式 | 中文解释 | 概念 | source_id |
|---|------|----------|------|-----------|
| F037 | `c = ⌊⟨s,b̃_n⟩/|b̃_n|²⌉` | Babai最近平面选择规则 | CVP/Babai | src_000025 |
| F038 | `s' = s − c·b_n` | Babai递归降维 | CVP/Babai | src_000025 |

## ⚠ 需要人工确认的公式

| # | 公式/问题 | 原因 |
|---|-----------|------|
| F-REVIEW-1 | PaCo中p vs Δ的确切数值和范围 | p=明文模数, Δ=编码因子, 需对照原文确认 |
| F-REVIEW-2 | PaCo Packing Relation: Pack(m₁,…,m_h)的精确公式 | 翻译材料中可能有细节遗漏 |
| F-REVIEW-3 | GHPS Duality / Good Bases 的具体定义 | 论文尚未导入源文件库(source_pending_import) |

## 来源

- 所有公式来源于 Batch-01 已读 27 个源文件
- 公式编号 F001-F038
- 标注 F-REVIEW 的项需要对照原文确认（详见 `需要人工确认清单.md`）
