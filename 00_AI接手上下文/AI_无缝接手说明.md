# AI 无缝接手说明

> 知识库根目录：`~/Desktop/Research-KB-RAG/`  
> 版本：RAG Baseline v1.0（2026-05-04 frozen）  
> 目标读者：任何接手这个知识库的 AI（DeepSeek / GPT / Claude 等）

---

## 1. 一句话说明

这是一个围绕 **FHE / CKKS / bootstrapping / 当前 Eurocrypt 2027 论文项目** 建立的研究知识库 + RAG 检索系统 + 论文协作区。

---

## 2. 当前状态

| 状态项 | 说明 |
|--------|------|
| 知识库批次整理 | 已完成 Batch-01/03/04/05a1/05a2/05a3/05b1/05c |
| 待处理批次 | Batch-05b2（噪声深入）暂停，等待作者决策；Batch-06（实验数据）暂缓 |
| RAG | Baseline v1.0 **已可用**：467 ready chunks，SiliconFlow BGE-M3，BM25+embedding+exact ID |
| 工程整理 | **已停止**，不继续 metadata 微调，不扩展新资料 |
| 当前论文 | proof/noise/security 仍需作者亲自确认 |
| AI 分工建议 | DeepSeek 日常整理；GPT 关键数学推导与 proof 检查 |

---

## 3. 最高优先级阅读顺序

接手后**必须先读**（按顺序）：

| # | 文件 | 用途 |
|---|------|------|
| 1 | `总流程入口.md` | **知识库唯一全局入口**，四层结构速览 |
| 2 | `00_AI接手上下文/README.md` → `USER.md` → `IDENTITY.md` → `SOUL.md` → `MEMORY.md` | 用户身份、研究偏好、Agent 人设、长期记忆 |
| 3 | `00_AI接手上下文/AI_无缝接手说明.md` | 你正在读的这份（已从根目录移入此目录） |
| 4 | `00_资料清单/资料清单总览.md` | 知识库全貌和各目录用途 |
| 5 | `06_RAG入库准备/RAG_Baseline_v1.0_记录.md` | RAG 技术配置和已知问题 |
| 6 | `06_RAG入库准备/RAG_日常使用说明.md` | RAG 查询命令和 LLM 提示词模板 |
| 7 | `06_RAG入库准备/RAG_更新决策树.md` | 什么时候需要/不需要重建 embedding |
| 8 | `04_论文协作区/01_当前论文项目/Eurocrypt2027/README.md` | 当前论文项目总览 |
| 9 | `04_论文协作区/01_当前论文项目/Eurocrypt2027/11_AI工作区/README.md` | 论文项目 AI 工作区入口 |
| 10 | `04_论文协作区/01_当前论文项目/Eurocrypt2027/11_AI工作区/00_接手说明/` | 项目总览、AI协作上下文、核心材料索引 |
| 11 | `04_论文协作区/01_当前论文项目/Eurocrypt2027/11_AI工作区/01_项目理解/` | 当前理解与主线、论文主线重构、写作路线 |
| 12 | `04_论文协作区/01_当前论文项目/Eurocrypt2027/11_AI工作区/02_待办与人工确认/` | 待办总表、待确认问题、待补证明缺口 |
| 13 | `04_论文协作区/01_当前论文项目/Eurocrypt2027/11_AI工作区/04_证明与噪声检查/` | 噪声分析决策、证明状态 |

---

## 4. 当前 RAG 使用方法

```bash
cd ~/Desktop/Research-KB-RAG/06_RAG入库准备/api_rag
python3 query_rag.py "你的问题" -k 5 -v
```

- 普通问题用 `-k 5`，复杂论文问题用 `-k 10`
- **主要看 top-3**，top-4/top-5 作为补充
- 支持 exact ID 查询：`Kxxx`、`Qxxx`、`src_xxxxxx`
- reduced-Q / leaf error / noise theorem 等不能直接给确定结论

**把 RAG 结果交给 LLM 时的标准提示词模板**见 `06_RAG入库准备/RAG_日常使用说明.md` §4。

---

## 5. 当前最重要的术语规则

**全库统一采用（PaCo semantic convention）：**

```
C2S = CoeffToSlot = coefficient representation → slot representation
S2C = SlotToCoeff = slot representation → coefficient representation
```

旧笔记中的相反 convention（SlotToCoeff 误关联 C2S、CoeffToSlot 误关联 S2C）只能作为 `original_note_convention` metadata 字段中的历史记录，不作为正式定义。

权威卡片：`K063__C2S与S2C方向统一：PaCo语义权威定义`

---

## 6. 当前论文项目边界

当前论文研究 **CKKS bootstrapping core 的子环/小环/leaf 分解结构**。

核心主线：
```
Split_k → ⊕(C2S_n → EvalMod_n → S2C_n) → Merge_k
```

