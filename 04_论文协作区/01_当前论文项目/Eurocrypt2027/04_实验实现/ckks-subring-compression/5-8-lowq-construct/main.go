package main

import (
	"flag"
	"fmt"
	"math"
	"math/bits"
	"math/rand"
	"time"
)

// This is a small, self-contained negacyclic RLWE prototype for testing the
// proposed constructive split:
//
//   R_N ciphertext under s(X)
//        -> two R_n module ciphertext rows under (s0(Y), s1(Y))
//        -> two low-q R_n ciphertexts under a fresh small-ring key s'(Y).
//
// It intentionally does not use Lattigo's CKKS stack. The purpose is to isolate
// the algebra and the modulus-drop issue before wiring the construction into the
// full bootstrapping implementation.

type poly []int64

type kskPair struct {
	b []poly
	a []poly
}

type ksk struct {
	base int64
	digs int
	s0   kskPair
	s1   kskPair
}

type ciphertext struct {
	c0 poly
	c1 poly
}

func main() {
	var (
		logN       = flag.Int("logN", 5, "large ring log dimension; N must be even")
		q          = flag.Int64("q", 134217757, "low output modulus q")
		qFactor    = flag.Int64("qFactor", 16, "large modulus is Q=q*qFactor; q must divide Q")
		scale      = flag.Int64("scale", 1<<16, "CKKS-like integer scale")
		msgBound   = flag.Int64("msgBound", 2, "sample message coefficients from [-msgBound,msgBound]")
		encErr     = flag.Int64("encErr", 1, "input encryption error coefficients from [-encErr,encErr]")
		ksErr      = flag.Int64("ksErr", 1, "key-switching key error coefficients from [-ksErr,ksErr]")
		base       = flag.Int64("base", 16, "gadget decomposition base")
		trials     = flag.Int("trials", 5, "number of random trials")
		seed       = flag.Int64("seed", time.Now().UnixNano(), "PRNG seed")
		printFirst = flag.Bool("printFirst", true, "print detailed first trial diagnostics")
	)
	flag.Parse()

	if *logN < 2 {
		panic("logN must be at least 2")
	}
	N := 1 << *logN
	if N%2 != 0 {
		panic("N must be even")
	}
	n := N / 2
	Q := (*q) * (*qFactor)
	if Q%*q != 0 {
		panic("q must divide Q")
	}
	if *base < 2 {
		panic("base must be at least 2")
	}
	checkMulSafety(Q)

	rng := rand.New(rand.NewSource(*seed))
	fmt.Printf("[config] N=%d n=%d Q=%d q=%d Q/q=%d scale=%d msgBound=%d encErr=%d ksErr=%d base=%d trials=%d seed=%d\n",
		N, n, Q, *q, *qFactor, *scale, *msgBound, *encErr, *ksErr, *base, *trials, *seed)
	fmt.Println("[goal] construct low-q small-ring ciphertexts for P0(Y), P1(Y) directly from one big-ring ciphertext encrypting P0(X^2)+X P1(X^2)")

	var sumBigErr, sumModule0, sumModule1, sumLeaf0, sumLeaf1 float64
	var sumBaseLeaf0, sumBaseLeaf1 float64
	var sumLowTime, sumBaseTime time.Duration
	bestPrec0, bestPrec1 := math.Inf(-1), math.Inf(-1)
	worstPrec0, worstPrec1 := math.Inf(1), math.Inf(1)
	bestBasePrec0, bestBasePrec1 := math.Inf(-1), math.Inf(-1)
	worstBasePrec0, worstBasePrec1 := math.Inf(1), math.Inf(1)
	successes := 0

	for trial := 0; trial < *trials; trial++ {
		res := runTrial(rng, N, n, Q, *q, *scale, *msgBound, *encErr, *ksErr, *base)
		if res.ok {
			successes++
		}

		sumBigErr += float64(res.bigErr)
		sumModule0 += float64(res.moduleErr0)
		sumModule1 += float64(res.moduleErr1)
		sumLeaf0 += float64(res.leafErr0)
		sumLeaf1 += float64(res.leafErr1)
		sumBaseLeaf0 += float64(res.baseLeafErr0)
		sumBaseLeaf1 += float64(res.baseLeafErr1)
		sumLowTime += res.lowTime
		sumBaseTime += res.baseTime
		bestPrec0 = math.Max(bestPrec0, res.prec0)
		bestPrec1 = math.Max(bestPrec1, res.prec1)
		worstPrec0 = math.Min(worstPrec0, res.prec0)
		worstPrec1 = math.Min(worstPrec1, res.prec1)
		bestBasePrec0 = math.Max(bestBasePrec0, res.basePrec0)
		bestBasePrec1 = math.Max(bestBasePrec1, res.basePrec1)
		worstBasePrec0 = math.Min(worstBasePrec0, res.basePrec0)
		worstBasePrec1 = math.Min(worstBasePrec1, res.basePrec1)

		if trial == 0 && *printFirst {
			fmt.Println()
			fmt.Println("=== first trial diagnostics ===")
			fmt.Printf("secret hamming: H(s)=%d, H(s0)=%d, H(s1)=%d, H(s')=%d\n", res.hS, res.hS0, res.hS1, res.hSP)
			fmt.Printf("big decrypt max error vs Delta*P:      %d\n", res.bigErr)
			fmt.Printf("module row 0 max error vs Delta*P0:   %d\n", res.moduleErr0)
			fmt.Printf("module row 1 max error vs Delta*P1:   %d\n", res.moduleErr1)
			fmt.Printf("leaf ct 0 max error vs Delta*P0:      %d  precision ~= %.2f bits\n", res.leafErr0, res.prec0)
			fmt.Printf("leaf ct 1 max error vs Delta*P1:      %d  precision ~= %.2f bits\n", res.leafErr1, res.prec1)
			fmt.Printf("baseline keep-Q leaf errors:          row0=%d row1=%d precision ~= %.2f/%.2f bits\n", res.baseLeafErr0, res.baseLeafErr1, res.basePrec0, res.basePrec1)
			fmt.Printf("timing: low-q construct=%s, keep-Q baseline=%s\n", res.lowTime, res.baseTime)
			fmt.Printf("wrap margin q/2 - max(|Delta*P_i + error|): row0=%d row1=%d\n", res.margin0, res.margin1)
			fmt.Printf("semantic check: %v\n", res.ok)
			fmt.Println("sample P0 coefficients:", res.p0Sample)
			fmt.Println("sample Dec(ct0')/scale:", res.dec0Sample)
			fmt.Println("sample P1 coefficients:", res.p1Sample)
			fmt.Println("sample Dec(ct1')/scale:", res.dec1Sample)
		}
	}

	den := float64(*trials)
	fmt.Println()
	fmt.Println("=== aggregate ===")
	fmt.Printf("successes: %d/%d\n", successes, *trials)
	fmt.Printf("avg big decrypt max error:    %.2f\n", sumBigErr/den)
	fmt.Printf("avg module row max error:     row0=%.2f row1=%.2f\n", sumModule0/den, sumModule1/den)
	fmt.Printf("avg leaf ciphertext max error: row0=%.2f row1=%.2f\n", sumLeaf0/den, sumLeaf1/den)
	fmt.Printf("precision range: row0=[%.2f, %.2f] bits row1=[%.2f, %.2f] bits\n", worstPrec0, bestPrec0, worstPrec1, bestPrec1)
	fmt.Printf("avg keep-Q baseline leaf error: row0=%.2f row1=%.2f\n", sumBaseLeaf0/den, sumBaseLeaf1/den)
	fmt.Printf("keep-Q baseline precision range: row0=[%.2f, %.2f] bits row1=[%.2f, %.2f] bits\n", worstBasePrec0, bestBasePrec0, worstBasePrec1, bestBasePrec1)
	fmt.Printf("avg timing: low-q construct=%s keep-Q baseline=%s speed ratio=%.2fx\n",
		sumLowTime/time.Duration(*trials), sumBaseTime/time.Duration(*trials), float64(sumBaseTime)/math.Max(1, float64(sumLowTime)))
	fmt.Printf("modulus reduction: log2(Q)=%.2f -> log2(q)=%.2f, saved %.2f bits, factor %.2fx\n",
		math.Log2(float64(Q)), math.Log2(float64(*q)), math.Log2(float64(Q))-math.Log2(float64(*q)), float64(Q)/float64(*q))
	fmt.Println()
	fmt.Println("Interpretation:")
	fmt.Println("- module row errors matching the big decrypt error mean the extraction equations are correct.")
	fmt.Println("- leaf errors add only module-to-ring key-switching noise; set -ksErr=0 to isolate exact algebra.")
	fmt.Println("- because q divides Q, dropping from Q to q preserves the plaintext scale instead of multiplying it by q/Q.")
}

