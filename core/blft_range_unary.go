// core/blft_range_unary.go v3
package core

import "math/big"

func (s blftState) affineXRange(xr Range) (Range, bool) {
	if !xr.Inside {
		return Range{}, false
	}
	if !isZeroCoeff(s.A) || !isZeroCoeff(s.C) || !isZeroCoeff(s.E) ||
		!isZeroCoeff(s.F) || !isZeroCoeff(s.G) {
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
	return orderedRangeFromEndpoints(lo, hi, true), true
}

func (s blftState) affineYRange(yr Range) (Range, bool) {
	if !yr.Inside {
		return Range{}, false
	}
	if !isZeroCoeff(s.A) || !isZeroCoeff(s.B) || !isZeroCoeff(s.E) ||
		!isZeroCoeff(s.F) || !isZeroCoeff(s.G) {
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
	return orderedRangeFromEndpoints(lo, hi, true), true
}

func (s blftState) lftXRange(xr Range) (Range, bool) {
	if !xr.Inside {
		return Range{}, false
	}
	if !isZeroCoeff(s.A) || !isZeroCoeff(s.C) || !isZeroCoeff(s.E) || !isZeroCoeff(s.G) {
		return Range{}, false
	}
	if isZeroCoeff(s.F) {
		return Range{}, false
	}

	lo, loOK := applyLFTEndpointX(xr.Lo, s.B, s.D, s.F, s.H)
	hi, hiOK := applyLFTEndpointX(xr.Hi, s.B, s.D, s.F, s.H)
	pole := NewRational(new(big.Int).Neg(coeffOrZero(s.H)), coeffOrZero(s.F))

	if rangeIncludesRational(xr, pole) {
		return outsideRangeFromEndpointValues(lo, loOK, hi, hiOK), true
	}
	if loOK && hiOK {
		return orderedRangeFromEndpoints(lo, hi, true), true
	}
	return outsideRangeFromEndpointValues(lo, loOK, hi, hiOK), true
}

func (s blftState) lftYRange(yr Range) (Range, bool) {
	if !yr.Inside {
		return Range{}, false
	}
	if !isZeroCoeff(s.A) || !isZeroCoeff(s.B) || !isZeroCoeff(s.E) || !isZeroCoeff(s.F) {
		return Range{}, false
	}
	if isZeroCoeff(s.G) {
		return Range{}, false
	}

	lo, loOK := applyLFTEndpointY(yr.Lo, s.C, s.D, s.G, s.H)
	hi, hiOK := applyLFTEndpointY(yr.Hi, s.C, s.D, s.G, s.H)
	pole := NewRational(new(big.Int).Neg(coeffOrZero(s.H)), coeffOrZero(s.G))

	if rangeIncludesRational(yr, pole) {
		return outsideRangeFromEndpointValues(lo, loOK, hi, hiOK), true
	}
	if loOK && hiOK {
		return orderedRangeFromEndpoints(lo, hi, true), true
	}
	return outsideRangeFromEndpointValues(lo, loOK, hi, hiOK), true
}

func applyAffineEndpointX(ep Endpoint, b, d, h *big.Int) Endpoint {
	num := ep.Value.Num()
	den := ep.Value.Den()

	scaledNum := new(big.Int).Mul(coeffOrZero(b), num)
	constantNum := new(big.Int).Mul(coeffOrZero(d), den)
	outNum := new(big.Int).Add(scaledNum, constantNum)
	outDen := new(big.Int).Mul(coeffOrZero(h), den)

	return Endpoint{
		Value: NewRational(outNum, outDen),
		Open:  ep.Open,
	}
}

func applyAffineEndpointY(ep Endpoint, c, d, h *big.Int) Endpoint {
	num := ep.Value.Num()
	den := ep.Value.Den()

	scaledNum := new(big.Int).Mul(coeffOrZero(c), num)
	constantNum := new(big.Int).Mul(coeffOrZero(d), den)
	outNum := new(big.Int).Add(scaledNum, constantNum)
	outDen := new(big.Int).Mul(coeffOrZero(h), den)

	return Endpoint{
		Value: NewRational(outNum, outDen),
		Open:  ep.Open,
	}
}

func applyLFTEndpointX(ep Endpoint, b, d, f, h *big.Int) (Endpoint, bool) {
	num := ep.Value.Num()
	den := ep.Value.Den()

	outNum := new(big.Int).Add(
		new(big.Int).Mul(coeffOrZero(b), num),
		new(big.Int).Mul(coeffOrZero(d), den),
	)
	outDen := new(big.Int).Add(
		new(big.Int).Mul(coeffOrZero(f), num),
		new(big.Int).Mul(coeffOrZero(h), den),
	)

	if outDen.Sign() == 0 {
		return Endpoint{}, false
	}

	return Endpoint{
		Value: NewRational(outNum, outDen),
		Open:  ep.Open,
	}, true
}

func applyLFTEndpointY(ep Endpoint, c, d, g, h *big.Int) (Endpoint, bool) {
	num := ep.Value.Num()
	den := ep.Value.Den()

	outNum := new(big.Int).Add(
		new(big.Int).Mul(coeffOrZero(c), num),
		new(big.Int).Mul(coeffOrZero(d), den),
	)
	outDen := new(big.Int).Add(
		new(big.Int).Mul(coeffOrZero(g), num),
		new(big.Int).Mul(coeffOrZero(h), den),
	)

	if outDen.Sign() == 0 {
		return Endpoint{}, false
	}

	return Endpoint{
		Value: NewRational(outNum, outDen),
		Open:  ep.Open,
	}, true
}

// core/blft_range_unary.go v3
