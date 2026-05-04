可以，这一步其实就是 **“乘法后出现 $s^2$，怎样把它重新变回只含 $(1,s)$ 的正常密文”**。

你现在看到的是：

$$
c(s)c'(s)
=(c_0+c_1s)(c_0'+c_1's)
=\alpha_0+\alpha_1 s+\alpha_2 s^2.
$$

问题不在乘法本身，而在于 **原来的解密器只会处理一次式**。
FV12 的标准密文是两项 $(c_0, c_1)$，解密时算

$$
c_0 + c_1s.
$$

所以乘完以后这个三项式
$$
\alpha_0+\alpha_1 s+\alpha_2 s^2
$$
不能直接当作“标准密文”继续用，必须把它压回
$$
\beta_0+\beta_1 s.
$$

这一步就是 **relinearization / key switching**。

---

## 1. 先看问题到底出在哪

把两个密文相乘后，得到扩展密文系数：

$$
(\alpha_0,\alpha_1,\alpha_2) = 
\begin{cases}
  \alpha_0 = c_0c_0' \\
  \alpha_1 = c_0c_1' + c_1c_0' \\
  \alpha_2 = c_1c_1'
\end{cases}
$$

于是它对应的“解密表达式”变成

$$
\alpha_0+\alpha_1 s+\alpha_2 s^2.
$$

你可以把它看成：

* 原来标准密钥是 $(1,s)$；
* 现在乘法后，实际上需要用扩展密钥 $(1,s,s^2)$。

也就是说，密文结构从

$$
\langle (c_0, c_1), (1,s) \rangle
$$

变成了

$$
\langle (\alpha_0, \alpha_1, \alpha_2), (1,s,s^2) \rangle.
$$

这当然不理想，因为每乘一次，次数都会再升。再乘一次会冒出 $s^3, s^4$。密文尺寸和解密密钥都会爆炸。

---

## 2. 核心想法：把 $s^2$ “换写”成一次式

文档里那句
$$
\alpha_0+\alpha_1 s+\alpha_2 s^2 = \beta_0+\beta_1 s
$$
本质意思不是说 **$s^2$ 在环里真的等于一次多项式**，而是说：

**我们预先发布一份“关于 $s^2$ 的加密信息”，使得任何 $\alpha_2 s^2$ 都能改写成某个 $\beta_0+\beta_1 s$ 的形式。**

换句话说，不是代数恒等变形，而是借助评估密钥完成的“密文层面的替换”。

---

## 3. 最朴素的直觉版

假设我们手里有一个“加密的 $s^2$”：

$$
\text{evk} \approx (\gamma_0, \gamma_1), \quad \text{满足} \quad \gamma_0 + \gamma_1 s \approx s^2.
$$

那就可以把 $\alpha_2 s^2$ 近似替换成

$$
\alpha_2 (\gamma_0 + \gamma_1 s).
$$

于是

$$
\alpha_0 + \alpha_1 s + \alpha_2 s^2
\approx
\alpha_0 + \alpha_1 s + \alpha_2 (\gamma_0 + \gamma_1 s)
$$

整理一下：

$$
(\alpha_0+\alpha_2\gamma_0)+(\alpha_1+\alpha_2\gamma_1) s.
$$

于是就得到

$$
\beta_0 = \alpha_0+\alpha_2\gamma_0, \qquad
\beta_1 = \alpha_1+\alpha_2\gamma_1,
$$

从而

$$
\alpha_0+\alpha_1 s+\alpha_2 s^2 \approx \beta_0 + \beta_1 s.
$$

这就是你看到那一步的本质。
它的思想非常简单：

**发布一个“$s^2$ 在密钥 $s$ 下的加密”，再用它把二次项吃掉。**

---

## 4. 为什么不能直接这么做

上面那个写法只是帮助理解，真正实现时不能直接拿一个普通的“$s^2$ 的加密”去乘 $\alpha_2$，因为：

* $\alpha_2$ 不是小数，而是模 $q$ 环里的大元素；
* 直接乘会让噪声炸掉；
* 所以要先把 $\alpha_2$ 按位分解，再逐位使用评估密钥。

