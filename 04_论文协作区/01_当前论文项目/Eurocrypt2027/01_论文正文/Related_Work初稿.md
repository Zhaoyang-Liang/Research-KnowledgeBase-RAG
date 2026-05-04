# Related Work 初稿

> paper_id: paper_2027_eurocrypt_fhe  
> 用途：新初稿 `Related Work` 章节的第一版文字。  
> 引用键暂用：`CKKS17`, `GHPS12RingSwitching`, `HERMES23`, `PaCo25`, `SubringSecret25`。

## Related Work

This work lies at the intersection of CKKS bootstrapping, ring switching, and subring-based acceleration techniques. We briefly review the closest lines of work and clarify how our contribution differs from them.

### CKKS and Conventional CKKS Bootstrapping

CKKS was introduced by Cheon, Kim, Kim, and Song as a homomorphic encryption scheme for approximate arithmetic over complex vectors \cite{CKKS17}. Its native plaintext semantics are based on encoding complex vectors into cyclotomic-ring plaintext polynomials via the canonical embedding, while homomorphic multiplication is followed by rescaling to control the growth of the scale and modulus consumption. This approximate plaintext semantics is the foundation of modern encrypted real/complex-number computation.

CKKS bootstrapping refreshes a ciphertext whose modulus has nearly been consumed. In the conventional bootstrapping pipeline, after modulus raising, the core computation is organized as

```text
CoeffToSlot -> EvalMod -> SlotToCoeff.
```

The CoeffToSlot transform moves coefficient-like values into encrypted slots, EvalMod homomorphically removes the modular term by applying an approximate scalar modular reduction, and SlotToCoeff maps the reduced values back to coefficient representation. A large part of the bootstrapping cost comes from the two linear transforms and the EvalMod approximation. Existing work has improved this pipeline through better linear transforms, better polynomial approximations, better parameter selection, or higher-precision bootstrapping.

Our work does not introduce a new CKKS scheme and does not replace the conventional EvalMod path. Instead, we keep the traditional `C2S -> EvalMod -> S2C` core and prove that this core is compatible with a coefficient/module split of the underlying plaintext polynomial. The main question we study is whether a large-ring CKKS bootstrapping core can be factored into multiple subring leaf cores, not how to redesign CKKS bootstrapping from first principles.

### Ring Switching and Coefficient/Module Decomposition

Ring switching was studied by Gentry, Halevi, Peikert, and Smart in the context of BGV-style homomorphic encryption \cite{GHPS12RingSwitching}. A central algebraic ingredient is the decomposition of a large-ring polynomial into coefficient residue classes. For `N = k n`, a polynomial can be written as

```text
P(X) = sum_{r=0}^{k-1} X^r P_r(X^k),
```

where each `P_r` belongs to a smaller ring. This decomposition is linear and allows the large-ring plaintext polynomial to be recovered by a simple merge map. At the ciphertext level, however, moving between rings is not a free coefficient rearrangement: it requires key switching or secret switching and introduces additional noise.

Our construction uses this coefficient/module decomposition as the semantic target of the subring split. The novelty is not the decomposition itself, nor ring switching as a ciphertext operation. Rather, we show that this decomposition is compatible with the `C2S -> EvalMod -> S2C` sandwich structure of CKKS bootstrapping. In particular, after coefficient split, we do not need to restore the original large-ring canonical slots. The leaf CoeffToSlot transforms already produce the coefficient-packed coordinates that the large-ring CoeffToSlot would have produced, grouped by residue class.

### Ring Packing and HERMES

HERMES gives an efficient ring packing framework using MLWE ciphertexts and applies it to transciphering into CKKS \cite{HERMES23}. It combines ring switching, bootstrapping, and intermediate MLWE representations to pack many LWE/MLWE-format ciphertexts into an RLWE/CKKS ciphertext suitable for subsequent homomorphic computation. HERMES is motivated by the complementary strengths of LWE-format ciphertexts, which provide fine granularity, and RLWE-format ciphertexts, which provide SIMD throughput.

