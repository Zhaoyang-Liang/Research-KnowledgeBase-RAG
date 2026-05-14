问题：Evalmod是怎么做到逐槽的呢，不是只能看到一个密文多项式吗，是把f也加密吗

不是把 (f) 加密。

现实里做到“逐槽”的核心原因是：

> CKKS 的消息空间虽然写成一个多项式环，但这个环在 canonical embedding 下等价于一组槽；环上的加法、乘法在槽坐标里就是逐坐标加法、逐坐标乘法。

也就是说，虽然你手里只有一个密文多项式，但它语义上携带了一个槽向量。

---

设某个密文加密

[
Q(X)\in \mathcal R_d.
]

它的槽向量是

[
\operatorname{can}_d(Q(X))
==========================

(Q(\omega_0),\ldots,Q(\omega_{d-1})).
]

也就是说：

[
\mathbf{ct}\vDash_d^{\mathrm{slot}}(z_0,\ldots,z_{d-1}),
]

其中

[
z_i=Q(\omega_i).
]

现在假设 EvalMod 里用的标量近似函数是一个多项式

[
f(t)=c_0+c_1t+c_2t^2+\cdots+c_mt^m.
]

同态上实际做的是：

[
f(Q(X))
=======

c_0+c_1Q(X)+c_2Q(X)^2+\cdots+c_mQ(X)^m.
]

这是在密文环里对加密的 (Q(X)) 做加法、乘法、乘明文常数。

但关键来了：对每个根 (\omega_i) 求值，有

[
f(Q(X))(\omega_i)
=================

f(Q(\omega_i)).
]

因为

[
Q(X)^2(\omega_i)=Q(\omega_i)^2,
]

[
Q(X)^3(\omega_i)=Q(\omega_i)^3,
]

所以整个多项式函数都会逐槽作用。

因此：

[
\operatorname{can}_d(f(Q(X)))
=============================

(f(z_0),\ldots,f(z_{d-1})).
]

这就是“逐槽”的来源。

---

所以不是说系统真的“看见每个槽”，也不是把每个槽单独拿出来算。

而是因为 canonical embedding 满足：

[
\operatorname{can}_d(A(X)+B(X))
===============================

\operatorname{can}_d(A(X))+\operatorname{can}_d(B(X)),
]

[
\operatorname{can}_d(A(X)B(X))
==============================

\operatorname{can}_d(A(X))\odot \operatorname{can}_d(B(X)),
]

其中 (\odot) 表示逐坐标乘法。

所以你在密文多项式环里做一次乘法，槽语义上等价于所有槽同时乘法。

---

举个最简单的例子。

假设某个密文的槽语义是：

[
\mathbf{ct}\vDash_d^{\mathrm{slot}}(z_0,z_1,z_2,z_3).
]

如果同态计算

[
\mathbf{ct}^2,
]

那么语义上得到：

[
(z_0^2,z_1^2,z_2^2,z_3^2).
]

如果同态计算

[
3\mathbf{ct}^2+2\mathbf{ct}+5,
]

那么语义上得到：

[
(3z_0^2+2z_0+5,\ldots,3z_3^2+2z_3+5).
]

这就是：

[
f(t)=3t^2+2t+5
]

被逐槽应用了。

---

回到 EvalMod。

现实中 EvalMod 的 (f) 通常是某个公开多项式近似，比如近似正弦型模约简函数：

[
f(t)\approx \frac{1}{2\pi}\sin(2\pi t).
]

但密文里不能直接算真正的 (\sin)，所以用多项式

[
p(t)
]

近似它。然后同态评估

[
p(Q(X)).
]

由于环运算在 canonical slots 里是逐槽的，所以自动得到：

[
(z_0,\ldots,z_{d-1})
\mapsto
(p(z_0),\ldots,p(z_{d-1})).
]

---

所以回答你的两个问题：

**1. 不是只能看到一个密文多项式吗？**

是的，算法只操作一个密文多项式。但这个多项式在 canonical embedding 下等价于一个槽向量。环上的加法乘法天然对应槽上的逐坐标加法乘法。

**2. 是把 (f) 也加密吗？**

不是。(f) 是公开的算法/电路。它的多项式系数

[
c_0,c_1,\ldots,c_m
]

