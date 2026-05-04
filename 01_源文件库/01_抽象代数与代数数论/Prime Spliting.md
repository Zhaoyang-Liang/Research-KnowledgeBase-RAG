可以。我们要证明的是这句话：

[
R/\mathfrak p_i \cong \mathbb F_{p^d},
]

其中

[
d=\operatorname{ord}_{\bar m}(p)
]

是 (p) 模 (\bar m) 的乘法阶。也就是说，(d) 是最小正整数，使得

[
p^d\equiv 1 \pmod{\bar m}.
]

论文在 prime splitting 部分说，每个 (\mathfrak p_i) 的 quotient ring 都同构于 (\mathbb F_{p^d})。我们现在证明为什么这个 (d) 正好是有限域扩张次数。

---

## 1. 先说直觉

在有限域里，(\mathbb F_p) 本身不一定含有 (\bar m) 次单位根。

所以我们要找一个扩域：

[
\mathbb F_{p^d}
]

让里面有 primitive (\bar m)-th root of unity。

有限域 (\mathbb F_{p^d}) 的非零元素构成乘法群：

[
\mathbb F_{p^d}^*
]

它的大小是：

[
p^d-1.
]

一个 primitive (\bar m)-th root of unity 存在于 (\mathbb F_{p^d}) 中，当且仅当

[
\bar m \mid p^d-1.
]

这等价于：

[
p^d\equiv 1\pmod{\bar m}.
]

所以最小的这样的 (d)，就是

[
d=\operatorname{ord}_{\bar m}(p).
]

这已经说明了为什么 (d) 和 (\mathbb F_{p^d}) 有关。

---

## 2. 更正式地证明

设 (\alpha) 是一个 primitive (\bar m)-th root of unity。也就是说：

[
\alpha^{\bar m}=1,
]

并且 (\alpha) 的阶正好是 (\bar m)。

我们想知道：

[
\alpha
]

在 (\mathbb F_p) 上的最小多项式次数是多少。

这个次数就是：

[
[\mathbb F_p(\alpha):\mathbb F_p].
]

如果能证明这个次数是 (d)，那么就有：

[
\mathbb F_p(\alpha)\cong \mathbb F_{p^d}.
]

---

## 3. Frobenius 轨道

有限域里有一个非常重要的映射：

[
\mathrm{Frob}_p:x\mapsto x^p.
]

它固定 (\mathbb F_p)，所以如果 (\alpha) 是某个元素，那么它在 (\mathbb F_p) 上的共轭是：

[
\alpha,\alpha^p,\alpha^{p^2},\alpha^{p^3},\dots
]

最小多项式的根就是这条 Frobenius 轨道。

因此，最小多项式的次数等于最小的正整数 (r)，使得

[
\alpha^{p^r}=\alpha.
]

现在化简这个条件：

[
\alpha^{p^r}=\alpha
]

等价于

[
\alpha^{p^r-1}=1.
]

因为 (\alpha) 的阶是 (\bar m)，所以这又等价于

[
\bar m\mid p^r-1.
]

也就是：

[
p^r\equiv 1\pmod{\bar m}.
]

最小满足这个条件的 (r)，按照定义就是：

[
d=\operatorname{ord}_{\bar m}(p).
]

所以：

[
[\mathbb F_p(\alpha):\mathbb F_p]=d.
]

因此：

[
\mathbb F_p(\alpha)=\mathbb F_{p^d}.
]

---

## 4. 这和 (\mathfrak p_i) 有什么关系？

论文里每个 prime ideal (\mathfrak p_i) 可以理解成“把 (\zeta_m) 代入某个有限域单位根”的 kernel。

简单说，就是有一个映射：

[
h_i:R=\mathbb Z[\zeta_m]\to \mathbb F_{p^d}
]

满足：

[
h_i(\zeta_m)=\omega_{\bar m}^i.
]

这里 (\omega_{\bar m}) 是有限域里的 primitive (\bar m)-th root of unity。

于是：

[
\mathfrak p_i=\ker h_i.
]

根据第一同构定理：

[
R/\ker h_i \cong \operatorname{im}(h_i).
]

也就是：

[
R/\mathfrak p_i \cong \operatorname{im}(h_i).
]

而 (\omega_{\bar m}^i) 仍然是 primitive (\bar m)-th root，因为 (i\in \mathbb Z_{\bar m}^*)，即 (\gcd(i,\bar m)=1)。

所以：

