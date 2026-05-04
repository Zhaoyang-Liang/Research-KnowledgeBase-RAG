# Batch-05b2-pre：噪声分析双层盘点

> 2026-05-04T16:45 | 仅盘点，不推导、不完成噪声分析

## Part A：通用噪声知识库盘点

### A.1 已导入源文件（通用噪声知识相关）

| source_id | 文件 | 类型 | 成熟度 | 适合长期知识库 |
|-----------|------|------|--------|----------------|
| src_000015 | 噪声控制-密钥切换.md | 用户笔记 | user_note_mature | ✅ 是 |
| src_000014 | 模数切换和重缩放的区别.md | 用户笔记 | user_note_mature | ✅ 是 |
| src_000013 | 密钥切换与重线性化.md | 用户笔记 | user_note_mature | ✅ 是 |
| src_000104 | 为什么以及什么时候需要密钥切换.md | 用户笔记 | user_note_mature | ✅ 是 |
| src_000105 | 模数切换与密钥交换区别.md | 用户笔记 | user_note_mature | ✅ 是 |
| src_000106 | 模数切换和重缩放的区别.md | 用户笔记 | user_note_mature | ✅ 是 |
| src_000059 | BGV-BFV Bootstrapping.pdf (229页) | 外部讲义 | standard_reference | ✅ 是（梁朝阳讲义） |
| src_000039 | ckks_parameter_tutorial.pdf | 外部论文 | standard_reference | ✅ 是 |

### A.2 现有知识卡片（噪声相关）

| card_id | 标题 | 覆盖 | 成熟度 |
|---------|------|------|--------|
| K008 | FHE方案对比 | 各方案噪声控制方式概览 | medium（用户笔记） |
| K037 | BGV/BFV Bootstrapping讲义总览 | bootstrapping noise通用 | standard_reference |
| K038 | BGV Bootstrapping完整流程 | digit extraction noise | standard_reference |
| K056 | FHE方案总览 | 各方案噪声公式 | medium（用户笔记总结外部论文） |
| K057 | BGV vs BFV 高低位 | ModSwitch/Rescale impact | medium + needs_review |
| K058 | 同态计算完整流程 | 噪声控制时机表 | medium |
| K059 | 密钥切换三场景 | KeySwitch噪声引入点 | medium |
| K060 | ModSwitch vs KeySwitch | 噪声影响区分 | medium |
| K061 | ModSwitch vs Rescaling | Rescale误差 vs ModSwitch降噪 | medium + needs_review |
| K062 | FHE方案层与论文关系 | CKKS vs BGV/BFV噪声差异 | medium + needs_review |

### A.3 现有原子问答（噪声相关）

| qa_id | 问题 | 领域 |
|-------|------|------|
| Q124 | CKKS编码中Δ的作用/误差 | CKKS编码（非噪声专项） |
| Q125 | CKKS complex slot vs CRT slot | 编码（非噪声） |
| Q126 | CKKS approximate性质 | 精度（非噪声界限） |
| Q160-Q175 | Batch-05b1 Q&A（KeySwitch/ModSwitch/Rescale/bootstrapping） | 方案操作 |

### A.4 通用噪声知识缺口（长期知识库应补充但暂缺）

| 缺口主题 | 现有材料覆盖度 | 原因 |
|-----------|---------------|------|
| CKKS噪声来源系统化总览 | 分散在 K056-K061 中，未整合 | 需要一张总览卡 |
| KeySwitch噪声的精确界 | K059提及场景，K060提噪声引入，但无通用界公式 | 可从 RLWE 文献补充 |
| Rescale/scale error的CKKS表达 | K061在三种方案context中，CKKS部分不突出 | 需要dedicated CKKS rescale error卡 |
| EvalMod近似误差的一般概念 | 仅在K062/Q171中提及，无独立卡 | 当前材料多在论文专属文件中 |
| Bootstrapping error的通用组成 | 不系统 | 需整理 generic composition |
| BGV/BFV/CKKS噪声机制系统化差异 | 分散在多张卡中 | 可整合为一张对比卡 |

### A.5 可生成知识卡片主题（候选）

| 候选主题 | 源材料 | 成熟度 | 建议 |
|-----------|--------|--------|------|
| CKKS噪声来源总览 | src_000015 + K056-K061 | user_note_mature | ✅ 可生成（medium confidence） |
| KeySwitch噪声来源与界 | src_000015 + K059 + K060 | user_note_mature | ✅ 可生成（medium confidence） |
| EvalMod近似误差 | src_000038（论文专属）+ CKKS 2017 | paper_candidate + standard | ❌ 暂不生成（混杂论文专属内容） |

### A.6 通用噪声知识进入长期知识库的准入标准

- ✅ 来自 CKKS 2017 / 成熟文献的稳定公式
- ✅ 来自用户成熟笔记（经多轮核对）
- ✅ 不依赖当前论文 construction details
- ❌ 包含当前论文 split/merge 特定噪声
- ❌ 假设当前论文的参数选择
- ❌ 混用 BGV/BFV 的噪声概念到 CKKS

---

## Part B：双层接口——通用知识如何服务于当前论文噪声分析

### B.1 可从长期知识库直接复用的内容

| 通用知识 | 在论文中的用途 |
|----------|---------------|
| CKKS加密噪声基础 | 论文输入噪声 e_0 的来源与分布 |
| KeySwitch噪声机制 | 论文中 Split/Merge/C2S/S2C 内部所有 KS 步骤的噪声模型 |
| Rescale误差 | 论文中 ScaleDown / leaf rescale 的 rounding error |
| CKKS乘法噪声增长 | 论文中 EvalMod 内 polynomial evaluation 的噪声累积 |
| 重线性化 = KS特例 | 论文中 EvalMod relinearization 的噪声模型 |
| CKKS bootstrapping通用流程 | 论文中 ScaleDown→ModUp→core→output 的框架 |

### B.2 通用知识无法覆盖、必须由论文专属分析的内容

| 论文专属问题 | 原因 |
|-------------|------|
| Split_k 系数重排是否引入新噪声 | 不是标准 CKKS 操作 |
| Merge_k 是否引入额外误差 | 不是标准 CKKS 操作 |
| leaf EvalMod 输入范围 K_leaf | 依赖论文的 secret distribution 和 split 策略 |
| leaf bootstrapping 参数 (d, r, Δ) 选择 | 需对 n = N/k 重新设计 |
| k 个 leaf 的 union bound | 取决于 k 和论文希望的安全级别 |
| γ_impl 是否进入数学 noise bound | 论文实现细节 |
| reduced-Q 可行性 | 论文特定参数探索 |
| 实验误差数据是否能支撑 theorem statement | 论文贡献判断 |

### B.3 不要污染长期知识库的红色线

1. ❌ 不要把 src_000038（噪声与成功概率分析.md）中的分析当成 validated knowledge 卡片
2. ❌ 不要把 Split/Merge noise 当作通用 FHE 知识
3. ❌ 不要把 γ_impl 补偿逻辑当成 CKKS 标准机制
4. ❌ 不要把「实验 5.03× speedup + 10⁻⁷ error」当作数学 noise bound 的证据
