// core/blft.go v4
package core

import "math/big"

type blftState TransformCoefficients

func (s blftState) IngestX(term PQTerm) blftState {
	return blftState{
		A: mulAdd(s.A, term.P, s.C),
		B: mulAdd(s.B, term.P, s.D),
		C: mul(s.A, term.Q),
		D: mul(s.B, term.Q),
		E: mulAdd(s.E, term.P, s.G),
		F: mulAdd(s.F, term.P, s.H),
		G: mul(s.E, term.Q),
		H: mul(s.F, term.Q),
	}
}

func (s blftState) IngestY(term PQTerm) blftState {
	return blftState{
		A: mulAdd(s.A, term.P, s.B),
		B: mul(s.A, term.Q),
		C: mulAdd(s.C, term.P, s.D),
		D: mul(s.C, term.Q),
		E: mulAdd(s.E, term.P, s.F),
		F: mul(s.E, term.Q),
		G: mulAdd(s.G, term.P, s.H),
		H: mul(s.G, term.Q),
	}
}

func mul(x, y *big.Int) *big.Int {
	if x == nil || y == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Mul(x, y)
}

func mulAdd(x, y, z *big.Int) *big.Int {
	out := mul(x, y)
	if z != nil {
		out.Add(out, z)
	}
	return out
}

// core/blft.go v4
