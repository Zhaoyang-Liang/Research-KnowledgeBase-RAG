# RAG Baseline v1.0 记录

> 冻结日期：2026-05-04  
> 状态：**FROZEN** — 不再做 metadata 微调、不再扩展新资料

---

## 1. Baseline 名称

**RAG Baseline v1.0**

---

## 2. 技术配置

| 配置项 | 值 |
|--------|-----|
| Embedding provider | SiliconFlow（硅基流动） |
| Embedding model | `Pro/BAAI/bge-m3` |
| Embedding dimension | 1024 |
| Retrieval 架构 | dense embedding + BM25 + exact ID match |
| Fusion 方式 | RRF（Reciprocal Rank Fusion） |
| Storage | 本地 JSONL（embeddings.jsonl）+ BM25 pickle（bm25_index.pkl） |
| Vector DB | 暂不使用 |
| Reranker | 暂不使用 |
| LLM generation | 暂不接入（仅返回检索结果） |

---

## 3. 当前数据状态

| 指标 | 数量 |
|------|------|
| ready chunks 总量 | **467**（302 card chunks + 165 QA chunks） |
| embeddings.jsonl 条数 | **467**（与 ready chunks 一一对应） |
| BM25 chunks 数量 | **467** |
| card_id 条目数 | **46**（cards）+ **47**（QA 中引用） |
| qa_id 条目数 | **165** |
| source_id 条目数 | **59**（cards）+ **64**（QA 中引用） |
| 知识卡片总数 | 69（其中 45 ready / 24 defer） |
| 原子问答总数 | 186（其中 165 ready / 21 defer） |

---

## 4. 已通过的关键测试

以下查询已通过 smoke test，top-3 返回正确/相关结果：

| # | 查询 | 状态 | 说明 |
|---|------|------|------|
| 1 | `C2S 和 S2C 的方向是什么？` | ✅ | K063 top-1/top-2，方向正确 |
| 2 | `CoeffToSlot 和 SlotToCoeff 分别是什么意思？` | ✅ | K044 Q131 + K063 权威定义 |
| 3 | `ModSwitch、Rescale、KeySwitch 三者有什么区别？` | ✅ | K060/K061 相关内容返回 |
| 4 | `K004` | ✅ | exact card_id 命中 |
| 5 | `K044` | ✅ | exact card_id 命中 |
| 6 | `K063` | ✅ | exact card_id 命中 |
| 7 | `src_000059` | ✅ | exact source_id 命中 |
| 8 | `CKKS 编码管线是什么？` | ✅ | 编码管线相关卡片返回 |

---

## 5. C2S/S2C 修复结论

**全库统一采用（PaCo semantic convention，GHPS Eurocrypt 2025）：**

```
C2S = CoeffToSlot = coefficient representation → slot representation
S2C = SlotToCoeff = slot representation → coefficient representation
```

**规则：**

- 旧笔记中相反 convention（SlotToCoeff 误关联 C2S、CoeffToSlot 误关联 S2C）只能作为 `original_note_convention` metadata 字段中的历史记录；
- 任何 ready chunk 正文 text 中**不得**出现完整旧错误表达；
- 权威术语卡片：**K063**（C2S/S2C方向统一 — PaCo语义权威定义）。

---

## 6. 已知但暂不修复的问题

以下问题已识别，但**不是当前 blocker**，暂不修复：

| 问题 | 影响 | 决策 |
|------|------|------|
| 部分 chunk metadata 不完整（source_id/path/domain/topic 空） | 不影响检索精度 | 暂不补 |
| 部分 chunk text 显示为 JSON dict 字符串 | 不影响语义检索 | 暂不美化 |
| top-4/top-5 偶尔出现轻微噪声 | 用户主要看 top-3 | 可接受 |
| 当前论文 open problem 只能检索背景材料，不能给确定结论 | 设计如此 | 非缺陷 |
| reduced-Q / leaf error / formal noise theorem | 属于 defer/candidate，不进入 ready | 等待作者决策 |
| 个别卡片来源说明不完美 | 不影响使用 | 暂不修 |
| K063 与 K063（噪声）编号冲突 | 两张卡片同名 K063 | 暂不解决，K063 语义卡片为术语权威 |

**原则：除非影响数学结论或检索答案，不再继续微调。**

### 6.1 K063 card_id 重号（详细说明）

当前存在两个 K063：

| # | 卡片 | 标题 |
|---|------|------|
| 1 | `K063__C2S与S2C方向统一：PaCo语义权威定义` | C2S/S2C 术语权威 |
| 2 | `K063__CKKS噪声来源总览：加密-乘法-重缩放-KeySwitch-Bootstrap` | CKKS 噪声概览 |

**当前处理策略：**
- 暂不重编号，不继续重建 RAG
- 日常查询 C2S/S2C 方向时，K063 默认指 C2S/S2C 权威术语卡片
- 查询 CKKS 噪声来源时，不要仅用 K063，应使用标题关键词（"CKKS 噪声来源" "KeySwitch 噪声" 等）
- 如果未来正式清理主键，建议保留 K063=C2S/S2C 权威定义，噪声卡片改编号为 K070 或后续未占用编号
- 在正式清理前，任何 AI 不应自动假设 K063 唯一

---

## 版本历史

| 版本 | 日期 | 变更 |
|------|------|------|
| v1.0 | 2026-05-04 | 初始冻结，467 ready chunks，SiliconFlow BGE-M3，BM25+embedding+exact ID |
