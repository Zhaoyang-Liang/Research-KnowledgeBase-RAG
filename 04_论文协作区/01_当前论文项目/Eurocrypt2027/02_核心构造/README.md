# 02_核心构造

## 用途
存放论文核心构造（construction/algorithm）：Split、Merge、C2S、EvalMod、S2C、ring split、同态拆分编码、多项式分解等相关材料。

## 已归位文件

### 核心算法与推导
| 文件 | 来源 | 说明 |
|------|------|------|
| 2算法.pdf / 2算法.tex | `idea-plaintext-Decop/重要！！/` | 算法核心构造 |
| 推导.pdf / 推导.tex | `idea-plaintext-Decop/重要！！/` | 核心推导 |
| 压缩.pdf / 压缩.tex | `idea-plaintext-Decop/重要！！/压缩/` | 压缩算法 |
| 2.pdf / 2.tex | `idea-plaintext-Decop/重要！！/压缩/` | 算法变体 |

### 构造说明（Markdown）
| 文件 | 来源 | 说明 |
|------|------|------|
| 1-总览.md | `idea-plaintext-Decop/重要！！/markdown/` | 多项式分解总览 |
| 2-正向.md | `idea-plaintext-Decop/重要！！/markdown/` | 正向分解 |
| 3-反向.md | `idea-plaintext-Decop/重要！！/markdown/` | 反向重构 |
| 4并行性.md | `idea-plaintext-Decop/重要！！/markdown/` | 并行性分析 |
| 3噪声分析.md | `idea-plaintext-Decop/重要！！/markdown/` | 噪声分析 |
| 离散矩阵的旋转加权和分解.md | `idea-plaintext-Decop/资料/重要定理-代数/` | 矩阵分解定理 |

### 拆分与视角
| 文件 | 来源 | 说明 |
|------|------|------|
| 同态下的FFT-based密文拆分.pdf | `AAA/` | FFT 密文拆分 |
| 补充视角-slot情况下的拆分.md | `AAA/` 和 `Filed Switching/` | slot 视角补充 |
| 逐坐标分解证明.pdf | `AAA/` | 逐坐标分解 |
| 逐分量成果.md | `Filed Switching/5月1号前最新成果/` | 最新分解成果 |

## 交叉引用
- `../P0-1_polynomial_split类型核实.md` — 多项式拆分类型核实
- `../P0-2_环参数_m_m0_r_核实.md` — 环参数量化
- `../P0-3_GoodBases条件核实.md` — GoodBases 条件
- `../P0-4_R_vs_Rdual_符号核实.md` — R/R^∨ 符号
- `../当前理解与主线.md` — 论文主线理解

## 外部核心参考（不在此目录）
- `Research-KB-RAG/01_源文件库/04_FHE核心技术模块/02_C2S_S2C/` — 13 个 C2S/S2C 详细推导文件
- `Research-KB-RAG/01_源文件库/04_FHE核心技术模块/03_SIMD_SlotPermutation/` — SIMD 与槽置换
- `Research-KB-RAG/01_源文件库/04_FHE核心技术模块/04_NTT_DFT_FFT/` — NTT/DFT/FFT
- `Research-KB-RAG/01_源文件库/04_FHE核心技术模块/trace操作.md` — Trace 操作

## 注意事项
- **算法状态**：多项式拆分/分解的算法仍为 candidate，未最终验证
- **3噪声分析.md** 涉及噪声部分同时属于 `05_参数与安全/`
- 核心构造依赖于 `02_C2S_S2C/` 参考材料中的 M→U→M^{-1} 管线
