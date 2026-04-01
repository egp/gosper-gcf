// core/blft_prefer_x_error_wb_test.go v1
package core

import "testing"

func TestWB_PreferXOnTie_MixedRangesReturnsErrorInsteadOfPanicking(t *testing.T) {
	xRange := Range{
		Lo:     Endpoint{Value: RationalFromInt64(5), Open: false},
		Hi:     Endpoint{Value: RationalFromInt64(2), Open: false},
		Inside: false,
	}
	yRange := Range{
		Lo:     Endpoint{Value: RationalFromInt64(0), Open: false},
		Hi:     Endpoint{Value: RationalFromInt64(1), Open: false},
		Inside: true,
	}

	_, err := preferXOnTie(xRange, yRange)
	if err == nil {
		t.Fatal("preferXOnTie error = nil, want error")
	}
}

// core/blft_prefer_x_error_wb_test.go v1
