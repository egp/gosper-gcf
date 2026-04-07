// core/sqrt_bb_test.go v4
package core_test

import (
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/egp/gosper-gcf/core"
)

const pendingTestSqrt = true

func TestBB_GCF_SqrtOfFourIsExactlyTwo(t *testing.T) {
	stream, status := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(4), Q: big.NewInt(1)},
			Range: sqrtExactRange(4, 1),
		},
	})
	if status != core.StatusOK {
		t.Fatalf("NewFinitePQStream status = %v, want %v", status, core.StatusOK)
	}

	g := core.Sqrt(stream)

	term, termStatus, err := nextRCFWithTimeoutSqrt(t, g, time.Second)
	if err != nil {
		t.Fatalf("first NextRCF error = %v", err)
	}
	if termStatus != core.StatusOK {
		t.Fatalf("first status = %v, want %v", termStatus, core.StatusOK)
	}
	if term.A().Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("first term = %v, want 2", term.A())
	}

	_, eofStatus, err := nextRCFWithTimeoutSqrt(t, g, time.Second)
	if err != nil {
		t.Fatalf("EOF NextRCF error = %v", err)
	}
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}
}

func TestBB_GCF_SqrtOfTwoMatchesKnownPrefix(t *testing.T) {
	if shouldSkipPendingSqrt() {
		t.Skip("pending public sqrt unary operation; set RUN_PENDING_TESTS=1 to run anyway")
	}

	stream, status := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(2), Q: big.NewInt(1)},
			Range: sqrtExactRange(2, 1),
		},
	})
	if status != core.StatusOK {
		t.Fatalf("NewFinitePQStream status = %v, want %v", status, core.StatusOK)
	}

	g := core.Sqrt(stream)
	assertRCFPrefixSqrt(t, g, []int64{1, 2, 2, 2, 2})
}

func TestBB_GCF_SqrtOfOneIsExactlyOne(t *testing.T) {
	stream, status := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(1), Q: big.NewInt(1)},
			Range: sqrtExactRange(1, 1),
		},
	})
	if status != core.StatusOK {
		t.Fatalf("NewFinitePQStream status = %v, want %v", status, core.StatusOK)
	}

	g := core.Sqrt(stream)

	term, termStatus, err := nextRCFWithTimeoutSqrt(t, g, time.Second)
	if err != nil {
		t.Fatalf("first NextRCF error = %v", err)
	}
	if termStatus != core.StatusOK {
		t.Fatalf("first status = %v, want %v", termStatus, core.StatusOK)
	}
	if term.A().Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("first term = %v, want 1", term.A())
	}

	_, eofStatus, err := nextRCFWithTimeoutSqrt(t, g, time.Second)
	if err != nil {
		t.Fatalf("EOF NextRCF error = %v", err)
	}
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}
}

func TestBB_GCF_SqrtOfOneQuarterIsOneHalf(t *testing.T) {
	stream, status := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(0), Q: big.NewInt(1)},
			Range: sqrtExactRange(1, 4),
		},
	})
	if status != core.StatusOK {
		t.Fatalf("NewFinitePQStream status = %v, want %v", status, core.StatusOK)
	}

	g := core.Sqrt(stream)

	term1, status1, err := nextRCFWithTimeoutSqrt(t, g, time.Second)
	if err != nil {
		t.Fatalf("first NextRCF error = %v", err)
	}
	if status1 != core.StatusOK {
		t.Fatalf("first status = %v, want %v", status1, core.StatusOK)
	}
	if term1.A().Cmp(big.NewInt(0)) != 0 {
		t.Fatalf("first term = %v, want 0", term1.A())
	}

	term2, status2, err := nextRCFWithTimeoutSqrt(t, g, time.Second)
	if err != nil {
		t.Fatalf("second NextRCF error = %v", err)
	}
	if status2 != core.StatusOK {
		t.Fatalf("second status = %v, want %v", status2, core.StatusOK)
	}
	if term2.A().Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("second term = %v, want 2", term2.A())
	}

	_, eofStatus, err := nextRCFWithTimeoutSqrt(t, g, time.Second)
	if err != nil {
		t.Fatalf("EOF NextRCF error = %v", err)
	}
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}
}

func TestBB_GCF_SqrtRejectsNegativeFiniteInput(t *testing.T) {
	stream, status := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(-4), Q: big.NewInt(1)},
			Range: sqrtExactRange(-4, 1),
		},
	})
	if status != core.StatusOK {
		t.Fatalf("NewFinitePQStream status = %v, want %v", status, core.StatusOK)
	}

	g := core.Sqrt(stream)

	_, _, err := nextRCFWithTimeoutSqrt(t, g, time.Second)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func shouldSkipPendingSqrt() bool {
	return pendingTestSqrt && os.Getenv("RUN_PENDING_TESTS") == ""
}

func assertRCFPrefixSqrt(t *testing.T, g *core.GCF, want []int64) {
	t.Helper()

	for i, w := range want {
		term, status, err := nextRCFWithTimeoutSqrt(t, g, time.Second)
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

func nextRCFWithTimeoutSqrt(t *testing.T, g *core.GCF, timeout time.Duration) (core.RCFTerm, core.Status, error) {
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
	case res := <-ch:
		return res.term, res.status, res.err
	case <-time.After(timeout):
		t.Fatalf("NextRCF() did not complete within %v", timeout)
		return core.NewRCFTerm(nil), core.StatusInvalidInput, nil
	}
}

func sqrtExactRange(num, den int64) core.Range {
	value := core.NewRational(big.NewInt(num), big.NewInt(den))
	return core.Range{
		Lo: core.Endpoint{
			Value: value,
			Open:  false,
		},
		Hi: core.Endpoint{
			Value: value,
			Open:  false,
		},
		Inside: true,
	}
}

// core/sqrt_bb_test.go v4
