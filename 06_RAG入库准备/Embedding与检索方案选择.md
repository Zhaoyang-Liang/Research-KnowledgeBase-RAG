# Embedding 与检索方案选择（效果优先版）

> 2026-05-04T18:45 | 约束变更后重新评估
>
> 隐私不再是主要限制 → 效果优先 + 易用性优先

## 一、新约束

| 维度 | 优先级 | 说明 |
|------|--------|------|
| 检索准确率 | ★★★★★ | 最重要 |
| 中文+数学+FHE术语 | ★★★★★ | 核心场景 |
| 易用性 | ★★★★ | 快速跑通 |
| 后续扩展 | ★★★★ | 多论文项目 |
| 工程复杂度 | ★★★ | 适中即可 |
| 成本 | ★★ | 小规模可接受 |
| 隐私 | ★ | 非主要限制 |

## 二、基线数据

| 指标 | 数量 |
|------|------|
| Ready 卡片 | 45 |
| Ready Q&A | 165 |
| Ready chunks | 464 |
| Defer/candidate/low | 已隔离（不进入检索） |

## 三、三种方案重新比较

### 方案 A: API Embedding + BM25 + 本地 JSON

```
OpenAI text-embedding-3-small
+ rank_bm25
+ embeddings.jsonl (本地 JSON lines)
+ metadata filter (代码层)
```

**评分：**

| 维度 | 评分 | 说明 |
|------|------|------|
| 检索效果 | ★★★★★ | OpenAI embedding 质量顶级，BM25 补缩写盲区 |
| 中文能力 | ★★★★ | text-embedding-3-small 多语言优秀 |
| FHE 术语 | ★★★★ | 纯 embedding 不认 CKKS/C2S，但 BM25 补齐 |
| 易用性 | ★★★★★ | 3 个 pip 包，~200 行代码，10 分钟跑通 |
| 扩展性 | ★★★★ | JSONL 按行追加即可增量 |
| 成本 | ★★★★★ | ~$0.0005 首次，~$0.000001/次查询 |
| 维护 | ★★★★ | 零运维，无数据库 |

**适用：** 当前 464 chunks 的最优选择。

---

### 方案 B: API Embedding + BM25 + LanceDB/Chroma

```
OpenAI text-embedding-3-small
+ rank_bm25
+ LanceDB (向量库 + metadata filter)
```

**评分：**

| 维度 | 评分 | 说明 |
|------|------|------|
| 检索效果 | ★★★★★ | 同方案 A |
| 易用性 | ★★★ | 多一个向量库依赖，API 学习成本 |
| 扩展性 | ★★★★★ | 原生支持增量、filter、分区 |
| 成本 | ★★★★★ | 同方案 A |
| 维护 | ★★★ | 需要管理 LanceDB 文件 |

**适用：** 当 chunks 增长到 5000+ 时明显优于方案 A。

**对比方案 A：**

| 差异点 | 方案 A (JSON) | 方案 B (LanceDB) |
|--------|-------------|-----------------|
| 向量存储 | JSONL 文件（2.8 MB） | LanceDB 列式存储 |
| 检索查询 | Python 遍历计算 cosine | LanceDB 内置 ANN |
| Metadata filter | 代码中手动过滤 | SQL WHERE 子句 |
| 多 paper 分区 | 手动管理多个 JSONL | 原生 partition |
| 1000 chunks 检索延迟 | ~5ms | ~2ms |
| 5000 chunks 检索延迟 | ~25ms | ~5ms |
| 50000 chunks 检索延迟 | ~250ms（需优化） | ~10ms |

**结论：** 当前 464 chunks → 方案 A 足够。chunks > 5000 → 迁移到方案 B。

---

### 方案 C: 本地 BGE + BM25 + LanceDB

```
BGE-large-zh-v1.5 (本地 1.3 GB)
+ rank_bm25
+ LanceDB
```

**评分：**

| 维度 | 评分 | 说明 |
|------|------|------|
| 检索效果 | ★★★★ | BGE 中文顶级，但略逊于 OpenAI |
| 易用性 | ★★★ | 需下载模型、管理 sentence-transformers |
| 扩展性 | ★★★★ | 同方案 B |
| 成本 | ★★★★★ | $0 |
| 维护 | ★★ | 模型更新、依赖管理 |

