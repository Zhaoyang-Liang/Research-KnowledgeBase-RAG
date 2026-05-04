#!/usr/bin/env python3
"""
build_bm25_index.py

读取 cards_ready_chunks.jsonl + qa_ready_chunks.jsonl，
构建 BM25 索引 + card_id/qa_id/source_id 精确查找表，
保存到 data/bm25_index.pkl。

用法:
  python3 build_bm25_index.py

输出:
  data/bm25_index.pkl — pickle 文件，包含:
    - bm25: BM25Okapi 实例
    - chunks: 原始 chunk 列表
    - card_id_index: card_id → [chunk_idx, ...]
    - qa_id_index:   qa_id   → [chunk_idx, ...]
    - source_id_index: source_id → [chunk_idx, ...]
"""

import json
import os
import pickle
import sys


def tokenize(text: str) -> list:
    """
    中文单字切分 + 英文缩写保持 + 数学符号保留。
    
    CKKS       → ["ckks"]
    slot-to-coeff → ["slot-to-coeff"]
    典范嵌入    → ["典","范","嵌","入"]
    σ⁻¹        → ["σ⁻¹"]
    K004       → ["k004"]
    """
    tokens = []
    buf = ""
    
    for ch in str(text):
        if '\u4e00' <= ch <= '\u9fff':  # 中文字符
            if buf:
                tokens.append(buf.lower())
                buf = ""
            tokens.append(ch)
        elif ch.isalnum() or ch in '-_./#':  # 英文/数字/允许的符号
            buf += ch
        elif ch in '⁰¹²³⁴⁵⁶⁷⁸⁹⁻⁺':  # Unicode 上下标
            buf += ch
        else:
            if buf:
                tokens.append(buf.lower())
                buf = ""
            if not ch.isspace():
                tokens.append(ch)
    
    if buf:
        tokens.append(buf.lower())
    
    return [t for t in tokens if t.strip()]


def main():
    KB = os.path.expanduser("~/Desktop/Research-KB-RAG")
    chunk_dir = f"{KB}/06_RAG入库准备/chunks"
    out_dir = os.path.dirname(os.path.abspath(__file__)) + "/data"
    os.makedirs(out_dir, exist_ok=True)
    out_path = f"{out_dir}/bm25_index.pkl"

    # ── 加载 chunks ────────────────────────────────────
    chunks = []
    for fname in ["cards_ready_chunks.jsonl", "qa_ready_chunks.jsonl"]:
        fp = f"{chunk_dir}/{fname}"
        if not os.path.exists(fp):
            print(f"WARNING: {fp} not found, skip")
            continue
        with open(fp) as f:
            for line in f:
                line = line.strip()
                if line:
                    c = json.loads(line)
                    if c.get("rag_status") == "ready":
                        chunks.append(c)

    print(f"✓ 加载 {len(chunks)} 个 ready chunks")

    # ── 构建 BM25 ──────────────────────────────────────
    # rank_bm25 需要先安装: pip install rank_bm25
    try:
        from rank_bm25 import BM25Okapi
    except ImportError:
        print("ERROR: 请先安装 rank_bm25: pip install rank_bm25")
        sys.exit(1)
    
    texts = [c["text"] for c in chunks]
    tokenized = [tokenize(t) for t in texts]
    bm25 = BM25Okapi(tokenized)

    # ── 构建精确 ID 查找表 ──────────────────────────────
    card_id_index = {}
    qa_id_index = {}
    source_id_index = {}

    for i, c in enumerate(chunks):
        cid = c.get("card_id", "")
        if cid:
            card_id_index.setdefault(cid, []).append(i)
        qid = c.get("qa_id", "")
        if qid:
            qa_id_index.setdefault(qid, []).append(i)
        sid = c.get("source_id", "")
        if sid:
            source_id_index.setdefault(sid, []).append(i)

    # ── 保存 ───────────────────────────────────────────
    index = {
        "bm25": bm25,
        "chunks": chunks,
        "card_id_index": card_id_index,
        "qa_id_index": qa_id_index,
        "source_id_index": source_id_index,
    }

    with open(out_path, "wb") as f:
        pickle.dump(index, f)

    print()
    print("=" * 50)
    print(f"✓ BM25 索引构建完成")
    print(f"  chunks:         {len(chunks)}")
    print(f"  card_id 条目:   {len(card_id_index)}")
    print(f"  qa_id 条目:     {len(qa_id_index)}")
    print(f"  source_id 条目: {len(source_id_index)}")
    print(f"  输出文件:       {out_path}")
    print("=" * 50)


if __name__ == "__main__":
    main()
