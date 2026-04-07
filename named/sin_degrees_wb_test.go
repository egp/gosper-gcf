// named/sin_degrees_wb_test.go v3
package named

import (
	"fmt"
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
		term, status, err := nextRCFWithTimeoutSinDegreesWB(t, g, time.Second)
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

	_, eofStatus, err := nextRCFWithTimeoutSinDegreesWB(t, g, time.Second)
	if err != nil {
		t.Fatalf("EOF NextRCF error = %v", err)
	}
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}
}

func assertRCFPrefixSinDegreesWB(t *testing.T, g *core.GCF, want []int64) {
	t.Helper()

	for i, w := range want {
		term, status, err := nextRCFWithTimeoutSinDegreesWB(t, g, time.Second)
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

func nextRCFWithTimeoutSinDegreesWB(t *testing.T, g *core.GCF, timeout time.Duration) (core.RCFTerm, core.Status, error) {
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

// TestWB_SinDegrees_RationalValuedAngles verifies the property:
// SinDegrees(d°) = exact rational RCF for all angles where sin is rational.
// Cases run individually so each can be activated as the implementation progresses.
func TestWB_SinDegrees_RationalValuedAngles(t *testing.T) {
	cases := []struct {
		degrees int64
		want    []int64 // exact RCF (finite)
	}{
		{0, []int64{0}},        // sin 0° = 0 — passes today
		{30, []int64{0, 2}},    // sin 30° = 1/2
		{90, []int64{1}},       // sin 90° = 1
		{150, []int64{0, 2}},   // sin 150° = 1/2
		{180, []int64{0}},      // sin 180° = 0
		{-30, []int64{-1, 2}},  // sin -30° = -1/2
		{-90, []int64{-1}},     // sin -90° = -1
		{-150, []int64{-1, 2}}, // sin -150° = -1/2
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%d°", tc.degrees), func(t *testing.T) {
			g := SinDegrees(core.PQStreamFromRational(core.RationalFromInt64(tc.degrees)))
			assertExactRCFSequenceSinDegreesWB(t, g, tc.want)
		})
	}
}

// named/sin_degrees_wb_test.go v3
