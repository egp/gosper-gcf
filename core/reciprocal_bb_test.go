// core/reciprocal_bb_test.go v2
package core_test

import (
	"math/big"
	"testing"

	"github.com/egp/gosper-gcf/core"
)

func TestBB_GCF_ReciprocalOfFiniteInput(t *testing.T) {
	stream, status := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(1), Q: big.NewInt(1)},
			Range: reciprocalExactRange(3, 2),
		},
		{
			Term:  core.PQTerm{P: big.NewInt(2), Q: big.NewInt(1)},
			Range: reciprocalExactRange(2, 1),
		},
	})
	if status != core.StatusOK {
		t.Fatalf("NewFinitePQStream status = %v, want %v", status, core.StatusOK)
	}

	g := core.Reciprocal(stream)

	term1, status1, err := g.NextRCF()
	if err != nil {
		t.Fatalf("first NextRCF error = %v", err)
	}
	if status1 != core.StatusOK {
		t.Fatalf("first status = %v, want %v", status1, core.StatusOK)
	}
	if term1.A().Cmp(big.NewInt(0)) != 0 {
		t.Fatalf("first term = %v, want 0", term1.A())
	}

	term2, status2, err := g.NextRCF()
	if err != nil {
		t.Fatalf("second NextRCF error = %v", err)
	}
	if status2 != core.StatusOK {
		t.Fatalf("second status = %v, want %v", status2, core.StatusOK)
	}
	if term2.A().Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("second term = %v, want 1", term2.A())
	}

	term3, status3, err := g.NextRCF()
	if err != nil {
		t.Fatalf("third NextRCF error = %v", err)
	}
	if status3 != core.StatusOK {
		t.Fatalf("third status = %v, want %v", status3, core.StatusOK)
	}
	if term3.A().Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("third term = %v, want 2", term3.A())
	}

	_, eofStatus, err := g.NextRCF()
	if err != nil {
		t.Fatalf("EOF NextRCF error = %v", err)
	}
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}
}

func reciprocalExactRange(num, den int64) core.Range {
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

// core/reciprocal_bb_test.go v2
