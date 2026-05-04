#!/usr/bin/env python3
"""
build_embeddings_openai.py

读取 cards_ready_chunks.jsonl + qa_ready_chunks.jsonl 中全部 ready chunks，
调用 OpenAI text-embedding-3-small API 生成向量，
本地保存到 data/embeddings.jsonl。

用法:
  source config.env
  python3 build_embeddings_openai.py

输出:
  data/embeddings.jsonl  — 每行 {chunk_id, embedding: [float,...]}

成本估算 (text-embedding-3-small, $0.02/1M tokens):
  464 chunks × ~200 chars × ~0.25 tok/char ≈ 23K tokens ≈ $0.0005
"""

import json
import os
import sys
import time


def main():
    # ── 配置 ──────────────────────────────────────────
    api_key = os.environ.get("OPENAI_API_KEY")
    if not api_key:
        print("ERROR: 请设置环境变量 OPENAI_API_KEY")
        print("  source config.env  # 或在 shell 中 export OPENAI_API_KEY=...")
        sys.exit(1)

    model = os.environ.get("EMBEDDING_MODEL", "text-embedding-3-small")

    KB = os.path.expanduser("~/Desktop/Research-KB-RAG")
    chunk_dir = f"{KB}/06_RAG入库准备/chunks"
    out_dir = os.path.dirname(os.path.abspath(__file__)) + "/data"
    os.makedirs(out_dir, exist_ok=True)
    out_path = f"{out_dir}/embeddings.jsonl"

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
                    chunks.append(json.loads(line))

    print(f"✓ 加载 {len(chunks)} 个 chunks")

    # ── 安全检查：只处理 rag_status=ready ──────────────
    ready_chunks = [c for c in chunks if c.get("rag_status") == "ready"]
    rejected = len(chunks) - len(ready_chunks)
    if rejected > 0:
        print(f"⚠ 已排除 {rejected} 个非 ready chunk")
    chunks = ready_chunks

    if not chunks:
        print("ERROR: 没有 ready chunk，退出。")
        sys.exit(1)

    # ── 调用 OpenAI Embedding API ──────────────────────
    from openai import OpenAI
    client = OpenAI(api_key=api_key)

    batch_size = 100  # 每批最多 2048 个 input
    embeddings = []
    total_tokens = 0

    for i in range(0, len(chunks), batch_size):
        batch = chunks[i : i + batch_size]
        texts = [c["text"][:8000] for c in batch]  # 截断超长文本

        response = client.embeddings.create(
            input=texts,
            model=model,
        )

        for j, emb_data in enumerate(response.data):
            chunk = batch[j]
            embeddings.append({
                "chunk_id": chunk["chunk_id"],
                "embedding": emb_data.embedding,
            })

        batch_tokens = response.usage.total_tokens
        total_tokens += batch_tokens
        progress = min(i + batch_size, len(chunks))
        print(f"  batch {i // batch_size + 1}: "
              f"{progress}/{len(chunks)} chunks, "
              f"{batch_tokens} tokens")

        # 速率控制
        if i + batch_size < len(chunks):
            time.sleep(0.3)

    # ── 保存 ───────────────────────────────────────────
    with open(out_path, "w") as f:
        for emb in embeddings:
            f.write(json.dumps(emb) + "\n")

    # ── 报告 ───────────────────────────────────────────
    cost_per_m = 0.02  # text-embedding-3-small
    est_cost = total_tokens / 1_000_000 * cost_per_m
    dim = len(embeddings[0]["embedding"]) if embeddings else 0

    print()
    print("=" * 50)
    print(f"✓ 完成")
    print(f"  chunks:       {len(embeddings)}")
    print(f"  向量维度:      {dim}")
    print(f"  总 tokens:     {total_tokens:,}")
    print(f"  估算成本:      ${est_cost:.6f}")
    print(f"  输出文件:      {out_path}")
    print("=" * 50)


if __name__ == "__main__":
    main()
