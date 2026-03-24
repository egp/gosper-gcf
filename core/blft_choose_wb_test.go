// core/blft_choose_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

func TestWB_BLFT_WidestCornerRangeChoosesXWhenXIsWider(t *testing.T) {
	xRange := Range{
		Lo: Endpoint{
			Value: RationalFromInt64(2),
			Open:  false,
		},
		Hi: Endpoint{
			Value: RationalFromInt64(8),
			Open:  false,
		},
		Inside: true,
	}
	yRange := Range{
		Lo: Endpoint{
			Value: RationalFromInt64(10),
			Open:  false,
		},
		Hi: Endpoint{
			Value: RationalFromInt64(12),
			Open:  false,
		},
		Inside: true,
	}

	if !chooseIngestX(xRange, yRange) {
		t.Fatal("chooseIngestX returned false, want true")
	}
}

func TestWB_BLFT_WidestCornerRangeChoosesYWhenYIsWider(t *testing.T) {
	xRange := Range{
		Lo: Endpoint{
			Value: RationalFromInt64(2),
			Open:  false,
		},
		Hi: Endpoint{
			Value: RationalFromInt64(4),
			Open:  false,
		},
		Inside: true,
	}
	yRange := Range{
		Lo: Endpoint{
			Value: RationalFromInt64(10),
			Open:  false,
		},
		Hi: Endpoint{
			Value: RationalFromInt64(15),
			Open:  false,
		},
		Inside: true,
	}

	if chooseIngestX(xRange, yRange) {
		t.Fatal("chooseIngestX returned true, want false")
	}
}

func TestWB_BLFT_WidestCornerRangeTieBreakGoesToX(t *testing.T) {
	xRange := Range{
		Lo: Endpoint{
			Value: NewRational(big.NewInt(2), big.NewInt(1)),
			Open:  false,
		},
		Hi: Endpoint{
			Value: NewRational(big.NewInt(5), big.NewInt(1)),
			Open:  false,
		},
		Inside: true,
	}
	yRange := Range{
		Lo: Endpoint{
			Value: NewRational(big.NewInt(10), big.NewInt(1)),
			Open:  false,
		},
		Hi: Endpoint{
			Value: NewRational(big.NewInt(13), big.NewInt(1)),
			Open:  false,
		},
		Inside: true,
	}

	if !chooseIngestX(xRange, yRange) {
		t.Fatal("chooseIngestX returned false on tie, want true")
	}
}

// core/blft_choose_wb_test.go v1