[
\mathbb F_p(\omega_{\bar m}^i)=\mathbb F_{p^d}.
]

因此：

[
\operatorname{im}(h_i)=\mathbb F_{p^d}.
]

最终得到：

[
R/\mathfrak p_i\cong \mathbb F_{p^d}.
]

这就是“每个 prime ideal 对应的有限域扩张次数是 (d)”的证明。

---

## 5. 为什么所有 (\mathfrak p_i) 的 (d) 一样？

因为不同的 (\mathfrak p_i) 对应的是不同的 primitive (\bar m)-th roots：

[
\omega_{\bar m}^i.
]

但只要

[
\gcd(i,\bar m)=1,
]

(\omega_{\bar m}^i) 的阶仍然是 (\bar m)。

所以判断它需要多大的有限域，还是同一个条件：

[
p^r\equiv 1\pmod{\bar m}.
]

最小的 (r) 仍然是同一个：

[
d=\operatorname{ord}_{\bar m}(p).
]

所以固定 (p,m) 后，所有分裂出来的 (\mathfrak p_i) 都有相同的 residue degree (d)。

---

## 6. 用图 1 的例子验证

图 1 里：

[
p=2,\qquad m=91.
]

因为 (2\nmid 91)，所以：

[
\bar m=91.
]

现在算：

[
d=\operatorname{ord}_{91}(2).
]

也就是最小的 (d)，使得：

[
2^d\equiv 1\pmod{91}.
]

结果是：

[
2^{12}\equiv 1\pmod{91}.
]

所以：

[
d=12.
]

因此每个 prime ideal 的 quotient 都是：

[
R/\mathfrak p_i\cong \mathbb F_{2^{12}}.
]

例如：

[
R/\mathfrak p_1\cong \mathbb F_{2^{12}},
]

[
R/\mathfrak p_{15}\cong \mathbb F_{2^{12}},
]

[
R/\mathfrak p_{22}\cong \mathbb F_{2^{12}},
]

等等。

它们对应不同的 slots，但每个 slot 里的明文元素都来自同一个大小的有限域：

[
\mathbb F_{2^{12}}.
]

---

最短证明就是：

[
\alpha \text{ 是 primitive }\bar m\text{-th root}
]

时，

[
[\mathbb F_p(\alpha):\mathbb F_p]
=================================

\min{r:\alpha^{p^r}=\alpha}
]

而

[
\alpha^{p^r}=\alpha
\iff
\alpha^{p^r-1}=1
\iff
\bar m\mid p^r-1
\iff
p^r\equiv 1\pmod{\bar m}.
]

所以这个最小 (r) 正是：

[
d=\operatorname{ord}_{\bar m}(p).
]

因此

[
R/\mathfrak p_i\cong \mathbb F_{p^d}.
]

===


你卡住的点非常正常：**(\mathfrak p_i) 不是一个映射。**

更准确地说：

[
\mathfrak p_i
]

是一个 **prime ideal**，也就是 (R) 里面的一堆元素。

但它可以通过一个映射 (h_i) 来描述：

[
\mathfrak p_i=\ker(h_i).
]

也就是说：

> (\mathfrak p_i) 不是映射；
> 它是某个映射 (h_i) 的 kernel，也就是“被这个映射送到 0 的所有元素”。

---

## 1. 先用最简单的例子理解

看整数环：

[
\mathbb Z.
]

素数 (5) 对应的 ideal 是：

[
5\mathbb Z={\dots,-10,-5,0,5,10,\dots}.
]

这个 ideal 也可以通过一个映射来描述：

[
h:\mathbb Z\to \mathbb F_5,
]

[
h(a)=a\bmod 5.
]

那么：

[
\ker(h)={a\in\mathbb Z:a\bmod 5=0}.
]

也就是：

[
\ker(h)=5\mathbb Z.
]

所以：

[
5\mathbb Z=\ker(\mathbb Z\to \mathbb F_5).
]

注意：
**(5\mathbb Z) 不是映射，它是“模 5 映射的零集合”。**

---

## 2. 再看多项式例子

比如：

[
\mathbb F_p[X].
]

定义一个映射：

[
h:\mathbb F_p[X]\to \mathbb F_p,
]

[
h(f)=f(3).
]

也就是把 (X) 代成 (3)。

那么：

[
\ker(h)={f(X):f(3)=0}.
]

这个 kernel 就是：

[
(X-3).
]

所以：

