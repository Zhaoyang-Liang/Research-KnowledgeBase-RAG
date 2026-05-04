好的，我将以**数学公式**为核心，系统地为你整理全同态加密的基础概念与主要方案的结构。每个部分都会先给出核心的数学定义，再说明其作用。力求简洁、准确、可推导。

---

## 1. LWE 与 RLWE 问题

### 1.1 LWE (Learning with Errors)

**参数**：  
- \(n\)：安全参数（维度）  
- \(q\)：模数（整数）  
- \(\chi\)：误差分布（常为离散高斯分布，标准差 \(\sigma\)）

**LWE 假设**：  
给定 \(m\) 个样本 \((a_i, b_i) \in \mathbb{Z}_q^n \times \mathbb{Z}_q\)，其中  
\[
b_i = \langle a_i, \mathbf{s} \rangle + e_i \pmod{q}
\]  
- \(\mathbf{s} \in \mathbb{Z}_q^n\) 是固定秘密（均匀随机）  
- \(a_i \leftarrow \mathbb{Z}_q^n\) 均匀随机  
- \(e_i \leftarrow \chi\) 小误差  

**计算性问题**：  
从样本中恢复 \(\mathbf{s}\) 是困难的（即使量子计算机）。  
**判定性 LWE**：区分 \((a_i, b_i)\) 与均匀随机的 \((a_i, u_i)\)。

### 1.2 RLWE (Ring-LWE)

**环**：  
\[
R = \mathbb{Z}[X] / (X^N + 1), \quad R_q = R / qR
\]  
其中 \(N\) 是 2 的幂，\(q\) 是整数。

**RLWE 样本**：  
\[
(a(x), b(x) = a(x) \cdot s(x) + e(x)) \in R_q \times R_q
\]  
- \(a(x) \leftarrow R_q\) 均匀  
- \(s(x) \leftarrow R_q\) 秘密（小系数）  
- \(e(x) \leftarrow \chi\) 小系数多项式

**优势**：密钥和密文尺寸从 \(O(n)\) 降到 \(O(N)\)，且乘法可用 NTT 快速计算。

---

## 2. 自举 (Bootstrapping) 的核心思想

设 \(\mathsf{Dec}_{sk}(c)\) 是解密函数，**同态解密**的目标是：给定加密后的密文 \(\overline{c} = \mathsf{Enc}_{pk}(c)\) 和加密后的密钥 \(\overline{sk} = \mathsf{Enc}_{pk}(sk)\)，计算：
\[
\mathsf{Enc}_{pk}\big(\mathsf{Dec}_{sk}(c)\big)
\]  
结果就是新鲜加密的明文，噪声重置。

**数学形式**（以简单 LWE 解密为例）：

LWE 密文：\(c = (a, b = \langle a, s \rangle + m \cdot \lfloor q/2 \rfloor + e)\)，解密为：
\[
m' = \mathsf{round}\left( \frac{2}{q} (b - \langle a, s \rangle) \right)
\]  
将其写作算术电路（含取整、乘除），然后同态执行该电路。  
Bootstrapping 的核心就是**同态计算这个解密电路**。

---

## 3. Gentry 09 方案（概念性结构）

### 3.1 Somewhat Homomorphic Encryption (SWHE)

加密：  
\[
c = m + 2e + \sum_{i} f_i \cdot r_i \quad (\text{理想格表示简化})
\]  
其中 \(f_i\) 是公钥中的元素，\(r_i\) 随机。  
解密：\(m = (c \mod 2)\)。

### 3.2 压缩 (Squashing)

将解密电路深度降低，使其可被 SWHE 同态计算。  
引入一个“提示向量” \(\mathbf{z}\)，使得解密变为：
\[
m = \sum_{j=1}^t z_j \cdot w_j \mod 2
\]  
其中 \(w_j\) 是密文的某个函数，深度很浅。

### 3.3 Bootstrapping 公式

设 \(\mathsf{Dec}_{sk}(c)\) 的解密电路深度为 \(D\)，SWHE 能计算深度 \(L\) 的电路。  
若 \(L \ge D\)，则：
\[
\mathsf{Enc}_{pk}(m) \xrightarrow{\text{同态解密}} \mathsf{Enc}_{pk}(m)
\]  
噪声重置。反复执行即可获得无限深度。

---

## 4. BGV 方案 (Brakerski-Gentry-Vaikuntanathan)

### 4.1 基本加密

明文空间：\(R_t = \mathbb{Z}_t[X]/(X^N+1)\)，\(t\) 为明文模数（小）。  
密文空间：\(R_q\) 上的向量（长度 2 或更多）。

加密 \(m \in R_t\)：
\[
c = (c_0, c_1) = (a \cdot s + t e + m, \; -a) \in R_q^2
\]  
其中 \(a \leftarrow R_q\)，\(e \leftarrow \chi\)。

解密：
\[
m = \big[ \langle c, (1, s) \rangle \big]_q \mod t
\]  
即：
\[
\mu = c_0 + c_1 \cdot s \pmod{q}
\]  
\[
m = \mu \mod t
\]

