# 抽象代数到全同态加密的数学关联图

> 从代数结构出发，追踪每条数学概念如何在 FHE 中被"实例化"。
> 已读范围: Batch-01 共 24 篇精读文件。
> 标注约定: **需要人工确认** 表示笔记表述或推导需对照原论文核查。

---

## 0. 概念-技术映射总表

| # | 抽象代数概念 | FHE 对应结构 | 技术模块 | 是否与当前论文有关 | 置信度 |
|---|------------|-------------|---------|-----------------|-------|
| 1 | 分圆多项式 Φₘ(X) | RLWE 环 R = Z[X]/Φₘ(X) | RLWE 安全性、明文空间 | ✅ 基础依赖 | 高 |
| 2 | 商环 Z[X]/(X^N+1) | 密文空间 R_q | 加密/解密/同态运算 | ✅ 基础依赖 | 高 |
| 3 | 多项式环不可约性 (Eisenstein) | 环无零因子 → 安全性 | RLWE 环构造 | ✅ 基础依赖 | 高 |
| 4 | 扩域 / 分裂域 | Slot = F_{p^d} 的 CRT 直积分量 | SIMD / Batching | ✅ field switching 核心 | 高 |
| 5 | Frobenius 自同构 x↦x^p | 有限域上的根轨道 → 不可约因子 F_i | Φₘ mod p 分解 → Slot 结构 | ✅ 与 slot 有关 | 高 |
| 6 | Galois 群 Gal(Q(ζₘ)/Q) ≅ (Z/mZ)^× | 自同构 σ_k: X↦X^k | Slot permutation / Rotation | ✅ rotation 操作 | 高 |
| 7 | 迹 Trace_{L/K} | Tr_{n→B} 降维操作 | Field switching / Packing / 系数降维 | ✅ field switching 核心操作 | 需确认具体映射 |
| 8 | 典范嵌入 σ: K→C^d | CKKS 编码 / 解码 | Encode / Decode | PaCo 中涉及 | 高 |
| 9 | 复嵌入 = 抽象自同构 ∘ 基础复嵌入 | σᵢ = ι ∘ τᵢ (两步分解) | CKKS 编码的槽视角 | 间接相关 | 中-高 |
| 10 | CRT (中国剩余定理) | R ≅ ⊕_{i} R/(F_i) | SIMD / Batching | ✅ slot 分解基础 | 高 |
| 11 | NTT / FFT (分奇偶 Cooley-Tukey) | 快速多项式乘法 / SlotToCoeff | 多项式乘法、编码转换 | ✅ CoeffToSlot/SlotToCoeff | 高 |
| 12 | 多项式代换 Y = X^{N/(2n)} | 从 R[Y]/(Y^{2n}+1) 嵌入到 R[X]/(X^N+1) | 编码环的统一嵌入 | ✅ packing 相关 | 高 |
| 13 | coef split: P(X)=A(X²)+X·B(X²) | 偶奇系数分离 | 当前论文 construction | ✅ 直接依赖 | 高 |
| 14 | slot-preserving split: Q⁺/Q⁻ | 保持槽义的拆分 | 当前论文 construction | ✅ 直接依赖 | 高 |
| 15 | 塔式分圆扩张 K₀ ⊂ K₁ ⊂ ... | Field switching 的代数基础 | 大环↔小环映射 | ✅ 核心依赖 | 高 |
| 16 | 对偶 (Duality) 与好基 | R^∨ 上 trace pairing | 噪声分析、基选择 | ✅ 论文 2.1.4-2.1.5 | 中 |
| 17 | 素理想分裂 pR = Π p_i^{e} | 明文模数 p 在分圆环中的分解 | 明文算术 / slot 结构 | ✅ 论文 2.1.3 | 高 |
| 18 | Gram-Schmidt 正交化 b̃_i | Babai 最近平面算法中的层间距 | CVP 近似求解 / 格基约简 | 间接 (格密码基础) | 中 |

---

## 1. 分圆多项式 → RLWE 环

### 代数表述

分圆多项式 Φₘ(X) 在 Q 上不可约。当 m = 2^{k+1} 时:
- Φ_{2^{k+1}}(X) = X^{2^k} + 1
- n = 2^k 是 2 的幂

