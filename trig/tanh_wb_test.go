// trig/tanh_wb_test.go v4
package trig

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/egp/gosper-gcf/core"
)

func TestWB_HalfInputForTanh_NilPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("halfInputForTanh(nil) did not panic")
		}
	}()

	_ = halfInputForTanh(nil)
}

func TestWB_HalfInputForTanh_ZeroIsExactlyZero(t *testing.T) {
	g := observePQAsRCFForTanh(halfInputForTanh(
		core.PQStreamFromRational(core.RationalFromInt64(0)),
	))
	assertExactRCFSequenceTanhWB(t, g, []int64{0})
}

func TestWB_HalfInputForTanh_OneIsExactlyOneHalf(t *testing.T) {
	g := observePQAsRCFForTanh(halfInputForTanh(
		core.PQStreamFromRational(core.RationalFromInt64(1)),
	))
	assertExactRCFSequenceTanhWB(t, g, []int64{0, 2})
}

func TestWB_HalfInputForTanh_MinusOneIsExactlyMinusOneHalf(t *testing.T) {
	g := observePQAsRCFForTanh(halfInputForTanh(
		core.PQStreamFromRational(core.RationalFromInt64(-1)),
	))
	assertExactRCFSequenceTanhWB(t, g, []int64{-1, 2})
}

func TestWB_TanhFromTanhHalf_NilPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("tanhFromTanhHalf(nil) did not panic")
		}
	}()

	_ = tanhFromTanhHalf(nil)
}

func TestWB_TanhFromTanhHalf_ZeroIsExactlyZero(t *testing.T) {
	g := tanhFromTanhHalf(exactTerminalGCFForTanh(
		[]int64{0},
		0, 1,
	))
	assertExactRCFSequenceTanhWB(t, g, []int64{0})
}

func TestWB_TanhFromTanhHalf_OneHalfIsExactlyFourFifths(t *testing.T) {
	g := tanhFromTanhHalf(exactTerminalGCFForTanh(
		[]int64{0, 2},
		1, 2,
	))
	assertExactRCFSequenceTanhWB(t, g, []int64{0, 1, 4})
}

func TestWB_TanhFromTanhHalf_OneIsExactlyOne(t *testing.T) {
	g := tanhFromTanhHalf(exactTerminalGCFForTanh(
		[]int64{1},
		1, 1,
	))
	assertExactRCFSequenceTanhWB(t, g, []int64{1})
}

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

func observePQAsRCFForTanh(x core.PQStream) *core.GCF {
	if x == nil {
		panic("observePQAsRCFForTanh: nil input")
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

func exactTerminalGCFForTanh(terms []int64, num, den int64) *core.GCF {
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

func assertExactRCFSequenceTanhWB(t *testing.T, g *core.GCF, want []int64) {
	t.Helper()

	for i, w := range want {
		term, status, err := nextRCFWithTimeoutTanhWB(t, g, time.Second)
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

	_, eofStatus, err := nextRCFWithTimeoutTanhWB(t, g, time.Second)
	if err != nil {
		t.Fatalf("EOF NextRCF error = %v", err)
	}
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}
}

func assertRCFPrefixTanhWB(t *testing.T, g *core.GCF, want []int64) {
	t.Helper()

	for i, w := range want {
		term, status, err := nextRCFWithTimeoutTanhWB(t, g, time.Second)
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

func nextRCFWithTimeoutTanhWB(t *testing.T, g *core.GCF, timeout time.Duration) (core.RCFTerm, core.Status, error) {
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

func TestWB_DoubleAngleFromHalfQuotient_NilReturnsError(t *testing.T) {
	g := doubleAngleFromHalfQuotient(nil)
	if g == nil {
		t.Fatal("doubleAngleFromHalfQuotient(nil) returned nil, want non-nil error-producing GCF")
	}

	_, _, err := nextRCFWithTimeoutTanhWB(t, g, time.Second)
	if err == nil {
		t.Fatal("doubleAngleFromHalfQuotient(nil) first NextRCF error = nil, want non-nil error")
	}
}

func TestWB_DoubleAngleFromHalfQuotient_ZeroIsExactlyZero(t *testing.T) {
	g := doubleAngleFromHalfQuotient(exactTerminalGCFForTanh(
		[]int64{0},
		0, 1,
	))
	assertExactRCFSequenceTanhWB(t, g, []int64{0})
}

func TestWB_DoubleAngleFromHalfQuotient_OneHalfIsExactlyFourFifths(t *testing.T) {
	g := doubleAngleFromHalfQuotient(exactTerminalGCFForTanh(
		[]int64{0, 2},
		1, 2,
	))
	assertExactRCFSequenceTanhWB(t, g, []int64{0, 1, 4})
}

func TestWB_DoubleAngleFromHalfQuotient_OneIsExactlyOne(t *testing.T) {
	g := doubleAngleFromHalfQuotient(exactTerminalGCFForTanh(
		[]int64{1},
		1, 1,
	))
	assertExactRCFSequenceTanhWB(t, g, []int64{1})
}

// TestWB_Tanh_FirstTermSignForPositiveInputs verifies the property:
// tanh(x) ∈ (0,1) for x > 0  → first RCF term = 0
// tanh(-x) ∈ (-1,0) for x > 0 → first RCF term = -1
func TestWB_Tanh_FirstTermSignForPositiveInputs(t *testing.T) {
	inputs := []core.Rational{
		core.NewRational(big.NewInt(1), big.NewInt(4)),
		core.NewRational(big.NewInt(1), big.NewInt(2)),
		core.RationalFromInt64(1),
		core.RationalFromInt64(2),
		core.RationalFromInt64(3),
	}

	for _, x := range inputs {
		name := fmt.Sprintf("%v/%v", x.Num(), x.Den())
		t.Run(name, func(t *testing.T) {
			// positive: tanh(x) ∈ (0,1), floor = 0
			pos := Tanh(core.PQStreamFromRational(x))
			term, status, err := nextRCFWithTimeoutTanhWB(t, pos, time.Second)
			if err != nil {
				t.Fatalf("tanh(%v) NextRCF error = %v", x, err)
			}
			if status != core.StatusOK {
				t.Fatalf("tanh(%v) status = %v, want %v", x, status, core.StatusOK)
			}
			if term.A().Sign() != 0 {
				t.Fatalf("tanh(%v) first term = %v, want 0", x, term.A())
			}

			// negative: tanh(-x) ∈ (-1,0), floor = -1
			negX := core.NewRational(new(big.Int).Neg(x.Num()), x.Den())
			neg := Tanh(core.PQStreamFromRational(negX))
			term2, status2, err2 := nextRCFWithTimeoutTanhWB(t, neg, time.Second)
			if err2 != nil {
				t.Fatalf("tanh(%v) NextRCF error = %v", negX, err2)
			}
			if status2 != core.StatusOK {
				t.Fatalf("tanh(%v) status = %v, want %v", negX, status2, core.StatusOK)
			}
			if term2.A().Cmp(big.NewInt(-1)) != 0 {
				t.Fatalf("tanh(%v) first term = %v, want -1", negX, term2.A())
			}
		})
	}
}

// trig/tanh_wb_test.go v4