这就是为什么真正的 relinearization / key switching 总是和 **digit decomposition / bit decomposition** 连在一起。文档前面也专门讲了 BitDecomp 和 Powersof2，这正是为这种步骤服务的。

---

## 5. 更正式一点的做法：分解 $\alpha_2$

设基数为 $T$（很多时候取 $2$ 或 $2^w$），把 $\alpha_2$ 展开成

$$
\alpha_2 = \sum_{i=0}^{\ell-1} \alpha_{2,i} T^i,
$$

其中每个 $\alpha_{2,i}$ 都是“小 digit”。

然后预先生成一组重线性化密钥：

$$
\text{rlk}_i = (a_i, b_i), \quad \text{满足} \quad b_i + a_i s \approx T^i s^2.
$$

注意这里“$\approx$”表示相差一个小噪声。
这组 $\text{rlk}_i$ 可以理解成：

**不是只发布一个 $s^2$ 的加密，而是发布很多个 $T^i s^2$ 的加密。**

这样就能处理任意系数 $\alpha_2$。

于是

$$
\alpha_2 s^2 = \left( \sum_{i=0}^{\ell-1} \alpha_{2,i} T^i \right) s^2
= \sum_{i=0}^{\ell-1} \alpha_{2,i} T^i s^2
\approx \sum_{i=0}^{\ell-1} \alpha_{2,i} (b_i + a_i s).
$$

所以

$$
\alpha_0+\alpha_1 s+\alpha_2 s^2
\approx
\alpha_0 + \alpha_1 s + \sum_{i=0}^{\ell-1} \alpha_{2,i} (b_i + a_i s).
$$

把常数项和 $s$ 项分组：

$$
\left( \alpha_0 + \sum_i \alpha_{2,i} b_i \right)
+ \left( \alpha_1 + \sum_i \alpha_{2,i} a_i \right) s.
$$

定义

$$
\beta_0 = \alpha_0 + \sum_i \alpha_{2,i} b_i, \qquad
\beta_1 = \alpha_1 + \sum_i \alpha_{2,i} a_i,
$$

就得到

$$
\alpha_0 + \alpha_1 s + \alpha_2 s^2
\approx
\beta_0 + \beta_1 s.
$$

这就是那一步的完整机制。

---

## 6. 这一步为什么叫 key switching

因为乘法后，你其实是在用“扩展密钥”

$$
(1, s, s^2)
$$

来解密扩展密文 $(\alpha_0, \alpha_1, \alpha_2)$。
而 relinearization 做的事，就是把它变成一个新密文 $(\beta_0, \beta_1)$，重新回到原来的密钥

$$
(1, s)
$$

下解密。

所以它叫：

* **relinearization**：因为从二次式降回一次式；
* **key switching**：因为从“关于 $(1, s, s^2)$ 的表示”切换回“关于 $(1, s)$ 的表示”。

---

## 7. 你可以这样记这一步

最简洁的记忆法是：

$$
(\alpha_0, \alpha_1, \alpha_2)
\quad\text{在密钥}\quad (1, s, s^2) \text{下解密}
$$

通过重线性化，变成

$$
(\beta_0, \beta_1)
\quad\text{在密钥}\quad (1, s) \text{下解密},
$$

并且保持解密值几乎不变：

$$
\alpha_0 + \alpha_1 s + \alpha_2 s^2
\approx
\beta_0 + \beta_1 s.
$$

其中“几乎”差的就是重新引入的一点小噪声。

---

## 8. 最后补一句：FV12 和 BGV 这里其实是同一个骨架

你现在看到的这步，在 FV12 里是

* 乘法后出现 $s^2$
* 用 relinearization 压回一次式

而在 BGV 那套语言里，本质也是一样的：
都是 **“乘法提升密文次数，重线性化把次数降回来”**。

你之所以会觉得这一句跳得快，是因为文档把“有一组加密的 $s^2$”这件事省略成了一句“最后一步是密钥交换/再线性化”。

下一步如果你要，我可以继续把 **FV12 的 relinearization key 具体怎么生成** 给你写成一套完整公式。
