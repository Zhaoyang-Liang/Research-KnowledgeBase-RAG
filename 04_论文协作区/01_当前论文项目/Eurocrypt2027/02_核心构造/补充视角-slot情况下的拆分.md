以下是将您提供的完整推导内容转换为标准 Markdown 语法的结果。所有数学表达式均已使用 LaTeX 格式包裹，矩阵、公式和代码块结构完整保留，标题、列表和引用均采用 Markdown 规范语法。


# 可以。我们慢慢来。你问的是这个问题：

> CKKS 里如果我想在 **slot 层面**做拆分，那么回到多项式 \(P(X)\) 层面，到底对应什么样的多项式分解？

我先给结论，然后一步步推。

---

# 0. 先说结论

有两种完全不同的“拆分”。

## 第一种：系数拆分

这是简单的：

\[
P(X)=A(X^2)+X B(X^2).
\]

其中：

\[
A(Y)=a_0+a_2Y+a_4Y^2+\cdots,
\]
\[
B(Y)=a_1+a_3Y+a_5Y^2+\cdots.
\]

这叫 **coefficient split**，也就是把多项式系数按偶数项、奇数项拆开。

---

## 第二种：slot 拆分

如果你想要的是：

\[
\text{大环 slots} \longrightarrow \text{小环 slots}
\]

比如想把一对大环 slots 变成两个小环密文的 slots，那么对应到多项式就不是简单的 \(A(Y),B(Y)\)。

它会变成：

\[
Q^+(Y)=A(Y)+h(Y)B(Y),
\]
\[
Q^-(Y)=A(Y)-h(Y)B(Y),
\]

其中 \(h(Y)\) 是一个“平方根选择多项式”，满足：

\[
h(\theta_j)=\alpha_j,\qquad \alpha_j^2=\theta_j.
\]

这才是我前面说的东西。

---

# 1. 先建立大环和小环

设大环是：

\[
R_N=\mathbb C[X]/(X^N+1).
\]

小环是：

\[
R_n=\mathbb C[Y]/(Y^n+1),
\]

其中：

\[
n=\frac{N}{2}.
\]

我们令：

\[
Y=X^2.
\]

所以大环里的一个多项式：

\[
P(X)
\]

可以写成：

\[
P(X)=A(X^2)+X B(X^2).
\]

这只是普通的奇偶项分解。

---

# 2. 用系数向量来看 coefficient split

设：

\[
P(X)=a_0+a_1X+a_2X^2+\cdots+a_{N-1}X^{N-1}.
\]

定义两个小环多项式：

\[
A(Y)=a_0+a_2Y+a_4Y^2+\cdots+a_{N-2}Y^{n-1},
\]
\[
B(Y)=a_1+a_3Y+a_5Y^2+\cdots+a_{N-1}Y^{n-1}.
\]

于是：

\[
P(X)=A(X^2)+X B(X^2).
\]

这一步非常简单。

它对应的是：

```text
P 的系数:
(a0, a1, a2, a3, a4, a5, ...)

拆成:

A 的系数:
(a0, a2, a4, ...)

B 的系数:
(a1, a3, a5, ...)
```

所以 coefficient split 是系数层面的奇偶拆分。

---

# 3. 现在看 slot：大环 roots 和小环 roots

大环：

\[
X^N+1=0
\]

的 roots 是：

\[
\alpha_0,\alpha_1,\dots,\alpha_{N-1},
\]

其中可以取：

\[
\alpha_j=\exp\left(\frac{(2j+1)\pi i}{N}\right).
\]

小环：

\[
Y^n+1=0
\]

的 roots 是：

\[
\theta_0,\theta_1,\dots,\theta_{n-1},
\]

其中：

\[
\theta_j=\exp\left(\frac{(2j+1)\pi i}{n}\right).
\]

因为：

\[
n=\frac N2,
\]

所以：

\[
\theta_j=\exp\left(\frac{(2j+1)2\pi i}{N}\right).
\]

注意：

\[
\alpha_j^2 = \exp\left(\frac{(2j+1)2\pi i}{N}\right) = \theta_j.
\]

所以：

\[
\boxed{\theta_j=\alpha_j^2.}
\]

同时：

\[
\alpha_{j+n}=-\alpha_j.
\]

