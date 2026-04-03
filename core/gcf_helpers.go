// core/gcf_helpers.go v9
package core

import "math/big"

func collapseUnaryXEOF(coeffs blftState) Rational {
	if coeffs.F != nil && coeffs.F.Sign() != 0 {
		return NewRational(coeffs.B, coeffs.F)
	}
	return NewRational(coeffs.D, coeffs.H)
}

// Independent of X means the BLFT already has the unary-in-Y form
// (C*y + D) / (G*y + H). In this codebase's unary slot layout, the active
// unary variable lives in B,D,F,H, so we remap C,D,G,H into B,D,F,H.
func collapseIndependentOfXToUnary(coeffs blftState) blftState {
	return blftState{
		A: big.NewInt(0),
		B: cloneBigInt(coeffs.C),
		C: big.NewInt(0),
		D: cloneBigInt(coeffs.D),
		E: big.NewInt(0),
		F: cloneBigInt(coeffs.G),
		G: big.NewInt(0),
		H: cloneBigInt(coeffs.H),
	}.Normalize()
}

// Independent of Y means the BLFT already has the unary-in-X form
// (B*x + D) / (F*x + H). In unary slot layout that stays in B,D,F,H.
func collapseIndependentOfYToUnary(coeffs blftState) blftState {
	return blftState{
		A: big.NewInt(0),
		B: cloneBigInt(coeffs.B),
		C: big.NewInt(0),
		D: cloneBigInt(coeffs.D),
		E: big.NewInt(0),
		F: cloneBigInt(coeffs.F),
		G: big.NewInt(0),
		H: cloneBigInt(coeffs.H),
	}.Normalize()
}

// Required by gcf_types.go for the lazy independent-of-X/Y constructor path.
// Do not remove without also updating that path and its WB tests.
func blftCoefficientsFromState(s blftState) BLFTCoefficients {
	return BLFTCoefficients{
		A: cloneBigInt(s.A),
		B: cloneBigInt(s.B),
		C: cloneBigInt(s.C),
		D: cloneBigInt(s.D),
		E: cloneBigInt(s.E),
		F: cloneBigInt(s.F),
		G: cloneBigInt(s.G),
		H: cloneBigInt(s.H),
	}
}

func isIndependentOfY(coeffs blftState) bool {
	return coeffIsZero(coeffs.A) && coeffIsZero(coeffs.C) && coeffIsZero(coeffs.E) && coeffIsZero(coeffs.G)
}

func isIndependentOfX(coeffs blftState) bool {
	return coeffIsZero(coeffs.A) && coeffIsZero(coeffs.B) && coeffIsZero(coeffs.E) && coeffIsZero(coeffs.F)
}

func coeffIsZero(x *big.Int) bool {
	return x == nil || x.Sign() == 0
}

// Gosper/HAKMEM-style X exhaustion leaves a unary transform in Y using A,B,E,F.
func collapseXEOFToUnary(coeffs blftState) blftState {
	if isIndependentOfX(coeffs) {
		return collapseIndependentOfXToUnary(coeffs)
	}
	return blftState{
		A: big.NewInt(0),
		B: cloneBigInt(coeffs.A),
		C: big.NewInt(0),
		D: cloneBigInt(coeffs.B),
		E: big.NewInt(0),
		F: cloneBigInt(coeffs.E),
		G: big.NewInt(0),
		H: cloneBigInt(coeffs.F),
	}.Normalize()
}

// Gosper/HAKMEM-style Y exhaustion leaves a unary transform in X using A,C,E,G.
func collapseYEOFToUnary(coeffs blftState) blftState {
	if isIndependentOfY(coeffs) {
		return collapseIndependentOfYToUnary(coeffs)
	}
	return blftState{
		A: big.NewInt(0),
		B: cloneBigInt(coeffs.A),
		C: big.NewInt(0),
		D: cloneBigInt(coeffs.C),
		E: big.NewInt(0),
		F: cloneBigInt(coeffs.E),
		G: big.NewInt(0),
		H: cloneBigInt(coeffs.G),
	}.Normalize()
}

func isExactIntegerRangeWithTerm(r Range, term RCFTerm) bool {
	if !r.Inside {
		return false
	}
	if r.Lo.Value.Cmp(r.Hi.Value) != 0 {
		return false
	}
	if r.Lo.Value.Den().Cmp(big.NewInt(1)) != 0 {
		return false
	}
	return r.Lo.Value.Num().Cmp(term.A()) == 0
}

func isEOFPQStream(x PQStream) bool {
	_, ok := x.(*eofPQStream)
	return ok
}

func exactRangeFromRational(r Rational) Range {
	return Range{
		Lo: Endpoint{
			Value: r,
			Open:  false,
		},
		Hi: Endpoint{
			Value: r,
			Open:  false,
		},
		Inside: true,
	}
}

func cloneRCFTerms(terms []RCFTerm) []RCFTerm {
	if terms == nil {
		return nil
	}
	out := make([]RCFTerm, len(terms))
	for i, term := range terms {
		out[i] = NewRCFTerm(term.A())
	}
	return out
}

// core/gcf_helpers.go v9
