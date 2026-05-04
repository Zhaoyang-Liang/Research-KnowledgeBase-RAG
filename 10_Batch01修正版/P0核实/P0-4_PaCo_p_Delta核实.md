# P0-4: PaCo 中 p 和 Δ 的精确区别

> 生成时间: 2026-05-04 | 来源: PaCo 原文 Section 2-3 + 翻译 K024 | 状态: 已从原文核实

## 问题

PaCo 中 p 和 Δ 分别是什么？是否在不同语境下使用？

## 原文核实结果

### 原文关键段落（PaCo Section 2）

> "Two Different Scaling Factors p and Δ.
> To support flexible parameterization, we distinguish between two scaling
> factors: **p** and **Δ**. The factor **p ≤ q** is used for homomorphic
> operations **outside bootstrapping**, whereas **Δ ≥ q** applies **during
> bootstrapping**. Starting from a ciphertext under the modulus q, we
> bootstrap over the moduli qp·Δ^ℓ for ℓ ∈ [0, L] to finally output a
> ciphertext under the modulus qp, supporting one additional multiplication
> with scaling factor p."

### p 的含义

- **p = 明文模数 (plaintext modulus)**
- p ≤ q（小于等于当前密文模数）
- 用于 bootstrapping **之外**的同态操作
- 候选值范围：p ≈ 2^30（实际取值取决于参数选择）
- 与 CKKS 编码的近似误差相关

### Δ 的含义

- **Δ = 编码因子 (encoding/scaling factor)**
- Δ ≥ q（大于等于模数）
- 用于 bootstrapping **内部**过程
- 作用：将复数向量放大为整数多项式
- 候选值范围：Δ ≈ 2^80（远大于 p）

### 两者如何配合

PaCo bootstrapping 的模数链结构：

```
密文模数 q
  → bootstrapping 内部: qp·Δ^L → qp·Δ^{L-1} → ... → qp·Δ^0
  → bootstrapping 输出: 模数 qp（支持一次额外乘法，缩放因子 p）
```

### ⚠ 注意

PaCo 原文的 Section 4.2 标题就是 "Two Different Scaling Factors p and ∆"，
说明作者明确意识到了这两个参数容易混淆，从而专门用一个小节区分。

## 对 Packing Relation 和盲旋转的影响

| 组件 | 使用 p 还是 Δ | 说明 |
|------|---------------|------|
| CKKS 编码 | Δ | 典范嵌入放大因子 |
| 同态乘法（bootstrapping 外） | p | 消息模数 |
| 盲旋转 test vector | Δ | 需要用 Δ 放大 test vector |
| 盲旋转结果 | 混合 | 解密结果在模数 qp·Δ^ℓ 下 |
| Packing Relation | p | 打包结果模 p |

## 对已有卡片的影响

| 卡片 | 需要修正的内容 |
|------|----------------|
| K024 | p 和 Δ 的区分已验证正确 ✅ |
| QA-K024-2 | "PaCo p vs Δ" 现在可标记为**已核实** |
| QA-K024-4 | "Packing Relation" 需确认是否使用 p 作为消息模数 |

## 推荐论文写作建议

在论文中建议明确：

1. 区分 "plaintext modulus p" 和 "scaling factor Δ"
2. 引用 PaCo Section 4.2 的定义
3. 如果论文使用了 PaCo 的参数约定，直接引用
4. 如果论文有自己的参数命名，在 Preliminaries 中明确定义

## 来源

- PaCo 原文 (Sec 4.2 "Two Different Scaling Factors p and ∆")
- K024 (PaCo 翻译)
- QA-K024-2
