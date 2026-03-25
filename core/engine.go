// core/engine.go v1
package core

type unaryEngine interface {
	UnaryRange(xRange Range) Range
	CanEmitRCFTerm(r Range) (RCFTerm, bool)
	EmitUnary(term RCFTerm) unaryEngine
	IngestUnaryX(term PQTerm) unaryEngine
	CollapseUnaryEOF() Rational
}

type binaryEngine interface {
	BinaryRange(xRange, yRange Range) Range
	CanEmitRCFTerm(r Range) (RCFTerm, bool)
	EmitBinary(term RCFTerm) binaryEngine
	IngestBinaryX(term PQTerm) binaryEngine
	IngestBinaryY(term PQTerm) binaryEngine
	CollapseBinaryXEOF() unaryEngine
	CollapseBinaryYEOF() unaryEngine
	CollapseBinaryBothEOF() Rational
	IndependentOfX() bool
	IndependentOfY() bool
}

// core/engine.go v1
