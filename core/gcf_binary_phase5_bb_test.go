// core/gcf_binary_phase5_bb_test.go v2
package core_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/egp/gosper-gcf/core"
)

type binaryPhase5CountingStream struct {
	steps []core.FinitePQStep
	calls *int
}

func (s *binaryPhase5CountingStream) NextPQ() (core.PQTerm, core.PQStream, core.Status, error) {
	if s.calls == nil {
		return core.PQTerm{}, s, core.StatusEOF, fmt.Errorf("binaryPhase5CountingStream.NextPQ: nil calls counter")
	}
	*s.calls++

	if len(s.steps) == 0 {
		return core.PQTerm{
			P: big.NewInt(0),
			Q: big.NewInt(0),
		}, s, core.StatusEOF, nil
	}

	head := s.steps[0]
	tail := &binaryPhase5CountingStream{
		steps: append([]core.FinitePQStep(nil), s.steps[1:]...),
		calls: s.calls,
	}

	return head.Term, tail, core.StatusOK, nil
}

func (s *binaryPhase5CountingStream) CurrentInterval() (core.Interval, error) {
	if len(s.steps) == 0 {
		return core.Interval{}, fmt.Errorf("binaryPhase5CountingStream.CurrentInterval: undefined on EOF PQStream")
	}
	return s.steps[0].Range, nil
}

func (s *binaryPhase5CountingStream) Range() (core.Range, error) {
	return s.CurrentInterval()
}

func TestBB_GCF_BinaryProjectXPassesThroughLeftInput(t *testing.T) {
	x, xStatus := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(3), Q: big.NewInt(1)},
			Range: binaryPhase5ExactRange(19, 5),
		},
		{
			Term:  core.PQTerm{P: big.NewInt(1), Q: big.NewInt(1)},
			Range: binaryPhase5ExactRange(5, 4),
		},
		{
			Term:  core.PQTerm{P: big.NewInt(4), Q: big.NewInt(1)},
			Range: binaryPhase5ExactRange(4, 1),
		},
	})
	if xStatus != core.StatusOK {
		t.Fatalf("left stream status = %v, want %v", xStatus, core.StatusOK)
	}

	y, yStatus := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(2), Q: big.NewInt(1)},
			Range: binaryPhase5ExactRange(2, 1),
		},
	})
	if yStatus != core.StatusOK {
		t.Fatalf("right stream status = %v, want %v", yStatus, core.StatusOK)
	}

	g := core.NewGCF2(binaryPhase5ProjectXCoeffs(), x, y)
	binaryPhase5AssertRCFSequence(t, g, []int64{3, 1, 4})
}

func TestBB_GCF_BinaryProjectYPassesThroughRightInput(t *testing.T) {
	x, xStatus := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(7), Q: big.NewInt(1)},
			Range: binaryPhase5ExactRange(7, 1),
		},
	})
	if xStatus != core.StatusOK {
		t.Fatalf("left stream status = %v, want %v", xStatus, core.StatusOK)
	}

	y, yStatus := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(2), Q: big.NewInt(1)},
			Range: binaryPhase5ExactRange(5, 2),
		},
		{
			Term:  core.PQTerm{P: big.NewInt(2), Q: big.NewInt(1)},
			Range: binaryPhase5ExactRange(2, 1),
		},
	})
	if yStatus != core.StatusOK {
		t.Fatalf("right stream status = %v, want %v", yStatus, core.StatusOK)
	}

	g := core.NewGCF2(binaryPhase5ProjectYCoeffs(), x, y)
	binaryPhase5AssertRCFSequence(t, g, []int64{2, 2})
}

func TestBB_GCF_BinaryCollapseContinuesAfterOneSideEOF(t *testing.T) {
	x, xStatus := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(1), Q: big.NewInt(1)},
			Range: binaryPhase5InsideRange(1, 2),
		},
		{
			Term:  core.PQTerm{P: big.NewInt(1), Q: big.NewInt(1)},
			Range: binaryPhase5ExactRange(2, 1),
		},
		{
			Term:  core.PQTerm{P: big.NewInt(1), Q: big.NewInt(1)},
			Range: binaryPhase5ExactRange(1, 1),
		},
	})
	if xStatus != core.StatusOK {
		t.Fatalf("left stream status = %v, want %v", xStatus, core.StatusOK)
	}

	y, yStatus := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(2), Q: big.NewInt(1)},
			Range: binaryPhase5ExactRange(2, 1),
		},
	})
	if yStatus != core.StatusOK {
		t.Fatalf("right stream status = %v, want %v", yStatus, core.StatusOK)
	}

	g := core.NewGCF2(binaryPhase5AddCoeffs(), x, y)
	binaryPhase5AssertRCFSequence(t, g, []int64{3, 2})
}

