# Preliminaries 章节草案

> paper_id: paper_2027_eurocrypt_fhe  
> 用途：新初稿 Section 2 的正文草案。  
> 写作原则：这里只建立 CKKS bootstrapping 与 subring decomposition 所需的公共语言；本文自己的 coefficient/slot 语义贡献放到下一章。

## Section Draft: Preliminaries

### 2.1 Rings, Messages, and Ciphertext Semantics

Let \(d\) be a power of two. We write

\[
\mathcal R_d=\mathbb C[X]/(X^d+1)
\]

for the complex message algebra used in the ideal message-level discussion, and

\[
\mathcal R_{q,d}=\mathbb Z_q[X]/(X^d+1)
\]

for the corresponding modular polynomial ring used by ciphertexts.

For a ciphertext \(\mathbf{ct}\) at ring dimension \(d\), scale \(\Delta\), and secret key \(s(X)\), we write

\[
\mathbf{ct}\vDash_{d,\Delta,s} P(X)
\]

to mean that \(\mathbf{ct}\) encrypts the message polynomial \(P(X)\in\mathcal R_d\) at scale \(\Delta\). Concretely, this means that decryption satisfies

\[
\operatorname{Dec}_{s(X)}(\mathbf{ct})
=
\Delta\cdot P(X)+e(X)
\quad \text{in } \mathcal R_{q,d},
\]

where \(e(X)\) is an error polynomial small enough for the current CKKS precision target.

When the scale, key, and modulus are clear from context, we abbreviate this as

\[
\mathbf{ct}\vDash_d P(X).
\]

This notation is a semantic judgment. It does not mean that \(P(X)\) is publicly available. All algorithms in this paper operate on ciphertexts; the judgment \(\vDash\) records the induced transformation on the encrypted plaintext message.

### 2.2 Canonical-Slot and Coefficient Coordinates

Let

\[
\Omega_d=\{\omega\in\mathbb C:\omega^d+1=0\}
\]

be the set of roots of \(X^d+1\). Fix an ordering of \(\Omega_d\). The full canonical embedding is the linear isomorphism

\[
\operatorname{can}_d:\mathcal R_d\to \mathbb C^d,
\qquad
P(X)\mapsto (P(\omega))_{\omega\in\Omega_d}.
\]

We call

\[
\operatorname{can}_d(P(X))
\]

the canonical-slot coordinate vector of \(P(X)\).

If

\[
P(X)=\sum_{i=0}^{d-1}a_iX^i,
\]

we define the coefficient coordinate vector

\[
\operatorname{coeff}_d(P(X))
=
(a_0,\ldots,a_{d-1})\in\mathbb C^d.
\]

Both

\[
\operatorname{can}_d:\mathcal R_d\to\mathsf{Can}_d
\]

and

\[
\operatorname{coeff}_d:\mathcal R_d\to\mathsf{Coeff}_d
\]

are fixed linear coordinate isomorphisms. Here \(\mathsf{Can}_d\) denotes the canonical-slot coordinate space and \(\mathsf{Coeff}_d\) denotes the coefficient-coordinate space.

In this section, \(d\) denotes the full complex embedding dimension of \(\mathbb C[X]/(X^d+1)\). This is not necessarily the number of user-visible CKKS complex slots under the standard half-slot convention. The half-slot convention can be recovered by restricting to the conjugate-compatible subspace and applying fixed projection and layout maps. The algebraic statements below are written in the full model to avoid layout-dependent distractions.

### 2.3 Slot Semantics for Ciphertexts

The semantic judgment

\[
\mathbf{ct}\vDash_d^{\mathrm{slot}} z
\]

means that \(\mathbf{ct}\) encrypts a polynomial \(P(X)\in\mathcal R_d\) whose canonical-slot coordinate vector is \(z\in\mathsf{Can}_d\). Formally,

\[
\mathbf{ct}\vDash_d^{\mathrm{slot}} z
\quad\Longleftrightarrow\quad
\mathbf{ct}\vDash_d P(X)
\text{ and }
z=\operatorname{can}_d(P(X)).
\]

Equivalently,

\[
\mathbf{ct}\vDash_d^{\mathrm{slot}} z
\quad\Longleftrightarrow\quad
\mathbf{ct}\vDash_d \operatorname{can}_d^{-1}(z).
\]

Thus "slots" are not a separate plaintext object. They are canonical coordinates of the encrypted message polynomial.

### 2.4 CoeffToSlot, SlotToCoeff, and Coefficient-Packed Slots

In the normalized message-level model, the ideal coefficient-to-slot and slot-to-coefficient transforms are the linear maps

\[
\mathsf{C2S}_d
=
\operatorname{coeff}_d\circ\operatorname{can}_d^{-1},
\qquad
\mathsf{S2C}_d
=
\operatorname{can}_d\circ\operatorname{coeff}_d^{-1}.
\]

At the ciphertext level, \(\mathsf{C2S}_d\) is an encrypted linear transform with the following semantic effect:

\[
\mathbf{ct}\vDash_d P(X)
\quad\Longrightarrow\quad
\mathsf{C2S}_d(\mathbf{ct})
\vDash_d^{\mathrm{slot}}
\operatorname{coeff}_d(P(X)).
\]

Equivalently,