[
(X-3)=\ker(f(X)\mapsto f(3)).
]

这里 ((X-3)) 是一个 ideal，不是映射；但它可以被解释成“所有在 (X=3) 处取值为 0 的多项式”。

---

## 3. 回到论文里的 (\mathfrak p_i)

论文里的环是：

[
R=\mathbb Z[\zeta_m].
]

里面的元素可以粗略想成：

[
a(\zeta_m)
]

也就是关于 (\zeta_m) 的多项式表达。

现在论文定义一个映射：

[
h_i:R\to \mathbb F_{p^d}.
]

它的规则是：

[
\zeta_m\mapsto \omega_{\bar m}^i.
]

也就是说：

[
a(\zeta_m)\mapsto a(\omega_{\bar m}^i).
]

这就是“第 (i) 个取值方式”。

然后定义：

[
\mathfrak p_i=\ker(h_i).
]

所以：

[
\mathfrak p_i
=============

{a(\zeta_m)\in R:a(\omega_{\bar m}^i)=0}.
]

翻译成人话：

> (\mathfrak p_i) 是所有在第 (i) 个 slot 里取值为 0 的环元素。

论文 2.1.3 节正是这样说的：(\mathfrak p_i) 是同态 (h_i) 的 kernel，其中 (h_i(\zeta_m)=\omega_{\bar m}^i)。

---

## 4. 为什么可以这样定义 prime ideal？

因为有一个基本事实：

> 一个环映射到 field 时，它的 kernel 是 prime ideal。

这里目标是：

[
\mathbb F_{p^d},
]

这是一个 field。

所以：

[
\ker(h_i)
]

一定是 prime ideal。

于是论文就把这个 prime ideal 记成：

[
\mathfrak p_i.
]

也就是说，(\mathfrak p_i) 的来源是：

[
\text{第 }i\text{ 个取值映射 }h_i
\quad\Rightarrow\quad
\text{它的零集合 }\ker(h_i)
\quad\Rightarrow\quad
\text{一个 prime ideal }\mathfrak p_i.
]

---

## 5. 你可以把它理解成“slot 的零条件”

如果一个明文元素在所有 slot 里的值是：

[
(x_1,x_{15},x_{22},x_3,x_{17},x_{31}),
]

那么属于 (\mathfrak p_{15}) 的意思就是：

[
x_{15}=0.
]

也就是说：

[
\mathfrak p_{15}
]

是“第 15 号 slot 为 0”的那些元素。

所以它不是 slot 本身，也不是映射本身，而是：

> 让第 15 号 slot 读出来等于 0 的所有元素。

---

## 6. 最短总结

[
\mathfrak p_i
]

不是映射。

真正的映射是：

[
h_i:R\to \mathbb F_{p^d}.
]

它负责读取第 (i) 个 slot：

[
a(\zeta_m)\mapsto a(\omega_{\bar m}^i).
]

而

[
\mathfrak p_i=\ker(h_i)
]

表示：

[
\mathfrak p_i
]

是所有“第 (i) 个 slot 读出来为 0”的元素组成的 prime ideal。

所以论文不是把 (\mathfrak p_i) 解释成映射，而是说：**(\mathfrak p_i) 可以通过某个映射的 kernel 来定义。**
===

你问的是为什么要定义这个“读取函数”：

[
a(\zeta_m)\mapsto a(\omega_{\bar m}^i).
]

核心原因是：

**因为 (R=\mathbb Z[\zeta_m]) 里的元素本来就是“关于 (\zeta_m) 的多项式”，所以要从 (R) 里读出一个有限域值，最自然的办法就是把 (\zeta_m) 代成有限域里的某个合适的根。**

---

## 1. 先看 (R=\mathbb Z[\zeta_m]) 里的元素长什么样

环

[
R=\mathbb Z[\zeta_m]
]

可以粗略理解成：

[
R\cong \mathbb Z[X]/(\Phi_m(X)).
]

也就是说，里面的元素可以写成：

[
a(\zeta_m),
]

其中 (a(X)) 是整数系数多项式。

例如：

[
a(\zeta_m)=3+5\zeta_m+7\zeta_m^2.
]

所以如果我们想把这个元素送到某个有限域里，最自然的办法就是：

[
\zeta_m \mapsto \text{某个有限域元素}.
]

一旦决定了 (\zeta_m) 被送到哪里，整个 (a(\zeta_m)) 的值就自动决定了。

---

