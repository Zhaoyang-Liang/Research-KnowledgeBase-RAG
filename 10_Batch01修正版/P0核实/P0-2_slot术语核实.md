# P0-2: CKKS slot 与 BGV/BFV CRT slot 术语核实

> 生成时间: 2026-05-04 | 来源: K028 + K003/K008/K009/K018/K019/K020 | 状态: 已从原文核实

## 问题

当前知识卡片和关系图中，凡是出现 "slot" 的地方，具体是 CKKS complex slot 还是 BGV/BFV CRT slot？

## 核实结论

两种 slot 的数学来源不同，已在 K028 中明确区分。以下是逐卡片核实：

| 卡片 | "slot" 出现位置 | 实际含义 | 判断 |
|------|----------------|----------|------|
| K003 | "CKKS 槽数=N/2" | CKKS complex slot | ✅ 上下文明确 |
| K008 | "三条线统一描述 slot 结构" | 泛指(代数层面的槽) | ⚠ 偏 BGV/BFV CRT slot |
| K009 | "每个 F_{p^d} 因子=一个 slot" | BGV/BFV CRT slot | ✅ 明确 |
| K018 | "CKKS 编码使用半槽" | CKKS complex slot | ✅ 明确 |
| K019 | "slot 复数向量→U_n→系数" | CKKS complex slot | ✅ 明确 |
| K020 | "slot 值恢复多项式系数" | CKKS complex slot | ✅ 明确 |
| K028 | "CKKS slot vs BGV/BFV CRT slot" | 专门区分卡片 | ✅ 明确 |

## 推荐术语 (论文写作)

| 上下文 | 推荐写法 | 避免写法 |
|--------|----------|----------|
| CKKS 方案 | "CKKS complex slot" 或 "复嵌入槽" | 不写 "CRT slot" |
| BGV/BFV 方案 | "CRT slot" 或 "CRT 槽" | 不写 "complex slot" |
| 泛指讨论 | "slot-like component" + 明确定义 | 不单独写 "slot" 无上下文 |
| 代数层面的槽 | "代数槽 (algebraic slot)" 来自 CRT 分解 | 不与 CKKS 槽混用 |
| 混合讨论 | 开篇定义两种 slot，后续用缩写 | 不交叉使用 |

## 论文中推荐的定义段落 (模板)

> 在本文中，我们区分两种 "slot" 概念：
> - **CRT slot**（代数槽）：在 BGV/BFV 方案中，通过中国剩余定理分解
>   R_p ≅ ∏ F_{p^d}，每个直积因子 F_{p^d} 称为一个 CRT slot，
>   用于存储一个有限域元素。
> - **CKKS complex slot**（复嵌入槽）：在 CKKS 方案中，通过典范嵌入
>   σ: R → C^{N/2}，每个坐标分量称为一个 complex slot，
>   用于存储一个复数。
>
> 当上下文不明确时，我们显式标注 slot 类型。

## K028 卡片本身

K028 专门区分了两种 slot，内容需保留并强化：
- ✅ 确认了 BGV/BFV CRT slot ← CRT 分解
- ✅ 确认了 CKKS slot ← 复嵌入
- ⚠ 需要人工确认：K028 中的 "Qclaw推断" 标签是否正确

## 来源

- 源文件: K003, K008, K009, K018, K019, K020, K028
- GHPS: 不涉及 CKKS slot 概念（使用 BGV 方案）
- PaCo: 使用 CKKS 编码，slot 为 CKKS complex slot
