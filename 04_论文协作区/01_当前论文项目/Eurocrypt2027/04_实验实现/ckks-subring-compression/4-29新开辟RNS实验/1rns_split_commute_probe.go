// rns_split_commute_probe.go
//
// A small, self-contained Go probe for the algebraic/RNS scheduling question:
//
//	Does coefficient/module Split_k commute with coefficient-wise RNS operations?
//
// This deliberately does NOT use Lattigo. It tests the necessary algebraic facts
// behind an RNS-aware subring/factored bootstrapping schedule:
//
//  1. Split_k(ModUp_N(a)) == (oplus ModUp_n)(Split_k(a))
//  2. Split_k(DropRNSPrimes_N(a)) == (oplus DropRNSPrimes_n)(Split_k(a))
//  3. Split_k(BasisConvert_N(a)) == (oplus BasisConvert_n)(Split_k(a))
//  4. Split_k(RoundScale_N(a)) == (oplus RoundScale_n)(Split_k(a))
//
// These equalities hold because the tested RNS operations are coefficient-wise,
// while Split_k is just residue-class coefficient/module decomposition:
//
//	P(X) = sum_{r=0}^{k-1} X^r P_r(X^k).
//
// This file is NOT a cryptographic security test and NOT a CKKS bootstrapping
// correctness test. It is a minimal sanity check for whether moving ModUp /
// basis conversion / modulus reduction across Split_k is algebraically legal.
//
// Example:
//
//	go run rns_split_commute_probe.go -N=64 -k=4 -trials=1000
package main

import (
	"flag"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"strings"
)

type RNS [][]uint64 // RNS[limb][coeff]

type LeafRNS []RNS // LeafRNS[leaf][limb][coeff]

