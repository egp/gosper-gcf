// core/blft_collapse_xindependent_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

func TestWB_BLFT_CollapseIndependentOfX_ProjectY_PreservesY(t *testing.T) {
	s := newBLFTState(BLFTCoefficients{
		A: big.NewInt(0),
		B: big.NewInt(0),
		C: big.NewInt(1),
		D: big.NewInt(0),
		E: big.NewInt(0),
		F: big.NewInt(0),
		G: big.NewInt(0),
		H: big.NewInt(1),
	})

	u := collapseIndependentOfXToUnary(s)

	got := exactRationalFromUnaryEngine(
		u,
		PQStreamFromRational(NewRational(big.NewInt(22), big.NewInt(7))),
	)

	want := NewRational(big.NewInt(22), big.NewInt(7))
	if got.Cmp(want) != 0 {
		t.Fatalf(
			"collapsed unary applied to 22/7 = %v/%v, want %v/%v",
			got.Num(), got.Den(), want.Num(), want.Den(),
		)
	}
}

// core/blft_collapse_xindependent_wb_test.go v1
