# RAG Smoke Test 清单

> RAG Baseline v1.0 (2026-05-04 frozen)  
> 每次重建 embeddings + BM25 后，必须运行以下测试

---

## 必测项目（5 条）

```bash
cd ~/Desktop/Research-KB-RAG/06_RAG入库准备/api_rag
```

### Test 1：C2S/S2C 方向正确性

```bash
python3 query_rag.py "C2S 和 S2C 的方向是什么？" -k 5 -v
```

验收标准：
- ✅ top-3 方向正确：`C2S = CoeffToSlot = coefficient → slot`
- ✅ top-3 方向正确：`S2C = SlotToCoeff = slot → coefficient`
- ✅ top-3 无完整旧错误表达
- ✅ at least one result from K063

### Test 2：ModSwitch/Rescale/KeySwitch 区分

```bash
python3 query_rag.py "ModSwitch、Rescale、KeySwitch 三者有什么区别？" -k 5 -v
```

验收标准：
- ✅ top-3 能区分三个概念的不同用途
- ✅ K060/K061 相关内容被返回
- ✅ 明确指出 ModSwitch 降噪不变量、Rescale 保持缩放因子、KeySwitch 换密钥

### Test 3：CKKS 编码管线

```bash
python3 query_rag.py "CKKS 编码管线是什么？" -k 5 -v
```

验收标准：
- ✅ top-3 包含编码流程（z → σ^{-1} → m → Δ 缩放 → 加密）
- ✅ 典范嵌入和 CKKS 编码的区别被正确表述

### Test 4：Exact ID 命中

```bash
python3 query_rag.py "K063" -k 5 -v
```

验收标准：
- ✅ top-1 为 K063 的 card_summary 或 card_knowledge_point chunk
- ✅ 文本内容正确（C2S/S2C 术语定义）

### Test 5：Open problem 不返回确定结论

```bash
python3 query_rag.py "reduced-Q leaf bootstrapping 是否成立？" -k 5 -v
```

验收标准：
- ✅ top-5 不包含「reduced-Q 一定成立」之类的 validated conclusion
- ✅ 如果返回相关内容，应来自 defer/candidate 或中性背景材料
- ✅ 不应出现幻觉式的确定答案

---

## 可选补充测试

```bash
# CKKS vs BGV/BFV slot 区别
python3 query_rag.py "CKKS 的 slot 和 BGV 的 CRT slot 有什么不同？" -k 5 -v

# KeySwitch 通用公式
python3 query_rag.py "KeySwitching 的通用公式是什么？" -k 5 -v

# R vs R^∨ 符号
python3 query_rag.py "R 和 R双 ∨ 在 GHPS 和当前论文中有什么区别？" -k 5 -v

# Exact source_id
python3 query_rag.py "src_000059" -k 5 -v

# Exact card_id
python3 query_rag.py "K004" -k 5 -v
```

---

## 验收总结

| 项目 | 标准 |
|------|------|
| top-3 相关性 | 应基本相关，允许少量语义映射偏差 |
| top-5 噪声 | 允许轻微噪声，但不应出现完全不相关的结果 |
| forbidden patterns | 零容忍（C2S/S2C 方向类查询） |
| exact ID | top-1 应命中 |
| open problem | 不返回 validated conclusion |
| 性能 | 查询应在几秒内返回 |

---

## 快速一命令全测

```bash
cd ~/Desktop/Research-KB-RAG/06_RAG入库准备/api_rag && \
echo "=== Test 1: C2S/S2C ===" && python3 query_rag.py "C2S 和 S2C 的方向是什么？" -k 3 && \
echo "=== Test 2: ModSwitch/Rescale/KeySwitch ===" && python3 query_rag.py "ModSwitch、Rescale、KeySwitch 三者有什么区别？" -k 3 && \
echo "=== Test 3: CKKS encoding ===" && python3 query_rag.py "CKKS 编码管线是什么？" -k 3 && \
echo "=== Test 4: K063 exact ID ===" && python3 query_rag.py "K063" -k 3 && \
echo "=== Test 5: reduced-Q open problem ===" && python3 query_rag.py "reduced-Q leaf bootstrapping 是否成立？" -k 3 && \
echo "=== ALL DONE ==="
```
