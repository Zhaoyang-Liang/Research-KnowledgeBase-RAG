# Hermes Split and Key Switching

你问得很关键。HERMES 的“拆”不是 slot 拆分，也不是简单把 ciphertext 数组切半；它本质上是：

[
\boxed{
\text{把大环 }R_{q,N}\text{ 看成小环 }R_{q,n}\text{ 上的 }k\text{ 维模，然后按 }X\text{ 的余数类拆系数。}
}
]

这里

[
N=kn,\qquad R_{q,N}=\mathbb Z_q[X]/(X^N+1),
]

[
R_{q,n}=\mathbb Z_q[Y]/(Y^n+1),
\qquad Y=X^k.
]

HERMES Section 3.2 正是这么定义 module decomposition map (\pi^k_{q,N}) 和 ring embedding (\iota^k_{q,N}) 的；它说 (R_{q,N}) 可以作为 (R_{q,N/k})-module，basis 是 (1,X,\ldots,X^{k-1})，并把一个大环元素拆成 (k) 个小环元素。

---

# 1. 先看大环怎么按小环拆

设

[
N=kn.
]

任意大环多项式

[
P(X)\in R_{q,N}
]

都可以唯一写成：

[
P(X)
====

P_0(X^k)+XP_1(X^k)+\cdots+X^{k-1}P_{k-1}(X^k).
]

令

[
Y=X^k,
]

就等价于：

[
P(X)
====

P_0(Y)+XP_1(Y)+\cdots+X^{k-1}P_{k-1}(Y),
]

其中每个

[
P_j(Y)\in R_{q,n}.
]

这个就是 coefficient split。

例如 (k=2) 时：

[
P(X)=P_0(X^2)+X P_1(X^2).
]

如果

[
P(X)=p_0+p_1X+p_2X^2+p_3X^3+\cdots,
]

那么

[
P_0(Y)=p_0+p_2Y+p_4Y^2+\cdots,
]

[
P_1(Y)=p_1+p_3Y+p_5Y^2+\cdots.
]

也就是：

[
\boxed{
\text{偶数次系数进 }P_0,\quad
\text{奇数次系数进 }P_1.
}
]

更一般地，(k) 个 leaf 对应 (X) 次数模 (k) 的余数类：

[
P_r(Y)=\sum_j p_{r+jk}Y^j.
]

这就是你文里的：

[
P(X)=\sum_{r=0}^{k-1}X^rP_r(X^k).
]

---

# 2. HERMES 的 Split 是对 ciphertext 的两个分量一起做这个拆分

CKKS/RLWE ciphertext 通常写作：

[
\ct=(b,a)\in R_{q,N}^2.
]

HERMES 采用的符号约定大概是：

[
b+a s=m.
]

有时也写成：

[
b=-as+m.
]

这两个是同一回事。

现在把 (a,b,m) 都按上面的 module decomposition 拆开：

[
a=\sum_{j=0}^{k-1}a_j\beta_j,
\qquad
b=\sum_{j=0}^{k-1}b_j\beta_j,
\qquad
m=\sum_{j=0}^{k-1}m_j\beta_j,
]

其中

[
\beta_j=X^j,
]

并且

[
a_j,b_j,m_j\in R_{q,n}.
]

所以大 ciphertext 的两个多项式分量被拆成：

[
a\mapsto(a_0,\ldots,a_{k-1}),
]

[
b\mapsto(b_0,\ldots,b_{k-1}).
]

HERMES 里说的 Split，本质上就是：

[
(b,a)
\mapsto
\bigl((b_0,a_0),\ldots,(b_{k-1},a_{k-1})\bigr).
]

但注意：**这一步只有在 secret 已经属于小环子环时，才会得到合法的小环 ciphertext。**

---

# 3. 为什么 secret 必须在子环里？

这是最关键的一点。

假设大环 ciphertext 满足：

[
b=-as+m.
]

如果 secret 是一个小环 secret (s_0(Y))，嵌入大环后是：

[
s=\iota(s_0)=s_0(X^k).
]

也就是说：

[
s=s_0\cdot 1.
]

这里的“(s=s_0\cdot 1)”意思是：在大环作为小环 module 的分解里，secret 只有第 0 个方向：

[
s=s_0+0\cdot X+0\cdot X^2+\cdots+0\cdot X^{k-1}.
]

