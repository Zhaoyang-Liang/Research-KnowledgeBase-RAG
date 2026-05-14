不要求小环的 ciphertext modulus \(Q\) 必须和大环完全一样。**代数 factorization theorem 本身是 message-level 的，不依赖 \(Q\)**；它只要求 split / C2S / EvalMod / S2C / merge 在明文语义上满足同一个规范化坐标语义。

但到 ciphertext-level 正确性时，\(Q\) 会影响三件事：

1. **语义兼容性**：小环降到更小的 \(Q_n\) 可以，但 \(\mathsf{RSDown}_{N\to n}\) 必须保证  
   \[
   \mathbf{ct}_{n,r}\vDash_n P_r(Y)
   \]
   在同一个 scale / normalization 下成立，或者额外的缩放因子必须被明确吸收到定义里。

2. **EvalMod 兼容性**：定理要求大环和小环 EvalMod 最终实现同一个标量函数  
   \[
   \mathsf{EvalMod}_N=f^{\oplus N},\qquad
   \mathsf{EvalMod}_n=f^{\oplus n}.
   \]
   所以小环 \(Q_n\) 可以更低，但不能导致 EvalMod 的输入 scale、period normalization、approximation polynomial 变成另一个函数 \(f_n\neq f\)，除非你重新证明带 scaling 的版本。

3. **安全性与噪声**：小环维度 \(n\) 更小，如果还用和大环一样大的 \(Q\)，RLWE 安全性通常更差。因此从安全分析角度，小环 \(Q_n\) 反而往往应该跟着降低。但降低多少要由 leaf-ring 安全估计和噪声预算共同决定。

所以结论是：

\[
Q_n \text{ 可以小于 } Q_N,
\]

但需要满足：

\[
\text{RSDown/RSUp 后 message scale 一致或可规范化，且 EvalMod 看到的是同一个 } f.
\]

最保守、最容易证明的版本是先写 ideal theorem 时假设 \(Q\) 不进入讨论；到参数与安全章节再写：实际实现允许 leaf modulus chain 降低，但必须证明 modulus switching / ring switching 引入的 rounding error 和 normalization 不破坏 Theorem 4.6 的假设。

如果你说的 \(Q\) 是之前 Section 2.5 里临时用的明文多项式 \(Q(X)\)，那它当然不会“不变”：split 后叶子里是 \(P_r(Y)\)，不是同一个大环多项式。为了避免歧义，正文里最好别用 \(Q\) 同时表示 modulus 和 message polynomial。

===


对，当前程序至少证明了一个工程事实：

\[
Q_{\text{leaf}} < Q_{\text{top}}
\]

时，no-\(B_k\) 路径仍可以正确工作到一定精度。比如 `leafCircuitLevels=9, leafNumP=3` 还能有约 `22.75 bits` 精度。

但要严格说：

- **数学 factorization theorem 不依赖 ciphertext modulus \(Q\)**。它是 message-coordinate 层面的等式。
- **实际 CKKS 正确性依赖 \(Q\)**，因为 \(Q\) 决定噪声预算、可用 level、rescale/modswitch 舍入误差、key-switching 误差和 EvalMod 可执行深度。
- **安全性依赖 \((n,Q,\chi_s,\chi_e)\)**。小环 \(n\) 更小，所以要降低 \(Q\) 或 \(QP\)，否则安全性下降。

模数降低主要通过下面这些数学关系影响噪声。

## 1. RLWE 安全性

RLWE 问题大致是：

\[
b(X)=a(X)s(X)+e(X)\pmod Q.
\]

安全性由这些量决定：

\[
n,\quad Q,\quad \sigma_e,\quad \chi_s.
\]

维度 \(n\) 越小，安全性越弱；模数 \(Q\) 越大，安全性也越弱。粗略说，攻击者面对的是：

\[
\frac{Q}{\sigma_e}
\]

越大，问题越容易。因此 leaf ring 变小后，如果仍保持大 \(Q\)，安全性会恶化。

所以你希望：

\[
(N,Q_N)\to(n,Q_n),
\qquad Q_n<Q_N.
\]

这是安全性层面的理由。

## 2. CKKS 正确性噪声预算

CKKS 解密形式可以写成：

\[
\operatorname{Dec}(\mathbf{ct}) = \Delta m + e \pmod Q.
\]

正确解码要求噪声不要太大：

\[
\|e\| \ll \Delta.
\]

更具体地，为了避免模 \(Q\) wrap-around，还需要：

\[
\|\Delta m + e\| < \frac{Q}{2}.
\]

