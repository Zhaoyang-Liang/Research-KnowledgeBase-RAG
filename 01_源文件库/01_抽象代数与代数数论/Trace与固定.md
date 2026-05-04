trace 就是“把一个数的所有共轭版本加起来”。
一个数如果是某个有理系数多项式的根，那么这个多项式的其他根，就是它的共轭。

我们证明这个命题：

\[
\operatorname{Tr}_{K/\mathbb Q}(a)=\sum_{\tau}\tau(a)
\]

为什么一定落在 $\mathbb Q$ 里；以及更一般地，

\[
\operatorname{Tr}_{K/K_0}(a)=\sum_{\tau\text{ fixes }K_0}\tau(a)
\]

为什么一定落在 $K_0$ 里。

GHPS 里正是这样定义 trace：对于 $K=\mathbb Q(\zeta_m)$，它有自同构 $\tau_i(\zeta_m)=\zeta_m^i$，并且因为 $K/\mathbb Q$ 是 Galois，trace 可以定义为这些自同构结果的和；对于 $K/K_0$，trace 是所有固定 $K_0$ 的自同构之和。

---

# 1. 先证明 $K/\mathbb Q$ 的情况

设：

\[
K=\mathbb Q(\zeta_m)
\]

它的自同构是：

\[
\tau_i:K\to K,\qquad \tau_i(\zeta_m)=\zeta_m^i
\]

其中：

\[
i\in \mathbb Z_m^*
\]

也就是 $i$ 和 $m$ 互素。

对任意：

\[
a\in K
\]

定义：

\[
T=\sum_{i\in\mathbb Z_m^*}\tau_i(a)
\]

我们要证明：

\[
T\in \mathbb Q
\]

---

## 第一步：证明 $T$ 被所有自同构固定

任取一个自同构：

\[
\tau_j
\]

其中：

\[
j\in\mathbb Z_m^*
\]

看它作用在 $T$ 上：

\[
\tau_j(T)
=
\tau_j\left(\sum_{i\in\mathbb Z_m^*}\tau_i(a)\right)
\]

因为 $\tau_j$ 保持加法，所以：

\[
\tau_j(T)
=
\sum_{i\in\mathbb Z_m^*}\tau_j(\tau_i(a))
\]

自同构复合后还是自同构，而且：

\[
\tau_j\circ\tau_i=\tau_{ji}
\]

所以：

\[
\tau_j(T)
=
\sum_{i\in\mathbb Z_m^*}\tau_{ji}(a)
\]

现在关键来了：当 $i$ 遍历所有 $\mathbb Z_m^*$ 时，$ji$ 也遍历所有 $\mathbb Z_m^*$。

原因是 $j\in\mathbb Z_m^*$，所以 $j$ 模 $m$ 可逆。乘以 $j$ 只是把集合 $\mathbb Z_m^*$ 重新排列了一遍。

因此：

\[
\sum_{i\in\mathbb Z_m^*}\tau_{ji}(a)
=
\sum_{i\in\mathbb Z_m^*}\tau_i(a)
= T
\]

所以：

\[
\tau_j(T)=T
\]

这对所有 $j\in\mathbb Z_m^*$ 都成立。

也就是说：

\[
T
\]

被所有自同构固定。

---

## 第二步：为什么“被所有自同构固定”就说明 $T\in\mathbb Q$？

因为在：

\[
K=\mathbb Q(\zeta_m)
\]

里，自同构的作用本质上就是把：

\[
\zeta_m
\]

替换成：

\[
\zeta_m^i,\qquad i\in\mathbb Z_m^*
\]

如果一个元素 $T\in K$ 在所有这些替换下都不变，那么它不依赖于选择哪个本原根。这样的元素只能来自基础域：

\[
\mathbb Q
\]

更正式地说，$K/\mathbb Q$ 是Galois扩张；Galois扩张里，被全部 $\mathbb Q$-自同构固定的元素，正好就是 $\mathbb Q$。GHPS 也明确使用了 $K/\mathbb Q$ 是Galois这一点。

所以：

