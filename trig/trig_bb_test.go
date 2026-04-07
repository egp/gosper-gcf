// trig/trig_bb_test.go v5
package trig_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/egp/gosper-gcf/core"
	"github.com/egp/gosper-gcf/trig"
)

func TestBB_Trig_Sin_ZeroIsExactlyZero(t *testing.T) {
	g := trig.Sin(core.PQStreamFromRational(core.RationalFromInt64(0)))
	assertExactRCFSequenceTrig(t, g, []int64{0})
}

func TestBB_Trig_Sin_OneHalfMatchesKnownPrefix(t *testing.T) {
	skipIfPending(t, "trig.Sin")

	g := trig.Sin(core.PQStreamFromRational(
		core.NewRational(big.NewInt(1), big.NewInt(2)),
	))
	assertRCFPrefixTrig(t, g, []int64{0, 2, 11, 1, 1, 1, 6, 2})
}

func TestBB_Trig_Tanh_ZeroIsExactlyZero(t *testing.T) {
	g := trig.Tanh(core.PQStreamFromRational(core.RationalFromInt64(0)))
	assertExactRCFSequenceTrig(t, g, []int64{0})
}

func TestBB_Trig_Tanh_OneHalfMatchesKnownPrefix(t *testing.T) {
	g := trig.Tanh(core.PQStreamFromRational(
		core.NewRational(big.NewInt(1), big.NewInt(2)),
	))
	assertRCFPrefixTrig(t, g, []int64{0, 2, 6, 10, 14, 18, 22, 26})
}

func TestBB_Trig_Tanh_OneMatchesKnownPrefix(t *testing.T) {
	g := trig.Tanh(core.PQStreamFromRational(core.RationalFromInt64(1)))
	assertRCFPrefixTrig(t, g, []int64{0, 1, 3, 5, 7, 9, 11, 13})
}

func TestBB_Trig_Tanh_MinusOneHalfMatchesKnownPrefix(t *testing.T) {
	g := trig.Tanh(core.PQStreamFromRational(
		core.NewRational(big.NewInt(-1), big.NewInt(2)),
	))
	assertRCFPrefixTrig(t, g, []int64{-1, 1, 1, 6, 10, 14, 18, 22})
}

func TestBB_Trig_Tanh_TwoMatchesKnownPrefix(t *testing.T) {
	g := trig.Tanh(core.PQStreamFromRational(core.RationalFromInt64(2)))
	assertRCFPrefixTrig(t, g, []int64{0, 1, 26, 1, 3, 1, 42, 2})
}

func assertExactRCFSequenceTrig(t *testing.T, g *core.GCF, want []int64) {
	t.Helper()

	for i, w := range want {
		term, status, err := nextRCFWithTimeoutTrig(t, g, time.Second)
		if err != nil {
			t.Fatalf("term %d NextRCF error = %v", i+1, err)
		}
		if status != core.StatusOK {
			t.Fatalf("term %d status = %v, want %v", i+1, status, core.StatusOK)
		}
		if term.A().Cmp(big.NewInt(w)) != 0 {
			t.Fatalf("term %d = %v, want %d", i+1, term.A(), w)
		}
	}

	_, eofStatus, err := nextRCFWithTimeoutTrig(t, g, time.Second)
	if err != nil {
		t.Fatalf("EOF NextRCF error = %v", err)
	}
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}
}

func assertRCFPrefixTrig(t *testing.T, g *core.GCF, want []int64) {
	t.Helper()

	for i, w := range want {
		term, status, err := nextRCFWithTimeoutTrig(t, g, time.Second)
		if err != nil {
			t.Fatalf("term %d NextRCF error = %v", i+1, err)
		}
		if status != core.StatusOK {
			t.Fatalf("term %d status = %v, want %v", i+1, status, core.StatusOK)
		}
		if term.A().Cmp(big.NewInt(w)) != 0 {
			t.Fatalf("term %d = %v, want %d", i+1, term.A(), w)
		}
	}
}

func nextRCFWithTimeoutTrig(t *testing.T, g *core.GCF, timeout time.Duration) (core.RCFTerm, core.Status, error) {
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

// trig/trig_bb_test.go v5