在 bootstrapping 语境中：Split/Merge 与 C2S/EvalMod/S2C 的兼容性。

**关键边界：**

| 能做 | 不能做 |
|------|--------|
| 用 RAG 检索背景材料 | 替作者完成正式 proof |
| 整理已有推导、指出缺口 | 宣称 reduced-Q leaf bootstrapping 成立 |
| 生成检查清单、辅助推导 | 自动生成最终 noise theorem |
| 整理 notation、对比方案 | 替作者做最终 security analysis |
| 汇总待确认点 | 把 candidate 当作结论 |

**关键数学结论必须作者确认。**

---

## 7. candidate / defer 保护墙

以下内容**不能**作为 validated conclusion 使用：

| 标记 | 处理规则 |
|------|----------|
| `candidate_knowledge` | 候选知识，不能说「已证明」 |
| `needs_human_review` | 只能说「以下来自待审核条目」 |
| `low` confidence | 只能说「以下仅为初步直觉」 |
| 当前论文 open problem | 只能说材料不足、待确认 |
| reduced-Q | 不能说成立 |
| leaf error aggregation | 不能给 final formula |
| formal noise theorem | 不能给最终陈述 |
| 主定理最终 statement | 不能说已确定 |

**如果被问到以上问题：**

> 「当前知识库没有 validated 结论，只能提供相关背景材料和待确认点。」

---

## 8. 已知特殊情况

| 问题 | 处理 |
|------|------|
| **K063 重号**（两个 K063：C2S/S2C 权威 + CKKS 噪声概览） | 暂不修复。C2S/S2C 查询时 K063 默认指术语卡片。查噪声时用关键词。详见 `09_工作进度/已知特殊情况与处理规则.md` |
| 部分 metadata/path/source_id 为空 | 不影响检索，暂不补 |
| 部分 chunk text 是 JSON dict 字符串 | 不影响语义检索，暂不美化 |
| `10_Batch01修正版/` 是 archive_readonly | 不是日常入口，只用于审计回溯 |
| top-5 检索偶尔有轻微噪声 | 主要看 top-3，可接受 |
| PaCo/GHPS 等部分来源表述 | 可能需要作者或原文再核对 |

**这些不是当前 blocker。除非影响数学结论或检索答案，不要继续微调。**

---

## 9. 如果未来新增资料

参考文件：
- `06_RAG入库准备/RAG_增量更新流程.md`（8 步标准流程）
- `09_工作进度/新资料入库SOP.md`
- `09_工作进度/知识库维护SOP.md`

核心规则：
- 只要 chunks text 内容变化 → 必须重建 embeddings
- 只要 ready chunks 数量变化 → 必须重建 BM25
- 重建后必须跑 `RAG_Smoke_Test_清单.md`

---

## 10. 给下一个 AI 的行为规范

| ✅ 应该做 | ❌ 不应该做 |
|-----------|------------|
| 先读 `总流程入口.md` | 擅自扩展目录 |
| 先读 `00_AI接手上下文/` 理解用户和 Agent | 跳过上下文直接操作 |
| 论文项目从 `11_AI工作区/` 进入 | 在 01–10 论文材料区乱翻 |
| 01–10 = 论文材料区，11_AI工作区 = AI 入口 | 把 AI 文件散落到论文材料区 |
| 用 RAG 检索后基于材料回答 | 生成大量新文件（除非用户要求） |
| 材料不足时明确说不足 | 过度整理 metadata |
| 不确定时先用 `python3 query_rag.py` 检索 | 凭空猜测技术结论 |
| 核心 proof/noise/security 建议用 GPT 做高置信推导检查 | 修改 RAG baseline |
| GPT → 关键数学推导；DeepSeek → 普通整理 | 用 DeepSeek 做严格数学证明 |
| 把 candidate/open problem 标记清楚 | 把 candidate 当正式知识 |
| 维护已知特殊情况记录 | 替作者做数学证明 |
| 遵循术语规则（C2S/S2C） | 继续 K063 重号修复 |
| 遵循只读归档规则 | 动 `10_Batch01修正版/` |

---

**快速入口速查：**

| 你想…… | 看这里 |
|----------|--------|
| 理解用户 & Agent 身份 | `00_AI接手上下文/README.md` → `USER.md` → `IDENTITY.md` → `SOUL.md` |
| 查资料 | `python3 query_rag.py "问题" -k 5 -v` |
| 了解 RAG | `06_RAG入库准备/RAG_Baseline_v1.0_记录.md` |
| 了解论文 | `04_论文协作区/01_当前论文项目/Eurocrypt2027/11_AI工作区/00_接手说明/项目总览.md` |
| 了解目录 | `00_资料清单/资料清单总览.md` |
| 了解已知问题 | `09_工作进度/已知特殊情况与处理规则.md` |
| 增量更新 | `06_RAG入库准备/RAG_增量更新流程.md` |
