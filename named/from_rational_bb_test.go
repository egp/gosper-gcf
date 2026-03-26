// named/from_rational_bb_test.go v1
package named_test

import (
	"math/big"
	"testing"

	"github.com/egp/gosper-gcf/core"
	"github.com/egp/gosper-gcf/named"
)

func TestBB_Named_FromRational_ProducesExpectedPrefixAndRange(t *testing.T) {
	src := named.FromRational(core.NewRational(big.NewInt(3), big.NewInt(2)))
	g := core.NewGCF1(identityUnaryCoeffsFromRational(), src)

	r0 := g.Range()
	want0 := core.NewRational(big.NewInt(3), big.NewInt(2))
	assertExactRangeEqualsFromRational(t, r0, want0, 0)

	term1, status1 := g.NextRCF()
	if status1 != core.StatusOK {
		t.Fatalf("first status = %v, want %v", status1, core.StatusOK)
	}
	if term1.A().Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("first term = %v, want 1", term1.A())
	}

	r1 := g.Range()
	want1 := core.NewRational(big.NewInt(2), big.NewInt(1))
	assertExactRangeEqualsFromRational(t, r1, want1, 1)

	term2, status2 := g.NextRCF()
	if status2 != core.StatusOK {
		t.Fatalf("second status = %v, want %v", status2, core.StatusOK)
	}
	if term2.A().Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("second term = %v, want 2", term2.A())
	}

	_, eofStatus := g.NextRCF()
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}
}

func identityUnaryCoeffsFromRational() core.BLFTCoefficients {
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

func assertExactRangeEqualsFromRational(t *testing.T, got core.Range, want core.Rational, step int) {
	t.Helper()

	if !got.Inside {
		t.Fatalf("range %d has Inside=false, want true", step)
	}
	if got.Lo.Open || got.Hi.Open {
		t.Fatalf("range %d openness wrong, want both closed", step)
	}
	if got.Lo.Value.Cmp(want) != 0 || got.Hi.Value.Cmp(want) != 0 {
		t.Fatalf("range %d = [%v/%v,%v/%v], want exact %v/%v",
			step,
			got.Lo.Value.Num(), got.Lo.Value.Den(),
			got.Hi.Value.Num(), got.Hi.Value.Den(),
			want.Num(), want.Den(),
		)
	}
}

// named/from_rational_bb_test.go v1
