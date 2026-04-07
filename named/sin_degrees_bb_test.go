// named/sin_degrees_bb_test.go v2
package named_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/egp/gosper-gcf/core"
	"github.com/egp/gosper-gcf/named"
)

func TestBB_Named_SinDegrees_ZeroIsExactlyZero(t *testing.T) {
	g := named.SinDegrees(core.PQStreamFromRational(core.RationalFromInt64(0)))
	assertExactRCFSequenceSinDegrees(t, g, []int64{0})
}

func TestBB_Named_SinDegrees_ThirtyIsExactlyOneHalf(t *testing.T) {
	skipIfPending(t, "named.SinDegrees")

	g := named.SinDegrees(core.PQStreamFromRational(core.RationalFromInt64(30)))
	assertExactRCFSequenceSinDegrees(t, g, []int64{0, 2})
}

func TestBB_Named_SinDegrees_NinetyIsExactlyOne(t *testing.T) {
	skipIfPending(t, "named.SinDegrees")

	g := named.SinDegrees(core.PQStreamFromRational(core.RationalFromInt64(90)))
	assertExactRCFSequenceSinDegrees(t, g, []int64{1})
}

func TestBB_Named_SinDegrees_MinusThirtyIsExactlyMinusOneHalf(t *testing.T) {
	skipIfPending(t, "named.SinDegrees")

	g := named.SinDegrees(core.PQStreamFromRational(core.RationalFromInt64(-30)))
	assertExactRCFSequenceSinDegrees(t, g, []int64{-1, 2})
}

func TestBB_Named_SinDegrees_SixtyNineMatchesKnownPrefix(t *testing.T) {
	skipIfPending(t, "named.SinDegrees")

	g := named.SinDegrees(core.PQStreamFromRational(core.RationalFromInt64(69)))
	assertRCFPrefixSinDegrees(t, g, []int64{0, 1, 14, 17, 1, 11, 1, 1, 5, 1})
}

func assertExactRCFSequenceSinDegrees(t *testing.T, g *core.GCF, want []int64) {
	t.Helper()

	for i, w := range want {
		term, status, err := nextRCFWithTimeoutSinDegrees(t, g, time.Second)
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

	_, eofStatus, err := nextRCFWithTimeoutSinDegrees(t, g, time.Second)
	if err != nil {
		t.Fatalf("EOF NextRCF error = %v", err)
	}
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}
}

func assertRCFPrefixSinDegrees(t *testing.T, g *core.GCF, want []int64) {
	t.Helper()

	for i, w := range want {
		term, status, err := nextRCFWithTimeoutSinDegrees(t, g, time.Second)
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

func nextRCFWithTimeoutSinDegrees(t *testing.T, g *core.GCF, timeout time.Duration) (core.RCFTerm, core.Status, error) {
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

// named/sin_degrees_bb_test.go v2
