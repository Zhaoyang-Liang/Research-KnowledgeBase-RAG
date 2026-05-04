# 由于对NTRU已经非常熟悉了，所以这里主要是一些代数上之前理解有误差的东西

# $(\mathbb Z[\zeta_n])$ 是 $(\mathbb Z)$ 的环扩张

对，这里最该先分清的是：

## ( \mathbb Z[\zeta_n] ) 不是域，通常是一个**环**

## ( \mathbb Q(\alpha) ) 是一个**域**

## ( \mathbb Z[\alpha]\subset \mathbb Q(\alpha) ) 的意思是：前者是后者里面的一个子环

我用最经典的例子先讲，你会马上明白。

---

# 1. 先看最熟悉的例子：(i=\sqrt{-1})

你肯定见过：

[
\mathbb Z[i]={a+bi: a,b\in \mathbb Z}
]

这叫 **Gaussian integers**，高斯整数。它是一个**环**：

* 可以加减乘
* 也有 1
* 但不是每个非零元素都可逆

比如 (2\in \mathbb Z[i])，但 (1/2\notin \mathbb Z[i])，所以它不是域。

再看

[
\mathbb Q(i)={a+bi: a,b\in \mathbb Q}
]

这就是一个**域**。因为现在分母允许是有理数了，所以非零元素都能求逆。

于是有

[
\mathbb Z[i]\subset \mathbb Q(i).
]

意思就是：

* 左边是“整数系数”的那批元素
* 右边是“有理数系数”的那批元素
* 左边包含在右边里面

---

# 2. 完全照搬到 ( \alpha ) 上

现在把 (i) 换成某个代数元素 (\alpha)。

## ( \mathbb Z[\alpha] )

定义为：

[
\mathbb Z[\alpha]
=================

{f(\alpha): f(x)\in \mathbb Z[x]}.
]

也就是：

> 用**整数系数多项式**去代入 (\alpha) 所得到的所有元素。

比如如果 (\alpha) 满足某个代数方程，那么这里面元素都长得像

[
a_0+a_1\alpha+\cdots+a_k\alpha^k,\qquad a_i\in\mathbb Z.
]

这是个**环**。

---

## ( \mathbb Q(\alpha) )

定义为：

[
\mathbb Q(\alpha)
=================

\text{the smallest field containing }\mathbb Q\text{ and }\alpha.
]

也可以理解成：

[
\mathbb Q(\alpha)
=================

{f(\alpha)/g(\alpha): f,g\in \mathbb Q[x],\ g(\alpha)\neq 0}.
]

如果 (\alpha) 是代数元，那它其实也可以写成

[
a_0+a_1\alpha+\cdots+a_{d-1}\alpha^{d-1},\qquad a_i\in\mathbb Q,
]

其中 (d) 是 (\alpha) 的最小多项式次数。

这是个**域**。

---

# 3. 所以 ( \mathbb Z[\alpha]\subset \mathbb Q(\alpha) ) 是啥意思

意思非常直接：

* ( \mathbb Z[\alpha] ) 里只允许**整数系数**
* ( \mathbb Q(\alpha) ) 里允许**有理数系数**
* 所以前者自然包含在后者里

就像

[
\mathbb Z[i]\subset \mathbb Q(i)
]

一样。

你可以把它看成：

> **同一个生成元 (\alpha)，但一个是“整数版”，一个是“有理数版”。**

---

# 4. 那 ( \mathbb Z[\zeta_n] ) 是什么

这就是把上面的 (\alpha) 换成 (\zeta_n)。

通常 (\zeta_n) 表示一个单位根，比如 primitive (n)-th root of unity。
于是

[
\mathbb Z[\zeta_n]
==================

{f(\zeta_n): f(x)\in\mathbb Z[x]}.
]

也就是：

> 所有用整数系数多项式代入 (\zeta_n) 得到的元素组成的环。

它是一个**环**，通常不是域。

相应地，

[
\mathbb Q(\zeta_n)
]

就是包含 (\mathbb Q) 和 (\zeta_n) 的最小域。

---

# 5. 为什么 PPT 里会从 ( \mathbb Z[x]/(x^n+1) ) 跳到 ( \mathbb Z[\zeta_n] )

因为作者想把“商环模型”换成“根生成模型”。

## 商环模型

[
\mathbb Z[x]/(x^n+1)
]

这里的元素是“模 (x^n+1) 的多项式类”。

## 根生成模型

找一个 (\alpha) 满足

[
\alpha^n=-1.
]

那么就可以把商环里的 (x) 的类看成这个 (\alpha)，于是元素

[
[p(x)]
]

对应成

[
p(\alpha).
]

这时得到的就是某个 (\mathbb Z[\alpha]) 里的元素。

---

# 6. 你最容易混的点：环和域到底各自干嘛

## ( \mathbb Z[\zeta_n] )

是**整数环/子环**视角。
适合谈：

