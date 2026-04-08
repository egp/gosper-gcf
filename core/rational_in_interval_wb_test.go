// core/rational_in_interval_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

// wbRat constructs an internal Rational for white-box tests.
func wbRat(num, den int64) Rational {
	return NewRational(big.NewInt(num), big.NewInt(den))
}

func wbCheckRat(t *testing.T, label string, got Rational, wantNum, wantDen int64) {
	t.Helper()
	if got.Num().Cmp(big.NewInt(wantNum)) != 0 || got.Den().Cmp(big.NewInt(wantDen)) != 0 {
		t.Fatalf("%s: got %v/%v, want %v/%v", label, got.Num(), got.Den(), wantNum, wantDen)
	}
}

// --- intCeilRational ---

func TestWB_intCeilRational_ClosedInteger(t *testing.T) {
	got := intCeilRational(wbRat(2, 1), false)
	wbCheckRat(t, "ceil(2,closed)", got, 2, 1)
}

func TestWB_intCeilRational_OpenInteger(t *testing.T) {
	got := intCeilRational(wbRat(2, 1), true)
	wbCheckRat(t, "ceil(2,open)", got, 3, 1)
}

func TestWB_intCeilRational_NonInteger(t *testing.T) {
	got := intCeilRational(wbRat(1, 3), false)
	wbCheckRat(t, "ceil(1/3,closed)", got, 1, 1)
}

func TestWB_intCeilRational_NegativeNonInteger(t *testing.T) {
	// ceil(-2/3) = 0 regardless of openness
	got := intCeilRational(wbRat(-2, 3), false)
	wbCheckRat(t, "ceil(-2/3,closed)", got, 0, 1)
}

// --- intFloorRational ---

func TestWB_intFloorRational_ClosedInteger(t *testing.T) {
	got := intFloorRational(wbRat(2, 1), false)
	wbCheckRat(t, "floor(2,closed)", got, 2, 1)
}

func TestWB_intFloorRational_OpenInteger(t *testing.T) {
	got := intFloorRational(wbRat(2, 1), true)
	wbCheckRat(t, "floor(2,open)", got, 1, 1)
}

func TestWB_intFloorRational_NonInteger(t *testing.T) {
	got := intFloorRational(wbRat(3, 2), false)
	wbCheckRat(t, "floor(3/2,closed)", got, 1, 1)
}

func TestWB_intFloorRational_NegativeNonInteger(t *testing.T) {
	// floor(-1/3) = -1
	got := intFloorRational(wbRat(-1, 3), false)
	wbCheckRat(t, "floor(-1/3,closed)", got, -1, 1)
}

// --- subtractRationals ---

func TestWB_subtractRationals(t *testing.T) {
	got := subtractRationals(wbRat(1, 2), wbRat(1, 3))
	wbCheckRat(t, "1/2 - 1/3", got, 1, 6)
}

func TestWB_subtractRationals_ResultZero(t *testing.T) {
	got := subtractRationals(wbRat(3, 4), wbRat(3, 4))
	wbCheckRat(t, "3/4 - 3/4", got, 0, 1)
}

// --- invertRational ---

func TestWB_invertRational_Positive(t *testing.T) {
	got, err := invertRational(wbRat(2, 3))
	if err != nil {
		t.Fatal(err)
	}
	wbCheckRat(t, "invert(2/3)", got, 3, 2)
}

func TestWB_invertRational_Negative(t *testing.T) {
	got, err := invertRational(wbRat(-3, 4))
	if err != nil {
		t.Fatal(err)
	}
	wbCheckRat(t, "invert(-3/4)", got, -4, 3)
}

func TestWB_invertRational_Zero(t *testing.T) {
	_, err := invertRational(wbRat(0, 1))
	if err == nil {
		t.Fatal("want error for invertRational(0), got nil")
	}
}

// --- smallestRationalInside: edge cases ---

// TestWB_SmallestRational_OpenLoInteger: (1, 3/2) → 4/3.
// Exercises the loShifted=0 branch: lo=1 is an integer and open,
// so lo-k=0 maps to +∞ in the reciprocal sub-problem.
func TestWB_SmallestRational_OpenLoInteger(t *testing.T) {
	got, err := smallestRationalInside(wbRat(1, 1), wbRat(3, 2), true, true)
	if err != nil {
		t.Fatal(err)
	}
	wbCheckRat(t, "(1,3/2)", got, 4, 3)
}

// TestWB_SmallestRational_FareyAdjacent: (2/5, 3/7) → 5/12.
// 2/5 and 3/7 are Farey-adjacent (2*7 - 3*5 = -1), so their mediant 5/12 is
// the simplest fraction strictly between them. Exercises deeper recursion.
func TestWB_SmallestRational_FareyAdjacent(t *testing.T) {
	got, err := smallestRationalInside(wbRat(2, 5), wbRat(3, 7), true, true)
	if err != nil {
		t.Fatal(err)
	}
	wbCheckRat(t, "(2/5,3/7)", got, 5, 12)
}

// core/rational_in_interval_wb_test.go v1
