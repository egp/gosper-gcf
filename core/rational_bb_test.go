package core_test

import (
	"math/big"
	"testing"

	"github.com/egp/gosper-gcf/core"
)

func TestBB_Rational_ExactIntegerConversion(t *testing.T) {
	r := core.FromInt64(7)

	if got := r.Num().Cmp(big.NewInt(7)); got != 0 {
		t.Fatalf("FromInt64(7).Num() = %v, want 7", r.Num())
	}

	if got := r.Den().Cmp(big.NewInt(1)); got != 0 {
		t.Fatalf("FromInt64(7).Den() = %v, want 1", r.Den())
	}
}
