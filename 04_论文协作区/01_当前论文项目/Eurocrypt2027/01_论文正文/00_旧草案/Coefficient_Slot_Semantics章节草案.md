# Coefficient-Slot Semantics 章节草案

> paper_id: paper_2027_eurocrypt_fhe  
> 用途：新初稿 Section 3 的正文草案。  
> 写作目标：形式化说明 coefficient split、slot-preserving split、coefficient-packed slots 与 \(B_k\) 的区别，并解释为什么 CKKS bootstrapping core 不需要先恢复 big-ring slots。

## Section Draft: Coefficient-Split Semantics

### 3.1 Ciphertext-Level Split and Message-Level Meaning

Let \(N=kn\). For every message polynomial

\[
P(X)\in\mathcal R_N=\mathbb C[X]/(X^N+1),
\]

there is a unique decomposition

\[
P(X)=\sum_{r=0}^{k-1}X^rP_r(X^k),
\qquad
P_r(Y)\in\mathcal R_n=\mathbb C[Y]/(Y^n+1).
\]

This decomposition defines the coefficient split of \(P(X)\):

\[
\operatorname{CoeffSplit}_k(P(X))
=
(P_0(Y),\ldots,P_{k-1}(Y)).
\]

If

\[
P(X)=\sum_{i=0}^{N-1}a_iX^i,
\]

then

\[
P_r(Y)=\sum_{j=0}^{n-1}a_{r+jk}Y^j.
\]

Thus coefficient split is residue-class grouping of polynomial coefficients.

At the ciphertext level, the split algorithm is not a public extraction of coefficients. It is a ciphertext operation, typically implemented by ring switching and key switching. Its intended semantic specification is:

\[
\mathbf{ct}_N\vDash_N P(X)
\quad\Longrightarrow\quad
\mathsf{Split}_k(\mathbf{ct}_N)
=
(\mathbf{ct}_{n,0},\ldots,\mathbf{ct}_{n,k-1}),
\]

with

\[
\mathbf{ct}_{n,r}\vDash_n P_r(Y)
\qquad
\text{for }r=0,\ldots,k-1.
\]

This is the central semantic convention of our construction: algorithms act on ciphertexts, while the notation above describes their induced action on encrypted message polynomials.

The inverse operation is merge. Given ciphertexts

\[
\mathbf{ct}_{n,r}\vDash_n P_r(Y),
\qquad r=0,\ldots,k-1,
\]

the intended semantic effect of \(\mathsf{Merge}_k\) is

\[
\mathsf{Merge}_k(\mathbf{ct}_{n,0},\ldots,\mathbf{ct}_{n,k-1})
\vDash_N
\sum_{r=0}^{k-1}X^rP_r(X^k).
\]

The precise ciphertext-level cost, error, and key material required for these operations are implementation issues. The present section fixes the message semantics that the later correctness theorem relies on.

### 3.2 Coefficient Split Is Not Slot Selection

The coefficient split of \(P(X)\) generally does not preserve the original big-ring canonical slots.

Consider the case \(k=2\), so \(N=2n\). Write

\[
P(X)=A(X^2)+XB(X^2),
\qquad
A(Y),B(Y)\in\mathcal R_n.
\]

Let \(\alpha\in\Omega_N\) be a root of \(X^N+1\), and let

\[
\theta=\alpha^2\in\Omega_n.
\]

Then the two big-ring slots over the same \(\theta\)-fiber satisfy

\[
P(\alpha)=A(\theta)+\alpha B(\theta),
\]

and

\[
P(-\alpha)=A(\theta)-\alpha B(\theta).
\]

Therefore the leaf slot values \(A(\theta)\) and \(B(\theta)\) are recovered from the paired big-ring slots by

\[
A(\theta)
=
\frac{P(\alpha)+P(-\alpha)}{2},
\]

and

\[
B(\theta)
=
\frac{P(\alpha)-P(-\alpha)}{2\alpha}.
\]

Thus coefficient split maps paired big-ring canonical slots to local linear combinations. It does not select a subset of the original big-ring slot values.

For general \(k\), fix \(\theta\in\Omega_n\). The fiber above \(\theta\) is

\[
\pi^{-1}(\theta)
=
\{\alpha\in\Omega_N:\alpha^k=\theta\}.
\]

For every \(\alpha\in\pi^{-1}(\theta)\),

\[
P(\alpha)=\sum_{r=0}^{k-1}\alpha^r P_r(\theta).
\]

Hence the vector of big-ring slot values over the fiber is obtained from

\[
(P_0(\theta),\ldots,P_{k-1}(\theta))
\]

by a \(k\times k\) Vandermonde matrix