type trialResult struct {
	bigErr       int64
	moduleErr0   int64
	moduleErr1   int64
	leafErr0     int64
	leafErr1     int64
	baseLeafErr0 int64
	baseLeafErr1 int64
	prec0        float64
	prec1        float64
	basePrec0    float64
	basePrec1    float64
	margin0      int64
	margin1      int64
	lowTime      time.Duration
	baseTime     time.Duration
	ok           bool
	hS           int
	hS0          int
	hS1          int
	hSP          int
	p0Sample     []int64
	p1Sample     []int64
	dec0Sample   []float64
	dec1Sample   []float64
}

func runTrial(rng *rand.Rand, N, n int, Q, q, scale, msgBound, encErr, ksErr, base int64) trialResult {
	sBigSigned := sampleTernary(rng, N)
	s0Signed := evenPartSigned(sBigSigned)
	s1Signed := oddPartSigned(sBigSigned)
	sPrimeSigned := sampleTernary(rng, n)

	p0Signed := sampleSmall(rng, n, msgBound)
	p1Signed := sampleSmall(rng, n, msgBound)
	PBigSigned := make(poly, N)
	for i := 0; i < n; i++ {
		PBigSigned[2*i] = p0Signed[i]
		PBigSigned[2*i+1] = p1Signed[i]
	}

	ct := encryptBig(rng, PBigSigned, sBigSigned, Q, scale, encErr)
	wantBigQ := scalePoly(PBigSigned, Q, scale)
	decBigQ := decrypt(ct.c0, ct.c1, signedToMod(sBigSigned, Q), N, Q)
	bigErr := maxCenteredDiff(decBigQ, wantBigQ, Q)

	startLow := time.Now()
	A := reduceMod(evenPart(ct.c0), q)
	C := reduceMod(oddPart(ct.c0), q)
	B := reduceMod(evenPart(ct.c1), q)
	D := reduceMod(oddPart(ct.c1), q)
	YD := mulByMonomial(D, 1, n, q)

	s0q := signedToMod(s0Signed, q)
	s1q := signedToMod(s1Signed, q)
	spq := signedToMod(sPrimeSigned, q)

	want0q := scalePoly(p0Signed, q, scale)
	want1q := scalePoly(p1Signed, q, scale)

	module0 := addMod(addMod(A, mulNegacyclic(B, s0q, n, q), q), mulNegacyclic(YD, s1q, n, q), q)
	module1 := addMod(addMod(C, mulNegacyclic(D, s0q, n, q), q), mulNegacyclic(B, s1q, n, q), q)
	moduleErr0 := maxCenteredDiff(module0, want0q, q)
	moduleErr1 := maxCenteredDiff(module1, want1q, q)

	evk := genKSK(rng, s0q, s1q, spq, n, q, base, ksErr)
	ct0 := keySwitchRow(A, B, YD, evk, n, q)
	ct1 := keySwitchRow(C, D, B, evk, n, q)
	lowTime := time.Since(startLow)
	dec0 := decrypt(ct0.c0, ct0.c1, spq, n, q)
	dec1 := decrypt(ct1.c0, ct1.c1, spq, n, q)
	leafErr0 := maxCenteredDiff(dec0, want0q, q)
	leafErr1 := maxCenteredDiff(dec1, want1q, q)

	prec0 := precisionBits(scale, leafErr0)
	prec1 := precisionBits(scale, leafErr1)
	margin0 := wrapMargin(dec0, q)
	margin1 := wrapMargin(dec1, q)

	startBase := time.Now()
	AQ := evenPart(ct.c0)
	CQ := oddPart(ct.c0)
	BQ := evenPart(ct.c1)
	DQ := oddPart(ct.c1)
	YDQ := mulByMonomial(DQ, 1, n, Q)
	s0Q := signedToMod(s0Signed, Q)
	s1Q := signedToMod(s1Signed, Q)
	spQ := signedToMod(sPrimeSigned, Q)
	want0Q := scalePoly(p0Signed, Q, scale)
	want1Q := scalePoly(p1Signed, Q, scale)
	baseEvk := genKSK(rng, s0Q, s1Q, spQ, n, Q, base, ksErr)
	baseCt0 := keySwitchRow(AQ, BQ, YDQ, baseEvk, n, Q)
	baseCt1 := keySwitchRow(CQ, DQ, BQ, baseEvk, n, Q)
	baseTime := time.Since(startBase)
	baseDec0 := decrypt(baseCt0.c0, baseCt0.c1, spQ, n, Q)
	baseDec1 := decrypt(baseCt1.c0, baseCt1.c1, spQ, n, Q)
	baseLeafErr0 := maxCenteredDiff(baseDec0, want0Q, Q)
	baseLeafErr1 := maxCenteredDiff(baseDec1, want1Q, Q)
	basePrec0 := precisionBits(scale, baseLeafErr0)
	basePrec1 := precisionBits(scale, baseLeafErr1)

	ok := bigErr <= encErr && moduleErr0 <= bigErr && moduleErr1 <= bigErr && margin0 > 0 && margin1 > 0

	return trialResult{
		bigErr:       bigErr,
		moduleErr0:   moduleErr0,
		moduleErr1:   moduleErr1,
		leafErr0:     leafErr0,
		leafErr1:     leafErr1,
		baseLeafErr0: baseLeafErr0,
		baseLeafErr1: baseLeafErr1,
		prec0:        prec0,
		prec1:        prec1,
		basePrec0:    basePrec0,
		basePrec1:    basePrec1,
		margin0:      margin0,
		margin1:      margin1,
		lowTime:      lowTime,
		baseTime:     baseTime,
		ok:           ok,
		hS:           hamming(sBigSigned),
		hS0:          hamming(s0Signed),
		hS1:          hamming(s1Signed),
		hSP:          hamming(sPrimeSigned),
		p0Sample:     samplePrefixInt(p0Signed, 8),
		p1Sample:     samplePrefixInt(p1Signed, 8),
		dec0Sample:   scaledPrefix(dec0, q, scale, 8),
		dec1Sample:   scaledPrefix(dec1, q, scale, 8),
	}
}

