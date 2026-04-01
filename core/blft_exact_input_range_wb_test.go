// core/blft_exact_input_range_wb_test.go v3
package core

import (
	"math/big"
	"testing"
	"time"
)

type exactInputRangeRCFStream struct {
	terms []RCFTerm
	index int
	rng   Range
}

func (s *exactInputRangeRCFStream) NextRCF() (RCFTerm, Status, error) {
	if s.index >= len(s.terms) {
		return NewRCFTerm(nil), StatusEOF, nil
	}
	term := s.terms[s.index]
	s.index++
	return term, StatusOK, nil
}

func (s *exactInputRangeRCFStream) Range() (Range, error) {
	return s.rng, nil
}

type exactInputRangePQStream struct {
	rng Range
}

func (s *exactInputRangePQStream) NextPQ() (PQTerm, PQStream, Status, error) {
	return PQTerm{
		P: big.NewInt(0),
		Q: big.NewInt(0),
	}, s, StatusEOF, nil
}

func (s *exactInputRangePQStream) Range() (Range, error) {
	return s.rng, nil
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
		Lo: Endpoint{
			Value: RationalFromInt64(180),
			Open:  false,
		},
		Hi: Endpoint{
			Value: RationalFromInt64(180),
			Open:  false,
		},
		Inside: true,
	}

	yr := Range{
		Lo: Endpoint{
			Value: RationalFromInt64(4),
			Open:  true,
		},
		Hi: Endpoint{
			Value: RationalFromInt64(3),
			Open:  false,
		},
		Inside: false,
	}

	got, err := s.BinaryRange(xr, yr)
	if err != nil {
		t.Fatalf("BinaryRange error = %v", err)
	}

	if got.Inside {
		t.Fatal("Inside = true, want false")
	}
	if !got.Lo.Open {
		t.Fatal("Lo.Open = false, want true")
	}
	if got.Hi.Open {
		t.Fatal("Hi.Open = true, want false")
	}
	if got.Lo.Value.Cmp(RationalFromInt64(4)) != 0 {
		t.Fatalf("Lo = %v/%v, want 4/1", got.Lo.Value.Num(), got.Lo.Value.Den())
	}
	if got.Hi.Value.Cmp(RationalFromInt64(3)) != 0 {
		t.Fatalf("Hi = %v/%v, want 3/1", got.Hi.Value.Num(), got.Hi.Value.Den())
	}
}

func TestWB_GCF_BinaryRange_DegreesScale_Exact180PreservesOutsideYRange(t *testing.T) {
	xsrc := &exactInputRangePQStream{
		rng: Range{
			Lo: Endpoint{
				Value: RationalFromInt64(180),
				Open:  false,
			},
			Hi: Endpoint{
				Value: RationalFromInt64(180),
				Open:  false,
			},
			Inside: true,
		},
	}

	ysrc := &exactInputRangePQStream{
		rng: Range{
			Lo: Endpoint{
				Value: RationalFromInt64(4),
				Open:  true,
			},
			Hi: Endpoint{
				Value: RationalFromInt64(3),
				Open:  false,
			},
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

	got, err := g.Range()
	if err != nil {
		t.Fatalf("Range error = %v", err)
	}

	if got.Inside {
		t.Fatal("Inside = true, want false")
	}
	if !got.Lo.Open {
		t.Fatal("Lo.Open = false, want true")
	}
	if got.Hi.Open {
		t.Fatal("Hi.Open = true, want false")
	}
	if got.Lo.Value.Cmp(RationalFromInt64(4)) != 0 {
		t.Fatalf("Lo = %v/%v, want 4/1", got.Lo.Value.Num(), got.Lo.Value.Den())
	}
	if got.Hi.Value.Cmp(RationalFromInt64(3)) != 0 {
		t.Fatalf("Hi = %v/%v, want 3/1", got.Hi.Value.Num(), got.Hi.Value.Den())
	}
}

func TestWB_GCF_UnaryIdentity_OverPQStreamFromRCF_OfDegreesScaleOutsideCase_DoesNotPanic(t *testing.T) {
	xsrc := PQStreamFromRational(RationalFromInt64(180))

	yrcf := &exactInputRangeRCFStream{
		terms: []RCFTerm{
			NewRCFTerm(big.NewInt(3)),
			NewRCFTerm(big.NewInt(7)),
			NewRCFTerm(big.NewInt(15)),
		},
		rng: Range{
			Lo: Endpoint{
				Value: RationalFromInt64(4),
				Open:  true,
			},
			Hi: Endpoint{
				Value: RationalFromInt64(3),
				Open:  false,
			},
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

	_, status, err := nextRCFWithTimeoutExactInputRange(t, observed, time.Second)
	if err != nil {
		t.Fatalf("NextRCF error = %v", err)
	}
	if status != StatusOK {
		t.Fatalf("first status = %v, want %v", status, StatusOK)
	}
}

func TestWB_ExactInputRange_PQHelper_ImplementsErrorChannelShape(t *testing.T) {
	src := &exactInputRangePQStream{
		rng: Range{
			Lo: Endpoint{
				Value: RationalFromInt64(9),
				Open:  false,
			},
			Hi: Endpoint{
				Value: RationalFromInt64(9),
				Open:  false,
			},
			Inside: true,
		},
	}

	_, _, status, err := src.NextPQ()
	if err != nil {
		t.Fatalf("NextPQ error = %v", err)
	}
	if status != StatusEOF {
		t.Fatalf("status = %v, want %v", status, StatusEOF)
	}

	rng, err := src.Range()
	if err != nil {
		t.Fatalf("Range error = %v", err)
	}
	if rng.Lo.Value.Cmp(RationalFromInt64(9)) != 0 || rng.Hi.Value.Cmp(RationalFromInt64(9)) != 0 {
		t.Fatalf("Range = [%v/%v,%v/%v], want exact 9/1",
			rng.Lo.Value.Num(), rng.Lo.Value.Den(),
			rng.Hi.Value.Num(), rng.Hi.Value.Den(),
		)
	}
}

func nextRCFWithTimeoutExactInputRange(t *testing.T, g *GCF, timeout time.Duration) (RCFTerm, Status, error) {
	t.Helper()

	type result struct {
		term   RCFTerm
		status Status
		err    error
	}

	ch := make(chan result, 1)
	go func() {
		term, status, err := g.NextRCF()
		ch <- result{term: term, status: status, err: err}
	}()

	select {
	case got := <-ch:
		return got.term, got.status, got.err
	case <-time.After(timeout):
		t.Fatalf("NextRCF() did not complete within %v", timeout)
		return NewRCFTerm(nil), StatusInvalidInput, nil
	}
}

// core/blft_exact_input_range_wb_test.go v3