所以我前面说：

[
s=s_0\cdot 1,
]

没有 (X,X^2,\ldots,X^{k-1}) 方向的分量。

不是说 (s) 是常数，而是说：

[
s=s_0(X^k)
]

只含有 (X^{0},X^k,X^{2k},\ldots) 这些次数。

---

# 4. 为什么这样就能拆成小环 ciphertext？

因为

[
a=\sum_{j=0}^{k-1}a_j\beta_j,
\qquad
s=s_0,
\qquad
m=\sum_{j=0}^{k-1}m_j\beta_j.
]

于是

[
as
==

# \left(\sum_j a_j\beta_j\right)s_0

\sum_j (a_js_0)\beta_j.
]

注意这里没有混合不同 (j)，因为 (s_0\in R_{q,n})，它只是在每个 leaf 里面乘。

所以：

[
b=-as+m
]

变成：

[
\sum_j b_j\beta_j
=================

-\sum_j a_js_0\beta_j
+
\sum_j m_j\beta_j.
]

由于 module decomposition 唯一，所以逐个 (j) 比较系数：

[
b_j=-a_js_0+m_j.
]

等价于：

[
b_j+a_js_0=m_j.
]

这说明：

[
(b_j,a_j)
]

就是一个小环 (R_{q,n}) 上的 RLWE ciphertext，secret 是 (s_0)，plaintext 是 (m_j)。

所以：

[
\boxed{
(b,a)\text{ 在大环下加密 }m
\quad\Rightarrow\quad
(b_j,a_j)\text{ 在小环下加密 }m_j.
}
]

这就是 HERMES Split 的核心。

HERMES 原文也说：如果 ciphertext 的 secret key belongs to a subring，那么 ciphertext 可以 split 成 (k) pieces，每一块 encrypts 一部分 plaintext；并把 conversion ((b,a)\mapsto(b_j,a_j)_j) 称为 Split，反方向称为 Combine。

---

# 5. 用 (N=8,k=2,n=4) 举一个非常具体的例子

令

[
R_{q,8}=\mathbb Z_q[X]/(X^8+1),
]

[
R_{q,4}=\mathbb Z_q[Y]/(Y^4+1),
\qquad Y=X^2.
]

任意大环元素都能写成：

[
a(X)=a_0(Y)+Xa_1(Y),
]

[
b(X)=b_0(Y)+Xb_1(Y),
]

[
m(X)=m_0(Y)+Xm_1(Y).
]

假设 secret 是小环 secret 嵌入：

[
s(X)=s_0(Y)=s_0(X^2).
]

大环 ciphertext 满足：

[
b(X)=-a(X)s(X)+m(X).
]

代入：

[
b_0(Y)+Xb_1(Y)
==============

-\bigl(a_0(Y)+Xa_1(Y)\bigr)s_0(Y)
+
m_0(Y)+Xm_1(Y).
]

展开：

[
b_0(Y)+Xb_1(Y)
==============

\bigl(-a_0(Y)s_0(Y)+m_0(Y)\bigr)
+
X\bigl(-a_1(Y)s_0(Y)+m_1(Y)\bigr).
]

因此：

[
b_0(Y)=-a_0(Y)s_0(Y)+m_0(Y),
]

[
b_1(Y)=-a_1(Y)s_0(Y)+m_1(Y).
]

所以：

[
(b_0,a_0)
]

是小环 ciphertext，加密 (m_0)；

[
(b_1,a_1)
]

是小环 ciphertext，加密 (m_1)。

这就是 split。

---

# 6. 如果 secret 不是子环 secret，会发生什么？

这是最容易误解的地方。

假设 (k=2)，但 secret 不是小环嵌入，而是一般大环 secret：

[
s(X)=s_0(Y)+X s_1(Y).
]

同时：

[
a(X)=a_0(Y)+Xa_1(Y).
]

那么：

[
a(X)s(X)
========

(a_0+Xa_1)(s_0+Xs_1).
]

展开：

[
a(X)s(X)
========

a_0s_0
+
X(a_0s_1+a_1s_0)
+
X^2a_1s_1.
]

因为

[
X^2=Y,
]

所以：

[
X^2a_1s_1
=========

Y a_1s_1
]

又回到第 0 个方向。

于是：

