// core/pqstream_bb_test.go v2
package core_test

import (
	"math/big"
	"testing"

	"github.com/egp/gosper-gcf/core"
)

func TestBB_GCFStream_ReadFiniteSequenceThenEOF(t *testing.T) {
	stream, status := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(1), Q: big.NewInt(1)},
			Range: pqStreamInsideRange(1, 2),
		},
		{
			Term:  core.PQTerm{P: big.NewInt(2), Q: big.NewInt(1)},
			Range: pqStreamInsideRange(1, 2),
		},
	})
	if status != core.StatusOK {
		t.Fatalf("constructor status = %v, want %v", status, core.StatusOK)
	}

	term1, tail1, status1, err := stream.NextPQ()
	if err != nil {
		t.Fatalf("first NextPQ error = %v", err)
	}
	if status1 != core.StatusOK {
		t.Fatalf("first status = %v, want %v", status1, core.StatusOK)
	}
	if term1.P.Cmp(big.NewInt(1)) != 0 || term1.Q.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("first term = (%v,%v), want (1,1)", term1.P, term1.Q)
	}

	term2, tail2, status2, err := tail1.NextPQ()
	if err != nil {
		t.Fatalf("second NextPQ error = %v", err)
	}
	if status2 != core.StatusOK {
		t.Fatalf("second status = %v, want %v", status2, core.StatusOK)
	}
	if term2.P.Cmp(big.NewInt(2)) != 0 || term2.Q.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("second term = (%v,%v), want (2,1)", term2.P, term2.Q)
	}

	_, _, eofStatus, err := tail2.NextPQ()
	if err != nil {
		t.Fatalf("EOF NextPQ error = %v", err)
	}
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}
}

func TestBB_GCFStream_FirstMalformedTermFails(t *testing.T) {
	stream, status := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(1), Q: big.NewInt(0)},
			Range: pqStreamInsideRange(1, 2),
		},
	})

	if status != core.StatusInvalidInput {
		t.Fatalf("status = %v, want %v", status, core.StatusInvalidInput)
	}
	if stream != nil {
		t.Fatal("stream != nil for invalid input")
	}
}

func TestBB_GCFStream_LaterMalformedTermFails(t *testing.T) {
	stream, status := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(1), Q: big.NewInt(1)},
			Range: pqStreamInsideRange(1, 2),
		},
		{
			Term:  core.PQTerm{P: big.NewInt(0), Q: big.NewInt(1)},
			Range: pqStreamInsideRange(1, 2),
		},
	})

	if status != core.StatusInvalidInput {
		t.Fatalf("status = %v, want %v", status, core.StatusInvalidInput)
	}
	if stream != nil {
		t.Fatal("stream != nil for invalid input")
	}
}

func TestBB_GCFStream_RangeTracksReturnedTail(t *testing.T) {
	stream, status := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(1), Q: big.NewInt(1)},
			Range: pqStreamInsideRange(3, 4),
		},
		{
			Term:  core.PQTerm{P: big.NewInt(2), Q: big.NewInt(1)},
			Range: pqStreamInsideRange(7, 9),
		},
	})
	if status != core.StatusOK {
		t.Fatalf("constructor status = %v, want %v", status, core.StatusOK)
	}

	r1, err := stream.Range()
	if err != nil {
		t.Fatalf("first Range error = %v", err)
	}
	if r1.Lo.Value.Cmp(core.RationalFromInt64(3)) != 0 || r1.Hi.Value.Cmp(core.RationalFromInt64(4)) != 0 {
		t.Fatalf("first range = [%v/%v, %v/%v], want [3/1, 4/1]",
			r1.Lo.Value.Num(), r1.Lo.Value.Den(), r1.Hi.Value.Num(), r1.Hi.Value.Den())
	}

	_, tail, nextStatus, err := stream.NextPQ()
	if err != nil {
		t.Fatalf("NextPQ error = %v", err)
	}
	if nextStatus != core.StatusOK {
		t.Fatalf("NextPQ status = %v, want %v", nextStatus, core.StatusOK)
	}

	r2, err := tail.Range()
	if err != nil {
		t.Fatalf("tail Range error = %v", err)
	}
	if r2.Lo.Value.Cmp(core.RationalFromInt64(7)) != 0 || r2.Hi.Value.Cmp(core.RationalFromInt64(9)) != 0 {
		t.Fatalf("tail range = [%v/%v, %v/%v], want [7/1, 9/1]",
			r2.Lo.Value.Num(), r2.Lo.Value.Den(), r2.Hi.Value.Num(), r2.Hi.Value.Den())
	}
}

func pqStreamInsideRange(lo, hi int64) core.Range {
	return core.Range{
		Lo: core.Endpoint{
			Value: core.RationalFromInt64(lo),
			Open:  false,
		},
		Hi: core.Endpoint{
			Value: core.RationalFromInt64(hi),
			Open:  false,
		},
		Inside: true,
	}
}

// core/pqstream_bb_test.go v2
