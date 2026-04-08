// core/rational_in_interval.go v1
package core

import (
	"errors"
	"math/big"
)

var (
	ErrEmptyInterval                  = errors.New("SmallestRationalInInterval: empty interval")
	ErrSmallestRationalDepthExceeded  = errors.New("SmallestRationalInInterval: recursion depth exceeded")
	ErrOutsideIntervalNotYetSupported = errors.New("SmallestRationalInInterval: outside intervals not yet supported")
)

// SmallestRationalInInterval returns the rational with the smallest denominator
// (and, among those, smallest absolute numerator) that lies within r.
// Uses the HAKMEM 101C Gosper method: recurse on the same-floor sub-interval
// until an integer is the answer, then reconstruct via k + 1/result.
//
// Roles: Rational Collapse (bounds isolate an exact rational), canonical witness
// generation in tests, and interval-debugging support.
//
// Inside ranges are fully supported. Outside ranges return ErrOutsideIntervalNotYetSupported.
func SmallestRationalInInterval(r Range) (Rational, error) {
	if r.Kind() == RangeFull {
		return RationalFromInt64(0), nil
	}
	if !r.Inside {
		return Rational{}, ErrOutsideIntervalNotYetSupported
	}
	return smallestRationalInside(r.Lo.Value, r.Hi.Value, r.Lo.Open, r.Hi.Open)
}

func smallestRationalInside(lo, hi Rational, loOpen, hiOpen bool) (Rational, error) {
	cmp := lo.Cmp(hi)
	if cmp > 0 {
		return Rational{}, ErrEmptyInterval
	}
	if cmp == 0 {
		if loOpen || hiOpen {
			return Rational{}, ErrEmptyInterval
		}
		return lo, nil
	}
	return smallestRationalHelper(lo, hi, loOpen, hiOpen, 0)
}

// smallestRationalHelper is the recursive HAKMEM 101C kernel.
// Invariant: lo < hi on entry (strict).
func smallestRationalHelper(lo, hi Rational, loOpen, hiOpen bool, depth int) (Rational, error) {
	if depth > 300 {
		return Rational{}, ErrSmallestRationalDepthExceeded
	}

	// Step 1: find the smallest integer that lies within the interval (respecting openness).
	loCeil := intCeilRational(lo, loOpen)
	hiFloor := intFloorRational(hi, hiOpen)
	if loCeil.Cmp(hiFloor) <= 0 {
		return loCeil, nil
	}

	// Step 2: no integer in the interval; lo and hi share the same floor k.
	kBig, _, err := floorQuoRemChecked(lo.Num(), lo.Den())
	if err != nil {
		return Rational{}, err
	}
	kRat := rationalInt(kBig)

	loShifted := subtractRationals(lo, kRat)
	hiShifted := subtractRationals(hi, kRat)

	// Invert the sub-interval: [lo-k, hi-k] → [1/(hi-k), 1/(lo-k)].
	// Openness swaps because inversion reverses order.
	newLo, err := invertRational(hiShifted)
	if err != nil {
		return Rational{}, err
	}
	newLoOpen := hiOpen

	var inner Rational
	if loShifted.Num().Sign() == 0 {
		// lo equals k (an integer) and is open; 1/(lo-k) → +∞.
		// The simplest integer above newLo suffices.
		inner = intCeilRational(newLo, newLoOpen)
	} else {
		newHi, err := invertRational(loShifted)
		if err != nil {
			return Rational{}, err
		}
		inner, err = smallestRationalHelper(newLo, newHi, newLoOpen, loOpen, depth+1)
		if err != nil {
			return Rational{}, err
		}
	}

	// Reconstruct: kRat + 1/inner.
	innerInv, err := invertRational(inner)
	if err != nil {
		return Rational{}, err
	}
	// kRat is an integer so kRat.Den() == 1; simplify accordingly.
	num := new(big.Int).Add(
		new(big.Int).Mul(kRat.Num(), innerInv.Den()),
		innerInv.Num(),
	)
	return NewRationalChecked(num, innerInv.Den())
}

// intCeilRational returns the smallest integer ≥ r (closed) or > r (open).
func intCeilRational(r Rational, open bool) Rational {
	q, rem, err := floorQuoRemChecked(r.Num(), r.Den())
	if err != nil {
		return RationalFromInt64(0)
	}
	if rem.Sign() > 0 || open {
		q.Add(q, big.NewInt(1))
	}
	return rationalInt(q)
}

// intFloorRational returns the largest integer ≤ r (closed) or < r (open).
func intFloorRational(r Rational, open bool) Rational {
	q, rem, err := floorQuoRemChecked(r.Num(), r.Den())
	if err != nil {
		return RationalFromInt64(0)
	}
	if rem.Sign() == 0 && open {
		q.Sub(q, big.NewInt(1))
	}
	return rationalInt(q)
}

// subtractRationals returns a − b, reduced.
func subtractRationals(a, b Rational) Rational {
	num := new(big.Int).Sub(
		new(big.Int).Mul(a.Num(), b.Den()),
		new(big.Int).Mul(b.Num(), a.Den()),
	)
	den := new(big.Int).Mul(a.Den(), b.Den())
	return NewRational(num, den)
}

// invertRational returns 1/r; returns an error when r is zero.
func invertRational(r Rational) (Rational, error) {
	if r.Num().Sign() == 0 {
		return Rational{}, errors.New("invertRational: zero")
	}
	return NewRationalChecked(r.Den(), r.Num())
}

// rationalInt constructs the rational n/1 from a *big.Int.
func rationalInt(n *big.Int) Rational {
	return Rational{
		num: new(big.Int).Set(n),
		den: big.NewInt(1),
	}
}

// core/rational_in_interval.go v1
