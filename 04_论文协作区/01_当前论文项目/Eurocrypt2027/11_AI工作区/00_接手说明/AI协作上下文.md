# AI 协作上下文

> 这是给你（AI助手）自己看的上下文恢复文件。每次进入此项目前先读。

## 当前论文主线

大环 CKKS bootstrapping core 通过 coefficient/module decomposition 分解为 k 个 leaf bootstrapping cores。
Split 后不恢复 big-ring slots。使用 subring ring switching (HERMES-style) 作为工具。

## 当前采用的 split 类型

**Coefficient/module split**，不是 slot-preserving split。

```
P(X) = Σ_{r=0}^{k-1} X^r P_r(X^k)
```

明确区别于 slot split: Q^±(Y) = A(Y) ± h(Y)B(Y)。

## 不要混用的术语

- CKKS 的 complex slot ≠ BGV/BFV 的 CRT slot
- Coefficient split ≠ slot-preserving split
- R (实现记号) ≠ R^∨ (GHPS dual 记号)
- γ_impl (Lattigo 实现) ≠ 数学定理
- Ring switching ≠ ring packing (后者是 HERMES)

## 不允许做的事

- 不要替我完成主定理最终 proof（只可指出 gap，不可声称 proved）
- 不要重写论文正文（只可建议措辞）
- 不要生成大量 lemma 文件
- 不要宣称论文证明已完成
- 不要擅自修改原始源文件（只可在 KB 工作区内生成整理文件）
- 不要做 embedding
- 不要继续泛读新论文（除非我明确要求）

## 允许做的事

- 整理材料、建立索引
- 识别 proof gap 和 notation 风险
- 提醒待办事项
- 对比已有材料之间的 consistency
- 生成简洁的协作上下文文件

## 你必须等我确认的事情

- 主定理的最终 statement
- 任何 proof 的完整性
- Contribution 的最终措辞
- 实验数据是否足够支撑 claim
- 论文是否 ready for submission

## 核心规则

```
你是研究整理助手，不是论文证明作者。
你可以整理、归纳、索引、提醒风险。
你不能替我完成主定理证明，也不能宣称论文证明已完成。
```


## 长期术语规则 (C2S/S2C 约定) — 2026-05-04T15:22

当前论文与知识库中统一使用 C2S / S2C 缩写：

| 缩写 | 全称 | 方向 | 矩阵 |
|------|------|------|------|
| C2S | CoeffToSlot | coefficient → slot | slot = U_n · coeff |
| S2C | SlotToCoeff | slot → coefficient | coeff = U_n^-1 · slot |

三层概念区分：
1. **CKKS 编码层**: encoding(σ^-1) / decoding(σ) — 消息与环之间的转换
2. **Bootstrapping 线性变换层**: C2S / S2C — 同一环元素的不同表示转换
3. **BGV/BFV CRT slot ≠ CKKS complex slot** — 不可混用

原始笔记中的 SlotToCoeff / CoeffToSlot 命名如有不同，必须标注为「note convention」，不可直接作为全库标准。
K019 和 K020 已修正方向并标记 needs_human_review。

---

## 目录入口（2026-05-04 归位后）

论文材料已归位至各子目录。进入项目时先读以下 README 了解各自目录内容：
- `01_论文正文/README.md` — 论文初稿与写作材料
- `02_核心构造/README.md` — Split/Merge/C2S/S2C 核心算法
- `03_证明系统/README.md` — 证明草稿与 gap 清单
- `04_实验实现/README.md` — 实验数据与验证
- `05_参数与安全/README.md` — 噪声/安全参数（含作者决策 A1/A2/A3）
- `06_Related_Work/README.md` — 文献定位与 PDF
- `07_图表与表格/README.md` — 交换图与表格
- `08_待确认问题/README.md` — P0-P3 待确认总览
- `09_知识沉淀候选/README.md` — 候选知识（全部待验证）
- `10_发布快照/README.md` — 未来快照说明

> ⚠️ 所有 candidate knowledge 不得作为 validated conclusion 使用。
> 噪声分析 A1/A2/A3 未决策前，噪声方向不确定。

---

## 实验材料规则 (2026-05-04 纳入 ckks-subring-compression)

如果 AI 需要分析实验结果，优先查看：
- `04_实验实现/ckks-subring-compression/实验索引.md`
- `04_实验实现/ckks-subring-compression/4-29/实验结果分析.md`
- `04_实验实现/ckks-subring-compression/4-29新开辟RNS实验/实验结果及其分析/`

**不要**默认使用 `old/` 或 `lattigo-v6-local/` 中的材料作为主实验结论。
这些目录未复制进论文项目，仅存在于原始仓库 `/Users/mac/ckks-subring-compression/`。
