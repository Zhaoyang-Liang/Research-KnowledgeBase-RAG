# 低模数构造式小环拆分原型

本目录验证一个最小构造：

\[
c_0(X)+s(X)c_1(X)=\Delta(P_0(X^2)+XP_1(X^2))+e(X)\pmod Q
\]

直接构造成两个低模数小环密文：

\[
ct_0'\vDash_{R_{n,q}}P_0(Y),\qquad ct_1'\vDash_{R_{n,q}}P_1(Y).
\]

这里 \(Y=X^2\)，且 \(q\mid Q\)。程序使用自包含的 toy negacyclic RLWE，不依赖 Lattigo bootstrapping，目的是单独验证代数构造和降模数语义。

## 算法结构

把 secret 和密文分量拆成奇偶部分：

\[
s(X)=s_0(Y)+Xs_1(Y),
\]

\[
c_0(X)=A(Y)+XC(Y),\qquad c_1(X)=B(Y)+XD(Y).
\]

则解密等式推出：

\[
A+s_0B+Ys_1D=\Delta P_0(Y)+\nu_0,
\]

\[
C+s_0D+s_1B=\Delta P_1(Y)+\nu_1.
\]

于是两行 module ciphertext 为：

\[
(A,B,YD),\qquad (C,D,B).
\]

再用辅助密钥

\[
\operatorname{KSK}_{(s_0,s_1)\to s'}
\]

切换到普通小环 secret \(s'\)。

## 运行命令

在父目录运行：

```bash
cd /Users/mac/Desktop/Research-KB-RAG/04_论文协作区/01_当前论文项目/Eurocrypt2027/04_实验实现/ckks-subring-compression
go run ./5-8-lowq-construct -trials=5 -seed=1
```

验证纯代数正确性：

```bash
go run ./5-8-lowq-construct -trials=3 -seed=1 -ksErr=0
```

提高 scale 的版本：

```bash
go run ./5-8-lowq-construct -trials=5 -seed=3 -q=1073741789 -qFactor=2 -scale=16777216 -ksErr=1 -base=16
```

## 当前观察

- module row 的误差与原始大环解密误差一致，说明抽取公式正确。
- 当 `ksErr=0` 时，最终 leaf ciphertext 只保留原始加密噪声，说明 module-to-ring key switching 代数正确。
- 当 `ksErr=1` 时，最终误差主要来自 toy key switching；这个原型没有实现真实 CKKS/RNS key switching 的 \(P\)-扩展与 ModDown，因此噪声比真实可调实现更粗糙。
- 因为使用 \(q\mid Q\)，程序做的是 exact level dropping，不是把 scale 乘上 \(q/Q\) 的近似 modulus switching。因此输出仍保持原来的 \(\Delta\)。

## 2026-05-08 初步对比结果

程序同时输出一个 keep-\(Q\) baseline：使用同样的构造和同样的 toy key switching，但目标模数保持 \(Q\)，用于对比“叶子密文继承大环 \(Q\)”和“构造时直接降到 \(q\)”。

| 参数 | 低模数构造精度 | keep-\(Q\) baseline 精度 | 平均时间：低模数构造 | 平均时间：keep-\(Q\) baseline | 模数降低 |
|---|---:|---:|---:|---:|---:|
| `N=32,n=16,Q/q=16,scale=2^16,ksErr=1` | 约 7.7--9.6 bits | 约 7.8--9.9 bits | 152 us | 166 us | 31 bits -> 27 bits，降低 4 bits，16x |
| `N=32,n=16,Q/q=2,scale=2^24,ksErr=1` | 约 15.6--17.5 bits | 约 15.6--17.0 bits | 29 us | 28 us | 31 bits -> 30 bits，降低 1 bit，2x |
| `N=1024,n=512,Q/q=16,scale=2^16,ksErr=1` | 约 5.0--5.6 bits | 约 4.7--5.9 bits | 18.8 ms | 20.2 ms | 31 bits -> 27 bits，降低 4 bits，16x |
| `N=1024,n=512,Q/q=16,scale=2^16,ksErr=0` | 16 bits | 16 bits | 19.2 ms | 20.1 ms | 31 bits -> 27 bits，降低 4 bits，16x |

解释：

- 低模数构造和 keep-\(Q\) baseline 的精度接近，说明降 \(q\) 本身没有破坏语义；精度主要由 key switching 噪声决定。
- toy 计时只能说明趋势，不能替代 Lattigo/RNS/NTT 实现的真实性能。这里 \(Q/q=16\) 时低模数构造略快，主要因为 gadget digit 少一位；\(Q/q=2\) 时差异被 Go toy 实现噪声淹没。
- 真正 Lattigo 默认 `SplitNew`/ring-packing 路线通常保持当前 RNS 链，也就是 leaf 继承大环 \(Q\)。本原型里的 keep-\(Q\) baseline 对应这一安全问题的抽象对照，但不是完整 Lattigo API 的逐行替代实现。

另外跑了一个完整 Lattigo 默认 bootstrapping 小维度 sanity baseline：

```bash
go run ./5-7 -run=baseline -logN=13 -reps=1 -q0Bits=36 -circuitPrimeBits=20 -defaultScale=25 -circuitLevels=1 -numP=1
```

结果：

| 方法 | 参数 | 总时间 | 精度 | \(Q\) 降低 |
|---|---|---:|---:|---:|
| Lattigo 默认 direct bootstrapping | `logN=13`, `logQP≈938` | 666 ms | Avg L1 precision 15.08 bits | 不降低 |

这个完整 Lattigo baseline 是 bootstrapping 总耗时，不能和 toy compression 原型的几十微秒/毫秒直接比较；它的作用是给出默认方法的精度和“不降低 \(Q\)”口径。
