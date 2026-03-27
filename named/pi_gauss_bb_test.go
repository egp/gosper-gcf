// named/pi_gauss_bb_test.go v2
package named_test

import (
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/egp/gosper-gcf/core"
	"github.com/egp/gosper-gcf/named"
)

const pendingTestNamedPiGauss = true

func TestBB_Named_Pi_SourcePQPrefix50(t *testing.T) {
	if shouldSkipPendingNamedPiGauss() {
		t.Skip("pending procedural Gauss Pi source; set RUN_PENDING_TESTS=1 to run anyway")
	}

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
	if shouldSkipPendingNamedPiGauss() {
		t.Skip("pending procedural Gauss Pi source; set RUN_PENDING_TESTS=1 to run anyway")
	}

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

func TestBB_Named_Pi_RCFPrefixMatchesKnownPrefix(t *testing.T) {
	if shouldSkipPendingNamedPiGauss() {
		t.Skip("pending procedural Gauss Pi source; set RUN_PENDING_TESTS=1 to run anyway")
	}

	g := core.NewGCF1(identityUnaryCoeffsPiGauss(), named.Pi())
	assertRCFPrefixPiGauss(t, g, []int64{3, 7, 15, 1, 292, 1, 1, 1})
}

func TestBB_Named_Pi_TakeAndRationalValidate(t *testing.T) {
	if shouldSkipPendingNamedPiGauss() {
		t.Skip("pending procedural Gauss Pi source; set RUN_PENDING_TESTS=1 to run anyway")
	}

	g := core.NewGCF1(identityUnaryCoeffsPiGauss(), named.Pi())

	taken := g.Take(4)

	_, tail2 := assertFinitePiStep(t, taken, 1, 3, 1, core.NewRational(big.NewInt(355), big.NewInt(113)))
	_, tail3 := assertFinitePiStep(t, tail2, 2, 7, 1, core.NewRational(big.NewInt(113), big.NewInt(16)))
	_, tail4 := assertFinitePiStep(t, tail3, 3, 15, 1, core.NewRational(big.NewInt(16), big.NewInt(1)))
	_, tail5 := assertFinitePiStep(t, tail4, 4, 1, 1, core.RationalFromInt64(1))
	assertFinitePiEOF(t, tail5, 5)

	g2 := core.NewGCF1(identityUnaryCoeffsPiGauss(), named.Pi())
	if got := g2.Rational(1); got.Cmp(core.RationalFromInt64(3)) != 0 {
		t.Fatalf("Rational(1) = %v/%v, want 3/1", got.Num(), got.Den())
	}
	g3 := core.NewGCF1(identityUnaryCoeffsPiGauss(), named.Pi())
	if got := g3.Rational(2); got.Cmp(core.NewRational(big.NewInt(22), big.NewInt(7))) != 0 {
		t.Fatalf("Rational(2) = %v/%v, want 22/7", got.Num(), got.Den())
	}
	g4 := core.NewGCF1(identityUnaryCoeffsPiGauss(), named.Pi())
	if got := g4.Rational(4); got.Cmp(core.NewRational(big.NewInt(355), big.NewInt(113))) != 0 {
		t.Fatalf("Rational(4) = %v/%v, want 355/113", got.Num(), got.Den())
	}
}

func shouldSkipPendingNamedPiGauss() bool {
	return pendingTestNamedPiGauss && os.Getenv("RUN_PENDING_TESTS") == ""
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
		t.Fatalf("step %d Lo = %v/%v, want %v/%v",
			step, got.Lo.Value.Num(), got.Lo.Value.Den(),
			want.Lo.Value.Num(), want.Lo.Value.Den(),
		)
	}
	if got.Hi.Value.Cmp(want.Hi.Value) != 0 {
		t.Fatalf("step %d Hi = %v/%v, want %v/%v",
			step, got.Hi.Value.Num(), got.Hi.Value.Den(),
			want.Hi.Value.Num(), want.Hi.Value.Den(),
		)
	}
}

func assertFinitePiStep(t *testing.T, src core.PQStream, step int, wantP, wantQ int64, wantRange core.Rational) (core.PQTerm, core.PQStream) {
	t.Helper()

	gotRange := src.Range()
	if !gotRange.Inside || gotRange.Lo.Open || gotRange.Hi.Open {
		t.Fatalf("step %d range openness/inside wrong", step)
	}
	if gotRange.Lo.Value.Cmp(wantRange) != 0 || gotRange.Hi.Value.Cmp(wantRange) != 0 {
		t.Fatalf("step %d range = [%v/%v,%v/%v], want exact %v/%v",
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

// named/pi_gauss_bb_test.go v2
