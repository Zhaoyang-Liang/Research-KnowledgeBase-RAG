# 04_实验实现

## 用途
存放实验日志、实现证据、benchmark、时间记录、正确性验证、实现相关材料。

## 已归位文件

| 文件 | 来源 | 说明 |
|------|------|------|
| time-logn16.txt | `重要实验结果/4-29/` | logN=16 实验时间记录 |
| 简单正确性验证.txt | `重要实验结果/4-29/` | 正确性验证数据 |
| 简单总结表格.md | `重要实验结果/4-29/` | 实验总结表格 |
| 复杂度分析初步.md | `安全性_复杂度_噪声分析/复杂度分析初步/` | 复杂度分析（含 logN=13/16/17 实验数据） |
| 实验与实现证据表.md | 项目根目录 → 本目录 | 实验与实现证据总表 |

## 交叉引用
- `../02_核心构造/` — 算法构造（实验验证对象）
- `../05_参数与安全/` — 噪声、安全参数

## 当前状态

### 已有实验
- logN=16 时间记录
- 正确性验证（简单正确性验证.txt）
- logN=13/16/17 复杂度分析数据（在复杂度分析初步.md 中）

### 待补充
- Benchmark 标准化（统一对比格式）
- 更多 logN 级别的实验
- 与 baseline 方案的对比实验

## ckks-subring-compression 实验仓库

**原始仓库**：`/Users/mac/ckks-subring-compression`

### 已整理入论文项目
1. **`ckks-subring-compression/4-29/`** — 主实验（main.go + N13/N16/N17 输出 + trace-version + 结果分析）
2. **`ckks-subring-compression/4-29新开辟RNS实验/`** — RNS 探针实验（4 个 probe + 3 篇分析）

### 未复制的目录（archive / dependency / old）
- `lattigo-v6-local/` — 本地依赖/fork，太大，需时追溯原仓库
- `old/` — 旧实验/历史代码
- `4-28-noBK/` — 旧探针/旁支实验
- `tex笔记/` — 旁支笔记
- root `main.go` — 待确认是否仍是当前入口

### 实验索引
👉 详见 `ckks-subring-compression/实验索引.md`

> ⚠️ 不运行实验，不修改代码。实验结果目前只能作为 implementation evidence / observation，不直接作为论文最终 claim。
> ⚠️ `lattigo-v6-local/` 已从项目副本中清理（误复制移除）。需要时追溯原始仓库 `/Users/mac/ckks-subring-compression/`。`go.mod` / `go.sum` 保留以记录实验环境。

## 注意事项
- 实验数据为临时记录，不是最终论文表格
- 复杂度分析初步.md 中表格数据需要确认是否来自正式 benchmark
- ckks-subring-compression 实验结论需作者确认后方可进入论文正文