\[
\mathsf{C2S}_d(\mathbf{ct})
\vDash_d
\operatorname{can}_d^{-1}(\operatorname{coeff}_d(P(X))).
\]

This motivates the following definition.

#### Definition 1: Coefficient-Packed Slots

Let \(\mathbf{ct}'\) be a ciphertext at ring dimension \(d\). We say that \(\mathbf{ct}'\) carries the coefficient-packed slots of \(P(X)\in\mathcal R_d\) if

\[
\mathbf{ct}'\vDash_d^{\mathrm{slot}}
\operatorname{coeff}_d(P(X)).
\]

Equivalently, \(\mathbf{ct}'\) encrypts the polynomial

\[
\operatorname{can}_d^{-1}(\operatorname{coeff}_d(P(X))).
\]

Thus coefficient-packed slots are still encrypted slots. The phrase means that the canonical slots of the current encrypted polynomial contain the coefficient vector of another polynomial \(P(X)\).

Conversely, \(\mathsf{S2C}_d\) has the semantic effect

\[
\mathbf{ct}'\vDash_d^{\mathrm{slot}} y
\quad\Longrightarrow\quad
\mathsf{S2C}_d(\mathbf{ct}')
\vDash_d
\operatorname{coeff}_d^{-1}(y),
\]

for every \(y\in\mathsf{Coeff}_d\), under the interpretation that \(y\) is a coefficient vector.

In particular, if

\[
\mathbf{ct}'\vDash_d^{\mathrm{slot}}
\operatorname{coeff}_d(Q(X)),
\]

then

\[
\mathsf{S2C}_d(\mathbf{ct}')
\vDash_d Q(X).
\]

### 2.5 The CKKS Bootstrapping Core

Ignoring modulus raising and ciphertext-level approximation errors, the core of CKKS bootstrapping has the semantic form

\[
\mathsf{C2S}_d
\quad\longrightarrow\quad
\mathsf{EvalMod}_d
\quad\longrightarrow\quad
\mathsf{S2C}_d.
\]

Let \(f:\mathbb C\to\mathbb C\) be the scalar approximation used by EvalMod under a fixed normalization. In the ideal normalized model, we assume

\[
\mathsf{EvalMod}_d=f^{\oplus d},
\]

meaning

\[
f^{\oplus d}(x_0,\ldots,x_{d-1})
=
(f(x_0),\ldots,f(x_{d-1})).
\]

Therefore, if

\[
\mathsf{C2S}_d(\mathbf{ct})
\vDash_d^{\mathrm{slot}}
(a_0,\ldots,a_{d-1}),
\]

then

\[
\mathsf{EvalMod}_d(\mathsf{C2S}_d(\mathbf{ct}))
\vDash_d^{\mathrm{slot}}
(f(a_0),\ldots,f(a_{d-1})).
\]

The important point is that EvalMod is component-wise on the slots after \(\mathsf{C2S}_d\). Since those slots are coefficient-packed slots, EvalMod is semantically applied component-wise to coefficient coordinates, not directly to the original canonical slots of \(P(X)\).

### 2.6 Subring and Module Decomposition

Let \(N=kn\). We use

\[
\mathcal R_N=\mathbb C[X]/(X^N+1),
\qquad
\mathcal R_n=\mathbb C[Y]/(Y^n+1).
\]

We identify \(\mathcal R_n\) with the subalgebra of \(\mathcal R_N\) generated by \(X^k\), via

\[
Y\mapsto X^k.
\]

This is well-defined because

\[
(X^k)^n=X^N=-1
\]

in \(\mathcal R_N\). It is injective because

\[
1,X^k,X^{2k},\ldots,X^{(n-1)k}
\]

are linearly independent in the monomial basis of \(\mathcal R_N\).

Every polynomial \(P(X)\in\mathcal R_N\) has a unique decomposition

\[
P(X)=\sum_{r=0}^{k-1}X^rP_r(X^k),
\qquad
P_r(Y)\in\mathcal R_n.
\]

If

\[
P(X)=\sum_{i=0}^{N-1}a_iX^i,
\]

then

\[
P_r(Y)=\sum_{j=0}^{n-1}a_{r+jk}Y^j.
\]

This decomposition is the algebraic basis for subring splitting. In this section it is only introduced as background. The next section explains its ciphertext-message semantics and its relation to slots.

### 2.7 Normalized Layout Convention

For

\[
p=\operatorname{coeff}_N(P(X))=(a_0,\ldots,a_{N-1}),
\]

define the residue-class grouping map

\[
\Pi_k:\mathsf{Coeff}_N\to\bigoplus_{r=0}^{k-1}\mathsf{Coeff}_n
\]

by

\[
\Pi_k p
=
\bigl(
(a_{0+jk})_{j=0}^{n-1},
(a_{1+jk})_{j=0}^{n-1},
\ldots,
(a_{k-1+jk})_{j=0}^{n-1}
\bigr).
\]

In the normalized theorem statements, \(\Pi_k\) is residue-class grouping, possibly followed by fixed coordinate permutations. It does not include arbitrary coordinate-dependent scalings. Any implementation-specific scaling must be normalized into the definitions of \(\mathsf{C2S}\), \(\mathsf{S2C}\), and the scalar EvalMod map \(f\), because nonlinear component-wise maps do not commute with arbitrary diagonal rescalings.

