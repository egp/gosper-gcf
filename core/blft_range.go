// core/blft_range.go v8
package core

import "math/big"

func (s blftState) CornerRange(xr, yr Range) Range {
	if special, ok := s.specialCaseRange(xr, yr); ok {
		return special
	}

	if !xr.Inside || !yr.Inside {
		panic("CornerRange currently supports only inside/inside ranges")
	}

	type cornerValue struct {
		value Rational
		open  bool
	}

	corners := []struct {
		x    Rational
		y    Rational
		open bool
	}{
		{xr.Lo.Value, yr.Lo.Value, xr.Lo.Open || yr.Lo.Open},
		{xr.Lo.Value, yr.Hi.Value, xr.Lo.Open || yr.Hi.Open},
		{xr.Hi.Value, yr.Lo.Value, xr.Hi.Open || yr.Lo.Open},
		{xr.Hi.Value, yr.Hi.Value, xr.Hi.Open || yr.Hi.Open},
	}

	values := make([]cornerValue, 0, 4)
	sawZeroDen := false
	sawPosDen := false
	sawNegDen := false

	for _, corner := range corners {
		num, den := evalBLFTNumDenAt(s, corner.x, corner.y)
		switch den.Sign() {
		case 0:
			sawZeroDen = true
			continue
		case 1:
			sawPosDen = true
		case -1:
			sawNegDen = true
		}

		values = append(values, cornerValue{
			value: NewRational(num, den),
			open:  corner.open,
		})
	}

	if sawZeroDen || (sawPosDen && sawNegDen) {
		raw := make([]Rational, 0, len(values))
		for _, v := range values {
			raw = append(raw, v.value)
		}
		return outsideRangeFromValues(raw)
	}

	if len(values) == 0 {
		return outsideRangeFromValues(nil)
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
		Lo: Endpoint{
			Value: lo,
			Open:  loOpen,
		},
		Hi: Endpoint{
			Value: hi,
			Open:  hiOpen,
		},
		Inside: true,
	}
}

func (s blftState) specialCaseRange(xr, yr Range) (Range, bool) {
	if r, ok := s.constantRange(); ok {
		return r, true
	}
	if r, ok := s.rangeWithExactX(xr, yr); ok {
		return r, true
	}
	if r, ok := s.rangeWithExactY(xr, yr); ok {
		return r, true
	}
	if r, ok := s.affineXRange(xr); ok {
		return r, true
	}
	if r, ok := s.affineYRange(yr); ok {
		return r, true
	}
	return Range{}, false
}

func (s blftState) constantRange() (Range, bool) {
	if !isZeroCoeff(s.A) || !isZeroCoeff(s.B) || !isZeroCoeff(s.C) ||
		!isZeroCoeff(s.E) || !isZeroCoeff(s.F) || !isZeroCoeff(s.G) {
		return Range{}, false
	}
	if isZeroCoeff(s.H) {
		return Range{}, false
	}

	value := NewRational(new(big.Int).Set(s.D), new(big.Int).Set(s.H))
	return exactRangeFromRational(value), true
}