\[
V_\theta=(\alpha^r)_{\alpha\in\pi^{-1}(\theta),\,0\le r<k}.
\]

Coefficient split therefore changes the slot basis inside each fiber. It is a coefficient/module decomposition, not a slot-preserving split.

### 3.3 Slot-Preserving Split and the Recovery Transform \(B_k\)

A split is slot-preserving if the leaf ciphertexts carry canonical-slot values that are selected from, or are a fixed permutation of, the original big-ring canonical-slot vector.

Specifically, we call a split operation slot-preserving if there exists a fixed coordinate permutation \(\sigma\) such that, for every \(P\in\mathcal R_N\), the canonical slots of the leaf messages satisfy \(\operatorname{can}_n(Q_r)_j = \operatorname{can}_N(P)_{\sigma^{-1}(r,j)}\).For example, when \(N=4\), \(k=2\), and \(n=2\), a slot-preserving split could map\(
\operatorname{can}_4(P)=(z_0,z_1,z_2,z_3)\)to two leaf messages \(Q_0,Q_1\in\mathcal R_2\) satisfying \(\operatorname{can}_2(Q_0)=(z_0,z_2)\), \(\operatorname{can}_2(Q_1)=(z_1,z_3)\).    
This split preserves slots because the leaf slots are exactly a fixed permutation of the original big-ring slots.

Coefficient split does not have this property. If one wants to recover the original big-ring canonical slots from coefficient-split leaf slots, one needs an additional slot-recovery transform.


Let
\[
u_r(\theta)=P_r(\theta),
\qquad r=0,\ldots,k-1,
\]
be the coefficient-split leaf canonical-slot values at \(\theta\in\Omega_n\).

Choose, for each \(\theta\in\Omega_n\), an ordering
\[
\pi^{-1}(\theta)=\{\alpha_0(\theta),\ldots,\alpha_{k-1}(\theta)\}.
\]
Define the slot-recovery transform
\[
B_k:\bigoplus_{r=0}^{k-1}\mathsf{Can}_n
\to
\bigoplus_{t=0}^{k-1}\mathsf{Can}_n
\]
fiber-wise by
\[
(B_k(u_0,\ldots,u_{k-1}))_t(\theta)
=
\sum_{r=0}^{k-1}\alpha_t(\theta)^r u_r(\theta).
\]

If \(u_r(\theta)=P_r(\theta)\), then the \(t\)-th recovered leaf slot at \(\theta\) is
\[
(B_k(u_0,\ldots,u_{k-1}))_t(\theta)
=
P(\alpha_t(\theta)).
\]
Thus \(B_k\) converts coefficient-split leaf slots into slot-preserving leaf slots.
In the special case \(k=2\), this recovery can be written as

\[
Q^+(Y)=A(Y)+h(Y)B(Y),
\]

\[
Q^-(Y)=A(Y)-h(Y)B(Y),
\]

where \(h(\theta)=\alpha\) for a chosen lift \(\alpha\) with \(\alpha^2=\theta\). Such an \(h(Y)\) exists over \(\mathbb C\) by interpolation on the finite root set \(\Omega_n\), after choosing one lift above each \(\theta\). The choice is not canonical.

The transform \(B_k\) is therefore a canonical-slot lifting or recovery transform. It is relevant when the goal is to reconstruct the original big-ring slots after coefficient split.



### 3.4 Coefficient-Packed Slots Are the Right Invariant for Bootstrapping

CKKS bootstrapping does not apply EvalMod directly to the original canonical slots of \(P(X)\). The bootstrapping core first applies \(\mathsf{C2S}_N\). Therefore the input to EvalMod is not

\[
\operatorname{can}_N(P(X)).
\]

It is

\[
\operatorname{coeff}_N(P(X)).
\]

At the ciphertext level, if

\[
\mathbf{ct}_N\vDash_N P(X),
\]

then

\[
\mathsf{C2S}_N(\mathbf{ct}_N)
\vDash_N^{\mathrm{slot}}
\operatorname{coeff}_N(P(X)).
\]

Thus the slots after \(\mathsf{C2S}_N\) are coefficient-packed slots.

Now apply coefficient split first. If

\[
\mathsf{Split}_k(\mathbf{ct}_N)
=
(\mathbf{ct}_{n,0},\ldots,\mathbf{ct}_{n,k-1})
\]

and

\[
\mathbf{ct}_{n,r}\vDash_n P_r(Y),
\]

then each leaf \(\mathsf{C2S}_n\) satisfies

\[
\mathsf{C2S}_n(\mathbf{ct}_{n,r})
\vDash_n^{\mathrm{slot}}
\operatorname{coeff}_n(P_r(Y)).
\]

If