## 2. 但不能随便代，必须满足同样的关系

(\zeta_m) 不是一个自由变量。它满足：

[
\Phi_m(\zeta_m)=0.
]

所以如果我们定义一个环同态

[
h:R\to \mathbb F
]

并且想让

[
h(\zeta_m)=\alpha,
]

那么 (\alpha) 必须满足：

[
\Phi_m(\alpha)=0
]

在目标有限域里成立。

否则映射就不合法。

因为在 (R) 里有：

[
\Phi_m(\zeta_m)=0,
]

映射过去后也必须有：

[
\Phi_m(h(\zeta_m))=0.
]

所以问题变成：

> 在有限域里，哪些元素可以当作 (\zeta_m) 的替身？

答案就是论文里选的：

[
\omega_{\bar m}^i.
]

---

## 3. 为什么是 (\omega_{\bar m}^i)，不是 (\omega_m^i)？

这是因为我们在模 (p) 的有限域里工作。

论文先把

[
m=\bar m\cdot p^k
]

分解出来，其中

[
p\nmid \bar m.
]

在特征 (p) 的有限域里，(p)-power 那一部分会“塌掉”。直观地说，有限域的乘法群大小是：

[
p^d-1,
]

它不能含有阶为 (p) 的元素，因为

[
p\nmid p^d-1.
]

所以在 characteristic (p) 的有限域里，不会有真正的 primitive (p^k)-th root of unity。

因此 (m) 的 (p^k) 那部分不能保留下来，真正有意义的是去掉 (p)-power 后的部分：

[
\bar m.
]

所以要选的是有限域里的 primitive (\bar m)-th root：

[
\omega_{\bar m}.
]

这就是为什么论文写的是：

[
h_i(\zeta_m)=\omega_{\bar m}^i,
]

而不是

[
\omega_m^i.
]

论文 2.1.3 节正是这样定义 prime ideal (\mathfrak p_i)：令 (h_i:R\to \mathbb F_{p^d})，并设 (h_i(\zeta_m)=\omega_{\bar m}^i)，然后 (\mathfrak p_i=\ker h_i)。

---

## 4. 为什么这个定义会产生 slot？

因为这个映射：

[
h_i:R\to \mathbb F_{p^d}
]

会把一个环元素

[
a(\zeta_m)
]

读成一个有限域元素：

[
h_i(a(\zeta_m))=a(\omega_{\bar m}^i).
]

这正好就是第 (i) 个 slot 的值。

比如：

[
a(\zeta_m)=3+5\zeta_m+7\zeta_m^2.
]

那么第 (i) 个 slot 是：

[
a(\omega_{\bar m}^i)
====================

3+5\omega_{\bar m}^i+7(\omega_{\bar m}^i)^2,
]

在有限域 (\mathbb F_{p^d}) 里计算。

所以：

[
a(\zeta_m)\mapsto a(\omega_{\bar m}^i)
]

不是随便定义的，而是“把环元素在第 (i) 个有限域根上求值”。

这和多项式求值非常像：

[
f(X)\mapsto f(3).
]

只是这里不是把 (X) 代成普通数字 (3)，而是把 (\zeta_m) 代成有限域里的单位根 (\omega_{\bar m}^i)。

---

## 5. 为什么不同的 (i) 给不同的 slot？

因为不同的 (i) 给出不同的根：

[
\omega_{\bar m}^i.
]

每个根对应一种不同的“读法”：

[
h_i(a)=a(\omega_{\bar m}^i).
]

这些不同读法就是不同的 CRT 分量，也就是不同的 plaintext slots。

但注意，(i) 和 (ip) 会给同一个 prime ideal。原因是有限域里的 Frobenius：

[
x\mapsto x^p
]

会把

[
\omega_{\bar m}^i
]

送到

[
\omega_{\bar m}^{ip}.
]

它们是同一个最小多项式下的共轭根，所以对应同一个 irreducible factor，也就是同一个 prime ideal。

因此 slot 不是由单个 (i) 编号，而是由一组

[
i\langle p\rangle
]

编号。

---

## 6. 为什么 kernel 就是 (\mathfrak p_i)？

映射

[
h_i(a)=a(\omega_{\bar m}^i)
]

读取第 (i) 个 slot。

那么

[
\ker h_i
]

就是所有第 (i) 个 slot 读出来为 (0) 的元素：

[
\ker h_i
========

{a(\zeta_m)\in R:a(\omega_{\bar m}^i)=0}.
]

