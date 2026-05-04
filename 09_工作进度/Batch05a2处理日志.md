# Batch-05a2 处理日志

> 2026-05-04T15:02 | SlotToCoeff / CoeffToSlot 算法整理

## 范围

- ✅ SlotToCoeff / CoeffToSlot 通用定义（BGV/BFV context）
- ✅ C2S/S2C线性变换T结构
- ✅ M→U_ℓ→M^-1三步分解管线
- ✅ U_ℓ块Fourier/Vandermonde结构与FFT-like分解
- ✅ CKKS bootstrapping中C2S/S2C的位置
- ✅ 当前论文中C2S/S2C与Split/Merge关系
- ❌ SMID泛化（→Batch-05a3）
- ❌ slot permutation（→Batch-05a3）
- ❌ NTT泛化（→Batch-05a3）

## 新增导入 (13 files)

| source_id | 文件 |
|-----------|------|
| src_000068 | 0定义.md |
| src_000069 | 1详细.md |
| src_000070 | 2为什么可以看做线性变换T.md |
| src_000071 | 3真实场景.md |
| src_000072 | 4-1局部槽内基与矩阵M.md |
| src_000073 | 4-2M矩阵的实现.md |
| src_000074 | 4-3M与forbeinus算子表示.md |
| src_000075 | 5-1k(v).md |
| src_000076 | 5-2综合例子.md |
| src_000077 | 6-1矩阵U具体.md |
| src_000078 | 6-2U矩阵展开具体.md |
| src_000079 | 6矩阵U.md |
| src_000080 | 7-再读论文之MUM.md |

## 产出

| 类别 | 数量 |
|------|------|
| 知识卡片 | 6 (K044~K049) |
| 原子问答 | 15 (Q131~Q145) |

## 新增知识卡片与交叉引用

| card_id | 标题 | 交叉引用 |
|---------|------|---------|
| K044 | SlotToCoeff/CoeffToSlot通用定义 | K019(CKKS C2S), K002(CRT) |
| K045 | C2S/S2C作为线性变换T | K044(定义), K046(管线) |
| K046 | M→U_ℓ→M^-1三步分解 | K047(M), K048(U_ℓ) |
| K047 | M矩阵：槽内基统一化 | K046(管线), K048(U_ℓ) |
| K048 | U_ℓ块Fourier/FFT结构 | K046(管线), K018(CKKS U_n), K019(CKKS C2S) |
| K049 | C2S/S2C在bootstrapping+论文位置 ⚠️ | K039(CKKS编码), K044(C2S定义) |

## 需要核原文

**无**。本批13个文件全部为用户的BGV/BFV slotToCoeff笔记（非翻译笔记）。
⚠️ K049标记needs_human_review：关于C2S/S2C在当前论文中位置的综合推导，需作者确认。

## 质量检查

- [x] 不把C2S/S2C和CKKS encoding混用
- [x] 不把M/U_ℓ(M^-1 pipeline和PaCo的U_n CT分解混用
- [x] 不把translation_note当原文（本批无翻译文件）
- [x] 与K018/K019/K030/K036/K039/K040/K042做去重
- [x] 当前论文相关内容标记needs_human_review


---

## C2S/S2C 术语方向审计后修正 (2026-05-04T15:22)

### K019/K020 方向冲突修正

- **K019**: 原标题"SlotToCoeff正向"，实际公式slot=U_n·coeff→归为C2S/CoeffToSlot。已重命名、降confidence、标needs_human_review
- **K020**: 原标题"CoeffToSlot反向"，实际公式coeff=U_n^-1·slot→归为S2C/SlotToCoeff。已重命名、降confidence、标needs_human_review

### 全库统一术语: C2S=coeff→slot, S2C=slot→coeff

详见 `C2S_S2C术语方向审计.md`
