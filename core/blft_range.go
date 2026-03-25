// core/blft_range.go v5
package core

import "math/big"

func (s blftState) CornerRange(xr, yr Range) Range {
	if !xr.Inside || !yr.Inside {
		panic("CornerRange currently supports only inside/inside ranges")
	}

	corners := [][2]Rational{
		{xr.Lo.Value, yr.Lo.Value},
		{xr.Lo.Value, yr.Hi.Value},
		{xr.Hi.Value, yr.Lo.Value},
		{xr.Hi.Value, yr.Hi.Value},
	}

	values := make([]Rational, 0, 4)
	sawZeroDen := false
	sawPosDen := false
	sawNegDen := false

	for _, corner := range corners {
		num, den := evalBLFTNumDenAt(s, corner[0], corner[1])

		switch den.Sign() {
		case 0:
			sawZeroDen = true
			continue
		case 1:
			sawPosDen = true
		case -1:
			sawNegDen = true
		}

		values = append(values, NewRational(num, den))
	}

	if sawZeroDen || (sawPosDen && sawNegDen) {
		return outsideRangeFromValues(values)
	}

	if len(values) == 0 {
		return outsideRangeFromValues(nil)
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
		Lo: Endpoint{
			Value: lo,
			Open:  false,
		},
		Hi: Endpoint{
			Value: hi,
			Open:  false,
		},
		Inside: true,
	}
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
			Lo: Endpoint{
				Value: zero,
				Open:  false,
			},
			Hi: Endpoint{
				Value: zero,
				Open:  false,
			},
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
		Lo: Endpoint{
			Value: lo,
			Open:  false,
		},
		Hi: Endpoint{
			Value: hi,
			Open:  false,
		},
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

// core/blft_range.go v5
