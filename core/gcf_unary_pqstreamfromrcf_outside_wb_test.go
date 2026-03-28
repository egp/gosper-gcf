// core/gcf_unary_pqstreamfromrcf_outside_wb_test.go v3
package core

import (
	"math/big"
	"testing"
	"time"
)

type advancingOutsideRCFStream struct {
	terms  []RCFTerm
	ranges []Range
	index  int
}

func (s *advancingOutsideRCFStream) NextRCF() (RCFTerm, Status) {
	if s.index >= len(s.terms) {
		return NewRCFTerm(nil), StatusEOF
	}
	term := s.terms[s.index]
	s.index++
	return term, StatusOK
}

func (s *advancingOutsideRCFStream) Range() Range {
	if len(s.ranges) == 0 {
		return exactRangeFromRational(RationalFromInt64(0))
	}
	if s.index >= len(s.ranges) {
		return s.ranges[len(s.ranges)-1]
	}
	return s.ranges[s.index]
}

func TestWB_GCF_UnaryIdentity_OverPQStreamFromRCF_WithOutsideRange_DoesNotPanic(t *testing.T) {
	src := &advancingOutsideRCFStream{
		terms: []RCFTerm{
			NewRCFTerm(big.NewInt(3)),
			NewRCFTerm(big.NewInt(7)),
		},
		ranges: []Range{
			{
				Lo:     Endpoint{Value: RationalFromInt64(3), Open: false},
				Hi:     Endpoint{Value: RationalFromInt64(4), Open: true},
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

	term, status := nextRCFWithTimeoutUnaryOutside(t, g, time.Second)
	if status != StatusOK {
		t.Fatalf("first status = %v, want %v", status, StatusOK)
	}
	if term.A().Cmp(big.NewInt(3)) != 0 {
		t.Fatalf("first term = %v, want 3", term.A())
	}
}

func nextRCFWithTimeoutUnaryOutside(t *testing.T, g *GCF, timeout time.Duration) (RCFTerm, Status) {
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

// core/gcf_unary_pqstreamfromrcf_outside_wb_test.go v3
