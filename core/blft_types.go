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

// core/blft.go V8

// --- appended from core/blft_coefficients.go ---
// core/blftcoefficients.go v1

type BLFTCoefficients struct {
	A *big.Int
	B *big.Int
	C *big.Int
	D *big.Int
	E *big.Int
	F *big.Int
	G *big.Int
	H *big.Int
}

func cloneBLFTCoefficients(tc BLFTCoefficients) BLFTCoefficients {
	return BLFTCoefficients{
		A: cloneBigInt(tc.A),
		B: cloneBigInt(tc.B),
		C: cloneBigInt(tc.C),
		D: cloneBigInt(tc.D),
		E: cloneBigInt(tc.E),
		F: cloneBigInt(tc.F),
		G: cloneBigInt(tc.G),
		H: cloneBigInt(tc.H),
	}
}

// core/blftcoefficients.go v1
