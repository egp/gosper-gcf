// trig/tanh_wb_test.go v1
package trig

import (
	"math/big"
	"testing"
	"time"

	"github.com/egp/gosper-gcf/core"
)

func TestWB_Tanh_ZeroIsExactlyZero(t *testing.T) {
	g := Tanh(core.PQStreamFromRational(core.RationalFromInt64(0)))
	assertExactRCFSequenceTanhWB(t, g, []int64{0})
}

func TestWB_Tanh_OneHalfMatchesKnownPrefix(t *testing.T) {
	g := Tanh(core.PQStreamFromRational(
		core.NewRational(big.NewInt(1), big.NewInt(2)),
	))
	assertRCFPrefixTanhWB(t, g, []int64{0, 2, 6, 10, 14, 18, 22, 26})
}

func TestWB_Tanh_OneMatchesKnownPrefix(t *testing.T) {
	g := Tanh(core.PQStreamFromRational(core.RationalFromInt64(1)))
	assertRCFPrefixTanhWB(t, g, []int64{0, 1, 3, 5, 7, 9, 11, 13})
}

func assertExactRCFSequenceTanhWB(t *testing.T, g *core.GCF, want []int64) {
	t.Helper()

	for i, w := range want {
		term, status := nextRCFWithTimeoutTanhWB(t, g, time.Second)
		if status != core.StatusOK {
			t.Fatalf("term %d status = %v, want %v", i+1, status, core.StatusOK)
		}
		if term.A().Cmp(big.NewInt(w)) != 0 {
			t.Fatalf("term %d = %v, want %d", i+1, term.A(), w)
		}
	}

	_, eofStatus := nextRCFWithTimeoutTanhWB(t, g, time.Second)
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}
}

func assertRCFPrefixTanhWB(t *testing.T, g *core.GCF, want []int64) {
	t.Helper()

	for i, w := range want {
		term, status := nextRCFWithTimeoutTanhWB(t, g, time.Second)
		if status != core.StatusOK {
			t.Fatalf("term %d status = %v, want %v", i+1, status, core.StatusOK)
		}
		if term.A().Cmp(big.NewInt(w)) != 0 {
			t.Fatalf("term %d = %v, want %d", i+1, term.A(), w)
		}
	}
}

func nextRCFWithTimeoutTanhWB(t *testing.T, g *core.GCF, timeout time.Duration) (core.RCFTerm, core.Status) {
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
		t.Fatalf("NextRCF() did not complete within %v", timeout)
		return core.NewRCFTerm(nil), core.StatusInvalidInput
	}
}

// trig/tanh_wb_test.go v1
