// named/pi_gauss_bb_test.go v3
package named_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/egp/gosper-gcf/core"
	"github.com/egp/gosper-gcf/named"
)

func TestBB_Named_Pi_SourcePQPrefix50(t *testing.T) {
	src := named.Pi()
	for i := 0; i < 50; i++ {
		wantP, wantQ := expectedPiGaussPQ(i)

		term, tail, status := src.NextPQ()
		if status != core.StatusOK {
			t.Fatalf("term %d status = %v, want %v", i+1, status, core.StatusOK)
		}
		if term.P.Cmp(big.NewInt(wantP)) != 0 {
			t.Fatalf("term %d P = %v, want %d", i+1, term.P, wantP)
		}
		if term.Q.Cmp(big.NewInt(wantQ)) != 0 {
			t.Fatalf("term %d Q = %v, want %d", i+1, term.Q, wantQ)
		}
		src = tail
	}
}

func TestBB_Named_Pi_SourceRangeLookahead50(t *testing.T) {
	src := named.Pi()
	for i := 0; i < 50; i++ {
		want := expectedPiGaussRange(i)
		assertPiRange(t, src.Range(), want, i+1)

		_, tail, status := src.NextPQ()
		if status != core.StatusOK {
			t.Fatalf("term %d status = %v, want %v", i+1, status, core.StatusOK)
		}
		src = tail
	}
}

func TestBB_Named_Pi_RCFPrefixMatchesOEISPrefix98(t *testing.T) {
	g := core.NewGCF1(identityUnaryCoeffsPiGauss(), named.Pi())
	assertRCFPrefixPiGauss(t, g, oeisPiRCFTerms98())
}

func TestBB_Named_Pi_Take8MatchesOEISFinitePrefix(t *testing.T) {
	terms := oeisPiRCFTerms98()[:8]
	suffixes := suffixRationalsFromRCFTermsPiBB(terms)

	g := core.NewGCF1(identityUnaryCoeffsPiGauss(), named.Pi())
	src := g.Take(len(terms))

	for i, wantTerm := range terms {
		_, tail := assertFinitePiStep(
			t,
			src,
			i+1,
			wantTerm,
			1,
			suffixes[i],
		)
		src = tail
	}

	assertFinitePiEOF(t, src, len(terms)+1)
}

func TestBB_Named_Pi_RationalMatchesOEISConvergents(t *testing.T) {
	terms := oeisPiRCFTerms98()
	depths := []int{1, 2, 4, 8, 12, 20}

	for _, depth := range depths {
		g := core.NewGCF1(identityUnaryCoeffsPiGauss(), named.Pi())
		got := g.Rational(depth)
		want := convergentFromRCFTermsPiBB(terms[:depth])

		if got.Cmp(want) != 0 {
			t.Fatalf(
				"Rational(%d) = %v/%v, want %v/%v",
				depth,
				got.Num(), got.Den(),
				want.Num(), want.Den(),
			)
		}
	}
}

func identityUnaryCoeffsPiGauss() core.BLFTCoefficients {
	return core.BLFTCoefficients{
		A: big.NewInt(0),
		B: big.NewInt(1),
		C: big.NewInt(0),
		D: big.NewInt(0),
		E: big.NewInt(0),
		F: big.NewInt(0),
		G: big.NewInt(0),
		H: big.NewInt(1),
	}
}

func expectedPiGaussPQ(index int) (int64, int64) {
	if index == 0 {
		return 0, 4
	}
	n := int64(index)
	return 2*n - 1, n * n
}

func expectedPiGaussRange(index int) core.Range {
	if index == 0 {
		return core.Range{
			Lo:     core.Endpoint{Value: core.RationalFromInt64(0), Open: true},
			Hi:     core.Endpoint{Value: core.RationalFromInt64(4), Open: true},
			Inside: true,
		}
	}

	n := int64(index)
	p := 2*n - 1
	q := n * n
	nextOdd := 2*n + 1

	return core.Range{
		Lo: core.Endpoint{
			Value: core.RationalFromInt64(p),
			Open:  true,
		},
		Hi: core.Endpoint{
			Value: core.NewRational(big.NewInt(p*nextOdd+q), big.NewInt(nextOdd)),
			Open:  true,
		},
		Inside: true,
	}
}

func oeisPiRCFTerms98() []int64 {
	return []int64{
		3, 7, 15, 1, 292, 1, 1, 1, 2, 1, 3, 1, 14, 2, 1, 1, 2, 2, 2, 2,
		1, 84, 2, 1, 1, 15, 3, 13, 1, 4, 2, 6, 6, 99, 1, 2, 2, 6, 3, 5,
		1, 1, 6, 8, 1, 7, 1, 2, 3, 7, 1, 2, 1, 1, 12, 1, 1, 1, 3, 1,
		1, 8, 1, 1, 2, 1, 6, 1, 1, 5, 2, 2, 3, 1, 2, 4, 4, 16, 1, 161,
		45, 1, 22, 1, 2, 2, 1, 4, 1, 2, 24, 1, 2, 1, 3, 1, 2, 1,
	}
}

