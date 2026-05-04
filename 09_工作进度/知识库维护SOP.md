# 知识库维护 SOP

> 2026-05-04T18:00 | 知识库工程完成阶段

## A. 新资料入库

详见 `新资料入库SOP.md` 完整流程。

**嵌入决策**：只当新资料生成 `rag_status=ready` 的卡片或 Q&A 时才重建 embeddings。详见 `../06_RAG入库准备/RAG_更新决策树.md`。

```
1. 复制到 01_源文件库/ 对应分类目录
2. 在 canonical_source_manifest.jsonl 中分配新 source_id（从 src_000119 开始）
3. 更新 原始路径到库内路径映射.jsonl
4. 判断 source_type：user_note / external_paper / translation_note / paper_planning
5. 判断 domain：abstract_algebra / rlwe_pqc / fhe / paper_project / implementation
6. 轻读判断是否需精读 → 需精读 → 生成卡片和 Q&A
7. 新卡片必须包含：card_id, title, source_id, confidence, knowledge_status, rag_status
8. 新 Q&A 必须包含：qa_id, question, answer, card_id, source_id, rag_status
9. 更新 canonical_card_manifest.jsonl 和 canonical_qa_manifest.jsonl
10. 更新 RAG ready/defer 清单
```

## B. 新论文项目

详见 `未来论文项目入库流程.md` 完整流程。

新论文项目默认采用双层结构（01–10 论文材料 + 11_AI工作区），模板在 `../04_论文协作区/02_未来论文项目模板/`。

### 论文标准接口

每个论文项目必须有：

| 字段 | 类型 | 说明 |
|------|------|------|
| paper_id | string | 唯一标识，如 `Eurocrypt2027` |
| paper_title | string | 论文标题 |
| paper_status | enum | `active` / `paused` / `archived` / `published` |
| main_topic | string | 主 topic，如 `CKKS bootstrapping factorization` |
| project_path | string | `04_论文协作区/论文项目名/` |
| core_idea_ids | [string] | 核心 idea ID 列表 |
| candidate_knowledge_ids | [string] | 关联的 candidate ID 列表 |
| related_source_ids | [string] | 关联的 source_id 列表 |

### 新论文启动流程

```
1. 创建 paper_id
2. 在 04_论文协作区/ 创建论文项目目录
3. 建立 7 个核心入口文件：
   项目总览.md
   AI协作上下文.md
   核心材料索引.md
   当前理解与主线.md
   待办与人工确认.md
   证明状态简表.md
   知识沉淀候选.md
4. 登记核心 idea（写入 idea_id + idea_type）
5. 登记相关 source（关联 source_id）
6. 明确 AI 协作边界（在 AI协作上下文.md 中）
7. 所有新 idea 先进入项目内的 知识沉淀候选.md
8. 同步登记到 11_知识入库接口/candidate_knowledge总表.md
9. 不直接进入正式知识库
10. 不直接进入 RAG ready
```

### 论文 idea → 长期知识库 完整流程

```
论文项目产生 idea
  ↓
写入论文项目的 知识沉淀候选.md
  ↓
同步登记到 11_知识入库接口/candidate_knowledge总表.md
  ↓
标记 paper_id / idea_id / needs_validation
  ↓
作者确认或后续 Batch 验证
  ↓
生成正式知识卡片（card_id）
  ↓
发布到长期知识库（01_抽象代数/02_RLWE/03_FHE 等）
  ↓
更新 canonical_card_manifest.jsonl
  ↓
更新 RAG ready/defer 状态（候选升级为 ready）
```

### 论文 idea 标准接口

```json
{
 "idea_id": "I-Eurocrypt2027-reduced_Q",
 "paper_id": "Eurocrypt2027",
 "idea_title": "reduced-Q leaf bootstrapping 可行性",
 "idea_type": "open_problem",
 "source_path": "04_论文协作区/Eurocrypt2027/知识沉淀候选.md",
 "status": "candidate_knowledge",
 "needs_human_review": true,
 "can_be_promoted_to_kb": false,
 "linked_card_ids": [],
 "linked_source_ids": ["src_000038"],
 "rag_status": "defer"
}
```

### idea_type 枚举

| 值 | 说明 | 举例 |
|----|------|------|
| `construction` | 新构造 / 新算法 | PaCo 结构化密钥 |
| `proof_trick` | 证明技巧 | trace 降维技巧 |
| `noise_analysis` | 噪声分析结论 | leaf error aggregation |
| `parameter_lesson` | 参数经验 | logN=16 的模数选择 |
| `implementation_lesson` | 实现经验 | RNS 分解顺序影响性能 |
| `negative_result` | 失败探索 | slot-preserving 分拆不可行 |
| `related_work_comparison` | 相关工作对比 | PaCo vs HERMES |
| `open_problem` | 未解决问题 | reduced-Q 是否成立 |

### 论文项目间知识复用

```
论文 A 的 validated 知识卡片 → 论文 B 可直接引用（source 可追溯）
论文 A 的 candidate_knowledge → 论文 B 不能直接引用 → 需独立验证
论文 A 的 negative_result → 论文 B 可借鉴（避免重复走弯路）
```

## C. 新卡片创建规则

```
1. card_id 格式：KNNN（3位编号，从 K070 继续）
2. 必须有 source_id（或 source_ids）
3. 必须有 confidence：high / medium / low
4. 必须有 knowledge_status：source_preserved / validated / candidate / support_note
5. 必须有 rag_status：ready / defer / exclude
6. 必须有 canonical_card_path（相对路径）
7. 如果 knowledge_status=candidate → rag_status 自动 = defer
8. 如果 confidence=low → rag_status 自动 = defer
9. 如果 needs_human_review=true → rag_status 自动 = defer
10. JSON+MD 双格式，MD 写知识点，JSON 写结构化 metadata
```

## D. 新 Q&A 创建规则

```
1. qa_id 格式：QNNN（3位编号，从 Q188 继续）
2. 必须有 card_id（尽量）
3. 如果有 card_id → card_id_status=explicit
4. 如果推断 card_id → card_id_status=inferred
5. 如果无法确定 → card_id_status=unresolved（不硬编）
6. 必须有 source_id
7. 必须有 rag_status
8. 不能因为 cross_reference 复制成两个实体
```

## E. RAG 前检查

```
1. 只读 canonical_card_manifest.jsonl 和 canonical_qa_manifest.jsonl
2. 排除 rag_status=defer 和 rag_status=exclude
3. 只使用 canonical_card_path → 排除 duplicate_paths
4. 不处理 card_id_status=unresolved 的 Q&A
5. 不处理 confidence=low 的内容
6. 检查 source_id 是否可回溯 → 不可回溯的不入 RAG
```

## F. 索引维护

```
1. canonical_source_manifest.jsonl：每次导入新资料更新
2. canonical_card_manifest.jsonl：每次新建/修改卡片更新
3. canonical_qa_manifest.jsonl：每次新建/修改 Q&A 更新
4. RAG_ready_defer清单.md：每次 rag_status 变更更新
5. 入库状态看板.md：每次 Batch 完成更新
```

## G. Manifest 一致性审计（每季度或每 50 条变更）

```
1. 检查 source_id 是否全部有 manifest 记录
2. 检查 card_id 是否全部有 canonical_card_path
3. 检查 qa_id 是否无重复
4. 检查 rag_status 是否与卡片字段一致
5. 检查是否有孤儿 source/card/QA
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
