// core/gcf_dlft_bb_test.go v1
package core_test

import (
	"math/big"
	"testing"

	"github.com/egp/gosper-gcf/core"
)

func TestBB_DLFT1_SquareOfFiniteInput(t *testing.T) {
	stream, status := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(1), Q: big.NewInt(1)},
			Range: dlftExactRange(3, 2),
		},
		{
			Term:  core.PQTerm{P: big.NewInt(2), Q: big.NewInt(1)},
			Range: dlftExactRange(2, 1),
		},
	})
	if status != core.StatusOK {
		t.Fatalf("NewFinitePQStream status = %v, want %v", status, core.StatusOK)
	}

	g := core.NewDLFT1(dlftSquareCoeffs(), stream)

	term1, status1 := g.NextRCF()
	if status1 != core.StatusOK {
		t.Fatalf("first status = %v, want %v", status1, core.StatusOK)
	}
	if term1.A().Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("first term = %v, want 2", term1.A())
	}

	term2, status2 := g.NextRCF()
	if status2 != core.StatusOK {
		t.Fatalf("second status = %v, want %v", status2, core.StatusOK)
	}
	if term2.A().Cmp(big.NewInt(4)) != 0 {
		t.Fatalf("second term = %v, want 4", term2.A())
	}

	_, eofStatus := g.NextRCF()
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}
}

func dlftSquareCoeffs() core.DLFTCoefficients {
	return core.DLFTCoefficients{
		A: big.NewInt(1),
		B: big.NewInt(0),
		C: big.NewInt(0),
		D: big.NewInt(0),
		E: big.NewInt(0),
		F: big.NewInt(1),
	}
}

func dlftExactRange(num, den int64) core.Range {
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

// core/gcf_dlft_bb_test.go v1
