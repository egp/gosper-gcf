// core/blft_choose_mixed_pending_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

func TestWB_ChooseIngestX_OutsideX_InsideY_PrefersY(t *testing.T) {
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

	if chooseIngestX(xRange, yRange) {
		t.Fatal("chooseIngestX returned true, want false")
	}
}
