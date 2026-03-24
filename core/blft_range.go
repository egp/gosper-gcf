// core/blft_range.go v3
package core

import "math/big"

func (s blftState) CornerRange(xr, yr Range) Range {
	if !xr.Inside || !yr.Inside {
		panic("CornerRange currently supports only inside/inside ranges")
	}

	corners := []Rational{
		evalBLFTAt(s, xr.Lo.Value, yr.Lo.Value),
		evalBLFTAt(s, xr.Lo.Value, yr.Hi.Value),
		evalBLFTAt(s, xr.Hi.Value, yr.Lo.Value),
		evalBLFTAt(s, xr.Hi.Value, yr.Hi.Value),
	}

	lo := corners[0]
	hi := corners[0]

	for _, c := range corners[1:] {
		if c.Cmp(lo) < 0 {
			lo = c
		}
		if c.Cmp(hi) > 0 {
			hi = c
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

func evalBLFTAt(s blftState, x, y Rational) Rational {
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

	return NewRational(num, den)
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

// core/blft_range.go v3
