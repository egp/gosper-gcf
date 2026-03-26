// named/constants_bb_test.go v4
package named_test

import (
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/egp/gosper-gcf/core"
	"github.com/egp/gosper-gcf/named"
)

const pendingTestNamedE = true
const pendingTestNamedPi = true

func TestBB_Named_E_MatchesKnownPrefix50(t *testing.T) {
	if shouldSkipPendingNamedE() {
		t.Skip("pending robust named E() source; set RUN_PENDING_TESTS=1 to run anyway")
	}

	g := core.NewGCF1(identityUnaryCoeffsNamedConstants(), named.E())
	assertRCFPrefixNamedConstants(t, g, eTermsNamedConstants(50))
}

func TestBB_Named_Pi_MatchesKnownPrefix(t *testing.T) {
	if shouldSkipPendingNamedPi() {
		t.Skip("pending named Pi() source; set RUN_PENDING_TESTS=1 to run anyway")
	}

	g := core.NewGCF1(identityUnaryCoeffsNamedConstants(), named.Pi())
	assertRCFPrefixNamedConstants(t, g, []int64{3, 7, 15, 1, 292, 1, 1, 1})
}

func shouldSkipPendingNamedE() bool {
	return pendingTestNamedE && os.Getenv("RUN_PENDING_TESTS") == ""
}

func shouldSkipPendingNamedPi() bool {
	return pendingTestNamedPi && os.Getenv("RUN_PENDING_TESTS") == ""
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

func eTermsNamedConstants(n int) []int64 {
	if n <= 0 {
		return nil
	}

	out := make([]int64, 0, n)
	out = append(out, 2)

	k := int64(1)
	for len(out) < n {
		out = append(out, 1)
		if len(out) >= n {
			break
		}
		out = append(out, 2*k)
		if len(out) >= n {
			break
		}
		out = append(out, 1)
		k++
	}

	return out
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
	case got := <-ch:
		return got.term, got.status
	case <-time.After(timeout):
		t.Fatalf("NextRCF timed out after %v", timeout)
		return core.NewRCFTerm(nil), core.StatusInvalidInput
	}
}

// named/constants_bb_test.go v4
