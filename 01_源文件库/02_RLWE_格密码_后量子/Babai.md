### 总结： s(目标点) 在 $\widetilde{b_n}$ 上的投影取整就是 c ， c 是使得c$\widetilde{b_n}$ + $span(b_1,b_2, b_{n-1})$离s最近的整数

这张图是在**用几何方式解释 Babai 的 nearest plane algorithm**。核心思想就一句话：

> **先决定“第 (n) 层”选哪一层，再把问题降到 (n-1) 维继续做。** 

下面我按图给你拆开讲。

---

## 1. 左图到底在画什么？

这里是一个 **rank 3 lattice** 的情形，也就是格由三个基向量
[
b_1,b_2,b_3
]
生成。

但 Babai 不是直接沿着 (b_3) 看“层”，而是沿着 **Gram-Schmidt 正交向量**
[
\widetilde b_3
]
来看。

所以左图里那一层一层**平行的平面**，其实是：

[
c\widetilde b_3+\operatorname{span}(b_1,b_2),\qquad c\in\mathbb Z.
]

也就是说，整个三维空间被切成很多张互相平行的“楼层”。
相邻两层的垂直距离就是

[
|\widetilde b_3|.
]

图里右侧那个大括号标的就是这个距离。

---

## 2. 为什么先把 (t) 投影成 (s)？

文中第一步写的是：

[
\text{Let } s \text{ be the projection of } t \text{ on } \operatorname{span}(b_1,\dots,b_n).
]

意思是：
如果 (t) 不在格张成的子空间里，那就先把它正交投影到这个子空间，得到 (s)。

### 为什么这样做不影响“最近格点”？

因为所有格点都在
[
\operatorname{span}(b_1,\dots,b_n)
]
里面。

把 (t) 分解成
[
t=s+u,
]
其中

* (s\in \operatorname{span}(b_1,\dots,b_n))
* (u\perp \operatorname{span}(b_1,\dots,b_n))

对任意格点 (y\in L(B))，都有
[
t-y=(s-y)+u,
]
而且 ((s-y)\perp u)。所以
[
|t-y|^2=|s-y|^2+|u|^2.
]

这里 (|u|^2) 对所有格点 (y) 都是一样的常数。
因此：

* 谁离 (s) 最近，
* 谁就离 (t) 最近。

所以先投影是合理的。文中下面那段话说的就是这个意思：**closest lattice vector to (s) is the same as closest lattice vector to (t)**。

---

## 3. 第二步“找最近的 hyperplane”是什么意思？

文中第二步是：

[
\text{Find } c \text{ such that } c\widetilde b_n+\operatorname{span}(b_1,\dots,b_{n-1})
]
is as close as possible to (s). 

拿这张图来说，就是：

[
c\widetilde b_3+\operatorname{span}(b_1,b_2)
]

这一族平行平面里，选一张**离 (s) 最近的**。

### 为什么是“选平面”而不是直接选点？

因为所有格点可以按“第 3 层”分组。

固定一个整数 (c) 之后，所有形如
[
cb_3+z,\qquad z\in L(b_1,b_2)
]
的格点，都会落在同一张平行平面里。
所以整格可以拆成一层一层的“格点切片”。

Babai 的策略就是：

* 先猜最近格点在哪一层；
* 再只在这一层里继续找。

图中**加粗的平面**就是“被选中的那一层”。

---

## 4. 这里为什么一会儿写 (c\widetilde b_3+\operatorname{span}(b_1,b_2))，一会儿又写 (cb_3+L(b_1,b_2))？

这是个很关键的点。

它们不是一回事，但彼此对应：

### (a) 平面本身

[
c\widetilde b_3+\operatorname{span}(b_1,b_2)
]

这是一个**连续的仿射平面**，里面有无穷多个实点。

### (b) 这一层上的格点集合

[
cb_3+L(b_1,b_2)
]

这是该平面里的**离散格点子集**。

---

为什么这两个东西会对应？

因为 Gram-Schmidt 分解里
[
b_3=\widetilde b_3+\mu_{3,1}\widetilde b_1+\mu_{3,2}\widetilde b_2.
]
后面那一坨就在 (\operatorname{span}(b_1,b_2)) 里，所以

[
cb_3+ \operatorname{span}(b_1,b_2)
==================================

c\widetilde b_3+\operatorname{span}(b_1,b_2).
]

也就是说：

* 用 (c\widetilde b_3) 来描述“第几层平面”更直观；
* 用 (cb_3+L(b_1,b_2)) 来描述“这层里的格点”更方便递归。

---

## 5. 第三步为什么要令 (s'=s-cb_3)？

文中第三步写：

[
s' = s-cb_n.
]

在这张图里就是
[
s'=s-cb_3.
]

意思是：

* 我们已经决定要在“第 (c) 层”里找；
* 那就把这一层整体平移回原点附近；
* 这样问题就变成一个 **rank 2** 的子问题。

因为“第 (c) 层”的格点集合是
[
cb_3+L(b_1,b_2).
]
如果整体减去 (cb_3)，就变成
[
L(b_1,b_2).
]

于是目标点也对应地变成
[
s'=s-cb_3.
]

这就是“降维递归”的本质。

---

## 6. 右图在表达什么？

右图就是把左图中选中的那一层，变成了一个 **二维问题**。

现在只剩下
[
L(b_1,b_2)
]
上的最近点问题了。

你可以理解成：

* 左图：先决定 3 维里选哪张“楼层”
* 右图：进入那张楼层后，只需要在这张二维平面里继续找最近格点

然后在二维里，再重复同样的逻辑：

* 用 (\widetilde b_2) 再切成一条一条平行线
* 选最近的一条
* 再降到 1 维

