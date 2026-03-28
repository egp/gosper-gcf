// core/blft_collapse_xeof_runtime_wb_test.go v2
package core

import (
	"math/big"
	"testing"
)

func TestWB_BLFT_CollapseBinaryXEOF_RuntimeAfterIngest_DegreesScale_BecomesIdentityOnY(t *testing.T) {
	s := newBLFTState(BLFTCoefficients{
		A: big.NewInt(1),
		B: big.NewInt(0),
		C: big.NewInt(0),
		D: big.NewInt(0),
		E: big.NewInt(0),
		F: big.NewInt(0),
		G: big.NewInt(0),
		H: big.NewInt(180),
	})

	s = s.IngestX(PQTerm{
		P: big.NewInt(180),
		Q: big.NewInt(1),
	})

	got := exactRationalFromUnaryEngine(
		s.CollapseBinaryXEOF(),
		PQStreamFromRational(NewRational(big.NewInt(22), big.NewInt(7))),
	)
	want := NewRational(big.NewInt(22), big.NewInt(7))

	if got.Cmp(want) != 0 {
		t.Fatalf(
			"collapsed result = %v/%v, want %v/%v",
			got.Num(), got.Den(),
			want.Num(), want.Den(),
		)
	}
}

// core/blft_collapse_xeof_runtime_wb_test.go v2