因为：

\[
\alpha_{j+n} = \exp\left(\frac{(2(j+n)+1)\pi i}{N}\right) = -\exp\left(\frac{(2j+1)\pi i}{N}\right).
\]

所以：

\[
\boxed{\alpha_j^2=\alpha_{j+n}^2=\theta_j.}
\]

这句话非常重要：

```text
小环的一个 slot root theta_j
对应大环的两个 slot roots:

alpha_j 和 -alpha_j
```

也就是：

\[
\alpha_j \longmapsto \theta_j,
\]
\[
-\alpha_j \longmapsto \theta_j.
\]

---

# 4. 大环 slots 是什么？

大环多项式 \(P(X)\) 的 slot values 是评价值：

\[
P(\alpha_0),P(\alpha_1),\dots,P(\alpha_{N-1}).
\]

为了成对看，我们把它们写成：

\[
z_j^+=P(\alpha_j),
\]
\[
z_j^-=P(-\alpha_j)=P(\alpha_{j+n}).
\]

其中：

\[
j=0,\dots,n-1.
\]

所以一对大环 slots 是：

\[
\left(P(\alpha_j),P(-\alpha_j)\right).
\]

---

# 5. 把 \(P(X)=A(X^2)+XB(X^2)\) 代入 slot

现在代入：

\[
P(X)=A(X^2)+X B(X^2).
\]

对 \(X=\alpha_j\)，有：

\[
P(\alpha_j)=A(\alpha_j^2)+\alpha_j B(\alpha_j^2).
\]

因为：

\[
\alpha_j^2=\theta_j,
\]

所以：

\[
P(\alpha_j)=A(\theta_j)+\alpha_j B(\theta_j).
\]

再对 \(X=-\alpha_j\)，有：

\[
P(-\alpha_j)=A((-\alpha_j)^2)+(-\alpha_j)B((-\alpha_j)^2).
\]

因为：

\[
(-\alpha_j)^2=\theta_j,
\]

所以：

\[
P(-\alpha_j)=A(\theta_j)-\alpha_j B(\theta_j).
\]

于是得到两条式子：

\[
P(\alpha_j)=A(\theta_j)+\alpha_j B(\theta_j),
\]
\[
P(-\alpha_j)=A(\theta_j)-\alpha_j B(\theta_j).
\]

---

# 6. 这两条式子说明什么？

令：

\[
u_j=A(\theta_j),
\]
\[
v_j=B(\theta_j).
\]

也就是 \(A,B\) 在小环 slots 上的值。

令：

\[
z_j^+=P(\alpha_j),
\]
\[
z_j^-=P(-\alpha_j).
\]

那么上面的两条式子就是：

\[
z_j^+=u_j+\alpha_j v_j,
\]
\[
z_j^-=u_j-\alpha_j v_j.
\]

写成矩阵：

\[
\begin{pmatrix}
z_j^+\\
z_j^-
\end{pmatrix} =
\begin{pmatrix}
1 & \alpha_j\\
1 & -\alpha_j
\end{pmatrix}
\begin{pmatrix}
u_j\\
v_j
\end{pmatrix}.
\]

反过来：

\[
\begin{pmatrix}
u_j\\
v_j
\end{pmatrix}
=
\begin{pmatrix}
\frac12 & \frac12\\
\frac{1}{2\alpha_j} & -\frac{1}{2\alpha_j}
\end{pmatrix}
\begin{pmatrix}
z_j^+\\
z_j^-
\end{pmatrix}.
\]

这就是最关键的矩阵关系。

---

# 7. 解释这个矩阵关系

coefficient split 输出的是：

\[
A(Y),\quad B(Y).
\]

它们的小环 slot 分别是：

\[
u_j=A(\theta_j),\quad v_j=B(\theta_j).
\]

但这两个值不是原来的大环 slots。

它们是：

\[
u_j=\frac{P(\alpha_j)+P(-\alpha_j)}{2},
\]
\[
v_j=\frac{P(\alpha_j)-P(-\alpha_j)}{2\alpha_j}.
\]

所以：

```text
coefficient split 在 slot 层面做的是：

一对大环 slots
(P(alpha_j), P(-alpha_j))

变成：

pairwise average:
(P(alpha_j)+P(-alpha_j))/2

weighted pairwise difference:
(P(alpha_j)-P(-alpha_j))/(2 alpha_j)
```

