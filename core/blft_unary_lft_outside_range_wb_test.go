// core/blft_unary_lft_outside_range_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

func TestWB_BLFT_UnaryRange_GeneralLFT_OutsideRange_PoleInExcludedGap_ReturnsConservativeHull(t *testing.T) {
	// z = 1 / (2x - 7)
	// Pole at x = 7/2 lies in the excluded gap (3,4), not in the included outside domain.
	s := newBLFTState(BLFTCoefficients{
		A: big.NewInt(0),
		B: big.NewInt(0),
		C: big.NewInt(0),
		D: big.NewInt(1),
		E: big.NewInt(0),
		F: big.NewInt(2),
		G: big.NewInt(0),
		H: big.NewInt(-7),
	})

	xr := Range{
		Lo:     Endpoint{Value: RationalFromInt64(3), Open: false},
		Hi:     Endpoint{Value: RationalFromInt64(4), Open: true},
		Inside: false,
	}

	got := s.UnaryRange(xr)

	if !got.Inside {
		t.Fatal("Inside = false, want true conservative hull")
	}

	wantLo := RationalFromInt64(-1)
	wantHi := RationalFromInt64(1)

	if got.Lo.Value.Cmp(wantLo) != 0 {
		t.Fatalf("Lo = %v/%v, want -1/1", got.Lo.Value.Num(), got.Lo.Value.Den())
	}
	if got.Hi.Value.Cmp(wantHi) != 0 {
		t.Fatalf("Hi = %v/%v, want 1/1", got.Hi.Value.Num(), got.Hi.Value.Den())
	}
}

// core/blft_unary_lft_outside_range_wb_test.go v1
