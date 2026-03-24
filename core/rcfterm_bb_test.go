package core_test

import (
	"math/big"
	"testing"

	"github.com/egp/gosper-gcf/core"
)

func TestBB_RCFTerm_HoldsBigIntExactly(t *testing.T) {
	want := new(big.Int)
	want.SetString("123456789012345678901234567890", 10)

	term := core.NewRCFTerm(want)

	if got := term.A(); got.Cmp(want) != 0 {
		t.Fatalf("RCFTerm.A() = %v, want %v", got, want)
	}
}