func (s blftState) rangeWithExactX(xr, yr Range) (Range, bool) {
	if !isExactClosedRangeBLFT(xr) {
		return Range{}, false
	}

	reduced := s.reduceWithExactX(xr.Lo.Value)

	if r, ok := reduced.constantRange(); ok {
		return r, true
	}
	if r, ok := reduced.affineYRange(yr); ok {
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

	reduced := s.reduceWithExactY(yr.Lo.Value)

	if r, ok := reduced.constantRange(); ok {
		return r, true
	}
	if r, ok := reduced.affineXRange(xr); ok {
		return r, true
	}
	if r, ok := reduced.lftXRange(xr); ok {
		return r, true
	}

	return Range{}, false
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
	return r.Inside &&
		!r.Lo.Open &&
		!r.Hi.Open &&
		r.Lo.Value.Cmp(r.Hi.Value) == 0
}

func (s blftState) affineXRange(xr Range) (Range, bool) {
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
	return orderedRangeFromEndpoints(lo, hi, xr.Inside), true
}

func (s blftState) affineYRange(yr Range) (Range, bool) {
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
	return orderedRangeFromEndpoints(lo, hi, yr.Inside), true
}

func (s blftState) lftXRange(xr Range) (Range, bool) {
	if !isZeroCoeff(s.A) || !isZeroCoeff(s.C) || !isZeroCoeff(s.E) || !isZeroCoeff(s.G) {
		return Range{}, false
	}
	if isZeroCoeff(s.F) {
		return Range{}, false
	}
	if isZeroCoeff(s.F) && isZeroCoeff(s.H) {
		return Range{}, false
	}

	lo, loOK := applyLFTEndpointX(xr.Lo, s.B, s.D, s.F, s.H)
	hi, hiOK := applyLFTEndpointX(xr.Hi, s.B, s.D, s.F, s.H)

	pole := NewRational(new(big.Int).Neg(coeffOrZero(s.H)), coeffOrZero(s.F))
	asymptote := NewRational(coeffOrZero(s.B), coeffOrZero(s.F))

	if xr.Inside {
		if rangeIncludesRational(xr, pole) {
			return outsideRangeFromEndpointValues(lo, loOK, hi, hiOK), true
		}
		if loOK && hiOK {
			return orderedRangeFromEndpoints(lo, hi, true), true
		}
		return outsideRangeFromEndpointValues(lo, loOK, hi, hiOK), true
	}

	if rangeIncludesRational(xr, pole) {
		return outsideRangeFromValues([]Rational{asymptote}), true
	}

	return insideHullRangeFromValues(collectRationalsFromEndpointsAndValue(lo, loOK, hi, hiOK, asymptote)), true
}

func (s blftState) lftYRange(yr Range) (Range, bool) {
	if !isZeroCoeff(s.A) || !isZeroCoeff(s.B) || !isZeroCoeff(s.E) || !isZeroCoeff(s.F) {
		return Range{}, false
	}
	if isZeroCoeff(s.G) {
		return Range{}, false
	}
	if isZeroCoeff(s.G) && isZeroCoeff(s.H) {
		return Range{}, false
	}

	lo, loOK := applyLFTEndpointY(yr.Lo, s.C, s.D, s.G, s.H)
	hi, hiOK := applyLFTEndpointY(yr.Hi, s.C, s.D, s.G, s.H)

	pole := NewRational(new(big.Int).Neg(coeffOrZero(s.H)), coeffOrZero(s.G))
	asymptote := NewRational(coeffOrZero(s.C), coeffOrZero(s.G))

	if yr.Inside {
		if rangeIncludesRational(yr, pole) {
			return outsideRangeFromEndpointValues(lo, loOK, hi, hiOK), true
		}
		if loOK && hiOK {
			return orderedRangeFromEndpoints(lo, hi, true), true
		}
		return outsideRangeFromEndpointValues(lo, loOK, hi, hiOK), true
	}

	if rangeIncludesRational(yr, pole) {
		return outsideRangeFromValues([]Rational{asymptote}), true
	}

	return insideHullRangeFromValues(collectRationalsFromEndpointsAndValue(lo, loOK, hi, hiOK, asymptote)), true
}

func applyAffineEndpointX(ep Endpoint, b, d, h *big.Int) Endpoint {
	num := ep.Value.Num()
	den := ep.Value.Den()

	scaledNum := new(big.Int).Mul(new(big.Int).Set(b), num)
	constantNum := new(big.Int).Mul(new(big.Int).Set(d), den)
	outNum := new(big.Int).Add(scaledNum, constantNum)
	outDen := new(big.Int).Mul(new(big.Int).Set(h), den)

	return Endpoint{
		Value: NewRational(outNum, outDen),
		Open:  ep.Open,
	}
}

func applyAffineEndpointY(ep Endpoint, c, d, h *big.Int) Endpoint {
	num := ep.Value.Num()
	den := ep.Value.Den()

	scaledNum := new(big.Int).Mul(new(big.Int).Set(c), num)
	constantNum := new(big.Int).Mul(new(big.Int).Set(d), den)
	outNum := new(big.Int).Add(scaledNum, constantNum)
	outDen := new(big.Int).Mul(new(big.Int).Set(h), den)

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

func rangeIncludesRational(r Range, x Rational) bool {
	cmpLo := x.Cmp(r.Lo.Value)
	cmpHi := x.Cmp(r.Hi.Value)

	if r.Inside {
		if cmpLo < 0 || cmpHi > 0 {
			return false
		}
		if cmpLo == 0 && r.Lo.Open {
			return false
		}
		if cmpHi == 0 && r.Hi.Open {
			return false
		}
		return true
	}

	leftIncluded := cmpLo < 0 || (cmpLo == 0 && !r.Lo.Open)
	rightIncluded := cmpHi > 0 || (cmpHi == 0 && !r.Hi.Open)
	return leftIncluded || rightIncluded
}

func collectRationalsFromEndpointsAndValue(lo Endpoint, loOK bool, hi Endpoint, hiOK bool, extra Rational) []Rational {
	values := make([]Rational, 0, 3)
	if loOK {
		values = append(values, lo.Value)
	}
	if hiOK {
		values = append(values, hi.Value)
	}
	values = append(values, extra)
	return values
}

func outsideRangeFromEndpointValues(lo Endpoint, loOK bool, hi Endpoint, hiOK bool) Range {
	values := make([]Rational, 0, 2)
	if loOK {
		values = append(values, lo.Value)
	}
	if hiOK {
		values = append(values, hi.Value)
	}
	return outsideRangeFromValues(values)
}

func insideHullRangeFromValues(values []Rational) Range {
	if len(values) == 0 {
		zero := RationalFromInt64(0)
		return Range{
			Lo:     Endpoint{Value: zero, Open: false},
			Hi:     Endpoint{Value: zero, Open: false},
			Inside: true,
		}
	}

	lo := values[0]
	hi := values[0]
	for _, v := range values[1:] {
		if v.Cmp(lo) < 0 {
			lo = v
		}
		if v.Cmp(hi) > 0 {
			hi = v
		}
	}

	return Range{
		Lo:     Endpoint{Value: lo, Open: false},
		Hi:     Endpoint{Value: hi, Open: false},
		Inside: true,
	}
}

func orderedRangeFromEndpoints(lo, hi Endpoint, inside bool) Range {
	switch cmp := lo.Value.Cmp(hi.Value); {
	case cmp < 0:
		return Range{Lo: lo, Hi: hi, Inside: inside}
	case cmp > 0:
		return Range{Lo: hi, Hi: lo, Inside: inside}
	default:
		return Range{
			Lo: Endpoint{
				Value: lo.Value,
				Open:  lo.Open && hi.Open,
			},
			Hi: Endpoint{
				Value: hi.Value,
				Open:  lo.Open && hi.Open,
			},
			Inside: inside,
		}
	}
}

func isZeroCoeff(x *big.Int) bool {
	return x == nil || x.Sign() == 0
}

func coeffOrZero(x *big.Int) *big.Int {
	if x == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Set(x)
}

func preferXOnTie(xRange, yRange Range) bool {
	if !xRange.Inside || !yRange.Inside {
		panic("preferXOnTie currently supports only inside/inside ranges")
	}

	xWidth := rangeWidth(xRange)
	yWidth := rangeWidth(yRange)
	return xWidth.Cmp(yWidth) <= 0
}

func evalBLFTNumDenAt(s blftState, x, y Rational) (*big.Int, *big.Int) {
	xn := x.Num()
	xd := x.Den()
	yn := y.Num()
	yd := y.Den()

	commonDen := new(big.Int).Mul(xd, yd)

	num := big.NewInt(0)
	num.Add(num, scaledTermXY(s.A, xn, yn))
	num.Add(num, scaledTermX(s.B, xn, yd))
	num.Add(num, scaledTermY(s.C, yn, xd))
	num.Add(num, scaledConstant(s.D, commonDen))

	den := big.NewInt(0)
	den.Add(den, scaledTermXY(s.E, xn, yn))
	den.Add(den, scaledTermX(s.F, xn, yd))
	den.Add(den, scaledTermY(s.G, yn, xd))
	den.Add(den, scaledConstant(s.H, commonDen))

	return num, den
}

func outsideRangeFromValues(values []Rational) Range {
	if len(values) == 0 {
		zero := RationalFromInt64(0)
		return Range{
			Lo:     Endpoint{Value: zero, Open: false},
			Hi:     Endpoint{Value: zero, Open: false},
			Inside: false,
		}
	}

	lo := values[0]
	hi := values[0]
	for _, v := range values[1:] {
		if v.Cmp(lo) < 0 {
			lo = v
		}
		if v.Cmp(hi) > 0 {
			hi = v
		}
	}

	return Range{
		Lo:     Endpoint{Value: lo, Open: false},
		Hi:     Endpoint{Value: hi, Open: false},
		Inside: false,
	}
}

func scaledTermXY(coeff, xn, yn *big.Int) *big.Int {
	if coeff == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Mul(coeff, new(big.Int).Mul(xn, yn))
}

func scaledTermX(coeff, xn, yd *big.Int) *big.Int {
	if coeff == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Mul(coeff, new(big.Int).Mul(xn, yd))
}

func scaledTermY(coeff, yn, xd *big.Int) *big.Int {
	if coeff == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Mul(coeff, new(big.Int).Mul(yn, xd))
}

func scaledConstant(coeff, commonDen *big.Int) *big.Int {
	if coeff == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Mul(coeff, commonDen)
}

func rangeWidth(r Range) Rational {
	hiNum := r.Hi.Value.Num()
	hiDen := r.Hi.Value.Den()
	loNum := r.Lo.Value.Num()
	loDen := r.Lo.Value.Den()

	widthNumLeft := new(big.Int).Mul(hiNum, loDen)
	widthNumRight := new(big.Int).Mul(loNum, hiDen)
	widthNum := new(big.Int).Sub(widthNumLeft, widthNumRight)
	widthDen := new(big.Int).Mul(hiDen, loDen)

	return NewRational(widthNum, widthDen)
}

// core/blft_range.go v8
