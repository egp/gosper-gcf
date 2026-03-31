// core/blft_range_helpers.go v2
package core

import (
	"fmt"
	"math/big"
)

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

func rangeIncludesRational(r Range, q Rational) bool {
	cmpLo := q.Cmp(r.Lo.Value)
	cmpHi := q.Cmp(r.Hi.Value)

	if r.Inside {
		leftIncluded := cmpLo > 0 || (cmpLo == 0 && !r.Lo.Open)
		rightIncluded := cmpHi < 0 || (cmpHi == 0 && !r.Hi.Open)
		return leftIncluded && rightIncluded
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

func formatBLFTDebug(s blftState) string {
	return fmt.Sprintf(
		"(A=%s B=%s C=%s D=%s E=%s F=%s G=%s H=%s)",
		formatBigIntDebug(s.A),
		formatBigIntDebug(s.B),
		formatBigIntDebug(s.C),
		formatBigIntDebug(s.D),
		formatBigIntDebug(s.E),
		formatBigIntDebug(s.F),
		formatBigIntDebug(s.G),
		formatBigIntDebug(s.H),
	)
}

func formatRangeDebug(r Range) string {
	return fmt.Sprintf(
		"{Inside=%t Lo=%s Hi=%s}",
		r.Inside,
		formatEndpointDebug(r.Lo),
		formatEndpointDebug(r.Hi),
	)
}

func formatEndpointDebug(ep Endpoint) string {
	return fmt.Sprintf(
		"{Value=%s Open=%t}",
		formatRationalDebug(ep.Value),
		ep.Open,
	)
}

func formatRationalDebug(r Rational) string {
	return fmt.Sprintf("%s/%s", r.Num().String(), r.Den().String())
}

func formatBigIntDebug(x *big.Int) string {
	if x == nil {
		return "nil"
	}
	return x.String()
}

// core/blft_range_helpers.go v2