\[
T\in\mathbb Q
\]

也就是：

\[
\operatorname{Tr}_{K/\mathbb Q}(a)\in\mathbb Q
\]

---

# 2. 用一个具体例子看这个证明

取：

\[
K=\mathbb Q(i)
\]

有两个自同构：

\[
\tau_1(i)=i
\]
\[
\tau_2(i)=-i
\]

设：

\[
a=3+2i
\]

trace 是：

\[
T=(3+2i)+(3-2i)=6
\]

现在看 $T$ 是否被所有自同构固定。

恒等自同构当然固定：

\[
\tau_1(6)=6
\]

复共轭自同构也固定：

\[
\tau_2(6)=6
\]

所以 $T$ 在两个替换下都不变。

而 $6$ 确实在：

\[
\mathbb Q
\]

里。

---

# 3. 再证明 $K/K_0$ 的情况

现在设：

\[
K_0\subset K
\]

在 GHPS 里：

\[
K=\mathbb Q(\zeta_m)
\]
\[
K_0=\mathbb Q(\zeta_{m_0})
\]

其中：

\[
m_0\mid m
\]

论文说明，$K/K_0$ 也是 Galois，并且固定 $K_0$ 的自同构正好是那些满足：

\[
i\equiv 1\pmod{m_0}
\]

的 $\tau_i$。于是定义：

\[
\operatorname{Tr}_{K/K_0}(a)
=
\sum_{i \equiv 1\pmod{m_0}} \tau_i(a)
\]

设：

\[
T=\sum_{\tau\text{ fixes }K_0}\tau(a)
\]

我们要证明：

\[
T\in K_0
\]

证明和刚才一样。

任取一个固定 $K_0$ 的自同构 $\rho$，则：

\[
\rho(T) = \rho\left(\sum_{\tau\text{ fixes }K_0}\tau(a)\right) = \sum_{\tau\text{ fixes }K_0}\rho(\tau(a)) = \sum_{\tau\text{ fixes }K_0}(\rho\circ\tau)(a)
\]

当 $\tau$ 遍历所有固定 $K_0$ 的自同构时，$\rho\circ\tau$ 也遍历同一批自同构，只是顺序变了。

所以：

\[
\rho(T)=T
\]

也就是说：

\[
T
\]

被所有固定 $K_0$ 的自同构固定。

在 Galois 扩张 $K/K_0$ 里，被这些自同构全部固定的元素，正好就是：

\[
K_0
\]

所以：

\[
T\in K_0
\]

也就是：

\[
\operatorname{Tr}_{K/K_0}(a)\in K_0
\]

---

# 4. 一个很具体的 $K/K_0$ 例子

取：

\[
K=\mathbb Q(\zeta_8)
\]
\[
K_0=\mathbb Q(i)
\]

因为：

\[
i=\zeta_8^2
\]
所以：

\[
K_0=\mathbb Q(\zeta_8^2)
\]

在 $K$ 里，考虑替换：

\[
\zeta_8\mapsto -\zeta_8
\]

这个替换固定 $i$，因为：

\[
(-\zeta_8)^2=\zeta_8^2=i
\]

所以相对 trace 是：

\[
\operatorname{Tr}_{K/K_0}(a)=a(\zeta_8)+a(-\zeta_8)
\]

取：

\[
a=3+2\zeta_8+5\zeta_8^2
\]

那么：

\[
a(-\zeta_8)=3-2\zeta_8+5\zeta_8^2
\]

相加：

\[
\operatorname{Tr}_{K/K_0}(a)
=
(3+2\zeta_8+5\zeta_8^2) + (3-2\zeta_8+5\zeta_8^2) = 6+10\zeta_8^2
\]

因为：

\[
\zeta_8^2=i
\]

所以：

\[
6+10\zeta_8^2=6+10i
\]

它属于：

\[
K_0=\mathbb Q(i)
\]

你可以看到：含 $\zeta_8$ 的部分被抵消掉了，只剩下含 $\zeta_8^2=i$ 的部分。

---

# 5. 这个有什么用？

