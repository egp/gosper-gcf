// core/blft_choose_outside_outside_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

func TestWB_ChooseIngestX_OutsideOutsideEqualWidth_ChoosesXOnTie(t *testing.T) {
	xRange := Range{
		Lo:     Endpoint{Value: RationalFromInt64(-1), Open: false},
		Hi:     Endpoint{Value: NewRational(big.NewInt(540), big.NewInt(301)), Open: false},
		Inside: false,
	}
	yRange := Range{
		Lo:     Endpoint{Value: RationalFromInt64(-1), Open: false},
		Hi:     Endpoint{Value: NewRational(big.NewInt(540), big.NewInt(301)), Open: false},
		Inside: false,
	}

	if !chooseIngestX(xRange, yRange) {
		t.Fatal("chooseIngestX returned false, want true")
	}
}