也就是说，coefficient split 不是“保留某个 slot”，而是“把一对 slots 做和差变换”。

---

# 8. 这就是为什么 coefficient split 不等于 slot split

如果你只是把系数拆成偶数和奇数：

\[
P(X)\mapsto A(Y),B(Y),
\]

那么在 slot 层面实际得到的是：

\[
\left(
\frac{P(\alpha_j)+P(-\alpha_j)}{2},
\frac{P(\alpha_j)-P(-\alpha_j)}{2\alpha_j}
\right).
\]

但如果你真正想要的是保留：

\[
P(\alpha_j)
\]

和：

\[
P(-\alpha_j)
\]

那你不应该直接输出 \(A(Y)\) 和 \(B(Y)\)。

你应该输出另外两个小环多项式：

\[
Q^+(Y),\quad Q^-(Y),
\]

使得：

\[
Q^+(\theta_j)=P(\alpha_j),
\]
\[
Q^-(\theta_j)=P(-\alpha_j).
\]

---

# 9. 如何构造 \(Q^+(Y)\) 和 \(Q^-(Y)\)？

我们已经知道：

\[
P(\alpha_j)=A(\theta_j)+\alpha_j B(\theta_j),
\]
\[
P(-\alpha_j)=A(\theta_j)-\alpha_j B(\theta_j).
\]

如果存在一个小环多项式 \(h(Y)\)，满足：

\[
h(\theta_j)=\alpha_j,
\]

那么就可以定义：

\[
Q^+(Y)=A(Y)+h(Y)B(Y),
\]
\[
Q^-(Y)=A(Y)-h(Y)B(Y).
\]

验证一下：

\[
Q^+(\theta_j)=A(\theta_j)+h(\theta_j)B(\theta_j)=A(\theta_j)+\alpha_j B(\theta_j)=P(\alpha_j).
\]

同理：

\[
Q^-(\theta_j)=A(\theta_j)-h(\theta_j)B(\theta_j)=A(\theta_j)-\alpha_j B(\theta_j)=P(-\alpha_j).
\]

所以：

\[
\boxed{Q^+(Y)=A(Y)+h(Y)B(Y)}
\]
\[
\boxed{Q^-(Y)=A(Y)-h(Y)B(Y)}
\]

就是 slot-preserving split 对应的小环多项式。

---

# 10. 那 \(h(Y)\) 是什么？

它是一个小环多项式，满足：

\[
h(\theta_j)=\alpha_j,\qquad \alpha_j^2=\theta_j.
\]

所以它在每个小环 root \(\theta_j\) 上，选择一个平方根 \(\alpha_j\)。

因此：

\[
h(\theta_j)^2=\theta_j.
\]

这意味着在商环：

\[
\mathbb C[Y]/(Y^n+1)
\]

里，

\[
h(Y)^2\equiv Y \pmod{Y^n+1}.
\]

所以可以把 \(h(Y)\) 理解成：

```text
Y 的一个平方根选择函数
```

但它通常不是简单的 \(Y^{1/2}\)，也不是一个低阶 monomial。

它是通过插值得到的多项式。

---

# 11. \(h(Y)\) 的矩阵形式

现在把它写成矩阵。

设小环 evaluation matrix 是：

\[
F_n=
\begin{pmatrix}
1 & \theta_0 & \theta_0^2 & \cdots & \theta_0^{n-1}\\
1 & \theta_1 & \theta_1^2 & \cdots & \theta_1^{n-1}\\
\vdots & \vdots & \vdots & & \vdots\\
1 & \theta_{n-1} & \theta_{n-1}^2 & \cdots & \theta_{n-1}^{n-1}
\end{pmatrix}.
\]

如果 \(h(Y)\) 的系数向量是：

\[
h=
\begin{pmatrix}
h_0\\
h_1\\
\vdots\\
h_{n-1}
\end{pmatrix},
\]

那么：

\[
F_n h =
\begin{pmatrix}
h(\theta_0)\\
h(\theta_1)\\
\vdots\\
h(\theta_{n-1})
\end{pmatrix}
=
\begin{pmatrix}
\alpha_0\\
\alpha_1\\
\vdots\\
\alpha_{n-1}
\end{pmatrix}.
\]

