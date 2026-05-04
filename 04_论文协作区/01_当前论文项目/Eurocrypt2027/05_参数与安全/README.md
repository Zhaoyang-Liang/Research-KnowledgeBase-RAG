# 05_参数与安全

## 用途
存放噪声分析、安全性分析、参数选择、reduced-Q leaf bootstrapping、γ_impl、ring switching noise、KeySwitch noise 等材料。

## 已归位文件

### 噪声分析
| 文件 | 来源 | 说明 |
|------|------|------|
| 噪声与成功概率分析.md | `安全性_复杂度_噪声分析/噪声与成功概率/` | 噪声与成功概率分析 |
| 3噪声分析.pdf | `idea-plaintext-Decop/重要！！/` | 算法相关的噪声分析 |
| 3噪声分析.tex | `idea-plaintext-Decop/重要！！/` | 噪声分析 TeX |
| 3噪声分析.md | `idea-plaintext-Decop/重要！！/markdown/` | 噪声分析 Markdown |
| 噪声分析作者决策页.md | 项目根目录 → 本目录 | **关键文件**：作者噪声决策 |
| 噪声与正确性缺口清单.md | 项目根目录 → 本目录 | 噪声与正确性缺口 |

### 归一化
| 文件 | 来源 | 说明 |
|------|------|------|
| 归一化原因.pdf | `安全性_复杂度_噪声分析/` | 归一化原因推导 |
| 归一化原因.tex | `安全性_复杂度_噪声分析/` | 归一化 TeX |
| 归一化原因进一步解释.md | `安全性_复杂度_噪声分析/` | 归一化详细解释 |

### 安全参数
| 文件 | 来源 | 说明 |
|------|------|------|
| ckks_parameter_tutorial.pdf | `安全性_复杂度_噪声分析/安全性分析/` | CKKS 参数教程 |
| ckks_parameter_tutorial.tex | `安全性_复杂度_噪声分析/安全性分析/` | 参数教程 TeX |
| 安全参数分析.md | `idea-plaintext-Decop/paper组内讨论初步/待完成/` | 安全参数分析 |
| 都要进行哪些安全性分析.md | `Filed Switching/5月1号前最新成果/` | 安全性分析规划 |

## ⚠️ 重要：未完成项

| 项目 | 状态 |
|------|------|
| reduced-Q leaf bootstrapping | 🔴 open / candidate |
| leaf error aggregation | 🔴 未完成 |
| formal noise theorem | 🔴 未完成 |
| 主定理噪声界 | 🟡 待作者决策 (A1/A2/A3) |

## 交叉引用
- `../../Research-KB-RAG/09_工作进度/Batch05b2_pre_噪声分析双层盘点.md` — 噪声双层盘点
- `../../Research-KB-RAG/11_知识入库接口/噪声分析知识沉淀候选.md` — 噪声知识沉淀

## 注意事项
- **本目录中的噪声/安全分析不能当作 validated conclusion 使用**
- 噪声分析作者决策页中的 A1/A2/A3 未决策前，噪声分析暂停
- 归一化（DIV-K vs RAW）相关结论需作者确认
