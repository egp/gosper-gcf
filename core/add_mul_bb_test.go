// core/add_mul_bb_test.go v2
package core_test

import (
	"math/big"
	"testing"

	"github.com/egp/gosper-gcf/core"
)

func TestBB_GCF_AddOfFiniteInputs(t *testing.T) {
	x, xStatus := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(1), Q: big.NewInt(1)},
			Range: addMulExactRange(3, 2),
		},
		{
			Term:  core.PQTerm{P: big.NewInt(2), Q: big.NewInt(1)},
			Range: addMulExactRange(2, 1),
		},
	})
	if xStatus != core.StatusOK {
		t.Fatalf("left stream status = %v, want %v", xStatus, core.StatusOK)
	}

	y, yStatus := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(2), Q: big.NewInt(1)},
			Range: addMulExactRange(2, 1),
		},
	})
	if yStatus != core.StatusOK {
		t.Fatalf("right stream status = %v, want %v", yStatus, core.StatusOK)
	}

	g := core.Add(x, y)

	term1, status1, err := g.NextRCF()
	if err != nil {
		t.Fatalf("first NextRCF error = %v", err)
	}
	if status1 != core.StatusOK {
		t.Fatalf("first status = %v, want %v", status1, core.StatusOK)
	}
	if term1.A().Cmp(big.NewInt(3)) != 0 {
		t.Fatalf("first term = %v, want 3", term1.A())
	}

	term2, status2, err := g.NextRCF()
	if err != nil {
		t.Fatalf("second NextRCF error = %v", err)
	}
	if status2 != core.StatusOK {
		t.Fatalf("second status = %v, want %v", status2, core.StatusOK)
	}
	if term2.A().Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("second term = %v, want 2", term2.A())
	}

	_, eofStatus, err := g.NextRCF()
	if err != nil {
		t.Fatalf("EOF NextRCF error = %v", err)
	}
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}
}

func TestBB_GCF_MulOfFiniteInputs(t *testing.T) {
	x, xStatus := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(1), Q: big.NewInt(1)},
			Range: addMulExactRange(3, 2),
		},
		{
			Term:  core.PQTerm{P: big.NewInt(2), Q: big.NewInt(1)},
			Range: addMulExactRange(2, 1),
		},
	})
	if xStatus != core.StatusOK {
		t.Fatalf("left stream status = %v, want %v", xStatus, core.StatusOK)
	}

	y, yStatus := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(2), Q: big.NewInt(1)},
			Range: addMulExactRange(2, 1),
		},
	})
	if yStatus != core.StatusOK {
		t.Fatalf("right stream status = %v, want %v", yStatus, core.StatusOK)
	}

	g := core.Mul(x, y)

	term1, status1, err := g.NextRCF()
	if err != nil {
		t.Fatalf("first NextRCF error = %v", err)
	}
	if status1 != core.StatusOK {
		t.Fatalf("first status = %v, want %v", status1, core.StatusOK)
	}
	if term1.A().Cmp(big.NewInt(3)) != 0 {
		t.Fatalf("first term = %v, want 3", term1.A())
	}

	_, eofStatus, err := g.NextRCF()
	if err != nil {
		t.Fatalf("EOF NextRCF error = %v", err)
	}
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}
}

// new regression coverage for the error-channel migration.
func TestBB_GCF_AddOfFiniteInputs_HasExactRange(t *testing.T) {
	x, xStatus := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(1), Q: big.NewInt(1)},
			Range: addMulExactRange(3, 2),
		},
		{
			Term:  core.PQTerm{P: big.NewInt(2), Q: big.NewInt(1)},
			Range: addMulExactRange(2, 1),
		},
	})
	if xStatus != core.StatusOK {
		t.Fatalf("left stream status = %v, want %v", xStatus, core.StatusOK)
	}

	y, yStatus := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(2), Q: big.NewInt(1)},
			Range: addMulExactRange(2, 1),
		},
	})
	if yStatus != core.StatusOK {
		t.Fatalf("right stream status = %v, want %v", yStatus, core.StatusOK)
	}

	g := core.Add(x, y)

	r, err := g.Range()
	if err != nil {
		t.Fatalf("Range error = %v", err)
	}
	want := core.NewRational(big.NewInt(7), big.NewInt(2))
	if !r.Inside {
		t.Fatal("Inside = false, want true")
	}
	if r.Lo.Value.Cmp(want) != 0 || r.Hi.Value.Cmp(want) != 0 {
		t.Fatalf("Range = [%v/%v,%v/%v], want exact 7/2",
			r.Lo.Value.Num(), r.Lo.Value.Den(),
			r.Hi.Value.Num(), r.Hi.Value.Den(),
		)
	}
}

func addMulExactRange(num, den int64) core.Range {
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

// core/add_mul_bb_test.go v2
