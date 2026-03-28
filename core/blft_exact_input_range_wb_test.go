// core/blft_exact_input_range_wb_test.go v2
package core

import (
	"math/big"
	"testing"
	"time"
)

type staticOutsideRCFStream struct {
	terms []RCFTerm
	index int
	rng   Range
}

func (s *staticOutsideRCFStream) NextRCF() (RCFTerm, Status) {
	if s.index >= len(s.terms) {
		return NewRCFTerm(nil), StatusEOF
	}
	term := s.terms[s.index]
	s.index++
	return term, StatusOK
}

func (s *staticOutsideRCFStream) Range() Range {
	return s.rng
}

func TestWB_BLFT_BinaryRange_DegreesScale_Exact180PreservesOutsideYRange(t *testing.T) {
	s := blftState{
		A: big.NewInt(1),
		B: big.NewInt(0),
		C: big.NewInt(0),
		D: big.NewInt(0),
		E: big.NewInt(0),
		F: big.NewInt(0),
		G: big.NewInt(0),
		H: big.NewInt(180),
	}

	xr := Range{
		Lo:     Endpoint{Value: RationalFromInt64(180), Open: false},
		Hi:     Endpoint{Value: RationalFromInt64(180), Open: false},
		Inside: true,
	}

	yr := Range{
		Lo:     Endpoint{Value: RationalFromInt64(3), Open: false},
		Hi:     Endpoint{Value: RationalFromInt64(4), Open: true},
		Inside: false,
	}

	got := s.BinaryRange(xr, yr)

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

func TestWB_GCF_BinaryRange_DegreesScale_Exact180PreservesOutsideYRange(t *testing.T) {
	xsrc := &staticRangePQStream{
		rng: Range{
			Lo:     Endpoint{Value: RationalFromInt64(180), Open: false},
			Hi:     Endpoint{Value: RationalFromInt64(180), Open: false},
			Inside: true,
		},
	}

	ysrc := &staticRangePQStream{
		rng: Range{
			Lo:     Endpoint{Value: RationalFromInt64(3), Open: false},
			Hi:     Endpoint{Value: RationalFromInt64(4), Open: true},
			Inside: false,
		},
	}

	g := NewGCF2(
		BLFTCoefficients{
			A: big.NewInt(1),
			B: big.NewInt(0),
			C: big.NewInt(0),
			D: big.NewInt(0),
			E: big.NewInt(0),
			F: big.NewInt(0),
			G: big.NewInt(0),
			H: big.NewInt(180),
		},
		xsrc,
		ysrc,
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

func TestWB_GCF_UnaryIdentity_OverPQStreamFromRCF_OfDegreesScaleOutsideCase_DoesNotPanic(t *testing.T) {
	xsrc := PQStreamFromRational(RationalFromInt64(180))

	yrcf := &staticOutsideRCFStream{
		terms: []RCFTerm{
			NewRCFTerm(big.NewInt(3)),
			NewRCFTerm(big.NewInt(7)),
			NewRCFTerm(big.NewInt(15)),
		},
		rng: Range{
			Lo:     Endpoint{Value: RationalFromInt64(3), Open: false},
			Hi:     Endpoint{Value: RationalFromInt64(4), Open: true},
			Inside: false,
		},
	}

	ysrc := PQStreamFromRCF(yrcf)

	radians := NewGCF2(
		BLFTCoefficients{
			A: big.NewInt(1),
			B: big.NewInt(0),
			C: big.NewInt(0),
			D: big.NewInt(0),
			E: big.NewInt(0),
			F: big.NewInt(0),
			G: big.NewInt(0),
			H: big.NewInt(180),
		},
		xsrc,
		ysrc,
	)

	observed := NewGCF1(
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
		PQStreamFromRCF(radians),
	)

	_, status := nextRCFWithTimeoutExactInputRange(t, observed, time.Second)
	if status != StatusOK {
		t.Fatalf("first status = %v, want %v", status, StatusOK)
	}
}

func nextRCFWithTimeoutExactInputRange(t *testing.T, g *GCF, timeout time.Duration) (RCFTerm, Status) {
	t.Helper()

	type result struct {
		term   RCFTerm
		status Status
	}

	ch := make(chan result, 1)
	go func() {
		term, status := g.NextRCF()
		ch <- result{term: term, status: status}
	}()

	select {
	case got := <-ch:
		return got.term, got.status
	case <-time.After(timeout):
		t.Fatalf("NextRCF() did not complete within %v", timeout)
		return NewRCFTerm(nil), StatusInvalidInput
	}
}

// core/blft_exact_input_range_wb_test.go v2
