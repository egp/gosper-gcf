// core/sub_div_bb_test.go v1
package core_test

import (
	"math/big"
	"testing"

	"github.com/egp/gosper-gcf/core"
)

func TestBB_GCF_SubOfFiniteInputs(t *testing.T) {
	x, xStatus := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(1), Q: big.NewInt(1)},
			Range: subDivExactRange(3, 2),
		},
		{
			Term:  core.PQTerm{P: big.NewInt(2), Q: big.NewInt(1)},
			Range: subDivExactRange(2, 1),
		},
	})
	if xStatus != core.StatusOK {
		t.Fatalf("left stream status = %v, want %v", xStatus, core.StatusOK)
	}

	y, yStatus := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(2), Q: big.NewInt(1)},
			Range: subDivExactRange(2, 1),
		},
	})
	if yStatus != core.StatusOK {
		t.Fatalf("right stream status = %v, want %v", yStatus, core.StatusOK)
	}

	g := core.Sub(x, y)

	term1, status1 := g.NextRCF()
	if status1 != core.StatusOK {
		t.Fatalf("first status = %v, want %v", status1, core.StatusOK)
	}
	if term1.A().Cmp(big.NewInt(-1)) != 0 {
		t.Fatalf("first term = %v, want -1", term1.A())
	}

	term2, status2 := g.NextRCF()
	if status2 != core.StatusOK {
		t.Fatalf("second status = %v, want %v", status2, core.StatusOK)
	}
	if term2.A().Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("second term = %v, want 2", term2.A())
	}

	_, eofStatus := g.NextRCF()
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}
}

func TestBB_GCF_DivOfFiniteInputs(t *testing.T) {
	x, xStatus := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(1), Q: big.NewInt(1)},
			Range: subDivExactRange(3, 2),
		},
		{
			Term:  core.PQTerm{P: big.NewInt(2), Q: big.NewInt(1)},
			Range: subDivExactRange(2, 1),
		},
	})
	if xStatus != core.StatusOK {
		t.Fatalf("left stream status = %v, want %v", xStatus, core.StatusOK)
	}

	y, yStatus := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(2), Q: big.NewInt(1)},
			Range: subDivExactRange(2, 1),
		},
	})
	if yStatus != core.StatusOK {
		t.Fatalf("right stream status = %v, want %v", yStatus, core.StatusOK)
	}

	g := core.Div(x, y)

	term1, status1 := g.NextRCF()
	if status1 != core.StatusOK {
		t.Fatalf("first status = %v, want %v", status1, core.StatusOK)
	}
	if term1.A().Cmp(big.NewInt(0)) != 0 {
		t.Fatalf("first term = %v, want 0", term1.A())
	}

	term2, status2 := g.NextRCF()
	if status2 != core.StatusOK {
		t.Fatalf("second status = %v, want %v", status2, core.StatusOK)
	}
	if term2.A().Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("second term = %v, want 1", term2.A())
	}

	term3, status3 := g.NextRCF()
	if status3 != core.StatusOK {
		t.Fatalf("third status = %v, want %v", status3, core.StatusOK)
	}
	if term3.A().Cmp(big.NewInt(3)) != 0 {
		t.Fatalf("third term = %v, want 3", term3.A())
	}

	_, eofStatus := g.NextRCF()
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}
}

func subDivExactRange(num, den int64) core.Range {
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

// core/sub_div_bb_test.go v1