func encryptBig(rng *rand.Rand, msgSigned, sSigned poly, mod, scale, errBound int64) ciphertext {
	N := len(msgSigned)
	c1 := uniformPoly(rng, N, mod)
	sMod := signedToMod(sSigned, mod)
	sc1 := mulNegacyclic(sMod, c1, N, mod)
	m := scalePoly(msgSigned, mod, scale)
	e := signedToMod(sampleSmall(rng, N, errBound), mod)
	c0 := subMod(addMod(m, e, mod), sc1, mod)
	return ciphertext{c0: c0, c1: c1}
}

func genKSK(rng *rand.Rand, s0, s1, sPrime poly, n int, q, base, errBound int64) ksk {
	digs := digitCount(q, base)
	return ksk{
		base: base,
		digs: digs,
		s0:   genKSKPair(rng, s0, sPrime, n, q, base, digs, errBound),
		s1:   genKSKPair(rng, s1, sPrime, n, q, base, digs, errBound),
	}
}

func genKSKPair(rng *rand.Rand, source, target poly, n int, q, base int64, digs int, errBound int64) kskPair {
	out := kskPair{b: make([]poly, digs), a: make([]poly, digs)}
	pow := int64(1)
	for i := 0; i < digs; i++ {
		a := uniformPoly(rng, n, q)
		e := signedToMod(sampleSmall(rng, n, errBound), q)
		msg := scalarMul(source, pow, q)
		ta := mulNegacyclic(target, a, n, q)
		b := subMod(addMod(msg, e, q), ta, q)
		out.a[i] = a
		out.b[i] = b
		pow = (pow * base) % q
	}
	return out
}

