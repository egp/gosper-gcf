// core/rcf_adapter_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

type fakeRCFStream struct {
	terms []RCFTerm
	rngs  []Range
	next  int
}

func (s *fakeRCFStream) NextRCF() (RCFTerm, Status) {
	if s.next >= len(s.terms) {
		return NewRCFTerm(nil), StatusEOF
	}
	term := s.terms[s.next]
	s.next++
	return term, StatusOK
}

func (s *fakeRCFStream) Range() Range {
	if s.next >= len(s.rngs) {
		panic("Range() on exhausted fakeRCFStream")
	}
	return s.rngs[s.next]
}

func TestWB_PQStreamFromRCF_MapsTermsToQEqualsOne(t *testing.T) {
	rcf := &fakeRCFStream{
		terms: []RCFTerm{
			NewRCFTerm(big.NewInt(3)),
			NewRCFTerm(big.NewInt(1)),
			NewRCFTerm(big.NewInt(4)),
		},
		rngs: []Range{
			exactRangeWB(19, 6),
			exactRangeWB(5, 4),
			exactRangeWB(4, 1),
		},
	}

	pq := PQStreamFromRCF(rcf)

	term1, tail1, status1 := pq.NextPQ()
	if status1 != StatusOK {
		t.Fatalf("first status = %v, want %v", status1, StatusOK)
	}
	if term1.P.Cmp(big.NewInt(3)) != 0 || term1.Q.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("first PQ term = (%v,%v), want (3,1)", term1.P, term1.Q)
	}

	term2, tail2, status2 := tail1.NextPQ()
	if status2 != StatusOK {
		t.Fatalf("second status = %v, want %v", status2, StatusOK)
	}
	if term2.P.Cmp(big.NewInt(1)) != 0 || term2.Q.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("second PQ term = (%v,%v), want (1,1)", term2.P, term2.Q)
	}

	term3, _, status3 := tail2.NextPQ()
	if status3 != StatusOK {
		t.Fatalf("third status = %v, want %v", status3, StatusOK)
	}
	if term3.P.Cmp(big.NewInt(4)) != 0 || term3.Q.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("third PQ term = (%v,%v), want (4,1)", term3.P, term3.Q)
	}
}

func TestWB_PQStreamFromRCF_ForwardsRemainingRange(t *testing.T) {
	rcf := &fakeRCFStream{
		terms: []RCFTerm{
			NewRCFTerm(big.NewInt(3)),
			NewRCFTerm(big.NewInt(1)),
		},
		rngs: []Range{
			exactRangeWB(4, 1),
			exactRangeWB(1, 1),
		},
	}

	pq := PQStreamFromRCF(rcf)

	r0 := pq.Range()
	assertExactRangeWB(t, r0, NewRational(big.NewInt(4), big.NewInt(1)), 0)

	_, tail, status := pq.NextPQ()
	if status != StatusOK {
		t.Fatalf("status = %v, want %v", status, StatusOK)
	}

	r1 := tail.Range()
	assertExactRangeWB(t, r1, NewRational(big.NewInt(1), big.NewInt(1)), 1)
}

func TestWB_PQStreamFromRCF_EOFMapsCleanly(t *testing.T) {
	rcf := &fakeRCFStream{
		terms: nil,
		rngs:  nil,
	}

	pq := PQStreamFromRCF(rcf)

	term, _, status := pq.NextPQ()
	if status != StatusEOF {
		t.Fatalf("status = %v, want %v", status, StatusEOF)
	}
	if term.P.Cmp(big.NewInt(0)) != 0 || term.Q.Cmp(big.NewInt(0)) != 0 {
		t.Fatalf("EOF PQ term = (%v,%v), want (0,0)", term.P, term.Q)
	}
}

func exactRangeWB(num, den int64) Range {
	value := NewRational(big.NewInt(num), big.NewInt(den))
	return Range{
		Lo:     Endpoint{Value: value, Open: false},
		Hi:     Endpoint{Value: value, Open: false},
		Inside: true,
	}
}

func assertExactRangeWB(t *testing.T, got Range, want Rational, step int) {
	t.Helper()

	if !got.Inside {
		t.Fatalf("range %d Inside=false, want true", step)
	}
	if got.Lo.Open || got.Hi.Open {
		t.Fatalf("range %d openness wrong, want both closed", step)
	}
	if got.Lo.Value.Cmp(want) != 0 || got.Hi.Value.Cmp(want) != 0 {
		t.Fatalf("range %d = [%v/%v,%v/%v], want exact %v/%v",
			step,
			got.Lo.Value.Num(), got.Lo.Value.Den(),
			got.Hi.Value.Num(), got.Hi.Value.Den(),
			want.Num(), want.Den(),
		)
	}
}

// core/rcf_adapter_wb_test.go v1
