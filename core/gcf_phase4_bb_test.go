// core/gcf_phase4_bb_test.go v1
package core_test

import (
	"math/big"
	"testing"

	"github.com/egp/gosper-gcf/core"
)

func TestBB_BLFT_UnaryIdentityPassesThroughRegularInput(t *testing.T) {
	stream, status := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(3), Q: big.NewInt(1)},
			Range: phase4InsideRange(3, 4),
		},
		{
			Term:  core.PQTerm{P: big.NewInt(1), Q: big.NewInt(1)},
			Range: phase4InsideRange(1, 2),
		},
		{
			Term:  core.PQTerm{P: big.NewInt(4), Q: big.NewInt(1)},
			Range: phase4InsideRange(4, 5),
		},
	})
	if status != core.StatusOK {
		t.Fatalf("NewFinitePQStream status = %v, want %v", status, core.StatusOK)
	}

	g := core.NewGCF1(phase4IdentityUnaryX(), stream)

	phase4AssertRCFSequence(t, g, []int64{3, 1, 4})
}

func TestBB_BLFT_GeneralizedInputHandlesQNotOne(t *testing.T) {
	stream, status := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(1), Q: big.NewInt(2)},
			Range: phase4InsideRange(1, 2),
		},
		{
			Term:  core.PQTerm{P: big.NewInt(3), Q: big.NewInt(1)},
			Range: phase4InsideRange(3, 4),
		},
	})
	if status != core.StatusOK {
		t.Fatalf("NewFinitePQStream status = %v, want %v", status, core.StatusOK)
	}

	g := core.NewGCF1(phase4IdentityUnaryX(), stream)

	phase4AssertRCFSequence(t, g, []int64{1, 1, 2})
}

func TestBB_BLFT_CollapseContinuesAfterEOF(t *testing.T) {
	stream, status := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(1), Q: big.NewInt(1)},
			Range: phase4InsideRange(1, 2),
		},
		{
			Term:  core.PQTerm{P: big.NewInt(1), Q: big.NewInt(1)},
			Range: phase4InsideRange(1, 2),
		},
	})
	if status != core.StatusOK {
		t.Fatalf("NewFinitePQStream status = %v, want %v", status, core.StatusOK)
	}

	g := core.NewGCF1(phase4IdentityUnaryX(), stream)

	phase4AssertRCFSequence(t, g, []int64{2})
}

func phase4AssertRCFSequence(t *testing.T, g *core.GCF, want []int64) {
	t.Helper()

	for i, w := range want {
		got, status := g.NextRCF()
		if status != core.StatusOK {
			t.Fatalf("term %d status = %v, want %v", i+1, status, core.StatusOK)
		}
		if got.A().Cmp(big.NewInt(w)) != 0 {
			t.Fatalf("term %d = %v, want %d", i+1, got.A(), w)
		}
	}

	_, eofStatus := g.NextRCF()
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}
}

func phase4IdentityUnaryX() core.TransformCoefficients {
	return core.TransformCoefficients{
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

func phase4InsideRange(lo, hi int64) core.Range {
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

// core/gcf_phase4_bb_test.go v1
