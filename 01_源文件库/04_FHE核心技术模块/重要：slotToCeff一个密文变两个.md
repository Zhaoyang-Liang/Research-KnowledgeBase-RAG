

可以。我把刚才那段**完整重写一遍**，这次做两件事：

1. **所有对象都写全**，比如明确写成 $p(Y)$、$m(X)$、$\operatorname{Ecd}(\cdot)$。
2. **完全不用 $\operatorname{Enc}^*(\cdot)$**。以后凡是原来写 $\operatorname{Enc}^*(\mathbf v)$ 的地方，我都改写成 $\operatorname{Enc}(\operatorname{Ecd}(\mathbf v))$。也就是说，它实际上加密的是“向量 $\mathbf v$ 经过 CKKS 编码后得到的那个多项式”。

---

# 1. 先把所有符号写清楚

设
$$
p(Y)=\sum_{i=0}^{2n-1} p_i Y^i \in \mathbb{R}[Y]/(Y^{2n}+1).
$$

把它拆成前后两半系数向量：
$$
\mathbf p_0=(p_0,p_1,\dots,p_{n-1})\in\mathbb{R}^n,\\
\mathbf p_1=(p_n,p_{n+1},\dots,p_{2n-1})\in\mathbb{R}^n.
$$

再定义复向量
$$
\mathbf w=\mathbf p_0+\mathbf I \mathbf p_1=(p_0+\mathbf I p_n,\, p_1+\mathbf I p_{n+1},\, \dots,\, p_{n-1}+\mathbf I p_{2n-1})\in\mathbb{C}^n.
$$

然后定义
$$
\mathbf z = U_n \mathbf w \in \mathbb{C}^n.
$$

根据论文里的关系，
$$
\mathbf z = \tau_n\bigl(p(Y)\bigr),
$$
也就是
$$
\tau_n^{-1}(\mathbf z)=p(Y).
$$

---

# 2. $\operatorname{Ecd}(\mathbf z)$ 到底是什么

CKKS 对向量 $\mathbf z\in\mathbb{C}^n$ 的编码是
$$
\operatorname{Ecd}(\mathbf z)
= \left\lfloor \Delta\cdot \tau_n^{-1}(\mathbf z)\right\rfloor.
$$

因为这里
$$
\tau_n^{-1}(\mathbf z)=p(Y),
$$
所以
$$
\operatorname{Ecd}(\mathbf z)
= \left\lfloor \Delta\cdot p(Y)\right\rfloor.
$$

但是 CKKS 的实际密文环是
$$
\mathbb{Z}[X]/(X^N+1),
$$
因此严格地说，需要把 $Y$ 替换成 $X^{N/(2n)}$。所以更完整地写是：
$$
\operatorname{Ecd}(\mathbf z)
= \left\lfloor \Delta\cdot p\bigl(X^{N/(2n)}\bigr)\right\rfloor
\in \mathbb{Z}[X]/(X^N+1).
$$

也就是说，如果
$$
p(Y)=\sum_{i=0}^{2n-1}p_iY^i,
$$
那么
$$
p\bigl(X^{N/(2n)}\bigr) = \sum_{i=0}^{2n-1} p_i X^{iN/(2n)}.
$$

所以你之前看到我写
$$
m\approx \Delta\cdot p
$$
确实太省略了。更准确的写法应该是：

$$
m(X)\approx \left\lfloor \Delta\cdot p\bigl(X^{N/(2n)}\bigr)\right\rfloor.
$$

或者如果只强调“编码语义”，就写

$$
m(X)=\operatorname{Ecd}(\mathbf z)
= \left\lfloor \Delta\cdot p\bigl(X^{N/(2n)}\bigr)\right\rfloor.
$$

---

# 3. 现在重写“反过来怎么得到两个关于 $p_i$ 的密文”

你问的是：

> 已知一个密文 $\mathbf{ct}$，它加密的是
> $$
> m(X)=\left\lfloor \Delta\cdot p\bigl(X^{N/(2n)}\bigr)\right\rfloor,
> $$
> 怎么同态地得到两个新密文，分别对应
> $$
> (p_0,\dots,p_{n-1})
> \quad\text{和}\quad
> (p_n,\dots,p_{2n-1})?
> $$

下面是完整写法。

---

## 第一步：输入密文到底加密什么

设输入密文为
$$
\mathbf{ct}=\operatorname{Enc}(m(X)),
$$
其中
$$
m(X)=\operatorname{Ecd}(\mathbf z)
= \left\lfloor \Delta\cdot p\bigl(X^{N/(2n)}\bigr)\right\rfloor,
$$
并且
$$
\mathbf z = U_n(\mathbf p_0+\mathbf I\mathbf p_1)=U_n\mathbf w.
$$

所以同一个输入密文 $\mathbf{ct}$ 也可以理解成：
$$
\mathbf{ct}
= \operatorname{Enc}\bigl(\operatorname{Ecd}(\mathbf z)\bigr).
$$

这就是原来写成 $\operatorname{Enc}^*(\mathbf z)$ 的地方；现在我已经把星号完全展开了。

---

## 第二步：同态乘上 $U_n^{-1}$

因为
$$
\mathbf z = U_n\mathbf w,
$$
所以
$$
\mathbf w = U_n^{-1}\mathbf z.
$$

于是我们对密文做矩阵-密文乘法：

$$
\mathbf{ct}' = \operatorname{Ecd}(U_n^{-1}) \times \mathbf{ct}.
$$

