// named/constants_range_bb_test.go v3
package named_test

import (
	"math/big"
	"os"
	"testing"

	"github.com/egp/gosper-gcf/core"
	"github.com/egp/gosper-gcf/named"
)

const pendingTestNamedERange = true

func TestBB_Named_E_RangeTracksBestLookaheadInterval50(t *testing.T) {
	if shouldSkipPendingNamedERange() {
		t.Skip("pending robust named E() range behavior; set RUN_PENDING_TESTS=1 to run anyway")
	}

	terms := eTermsNamedRange(51)
	g := core.NewGCF1(identityUnaryCoeffsNamedRange(), named.E())

	for i := 0; i < 50; i++ {
		wantRange := namedLookaheadRange(terms[i], terms[i+1])
		assertIntervalRangeEqualsNamed(t, g.Range(), wantRange, i+1)

		got, status := g.NextRCF()
		if status != core.StatusOK {
			t.Fatalf("term %d status = %v, want %v", i+1, status, core.StatusOK)
		}

		if got.A().Cmp(big.NewInt(terms[i])) != 0 {
			t.Fatalf("term %d = %v, want %d", i+1, got.A(), terms[i])
		}
	}
}

func TestBB_Named_Pi_RangeTracksRemainingExactSuffix(t *testing.T) {
	terms := []int64{3, 7, 15, 1, 292, 1, 1, 1, 2}
	suffixes := namedSuffixRationals(terms)
	g := core.NewGCF1(identityUnaryCoeffsNamedRange(), named.Pi())

	for i, wantTerm := range terms {
		assertExactRangeEqualsNamed(t, g.Range(), suffixes[i], i+1)

		got, status := g.NextRCF()
		if status != core.StatusOK {
			t.Fatalf("term %d status = %v, want %v", i+1, status, core.StatusOK)
		}

		if got.A().Cmp(big.NewInt(wantTerm)) != 0 {
			t.Fatalf("term %d = %v, want %d", i+1, got.A(), wantTerm)
		}
	}

	_, eofStatus := g.NextRCF()
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}
}

func shouldSkipPendingNamedERange() bool {
	return pendingTestNamedERange && os.Getenv("RUN_PENDING_TESTS") == ""
}

func identityUnaryCoeffsNamedRange() core.BLFTCoefficients {
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

func eTermsNamedRange(n int) []int64 {
	if n <= 0 {
		return nil
	}

	out := make([]int64, 0, n)
	out = append(out, 2)

	k := int64(1)
	for len(out) < n {
		out = append(out, 1)
		if len(out) >= n {
			break
		}
		out = append(out, 2*k)
		if len(out) >= n {
			break
		}
		out = append(out, 1)
		k++
	}

	return out
}

func namedLookaheadRange(a, next int64) core.Range {
	loNum := big.NewInt(a*(next+1) + 1)
	loDen := big.NewInt(next + 1)

	hiNum := big.NewInt(a*next + 1)
	hiDen := big.NewInt(next)

	return core.Range{
		Lo: core.Endpoint{
			Value: core.NewRational(loNum, loDen),
			Open:  false,
		},
		Hi: core.Endpoint{
			Value: core.NewRational(hiNum, hiDen),
			Open:  true,
		},
		Inside: true,
	}
}

func namedSuffixRationals(terms []int64) []core.Rational {
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

func assertExactRangeEqualsNamed(t *testing.T, got core.Range, want core.Rational, step int) {
	t.Helper()

	if !got.Inside {
		t.Fatalf("range before term %d has Inside=false, want true", step)
	}
	if got.Lo.Open {
		t.Fatalf("range before term %d has Lo.Open=true, want false", step)
	}
	if got.Hi.Open {
		t.Fatalf("range before term %d has Hi.Open=true, want false", step)
	}
	if got.Lo.Value.Cmp(want) != 0 {
		t.Fatalf(
			"range before term %d Lo = %v/%v, want %v/%v",
			step, got.Lo.Value.Num(), got.Lo.Value.Den(), want.Num(), want.Den(),
		)
	}
	if got.Hi.Value.Cmp(want) != 0 {
		t.Fatalf(
			"range before term %d Hi = %v/%v, want %v/%v",
			step, got.Hi.Value.Num(), got.Hi.Value.Den(), want.Num(), want.Den(),
		)
	}
}

func assertIntervalRangeEqualsNamed(t *testing.T, got core.Range, want core.Range, step int) {
	t.Helper()

	if got.Inside != want.Inside {
		t.Fatalf("range before term %d Inside = %v, want %v", step, got.Inside, want.Inside)
	}
	if got.Lo.Open != want.Lo.Open {
		t.Fatalf("range before term %d Lo.Open = %v, want %v", step, got.Lo.Open, want.Lo.Open)
	}
	if got.Hi.Open != want.Hi.Open {
		t.Fatalf("range before term %d Hi.Open = %v, want %v", step, got.Hi.Open, want.Hi.Open)
	}
	if got.Lo.Value.Cmp(want.Lo.Value) != 0 {
		t.Fatalf(
			"range before term %d Lo = %v/%v, want %v/%v",
			step, got.Lo.Value.Num(), got.Lo.Value.Den(), want.Lo.Value.Num(), want.Lo.Value.Den(),
		)
	}
	if got.Hi.Value.Cmp(want.Hi.Value) != 0 {
		t.Fatalf(
			"range before term %d Hi = %v/%v, want %v/%v",
			step, got.Hi.Value.Num(), got.Hi.Value.Den(), want.Hi.Value.Num(), want.Hi.Value.Den(),
		)
	}
}

// named/constants_range_bb_test.go v3