func keySwitchRow(t, u, v poly, evk ksk, n int, q int64) ciphertext {
	uDigits := decompose(u, evk.base, evk.digs)
	vDigits := decompose(v, evk.base, evk.digs)
	d0 := clone(t)
	d1 := make(poly, n)
	for i := 0; i < evk.digs; i++ {
		d0 = addMod(d0, mulNegacyclic(uDigits[i], evk.s0.b[i], n, q), q)
		d0 = addMod(d0, mulNegacyclic(vDigits[i], evk.s1.b[i], n, q), q)
		d1 = addMod(d1, mulNegacyclic(uDigits[i], evk.s0.a[i], n, q), q)
		d1 = addMod(d1, mulNegacyclic(vDigits[i], evk.s1.a[i], n, q), q)
	}
	return ciphertext{c0: d0, c1: d1}
}

func decrypt(c0, c1, s poly, N int, mod int64) poly {
	return addMod(c0, mulNegacyclic(s, c1, N, mod), mod)
}

func decompose(p poly, base int64, digs int) []poly {
	out := make([]poly, digs)
	for i := range out {
		out[i] = make(poly, len(p))
	}
	for j, coeff := range p {
		x := coeff
		for i := 0; i < digs; i++ {
			out[i][j] = x % base
			x /= base
		}
	}
	return out
}

