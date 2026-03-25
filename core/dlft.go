// core/dlft.go v4
package core

import "math/big"

type dlftState DLFTCoefficients

func newDLFTState(coeffs DLFTCoefficients) dlftState {
	return dlftState(cloneDLFTCoefficients(coeffs))
}

func (s dlftState) IngestX(term PQTerm) dlftState {
	p := term.P
	q := term.Q

	pp := mul(p, p)
	qq := mul(q, q)

	twoAP := mul(big.NewInt(2), mul(s.A, p))
	twoDP := mul(big.NewInt(2), mul(s.D, p))

	return dlftState{
		A: add3(mul(s.A, pp), mul(s.B, p), cloneBigIntOrZero(s.C)),
		B: mul(add2(twoAP, cloneBigIntOrZero(s.B)), q),
		C: mul(s.A, qq),
		D: add3(mul(s.D, pp), mul(s.E, p), cloneBigIntOrZero(s.F)),
		E: mul(add2(twoDP, cloneBigIntOrZero(s.E)), q),
		F: mul(s.D, qq),
	}
}

func (s dlftState) Emit(term RCFTerm) dlftState {
	n := term.A()

	return dlftState{
		A: cloneBigIntOrZero(s.D),
		B: cloneBigIntOrZero(s.E),
		C: cloneBigIntOrZero(s.F),
		D: subMul(s.A, n, s.D),
		E: subMul(s.B, n, s.E),
		F: subMul(s.C, n, s.F),
	}
}

