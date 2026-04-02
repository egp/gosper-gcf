// core/gcf_unary_outside_range_wb_test.go v3
package core

import (
	"errors"
	"math/big"
	"testing"
)

type staticOutsideRCFStream struct {
	terms []RCFTerm
	rng   Range
	index int
}

func (s *staticOutsideRCFStream) NextRCF() (RCFTerm, Status, error) {
	if s.index >= len(s.terms) {
		return NewRCFTerm(nil), StatusEOF, nil
	}
	term := s.terms[s.index]
	s.index++
	return term, StatusOK, nil
}

func (s *staticOutsideRCFStream) CurrentInterval() (Interval, error) {
	return s.rng, nil
}

func (s *staticOutsideRCFStream) Range() (Range, error) {
	return s.CurrentInterval()
}

func TestWB_BLFT_UnaryRange_Identity_OutsideRangeReturnsUnsupportedError(t *testing.T) {
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
		Lo:     Endpoint{Value: RationalFromInt64(4), Open: true},
		Hi:     Endpoint{Value: RationalFromInt64(3), Open: false},
		Inside: false,
	}

	_, err := s.UnaryRange(xr)
	if !errors.Is(err, ErrUnsupportedRangeCase) {
		t.Fatalf("UnaryRange error = %v, want ErrUnsupportedRangeCase", err)
	}
}

func TestWB_BLFT_UnaryRange_ScaleHalf_OutsideRangeReturnsUnsupportedError(t *testing.T) {
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
		Lo:     Endpoint{Value: RationalFromInt64(4), Open: true},
		Hi:     Endpoint{Value: RationalFromInt64(3), Open: false},
		Inside: false,
	}

	_, err := s.UnaryRange(xr)
	if !errors.Is(err, ErrUnsupportedRangeCase) {
		t.Fatalf("UnaryRange error = %v, want ErrUnsupportedRangeCase", err)
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
		Lo:     Endpoint{Value: RationalFromInt64(4), Open: true},
		Hi:     Endpoint{Value: RationalFromInt64(3), Open: false},
		Inside: false,
	}

	got, err := s.UnaryRange(xr)
	if err != nil {
		t.Fatalf("UnaryRange error = %v", err)
	}
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

func TestWB_GCF_Range_UnaryIdentity_OverPQStreamFromRCF_WithOutsideRangeReturnsError(t *testing.T) {
	src := &staticOutsideRCFStream{
		terms: []RCFTerm{
			NewRCFTerm(big.NewInt(3)),
			NewRCFTerm(big.NewInt(7)),
		},
		rng: Range{
			Lo:     Endpoint{Value: RationalFromInt64(4), Open: true},
			Hi:     Endpoint{Value: RationalFromInt64(3), Open: false},
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

	_, err := g.Range()
	if !errors.Is(err, ErrUnsupportedRangeCase) {
		t.Fatalf("Range error = %v, want ErrUnsupportedRangeCase", err)
	}
}

// core/gcf_unary_outside_range_wb_test.go v3
