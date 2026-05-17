# 5-17 Server Experiment Package Plan

This package contains two independent offline Go modules. Upload both to the
server and run the scripts from inside each subdirectory.

## Experiment Group 1: Lattigo Split Level Sweeps

Directory:

```text
baseline_top16_Lout1-16实验/
```

Script:

```bash
cd baseline_top16_Lout1-16实验
chmod +x run_split_level_sweeps.sh
REPS=1 ./run_split_level_sweeps.sh
```

Runs:

```text
top16 -> leaf15, Lout = 1..3
top17 -> leaf15, Lout = 1..3
top17 -> leaf16, Lout = 1..15
```

Outputs:

```text
logs/split_level_sweeps_YYYYMMDD_HHMMSS/master.log
results/split_level_sweeps_YYYYMMDD_HHMMSS/split_level_sweeps.tsv
```

## Experiment Group 2: Subkey Kernel Compatibility

Directory:

```text
5-16-内核subkey/
```

Main script:

```bash
cd 5-16-内核subkey
chmod +x run_subkey_top17_i1_group.sh
REPS=1 ./run_subkey_top17_i1_group.sh
```

Runs:

```text
top17 direct subkey i1, no split, Lout = 15
top17 -> leaf16 split + subkey i1, Lout = 15
```

Optional exploratory leaf15 subkey attempts:

```bash
RUN_LEAF15=1 REPS=1 ./run_subkey_top17_i1_group.sh
```

Outputs:

```text
logs/subkey_top17_i1_group_YYYYMMDD_HHMMSS/master.log
results/subkey_top17_i1_group_YYYYMMDD_HHMMSS/subkey_top17_i1_group.tsv
```

## Upload Package

From the parent directory:

```bash
chmod +x make_server_package_5_17.sh
./make_server_package_5_17.sh
```

This creates:

```text
server_packages/ckks_split_server_experiments_5_17.tar.gz
```

The package excludes local `logs/`, `results/`, and Go build caches, but keeps:

```text
go.mod
go.sum
vendor/
deps/subring_bts-main/
main.go
anyu25.go
run_*.sh
README.md
scripts/
```

The server does not need GitHub access because each module uses:

```text
replace github.com/tuneinsight/lattigo/v6 => ./deps/subring_bts-main
GOFLAGS=-mod=vendor
```

Use Go 1.25.x if possible. Go 1.18.1 is too old for the vendored dependencies.