所以 \(Q\) 降低会让可容忍范围变小：

\[
Q_n \downarrow
\quad\Longrightarrow\quad
\text{wrap-around margin decreases}.
\]

但只要 bootstrapping 中间值和噪声仍满足：

\[
\|\Delta m + e\| < Q_n/2,
\]

就还能正确。

## 3. DropLevel / Modulus Switching 的误差

你现在做的是把 leaf ciphertext 从 top level 截到 leaf level。理想上是：

\[
\mathbf{ct}\bmod Q_{\text{top}}
\quad\mapsto\quad
\mathbf{ct}\bmod Q_{\text{leaf}}.
\]

如果只是 drop 掉高层 primes，不做 rescale，那么 message scale \(\Delta\) 不变，但可用 modulus 变小。它主要影响的是剩余预算：

\[
\log Q_{\text{leaf}} - \log \Delta - \log \|m\| - \log \|e\|.
\]

这个量越小，越容易失败。

## 4. Key Switching 噪声

Split/Merge 和 bootstrapping 里大量用 key switching。典型 key switching 后噪声可以抽象成：

\[
e_{\text{out}}
=
e_{\text{in}}
+
e_{\text{ks}}.
\]

其中 \(e_{\text{ks}}\) 跟 gadget decomposition、\(P\)、分解基、密钥分布有关。粗略可写成：

\[
\|e_{\text{ks}}\|
\approx
O(d_{\mathrm{ks}}\cdot B_{\mathrm{ks}}\cdot \sigma_e\cdot \|s\|).
\]

这里：
- \(d_{\mathrm{ks}}\) 是 decomposition 长度；
- \(B_{\mathrm{ks}}\) 是 gadget base；
- \(P\) 是辅助模数，影响 key switching 精度和噪声；
- \(\|s\|\) 与 secret 稀疏/稠密有关。

所以降低 `leafNumP` 会减少 \(QP\) 安全压力，但可能增加 key-switching 误差或让 key switching 质量下降。你测试里 `leafNumP=3` 仍然能保持 `22.75 bits`，说明这一档还可行。

## 5. EvalMod 深度和 approximation error

EvalMod 本身会消耗很多 level。可以抽象成：

\[
\mathsf{EvalMod}(x)=f(x)+\epsilon_{\mathrm{approx}}+\epsilon_{\mathrm{eval}}.
\]

其中：
- \(\epsilon_{\mathrm{approx}}\)：多项式近似误差；
- \(\epsilon_{\mathrm{eval}}\)：同态乘法、rescale、rounding、key switching 误差。

每次乘法/rescale 都消耗 modulus primes。若 leaf \(Q_n\) 太短，就会出现：

\[
\text{remaining levels} < \text{EvalMod required levels}.
\]

或者虽然能跑完，但最后噪声过大，精度下降。

这正是你实验看到的趋势：

\[
\begin{array}{c|c}
\text{leafCircuitLevels} & \text{Avg L1 precision}\\
\hline
10 & 26.13 \text{ bits}\\
9 & 22.75 \text{ bits}\\
8 & 16.60 \text{ bits}\\
7 & 11.81 \text{ bits}
\end{array}
\]

这说明每减少一档 leaf modulus chain，正确性预算明显减少。

## 6. 论文里应该怎么表述

可以写成：

\[
\text{factorization correctness is modulus-agnostic at the message level,}
\]

but ciphertext correctness requires that the concrete parameters satisfy:

\[
B_{\mathrm{msg}} + B_{\mathrm{noise}} + B_{\mathrm{eval}} < \log Q_{\mathrm{leaf}}.
\]

更具体地：

\[
\log Q_{\mathrm{leaf}}
>
\log \Delta
+
\log \|m\|
+
\log \|e_{\mathrm{split}}\|
+
\log \|e_{\mathrm{C2S}}\|
+
\log \|e_{\mathrm{EvalMod}}\|
+
\log \|e_{\mathrm{S2C}}\|
+
\log \|e_{\mathrm{merge}}\|
+
\lambda_{\mathrm{margin}}.
\]

这不是精确 tight bound，但作为论文参数章节的结构非常好。

核心结论可以这样说：

\[
\boxed{
Q_{\mathrm{leaf}}\text{ 不出现在理想代数等式中，但它控制 ciphertext-level correctness margin 和 RLWE security.}
}
\]

所以数学推导允许 \(Q_n\) 跟随降低；实验说明降低到某个范围仍然正确；最终安全章节需要用 estimator 证明降低后的 \((n,Q_n,P_n)\) 达到目标安全级别。