所以 Babai 本质上就是：

[
n\text{维} \to (n-1)\text{维} \to (n-2)\text{维} \to \cdots \to 1\text{维}.
]

---

## 7. 图里的粗线/粗平面是什么意思？

文中 caption 说：

> The chosen hyperplanes are thicker. 

意思是：

* 那些画得更粗的平面/线，
* 就是算法在这一层递归里**选中的那一层**。

也就是 Babai 认为“最近格点大概率在这一层”。

---

## 8. 这一页真正想表达的核心

这页不是在证明精度，而是在解释 **Babai 的递归几何结构**：

### 第一步

先把目标点 (t) 投影到格张成的子空间，得到 (s)。

### 第二步

用 (\widetilde b_n) 方向把空间分层，选离 (s) 最近的一层。

### 第三步

把这一层平移回去，问题变成低一维子格上的最近点问题。

### 第四步

递归求解，再加回 (cb_n)。

---

## 9. 你可以把这张图记成一句话

**Babai = 先选最近的“楼层”，再在楼层里继续找。**

其中：

* “楼层间距”由 (|\widetilde b_n|) 决定
* “楼层编号”由整数 (c) 决定
* 选完楼层后，就把问题降到 (L(b_1,\dots,b_{n-1}))

---

## 10. 你这里最容易卡住的点

你大概率卡在这一句：

[
c\widetilde b_n+\operatorname{span}(b_1,\dots,b_{n-1})
]

它不是说“最近格点等于 (c\widetilde b_n)”；

它只是说：

> **先按 (\widetilde b_n) 方向，把空间切成很多平行层；最近格点应该在某一层里。**

然后真正的格点是在这层中的
[
cb_n+L(b_1,\dots,b_{n-1})
]
里面继续找。

---

如果你愿意，我下一条可以继续给你讲：

## 为什么第二步里选 (c) 的公式正好是

[
c=\left\lfloor \frac{\langle s,\widetilde b_n\rangle}{|\widetilde b_n|^2}\right\rceil
]
也就是它和“最近平面”之间到底怎么对应。


---

对，**就是这个意思**。

这里的 (c) 就是要选一个整数，使得这张平面

[
H_c:=c\widetilde b_n+\operatorname{span}(b_1,\dots,b_{n-1})
]

离 (s) 最近。讲义第 2 步说的就是这个。

---

## 为什么？

关键点只有一个：

[
\widetilde b_n \perp \operatorname{span}(b_1,\dots,b_{n-1}).
]

也就是说，(\widetilde b_n) 是“竖直方向”，而
(\operatorname{span}(b_1,\dots,b_{n-1})) 是“底面”。

所以 (s) 可以唯一分解成

[
s=u+\alpha \widetilde b_n,
\qquad
u\in \operatorname{span}(b_1,\dots,b_{n-1}),
]
其中
[
\alpha=\frac{\langle s,\widetilde b_n\rangle}{|\widetilde b_n|^2}.
]

这个 (\alpha) 的含义就是：

> (s) 在“竖直方向” (\widetilde b_n) 上，位于第几层附近。

---

## 第 (c) 张平面长什么样？

第 (c) 张平面是

[
H_c=c\widetilde b_n+\operatorname{span}(b_1,\dots,b_{n-1}).
]

它里面离 (s) 最近的点，其实就是

[
u+c\widetilde b_n.
]

因为 (u) 已经是底面方向上的那部分了，你只需要把“高度”从 (\alpha) 改成 (c)。

所以 (s) 到这张平面的距离就是

[
|s-(u+c\widetilde b_n)|
=======================

# |(\alpha-c)\widetilde b_n|

|\alpha-c|,|\widetilde b_n|.
]

---

## 于是问题就变成了什么？

因为 (|\widetilde b_n|) 是常数，所以要让距离最小，只需要让

[
|\alpha-c|
]

最小。

而 (c) 必须是整数，因此最优选择就是：

[
c=\lfloor \alpha \rceil
=======================

\left\lfloor
\frac{\langle s,\widetilde b_n\rangle}{|\widetilde b_n|^2}
\right\rceil.
]

这就是为什么 Babai 里要取

[
c_j=
\left\lfloor
\frac{\langle b,\widetilde b_j\rangle}{|\widetilde b_j|^2}
\right\rceil.
]

---

## 三维图像版理解

你可以把这些平面想成一层一层的“楼板”：

* 第 (0) 层：(\operatorname{span}(b_1,\dots,b_{n-1}))
* 第 (1) 层：(\widetilde b_n+\operatorname{span}(b_1,\dots,b_{n-1}))
* 第 (2) 层：(2\widetilde b_n+\operatorname{span}(b_1,\dots,b_{n-1}))

而 (s) 的“楼层高度”就是 (\alpha)。

所以：

* 如果 (\alpha=2.2)，最近的是第 (2) 层
* 如果 (\alpha=2.8)，最近的是第 (3) 层
* 如果 (\alpha=-1.3)，最近的是第 (-1) 层

也就是**四舍五入**。

---

## 一句话总结

是的，(c) 就是那个让

[
c\widetilde b_n+\operatorname{span}(b_1,\dots,b_{n-1})
]

离 (s) 最近的整数。
原因是：**这些平面只在 (\widetilde b_n) 方向上彼此错开，所以比较“哪张平面最近”，等价于比较 (s) 在 (\widetilde b_n) 方向上的坐标离哪个整数最近。**

如果你愿意，我下一条可以继续讲：

## 为什么找平面时用的是 (c\widetilde b_n)，但真正递归平移时却减的是 (cb_n) 而不是 (c\widetilde b_n)。
