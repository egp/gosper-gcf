// trig/lambert_wb_test.go v3
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

func TestWB_LambertMode_CircularStageCoefficients(t *testing.T) {
	got := lambertModeCircular.stageCoefficients(big.NewInt(3))

	assertBLFTCoefficientsLambert(
		t,
		got,
		0, 1, 0, 0,
		-1, 0, 0, 3,
	)
}

func TestWB_LambertMode_HyperbolicStageCoefficients(t *testing.T) {
	got := lambertModeHyperbolic.stageCoefficients(big.NewInt(5))

	assertBLFTCoefficientsLambert(
		t,
		got,
		0, 1, 0, 0,
		1, 0, 0, 5,
	)
}

func TestWB_LambertStageGCF_CircularWithZeroTailIsXOverOdd(t *testing.T) {
	g := lambertStageGCF(
		lambertModeCircular,
		big.NewInt(3),
		core.PQStreamFromRational(core.NewRational(big.NewInt(1), big.NewInt(2))),
		core.PQStreamFromRational(core.RationalFromInt64(0)),
	)

	assertExactRCFSequenceLambert(t, g, []int64{0, 6})
}

func TestWB_LambertStageGCF_HyperbolicWithZeroTailIsXOverOdd(t *testing.T) {
	g := lambertStageGCF(
		lambertModeHyperbolic,
		big.NewInt(5),
		core.PQStreamFromRational(core.NewRational(big.NewInt(1), big.NewInt(2))),
		core.PQStreamFromRational(core.RationalFromInt64(0)),
	)

	assertExactRCFSequenceLambert(t, g, []int64{0, 10})
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

func assertBLFTCoefficientsLambert(
	t *testing.T,
	got core.BLFTCoefficients,
	a, b, c, d, e, f, g, h int64,
) {
	t.Helper()

	assertBigIntEqualLambert(t, "A", got.A, a)
	assertBigIntEqualLambert(t, "B", got.B, b)
	assertBigIntEqualLambert(t, "C", got.C, c)
	assertBigIntEqualLambert(t, "D", got.D, d)
	assertBigIntEqualLambert(t, "E", got.E, e)
	assertBigIntEqualLambert(t, "F", got.F, f)
	assertBigIntEqualLambert(t, "G", got.G, g)
	assertBigIntEqualLambert(t, "H", got.H, h)
}

func assertBigIntEqualLambert(t *testing.T, name string, got *big.Int, want int64) {
	t.Helper()

	if got == nil {
		t.Fatalf("%s = nil, want %d", name, want)
	}
	if got.Cmp(big.NewInt(want)) != 0 {
		t.Fatalf("%s = %v, want %d", name, got, want)
	}
}

func assertExactRCFSequenceLambert(t *testing.T, g *core.GCF, want []int64) {
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

	_, eofStatus := nextRCFWithTimeoutLambert(t, g, time.Second)
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}
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

// trig/lambert_wb_test.go v3
