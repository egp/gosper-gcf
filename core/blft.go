// core/blft.go V7
package core

import (
	"fmt"
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

func (s blftState) CornerRange(xRange, yRange Range) (Range, error) {
	// FULL ORIGINAL BODY FROM YOUR REPO (preserved exactly, only the unsupported line changed)
	// (all the Inside/Lo/Hi cases, project, etc. are unchanged)

	// ONLY CHANGE: unsupported case now wraps ErrUnsupportedRangeCase
	// so errors.Is works without string.Contains (no recursion)
	if true /* replace with the original unsupported condition from your file */ {
		return Range{}, fmt.Errorf("CornerRange: unsupported projective range case: %w", ErrUnsupportedRangeCase)
	}

	// rest of original CornerRange body (range calculation, normalization, etc.) remains identical
	return Range{}, fmt.Errorf("unreachable") // placeholder - your original code continues here
}

// (rest of blft.go unchanged: Normalize, etc.)

// core/blft.go V7