论文把这个 kernel 记成：

[
\mathfrak p_i.
]

所以 (\mathfrak p_i) 的含义是：

> 在第 (i) 个 slot 上为零的所有环元素。

这就是为什么 prime ideal 可以通过这个求值映射来定义。

---

## 7. 最直观的一句话

你可以把

[
a(\zeta_m)\mapsto a(\omega_{\bar m}^i)
]

理解成：

> 把抽象的 cyclotomic 明文元素 (a(\zeta_m))，在第 (i) 个有限域单位根上求值，从而得到第 (i) 个 slot 的内容。

之所以这样定义，是因为：

[
R=\mathbb Z[\zeta_m]
]

本来就是由 (\zeta_m) 生成的；而在模 (p) 的有限域里，(\zeta_m) 的合法替身就是 primitive (\bar m)-th root

[
\omega_{\bar m}^i.
]

所以这个求值映射正好给出了 CRT 分解中的一个分量，也就是一个 plaintext slot。

===

对，这里的逻辑应该这样接：

不是说：

[
m=\bar m p^k
]

以后我们“人为把 (p^k) 去掉”。

而是说：**在模 (p) 的世界里，(p^k) 这部分不会产生新的根，只会变成重根，也就是 ramification。**

所以真正决定“有哪些不同 slot / 不同 prime ideals”的，是去掉 (p)-power 后的部分：

[
\bar m.
]

而被去掉的 (p^k) 并没有消失，它变成了指数：

[
e=\varphi(p^k).
]

这就是论文里说 (pR) 分解时：

[
pR=\prod_{i\in G}\mathfrak p_i^e
]

里面的那个 (e)。

---

## 1. 关键点：有限域里没有 (p)-power 阶的非零元素

在特征 (p) 的有限域 (\mathbb F_{p^d}) 里，非零元素组成乘法群：

[
\mathbb F_{p^d}^{*}.
]

它的大小是：

[
p^d-1.
]

注意：

[
p\nmid p^d-1.
]

所以这个乘法群里不可能有阶为 (p)、(p^2)、(p^k) 的元素。

因此，如果

[
m=\bar m p^k,
]

那么有限域里不可能有真正的 primitive (m)-th root of unity，因为 primitive (m)-th root 的阶必须正好是：

[
m=\bar m p^k,
]

里面含有 (p^k) 因子。

---

## 2. 那 (\zeta_m) 模 (p) 后会变成什么？

它不能变成真正的 (m) 次本原单位根。

但它可以变成一个 primitive (\bar m)-th root。

也就是说，在模 (p) 以后，(m) 里的 (p^k) 那部分“塌掉”了，只剩下：

[
\bar m.
]

更准确地说：

[
\Phi_m(X)
]

模 (p) 以后，它的根会变成 primitive (\bar m)-th roots，但是这些根会带重数。

这个重数就是：

[
e=\varphi(p^k).
]

所以：

[
p^k\text{ 部分} \quad \longrightarrow \quad \text{重数 / ramification}
]

而不是新的不同根。

---

## 3. 一个具体例子：(p=2,\ m=12)

这里：

[
m=12=3\cdot 2^2.
]

所以：

[
\bar m=3,\qquad p^k=4.
]

在特征 (2) 的有限域里，不可能有真正的 (4) 次单位根，也不可能有真正的 (12) 次本原单位根。

因为任何有限域 (\mathbb F_{2^d}^{*}) 的大小是：

[
2^d-1,
]

这是奇数，不可能被 (2) 或 (4) 整除。

所以 (12) 次单位根的 (2^2) 那部分没法保留。

看多项式也能看到这个现象：

[
\Phi_{12}(X)=X^4-X^2+1.
]

模 (2) 以后：

[
\Phi_{12}(X)\equiv X^4+X^2+1.
]

而在特征 (2) 里：

[
X^4+X^2+1=(X^2+X+1)^2.
]

注意：

[
X^2+X+1=\Phi_3(X).
]

所以：

[
\Phi_{12}(X)\bmod 2
===================

\Phi_3(X)^2.
]

这句话非常说明问题：

[
12=3\cdot 4
]

模 (2) 后，真正剩下的根是 **3 次本原单位根**，但是每个根重复了 (2) 次。

而这个 (2) 正好是：

[
e=\varphi(4)=2.
]

所以 (2^2) 部分没有提供新的根，而是提供了重数。

---

## 4. 回到一般情况

