# Batch-05b1 处理日志

> 2026-05-04T16:00 完成 | 2026-05-04T16:15 收尾检查

## 范围
- ✅ FHE 方案总览（Gentry09, BGV, BFV, CKKS）
- ✅ BGV vs BFV 本质区别（消息低位高位、ModSwitch vs Rescale）
- ✅ FHE 同态计算完整流程（乘法→重线性化→噪声控制→解密）
- ✅ 密钥切换三场景
- ✅ ModSwitch vs KeySwitch 区别
- ✅ ModSwitch vs Rescaling 区别
- ✅ 与当前论文 CKKS bootstrapping 关系
- ❌ 噪声/安全/参数/实验/GSW 深入

## 新增导入 (8 files)

| source_id | 文件 | 标注 |
|-----------|------|------|
| src_000101 | gentry-BGFV-CKKS.md | 方案总览笔记 |
| src_000102 | BFV与BGV的高低位明文.md | BGV vs BFV |
| src_000103 | BGV-CKKS完整流程举例.md | 计算流程 |
| src_000104 | 为什么以及什么时候需要密钥切换.md | KeySwitch |
| src_000105 | 模数切换与密钥交换区别.md | ModSwitch vs KeySwitch |
| src_000106 | 模数切换和重缩放的区别.md | ModSwitch vs Rescale |
| src_000107 | Gentry 2009 原论文.pdf | 外部论文-参考 |
| src_000108 | CKKS 2017 原论文.pdf | 外部论文-参考 |

## source_id 跳号说明

Batch-05a3 结束时 source_id 到 src_000094。
Batch-05b1 新增从 src_000101 开始，跳过了 src_000095~src_000100 共 6 个编号。

- **原因**：编号脚本在生成 Batch-05b1 的 source_id 时，手动从 101 开始（预留 095-100 为 Batch-05a4 或其他子批次使用）。
- **状态**：src_000095~100 不存在。不是丢失文件，不是导入失败，纯编号跳号。
- **影响**：无。所有 8 个 Batch-05b1 文件已正确导入并记录在 `00_导入记录/原始路径到库内路径映射.jsonl` 中。

## 产出

| 类别 | 数量 |
|------|------|
| 知识卡片 | 7 (K056~K062) |
| 原子问答 | 15 (Q160~Q175，含 Q175 原名 Q161_2) |

## 知识卡片

| card_id | 标题 | confidence | review | card_type |
|---------|------|------------|--------|-----------|
| K056 | FHE方案总览：Gentry09/BGV/BFV/CKKS | medium | ✓ | 用户笔记总结外部论文 |
| K057 | BGV vs BFV本质区别 | medium | ⚠️ | 用户笔记 + BGV/BFV context（见下） |
| K058 | FHE同态计算完整流程 | medium | ✓ | 用户笔记 |
| K059 | 密钥切换三场景 | medium | ✓ | 用户笔记 |
| K060 | ModSwitch vs KeySwitch | medium | ✓ | 用户笔记 |
| K061 | ModSwitch vs Rescaling | medium | ⚠️ | 用户笔记 + 跨方案综合（见下） |
| K062 | 方案层与论文CKKS Bootstrapping关系 | medium | ⚠️ | 跨论文综合 + candidate_knowledge |

## Q&A 编号修复

- Q161_2（非标准编号）→ Q175（统一连续编号）
- 全局索引已更新，0 冲突，0 重复

## 方案比较 caveat

### K057 的风险

K057 标题和内容在 BGV/BFV context 中讨论 ModSwitch 和 Rescale。虽已明确标注 BGV/BFV，但整卡多处出现「重缩放」和「模数切换」术语，读者可能误套到 CKKS。

已在知识卡片末尾追加：
> ⚠️ Caveat：CKKS 的 Rescale（by p）≠ BGV 的 ModSwitch，≠ BFV 的 Rescale（by t）。当前论文在 CKKS context。

已标 needs_human_review = true。

### K061 的风险

K061 跨 BGV/BFV/CKKS 三种 context 对比 ModSwitch 和 Rescale。存在读者将 CKKS Rescale 等同于 BGV ModSwitch 的风险。

已在知识卡片末尾追加：
> ⚠️ Caveat：CKKS 的 Rescale 含精度管理含义，≠ BGV 的 ModSwitch。引用时必须在明确方案 context 下使用。

已标 needs_human_review = true。

### K062

已从创建时标 candidate_knowledge + needs_human_review = true。不判断论文 contribution。

## 质量检查
- [x] Exact FHE (BGV/BFV) vs Approximate FHE (CKKS) 区分
- [x] BGV message in low bits, BFV message in high bits 区分
- [x] ModSwitch (BGV, arbitrary scale) vs Rescaling (BFV/CKKS, fixed scale) 区分
- [x] ModSwitch (change q, same s) vs KeySwitch (change s, same q) 区分
- [x] Relinearization = KeySwitch special case (s^2 -> s)
- [x] 不把 BGV/BFV bootstrapping 与 CKKS bootstrapping 混用
- [x] 当前论文的 C2S/S2C (complex domain) ≠ BGV/BFV 的 C2S/S2C (finite field)
- [x] 术语规则遵守（C2S=coeff->slot, S2C=slot->coeff）
- [x] 收尾检查：K057/K061 补 CKKS Rescale caveat（2026-05-04T16:15）
- [x] 收尾检查：Q161_2 → Q175 修复

## 收尾检查记录（2026-05-04T16:15）

| 检查项 | 结果 |
|--------|------|
| src_000095~100 | 跳过（编号 gap），无文件丢失，已说明 |
| Q161_2 → Q175 | 已修复，0 冲突 |
| K057 CKKS caveat | 已补，needs_human_review=true |
| K061 CKKS caveat | 已补，needs_human_review=true |
| K062 candidate_knowledge | 保持 needs_human_review=true |
| QA 索引 | 已更新，174 条 |
| 卡片索引 | 已更新，62 张 |