func TestBB_GCF_BinaryTieBreakConsumesXFirst(t *testing.T) {
	xCalls := 0
	yCalls := 0

	x := &binaryPhase5CountingStream{
		steps: []core.FinitePQStep{
			{
				Term:  core.PQTerm{P: big.NewInt(3), Q: big.NewInt(1)},
				Range: binaryPhase5InsideRange(3, 4),
			},
			{
				Term:  core.PQTerm{P: big.NewInt(2), Q: big.NewInt(1)},
				Range: binaryPhase5ExactRange(2, 1),
			},
		},
		calls: &xCalls,
	}

	y := &binaryPhase5CountingStream{
		steps: []core.FinitePQStep{
			{
				Term:  core.PQTerm{P: big.NewInt(10), Q: big.NewInt(1)},
				Range: binaryPhase5InsideRange(10, 11),
			},
		},
		calls: &yCalls,
	}

	g := core.NewGCF2(binaryPhase5ProjectXCoeffs(), x, y)

	term, status, err := g.NextRCF()
	if err != nil {
		t.Fatalf("first NextRCF error = %v", err)
	}
	if status != core.StatusOK {
		t.Fatalf("first status = %v, want %v", status, core.StatusOK)
	}
	if term.A().Cmp(big.NewInt(3)) != 0 {
		t.Fatalf("first term = %v, want 3", term.A())
	}

	if xCalls == 0 {
		t.Fatal("left stream was never consumed, want X to be consumed first on tie")
	}
	if yCalls != 0 {
		t.Fatalf("right stream consumed %d times, want 0 before first emission", yCalls)
	}
}

func binaryPhase5AssertRCFSequence(t *testing.T, g *core.GCF, want []int64) {
	t.Helper()

	for i, w := range want {
		got, status, err := g.NextRCF()
		if err != nil {
			t.Fatalf("term %d NextRCF error = %v", i+1, err)
		}
		if status != core.StatusOK {
			t.Fatalf("term %d status = %v, want %v", i+1, status, core.StatusOK)
		}
		if got.A().Cmp(big.NewInt(w)) != 0 {
			t.Fatalf("term %d = %v, want %d", i+1, got.A(), w)
		}
	}

	_, eofStatus, err := g.NextRCF()
	if err != nil {
		t.Fatalf("EOF NextRCF error = %v", err)
	}
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}
}

func binaryPhase5ProjectXCoeffs() core.BLFTCoefficients {
	return core.BLFTCoefficients{
		A: big.NewInt(0),
		B: big.NewInt(1),
		C: big.NewInt(0),
		D: big.NewInt(0),
		E: big.NewInt(0),
		F: big.NewInt(0),
		G: big.NewInt(0),
		H: big.NewInt(1),
	}
}

func binaryPhase5ProjectYCoeffs() core.BLFTCoefficients {
	return core.BLFTCoefficients{
		A: big.NewInt(0),
		B: big.NewInt(0),
		C: big.NewInt(1),
		D: big.NewInt(0),
		E: big.NewInt(0),
		F: big.NewInt(0),
		G: big.NewInt(0),
		H: big.NewInt(1),
	}
}

func binaryPhase5AddCoeffs() core.BLFTCoefficients {
	return core.BLFTCoefficients{
		A: big.NewInt(0),
		B: big.NewInt(1),
		C: big.NewInt(1),
		D: big.NewInt(0),
		E: big.NewInt(0),
		F: big.NewInt(0),
		G: big.NewInt(0),
		H: big.NewInt(1),
	}
}

func binaryPhase5ExactRange(num, den int64) core.Range {
	value := core.NewRational(big.NewInt(num), big.NewInt(den))
	return core.Range{
		Lo: core.Endpoint{
			Value: value,
			Open:  false,
		},
		Hi: core.Endpoint{
			Value: value,
			Open:  false,
		},
		Inside: true,
	}
}

func binaryPhase5InsideRange(lo, hi int64) core.Range {
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

// core/gcf_binary_phase5_bb_test.go v2
