#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"
mkdir -p results logs

LOG="logs/lowq_lattigo_p35_sec128_$(date +%Y%m%d_%H%M%S).log"

echo "[run] 128-bit leaf-secure higher-precision low-output-level top BTS vs split BTS"
echo "[log] ${LOG}"
echo "[note] top logN=17, split leaf logN=16; estimator gives leaf logQP≈962 well above 128 bits."

go run . \
  -logN=17 \
  -layers=1 \
  -run=all \
  -leafEngine=lattigo \
  -q0Bits=45 \
  -circuitLevels=1 \
  -circuitPrimeBits=35 \
  -numP=1 \
  -leafNumP=1 \
  -defaultScale=35 \
  -splitFirstLeafPreprocess=true \
  -parallelLeaves=true \
  -leafWorkers=2 \
  -reps="${REPS:-1}" \
  -outputTSV=results/lowq_lattigo_p35_sec128.tsv \
  2>&1 | tee "${LOG}"

