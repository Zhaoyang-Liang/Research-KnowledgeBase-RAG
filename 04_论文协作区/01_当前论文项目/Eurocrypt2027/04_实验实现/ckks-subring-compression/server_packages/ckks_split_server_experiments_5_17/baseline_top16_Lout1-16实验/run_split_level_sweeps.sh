#!/usr/bin/env bash
set -euo pipefail

REPS="${REPS:-1}"
export GOFLAGS="${GOFLAGS:--mod=vendor}"

STAMP="$(date +%Y%m%d_%H%M%S)"
LOG_ROOT="${LOG_ROOT:-logs/split_level_sweeps_${STAMP}}"
RESULT_ROOT="${RESULT_ROOT:-results/split_level_sweeps_${STAMP}}"
MASTER_LOG="${LOG_ROOT}/master.log"
SUMMARY_TSV="${RESULT_ROOT}/split_level_sweeps.tsv"

mkdir -p "${LOG_ROOT}" "${RESULT_ROOT}"

exec > >(tee -a "${MASTER_LOG}") 2>&1

run_one() {
  local name="$1"
  shift
  local log_file="${LOG_ROOT}/${name}.log"

  echo "----------------------------------------------------------------------"
  echo "[run] ${name}"
  echo "[run] log file = ${log_file}"
  echo "----------------------------------------------------------------------"

  go run . "$@" \
    -reps="${REPS}" \
    -outputTSV="${SUMMARY_TSV}" \
    2>&1 | tee "${log_file}"

  echo "[done] ${name}"
  echo
}

echo "=== split level sweeps ==="
echo "timestamp  : ${STAMP}"
echo "reps       : ${REPS}"
echo "GOFLAGS    : ${GOFLAGS}"
echo "master log : ${MASTER_LOG}"
echo "summary TSV: ${SUMMARY_TSV}"
echo

echo "=== sweep A: top16 -> leaf15, output level 1..3 ==="
for level in 1 2 3; do
  run_one "split_top16_leaf15_lout${level}" \
    -logN=16 \
    -layers=1 \
    -run=ours \
    -leafEngine=lattigo \
    -lattigoPreset="n15lout${level}" \
    -splitFirstLeafPreprocess=true \
    -parallelLeaves=true \
    -leafWorkers=2
done

echo "=== sweep B: top17 -> leaf15, output level 1..3 ==="
for level in 1 2 3; do
  run_one "split_top17_leaf15_lout${level}" \
    -logN=17 \
    -layers=2 \
    -run=ours \
    -leafEngine=lattigo \
    -lattigoPreset="n15lout${level}" \
    -splitFirstLeafPreprocess=true \
    -parallelLeaves=true \
    -leafWorkers=4
done

echo "=== sweep C: top17 -> leaf16, output level 1..15 ==="
for level in $(seq 1 15); do
  run_one "split_top17_leaf16_lout${level}" \
    -logN=17 \
    -layers=1 \
    -run=ours \
    -leafEngine=lattigo \
    -q0Bits=45 \
    -circuitLevels="${level}" \
    -circuitPrimeBits=35 \
    -numP=5 \
    -leafNumP=5 \
    -defaultScale=35 \
    -splitFirstLeafPreprocess=true \
    -parallelLeaves=true \
    -leafWorkers=2
done

echo "=== split level sweeps finished ==="
echo "master log : ${MASTER_LOG}"
echo "summary TSV: ${SUMMARY_TSV}"
