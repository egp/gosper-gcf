// core/pqstream_procedural_bb_test.go v2
package core_test

import (
	"math/big"
	"testing"

	"github.com/egp/gosper-gcf/core"
)

func TestBB_GCFStream_ProceduralStreamCanFailLate(t *testing.T) {
	first := core.FinitePQStep{
		Term: core.PQTerm{
			P: big.NewInt(1),
			Q: big.NewInt(1),
		},
		Range: insideRange(1, 2),
	}

	next := func() (core.FinitePQStep, core.ProceduralPQNext, core.Status) {
		return core.FinitePQStep{
			Term: core.PQTerm{
				P: big.NewInt(2),
				Q: big.NewInt(0),
			},
			Range: insideRange(1, 2),
		}, nil, core.StatusOK
	}

	stream, status := core.NewProceduralPQStream(first, next)
	if status != core.StatusOK {
		t.Fatalf("constructor status = %v, want %v", status, core.StatusOK)
	}

	term1, tail, status1, err := stream.NextPQ()
	if err != nil {
		t.Fatalf("first NextPQ error = %v", err)
	}
	if status1 != core.StatusOK {
		t.Fatalf("first NextPQ status = %v, want %v", status1, core.StatusOK)
	}
	if term1.P.Cmp(big.NewInt(1)) != 0 || term1.Q.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("first term = (%v,%v), want (1,1)", term1.P, term1.Q)
	}

	_, _, status2, err := tail.NextPQ()
	if err != nil {
		t.Fatalf("second NextPQ error = %v", err)
	}
	if status2 != core.StatusInvalidInput {
		t.Fatalf("second NextPQ status = %v, want %v", status2, core.StatusInvalidInput)
	}
}

// core/pqstream_procedural_bb_test.go v2