func digitCount(q, base int64) int {
	digs := 0
	x := q - 1
	for x > 0 {
		digs++
		x /= base
	}
	return digs
}

func sampleTernary(rng *rand.Rand, n int) poly {
	out := make(poly, n)
	for i := range out {
		out[i] = int64(rng.Intn(3) - 1)
	}
	return out
}

func sampleSmall(rng *rand.Rand, n int, bound int64) poly {
	out := make(poly, n)
	if bound <= 0 {
		return out
	}
	for i := range out {
		out[i] = int64(rng.Int63n(2*bound+1)) - bound
	}
	return out
}

func uniformPoly(rng *rand.Rand, n int, mod int64) poly {
	out := make(poly, n)
	for i := range out {
		out[i] = rng.Int63n(mod)
	}
	return out
}

func signedToMod(p poly, mod int64) poly {
	out := make(poly, len(p))
	for i, x := range p {
		out[i] = modNorm(x, mod)
	}
	return out
}

func scalePoly(p poly, mod, scale int64) poly {
	out := make(poly, len(p))
	for i, x := range p {
		out[i] = modNorm(x*scale, mod)
	}
	return out
}

func reduceMod(p poly, mod int64) poly {
	out := make(poly, len(p))
	for i, x := range p {
		out[i] = modNorm(x, mod)
	}
	return out
}

func evenPart(p poly) poly {
	out := make(poly, len(p)/2)
	for i := range out {
		out[i] = p[2*i]
	}
	return out
}

