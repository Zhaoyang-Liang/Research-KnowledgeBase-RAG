太好了，既然这两句话帮你打通了理解，我就帮你整理成一份**复习笔记**，你可以随时回顾。

---

## 复习笔记：如何把“找短原像”变成“几何问题（CVP）”

### 1. 问题设定
- 已知矩阵 \(\mathbf{A} \in \mathbb{Z}_Q^{n \times m}\)，目标向量 \(\mathbf{u} \in \mathbb{Z}_Q^n\)。
- 想找一个 **短** 的整数向量 \(\mathbf{x} \in \mathbb{Z}^m\) 满足：
  \[
  \mathbf{A} \mathbf{x} \equiv \mathbf{u} \pmod{Q}
  \]

### 2. 解的结构（陪集）
- 任取一个特解 \(\mathbf{x}_0\)（满足方程，不一定短）。
- 那么所有解可以写成：
  \[
  \mathbf{x} = \mathbf{x}_0 + \mathbf{y}, \quad \text{其中 } \mathbf{y} \in \Lambda_Q^\perp(\mathbf{A})
  \]
- 这里 \(\Lambda_Q^\perp(\mathbf{A}) = \{\mathbf{y} \in \mathbb{Z}^m : \mathbf{A}\mathbf{y} \equiv \mathbf{0} \pmod{Q}\}\) 是一个**格**（称为核格）。
- 因此，解集是格 \(\Lambda_Q^\perp(\mathbf{A})\) 的一个**陪集**：\(\mathbf{x}_0 + \Lambda_Q^\perp(\mathbf{A})\)。

### 3. 找短解等价于什么？
- 我们要找解 \(\mathbf{x} = \mathbf{x}_0 + \mathbf{y}\) 使得 \(\|\mathbf{x}\|\) 最小。
- 即：在陪集 \(\mathbf{x}_0 + \Lambda_Q^\perp(\mathbf{A})\) 中找一个离原点最近的向量。
- 这正是一个 **CVP**（最近向量问题）：
  > 给定格 \(\mathcal{L} = \Lambda_Q^\perp(\mathbf{A})\) 和目标点 \(\mathbf{t} = -\mathbf{x}_0\)，找一个格点 \(\mathbf{y} \in \mathcal{L}\) 使得 \(\|\mathbf{y} - \mathbf{t}\| = \|\mathbf{y} + \mathbf{x}_0\|\) 最小。  
  > 那么 \(\mathbf{x} = \mathbf{x}_0 + \mathbf{y}\) 就是原方程的一个最短解。

### 4. 为什么重要？
- 把**代数方程求短解**（SIS 求逆）转化成了**几何问题**（CVP）。
- 可以用格基约简、枚举等算法求解。
- 同时，如果知道格 \(\Lambda_Q^\perp(\mathbf{A})\) 的“好基”（陷门），就能快速解 CVP，从而快速签名。

### 5. 关键一句话
> **解集是格的一个陪集。找短解，就是在该陪集中找离原点最近的点——即 CVP。**

---


# 格密码复习笔记：解集作为陪集与CVP转化

## 1. 解集为什么是陪集？

设 **A** 是 \(n \times m\) 整数矩阵，模 \(Q\)。考虑所有整数向量 \(\mathbf{x} \in \mathbb{Z}^m\) 满足：

\[
\mathbf{A} \mathbf{x} \equiv \mathbf{u} \pmod{Q}
\]

假设 \(\mathbf{x}_0\) 是其中一个解（比如通过某种方法找到一个）。那么：

- 如果 \(\mathbf{x}\) 是任意解，则 \(\mathbf{A}(\mathbf{x} - \mathbf{x}_0) \equiv \mathbf{0} \pmod{Q}\)，所以  
  \(\mathbf{x} - \mathbf{x}_0 \in \Lambda_Q^\perp(\mathbf{A}) = \{ \mathbf{y} \in \mathbb{Z}^m : \mathbf{A} \mathbf{y} \equiv \mathbf{0} \pmod{Q} \}\)。
- 反过来，任何 \(\mathbf{x} = \mathbf{x}_0 + \mathbf{y}\)，其中 \(\mathbf{y} \in \Lambda_Q^\perp(\mathbf{A})\)，都是解。

因此所有解组成的集合是：

\[
\mathbf{x}_0 + \Lambda_Q^\perp(\mathbf{A})
\]

这就是一个**陪集**（子格 \(\Lambda_Q^\perp(\mathbf{A})\) 的平移）。

> **注意**：\(\Lambda_Q^\perp(\mathbf{A})\) 是一个格（整数点阵，加法子群）。陪集就是把这个格整个平移 \(\mathbf{x}_0\)。

---

## 2. “找短解”为什么变成“在陪集里找离原点最近的向量”？

“短解”是指 \(\mathbf{x}\) 的欧几里得长度 \(\|\mathbf{x}\|\) 小。我们在所有解中找长度最小的那个。

所有解是 \(\mathbf{x} = \mathbf{x}_0 + \mathbf{y}\)，\(\mathbf{y} \in \Lambda_Q^\perp(\mathbf{A})\)。所以：

\[
\min_{\mathbf{x} \in \text{solutions}} \|\mathbf{x}\| = \min_{\mathbf{y} \in \Lambda_Q^\perp(\mathbf{A})} \|\mathbf{x}_0 + \mathbf{y}\|
\]

这个最小值就是 **点 \(\mathbf{x}_0\) 到格 \(\Lambda_Q^\perp(\mathbf{A})\) 的距离**。因为我们要在格中找一个点 \(\mathbf{y}\) 使得 \(\mathbf{x}_0 + \mathbf{y}\) 离原点最近。等价于：

\[
\|\mathbf{x}_0 + \mathbf{y}\| = \|\mathbf{y} - (-\mathbf{x}_0)\|
\]

即在格 \(\Lambda_Q^\perp(\mathbf{A})\) 中找一个点 \(\mathbf{y}\) 使得它离 \(-\mathbf{x}_0\) 最近。

这就是**最近向量问题（CVP）**：
- 目标点 \(\mathbf{t} = -\mathbf{x}_0\)
- 格为 \(\Lambda_Q^\perp(\mathbf{A})\)
- 找最近的格点 \(\mathbf{y}^*\)
- 那么短解就是 \(\mathbf{x}^* = \mathbf{x}_0 + \mathbf{y}^*\)

---

## 关键结论

> **解集是格的一个陪集。找短解，就是在该陪集中找离原点最近的点——即 CVP。**