是公开常数，通常作为明文常数参与同态计算。真正被加密的是输入 (Q(X))，不是函数 (f)。

---

一句话总结：

> EvalMod 逐槽，不是因为它能解开密文看见每个槽，而是因为 CKKS 的消息环在 canonical embedding 下本来就是 SIMD 结构；对加密环元素评估公开多项式 (f)，槽语义上自动变成对每个槽分别评估同一个 (f)。

===

2. 问： 那这个f怎么直接能跨越到Can空间呢

不是 (f) 自己“跨越到 (\mathsf{Can}_d) 空间”。

更准确地说：

[
f:\mathbb C\to\mathbb C
]

只是一个**标量函数**，只能吃一个复数，吐一个复数。

但 canonical-slot 向量是

[
z=(z_0,\ldots,z_{d-1})\in \mathsf{Can}_d\cong \mathbb C^d.
]

所以定义

[
f^{\oplus d}:\mathsf{Can}_d\to\mathsf{Can}_d
]

为逐坐标作用：

[
f^{\oplus d}(z_0,\ldots,z_{d-1})
================================

(f(z_0),\ldots,f(z_{d-1})).
]

也就是说，真正作用在 (\mathsf{Can}_d) 上的不是 (f)，而是它的逐槽提升：

[
f^{\oplus d}.
]

---

放到密文语义里：

如果

[
\mathbf{ct}\vDash_d^{\mathrm{slot}} z
]

意思是输出密文的槽向量是

[
z=(z_0,\ldots,z_{d-1}).
]

那么 EvalMod 的理想语义是：

[
\mathsf{EvalMod}_d(\mathbf{ct})
\vDash_d^{\mathrm{slot}}
f^{\oplus d}(z).
]

即：

[
(z_0,\ldots,z_{d-1})
\mapsto
(f(z_0),\ldots,f(z_{d-1})).
]

所以它不是在多项式系数上天然作用，而是在**槽坐标**上逐坐标作用。

---

更严格地说，如果你不想使用 slot notation，而是只在多项式空间 (\mathcal R_d) 里写，那么 EvalMod 对多项式的理想作用应该写成：

[
P(X)
\mapsto
\operatorname{can}_d^{-1}
\left(
f^{\oplus d}(\operatorname{can}_d(P(X)))
\right).
]

也就是：

[
\mathcal R_d
\xrightarrow{\operatorname{can}_d}
\mathsf{Can}_d
\xrightarrow{f^{\oplus d}}
\mathsf{Can}_d
\xrightarrow{\operatorname{can}_d^{-1}}
\mathcal R_d.
]

图示：

```text
多项式 P(X)
   │
   │ can_d
   ▼
槽向量 (z_0,...,z_{d-1})
   │
   │ f^{⊕d}
   ▼
槽向量 (f(z_0),...,f(z_{d-1}))
   │
   │ can_d^{-1}
   ▼
新的多项式
```

所以 (f) 没有直接跨到 (\mathcal R_d) 或 (\mathsf{Can}_d)。
是先选定 canonical-slot 坐标，然后把 (f) 逐槽应用。

---

现在回到你关心的 C2S 场景。

原来：

[
\mathbf{ct}\vDash_d P(X),
]

其中

[
P(X)=a_0+a_1X+\cdots+a_{d-1}X^{d-1}.
]

做完 C2S 后：

[
\mathsf{C2S}*d(\mathbf{ct})
\vDash_d^{\mathrm{slot}}
(a_0,\ldots,a*{d-1}).
]

这时 EvalMod 看到的槽向量就是

[
(a_0,\ldots,a_{d-1}).
]

所以它逐槽作用：

[
(a_0,\ldots,a_{d-1})
\mapsto
(f(a_0),\ldots,f(a_{d-1})).
]

因此：

[
\mathsf{EvalMod}_d(\mathsf{C2S}*d(\mathbf{ct}))
\vDash_d^{\mathrm{slot}}
(f(a_0),\ldots,f(a*{d-1})).
]

这里 (f) 没有“直接作用在系数空间”或者“直接作用在多项式空间”。更准确地说：

> C2S 先把原多项式的系数向量放进 canonical slots；然后 EvalMod 在 canonical slots 上逐槽应用 (f)。

---

一句话总结：

[
f
]

只是标量函数；

[
f^{\oplus d}
]

