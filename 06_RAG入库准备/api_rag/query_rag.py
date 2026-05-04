#!/usr/bin/env python3
"""
query_rag.py

SiliconFlow 直连版 Hybrid RAG 检索脚本。

检索方式：
- SiliconFlow embedding 语义检索
- BM25 关键词检索
- card_id / qa_id / source_id 精确匹配
- RRF 融合
- 只返回 rag_status=ready 的 chunks

用法：
  python3 query_rag.py "C2S 和 S2C 的区别"
  python3 query_rag.py "CKKS 编码管线是什么" -v
  python3 query_rag.py "K004" --ids-only
  python3 query_rag.py "LWE 和 RLWE" --json
"""

import argparse
import json
import os
import pickle
import sys
import urllib.request
import urllib.error
from pathlib import Path

import numpy as np


API_URL = "https://api.siliconflow.cn/v1/embeddings"


def read_config_file(config_path: Path) -> dict:
    values = {}

    if not config_path.exists():
        return values

    for raw_line in config_path.read_text(encoding="utf-8").splitlines():
        line = raw_line.strip()

        if not line or line.startswith("#"):
            continue

        if line.startswith("export "):
            line = line[len("export "):].strip()

        if "=" not in line:
            continue

        key, value = line.split("=", 1)
        key = key.strip()
        value = value.strip()

        if " #" in value:
            value = value.split(" #", 1)[0].strip()

        value = value.strip().strip('"').strip("'")

        if key:
            values[key] = value

    return values


def get_config() -> tuple[str, str]:
    script_dir = Path(__file__).resolve().parent
    cfg_file = script_dir / "config.env"
    file_cfg = read_config_file(cfg_file)

    provider = os.environ.get("EMBEDDING_PROVIDER") or file_cfg.get("EMBEDDING_PROVIDER") or "siliconflow"
    model = os.environ.get("EMBEDDING_MODEL") or file_cfg.get("EMBEDDING_MODEL") or "Pro/BAAI/bge-m3"
    api_key = os.environ.get("EMBEDDING_API_KEY") or file_cfg.get("EMBEDDING_API_KEY") or ""

    provider = provider.strip()
    model = model.strip()
    api_key = api_key.strip().strip('"').strip("'")

    if provider != "siliconflow":
        print(f"ERROR: 这个脚本只支持 siliconflow，当前 EMBEDDING_PROVIDER={provider}")
        sys.exit(1)

    if not api_key:
        print("ERROR: 没有读到 EMBEDDING_API_KEY。请检查 config.env。")
        sys.exit(1)

    return model, api_key


def embed_text(text: str, model: str, api_key: str) -> np.ndarray:
    payload = {
        "model": model,
        "input": text,
        "encoding_format": "float",
    }

    data = json.dumps(payload, ensure_ascii=False).encode("utf-8")

    headers = {
        "Authorization": f"Bearer {api_key}",
        "Content-Type": "application/json",
        "Accept": "application/json",
        "User-Agent": "Research-KB-RAG/1.0",
    }

    req = urllib.request.Request(
        API_URL,
        data=data,
        headers=headers,
        method="POST",
    )

    try:
        with urllib.request.urlopen(req, timeout=90) as resp:
            raw = resp.read().decode("utf-8")
            obj = json.loads(raw)
            emb = obj["data"][0]["embedding"]
            return np.array(emb, dtype=np.float32)

    except urllib.error.HTTPError as e:
        body = e.read().decode("utf-8", errors="replace")
        raise RuntimeError(f"HTTP {e.code}: {body}")


def tokenize(text: str) -> list[str]:
    """
    中文单字切分 + 英文缩写保持 + 数学符号保留。

    CKKS -> ckks
    C2S  -> c2s
    K004 -> k004
    σ⁻¹ -> σ⁻¹
    """
    tokens = []
    buf = ""

    for ch in str(text):
        if "\u4e00" <= ch <= "\u9fff":
            if buf:
                tokens.append(buf.lower())
                buf = ""
            tokens.append(ch)

        elif ch.isalnum() or ch in "-_./#":
            buf += ch

        elif ch in "⁰¹²³⁴⁵⁶⁷⁸⁹⁻⁺":
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


