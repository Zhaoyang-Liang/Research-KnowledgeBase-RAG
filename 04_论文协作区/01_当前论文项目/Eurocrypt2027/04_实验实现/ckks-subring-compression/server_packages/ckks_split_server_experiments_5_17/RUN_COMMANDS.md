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
