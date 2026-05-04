#!/usr/bin/env python3
"""
build_embeddings.py

SiliconFlow 直连版 embedding 构建脚本。

特点：
- 不使用 OpenAI SDK
- 不依赖 source config.env 是否 export
- 自动读取当前目录下的 config.env
- 只处理 rag_status=ready 的 chunks
- 一次请求一个 chunk，最稳
- 输出 data/embeddings.jsonl

用法：
  cd ~/Desktop/Research-KB-RAG/06_RAG入库准备/api_rag
  python3 build_embeddings.py
"""

import json
import os
import sys
import time
import urllib.request
import urllib.error
from pathlib import Path


API_URL = "https://api.siliconflow.cn/v1/embeddings"


def read_config_file(config_path: Path) -> dict:
    """
    读取 config.env。支持：
      KEY=value
      export KEY=value
      KEY="value"
      KEY='value'
    """
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

        # 去掉行尾注释，但避免破坏 key 本身
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
        print("请在 config.env 中设置：EMBEDDING_PROVIDER=siliconflow")
        sys.exit(1)

    if not api_key:
        print("ERROR: 没有读到 EMBEDDING_API_KEY。")
        print("请检查当前目录下的 config.env 是否包含：")
        print("  EMBEDDING_API_KEY=你的SiliconFlow密钥")
        sys.exit(1)

    if api_key in {"sk-your-key-here", "your-key", "YOUR_KEY"}:
        print("ERROR: EMBEDDING_API_KEY 仍然是占位符，请填入真实 SiliconFlow API key。")
        sys.exit(1)

    return model, api_key


def siliconflow_embed_one(text: str, model: str, api_key: str, retries: int = 3) -> tuple[list[float], int]:
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

    last_error = None

    for attempt in range(1, retries + 1):
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
                tokens = int(obj.get("usage", {}).get("total_tokens", 0) or 0)

                if not isinstance(emb, list) or not emb:
                    raise RuntimeError(f"返回 embedding 为空：{obj}")

                return emb, tokens

        except urllib.error.HTTPError as e:
            body = e.read().decode("utf-8", errors="replace")
            last_error = f"HTTP {e.code}: {body}"

            if e.code in {429, 500, 502, 503, 504} and attempt < retries:
                time.sleep(1.5 * attempt)
                continue

            raise RuntimeError(last_error)

        except Exception as e:
            last_error = repr(e)

            if attempt < retries:
                time.sleep(1.5 * attempt)
                continue

            raise RuntimeError(last_error)

    raise RuntimeError(last_error or "unknown error")


def load_ready_chunks() -> list[dict]:
    kb = Path.home() / "Desktop" / "Research-KB-RAG"
    chunk_dir = kb / "06_RAG入库准备" / "chunks"

    chunks = []

    for fname in ["cards_ready_chunks.jsonl", "qa_ready_chunks.jsonl"]:
        fp = chunk_dir / fname

        if not fp.exists():
            print(f"WARNING: 找不到 {fp}，跳过。")
            continue

        with fp.open("r", encoding="utf-8") as f:
            for line in f:
                line = line.strip()

                if not line:
                    continue

                obj = json.loads(line)

                if obj.get("rag_status") == "ready":
                    chunks.append(obj)

    return chunks


def main() -> None:
    model, api_key = get_config()

    print("Provider: SiliconFlow")
    print(f"Model:    {model}")
    print(f"API URL:  {API_URL}")
    print("Key:      已读取，不显示")
    print()

    script_dir = Path(__file__).resolve().parent
    out_dir = script_dir / "data"
    out_dir.mkdir(exist_ok=True)

    out_path = out_dir / "embeddings.jsonl"
    tmp_path = out_dir / "embeddings.jsonl.tmp"

    chunks = load_ready_chunks()
    print(f"✓ 加载 {len(chunks)} 个 ready chunks")

    if not chunks:
        print("ERROR: 没有 ready chunks。")
        sys.exit(1)

    success = 0
    failed = 0
    total_tokens = 0
    dim = None

    with tmp_path.open("w", encoding="utf-8") as out:
        for idx, chunk in enumerate(chunks, start=1):
            chunk_id = chunk.get("chunk_id", "")
            text = chunk.get("text", "")

            if not chunk_id:
                print(f"  ✗ 第 {idx} 条缺少 chunk_id，跳过")
                failed += 1
                continue

            if not text.strip():
                print(f"  ✗ {chunk_id} text 为空，跳过")
                failed += 1
                continue

            # BGE-M3 支持较长输入；这里保守截断超长文本，避免异常。
            # 当前 chunk 通常很短，不会触发。
            safe_text = text[:12000]

            try:
                embedding, tokens = siliconflow_embed_one(
                    text=safe_text,
                    model=model,
                    api_key=api_key,
                )

                dim = dim or len(embedding)
                total_tokens += tokens

                row = {
                    "chunk_id": chunk_id,
                    "embedding": embedding,
                    "provider": "siliconflow",
                    "model": model,
                }

                out.write(json.dumps(row, ensure_ascii=False) + "\n")
                out.flush()

                success += 1

                if idx % 25 == 0 or idx == len(chunks):
                    print(f"  {idx}/{len(chunks)} chunks, success={success}, failed={failed}")

                time.sleep(0.05)

            except Exception as e:
                failed += 1
                print(f"  ✗ chunk {idx}/{len(chunks)} FAILED: {chunk_id}")
                print(f"    {e}")

                if "401" in str(e):
                    print()
                    print("认证失败。请检查 config.env 中的 EMBEDDING_API_KEY。")
                    print("注意：curl 成功但脚本失败时，通常是 config.env 中的 key 和 curl 实际使用的 key 不一致。")
                    sys.exit(1)

                if failed >= 10:
                    print()
                    print("连续失败较多，停止。")
                    sys.exit(1)

    if success == 0:
        print("ERROR: 没有成功生成任何 embedding。")
        sys.exit(1)

    tmp_path.replace(out_path)

    print()
    print("=" * 60)
    print("✓ 完成 — SiliconFlow embeddings")
    print(f"  chunks 成功:   {success}")
    print(f"  chunks 失败:   {failed}")
    print(f"  向量维度:      {dim}")
    print(f"  总 tokens:     {total_tokens}")
    print(f"  输出文件:      {out_path}")
    print("=" * 60)


if __name__ == "__main__":
    main()