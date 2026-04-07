// named/square_sqrt2_bb_test.go v4
package named_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/egp/gosper-gcf/core"
	"github.com/egp/gosper-gcf/named"
)

func TestBB_SquareOfSqrt2IsExactlyTwo(t *testing.T) {
	skipIfPending(t, "infinite algebraic square(sqrt2) certification gap")

	g := core.Square(named.Sqrt2())

	term, status, err := nextRCFWithTimeoutSquareSqrt2(t, g, time.Second)
	if err != nil {
		t.Fatalf("first NextRCF error = %v", err)
	}
	if status != core.StatusOK {
		t.Fatalf("first status = %v, want %v", status, core.StatusOK)
	}
	if term.A().Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("first term = %v, want 2", term.A())
	}

	_, eofStatus, err := nextRCFWithTimeoutSquareSqrt2(t, g, time.Second)
	if err != nil {
		t.Fatalf("EOF NextRCF error = %v", err)
	}
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}
}

func nextRCFWithTimeoutSquareSqrt2(t *testing.T, g *core.GCF, timeout time.Duration) (core.RCFTerm, core.Status, error) {
	t.Helper()

	type result struct {
		term   core.RCFTerm
		status core.Status
		err    error
	}

	ch := make(chan result, 1)

	go func() {
		term, status, err := g.NextRCF()
		ch <- result{term: term, status: status, err: err}
	}()

	select {
	case got := <-ch:
		return got.term, got.status, got.err
	case <-time.After(timeout):
		t.Fatalf("NextRCF() did not complete within %v", timeout)
		return core.NewRCFTerm(nil), core.StatusInvalidInput, nil
	}
}

// named/square_sqrt2_bb_test.go v4
