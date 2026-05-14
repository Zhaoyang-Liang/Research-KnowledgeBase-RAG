logQP 的差距（约 95 bits）是因为他们把 StC/EvalMod/CtS 的 dedicated primes 全部算在了 logQP 里，而 Lattigo 默认 bootstrapping 的 chain 不同。这个对比目前是合理的，等实验跑出来看结果再决定要不要进一步对齐 prime chain。

==

mac@budongjishubu:~/Desktop/Research-KB-RAG/04_论文协作区/01_当前论文项目/Eurocrypt2027/04_实验实现/ckks-subring-compression % go run ./4-29/. -logN=16 -layers=1 -run=all -paramSet=anyu25 -reps=1
[paramSet=anyu25] Xs = uniform ternary P=2/3 (H≈43690 for N=2^16)
Bootstrapping parameters: logN=16, logSlots=15, H(43691; 0), ExtDegree=1, sigma={3.2 19.2}, logQP=1635.000045, levels=31, scale=2^35

Generating bootstrapping evaluation keys...
Generating direct top bootstrapping keys for LogN=16, bootLogQP=1635.0, maxLevel=30...
Generating no-B_k EXACT-BASIS leaf bootstrapping keys for LogN=15, logQP=1635.0, maxLevel=30, workers=2...
Done

Precision of values vs. ciphertext

Level: 0 (logQ = 46)
Scale: 2^35.000000
ValuesTest: (-1.00000006710065-0.31250027115065i) (-0.87499998113434-0.24999997914564i) (-0.74999980167898-0.18749893551086i) (-0.62500035804813-0.12500036809500i)...
ValuesWant: (-1.00000000000000-0.31250000000000i) (-0.87500000000000-0.25000000000000i) (-0.75000000000000-0.18750000000000i) (-0.62500000000000-0.12500000000000i)...
L1 Err = 22.076606

┌─────────┬───────┬───────┬───────┐
│    Log2 │ REAL  │ IMAG  │ L2    │
├─────────┼───────┼───────┼───────┤
│MIN Prec │ 18.70 │ 18.55 │ 18.20 │
│MAX Prec │ 35.00 │ 35.00 │ 35.00 │
│AVG Prec │ 22.90 │ 22.91 │ 22.40 │
│MED Prec │ 22.61 │ 22.61 │ 22.11 │
│STD Prec │  1.84 │  1.85 │  1.84 │
├─────────┼───────┼───────┼───────┤
│MIN Err  │  0.00 │  0.00 │  0.00 │
│MAX Err  │ 16.30 │ 16.45 │ 16.80 │
│AVG Err  │ 12.10 │ 12.09 │ 12.60 │
│MED Err  │ 12.39 │ 12.39 │ 12.89 │
│STD Err  │  1.84 │  1.85 │  1.84 │
└─────────┴───────┴───────┴───────┘



=== Baseline: direct single-ring bootstrapping (k=1) ===
Bootstrapping...
Bootstrapping finished in 1m6.122002209s
StC: 15.558303459s, ModUp: 1.607194042s, CtS: 43.964068458s, EvalMod: 4.99243625s
Done
Output ciphertext level = 15, circuit moduli = [35 35 35 35 35 35 35 35 35 35 35 35 35 35 35]

Precision of ciphertext vs. Bootstrap(ciphertext)

Level: 15 (logQ = 571)
Scale: 2^35.000000
ValuesTest: (-1.00000021888171-0.31250029474607i) (-0.87500007770224-0.25000006852427i) (-0.74999855873272-0.18749865643960i) (-0.62500011157222-0.12499992002579i)...
ValuesWant: (-1.00000006710065-0.31250027115065i) (-0.87499998113434-0.24999997914564i) (-0.74999980167898-0.18749893551086i) (-0.62500035804813-0.12500036809500i)...
L1 Err = 21.962239