在 Z[X] 上:
- n 为 2 的幂 → X^n + 1 不可约 (Eisenstein 判别法，取变元替换 x=y+1，素数 p=2)
- n 非 2 的幂 → X^n + 1 可约 (含因子 X+1)

### FHE 实例化

```
代数层: Z[X]/(Φₘ(X)) — 整环，无零因子
               ↓ 模 q
实例层: R_q = Z_q[X]/(X^N+1) — 这是 RLWE 的标准环
```

**技术意义**:
- 安全性: RLWE 假设的困难性依赖于这个环结构
- 计算效率: X^N ≡ -1，模约化只需一次减法
- NTT 条件: N|(q-1) 时可在 F_q 中做 NTT，复杂度 O(N log N)

**来源文件**: `01_源文件库/01_抽象代数与代数数论/为什么选择分圆多项式-好处汇总.md`, `01_源文件库/03_FHE总论与基础方案/gentry-BGFV-CKKS.md`
**与当前论文关系**: 论文基于 BGV 方案，BGV 的 RLWE 环正是 Z[X]/(X^N+1)
**置信度**: 高

---

## 2. 多项式环 / 商环 → 明文空间 / 密文空间

### 代数表述

- 明文空间: R_t = Z_t[X]/(X^N+1)
- 密文空间: R_q²
- BGV: m 在低位，c = (a·s + t·e + m, -a)
- BFV: m 在高位，c = (a·s + Δ·m + e, -a), Δ = ⌊q/t⌋

### 代数机制

```
R = Z[X]/(X^N+1)  ←── 这是统一的代数对象
   ↓ mod q               ↓ mod t
R_q = Z_q[X]/(X^N+1)   R_t = Z_t[X]/(X^N+1)
  (密文环)              (明文环)
```

### 关键区别

BGV 和 BFV 共享同一个环 R，但消息嵌入位置不同:
- BGV: 消息在 Z_t 的"低比特"，模 t 直接提取
- BFV: 消息乘以大因子 Δ = ⌊q/t⌋ 放在"高比特"，解密时除以 Δ 取整

**为什么这是重要的代数区别?**
模数切换 (BGV) 和重缩放 (BFV) 的代数含义完全不同:
- ModSwitch: 商环的同态 (改变模数 q)
- Rescale: 缩放因子的同步调整 (改变嵌入尺度)

**来源文件**: `01_源文件库/03_FHE总论与基础方案/BFV与BGV的高低位明文.md`, `01_源文件库/03_FHE总论与基础方案/模数切换和重缩放的区别.md`
**与当前论文关系**: 论文基于 BGV，使用模数切换而非重缩放
**置信度**: 高

---

## 3. 分裂域 / 扩域 → Slot 分解

### 代数表述

在 F_p 上，分圆多项式分解为:
Φₘ(X) ≡ Π_{i=1}^L F_i(X) (mod p)
其中 deg(F_i) = d = ordₘ(p), L = φ(m)/d

### FHE 实例化

```
代数: Φₘ(X) = Π F_i(X)  (在 F_p 上)
           ↓ CRT
FHE:  R_t ≅ ⊕_{i=1}^L F_p[X]/(F_i(X))  ≅ ⊕_{i=1}^L F_{p^d}
      每个直积分量 = 一个 slot
```

**slot 数量和维度的决定性公式**:
```
d = ordₘ(p) = min{k > 0 : p^k ≡ 1 (mod m)}
L = φ(m)/d
```

**Frobenius 轨道解释**:
- 每个不可约因子 F_i 的根集 = {ζ^h, ζ^{hp}, ..., ζ^{hp^{d-1}}}
- 这是一条 Frobenius 轨道
- 轨道数 = L，轨道长度 = d
- 每条轨道 = 一个 slot

**特殊情况**:
- d=1 (p ≡ 1 mod m): L=φ(m)，每个 slot = F_p，SIMD 并行度最大
- d=φ(m) (p 是模 m 的原根): L=1，只有一个 slot = F_{p^{φ(m)}}

**来源文件**: `01_源文件库/01_抽象代数与代数数论/分圆多项式的F_p分解.md`, `01_源文件库/01_抽象代数与代数数论/SIMD代数版本总结.md`
**与当前论文关系**: Field Switching 的核心就是在大环的 slot 结构和小环的 slot 结构之间建立映射
**置信度**: 高

---

