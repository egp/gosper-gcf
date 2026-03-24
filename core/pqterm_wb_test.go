package core

import (
	"math/big"
	"testing"
)

func TestWB_PQTerm_AllowsNegativeFirstTermConstruction(t *testing.T) {
	term := PQTerm{
		P: big.NewInt(-3),
		Q: big.NewInt(5),
	}

	if term.P == nil || term.P.Cmp(big.NewInt(-3)) != 0 {
		t.Fatalf("P = %v, want -3", term.P)
	}
	if term.Q == nil || term.Q.Cmp(big.NewInt(5)) != 0 {
		t.Fatalf("Q = %v, want 5", term.Q)
	}
}
