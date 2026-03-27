// trig/lambert_wb_test.go v2
package trig

import (
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/egp/gosper-gcf/core"
)

const pendingTestLambertKernel = true

func TestWB_LambertMode_StageOddSequence(t *testing.T) {
	mode := lambertModeCircular
	want := []int64{1, 3, 5, 7, 9}

	for stage, w := range want {
		got := mode.stageOdd(stage)
		if got.Cmp(big.NewInt(w)) != 0 {
			t.Fatalf("stageOdd(%d) = %v, want %d", stage, got, w)
		}
	}
}

func TestWB_LambertMode_XYSignByMode(t *testing.T) {
	if got := lambertModeCircular.xySign(); got != -1 {
		t.Fatalf("circular xySign = %d, want -1", got)
	}
	if got := lambertModeHyperbolic.xySign(); got != 1 {
		t.Fatalf("hyperbolic xySign = %d, want 1", got)
	}
}

func TestWB_Lambert_CircularHalfAngleKernel_OneHalfMatchesKnownPrefix(t *testing.T) {
	if shouldSkipPendingLambertKernel() {
		t.Skip("pending Lambert kernel implementation; set RUN_PENDING_TESTS=1 to run anyway")
	}

	g := circularHalfAngleKernel(
		core.PQStreamFromRational(core.NewRational(big.NewInt(1), big.NewInt(2))),
	)

	assertRCFPrefixLambert(t, g, []int64{0, 1, 1, 4, 1, 8, 1, 12})
}

func TestWB_Lambert_HyperbolicHalfAngleKernel_OneHalfMatchesKnownPrefix(t *testing.T) {
	if shouldSkipPendingLambertKernel() {
		t.Skip("pending Lambert kernel implementation; set RUN_PENDING_TESTS=1 to run anyway")
	}

	g := hyperbolicHalfAngleKernel(
		core.PQStreamFromRational(core.NewRational(big.NewInt(1), big.NewInt(2))),
	)

	assertRCFPrefixLambert(t, g, []int64{0, 2, 6, 10, 14, 18, 22, 26})
}

func TestWB_Lambert_HyperbolicHalfAngleKernel_OneMatchesKnownPrefix(t *testing.T) {
	if shouldSkipPendingLambertKernel() {
		t.Skip("pending Lambert kernel implementation; set RUN_PENDING_TESTS=1 to run anyway")
	}

	g := hyperbolicHalfAngleKernel(
		core.PQStreamFromRational(core.RationalFromInt64(1)),
	)

	assertRCFPrefixLambert(t, g, []int64{0, 1, 3, 5, 7, 9, 11, 13})
}

func shouldSkipPendingLambertKernel() bool {
	return pendingTestLambertKernel && os.Getenv("RUN_PENDING_TESTS") == ""
}

func assertRCFPrefixLambert(t *testing.T, g *core.GCF, want []int64) {
	t.Helper()

	for i, w := range want {
		term, status := nextRCFWithTimeoutLambert(t, g, time.Second)
		if status != core.StatusOK {
			t.Fatalf("term %d status = %v, want %v", i+1, status, core.StatusOK)
		}
		if term.A().Cmp(big.NewInt(w)) != 0 {
			t.Fatalf("term %d = %v, want %d", i+1, term.A(), w)
		}
	}
}

func nextRCFWithTimeoutLambert(t *testing.T, g *core.GCF, timeout time.Duration) (core.RCFTerm, core.Status) {
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

// trig/lambert_wb_test.go v2