## 4. Frobenius 轨道 → Slot 内部结构

### 代数表述

Frobenius 映射 Fr_p: x ↦ x^p 是 Gal(F_{p^d}/F_p) 的生成元。
对分圆域的每个本原 m 次根 ζ，其 Frobenius 轨道为:
{ζ, ζ^p, ζ^{p²}, ..., ζ^{p^{d-1}}}

### FHE 实例化

在 BGV/BFV 中，每个 slot 是 F_{p^d} (或 Galois 环的扩环)。
Frobenius 在每个 slot 内部作用于该 slot 的元素:
```
a(ζ) ↦ a(ζ^p)
```

这正是 SlotToCoeff 和 CoeffToSlot 变换中需要处理的"槽内自同构"。

**来源文件**: `01_源文件库/01_抽象代数与代数数论/SIMD代数版本总结.md`, `01_源文件库/01_抽象代数与代数数论/分圆多项式的F_p分解.md`
**与当前论文关系**: 在 field switching 中，需要理解不同环中 Frobenius 轨道如何对应
**置信度**: 高

---

## 5. 自同构 (Automorphism) → Slot Permutation / Rotation

### 代数表述

Gal(Q(ζₘ)/Q) ≅ (Z/mZ)^× (取 m 次本原根 ζₘ)

自同构 τ_k: ζₘ ↦ ζₘ^k (k ∈ Z_m^×)

### FHE 实例化

在 CKKS 中，自同构 σ_k: X ↦ X^k 对应于:
- 密文层面: 旋转操作 Rot_j(ct)
- 明文层面: 对复向量的槽位进行置换

关键公式 (CKKS 中):
- σ_{1+2n}: X ↦ X^{1+2n} → Y = X^{N/(2n)} 映射到 -Y
- 这用于偶奇系数分解: P(-Y) = σ_{1+2n}(P(Y))

### 自同构的三个层次

```
抽象代数层: Gal(Q(ζₘ)/Q) ≅ (Z/mZ)^×
    ↓ 限制到分圆环
环自同构层: σ_k: R → R, X ↦ X^k
    ↓ CKKS 的同态实现
操作层: Rot_j, Conj (需要自同构密钥)
```

**在 Field Switching (GHPS) 中**:
- trace = Σ_{σ∈Gal(L/K)} σ(x) 是 Galois 自同构的和
- 通过自同构 + 加法的组合实现降维

**来源文件**: `01_源文件库/01_抽象代数与代数数论/自同构的个数.md`, `01_源文件库/01_抽象代数与代数数论/共轭-多项式角度.md`
**与当前论文关系**: ✅ 核心: 论文的 trace/旋转操作基于 Galois 自同构
**置信度**: 高

---

## 6. Trace → Field Switching / 降维 / Packing

### 代数表述

设 L/K 是域扩张，[L:K] = n。
迹 Tr_{L/K}: L → K 定义为:
Tr_{L/K}(x) = Σ_{σ: L→K̄, σ|_K=id} σ(x)

迹的值落在基域 K 中，且 Tr 是 K-线性的。

### FHE 实例化

在 Field Switching (GHPS) 中:
- 大环 ↔ 扩域 L
- 小环 ↔ 子域 K
- Tr_{L/K} 将扩域元素映射到子域

同态计算: Tr_{n→B} = (id + Rot_B) ∘ ... ∘ (id + Rot_{n/2})
- 复杂度: O(log(n/B)) 次旋转+加法
- 效果: 将 n 个槽的向量降维为 B 个槽的向量

### 迹落入基域的证明

核心论证: 对任何 τ ∈ Gal(L/K),
τ(Tr(x)) = τ(Σ_σ σ(x)) = Σ_σ (τ∘σ)(x) = Σ_σ σ(x) = Tr(x)

因为所有自同构的复合只是对求和指标的置换。
所以 Tr(x) 被 Gal(L/K) 固定 → Tr(x) ∈ K (Galois 理论: 固定域 = K)。

**来源文件**: `01_源文件库/01_抽象代数与代数数论/Trace与固定.md`, `01_源文件库/01_抽象代数与代数数论/复嵌入的两个步骤.md`
**与当前论文关系**: ✅ 核心: trace 是 field switching 的代数基础
**置信度**: 中-高 (需确认笔记中的同态实现是否对应原始 GHPS 论文)
**需人工确认**: trace 在 CKKS 密文上的同态实现公式是否为原始 GHPS 论文的标准表述

