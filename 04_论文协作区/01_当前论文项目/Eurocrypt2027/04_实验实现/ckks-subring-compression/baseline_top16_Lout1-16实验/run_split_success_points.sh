#!/usr/bin/env bash
set -euo pipefail

REPS="${REPS:-1}"
export GOFLAGS="${GOFLAGS:--mod=vendor}"

STAMP="$(date +%Y%m%d_%H%M%S)"
LOG_ROOT="${LOG_ROOT:-logs/split_success_points_${STAMP}}"
RESULT_ROOT="${RESULT_ROOT:-results/split_success_points_${STAMP}}"
MASTER_LOG="${LOG_ROOT}/master.log"
SUMMARY_TSV="${RESULT_ROOT}/split_success_points.tsv"

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

echo "=== split success points ==="
echo "timestamp  : ${STAMP}"
echo "reps       : ${REPS}"
echo "GOFLAGS    : ${GOFLAGS}"
echo "master log : ${MASTER_LOG}"
echo "summary TSV: ${SUMMARY_TSV}"
echo

run_one "split_top16_leaf15_lout3_n15qp849" \
  -logN=16 \
  -layers=1 \
  -run=ours \
  -leafEngine=lattigo \
  -lattigoPreset=n15qp849 \
  -splitFirstLeafPreprocess=true \
  -parallelLeaves=true \
  -leafWorkers=2

run_one "split_top17_leaf15_lout3_n15qp849" \
  -logN=17 \
  -layers=2 \
  -run=ours \
  -leafEngine=lattigo \
  -lattigoPreset=n15qp849 \
  -splitFirstLeafPreprocess=true \
  -parallelLeaves=true \
  -leafWorkers=4

run_one "split_top17_leaf16_lout15_default" \
  -logN=17 \
  -layers=1 \
  -run=ours \
  -leafEngine=lattigo \
  -q0Bits=45 \
  -circuitLevels=15 \
  -circuitPrimeBits=35 \
  -numP=5 \
  -leafNumP=5 \
  -defaultScale=35 \
  -splitFirstLeafPreprocess=true \
  -parallelLeaves=true \
  -leafWorkers=2

echo "=== split runs finished ==="
echo "master log : ${MASTER_LOG}"
echo "summary TSV: ${SUMMARY_TSV}"
