// core/blft_range_wb_test.go v3
package core

import (
	"errors"
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

	got, err := s.CornerRange(xr, yr)
	if err != nil {
		t.Fatalf("CornerRange error = %v", err)
	}
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

	got, err := s.CornerRange(xr, yr)
	if err != nil {
		t.Fatalf("CornerRange error = %v", err)
	}
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

func TestWB_BLFT_CornerRange_OutsideXRangeReturnsUnsupportedError(t *testing.T) {
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
			Value: RationalFromInt64(4),
			Open:  false,
		},
		Hi: Endpoint{
			Value: NewRational(big.NewInt(5), big.NewInt(2)),
			Open:  true,
		},
		Inside: false,
	}

	yr := Range{
		Lo: Endpoint{
			Value: RationalFromInt64(0),
			Open:  false,
		},
		Hi: Endpoint{
			Value: RationalFromInt64(0),
			Open:  false,
		},
		Inside: true,
	}

	_, err := s.CornerRange(xr, yr)
	if !errors.Is(err, ErrUnsupportedRangeCase) {
		t.Fatalf("CornerRange error = %v, want ErrUnsupportedRangeCase", err)
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

	got, err := preferXOnTie(xRange, yRange)
	if err != nil {
		t.Fatalf("preferXOnTie error = %v", err)
	}
	if !got {
		t.Fatal("preferXOnTie returned false, want true")
	}
}

func TestWB_BLFT_TieBreakGoesToYWhenYNarrower(t *testing.T) {
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
			Value: RationalFromInt64(11),
			Open:  false,
		},
		Inside: true,
	}

	got, err := preferXOnTie(xRange, yRange)
	if err != nil {
		t.Fatalf("preferXOnTie error = %v", err)
	}
	if got {
		t.Fatal("preferXOnTie returned true, want false")
	}
}

// core/blft_range_wb_test.go v3
