#!/usr/bin/env bash
set -euo pipefail

export GOFLAGS="${GOFLAGS:--mod=vendor}"

echo "=== subkey/Anyu leaf-kernel parameter dry runs ==="
echo "GOFLAGS: ${GOFLAGS}"
echo

run_one() {
  local name="$1"
  shift
  echo "----------------------------------------------------------------------"
  echo "[dry-run] ${name}"
  echo "----------------------------------------------------------------------"
  go run . -dryRunParams=true "$@"
  echo
}

run_one "top16 -> leaf15, Anyu lowq, current leaf15 PQ target logQP=868, P=0" \
  -logN=16 \
  -layers=1 \
  -run=ours \
  -leafEngine=anyu \
  -anyuProfile=lowq \
  -anyuMaxLogQP=868 \
  -anyuNumP=0 \
  -splitFirstLeafPreprocess=true

run_one "top17 -> leaf15, Anyu lowq, current leaf15 PQ target logQP=868, P=0" \
  -logN=17 \
  -layers=2 \
  -run=ours \
  -leafEngine=anyu \
  -anyuProfile=lowq \
  -anyuMaxLogQP=868 \
  -anyuNumP=0 \
  -splitFirstLeafPreprocess=true

run_one "top17 -> leaf16, Anyu i1, current high-output profile" \
  -logN=17 \
  -layers=1 \
  -run=ours \
  -leafEngine=anyu \
  -anyuProfile=i1 \
  -splitFirstLeafPreprocess=true

echo "=== dry runs finished ==="
