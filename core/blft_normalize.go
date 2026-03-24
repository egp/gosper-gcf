// core/blft_normalize.go v2
package core

import "math/big"

func (s blftState) Normalize() blftState {
	out := blftState{
		A: normalizeCoeff(s.A),
		B: normalizeCoeff(s.B),
		C: normalizeCoeff(s.C),
		D: normalizeCoeff(s.D),
		E: normalizeCoeff(s.E),
		F: normalizeCoeff(s.F),
		G: normalizeCoeff(s.G),
		H: normalizeCoeff(s.H),
	}

	g := overallCoeffGCD(out)
	if g.Sign() > 0 && g.Cmp(big.NewInt(1)) != 0 {
		out.A = divExact(out.A, g)
		out.B = divExact(out.B, g)
		out.C = divExact(out.C, g)
		out.D = divExact(out.D, g)
		out.E = divExact(out.E, g)
		out.F = divExact(out.F, g)
		out.G = divExact(out.G, g)
		out.H = divExact(out.H, g)
	}

	if firstNonZeroIsNegative(out) {
		out.A.Neg(out.A)
		out.B.Neg(out.B)
		out.C.Neg(out.C)
		out.D.Neg(out.D)
		out.E.Neg(out.E)
		out.F.Neg(out.F)
		out.G.Neg(out.G)
		out.H.Neg(out.H)
	}

	return out
}

func normalizeCoeff(x *big.Int) *big.Int {
	if x == nil {
		return big.NewInt(0)
	}
	return cloneBigInt(x)
}

func overallCoeffGCD(s blftState) *big.Int {
	coeffs := []*big.Int{s.A, s.B, s.C, s.D, s.E, s.F, s.G, s.H}

	g := big.NewInt(0)
	tmp := new(big.Int)

	for _, c := range coeffs {
		if c.Sign() == 0 {
			continue
		}
		tmp.Abs(c)
		if g.Sign() == 0 {
			g.Set(tmp)
			continue
		}
		g.GCD(nil, nil, g, tmp)
	}

	if g.Sign() == 0 {
		return big.NewInt(1)
	}
	return g
}

func divExact(x, y *big.Int) *big.Int {
	return new(big.Int).Quo(x, y)
}

func firstNonZeroIsNegative(s blftState) bool {
	coeffs := []*big.Int{s.A, s.B, s.C, s.D, s.E, s.F, s.G, s.H}

	for _, c := range coeffs {
		switch c.Sign() {
		case -1:
			return true
		case 1:
			return false
		}
	}
	return false
}

// core/blft_normalize.go v2