┌─────────┬───────┬───────┬───────┐
│    Log2 │ REAL  │ IMAG  │ L2    │
├─────────┼───────┼───────┼───────┤
│MIN Prec │ 18.45 │ 18.60 │ 17.95 │
│MAX Prec │ 35.00 │ 35.00 │ 35.00 │
│AVG Prec │ 22.70 │ 22.69 │ 22.20 │
│MED Prec │ 22.42 │ 22.42 │ 21.92 │
│STD Prec │  1.74 │  1.72 │  1.74 │
├─────────┼───────┼───────┼───────┤
│MIN Err  │  0.00 │  0.00 │  0.00 │
│MAX Err  │ 16.55 │ 16.40 │ 17.05 │
│AVG Err  │ 12.30 │ 12.31 │ 12.80 │
│MED Err  │ 12.58 │ 12.58 │ 13.08 │
│STD Err  │  1.74 │  1.72 │  1.74 │
└─────────┴───────┴───────┴───────┘


Avg L1 error is 21.962239
Avg time... ModUp: 1.607194042s, CtS: 43.964068458s, EvalMod: 4.99243625s, StC: 15.558303459s, Total: 1m6.122002209s

=== Proposed: factored no-B_k bootstrapping (k=2) ===
Bootstrapping...
Bootstrapping finished in 35.088239458s
StC: 8.032731291s, ModUp: 744.608042ms, CtS: 23.828084625s, EvalMod: 2.4828155s
Done
Output ciphertext level = 15, circuit moduli = [35 35 35 35 35 35 35 35 35 35 35 35 35 35 35]

Precision of ciphertext vs. Bootstrap(ciphertext)

Level: 15 (logQ = 571)
Scale: 2^35.000000
ValuesTest: (-1.00000032031284-0.31250028435899i) (-0.87499983210603-0.24999973976565i) (-0.74999961139382-0.18749918210324i) (-0.62500027298352-0.12500040136705i)...
ValuesWant: (-1.00000006710065-0.31250027115065i) (-0.87499998113434-0.24999997914564i) (-0.74999980167898-0.18749893551086i) (-0.62500035804813-0.12500036809500i)...
L1 Err = 22.644398

┌─────────┬───────┬───────┬───────┐
│    Log2 │ REAL  │ IMAG  │ L2    │
├─────────┼───────┼───────┼───────┤
│MIN Prec │ 19.48 │ 19.52 │ 18.98 │
│MAX Prec │ 35.00 │ 35.00 │ 35.00 │
│AVG Prec │ 23.33 │ 23.31 │ 22.83 │
│MED Prec │ 23.02 │ 23.01 │ 22.52 │
│STD Prec │  1.69 │  1.67 │  1.69 │
├─────────┼───────┼───────┼───────┤
│MIN Err  │  0.00 │  0.00 │  0.00 │
│MAX Err  │ 15.52 │ 15.48 │ 16.02 │
│AVG Err  │ 11.67 │ 11.69 │ 12.17 │
│MED Err  │ 11.98 │ 11.99 │ 12.48 │
│STD Err  │  1.69 │  1.67 │  1.69 │
└─────────┴───────┴───────┴───────┘


Avg L1 error is 22.644398
Avg time... ModUp: 744.608042ms, CtS: 23.828084625s, EvalMod: 2.4828155s, StC: 8.032731291s, Total: 35.088239458s

=== subKey25: AnyuWang 2025 subring-secret bootstrapping (I1 params) ===
[subKey25] running: main -eps=-16 -scale=35 -degree=127 -r=3 -scalestc=33 -scalects=52 -repeat=1

old K is 101
Failure prob is 2^-15.820882 for K = 102
{35 [45] [33 33 33] [52 52 52 52] [61 61 61 61 61] 60 [1 1 1] 8 102 127 3 2 8 16 {0.6666666666666666 0 0} false 0}
Number of circuit primes = 15
Primes length: [45 35 35 35 35 35 35 35 35 35 35 35 35 35 35 35 33 33 33 60 60 60 60 60 60 60 60 60 60 52 52 52]
Bootstrapping parameters: logN=16, logSlots=15, H(43691; 0), ExtDegree=16, sigma={3.200000 19.200000 %!f(int=0)}, logQP=1730.000214, levels=32, scale=2^35

Generating bootstrapping evaluation keys...
Done

