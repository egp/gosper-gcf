// core/gcf_helpers.go v1
package core

import "math/big"

func collapseUnaryXEOF(coeffs blftState) Rational {
	if coeffs.F != nil && coeffs.F.Sign() != 0 {
		return NewRational(coeffs.B, coeffs.F)
	}
	return NewRational(coeffs.D, coeffs.H)
}

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

func collapseIndependentOfYToUnary(coeffs blftState) blftState {
	return blftState(coeffs.CollapseY()).Normalize()
}

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

func isIndependentOfY(coeffs blftState) bool {
	return coeffIsZero(coeffs.A) && coeffIsZero(coeffs.C) && coeffIsZero(coeffs.E) && coeffIsZero(coeffs.G)
}

func isIndependentOfX(coeffs blftState) bool {
	return coeffIsZero(coeffs.A) && coeffIsZero(coeffs.B) && coeffIsZero(coeffs.E) && coeffIsZero(coeffs.F)
}

func coeffIsZero(x *big.Int) bool {
	return x == nil || x.Sign() == 0
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

// core/gcf_helpers.go v1
