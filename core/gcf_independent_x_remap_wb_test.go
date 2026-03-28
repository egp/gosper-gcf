// core/gcf_independent_x_remap_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

func TestWB_CollapseIndependentOfXToUnary_MapsCDGHIntoUnarySlots(t *testing.T) {
	got := collapseIndependentOfXToUnary(blftState{
		A: big.NewInt(0),
		B: big.NewInt(0),
		C: big.NewInt(11),
		D: big.NewInt(13),
		E: big.NewInt(0),
		F: big.NewInt(0),
		G: big.NewInt(17),
		H: big.NewInt(19),
	})

	if got.B.Cmp(big.NewInt(11)) != 0 {
		t.Fatalf("B = %v, want 11", got.B)
	}
	if got.D.Cmp(big.NewInt(13)) != 0 {
		t.Fatalf("D = %v, want 13", got.D)
	}
	if got.F.Cmp(big.NewInt(17)) != 0 {
		t.Fatalf("F = %v, want 17", got.F)
	}
	if got.H.Cmp(big.NewInt(19)) != 0 {
		t.Fatalf("H = %v, want 19", got.H)
	}
}

// core/gcf_independent_x_remap_wb_test.go v1
