// core/gcf_helpers_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

func TestWB_BLFTCoefficientsFromState_PreservesAllSlots(t *testing.T) {
	state := blftState{
		A: big.NewInt(2),
		B: big.NewInt(3),
		C: big.NewInt(5),
		D: big.NewInt(7),
		E: big.NewInt(11),
		F: big.NewInt(13),
		G: big.NewInt(17),
		H: big.NewInt(19),
	}

	got := blftCoefficientsFromState(state)

	if got.A.Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("A = %v, want 2", got.A)
	}
	if got.B.Cmp(big.NewInt(3)) != 0 {
		t.Fatalf("B = %v, want 3", got.B)
	}
	if got.C.Cmp(big.NewInt(5)) != 0 {
		t.Fatalf("C = %v, want 5", got.C)
	}
	if got.D.Cmp(big.NewInt(7)) != 0 {
		t.Fatalf("D = %v, want 7", got.D)
	}
	if got.E.Cmp(big.NewInt(11)) != 0 {
		t.Fatalf("E = %v, want 11", got.E)
	}
	if got.F.Cmp(big.NewInt(13)) != 0 {
		t.Fatalf("F = %v, want 13", got.F)
	}
	if got.G.Cmp(big.NewInt(17)) != 0 {
		t.Fatalf("G = %v, want 17", got.G)
	}
	if got.H.Cmp(big.NewInt(19)) != 0 {
		t.Fatalf("H = %v, want 19", got.H)
	}
}

func TestWB_CollapseIndependentOfXToUnary_MapsCDGHIntoUnarySlots1(t *testing.T) {
	state := blftState{
		A: big.NewInt(0),
		B: big.NewInt(0),
		C: big.NewInt(11),
		D: big.NewInt(13),
		E: big.NewInt(0),
		F: big.NewInt(0),
		G: big.NewInt(17),
		H: big.NewInt(19),
	}

	got := collapseIndependentOfXToUnary(state)

	if got.A.Sign() != 0 {
		t.Fatalf("A = %v, want 0", got.A)
	}
	if got.B.Cmp(big.NewInt(11)) != 0 {
		t.Fatalf("B = %v, want 11", got.B)
	}
	if got.C.Sign() != 0 {
		t.Fatalf("C = %v, want 0", got.C)
	}
	if got.D.Cmp(big.NewInt(13)) != 0 {
		t.Fatalf("D = %v, want 13", got.D)
	}
	if got.E.Sign() != 0 {
		t.Fatalf("E = %v, want 0", got.E)
	}
	if got.F.Cmp(big.NewInt(17)) != 0 {
		t.Fatalf("F = %v, want 17", got.F)
	}
	if got.G.Sign() != 0 {
		t.Fatalf("G = %v, want 0", got.G)
	}
	if got.H.Cmp(big.NewInt(19)) != 0 {
		t.Fatalf("H = %v, want 19", got.H)
	}
}

func TestWB_CollapseIndependentOfYToUnary_MapsBDFHIntoUnarySlots(t *testing.T) {
	state := blftState{
		A: big.NewInt(0),
		B: big.NewInt(11),
		C: big.NewInt(0),
		D: big.NewInt(13),
		E: big.NewInt(0),
		F: big.NewInt(17),
		G: big.NewInt(0),
		H: big.NewInt(19),
	}

	got := collapseIndependentOfYToUnary(state)

	if got.A.Sign() != 0 {
		t.Fatalf("A = %v, want 0", got.A)
	}
	if got.B.Cmp(big.NewInt(11)) != 0 {
		t.Fatalf("B = %v, want 11", got.B)
	}
	if got.C.Sign() != 0 {
		t.Fatalf("C = %v, want 0", got.C)
	}
	if got.D.Cmp(big.NewInt(13)) != 0 {
		t.Fatalf("D = %v, want 13", got.D)
	}
	if got.E.Sign() != 0 {
		t.Fatalf("E = %v, want 0", got.E)
	}
	if got.F.Cmp(big.NewInt(17)) != 0 {
		t.Fatalf("F = %v, want 17", got.F)
	}
	if got.G.Sign() != 0 {
		t.Fatalf("G = %v, want 0", got.G)
	}
	if got.H.Cmp(big.NewInt(19)) != 0 {
		t.Fatalf("H = %v, want 19", got.H)
	}
}

// core/gcf_helpers_wb_test.go v1
