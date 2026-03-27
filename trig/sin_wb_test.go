// trig/sin_wb_test.go v1
package trig

import (
	"math/big"
	"testing"
	"time"

	"github.com/egp/gosper-gcf/core"
)

func TestWB_SinFromTanHalf_ZeroIsExactlyZero(t *testing.T) {
	g := sinFromTanHalf(exactTerminalGCFForSin(
		[]int64{0},
		0, 1,
	))

	assertExactRCFSequenceSin(t, g, []int64{0})
}

func TestWB_SinFromTanHalf_OneHalfIsExactlyFourFifths(t *testing.T) {
	g := sinFromTanHalf(exactTerminalGCFForSin(
		[]int64{0, 2},
		1, 2,
	))

	assertExactRCFSequenceSin(t, g, []int64{0, 1, 4})
}

func TestWB_SinFromTanHalf_OneIsExactlyOne(t *testing.T) {
	g := sinFromTanHalf(exactTerminalGCFForSin(
		[]int64{1},
		1, 1,
	))

	assertExactRCFSequenceSin(t, g, []int64{1})
}

func exactTerminalGCFForSin(terms []int64, num, den int64) *core.GCF {
	rcfTerms := make([]core.RCFTerm, 0, len(terms))
	for _, term := range terms {
		rcfTerms = append(rcfTerms, core.NewRCFTerm(big.NewInt(term)))
	}

	value := core.NewRational(big.NewInt(num), big.NewInt(den))
	return core.NewExactTerminalGCF(
		rcfTerms,
		core.Range{
			Lo:     core.Endpoint{Value: value, Open: false},
			Hi:     core.Endpoint{Value: value, Open: false},
			Inside: true,
		},
	)
}

func assertExactRCFSequenceSin(t *testing.T, g *core.GCF, want []int64) {
	t.Helper()

	for i, w := range want {
		term, status := nextRCFWithTimeoutSin(t, g, time.Second)
		if status != core.StatusOK {
			t.Fatalf("term %d status = %v, want %v", i+1, status, core.StatusOK)
		}
		if term.A().Cmp(big.NewInt(w)) != 0 {
			t.Fatalf("term %d = %v, want %d", i+1, term.A(), w)
		}
	}

	_, eofStatus := nextRCFWithTimeoutSin(t, g, time.Second)
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}
}

func nextRCFWithTimeoutSin(t *testing.T, g *core.GCF, timeout time.Duration) (core.RCFTerm, core.Status) {
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

// trig/sin_wb_test.go v1
