// core/blft_range_exact.go v2
package core

import "math/big"

func (s blftState) constantRange() (Range, bool) {
	if !isZeroCoeff(s.A) || !isZeroCoeff(s.B) || !isZeroCoeff(s.C) || !isZeroCoeff(s.E) || !isZeroCoeff(s.F) || !isZeroCoeff(s.G) {
		return Range{}, false
	}
	if isZeroCoeff(s.H) {
		return Range{}, false
	}

	value := NewRational(coeffOrZero(s.D), coeffOrZero(s.H))
	return exactRangeFromRational(value), true
}

func (s blftState) rangeWithExactX(xr, yr Range) (Range, bool) {
	if !isExactClosedRangeBLFT(xr) {
		return Range{}, false
	}
	if !s.dependsOnX() {
		return Range{}, false
	}

	reduced := s.reduceWithExactX(xr.Lo.Value)
	if r, ok := reduced.constantRange(); ok {
		return r, true
	}
	if r, ok := reduced.affineYRangePreserveOutside(yr); ok {
		return r, true
	}
	if r, ok := reduced.lftYRange(yr); ok {
		return r, true
	}
	return Range{}, false
}

func (s blftState) rangeWithExactY(xr, yr Range) (Range, bool) {
	if !isExactClosedRangeBLFT(yr) {
		return Range{}, false
	}
	if !s.dependsOnY() {
		return Range{}, false
	}

	reduced := s.reduceWithExactY(yr.Lo.Value)
	if r, ok := reduced.constantRange(); ok {
		return r, true
	}
	if r, ok := reduced.affineXRangePreserveOutside(xr); ok {
		return r, true
	}
	if r, ok := reduced.lftXRange(xr); ok {
		return r, true
	}
	return Range{}, false
}

func (s blftState) dependsOnX() bool {
	return !isZeroCoeff(s.A) || !isZeroCoeff(s.B) || !isZeroCoeff(s.E) || !isZeroCoeff(s.F)
}

func (s blftState) dependsOnY() bool {
	return !isZeroCoeff(s.A) || !isZeroCoeff(s.C) || !isZeroCoeff(s.E) || !isZeroCoeff(s.G)
}

func (s blftState) reduceWithExactX(x Rational) blftState {
	xn := x.Num()
	xd := x.Den()

	return blftState{
		A: big.NewInt(0),
		B: big.NewInt(0),
		C: addScaledCoeffsBLFT(s.A, xn, s.C, xd),
		D: addScaledCoeffsBLFT(s.B, xn, s.D, xd),
		E: big.NewInt(0),
		F: big.NewInt(0),
		G: addScaledCoeffsBLFT(s.E, xn, s.G, xd),
		H: addScaledCoeffsBLFT(s.F, xn, s.H, xd),
	}
}

func (s blftState) reduceWithExactY(y Rational) blftState {
	yn := y.Num()
	yd := y.Den()

	return blftState{
		A: big.NewInt(0),
		B: addScaledCoeffsBLFT(s.A, yn, s.B, yd),
		C: big.NewInt(0),
		D: addScaledCoeffsBLFT(s.C, yn, s.D, yd),
		E: big.NewInt(0),
		F: addScaledCoeffsBLFT(s.E, yn, s.F, yd),
		G: big.NewInt(0),
		H: addScaledCoeffsBLFT(s.G, yn, s.H, yd),
	}
}

func addScaledCoeffsBLFT(leftCoeff, leftScale, rightCoeff, rightScale *big.Int) *big.Int {
	left := big.NewInt(0)
	right := big.NewInt(0)

	if leftCoeff != nil {
		left = new(big.Int).Mul(new(big.Int).Set(leftCoeff), new(big.Int).Set(leftScale))
	}
	if rightCoeff != nil {
		right = new(big.Int).Mul(new(big.Int).Set(rightCoeff), new(big.Int).Set(rightScale))
	}

	return new(big.Int).Add(left, right)
}

func isExactClosedRangeBLFT(r Range) bool {
	return r.Inside && !r.Lo.Open && !r.Hi.Open && r.Lo.Value.Cmp(r.Hi.Value) == 0
}

func isDegenerateClosedRangeBLFT(r Range) bool {
	return !r.Lo.Open && !r.Hi.Open && r.Lo.Value.Cmp(r.Hi.Value) == 0
}

func (s blftState) rangeWithDegenerateClosedX(xr, yr Range) (Range, bool) {
	if !isDegenerateClosedRangeBLFT(xr) || isExactClosedRangeBLFT(xr) {
		return Range{}, false
	}

	reduced := s.reduceWithExactX(xr.Lo.Value)
	if r, ok := reduced.constantRange(); ok {
		return r, true
	}
	if r, ok := reduced.affineYRangePreserveOutside(yr); ok {
		return r, true
	}
	if r, ok := reduced.lftYRange(yr); ok {
		return r, true
	}
	return Range{}, false
}

func (s blftState) rangeWithDegenerateClosedY(xr, yr Range) (Range, bool) {
	if !isDegenerateClosedRangeBLFT(yr) || isExactClosedRangeBLFT(yr) {
		return Range{}, false
	}

	reduced := s.reduceWithExactY(yr.Lo.Value)
	if r, ok := reduced.constantRange(); ok {
		return r, true
	}
	if r, ok := reduced.affineXRangePreserveOutside(xr); ok {
		return r, true
	}
	if r, ok := reduced.lftXRange(xr); ok {
		return r, true
	}
	return Range{}, false
}

func (s blftState) affineXRangePreserveOutside(xr Range) (Range, bool) {
	if !isZeroCoeff(s.A) || !isZeroCoeff(s.C) || !isZeroCoeff(s.E) || !isZeroCoeff(s.F) || !isZeroCoeff(s.G) {
		return Range{}, false
	}
	if isZeroCoeff(s.H) {
		return Range{}, false
	}
	if isZeroCoeff(s.B) && !isZeroCoeff(s.D) {
		return Range{}, false
	}

	lo := applyAffineEndpointX(xr.Lo, s.B, s.D, s.H)
	hi := applyAffineEndpointX(xr.Hi, s.B, s.D, s.H)
	return affineRangeFromEndpoints(lo, hi, xr.Inside, s.B, s.H), true
}

func (s blftState) affineYRangePreserveOutside(yr Range) (Range, bool) {
	if !isZeroCoeff(s.A) || !isZeroCoeff(s.B) || !isZeroCoeff(s.E) || !isZeroCoeff(s.F) || !isZeroCoeff(s.G) {
		return Range{}, false
	}
	if isZeroCoeff(s.H) {
		return Range{}, false
	}
	if isZeroCoeff(s.C) && !isZeroCoeff(s.D) {
		return Range{}, false
	}

	lo := applyAffineEndpointY(yr.Lo, s.C, s.D, s.H)
	hi := applyAffineEndpointY(yr.Hi, s.C, s.D, s.H)
	return affineRangeFromEndpoints(lo, hi, yr.Inside, s.C, s.H), true
}

// core/blft_range_exact.go v2
