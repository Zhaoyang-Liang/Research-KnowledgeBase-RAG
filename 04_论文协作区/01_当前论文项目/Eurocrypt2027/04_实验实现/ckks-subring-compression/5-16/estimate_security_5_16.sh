#!/usr/bin/env bash
set -euo pipefail

ESTIMATOR_DIR="${ESTIMATOR_DIR:-/Users/mac/opt/lattice-estimator}"
SAGE_PY="${SAGE_PY:-/opt/anaconda3/envs/sage/bin/python}"
SAGE_BIN_DIR="$(dirname "${SAGE_PY}")"

cd "${ESTIMATOR_DIR}"

run_one() {
  local tag="$1"
  local n="$2"
  local logqp="$3"
  echo
  echo "===== ${tag}: n=${n}, logQP=${logqp} ====="
  PATH="${SAGE_BIN_DIR}:$PATH" "${SAGE_PY}" fhe_estimate.py \
    --n "${n}" \
    --logQP "${logqp}" \
    --sigma 3.2 \
    --tag "${tag}" \
    --attacks primal_usvp,primal_bdd,dual_hybrid
}

run_one "leaf_logN15_lowq_p20_actual_QP" 32767 936
run_one "leaf_logN15_p35_Q_only_diagnostic" 32767 901
run_one "leaf_logN15_p35_actual_QP" 32767 962
run_one "leaf_logN15_threshold_128" 32767 868
run_one "leaf_logN16_p35_actual_QP_sec128_adjusted" 65535 962

