// core/blft_unary_lft_outside_range_wb_test.go v3
package core

import (
	"math/big"
	"testing"
)

func TestWB_BLFT_UnaryRange_GeneralLFT_OutsideRange_ReturnsInsideRange(t *testing.T) {
	// z = 1 / (2x - 7); pole at x = 7/2 = 3.5
	// Outside range {Lo=4, Hi=3} = x ≥ 4 OR x ≤ 3 (pole excluded).
	// True image: (0, 1] from x ≥ 4 side, (-1, 0) from x ≤ 3 side.
	// Conservative result: [-1, 1] (Lo open, Hi closed).
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
		Lo:     Endpoint{Value: RationalFromInt64(4), Open: false},
		Hi:     Endpoint{Value: RationalFromInt64(3), Open: true},
		Inside: false,
	}

	got, err := s.UnaryRange(xr)
	if err != nil {
		t.Fatalf("UnaryRange error = %v, want nil", err)
	}
	if !got.Inside {
		t.Fatalf("Inside = false, want true")
	}
	wantLo := RationalFromInt64(-1)
	wantHi := RationalFromInt64(1)
	if got.Lo.Value.Cmp(wantLo) != 0 {
		t.Errorf("Lo = %v/%v, want -1/1", got.Lo.Value.Num(), got.Lo.Value.Den())
	}
	if !got.Lo.Open {
		t.Errorf("Lo.Open = false, want true (x=3 boundary is open)")
	}
	if got.Hi.Value.Cmp(wantHi) != 0 {
		t.Errorf("Hi = %v/%v, want 1/1", got.Hi.Value.Num(), got.Hi.Value.Den())
	}
	if got.Hi.Open {
		t.Errorf("Hi.Open = true, want false (x=4 boundary is closed)")
	}
}

// core/blft_unary_lft_outside_range_wb_test.go v3