[
as
==

\underbrace{(a_0s_0+Ya_1s_1)}*{\text{第 0 个方向}}
+
X\underbrace{(a_0s_1+a_1s_0)}*{\text{第 1 个方向}}.
]

这时候第 0 个方向不只是 (a_0s_0)，还混进了：

[
Ya_1s_1.
]

第 1 个方向也不只是 (a_1s_0)，还混进了：

[
a_0s_1.
]

所以如果你直接 split：

[
(b,a)\mapsto (b_0,a_0),(b_1,a_1),
]

那么一般不会有：

[
b_0=-a_0s_0+m_0,
]

[
b_1=-a_1s_0+m_1.
]

换句话说：

[
(b_0,a_0)
]

和

[
(b_1,a_1)
]

不是合法的小环 ciphertext。

所以必须先保证：

[
s=s_0(X^k)
]

只在子环里。

---

# 7. 那原始 ciphertext 的 secret 不在子环怎么办？

实际情况通常是：top ciphertext 一开始在普通大环 secret (s_N) 下：

[
\ct_N=(b,a),\qquad b+a s_N=m.
]

这个 (s_N) 一般不是子环 secret。

所以 HERMES / ring switching 的降环方向是：

[
\boxed{
\text{先 key switch 到嵌入的子环 secret，再 Split。}
}
]

也就是：

[
\ct_N(s_N)
]

先通过 switching key 变成：

[
\ct_N(\iota(s_n)),
]

其中

[
s_n\in R_{q,n}
]

是小环 secret，

[
\iota(s_n)=s_n(X^k)\in R_{q,N}.
]

然后再做 coefficient Split：

[
\ct_N(\iota(s_n))
\mapsto
(\ct_{n,0}(s_n),\ldots,\ct_{n,k-1}(s_n)).
]

这一步才会得到合法 leaf ciphertexts。

HERMES 原文也明确说，ring switching 到 lower degree 包括两步：先把 large-degree ciphertext key switch 到 small-degree key，这个 key belongs to the subring，然后 applying Split。

---

# 8. Combine 是 Split 的反方向

如果你已经有 (k) 个小环 ciphertext：

[
(b_j,a_j)\in R_{q,n}^2,
]

都在同一个小环 secret (s_0) 下：

[
b_j+a_js_0=m_j,
]

那么可以 Combine：

[
b=\sum_j b_jX^j,
\qquad
a=\sum_j a_jX^j,
\qquad
m=\sum_j m_jX^j.
]

因为

[
b+as
====

\sum_j b_jX^j
+
\left(\sum_j a_jX^j\right)s_0
]

# [

# \sum_j(b_j+a_js_0)X^j

# \sum_jm_jX^j

m.
]

所以 Combine 得到一个大环 ciphertext，但它的 secret 是嵌入的：

[
s=s_0(X^k).
]

如果你最后想回到普通 top secret (s_N)，就再做一次 key switch：

[
\ct_N(\iota(s_n))
\to
\ct_N(s_N).
]

所以升环方向一般是：

[
(\ct_{n,0},\ldots,\ct_{n,k-1})
\xrightarrow{\mathrm{Combine}}
\ct_N(\iota(s_n))
\xrightarrow{\mathrm{KS}}
\ct_N(s_N).
]

HERMES 也说，ring switching to higher degree is made of Combine to go from lower-degree RLWE ciphertexts to a higher-degree RLWE ciphertext with the same key, and then key switch to obtain a higher-degree key。

---

# 9. HERMES 里为什么说 coefficients-encoded CKKS 里 Split/Combine 很便宜？

因为 Split/Combine 只是 rearrangement of coefficients。

例如 (k=2)：

[
a(X)=a_0(X^2)+Xa_1(X^2).
]

Split 就是把 (a) 的偶数系数放到 (a_0)，奇数系数放到 (a_1)。

对 (b) 同理。

这不需要复杂 homomorphic computation，也不需要乘 plaintext matrix。

HERMES 原文特别说：在 coefficients-encoded CKKS 里，((b_j,a_j)_j) 和 ((b,a)) 的 bijection 是 direct sum，只是 coefficient rearrangement，因此 Split 和 Combine 几乎没有计算成本；这也是他们 revisiting ring switching 的主要改进点。

这和你论文里的观点完全吻合：

[
\text{Hermes Split 是 coefficient split，不是 slot split。}
]

---

# 10. 和你当前 no-(B_k) bootstrapping 的关系

