// core/blft_sin_from_tan_half_post_ingest_both_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

func TestWB_BLFT_BinaryRange_SinFromTanHalf_AfterUnitIngestBoth_OutsideOutside_ReturnsUnitHull(t *testing.T) {
	s := blftState{
		A: big.NewInt(2),
		B: big.NewInt(0),
		C: big.NewInt(2),
		D: big.NewInt(0),
		E: big.NewInt(2),
		F: big.NewInt(1),
		G: big.NewInt(1),
		H: big.NewInt(1),
	}

	xr := Range{
		Lo:     Endpoint{Value: RationalFromInt64(-1), Open: false},
		Hi:     Endpoint{Value: NewRational(big.NewInt(540), big.NewInt(301)), Open: false},
		Inside: false,
	}
	yr := Range{
		Lo:     Endpoint{Value: RationalFromInt64(-1), Open: false},
		Hi:     Endpoint{Value: NewRational(big.NewInt(540), big.NewInt(301)), Open: false},
		Inside: false,
	}

	got := s.BinaryRange(xr, yr)

	if !got.Inside {
		t.Fatal("Inside = false, want true")
	}
	if got.Lo.Value.Cmp(RationalFromInt64(-1)) != 0 {
		t.Fatalf("Lo = %v/%v, want -1/1", got.Lo.Value.Num(), got.Lo.Value.Den())
	}
	if got.Hi.Value.Cmp(RationalFromInt64(1)) != 0 {
		t.Fatalf("Hi = %v/%v, want 1/1", got.Hi.Value.Num(), got.Hi.Value.Den())
	}
}
