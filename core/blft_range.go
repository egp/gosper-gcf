// core/blft_range.go v5
package core

import "math/big"

func (s blftState) CornerRange(xr, yr Range) Range {
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
