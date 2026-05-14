如果固定 `h`，安全性分析要分两层：

1. **LWE/RLWE 的格攻击安全性**：primal、dual、hybrid、BDD 等。
2. **fixed-weight / sparse secret 额外攻击面**：攻击者知道 secret 只有 `h` 个非零项，会带来组合搜索、MITM、hybrid sparse-secret 攻击。

固定 `h` 本身不是不安全；关键看 `h/N` 是不是很小。

---

# 1. 先定义 fixed-h secret 模型

通常 fixed Hamming weight ternary secret 是：

```text
s ∈ {-1, 0, 1}^N

exactly h coefficients are nonzero
nonzero coefficients are ±1
```

也就是：

```text
wt(s) = h
```

论文里也把这种分布称为 **fixed Hamming weight secret**，并说明具体非零项采样方法可能依实现而不同。

如果 `h` 很小，比如论文说的 `h < 0.25 · n`，这种 secret 就属于 **sparse secret keys**。

---

# 2. 固定 h 对安全性的第一影响：secret space 变小

如果 secret 是普通 ternary：

```text
s_i ∈ {-1,0,1}
```

每个位置都有三种可能，secret 空间约为：

```text
3^N
```

如果固定 Hamming weight 为 `h`，攻击者知道只有 `h` 个非零位置，那么 secret 空间变成：

```text
|S_h| = C(N, h) · 2^h
```

其中：

* `C(N,h)`：选择哪 `h` 个位置非零；
* `2^h`：每个非零位置选 `+1` 或 `-1`。

所以 secret 熵是：

```text
H_fixed = log2 C(N,h) + h
```

这是最基础的安全下界检查。

如果要达到 `λ` bit 安全，至少要满足：

```text
log2 C(N,h) + h  >>  λ
```

但注意：

> **secret 熵足够大只是必要条件，不是充分条件。**

因为格攻击和 hybrid 攻击可以比暴力枚举更快。

---

# 3. 固定 h 对 bootstrapping 正确性是好事，但对安全性可能是坏事

在 CKKS bootstrapping failure probability 里，固定 `h` 的好处是：

```text
I(Y) 的分布确定
  ↓
ffail(K,h,M) 可以精确估计
  ↓
K 不需要为随机 h 的上尾波动加 correction factor
```

但是安全性上，攻击者也多知道了一件事：

```text
wt(s) = h
```

如果 `h` 很小，攻击者可以利用这个结构。

论文明确说，很多攻击和分析会利用 sparse secret 的性质，因此可能适用于使用 sparse secret 的 FHE 参数集；而本文没有纳入 sparse secret 参数表，因为 Lattice Estimator 当前不支持这些 sparse-secret 攻击成本估计。

所以：

```text
固定 h 大且稠密：
    主要影响是熵略变小，安全性接近普通 ternary。

固定 h 小且稀疏：
    会引入专门 sparse-secret 攻击，需要额外分析。
```

---

# 4. 安全性分析流程

如果你固定 `h`，我建议按下面流程分析。

---

## Step 1：把 FHE 参数转成 LWE/RLWE 参数

先明确底层实例：

```text
N      = ring dimension
q      = ciphertext modulus 或最大模数
σ      = error standard deviation
χ_s    = fixed-weight ternary secret
χ_e    = error distribution
m      = available LWE samples
```

如果是 GLWE，则常用：

```text
n = kN
```

论文也说 RLWE 和 GLWE 实例可以解释成 LWE 实例，所以它们集中估计 LWE 具体安全性。

---

## Step 2：跑普通 LWE lattice attack 估计

这一层仍然需要做。

也就是估计：

* primal attack；
* dual attack；
* decoding / BDD；
* hybrid attack；
* Coded-BKW 等。

论文使用 Lattice Estimator，并说明其支持 primal、dual、decoding、Coded-BKW、部分 hybrid combinatorial/lattice attacks。

但是这里要小心：

> 如果你的 fixed-h secret 很稀疏，不能只相信普通 Lattice Estimator 的结果。

因为论文特别指出，一些适用于 sparse secret 的攻击估计并不被当前 Lattice Estimator 支持，所以他们没有给 sparse secret 参数表。

实践上可以做两种近似：

```text
如果 h/N 不小：
    用 ternary / small secret 模型近似，跑普通 LWE estimator。

如果 h/N 很小：
    普通 estimator 只能作为参考；
    必须额外估计 sparse-secret 攻击。
```

---

## Step 3：检查暴力搜索复杂度

攻击者可以直接枚举所有 fixed-weight secrets。

成本约为：

```text
C_enum ≈ log2 C(N,h) + h
```

要求：

```text
C_enum ≥ λ
```

最好还要有安全余量。

如果这个值低于目标安全等级，直接不安全。

例如：

```text
N = 2^16, h = 64
```

虽然 `h` 很小，但：

```text
log2 C(N,64) + 64
```

仍然可能远大于 128。

但这不代表安全，因为 hybrid / MITM 攻击可能远快于完全枚举。

---

## Step 4：检查 meet-in-the-middle / hybrid sparse-secret 攻击

