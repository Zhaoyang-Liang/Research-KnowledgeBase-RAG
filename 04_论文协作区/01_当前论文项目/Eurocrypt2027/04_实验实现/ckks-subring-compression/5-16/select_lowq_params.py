#!/usr/bin/env python3
"""
Scan low-output-level CKKS bootstrapping parameters.

The script calls the local Go experiment with -dryRunParams=true and parses the
actual Lattigo bootstrapping logQ/logQP. This is intentionally a dry scanner:
it helps choose candidates before spending minutes on full key generation and
bootstrapping.

Example:

  ./select_lowq_params.py \
    --top-logN 16 --layers 1 \
    --q0 26:38 --scale 20:35 --circuit-prime 18:35 \
    --numP 1 --security-modulus Q --max-log-sec 868 \
    --top 20

For conservative accounting, use:

  --security-modulus QP
"""

from __future__ import annotations

import argparse
import math
import re
import subprocess
from dataclasses import dataclass
from pathlib import Path


ROOT = Path(__file__).resolve().parent


TOP_RE = re.compile(
    r"top:\s+logN=(?P<logN>\d+)\s+N=(?P<N>\d+)\s+qCount=(?P<qCount>\d+)\s+"
    r"pCount=(?P<pCount>\d+)\s+logQ=(?P<logQ>[0-9.]+)\s+logQP=(?P<logQP>[0-9.]+)"
)
LEAF_RE = re.compile(
    r"leaf:\s+logN=(?P<logN>\d+)\s+N=(?P<N>\d+)\s+qCount=(?P<qCount>\d+)\s+"
    r"pCount=(?P<pCount>\d+)\s+logQ=(?P<logQ>[0-9.]+)\s+logQP=(?P<logQP>[0-9.]+)"
)


@dataclass
class Candidate:
    q0: int
    scale: int
    circuit_prime: int
    num_p: int
    top_logq: float
    top_logqp: float
    top_q_count: int
    top_p_count: int
    leaf_logq: float
    leaf_logqp: float
    leaf_q_count: int
    leaf_p_count: int
    scale_gap: int
    security_metric: float
    security_margin: float
    score: float


def parse_range(spec: str) -> list[int]:
    if "," in spec:
        return [int(x) for x in spec.split(",") if x.strip()]
    if ":" in spec:
        parts = [int(x) for x in spec.split(":")]
        if len(parts) == 2:
            start, stop = parts
            step = 1
        elif len(parts) == 3:
            start, stop, step = parts
        else:
            raise ValueError(f"bad range: {spec}")
        return list(range(start, stop + 1, step))
    return [int(spec)]


def dry_run(top_logn: int, layers: int, q0: int, scale: int, circuit_prime: int, num_p: int) -> Candidate | None:
    cmd = [
        "go",
        "run",
        ".",
        "-dryRunParams=true",
        f"-logN={top_logn}",
        f"-layers={layers}",
        f"-q0Bits={q0}",
        "-circuitLevels=1",
        f"-circuitPrimeBits={circuit_prime}",
        f"-defaultScale={scale}",
        f"-numP={num_p}",
        f"-leafNumP={num_p}",
    ]

    proc = subprocess.run(cmd, cwd=ROOT, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    if proc.returncode != 0:
        return None

    top = TOP_RE.search(proc.stdout)
    leaf = LEAF_RE.search(proc.stdout)
    if not top or not leaf:
        return None

    top_logq = float(top.group("logQ"))
    top_logqp = float(top.group("logQP"))
    leaf_logq = float(leaf.group("logQ"))
    leaf_logqp = float(leaf.group("logQP"))

    return Candidate(
        q0=q0,
        scale=scale,
        circuit_prime=circuit_prime,
        num_p=num_p,
        top_logq=top_logq,
        top_logqp=top_logqp,
        top_q_count=int(top.group("qCount")),
        top_p_count=int(top.group("pCount")),
        leaf_logq=leaf_logq,
        leaf_logqp=leaf_logqp,
        leaf_q_count=int(leaf.group("qCount")),
        leaf_p_count=int(leaf.group("pCount")),
        scale_gap=q0 - scale,
        security_metric=math.nan,
        security_margin=math.nan,
        score=math.nan,
    )


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--top-logN", type=int, default=16)
    parser.add_argument("--layers", type=int, default=1)
    parser.add_argument("--q0", default="26:38")
    parser.add_argument("--scale", default="20:35")
    parser.add_argument("--circuit-prime", default="18:35")
    parser.add_argument("--numP", type=int, default=1)
    parser.add_argument("--security-modulus", choices=["Q", "QP"], default="QP")
    parser.add_argument("--max-log-sec", type=float, default=868.0)
    parser.add_argument("--target-logq", type=float, default=868.0)
    parser.add_argument("--min-scale-gap", type=int, default=8)
    parser.add_argument("--min-scale", type=int, default=0)
    parser.add_argument("--top", type=int, default=20)
    args = parser.parse_args()

    candidates: list[Candidate] = []

    for q0 in parse_range(args.q0):
        for scale in parse_range(args.scale):
            if scale < args.min_scale:
                continue
            for circuit_prime in parse_range(args.circuit_prime):
                cand = dry_run(args.top_logN, args.layers, q0, scale, circuit_prime, args.numP)
                if cand is None:
                    continue

                metric = cand.leaf_logq if args.security_modulus == "Q" else cand.leaf_logqp
                margin = args.max_log_sec - metric
                scale_penalty = max(0, args.min_scale_gap - cand.scale_gap) * 1000.0
                security_penalty = max(0.0, -margin) * 100.0
                target_penalty = abs(cand.leaf_logq - args.target_logq)

                cand.security_metric = metric
                cand.security_margin = margin
                cand.score = scale_penalty + security_penalty + target_penalty
                candidates.append(cand)

    candidates.sort(key=lambda c: (c.score, abs(c.leaf_logq - args.target_logq), -c.scale))

    print(
        "q0\tscale\tgap\tcp\tnumP\tleafLogQ\tleafLogQP\tsecMetric\tsecMargin\t"
        "scaleOK\tsecOK\tscore\tfullRunCommand"
    )
    for cand in candidates[: args.top]:
        scale_ok = cand.scale_gap >= args.min_scale_gap
        sec_ok = cand.security_margin >= 0
        cmd = (
            f"go run . -logN={args.top_logN} -layers={args.layers} -run=all "
            f"-leafEngine=lattigo -q0Bits={cand.q0} -circuitLevels=1 "
            f"-circuitPrimeBits={cand.circuit_prime} -numP={cand.num_p} "
            f"-leafNumP={cand.num_p} -defaultScale={cand.scale} "
            f"-splitFirstLeafPreprocess=true -parallelLeaves=true -leafWorkers=2 -reps=1"
        )
        print(
            f"{cand.q0}\t{cand.scale}\t{cand.scale_gap}\t{cand.circuit_prime}\t{cand.num_p}\t"
            f"{cand.leaf_logq:.1f}\t{cand.leaf_logqp:.1f}\t{cand.security_metric:.1f}\t"
            f"{cand.security_margin:.1f}\t{scale_ok}\t{sec_ok}\t{cand.score:.2f}\t{cmd}"
        )


if __name__ == "__main__":
    main()