**适用：** 当 AP I不可用、或需要完全离线时。

**对比方案 A：**

| 差异点 | 方案 A (OpenAI) | 方案 C (BGE) |
|--------|----------------|-------------|
| 效果 | ★★★★★ | ★★★★ |
| 启动时间 | 0（无模型下载） | 首次 ~2 min 下载 |
| 首次 embedding 耗时 | ~5s（API 调用） | ~10s（本地 CPU） |
| 每次查询 embedding | ~0.2s（API 网络） | ~0.05s（本地） |
| 离线能力 | ❌ | ✅ |
| 中文 FHE 术语 | 通用模型，不确定 | BGE 可 fine-tune |
| 模型升级 | API 自动升级 | 手动下载新版本 |

---

## 四、最终推荐（效果优先约束下）

### ⭐ 推荐：方案 A（API Embedding + BM25 + 本地 JSON）

```
OpenAI text-embedding-3-small + rank_bm25 + embeddings.jsonl
```

### 逐项回答

**1. API 还是本地 embedding？**

→ **API（OpenAI text-embedding-3-small）**。效果最佳、零模型管理、即开即用。成本可忽略。

**2. 当前 464 chunks 是否需要 LanceDB？**

→ **不需要**。JSON 遍历 cosine similarity 足够（~5ms）。等到 5000+ chunks 再迁移。

**3. 可以先用本地 JSON 保存向量吗？**

→ **可以，推荐**。每行 `{"chunk_id": "C000001", "embedding": [0.123, -0.456, ...]}`。464 chunks × 1536 维 ≈ 2.8 MB。增量追加即可。

**4. BM25 是否仍然保留？**

→ **必须保留**。CKKS/C2S/slot-to-coeff/K004/src_000018/τ_n⁻¹ 这些缩写和 ID，纯 embedding 无法精确匹配。

**5. 是否需要 reranker？**

→ **暂时不需要**。464 chunks 量级太小，BM25 + embedding 的 RRF fusion 已足够。等到 5000+ chunks 且 top-10 中出现大量相近概念（如分式域 vs 分裂域）时再加。

**6. 是否应该先只返回 top-k chunks？**

→ **是**。当前阶段只验证检索质量。不在检索达标前接聊天模型。

**7. 下一步应该生成哪些脚本？**

→ 已生成（见下方文件清单）。

---

## 五、已生成文件

```
06_RAG入库准备/api_rag/
├── README.md                         # 快速开始
├── config.example.env                # API key 模板
├── build_embeddings_openai.py        # 调用 OpenAI API 生成向量
├── build_bm25_index.py               # 构建 BM25 索引
├── query_rag.py                      # 检索入口
├── API_RAG接入说明.md                # 详细文档
└── data/                             # 生成后
    ├── embeddings.jsonl              # chunk_id + embedding
    └── bm25_index.pkl                # BM25 索引 + ID 查找表
```

---

## 六、成本估算

| 操作 | tokens | 成本 |
|------|--------|------|
| 首次构建 (464 chunks) | ~23K | ~$0.0005 |
| 每次查询 | ~50 | ~$0.000001 |
| 100 次查询 | ~5K | ~$0.0001 |
| 全部重建 (5000 chunks) | ~250K | ~$0.005 |

模型：text-embedding-3-small, $0.02/1M tokens

---

## 七、迁移路径

```
当前 (464 chunks)
  → 方案 A: API + BM25 + JSON    ← 现在
  ↓
未来 (5000+ chunks)
  → 方案 B: API + BM25 + LanceDB ← 数据量大时迁移
  ↓
如果需要离线
  → 方案 C: BGE + BM25 + LanceDB ← 备选
```

迁移成本极低：方案 A → 方案 B 只需把 embeddings.jsonl 写入 LanceDB（~10 行代码）。

---

## 八、何时重新评估

| 触发条件 | 评估内容 |
|---------|---------|
| chunks > 5000 | 考虑迁移到 LanceDB |
| PKB 中出现大量相近概念混淆 | 加 reranker |
| 需要离线部署 | 切换到本地 BGE |
| 检索 Recall@5 < 0.8 | 检查 chunk 切分策略、query expansion |
| 新论文项目启动 | 添加 paper_id 分区 |
