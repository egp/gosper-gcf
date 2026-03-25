// core/dlft_range_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

func TestWB_DLFT_CandidateRange_ExactMonotoneEndpoints(t *testing.T) {
	s := dlftState{
		A: big.NewInt(1),
		B: big.NewInt(0),
		C: big.NewInt(0),
		D: big.NewInt(0),
		E: big.NewInt(0),
		F: big.NewInt(1),
	}

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

	got := s.CandidateRange(xRange)

	wantLo := RationalFromInt64(4)
	wantHi := RationalFromInt64(25)

	if !got.Inside {
		t.Fatal("CandidateRange().Inside = false, want true")
	}
	if got.Lo.Value.Cmp(wantLo) != 0 {
		t.Fatalf("Lo = %v/%v, want 4/1", got.Lo.Value.Num(), got.Lo.Value.Den())
	}
	if got.Hi.Value.Cmp(wantHi) != 0 {
		t.Fatalf("Hi = %v/%v, want 25/1", got.Hi.Value.Num(), got.Hi.Value.Den())
	}
}

func TestWB_DLFT_CandidateRange_DetectsInteriorPoleAsOutside(t *testing.T) {
	s := dlftState{
		A: big.NewInt(0),
		B: big.NewInt(0),
		C: big.NewInt(1),
		D: big.NewInt(1),
		E: big.NewInt(0),
		F: big.NewInt(-2),
	}

	xRange := Range{
		Lo: Endpoint{
			Value: RationalFromInt64(1),
			Open:  false,
		},
		Hi: Endpoint{
			Value: RationalFromInt64(2),
			Open:  false,
		},
		Inside: true,
	}

	got := s.CandidateRange(xRange)

	wantLo := RationalFromInt64(-1)
	wantHi := NewRational(big.NewInt(1), big.NewInt(2))

	if got.Inside {
		t.Fatal("CandidateRange().Inside = true, want false because denominator has an interior root")
	}
	if got.Lo.Value.Cmp(wantLo) != 0 {
		t.Fatalf("Lo = %v/%v, want -1/1", got.Lo.Value.Num(), got.Lo.Value.Den())
	}
	if got.Hi.Value.Cmp(wantHi) != 0 {
		t.Fatalf("Hi = %v/%v, want 1/2", got.Hi.Value.Num(), got.Hi.Value.Den())
	}
}

func TestWB_DLFT_CandidateRange_UsesInteriorCriticalPoint(t *testing.T) {
	s := dlftState{
		A: big.NewInt(1),
		B: big.NewInt(0),
		C: big.NewInt(0),
		D: big.NewInt(0),
		E: big.NewInt(0),
		F: big.NewInt(1),
	}

	xRange := Range{
		Lo: Endpoint{
			Value: RationalFromInt64(-1),
			Open:  false,
		},
		Hi: Endpoint{
			Value: RationalFromInt64(1),
			Open:  false,
		},
		Inside: true,
	}

	got := s.CandidateRange(xRange)

	wantLo := RationalFromInt64(0)
	wantHi := RationalFromInt64(1)

	if !got.Inside {
		t.Fatal("CandidateRange().Inside = false, want true")
	}
	if got.Lo.Value.Cmp(wantLo) != 0 {
		t.Fatalf("Lo = %v/%v, want 0/1", got.Lo.Value.Num(), got.Lo.Value.Den())
	}
	if got.Hi.Value.Cmp(wantHi) != 0 {
		t.Fatalf("Hi = %v/%v, want 1/1", got.Hi.Value.Num(), got.Hi.Value.Den())
	}
}

// core/dlft_range_wb_test.go v1
