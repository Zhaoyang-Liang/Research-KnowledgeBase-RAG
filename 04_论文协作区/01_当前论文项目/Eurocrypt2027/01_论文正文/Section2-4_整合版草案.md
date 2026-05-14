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
\operatorname{Encode}_{\Delta,d}(P(X))+e(X)
\quad\text{in }\mathcal R_{q,d},
\]

where \(e(X)\) is small enough for the intended CKKS precision.

Here \(\operatorname{Encode}_{\Delta,d}(P(X))\) denotes the scaled-and-rounded CKKS plaintext representative of \(P(X)\) in \(\mathcal R_{q,d}\). Throughout the ideal message-level discussion, we suppress this explicit encoding map and write \(\Delta\cdot P(X)\) as shorthand for the same representative when no confusion can arise.

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

This component-wise behavior does not require the algorithm to decrypt slots or to encrypt the function \(f\). The function \(f\) is public. In practice, EvalMod evaluates a public polynomial approximation

\[
p(t)=\sum_{\ell=0}^{m}c_\ell t^\ell
\]

on an encrypted ring element. If a ciphertext encrypts \(Q(X)\in\mathcal R_d\), then the homomorphic computation produces an encryption of

\[
p(Q(X))=\sum_{\ell=0}^{m}c_\ell Q(X)^\ell.
\]

The canonical embedding is compatible with ring addition and multiplication:

\[
\operatorname{can}_d(A(X)+B(X))
=
\operatorname{can}_d(A(X))+\operatorname{can}_d(B(X)),
\]

\[
\operatorname{can}_d(A(X)B(X))
=
\operatorname{can}_d(A(X))\odot\operatorname{can}_d(B(X)),
\]

where \(\odot\) denotes component-wise multiplication. Therefore, for \(z=\operatorname{can}_d(Q(X))=(z_0,\ldots,z_{d-1})\),

\[
\operatorname{can}_d(p(Q(X)))
=
(p(z_0),\ldots,p(z_{d-1}))
=
p^{\oplus d}(z).
\]

Thus \(f^{\oplus d}\) is the slot-coordinate semantics of evaluating a public scalar approximation on one encrypted polynomial. The theorem later assumes that the same normalized scalar map is used in the large ring and in every leaf ring.

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

It is often useful to view \(p\) as an \(n\times k\) coefficient rectangle whose row \(j\) and column \(r\) entry is

\[
a_{r+jk}.
\]

Then \(S_rp\) is exactly the \(r\)-th column of this rectangle, and \(P_r(Y)\) is the small-ring polynomial collecting all coefficients of \(P(X)\) whose indices are congruent to \(r\) modulo \(k\).

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
| \(\mathsf{C2S}_d\) | CoeffToSlot; as a coordinate map, \(\mathsf{Can}_d\to\mathsf{Coeff}_d\), making the canonical slots of the output ciphertext carry the coefficient vector of the input message polynomial |
| \(\mathsf{S2C}_d\) | SlotToCoeff; as a coordinate map, \(\mathsf{Coeff}_d\to\mathsf{Can}_d\), interpreting slot values as coefficients |
| \(f^{\oplus d}\) | component-wise lift of a public scalar EvalMod approximation to \(d\) slot coordinates |
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

Equivalently, the map

\[
\pi:\Omega_N\to\Omega_n,
\qquad
\pi(\alpha)=\alpha^k
\]

organizes the \(N\) big-ring roots into \(n\) fibers, each of size \(k\). The small-ring slot \(\theta\) is therefore the base point of a local fiber of big-ring slots, not one selected big-ring slot.

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

fiber-wise as follows. Here the subscript \(t\) indexes the output leaf, and \(\theta\in\Omega_n\) indexes the slot inside that leaf:

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

We now formalize this semantic observation as an algebraic factorization theorem.

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

The theorem does not require \(\Phi_N\) or \(\Phi_n\) to be linear. The only nonlinear requirement is compatibility (2), which will follow from component-wise EvalMod on coefficient-packed slots. We now verify conditions (1)-(3) for the CKKS instantiation in Lemmas 4.3-4.5 below, with Lemma 4.2 providing the component-wise compatibility used by EvalMod.

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

