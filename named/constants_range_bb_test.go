// named/constants_range_bb_test.go v1
package named_test

import (
	"math/big"
	"testing"

	"github.com/egp/gosper-gcf/core"
	"github.com/egp/gosper-gcf/named"
)

func TestBB_Named_E_RangeTracksRemainingExactSuffix(t *testing.T) {
	terms := []int64{2, 1, 2, 1, 1, 4, 1, 1, 6}
	suffixes := namedSuffixRationals(terms)

	g := core.NewGCF1(identityUnaryCoeffsNamedRange(), named.E())

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

func namedSuffixRationals(terms []int64) []core.Rational {
	n := len(terms)
	out := make([]core.Rational, n)

	current := core.NewRational(big.NewInt(terms[n-1]), big.NewInt(1))
	out[n-1] = current

	for i := n - 2; i >= 0; i-- {
		a := big.NewInt(terms[i])

		num := new(big.Int).Mul(a, current.Num())
		num.Add(num, current.Den())

		den := current.Num()
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
		t.Fatalf("range before term %d Lo = %v/%v, want %v/%v",
			step, got.Lo.Value.Num(), got.Lo.Value.Den(), want.Num(), want.Den())
	}
	if got.Hi.Value.Cmp(want) != 0 {
		t.Fatalf("range before term %d Hi = %v/%v, want %v/%v",
			step, got.Hi.Value.Num(), got.Hi.Value.Den(), want.Num(), want.Den())
	}
}

// named/constants_range_bb_test.go v1
