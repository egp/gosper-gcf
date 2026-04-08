// core/rational_in_interval_bb_test.go v1
package core_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/egp/gosper-gcf/core"
)

// sriRat constructs a Rational from two int64 values for SmallestRationalInInterval tests.
func sriRat(num, den int64) core.Rational {
	return core.NewRational(big.NewInt(num), big.NewInt(den))
}

// sriRange builds an inside Range with explicit openness flags.
func sriRange(lo, hi core.Rational, loOpen, hiOpen bool) core.Range {
	return core.Range{
		Lo:     core.Endpoint{Value: lo, Open: loOpen},
		Hi:     core.Endpoint{Value: hi, Open: hiOpen},
		Inside: true,
	}
}

// checkSRI asserts that got equals the rational wantNum/wantDen.
func checkSRI(t *testing.T, label string, got core.Rational, wantNum, wantDen int64) {
	t.Helper()
	if got.Num().Cmp(big.NewInt(wantNum)) != 0 || got.Den().Cmp(big.NewInt(wantDen)) != 0 {
		t.Fatalf("%s: got %v/%v, want %v/%v", label, got.Num(), got.Den(), wantNum, wantDen)
	}
}

// TestBB_SmallestRational_ClosedSymmetric: [1/4, 3/4] → 1/2.
func TestBB_SmallestRational_ClosedSymmetric(t *testing.T) {
	r := sriRange(sriRat(1, 4), sriRat(3, 4), false, false)
	got, err := core.SmallestRationalInInterval(r)
	if err != nil {
		t.Fatal(err)
	}
	checkSRI(t, "[1/4,3/4]", got, 1, 2)
}

// TestBB_SmallestRational_OpenFractions: (1/3, 1/2) → 2/5.
func TestBB_SmallestRational_OpenFractions(t *testing.T) {
	r := sriRange(sriRat(1, 3), sriRat(1, 2), true, true)
	got, err := core.SmallestRationalInInterval(r)
	if err != nil {
		t.Fatal(err)
	}
	checkSRI(t, "(1/3,1/2)", got, 2, 5)
}

// TestBB_SmallestRational_ClosedLoBoundary: [1/3, 5/12] → 1/3.
// 1/2 is not in the interval (0.5 > 5/12 ≈ 0.417), so denom-3 fraction 1/3 wins.
func TestBB_SmallestRational_ClosedLoBoundary(t *testing.T) {
	r := sriRange(sriRat(1, 3), sriRat(5, 12), false, false)
	got, err := core.SmallestRationalInInterval(r)
	if err != nil {
		t.Fatal(err)
	}
	checkSRI(t, "[1/3,5/12]", got, 1, 3)
}

// TestBB_SmallestRational_ClosedIntegerRange: [2, 5] → 2.
func TestBB_SmallestRational_ClosedIntegerRange(t *testing.T) {
	r := sriRange(sriRat(2, 1), sriRat(5, 1), false, false)
	got, err := core.SmallestRationalInInterval(r)
	if err != nil {
		t.Fatal(err)
	}
	checkSRI(t, "[2,5]", got, 2, 1)
}

// TestBB_SmallestRational_OpenIntegerRange: (2, 5) → 3.
func TestBB_SmallestRational_OpenIntegerRange(t *testing.T) {
	r := sriRange(sriRat(2, 1), sriRat(5, 1), true, true)
	got, err := core.SmallestRationalInInterval(r)
	if err != nil {
		t.Fatal(err)
	}
	checkSRI(t, "(2,5)", got, 3, 1)
}

// TestBB_SmallestRational_SpanningZero: [-1/2, 1/2] → 0.
func TestBB_SmallestRational_SpanningZero(t *testing.T) {
	r := sriRange(sriRat(-1, 2), sriRat(1, 2), false, false)
	got, err := core.SmallestRationalInInterval(r)
	if err != nil {
		t.Fatal(err)
	}
	checkSRI(t, "[-1/2,1/2]", got, 0, 1)
}

// TestBB_SmallestRational_SinglePoint: [3/7, 3/7] → 3/7.
func TestBB_SmallestRational_SinglePoint(t *testing.T) {
	r := sriRange(sriRat(3, 7), sriRat(3, 7), false, false)
	got, err := core.SmallestRationalInInterval(r)
	if err != nil {
		t.Fatal(err)
	}
	checkSRI(t, "[3/7,3/7]", got, 3, 7)
}

// TestBB_SmallestRational_NegativeInterval: [-2/3, -1/3] → -1/2.
func TestBB_SmallestRational_NegativeInterval(t *testing.T) {
	r := sriRange(sriRat(-2, 3), sriRat(-1, 3), false, false)
	got, err := core.SmallestRationalInInterval(r)
	if err != nil {
		t.Fatal(err)
	}
	checkSRI(t, "[-2/3,-1/3]", got, -1, 2)
}

// TestBB_SmallestRational_EmptyOpenSinglePoint: (1/2, 1/2) → ErrEmptyInterval.
func TestBB_SmallestRational_EmptyOpenSinglePoint(t *testing.T) {
	r := sriRange(sriRat(1, 2), sriRat(1, 2), true, false)
	_, err := core.SmallestRationalInInterval(r)
	if !errors.Is(err, core.ErrEmptyInterval) {
		t.Fatalf("want ErrEmptyInterval, got %v", err)
	}
}

// TestBB_SmallestRational_OutsideReturnsError: outside [1,2] → ErrOutsideIntervalNotYetSupported.
func TestBB_SmallestRational_OutsideReturnsError(t *testing.T) {
	r := core.Range{
		Lo:     core.Endpoint{Value: sriRat(1, 1), Open: false},
		Hi:     core.Endpoint{Value: sriRat(2, 1), Open: false},
		Inside: false,
	}
	_, err := core.SmallestRationalInInterval(r)
	if !errors.Is(err, core.ErrOutsideIntervalNotYetSupported) {
		t.Fatalf("want ErrOutsideIntervalNotYetSupported, got %v", err)
	}
}

// core/rational_in_interval_bb_test.go v1
