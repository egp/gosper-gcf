// core/sqrt_cycle3_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

func TestWB_SqrtController_RejectsWhollyNegativeRadicandRange(t *testing.T) {
	x := &sqrtControllerCountingPQStream{
		rng: exactRangeFromRational(NewRational(big.NewInt(-4), big.NewInt(1))),
	}

	expectPanicWBSqrt(t, func() {
		_ = newSqrtController(x)
	})
}

func TestWB_SqrtController_FeedCertifiedTerm_AppendsIntoOurorobos(t *testing.T) {
	x := &sqrtControllerCountingPQStream{
		rng: exactRangeFromRational(RationalFromInt64(2)),
	}

	c := newSqrtController(x)

	c.feedCertifiedTerm(
		NewRCFTerm(big.NewInt(1)),
		exactRangeFromRational(RationalFromInt64(1)),
	)

	term, _, status := c.ourorobosApproximation().NextPQ()
	if status != StatusOK {
		t.Fatalf("status = %v, want %v", status, StatusOK)
	}
	if term.P.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("term.P = %v, want 1", term.P)
	}
	if term.Q.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("term.Q = %v, want 1", term.Q)
	}
}

func TestWB_Sqrt_FinitePositiveInput_IsNotConstructorTimeTerminalShortcut(t *testing.T) {
	stream, status := NewFinitePQStream([]FinitePQStep{
		{
			Term:  PQTerm{P: big.NewInt(4), Q: big.NewInt(1)},
			Range: exactRangeFromRational(RationalFromInt64(4)),
		},
	})
	if status != StatusOK {
		t.Fatalf("NewFinitePQStream status = %v, want %v", status, StatusOK)
	}

	g := Sqrt(stream)
	if g == nil {
		t.Fatalf("Sqrt returned nil")
	}
	if g.terminal != nil {
		t.Fatalf("Sqrt created terminal result at construction; want live sqrt path")
	}
}

func expectPanicWBSqrt(t *testing.T, fn func()) {
	t.Helper()

	defer func() {
		if recover() == nil {
			t.Fatalf("expected panic, got none")
		}
	}()

	fn()
}

// core/sqrt_cycle3_wb_test.go v1