### 4.2 同态加法与乘法

**加法**：对应分量相加。  
**乘法**：张量积后重线性化（key switching）。

设 \(c = (c_0, c_1)\)，\(d = (d_0, d_1)\)，则乘积密文（未线性化）为：
\[
c_{\times} = (c_0 d_0, \; c_0 d_1 + c_1 d_0, \; c_1 d_1)
\]  
维数增长，需用密钥切换回 2 维。

### 4.3 模数切换 (Modulus Switching)

给定密文 \(c \in R_q^2\)，噪声 \(\|e\|\)。切换到更小的模数 \(q'\)：
\[
c' = \left\lfloor \frac{q'}{q} c \right\rceil \in R_{q'}^2
\]  
效果：噪声约化为 \(\|e'\| \approx \frac{q'}{q} \|e\|\)。  
这使得噪声增长被控制，实现分层同态。

---

## 5. BFV 方案 (Brakerski-Fan-Vercauteren)

### 5.1 加密与解密

明文空间 \(R_t\)，密文空间 \(R_q^2\)。  
加密：
\[
c = (c_0, c_1) = \big( a \cdot s + \Delta m + e, \; -a \big) \in R_q^2
\]  
其中 \(\Delta = \lfloor q/t \rfloor\)（缩放因子）。

解密：
\[
\mu = c_0 + c_1 \cdot s \pmod{q}
\]  
\[
m = \left\lfloor \frac{t}{q} \mu \right\rceil
\]

### 5.2 乘法与重缩放 (Rescaling)

