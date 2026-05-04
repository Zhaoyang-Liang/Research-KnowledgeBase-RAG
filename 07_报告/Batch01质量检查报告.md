# Batch-01 质量检查报告

> 检查范围: Batch-01 生成的 14 张知识卡片 + 20 条原子问答
> 检查时间: 2026-05-04

---

## 一、总体评估

| 维度 | 评分 | 说明 |
|------|------|------|
| 代数严谨性 | B | 部分卡片公式正确但上下文过于简略 |
| 术语准确性 | B+ | 主要术语正确，但少数需要区分方案 (BGV vs BFV vs CKKS) |
| 来源追溯性 | A- | 大部分卡片有 source_file 字段 |
| 知识覆盖度 | B | 14 张卡片覆盖了核心领域，但深度不均 |
| 原子问答质量 | B | 20 条问题覆盖主要概念，但部分答案过于简略 |

---

## 二、逐卡检查

### 卡 1: Gentry09/BGV/BFV/CKKS 四方案全对比
- 文件: 02_知识卡片/FHE基础/gentry-BGV-BFV-CKKS对比.json
- 质量: [OK] 保留
- 问题: CKKS 编码部分的 knowledge_point 将 "Encode/Decode" 合在一起表述，可能让初学者困惑 CKKS 的 Decode 与 BGV/BFV 的解密有何不同
- 建议: 拆分 CKKS 的 encoding 和 decryption 为两个知识点

### 卡 2: BGV与CKKS完整同态计算流程
- 文件: 02_知识卡片/FHE基础/BGV-CKKS完整流程.json
- 质量: [WARN] 需要扩写
- 问题:
  1. knowledge_points 只有 3 条，缺少 BFV 流程的独立知识点
  2. 未包含"噪声评估"这一步的判定标准 (如何判断是否超过阈值?)
  3. 文件名说 "BGV与CKKS" 但 knowledge_points 提到了 BFV 却未在 title 中
- 建议: 扩展为 5 条知识点，补 BFV 专用流程

### 卡 3: 密钥切换通用公式与重线性化特例
- 文件: 02_知识卡片/FHE基础/密钥切换通用公式.json
- 质量: [OK] 保留，需标注
- 问题: "c_0' = c_0 + Sum c_j^{(i)} b_{i,j}" 公式中 c_0' 和原始 c_0 用同一符号可能造成混淆 — 这是两轮修改后的 c_0 (加了分解求和)
- 标注建议: 区分 "c_0^{new}" vs "c_0^{old}"

### 卡 4: 模数切换(ModSwitch)与重缩放(Rescale)的区别
- 文件: 02_知识卡片/FHE基础/模数切换vs重缩放.json
- 质量: [WARN] 需要人工确认
- 问题:
  1. 断言 "BGV专用ModSwitch，BFV/CKKS专用Rescale" 过于绝对。CKKS 在某些变体中也使用 ModSwitch
  2. "p = Delta = t" (BFV) 不够精确: BFV 的 Rescale 因子通常是 t (明文模数) 但其是否等于 Delta 取决于具体参数
  3. CKKS 的 Rescale 因子不一定是 Delta，PaCo 论文区分了 "两个不同的缩放因子 p 和 Delta"
- 建议: 标注【需要人工确认】，补充 PaCo 论文中的 p/Delta 区分的注释

### 卡 5: BGV与BFV的核心区别
- 文件: 02_知识卡片/FHE基础/BGV-vs-BFV高位低位.json
- 质量: [OK] 保留
- 问题: 无重大技术错误。"消息在低位/高位"的直觉比喻清晰
- 小建议: 补充 "BGV的ModSwitch链: q_0, q_1 = q_0*Delta, ..." 与 "BFV的Rescale链: q, q/t, q/t^2, ..." 的对比

### 卡 6: 噪声控制、重线性化、密钥切换三者关系
- 文件: 02_知识卡片/FHE基础/噪声控制与密钥切换关系.json
- 质量: [OK] 保留
- 问题: knowledge_points 只有 2 条，过于精简
- 建议: 补充 "ModSwitch vs Rescale vs Bootstrapping 的关系" 和 "为什么噪声控制是目标而密钥切换是工具"

### 卡 7: NTT与Cooley-Tukey快速算法
- 文件: 02_知识卡片/FHE基础/NTT-Cooley-Tukey.json
- 质量: [OK] 保留
- 问题: 缺少 "NTT 条件: n|(q-1)" 这个关键约束
- 建议: 补充此条件对参数选择的影响

### 卡 8: Babai最近平面算法
- 文件: 02_知识卡片/格理论与安全性/Babai最近平面算法.json
- 质量: [WARN] 需要扩写
- 问题:
  1. knowledge_points 仅 2 条，且缺乏与 FHE 的关联说明
  2. 未说明 Babai 算法在 FHE 中主要用于什么 (安全证明? 密钥生成? 还是其他?)
- 建议: 明确 Babai 与 FHE 的关系 (主要用于理解格困难问题, 非直接用于 FHE 构造)

### 卡 9: 分圆多项式在F_p上的分解与槽结构
- 文件: 02_知识卡片/代数基础/分圆多项式Fp分解.json
- 质量: [OK] 保留
- 问题: knowledge_points 中 "d=2: 每个槽=F_{p^2}" 不应写成 "槽中可存p^2范围的值" (准确的说是 p^2 个不同的元素, 不是 "p^2 范围")
- 标注: 【术语不严谨, 需要人工确认】

