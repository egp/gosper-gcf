// core/pqstream_wb_test.go v2
package core

import (
	"errors"
	"math/big"
	"testing"
)

func TestWB_GCFStream_FirstTermMayBeNegative(t *testing.T) {
	stream, status := NewFinitePQStream([]FinitePQStep{
		{
			Term: PQTerm{
				P: big.NewInt(-3),
				Q: big.NewInt(-5),
			},
			Range: testInsideRange(0, 1),
		},
	})
	if status != StatusOK {
		t.Fatalf("NewFinitePQStream status = %v, want %v", status, StatusOK)
	}
	if stream == nil {
		t.Fatal("stream is nil")
	}

	term, _, nextStatus, err := stream.NextPQ()
	if err != nil {
		t.Fatalf("NextPQ error = %v", err)
	}
	if nextStatus != StatusOK {
		t.Fatalf("NextPQ status = %v, want %v", nextStatus, StatusOK)
	}
	if term.P.Cmp(big.NewInt(-3)) != 0 {
		t.Fatalf("first P = %v, want -3", term.P)
	}
	if term.Q.Cmp(big.NewInt(-5)) != 0 {
		t.Fatalf("first Q = %v, want -5", term.Q)
	}
}

func TestWB_GCFStream_LaterNegativePRejected(t *testing.T) {
	stream, status := NewFinitePQStream([]FinitePQStep{
		{
			Term:  PQTerm{P: big.NewInt(1), Q: big.NewInt(1)},
			Range: testInsideRange(1, 2),
		},
		{
			Term:  PQTerm{P: big.NewInt(-2), Q: big.NewInt(1)},
			Range: testInsideRange(1, 2),
		},
	})
	if status != StatusInvalidInput {
		t.Fatalf("status = %v, want %v", status, StatusInvalidInput)
	}
	if stream != nil {
		t.Fatal("stream != nil for invalid input")
	}
}

func TestWB_GCFStream_LaterNegativeQRejected(t *testing.T) {
	stream, status := NewFinitePQStream([]FinitePQStep{
		{
			Term:  PQTerm{P: big.NewInt(1), Q: big.NewInt(1)},
			Range: testInsideRange(1, 2),
		},
		{
			Term:  PQTerm{P: big.NewInt(2), Q: big.NewInt(-1)},
			Range: testInsideRange(1, 2),
		},
	})
	if status != StatusInvalidInput {
		t.Fatalf("status = %v, want %v", status, StatusInvalidInput)
	}
	if stream != nil {
		t.Fatal("stream != nil for invalid input")
	}
}

func TestWB_GCFStream_ZeroQRejected(t *testing.T) {
	stream, status := NewFinitePQStream([]FinitePQStep{
		{
			Term:  PQTerm{P: big.NewInt(1), Q: big.NewInt(0)},
			Range: testInsideRange(1, 2),
		},
	})
	if status != StatusInvalidInput {
		t.Fatalf("status = %v, want %v", status, StatusInvalidInput)
	}
	if stream != nil {
		t.Fatal("stream != nil for invalid input")
	}
}

func TestWB_GCFStream_EOFReturnsStatusEOF(t *testing.T) {
	stream, status := NewFinitePQStream(nil)
	if status != StatusOK {
		t.Fatalf("constructor status = %v, want %v", status, StatusOK)
	}

	_, tail, nextStatus, err := stream.NextPQ()
	if err != nil {
		t.Fatalf("NextPQ error = %v", err)
	}
	if nextStatus != StatusEOF {
		t.Fatalf("NextPQ status = %v, want %v", nextStatus, StatusEOF)
	}
	if tail == nil {
		t.Fatal("EOF tail is nil, want reusable EOF stream")
	}
}

func TestWB_GCFStream_TailAdvancesCorrectly(t *testing.T) {
	stream, status := NewFinitePQStream([]FinitePQStep{
		{
			Term:  PQTerm{P: big.NewInt(1), Q: big.NewInt(1)},
			Range: testInsideRange(1, 2),
		},
		{
			Term:  PQTerm{P: big.NewInt(2), Q: big.NewInt(1)},
			Range: testInsideRange(1, 2),
		},
	})
	if status != StatusOK {
		t.Fatalf("constructor status = %v, want %v", status, StatusOK)
	}

	first, tail, firstStatus, err := stream.NextPQ()
	if err != nil {
		t.Fatalf("first NextPQ error = %v", err)
	}
	if firstStatus != StatusOK {
		t.Fatalf("first NextPQ status = %v, want %v", firstStatus, StatusOK)
	}

	second, tail2, secondStatus, err := tail.NextPQ()
	if err != nil {
		t.Fatalf("second NextPQ error = %v", err)
	}
	if secondStatus != StatusOK {
		t.Fatalf("second NextPQ status = %v, want %v", secondStatus, StatusOK)
	}

	if first.P.Cmp(big.NewInt(1)) != 0 || first.Q.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("first term = (%v,%v), want (1,1)", first.P, first.Q)
	}
	if second.P.Cmp(big.NewInt(2)) != 0 || second.Q.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("second term = (%v,%v), want (2,1)", second.P, second.Q)
	}

	_, _, eofStatus, err := tail2.NextPQ()
	if err != nil {
		t.Fatalf("EOF NextPQ error = %v", err)
	}
	if eofStatus != StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, StatusEOF)
	}
}

func TestWB_GCFStream_RangeAvailableOnLiveStream(t *testing.T) {
	stream, status := NewFinitePQStream([]FinitePQStep{
		{
			Term:  PQTerm{P: big.NewInt(1), Q: big.NewInt(1)},
			Range: testInsideRange(3, 5),
		},
	})
	if status != StatusOK {
		t.Fatalf("constructor status = %v, want %v", status, StatusOK)
	}

	r, err := stream.Range()
	if err != nil {
		t.Fatalf("Range error = %v", err)
	}
	if !r.Inside {
		t.Fatal("Range.Inside = false, want true")
	}
	if r.Lo.Value.Cmp(RationalFromInt64(3)) != 0 {
		t.Fatalf("Lo = %v/%v, want 3/1", r.Lo.Value.Num(), r.Lo.Value.Den())
	}
	if r.Hi.Value.Cmp(RationalFromInt64(5)) != 0 {
		t.Fatalf("Hi = %v/%v, want 5/1", r.Hi.Value.Num(), r.Hi.Value.Den())
	}
}

func TestWB_GCFStream_EOFRangeReturnsTypedError(t *testing.T) {
	stream, status := NewFinitePQStream(nil)
	if status != StatusOK {
		t.Fatalf("constructor status = %v, want %v", status, StatusOK)
	}

	_, err := stream.Range()
	if !errors.Is(err, ErrUndefinedRangeOnEOFStream) {
		t.Fatalf("Range error = %v, want ErrUndefinedRangeOnEOFStream", err)
	}
}

func testInsideRange(lo, hi int64) Range {
	return Range{
		Lo: Endpoint{
			Value: RationalFromInt64(lo),
			Open:  false,
		},
		Hi: Endpoint{
			Value: RationalFromInt64(hi),
			Open:  false,
		},
		Inside: true,
	}
}

// core/pqstream_wb_test.go v2
