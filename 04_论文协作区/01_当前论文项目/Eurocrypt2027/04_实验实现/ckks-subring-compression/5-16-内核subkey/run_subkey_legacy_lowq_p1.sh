#!/usr/bin/env bash
set -euo pipefail

REPS="${REPS:-1}"
export GOFLAGS="${GOFLAGS:--mod=vendor}"

STAMP="$(date +%Y%m%d_%H%M%S)"
LOG_ROOT="${LOG_ROOT:-logs/subkey_legacy_lowq_p1_${STAMP}}"
RESULT_ROOT="${RESULT_ROOT:-results/subkey_legacy_lowq_p1_${STAMP}}"
MASTER_LOG="${LOG_ROOT}/master.log"
SUMMARY_TSV="${RESULT_ROOT}/subkey_legacy_lowq_p1.tsv"

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

echo "=== legacy lowq subkey compatibility fallback ==="
echo "WARNING: this uses P=1 and leaf15 logQP ~= 929, so it is not PQ-128 under the current strict leaf15 target."
echo "Use only as backend-compatibility evidence if the current P=0 attempt fails."
echo

run_one "subkey_top16_leaf15_legacy_lowq_p1_qp929" \
  -logN=16 \
  -layers=1 \
  -run=ours \
  -leafEngine=anyu \
  -anyuProfile=lowq \
  -splitFirstLeafPreprocess=true \
  -parallelLeaves=true \
  -leafWorkers=2

run_one "subkey_top17_leaf15_legacy_lowq_p1_qp929" \
  -logN=17 \
  -layers=2 \
  -run=ours \
  -leafEngine=anyu \
  -anyuProfile=lowq \
  -splitFirstLeafPreprocess=true \
  -parallelLeaves=true \
  -leafWorkers=4

echo "=== legacy fallback runs finished ==="
echo "master log : ${MASTER_LOG}"
echo "summary TSV: ${SUMMARY_TSV}"