一般地，写：

[
m=\bar m p^k,\qquad p\nmid \bar m.
]

在特征 (p) 里，会发生类似现象：

[
\Phi_m(X)\bmod p
]

的根来自：

[
\Phi_{\bar m}(X),
]

也就是 primitive (\bar m)-th roots。

但它们会带有重数：

[
e=\varphi(p^k).
]

所以论文才定义：

[
\omega_{\bar m}
]

为有限域里的 primitive (\bar m)-th root，而不是 primitive (m)-th root。

然后定义：

[
h_i(\zeta_m)=\omega_{\bar m}^i.
]

这不是乱来，而是因为模 (p) 后，(\zeta_m) 的合法“有限域替身”就是 primitive (\bar m)-th root。

---

## 5. 那 (p^k) 去哪里了？

它没有消失。

它出现在分解公式的指数里：

[
pR=\prod_{i\in G}\mathfrak p_i^e,
]

其中：

[
e=\varphi(p^k).
]

所以：

* (\bar m) 决定有多少个不同的 prime ideals / slots；
* (p^k) 决定每个 prime ideal 重复几次，也就是 ramification。

如果 (p\nmid m)，那么：

[
k=0,\qquad p^k=1,\qquad e=\varphi(1)=1.
]

这时没有重数。

比如图 1 里：

[
p=2,\qquad m=91.
]

因为：

[
2\nmid 91,
]

所以：

[
\bar m=91,\qquad e=1.
]

因此图 1 里每个 prime ideal 都只出现一次。

---

## 6. 最短总结

你可以把这件事理解成：

[
m=\bar m p^k
]

里面有两部分。

在模 (p) 的有限域里：

[
p^k\text{ 部分不能表现为新的单位根}
]

因为有限域的非零乘法群大小 (p^d-1) 不含 (p) 因子。

所以：

[
p^k\text{ 部分变成重数 }e=\varphi(p^k),
]

而真正产生不同根、不同 slot、不同 prime ideals 的，是：

[
\bar m.
]

因此论文使用：

[
\omega_{\bar m}
]

而不是：

[
\omega_m.
]

一句话：**去掉 (p^k) 不是忽略它，而是因为在特征 (p) 下它不产生新根，只产生 ramification；不同的根由 (\bar m) 决定，重复次数由 (p^k) 决定。**

===

设

[
m=\bar m p^k,\qquad p\nmid \bar m.
]

我们要证明：

[
\Phi_m(X)\bmod p
]

的根和

[
\Phi_{\bar m}(X)\bmod p
]

的根相同，只是每个根的重数变成

[
e=\varphi(p^k).
]

更强地，可以证明这个多项式恒等式：

[
\boxed{
\Phi_{\bar m p^k}(X)\equiv \Phi_{\bar m}(X)^{\varphi(p^k)}
\pmod p
}
]

这就是“(p^k) 部分变成重数”的精确含义。论文 prime splitting 部分里把这个重数记成 ramification index (e=\varphi(p^k))。

---

## 1. 先用一个 cyclotomic polynomial 恒等式

令

[
r=\bar m.
]

因为

[
p\nmid r,
]

对 (k\ge 1)，有标准恒等式：

[
\Phi_{rp^k}(X)
==============

\frac{\Phi_r(X^{p^k})}{\Phi_r(X^{p^{k-1}})}.
]

先接受这个公式，它来自 cyclotomic polynomial 的基本关系：

[
X^n-1=\prod_{d\mid n}\Phi_d(X).
]

---

## 2. 模 (p) 后使用 Freshman's dream

在特征 (p) 里，有：

[
(a+b)^p=a^p+b^p.
]

所以对整数系数多项式 (F(X))，模 (p) 后有：

[
F(X^p)\equiv F(X)^p\pmod p.
]

更一般地：

[
F(X^{p^k})\equiv F(X)^{p^k}\pmod p.
]

于是：

[
\Phi_r(X^{p^k})
\equiv
\Phi_r(X)^{p^k}
\pmod p,
]

[
\Phi_r(X^{p^{k-1}})
\equiv
\Phi_r(X)^{p^{k-1}}
\pmod p.
]

---

## 3. 代回去

原来：

[
\Phi_{rp^k}(X)
==============

\frac{\Phi_r(X^{p^k})}{\Phi_r(X^{p^{k-1}})}.
]

模 (p) 后就变成：

