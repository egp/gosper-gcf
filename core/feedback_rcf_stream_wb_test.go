// core/feedback_rcf_stream_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

func TestWB_FeedbackRCFStream_AppendedTermsAreReadInOrder(t *testing.T) {
	s := newFeedbackRCFStream()

	s.Append(NewRCFTerm(big.NewInt(1)), feedbackRCFExactRange(3, 2))
	s.Append(NewRCFTerm(big.NewInt(2)), feedbackRCFExactRange(2, 1))
	s.Close()

	term1, status1 := s.NextRCF()
	if status1 != StatusOK {
		t.Fatalf("first status = %v, want %v", status1, StatusOK)
	}
	if term1.A().Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("first term = %v, want 1", term1.A())
	}

	term2, status2 := s.NextRCF()
	if status2 != StatusOK {
		t.Fatalf("second status = %v, want %v", status2, StatusOK)
	}
	if term2.A().Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("second term = %v, want 2", term2.A())
	}

	_, status3 := s.NextRCF()
	if status3 != StatusEOF {
		t.Fatalf("third status = %v, want %v", status3, StatusEOF)
	}
}

func TestWB_FeedbackRCFStream_PQStreamFromRCFMapsQToOne(t *testing.T) {
	s := newFeedbackRCFStream()

	s.Append(NewRCFTerm(big.NewInt(7)), feedbackRCFExactRange(7, 1))
	s.Close()

	pq := PQStreamFromRCF(s)

	gotRange := pq.Range()
	assertFeedbackRCFExactRange(t, gotRange, NewRational(big.NewInt(7), big.NewInt(1)))

	term, _, status := pq.NextPQ()
	if status != StatusOK {
		t.Fatalf("status = %v, want %v", status, StatusOK)
	}
	if term.P.Cmp(big.NewInt(7)) != 0 {
		t.Fatalf("term.P = %v, want 7", term.P)
	}
	if term.Q.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("term.Q = %v, want 1", term.Q)
	}

	_, _, eofStatus := pq.NextPQ()
	if eofStatus != StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, StatusEOF)
	}
}

func TestWB_FeedbackRCFStream_RangeTracksRemainingSuffix(t *testing.T) {
	s := newFeedbackRCFStream()

	s.Append(NewRCFTerm(big.NewInt(1)), feedbackRCFExactRange(3, 2))
	s.Append(NewRCFTerm(big.NewInt(2)), feedbackRCFExactRange(2, 1))
	s.Close()

	initial := s.Range()
	assertFeedbackRCFExactRange(t, initial, NewRational(big.NewInt(3), big.NewInt(2)))

	_, status := s.NextRCF()
	if status != StatusOK {
		t.Fatalf("first status = %v, want %v", status, StatusOK)
	}

	remaining := s.Range()
	assertFeedbackRCFExactRange(t, remaining, NewRational(big.NewInt(2), big.NewInt(1)))
}

func feedbackRCFExactRange(num, den int64) Range {
	value := NewRational(big.NewInt(num), big.NewInt(den))
	return Range{
		Lo: Endpoint{
			Value: value,
			Open:  false,
		},
		Hi: Endpoint{
			Value: value,
			Open:  false,
		},
		Inside: true,
	}
}

func assertFeedbackRCFExactRange(t *testing.T, got Range, want Rational) {
	t.Helper()

	if !got.Inside {
		t.Fatalf("Inside = false, want true")
	}
	if got.Lo.Open || got.Hi.Open {
		t.Fatalf("range openness = (%v,%v), want both closed", got.Lo.Open, got.Hi.Open)
	}
	if got.Lo.Value.Cmp(want) != 0 || got.Hi.Value.Cmp(want) != 0 {
		t.Fatalf(
			"range = [%v/%v,%v/%v], want exact %v/%v",
			got.Lo.Value.Num(), got.Lo.Value.Den(),
			got.Hi.Value.Num(), got.Hi.Value.Den(),
			want.Num(), want.Den(),
		)
	}
}

// core/feedback_rcf_stream_wb_test.go v1