trace 的用处可以分成两个层次。

---

## 用处一：把“大域元素”变成“小域元素”

这是最基本的作用。

\[
\operatorname{Tr}_{K/K_0}:K\to K_0
\]

它把大域 $K$ 里的元素变成小域 $K_0$ 里的元素。

在刚才例子里：

\[
3+2\zeta_8+5\zeta_8^2\in K=\mathbb Q(\zeta_8)
\]

trace 后：

\[
6+10i\in K_0=\mathbb Q(i)
\]

这就是“降到小域”。

在 GHPS 论文里，这正是 field switching 的核心：最后一步是对 ciphertext 中的 $K$-元素取 trace，从而得到 $K_0$ 上的 ciphertext。论文贡献部分也明确说：他们最后对密文中的 $K$-元素取 trace，得到子域 $K_0$ 上的输出密文。

---

## 用处二：可以做“投影”或“抽取”

在二次例子里：

\[
F(X)=a_0+a_1X+a_2X^2+a_3X^3
\]

如果小域是由：

\[
Y=X^2
\]

生成，那么固定小域的替换是：

\[
X\mapsto -X
\]

trace 是：

\[
F(X)+F(-X)
\]

计算：

\[
F(-X)=a_0-a_1X+a_2X^2-a_3X^3
\]

所以：

\[
F(X)+F(-X)=2a_0+2a_2X^2
\]

奇次项消失，只剩偶次项。

因此：

\[
\frac{F(X)+F(-X)}{2} = a_0+a_2X^2
\]

这就是把 $F$ 的偶部抽出来。

如果做：

\[
\frac{F(X)-F(-X)}{2} = a_1X+a_3X^3
\]

则抽出奇部乘 $X$ 的形式。

这正是你之前两个文档里偶奇拆分公式的来源。

---

## 用处三：在 GHPS 里，trace 是 field switching 的核心操作

GHPS 的 field switching 目标是：

\[
\text{big-field ciphertext over }K
\longrightarrow
\text{small-field ciphertext over }K_0
\]

论文的三步是：

1. 先把密文切换到小域私钥；
2. 乘一个合适的元素 $r$；
3. 对密文分量取 trace，得到小域密文。

第 3 步就是：

\[
c_i'=\operatorname{Tr}_{K/K_0}(r c_i)
\]

这个式子为什么合理？因为：

\[
r c_i\in K
\]

而：

\[
\operatorname{Tr}_{K/K_0}(r c_i)\in K_0
\]

所以密文分量从 $K$ 里的东西变成了 $K_0$ 里的东西。

---

## 用处四：乘 $r$ 再 trace 可以实现不同的线性选择

只取 trace：

\[
\operatorname{Tr}_{K/K_0}(a)
\]

是一种固定操作。

但如果先乘不同的 $r$：

\[
\operatorname{Tr}_{K/K_0}(r a)
\]

就可以得到不同的结果。

GHPS 后面正是利用这一点：他们研究哪些函数可以写成

\[
L(a)=\operatorname{Tr}_{K/K_0}(r a)
\]

并说明通过合适的 $r$，可以实现某些 plaintext slots 的选择或线性函数。论文在 2.1.3 里说，field switching 是围绕对 big-field ciphertext 应用 trace 展开的，并且要刻画哪些函数可以表示成 $L(a)=\operatorname{Tr}_{K/K_0}(r\cdot a)$。

---

# 6. 最后总结

你要证明的核心就是这两步：

第一：

\[
T=\sum_{\tau}\tau(a)
\]

对任意自同构 $\rho$，都有：

\[
\rho(T)=T
\]

因为 $\rho$ 只是把求和里的项重新排列了一遍。

第二：

如果一个元素被所有相关自同构固定，那么它属于被这些自同构共同固定的域。

所以：

\[
\sum_{\tau\in\operatorname{Aut}(K/\mathbb Q)}\tau(a)\in\mathbb Q
\]

以及：

\[
\sum_{\tau\in\operatorname{Aut}(K/K_0)}\tau(a)\in K_0
\]

