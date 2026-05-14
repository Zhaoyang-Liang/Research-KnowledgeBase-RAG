下面详细讲附录 A：**CKKS bootstrapping failure probability**。这个附录讨论的不是 CKKS 的 LWE 安全性，而是 **bootstrapping 正确性失败概率**：也就是 bootstrapping 运行完后，输出明文是否可能被彻底“刷坏”。

---

# 1. 这个附录在讲什么问题？

CKKS bootstrapping 的核心目标是：

> 把一个噪声快满、层数快耗尽的 CKKS 密文，刷新成一个新的低噪声密文。

典型 CKKS bootstrapping 可以粗略分成几步：

```text
ModRaise
  ↓
CoeffToSlots
  ↓
EvalMod
  ↓
SlotsToCoeff
```

附录 A 专门讨论其中最关键也最危险的一步：**EvalMod**。

论文说，CKKS bootstrapping 的 failure probability 和 EvalMod 步骤密切相关。

---

# 2. EvalMod 到底在干什么？

CKKS bootstrapping 中，经过 modulus raising 之后，密文解密出来的东西大致形如：

```text
I(Y) · Q + Δm(Y)
```

论文中写作：

```text
I(Y) · Q + Δm(Y)
```

其中：

* `m(Y)` 是真正想保留的消息；
* `Δ` 是 CKKS scaling factor；
* `Q` 是之前的 ciphertext modulus；
* `I(Y)` 是某个整数多项式；
* `Y = X^{N / 2M}`；
* `M` 是 complex slots 的数量。

直觉上可以把它看成：

```text
解密值 = 整数倍的大模数部分 + 真正消息部分
```

EvalMod 的任务就是把那个整数倍部分 `I(Y) · Q` 去掉，只留下近似消息 `Δm(Y)`。

也就是：

```text
I(Y) · Q + Δm(Y)
        ↓ EvalMod
Δm(Y)
```

它靠同态计算一个类似模函数的东西完成：

```text
fmod(x) = x mod 1
```

论文说 EvalMod 会在一组区间上近似这个函数：

```text
⋃_{i=-K}^{K} [i - ε, i + ε]
```

意思是：它假设输入靠近某个整数 `i`，并且这个整数 `i` 的范围不会超过 `[-K, K]`。

---

# 3. 为什么会失败？

关键点是：

> EvalMod 只在有限区间 `[-K, K]` 范围内近似 `x mod 1`。如果真实的整数部分 `I(Y)` 的某个系数超过了这个范围，bootstrapping 就会失败。

论文明确说：

```text
如果 ||I(Y)|| > K，则 EvalMod 返回不可用的 corrupted plaintext。
```

也就是说，失败事件是：

```text
||I(Y)|| > K
```

所以 bootstrapping failure probability 定义为：

```text
ffail(K, h, M) = Pr[ ||I(Y)|| > K ]
```

其中：

* `K` 是 EvalMod 支持的整数区间范围；
* `h` 是 secret 的 Hamming weight；
* `M` 是 complex slots 数量。

---

# 4. `I(Y)` 为什么和 secret 的 Hamming weight 有关？

论文说，`I(Y)` 的每个系数可以看成：

```text
h + 1 个 [-0.5, 0.5) 上均匀随机变量的和
```

其中 `h` 是 secret 的 Hamming weight。

也就是设：

```text
U_1, ..., U_{h+1} ~ Uniform[-0.5, 0.5)
```

那么某个系数近似服从：

```text
S = U_1 + U_2 + ... + U_{h+1}
```

这就把问题变成了一个概率问题：

> `2M` 个这样的系数里，有没有某一个绝对值超过 `K`？

如果有，则：

```text
||I(Y)|| > K
```

bootstrapping 失败。

---

# 5. 公式从哪里来？

论文给出的失败概率是：

```text
ffail(K, h, M) = 1 - Pr[ ||I(Y)|| ≤ K ]
```

因为每个系数是 `h + 1` 个均匀变量的和，所以它服从 **Irwin–Hall 分布** 的平移版本。论文说这个失败概率可以通过改写 Irwin–Hall 累积分布函数计算。

我们分解一下。

令：

```text
S = U_1 + ... + U_{h+1}
```

其中：

```text
U_i ~ Uniform[-0.5, 0.5)
```

那么：

```text
Pr[ |S| ≤ K ]
```

表示一个系数落在安全范围内。

由于 `S` 的分布关于 0 对称：

```text
Pr[ |S| ≤ K ] = 2 Pr[S ≤ K] - 1
```

而 `Pr[S ≤ K]` 可以用 Irwin–Hall CDF 算出来。

所以单个系数安全的概率是：

```text
2F(K) - 1
```

如果一共有 `2M` 个相关系数需要同时落在范围内，论文使用：

```text
Pr[ ||I(Y)|| ≤ K ] = (2F(K) - 1)^{2M}
```

因此失败概率是：

```text
ffail(K, h, M) = 1 - (2F(K) - 1)^{2M}
```

这就是附录公式的结构。

---

# 6. 直觉解释：`K, h, M` 分别怎样影响失败概率？