---

## 7. 典范嵌入 / 复嵌入 → CKKS Encoding

### 代数表述

设 K = Q(ζₘ)，[K:Q] = d = φ(m)。

典范嵌入 σ: K → C^d:
σ(a) = (σ_1(a), ..., σ_d(a))

其中 σ_i: K → C 是 d 个嵌入。

### 两步分解

σ_i = ι ∘ τ_i:
- τ_i: 抽象自同构 (Galois 群元素)
- ι: 基础复嵌入 (选定一个"根"的复数值)

### FHE 实例化

CKKS 的编码:
1. 取复向量 z ∈ C^{N/2}
2. 通过 τ_n^{-1}: C^{N/2} → R[Y]/(Y^{N/2}+1) 完成逆典范嵌入
3. m = ⌊Δ · τ_n^{-1}(z)⌋ 得整数多项式
4. 通过 Y = X^{N/(2n)} 嵌入到大环

**关键公式**:
τ_n: R[Y]/(Y^{2n}+1) → C^n, p ↦ (p(ζ_{n,i}))_{i∈[0,n)}
其中 ζ_{n,i} = exp(2πI · 5^i / (4n))

**来源文件**: `01_源文件库/01_抽象代数与代数数论/复嵌入的两个步骤.md`, `01_源文件库/04_FHE核心技术模块/典范嵌入极简与U矩阵表示推导.md`
**与当前论文关系**: CKKS 的编码基础; PaCo 论文中涉及多次编码操作
**置信度**: 高

---

## 8. CRT 分解 → SIMD / Batching

### 代数表述

中国剩余定理: 若理想 I_i 两两互素，则
R/(Π I_i) ≅ ⊕ R/I_i

### FHE 实例化

在 BGV/BFV 中:
```
R_t = Z_t[X]/(X^N+1)  (N = 2^k)
     ≅ ⊕_{i=1}^L Z_t[X]/(F_i(X))  (F_i 是 X^N+1 在 Z_t 上的不可约因子)
```

**BGV 的 slot 参数**:
- d = ord_{2N}(t) = min{k: t^k ≡ 1 (mod 2N)}
- L = φ(2N)/d = N/d
- 每个 slot 的结构: 若 t = p^r (素数幂) 则为 Galois 环 GR(p^r, d)

**CKKS 的 "slot"**:
CKKS 没有 CRT 意义上的 slot (因为 CKKS 在范数意义上工作于实数/复数)。
CKKS 的 "slot" 来自典范嵌入的分量: 每个嵌入分量 σ_i(a) 提供了一个"slot 视角"。

**来源文件**: `01_源文件库/01_抽象代数与代数数论/Prime Spliting.md`, `01_源文件库/01_抽象代数与代数数论/SIMD代数版本总结.md`
**与当前论文关系**: ✅ 核心: 论文需要在不同 slot 配置的两个环之间做 field switching
**置信度**: 高

---

## 9. NTT / FFT → SlotToCoeff / CoeffToSlot

### 代数表述

离散傅里叶变换 DFT_n: C^n → C^n, (a_i) ↦ (Σ a_i ω^{ik})
Cooley-Tukey: 分奇偶递归 → O(n log n)
蝴蝶: x' = x + ω^k y, y' = x - ω^k y

### FHE 实例化

DFT 在 F_q 上 → NTT (ω 是 q-1 次单位根中的 n 次本原根)
NTT 的有限域版本用于 BGV/BFV 的多项式乘法和 Slot-to-Coeff 变换。

**具体到 PaCo 的 SlotToCoeff**:
- U_n = Π^{(n)} · (Π_{ℓ=0}^{log n - 1} D_{n, n/2^{ℓ+1}})  (Cooley-Tukey 分解)
- 每个 D_{n,2^ℓ} 是三对角的块对角矩阵
- 比特反转 Π^{(n)} 被"吸收"到 E 矩阵中
- 同态实现通过与稀疏矩阵的明文-密文乘法

**来源文件**: `01_源文件库/01_抽象代数与代数数论/变换与NTT-Cooley-Tukey.md`, `PaCo翻译§4`
**与当前论文关系**: ✅ 论文中 SlotToCoeff/CoeffToSlot 依赖 NTT 的矩阵分解
**置信度**: 高

