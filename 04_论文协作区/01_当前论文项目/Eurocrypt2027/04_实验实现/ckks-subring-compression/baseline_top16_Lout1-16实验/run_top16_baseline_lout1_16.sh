#!/usr/bin/env bash
set -euo pipefail

REPS="${REPS:-1}"
START_LEVEL="${START_LEVEL:-1}"
END_LEVEL="${END_LEVEL:-16}"

Q0_BITS="${Q0_BITS:-45}"
CIRCUIT_PRIME_BITS="${CIRCUIT_PRIME_BITS:-35}"
NUM_P="${NUM_P:-5}"
DEFAULT_SCALE="${DEFAULT_SCALE:-35}"
export GOFLAGS="${GOFLAGS:--mod=vendor}"

STAMP="$(date +%Y%m%d_%H%M%S)"
LOG_ROOT="${LOG_ROOT:-logs/top16_baseline_lout1_16_${STAMP}}"
RESULT_ROOT="${RESULT_ROOT:-results/top16_baseline_lout1_16_${STAMP}}"
MASTER_LOG="${LOG_ROOT}/master.log"
SUMMARY_TSV="${RESULT_ROOT}/summary.tsv"
RAW_TSV="${RESULT_ROOT}/raw_go_summary.tsv"

mkdir -p "${LOG_ROOT}" "${RESULT_ROOT}"

exec > >(tee -a "${MASTER_LOG}") 2>&1

echo "=== top16 direct baseline sweep: output level ${START_LEVEL}..${END_LEVEL} ==="
echo "timestamp       : ${STAMP}"
echo "reps            : ${REPS}"
echo "q0Bits          : ${Q0_BITS}"
echo "circuitPrimeBits: ${CIRCUIT_PRIME_BITS}"
echo "numP            : ${NUM_P}"
echo "defaultScale    : ${DEFAULT_SCALE}"
echo "GOFLAGS         : ${GOFLAGS}"
echo "master log      : ${MASTER_LOG}"
echo "summary TSV     : ${SUMMARY_TSV}"
echo "raw Go TSV      : ${RAW_TSV}"
echo

for level in $(seq "${START_LEVEL}" "${END_LEVEL}"); do
  level_log="${LOG_ROOT}/baseline_top16_Lout${level}.log"

  echo "----------------------------------------------------------------------"
  echo "[run] target output level = ${level}"
  echo "[run] log file = ${level_log}"
  echo "----------------------------------------------------------------------"

  go run . \
    -logN=16 \
    -layers=1 \
    -run=baseline \
    -leafEngine=lattigo \
    -q0Bits="${Q0_BITS}" \
    -circuitLevels="${level}" \
    -circuitPrimeBits="${CIRCUIT_PRIME_BITS}" \
    -numP="${NUM_P}" \
    -defaultScale="${DEFAULT_SCALE}" \
    -reps="${REPS}" \
    -outputTSV="${RAW_TSV}" \
    2>&1 | tee "${level_log}"

  python3 scripts/summarize_baseline_log.py \
    --target-level "${level}" \
    --log "${level_log}" \
    --summary "${SUMMARY_TSV}"

  echo "[done] target output level = ${level}"
  tail -n 1 "${SUMMARY_TSV}"
  echo
done

echo "=== sweep finished ==="
echo "master log : ${MASTER_LOG}"
echo "summary TSV: ${SUMMARY_TSV}"
echo "raw Go TSV : ${RAW_TSV}"