失败概率：

```text
ffail(K, h, M) = Pr[ ||I(Y)|| > K ]
```

受三个参数控制。

---

## 6.1 `K` 越大，失败概率越小

`K` 是 EvalMod 可处理的整数范围。

```text
K 小：只近似很窄范围，容易失败
K 大：覆盖更大整数范围，不容易失败
```

但 `K` 不能无限大，因为 EvalMod 要近似周期模函数。近似范围越大，多项式近似越难，通常会影响：

* bootstrapping level cost；
* 近似误差；
* 最终精度；
* 参数复杂度。

所以 `K` 是 correctness 和 efficiency 之间的 trade-off。

---

## 6.2 `h` 越大，失败概率越大

`h` 是 secret 的 Hamming weight。

如果：

```text
S = U_1 + ... + U_{h+1}
```

那么 `h` 越大，求和项越多，`S` 的方差越大。

直觉上：

```text
h 小：S 集中在 0 附近
h 大：S 更容易偏离 0
```

所以 `h` 越大，`|S| > K` 的概率越高。

这也是为什么很多 CKKS bootstrapping 工作喜欢用 sparse secret：秘密越稀疏，`h` 越小，EvalMod 的整数部分越容易控制。

---

## 6.3 `M` 越大，失败概率越大

`M` 是 complex slots 数量。

论文公式里有：

```text
( ... )^{2M}
```

这表示要所有 `2M` 个系数都不失败。

即使单个系数失败概率很小，系数数量多了以后，总失败概率也会上升。

这和常见 union bound 的直觉一致：

```text
单点失败概率小 ≠ 整个向量失败概率小
```

slot 越多，需要同时正确的系数越多，因此整体 failure probability 更高。

---

# 7. 为什么固定 Hamming weight secret 更容易估计？

论文说，通常 bootstrapping 参数使用 fixed Hamming weight secret，这样可以精确知道 `h`，于是可以精确估计：

```text
ffail(K, h, M)
```

然后选一个合适的 `K` 使失败概率低于目标值。

比如你想要：

```text
ffail(K, h, M) ≤ 2^{-40}
```

就可以二分搜索 `K`。

但是论文这里遇到的问题是：

> 它使用的是 ternary secret，不是 fixed Hamming weight secret。

---

# 8. Ternary secret 下为什么麻烦？

论文说，它们的 ternary secret 每个系数按照概率：

```text
[p/2, 1 - p, p/2]
```

采样，并且 `p = 2/3`。

也就是每个系数：

```text
Pr[s_i = -1] = p/2
Pr[s_i = 0]  = 1 - p
Pr[s_i = 1]  = p/2
```

当 `p = 2/3` 时：

```text
Pr[s_i ≠ 0] = 2/3
```

所以 Hamming weight `h` 不是固定值，而是随机变量：

```text
h ~ Binomial(N, p)
```

其期望是：

```text
E[h] = Np
```

问题是：

> 如果只用 `E[h]` 来选 `K`，实际采样出来的 `h` 可能比期望更大，从而失败概率超过目标。

所以论文提出了一个修正方法。

---

# 9. 附录提出的 `K` 选择流程

目标是：

```text
找到 K，使得 Pr[ ||I(Y)|| > K ] ≤ 2^δ
```

其中 `δ < 0`。

例如：

```text
δ = -40
```

表示目标失败概率不超过：

```text
2^{-40}
```

论文给了三步。

---

## Step 1：用 `E[h]` 先估计一个基础 `K`

先把随机的 `h` 替换成期望：

```text
E[h] = Np
```

然后通过二分搜索找到一个 `K`，使得：

```text
ffail(K, E[h], M) ≤ 2^δ
```

这是第一步估计。

论文说这一步是 straightforward，可以通过 binary search 反复计算 `ffail` 完成。

---

## Step 2：估计修正因子 `K'`

因为真实 `h` 可能大于 `E[h]`，所以还要给 `K` 加一个 buffer：

```text
K_final = K + K'
```

论文根据正态近似估计 `h` 的波动。

因为：

```text
h ~ Binomial(N, p)
```

所以：

```text
E[h] = Np
σ_h = sqrt(Np(1-p))
```

论文写：

```text
σ_h = sqrt(Np(1-p))
```

然后它使用一个参数 `d`，考虑：

```text
h ≤ E[h] + dσ_h
```

这个事件的概率近似由误差函数 erf 给出：

```text
Pr[h ≤ E[h] + dσ_h] = erf(d / sqrt(2))
```

于是希望尾部概率小于目标：

```text
1 - erf(d / sqrt(2)) ≤ 2^δ
```

找到这样的 `d` 之后，就能估计 `K'`。

---

## Step 3：最终设置

最后：

```text
K := K + K'
```

也就是：

```text
K_final = 基于 E[h] 的 K + 考虑 h 波动的安全余量
```

这样就能处理 ternary secret 的随机 Hamming weight。

---

# 10. `K'` 的近似公式怎么理解？

论文中有一段推导：

```text
K = ceil(κ · sqrt(E[h] + 1))
```

这个形式来自中心极限定理直觉：

如果：

