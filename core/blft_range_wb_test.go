// core/blft_range_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

func TestWB_BLFT_CornerRangeIdentityOverInsideRectangle(t *testing.T) {
	s := blftState{
		A: big.NewInt(0),
		B: big.NewInt(1),
		C: big.NewInt(0),
		D: big.NewInt(0),
		E: big.NewInt(0),
		F: big.NewInt(0),
		G: big.NewInt(0),
		H: big.NewInt(1),
	}

	xr := Range{
		Lo: Endpoint{
			Value: RationalFromInt64(2),
			Open:  false,
		},
		Hi: Endpoint{
			Value: RationalFromInt64(5),
			Open:  false,
		},
		Inside: true,
	}
	yr := Range{
		Lo: Endpoint{
			Value: RationalFromInt64(7),
			Open:  false,
		},
		Hi: Endpoint{
			Value: RationalFromInt64(11),
			Open:  false,
		},
		Inside: true,
	}

	got := s.CornerRange(xr, yr)

	if !got.Inside {
		t.Fatal("Inside = false, want true")
	}
	if got.Lo.Value.Cmp(RationalFromInt64(2)) != 0 {
		t.Fatalf("Lo = %v/%v, want 2/1", got.Lo.Value.Num(), got.Lo.Value.Den())
	}
	if got.Hi.Value.Cmp(RationalFromInt64(5)) != 0 {
		t.Fatalf("Hi = %v/%v, want 5/1", got.Hi.Value.Num(), got.Hi.Value.Den())
	}
}

func TestWB_BLFT_CornerRangeConstantFunctionIsExact(t *testing.T) {
	s := blftState{
		A: big.NewInt(0),
		B: big.NewInt(0),
		C: big.NewInt(0),
		D: big.NewInt(3),
		E: big.NewInt(0),
		F: big.NewInt(0),
		G: big.NewInt(0),
		H: big.NewInt(2),
	}

	xr := Range{
		Lo: Endpoint{
			Value: RationalFromInt64(2),
			Open:  false,
		},
		Hi: Endpoint{
			Value: RationalFromInt64(5),
			Open:  false,
		},
		Inside: true,
	}
	yr := Range{
		Lo: Endpoint{
			Value: RationalFromInt64(7),
			Open:  false,
		},
		Hi: Endpoint{
			Value: RationalFromInt64(11),
			Open:  false,
		},
		Inside: true,
	}

	got := s.CornerRange(xr, yr)
	want := NewRational(big.NewInt(3), big.NewInt(2))

	if !got.Inside {
		t.Fatal("Inside = false, want true")
	}
	if got.Lo.Value.Cmp(want) != 0 {
		t.Fatalf("Lo = %v/%v, want 3/2", got.Lo.Value.Num(), got.Lo.Value.Den())
	}
	if got.Hi.Value.Cmp(want) != 0 {
		t.Fatalf("Hi = %v/%v, want 3/2", got.Hi.Value.Num(), got.Hi.Value.Den())
	}
}

func TestWB_BLFT_TieBreakGoesToX(t *testing.T) {
	xRange := Range{
		Lo: Endpoint{
			Value: RationalFromInt64(2),
			Open:  false,
		},
		Hi: Endpoint{
			Value: RationalFromInt64(5),
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
			Value: RationalFromInt64(13),
			Open:  false,
		},
		Inside: true,
	}

	if !preferXOnTie(xRange, yRange) {
		t.Fatal("preferXOnTie returned false, want true")
	}
}

// core/blft_range_wb_test.go v1
