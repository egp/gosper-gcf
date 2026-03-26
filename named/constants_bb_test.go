// named/constants_bb_test.go v2
package named_test

import (
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/egp/gosper-gcf/core"
	"github.com/egp/gosper-gcf/named"
)

const pendingTestNamedConstants = true

func TestBB_Named_E_MatchesKnownPrefix(t *testing.T) {
	if shouldSkipPendingNamedConstants() {
		t.Skip("pending named E() source; set RUN_PENDING_TESTS=1 to run anyway")
	}

	g := core.NewGCF1(identityUnaryCoeffsNamedConstants(), named.E())
	assertRCFPrefixNamedConstants(t, g, []int64{2, 1, 2, 1, 1, 4, 1})
}

func TestBB_Named_Pi_MatchesKnownPrefix(t *testing.T) {
	if shouldSkipPendingNamedConstants() {
		t.Skip("pending named Pi() source; set RUN_PENDING_TESTS=1 to run anyway")
	}

	g := core.NewGCF1(identityUnaryCoeffsNamedConstants(), named.Pi())
	assertRCFPrefixNamedConstants(t, g, []int64{3, 7, 15, 1, 292, 1, 1, 1})
}

func shouldSkipPendingNamedConstants() bool {
	return pendingTestNamedConstants && os.Getenv("RUN_PENDING_TESTS") == ""
}

func identityUnaryCoeffsNamedConstants() core.BLFTCoefficients {
	return core.BLFTCoefficients{
		A: big.NewInt(0),
		B: big.NewInt(1),
		C: big.NewInt(0),
		D: big.NewInt(0),
		E: big.NewInt(0),
		F: big.NewInt(0),
		G: big.NewInt(0),
		H: big.NewInt(1),
	}
}

func assertRCFPrefixNamedConstants(t *testing.T, g *core.GCF, want []int64) {
	t.Helper()

	for i, w := range want {
		term, status := nextRCFWithTimeoutNamedConstants(t, g, time.Second)
		if status != core.StatusOK {
			t.Fatalf("term %d status = %v, want %v", i+1, status, core.StatusOK)
		}
		if term.A().Cmp(big.NewInt(w)) != 0 {
			t.Fatalf("term %d = %v, want %d", i+1, term.A(), w)
		}
	}
}

func nextRCFWithTimeoutNamedConstants(t *testing.T, g *core.GCF, timeout time.Duration) (core.RCFTerm, core.Status) {
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

// named/constants_bb_test.go v2