所以：

\[
h = F_n^{-1}
\begin{pmatrix}
\alpha_0\\
\alpha_1\\
\vdots\\
\alpha_{n-1}
\end{pmatrix}.
\]

这就是 \(h(Y)\) 的插值公式。

---

# 12. 全局矩阵形式：coefficient split

现在我们写完整矩阵。

令大环多项式系数向量为：

\[
p=
\begin{pmatrix}
a_0\\
a_1\\
a_2\\
a_3\\
\vdots\\
a_{N-1}
\end{pmatrix}.
\]

定义两个选择矩阵：

\[
E=
\begin{pmatrix}
1&0&0&0&\cdots\\
0&0&1&0&\cdots\\
0&0&0&0&1&\cdots\\
\vdots&&&&
\end{pmatrix},
\]

它取偶数系数：

\[
E p=
\begin{pmatrix}
a_0\\
a_2\\
a_4\\
\vdots
\end{pmatrix}.
\]

再定义：

\[
O=
\begin{pmatrix}
0&1&0&0&\cdots\\
0&0&0&1&\cdots\\
0&0&0&0&0&1&\cdots\\
\vdots&&&&
\end{pmatrix},
\]

它取奇数系数：

\[
O p=
\begin{pmatrix}
a_1\\
a_3\\
a_5\\
\vdots
\end{pmatrix}.
\]

于是：

\[
A\text{ 的系数向量}=E p,
\]
\[
B\text{ 的系数向量}=O p.
\]

这就是 coefficient split 的矩阵形式：

\[
p \longmapsto (Ep,Op).
\]

---

# 13. 全局矩阵形式：slot values

大环 evaluation matrix 是：

\[
F_N=
\begin{pmatrix}
1 & \alpha_0 & \alpha_0^2 & \cdots & \alpha_0^{N-1}\\
1 & \alpha_1 & \alpha_1^2 & \cdots & \alpha_1^{N-1}\\
\vdots & \vdots & \vdots & & \vdots\\
1 & \alpha_{N-1} & \alpha_{N-1}^2 & \cdots & \alpha_{N-1}^{N-1}
\end{pmatrix}.
\]

那么大环 slot vector 是：

\[
z=F_N p.
\]

也就是说：

\[
z_j=P(\alpha_j).
\]

小环 \(A(Y)\) 的 slot vector 是：

\[
u=F_n E p.
\]

小环 \(B(Y)\) 的 slot vector 是：

\[
v=F_n O p.
\]

所以 coefficient split 在 slot 层面的矩阵形式是：

\[
p \longmapsto \left(F_n E p,\ F_n O p\right).
\]

如果用大环 slots \(z\) 表示，由于：

\[
p=F_N^{-1}z,
\]

所以：

\[
u=F_n E F_N^{-1}z,
\]
\[
v=F_n O F_N^{-1}z.
\]

这就是我之前说的：

\[
\operatorname{DFT}_{n} \circ \mathrm{EvenCoeff} \circ \operatorname{iDFT}_{N}.
\]

用标准矩阵写就是：

\[
u=F_n E F_N^{-1}z,\qquad v=F_n O F_N^{-1}z.
\]

这说明 coefficient split 对 slots 是一个线性变换，不是简单选择坐标。

---

# 14. 全局矩阵形式：pairwise block

如果我们把大环 slots 按照配对排列：

\[
z^+=
\begin{pmatrix}
P(\alpha_0)\\
P(\alpha_1)\\
\vdots\\
P(\alpha_{n-1})
\end{pmatrix},
\qquad
z^-=
\begin{pmatrix}
P(-\alpha_0)\\
P(-\alpha_1)\\
\vdots\\
P(-\alpha_{n-1})
\end{pmatrix}.
\]

定义对角矩阵：

\[
D=
\begin{pmatrix}
\alpha_0 & & \\
& \alpha_1 & \\
& & \ddots\\
& & & \alpha_{n-1}
\end{pmatrix}.
\]

那么我们有：

\[
z^+=u+Dv,\qquad z^-=u-Dv.
\]

矩阵写成：