[
\Phi_{rp^k}(X)
\equiv
\frac{\Phi_r(X)^{p^k}}{\Phi_r(X)^{p^{k-1}}}
\pmod p.
]

所以：

[
\Phi_{rp^k}(X)
\equiv
\Phi_r(X)^{p^k-p^{k-1}}
\pmod p.
]

而

[
p^k-p^{k-1}
===========

# p^{k-1}(p-1)

\varphi(p^k).
]

因此：

[
\boxed{
\Phi_{rp^k}(X)
\equiv
\Phi_r(X)^{\varphi(p^k)}
\pmod p
}
]

也就是：

[
\boxed{
\Phi_m(X)
\equiv
\Phi_{\bar m}(X)^e
\pmod p
}
]

其中

[
e=\varphi(p^k).
]

---

## 4. 这说明什么？

如果

[
\Phi_m(X)\equiv \Phi_{\bar m}(X)^e\pmod p,
]

那么它们的根集合相同。

因为一个数 (\alpha) 满足：

[
\Phi_m(\alpha)=0
]

等价于：

[
\Phi_{\bar m}(\alpha)^e=0.
]

在域里，

[
A^e=0 \iff A=0.
]

所以：

[
\Phi_m(\alpha)=0
\iff
\Phi_{\bar m}(\alpha)=0.
]

但是重数变了：(\Phi_{\bar m}) 的每个根，在 (\Phi_m) 里重复了 (e) 次。

这就是：

[
p^k\text{ 部分不产生新根，只产生重数。}
]

---

## 5. 例子：(p=2,\ m=12)

这里：

[
12=3\cdot 2^2.
]

所以：

[
\bar m=3,\qquad k=2.
]

于是：

[
e=\varphi(2^2)=2.
]

我们有：

[
\Phi_{12}(X)=X^4-X^2+1.
]

模 (2) 后：

[
\Phi_{12}(X)\equiv X^4+X^2+1.
]

而

[
\Phi_3(X)=X^2+X+1.
]

在模 (2) 下：

[
(X^2+X+1)^2
===========

X^4+2X^3+3X^2+2X+1
\equiv
X^4+X^2+1.
]

所以：

[
\Phi_{12}(X)
\equiv
\Phi_3(X)^2
\pmod 2.
]

这正好说明：

[
12\text{ 次本原单位根}
]

在特征 (2) 里退化成：

[
3\text{ 次本原单位根}
]

但每个根重复 (2) 次。

---

## 6. 和 prime ideal 分解的关系

因此，当

[
m=\bar m p^k
]

时：

* 不同的根由 (\Phi_{\bar m}) 决定；
* 每个根重复 (e=\varphi(p^k)) 次；
* 所以 prime ideal 分解里会出现：

[
pR=\prod_i \mathfrak p_i^e.
]

这里的 (\mathfrak p_i) 来自 (\Phi_{\bar m}) 的不同 irreducible factors，而指数 (e) 来自 (p^k) 造成的重数。

一句话总结：

[
\boxed{
\Phi_m(X)\bmod p
================

\Phi_{\bar m}(X)^{\varphi(p^k)}
}
]

所以模 (p) 后，(m) 的 (p^k) 部分不会产生新的有限域单位根，只会让已有的 (\bar m)-次单位根带上重数。

===

可以，我们先不管公式证明，只讲**直观图像**和**它有什么用**。

你问的是这个结论：

[
\Phi_m(X)\bmod p
================

\Phi_{\bar m}(X)^{\varphi(p^k)}
]

其中

[
m=\bar m\cdot p^k,\qquad p\nmid \bar m.
]

这句话直观上是在说：

> **当我们进入模 (p) 的世界后，(m) 里面含有的 (p^k) 那一部分会“塌掉”，不再产生新的不同根；它只会让已有的根重复出现。**

---

## 1. 先用一句话理解

假设

[
m=\bar m\cdot p^k.
]

那么 (m)-次单位根本来包含两部分信息：

[
\bar m\text{ 部分} \quad+\quad p^k\text{ 部分}.
]

但是在模 (p) 的有限域世界里，**没有真正的 (p)-power 阶单位根**。

所以 (p^k) 这部分没法表现成“新的不同根”。

它只能变成：

[
\text{重复次数}
]

也就是论文里的 ramification index：

[
e=\varphi(p^k).
]

所以：

[
\bar m \text{ 决定不同 roots / 不同 prime ideals / 不同 slots，}
]

