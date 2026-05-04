# C2S / S2C 术语方向说明：以 PaCo 为准的统一规则

> 目的：修复知识库中 `SlotToCoeff(C2S)` / `CoeffToSlot(S2C)` 等反向写法，避免后续 RAG 和论文协作继续混淆。
>
> 结论先行：
>
> - **CoeffToSlot，简称 C2S**：coefficient representation → slot representation。
> - **SlotToCoeff，简称 S2C**：slot representation → coefficient representation。
> - 但是在 CKKS / PaCo 中，`U_n` 或 `U_n^{-1}` 的出现不能单独决定操作名，因为操作名取决于"语义层输入输出"，而不只取决于某个矩阵在中间向量空间里的方向。

---

## 0. 一句话规范

统一采用如下命名：

\[
\boxed{\mathrm{CoeffToSlot}\;(\mathrm{C2S}):\quad \text{coefficient data} \longrightarrow \text{slot data}}
\]

\[
\boxed{\mathrm{SlotToCoeff}\;(\mathrm{S2C}):\quad \text{slot data} \longrightarrow \text{coefficient data}}
\]

因此，知识库中如果出现：

```text
SlotToCoeff(C2S)
CoeffToSlot(S2C)
```

应视为**主术语方向错误**。除非它明确标在 `original_note_convention` 或"原笔记旧命名"字段中，否则不应作为 validated / ready 知识进入 RAG。

---

## 1. 为什么 PaCo 支持这个命名？

PaCo 在介绍传统 CKKS bootstrapping 时，明确把流程分成：

\[
\mathrm{ModRaise}\;\rightarrow\;\mathrm{CoeffToSlot}\;\rightarrow\;\mathrm{EvalMod}\;\rightarrow\;\mathrm{SlotToCoeff}.
\]

其语义是：

1. **CoeffToSlot**：把 ModRaise 后多项式明文
 \[
 m'=m+qJ=\sum_{i=0}^{N-1}(m_i+qJ_i)X^i
 \]
 的系数搬到 slot 中，使 slot 中保存整数 \(m_i+qJ_i\)。

2. **EvalMod**：在 slot 中逐点近似去掉 \(qJ_i\)，得到近似的 \(m_i\)。

3. **SlotToCoeff**：把 slot 中近似得到的系数 \(m_i\) 搬回明文多项式的 coefficient 位置，得到刷新后的 ciphertext。

所以从 bootstrapping 语义看：

\[
\mathrm{CoeffToSlot}:\quad (m_i+qJ_i)_i\;\text{as polynomial coefficients}
\longmapsto
(m_i+qJ_i)_i\;\text{as encrypted slots},
\]

\[
\mathrm{SlotToCoeff}:\quad (m_i)_i\;\text{as encrypted slots}
\longmapsto
\sum_i m_iX^i\;\text{as plaintext polynomial}.
\]

这正是 `CoeffToSlot = coefficient → slot`，`SlotToCoeff = slot → coefficient`。

---

## 2. CKKS 编码/解码与 C2S/S2C 不是一回事

PaCo 的 CKKS preliminaries 使用如下映射。令

\[
R_n = \mathbb{R}[Y]/(Y^{2n}+1),
\]

并定义典范嵌入/求值映射

\[
\tau_n:
R_n\longrightarrow \mathbb{C}^n,
\qquad
p\longmapsto (p(\zeta_{n,i}))_{i\in[0,n)}.
\]

CKKS 编码为

\[
\mathrm{Ecd}(z)=\left\lfloor \Delta\tau_n^{-1}(z)\right\rceil
\in \mathbb{Z}[Y]/(Y^{2n}+1),
\]

解码为

\[
\mathrm{Dcd}_n(m)=\frac{\tau_n(m)}{\Delta}\in\mathbb{C}^n.
\]

这里：

- CKKS **encoding** 是外部消息向量 \(z\in\mathbb{C}^n\) 到 plaintext polynomial 的构造；
- CKKS **decoding** 是 plaintext polynomial 到近似 complex vector 的读取；
- C2S / S2C 是 bootstrapping 内部的 homomorphic linear transform，用来在"系数数据"和"slot 数据"之间搬运同一个 plaintext 信息。

因此不要把下面两组概念混为一谈：

