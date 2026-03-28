// core/gcf_unary_outside_range_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

func TestWB_BLFT_UnaryRange_Identity_PreservesOutsideRange(t *testing.T) {
	s := newBLFTState(BLFTCoefficients{
		A: big.NewInt(0),
		B: big.NewInt(1),
		C: big.NewInt(0),
		D: big.NewInt(0),
		E: big.NewInt(0),
		F: big.NewInt(0),
		G: big.NewInt(0),
		H: big.NewInt(1),
	})

	xr := Range{
		Lo:     Endpoint{Value: RationalFromInt64(3), Open: false},
		Hi:     Endpoint{Value: RationalFromInt64(4), Open: true},
		Inside: false,
	}

	got := s.UnaryRange(xr)

	if got.Inside {
		t.Fatal("Inside = true, want false")
	}
	if got.Lo.Open {
		t.Fatal("Lo.Open = true, want false")
	}
	if !got.Hi.Open {
		t.Fatal("Hi.Open = false, want true")
	}
	if got.Lo.Value.Cmp(RationalFromInt64(3)) != 0 {
		t.Fatalf("Lo = %v/%v, want 3/1", got.Lo.Value.Num(), got.Lo.Value.Den())
	}
	if got.Hi.Value.Cmp(RationalFromInt64(4)) != 0 {
		t.Fatalf("Hi = %v/%v, want 4/1", got.Hi.Value.Num(), got.Hi.Value.Den())
	}
}

func TestWB_BLFT_UnaryRange_ScaleHalf_PreservesOutsideRangeAndOpenness(t *testing.T) {
	s := newBLFTState(BLFTCoefficients{
		A: big.NewInt(0),
		B: big.NewInt(1),
		C: big.NewInt(0),
		D: big.NewInt(0),
		E: big.NewInt(0),
		F: big.NewInt(0),
		G: big.NewInt(0),
		H: big.NewInt(2),
	})

	xr := Range{
		Lo:     Endpoint{Value: RationalFromInt64(3), Open: false},
		Hi:     Endpoint{Value: RationalFromInt64(4), Open: true},
		Inside: false,
	}

	got := s.UnaryRange(xr)

	wantLo := NewRational(big.NewInt(3), big.NewInt(2))
	wantHi := RationalFromInt64(2)

	if got.Inside {
		t.Fatal("Inside = true, want false")
	}
	if got.Lo.Open {
		t.Fatal("Lo.Open = true, want false")
	}
	if !got.Hi.Open {
		t.Fatal("Hi.Open = false, want true")
	}
	if got.Lo.Value.Cmp(wantLo) != 0 {
		t.Fatalf("Lo = %v/%v, want 3/2", got.Lo.Value.Num(), got.Lo.Value.Den())
	}
	if got.Hi.Value.Cmp(wantHi) != 0 {
		t.Fatalf("Hi = %v/%v, want 2/1", got.Hi.Value.Num(), got.Hi.Value.Den())
	}
}

func TestWB_BLFT_UnaryRange_Constant_IgnoresOutsideInputAndReturnsExact(t *testing.T) {
	s := newBLFTState(BLFTCoefficients{
		A: big.NewInt(0),
		B: big.NewInt(0),
		C: big.NewInt(0),
		D: big.NewInt(7),
		E: big.NewInt(0),
		F: big.NewInt(0),
		G: big.NewInt(0),
		H: big.NewInt(1),
	})

	xr := Range{
		Lo:     Endpoint{Value: RationalFromInt64(3), Open: false},
		Hi:     Endpoint{Value: RationalFromInt64(4), Open: true},
		Inside: false,
	}

	got := s.UnaryRange(xr)

	if !got.Inside {
		t.Fatal("Inside = false, want true")
	}
	if got.Lo.Open {
		t.Fatal("Lo.Open = true, want false")
	}
	if got.Hi.Open {
		t.Fatal("Hi.Open = true, want false")
	}
	if got.Lo.Value.Cmp(RationalFromInt64(7)) != 0 {
		t.Fatalf("Lo = %v/%v, want 7/1", got.Lo.Value.Num(), got.Lo.Value.Den())
	}
	if got.Hi.Value.Cmp(RationalFromInt64(7)) != 0 {
		t.Fatalf("Hi = %v/%v, want 7/1", got.Hi.Value.Num(), got.Hi.Value.Den())
	}
}

func TestWB_GCF_Range_UnaryIdentity_OverPQStreamFromRCF_WithOutsideRange_PreservesRange(t *testing.T) {
	src := &staticOutsideRCFStream{
		terms: []RCFTerm{
			NewRCFTerm(big.NewInt(3)),
			NewRCFTerm(big.NewInt(7)),
		},
		rng: Range{
			Lo:     Endpoint{Value: RationalFromInt64(3), Open: false},
			Hi:     Endpoint{Value: RationalFromInt64(4), Open: true},
			Inside: false,
		},
	}

	g := NewGCF1(
		BLFTCoefficients{
			A: big.NewInt(0),
			B: big.NewInt(1),
			C: big.NewInt(0),
			D: big.NewInt(0),
			E: big.NewInt(0),
			F: big.NewInt(0),
			G: big.NewInt(0),
			H: big.NewInt(1),
		},
		PQStreamFromRCF(src),
	)

	got := g.Range()

	if got.Inside {
		t.Fatal("Inside = true, want false")
	}
	if got.Lo.Open {
		t.Fatal("Lo.Open = true, want false")
	}
	if !got.Hi.Open {
		t.Fatal("Hi.Open = false, want true")
	}
	if got.Lo.Value.Cmp(RationalFromInt64(3)) != 0 {
		t.Fatalf("Lo = %v/%v, want 3/1", got.Lo.Value.Num(), got.Lo.Value.Den())
	}
	if got.Hi.Value.Cmp(RationalFromInt64(4)) != 0 {
		t.Fatalf("Hi = %v/%v, want 4/1", got.Hi.Value.Num(), got.Hi.Value.Den())
	}
}

// core/gcf_unary_outside_range_wb_test.go v1
