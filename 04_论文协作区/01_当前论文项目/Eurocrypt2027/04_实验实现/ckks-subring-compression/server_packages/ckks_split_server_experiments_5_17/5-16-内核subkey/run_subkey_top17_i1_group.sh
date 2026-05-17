#!/usr/bin/env bash
set -euo pipefail

REPS="${REPS:-1}"
RUN_LEAF15="${RUN_LEAF15:-0}"
export GOFLAGS="${GOFLAGS:--mod=vendor}"

STAMP="$(date +%Y%m%d_%H%M%S)"
LOG_ROOT="${LOG_ROOT:-logs/subkey_top17_i1_group_${STAMP}}"
RESULT_ROOT="${RESULT_ROOT:-results/subkey_top17_i1_group_${STAMP}}"
MASTER_LOG="${LOG_ROOT}/master.log"
SUMMARY_TSV="${RESULT_ROOT}/subkey_top17_i1_group.tsv"

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

echo "=== top17 i1 subkey compatibility group ==="
echo "timestamp   : ${STAMP}"
echo "reps        : ${REPS}"
echo "run leaf15? : ${RUN_LEAF15}"
echo "GOFLAGS     : ${GOFLAGS}"
echo "master log  : ${MASTER_LOG}"
echo "summary TSV : ${SUMMARY_TSV}"
echo

echo "=== A. direct subkey baseline / no split, top17, Lout=15 ==="
run_one "top17_direct_subkey_i1_lout15" \
  -logN=17 \
  -layers=1 \
  -run=baseline \
  -leafEngine=anyu \
  -anyuProfile=i1 \
  -anyuSubringLogN=12 \
  -splitFirstLeafPreprocess=true

echo "=== B. split + subkey leaf kernel, top17 -> leaf16, Lout=15 ==="
run_one "top17_split_leaf16_subkey_i1_lout15" \
  -logN=17 \
  -layers=1 \
  -run=ours \
  -leafEngine=anyu \
  -anyuProfile=i1 \
  -anyuSubringLogN=12 \
  -splitFirstLeafPreprocess=true \
  -parallelLeaves=true \
  -leafWorkers=2

if [[ "${RUN_LEAF15}" == "1" ]]; then
  echo "=== C. optional exploratory leaf15 lowq/P=0 attempts ==="
  echo "These are not paper-faithful AnyuWang parameters; use only as exploratory checks."

  run_one "top16_split_leaf15_subkey_lowq_p0_qp868" \
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

  run_one "top17_split_leaf15_subkey_lowq_p0_qp868" \
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
fi

echo "=== top17 i1 subkey compatibility group finished ==="
echo "master log  : ${MASTER_LOG}"
echo "summary TSV : ${SUMMARY_TSV}"
