# P0-2: 环参数 m, m₀, r 核实

> 基于: 初稿.pdf (src_000031) + GHPS 2013 (src_000028)

---

## 一、初稿中的环记号

初稿 Section 2.1（P5）定义:

```
R_{q,N} = Z_q[X]/(X^N + 1)
N 通常为二的幂
ct = (c_0, c_1) ∈ R_{q,N}^2
N = kn
R_{q,n} = Z_q[Y]/(Y^n + 1)
Y = X^k
```

## 二、与 GHPS 分圆记号对齐

| 当前论文 | GHPS | 说明 |
|---------|------|------|
| R_{q,N} = Z_q[X]/(X^N+1) | R_q = Z_q[ζ_m] (mod q) | power-of-2 negacyclic ring |
| N (power-of-2) | m = 2N | 对应 2N-th cyclotomic |
| R_{q,n} | R₀ = Z_q[ζ_{m₀}] | 子环 |
| n (power-of-2) | m₀ = 2n | 对应 2n-th cyclotomic |
| k = N/n | [GHPS: m/m₀] | 扩张度 |
| — | r = rad(m)/rad(m₀) | 不同素因子比 |

## 三、对应规则

由于当前论文使用 **power-of-two** cyclotomic (X^N+1):

```
X^N + 1 是第 2N-th 分圆多项式
  → m = 2N
  → N = φ(m)/2 = m/2 ✅

同理: m₀ = 2n, n = m₀/2 ✅
```

所以:
```
m = 2N,  m₀ = 2n
m₀ | m  (因为 n | N, 所以 2n | 2N) ✅
```

## 四、r 的计算

```
rad(m) = rad(2N) = rad(2·2^a) = 2  （N是2的幂, 唯一素因子=2）
rad(m₀) = rad(2n) = 2

r = rad(m) / rad(m₀) = 2/2 = 1
```

## 五、与 Geelen 的对齐

Geelen (src_000030) 使用:
```
R = Z[X]/(X^N+1), N = m/2
R' = Z[Y]/(Y^n+1) (subring), Y = X^{N/n}
```

这与当前论文的 R_{q,N}/R_{q,n} 完全一致。

## 六、结论

| 问题 | 答案 |
|------|------|
| m的值 | m = 2N (power-of-2 → 2N-th cyclotomic) |
| m₀的值 | m₀ = 2n |
| m₀ \| m ? | ✅ 是（2n \| 2N 等价于 n \| N） |
| rad(m) | 2 |
| rad(m₀) | 2 |
| r = rad(m)/rad(m₀) | **r = 1** |
| Field switching 方向 | K = Q(ζ_{2N}) → K₀ = Q(ζ_{2n}) |
| 与 GHPS 对齐 | ✅ 可通过 m=2N, m₀=2n 建立对应 |
| 状态 | ✅ 已核实 (high confidence) |

## 七、未确认项

- 初稿中未显式定义 K = Q(ζ_{2N})（数域视角）
- 如果论文使用 "2N-th cyclotomic" 优于 "ring Z[X]/(X^N+1)" 的表述，需在 Preliminaries 中说明