| 层次 | 操作 | 方向 |
|---|---|---|
| CKKS 编码层 | encoding \(\mathrm{Ecd}\) | external complex vector \(z\) → plaintext polynomial |
| CKKS 解码层 | decoding \(\mathrm{Dcd}\) | plaintext polynomial → complex vector |
| bootstrapping 线性变换层 | CoeffToSlot / C2S | polynomial coefficients → encrypted slots |
| bootstrapping 线性变换层 | SlotToCoeff / S2C | encrypted slots → polynomial coefficients |

---

## 3. 矩阵 \(U_n\) 的方向：为什么它容易误导？

PaCo 第 3.6 节回顾 conventional SlotToCoeff。给定

\[
p(Y)=\sum_{i=0}^{2n-1}p_iY^i=\tau_n^{-1}(z)
\in \mathbb{R}[Y]/(Y^{2n}+1),
\]

把系数拆成两半：

\[
p_0=(p_i)_{i\in[0,n)},
\qquad
p_1=(p_{i+n})_{i\in[0,n)}.
\]

定义 Vandermonde 矩阵

\[
U_n=(\zeta_{n,i}^{\,j})_{i,j\in[0,n)}.
\]

则有

\[
\boxed{z=U_n(p_0+I p_1)}.
\]

从裸矩阵角度看：

\[
U_n: \text{coefficient vector}\longrightarrow\text{evaluation/slot vector}.
\]

这很容易让人误以为：

```text
用了 U_n，所以一定是 CoeffToSlot / C2S。
```

但 PaCo 在 conventional SlotToCoeff 中的语义是：输入 ciphertext 的 slots 里已经装着 coefficient vector \(p_0,p_1\)，通过同态矩阵乘法后得到一个 ciphertext，它在 CKKS 语义下加密 \(z\)，而 \(z=\tau_n(p)\) 正是 polynomial \(p\) 的 CKKS slot view；同时这个 ciphertext 也可记作 \(\mathrm{Enc}(m)\)，即它恢复为以 \(p_i\) 为 coefficients 的 plaintext polynomial。

所以这里有两个不同层次：

| 层次 | 看见的方向 | 解释 |
|---|---|---|
| 矩阵层 | \(U_n: p_0+Ip_1\mapsto z\) | coefficient vector → evaluation vector |
| 操作语义层 | SlotToCoeff | input slots contain coefficients → output ciphertext represents coefficient polynomial |

这就是混乱的根源：**操作名不是单纯由 \(U_n\) 或 \(U_n^{-1}\) 决定，而是由 bootstrapping 语义中的输入输出决定。**

---

## 4. PaCo 中 partial CoeffToSlot 的方向

PaCo 的新方法核心叫 **Partial CoeffToSlot**。它不是 full CoeffToSlot，而是只使用 CoeffToSlot 分解中的一部分矩阵。

PaCo 的目标之一是：在 blind rotation 之后，让某些多项式

\[
b'_v(Z)=\sum_{i=0}^{2B-1}b'_{v,i}Z^i
\]

