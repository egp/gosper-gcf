// core/blft_range_mixed_error_pending_wb_test.go v1
package core

import (
	"math/big"
	"strings"
	"testing"
)

func TestWB_CornerRange_OutsideInsideReturnsDiagnosticErrorInsteadOfPanicking(t *testing.T) {
	s := blftState{
		A: big.NewInt(0),
		B: big.NewInt(2),
		C: big.NewInt(0),
		D: big.NewInt(2),
		E: big.NewInt(1),
		F: big.NewInt(1),
		G: big.NewInt(1),
		H: big.NewInt(0),
	}

	xRange := Range{
		Lo:     Endpoint{Value: RationalFromInt64(-1), Open: false},
		Hi:     Endpoint{Value: NewRational(big.NewInt(540), big.NewInt(301)), Open: false},
		Inside: false,
	}
	yRange := Range{
		Lo:     Endpoint{Value: RationalFromInt64(0), Open: true},
		Hi:     Endpoint{Value: NewRational(big.NewInt(841), big.NewInt(540)), Open: true},
		Inside: true,
	}

	_, err := s.CornerRange(xRange, yRange)
	if err == nil {
		t.Fatal("CornerRange error = nil, want diagnostic error")
	}

	for _, want := range []string{"CornerRange", "BLFT=", "xRange=", "yRange="} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("CornerRange error %q missing %q", err.Error(), want)
		}
	}
}

// core/blft_range_mixed_error_pending_wb_test.go v1