它的用处是：**trace 是一个从大域到小域的标准降维操作**。在 GHPS 里，它就是把 big-field ciphertext 变成 small-field ciphertext 的最后一步；在你关心的 $X\mapsto -X$ 情况里，它表现为偶奇分解：

\[
F(X)+F(-X)
\]

会消掉奇次项，只留下小环里的偶次部分。


===

这个条件：

\[
i=1\pmod{m_0}
\]

意思是：

\[
i
\]

除以

\[
m_0
\]

以后余数是 1。

例如：

\[
i=1\pmod 4
\]

表示：

\[
i=1,5,9,13,\ldots
\]

在论文里，作者写的是：$K/K_0$ 中那些 **固定 $K_0$** 的自同构，正好是满足 $i=1\pmod{m_0}$ 的 $\tau_i$。论文也给了原因：因为 $\tau_i(\zeta_{m_0})=\zeta_{m_0}^{i\bmod m_0}$。

---

# 1. 什么叫“固定”？

一个自同构：

\[
\tau:K\to K
\]

说它 **固定 $K_0$**，意思是：

\[
\forall x\in K_0,\quad \tau(x)=x
\]

也就是 $K_0$ 里面的每个元素都不变。

这叫 **fix $K_0$ pointwise**。

注意，不是说“把 $K_0$ 映射回 $K_0$”就够了，而是要每个元素都原封不动。

例如：

\[
K_0=\mathbb Q(i)
\]

复共轭：

\[
i\mapsto -i
\]

虽然会把 $\mathbb Q(i)$ 映射到 $\mathbb Q(i)$，但它没有固定 $\mathbb Q(i)$，因为：

\[
i\neq -i
\]

所以它不叫 fix $K_0$ pointwise。

---

# 2. 为什么只需要看 $\zeta_{m_0}$？

因为：

\[
K_0=\mathbb Q(\zeta_{m_0})
\]

意思是 $K_0$ 里所有元素都是由有理数和 $\zeta_{m_0}$ 组合出来的，比如：

\[
a_0+a_1\zeta_{m_0}+a_2\zeta_{m_0}^2+\cdots
\]

自同构本来就固定所有有理数，所以要让整个 $K_0$ 都不变，只需要让生成元不变：

\[
\tau(\zeta_{m_0})=\zeta_{m_0}
\]

如果 $\zeta_{m_0}$ 不变，那么由它构造出来的所有元素都会不变。

---

# 3. 为什么条件是 $i=1\pmod{m_0}$？

在 GHPS 里：

\[
K=\mathbb Q(\zeta_m)
\]
\[
K_0=\mathbb Q(\zeta_{m_0})
\]

并且：

\[
m_0\mid m
\]

所以可以把 $\zeta_{m_0}$ 看成：

\[
\zeta_{m_0}=\zeta_m^{m/m_0}
\]

这是因为：

\[
\left(\zeta_m^{m/m_0}\right)^{m_0}
=
\zeta_m^{m} = 1
\]

它确实是一个 $m_0$ 次单位根。

现在 $K$ 的自同构是：

\[
\tau_i(\zeta_m)=\zeta_m^i
\]

那么它作用在 $\zeta_{m_0}$ 上：

\[
\tau_i(\zeta_{m_0}) = \tau_i\left(\zeta_m^{m/m_0}\right) = (\tau_i(\zeta_m))^{m/m_0} = (\zeta_m^i)^{m/m_0} = \zeta_m^{i(m/m_0)}
\]
而
\[
\zeta_m^{m/m_0}=\zeta_{m_0}
\]
所以
\[
    \zeta_m^{i(m/m_0)} = (\zeta_m^{m/m_0})^i = \zeta_{m_0}^i
\]

因此：

\[
\tau_i(\zeta_{m_0})=\zeta_{m_0}^i
\]

要让 $\tau_i$ 固定 $K_0$，就必须有：

\[
\tau_i(\zeta_{m_0}) = \zeta_{m_0}
\]

也就是：

\[
\zeta_{m_0}^i = \zeta_{m_0}
\]