你的流程大概是：

[
\ct_N(s_N)
\to
\text{RSDown/Split}
\to
(\ct_{n,0},\ldots,\ct_{n,k-1})
\to
\text{leaf BTS}
\to
\text{Merge/RSUp}.
]

按照 HERMES 语义，RSDown/Split 实际上应该理解为：

[
\ct_N(s_N)
\xrightarrow{\mathrm{KS}*{s_N\to \iota(s_n)}}
\ct_N(\iota(s_n))
\xrightarrow{\mathrm{Split}}
(\ct*{n,0}(s_n),\ldots,\ct_{n,k-1}(s_n)).
]

Merge/RSUp 则是：

[
(\ct_{n,0}(s_n),\ldots,\ct_{n,k-1}(s_n))
\xrightarrow{\mathrm{Combine}}
\ct_N(\iota(s_n))
\xrightarrow{\mathrm{KS}_{\iota(s_n)\to s_N}}
\ct_N(s_N).
]

你代码里用 Lattigo 的 RingPackingEvaluator，底层大概率就是在做这种 key-switch + Split/Combine 语义。

---

# 11. 这和你前面证明的 coefficient split theorem 怎么接？

你的 theorem 主要在 message level 上说：

[
P(X)
====

\sum_{r=0}^{k-1}X^rP_r(X^k),
]

然后 leaf core 分别作用在 (P_r) 上，最后 Merge 回来。

HERMES 解释的是 ciphertext level 如何实现：

[
\ct_N(P)
\to
(\ct_n(P_0),\ldots,\ct_n(P_{k-1})).
]

二者对应关系是：

[
\text{HERMES ciphertext Split}
\quad\text{implements}\quad
\text{你的 message coefficient split}.
]

但是前提是：

[
\text{ciphertext 已经被 key switched 到 embedded subring secret。}
]

所以你论文中最好明确说：

```tex
At the ciphertext level, \(\mathsf{Split}_k\) denotes the ideal message
semantics induced by the HERMES-style RSDown operation: first switch the
large-ring ciphertext to an embedded subring secret, and then apply the
coefficient/module decomposition.
```

这样就非常严谨。

---

# 12. 一个总图

你可以把 HERMES-style split/merge 记成：

```text
Down direction:

ct_N under s_N
    |
    | key switch to embedded subring secret iota(s_n)
    v
ct_N under iota(s_n)
    |
    | coefficient/module Split
    v
(ct_{n,0}, ..., ct_{n,k-1}) under s_n
```

反方向：

```text
Up direction:

(ct_{n,0}, ..., ct_{n,k-1}) under s_n
    |
    | coefficient/module Combine
    v
ct_N under iota(s_n)
    |
    | key switch to top secret s_N
    v
ct_N under s_N
```

这就是你现在实验里所谓的 RSDown / RSUp 的理想解释。

---

# 13. 最后一句直观总结

HERMES 的拆分不是说：

[
\text{大 ciphertext 的 slots 切成几段。}
]

它说的是：

[
\text{先把 secret 切到一个嵌入小环的 secret，使大环 ciphertext 在代数上可视作 }k\text{ 个小环 ciphertext 的 direct sum。}
]

然后 coefficient decomposition：

[
a(X)=\sum_jX^ja_j(X^k),
\qquad
b(X)=\sum_jX^jb_j(X^k)
]

自然给出：

[
(b_j,a_j)
]

这些小环 ciphertext。

这就是 HERMES Split 的核心。

===

可以。HERMES 里面有两层 key switching：

