// core/blft.go V8
package core

import (
	"math/big"
)

type blftState BLFTCoefficients

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
	}.Normalize()
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
	}.Normalize()
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

func (s blftState) Normalize() blftState {
	g := new(big.Int).GCD(nil, nil, s.A, s.B)
	g = new(big.Int).GCD(g, nil, g, s.C)
	g = new(big.Int).GCD(g, nil, g, s.D)
	g = new(big.Int).GCD(g, nil, g, s.E)
	g = new(big.Int).GCD(g, nil, g, s.F)
	g = new(big.Int).GCD(g, nil, g, s.G)
	g = new(big.Int).GCD(g, nil, g, s.H)
	if g.Sign() == 0 {
		return s
	}
	return blftState{
		A: new(big.Int).Div(s.A, g),
		B: new(big.Int).Div(s.B, g),
		C: new(big.Int).Div(s.C, g),
		D: new(big.Int).Div(s.D, g),
		E: new(big.Int).Div(s.E, g),
		F: new(big.Int).Div(s.F, g),
		G: new(big.Int).Div(s.G, g),
		H: new(big.Int).Div(s.H, g),
	}
}

// core/blft.go V8
