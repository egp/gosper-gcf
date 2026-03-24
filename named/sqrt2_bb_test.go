package named_test

import (
	"math/big"
	"testing"

	"github.com/egp/gosper-gcf/core"
	"github.com/egp/gosper-gcf/named"
)

func TestBB_Sqrt2_FirstTermHasExpectedPrefix(t *testing.T) {
	s := named.Sqrt2()
	if s == nil {
		t.Fatal("Sqrt2() returned nil")
	}

	term, tail, status := s.NextPQ()

	if status != core.StatusOK {
		t.Fatalf("NextPQ() status = %v, want %v", status, core.StatusOK)
	}

	if tail == nil {
		t.Fatal("NextPQ() tail is nil, want non-nil tail stream")
	}

	if term.P == nil || term.Q == nil {
		t.Fatal("NextPQ() returned nil big.Int field")
	}

	if term.P.Cmp(big.NewInt(1)) != 0 || term.Q.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("first sqrt2 term = (%v,%v), want (1,1)", term.P, term.Q)
	}
}