def load_data() -> tuple[dict[str, np.ndarray], dict]:
    script_dir = Path(__file__).resolve().parent
    data_dir = script_dir / "data"

    emb_path = data_dir / "embeddings.jsonl"
    bm25_path = data_dir / "bm25_index.pkl"

    if not emb_path.exists():
        print("ERROR: embeddings.jsonl 不存在。请先运行：")
        print("  python3 build_embeddings.py")
        sys.exit(1)

    if not bm25_path.exists():
        print("ERROR: bm25_index.pkl 不存在。请先运行：")
        print("  python3 build_bm25_index.py")
        sys.exit(1)

    embeddings = {}

    with emb_path.open("r", encoding="utf-8") as f:
        for line in f:
            line = line.strip()

            if not line:
                continue

            obj = json.loads(line)
            chunk_id = obj["chunk_id"]
            emb = np.array(obj["embedding"], dtype=np.float32)
            embeddings[chunk_id] = emb

    with bm25_path.open("rb") as f:
        bm25_index = pickle.load(f)

    return embeddings, bm25_index


def cosine_scores(query_emb: np.ndarray, embeddings: dict[str, np.ndarray]) -> dict[str, float]:
    q_norm = np.linalg.norm(query_emb)

    if q_norm == 0:
        return {cid: 0.0 for cid in embeddings}

    q = query_emb / (q_norm + 1e-12)
    scores = {}

    for cid, emb in embeddings.items():
        if len(emb) != len(query_emb):
            raise RuntimeError(
                f"向量维度不一致：query={len(query_emb)}, chunk={len(emb)}。"
                "请删除 data/embeddings.jsonl 后重新运行 build_embeddings.py。"
            )

        e = emb / (np.linalg.norm(emb) + 1e-12)
        scores[cid] = float(np.dot(q, e))

    return scores


def hybrid_search(query: str, top_k: int = 10, k_rrf: int = 60) -> list[dict]:
    model, api_key = get_config()
    embeddings, index = load_data()

    chunks = index["chunks"]
    bm25 = index["bm25"]

    chunk_id_to_idx = {c["chunk_id"]: i for i, c in enumerate(chunks)}

    # 1. Embedding 语义检索
    query_emb = embed_text(query, model, api_key)
    emb_scores = cosine_scores(query_emb, embeddings)
    emb_ranked = sorted(emb_scores.items(), key=lambda x: x[1], reverse=True)

    # 2. BM25 关键词检索
    tokenized_query = tokenize(query)
    bm25_raw_scores = bm25.get_scores(tokenized_query)
    bm25_max = float(max(bm25_raw_scores)) if len(bm25_raw_scores) and max(bm25_raw_scores) > 0 else 1.0

    bm25_scores = {
        chunks[i]["chunk_id"]: float(bm25_raw_scores[i]) / bm25_max
        for i in range(len(chunks))
    }

    bm25_ranked = sorted(bm25_scores.items(), key=lambda x: x[1], reverse=True)

    # 3. 精确 ID 匹配
    exact_boost = 2.0
    exact_matches = {}
    q_lower = query.lower()

    for card_id, idxs in index.get("card_id_index", {}).items():
        if card_id.lower() in q_lower:
            for idx in idxs:
                exact_matches[chunks[idx]["chunk_id"]] = exact_boost

    for qa_id, idxs in index.get("qa_id_index", {}).items():
        if qa_id.lower() in q_lower:
            for idx in idxs:
                exact_matches[chunks[idx]["chunk_id"]] = exact_boost

    for source_id, idxs in index.get("source_id_index", {}).items():
        if source_id.lower() in q_lower:
            for idx in idxs:
                exact_matches[chunks[idx]["chunk_id"]] = exact_boost

    # 4. RRF 融合
    rrf_scores = {}

    for rank, (chunk_id, _) in enumerate(emb_ranked):
        rrf_scores[chunk_id] = rrf_scores.get(chunk_id, 0.0) + 1.0 / (k_rrf + rank + 1)

    for rank, (chunk_id, _) in enumerate(bm25_ranked):
        rrf_scores[chunk_id] = rrf_scores.get(chunk_id, 0.0) + 1.0 / (k_rrf + rank + 1)

    for chunk_id, boost in exact_matches.items():
        rrf_scores[chunk_id] = rrf_scores.get(chunk_id, 0.0) + boost

    # 5. metadata filter：只保留 ready
    ready_set = {
        c["chunk_id"]
        for c in chunks
        if c.get("rag_status") == "ready"
    }

    rrf_scores = {
        chunk_id: score
        for chunk_id, score in rrf_scores.items()
        if chunk_id in ready_set
    }

    # 6. 排序 + 同一卡片知识点去重
    sorted_items = sorted(rrf_scores.items(), key=lambda x: x[1], reverse=True)

    seen_card_ids = set()
    deduped = []

    for chunk_id, score in sorted_items:
        idx = chunk_id_to_idx.get(chunk_id)

        if idx is None:
            continue

        chunk = chunks[idx]
        card_id = chunk.get("card_id", "")

        if card_id and chunk.get("chunk_type") == "card_knowledge_point":
            if card_id in seen_card_ids:
                continue
            seen_card_ids.add(card_id)

        deduped.append((chunk_id, score))

    # 7. 构造结果
    results = []

    emb_rank_pos = {cid: i for i, (cid, _) in enumerate(emb_ranked)}
    bm25_rank_pos = {cid: i for i, (cid, _) in enumerate(bm25_ranked)}

    for rank, (chunk_id, score) in enumerate(deduped[:top_k], start=1):
        idx = chunk_id_to_idx.get(chunk_id)

        if idx is None:
            continue

        chunk = chunks[idx]

        if chunk_id in exact_matches:
            retrieval_source = "exact_id_match"
        else:
            emb_r = emb_rank_pos.get(chunk_id, 999999)
            bm25_r = bm25_rank_pos.get(chunk_id, 999999)

            if emb_r < min(top_k * 2, 20) and bm25_r < min(top_k * 2, 20):
                retrieval_source = "hybrid"
            elif emb_r < bm25_r:
                retrieval_source = "embedding"
            else:
                retrieval_source = "bm25"

        results.append(
            {
                "rank": rank,
                "score": round(float(score), 4),
                "retrieval_source": retrieval_source,
                "chunk_id": chunk.get("chunk_id", ""),
                "chunk_type": chunk.get("chunk_type", ""),
                "card_id": chunk.get("card_id", ""),
                "qa_id": chunk.get("qa_id", ""),
                "source_id": chunk.get("source_id", ""),
                "domain": chunk.get("domain", ""),
                "topic": chunk.get("topic", ""),
                "confidence": chunk.get("confidence", ""),
                "canonical_card_path": chunk.get("canonical_card_path", ""),
                "text": chunk.get("text", ""),
            }
        )

    return results