func suffixRationalsFromRCFTermsPiBB(terms []int64) []core.Rational {
	n := len(terms)
	out := make([]core.Rational, n)

	current := core.NewRational(big.NewInt(terms[n-1]), big.NewInt(1))
	out[n-1] = current

	for i := n - 2; i >= 0; i-- {
		a := big.NewInt(terms[i])

		num := new(big.Int).Mul(a, current.Num())
		num.Add(num, current.Den())

		den := new(big.Int).Set(current.Num())

		current = core.NewRational(num, den)
		out[i] = current
	}

	return out
}

func convergentFromRCFTermsPiBB(terms []int64) core.Rational {
	suffixes := suffixRationalsFromRCFTermsPiBB(terms)
	return suffixes[0]
}

func assertRCFPrefixPiGauss(t *testing.T, g *core.GCF, want []int64) {
	t.Helper()

	for i, w := range want {
		term, status := nextRCFWithTimeoutPiGauss(t, g, time.Second)
		if status != core.StatusOK {
			t.Fatalf("term %d status = %v, want %v", i+1, status, core.StatusOK)
		}
		if term.A().Cmp(big.NewInt(w)) != 0 {
			t.Fatalf("term %d = %v, want %d", i+1, term.A(), w)
		}
	}
}

func nextRCFWithTimeoutPiGauss(t *testing.T, g *core.GCF, timeout time.Duration) (core.RCFTerm, core.Status) {
	t.Helper()

	type result struct {
		term   core.RCFTerm
		status core.Status
	}

	ch := make(chan result, 1)
	go func() {
		term, status := g.NextRCF()
		ch <- result{term: term, status: status}
	}()

	select {
	case got := <-ch:
		return got.term, got.status
	case <-time.After(timeout):
		t.Fatalf("NextRCF timed out after %v", timeout)
		return core.NewRCFTerm(nil), core.StatusInvalidInput
	}
}

func assertPiRange(t *testing.T, got core.Range, want core.Range, step int) {
	t.Helper()

	if got.Inside != want.Inside {
		t.Fatalf("step %d Inside = %v, want %v", step, got.Inside, want.Inside)
	}
	if got.Lo.Open != want.Lo.Open {
		t.Fatalf("step %d Lo.Open = %v, want %v", step, got.Lo.Open, want.Lo.Open)
	}
	if got.Hi.Open != want.Hi.Open {
		t.Fatalf("step %d Hi.Open = %v, want %v", step, got.Hi.Open, want.Hi.Open)
	}
	if got.Lo.Value.Cmp(want.Lo.Value) != 0 {
		t.Fatalf(
			"step %d Lo = %v/%v, want %v/%v",
			step, got.Lo.Value.Num(), got.Lo.Value.Den(),
			want.Lo.Value.Num(), want.Lo.Value.Den(),
		)
	}
	if got.Hi.Value.Cmp(want.Hi.Value) != 0 {
		t.Fatalf(
			"step %d Hi = %v/%v, want %v/%v",
			step, got.Hi.Value.Num(), got.Hi.Value.Den(),
			want.Hi.Value.Num(), want.Hi.Value.Den(),
		)
	}
}

func assertFinitePiStep(
	t *testing.T,
	src core.PQStream,
	step int,
	wantP int64,
	wantQ int64,
	wantRange core.Rational,
) (core.PQTerm, core.PQStream) {
	t.Helper()

	gotRange := src.Range()
	if !gotRange.Inside || gotRange.Lo.Open || gotRange.Hi.Open {
		t.Fatalf("step %d range openness/inside wrong", step)
	}
	if gotRange.Lo.Value.Cmp(wantRange) != 0 || gotRange.Hi.Value.Cmp(wantRange) != 0 {
		t.Fatalf(
			"step %d range = [%v/%v,%v/%v], want exact %v/%v",
			step,
			gotRange.Lo.Value.Num(), gotRange.Lo.Value.Den(),
			gotRange.Hi.Value.Num(), gotRange.Hi.Value.Den(),
			wantRange.Num(), wantRange.Den(),
		)
	}

	term, tail, status := src.NextPQ()
	if status != core.StatusOK {
		t.Fatalf("step %d status = %v, want %v", step, status, core.StatusOK)
	}
	if term.P.Cmp(big.NewInt(wantP)) != 0 {
		t.Fatalf("step %d P = %v, want %d", step, term.P, wantP)
	}
	if term.Q.Cmp(big.NewInt(wantQ)) != 0 {
		t.Fatalf("step %d Q = %v, want %d", step, term.Q, wantQ)
	}

	return term, tail
}

func assertFinitePiEOF(t *testing.T, src core.PQStream, step int) {
	t.Helper()

	_, _, status := src.NextPQ()
	if status != core.StatusEOF {
		t.Fatalf("step %d status = %v, want %v", step, status, core.StatusEOF)
	}
}

// named/pi_gauss_bb_test.go v3
