#!/usr/bin/env bash
set -euo pipefail

REPS="${REPS:-1}"
export GOFLAGS="${GOFLAGS:--mod=vendor}"

STAMP="$(date +%Y%m%d_%H%M%S)"
LOG_ROOT="${LOG_ROOT:-logs/subkey_current_params_${STAMP}}"
RESULT_ROOT="${RESULT_ROOT:-results/subkey_current_params_${STAMP}}"
MASTER_LOG="${LOG_ROOT}/master.log"
SUMMARY_TSV="${RESULT_ROOT}/subkey_current_params.tsv"

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

echo "=== split framework with subkey/Anyu leaf kernel ==="
echo "timestamp  : ${STAMP}"
echo "reps       : ${REPS}"
echo "GOFLAGS    : ${GOFLAGS}"
echo "master log : ${MASTER_LOG}"
echo "summary TSV: ${SUMMARY_TSV}"
echo

# Leaf15 exploratory current-PQ attempt: exact leaf logQP=868 with P=0.
# P=0 means LogP=[] / no key-switching auxiliary primes. This is not an
# AnyuWang paper parameter; it is only the closest match to the current strict
# leaf15 security target.
# If this fails during keygen/runtime, use run_subkey_legacy_lowq_p1.sh as
# a compatibility-only fallback and document that it is not PQ-128 under leaf15.
run_one "subkey_top16_leaf15_lowq_p0_qp868" \
  -logN=16 \
  -layers=1 \
  -run=ours \
  -leafEngine=anyu \
  -anyuProfile=lowq \
  -anyuMaxLogQP=868 \
  -anyuNumP=0 \
  -splitFirstLeafPreprocess=true \
  -parallelLeaves=true \
  -leafWorkers=2

run_one "subkey_top17_leaf15_lowq_p0_qp868" \
  -logN=17 \
  -layers=2 \
  -run=ours \
  -leafEngine=anyu \
  -anyuProfile=lowq \
  -anyuMaxLogQP=868 \
  -anyuNumP=0 \
  -splitFirstLeafPreprocess=true \
  -parallelLeaves=true \
  -leafWorkers=4

# Leaf16 current high-output subkey profile. Dry-run shows logQP ~= 1730
# and it passes the local leaf16 128-bit PQ reference.
run_one "subkey_top17_leaf16_i1" \
  -logN=17 \
  -layers=1 \
  -run=ours \
  -leafEngine=anyu \
  -anyuProfile=i1 \
  -splitFirstLeafPreprocess=true \
  -parallelLeaves=true \
  -leafWorkers=2

echo "=== subkey current-parameter runs finished ==="
echo "master log : ${MASTER_LOG}"
echo "summary TSV: ${SUMMARY_TSV}"
