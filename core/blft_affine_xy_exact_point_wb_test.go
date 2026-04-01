// core/blft_affine_xy_exact_point_wb_test.go v2
package core

import (
	"math/big"
	"testing"
)

func TestWB_BLFT_BinaryRange_AffineAverage_ExactPointX_InsideY(t *testing.T) {
	s := newBLFTState(BLFTCoefficients{
		A: big.NewInt(0),
		B: big.NewInt(1),
		C: big.NewInt(1),
		D: big.NewInt(0),
		E: big.NewInt(0),
		F: big.NewInt(0),
		G: big.NewInt(0),
		H: big.NewInt(2),
	})

	xr := Range{
		Lo: Endpoint{
			Value: NewRational(big.NewInt(218340), big.NewInt(75601)),
			Open:  false,
		},
		Hi: Endpoint{
			Value: NewRational(big.NewInt(218340), big.NewInt(75601)),
			Open:  false,
		},
		Inside: false,
	}

	yr := Range{
		Lo: Endpoint{
			Value: RationalFromInt64(0),
			Open:  true,
		},
		Hi: Endpoint{
			Value: NewRational(big.NewInt(75601), big.NewInt(218340)),
			Open:  true,
		},
		Inside: true,
	}

	got, err := s.BinaryRange(xr, yr)
	if err != nil {
		t.Fatalf("BinaryRange error = %v", err)
	}
	if !got.Inside {
		t.Fatal("Inside = false, want true")
	}
}

// core/blft_affine_xy_exact_point_wb_test.go v2