func oddPart(p poly) poly {
	out := make(poly, len(p)/2)
	for i := range out {
		out[i] = p[2*i+1]
	}
	return out
}

func evenPartSigned(p poly) poly { return evenPart(p) }
func oddPartSigned(p poly) poly  { return oddPart(p) }

func clone(p poly) poly {
	out := make(poly, len(p))
	copy(out, p)
	return out
}

func addMod(a, b poly, mod int64) poly {
	out := make(poly, len(a))
	for i := range a {
		out[i] = modNorm(a[i]+b[i], mod)
	}
	return out
}

func subMod(a, b poly, mod int64) poly {
	out := make(poly, len(a))
	for i := range a {
		out[i] = modNorm(a[i]-b[i], mod)
	}
	return out
}

func scalarMul(a poly, scalar, mod int64) poly {
	out := make(poly, len(a))
	for i := range a {
		out[i] = modNorm(a[i]*scalar, mod)
	}
	return out
}

func mulByMonomial(a poly, shift int, N int, mod int64) poly {
	out := make(poly, N)
	for i, coeff := range a {
		j := i + shift
		sign := int64(1)
		for j >= N {
			j -= N
			sign = -sign
		}
		if sign > 0 {
			out[j] = modNorm(out[j]+coeff, mod)
		} else {
			out[j] = modNorm(out[j]-coeff, mod)
		}
	}
	return out
}

func mulNegacyclic(a, b poly, N int, mod int64) poly {
	out := make(poly, N)
	for i := 0; i < N; i++ {
		ai := a[i]
		if ai == 0 {
			continue
		}
		for j := 0; j < N; j++ {
			bj := b[j]
			if bj == 0 {
				continue
			}
			k := i + j
			prod := (ai * bj) % mod
			if k >= N {
				k -= N
				out[k] = modNorm(out[k]-prod, mod)
			} else {
				out[k] = modNorm(out[k]+prod, mod)
			}
		}
	}
	return out
}

func modNorm(x, mod int64) int64 {
	x %= mod
	if x < 0 {
		x += mod
	}
	return x
}

func center(x, mod int64) int64 {
	x = modNorm(x, mod)
	if x > mod/2 {
		x -= mod
	}
	return x
}

func maxCenteredDiff(got, want poly, mod int64) int64 {
	var max int64
	for i := range got {
		d := abs64(center(got[i]-want[i], mod))
		if d > max {
			max = d
		}
	}
	return max
}

func precisionBits(scale, err int64) float64 {
	if err <= 0 {
		return 60
	}
	return math.Log2(float64(scale) / float64(err))
}

func wrapMargin(p poly, mod int64) int64 {
	limit := mod / 2
	max := int64(0)
	for _, x := range p {
		a := abs64(center(x, mod))
		if a > max {
			max = a
		}
	}
	return limit - max
}

func hamming(p poly) int {
	h := 0
	for _, x := range p {
		if x != 0 {
			h++
		}
	}
	return h
}

func samplePrefixInt(p poly, limit int) []int64 {
	if len(p) < limit {
		limit = len(p)
	}
	out := make([]int64, limit)
	copy(out, p[:limit])
	return out
}

func scaledPrefix(p poly, mod, scale int64, limit int) []float64 {
	if len(p) < limit {
		limit = len(p)
	}
	out := make([]float64, limit)
	for i := 0; i < limit; i++ {
		out[i] = float64(center(p[i], mod)) / float64(scale)
	}
	return out
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func checkMulSafety(mod int64) {
	if mod <= 0 {
		panic("modulus must be positive")
	}
	const maxInt64 = uint64(1<<63 - 1)
	hi, lo := bits.Mul64(uint64(mod-1), uint64(mod-1))
	if hi != 0 || lo > maxInt64 {
		panic("modulus too large for this int64 toy implementation")
	}
}