Although HERMES and our work both use ring switching ideas in a CKKS-related setting, the computational tasks are different. HERMES starts from many LWE or MLWE ciphertexts and produces an RLWE/CKKS ciphertext, with ring packing and transciphering as the main goals. In contrast, our input is already a CKKS/RLWE ciphertext, and our goal is to factor its own bootstrapping core into independent subring computations. We do not solve LWE-to-RLWE packing, and we do not rely on the MLWE midpoint or column/row packing methods used in HERMES. Our result is instead a compatibility theorem for the traditional CKKS bootstrapping core under coefficient/module decomposition.

### Subring Secret Encapsulation

A closely related subring-based acceleration technique is subring secret encapsulation, recently proposed for practical dense-key bootstrapping \cite{SubringSecret25}. That work switches the secret key to a dense secret in a subring before bootstrapping. The smaller Hamming weight and algebraic structure of the subring secret improve several bootstrapping components: EvalMod or digit removal benefits from a smaller approximation range, and CoeffsToSlots/SlotsToCoeffs benefit from hoisted key switching under the subring secret.

This line of work is close in spirit because it also exploits subring structure to accelerate bootstrapping. However, the mechanism is different. Subring secret encapsulation is a secret-side optimization: it keeps the ciphertext as a single bootstrapping object but changes the secret under which the bootstrapping is performed. Our method is a ciphertext/message-side decomposition: it ring-switches the ciphertext semantics into multiple leaf ciphertexts, runs leaf bootstrapping cores independently, and merges the result back. Thus, subring secret encapsulation is best understood as subring-secret acceleration, whereas our work is subring-ciphertext core factorization.

This difference is closely related to a point made in \cite{SubringSecret25}: for SIMD-packed ciphertexts, a traditional ring-switching approach would typically unpack the ciphertext into low-dimensional sub-ciphertexts, apply a multi-ciphertext linear transformation to recover the original slot values, bootstrap the sub-ciphertexts, and finally repack them. In our terminology, the slot-recovery linear transformation is precisely the kind of fiber-wise Vandermonde transform denoted by \(B_k\): it maps coefficient-split leaf canonical slots back to the original large-ring canonical slots. Our observation is that this recovery target is unnecessary for the conventional CKKS bootstrapping core. Since EvalMod is applied after CoeffToSlot, the relevant invariant is not the original slot vector but the coefficient-packed slot vector. Thus, instead of implementing or optimizing the \(B_k\)-type recovery transform, we prove that the `C2S -> EvalMod -> S2C` core factors through the coefficient split directly.

This distinction also matters for security claims. The dense-subring-secret analysis in \cite{SubringSecret25} does not automatically imply security for our leaf ciphertexts or ring-switching keys. In our setting, the top ciphertexts, leaf ciphertexts, ring-switching keys, and leaf bootstrapping keys must be considered under their own ring dimensions and public moduli.

### New CKKS Bootstrapping Designs: PaCo

PaCo introduces a new CKKS bootstrapping procedure based on partial CoeffToSlot \cite{PaCo25}. It reformulates the CKKS decryption equation using blind rotations and modular additions, uses the circle group in the complex plane to simulate modular addition, and relies on alternative polynomial ring structures. PaCo therefore changes the internal bootstrapping circuit and does not simply optimize the conventional `CoeffToSlot -> EvalMod -> SlotToCoeff` path.

Our approach is orthogonal to PaCo. We do not reformulate the CKKS decryption equation, do not use blind rotations as the main bootstrapping mechanism, and do not replace EvalMod by circle-group modular additions. Instead, we preserve the conventional CKKS bootstrapping core and change its execution granularity. The central observation is that coefficient split followed by leaf CoeffToSlot produces the same coefficient-packed data as large-ring CoeffToSlot followed by residue-class grouping. This allows EvalMod to be evaluated independently on the leaves because it is coordinate-wise in the coefficient-packed slots.

### Positioning of This Work

The closest conceptual point of comparison is the following: ring switching and HERMES show that ciphertexts and plaintext polynomials can be moved between rings; subring secret encapsulation shows that subring structure can accelerate bootstrapping through the secret; PaCo shows that alternative CoeffToSlot-like structures can support new bootstrapping designs. Our work focuses on a different compatibility question:

```text
Does the conventional CKKS bootstrapping core itself factor through coefficient/module decomposition?
```