的 coefficients \(b'_{v,i}\) 出现在 encrypted vector 的 slots 中。

也就是说 partial CoeffToSlot 的语义是：

\[
\boxed{
(b'_{v,i})_{v,i}\;\text{as polynomial coefficients}
\longrightarrow
(b'_{v,i})_{v,i}\;\text{in encrypted slots}
}
\]

因此它仍然是 CoeffToSlot 家族：

\[
\mathrm{partial\;CoeffToSlot}\subseteq \mathrm{CoeffToSlot\;style\;maps}.
\]

PaCo 后续再利用这些 slot 中的 \(b'_{v,i}\) 计算近似系数 \(\tilde m_i\)，最后用 conventional/full SlotToCoeff 把它们放回 plaintext polynomial。

PaCo 高层流程可写为：

\[
\text{coefficient data}
\xrightarrow{\text{partial CoeffToSlot}}
\text{slot data containing }b'_{v,i}
\xrightarrow{\text{slotwise operations}}
\text{slot data containing }\tilde m_i
\xrightarrow{\text{SlotToCoeff}}
\text{coefficient polynomial }\tilde m.
\]

---

## 5. 推荐在知识库中采用的字段规范

为了避免以后继续被 \(U_n\)、\(U_n^{-1}\)、encoding/decoding、C2S/S2C 四套语言搅乱，建议每张相关卡片同时写三个字段。

### 5.1 主术语字段

```json
{
 "canonical_operation": "CoeffToSlot",
 "canonical_abbrev": "C2S",
 "semantic_direction": "coefficient representation -> slot representation"
}
```

或：

```json
{
 "canonical_operation": "SlotToCoeff",
 "canonical_abbrev": "S2C",
 "semantic_direction": "slot representation -> coefficient representation"
}
```

### 5.2 矩阵实现字段

```json
{
 "matrix_layer": {
 "matrix": "U_n or U_n^{-1}",
 "vector_direction": "coefficient vector -> evaluation vector / evaluation vector -> coefficient vector",
 "warning": "matrix direction alone does not determine the semantic operation name"
 }
}
```

### 5.3 原笔记旧命名字段

如果旧笔记用了相反命名，不删除，但必须降级为历史信息：

```json
{
 "original_note_convention": {
 "old_name": "SlotToCoeff正向",
 "old_interpretation": "U_n maps coefficient vector to slot/evaluation vector",
 "status": "legacy_note_only_not_canonical"
 }
}
```

---

## 6. 审计规则

### 6.1 必须修正的表达

以下表达在 canonical 主字段中一律错误：

```text
SlotToCoeff(C2S)
CoeffToSlot(S2C)
SlotToCoeff = coefficient -> slot
CoeffToSlot = slot -> coefficient
```

应改为：

```text
CoeffToSlot(C2S): coefficient representation -> slot representation
SlotToCoeff(S2C): slot representation -> coefficient representation
```

### 6.2 可以保留但必须降级的表达

以下表达可以保留在 `original_note_convention` 或 `matrix_layer` 中，但不能作为主结论：

```text
U_n maps coefficient vector to slot/evaluation vector.
U_n^{-1} maps slot/evaluation vector to coefficient vector.
```

因为它们描述的是 **matrix layer**，不是直接描述 **semantic operation layer**。

---

## 7. 对 K019/K020/K044/K046 的建议

### K019 / K020

如果 K019/K020 当前依据是

\[
\text{slot}=U_n\cdot\text{coeff},
\qquad
\text{coeff}=U_n^{-1}\cdot\text{slot},
\]

则它们应明确标注：

- 这是 `matrix_layer` 的方向；
- 不等同于 PaCo/bootstrapping 语义中的 operation 名字；
- 若要进入 RAG ready，应重写成"矩阵方向说明卡"，而不是直接命名为 SlotToCoeff/CoeffToSlot。

### K044 / K046

如果 K044/K046 是 bootstrapping 语义卡，则应采用：

```text
CoeffToSlot(C2S): polynomial coefficients are moved into encrypted slots.
SlotToCoeff(S2C): encrypted slots holding coefficients are moved back to plaintext polynomial coefficients.
```

尤其要删除：

```text
SlotToCoeff(C2S)
CoeffToSlot(S2C)
```

---

## 8. 最终可复制到 AI 协作上下文的短规则

```text
C2S/S2C canonical convention:

- C2S = CoeffToSlot = coefficient representation -> slot representation.
- S2C = SlotToCoeff = slot representation -> coefficient representation.

PaCo evidence:
- Conventional CKKS bootstrapping applies CoeffToSlot after ModRaise to move polynomial coefficients mi+qJi into encrypted slots.
- EvalMod acts slotwise.
- Then SlotToCoeff moves the approximate coefficients mi from encrypted slots back into plaintext polynomial coefficients.
- PaCo's partial CoeffToSlot is a partial version of the CoeffToSlot family: it embeds coefficients of rotated polynomials into slots.
- The final PaCo step is a conventional/full SlotToCoeff, returning from slot-held coefficients to a bootstrapped plaintext polynomial.

Important caveat:
- U_n maps coefficient vectors to evaluation/slot vectors at the matrix layer.
- However, operation names should be determined by bootstrapping semantic input/output, not by matrix direction alone.
- Therefore, any card saying SlotToCoeff(C2S) or CoeffToSlot(S2C) in canonical fields must be corrected or deferred.
```

---

## 9. 推荐结论

在你的知识库和论文协作区中，建议统一写成：

\[
\boxed{\mathrm{C2S}=\mathrm{CoeffToSlot}:\text{coefficients}\rightarrow\text{slots}}
\]

\[
\boxed{\mathrm{S2C}=\mathrm{SlotToCoeff}:\text{slots}\rightarrow\text{coefficients}}
\]

同时额外保留一个技术备注：

> 在 PaCo 的 decomposed SlotToCoeff 中，矩阵实现会出现 \(U_n\) 及其 Cooley--Tukey 分解；\(U_n\) 本身把 coefficient vector 映到 evaluation/slot vector。但这只是内部矩阵层方向，不改变 bootstrapping 语义上 SlotToCoeff = slot-held coefficients → coefficient polynomial 的命名。
