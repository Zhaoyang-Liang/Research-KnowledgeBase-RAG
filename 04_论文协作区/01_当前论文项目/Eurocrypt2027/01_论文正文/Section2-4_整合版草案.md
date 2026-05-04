# Section 2-4 整合版草案

> paper_id: paper_2027_eurocrypt_fhe  
> 用途：将 `Preliminaries`、`Coefficient-Split Semantics`、`Core Factorization Theorem` 整合为连续可读的论文正文草案。  
> 写作目标：统一符号；明确所有算法作用在 ciphertext 上；用 \(\vDash\) 描述诱导的 message semantics；让 no-\(B_k\)、coefficient-packed slots 和 core factorization theorem 形成一条主线。

## 2. Preliminaries

This section fixes the notation and the ciphertext-message semantics used throughout the paper. We intentionally separate standard CKKS and ring-switching background from the coefficient-split semantic observation developed in Section 3.

### 2.1 Rings and Ciphertext Semantics

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

to mean that \(\mathbf{ct}\) encrypts the message polynomial \(P(X)\in\mathcal R_d\) at scale \(\Delta\). Concretely, this means

\[
\operatorname{Dec}_{s(X)}(\mathbf{ct})
=
\Delta\cdot P(X)+e(X)
\quad\text{in }\mathcal R_{q,d},
\]

where \(e(X)\) is small enough for the intended CKKS precision.

Here \(\Delta\cdot P(X)\) denotes the scaled-and-rounded CKKS plaintext representative of \(P(X)\) in \(\mathcal R_{q,d}\). Equivalently, one may write

\[
\operatorname{Dec}_{s(X)}(\mathbf{ct})
=
\operatorname{Encode}_{\Delta,d}(P(X))+e(X)
\quad\text{in }\mathcal R_{q,d}.
\]

Throughout the ideal message-level discussion, we suppress the explicit encoding map and write \(\Delta\cdot P(X)\) for this representative.

When \(\Delta\), \(s(X)\), and the modulus are clear from context, we abbreviate the semantic judgment as

\[
\mathbf{ct}\vDash_d P(X).
\]

The judgment \(\vDash\) is semantic. It does not mean that \(P(X)\) is publicly available. All algorithms in this paper operate on ciphertexts; the notation records the induced transformation on the encrypted message.

### 2.2 Canonical-Slot and Coefficient Coordinates

Let

\[
\Omega_d=\{\omega\in\mathbb C:\omega^d+1=0\}
\]

be the root set of \(X^d+1\), with a fixed ordering. The full canonical embedding is

\[
\operatorname{can}_d:\mathcal R_d\to\mathsf{Can}_d,
\qquad
P(X)\mapsto (P(\omega))_{\omega\in\Omega_d}.
\]

Here \(\mathsf{Can}_d\) denotes the canonical-slot coordinate space.

If

\[
P(X)=\sum_{i=0}^{d-1}a_iX^i,
\]

then its coefficient coordinate vector is

\[
\operatorname{coeff}_d(P(X))
=
(a_0,\ldots,a_{d-1})\in\mathsf{Coeff}_d,
\]

where \(\mathsf{Coeff}_d\) denotes the coefficient-coordinate space. Both

\[
\operatorname{can}_d:\mathcal R_d\to\mathsf{Can}_d
\]

and

\[
\operatorname{coeff}_d:\mathcal R_d\to\mathsf{Coeff}_d
\]

are fixed linear coordinate isomorphisms.

In this paper, \(d\) denotes the full complex embedding dimension of \(\mathbb C[X]/(X^d+1)\), not necessarily the number of user-visible CKKS complex slots under the half-slot convention. The usual CKKS half-slot model can be recovered by restricting to the conjugate-compatible subspace and applying fixed projection/layout maps.

### 2.3 Slot Semantics for Ciphertexts

The slot semantic judgment

\[
\mathbf{ct}\vDash_d^{\mathrm{slot}}z
\]

means that \(\mathbf{ct}\) encrypts a polynomial \(P(X)\in\mathcal R_d\) whose canonical-slot vector is \(z\in\mathsf{Can}_d\). Formally,

\[
\mathbf{ct}\vDash_d^{\mathrm{slot}}z
\quad\Longleftrightarrow\quad
\mathbf{ct}\vDash_d P(X)
\text{ and }
z=\operatorname{can}_d(P(X)).
\]

Equivalently,

\[
\mathbf{ct}\vDash_d^{\mathrm{slot}}z
\quad\Longleftrightarrow\quad
\mathbf{ct}\vDash_d\operatorname{can}_d^{-1}(z).
\]

