# API-first RAG 接入说明

> 2026-05-04 | 检索验证阶段
> 
> ⚠️ 当前仅做 chunk 检索，**不调用聊天模型生成最终答案**。

## 一、设计目标

当前阶段的唯一目标：**验证检索质量**。

```
用户提问 → hybrid 检索 → 返回 top-k chunks（不生成答案）
```

后续在检索质量达标后，再接入聊天模型做 RAG 问答。

## 二、检索流程详解

### Step 1: 查询预处理

用户输入查询，不做改写、不做 query expansion，原样送入三个检索通道。

### Step 2: 三通道检索

| 通道 | 方法 | 特点 |
|------|------|------|
| Embedding | OpenAI text-embedding-3-small → cosine similarity | 语义匹配，中文友好 |
| BM25 | rank_bm25 (中文单字 + 英文缩写 tokenizer) | 精确匹配 CKKS/C2S/K004 |
| Exact ID | card_id / qa_id / source_id 字符串包含检测 | 直接命中指定文档 |

### Step 3: RRF Fusion

```
RRF_score(chunk) = Σ 1/(k + rank_i)  其中 k=60
```

三个通道各自排序后，按 RRF 公式融合分数。精确 ID 匹配给予 +2.0 boost。

### Step 4: Metadata Filter

```sql
WHERE rag_status = 'ready'
  AND needs_human_review = false
  AND confidence != 'low'
```

此过滤在 chunk 写入阶段就已保证，检索时再做一次 safety check。

### Step 5: 去重

同一 card_id 的 `card_knowledge_point` 类 chunk 只保留 RRF 分数最高的 1 条。`card_summary` 和 `qa_pair` 不受此限制。

## 三、Tokenizer 设计

中文 + 数学 + 缩写混合场景的 tokenizer：

```
输入: "CKKS 中 σ⁻¹ 将复数向量编码为多项式"
输出: ["ckks", "中", "σ⁻¹", "将", "复", "数", "向", "量", "编", "码", "为", "多", "项", "式"]
```

关键设计：
- **英文缩写保持完整**：`CKKS` → `["ckks"]`，不切分为 `["c","k","k","s"]`
- **中文逐字切分**：`编码` → `["编","码"]`（BM25 对单字 n-gram 效果更好）
- **数学符号保持**：`σ⁻¹` → `["σ⁻¹"]`，Unicode 上下标视为符号的一部分
- **ID 保持完整**：`K004` → `["k004"]`，`src_000018` → `["src_000018"]`

## 四、输出格式

### 默认输出（表格）

```
---  ---
  #    score  source           chunk_id  card/qa       source_id       text
---  ---
  1   0.0856  hybrid           C000012   K004          src_000001      分式域与分裂域的区别: 分式域是包含一个整环...
  2   0.0672  embedding        C000088   Q108          src_000015      Q: CKKS中Trace操作如何用来计算同态内积？...
  3   0.0543  bm25             C000201   K005          src_000002      分裂域的构造: 设f(x)为不可约多项式...
---  ---
```

### JSON 输出（`--json`）

```json
[
  {
    "rank": 1,
    "score": 0.0856,
    "retrieval_source": "hybrid",
    "chunk_id": "C000012",
    "chunk_type": "card_knowledge_point",
    "card_id": "K004",
    "qa_id": "",
    "source_id": "src_000001",
    "domain": "抽象代数",
    "topic": "域论/分式域与分裂域",
    "confidence": "high",
    "canonical_card_path": "03_FHE知识库/文档知识卡片/K004__...json",
    "text": "分式域与分裂域的区别: ..."
  }
]
```

### ID-only 输出（`--ids-only`）

```
C000012  card=K004  qa=-
C000088  card=-     qa=Q108
C000201  card=K005  qa=-
```

## 五、使用场景

### 场景 1：快速验证某个概念是否在 KB 中

```bash
python3 query_rag.py "field switching trace" --ids-only
# 检查返回的 card_id 是否覆盖所需知识
```

### 场景 2：论文写作时查找相关材料

```bash
python3 query_rag.py "CKKS 模数切换和重缩放的区别是什么" -k 20 -v
# 详细查看每条结果的完整 text
```

### 场景 3：批量测试检索质量

```bash
for q in "分式域和分裂域区别" "CKKS编码" "C2S方向" "BGV vs BFV" "重线性化必要性"; do
  echo "=== $q ==="
  python3 query_rag.py "$q" -k 5
  echo
done
```

### 场景 4：脚本调用

```bash
# Python 脚本中调用
result=$(python3 query_rag.py "LWE RLWE 区别" --json)
card_ids=$(echo "$result" | python3 -c "import json,sys; [print(d['card_id']) for d in json.load(sys.stdin)]")
```

## 六、安全规则实施

| 规则 | 实施方式 |
|------|---------|
| API key 不写进代码 | `os.environ.get("OPENAI_API_KEY")` |
| 不打印 API key | 代码中无任何 print(api_key) |
| 不写入日志 | 使用 print 而非 logging |
| 只处理 ready chunks | `chunk["rag_status"] == "ready"` |
| defer 保护 | 加载时即过滤，检索后二次检查 |
| candidate 保护 | 同 defer |
| low confidence 保护 | 同 defer |

## 七、常见问题

**Q: 为什么不直接生成答案？**
A: 当前阶段目标是验证检索质量。等确认 top-5 的相关性足够好后，再加聊天模型。

**Q: 检索结果不相关怎么办？**
A: 检查 `retrieval_source` 字段：
- 如果都是 `bm25` → 可能需要调整 query 的词，或者 embedding 未生效
- 如果都是 `embedding` → 语义匹配可能不如预期，尝试加更多关键词
- 如果 score 都 < 0.01 → 可能 KB 中确实没有相关内容

**Q: 增量更新后需要重新生成吗？**
A: 
- 新增 chunks → 只需对新 chunks 调用 embedding API，追加到 `embeddings.jsonl`；重新运行 `build_bm25_index.py`
- defer → ready 升级 → 同上，对新 ready 的 chunks 生成 embedding

**Q: 可以换模型吗？**
A: 可以。修改 `config.env` 中的 `EMBEDDING_MODEL`，然后重新运行 `build_embeddings_openai.py`。注意不同模型维度不同，需要全部重新生成。

## 八、下一步

检索质量验证达标后：

1. 添加聊天模型（OpenAI / Claude / DeepSeek）
2. 将 top-k chunks 注入 prompt context
3. 生成带 citation 的最终答案
4. 支持对话历史
5. 添加 RAG 质量评估（人工标注 + 自动指标）
