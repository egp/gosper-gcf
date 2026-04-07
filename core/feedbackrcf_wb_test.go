// core/feedbackrcf_wb_test.go v2
package core

import (
	"math/big"
	"testing"
)

func TestWB_FeedbackRCFBuffer_ForkedCursorsReadSamePrefixIndependently(t *testing.T) {
	buf := newFeedbackRCFBuffer(func(target int) {
		_ = target
	})

	r0 := exactRangeFromRational(NewRational(big.NewInt(3), big.NewInt(1)))
	r1 := exactRangeFromRational(NewRational(big.NewInt(7), big.NewInt(2)))

	if err := buf.Append(NewRCFTerm(big.NewInt(3)), r0); err != nil {
		t.Fatalf("Append term 0 error = %v", err)
	}
	if err := buf.Append(NewRCFTerm(big.NewInt(2)), r1); err != nil {
		t.Fatalf("Append term 1 error = %v", err)
	}
	if err := buf.Close(); err != nil {
		t.Fatalf("Close error = %v", err)
	}

	left := buf.Cursor()
	right := buf.Cursor()

	t0, s0, err := left.NextRCF()
	if err != nil {
		t.Fatalf("left first NextRCF error = %v", err)
	}
	if s0 != StatusOK {
		t.Fatalf("left first status = %v, want %v", s0, StatusOK)
	}
	if t0.A().Cmp(big.NewInt(3)) != 0 {
		t.Fatalf("left first term = %v, want 3", t0.A())
	}

	t1, s1, err := right.NextRCF()
	if err != nil {
		t.Fatalf("right first NextRCF error = %v", err)
	}
	if s1 != StatusOK {
		t.Fatalf("right first status = %v, want %v", s1, StatusOK)
	}
	if t1.A().Cmp(big.NewInt(3)) != 0 {
		t.Fatalf("right first term = %v, want 3", t1.A())
	}

	t2, s2, err := right.NextRCF()
	if err != nil {
		t.Fatalf("right second NextRCF error = %v", err)
	}
	if s2 != StatusOK {
		t.Fatalf("right second status = %v, want %v", s2, StatusOK)
	}
	if t2.A().Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("right second term = %v, want 2", t2.A())
	}

	t3, s3, err := left.NextRCF()
	if err != nil {
		t.Fatalf("left second NextRCF error = %v", err)
	}
	if s3 != StatusOK {
		t.Fatalf("left second status = %v, want %v", s3, StatusOK)
	}
	if t3.A().Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("left second term = %v, want 2", t3.A())
	}
}

func TestWB_FeedbackRCFBuffer_LaggingCursorSeesEarlierRemainingRange(t *testing.T) {
	buf := newFeedbackRCFBuffer(func(target int) {
		_ = target
	})

	r0 := exactRangeFromRational(NewRational(big.NewInt(5), big.NewInt(2)))
	r1 := exactRangeFromRational(NewRational(big.NewInt(7), big.NewInt(3)))

	if err := buf.Append(NewRCFTerm(big.NewInt(2)), r0); err != nil {
		t.Fatalf("Append term 0 error = %v", err)
	}
	if err := buf.Append(NewRCFTerm(big.NewInt(3)), r1); err != nil {
		t.Fatalf("Append term 1 error = %v", err)
	}
	if err := buf.Close(); err != nil {
		t.Fatalf("Close error = %v", err)
	}

	lead := buf.Cursor()
	lag := buf.Cursor()

	_, status, err := lead.NextRCF()
	if err != nil {
		t.Fatalf("lead NextRCF error = %v", err)
	}
	if status != StatusOK {
		t.Fatalf("lead status = %v, want %v", status, StatusOK)
	}

	leadRange, err := lead.Range()
	if err != nil {
		t.Fatalf("lead Range error = %v", err)
	}
	assertFeedbackExactRange(t, leadRange, r1, "lead")

	lagRange, err := lag.Range()
	if err != nil {
		t.Fatalf("lag Range error = %v", err)
	}
	assertFeedbackExactRange(t, lagRange, r0, "lag")
}

func TestWB_FeedbackRCFBuffer_PQAdapterTracksRemainingRange(t *testing.T) {
	buf := newFeedbackRCFBuffer(func(target int) {
		_ = target
	})

	r0 := exactRangeFromRational(NewRational(big.NewInt(8), big.NewInt(3)))
	r1 := exactRangeFromRational(NewRational(big.NewInt(5), big.NewInt(2)))

	if err := buf.Append(NewRCFTerm(big.NewInt(2)), r0); err != nil {
		t.Fatalf("Append term 0 error = %v", err)
	}
	if err := buf.Append(NewRCFTerm(big.NewInt(1)), r1); err != nil {
		t.Fatalf("Append term 1 error = %v", err)
	}
	if err := buf.Close(); err != nil {
		t.Fatalf("Close error = %v", err)
	}

	pq := PQStreamFromRCF(buf.Cursor())

	got0, err := pq.Range()
	if err != nil {
		t.Fatalf("initial PQ Range error = %v", err)
	}
	assertFeedbackExactRange(t, got0, r0, "initial pq")

	term, tail, status, err := pq.NextPQ()
	if err != nil {
		t.Fatalf("PQ NextPQ error = %v", err)
	}
	if status != StatusOK {
		t.Fatalf("PQ status = %v, want %v", status, StatusOK)
	}
	if term.P.Cmp(big.NewInt(2)) != 0 || term.Q.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("PQ term = (%v,%v), want (2,1)", term.P, term.Q)
	}

	got1, err := tail.Range()
	if err != nil {
		t.Fatalf("tail PQ Range error = %v", err)
	}
	assertFeedbackExactRange(t, got1, r1, "tail pq")
}

func TestWB_FeedbackRCFBuffer_PQAdapterEOFMapsCleanly(t *testing.T) {
	buf := newFeedbackRCFBuffer(func(target int) {
		_ = target
	})
	if err := buf.Close(); err != nil {
		t.Fatalf("Close error = %v", err)
	}

	pq := PQStreamFromRCF(buf.Cursor())

	term, _, status, err := pq.NextPQ()
	if err != nil {
		t.Fatalf("NextPQ error = %v", err)
	}
	if status != StatusEOF {
		t.Fatalf("status = %v, want %v", status, StatusEOF)
	}
	if term.P.Cmp(big.NewInt(0)) != 0 || term.Q.Cmp(big.NewInt(0)) != 0 {
		t.Fatalf("EOF term = (%v,%v), want (0,0)", term.P, term.Q)
	}
}

func assertFeedbackExactRange(t *testing.T, got Range, want Range, label string) {
	t.Helper()

	if got.Inside != want.Inside {
		t.Fatalf("%s Inside = %v, want %v", label, got.Inside, want.Inside)
	}
	if got.Lo.Open != want.Lo.Open || got.Hi.Open != want.Hi.Open {
		t.Fatalf("%s openness = (%v,%v), want (%v,%v)", label, got.Lo.Open, got.Hi.Open, want.Lo.Open, want.Hi.Open)
	}
	if got.Lo.Value.Cmp(want.Lo.Value) != 0 || got.Hi.Value.Cmp(want.Hi.Value) != 0 {
		t.Fatalf(
			"%s range = [%v/%v,%v/%v], want [%v/%v,%v/%v]",
			label,
			got.Lo.Value.Num(), got.Lo.Value.Den(),
			got.Hi.Value.Num(), got.Hi.Value.Den(),
			want.Lo.Value.Num(), want.Lo.Value.Den(),
			want.Hi.Value.Num(), want.Hi.Value.Den(),
		)
	}
}

// core/feedbackrcf_wb_test.go v2
