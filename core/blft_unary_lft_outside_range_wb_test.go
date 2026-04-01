// core/blft_unary_lft_outside_range_wb_test.go v2
package core

import (
	"errors"
	"math/big"
	"testing"
)

func TestWB_BLFT_UnaryRange_GeneralLFT_OutsideRangeReturnsUnsupportedError(t *testing.T) {
	// z = 1 / (2x - 7)
	// Outside-range support is not implemented generically yet, so the current
	// contract is to return a typed error instead of panicking.
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

	_, err := s.UnaryRange(xr)
	if !errors.Is(err, ErrUnsupportedRangeCase) {
		t.Fatalf("UnaryRange error = %v, want ErrUnsupportedRangeCase", err)
	}
}

// core/blft_unary_lft_outside_range_wb_test.go v2