Thus slots are not separate plaintext objects. They are coordinates of the encrypted message polynomial under the canonical embedding.

### 2.4 CoeffToSlot, SlotToCoeff, and Coefficient-Packed Slots

We write \(\mathsf{C2S}_d\) for CoeffToSlot and \(\mathsf{S2C}_d\) for SlotToCoeff. In the normalized message-level model,

\[
\mathsf{C2S}_d
=
\operatorname{coeff}_d\circ\operatorname{can}_d^{-1},
\qquad
\mathsf{S2C}_d
=
\operatorname{can}_d\circ\operatorname{coeff}_d^{-1}.
\]

At the ciphertext level, \(\mathsf{C2S}_d\) is an encrypted linear transform with semantic effect

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

This motivates the central terminology.

#### Definition 2.1: Coefficient-Packed Slots

A ciphertext \(\mathbf{ct}'\) carries the coefficient-packed slots of \(P(X)\in\mathcal R_d\) if

\[
\mathbf{ct}'\vDash_d^{\mathrm{slot}}
\operatorname{coeff}_d(P(X)).
\]

Equivalently, \(\mathbf{ct}'\) encrypts

\[
\operatorname{can}_d^{-1}(\operatorname{coeff}_d(P(X))).
\]

Thus coefficient-packed slots are still encrypted slots. The phrase means that the canonical slots of the current encrypted polynomial contain the coefficient vector of another polynomial \(P(X)\).

Conversely, \(\mathsf{S2C}_d\) has semantic effect

\[
\mathbf{ct}'\vDash_d^{\mathrm{slot}}y
\quad\Longrightarrow\quad
\mathsf{S2C}_d(\mathbf{ct}')
\vDash_d
\operatorname{coeff}_d^{-1}(y),
\]

where \(y\in\mathsf{Coeff}_d\) is interpreted as a coefficient vector. In particular, if

\[
\mathbf{ct}'\vDash_d^{\mathrm{slot}}\operatorname{coeff}_d(Q(X)),
\]

then

\[
\mathsf{S2C}_d(\mathbf{ct}')\vDash_d Q(X).
\]

### 2.5 The CKKS Bootstrapping Core

Ignoring modulus raising and implementation errors, the conventional CKKS bootstrapping core has the form

\[
\mathsf{C2S}_d
\quad\longrightarrow\quad
\mathsf{EvalMod}_d
\quad\longrightarrow\quad
\mathsf{S2C}_d.
\]

Let \(f:\mathbb C\to\mathbb C\) be the scalar EvalMod approximation under a fixed normalization. In the ideal model,

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
\mathsf{C2S}_d(\mathbf{ct})\vDash_d^{\mathrm{slot}}(a_0,\ldots,a_{d-1}),
\]

then

\[
\mathsf{EvalMod}_d(\mathsf{C2S}_d(\mathbf{ct}))
\vDash_d^{\mathrm{slot}}
(f(a_0),\ldots,f(a_{d-1})).
\]

The important point is semantic: EvalMod is component-wise on the slots after \(\mathsf{C2S}_d\). Since those slots are coefficient-packed slots, EvalMod acts on coefficient coordinates, not directly on the original canonical slots of \(P(X)\).

### 2.6 Subring Decomposition and Layout

Let \(N=kn\). We use

\[
\mathcal R_N=\mathbb C[X]/(X^N+1),
\qquad
\mathcal R_n=\mathbb C[Y]/(Y^n+1).
\]

We identify \(\mathcal R_n\) with the subalgebra of \(\mathcal R_N\) generated by \(X^k\), via \(Y\mapsto X^k\). This is well-defined because

\[
(X^k)^n=X^N=-1
\]

in \(\mathcal R_N\), and it is injective because

\[
1,X^k,X^{2k},\ldots,X^{(n-1)k}
\]

are linearly independent in the monomial basis of \(\mathcal R_N\).

Every \(P(X)\in\mathcal R_N\) has a unique decomposition

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

For a coefficient vector

\[
p=\operatorname{coeff}_N(P(X))=(a_0,\ldots,a_{N-1}),
\]

define

\[
S_rp=(a_r,a_{r+k},a_{r+2k},\ldots,a_{r+(n-1)k})
\in\mathsf{Coeff}_n
\]

and

\[
\Pi_k:\mathsf{Coeff}_N\to\bigoplus_{r=0}^{k-1}\mathsf{Coeff}_n,
\qquad
\Pi_kp=(S_0p,\ldots,S_{k-1}p).
\]

In the normalized theorem statements, \(\Pi_k\) is residue-class grouping, possibly followed by fixed coordinate permutations. It does not include arbitrary diagonal scalings. Any implementation-specific scaling must be normalized into \(\mathsf{C2S}\), \(\mathsf{S2C}\), or the scalar EvalMod map \(f\), since nonlinear component-wise maps do not commute with arbitrary diagonal rescalings.

### 2.7 Notation Summary

The following notation is used consistently in Sections 2-4.

| Symbol | Meaning |
|---|---|
| \(\mathbf{ct}\vDash_d P(X)\) | ciphertext \(\mathbf{ct}\) encrypts message polynomial \(P(X)\in\mathcal R_d\), up to scale and error |
| \(\mathbf{ct}\vDash_d^{\mathrm{slot}}z\) | ciphertext \(\mathbf{ct}\) encrypts a polynomial whose canonical-slot vector is \(z\) |
| \(\mathsf{Can}_d\) | canonical-slot coordinate space |
| \(\mathsf{Coeff}_d\) | coefficient-coordinate space |
| \(\operatorname{can}_d(P(X))\) | canonical-slot vector of \(P(X)\) |
| \(\operatorname{coeff}_d(P(X))\) | coefficient vector of \(P(X)\) |
| \(\mathsf{C2S}_d\) | CoeffToSlot; as a coordinate map, \(\mathsf{Can}_d\to\mathsf{Coeff}_d\), placing message coefficients into ciphertext slots |
| \(\mathsf{S2C}_d\) | SlotToCoeff; as a coordinate map, \(\mathsf{Coeff}_d\to\mathsf{Can}_d\), interpreting slot values as coefficients |
| \(\Pi_k\) | residue-class grouping of coefficient coordinates |
| \(\mathsf{Split}_k\) | ciphertext split algorithm, or its induced message-coordinate split when the type is clear |
| \(\mathsf{Merge}_k\) | ciphertext merge algorithm, or its induced message-coordinate merge when the type is clear |
| \(B_k\) | optional slot-recovery transform from coefficient-split leaf slots to a slot-preserving leaf representation of big-ring slots |

## 3. Coefficient-Split Semantics

Section 2 defined the standard coordinates and ciphertext semantics. This section explains the semantic point specific to this paper: coefficient split is not slot-preserving, but slot preservation is not the invariant needed by the CKKS bootstrapping core.

### 3.1 Ciphertext Split and Message Meaning

The coefficient split of \(P(X)\in\mathcal R_N\) is

\[
\operatorname{CoeffSplit}_k(P(X))
=
(P_0(Y),\ldots,P_{k-1}(Y)),
\]

where

\[
P(X)=\sum_{r=0}^{k-1}X^rP_r(X^k).
\]

At the ciphertext level, \(\mathsf{Split}_k\) is not public coefficient extraction. It is a ciphertext operation, usually implemented through ring switching and key switching. Its intended semantic specification is

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

The inverse semantic operation is merge:

\[
\mathbf{ct}_{n,r}\vDash_n P_r(Y)
\quad (r=0,\ldots,k-1)
\]

implies

\[
\mathsf{Merge}_k(\mathbf{ct}_{n,0},\ldots,\mathbf{ct}_{n,k-1})
\vDash_N
\sum_{r=0}^{k-1}X^rP_r(X^k).
\]

The proof in Section 4 uses these as ideal message-coordinate maps. A ciphertext-level theorem must additionally bound the key-switching and rounding errors of the concrete \(\mathsf{Split}_k\) and \(\mathsf{Merge}_k\) implementation.

### 3.2 Coefficient Split Is Not Slot Selection

Coefficient split does not generally preserve the original big-ring canonical slots.

For \(k=2\), write

\[
P(X)=A(X^2)+XB(X^2),
\qquad
A(Y),B(Y)\in\mathcal R_n.
\]

Let \(\alpha\in\Omega_N\) and \(\theta=\alpha^2\in\Omega_n\). Then

\[
P(\alpha)=A(\theta)+\alpha B(\theta),
\qquad
P(-\alpha)=A(\theta)-\alpha B(\theta).
\]

Thus

\[
A(\theta)=\frac{P(\alpha)+P(-\alpha)}{2},
\qquad
B(\theta)=\frac{P(\alpha)-P(-\alpha)}{2\alpha}.
\]

The leaf slot values \(A(\theta)\) and \(B(\theta)\) are local linear combinations of paired big-ring slots, not selected original slots.

For general \(k\), the fiber above \(\theta\in\Omega_n\) is

\[
\pi^{-1}(\theta)=\{\alpha\in\Omega_N:\alpha^k=\theta\}.
\]

For every \(\alpha\in\pi^{-1}(\theta)\),

\[
P(\alpha)=\sum_{r=0}^{k-1}\alpha^rP_r(\theta).
\]

Since the elements of \(\pi^{-1}(\theta)\) are distinct, the corresponding Vandermonde matrix is nonsingular. Thus the relation above is a genuine change of basis inside each fiber.

Hence coefficient split changes the slot basis inside each fiber through an invertible Vandermonde relation. It is a coefficient/module decomposition, not a slot-preserving split.

### 3.3 Slot-Recovery and the Transform \(B_k\)

A split is slot-preserving if the leaf ciphertexts carry canonical-slot values that are selected from, or are a fixed permutation of, the original big-ring canonical-slot vector.

Equivalently, suppose a split maps a message \(P(X)\in\mathcal R_N\) to leaf messages \(Q_0(Y),\ldots,Q_{k-1}(Y)\in\mathcal R_n\). The split is slot-preserving if there exists a fixed bijection

\[
\sigma:\{0,\ldots,N-1\}
\to
\{0,\ldots,k-1\}\times\{0,\ldots,n-1\}
\]

such that, for every \(P(X)\in\mathcal R_N\), if \(\sigma(i)=(r,j)\), then

\[
(\operatorname{can}_n(Q_r(Y)))_j
=
(\operatorname{can}_N(P(X)))_i.
\]

For example, if \(N=4\), \(k=2\), \(n=2\), and

\[
\operatorname{can}_4(P(X))=(z_0,z_1,z_2,z_3),
\]

then a slot-preserving split may output \(Q_0(Y),Q_1(Y)\in\mathcal R_2\) with

\[
\operatorname{can}_2(Q_0(Y))=(z_0,z_2),
\qquad
\operatorname{can}_2(Q_1(Y))=(z_1,z_3).
\]

The leaf slots are then just a fixed regrouping of the original big-ring slots.

Coefficient split is not slot-preserving. If one wants to recover the original big-ring slots from coefficient-split leaf slots, one needs an additional slot-recovery transform.

Choose an ordering

\[
\pi^{-1}(\theta)=\{\alpha_0(\theta),\ldots,\alpha_{k-1}(\theta)\}
\]

for every \(\theta\in\Omega_n\). Given leaf slot functions \(u_r(\theta)\), define

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

If \(u_r(\theta)=P_r(\theta)\), then

\[
(B_k(u_0,\ldots,u_{k-1}))_t(\theta)
=
P(\alpha_t(\theta)).
\]

Thus \(B_k\) converts coefficient-split leaf slots into a slot-preserving leaf representation of the original big-ring slots. Its output is still organized as \(k\) small-ring slot vectors; it is not yet the single large-ring canonical vector \(\operatorname{can}_N(P(X))\) unless one additionally fixes an ordering and flattens the leaves.

In the case \(k=2\), this corresponds to

\[
Q^+(Y)=A(Y)+h(Y)B(Y),
\qquad
Q^-(Y)=A(Y)-h(Y)B(Y),
\]

where \(h(\theta)=\alpha\) for a chosen lift \(\alpha\) satisfying \(\alpha^2=\theta\). Such an \(h(Y)\) exists over \(\mathbb C\) by interpolation on the finite root set, after choosing one lift above each \(\theta\). It is branch-dependent and noncanonical.

The relevance of \(B_k\) is conceptual. It formalizes the extra multi-ciphertext linear transformation required if one insists on recovering the original big-ring slots after coefficient split. The main observation of this paper is that CKKS bootstrapping does not require this recovery.

### 3.4 Coefficient-Packed Slots Are the Correct Invariant

If

\[
\mathbf{ct}_N\vDash_N P(X),
\]

then the large-ring CoeffToSlot operation gives

\[
\mathsf{C2S}_N(\mathbf{ct}_N)
\vDash_N^{\mathrm{slot}}
\operatorname{coeff}_N(P(X)).
\]

Thus the input to EvalMod is not the original canonical-slot vector \(\operatorname{can}_N(P(X))\). It is the coefficient-packed slot vector \(\operatorname{coeff}_N(P(X))\).

Now split first. If

\[
\mathsf{Split}_k(\mathbf{ct}_N)
=
(\mathbf{ct}_{n,0},\ldots,\mathbf{ct}_{n,k-1})
\]

and

\[
\mathbf{ct}_{n,r}\vDash_n P_r(Y),
\]

then leaf CoeffToSlot gives

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
(a_r,a_{r+k},a_{r+2k},\ldots,a_{r+(n-1)k}).
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

This is the key semantic point. The invariant required by CKKS bootstrapping is not preservation of the original big-ring canonical slots. The invariant is preservation of the \(\mathsf{C2S}\)-output coefficient-packed coordinates, up to residue-class grouping.

### 3.5 No-\(B_k\) Route

The slot-recovery transform \(B_k\) solves

\[
\text{coefficient-split leaf canonical slots}
\longrightarrow
\text{slot-preserving leaf representation of big-ring slots}.
\]

The output leaves of \(B_k\) contain the original big-ring slot values, but they are still organized as \(k\) small-ring slot vectors. This distinction keeps \(B_k\) separate from \(\mathsf{Merge}_k\), whose role is to assemble message polynomials across rings.

The bootstrapping core does not need to solve this problem. It follows

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

After \(\mathsf{Split}_k\), the leaf ciphertexts encrypt \(P_r(Y)\), whose leaf canonical slots are not original big-ring slots. But after leaf \(\mathsf{C2S}_n\), they carry exactly the residue-class coefficient sub-vectors of \(P(X)\). Since EvalMod is component-wise on coefficient-packed slots, it can be evaluated independently on each leaf. Leaf \(\mathsf{S2C}_n\) then returns each reduced coefficient vector to leaf canonical form, and \(\mathsf{Merge}_k\) assembles the corresponding large-ring message.

Thus the no-\(B_k\) route does not claim that coefficient split preserves big-ring slots. It claims the precise invariant needed for the core:

\[
\text{coefficient split preserves the C2S-output coefficient coordinates up to residue grouping.}
\]

## 4. Core Factorization Theorem

This section proves the algebraic correctness of the no-\(B_k\) route. The proof is stated at the ideal message-coordinate level using the notation of Sections 2 and 3. Ciphertext noise, key-switching error, rescaling error, and implementation-specific normalization are deferred to the parameter and implementation sections.

### 4.1 Abstract Split-Compatible Sandwich Theorem

We first prove an abstract theorem for a split-compatible sandwich

\[
\text{linear map}
\quad\to\quad
\text{component-wise map}
\quad\to\quad
\text{linear map}.
\]

#### Theorem 4.1: Split-Compatible Sandwich Theorem

Let \(V_N,W_N,V_n,W_n\) be vector spaces. Let

\[
L_N^{(1)}:V_N\to W_N,
\qquad
L_n^{(1)}:V_n\to W_n,
\]

and

\[
L_N^{(2)}:W_N\to V_N,
\qquad
L_n^{(2)}:W_n\to V_n
\]

be linear maps. Let

\[
\Phi_N:W_N\to W_N,
\qquad
\Phi_n:W_n\to W_n
\]

be arbitrary maps, not necessarily linear.

For possibly nonlinear \(\Phi_n\), the notation

\[
\bigoplus_{r=0}^{k-1}\Phi_n
\]

means

\[
(w_0,\ldots,w_{k-1})
\mapsto
(\Phi_n(w_0),\ldots,\Phi_n(w_{k-1})).
\]

Assume there exist maps

\[
\mathsf{Split}:V_N\to\bigoplus_{r=0}^{k-1}V_n,
\qquad
\mathsf{Merge}:\bigoplus_{r=0}^{k-1}V_n\to V_N,
\]

and

\[
\Pi:W_N\to\bigoplus_{r=0}^{k-1}W_n
\]

such that

\[
\left(\bigoplus_{r=0}^{k-1}L_n^{(1)}\right)\circ\mathsf{Split}
=
\Pi\circ L_N^{(1)},
\tag{1}
\]

\[
\left(\bigoplus_{r=0}^{k-1}\Phi_n\right)\circ\Pi
=
\Pi\circ\Phi_N,
\tag{2}
\]

and

\[
\mathsf{Merge}\circ
\left(\bigoplus_{r=0}^{k-1}L_n^{(2)}\right)
\circ\Pi
=
L_N^{(2)}.
\tag{3}
\]

Define

\[
F_N=L_N^{(2)}\circ\Phi_N\circ L_N^{(1)}
\]

and

\[
F_n=L_n^{(2)}\circ\Phi_n\circ L_n^{(1)}.
\]

Then

\[
F_N
=
\mathsf{Merge}\circ
\left(\bigoplus_{r=0}^{k-1}F_n\right)
\circ
\mathsf{Split}.
\]

#### Proof

Starting from the right-hand side,

\[
\begin{aligned}
&\mathsf{Merge}\circ
\left(\bigoplus_{r=0}^{k-1}F_n\right)
\circ\mathsf{Split}
\\
&=
\mathsf{Merge}\circ
\left(\bigoplus_{r=0}^{k-1}
L_n^{(2)}\circ\Phi_n\circ L_n^{(1)}
\right)
\circ\mathsf{Split}
\\
&=
\mathsf{Merge}\circ
\left(\bigoplus_{r=0}^{k-1}L_n^{(2)}\right)
\circ
\left(\bigoplus_{r=0}^{k-1}\Phi_n\right)
\circ
\left(\bigoplus_{r=0}^{k-1}L_n^{(1)}\right)
\circ\mathsf{Split}.
\end{aligned}
\]

Using (1), this becomes

\[
\mathsf{Merge}\circ
\left(\bigoplus_{r=0}^{k-1}L_n^{(2)}\right)
\circ
\left(\bigoplus_{r=0}^{k-1}\Phi_n\right)
\circ
\Pi
\circ
L_N^{(1)}.
\]

Using (2), we get

\[
\mathsf{Merge}\circ
\left(\bigoplus_{r=0}^{k-1}L_n^{(2)}\right)
\circ
\Pi
\circ
\Phi_N
\circ
L_N^{(1)}.
\]

Using (3), this equals

\[
L_N^{(2)}\circ\Phi_N\circ L_N^{(1)}
=
F_N.
\]

This proves the theorem. \(\square\)

The theorem does not require \(\Phi_N\) or \(\Phi_n\) to be linear. The only nonlinear requirement is compatibility (2), which will follow from component-wise EvalMod on coefficient-packed slots.

### 4.2 Component-Wise Compatibility

#### Lemma 4.2: Component-Wise Maps Commute with Residue Grouping

Let

\[
\Phi_N=f^{\oplus N}:\mathbb C^N\to\mathbb C^N
\]

and

\[
\Phi_n=f^{\oplus n}:\mathbb C^n\to\mathbb C^n
\]

for the same scalar function \(f:\mathbb C\to\mathbb C\). If \(\Pi_k\) is residue-class grouping up to fixed coordinate permutations, then

\[
\left(\bigoplus_{r=0}^{k-1}\Phi_n\right)\circ\Pi_k
=
\Pi_k\circ\Phi_N.
\]

#### Proof

Let \(x=(x_0,\ldots,x_{N-1})\). The coordinate in leaf \(r\) and position \(j\) is \(x_{r+jk}\). Applying the leaf map gives \(f(x_{r+jk})\). Applying \(f^{\oplus N}\) first and then grouping also gives \(f(x_{r+jk})\). If \(\Pi_k\) includes fixed coordinate permutations, the same argument holds after reindexing, since \(f^{\oplus d}\) commutes with every coordinate permutation. All coordinates agree. \(\square\)

This lemma uses neither additivity nor multiplicativity of \(f\). It uses only that the same scalar function is applied independently to each coordinate, and that \(\Pi_k\) contains no unnormalized coordinate-dependent scaling.

### 4.3 Instantiation to the CKKS Core

The large-ring bootstrapping core is

\[
\mathsf{BTS}^{\mathrm{core}}_N
=
\mathsf{S2C}_N\circ
\mathsf{EvalMod}_N\circ
\mathsf{C2S}_N.
\]

The leaf core is

\[
\mathsf{BTS}^{\mathrm{core}}_n
=
\mathsf{S2C}_n\circ
\mathsf{EvalMod}_n\circ
\mathsf{C2S}_n.
\]

In Theorem 4.1 we instantiate

\[
V_N=\mathsf{Can}_N,\quad W_N=\mathsf{Coeff}_N,
\qquad
V_n=\mathsf{Can}_n,\quad W_n=\mathsf{Coeff}_n,
\]

\[
L_N^{(1)}=\mathsf{C2S}_N,
\qquad
L_n^{(1)}=\mathsf{C2S}_n,
\]

\[
L_N^{(2)}=\mathsf{S2C}_N,
\qquad
L_n^{(2)}=\mathsf{S2C}_n,
\]

and

\[
\Phi_N=\mathsf{EvalMod}_N,
\qquad
\Phi_n=\mathsf{EvalMod}_n.
\]

The split, merge, and layout maps are \(\mathsf{Split}_k\), \(\mathsf{Merge}_k\), and \(\Pi_k\).

### 4.4 Split-C2S Compatibility

#### Lemma 4.3: Split-C2S Compatibility

Under the normalized layout convention,

\[
\left(\bigoplus_{r=0}^{k-1}\mathsf{C2S}_n\right)
\circ
\mathsf{Split}_k
=
\Pi_k\circ\mathsf{C2S}_N.
\]

#### Proof

Let

\[
P(X)=\sum_{i=0}^{N-1}a_iX^i\in\mathcal R_N,
\]

and write

\[
z=\operatorname{can}_N(P(X)),
\qquad
p=\operatorname{coeff}_N(P(X))=(a_0,\ldots,a_{N-1}).
\]

By coefficient split,

\[
P(X)=\sum_{r=0}^{k-1}X^rP_r(X^k),
\qquad
P_r(Y)=\sum_{j=0}^{n-1}a_{r+jk}Y^j.
\]

Therefore

\[
\operatorname{coeff}_n(P_r(Y))=S_rp.
\]

The left-hand side maps

\[
z
\mapsto
(\operatorname{can}_n(P_0(Y)),\ldots,\operatorname{can}_n(P_{k-1}(Y)))
\]

and then applies leaf \(\mathsf{C2S}_n\), giving

\[
(\operatorname{coeff}_n(P_0(Y)),\ldots,\operatorname{coeff}_n(P_{k-1}(Y)))
=
(S_0p,\ldots,S_{k-1}p)
=
\Pi_kp.
\]

The right-hand side first maps \(z\) by \(\mathsf{C2S}_N\) to \(p\), and then applies \(\Pi_k\), also giving \(\Pi_kp\). Since \(\operatorname{can}_N\) is a linear coordinate isomorphism, every \(z\in\mathsf{Can}_N\) is uniquely of the form \(\operatorname{can}_N(P(X))\). Hence the maps are equal. \(\square\)

In ciphertext terms, this lemma formalizes the semantic statement from Section 3. If

\[
\mathbf{ct}_N\vDash_N P(X)
\]

and

\[
\mathsf{Split}_k(\mathbf{ct}_N)
=
(\mathbf{ct}_{n,0},\ldots,\mathbf{ct}_{n,k-1}),
\qquad
\mathbf{ct}_{n,r}\vDash_n P_r(Y),
\]

then

\[
\bigoplus_r\mathsf{C2S}_n(\mathbf{ct}_{n,r})
\]

carries exactly the residue-class grouping of the coefficient-packed slots carried by \(\mathsf{C2S}_N(\mathbf{ct}_N)\).

### 4.5 EvalMod Compatibility

#### Lemma 4.4: EvalMod Compatibility

Assume that, after normalized \(\mathsf{C2S}\) output scaling,

\[
\mathsf{EvalMod}_N=f^{\oplus N},
\qquad
\mathsf{EvalMod}_n=f^{\oplus n}
\]

for the same scalar function \(f:\mathbb C\to\mathbb C\), with the same approximation polynomial, period normalization, input scaling, and output scaling. Then

\[
\left(\bigoplus_{r=0}^{k-1}\mathsf{EvalMod}_n\right)
\circ
\Pi_k
=
\Pi_k\circ\mathsf{EvalMod}_N.
\]

#### Proof

This is Lemma 4.2 with \(\Phi_N=\mathsf{EvalMod}_N\) and \(\Phi_n=\mathsf{EvalMod}_n\). \(\square\)

The point is that EvalMod may be nonlinear. Nonlinearity is harmless here because it is applied component-wise to coefficient-packed slots.

### 4.6 S2C-Merge Compatibility

#### Lemma 4.5: S2C-Merge Compatibility

Under the normalized layout convention,

\[
\mathsf{Merge}_k
\circ
\left(\bigoplus_{r=0}^{k-1}\mathsf{S2C}_n\right)
\circ
\Pi_k
=
\mathsf{S2C}_N.
\]

#### Proof

Let

\[
y=(y_0,\ldots,y_{N-1})\in\mathsf{Coeff}_N.
\]

After \(\Pi_k\), the \(r\)-th leaf receives

\[
S_ry=(y_r,y_{r+k},\ldots,y_{r+(n-1)k}).
\]

Interpret this as the coefficient vector of

\[
Y_r(Y)=\sum_{j=0}^{n-1}y_{r+jk}Y^j.
\]

Applying \(\mathsf{S2C}_n\) gives \(\operatorname{can}_n(Y_r(Y))\). Merging the leaf messages gives the large-ring polynomial

\[
Y(X)=\sum_{r=0}^{k-1}X^rY_r(X^k).
\]

Expanding,

\[
Y(X)
=
\sum_{r=0}^{k-1}\sum_{j=0}^{n-1}y_{r+jk}X^{r+jk}
=
\sum_{i=0}^{N-1}y_iX^i.
\]

Therefore the leaf route outputs the canonical-slot vector of the large-ring polynomial whose coefficient vector is \(y\), which is exactly \(\mathsf{S2C}_N(y)\). \(\square\)

### 4.7 Ideal Normalized CKKS Core Factorization

#### Theorem 4.6: Ideal Normalized CKKS Core Factorization under Coefficient-Split Semantics

Assume the normalized message-level model of Section 2:

1. \(\operatorname{can}_d\) and \(\operatorname{coeff}_d\) are fixed linear coordinate isomorphisms for \(d\in\{N,n\}\).

2. \(\mathsf{C2S}_d=\operatorname{coeff}_d\circ\operatorname{can}_d^{-1}\) and \(\mathsf{S2C}_d=\operatorname{can}_d\circ\operatorname{coeff}_d^{-1}\).

3. \(\mathsf{Split}_k\) and \(\mathsf{Merge}_k\) are induced by \(P(X)=\sum_{r=0}^{k-1}X^rP_r(X^k)\).

4. \(\Pi_k\) is residue-class grouping, up to fixed coordinate permutations only, and contains no unnormalized diagonal scaling.

5. The large-ring and leaf-ring EvalMod stages apply the same component-wise scalar map under the same normalization:

\[
\mathsf{EvalMod}_N=f^{\oplus N},
\qquad
\mathsf{EvalMod}_n=f^{\oplus n}.
\]

Then

\[
\mathsf{S2C}_N\circ
\mathsf{EvalMod}_N\circ
\mathsf{C2S}_N
=
\mathsf{Merge}_k
\circ
\left(
\bigoplus_{r=0}^{k-1}
\mathsf{S2C}_n\circ
\mathsf{EvalMod}_n\circ
\mathsf{C2S}_n
\right)
\circ
\mathsf{Split}_k.
\]

#### Proof

The three compatibility conditions of Theorem 4.1 are exactly Lemmas 4.3, 4.4, and 4.5:

\[
\left(\bigoplus_r\mathsf{C2S}_n\right)\circ\mathsf{Split}_k
=
\Pi_k\circ\mathsf{C2S}_N,
\]

\[
\left(\bigoplus_r\mathsf{EvalMod}_n\right)\circ\Pi_k
=
\Pi_k\circ\mathsf{EvalMod}_N,
\]

and

\[
\mathsf{Merge}_k
\circ
\left(\bigoplus_r\mathsf{S2C}_n\right)
\circ
\Pi_k
=
\mathsf{S2C}_N.
\]

Applying Theorem 4.1 with

\[
L^{(1)}=\mathsf{C2S},
\qquad
\Phi=\mathsf{EvalMod},
\qquad
L^{(2)}=\mathsf{S2C}
\]

gives the stated factorization. \(\square\)

### 4.8 Ciphertext-Level Reading

Theorem 4.6 is an ideal message-coordinate equality. Its ciphertext-level interpretation is:

\[
\mathbf{ct}_N\vDash_N P(X)
\]

and

\[
\mathsf{Split}_k(\mathbf{ct}_N)
=
(\mathbf{ct}_{n,0},\ldots,\mathbf{ct}_{n,k-1}),
\qquad
\mathbf{ct}_{n,r}\vDash_n P_r(Y),
\]

then the leaf route

\[
\mathsf{Split}_k
\to
\bigoplus_r\mathsf{C2S}_n
\to
\bigoplus_r\mathsf{EvalMod}_n
\to
\bigoplus_r\mathsf{S2C}_n
\to
\mathsf{Merge}_k
\]

has the same ideal message effect as

\[
\mathsf{C2S}_N
\to
\mathsf{EvalMod}_N
\to
\mathsf{S2C}_N.
\]

The equality is not a statement that \(\mathsf{Split}_k\) preserves big-ring slots. It is a statement that, after \(\mathsf{C2S}\), the relevant encrypted slots are coefficient-packed slots, and those coordinates decompose by residue class. This is the formal reason the core can be factored without the \(B_k\) slot-recovery transform.

A full ciphertext-level correctness theorem must separately show that the concrete \(\mathsf{RSDown}_{N\to n}\) and \(\mathsf{RSUp}_{n\to N}\) operations implement \(\mathsf{Split}_k\) and \(\mathsf{Merge}_k\) up to controlled error and the same normalization. It must also account for key-switching error, EvalMod approximation error, C2S/S2C linear transform error, RNS rounding, rescaling, implementation layout, and security parameters for the top and leaf rings.