func (s dlftState) CandidateRange(xRange Range) Range {
	if !xRange.Inside {
		panic("CandidateRange currently supports only inside ranges")
	}

	points := []Rational{
		xRange.Lo.Value,
		xRange.Hi.Value,
	}

	// Add rational denominator roots, if any.
	points = append(points, rationalRootsQuadratic(
		cloneBigIntOrZero(s.D),
		cloneBigIntOrZero(s.E),
		cloneBigIntOrZero(s.F),
		xRange,
	)...)

	// Add rational critical points, if any.
	critA := new(big.Int).Sub(
		mul(s.A, s.E),
		mul(s.B, s.D),
	)
	critB := mul(
		big.NewInt(2),
		new(big.Int).Sub(
			mul(s.A, s.F),
			mul(s.C, s.D),
		),
	)
	critC := new(big.Int).Sub(
		mul(s.B, s.F),
		mul(s.C, s.E),
	)

	points = append(points, rationalRootsQuadratic(
		critA,
		critB,
		critC,
		xRange,
	)...)

	values := make([]Rational, 0, len(points))
	sawZeroDen := false
	sawPosDen := false
	sawNegDen := false

	for _, x := range dedupeRationals(points) {
		num, den := evalDLFTNumDenAtPoint(s, x)

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

func (s dlftState) CollapseToRational() Rational {
	switch {
	case s.D != nil && s.D.Sign() != 0:
		return NewRational(s.A, s.D)
	case s.E != nil && s.E.Sign() != 0:
		return NewRational(s.B, s.E)
	default:
		return NewRational(s.C, s.F)
	}
}

func (s dlftState) UnaryRange(xRange Range) Range {
	return s.CandidateRange(xRange)
}

func (s dlftState) CanEmitRCFTerm(r Range) (RCFTerm, bool) {
	return canEmitRCFTermFromRange(r)
}

func (s dlftState) EmitUnary(term RCFTerm) unaryEngine {
	return s.Emit(term)
}

func (s dlftState) IngestUnaryX(term PQTerm) unaryEngine {
	return s.IngestX(term)
}

func (s dlftState) CollapseUnaryEOF() Rational {
	return s.CollapseToRational()
}

func add2(x, y *big.Int) *big.Int {
	out := cloneBigIntOrZero(x)
	out.Add(out, cloneBigIntOrZero(y))
	return out
}

func add3(x, y, z *big.Int) *big.Int {
	out := cloneBigIntOrZero(x)
	out.Add(out, cloneBigIntOrZero(y))
	out.Add(out, cloneBigIntOrZero(z))
	return out
}

func evalDLFTNumDenAtPoint(s dlftState, x Rational) (*big.Int, *big.Int) {
	xn := x.Num()
	xd := x.Den()

	xn2 := mul(xn, xn)
	xd2 := mul(xd, xd)

	num := big.NewInt(0)
	num.Add(num, scaledDLFTQuadratic(s.A, xn2))
	num.Add(num, scaledDLFTLinear(s.B, xn, xd))
	num.Add(num, scaledDLFTConstant(s.C, xd2))

	den := big.NewInt(0)
	den.Add(den, scaledDLFTQuadratic(s.D, xn2))
	den.Add(den, scaledDLFTLinear(s.E, xn, xd))
	den.Add(den, scaledDLFTConstant(s.F, xd2))

	return num, den
}

func scaledDLFTQuadratic(coeff, xn2 *big.Int) *big.Int {
	if coeff == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Mul(coeff, xn2)
}

func scaledDLFTLinear(coeff, xn, xd *big.Int) *big.Int {
	if coeff == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Mul(coeff, new(big.Int).Mul(xn, xd))
}

func scaledDLFTConstant(coeff, xden2 *big.Int) *big.Int {
	if coeff == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Mul(coeff, xden2)
}

func rationalRootsQuadratic(a, b, c *big.Int, xRange Range) []Rational {
	if a.Sign() == 0 {
		return rationalRootsLinear(b, c, xRange)
	}

	// discriminant = b^2 - 4ac
	disc := new(big.Int).Sub(
		mul(b, b),
		mul(big.NewInt(4), mul(a, c)),
	)
	if disc.Sign() < 0 {
		return nil
	}

	sqrtDisc, ok := perfectSquareRoot(disc)
	if !ok {
		return nil
	}

	twoA := mul(big.NewInt(2), a)
	negB := new(big.Int).Neg(cloneBigIntOrZero(b))

	r1 := NewRational(new(big.Int).Sub(cloneBigIntOrZero(negB), sqrtDisc), twoA)
	r2 := NewRational(new(big.Int).Add(cloneBigIntOrZero(negB), sqrtDisc), twoA)

	out := make([]Rational, 0, 2)
	if rationalInClosedInsideRange(r1, xRange) {
		out = append(out, r1)
	}
	if r2.Cmp(r1) != 0 && rationalInClosedInsideRange(r2, xRange) {
		out = append(out, r2)
	}
	return out
}

func rationalRootsLinear(a, b *big.Int, xRange Range) []Rational {
	if a.Sign() == 0 {
		return nil
	}

	root := NewRational(new(big.Int).Neg(cloneBigIntOrZero(b)), a)
	if rationalInClosedInsideRange(root, xRange) {
		return []Rational{root}
	}
	return nil
}

func rationalInClosedInsideRange(x Rational, r Range) bool {
	if !r.Inside {
		return false
	}
	return x.Cmp(r.Lo.Value) >= 0 && x.Cmp(r.Hi.Value) <= 0
}

func dedupeRationals(xs []Rational) []Rational {
	out := make([]Rational, 0, len(xs))
	for _, x := range xs {
		found := false
		for _, y := range out {
			if x.Cmp(y) == 0 {
				found = true
				break
			}
		}
		if !found {
			out = append(out, x)
		}
	}
	return out
}

func perfectSquareRoot(n *big.Int) (*big.Int, bool) {
	if n.Sign() < 0 {
		return nil, false
	}
	if n.Sign() == 0 {
		return big.NewInt(0), true
	}

	one := big.NewInt(1)
	two := big.NewInt(2)

	lo := big.NewInt(0)
	hi := cloneBigIntOrZero(n)

	for lo.Cmp(hi) <= 0 {
		mid := new(big.Int).Add(lo, hi)
		mid.Quo(mid, two)

		sq := new(big.Int).Mul(mid, mid)
		cmp := sq.Cmp(n)

		if cmp == 0 {
			return mid, true
		}
		if cmp < 0 {
			lo = new(big.Int).Add(mid, one)
		} else {
			hi = new(big.Int).Sub(mid, one)
		}
	}

	return nil, false
}

// core/dlft.go v4
