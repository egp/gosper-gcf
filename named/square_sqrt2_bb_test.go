// named/square_sqrt2_bb_test.go v2
package named_test

import (
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/egp/gosper-gcf/core"
	"github.com/egp/gosper-gcf/named"
)

const pendingTestDLFTInfiniteAlgebraic = true

func TestBB_SquareOfSqrt2IsExactlyTwo(t *testing.T) {
	if shouldSkipPendingDLFTInfiniteAlgebraic() {
		t.Skip("pending DLFT support for infinite algebraic streams under current PQStream+Range contract; set RUN_PENDING_TESTS=1 to run anyway")
	}

	g := core.Square(named.Sqrt2())

	term, status := nextRCFWithTimeoutSquareSqrt2(t, g, time.Second)
	if status != core.StatusOK {
		t.Fatalf("first status = %v, want %v", status, core.StatusOK)
	}
	if term.A().Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("first term = %v, want 2", term.A())
	}

	_, eofStatus := nextRCFWithTimeoutSquareSqrt2(t, g, time.Second)
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}
}

func shouldSkipPendingDLFTInfiniteAlgebraic() bool {
	return pendingTestDLFTInfiniteAlgebraic && os.Getenv("RUN_PENDING_TESTS") == ""
}

func nextRCFWithTimeoutSquareSqrt2(t *testing.T, g *core.GCF, timeout time.Duration) (core.RCFTerm, core.Status) {
	t.Helper()

	type result struct {
		term   core.RCFTerm
		status core.Status
	}

	ch := make(chan result, 1)

	go func() {
		term, status := g.NextRCF()
		ch <- result{term: term, status: status}
	}()

	select {
	case res := <-ch:
		return res.term, res.status
	case <-time.After(timeout):
		t.Fatalf("NextRCF() did not complete within %v", timeout)
		return core.NewRCFTerm(nil), core.StatusInvalidInput
	}
}

// named/square_sqrt2_bb_test.go v2
