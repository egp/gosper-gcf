// core/gcf_collapse_eof_wb_test.go v2
package core

import (
	"math/big"
	"testing"
)

func TestWB_CollapseBinaryXEOF_IndependentOfX_ProjectY_PreservesY(t *testing.T) {
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

	got, err := exactRationalFromUnaryEngine(
		s.CollapseBinaryXEOF(),
		PQStreamFromRational(NewRational(big.NewInt(22), big.NewInt(7))),
	)
	if err != nil {
		t.Fatalf("exactRationalFromUnaryEngine error = %v", err)
	}

	want := NewRational(big.NewInt(22), big.NewInt(7))
	if got.Cmp(want) != 0 {
		t.Fatalf(
			"collapsed result = %v/%v, want %v/%v",
			got.Num(), got.Den(),
			want.Num(), want.Den(),
		)
	}
}

func TestWB_CollapseBinaryXEOF_RuntimeAfterIngest_DegreesScale_BecomesIdentityOnY(t *testing.T) {
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

	got, err := exactRationalFromUnaryEngine(
		s.CollapseBinaryXEOF(),
		PQStreamFromRational(NewRational(big.NewInt(22), big.NewInt(7))),
	)
	if err != nil {
		t.Fatalf("exactRationalFromUnaryEngine error = %v", err)
	}

	want := NewRational(big.NewInt(22), big.NewInt(7))
	if got.Cmp(want) != 0 {
		t.Fatalf(
			"collapsed result = %v/%v, want %v/%v",
			got.Num(), got.Den(),
			want.Num(), want.Den(),
		)
	}
}

// core/gcf_collapse_eof_wb_test.go v2