Precision of values vs. ciphertext

Level: 0 (logQ = 46)
Scale: 2^35.000000
ValuesTest: (0.18541436740159+0.73795896646299i) (-0.56154232131226-0.23212438975543i) (-0.57477936991010+0.96315652332938i) (-0.41588939389989+0.35419215656848i)...
ValuesWant: (0.18541464762849+0.73795865957778i) (-0.56154227218289-0.23212415992769i) (-0.57477944047169+0.96315660259091i) (-0.41588932077898+0.35419201977790i)...
L1 Err = 21.583774

┌─────────┬───────┬───────┬───────┐
│    Log2 │ REAL  │ IMAG  │ L2    │
├─────────┼───────┼───────┼───────┤
│MIN Prec │ 18.68 │ 18.56 │ 18.18 │
│MAX Prec │ 35.00 │ 35.00 │ 35.00 │
│AVG Prec │ 22.92 │ 22.92 │ 22.42 │
│MED Prec │ 22.60 │ 22.61 │ 22.10 │
│STD Prec │  1.86 │  1.86 │  1.86 │
├─────────┼───────┼───────┼───────┤
│MIN Err  │  0.00 │  0.00 │  0.00 │
│MAX Err  │ 16.32 │ 16.44 │ 16.82 │
│AVG Err  │ 12.08 │ 12.08 │ 12.58 │
│MED Err  │ 12.40 │ 12.39 │ 12.90 │
│STD Err  │  1.86 │  1.86 │  1.86 │
└─────────┴───────┴───────┴───────┘


Bootstrapping...
Subring CtS time: 239.428ms
Bootstrapping finished in 36.002699875s
StC: 11.958110125s, ModUp: 120.548084ms, CtS: 14.995803083s, EvalMod: 8.928238583s
Done
Output ciphertext level = 15, circuit moduli = [35 35 35 35 35 35 35 35 35 35 35 35 35 35 35]

Precision of ciphertext vs. Bootstrap(ciphertext)

Level: 15 (logQ = 571)
Scale: 2^35.000000
ValuesTest: (0.18541607670605+0.73795989496107i) (-0.56154179531788-0.23212355908978i) (-0.57477720501694+0.96315795750848i) (-0.41588793969530+0.35419253081850i)...
ValuesWant: (0.18541436740159+0.73795896646299i) (-0.56154232131226-0.23212438975543i) (-0.57477936991010+0.96315652332938i) (-0.41588939389989+0.35419215656848i)...
L1 Err = 19.489173

┌─────────┬───────┬───────┬───────┐
│    Log2 │ REAL  │ IMAG  │ L2    │
├─────────┼───────┼───────┼───────┤
│MIN Prec │ 17.44 │ 17.52 │ 16.94 │
│MAX Prec │ 33.39 │ 34.16 │ 32.89 │
│AVG Prec │ 20.58 │ 20.58 │ 20.08 │
│MED Prec │ 20.24 │ 20.22 │ 19.74 │
│STD Prec │  1.59 │  1.61 │  1.59 │
├─────────┼───────┼───────┼───────┤
│MIN Err  │  1.61 │  0.84 │  2.11 │
│MAX Err  │ 17.56 │ 17.48 │ 18.06 │
│AVG Err  │ 14.42 │ 14.42 │ 14.92 │
│MED Err  │ 14.76 │ 14.78 │ 15.26 │
│STD Err  │  1.59 │  1.61 │  1.59 │
└─────────┴───────┴───────┴───────┘


Avg L1 error is 19.489173
Avg time... ModUp: 120.548084ms, CtS: 14.995803083s, EvalMod: 8.928238583s, StC: 11.958110125s, Total: 36.002699875s

=== Summary ===
Baseline  (k=1)      Total: 1m6.122002209s, Avg L1 prec: 21.96 bits
Proposed  (k=2)      Total: 35.088239458s, Avg L1 prec: 22.64 bits
subKey25  (AnyuWang) Total: 36.002699875s, Avg L1 prec: 19.49 bits
Speedup baseline→proposed: 1.88x
Speedup subKey25→proposed: 1.03x