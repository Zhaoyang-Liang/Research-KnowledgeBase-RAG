# API-first Hybrid RAG 检索

> 基线：464 个 ready chunks（45 张卡片 + 165 条 Q&A）
>
> 架构：Embedding API（多 provider）+ BM25 + 精确 ID 匹配 → RRF Fusion

## 支持的 Provider

| Provider | 默认模型 | 维度 | 价格 | 推荐 |
|----------|---------|------|------|------|
| **siliconflow** | BGE-M3 (多语言) | 1024 | ¥1.4/1M | ⭐ 推荐（直连 HTTP） |
| siliconflow 备选 | BGE-large-zh-v1.5 | 1024 | ¥0.7/1M | 纯中文更便宜 |
| dashscope | text-embedding-v3 | 1024 | ¥0.7/1M | |
| zhipu | embedding-2 | 1024 | ¥0.5/1M | |
| openai | text-embedding-3-small | 1536 | $0.02/1M | |

## 快速开始

### 1. 安装依赖

```bash
pip install openai numpy scipy rank_bm25
```

### 2. 设置 API Key（以硅基流动为例）

```bash
# 1. 注册 https://cloud.siliconflow.cn/
# 2. 获取 API Key: https://cloud.siliconflow.cn/account/ak
cp config.example.env config.env
# 编辑 config.env:
#   EMBEDDING_PROVIDER=siliconflow
#   EMBEDDING_API_KEY=sk-xxxxxx
source config.env
```

### 3. 构建索引

```bash
pip install openai numpy scipy rank_bm25

# 生成 embeddings
python3 build_embeddings.py

# 构建 BM25 索引
python3 build_bm25_index.py
```

**成本：464 chunks ≈ ¥0.03（硅基流动 BGE-M3）**

### 4. 查询

```bash
# 交互模式
python3 query_rag.py

# 单次查询
python3 query_rag.py "CKKS 的编码管线是什么"

# 指定 top-k
python3 query_rag.py "同态乘法后为什么要重线性化" -k 5

# JSON 输出
python3 query_rag.py "C2S 和 S2C 的区别" --json

# 只返回 chunk_id
python3 query_rag.py "slot-to-coeff" --ids-only

# 详细输出（看完整 text）
python3 query_rag.py "LWE和RLWE" -v
```

## 切换 Provider

修改 `config.env`：

```bash
# 硅基流动（推荐）
EMBEDDING_PROVIDER=siliconflow
EMBEDDING_API_KEY=sk-xxxxxx

# 阿里百炼
# EMBEDDING_PROVIDER=dashscope
# EMBEDDING_API_KEY=sk-xxxxxx

# 智谱AI
# EMBEDDING_PROVIDER=zhipu
# EMBEDDING_API_KEY=xxxxxx.xxxxxx

# OpenAI
# EMBEDDING_PROVIDER=openai
# EMBEDDING_API_KEY=sk-xxxxxx
```

切换后需要重新运行 `python3 build_embeddings.py`。

## 架构

```
Query → Embedding API + BM25 + Exact ID
           ↓
        RRF Fusion (k=60)
           ↓
     Metadata filter (ready only)
           ↓
        Top-k chunks
```

## 安全规则

- API key 从环境变量读取（`EMBEDDING_API_KEY`）
- 不打印、不写入、不硬编码
- 只检索 `rag_status=ready` 的 chunks
- defer / candidate / low confidence / needs_human_review 全部排除

## 目录结构

```
api_rag/
├── README.md
├── config.example.env
├── build_embeddings.py            ← 多 provider embedding 生成
├── build_bm25_index.py            ← BM25 索引构建
├── query_rag.py                   ← 检索入口
├── API_RAG接入说明.md
└── data/
    ├── embeddings.jsonl           ← 生成后
    └── bm25_index.pkl             ← 生成后
```