固定小 `h` 时，攻击者不一定枚举完整 secret，而是利用：

```text
只有 h 个位置非零
```

来做：

* meet-in-the-middle；
* hybrid lattice + guessing；
* sparse secret recovery；
* secret-position recovery；
* combinatorial filtering。

论文列出了一系列会利用 sparse secrets 的攻击和分析，并指出这些攻击可能适用于 FHE 参数集。

所以对 fixed-h secret，尤其是 sparse fixed-h，需要做额外问题：

```text
是否存在比普通 primal/dual 更快的 sparse-secret attack？
```

这一步不能只看普通 LWE security bit。

---

## Step 5：检查 fixed h 是否只是“去随机波动”，还是变成 sparse

这是最重要的判断。

设：

```text
p = h / N
```

如果 `p` 接近普通 ternary 的非零率，比如 `p ≈ 2/3`：

```text
h ≈ 2N/3
```

那么 fixed-h 只是把随机 ternary secret 条件化为固定重量。

安全影响通常不大，主要是：

```text
secret 熵少了大约 O(log N)
h 没有上尾波动
K 可以略微更紧
```

但如果：

```text
p << 1
```

比如：

```text
h/N < 0.25
```

论文就把它归入 sparse secret 的范围。

这时安全分析要明显升级。

---

# 5. 固定 h 的安全性可以这样分类

## 情况 A：固定 dense h

例如：

```text
h = 2N/3
```

这接近 ternary secret 的平均非零数。

安全分析：

```text
1. 计算 log2 C(N,h) + h
2. 用 LWE estimator 近似 small-secret / ternary security
3. 确认没有使用 sparse-secret 参数假设
4. bootstrapping failure 用精确 h 计算 K
```

这类 fixed-h 通常主要是为了：

```text
去掉 h 的随机波动
让 K 更精确
```

不是为了大幅降低 bootstrapping 成本。

---

## 情况 B：固定 sparse h

例如：

```text
h = 64, 128, 192
```

而 `N = 2^15` 或 `2^16`。

安全分析：

```text
1. 计算 secret entropy：log2 C(N,h) + h
2. 跑普通 LWE estimator，但只作为 baseline
3. 额外分析 sparse-secret attacks
4. 额外分析 hybrid guessing attacks
5. 避免直接套用论文 Table 5.2 / 5.3
```

这类 fixed-h 的好处是：

```text
K 明显变小
EvalMod 更容易
bootstrapping 更轻
```

但代价是：

```text
安全模型更复杂
攻击面更大
估计工具不够成熟
```

论文也正因为这个原因没有把 sparse secret 参数集放进本文参数表。

---

# 6. 一个实用判断标准

可以先看比例：

```text
p = h / N
```

然后粗略判断：

```text
p ≈ 1/2 或 2/3：
    dense fixed-h
    安全性接近普通 small ternary secret
    可用普通 LWE 估计 + 熵检查

p < 0.25：
    sparse fixed-h
    不能只用普通 LWE 估计
    必须考虑 sparse-secret 专门攻击

p 极小：
    高风险，需要非常谨慎
```

论文中也明确把 `h < 0.25 · n` 称为 sparse secret key 范围。

---

# 7. 和 K 的关系

固定 `h` 后，正确性参数 `K` 可以按：

```text
ffail(K,h,M) ≤ 2^δ
```

直接选。

但安全性不能只看：

```text
K 小不小
```

而要看：

```text
h 小不小
```

因为：

```text
h 小
  ↓
K 小
  ↓
bootstrapping 更便宜

但同时：

h 小
  ↓
secret 更稀疏
  ↓
sparse-secret 攻击更强
```

所以真正的 trade-off 是：

```text
减小 h 可以降低 K 和 bootstrapping 成本，
但会增加 fixed-weight / sparse-secret 安全风险。
```

---

# 8. 总结成一套安全分析 checklist

如果你固定 `h`，就按这个表检查：

```text
[1] 明确 secret 分布
    s ∈ {-1,0,1}^N, wt(s)=h

[2] 计算 secret 熵
    log2 C(N,h) + h

[3] 转成 LWE/RLWE 参数
    N, q, σ, χ_s, χ_e, m

[4] 跑普通 lattice attack 估计
    primal / dual / BDD / hybrid

[5] 判断 h/N 是否 sparse
    如果 h/N < 0.25，要额外分析 sparse-secret attacks

[6] 不要直接套普通 ternary 参数表
    sparse fixed-h 不在本文参数表覆盖范围内

[7] 正确性单独分析
    用 ffail(K,h,M) 选 K

[8] 最终取所有攻击成本的最小值
    security = min(
        lattice attack cost,
        enumeration cost,
        hybrid sparse attack cost,
        implementation-specific attack cost
    )
```

最重要的一句话：

> **固定 h 后，bootstrapping failure 更容易精确控制；但安全性要把 secret 分布从“普通 ternary”改成“fixed-weight ternary”重新估计。若 h 很小，它就是 sparse secret，不能只依赖普通 LWE Estimator 或论文中的通用参数表。**