The right-hand side first maps \(z\) by \(\mathsf{C2S}_N\) to \(p\), and then applies \(\Pi_k\), also giving \(\Pi_kp\). Since \(\operatorname{can}_N\) is a bijection, the identity holds for all \(z\in\mathsf{Can}_N\). \(\square\)

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

#### Remark 4.7: What Remains Outside the Ideal Theorem

Theorem 4.6 is a message-coordinate equality. A full ciphertext-level correctness theorem must separately show that the concrete \(\mathsf{RSDown}_{N\to n}\) and \(\mathsf{RSUp}_{n\to N}\) operations implement \(\mathsf{Split}_k\) and \(\mathsf{Merge}_k\) up to controlled error and the same normalization.

The parameter and correctness analysis must account for at least the following gaps: \(\mathsf{RSDown}/\mathsf{RSUp}\) key-switching error; \(\mathsf{C2S}/\mathsf{S2C}\) linear transform approximation error; EvalMod polynomial approximation error; RNS rounding and rescaling error; implementation-specific layout or normalization factors; and the security parameters for both the top ring and the leaf rings.

### 4.9 Incorporating ScaleDown and ModRaise

The factorization theorem above isolates the algebraic bootstrapping core

\[
\mathsf{C2S}\to\mathsf{EvalMod}\to\mathsf{S2C}.
\]

Concrete CKKS bootstrapping also contains a preprocessing stage, usually consisting of scale normalization and modulus raising. We write this stage abstractly as

\[
\mathsf{Pre}_d
=
\mathsf{ModRaise}_d\circ\mathsf{ScaleDown}_d.
\]

At the message level, \(\mathsf{ScaleDown}_d\) applies a public normalization to the encrypted plaintext representative. Under a fixed bootstrapping parameter set, we model its semantic effect as

\[
\mathbf{ct}\vDash_{d,\Delta,Q_0}P
\quad\Longrightarrow\quad
\mathsf{ScaleDown}_d(\mathbf{ct})
\vDash_{d,\Delta',Q'_0}
\rho P,
\]

where \(\rho\) is the public bootstrapping normalization factor. In many implementations \(\rho\) is absorbed into the subsequent EvalMod normalization, so that the full bootstrapping map remains semantically the intended refresh map.