* NTRU 里的整数系数对象
* 格
* 模 (q) 运算
* “元素本身是整数系数多项式”

## ( \mathbb Q(\zeta_n) )

是**数域**视角。
适合谈：

* 基底
* 线性变换
* 矩阵表示
* 共轭
* 范数、迹之类线性代数/数论操作

所以 PPT 里前一页偏环，后一页偏域。
这不是乱写，而是作者在**切换模型**。

---

# 7. 一个你现在可以直接记住的公式

如果 (\alpha) 的最小多项式次数是 (d)，那么：

[
\mathbb Q(\alpha)
=================

{a_0+a_1\alpha+\cdots+a_{d-1}\alpha^{d-1}: a_i\in\mathbb Q},
]

而

[
\mathbb Z[\alpha]
=================

{a_0+a_1\alpha+\cdots+a_{d-1}\alpha^{d-1}: a_i\in\mathbb Z}
]

通常是它里面的一个子环。

---

# 8. 放回你现在这页的语境里

你可以把这两页理解成：

## 在 NTRU 那页

作者说：

* 我们的对象活在
  [
  R=\mathbb Z[x]/(x^n+1)
  ]
* 这也可以看成某个
  [
  \mathbb Z[\zeta_n]
  ]
  类型的环

## 在 Matrix representations 那页

作者进一步把它放到更大的域

[
K_n=\mathbb Q(\zeta_n)
]

里，然后说：

* 这个域有一组基
* 所以“乘以 (f)”可以写成矩阵

也就是从“环里的元素”切到“域上的线性算子”。

---

# 9. 一句最直白的话

[
\mathbb Z[\zeta_n]
]

就是“用整数系数拿 (\zeta_n) 拼出来的所有东西”，是个环。

[
\mathbb Q(\alpha)
]

就是“用有理数系数拿 (\alpha) 拼出来的所有东西”，而且还允许做除法，所以是个域。

[
\mathbb Z[\alpha]\subset \mathbb Q(\alpha)
]

就是说：**整数系数版包含在有理数系数版里面。**

---

# 10. 最后给你一个一眼就懂的对应表

[
\mathbb Z[i]\subset \mathbb Q(i)
]

完全对应

[
\mathbb Z[\alpha]\subset \mathbb Q(\alpha).
]

所以你以后看到 ( \mathbb Z[\zeta_n] )，就把它先脑补成“(\zeta_n) 版的高斯整数”就行。

下一步如果继续，我建议我就只讲一个问题：

**为什么到了 ( \mathbb Q(\zeta_n) ) 里，‘乘以 (f)’ 自然就是一个矩阵。**

> 

可以这么**粗略理解**，但要稍微修正一下说法：

## 对，直觉上就是“给 (\mathbb Z) 加进一个新元素”

比如加进 (\alpha)，并要求它满足某个代数关系，比如

[
\alpha^n=-1.
]

于是得到

[
\mathbb Z[\alpha].
]

这可以理解成：

> 从 (\mathbb Z) 出发，把 (\alpha) 也允许拿来做加法和乘法，于是生成一个更大的环。

---

## 但更准确地说，它是 **(\mathbb Z) 的环扩张**

不是通常说的“域扩张”。

因为：

* (\mathbb Z) 不是域
* (\mathbb Z[\alpha]) 通常也不是域
* 所以这里更准确的术语是
  **ring extension（环扩张）**，或者说
  **adjoin 一个元素 (\alpha)**

也就是：

[
\mathbb Z \subset \mathbb Z[\alpha].
]

---

## 而

[
\mathbb Q(\alpha)
]
才更像“(\mathbb Q) 的域扩张”

因为：

* (\mathbb Q) 是域
* (\mathbb Q(\alpha)) 也是域
* 所以这里才是标准意义上的 **field extension**

即

[
\mathbb Q \subset \mathbb Q(\alpha).
]

---

# 你可以这样记

## 1. 整数版

[
\mathbb Z[\alpha]
]
= 给 (\mathbb Z) 加进 (\alpha) 得到的**环**

## 2. 有理数版

[
\mathbb Q(\alpha)
]
= 给 (\mathbb Q) 加进 (\alpha) 得到的**域**

---

# 最经典类比

## 环扩张

[
\mathbb Z \subset \mathbb Z[i]
]

## 域扩张

[
\mathbb Q \subset \mathbb Q(i)
]

所以你这句“说白了就是 (\mathbb Z) 的扩张”，在直觉上是对的；
只是更严谨地说：

> **(\mathbb Z[\zeta_n]) 是 (\mathbb Z) 的环扩张，(\mathbb Q(\zeta_n)) 是 (\mathbb Q) 的域扩张。**

如果你愿意，我下一条可以顺着这个，直接讲清楚为什么
[
\mathbb Z[x]/(x^n+1)\cong \mathbb Z[\alpha].
]
