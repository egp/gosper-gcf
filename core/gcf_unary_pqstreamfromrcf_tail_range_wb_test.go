// core/gcf_unary_pqstreamfromrcf_tail_range_wb_test.go v3
package core

import (
	"math/big"
	"testing"
	"time"
)

type steppingRCFStream struct {
	terms  []RCFTerm
	ranges []Range
	index  int
}

func (s *steppingRCFStream) NextRCF() (RCFTerm, Status, error) {
	if s.index >= len(s.terms) {
		return NewRCFTerm(nil), StatusEOF, nil
	}
	term := s.terms[s.index]
	s.index++
	return term, StatusOK, nil
}

func (s *steppingRCFStream) CurrentInterval() (Interval, error) {
	if len(s.ranges) == 0 {
		return exactRangeFromRational(RationalFromInt64(0)), nil
	}
	if s.index >= len(s.ranges) {
		return s.ranges[len(s.ranges)-1], nil
	}
	return s.ranges[s.index], nil
}

func (s *steppingRCFStream) Range() (Range, error) {
	return s.CurrentInterval()
}

func TestWB_BLFT_UnaryRange_IdentityAfterOneIngest_WithExactTailRange_IsExact22Over7(t *testing.T) {
	engine := newBLFTState(BLFTCoefficients{
		A: big.NewInt(0),
		B: big.NewInt(1),
		C: big.NewInt(0),
		D: big.NewInt(0),
		E: big.NewInt(0),
		F: big.NewInt(0),
		G: big.NewInt(0),
		H: big.NewInt(1),
	})
	engine = engine.IngestUnaryX(PQTerm{
		P: big.NewInt(3),
		Q: big.NewInt(1),
	}).(blftState)

	got, err := engine.UnaryRange(exactRangeFromRational(RationalFromInt64(7)))
	if err != nil {
		t.Fatalf("UnaryRange error = %v", err)
	}

	want := exactRangeFromRational(NewRational(big.NewInt(22), big.NewInt(7)))
	assertExactSameRangeTailWB(t, got, want)
}

func TestWB_PQStreamFromRCF_RangeTracksUpdatedUnderlyingTailRange(t *testing.T) {
	src := &steppingRCFStream{
		terms: []RCFTerm{
			NewRCFTerm(big.NewInt(3)),
			NewRCFTerm(big.NewInt(7)),
		},
		ranges: []Range{
			{
				Lo:     Endpoint{Value: RationalFromInt64(4), Open: true},
				Hi:     Endpoint{Value: RationalFromInt64(3), Open: false},
				Inside: false,
			},
			exactRangeFromRational(RationalFromInt64(7)),
			exactRangeFromRational(RationalFromInt64(0)),
		},
	}

	pq := PQStreamFromRCF(src)

	got0, err := pq.CurrentInterval()
	if err != nil {
		t.Fatalf("initial CurrentInterval error = %v", err)
	}
	if got0.Inside {
		t.Fatal("initial CurrentInterval().Inside = true, want false")
	}

	_, tail, status, err := pq.NextPQ()
	if err != nil {
		t.Fatalf("first NextPQ error = %v", err)
	}
	if status != StatusOK {
		t.Fatalf("first status = %v, want %v", status, StatusOK)
	}

	got1, err := tail.CurrentInterval()
	if err != nil {
		t.Fatalf("tail CurrentInterval error = %v", err)
	}
	want1 := exactRangeFromRational(RationalFromInt64(7))
	assertExactSameRangeTailWB(t, got1, want1)
}

func TestWB_GCF_UnaryIdentity_OverPQStreamFromRCF_WithUpdatingTailRange_EmitsThree(t *testing.T) {
	src := &steppingRCFStream{
		terms: []RCFTerm{
			NewRCFTerm(big.NewInt(3)),
			NewRCFTerm(big.NewInt(7)),
		},
		ranges: []Range{
			{
				Lo:     Endpoint{Value: RationalFromInt64(4), Open: true},
				Hi:     Endpoint{Value: RationalFromInt64(3), Open: false},
				Inside: false,
			},
			exactRangeFromRational(RationalFromInt64(7)),
			exactRangeFromRational(RationalFromInt64(0)),
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

	term, status, err := nextRCFWithTimeoutUnaryTailWB(t, g, time.Second)
	if err != nil {
		t.Fatalf("NextRCF error = %v", err)
	}
	if status != StatusOK {
		t.Fatalf("first status = %v, want %v", status, StatusOK)
	}
	if term.A().Cmp(big.NewInt(3)) != 0 {
		t.Fatalf("first term = %v, want 3", term.A())
	}
}

func nextRCFWithTimeoutUnaryTailWB(t *testing.T, g *GCF, timeout time.Duration) (RCFTerm, Status, error) {
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

func assertExactSameRangeTailWB(t *testing.T, got, want Range) {
	t.Helper()

	if got.Inside != want.Inside {
		t.Fatalf("Inside = %v, want %v", got.Inside, want.Inside)
	}
	if got.Lo.Open != want.Lo.Open {
		t.Fatalf("Lo.Open = %v, want %v", got.Lo.Open, want.Lo.Open)
	}
	if got.Hi.Open != want.Hi.Open {
		t.Fatalf("Hi.Open = %v, want %v", got.Hi.Open, want.Hi.Open)
	}
	if got.Lo.Value.Cmp(want.Lo.Value) != 0 {
		t.Fatalf(
			"Lo = %v/%v, want %v/%v",
			got.Lo.Value.Num(), got.Lo.Value.Den(),
			want.Lo.Value.Num(), want.Lo.Value.Den(),
		)
	}
	if got.Hi.Value.Cmp(want.Hi.Value) != 0 {
		t.Fatalf(
			"Hi = %v/%v, want %v/%v",
			got.Hi.Value.Num(), got.Hi.Value.Den(),
			want.Hi.Value.Num(), want.Hi.Value.Den(),
		)
	}
}

// core/gcf_unary_pqstreamfromrcf_tail_range_wb_test.go v3