1. **普通 RLWE key switching**：把同一个 ring 上的 ciphertext 从 secret (s) 切到 (s')；
2. **ring switching / MLWE key switching**：借助普通 RLWE key switching，在不同 ring degree 或不同 module rank 之间转换。

你现在最需要理解的是第 1 层，因为 HERMES 的 Split/Combine、ModPack、MLWE key switching 都是围绕它搭起来的。

---

# 1. HERMES 的普通 RLWE key switching

HERMES Section 2.1 用的是 RLWE key switching 的标准形式。设原 ciphertext 是

[
\ct=(b,a)\in R_{q,N}^2,
]

在 secret (s) 下解密：

[
b+a s=m.
]

目标是得到一个新 ciphertext

[
\ct'=(b',a')
]

使得它在新 secret (s') 下解密仍然是同一个 (m)：

[
b'+a's'\approx m.
]

HERMES 把 key switching 分成三个子步骤：`ModUp`、`MultSwk`、`ModDown`。它先把 (a\in R_{q,N}) 提升到 (R_{qp,N})，然后乘 switching key，最后 approximate divide by (p) 回到 (R_{q,N})。论文里明确写了这三个步骤，以及 dnum=1 情况下的 switching key 形式。

---

## 1.1 switching key 是什么

HERMES 定义从 (s) 到 (s') 的 switching key 为：

[
\mathsf{swk}_{s\to s'}
======================

(-\hat u s' + p s+\hat e,\ \hat u)
\in R_{qp,N}^2,
]

其中：

* (\hat u\leftarrow R_{qp,N}) 是均匀随机；
* (\hat e) 是小误差；
* (p) 是 auxiliary modulus；
* (qp) 是扩展模数。

你可以把它理解成：

[
\mathsf{swk}*{s\to s'}
\approx
\operatorname{Enc}*{s'}(p s).
]

也就是：**用目标 secret (s') 加密源 secret (s) 的 (p)-倍。**

为什么要乘 (p)？因为后面会 `ModDown`，近似除以 (p)。乘 (p) 是为了让 (s) 在除以 (p) 后恢复出来。

---

## 1.2 Key switching 的计算

给定

[
\ct=(b,a),
\qquad b+as=m,
]

先做：

[
\hat a=\operatorname{ModUp}(a)\in R_{qp,N}.
]

然后计算：

[
\hat a\cdot \mathsf{swk}_{s\to s'}
==================================

\left(
\hat a(-\hat u s'+ps+\hat e),
\hat a\hat u
\right).
]

再做 `ModDown`，近似除以 (p)，最后加上 ((b,0))：

[
\operatorname{KS}_{s\to s'}(\ct)
================================

\operatorname{ModDown}\left(
\hat a\cdot \mathsf{swk}_{s\to s'}
\right)
+
(b,0).
]

HERMES 公式就是：

[
\operatorname{KS}(\ct;\mathsf{swk})
===================================

\operatorname{ModDown}
\bigl(
\operatorname{MultSwk}(\operatorname{ModUp}(a),\mathsf{swk})
\bigr)
+
(b,0).
]



---

## 1.3 为什么它正确

忽略 rounding / ModDown 误差，`ModDown` 可以看成除以 (p)。于是输出近似为：

[
b'
==

b+\frac{\hat a(-\hat u s'+ps+\hat e)}{p},
]

[
a'
==

\frac{\hat a\hat u}{p}.
]

用新 secret (s') 解密：

[
b'+a's'
=======

b+
\frac{\hat a(-\hat u s'+ps+\hat e)}{p}
+
\frac{\hat a\hat u}{p}s'.
]

中间两项抵消：

[
\frac{-\hat a\hat u s'}{p}
+
\frac{\hat a\hat u s'}{p}
=0.
]

剩下：

[
b'+a's'
=======

b+\hat a s+\frac{\hat a\hat e}{p}
+\text{rounding error}.
]

因为 (\hat a) 是 (a) 的 ModUp，所以近似是：

[
b'+a's'
=======

b+a s
+
\text{key-switching noise}.
]

而原来：

[
b+as=m.
]

所以：

[
b'+a's'
=======

m+e_{\mathrm{KS}}.
]

这就是 key switching 的正确性。

一句话：

[
\boxed{
\text{switching key 里的 }-\hat u s'\text{ 和新密文里的 }+\hat u s'\text{ 抵消，留下 }ps。
}
]

除以 (p) 后就恢复出原来的 (as)。

---

# 2. HERMES 里的 ring switching 怎么用这个 key switching

HERMES Section 3.2 先把大环

[
R_{q,N}
]

看成小环

[
R_{q,n},\qquad n=N/k
]

上的 rank-(k) module。它定义了 module decomposition：

[
\pi^k_{q,N}:R_{q,N}\to (R_{q,n})^k,
]

以及 embedding：

[
\iota^k_{q,N}:R_{q,n}\to R_{q,N}.
]

也就是：

[
Y\mapsto X^k.
]

任意大环元素可以唯一写成：

[
a(X)=\sum_{j=0}^{k-1}X^j a_j(X^k).
]

HERMES 的 Split/Combine 正是基于这个分解。

---

## 2.1 降环：从大环切到小环

假设原始 ciphertext 在普通大环 secret (s_N) 下：

[
\ct_N=(b,a),
\qquad
b+a s_N=m.
]

为了 Split 成 (k) 个小环 ciphertext，不能直接拆。必须先把 secret key switch 到嵌入的小环 secret：

[
\iota(s_n)=s_n(X^k)\in R_{q,N}.
]

流程是：

[
\ct_N(s_N)
\xrightarrow{\operatorname{KS}*{s_N\to \iota(s_n)}}
\ct_N(\iota(s_n))
\xrightarrow{\operatorname{Split}}
(\ct*{n,0}(s_n),\ldots,\ct_{n,k-1}(s_n)).
]

HERMES 原文说：ring switching to a lower degree consists in switching the key of a large-degree RLWE ciphertext to a small-degree key, which belongs to the subring, and then applying Split。

这句话对你非常关键：
**Split 本身只是 coefficient rearrangement；真正让它变成合法 leaf ciphertext 的，是 Split 之前的 key switch 到 subring secret。**

---

## 2.2 为什么 Split 后每个 leaf 是合法 ciphertext

假设 key switching 后：

[
b=-a\cdot \iota(s_n)+m.
]

把

[
a=\sum_j X^j a_j(X^k),\quad
b=\sum_j X^j b_j(X^k),\quad
m=\sum_j X^j m_j(X^k)
]

代进去。因为

[
\iota(s_n)=s_n(X^k)
]

只在小环方向上，不会混合 (j) 的余数类，所以：

[
b_j=-a_js_n+m_j.
]

于是每个

[
(b_j,a_j)\in R_{q,n}^2
]

都是小环 ciphertext，解密到 (m_j)。

这就是 HERMES 的 coefficient Split。

---

## 2.3 升环：从小环合回大环

如果有 (k) 个小环 ciphertext：

[
(b_j,a_j),\qquad b_j+a_js_n=m_j,
]

先 Combine：

[
b=\sum_jX^j b_j(X^k),
\qquad
a=\sum_jX^j a_j(X^k).
]

这会得到一个大环 ciphertext，但 secret 是嵌入的：

[
\iota(s_n).
]

也就是：

[
b+a\iota(s_n)=m.
]

如果最终要回到普通 top secret (s_N)，再做一次 key switching：

[
\ct_N(\iota(s_n))
\xrightarrow{\operatorname{KS}_{\iota(s_n)\to s_N}}
\ct_N(s_N).
]

HERMES 对升环方向的描述是：先 Combine lower-degree RLWE ciphertexts 到 higher-degree RLWE ciphertext with same key，再 key switch 到 higher-degree key。

---

# 3. HERMES 的 MLWE key switching 是什么

HERMES 还把这个机制推广到 MLWE。这个主要用于它的 ModPack / BaseHERMES，不一定是你 no-(B_k) leaf bootstrapping 证明的核心，但理解它能帮助你理解 HERMES 的整体设计。

MLWE ciphertext 是：

[
c=(b,a_1,\ldots,a_k)\in (R_{q,n})^{k+1}
]

在 secret

[
s=(s_1,\ldots,s_k)
]

下满足：

[
b+\langle a,s\rangle=m.
]

HERMES 的 MLWE key switching 做三步：

1. 把 MLWE ciphertext **嵌入**成一个更高 degree 的 RLWE ciphertext；
2. 对这个 RLWE ciphertext 做普通 RLWE key switching；
3. 再从结果中 **抽取** 出 MLWE ciphertext。

论文 Section 4.2 明确写了这个三步流程：embed input MLWE ciphertext as part of an RLWE ciphertext, switch the key of the RLWE ciphertext, then extract an MLWE ciphertext that encrypts valid data。

它的核心技巧是一个 twist map。twist 的作用是让大环乘法中的某个 coefficient extraction 等于小环 MLWE 的 inner product：

[
\epsilon^k_{q,N}\left(a^{\mathrm{tw},k}\cdot s\right)
=====================================================

\langle \pi^k_{q,N}(a),\pi^k_{q,N}(s)\rangle.
]

直观地说：

[
\text{MLWE 内积}
\quad\leftrightarrow\quad
\text{嵌入到大环后的某个系数/分量}.
]

然后 HERMES Theorem 1 说：

[
\operatorname{Extract}
\circ
\operatorname{KS}_{S\to S'}
\circ
\operatorname{Embed}(c)
=======================

\operatorname{MLWE.Enc}_{s'}(m).
]

也就是：嵌入成 RLWE、做普通 key switch、再抽取，整体等价于 MLWE key switching。

---

# 4. HERMES 如何分析安全性

这里要说清楚：**HERMES 没有给一个像“完整 UC 安全证明”那样的长安全定理。** 它的安全处理更像工程 FHE 论文的参数化安全说明：

1. 依赖 LWE / RLWE / MLWE 标准假设；
2. 所有 switching keys 都是目标 secret 下的 RLWE encryptions；
3. 降低 ring degree 只允许在 lowered degree 与当前 modulus 下仍满足安全；
4. 参数表中给出能维持目标安全的最大 modulus / degree。

---

## 4.1 key switching key 的安全性

switching key 是：

[
\mathsf{swk}_{s\to s'}
======================

(-\hat u s'+ps+\hat e,\hat u).
]

这就是 target secret (s') 下对 plaintext (ps) 的 RLWE encryption。

所以它的安全依赖于：

[
(R_{q p,N},\ s',\ \hat e)
]

这组 RLWE 参数。

如果 (s') 是独立 target secret，那么可以把 switching key 看成普通 RLWE encryptions of some plaintext。只要目标参数安全，公开它不会泄露 (s)。

如果 source secret 和 target secret 有循环依赖，比如某些 relinearization key 加密 (s^2) 或同一个 secret 的函数，那就需要标准 HE 里的 key-switching/KDM/circular-style 假设。HERMES 本文没有展开这个问题，而是按常规 FHE key switching 处理。

对你论文来说，建议写得比 HERMES 更明确：

[
\boxed{
\text{RSDown/RSUp keys are RLWE encryptions under the target secret; their security is estimated at the target ring/modulus.}
}
]

---

## 4.2 ring degree 可以降低，但不能免费降低

HERMES 讲得很直白：

> We can lower the degree as long as RLWE remains secure under the current modulus. In particular, the only restriction is the security for the switching keys with respect to their moduli and lowered degrees.

也就是说，HERMES 并不是说：

[
N\text{ 安全}\Rightarrow n\text{ 自动安全}.
]

它说的是：

[
\text{可以降到 }n,\quad
\text{但前提是 }(n,q,\chi_s,\chi_e,m)\text{ 仍安全。}
]



这正是我们之前一直说的：

[
\text{leaf security 按 leaf degree }n\text{ 算。}
]

---

## 4.3 为什么 HERMES 能用小 degree

HERMES 的 base ring packing 是在很低的 modulus (Q_{\mathrm{Enc}}) 下做的，而不是在 bootstrapping 的大 modulus (Q_{\mathrm{top}}) 或 (QP) 下做。

它的流程是：

1. 低 modulus 下做 base ring packing；
2. 用 ring switching 把多个小 degree RLWE 合成大 degree RLWE；
3. 再 HalfBTS 提到高 modulus。

论文 Section 3.1–3.2 就是这个策略：先在最低 modulus (Q_{\mathrm{Enc}}) 做 coefficient-encoded ring packing，再 HalfBTS；这样 base ring packing 可以用更小 ring degree，同时仍然保持安全。

这点和你的 same-(Q) leaf bootstrapping 不完全一样。

你的当前实验是：

[
\text{leaf ring degree 变小，但 leaf bootstrapping 仍在 large }QP\text{ 下。}
]

所以你不能直接套 HERMES 的“小 degree 安全”直觉。HERMES 小 degree 能成立，关键原因之一是它在低 modulus 下做 base packing。

---

## 4.4 HERMES 的参数表怎么体现安全

HERMES Table 2 给了一个 HEaaN FGb bootstrapping 参数，例如：

[
N=2^{16},\qquad \log_2(QP)=1555.
]

表中说明 (\log_2(QP)) 是在保持 desired security 的情况下可使用的最大 ciphertext modulus。

Table 7 另给了一个比较 Pegasus 的参数：

[
N=2^{15},\qquad (h,\tilde h)=(256,32),\qquad \log_2(QP)=820.
]

表注说明这里 (h,\tilde h) 是 dense/sparse secret 的 Hamming weights，而 (\log_2(QP)) 也是在维持 desired security 下的最大可用 modulus。

这说明他们的安全处理方式是：

[
\text{根据 ring degree、modulus、secret distribution 选参数，使其约 }128\text{-bit secure}.
]

Table 11 也直接说其中比较的 transciphering 参数集达到约 128-bit security。

但他们没有在正文里展开 lattice-estimator 的全部攻击成本表。

---

# 5. 这对你的论文意味着什么

你现在用 HERMES-style RingPacking split/merge，可以这样理解：

## RSDown

[
\ct_N(s_N)
\to
\ct_N(\iota(s_n))
\to
(\ct_{n,0}(s_n),\ldots,\ct_{n,k-1}(s_n)).
]

第一箭头是 key switching：

[
\operatorname{KS}_{s_N\to \iota(s_n)}.
]

第二箭头是 coefficient Split。

## RSUp

[
(\ct_{n,0}(s_n),\ldots,\ct_{n,k-1}(s_n))
\to
\ct_N(\iota(s_n))
\to
\ct_N(s_N).
]

第一箭头是 Combine。

第二箭头是 key switching：

[
\operatorname{KS}_{\iota(s_n)\to s_N}.
]

所以你的安全表里至少要估：

| 对象                        | 怎么估                                                                           |
| ------------------------- | ----------------------------------------------------------------------------- |
| top ciphertext / top keys | degree (N), top modulus                                                       |
| RSDown key                | target 是 embedded subring secret；不能白嫖 (N)，要按 subring/leaf security accounting |
| leaf ciphertexts          | degree (n), leaf modulus                                                      |
| leaf bootstrapping keys   | degree (n), leaf (q_{\max})                                                   |
| RSUp key                  | 通常 target 回到 top secret，按 top-side switching key modulus 估                    |
| output ciphertext         | degree (N)                                                                    |

更准确地说，RSDown switching key 形式上活在大环 (R_{qp,N})，但目标 secret 是嵌入的小环 secret (\iota(s_n))。所以它不是普通 dense (N)-degree secret 的 RLWE sample。论文安全章节中应该保守写成：**RSDown 安全至少要按 target subring secret 的有效维度 (n)、实际 modulus (qp)、secret distribution 和 sample count 估计**。这也符合 HERMES 自己说的“lowered degree and current modulus must remain secure”。

---

# 6. 你可以直接放进论文的表述

```tex
At the ciphertext level, our \(\mathsf{RSDown}\) operation follows the
HERMES-style ring-switching view. Let \(N=kn\). We first key-switch the
large-ring ciphertext from the top secret \(s_N\) to an embedded subring
secret \(\iota(s_n)\in R_{q,N}\). Once the ciphertext is under
\(\iota(s_n)\), the coefficient/module decomposition
\[
R_{q,N}\simeq \bigoplus_{r=0}^{k-1} X^r R_{q,n}
\]
turns it into \(k\) valid \(R_{q,n}\)-ciphertexts under \(s_n\).
The reverse operation first combines the leaf ciphertexts into a
large-ring ciphertext under \(\iota(s_n)\), and then key-switches it back
to the top secret \(s_N\).
```

安全段可以写：

```tex
The HERMES-style ring-switching backend does not make security inherit
from the top ring. Each exposed object must be estimated at its own
effective ring degree and maximum public modulus. In particular, the
leaf ciphertexts and leaf bootstrapping keys are estimated at degree
\(n\), while switching keys are estimated according to their target
secret, public modulus, and sample count.
```

---

# 7. 一句话总结

HERMES 的 key switching 是：

[
\mathsf{swk}_{s\to s'}
======================

(-\hat u s'+ps+\hat e,\hat u),
]

它让

[
(b,a)\text{ under }s
]

变成

[
(b',a')\text{ under }s'
]

且 plaintext 不变，只增加 key-switching noise。

HERMES 的 ring switching 是：

[
\text{key switch 到嵌入子环 secret}
+
\text{coefficient Split}
]

或者反向：

[
\text{Combine}
+
\text{key switch 回目标大环 secret}.
]

HERMES 的安全分析不是“top 安全自动推出 leaf 安全”，而是：**只要 lowered degree 在当前 modulus 下仍然满足 RLWE/key-switching-key 安全，就可以降 degree**。这正是你论文安全章节需要继承并加强的地方。
