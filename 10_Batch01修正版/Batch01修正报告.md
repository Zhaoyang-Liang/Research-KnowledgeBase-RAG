# Batch-01 修正版完成报告

> 生成时间: 2026-05-04 02:30 CST | 版本: v2

## 一、完成统计

| 指标 | v1 (旧版) | v2 (修正版) | 变化 |
|------|-----------|-------------|------|
| 知识卡片 (JSON) | 10 | 28 | +18 (含1张新增专项卡片) |
| 知识卡片 (MD) | 0 | 28 | +28 (新增格式) |
| 原子问答 | 16 | 100 | +84 |
| 来源映射 | 10 | 27+1 | +18 |
| 术语表 | 0 | 1 (35术语) | 新增 |
| 公式表 | 0 | 1 (38公式) | 新增 |
| 人工确认清单 | 0 | 1 (12项) | 新增 |

## 二、卡片清单 (28张)

### 代数基础 (K001-K011)
K001 分式域与分裂域 | K002 共轭-多项式角度 | K003 复嵌入两步分解 | K004 Trace落入基域 | K005 自同构个数 | K006 Prime Splitting | K007 代数扩张与代值同态 | K008 SIMD代数总结 | K009 分圆多项式F_p分解(修正) | K010 分圆多项式十大优势(修正) | K011 NTT-Cooley-Tukey(修正)

### FHE总论 (K012-K017)
K012 四方案对比(增强) | K013 BGV/CKKS完整流程(修正) | K014 密钥切换通用公式(修正) | K015 ModSwitch vs Rescale(增强) | K016 噪声控制层次(修正) | K017 BGV vs BFV高低位(增强)

### FHE核心技术 (K018-K026)
K018 典范嵌入与U矩阵 | K019 SlotToCoeff正向 | K020 CoeffToSlot反向 | K021 Trace操作 | K022 多项式幂次代换 | K023 多项式拆分对比(需确认) | K024 PaCo全文翻译 | K025 多项式分解总览(需确认) | K026 正向分解实现(需确认)

### 格密码 (K027) + 跨源合成 (K028)
K027 Babai最近平面算法(修正) | K028 CKKS slot vs BGV/BFV CRT slot(新增,需确认)

## 三、重点修正主题状态

| 主题 | 状态 |
|------|------|
| BGV ModSwitch / BFV Rescale / CKKS Rescale 三方案区分 | 已修正 |
| CKKS slot vs BGV/BFV CRT slot 数学来源差异 | 新增K028,需确认 |
| PaCo中 p vs Delta 区别 | 已标注需确认 |
| Trace 与 field switching 关系 | 已修正 |
| Frobenius轨道与CRT slot关系 | 已修正 |
| 自同构与rotation/slot permutation关系 | 已修正 |
| Coeff split vs Slot-preserving split | 已标注需确认 |
| SlotToCoeff/CoeffToSlot与NTT/DFT关系 | 已修正 |

## 四、最高优先级(P0)需确认问题

1. 论文 construction 使用的多项式拆分类型 (K023/K025/K026)
2. CKKS slot vs CRT slot 术语是否混用 (K028)
3. BFV/CKKS Rescale 术语在论文中的定义 (K012/K015)
4. PaCo 中 p 和 Delta 的精确值 (K024)

## 五、是否建议先精读 GHPS Duality / Good Bases？

建议: 是。K004/K006/K021 都依赖 i=1(mod m0) 条件。
Field switching 的核心机制需要理解 Duality。
Good Bases 是 field switching 实现正确性的关键前提。
建议在开始 Batch-02 之前先导入 GHPS 论文并精读。

## 六、是否可开始准备 embedding 输入样例？

可以，但建议 P0 确认项完成1-2项后再做。
当前: 卡片格式规范、source引用完整、路径干净。
但有 4-6 张卡片标注 needs_human_review，GHPS/PaCo source_pending_import。

## 七、产出文件

```
10_Batch01修正版/
 修正版知识卡片_JSON/     28 个 .json
 修正版知识卡片_MD/       28 个 .md
 修正版原子问答.jsonl     100 条
 修正版来源映射.jsonl     28 条
 修正版术语表.md          35 术语
 修正版公式表.md          38 公式
 需要人工确认清单.md       12 项(P0/P1/P2)
 Batch01修正报告.md        本文件
```


## 九、文件名增强 (2026-05-04)

56 个文件已重命名为 `KXXX__语义标题.扩展名` 格式：
- JSON 卡片：28 个重命名
- Markdown 卡片：28 个重命名
- 所有 `card_id` 不变
- 新增 `卡片索引.md` 和 `card_manifest.jsonl`
