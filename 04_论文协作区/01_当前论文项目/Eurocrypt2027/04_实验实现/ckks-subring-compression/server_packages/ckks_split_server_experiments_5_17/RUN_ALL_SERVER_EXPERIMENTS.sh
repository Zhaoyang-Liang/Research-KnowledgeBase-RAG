#!/usr/bin/env bash
set -euo pipefail

REPS="${REPS:-1}"
RUN_LEAF15="${RUN_LEAF15:-0}"
STAMP="$(date +%Y%m%d_%H%M%S)"
TOP_LOG="overnight_run_${STAMP}.log"

exec > >(tee -a "${TOP_LOG}") 2>&1

echo "=== CKKS split server experiments: overnight run ==="
echo "timestamp  : ${STAMP}"
echo "reps       : ${REPS}"
echo "run leaf15?: ${RUN_LEAF15}"
echo "top log    : ${TOP_LOG}"
echo

echo "=== environment ==="
pwd
go version
echo

echo "=== group 1/2: Lattigo split level sweeps ==="
(
  cd "baseline_top16_Lout1-16实验"
  chmod +x run_split_level_sweeps.sh
  REPS="${REPS}" ./run_split_level_sweeps.sh
)
echo "=== group 1/2 finished ==="
echo

echo "=== group 2/2: top17 i1 subkey compatibility ==="
(
  cd "5-16-内核subkey"
  chmod +x run_subkey_top17_i1_group.sh
  REPS="${REPS}" RUN_LEAF15="${RUN_LEAF15}" ./run_subkey_top17_i1_group.sh
)
echo "=== group 2/2 finished ==="
echo

echo "=== result TSV files ==="
find . -path './*/results/*/*.tsv' -print | sort
echo
echo "=== all server experiments finished ==="
echo "top log: ${TOP_LOG}"
