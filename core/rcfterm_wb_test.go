package core

import (
	"math/big"
	"testing"
)

func TestWB_RCFTerm_HoldsBigIntExactly(t *testing.T) {
	want := new(big.Int)
	want.SetString("123456789012345678901234567890", 10)

	term := NewRCFTerm(want)

	if term.a == nil {
		t.Fatal("internal RCFTerm storage is nil")
	}
	if term.a.Cmp(want) != 0 {
		t.Fatalf("internal RCFTerm storage = %v, want %v", term.a, want)
	}
}
