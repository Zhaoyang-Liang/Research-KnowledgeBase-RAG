# RAG 增量更新流程

> RAG Baseline v1.0 (2026-05-04 frozen)  
> 适用于后续新增资料、卡片或 Q&A 后的 RAG 更新

---

## 标准更新流程（8 步）

### Step 1：新资料进入源文件库

```
~/Desktop/Research-KB-RAG/01_源文件库/
```

- 分配新的 `source_id`（格式 `src_XXXXXX`）
- 更新 `00_导入记录/` 中的导入映射

### Step 2：生成或更新知识卡片 / Q&A

```
~/Desktop/Research-KB-RAG/02_知识卡片/  ← 知识卡片 JSON
~/Desktop/Research-KB-RAG/03_原子问答/  ← 原子问答
```

- 知识卡片按 Batch 规范生成
- 原子问答关联到对应 card_id 和 source_id

### Step 3：判断 ready / defer / exclude

修改知识卡片 JSON 中的 `rag_status` 字段：

| rag_status | 含义 | RAG 行为 |
|------------|------|----------|
| `ready` | 已验证，可检索 | 进入 ready chunks |
| `defer` | 待验证/等待决策 | 不进入 ready chunks |
| `exclude` | 排除 | 不进入 RAG |

### Step 4：更新 canonical manifests

修改以下文件：
- `canonical_card_manifest.jsonl` — 新增/修改卡片
- `canonical_qa_manifest.jsonl` — 新增/修改问答

### Step 5：重建 ready chunks

更新以下 chunk 文件，确保 `rag_status` 字段正确：
- `06_RAG入库准备/chunks/cards_ready_chunks.jsonl`
- `06_RAG入库准备/chunks/qa_ready_chunks.jsonl`

**注意：** 同时同步到 `06_RAG入库准备/api_rag/` 目录（build 脚本从 `chunks/` 读取）。

### Step 6：重建 embeddings + BM25

```bash
cd ~/Desktop/Research-KB-RAG/06_RAG入库准备/api_rag
python3 build_embeddings.py
python3 build_bm25_index.py
```

### Step 7：Smoke test

```bash
python3 query_rag.py "C2S 和 S2C 的方向是什么？" -k 5 -v
python3 query_rag.py "ModSwitch、Rescale、KeySwitch 三者有什么区别？" -k 5 -v
python3 query_rag.py "K063" -k 5 -v
python3 query_rag.py "reduced-Q leaf bootstrapping 是否成立？" -k 5 -v
```

验收标准见 `RAG_Smoke_Test_清单.md`。

### Step 8：更新版本记录

- 如果只是小更新（几张卡片），记录在当前 `RAG_Baseline_v1.0_记录.md` 的版本历史中
- 如果是大更新（新 Batch、大量新卡片），生成新版本记录（如 `RAG_Baseline_v1.1_记录.md`）

---

## 重建触发条件

| 条件 | 需要重建 |
|------|----------|
| chunk 正文 text 变化 | embeddings **必须重建** |
| ready chunks 数量变化 | BM25 **必须重建** |
| 仅 metadata 字段变化（不涉及 text） | embeddings 不需重建 |
| 仅文档说明变化（.md 文件） | 都不需重建 |
| 新增 ready chunk | embeddings + BM25 都需重建 |
| 删除/改 defer 已有 chunk | embeddings + BM25 都需重建 |

---

## 快速命令速查

```bash
# 进入工具目录
cd ~/Desktop/Research-KB-RAG/06_RAG入库准备/api_rag

# 重建 embeddings（需要较长时间，~2-3 分钟）
python3 build_embeddings.py

# 重建 BM25（很快，~1 秒）
python3 build_bm25_index.py

# Smoke test
python3 query_rag.py "C2S 和 S2C 的方向是什么？" -k 5 -v
python3 query_rag.py "K063" -k 5 -v
```

---

## 注意事项

1. **build 脚本从 `chunks/` 目录读取**，不是从 `api_rag/`。确保 `chunks/` 下的 chunk 文件是最新的。
2. **`rag_status` 字段必须存在**：chunk JSON 中缺少此字段会被 build 脚本跳过。
3. **双份同步**：`chunks/` 和 `api_rag/` 下的 chunk 文件需保持一致（通常先更新 `chunks/`，build 后 `api_rag/` 下的 data/ 自动更新）。
4. **备份**：重建 embedding 前建议备份 `data/embeddings.jsonl`。


---

## 关联文档

> 本次全局流程规整 (2026-05-04) 新增/更新了以下核心文档：

| 文档 | 用途 |
|------|------|
| `总流程入口.md` | 知识库唯一全局入口，四层结构速览 |
| `09_工作进度/文档同步总表.md` | 变化类型 → 同步文档 → 是否需要 embedding |
| `06_RAG入库准备/RAG_更新决策树.md` | 什么时候需要/不需要重建 embedding |
| `00_AI接手上下文/AI_无缝接手说明.md` | 全局 AI 接手阅读顺序 |
| `04_论文协作区/02_未来论文项目模板/新论文项目模板.md` | 新论文项目标准结构（含 11_AI工作区） |

> 论文项目采用双层结构：01–10 论文材料区 + 11_AI工作区（AI 接手协作区）。
> RAG 重建规则：只有 ready chunk 文本变化时才需要重建 embedding + BM25。