func must(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

func parsePrimes(s string) []uint64 {
	parts := strings.Split(s, ",")
	out := make([]uint64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		z := new(big.Int)
		if _, ok := z.SetString(p, 10); !ok {
			panic("bad prime/modulus: " + p)
		}
		if !z.IsUint64() {
			panic("modulus does not fit uint64: " + p)
		}
		out = append(out, z.Uint64())
	}
	must(len(out) > 0, "empty modulus list")
	return out
}

func modInt64(x int64, q uint64) uint64 {
	qq := int64(q)
	r := x % qq
	if r < 0 {
		r += qq
	}
	return uint64(r)
}

func modBig(x *big.Int, q uint64) uint64 {
	m := new(big.Int).SetUint64(q)
	r := new(big.Int).Mod(x, m)
	return r.Uint64()
}

func product(moduli []uint64) *big.Int {
	out := big.NewInt(1)
	for _, q := range moduli {
		out.Mul(out, new(big.Int).SetUint64(q))
	}
	return out
}

func randomCenteredCoeffs(N int, bound int64, rng *rand.Rand) []*big.Int {
	out := make([]*big.Int, N)
	span := 2*bound + 1
	for i := 0; i < N; i++ {
		x := rng.Int63n(span) - bound
		out[i] = big.NewInt(x)
	}
	return out
}

func modUp(coeffs []*big.Int, moduli []uint64) RNS {
	out := make(RNS, len(moduli))
	for l, q := range moduli {
		out[l] = make([]uint64, len(coeffs))
		for i, x := range coeffs {
			out[l][i] = modBig(x, q)
		}
	}
	return out
}

func splitCoeffs(coeffs []*big.Int, k int) [][]*big.Int {
	N := len(coeffs)
	must(N%k == 0, "N must be divisible by k")
	n := N / k
	leaves := make([][]*big.Int, k)
	for r := 0; r < k; r++ {
		leaves[r] = make([]*big.Int, n)
		for j := 0; j < n; j++ {
			leaves[r][j] = new(big.Int).Set(coeffs[r+j*k])
		}
	}
	return leaves
}

func splitRNS(x RNS, k int) LeafRNS {
	must(len(x) > 0, "empty RNS")
	N := len(x[0])
	must(N%k == 0, "N must be divisible by k")
	n := N / k
	L := len(x)
	leaves := make(LeafRNS, k)
	for r := 0; r < k; r++ {
		leaves[r] = make(RNS, L)
		for l := 0; l < L; l++ {
			must(len(x[l]) == N, "inconsistent RNS limb length")
			leaves[r][l] = make([]uint64, n)
			for j := 0; j < n; j++ {
				leaves[r][l][j] = x[l][r+j*k]
			}
		}
	}
	return leaves
}

func modUpLeaves(leaves [][]*big.Int, moduli []uint64) LeafRNS {
	out := make(LeafRNS, len(leaves))
	for r := range leaves {
		out[r] = modUp(leaves[r], moduli)
	}
	return out
}

func dropLimbs(x RNS, keep int) RNS {
	must(keep >= 0 && keep <= len(x), "bad keep")
	out := make(RNS, keep)
	for l := 0; l < keep; l++ {
		out[l] = append([]uint64(nil), x[l]...)
	}
	return out
}

func dropLimbsLeaves(x LeafRNS, keep int) LeafRNS {
	out := make(LeafRNS, len(x))
	for r := range x {
		out[r] = dropLimbs(x[r], keep)
	}
	return out
}

func equalLeafRNS(a, b LeafRNS) (bool, string) {
	if len(a) != len(b) {
		return false, fmt.Sprintf("leaf count mismatch: %d vs %d", len(a), len(b))
	}
	for r := range a {
		if len(a[r]) != len(b[r]) {
			return false, fmt.Sprintf("leaf %d limb count mismatch: %d vs %d", r, len(a[r]), len(b[r]))
		}
		for l := range a[r] {
			if len(a[r][l]) != len(b[r][l]) {
				return false, fmt.Sprintf("leaf %d limb %d coeff count mismatch", r, l)
			}
			for j := range a[r][l] {
				if a[r][l][j] != b[r][l][j] {
					return false, fmt.Sprintf("mismatch leaf=%d limb=%d coeff=%d: %d vs %d", r, l, j, a[r][l][j], b[r][l][j])
				}
			}
		}
	}
	return true, ""
}

func crtReconstructCentered(residues []uint64, moduli []uint64) *big.Int {
	must(len(residues) == len(moduli), "residue/moduli length mismatch")
	M := product(moduli)
	x := big.NewInt(0)

	for i, qi := range moduli {
		q := new(big.Int).SetUint64(qi)
		Mi := new(big.Int).Div(new(big.Int).Set(M), q)
		inv := new(big.Int).ModInverse(Mi, q)
		if inv == nil {
			panic("moduli are not pairwise coprime or inverse does not exist")
		}
		term := new(big.Int).SetUint64(residues[i])
		term.Mul(term, Mi)
		term.Mul(term, inv)
		x.Add(x, term)
	}
	x.Mod(x, M)

	half := new(big.Int).Rsh(new(big.Int).Set(M), 1)
	if x.Cmp(half) > 0 {
		x.Sub(x, M)
	}
	return x
}

func basisConvert(x RNS, srcModuli, dstModuli []uint64) RNS {
	must(len(x) == len(srcModuli), "source RNS limb count does not match source moduli")
	N := len(x[0])
	out := make(RNS, len(dstModuli))
	for l := range out {
		out[l] = make([]uint64, N)
	}

	residues := make([]uint64, len(srcModuli))
	for i := 0; i < N; i++ {
		for l := range srcModuli {
			residues[l] = x[l][i]
		}
		centered := crtReconstructCentered(residues, srcModuli)
		for l, q := range dstModuli {
			out[l][i] = modBig(centered, q)
		}
	}
	return out
}

func basisConvertLeaves(x LeafRNS, srcModuli, dstModuli []uint64) LeafRNS {
	out := make(LeafRNS, len(x))
	for r := range x {
		out[r] = basisConvert(x[r], srcModuli, dstModuli)
	}
	return out
}

func roundDivBig(x *big.Int, divisor int64) *big.Int {
	must(divisor > 0, "divisor must be positive")
	d := big.NewInt(divisor)
	half := big.NewInt(divisor / 2)

	if x.Sign() >= 0 {
		y := new(big.Int).Add(x, half)
		return y.Div(y, d)
	}

	y := new(big.Int).Neg(x)
	y.Add(y, half)
	y.Div(y, d)
	y.Neg(y)
	return y
}

func roundScaleCoeffs(coeffs []*big.Int, divisor int64) []*big.Int {
	out := make([]*big.Int, len(coeffs))
	for i, x := range coeffs {
		out[i] = roundDivBig(x, divisor)
	}
	return out
}

func runTrial(N, k int, coeffBound int64, srcModuli, dstModuli []uint64, keep int, roundDiv int64, rng *rand.Rand) error {
	coeffs := randomCenteredCoeffs(N, coeffBound, rng)

	// Test 1: Split(ModUp_N(a)) == ModUp_leaves(Split(a)).
	topModUp := modUp(coeffs, srcModuli)
	left1 := splitRNS(topModUp, k)
	right1 := modUpLeaves(splitCoeffs(coeffs, k), srcModuli)
	if ok, msg := equalLeafRNS(left1, right1); !ok {
		return fmt.Errorf("ModUp/Split commute failed: %s", msg)
	}

	// Test 2: Dropping RNS limbs commutes with Split.
	left2 := splitRNS(dropLimbs(topModUp, keep), k)
	right2 := dropLimbsLeaves(splitRNS(topModUp, k), keep)
	if ok, msg := equalLeafRNS(left2, right2); !ok {
		return fmt.Errorf("DropRNSPrimes/Split commute failed: %s", msg)
	}

	// Test 3: Basis conversion / modulus extension-restriction commutes with Split.
	topConverted := basisConvert(topModUp, srcModuli, dstModuli)
	left3 := splitRNS(topConverted, k)
	right3 := basisConvertLeaves(splitRNS(topModUp, k), srcModuli, dstModuli)
	if ok, msg := equalLeafRNS(left3, right3); !ok {
		return fmt.Errorf("BasisConvert/Split commute failed: %s", msg)
	}

	// Test 4: Coefficient-wise rounding/rescale commutes with Split.
	topRounded := modUp(roundScaleCoeffs(coeffs, roundDiv), dstModuli)
	left4 := splitRNS(topRounded, k)
	leafRounded := make([][]*big.Int, k)
	for r, leaf := range splitCoeffs(coeffs, k) {
		leafRounded[r] = roundScaleCoeffs(leaf, roundDiv)
	}
	right4 := modUpLeaves(leafRounded, dstModuli)
	if ok, msg := equalLeafRNS(left4, right4); !ok {
		return fmt.Errorf("RoundScale/Split commute failed: %s", msg)
	}

	return nil
}

func main() {
	N := flag.Int("N", 64, "top coefficient length / ring degree")
	k := flag.Int("k", 4, "split factor; must divide N")
	trials := flag.Int("trials", 100, "random trials")
	seed := flag.Int64("seed", 1, "PRNG seed")
	bound := flag.Int64("bound", 1<<20, "absolute bound for random centered coefficients")
	src := flag.String("src", "2305843009211596801,2305843009211072513,2305843009208975361", "comma-separated source RNS moduli")
	dst := flag.String("dst", "2305843009206880257,2305843009204783105", "comma-separated destination RNS moduli")
	keep := flag.Int("keep", 2, "number of source limbs to keep in DropRNSPrimes test")
	roundDiv := flag.Int64("roundDiv", 257, "divisor for coefficient-wise rounding/rescale test")
	flag.Parse()

	if *N <= 0 || *k <= 0 || *N%*k != 0 {
		fmt.Fprintf(os.Stderr, "invalid N,k: N=%d k=%d; require k>0 and k|N\n", *N, *k)
		os.Exit(2)
	}

	srcModuli := parsePrimes(*src)
	dstModuli := parsePrimes(*dst)
	if *keep < 0 || *keep > len(srcModuli) {
		fmt.Fprintf(os.Stderr, "invalid keep=%d; source limb count=%d\n", *keep, len(srcModuli))
		os.Exit(2)
	}

	rng := rand.New(rand.NewSource(*seed))

	fmt.Println("=== RNS / coefficient-split commutation probe ===")
	fmt.Printf("N=%d k=%d n=%d trials=%d seed=%d\n", *N, *k, *N / *k, *trials, *seed)
	fmt.Printf("source moduli count=%d: %v\n", len(srcModuli), srcModuli)
	fmt.Printf("target moduli count=%d: %v\n", len(dstModuli), dstModuli)
	fmt.Printf("drop keep=%d, roundDiv=%d, coeffBound=%d\n", *keep, *roundDiv, *bound)
	fmt.Println()
	fmt.Println("Testing:")
	fmt.Println("  [1] Split_k o ModUp_N            == (oplus ModUp_n) o Split_k")
	fmt.Println("  [2] Split_k o DropRNSPrimes_N    == (oplus DropRNSPrimes_n) o Split_k")
	fmt.Println("  [3] Split_k o BasisConvert_N     == (oplus BasisConvert_n) o Split_k")
	fmt.Println("  [4] Split_k o RoundScale_N       == (oplus RoundScale_n) o Split_k")
	fmt.Println()

	for t := 0; t < *trials; t++ {
		if err := runTrial(*N, *k, *bound, srcModuli, dstModuli, *keep, *roundDiv, rng); err != nil {
			fmt.Printf("FAIL at trial %d: %v\n", t, err)
			os.Exit(1)
		}
	}

	fmt.Println("PASS: all tested coefficient-wise RNS operations commute with coefficient/module Split_k.")
	fmt.Println()
	fmt.Println("Interpretation:")
	fmt.Println("  This supports the algebraic part of moving RNS ModUp / basis conversion /")
	fmt.Println("  modulus-chain restriction across HERMES-style coefficient decomposition.")
	fmt.Println("  It does NOT prove CKKS bootstrapping correctness under reduced Q; that still")
	fmt.Println("  requires a ciphertext-level experiment measuring KS noise, EvalMod range,")
	fmt.Println("  scale conventions, output precision, and security for the reduced QP.")
}