---

## 10. 多项式代换 Y = X^{N/(2n)} → 环嵌入与自同构

### 代数表述

设 k = N/(2n)，则 Y = X^k。
- Y^{2n} = X^{2nk} = X^N = -1 → Y 满足 Y^{2n} + 1 = 0
- 大小环关系: R[Y]/(Y^{2n}+1) ↪ R[X]/(X^N+1) ← 这是一个环嵌入

### FHE 实例化

在 CKKS 中:
1. 先将 n 维复向量编码到小环 Z[Y]/(Y^{2n}+1)
2. 通过 Y = X^k 嵌入到大环 Z[X]/(X^N+1)
3. 加密在大环中

**偶奇分解中的应用**:
- σ_{1+2n}(X^k) = X^{k(1+2n)} = X^{k+2nk} = -X^k → σ_{1+2n}(Y) = -Y
- 因此: ct_- = σ_{1+2n}(ct_P) 得到加密 P(-Y) 的密文
- 然后: ct_even = (ct_P + ct_-)/2 得到 P_even(Y^2) 的密文

**来源文件**: `01_源文件库/04_FHE核心技术模块/多项式幂次代换与提升.md`, `01_源文件库/04_FHE核心技术模块/2-正向.md`
**与当前论文关系**: ✅ 直接: 论文的 polynomial decomposition 基于此代换
**置信度**: 高

---

## 11. Coefficient Split vs Slot-Preserving Split → 当前论文核心

### 三种多项式拆分

**(1) Coefficient Split (系数拆分)**:
P(X) = A(X²) + X · B(X²)
在 CKKS 中间态完成 (Galois 自同构 + 加减):
ct_- = σ_{1+2n}(ct_P)
ct_A = (ct_P + ct_-) / 2  → 加密 A(X²)
ct_B = (ct_P - ct_-) / 2  → 加密 X·B(X²)  (可选去 X 得 B(X²))

**(2) Slot-Preserving Split (保持槽义的拆分)**:
Q^+(Y) = A(Y) + h(Y)·B(Y)
Q^-(Y) = A(Y) - h(Y)·B(Y)
要求 h(θ_j) = α_j (在槽 j 上的值匹配)

**(3) 一般 Slot Selection**:
q = F_n^{-1} · S · F_N · p
其中 F_N IF DFT 矩阵, S 是选择矩阵 (选出想要的槽)

### 矩阵视角

| 变换 | 矩阵 | 密集度 | 级别消耗 |
|------|------|--------|---------|
| Coeff split | 对角/稀疏 | 低 | 0-1 |
| Slot-preserving split | 对角+槽自同构 | 中 | 1-2 |
| 一般 slot selection | 密集 | 高 | 多 |

**来源文件**: `01_源文件库/04_FHE核心技术模块/补充视角-slot情况下的拆分.md`, `01_源文件库/04_FHE核心技术模块/1-总览.md`, `01_源文件库/04_FHE核心技术模块/2-正向.md`
**与当前论文关系**: ✅ 核心: 论文的分解步骤直接基于此
**置信度**: 高

---

## 12. 关联强度总结

```
最高关联 (当前论文直接依赖):
├── trace → field switching 降维
├── Frobenius 轨道 → slot 结构
├── coefficient/slot split → 论文 construction
├── 多项式代换 → 环嵌入
├── tower of cyclotomics → 大小环转换
└── CRT → SIMD slot 分解

高关联 (FHE 基础):
├── 分圆多项式 → RLWE 环
├── 自同构 → rotation
├── NTT → SlotToCoeff/CoeffToSlot
├── 典范嵌入 → CKKS encoding
└── prime splitting → slot 参数

中关联 (间接依赖):
├── Gram-Schmidt → Babai CVP
├── 对偶 → trace pairing
└── 复嵌入两步分解 → 编码理论
```


> **pending_import**: 'GHPS 2013 Field Switching 论文' (GHPS13) 已于 2026-05-04 导入源文件库 (src_000028: GHPS 2013, src_000029: PaCo)。

> **note**: 关系图中的 'PaCo翻译§N' 引用指向 01_源文件库/04_FHE核心技术模块/翻译.md 的对应章节，section reference 而非文件路径。