```text
S = U_1 + ... + U_{h+1}
```

那么 `S` 的典型大小随：

```text
sqrt(h+1)
```

增长。

所以为了保持相似的尾概率，`K` 大致应该和：

```text
sqrt(h+1)
```

成比例。

于是：

```text
K ≈ κ sqrt(h+1)
```

如果真实 `h` 比 `E[h]` 大了 `dσ_h`，那么 `K` 也要相应增加。

论文给出的近似修正是：

```text
K' = ceil(dκ sqrt(1-p))
```

这里的直觉是：

* `κ` 控制当前目标 failure probability 下的安全半径；
* `d` 控制 Hamming weight 偏离多少期望标准差；
* `sqrt(1-p)` 来自二项分布标准差和 `sqrt(h)` 缩放之间的近似关系。

---

# 11. 附录中的两个 helper function

论文最后说他们实现了两个辅助函数。

---

## 11.1 `Probability`

```text
Probability(Xs, K, log2(N), log2(M)) → δ
```

输入：

* `Xs`：secret distribution；
* `K`；
* `log2(N)`；
* `log2(M)`。

输出：

```text
δ = log2(Pr[ ||I(Y)|| > K ])
```

也就是返回 failure probability 的 log2 值。

例如如果返回：

```text
δ = -37.65
```

表示失败概率大约是：

```text
2^{-37.65}
```

这正是论文 Table 5.8 里 CKKS bootstrapping 参数表给出的量：`log2(Pr[||I(X)|| > K]) = -37.65`，对应 `K = 512`。

---

## 11.2 `FindSuitableK`

```text
FindSuitableK(Xs, log2(N), log2(M), δ) → K
```

输入：

* secret distribution；
* `N`；
* `M`；
* 目标 failure exponent `δ`。

输出：

```text
K
```

使得：

```text
Pr[ ||I(Y)|| > K ] ≤ 2^δ
```

如果 `Xs` 是概率分布而不是 fixed Hamming weight，它会自动考虑上面说的修正因子 `K'`。

---

# 12. 为什么计算公式很贵？

附录最后还有一个非常工程化的提醒。

公式需要计算：

```text
Σ (-1)^i binomial(h+1, i) (...)
```

这是一个交错和，而且里面有高次幂：

```text
(...)^{h+1}
```

论文说 Equation 1 需要大约 `2h` 精度的 arbitrary precision arithmetic，否则数值不稳定；直接计算复杂度大约是：

```text
O(h^3)
```

当 `h` 很大时非常昂贵。

所以实践中他们建议：

> 对一个足够大的固定 `h`，预计算 `(K, δ)` 表，然后用比例关系
> `K / sqrt(h+1) ≈ K' / sqrt(h'+1)`
> 来近似其他 `h'` 的结果。

这其实很重要：附录不是纯理论推导，而是为了让 CKKS bootstrapping 参数能被实际生成。

---

# 13. 用一句话总结附录 A

附录 A 的核心是：

> CKKS bootstrapping 的 EvalMod 只在有限整数范围 `[-K, K]` 上正确近似模函数；而这个整数部分 `I(Y)` 的大小由 secret Hamming weight 和 slots 数量决定。因此需要用 Irwin–Hall 分布估计 `Pr[||I(Y)|| > K]`，并为随机 ternary secret 的 Hamming weight 波动加入修正因子，从而选择足够安全但不过度保守的 `K`。

---

# 14. 和 CKKS 参数选择的关系

这部分其实连接了三个参数层面：

```text
secret distribution χ_s
        ↓
Hamming weight h
        ↓
I(Y) 的系数范围
        ↓
EvalMod 需要覆盖的 K
        ↓
bootstrapping failure probability
        ↓
bootstrapping 成本与最终精度
```

所以 CKKS bootstrapping 参数不是只看 LWE 安全表就够了。

你还必须看：

* secret 是否 sparse / ternary / fixed weight；
* slot 数量 `M`；
* EvalMod 近似区间 `K`；
* 目标 failure probability；
* bootstrapping 消耗层数；
* 输出 precision bits。

这也是为什么论文 Table 5.8 同时列出：

* `K = 512`；
* `log2(Pr[||I(X)|| > K]) = -37.65`；
* EvalMod level cost；
* Precision Bits；
* Scaling Factor；
* `log2(PQ)` 等参数。

---

# 15. 你可以这样记

```text
EvalMod 是 CKKS bootstrapping 的危险点。
它要消掉整数部分 I(Y)。

如果 I(Y) 的系数落在 [-K, K] 内：
    bootstrapping 正常。

如果某个系数超过 K：
    EvalMod 近似区间不覆盖它，
    输出 corrupted plaintext。

所以 failure probability = Pr[||I(Y)|| > K]。

I(Y) 的系数 ≈ h+1 个 Uniform[-1/2,1/2) 的和。
因此用 Irwin–Hall 分布估计尾概率。
```

最关键的 trade-off 是：

```text
K 越大，失败概率越低；
但 EvalMod 近似越难，bootstrapping 越贵，精度可能越差。
```

这就是附录 A 想解决的问题。
