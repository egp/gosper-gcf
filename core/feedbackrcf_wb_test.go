// core/feedbackrcf_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

type feedbackRCFStep struct {
	term int64
	num  int64
	den  int64
}

type scriptedFeedbackProducer struct {
	steps []feedbackRCFStep
	calls int
	buf   *feedbackRCFBuffer
}

func newScriptedFeedbackBuffer(steps []feedbackRCFStep) (*feedbackRCFBuffer, *scriptedFeedbackProducer) {
	producer := &scriptedFeedbackProducer{
		steps: append([]feedbackRCFStep(nil), steps...),
	}
	buf := newFeedbackRCFBuffer(producer.ensure)
	producer.buf = buf
	return buf, producer
}

func (p *scriptedFeedbackProducer) ensure(target int) {
	p.calls++

	if len(p.steps) == 0 {
		p.buf.Close()
		return
	}

	step := p.steps[0]
	p.steps = p.steps[1:]

	p.buf.Append(
		NewRCFTerm(big.NewInt(step.term)),
		feedbackExactRange(step.num, step.den),
	)
}

func TestWB_FeedbackRCFBuffer_ReplaysSharedTermsWithoutDuplicateProduction(t *testing.T) {
	buf, producer := newScriptedFeedbackBuffer([]feedbackRCFStep{
		{term: 3, num: 19, den: 6},
		{term: 1, num: 4, den: 1},
	})

	left := buf.Cursor()
	right := buf.Cursor()

	term1, status1 := left.NextRCF()
	if status1 != StatusOK {
		t.Fatalf("left first status = %v, want %v", status1, StatusOK)
	}
	if term1.A().Cmp(big.NewInt(3)) != 0 {
		t.Fatalf("left first term = %v, want 3", term1.A())
	}
	if producer.calls != 1 {
		t.Fatalf("producer calls after left first term = %d, want 1", producer.calls)
	}

	term2, status2 := right.NextRCF()
	if status2 != StatusOK {
		t.Fatalf("right first status = %v, want %v", status2, StatusOK)
	}
	if term2.A().Cmp(big.NewInt(3)) != 0 {
		t.Fatalf("right first term = %v, want 3", term2.A())
	}
	if producer.calls != 1 {
		t.Fatalf("producer calls after right first term = %d, want still 1", producer.calls)
	}

	term3, status3 := right.NextRCF()
	if status3 != StatusOK {
		t.Fatalf("right second status = %v, want %v", status3, StatusOK)
	}
	if term3.A().Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("right second term = %v, want 1", term3.A())
	}
	if producer.calls != 2 {
		t.Fatalf("producer calls after right second term = %d, want 2", producer.calls)
	}
}

func TestWB_FeedbackRCFBuffer_CachesPerPositionRanges(t *testing.T) {
	buf, producer := newScriptedFeedbackBuffer([]feedbackRCFStep{
		{term: 3, num: 19, den: 6},
		{term: 1, num: 4, den: 1},
	})

	lead := buf.Cursor()
	lag := buf.Cursor()

	_, status := lead.NextRCF()
	if status != StatusOK {
		t.Fatalf("lead first status = %v, want %v", status, StatusOK)
	}
	if producer.calls != 1 {
		t.Fatalf("producer calls after lead first term = %d, want 1", producer.calls)
	}

	leadRange := lead.Range()
	assertFeedbackExactRange(t, leadRange, NewRational(big.NewInt(4), big.NewInt(1)), 1)

	if producer.calls != 2 {
		t.Fatalf("producer calls after lead next-range = %d, want 2", producer.calls)
	}

	lagRange := lag.Range()
	assertFeedbackExactRange(t, lagRange, NewRational(big.NewInt(19), big.NewInt(6)), 0)

	if producer.calls != 2 {
		t.Fatalf("producer calls after lag cached-range = %d, want still 2", producer.calls)
	}
}

func TestWB_FeedbackRCFBuffer_IntegratesWithPQStreamFromRCF(t *testing.T) {
	buf, producer := newScriptedFeedbackBuffer([]feedbackRCFStep{
		{term: 3, num: 19, den: 6},
	})

	pq := PQStreamFromRCF(buf.Cursor())

	r0 := pq.Range()
	assertFeedbackExactRange(t, r0, NewRational(big.NewInt(19), big.NewInt(6)), 0)

	term, _, status := pq.NextPQ()
	if status != StatusOK {
		t.Fatalf("status = %v, want %v", status, StatusOK)
	}
	if term.P.Cmp(big.NewInt(3)) != 0 || term.Q.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("PQ term = (%v,%v), want (3,1)", term.P, term.Q)
	}

	if producer.calls != 1 {
		t.Fatalf("producer calls = %d, want 1", producer.calls)
	}
}

func TestWB_FeedbackRCFBuffer_ReplaysEOFCleanly(t *testing.T) {
	buf, producer := newScriptedFeedbackBuffer(nil)

	left := buf.Cursor()
	right := buf.Cursor()

	_, status1 := left.NextRCF()
	if status1 != StatusEOF {
		t.Fatalf("left EOF status = %v, want %v", status1, StatusEOF)
	}
	if producer.calls != 1 {
		t.Fatalf("producer calls after first EOF = %d, want 1", producer.calls)
	}

	_, status2 := right.NextRCF()
	if status2 != StatusEOF {
		t.Fatalf("right EOF status = %v, want %v", status2, StatusEOF)
	}
	if producer.calls != 1 {
		t.Fatalf("producer calls after replayed EOF = %d, want still 1", producer.calls)
	}
}

func feedbackExactRange(num, den int64) Range {
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

func assertFeedbackExactRange(t *testing.T, got Range, want Rational, step int) {
	t.Helper()

	if !got.Inside {
		t.Fatalf("range %d Inside=false, want true", step)
	}
	if got.Lo.Open || got.Hi.Open {
		t.Fatalf("range %d openness wrong, want both closed", step)
	}
	if got.Lo.Value.Cmp(want) != 0 || got.Hi.Value.Cmp(want) != 0 {
		t.Fatalf(
			"range %d = [%v/%v,%v/%v], want exact %v/%v",
			step,
			got.Lo.Value.Num(), got.Lo.Value.Den(),
			got.Hi.Value.Num(), got.Hi.Value.Den(),
			want.Num(), want.Den(),
		)
	}
}

// core/feedbackrcf_wb_test.go v1
