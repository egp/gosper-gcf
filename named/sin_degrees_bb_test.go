// named/sin_degrees_bb_test.go v1
package named_test

import (
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/egp/gosper-gcf/core"
	"github.com/egp/gosper-gcf/named"
)

const pendingTestNamedSinDegrees = true

func TestBB_Named_SinDegrees_ZeroIsExactlyZero(t *testing.T) {
	g := named.SinDegrees(core.PQStreamFromRational(core.RationalFromInt64(0)))
	assertExactRCFSequenceSinDegrees(t, g, []int64{0})
}

func TestBB_Named_SinDegrees_ThirtyIsExactlyOneHalf(t *testing.T) {
	if shouldSkipPendingNamedSinDegrees() {
		t.Skip("pending named.SinDegrees implementation; set RUN_PENDING_TESTS=1 to run anyway")
	}

	g := named.SinDegrees(core.PQStreamFromRational(core.RationalFromInt64(30)))
	assertExactRCFSequenceSinDegrees(t, g, []int64{0, 2})
}

func TestBB_Named_SinDegrees_NinetyIsExactlyOne(t *testing.T) {
	if shouldSkipPendingNamedSinDegrees() {
		t.Skip("pending named.SinDegrees implementation; set RUN_PENDING_TESTS=1 to run anyway")
	}

	g := named.SinDegrees(core.PQStreamFromRational(core.RationalFromInt64(90)))
	assertExactRCFSequenceSinDegrees(t, g, []int64{1})
}

func TestBB_Named_SinDegrees_MinusThirtyIsExactlyMinusOneHalf(t *testing.T) {
	if shouldSkipPendingNamedSinDegrees() {
		t.Skip("pending named.SinDegrees implementation; set RUN_PENDING_TESTS=1 to run anyway")
	}

	g := named.SinDegrees(core.PQStreamFromRational(core.RationalFromInt64(-30)))
	assertExactRCFSequenceSinDegrees(t, g, []int64{-1, 2})
}

func TestBB_Named_SinDegrees_SixtyNineMatchesKnownPrefix(t *testing.T) {
	if shouldSkipPendingNamedSinDegrees() {
		t.Skip("pending named.SinDegrees implementation; set RUN_PENDING_TESTS=1 to run anyway")
	}

	g := named.SinDegrees(core.PQStreamFromRational(core.RationalFromInt64(69)))
	assertRCFPrefixSinDegrees(t, g, []int64{0, 1, 14, 17, 1, 11, 1, 1, 5, 1})
}

func shouldSkipPendingNamedSinDegrees() bool {
	return pendingTestNamedSinDegrees && os.Getenv("RUN_PENDING_TESTS") == ""
}

func assertExactRCFSequenceSinDegrees(t *testing.T, g *core.GCF, want []int64) {
	t.Helper()

	for i, w := range want {
		term, status := nextRCFWithTimeoutSinDegrees(t, g, time.Second)
		if status != core.StatusOK {
			t.Fatalf("term %d status = %v, want %v", i+1, status, core.StatusOK)
		}
		if term.A().Cmp(big.NewInt(w)) != 0 {
			t.Fatalf("term %d = %v, want %d", i+1, term.A(), w)
		}
	}

	_, eofStatus := nextRCFWithTimeoutSinDegrees(t, g, time.Second)
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}
}

func assertRCFPrefixSinDegrees(t *testing.T, g *core.GCF, want []int64) {
	t.Helper()

	for i, w := range want {
		term, status := nextRCFWithTimeoutSinDegrees(t, g, time.Second)
		if status != core.StatusOK {
			t.Fatalf("term %d status = %v, want %v", i+1, status, core.StatusOK)
		}
		if term.A().Cmp(big.NewInt(w)) != 0 {
			t.Fatalf("term %d = %v, want %d", i+1, term.A(), w)
		}
	}
}

func nextRCFWithTimeoutSinDegrees(t *testing.T, g *core.GCF, timeout time.Duration) (core.RCFTerm, core.Status) {
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

// named/sin_degrees_bb_test.go v1