\[
\begin{pmatrix}
z^+\\
z^-
\end{pmatrix}
=
\begin{pmatrix}
I & D\\
I & -D
\end{pmatrix}
\begin{pmatrix}
u\\
v
\end{pmatrix}.
\]

反过来：

\[
u=\frac{z^++z^-}{2},\qquad v=\frac{1}{2}D^{-1}(z^+-z^-).
\]

矩阵写成：

\[
\begin{pmatrix}
u\\
v
\end{pmatrix}
=
\begin{pmatrix}
\frac12 I & \frac12 I\\
\frac12 D^{-1} & -\frac12 D^{-1}
\end{pmatrix}
\begin{pmatrix}
z^+\\
z^-
\end{pmatrix}.
\]

这就是最清楚的 slot 层关系。

---

# 15. slot-preserving split 的矩阵形式

现在假设你的目标不是 \((u,v)\)，而是保留原来的两个 slot 分支：

\[
z^+=P(\alpha_j),\quad z^-=P(-\alpha_j).
\]

我们想构造两个小环多项式：

\[
Q^+(Y),\quad Q^-(Y),
\]

使得：

\[
Q^+(\theta_j)=z_j^+,\quad Q^-(\theta_j)=z_j^-.
\]

也就是：

\[
Q^+\text{ 的 slot vector}=z^+,\quad Q^-\text{ 的 slot vector}=z^-.
\]

它们的系数向量是：

\[
q^+=F_n^{-1}z^+,\quad q^-=F_n^{-1}z^-.
\]

但是：

\[
z^+=u+Dv,\quad z^-=u-Dv.
\]

又因为：

\[
u=F_n E p,\quad v=F_n O p.
\]

所以：

\[
q^+ = F_n^{-1}(u+Dv) = F_n^{-1}(F_n E p+D F_n O p) = E p+F_n^{-1}D F_n O p.
\]

同理：

\[
q^- = E p-F_n^{-1}D F_n O p.
\]

令：

\[
H=F_n^{-1}D F_n.
\]

那么：

\[
\boxed{q^+=E p+H O p}
\]
\[
\boxed{q^-=E p-H O p}
\]

这就是 slot-preserving split 的系数矩阵形式。

---

# 16. 矩阵 \(H\) 是什么？

\[
H=F_n^{-1}D F_n.
\]

这里：

\[
D=\operatorname{diag}(\alpha_0,\alpha_1,\dots,\alpha_{n-1}).
\]

在 slot/evaluation 坐标中，乘以 \(h(Y)\) 就是逐 slot 乘以：

\[
h(\theta_j)=\alpha_j.
\]

也就是说，在 slot 坐标里：

\[
B(\theta_j) \longmapsto \alpha_j B(\theta_j)
\]

是：

\[
v\longmapsto Dv.
\]

回到系数坐标，就是：

\[
b\longmapsto F_n^{-1}D F_n b.
\]

所以：

\[
H
\]

就是“乘以 \(h(Y)\)”这个操作在系数坐标里的矩阵。

因此：

\[
q^+=E p+H O p
\]

对应多项式：

\[
Q^+(Y)=A(Y)+h(Y)B(Y).
\]

\[
q^-=E p-H O p
\]

对应：

\[
Q^-(Y)=A(Y)-h(Y)B(Y).
\]

---

# 17. 任意 slot selection 的矩阵形式

现在更一般一点。

假设大环 slot vector 是：

\[
z=F_N p.
\]

你想要小环 slot vector：

\[
y=S z.
\]

这里 \(S\) 是你想要的 slot 选择或线性变换矩阵。

例如：

* 取前一半 slots；
* 取偶数 slots；
* 两两求和；
* 任意线性组合。

小环多项式系数 \(q\) 满足：

\[
F_n q=y.
\]

所以：

\[
q=F_n^{-1}y.
\]

代入：

\[
y=S z=S F_N p.
\]

因此：

\[
\boxed{q=F_n^{-1}S F_N p.}
\]

这就是最一般的公式。

它说明：

```text
slot selection 回到 polynomial coefficient 层面，
通常是一个稠密线性变换。
```

不是简单的：

\[
p\mapsto Ep
\]

或者：

\[
p\mapsto Op.
\]

---

# 18. 三种拆分的最终对比

## 18.1 coefficient split

目标：

