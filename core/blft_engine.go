// core/blft_engine.go v1
package core

func newBLFTState(coeffs BLFTCoefficients) blftState {
	return blftState(cloneBLFTCoefficients(coeffs))
}

func (s blftState) UnaryRange(xRange Range) Range {
	return s.CornerRange(
		xRange,
		exactRangeFromRational(RationalFromInt64(0)),
	)
}

func (s blftState) BinaryRange(xRange, yRange Range) Range {
	return s.CornerRange(xRange, yRange)
}

func (s blftState) EmitUnary(term RCFTerm) unaryEngine {
	return s.Emit(term)
}

func (s blftState) EmitBinary(term RCFTerm) binaryEngine {
	return s.Emit(term)
}

func (s blftState) IngestUnaryX(term PQTerm) unaryEngine {
	return s.IngestX(term)
}

func (s blftState) IngestBinaryX(term PQTerm) binaryEngine {
	return s.IngestX(term)
}

func (s blftState) IngestBinaryY(term PQTerm) binaryEngine {
	return s.IngestY(term)
}

func (s blftState) CollapseUnaryEOF() Rational {
	return collapseUnaryXEOF(s)
}

func (s blftState) CollapseBinaryXEOF() unaryEngine {
	return collapseXEOFToUnary(s)
}

func (s blftState) CollapseBinaryYEOF() unaryEngine {
	return collapseYEOFToUnary(s)
}

func (s blftState) CollapseBinaryBothEOF() Rational {
	return s.CollapseToRational()
}

func (s blftState) IndependentOfX() bool {
	return isIndependentOfX(s)
}

func (s blftState) IndependentOfY() bool {
	return isIndependentOfY(s)
}

// core/blft_engine.go v1