### 卡 10: 选择分圆多项式的十大理由
- 文件: 02_知识卡片/代数基础/分圆多项式十大优势.json
- 质量: [OK] 保留
- 问题: 无重大错误。但 "判别式小 -> 噪声扩张因子小" 在代数数论中是对的，但在 FHE 的安全分析中可能需要更精确的引用

### 卡 11-14: CKKS 相关卡片
- 典范嵌入与U矩阵.json, SlotToCoeff正向.json, CoeffToSlot反向.json, Trace操作.json
- 质量: [OK] 保留 (这些是用户最核心的翻译材料, 知识卡片作为索引)
- 共享问题: knowledge_points 中的公式直接来自中文翻译, 需要标注哪些公式来自原论文 vs 笔记独立推导

---

## 三、重点问题汇总

### 问题 A: 术语混用 (需要人工确认)

以下术语在卡片中可能被混用，需要逐个确认:

1. "模数切换" vs "Modulus Switching" vs "ModSwitch" — BGV 专用术语，但 CKKS 中也存在 "modulus switching" 作为 rescaling 的一部分
2. "重缩放" vs "Rescaling" vs "Rescale" — BFV 和 CKKS 都用这个词但含义不同 (BFV: 除以 t; CKKS: 除以编码缩放因子)
3. "加密向量" (PaCo 表1/表2) — 这是 CKKS 特有的说法 (加密的是复向量而非整数多项式)
4. "槽" vs "slot" — BGV/BFV 的 CRT slot 与 CKKS 的 canonical embedding slot 在数学上是不同的概念

### 问题 B: 过度简化 (需要修正)

1. 卡片将 BGV/BFV/CKKS 的噪声控制分为两派 (ModSwitch vs Rescale)，但在 PaCo 论文中 CKKS 使用两个不同的缩放因子 p 和 Delta，这与基本 CKKS 的简化模型不同
2. "乘法后必须重线性化" — 实际上在某些优化实现中，可能推迟重线性化以批量处理
3. 卡片没有区分 "CKKS 的 Rescale" 和 "BFV 的 Rescale" 在数学含义上的差异

### 问题 C: 缺少来源位置

大部分卡片有 source_file 但缺少行号或章节号。这在后续检索时会降低可追溯性。

### 问题 D: 数学表述不严谨

部分 knowledge_points 的 formula 字段混杂了中文和公式，如:
- "c=(a*s+t*e+m, -a), ModSwitch: c'=floor((q'/q)*c) (任意q'<q)"
- 这不是可执行的数学公式，而是混合表述

建议: formula 字段应使用纯数学符号，meaning 字段用中文解释

---

## 四、原子问答质量检查

### 总体评价

20 条原子问答中:
- 优秀 (概念清晰、公式正确): 约 12 条
- 良好 (概念正确但需要补充): 约 5 条
- 需要修改: 约 3 条

### 需要修改的原子问答

1. QA-FHE-003: "CKKS 为什么比 BGV/BFV 的 Bootstrapping 更高效" — 答案提到 "CKKS 利用近似多项式近似取模操作 (如 sin 函数近似)"，这在 CKKS 最初的 bootstrapping 论文中正确，但在 PaCo 论文中使用了不同的方法 (blind rotation + 圆群近似)。需要区分
2. QA-FHE-009: 比喻 "像放大 Delta 倍后写在尺子高端" — 比喻生动但有误导性，因为 "放大的不是消息而是编码因子"
3. QA-MATH-004: 需要补充 "x^n+1 仅两项" 在模约化时仅适用于取模运算本身，多项式乘法中的模约化成本还与乘法的实现有关

---

## 五、建议修正动作

| 优先级 | 动作 | 影响的卡片 |
|--------|------|-----------|
| P0 | 标注【需要人工确认】的 5 张卡片 | 模数切换vs重缩放, SlotToCoeff, CoeffToSlot, BGV-vs-BFV, 噪声控制 |
| P1 | 扩写知识卡片 (从 2-3 条 knowledge_point 扩到 4-5 条) | BGV-CKKS流程, 噪声控制, Babai |
| P1 | 区分 CKKS/BFV 中 Rescale 的不同含义 | 模数切换vs重缩放, BGV-vs-BFV |
| P2 | 为所有卡片的 formula 字段统一使用纯数学符号 | 全部 14 张卡片 |
| P2 | 补充 source_file 的章节/行号 | 全部 14 张卡片 |
| P3 | 新增 CKKS vs BGV/BFV 详细对比卡片 | — |
| P3 | 新增 Blind Rotation + GSW 基础卡片 | — |

---

## 六、建议下一步

1. 不继续扩大阅读范围 (按用户要求)
2. 先修正 P0 和 P1 质量问题的卡片
3. 为 14 张卡片生成 Markdown 版本 (便于人类阅读)
4. 开始从原子问答生成 JSONL 格式的 RAG 入库知识片段
5. 待用户确认后，再考虑继续阅读 GHPS 论文的 2.1.4-2.1.5 节 (Duality/Good Bases)
