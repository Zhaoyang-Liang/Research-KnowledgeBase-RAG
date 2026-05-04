# Batch-05c 处理日志 (收尾版)

> 2026-05-04T17:30 | RLWE / 格密码 / CVP / Babai / 重线性化 | 收尾完成

## 审核规则变更

采用用户指定的三层审核规则：
- **不强制 human_review**：背景知识卡片，有完整源文件，不直接影响当前论文 → `review_policy: review_when_used_for_paper`
- **强制 human_review**：影响当前论文 construction/proof/noise/security/contribution，或跨源推断，或源文件缺失
- **Falcon/NTRU/签名**：`document_role: side_material`，暂不展开

## Batch-05c 卡片审核状态

| card_id | needs_review | policy | 理由 |
|---------|-------------|--------|------|
| K065 | false | review_when_used | RLWE背景知识，源文件完整 |
| K066 | false | review_when_used | CVP/Babai几何背景，源文件完整 |
| K067 | false | review_when_used | LWE→SVP归约背景，源文件完整 |
| K068 | false | review_when_used | 重线性化标准机制，源文件完整 |
| K069 | false | review_when_used | Ideal Lattice背景，源文件完整 |

## Q185 降级

Q185 (Fourier bias) 从 `needs_human_review: true, confidence: low` 降级为 `false, medium`。
理由：LWE_SVP_reduc.tex 已完整导入，Q185 是后量子背景知识，不直接影响当前论文证明。

## Q186 归属修正

Q186 (平行六面体泄漏) 主归属：02_RLWE_格密码_后量子知识库（签名安全话题）。
不在 FHE 知识库中独立发布，仅通过 cross_reference 关联 GGH 历史。

## Falcon / NTRU / 签名材料

| 材料 | 标记 |
|------|------|
| src_000116 (Falcon.tex) | side_material, future: Post-Quantum Signatures |
| src_000118 (NTRU杂项) | support_note, review_when_used |
| src_000114 (Klein-GPV) | side_material, future: Post-Quantum Signatures |

全部 `_METADATA.json` 已写入对应目录。

## 产出汇总

| 类别 | 本批新增 | 累计 |
|------|---------|------|
| source_id | 11 (src_000109-119) | 119 |
| 知识卡片 | 5 (K065-K069) | 69 |
| 原子问答 | 12 (Q176-Q187) | 187 |
| needs_human_review 卡片 | 0 新增 | 保持不变 |
| side_material 标记 | 3 | 3 |

## 质量检查
- [x] 新审核规则已应用
- [x] K065-K069 全部 review_when_used_for_paper
- [x] Q185 已降级
- [x] Q186 归属已修正
- [x] Falcon/NTRU side_material 已标记
- [x] 无重复问答实体
- [x] 未展开签名专题
- [x] 未进入 Batch-06
- [x] 未做 embedding
- [x] 所有索引已更新

## Batch-05c 完全收尾 ✅