这个输出密文加密的是

$$
\operatorname{Ecd}(U_n^{-1}\mathbf z)
= \operatorname{Ecd}(\mathbf w)
= \operatorname{Ecd}(\mathbf p_0+\mathbf I\mathbf p_1).
$$

因此

$$
\mathbf{ct}'
= \operatorname{Enc}\bigl(\operatorname{Ecd}(\mathbf p_0+\mathbf I\mathbf p_1)\bigr).
$$

这一步之后，密文 $\mathbf{ct}'$ 在槽语义上对应的向量就是

$$
\mathbf p_0+\mathbf I\mathbf p_1
= (p_0+\mathbf I p_n,\, p_1+\mathbf I p_{n+1},\, \dots,\, p_{n-1}+\mathbf I p_{2n-1}).
$$

---

## 第三步：做一次共轭

对 $\mathbf{ct}'$ 做同态共轭，得到

$$
\operatorname{Conj}(\mathbf{ct}') = \operatorname{Enc}\bigl(\operatorname{Ecd}(\mathbf p_0-\mathbf I\mathbf p_1)\bigr).
$$

因为
$$
\overline{\mathbf p_0+\mathbf I\mathbf p_1} = \mathbf p_0-\mathbf I\mathbf p_1.
$$

---

## 第四步：取实部，得到第一个密文

定义
$$
\mathbf{ct}_0 = \frac12\bigl(\mathbf{ct}'+\operatorname{Conj}(\mathbf{ct}')\bigr).
$$

那么它加密的是

$$
\operatorname{Ecd}\left(
\frac12\bigl((\mathbf p_0+\mathbf I\mathbf p_1)+(\mathbf p_0-\mathbf I\mathbf p_1)\bigr)
\right)
= \operatorname{Ecd}(\mathbf p_0).
$$

所以

$$
\mathbf{ct}_0 = \operatorname{Enc}\bigl(\operatorname{Ecd}(\mathbf p_0)\bigr).
$$

也就是说，$\mathbf{ct}_0$ 对应的槽向量是

$$
\mathbf p_0=(p_0,p_1,\dots,p_{n-1}).
$$

---

## 第五步：取虚部，得到第二个密文

定义
$$
\mathbf{ct}_1 = \frac{1}{2\mathbf I}\bigl(\mathbf{ct}'-\operatorname{Conj}(\mathbf{ct}')\bigr).
$$

那么它加密的是

$$
\operatorname{Ecd}\left(
\frac{1}{2\mathbf I}
\bigl((\mathbf p_0+\mathbf I\mathbf p_1)-(\mathbf p_0-\mathbf I\mathbf p_1)\bigr)
\right)
= \operatorname{Ecd}(\mathbf p_1).
$$

所以

$$
\mathbf{ct}_1 = \operatorname{Enc}\bigl(\operatorname{Ecd}(\mathbf p_1)\bigr).
$$

也就是说，$\mathbf{ct}_1$ 对应的槽向量是

$$
\mathbf p_1=(p_n,p_{n+1},\dots,p_{2n-1}).
$$

---

# 4. 整条链不省略地写一遍

从输入到输出，完整地就是：

$$
\mathbf{ct}
=
\operatorname{Enc}\left(
\left\lfloor \Delta\cdot p\bigl(X^{N/(2n)}\bigr)\right\rfloor
\right)
= \operatorname{Enc}\bigl(\operatorname{Ecd}(\mathbf z)\bigr),
$$
其中
$$
\mathbf z=U_n(\mathbf p_0+\mathbf I\mathbf p_1).
$$

然后做

$$
\mathbf{ct}'
= \operatorname{Ecd}(U_n^{-1})\times \mathbf{ct}
= \operatorname{Enc}\bigl(\operatorname{Ecd}(\mathbf p_0+\mathbf I\mathbf p_1)\bigr).
$$

再做

$$
\mathbf{ct}_0 = \frac12\bigl(\mathbf{ct}'+\operatorname{Conj}(\mathbf{ct}')\bigr)
= \operatorname{Enc}\bigl(\operatorname{Ecd}(\mathbf p_0)\bigr),
$$

$$
\mathbf{ct}_1 = \frac{1}{2\mathbf I}\bigl(\mathbf{ct}'-\operatorname{Conj}(\mathbf{ct}')\bigr)
= \operatorname{Enc}\bigl(\operatorname{Ecd}(\mathbf p_1)\bigr).
$$

---

# 5. 最后一句话概括它到底在做什么

它做的不是“把多项式系数直接从密文里掏出来”。
它做的是：

1. 输入密文原本加密的是
   $$
   \left\lfloor \Delta\cdot p\bigl(X^{N/(2n)}\bigr)\right\rfloor;
   $$
2. 利用 CKKS 的解码语义，把这个密文看成在编码某个槽向量 $\mathbf z=\tau_n(p(Y))$；
3. 同态乘上 $U_n^{-1}$，把槽向量变成
   $$
   \mathbf p_0+\mathbf I\mathbf p_1;
   $$
4. 再用共轭把它拆成两个密文：
   $$
   \operatorname{Enc}(\operatorname{Ecd}(\mathbf p_0))
   \quad\text{和}\quad
   \operatorname{Enc}(\operatorname{Ecd}(\mathbf p_1)).
   $$

---

如果你愿意，我下一条可以继续把**正向的 SlotToCoeff** 也用同样的风格重写一遍：同样不使用 $\operatorname{Enc}^*(\cdot)$，把每个对象都写全。