乘积同 BGV，得到三维密文 \(c_{\times} = (c_0', c_1', c_2')\)。  
重线性化（key switching）后得到二维密文 \(c_{\text{mult}}\)，但此时密文模数仍为 \(q\)，而噪声变大。

重缩放：除以缩放因子 \(\Delta\)（或更一般地，除以 \(q\) 的一部分）：
\[
c_{\text{resc}} = \left\lfloor \frac{q'}{q} c_{\text{mult}} \right\rceil
\]  
其中 \(q' = q / p\)。BFV 中常用 \(p = t\)，即缩放后模数变为 \(q/t\)，噪声相应减小。

---

## 6. CKKS 方案 (Cheon-Kim-Kim-Song)

### 6.1 编码 (Encoding)

将复数向量 \(\mathbf{z} \in \mathbb{C}^{N/2}\) 编码为多项式 \(m(X) \in R_q\)。  
使用规范嵌入（canonical embedding）与逆 FFT：
\[
m = \mathsf{Encode}(\mathbf{z}, \Delta)
\]  
其中 \(\Delta\) 是缩放因子（保留精度）。

### 6.2 加密与解密

与 BFV 形式相同：
\[
c = (c_0, c_1) = (a \cdot s + m + e, \; -a)
\]  
注意这里没有乘以 \(\Delta\)（因为明文已经包含缩放）。

解密：
\[
\mu = c_0 + c_1 \cdot s \pmod{q}
\]  
\[
\mathbf{z}' = \mathsf{Decode}(\mu)
\]  
解码得到近似复数向量 \(\mathbf{z}' \approx \mathbf{z}\)。

### 6.3 乘法与重缩放

乘法后同样需要重线性化。  
重缩放：
\[
c' = \left\lfloor \frac{1}{p} c \right\rceil \in R_{q/p}^2
\]  
其中 \(p\) 是缩放因子（例如 \(p=2^{L}\)）。  
这相当于保持浮点精度，同时控制噪声增长。

---

## 7. 总结对照表（数学视角）

| 方案 | 密文形式 | 解密公式 | 噪声控制 | 明文域 |
|------|----------|----------|----------|--------|
| Gentry09 | 理想格元素 | \(m = (c \mod 2)\) | Bootstrapping | 比特 |
| BGV | \(c = (a s + t e + m, -a)\) | \(m = (c_0 + c_1 s \mod q) \mod t\) | 模数切换 | \(R_t\) |
| BFV | \(c = (a s + \Delta m + e, -a)\) | \(m = \lfloor t/q \cdot (c_0 + c_1 s) \rceil\) | 重缩放 | \(R_t\) |
| CKKS | \(c = (a s + m + e, -a)\) | \(\mathbf{z} = \mathsf{Decode}(c_0 + c_1 s)\) | 重缩放 | 复数向量 |

---

如果需要某个具体运算的完整推导（比如 BGV 的密钥切换公式、CKKS 的编码具体矩阵），我可以继续补充。


---


下面我们以之前的总结对照表为基础，扩展为一个更全面的对比表格，并补充各方案的区别、优劣与特点。

---

## 扩展对比表：Gentry09, BGV, BFV, CKKS

| 特性 | Gentry09 | BGV | BFV | CKKS |
|------|----------|-----|-----|------|
| **提出年份** | 2009 | 2012 | 2012 (Brakerski) / 2013 (Fan-Vercauteren) | 2017 |
| **数学基础** | 理想格（Ideal Lattices） | RLWE | RLWE | RLWE |
| **密文形式** | 理想格元素（单元素） | 向量 `(c₀, c₁) ∈ R_q²` | 向量 `(c₀, c₁) ∈ R_q²` | 向量 `(c₀, c₁) ∈ R_q²` |
| **明文域** | 比特（0/1） | `R_t = ℤ_t[X]/(X^N+1)` (整数多项式) | `R_t` (整数多项式) | 复数向量 `ℂ^{N/2}` |
| **编码方式** | 直接映射比特到理想格元素 | 系数映射（幂基） | 系数映射（高位放置） | 规范嵌入 + FFT 编码 |
| **加密公式** | 复杂理想格运算 | `c₀ = a·s + t·e + m`, `c₁ = -a` | `c₀ = a·s + Δ·m + e`, `c₁ = -a`，Δ = ⌊q/t⌋ | `c₀ = a·s + m + e`, `c₁ = -a` |
| **解密公式** | `m = (c mod 2)` | `m = [ [c₀ + c₁·s]_q ]_t` | `m = ⌊ t/q · [c₀ + c₁·s]_q ⌉` | `z = Decode(c₀ + c₁·s mod q)` |
| **噪声控制** | 仅靠 Bootstrapping | 模数切换 (Modulus Switching) | 重缩放 (Rescaling) | 重缩放 (Rescaling) |
| **是否精确** | 精确（比特） | 精确（整数多项式） | 精确（整数多项式） | 近似（允许误差） |
| **缩放因子** | 无 | 无（明文在低位） | 有 Δ = ⌊q/t⌋ | 有（编码时固定） |
| **乘法后维数** | 保持 | 需重线性化 (2→3→2) | 需重线性化 | 需重线性化 |
| **密钥切换** | 需要 | 需要（标准操作） | 需要 | 需要 |
| **Bootstrapping 复杂度** | 极高（首次实现） | 中等（有优化） | 中等（类似BGV） | 较低（因近似计算） |
| **典型应用** | 理论证明 | 整数算术、金融、精确统计 | 整数算术、与BGV类似 | 机器学习、科学计算、近似统计 |
| **实现库** | 无实用实现 | HElib, PALISADE, OpenFHE | SEAL, OpenFHE, Lattigo | SEAL, HElib (CKKS), OpenFHE, Lattigo |
| **性能特点** | 极慢，仅理论 | 适合大模数、大量插槽 | 常数模数，重缩放开销较小 | 精度可调，速度快，适合实数和复数 |

---

## 各方案的区别、优劣、特点详解

### 1. Gentry09
- **区别**：历史上第一个FHE方案，基于理想格，不使用RLWE。
- **优势**：开创性，证明了FHE的存在性。
- **劣势**：效率极低（解密电路深度大），无法实用。
- **特点**：引入了自举（Bootstrapping）和压缩（Squashing）两个核心思想。

### 2. BGV (Brakerski-Gentry-Vaikuntanathan)
- **区别**：明文放在系数**低位**，噪声增长通过**模数切换**控制。
- **优势**：模数切换能有效降低噪声绝对值；支持大量SIMD槽（通过CRT分解）；适合大整数运算。
- **劣势**：需要多个模数链，参数调整复杂；模数切换会消耗模数深度。
- **特点**：HElib 实现的主力方案；与BFV相比，乘法后需模数切换，但重线性化成本略低。

### 3. BFV (Brakerski-Fan-Vercauteren)
- **区别**：明文放在系数**高位**（乘以Δ），噪声控制通过**重缩放**（除以Δ）。
- **优势**：密文模数q在运算过程中保持不变（尺度不变），更易实现；重缩放操作简单，适合硬件实现。
- **劣势**：乘法后重缩放会引入舍入误差（但仍然是精确整数运算，误差可控）；需要精确管理Δ。
- **特点**：SEAL库默认方案；与BGV相比，在相同安全参数下性能相近，但实现更直观。

### 4. CKKS (Cheon-Kim-Kim-Song)
- **区别**：明文是**复数向量**，结果允许近似；使用**规范嵌入**编码。
- **优势**：效率极高，支持浮点运算，适合深度学习、科学计算；重缩放自然对应浮点精度控制。
- **劣势**：不是精确的，只能近似计算；参数选择更复杂（需要平衡精度和性能）。
- **特点**：是目前最受欢迎的FHE方案之一（尤其在ML领域）；自举效率高于BGV/BFV（因近似特性）。

---

## 直观对比（一句话总结）

| 方案 | 一句话特点 |
|------|------------|
| Gentry09 | 理论奠基，极慢，不实用 |
| BGV | 整数精确，模数切换控噪，适合大整数 |
| BFV | 整数精确，尺度不变，实现简单 |
| CKKS | 复数近似，快速高效，适合ML |

---

如果需要进一步比较它们的**自举方法差异**（例如BGV自举中的NTT分解 vs CKKS自举中的FFT），我可以继续补充。