We answer this question affirmatively in an ideal algebraic model. The resulting factorization is not a slot-preserving split. Coefficient split generally maps large-ring slots to local linear combinations on each fiber, and recovering the original large-ring slots would require an additional slot-lifting transform. The key point is that CKKS bootstrapping does not require such recovery after the split: EvalMod is applied after CoeffToSlot, and the relevant coordinates are coefficient-packed slots. This is why the proposed no-`B_k` route can factor the core without explicitly restoring the large-ring canonical slots.

The present version should therefore be read as an algebraic and implementation-oriented step toward subring-factored CKKS bootstrapping. A complete treatment of noise, key-switching error, final security parameters, and reduced-modulus leaf bootstrapping is left to the parameter and security analysis.

## 中文写作版

本文位于 CKKS bootstrapping、ring switching 和 subring-based acceleration 的交叉处。下面按相关程度说明最接近的几条工作线，并明确本文与它们的区别。

### CKKS 与传统 CKKS bootstrapping

CKKS 由 Cheon、Kim、Kim 和 Song 提出，用于支持复数向量上的近似同态计算 \cite{CKKS17}。其明文语义基于 canonical embedding：复数向量被编码为 cyclotomic ring 中的 plaintext polynomial，同态乘法之后通过 rescale 控制 scale 和 modulus chain 的增长。本文完全建立在 CKKS 的近似明文语义之上，但并不提出新的 CKKS scheme。

传统 CKKS bootstrapping 的核心通常可写成：

```text
CoeffToSlot -> EvalMod -> SlotToCoeff.
```

其中 CoeffToSlot 将 coefficient-like values 移入 encrypted slots，EvalMod 对这些 slots 逐坐标执行近似模约减，SlotToCoeff 再将结果转回 coefficient representation。本文保留这一传统 core，不替换 EvalMod，也不重新设计 CKKS 的 decryption equation。本文研究的问题是：这个大环上的传统 core 是否可以沿 coefficient/module decomposition 分解为多个小环上的 leaf cores。

### Ring switching 与 coefficient/module decomposition

Gentry、Halevi、Peikert 和 Smart 的 ring switching 工作给出了从大环到小环的经典工具 \cite{GHPS12RingSwitching}。其核心代数结构是：当 `N = k n` 时，大环 polynomial 可以按 coefficient residue class 写成

```text
P(X) = sum_{r=0}^{k-1} X^r P_r(X^k).
```

该分解是线性的，并且可以通过 merge map 恢复原始 polynomial。本文使用的 `Split_k/Merge_k` 正是沿这一 module decomposition 的语义展开。因此，ring switching 本身不是本文的新贡献；本文也不能把 ciphertext-level split 描述成免费的公开系数重排，因为真实密文转换需要 key switching 并引入噪声。

本文的新点是证明这一 coefficient/module split 与 CKKS bootstrapping core 的 `C2S -> EvalMod -> S2C` sandwich structure 兼容。换言之，split 后不需要先恢复原始 big-ring canonical slots；leaf CoeffToSlot 已经会产生 top CoeffToSlot 输出的 coefficient residue-class 子序列。

### HERMES 与 ring packing

HERMES 使用 MLWE ciphertexts 给出了高效 ring packing 框架，并将其用于 transciphering 到 CKKS \cite{HERMES23}。它结合 ring switching、bootstrapping 和 MLWE midpoint，将许多 LWE/MLWE-format ciphertexts 打包成一个 RLWE/CKKS ciphertext，以便后续进行 SIMD 同态计算。

本文与 HERMES 都使用 ring switching 思想，但任务不同。HERMES 的输入是许多 LWE/MLWE ciphertexts，目标是 ring packing / transciphering；本文的输入已经是 CKKS/RLWE ciphertext，目标是刷新它自身，并将其 bootstrapping core 分解到多个 subring leaf cores 上。本文不解决 LWE-to-RLWE packing，也不使用 HERMES 的 MLWE midpoint、column method 或 row method。本文给出的是传统 CKKS bootstrapping core 在 coefficient/module split 下的兼容性定理。

### Subring secret encapsulation

