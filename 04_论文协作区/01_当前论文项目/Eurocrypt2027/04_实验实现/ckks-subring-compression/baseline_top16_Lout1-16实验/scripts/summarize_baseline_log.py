#!/usr/bin/env python3
import argparse
import csv
import os
import re


def parse_duration_seconds(text: str) -> float:
    text = text.strip()
    units = {
        "h": 3600.0,
        "m": 60.0,
        "s": 1.0,
        "ms": 1e-3,
        "us": 1e-6,
        "µs": 1e-6,
        "ns": 1e-9,
    }
    parts = re.findall(r"([0-9.]+)(h|ms|us|µs|ns|m|s)", text)
    if not parts or "".join(value + unit for value, unit in parts) != text:
        raise ValueError(f"unsupported duration format: {text}")
    return sum(float(value) * units[unit] for value, unit in parts)


def find(pattern: str, data: str, default: str = "") -> str:
    match = re.search(pattern, data, re.MULTILINE)
    return match.group(1) if match else default


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--target-level", type=int, required=True)
    parser.add_argument("--log", required=True)
    parser.add_argument("--summary", required=True)
    args = parser.parse_args()

    with open(args.log, "r", encoding="utf-8", errors="replace") as f:
        data = f.read()

    logqp = find(r"Bootstrapping parameters:.*logQP=([0-9.]+)", data)
    q_count = find(r"Bootstrapping parameters:.*levels=([0-9]+)", data)
    actual_level = find(r"Output ciphertext level = ([0-9]+)", data)
    precision = find(r"Baseline\s+\(k=1\).*Avg L1 prec:\s+([0-9.]+) bits", data)
    total_raw = find(r"Bootstrapping finished in ([^\n]+)", data)
    stc_raw = find(r"StC:\s+([^,]+)", data)
    modup_raw = find(r"ModUp:\s+([^,]+)", data)
    cts_raw = find(r"CtS:\s+([^,]+)", data)
    evalmod_raw = find(r"EvalMod:\s+([^\n]+)", data)

    row = {
        "target_lout": args.target_level,
        "actual_output_level": actual_level,
        "logQP": logqp,
        "qCount": q_count,
        "total_sec": f"{parse_duration_seconds(total_raw):.9f}" if total_raw else "",
        "stc_sec": f"{parse_duration_seconds(stc_raw):.9f}" if stc_raw else "",
        "modup_sec": f"{parse_duration_seconds(modup_raw):.9f}" if modup_raw else "",
        "cts_sec": f"{parse_duration_seconds(cts_raw):.9f}" if cts_raw else "",
        "evalmod_sec": f"{parse_duration_seconds(evalmod_raw):.9f}" if evalmod_raw else "",
        "avg_l1_prec_bits": precision,
        "log_file": args.log,
    }

    os.makedirs(os.path.dirname(args.summary), exist_ok=True)
    need_header = not os.path.exists(args.summary) or os.path.getsize(args.summary) == 0
    with open(args.summary, "a", encoding="utf-8", newline="") as f:
        writer = csv.DictWriter(f, fieldnames=list(row.keys()), delimiter="\t")
        if need_header:
            writer.writeheader()
        writer.writerow(row)


if __name__ == "__main__":
    main()