def print_results(results: list[dict], verbose: bool = False, ids_only: bool = False) -> None:
    if ids_only:
        for r in results:
            print(f"{r['chunk_id']}  card={r['card_id'] or '-'}  qa={r['qa_id'] or '-'}")
        return

    header = (
        f"{'#':>3}  {'score':>8}  {'source':<16}  "
        f"{'chunk_id':<8}  {'card/qa':<12}  {'source_id':<12}  text"
    )

    print("-" * len(header))
    print(header)
    print("-" * len(header))

    for r in results:
        card_qa = r["card_id"] or r["qa_id"] or "-"
        preview = r["text"].replace("\n", " ")[:90]

        print(
            f"{r['rank']:>3}  {r['score']:>8.4f}  "
            f"{r['retrieval_source']:<16}  {r['chunk_id']:<8}  "
            f"{card_qa:<12}  {r['source_id']:<12}  {preview}"
        )

    print("-" * len(header))
    print(f"共 {len(results)} 条结果")

    if verbose:
        print()

        for r in results:
            print(f"--- [{r['chunk_id']}] {r['retrieval_source']} score={r['score']} ---")
            print(f"card_id={r['card_id']} qa_id={r['qa_id']} source_id={r['source_id']}")
            print(f"domain={r['domain']} topic={r['topic']} confidence={r['confidence']}")
            print(f"path={r['canonical_card_path']}")
            print(r["text"])
            print()


def main() -> None:
    parser = argparse.ArgumentParser(description="SiliconFlow Hybrid RAG 检索")
    parser.add_argument("query", nargs="?", help="查询文本")
    parser.add_argument("-k", "--top-k", type=int, default=10)
    parser.add_argument("-v", "--verbose", action="store_true")
    parser.add_argument("--json", action="store_true")
    parser.add_argument("--ids-only", action="store_true")

    args = parser.parse_args()

    if not args.query:
        print("用法: python3 query_rag.py '你的查询'")
        sys.exit(0)

    results = hybrid_search(args.query, top_k=args.top_k)

    if args.json:
        print(json.dumps(results, ensure_ascii=False, indent=2))
    else:
        print_results(results, verbose=args.verbose, ids_only=args.ids_only)


if __name__ == "__main__":
    main()