Subring secret encapsulation 是另一个非常接近的 subring-based bootstrapping acceleration 技术 \cite{SubringSecret25}。该方法在 bootstrapping 前将 secret 切换到 dense subring secret。由于 subring secret 的 Hamming weight 更小，EvalMod 或 digit removal 的近似范围降低；同时 subring secret 的代数结构也使 CoeffsToSlots 和 SlotsToCoeffs 中的 key switching 可以 hoist，从而提升 bootstrapping 吞吐。

本文与该工作最关键的区别在于：subring secret encapsulation 是 secret-side optimization，而本文是 ciphertext/message-side decomposition。前者仍然围绕单个 ciphertext 做 bootstrapping，只是切换 secret；后者将 ciphertext 的消息语义拆成多个 leaf ciphertexts，在每个 leaf 上独立执行 leaf bootstrapping core，再 merge 回大环。因此，前者可以称为 subring-secret acceleration，本文更准确地说是 subring-ciphertext core factorization。

这一点也对应 \cite{SubringSecret25} 对传统 packed ring switching 的一个判断：对于 SIMD-packed ciphertext，若直接用 ring switching 加速 bootstrapping，通常需要先 unpack 成多个低维 sub-ciphertexts，再施加一个 multi-ciphertext linear transformation 恢复原始 slot values，随后分别 bootstrapping，最后再 repack。用本文的记号看，这个恢复原始 slots 的线性变换正是 \(B_k\) 这一类 fiber-wise Vandermonde transform：它把 coefficient-split 后的 leaf canonical slots 重新组合为 large-ring canonical slots。本文的关键观察是：传统 CKKS bootstrapping core 并不需要这个恢复目标。因为 EvalMod 作用在 CoeffToSlot 之后，真正需要保持的是 coefficient-packed slot vector，而不是原始 canonical slot vector。因此，本文不是实现或优化 \(B_k\)，而是证明 `C2S -> EvalMod -> S2C` core 可以直接穿过 coefficient split。

这一差异也影响安全表述。Subring secret encapsulation 的 dense subring secret 安全分析不能直接迁移到本文。本文最终仍需分别分析 top ciphertext、leaf ciphertexts、RSDown/RSUp keys 和 leaf bootstrapping keys 的 ring dimension 与 public modulus。

### PaCo 与新式 CKKS bootstrapping

PaCo 提出了一条新的 CKKS bootstrapping 路线，通过 partial CoeffToSlot、blind rotations、circle-group modular additions 和 alternative ring structures 来重写 bootstrapping 逻辑 \cite{PaCo25}。它不是简单优化传统 CoeffToSlot 或 EvalMod，而是在更深层次上改写了 bootstrapping circuit。

本文与 PaCo 是正交关系。本文不重新表述 CKKS decryption equation，不用 blind rotations 作为主要机制，也不替换 EvalMod。本文保留传统的 `CoeffToSlot -> EvalMod -> SlotToCoeff` core，只改变其执行粒度：先做 coefficient/module split，再在 leaf 上分别执行 `CoeffToSlot -> EvalMod -> SlotToCoeff`，最后 merge。

### 本文定位

综上，已有 ring switching 和 HERMES 说明 ciphertexts / plaintext polynomials 可以在不同环之间转换；subring secret encapsulation 说明 subring secret 可以加速 bootstrapping；PaCo 说明 alternative CoeffToSlot-like structures 可以支持新的 bootstrapping 设计。本文回答的是另一个问题：

```text
传统 CKKS bootstrapping core 本身是否可以沿 coefficient/module decomposition 分解？
```

本文的回答是肯定的。该分解不是 slot-preserving split。Coefficient split 在 slot 视角下通常会将 big-ring slots 做局部线性组合；若要恢复原始 big-ring slots，需要额外的 slot-lifting transform。本文的关键观察是：CKKS bootstrapping 并不需要 split 后恢复 big-ring slots，因为 EvalMod 作用在 CoeffToSlot 之后的 coefficient-packed slots 上。正因如此，本文的 no-`B_k` 路线可以不显式恢复 large-ring canonical slots，而直接完成 core factorization。

当前版本应理解为 subring-factored CKKS bootstrapping 的代数和实现证据。完整 noise bound、key-switching error、最终安全参数和 reduced-modulus leaf bootstrapping 留待参数与安全章节处理。
