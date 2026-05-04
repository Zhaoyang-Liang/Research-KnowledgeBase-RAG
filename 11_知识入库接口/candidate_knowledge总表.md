# candidate_knowledge 总表

> 2026-05-04T18:10 | 知识库工程完成阶段（补充未来论文支持）

## 汇总 (10 条)

| ID | paper_id | 来源 | 类型 | needs_review | 能否入库 | 建议动作 |
|----|----------|------|------|-------------|----------|----------|
| C-NOISE-001 | Eurocrypt2027 | 噪声分析 | candidate | ✓ | ⏸ | 尚无 formal bound |
| C-NOISE-002 | Eurocrypt2027 | 噪声分析 | candidate | ✓ | ⏸ | 等待论文 noise analysis |
| C-NOISE-003 | Eurocrypt2027 | 噪声分析 | candidate | ✓ | ⏸ | 等待论文 noise analysis |
| C-NOISE-004 | Eurocrypt2027 | 噪声分析 | candidate | ✓ | ⏸ | 等待作者确认是否进入 theorem |
| C-NOISE-005 | Eurocrypt2027 | 噪声分析 | open_problem | ✓ | ⏸ | 等待 reduced-Q 决策 |
| C-NOISE-006 | Eurocrypt2027 | 噪声分析 | checklist | ✓ | ⏸ | 等待 proof completion |
| C-NOISE-007 | Eurocrypt2027 | 噪声分析 | candidate | ✓ | ⏸ | 需写入论文 preliminaries |
| C-NOISE-008 | Eurocrypt2027 | 噪声分析 | candidate | ✓ | ⏸ | 论文写作时确认 |
| C-NOISE-009 | Eurocrypt2027 | 噪声分析 | notation_bridge | ⚠️ | ⏸ | 需作者最终决定 |
| C-NOISE-010 | Eurocrypt2027 | 噪声分析 | candidate | ✓ | ⏸ | 论文写作时写入 |

## 按状态分类

### 需作者确认（2 条）

| ID | 内容 | 阻塞 |
|----|------|------|
| C-NOISE-005 | Reduced-Q Open Problem | 是 |
| C-NOISE-009 | R vs R^∨ Notation Bridge | 是 |

### 等待论文写作（8 条）

| ID | 内容 | 阶段 |
|----|------|------|
| C-NOISE-001~004 | Ring switching / Split / Leaf / γ_impl | proof/construction |
| C-NOISE-006~008 | Checklist / Slot区分 / 术语统一 | construction |
| C-NOISE-010 | 多项式拆分类型确认 | preliminaries |

## 入库决策

| 决策 | 数量 | 说明 |
|------|------|------|
| 论文写完后可入库 | 8 | 一旦论文中确认，可转为稳定知识 |
| 需作者决定是否入库 | 2 | C-NOISE-005 和 C-NOISE-009 |

## 未来论文项目接入

### 新论文 candidate 注册规则

1. 新论文项目的每个 idea 自动获得 `idea_id`（格式：`I-{paper_id}-{slug}`）
2. idea 先入论文项目内部的 `知识沉淀候选.md`
3. 同步登记到本总表：
   - candidate_id 格式：`C-{paper_id}-{序号}`，如 `C-ICITS2028-001`
   - paper_id 明确标注
   - status 初始 = `candidate_knowledge`
   - can_be_promoted_to_kb = false
   - rag_status = defer
4. 不直接进入正式知识库
5. 不直接进入 RAG ready

### 多论文项目共存

```
Eurocrypt2027（活跃）
  → C-NOISE-001 ~ C-NOISE-010

未来论文 A（如 ICITS 2028）
  → C-ICITS2028-001 ~ ...
  → 目录：04_论文协作区/ICITS2028/

未来论文 B（如 Crypto 2029）
  → C-Crypto2029-001 ~ ...
  → 目录：04_论文协作区/Crypto2029/
```

### 跨论文知识复用

| 来源论文 | 知识状态 | 下一论文可用性 |
|----------|---------|---------------|
| Eurocrypt2027 | validated 卡片 | ✅ 可直接引用（source 可追溯） |
| Eurocrypt2027 | candidate | ❌ 需独立验证 |
| Eurocrypt2027 | negative_result | ✅ 可借鉴（避免重复弯路） |
| Eurocrypt2027 | deprecated | ⚠️ 仅历史参考 |

---

⚠️ 所有 candidate 当前均未 validated，未进入正式知识库。
未来论文的 candidate 也遵循同样规则：先从 `candidate` 开始，经作者确认后再进入长期知识库。
