// core/pqstream_procedural_wb_test.go v2
package core

import (
	"math/big"
	"testing"
)

func TestWB_GCFStream_IncrementalValidationTriggersAtBadTerm(t *testing.T) {
	first := FinitePQStep{
		Term: PQTerm{
			P: big.NewInt(1),
			Q: big.NewInt(1),
		},
		Range: testInsideRange(1, 2),
	}

	next := func() (FinitePQStep, ProceduralPQNext, Status) {
		return FinitePQStep{
			Term: PQTerm{
				P: big.NewInt(-2),
				Q: big.NewInt(1),
			},
			Range: testInsideRange(1, 2),
		}, nil, StatusOK
	}

	stream, status := NewProceduralPQStream(first, next)
	if status != StatusOK {
		t.Fatalf("constructor status = %v, want %v", status, StatusOK)
	}

	term1, tail, status1, err := stream.NextPQ()
	if err != nil {
		t.Fatalf("first NextPQ error = %v", err)
	}
	if status1 != StatusOK {
		t.Fatalf("first NextPQ status = %v, want %v", status1, StatusOK)
	}
	if term1.P.Cmp(big.NewInt(1)) != 0 || term1.Q.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("first term = (%v,%v), want (1,1)", term1.P, term1.Q)
	}

	_, _, status2, err := tail.NextPQ()
	if err != nil {
		t.Fatalf("second NextPQ error = %v", err)
	}
	if status2 != StatusInvalidInput {
		t.Fatalf("second NextPQ status = %v, want %v", status2, StatusInvalidInput)
	}
}

// core/pqstream_procedural_wb_test.go v2