两边除以 $\zeta_{m_0}$：

\[
\zeta_{m_0}^{i-1}=1
\]

因为 $\zeta_{m_0}$ 的阶是 $m_0$，所以：

\[
\zeta_{m_0}^{i-1}=1
\]

当且仅当：

\[
m_0\mid i-1
\]

也就是：

\[
i=1\pmod{m_0}
\]

这就是这个条件的来源。

它不是随便来的，而是因为：

\[
\zeta_{m_0}^i
\]

要等于：

\[
\zeta_{m_0}
\]

所以指数 $i$ 必须和 $1$ 在模 $m_0$ 下相同。

---

# 4. 具体例子：$m=8,\ m_0=4$

设：

\[
K=\mathbb Q(\zeta_8)
\]
\[
K_0=\mathbb Q(\zeta_4)=\mathbb Q(i)
\]

因为：

\[
\zeta_4=\zeta_8^2
\]

$K$ 的自同构由：

\[
i\in \mathbb Z_8^*=\{1,3,5,7\}
\]

决定：

\[
\tau_i(\zeta_8)=\zeta_8^i
\]

现在看它们是否固定 $K_0$，也就是是否固定：

\[
\zeta_4=\zeta_8^2=i
\]

---

## $\tau_1$

\[
\tau_1(\zeta_4)=\zeta_4^1=\zeta_4
\]

固定。

因为：

\[
1=1\pmod 4
\]

---

## $\tau_3$

\[
\tau_3(\zeta_4)=\zeta_4^3
\]

而：

\[
\zeta_4=i
\]
\[
\zeta_4^3=i^3=-i
\]

所以：

\[
\tau_3(\zeta_4)=-i\neq i
\]

不固定。

因为：

\[
3\neq 1\pmod 4
\]

---

## $\tau_5$

\[
\tau_5(\zeta_4)=\zeta_4^5
\]

因为：

\[
5=1\pmod 4
\]

所以：

\[
\zeta_4^5=\zeta_4
\]

固定。

---

## $\tau_7$

\[
7=3\pmod 4
\]

所以：

\[
\tau_7(\zeta_4)=\zeta_4^7=\zeta_4^3=-i
\]

不固定。

---

因此固定 $K_0=\mathbb Q(i)$ 的只有：

\[
\tau_1,\tau_5
\]

它们对应：

\[
i=1,5
\]

正好满足：

\[
i=1\pmod 4
\]

---

# 5. 为什么相对 trace 要对这些 $i$ 求和？

因为相对 trace：

\[
\operatorname{Tr}_{K/K_0}:K\to K_0
\]

要把 $K$ 里的元素变成 $K_0$ 里的元素。

所以它只对那些 **不动 $K_0$** 的自同构求和：

\[
\operatorname{Tr}_{K/K_0}(a)
=
\sum_{i\equiv 1\pmod{m_0}}
\tau_i(a)
\]

这样求和结果会落在 $K_0$ 里。论文就是这样定义 intermediate trace 的。

在刚才例子里：

\[
\operatorname{Tr}_{\mathbb Q(\zeta_8)/\mathbb Q(i)}(a)
=
\tau_1(a)+\tau_5(a)
\]

因为只有 $\tau_1,\tau_5$ 固定 $\mathbb Q(i)$。

---

# 6. 一句话总结

\[
i=1\pmod{m_0}
\]

这个条件的意思是：

> 自同构 $\tau_i$ 把 $\zeta_m$ 送到 $\zeta_m^i$，但它必须让子域生成元 $\zeta_{m_0}$ 不变。

由于：

\[
\zeta_{m_0}=\zeta_m^{m/m_0}
\]

所以：

\[
\tau_i(\zeta_{m_0})=\zeta_{m_0}^i
\]

要让它等于原来的：

\[
\zeta_{m_0}
\]

就必须：

\[
\zeta_{m_0}^i=\zeta_{m_0}
\]

等价于：

\[
i=1\pmod{m_0}
\]

所以这个条件本质上是在说：

> 这个自同构不能改变小域 $K_0$ 的生成元。
