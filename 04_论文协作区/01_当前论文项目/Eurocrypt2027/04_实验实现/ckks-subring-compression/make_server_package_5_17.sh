#!/usr/bin/env bash
set -euo pipefail

PKG_NAME="ckks_split_server_experiments_5_17"
PKG_ROOT="server_packages/${PKG_NAME}"
ARCHIVE="server_packages/${PKG_NAME}.tar.gz"

rm -rf "${PKG_ROOT}" "${ARCHIVE}"
mkdir -p "${PKG_ROOT}"

copy_module() {
  local src="$1"
  local dst="${PKG_ROOT}/${src}"
  mkdir -p "${dst}"
  rsync -a \
    --exclude 'logs/' \
    --exclude 'results/' \
    --exclude '__pycache__/' \
    --exclude '.DS_Store' \
    --exclude '*.log' \
    --exclude '*.tsv' \
    "${src}/" "${dst}/"
}

copy_module "baseline_top16_Lout1-16实验"
copy_module "5-16-内核subkey"

cp SERVER_PACKAGE_PLAN.md "${PKG_ROOT}/README_SERVER.md"

cat > "${PKG_ROOT}/RUN_COMMANDS.md" <<'EOF'
# Server Run Commands

## One command for overnight run

Run both experiment groups sequentially:

```bash
chmod +x RUN_ALL_SERVER_EXPERIMENTS.sh
nohup bash RUN_ALL_SERVER_EXPERIMENTS.sh > overnight.nohup.log 2>&1 &
```

Check progress:

```bash
tail -f overnight.nohup.log
```

Optional leaf15 subkey exploratory checks:

```bash
RUN_LEAF15=1 nohup bash RUN_ALL_SERVER_EXPERIMENTS.sh > overnight.nohup.log 2>&1 &
```

## 1. Lattigo split level sweeps

```bash
cd baseline_top16_Lout1-16实验
chmod +x run_split_level_sweeps.sh
REPS=1 ./run_split_level_sweeps.sh
```

## 2. Subkey top17 i1 compatibility

```bash
cd 5-16-内核subkey
chmod +x run_subkey_top17_i1_group.sh
REPS=1 ./run_subkey_top17_i1_group.sh
```

Optional leaf15 exploratory checks:

```bash
RUN_LEAF15=1 REPS=1 ./run_subkey_top17_i1_group.sh
```

## Output files

Each script writes terminal output and logs at the same time.

```text
logs/.../master.log
results/.../*.tsv
```
EOF

cat > "${PKG_ROOT}/RUN_ALL_SERVER_EXPERIMENTS.sh" <<'EOF'
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
EOF

chmod +x "${PKG_ROOT}/RUN_ALL_SERVER_EXPERIMENTS.sh"

tar -czf "${ARCHIVE}" -C server_packages "${PKG_NAME}"

echo "Created package:"
echo "${ARCHIVE}"
du -sh "${ARCHIVE}"