The modulus raising stage changes the ciphertext modulus from a low-level modulus \(Q'_0\) to a bootstrapping modulus \(Q_{\mathrm{boot}}\), without changing the encrypted message in the ideal model:

\[
\mathsf{ModRaise}_d:
\quad
\mathbf{ct}\vDash_{d,\Delta',Q'_0}\rho P
\Longrightarrow
\mathsf{ModRaise}_d(\mathbf{ct})
\vDash_{d,\Delta',Q_{\mathrm{boot}}}\rho P.
\]

Thus \(\mathsf{Pre}_d\) is message-coordinate linear up to the fixed scalar normalization \(\rho\). If the same normalization is used in the large ring and in each leaf ring, then \(\mathsf{Pre}\) commutes with coefficient split at the ideal message level:

\[
\left(\bigoplus_{r=0}^{k-1}\mathsf{Pre}_n\right)\circ\mathsf{Split}_k
=
\mathsf{Split}_k\circ\mathsf{Pre}_N.
\tag{4}
\]

Indeed, if

\[
P(X)=\sum_{r=0}^{k-1}X^rP_r(X^k),
\]

then

\[
\rho P(X)
=
\sum_{r=0}^{k-1}X^r(\rho P_r)(X^k),
\]

so applying a uniform normalization before or after the coefficient split gives the same leaf messages. Modulus raising does not alter the message polynomial in the ideal semantics, hence it also commutes with split at this level.

In ciphertext implementations, (4) becomes an approximate compatibility condition:

\[
\left(\bigoplus_r\mathsf{Pre}_n\right)(\mathsf{Split}_k(\mathbf{ct}))_r
\approx
\mathsf{Split}_k(\mathsf{Pre}_N(\mathbf{ct}))_r,
\]

where the discrepancy consists of scale-down rounding, RNS basis-extension error, ring-switching error, and key-switching error. Therefore, a concrete implementation must ensure that the accumulated error remains below the CKKS precision budget and that the EvalMod stage sees the same normalized scalar input in the large-ring and leaf-ring routes.

### 4.10 Full Bootstrapping Factorization with Leaf Preprocessing

Define the full ideal CKKS bootstrapping map

\[
\mathsf{BTS}_d
=
\mathsf{S2C}_d
\circ
\mathsf{EvalMod}_d
\circ
\mathsf{C2S}_d
\circ
\mathsf{Pre}_d.
\]

Assume the hypotheses of Theorem 4.6 and, additionally, the preprocessing compatibility (4). Then

\[
\mathsf{BTS}_N
=
\mathsf{Merge}_k
\circ
\left(
\bigoplus_{r=0}^{k-1}
\mathsf{BTS}_n
\right)
\circ
\mathsf{Split}_k.
\tag{5}
\]

Equivalently, the split-first bootstrapping route is

\[
\mathsf{Split}_k
\to
\bigoplus_r\mathsf{ScaleDown}_n
\to
\bigoplus_r\mathsf{ModRaise}_n
\to
\bigoplus_r\mathsf{C2S}_n
\to
\bigoplus_r\mathsf{EvalMod}_n
\to
\bigoplus_r\mathsf{S2C}_n
\to
\mathsf{Merge}_k.
\]

This is the most natural ciphertext-level route when the leaf engine is a standard CKKS bootstrapping implementation: \(\mathsf{ScaleDown}_n\) and \(\mathsf{ModRaise}_n\) are executed as a matched preprocessing pair inside each leaf ring. Applying only \(\mathsf{ModRaise}_n\) after a large-ring \(\mathsf{ScaleDown}_N\) may break the implementation-specific normalization expected by the leaf EvalMod stage.

Equation (5) should be read as a modular design principle. The bootstrapping core is one instance of the split-compatible sandwich theorem. More generally, any encrypted computation of the form

\[
L_d^{(2)}\circ\Phi_d\circ L_d^{(1)}\circ\mathsf{Pre}_d
\]

can be factored through coefficient split whenever the linear maps, the component-wise nonlinear map, the preprocessing normalization, and the merge map satisfy the compatibility conditions in Theorem 4.1 and (4).

## 5. Implementation, Parameter Regime, and Preliminary Evaluation

This section records the implementation interpretation of the preceding theorems and the current parameter regime. The experimental numbers below are preliminary and are included to document the present prototype; they are not intended to be the final experimental evaluation.

### 5.1 Implementation Strategy

The implementation viewpoint is deliberately modular. The coefficient-split theorem does not require a new EvalMod algorithm; it requires a split-compatible placement of the bootstrapping pipeline. We therefore treat the leaf bootstrapping engine as a replaceable module. The abstract route is

\[
\mathsf{Split}
\to
\bigoplus_r \mathsf{BTS}_n
\to
\mathsf{Merge}.
\]

For CKKS bootstrapping, the currently preferred route is the split-first leaf-preprocessing route:

\[
\mathsf{Split}
\to
\bigoplus_r\mathsf{ScaleDown}_n
\to
\bigoplus_r\mathsf{ModRaise}_n
\to
\bigoplus_r\mathsf{C2S}_n
\to
\bigoplus_r\mathsf{EvalMod}_n
\to
\bigoplus_r\mathsf{S2C}_n
\to
\mathsf{Merge}.
\]

This differs from the earlier engineering route

\[
\mathsf{ScaleDown}_N
\to
\mathsf{ModRaise}_N
\to
\mathsf{Split}
\to
\bigoplus_r\mathsf{C2S}_n
\to
\cdots,
\]

which leaves the preprocessing cost in the large ring. It also differs from the hybrid route

\[
\mathsf{ScaleDown}_N
\to
\mathsf{Split}
\to
\bigoplus_r\mathsf{ModRaise}_n
\to
\cdots,
\]

which may be algebraically reasonable but can mismatch the leaf bootstrapping engine's expected normalization. In the current prototype, moving both \(\mathsf{ScaleDown}\) and \(\mathsf{ModRaise}\) to the leaf ring substantially improves the output precision compared with moving only \(\mathsf{ModRaise}\). This matches the theoretical picture of Section 4.9: \(\mathsf{ScaleDown}\) fixes the normalized input seen by EvalMod, while \(\mathsf{ModRaise}\) provides the modulus extension needed to evaluate the bootstrapping polynomial.

The present implementation uses a standard Lattigo CKKS bootstrapping evaluator as the leaf engine. This is deliberately conservative. The factorization theorem is independent of a particular bootstrapping engine, so optimized leaf bootstrapping methods can be substituted later as long as they implement the same normalized scalar EvalMod map and satisfy the compatibility conditions above. In this sense, the split construction is an outer factorization layer: improvements to the leaf engine should compound with the split route rather than compete with it.

### 5.2 Target Workloads: Capacity-Bound, Moderate-Depth Computation

The factorization is not intended to make every bootstrapping parameter set faster. Its natural target is a capacity-bound and moderate-depth regime. Let \(S(d)\) denote the number of user-visible CKKS slots available at ring dimension \(d\), and let \(M\) be the number of scalar values that an application wants to pack in one ciphertext. The first condition is

\[
S(n)<M\le S(N),
\qquad N=kn,
\tag{6}
\]

so that the large ring is required for packing capacity, while each leaf ring is too small to hold the whole packed object by itself. The second condition is that the computation between refreshes fits within the secure leaf modulus budget:

\[
D_{\mathrm{cycle}}\le L(Q_{\mathrm{leaf}}),
\tag{7}
\]

where \(D_{\mathrm{cycle}}\) is the multiplicative depth in one refresh cycle and \(L(Q_{\mathrm{leaf}})\) is the number of levels supported by the leaf modulus chain. In this regime, the large ring is chosen for slot capacity rather than for supporting the largest possible modulus.

Several encrypted workloads have this shape.

First, in batched SIMD inference, a ciphertext may pack

\[
M=B\cdot F
\]

values, where \(B\) is the batch size and \(F\) is the number of features, channels, or tensor entries per sample. A larger ring may be needed because \(B\cdot F>S(n)\). However, the depth between activations or refresh points may remain modest, especially when the network uses low-degree polynomial activations. If

\[
\log(PQ)_{\mathrm{leaf}}\le B_{\mathrm{sec}}(n)
\]

where \(B_{\mathrm{sec}}(n)\) is the secure modulus budget at the leaf dimension, then the leaf bootstrapping route can exploit large-ring packing while using leaf-ring bootstrapping parameters. The point is not that inference requires a huge depth in each ciphertext; it is that the tensor layout can require a large number of slots.

Second, encrypted database analytics and columnar aggregation are often dominated by the number of packed records rather than by multiplicative depth. For a column with \(R\) records and \(a\) slots per record,

\[
M=R\cdot a.
\]

Linear statistics such as sums and means consume essentially no multiplicative levels, while quantities such as variance require only a small number of multiplications:

\[
\operatorname{var}(x)
=
\frac{1}{R}\sum_i x_i^2
-
\left(\frac{1}{R}\sum_i x_i\right)^2.
\]

Predicate approximations or smooth filters may use a low-degree polynomial

\[
\mathbf{1}_{x>t}\approx p_d(x),
\]

whose depth is roughly \(\lceil\log_2 d\rceil\) under balanced evaluation. These are examples where \(M\) can force a large ring while \(D_{\mathrm{cycle}}\) remains moderate.

Third, large-vector analytics and many-model workloads may pack a large number of independent coordinates, genomic markers, or small models. If \(G\) independent tasks each require \(p\) packed values, then

\[
M=G\cdot p.
\]

The ring dimension is driven by \(G\), while the depth of each task is driven by the local arithmetic. For example, a polynomial approximation to a nonlinear function contributes depth according to its degree, not according to \(G\).

Fourth, batched activation or polynomial-refresh stages in encrypted neural networks are particularly close to the bootstrapping setting. If an activation is applied to \(M\) slots as

\[
y_i=f(x_i),
\qquad i=1,\ldots,M,
\]

and is evaluated by a polynomial approximation

\[
f(x)\approx p_d(x)=\sum_{\ell=0}^d c_\ell x^\ell,
\]

then the ring dimension is determined by \(M\), whereas the depth is determined by the evaluation strategy for \(p_d\). The split route is well matched when \(M\) requires the large ring but the activation/refresh cycle fits within the leaf modulus budget.

These examples also clarify what the method does not target. If the workload is depth-bound and genuinely needs the large-ring modulus budget, then forcing the computation into the smaller leaf modulus may increase the number of required bootstraps and may not be beneficial. The intended application argument is therefore not "smaller \(Q\) is always better." It is: when \(N\) is forced by capacity and \(Q_{\mathrm{leaf}}\) is already sufficient for the per-refresh computation, bootstrapping in leaf rings can reduce the dominant bootstrapping cost without changing the application layout.

### 5.3 Security Accounting

Once ciphertexts and evaluation keys exist in both the top ring and the leaf rings, the security level is determined by the weakest exposed RLWE instance:

\[
\lambda_{\mathrm{scheme}}
=
\min\{
\lambda(N,Q_{\mathrm{top}},P_{\mathrm{top}},\chi_s,\chi_e),
\lambda(n,Q_{\mathrm{leaf}},P_{\mathrm{leaf}},\chi_s,\chi_e),
\lambda_{\mathrm{rs/ks}}
\}.
\]

Therefore, the large ring's security margin cannot be used to hide an insecure leaf ring. If the split uses \(N=kn\), then the leaf modulus budget must be checked against the leaf dimension \(n\), not the top dimension \(N\).

This point is important for interpreting the apparent surplus of the large ring. A top-ring parameter such as \(N=2^{16}\) or \(N=2^{17}\) may permit a much larger modulus under a standard security table. After coefficient split, however, the scheme also exposes leaf ciphertexts and leaf evaluation keys at dimension \(n=N/k\). Using a modulus that is safe only for \(N\) but not for \(n\) would make the leaf instances the security bottleneck. Conversely, reducing \(Q\) may increase the top-ring security margin, but the scheme security remains the minimum over all rings.

In the current prototype, one representative parameter point is

\[
N=2^{16},\qquad k=2,\qquad n=2^{15},
\]

with dense ternary secret and approximate leaf modulus budget

\[
\log(PQ)\approx 946.
\]

A preliminary estimator calculation gives approximately \(117\)-bit security for the leaf-ring instance. This is below the conventional \(128\)-bit target but close enough to be useful for exploring the implementation path. The final paper should replace this placeholder with a complete estimator table covering the exact secret distribution, \(Q\), \(P\), and all evaluation-key instances.

The distinction between \(\log Q\) and \(\log(PQ)\) is important. Some HE security tables report a maximum ciphertext modulus budget, while practical key switching also introduces special primes \(P\). For conservative security accounting, the estimator should include the full modulus used by the relevant RLWE samples and switching keys.

### 5.4 Noise and Correctness Budget

The cost of choosing a smaller secure leaf modulus is paid in correctness margin. At the ciphertext level, a leaf route has total error of the schematic form

\[
e_{\mathrm{out}}
=
e_{\mathrm{enc}}
+e_{\mathrm{split}}
+e_{\mathrm{sd}}
+e_{\mathrm{mr}}
+e_{\mathrm{C2S}}
+e_{\mathrm{EvalMod}}
+e_{\mathrm{S2C}}
+e_{\mathrm{merge}}.
\]

Correct CKKS decoding requires this error to remain small relative to the output scale:

\[
\|e_{\mathrm{out}}\|\ll \Delta_{\mathrm{out}},
\]

and intermediate values must avoid modular wrap-around:

\[
\|\Delta m+e\| < Q/2.
\]

The reported CKKS precision is essentially a logarithmic error measure. If the average error is \(E\), then an \(L_1\)-precision value of \(p\) bits means heuristically

\[
E\approx 2^{-p}
\]

after the chosen normalization. Thus a drop from about \(22\) bits to about \(11\) bits is not a small cosmetic change; it corresponds to an error scale larger by roughly \(2^{11}\).

Lowering \(Q\) improves RLWE security but reduces the correctness margin in three ways. First, it reduces the wrap-around margin \(Q/2\). Second, it reduces the number of usable levels:

\[
L(Q)\approx \frac{\log Q-\log Q_{\mathrm{base}}}{b},
\]

where \(b\) is the typical bit-size consumed by a rescale. Third, it can force a smaller scale or a tighter first prime, which directly limits the attainable numerical precision before bootstrapping even begins. This explains why lowering \(Q\) can make EvalMod fail even when the ideal message-level theorem still holds.

The correctness target can be summarized as

\[
E_{\mathrm{split}}
+E_{\mathrm{pre}}
+E_{\mathrm{core}}
+E_{\mathrm{merge}}
\le
c\cdot \Delta_{\mathrm{out}},
\]

for a small application-dependent constant \(c\). The split theorem only identifies the ideal message map; it does not remove the need to budget all implementation errors. Therefore the practical parameter search must satisfy both the security condition

\[
\lambda(n,Q_{\mathrm{leaf}},P_{\mathrm{leaf}},\chi_s,\chi_e)
\ge
\lambda_{\mathrm{target}}
\]

and the correctness condition \(\|e_{\mathrm{out}}\|\ll \Delta_{\mathrm{out}}\). These two inequalities pull \(Q_{\mathrm{leaf}}\) in opposite directions.

If \(D_{\mathrm{cycle}}>L(Q_{\mathrm{leaf}})\), lowering the modulus may increase the required number of bootstraps. In that case the relevant comparison is

\[
B(Q_{\mathrm{leaf}})
\cdot
T_{\mathrm{boot}}^{\mathrm{split}}(N,n,Q_{\mathrm{leaf}})
<
B(Q_{\mathrm{top}})
\cdot
T_{\mathrm{boot}}^{\mathrm{direct}}(N,Q_{\mathrm{top}}).
\]

This inequality is a performance condition, not a theorem. It must be validated for each target application and implementation.

### 5.5 Cost Model and Benefit Boundary

The split route has the following schematic cost:

\[
T_{\mathrm{boot}}^{\mathrm{split}}
=
T_{\mathrm{Split}}
+
T_{\mathrm{Merge}}
+
T_{\mathrm{leaf-pre}}
+
T_{\mathrm{leaf-core}},
\]

where

\[
T_{\mathrm{leaf-pre}}
\approx
\max_{r}
T_{\mathrm{ScaleDown}_n+\mathrm{ModRaise}_n}^{(r)}
\]

and

\[
T_{\mathrm{leaf-core}}
\approx
\max_{r}
T_{\mathrm{C2S}_n+\mathrm{EvalMod}_n+\mathrm{S2C}_n}^{(r)}
\]

under ideal leaf parallelism. Without parallelism, these terms scale roughly by \(k\). The direct large-ring route costs

\[
T_{\mathrm{boot}}^{\mathrm{direct}}
=
T_{\mathrm{ScaleDown}_N+\mathrm{ModRaise}_N}
+
T_{\mathrm{C2S}_N+\mathrm{EvalMod}_N+\mathrm{S2C}_N}.
\]

Since NTTs, rotations, and key-switching operations scale roughly with the ring dimension and the number of moduli, a coarse model is

\[
T_{\mathrm{op}}(d,Q)
\approx
C_{\mathrm{op}}\cdot d\log d\cdot \#\mathrm{primes}(Q).
\]

Thus moving the dominant linear transforms and EvalMod work from \(N\) to \(n=N/k\) can be beneficial, especially when the leaf computations run in parallel. More explicitly, if the large-ring core has cost roughly

\[
C\cdot N\log N\cdot \#\mathrm{primes}(Q_{\mathrm{top}}),
\]

while the leaf route has \(k\) leaf cores of size \(n=N/k\), then perfect parallelism suggests a leading term closer to

\[
C\cdot n\log n\cdot \#\mathrm{primes}(Q_{\mathrm{leaf}}),
\]

plus split and merge. Without parallelism, the \(k\) leaf cores contribute about

\[
k\cdot C\cdot n\log n\cdot \#\mathrm{primes}(Q_{\mathrm{leaf}})
=
C\cdot N\log n\cdot \#\mathrm{primes}(Q_{\mathrm{leaf}}),
\]

which can still be smaller than the large-ring route when \(\#\mathrm{primes}(Q_{\mathrm{leaf}})\) is smaller or when the leaf engine is substantially optimized. However, the benefit is reduced by split/merge overhead, memory bandwidth, scheduling overhead, and the cost of additional ring-switching and leaf bootstrapping keys.

The correct end-to-end comparison must include bootstrapping frequency. If a smaller modulus decreases the number of levels per refresh, the number of bootstraps may increase. A sufficient performance condition is

\[
B(Q_{\mathrm{leaf}})
\cdot
T_{\mathrm{boot}}^{\mathrm{split}}(N,n,Q_{\mathrm{leaf}})
<
B(Q_{\mathrm{top}})
\cdot
T_{\mathrm{boot}}^{\mathrm{direct}}(N,Q_{\mathrm{top}}),
\tag{8}
\]

where \(B(Q)\) denotes the number of bootstraps required by the application under modulus budget \(Q\). In the best target regime, the bootstrapping placement is determined by algorithmic structure, such as activation boundaries or refresh cycles, and

\[
B(Q_{\mathrm{leaf}})=B(Q_{\mathrm{top}}).
\]

Then the comparison reduces to the per-bootstrap cost. In a depth-bound regime where

\[
B(Q_{\mathrm{leaf}})>B(Q_{\mathrm{top}}),
\]

the split route must compensate by making each bootstrap sufficiently cheaper.

This is also where future leaf-bootstrapping optimizations matter. In the cost expression above, the ordinary leaf engine contributes \(T_{\mathrm{leaf-core}}\). If an optimized leaf method replaces it by \(T_{\mathrm{leaf-core}}^{\mathrm{opt}}\), then the same split theorem gives the condition

\[
B(Q_{\mathrm{leaf}})
\cdot
\left(
T_{\mathrm{Split}}+T_{\mathrm{Merge}}+T_{\mathrm{leaf-pre}}
+T_{\mathrm{leaf-core}}^{\mathrm{opt}}
\right)
<
B(Q_{\mathrm{top}})
\cdot
T_{\mathrm{boot}}^{\mathrm{direct}}.
\]

Therefore, the current prototype should be viewed as a conservative instantiation rather than as the performance ceiling of the method.

Another natural baseline is to use several smaller ciphertexts instead of one large ciphertext. If \(c\) ciphertexts at dimension \(n\) replace one ciphertext at dimension \(N\), then

\[
T_{\mathrm{multi-ct}}
\approx
c\cdot T_{\mathrm{boot}}(n)
+
T_{\mathrm{cross\text{-}ct}},
\]

where \(T_{\mathrm{cross\text{-}ct}}\) accounts for cross-ciphertext layout, aggregation, or scheduling overhead. The split route is most attractive when the application benefits from keeping a single large-ring ciphertext layout, or when cross-ciphertext communication would be expensive.

These observations delimit the claim. The method is not a universal replacement for large-ring bootstrapping. It is a modular factorization that is useful when the large ring is needed for packing, the leaf modulus is sufficient for the refresh cycle, and the saved leaf computation outweighs split/merge and key overhead.

### 5.6 Anticipated Reviewer Questions and Design Responses

A skeptical reviewer may first ask why one should reduce the modulus at all, since bootstrapping is meant to restore levels and a smaller \(Q\) can force more frequent bootstrapping. The correct answer is that the method is not designed for depth-bound circuits. It should be evaluated under condition (8). If reducing \(Q\) changes \(B(Q)\) unfavorably, then the split route must compensate by a sufficiently cheaper per-bootstrap cost; otherwise it is not the right parameter choice.

A second objection is that the large ring can support a much larger modulus, so using a leaf-secure modulus appears to waste the large ring. This is a real tradeoff, but not a contradiction. In the intended regime, \(N\) is selected because \(S(n)<M\le S(N)\), not because the application needs the largest modulus allowed by \(N\). A larger \(Q\) is only useful if it reduces bootstrapping frequency or improves the required output precision. If it only increases arithmetic and key-switching cost, it is not an application-level resource being wasted.

A third objection is sharper: why not keep the large \(Q\) after splitting, since that may improve precision? This is generally not a secure leaf-ring parameter. Once leaf ciphertexts and leaf evaluation keys are exposed, the leaf dimension \(n\) determines the relevant RLWE security bound. Therefore \(Q_{\mathrm{leaf}}\) must be checked at dimension \(n\), regardless of the top-ring dimension. A top-ring parameter that looks highly secure at \(N\) can still be insecure for the leaf instances.

A fourth objection is that one could simply use several small ciphertexts instead of one large ciphertext. This is a valid baseline and should be included experimentally. The split route is useful when a single large-ring layout simplifies packing, rotations, aggregation, or interaction with surrounding computation, or when cross-ciphertext communication would dominate. It is not meant to replace the multi-ciphertext baseline in workloads where independent small ciphertexts are naturally sufficient.

A fifth objection is that coefficient split is not slot-preserving, so it may break the SIMD semantics of the application. This is precisely why the paper does not claim slot preservation. The claim is narrower: after \(\mathsf{C2S}\), the relevant bootstrapping coordinates are coefficient-packed slots, and these coordinates split by residue class. The construction is therefore valid for the bootstrapping core, while arbitrary slot-wise application logic would require a separate layout argument or the \(B_k\) recovery transform.

A sixth objection is that the theorem is message-level while CKKS implementations contain scale, rounding, RNS extension, and key-switching errors. This is why Section 4 separates the ideal factorization from ciphertext-level compatibility. The implementation must verify that preprocessing, split/merge, C2S/S2C, EvalMod, and rescaling errors fit within the CKKS precision budget.

A final objection is that the current implementation uses an ordinary leaf bootstrapping engine, so the measured speedup may be far from optimal. This is a limitation of the prototype rather than of the factorization. The leaf engine is a black-box module in the theorem; optimized bootstrapping methods can be inserted at the leaf level as long as they implement the same normalized scalar EvalMod semantics.

### 5.7 Preliminary Experimental Snapshot

The following numbers are current prototype measurements. They are included to summarize the parameter exploration and to motivate the final experimental design.

| Experiment | Top \(N\) | Leaf \(n\) | Route | Approx. \(\log(PQ)\) | Input L1 precision | Output L1 precision | Time |
|---|---:|---:|---|---:|---:|---:|---:|
| Low-scale leaf preprocessing | \(2^{16}\) | \(2^{15}\) | Split \(\to\) leaf ScaleDown/ModRaise/C2S | \(936\) | \(12.09\) bits | \(11.50\) bits | \(4.45\)s |
| Higher-scale leaf preprocessing | \(2^{16}\) | \(2^{15}\) | Split \(\to\) leaf ScaleDown/ModRaise/C2S | \(946\) | \(22.09\) bits | \(21.45\) bits | \(4.57\)s |
| Only ModRaise moved to leaf | \(2^{16}\) | \(2^{15}\) | top ScaleDown \(\to\) Split \(\to\) leaf ModRaise | \(936\) | \(12.08\) bits | \(0.96\) bits | \(4.67\)s |
| Earlier safe-window candidate | \(2^{17}\) | \(2^{16}\) | top preprocessing \(\to\) Split \(\to\) leaf core | \(1696\) | \(21.08\) bits | \(20.68\) bits | \(68.00\)s |

The comparison between the second and third rows is the most informative implementation observation so far. Moving only \(\mathsf{ModRaise}\) to the leaf ring causes a severe normalization mismatch. Moving the paired preprocessing stage \(\mathsf{ScaleDown}\circ\mathsf{ModRaise}\) to the leaf ring restores the expected precision behavior.

These results also show that a low output precision may simply reflect a low input precision caused by a small CKKS scale. In the low-scale run, the input ciphertext itself has only about \(12\) bits of L1 precision, while the bootstrapping route loses less than one additional bit. Raising the scale from \(2^{25}\) to \(2^{35}\), while increasing the corresponding \(q_0\), restores the output precision to above \(21\) bits.

### 5.8 Current Limitations and Next Steps

The current implementation establishes the feasibility of split-first leaf preprocessing, but several points remain open.

First, the security estimate around \(117\) bits is preliminary and should be replaced by a full estimator table. Second, the present leaf engine is an ordinary CKKS bootstrapping implementation; it has not yet incorporated state-of-the-art leaf-level bootstrapping optimizations. Third, the exact scale and modulus normalization conditions should be stated explicitly in a ciphertext-level correctness theorem. Finally, key size and memory consumption for split/merge keys and leaf bootstrapping keys should be reported separately from runtime.

The intended final claim is therefore conservative: the coefficient-split factorization is a general semantic and algebraic method; CKKS bootstrapping is the first concrete instance; and the practical benefit appears in parameter regimes where the large ring is required by packing capacity while the leaf-ring modulus budget remains sufficient for correctness and security.
