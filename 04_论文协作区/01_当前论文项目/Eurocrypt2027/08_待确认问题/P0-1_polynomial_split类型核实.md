# P0-1: Polynomial Split 类型核实

> 基于: 初稿.pdf (src_000031) + 补充视角-slot情况下的拆分.md (src_000034)
> 结论: coefficient/module split，非 slot-preserving split

---

## 一、初始判断

本论文主路线采用 **coefficient/module split**，不是 slot-preserving split。

### 源自初稿（src_000031）的证据

1. **Section 1.2 核心观察**（P3）:
   > "关键观察不是这个 split保留了大环 slots；事实上它并不保留。"
   
   这是论文最核心的自述：split 不保留 big-ring slots。

2. **Section 2.5** 标题:
   > "Coefficient split 与 slot-preserving split 的最小区别"

3. **贡献一**（P4）:
   > "split之后不需要先恢复原始 big-ring slots"

4. **Section 3.2 Coefficient/module split**:
   > P(X) = Σ_{r=0}^{k-1} X^r P_r(X^k)

5. **Section 3.4**:
   > "Coefficient split 在 slot视角下发生什么" — 解释为何不保持slots

### 源自补充视角文件（src_000034）的证据

明确区分了两种split:
- **类型1: Coefficient split**: P(X) = A(X²) + X·B(X²)
- **类型2: Slot split**: Q⁺(Y) = A(Y) + h(Y)B(Y), Q⁻(Y) = A(Y) − h(Y)B(Y)

初稿选择的是类型1。

## 二、初稿中是否始终保持此说法？

✅ 一致。核心主张"split 不保持 big-ring slots"贯穿 Introduction → Section 3 → Theorem 4。

## 三、与逐坐标证明的一致性

待精读 `逐坐标分解证明.pdf` (src_000032) 确认，但从初稿的定理陈述来看，逐坐标分解也是基于 coefficient split 的 residue class 分组。

## 四、是否需加入"不恢复big-ring slots"强调？

**建议加入。** 已在 Introduction 和 Section 3.9 有说明，但 Section 4 定理陈述前可再加一句:
> "We emphasize that Split_k does not recover the original big-ring slot values."

## 五、是否需固定 Split_k 的定义？

**必须固定。** Notations (Section 4.1) 应包含:
```
Split_k: R_{q,N} → (R_{q,n})^k
ct(X) ↦ (ct_0(Y), ..., ct_{k-1}(Y))
where each leaf is obtained by coefficient/module decomposition
ct_r(Y) = coeff of X^r in ct(X) under X^N+1 = (Y)^k + 1
```

## 六、最终结论

| 问题 | 答案 |
|------|------|
| Split类型 | **coefficient/module split** |
| 是否slot-preserving | **否** |
| 是否混用两种split | **否**，初稿明确区分 |
| 是否需要固定notation | **是**（Section 4.1） |
| 是否影响correctness | **不直接影响**（split是正确的polynomial identity） |
| 是否影响related work | **需要与HERMES/subring-encapsulation对比** |
| 状态 | ✅ 已核实 (high confidence) |

## 七、仍需人工确认

1. Section 3.6 ("Slot split与slot-preserving split") 是否讨论过度（可能使reviewer误认为论文也做slot split）
2. Related Work 中是否明确与 Geelen 的slot-to-coefficient 做区分
