# 完整论文工作区

> paper_id: paper_2027_eurocrypt_fhe  
> 创建日期：2026-05-09  
> 用途：从已有 Section 2-4 草案、Related Work 草案和 5-7 实验表格开始，整理一份可以继续扩写为完整论文的主稿。

## 当前文件

| 文件 | 用途 |
|---|---|
| `论文主稿.md` | 当前整合版主稿启动稿。按完整论文结构组织，已合并理论主线、related work 定位和主实验结论。 |

## 来源材料

| 来源 | 进入主稿的位置 |
|---|---|
| `../Section2-4_整合版草案.md` | Preliminaries、Coefficient-Split Semantics、Core Factorization Theorem、Implementation/Parameter Regime 的理论骨架。 |
| `../Related_Work初稿.md` | Related Work 章节的主要文字来源。 |
| `../../06_Related_Work/Related_Work定位表.md` | Related Work 的定位边界：哪些不能 claim，哪些可以 claim。 |
| `../../04_实验实现/ckks-subring-compression/5-7/成功实验正式表格.tex` | 实验章节中的主表结论、top17 主推参数、top16 边界实验、q-sweep 结论。 |
| `../../会话迁移总结_2026-05-08.md` | 当前项目主线、5-8 构造式低模数 split 的定位、后续待办。 |

## 写作口径

当前主稿采用 **英文正文 + 中文写作备注** 的方式。原因是最终投稿需要英文，但当前仍有若干作者决策点需要保留中文提示。

主推实验口径：

\[
\log N_{\mathrm{top}}=17,\quad
\log N_{\mathrm{leaf}}=16,\quad
\mathrm{layers}=1,\quad
\text{leaf subKey Bootstrapping I1}.
\]

当前单次结果：

```text
logQ ≈ 1425
log(PQ) ≈ 1730
leaf16 128-bit reference ≈ 1747
security margin ≈ +17 bits
Avg L1 precision ≈ 19.40 bits
output levels = 15
total time ≈ 52.03s
speedup over direct ordinary B17 ≈ 2.75x
speedup over only-subKey ≈ 2.60x
speedup over only-split OS17 ≈ 1.21x
```

## 需要继续补的关键项

1. 把 `论文主稿.md` 中的 theorem/proof sketch 替换为 `Section2-4_整合版草案.md` 中的正式证明细节。
2. 给所有 citations 建立统一 bib key 和 `.bib` 文件。
3. 对 L17-I1、OS17、B17、only-subKey 做重复运行，报告均值、方差和机器配置。
4. 补 evaluation key size / memory 表。审稿人很可能会问加速是否换来了过大的 key material。
5. 明确安全估计方法。当前表格中的 `log(PQ)` 上限是粗略 128-bit reference，最终应换成 estimator 表。
6. 决定 5-8 low-q constructive split 是否放主文，还是作为 future work / appendix。

