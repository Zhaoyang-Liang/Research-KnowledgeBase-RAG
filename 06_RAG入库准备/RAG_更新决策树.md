# RAG 更新决策树

> 版本: v1 | 生成: 2026-05-04
> 面向: 不太懂技术的作者 / 接手 AI
> 一句话: "做了某个操作之后，要不要重建 embedding / BM25？看这里。"

---

## 决策树（逐条判断）

### ❌ 不需要重建 embedding

以下操作**不会**影响 RAG 检索结果，无需重建：

- 写 README
- 写 SOP
- 移动论文项目文件
- 复制论文项目文件
- 更新 AI 工作区（11_AI工作区/）
- 新增 candidate_knowledge
- 新增未验证论文材料
- 新增实验原始日志（.txt, .go output）
- 整理目录结构
- 更新流程文档（09_工作进度/）
- 修改模板（02_未来论文项目模板/）

---

### ✅ 必须重建 embedding + BM25

以下操作**会**改变 ready chunk 内容或数量，必须重建：

| 操作 | 原因 |
|------|------|
| ready card 正文变了 | chunk 文本变了 → embedding 变了 |
| ready Q&A 正文变了 | chunk 文本变了 → embedding 变了 |
| ready chunks 重新生成了 | 内容或数量变了 |
| 新增了 ready card | 新 chunks 需要 embedding |
| 新增了 ready Q&A | 新 chunks 需要 embedding |
| 删除或排除了 ready chunk | chunks 数量变了 |
| 修正了会影响检索答案的术语错误 | chunk 文本变了 |
| 修改了 C2S/S2C 这类核心术语的 chunk | 会影响所有相关检索 |

---

### 🟡 特殊情况

| 操作 | 说明 |
|------|------|
| 只改 embeddings 模型（如从 BGE-M3 换到 text-embedding-3-small） | ✅ 需要重建 embeddings，BM25 不必重建 |
| 只改 `query_rag.py` 的排序/融合逻辑 | ❌ embeddings 和 BM25 都不必重建 |
| 只改 RRF 参数（如 k 值） | ❌ 都不必重建 |
| chunks 不变，只修 metadata 字段 | ❌ 不必重建（metadata 不影响向量检索） |

**经验法则：看有没有 .txt 文件变了。** 只有 ready chunks 的文本内容变化才需要重建。

---

## 标准重建命令

```bash
cd ~/Desktop/Research-KB-RAG/06_RAG入库准备/api_rag

# 1. 重建 embeddings（SiliconFlow BGE-M3）
python3 build_embeddings.py

# 2. 重建 BM25 索引
python3 build_bm25_index.py
```

---

## 重建后必须 smoke test

```bash
cd ~/Desktop/Research-KB-RAG/06_RAG入库准备/api_rag

# 术语方向（验证 C2S/S2C 正确）
python3 query_rag.py "C2S 和 S2C 的方向是什么？" -k 5 -v

# 核心概念区分（验证噪声/模态/密钥术语）
python3 query_rag.py "ModSwitch、Rescale、KeySwitch 三者有什么区别？" -k 5 -v

# exact ID 查询（验证卡片可检索）
python3 query_rag.py "K063" -k 5 -v

# 安全墙（验证 defer/candidate 被正确拦截）
python3 query_rag.py "reduced-Q leaf bootstrapping 是否成立？" -k 5 -v
```

如果上面 4 条都通过，RAG 即恢复正常。

---

## 什么时候重建之后的 RAG baseline 需要更新版本号？

| 变化程度 | 处理 |
|----------|------|
| 新增 < 5 个 ready chunks | 在 `RAG_Baseline_v1.0_记录.md` 末尾追加一条记录即可 |
| 新增 ≥ 5 个 ready chunks | 更新版本号为 v1.1，写变更摘要 |
| 重建全部 embeddings（模型更换等） | 更新版本号为 v2.0，完整记录 |