\[
P(X)=\sum_{i=0}^{N-1}a_iX^i,
\]

then

\[
\operatorname{coeff}_n(P_r(Y))
=
(a_{r},a_{r+k},a_{r+2k},\ldots,a_{r+(n-1)k}).
\]

Therefore

\[
(\operatorname{coeff}_n(P_0(Y)),\ldots,\operatorname{coeff}_n(P_{k-1}(Y)))
=
\Pi_k\operatorname{coeff}_N(P(X)).
\]

Equivalently,

\[
\left(\bigoplus_{r=0}^{k-1}\mathsf{C2S}_n\right)
\circ
\mathsf{Split}_k
\quad\text{has the same message-coordinate effect as}\quad
\Pi_k\circ\mathsf{C2S}_N.
\]

This is the key semantic point. The object that must be preserved across the split is not the original big-ring canonical-slot vector. The object that must be preserved is the coefficient-packed slot vector after \(\mathsf{C2S}\), up to residue-class grouping.

### 3.5 EvalMod Compatibility Under Coefficient-Packed Semantics

Let

\[
\mathsf{EvalMod}_N=f^{\oplus N},
\qquad
\mathsf{EvalMod}_n=f^{\oplus n}
\]

for the same scalar function \(f:\mathbb C\to\mathbb C\) under the same normalization.

For

\[
x=\operatorname{coeff}_N(P(X))=(a_0,\ldots,a_{N-1}),
\]

we have

\[
\Pi_k f^{\oplus N}(x)
=
\left(
(f(a_{0+jk}))_{j=0}^{n-1},
\ldots,
(f(a_{k-1+jk}))_{j=0}^{n-1}
\right).
\]

On the other hand,

\[
\left(\bigoplus_{r=0}^{k-1}f^{\oplus n}\right)\Pi_k x
=
\left(
(f(a_{0+jk}))_{j=0}^{n-1},
\ldots,
(f(a_{k-1+jk}))_{j=0}^{n-1}
\right).
\]

Hence

\[
\left(\bigoplus_{r=0}^{k-1}f^{\oplus n}\right)\circ\Pi_k
=
\Pi_k\circ f^{\oplus N}.
\]

This compatibility does not require \(f\) to be linear. It only requires that the same scalar function be applied independently to each coefficient-packed slot, and that the layout map contain no unnormalized coordinate-dependent scaling.

### 3.6 Why the No-\(B_k\) Route Is Correct for the Core

The slot-recovery transform \(B_k\) solves the following problem:

\[
\text{leaf canonical slots}
\longrightarrow
\text{big-ring canonical slots}.
\]

But the bootstrapping core does not need to solve this problem immediately after split. Instead, it follows the route

\[
\mathsf{Split}_k
\longrightarrow
\bigoplus_r\mathsf{C2S}_n
\longrightarrow
\bigoplus_r\mathsf{EvalMod}_n
\longrightarrow
\bigoplus_r\mathsf{S2C}_n
\longrightarrow
\mathsf{Merge}_k.
\]

After \(\mathsf{Split}_k\), the leaf ciphertexts carry \(P_r(Y)\). These leaf canonical slots are not original big-ring slots. However, after leaf \(\mathsf{C2S}_n\), they carry

\[
\operatorname{coeff}_n(P_r(Y)),
\]

which are exactly the residue-class sub-vectors of \(\operatorname{coeff}_N(P(X))\).

Since EvalMod acts component-wise on coefficient-packed slots, it can be executed independently on each leaf. Finally, leaf \(\mathsf{S2C}_n\) converts the reduced coefficient-packed slots back to leaf canonical representation, and \(\mathsf{Merge}_k\) assembles the corresponding large-ring polynomial.

Thus the no-\(B_k\) route is not claiming that coefficient split preserves big-ring slots. It claims something more precise:

\[
\text{coefficient split preserves the C2S-output coefficient coordinates up to residue grouping.}
\]

This is exactly the invariant required by the CKKS bootstrapping core.

### 3.7 Section Summary

The semantic conclusions of this section are:

1. \(\mathsf{Split}_k\) is a ciphertext operation whose intended message effect is coefficient/module decomposition of \(P(X)\), not public coefficient extraction.

2. Coefficient split is not slot-preserving. In canonical-slot coordinates, it applies a fiber-wise change of basis.

3. The transform \(B_k\) recovers big-ring canonical slots from coefficient-split leaf slots. It is necessary for slot-preserving recovery, but it is not necessary for the bootstrapping core.

4. CKKS bootstrapping applies EvalMod to coefficient-packed slots after \(\mathsf{C2S}\). Therefore the correct invariant is preservation of coefficient coordinates after \(\mathsf{C2S}\), not preservation of original canonical slots.