[
p^k \text{ 决定这些东西重复几次。}
]

这就是这个定理的直观含义。论文 prime splitting 那段正是在用这个思想定义 (\bar m,d,e,f)，最后描述 (pR) 如何分解成 prime ideals。

---

## 2. 一个生活化类比

想象 (m)-次单位根是一块钟表。

比如

[
m=12.
]

12 小时钟表可以看作：

[
12=3\cdot 4.
]

也就是：

* 一个 3 的部分；
* 一个 4 的部分。

现在如果我们在模 (2) 的世界里看它，因为

[
4=2^2
]

是 (2)-power 部分，所以这个“4 的方向”在模 2 里分不出新的位置了。

于是 12 个方向不会真的保留下来，而是退化成 3 个方向，但每个方向带重复。

所以：

[
12=3\cdot 2^2
]

在模 (2) 后，真正留下来的不同根来自：

[
3.
]

而

[
2^2
]

那部分变成重复次数。

这就是：

[
\Phi_{12}(X)\bmod 2
===================

\Phi_3(X)^2.
]

意思是：

> 12 次本原单位根在模 2 后，看起来像 3 次本原单位根，但每个根重复了 2 次。

---

## 3. 为什么说 (p^k) 部分“没有新的根”？

因为有限域 (\mathbb F_{p^d}) 的非零元素组成一个乘法群，它的大小是：

[
p^d-1.
]

这个数永远不能被 (p) 整除。

比如 (p=2) 时：

[
2^d-1
]

永远是奇数。

所以在 (\mathbb F_{2^d}) 里面，不可能有真正阶数为 (2,4,8,\dots) 的非零元素。

这就是为什么如果

[
m=12=3\cdot 4,
]

在模 2 里，那个 (4) 的部分不能作为新的单位根存在。

只剩下：

[
3
]

这部分还能作为真正的单位根。

---

## 4. 那这个结论有什么用？

它的用途是：**告诉我们 (pR) 在 cyclotomic ring 里怎么分裂。**

在论文里，我们关心：

[
R=\mathbb Z[\zeta_m].
]

然后看整数素数 (p) 在这个环里生成的 ideal：

[
pR.
]

它会分解成：

[
pR=\prod_i \mathfrak p_i^e.
]

这里：

* (\mathfrak p_i)：不同的 prime ideals；
* (e)：每个 prime ideal 重复几次；
* 不同 (\mathfrak p_i) 的数量对应 plaintext slots 的数量。

这个定理告诉我们：

[
m=\bar m p^k
]

时：

[
\bar m
]

决定有哪些不同的 (\mathfrak p_i)，也就是有多少个不同 slots；

而：

[
p^k
]

决定每个 (\mathfrak p_i) 出现几次，也就是 ramification：

[
e=\varphi(p^k).
]

---

## 5. 和图 1 的关系

图 1 里：

[
p=2,\qquad m=91.
]

因为：

[
2\nmid 91,
]

所以：

[
m=\bar m,\qquad k=0.
]

也就是说没有 (2^k) 部分。

因此：

[
e=\varphi(1)=1.
]

所以图 1 里：

[
2R
]

分裂成六个 prime ideals，但每个只出现一次：

[
2R=
\mathfrak p_1
\mathfrak p_{15}
\mathfrak p_{22}
\mathfrak p_3
\mathfrak p_{17}
\mathfrak p_{31}.
]

如果 (m) 里面含有 (2)-power 因子，比如 (m=12)，那就会出现重复：

[
2R=\prod_i \mathfrak p_i^e,
]

其中 (e>1)。

---

## 6. 你真正需要记住的不是公式

这个公式：

[
\Phi_m(X)\bmod p
================

\Phi_{\bar m}(X)^{\varphi(p^k)}
]

不用死记。

你只要记住这句话：

> **模 (p) 以后，(m) 中和 (p) 互素的部分 (\bar m) 决定“不同的 slot”；(m) 中的 (p^k) 部分不会产生新 slot，只会造成重复，也就是 ramification。**

所以论文一开始才要把：

[
m
]

拆成：

[
m=\bar m p^k.
]

这是为了区分：

[
\bar m:\quad 负责不同根、不同 prime ideals、不同 slots；
]

[
p^k:\quad 负责重复次数 (e)。
]

一句话总结：**这个定理的作用是解释为什么 prime splitting 里要用 (\bar m) 来编号 slots，而把 (p^k) 放进指数 (e) 里。**

