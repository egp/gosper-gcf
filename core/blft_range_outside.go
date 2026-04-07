// core/blft_range_outside.go v2
package core

import (
	"fmt"
	"math/big"
)

// cornerRangeOutsideXInsideY computes the image of the BLFT over an outside
// x-range and an inside y-range. It evaluates the function at the four finite
// corners plus the two x→±∞ limits (which reduce to (Ay+B)/(Ey+F)).
//
// The outside x-range has two components, (-∞, xLo] and [xHi, +∞). Any pole
// in x lies between xLo and xHi (the excluded interior), so denominator sign
// may differ between the two components — that is expected and correct.
// We only bail if a denominator is exactly zero at a corner or ∞-limit.
func (s blftState) cornerRangeOutsideXInsideY(xr, yr Range) (Range, error) {
	type cv struct {
		value Rational
		open  bool
	}

	values := make([]cv, 0, 6)
	sawZeroDen := false

	addCorner := func(num, den *big.Int, open bool) {
		if den.Sign() == 0 {
			sawZeroDen = true
			return
		}
		values = append(values, cv{value: NewRational(num, den), open: open})
	}

	// Four finite corners
	finiteCorners := [4]struct {
		x, y Rational
		open bool
	}{
		{xr.Lo.Value, yr.Lo.Value, xr.Lo.Open || yr.Lo.Open},
		{xr.Lo.Value, yr.Hi.Value, xr.Lo.Open || yr.Hi.Open},
		{xr.Hi.Value, yr.Lo.Value, xr.Hi.Open || yr.Lo.Open},
		{xr.Hi.Value, yr.Hi.Value, xr.Hi.Open || yr.Hi.Open},
	}
	for _, c := range finiteCorners {
		num, den := evalBLFTNumDenAt(s, c.x, c.y)
		addCorner(num, den, c.open)
	}

	// Two x→±∞ limits: both approach the same projective point ∞,
	// giving limit value (Ay+B)/(Ey+F) for each y endpoint.
	for _, yep := range [2]struct {
		y    Rational
		open bool
	}{{yr.Lo.Value, yr.Lo.Open}, {yr.Hi.Value, yr.Hi.Open}} {
		yn := yep.y.Num()
		yd := yep.y.Den()
		infNum := new(big.Int).Add(
			new(big.Int).Mul(coeffOrZero(s.A), yn),
			new(big.Int).Mul(coeffOrZero(s.B), yd),
		)
		infDen := new(big.Int).Add(
			new(big.Int).Mul(coeffOrZero(s.E), yn),
			new(big.Int).Mul(coeffOrZero(s.F), yd),
		)
		addCorner(infNum, infDen, true) // limit is always open
	}

	if sawZeroDen {
		return Range{}, fmt.Errorf(
			"CornerRange: %w\nBLFT=%s\nxRange=%s\nyRange=%s",
			ErrUnsupportedRangeCase,
			formatBLFTDebug(s),
			formatRangeDebug(xr),
			formatRangeDebug(yr),
		)
	}

	if len(values) == 0 {
		return outsideRangeFromValues(nil), nil
	}

	lo := values[0].value
	hi := values[0].value
	loOpen := values[0].open
	hiOpen := values[0].open

	for _, v := range values[1:] {
		switch cmp := v.value.Cmp(lo); {
		case cmp < 0:
			lo = v.value
			loOpen = v.open
		case cmp == 0:
			loOpen = loOpen && v.open
		}
		switch cmp := v.value.Cmp(hi); {
		case cmp > 0:
			hi = v.value
			hiOpen = v.open
		case cmp == 0:
			hiOpen = hiOpen && v.open
		}
	}

	return Range{
		Lo:     Endpoint{Value: lo, Open: loOpen},
		Hi:     Endpoint{Value: hi, Open: hiOpen},
		Inside: true,
	}, nil
}

// cornerRangeInsideXOutsideY is the symmetric counterpart: outside y-range,
// inside x-range. The x→±∞ limit becomes (Ax+C)/(Ex+G).
func (s blftState) cornerRangeInsideXOutsideY(xr, yr Range) (Range, error) {
	// Swap x↔y by transposing: A→A, B↔C, D→D, E→E, F↔G, H→H
	transposed := blftState{
		A: s.A, B: s.C, C: s.B, D: s.D,
		E: s.E, F: s.G, G: s.F, H: s.H,
	}
	return transposed.cornerRangeOutsideXInsideY(yr, xr)
}

// core/blft_range_outside.go v2
