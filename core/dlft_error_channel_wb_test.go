// core/dlft_error_channel_wb_test.go v1
package core

import (
	"errors"
	"math/big"
	"testing"
)

func TestWB_DLFTUnaryRange_OutsideInputReturnsError(t *testing.T) {
	s := newDLFTState(DLFTCoefficients{
		A: big.NewInt(1),
		B: big.NewInt(0),
		C: big.NewInt(0),
		D: big.NewInt(1),
		E: big.NewInt(0),
		F: big.NewInt(1),
	})

	xRange := Range{
		Lo:     Endpoint{Value: RationalFromInt64(5), Open: false},
		Hi:     Endpoint{Value: RationalFromInt64(2), Open: false},
		Inside: false,
	}

	_, err := s.UnaryRange(xRange)
	if !errors.Is(err, ErrUnsupportedRangeCase) {
		t.Fatalf("UnaryRange error = %v, want ErrUnsupportedRangeCase", err)
	}
}

// core/dlft_error_channel_wb_test.go v1