\[
P(X)\mapsto A(Y),B(Y).
\]

公式：

\[
P(X)=A(X^2)+XB(X^2).
\]

矩阵：

\[
A\text{ coeff}=Ep,\quad B\text{ coeff}=Op.
\]

优点：非常简单，ring switching 可以高效实现。

缺点：在 slot 语义下不是保留原 slots，而是 pairwise sum/difference。

---

## 18.2 pairwise slot-preserving split

目标：

\[
Q^+(\theta_j)=P(\alpha_j),\quad Q^-(\theta_j)=P(-\alpha_j).
\]

公式：

\[
Q^+(Y)=A(Y)+h(Y)B(Y),\quad Q^-(Y)=A(Y)-h(Y)B(Y).
\]

矩阵：

\[
q^+=Ep+HOp,\quad q^-=Ep-HOp,
\]

其中：

\[
H=F_n^{-1}D F_n.
\]

这比 coefficient split 多了一个 \(H\)，也就是乘以平方根选择多项式 \(h(Y)\)。

---

## 18.3 arbitrary slot selection

目标：

\[
y=S z.
\]

其中：

\[
z=F_N p.
\]

公式：

\[
q=F_n^{-1}S F_N p.
\]

这是最一般的 slot-level 线性变换，通常是稠密的。

---

# 19. 一个最小例子：\(N=4,n=2\)

设：

\[
P(X)=a_0+a_1X+a_2X^2+a_3X^3.
\]

则：

\[
A(Y)=a_0+a_2Y,\quad B(Y)=a_1+a_3Y.
\]

大环 roots 成对：

\[
\alpha_0,\ -\alpha_0,\quad \alpha_1,\ -\alpha_1.
\]

小环 roots：

\[
\theta_0=\alpha_0^2,\quad \theta_1=\alpha_1^2.
\]

一对 slots：

\[
P(\alpha_0)=A(\theta_0)+\alpha_0B(\theta_0),\quad P(-\alpha_0)=A(\theta_0)-\alpha_0B(\theta_0).
\]

所以：

\[
A(\theta_0)=\frac{P(\alpha_0)+P(-\alpha_0)}{2},\quad B(\theta_0)=\frac{P(\alpha_0)-P(-\alpha_0)}{2\alpha_0}.
\]

同理对 \(\alpha_1\)。

如果你想保留原 slots，则需要：

\[
Q^+(\theta_0)=P(\alpha_0),\quad Q^+(\theta_1)=P(\alpha_1).
\]

那么：

\[
Q^+(Y)=A(Y)+h(Y)B(Y),
\]

其中：

\[
h(\theta_0)=\alpha_0,\quad h(\theta_1)=\alpha_1.
\]

这就是插值出来的多项式。

---

# 20. 这对你的研究意味着什么？

如果你做的是：

\[
P(X)\mapsto A(Y),B(Y),
\]

那是 coefficient split。

它高效，但 slot 语义是：

\[
(P(\alpha),P(-\alpha)) \mapsto \left(
\frac{P(\alpha)+P(-\alpha)}{2},\ \frac{P(\alpha)-P(-\alpha)}{2\alpha}
\right).
\]

如果你想要：

\[
(P(\alpha),P(-\alpha)) \mapsto (P(\alpha),P(-\alpha))
\]

但放到两个小环密文里，那么需要：

\[
Q^+=A+hB,\quad Q^-=A-hB.
\]

这一步比 coefficient split 更复杂，因为要乘 \(h(Y)\)。

如果你想要任意选择 slots，那更一般地是：

\[
q=F_n^{-1}S F_N p.
\]

这通常就是一个 homomorphic linear transform 问题，而不是简单 ring split。

---

一句话总结：

\[
\boxed{\text{coefficient split 是 }p\mapsto(Ep,Op).}
\]

\[
\boxed{\text{pairwise slot-preserving split 是 }p\mapsto(Ep+HOp,\ Ep-HOp).}
\]

\[
\boxed{\text{任意 slot split 是 }p\mapsto F_n^{-1}S F_Np.}
\]

其中：

\[
H=F_n^{-1}D F_n.
\]

这就是完整的矩阵形式。


以上内容完全保留了您提供的推导结构和所有数学表达式，并确保了正确的 Markdown 渲染效果。