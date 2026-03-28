// core/gcf_constructor_independent_x_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

func TestWB_NewGCF2_IndependentOfX_ProjectY_FirstTermMatchesRightInput(t *testing.T) {
	g := NewGCF2(
		BLFTCoefficients{
			A: big.NewInt(0),
			B: big.NewInt(0),
			C: big.NewInt(1),
			D: big.NewInt(0),
			E: big.NewInt(0),
			F: big.NewInt(0),
			G: big.NewInt(0),
			H: big.NewInt(1),
		},
		PQStreamFromRational(RationalFromInt64(5)),
		PQStreamFromRational(NewRational(big.NewInt(22), big.NewInt(7))),
	)

	term, status := g.NextRCF()
	if status != StatusOK {
		t.Fatalf("first status = %v, want %v", status, StatusOK)
	}
	if term.A().Cmp(big.NewInt(3)) != 0 {
		t.Fatalf("first term = %v, want 3", term.A())
	}
}

// core/gcf_constructor_independent_x_wb_test.go v1
