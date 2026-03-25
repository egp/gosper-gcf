// core/dlft.go v1
package core

type dlftState DLFTCoefficients

func newDLFTState(coeffs DLFTCoefficients) dlftState {
	return dlftState(cloneDLFTCoefficients(coeffs))
}

// Shape-only stub.
// Later this will implement the real DLFT ingestion formula for x = p + q/x'.
func (s dlftState) IngestX(term PQTerm) dlftState {
	_ = term
	return s
}

// Shape-only stub.
// Later this will implement the DLFT emission update z' = 1 / (z - n).
func (s dlftState) Emit(term RCFTerm) dlftState {
	_ = term
	return s
}

// Shape-only stub.
// Later this will evaluate candidate unary range, including endpoint / pole / critical-point logic.
func (s dlftState) CandidateRange(xRange Range) Range {
	_ = xRange
	return exactRangeFromRational(RationalFromInt64(0))
}

// Shape-only stub.
// Later this will collapse by highest surviving degree: A/D, else B/E, else C/F.
func (s dlftState) CollapseToRational() Rational {
	return RationalFromInt64(0)
}

// core/dlft.go v1