才是它在槽空间上的逐坐标提升；

而 EvalMod 的理想语义是：

[
\operatorname{can}_d^{-1}
\circ
f^{\oplus d}
\circ
\operatorname{can}_d.
]

在 C2S 之后，因为槽里正好装着系数 ((a_0,\ldots,a_{d-1}))，所以 EvalMod 才表现得像是在对系数逐个做 (f)。

===


3. 问：C2Sd​(ct)⊨d​cand−1​(coeffd​(P(X)))。这个怎么理解

这个式子：

[
\mathsf{C2S}_d(\mathbf{ct})
\vDash_d
\operatorname{can}_d^{-1}(\operatorname{coeff}_d(P(X)))
]

是在说：

> 如果原来 (\mathbf{ct}) 加密的是 (P(X))，那么做完 (\mathsf{C2S}_d) 后，新密文加密的**不再是原来的 (P(X))**，而是另一个多项式 (Q(X))。这个 (Q(X)) 的标准槽值正好等于 (P(X)) 的系数向量。

也就是令

[
Q(X)
====

\operatorname{can}_d^{-1}(\operatorname{coeff}_d(P(X))),
]

那么

[
\mathsf{C2S}_d(\mathbf{ct})\vDash_d Q(X).
]

---

关键是要区分两个层次。

原始密文：

[
\mathbf{ct}\vDash_d P(X)
]

表示它加密的是多项式

[
P(X)=a_0+a_1X+\cdots+a_{d-1}X^{d-1}.
]

所以

[
\operatorname{coeff}*d(P(X))=(a_0,a_1,\ldots,a*{d-1}).
]

做 (\mathsf{C2S}_d) 后，我们想让这些系数出现在**槽里**，即：

[
\mathsf{C2S}*d(\mathbf{ct})
\vDash_d^{\mathrm{slot}}
(a_0,a_1,\ldots,a*{d-1}).
]

而槽语义定义是：

[
\mathbf{ct}'\vDash_d^{\mathrm{slot}} z
\quad\Longleftrightarrow\quad
\mathbf{ct}'\vDash_d \operatorname{can}_d^{-1}(z).
]

所以把

[
z=\operatorname{coeff}_d(P(X))
]

代进去，就得到：

[
\mathsf{C2S}_d(\mathbf{ct})
\vDash_d
\operatorname{can}_d^{-1}(\operatorname{coeff}_d(P(X))).
]

这就是你问的式子。它只是把“槽里放的是 (P) 的系数”改写成“密文加密的是某个多项式”的形式。

---

最容易混淆的地方是：

[
\operatorname{coeff}_d(P(X))
]

本来是 **输入多项式 (P(X)) 的系数向量**。

但在

[
\operatorname{can}_d^{-1}(\operatorname{coeff}_d(P(X)))
]

里面，这个向量被当成了**输出多项式的槽值向量**。

所以同一个向量

[
(a_0,\ldots,a_{d-1})
]

有两种身份：

[
\text{作为 }P(X)\text{ 的系数}
]

以及

[
\text{作为新多项式 }Q(X)\text{ 的槽值}.
]

因此

[
Q(X)=\operatorname{can}*d^{-1}(a_0,\ldots,a*{d-1})
]

一般并不等于

[
P(X)=a_0+a_1X+\cdots+a_{d-1}X^{d-1}.
]

---

用一句话概括：

[
\mathsf{C2S}_d
]

不是“把密文解密后取系数”，而是在密文上做一个线性变换，使得输出密文的**canonical slots** 语义上等于原消息多项式的**coefficient vector**。

所以

[
\mathsf{C2S}_d(\mathbf{ct})
\vDash_d
\operatorname{can}_d^{-1}(\operatorname{coeff}_d(P(X)))
]

等价于更直观的说法：

[
\mathsf{C2S}_d(\mathbf{ct})
\vDash_d^{\mathrm{slot}}
\operatorname{coeff}_d(P(X)).
]

后者更好理解：
**做完 C2S 后，槽里装的是原多项式的系数。**

===

4. 问：S2C理解：

这个式子：

