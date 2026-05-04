
# 不做M会发生什么

可以，我就**直接乘**给你看，不讲别的。



为了看清楚，我取最小但还能看出问题的情形：

* slot 数 (\ell=2)
* 每个 slot 内部维度 (d=2)

这时矩阵 (U_\ell) 就是一个 (2\times 2) 矩阵，第 (i) 行第 (k) 列是 (\zeta_{m,i}^{dk})。所以

[
U_\ell=
\begin{bmatrix}
1 & \zeta_{m,0}^{2}\
1 & \zeta_{m,1}^{2}
\end{bmatrix}.
]

这和论文里 (U_\ell=(\zeta_{m,i}^{d\cdot j})) 的定义是一致的。

---

## 1. 先故意“不做 (M)”，直接拿原来的 fully packed 向量去乘

这时第 0 个 slot 用自己的局部基 (\zeta_{m,0})，第 1 个 slot 用自己的局部基 (\zeta_{m,1})。

所以向量是

[
\vec m=
\begin{bmatrix}
a_0+a_1\zeta_{m,0}\
b_0+b_1\zeta_{m,1}
\end{bmatrix}.
]

这就是论文 fully packed 的原始形式
[
\left(\sum_{j=0}^{d-1}m_{i,j}\zeta_{m,i}^j\right)_{0\le i<\ell}
]
在 (\ell=2,d=2) 的具体样子。

---

## 2. 现在直接乘 (U_\ell)

做普通矩阵乘法：

[
U_\ell \vec m
=============

\begin{bmatrix}
1 & \zeta_{m,0}^{2}\
1 & \zeta_{m,1}^{2}
\end{bmatrix}
\begin{bmatrix}
a_0+a_1\zeta_{m,0}\
b_0+b_1\zeta_{m,1}
\end{bmatrix}.
]

第一行乘出来是

[
(\star)_0
=========

(a_0+a_1\zeta_{m,0})
+
\zeta_{m,0}^{2}(b_0+b_1\zeta_{m,1}).
]

展开：

[
(\star)_0
=========

a_0+a_1\zeta_{m,0}+b_0\zeta_{m,0}^{2}+b_1\zeta_{m,0}^{2}\zeta_{m,1}.
]

第二行乘出来是

[
(\star)_1
=========

(a_0+a_1\zeta_{m,0})
+
\zeta_{m,1}^{2}(b_0+b_1\zeta_{m,1}),
]

展开：

[
(\star)_1
=========

a_0+a_1\zeta_{m,0}+b_0\zeta_{m,1}^{2}+b_1\zeta_{m,1}^{3}.
]

所以直接乘 (U_\ell) 的结果是

[
U_\ell \vec m=
\begin{bmatrix}
a_0+a_1\zeta_{m,0}+b_0\zeta_{m,0}^{2}+b_1\zeta_{m,0}^{2}\zeta_{m,1}[4pt]
a_0+a_1\zeta_{m,0}+b_0\zeta_{m,1}^{2}+b_1\zeta_{m,1}^{3}
\end{bmatrix}.
]

---

## 3. 问题出在哪

关键看第一个分量：

[
a_0+a_1\zeta_{m,0}+b_0\zeta_{m,0}^{2}+b_1\zeta_{m,0}^{2}\zeta_{m,1}.
]

你本来希望最后能整理成“第 0 个 slot 自己的局部基”：

[
c_0+c_1\zeta_{m,0}+c_2\zeta_{m,0}^2+c_3\zeta_{m,0}^3.
]

前三项还勉强像那么回事，但最后一项是

[
b_1\zeta_{m,0}^{2}\zeta_{m,1}.
]

它**卡住了**，因为这里混着两个不同 slot 的根：

* (\zeta_{m,0})
* (\zeta_{m,1})

你没法把它直接看成“只用 (\zeta_{m,0}) 的幂展开”的一项。
这就是“不能直接乘 (U)”的本质原因：

> **一乘完以后，不同 slot 的局部基混在一起了。**

---

## 4. 现在看“先做 (M)”会发生什么

先做 (M) 以后，两个 slot 都统一改写成同一个标准基 (\zeta_m)：

[
\vec m'=
\begin{bmatrix}
a_0+a_1\zeta_m\
b_0+b_1\zeta_m
\end{bmatrix}.
]

这正是论文第 1 步。

再乘同一个 (U_\ell)：

[
U_\ell \vec m'
==============

\begin{bmatrix}
1 & \zeta_{m,0}^{2}\
1 & \zeta_{m,1}^{2}
\end{bmatrix}
\begin{bmatrix}
a_0+a_1\zeta_m\
b_0+b_1\zeta_m
\end{bmatrix}.
]

第一行变成

[
(\vec m'')_0
============

(a_0+a_1\zeta_m)
+
\zeta_{m,0}^{2}(b_0+b_1\zeta_m)
]

也就是

[
(\vec m'')_0
============

a_0+a_1\zeta_m+b_0\zeta_{m,0}^{2}+b_1\zeta_{m,0}^{2}\zeta_m.
]

注意这里最后一项已经不是
[
\zeta_{m,0}^{2}\zeta_{m,1}
]
了，而是
[
\zeta_{m,0}^{2}\zeta_m.
]

也就是说，现在“内部基”统一成了同一个 (\zeta_m)，不再有两种不同的局部根混在一起。

---

## 5. 再做 (M^{-1}) 为什么就通了

第 3 步 (M^{-1}) 会把第 0 个分量里的 (\zeta_m) 变回 (\zeta_{m,0})，所以

[
a_0+a_1\zeta_m+b_0\zeta_{m,0}^{2}+b_1\zeta_{m,0}^{2}\zeta_m
]

会变成

[
a_0+a_1\zeta_{m,0}+b_0\zeta_{m,0}^{2}+b_1\zeta_{m,0}^{2}\zeta_{m,0}.
]

于是最后一项就能合并成

[
b_1\zeta_{m,0}^{3}.
]

整个第 0 个分量就整理成

[
a_0+a_1\zeta_{m,0}+b_0\zeta_{m,0}^{2}+b_1\zeta_{m,0}^{3},
]

这正是你想要的“统一按指数 (j+dk) 排列”的形式。论文第 3 步也正是得到
[
\sum_{k,j} m_{k,j}\zeta_{m,i}^{j+dk}.
]

---

## 6. 最短总结

你要的“直接乘 (U)”结果就是：

[
(a_0+a_1\zeta_{m,0})+\zeta_{m,0}^2(b_0+b_1\zeta_{m,1})
======================================================

a_0+a_1\zeta_{m,0}+b_0\zeta_{m,0}^2+b_1\zeta_{m,0}^2\zeta_{m,1}.
]

问题就出在最后那项：

[
\zeta_{m,0}^2\zeta_{m,1}.
]

它把两个不同 slot 的局部根混到了一起，所以你没法把结果继续看成“第 0 个 slot 自己的幂基展开”。

而先做 (M) 的作用，就是把 (\zeta_{m,1}) 先统一成公共的 (\zeta_m)，这样再经过 (M^{-1}) 才能在每个 slot 内部重新合并指数。

这就是为什么**不能直接乘 (U)**。
