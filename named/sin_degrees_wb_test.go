// named/sin_degrees_wb_test.go v1
package named

import (
	"math/big"
	"testing"
	"time"

	"github.com/egp/gosper-gcf/core"
)

func TestWB_DegreesToRadiansForSin_ZeroIsExactlyZero(t *testing.T) {
	g := observePQAsRCFSinDegrees(
		degreesToRadiansForSin(core.PQStreamFromRational(core.RationalFromInt64(0))),
	)
	assertExactRCFSequenceSinDegreesWB(t, g, []int64{0})
}

func TestWB_DegreesToRadiansForSin_OneEightyMatchesPiPrefix(t *testing.T) {
	g := observePQAsRCFSinDegrees(
		degreesToRadiansForSin(core.PQStreamFromRational(core.RationalFromInt64(180))),
	)
	assertRCFPrefixSinDegreesWB(t, g, []int64{3, 7, 15, 1, 292, 1, 1, 1})
}

func observePQAsRCFSinDegrees(x core.PQStream) *core.GCF {
	if x == nil {
		panic("observePQAsRCFSinDegrees: nil input")
	}

	return core.NewGCF1(
		core.BLFTCoefficients{
			A: big.NewInt(0),
			B: big.NewInt(1),
			C: big.NewInt(0),
			D: big.NewInt(0),
			E: big.NewInt(0),
			F: big.NewInt(0),
			G: big.NewInt(0),
			H: big.NewInt(1),
		},
		x,
	)
}

func assertExactRCFSequenceSinDegreesWB(t *testing.T, g *core.GCF, want []int64) {
	t.Helper()

	for i, w := range want {
		term, status := nextRCFWithTimeoutSinDegreesWB(t, g, time.Second)
		if status != core.StatusOK {
			t.Fatalf("term %d status = %v, want %v", i+1, status, core.StatusOK)
		}
		if term.A().Cmp(big.NewInt(w)) != 0 {
			t.Fatalf("term %d = %v, want %d", i+1, term.A(), w)
		}
	}

	_, eofStatus := nextRCFWithTimeoutSinDegreesWB(t, g, time.Second)
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}
}

func assertRCFPrefixSinDegreesWB(t *testing.T, g *core.GCF, want []int64) {
	t.Helper()

	for i, w := range want {
		term, status := nextRCFWithTimeoutSinDegreesWB(t, g, time.Second)
		if status != core.StatusOK {
			t.Fatalf("term %d status = %v, want %v", i+1, status, core.StatusOK)
		}
		if term.A().Cmp(big.NewInt(w)) != 0 {
			t.Fatalf("term %d = %v, want %d", i+1, term.A(), w)
		}
	}
}

func nextRCFWithTimeoutSinDegreesWB(t *testing.T, g *core.GCF, timeout time.Duration) (core.RCFTerm, core.Status) {
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

// named/sin_degrees_wb_test.go v1