[
\mathbf{ct}'\vDash_d^{\mathrm{slot}}y
\quad\Longrightarrow\quad
\mathsf{S2C}_d(\mathbf{ct}')
\vDash_d
\operatorname{coeff}_d^{-1}(y)
]

是在说 **SlotToCoeff 的语义效果**。

直观说：

> 如果当前密文 (\mathbf{ct}') 的槽里装着一个向量 (y)，那么做完 (\mathsf{S2C}_d) 后，输出密文加密的是“把 (y) 当作系数向量”得到的那个多项式。

---

设

[
y=(y_0,\ldots,y_{d-1})\in\mathsf{Coeff}_d.
]

虽然左边写的是

[
\mathbf{ct}'\vDash_d^{\mathrm{slot}}y,
]

意思是：

[
\mathbf{ct}'\text{ 的 canonical slots 是 }y。
]

也就是说，(\mathbf{ct}') 加密的是某个多项式 (Q(X))，并且

[
\operatorname{can}_d(Q(X))=y.
]

现在做 (\mathsf{S2C}_d)。它的作用是把“槽里的值”解释成系数。

所以输出密文加密的是

[
y_0+y_1X+\cdots+y_{d-1}X^{d-1}.
]

这正是

[
\operatorname{coeff}_d^{-1}(y).
]

因此：

[
\mathsf{S2C}_d(\mathbf{ct}')
\vDash_d
\operatorname{coeff}_d^{-1}(y).
]

---

它和前面的 C2S 正好相反。

C2S 是：

[
P(X)
\longmapsto
\text{把 }P(X)\text{ 的系数放到槽里}.
]

也就是：

[
\mathbf{ct}\vDash_d P(X)
\Rightarrow
\mathsf{C2S}_d(\mathbf{ct})
\vDash_d^{\mathrm{slot}}
\operatorname{coeff}_d(P(X)).
]

S2C 是：

[
\text{槽里有 }y
\longmapsto
\text{把 }y\text{ 当作系数还原成多项式}.
]

也就是：

[
\mathbf{ct}'\vDash_d^{\mathrm{slot}}y
\Rightarrow
\mathsf{S2C}_d(\mathbf{ct}')
\vDash_d
\operatorname{coeff}_d^{-1}(y).
]

所以如果

[
y=\operatorname{coeff}_d(P(X)),
]

那么

[
\operatorname{coeff}_d^{-1}(y)=P(X),
]

于是有：

[
\mathbf{ct}'\vDash_d^{\mathrm{slot}}\operatorname{coeff}_d(P(X))
\Rightarrow
\mathsf{S2C}_d(\mathbf{ct}')\vDash_d P(X).
]

这就是为什么说 (\mathsf{S2C}_d) 可以把 coefficient-packed slots 变回原来的消息多项式。

一句话总结：

[
\mathsf{S2C}_d
]

就是把**槽坐标里的向量**重新解释成**系数坐标里的向量**，然后得到对应的多项式。

5. 出现了一个未定义的a：
这里的

[
(a_0,\ldots,a_{d-1})
]

就是原始消息多项式

[
P(X)=a_0+a_1X+\cdots+a_{d-1}X^{d-1}
]

的**系数向量**。

也就是说：

[
a_i
]

是 (X^i) 前面的系数。

例如：

[
P(X)=3+5X-2X^2+7X^3
]

那么

[
a_0=3,\quad a_1=5,\quad a_2=-2,\quad a_3=7.
]

所以

[
\operatorname{coeff}*d(P(X))=(a_0,\ldots,a*{d-1}).
]

---

在你这句里：

[
\mathsf{C2S}*d(\mathbf{ct})\vDash_d^{\mathrm{slot}}(a_0,\ldots,a*{d-1})
]

意思是：

> 做完 (\mathsf{C2S}_d) 以后，输出密文的槽里装的是原多项式 (P(X)) 的系数。

然后 EvalMod 是逐槽作用的，所以：

[
\mathsf{EvalMod}_d(\mathsf{C2S}*d(\mathbf{ct}))
\vDash_d^{\mathrm{slot}}
(f(a_0),\ldots,f(a*{d-1})).
]

意思就是：

> 每个槽里的系数 (a_i) 被单独送进同一个函数 (f)。

所以 (a_i) 不是新的变量，也不是槽本身的名字，而是原始多项式 (P(X)) 的第 (i) 个系数。

一句话：

[
a_i = P(X)\text{ 在 }X^i\text{ 项上的系数。}
]

