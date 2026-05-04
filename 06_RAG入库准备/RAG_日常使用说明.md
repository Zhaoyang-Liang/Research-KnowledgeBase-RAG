# RAG 日常使用说明

> RAG Baseline v1.0 (2026-05-04 frozen)  
> 知识库根目录：`~/Desktop/Research-KB-RAG/`

---

## 1. 基本查询命令

```bash
cd ~/Desktop/Research-KB-RAG/06_RAG入库准备/api_rag
python3 query_rag.py "你的问题" -k 5 -v
```

参数说明：
- `-k N`：返回前 N 条结果（默认 5）
- `-v`：显示每条结果的完整 text

---

## 2. 推荐使用方式

| 场景 | 参数 | 说明 |
|------|------|------|
| 普通概念查询 | `-k 5` | top-3 为主要参考 |
| 复杂论文问题 | `-k 10` | 需要更多上下文 |
| 查具体卡片 | 直接输入 `Kxxx` | exact ID match，例如 `K004`、`K063` |
| 查具体问答 | 直接输入 `Qxxx` | exact ID match，例如 `Q131` |
| 查具体来源 | 直接输入 `src_xxxxxx` | exact ID match，例如 `src_000059` |
| 结果深度查看 | `-v` | 显示完整 chunk text |

**结果解读原则：**
- **主要参考 top-3**：通常已包含最相关的知识点
- **top-4/top-5 作为补充**：可能包含相关背景或交叉引用
- top-5 允许少量噪声，top-3 应基本相关

---

## 3. 示例查询

### 概念类
```bash
python3 query_rag.py "C2S 和 S2C 的方向是什么？" -k 5 -v
python3 query_rag.py "CKKS 编码管线是什么？" -k 5 -v
python3 query_rag.py "ModSwitch、Rescale、KeySwitch 三者有什么区别？" -k 5 -v
python3 query_rag.py "CKKS 的 slot 和 BGV/BFV 的 CRT slot 有什么不同？" -k 5 -v
python3 query_rag.py "KeySwitching 的通用公式是什么？" -k 5 -v
```

### 论文相关
```bash
python3 query_rag.py "R 和 R^∨ 在 GHPS 和当前论文中有什么区别？" -k 5 -v
python3 query_rag.py "reduced-Q leaf bootstrapping 是否成立？" -k 5 -v
python3 query_rag.py "coefficient split 和 slot-preserving split 的区别" -k 5 -v
```

### Exact ID 查询
```bash
python3 query_rag.py "K063" -k 5 -v
python3 query_rag.py "src_000059" -k 5 -v
python3 query_rag.py "K004" -k 5 -v
```

---

## 4. 如何把 RAG 结果交给 LLM

将 `query_rag.py` 的输出（top-k 结果）粘贴给 LLM 时，使用以下标准提示词模板：

```
请基于以下 RAG 检索结果回答问题。

约束：
- 不要使用检索结果之外的结论；
- 如果材料不足，请明确说「材料不足」；
- 如果问题涉及当前论文 candidate/open problem，不要给确定结论，
  只能总结已有材料和待确认点；
- 如果是 needs_human_review 的内容，需标注「以下结论来自待审核条目」；
- 如果是 low confidence 的内容，需标注「以下仅为初步直觉，未验证」。

RAG 检索结果：
[粘贴 query_rag.py 输出]
```

---

## 5. 当前保护规则

以下内容**不能**作为确定结论使用：

| 标记 | 含义 | 使用规则 |
|------|------|----------|
| `candidate_knowledge` | 候选知识，未验证 | 不作为 validated conclusion |
| `needs_human_review` | 需人工审核 | 不作为确定结论，需标注 |
| `low` confidence | 低置信度 | 只能作为 intuition |
| 当前论文 open problem | 论文未决问题 | 必须提示「未验证」 |
| `original_note_convention` | 旧笔记 convention | 仅历史参考，不用作正式定义 |

---

## 6. 快速命令速查

```bash
# 进入工具目录
cd ~/Desktop/Research-KB-RAG/06_RAG入库准备/api_rag

# 概念查询（默认 top-5）
python3 query_rag.py "你的问题"

# 详细查询（top-10 + 完整文本）
python3 query_rag.py "你的问题" -k 10 -v

# Exact ID 查询
python3 query_rag.py "K063" -k 5 -v

# 重建索引（新增内容后）
python3 build_embeddings.py
python3 build_bm25_index.py
```


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
