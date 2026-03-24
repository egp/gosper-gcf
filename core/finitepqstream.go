// core/finitepqstream.go v1
package core

import "math/big"

type FinitePQStep struct {
	Term  PQTerm
	Range Range
}

type finitePQStream struct {
	step FinitePQStep
	tail PQStream
}

type eofPQStream struct{}

var finitePQEOF PQStream = &eofPQStream{}

func NewFinitePQStream(steps []FinitePQStep) (PQStream, Status) {
	tail := finitePQEOF

	for i := len(steps) - 1; i >= 0; i-- {
		if !isValidFinitePQTerm(steps[i].Term, i == 0) {
			return nil, StatusInvalidInput
		}

		tail = &finitePQStream{
			step: cloneFinitePQStep(steps[i]),
			tail: tail,
		}
	}

	return tail, StatusOK
}

func (s *finitePQStream) NextPQ() (PQTerm, PQStream, Status) {
	if s == nil {
		panic("finitePQStream receiver is nil")
	}

	return clonePQTerm(s.step.Term), s.tail, StatusOK
}

func (s *finitePQStream) Range() Range {
	if s == nil {
		panic("finitePQStream receiver is nil")
	}

	return cloneRange(s.step.Range)
}

func (s *eofPQStream) NextPQ() (PQTerm, PQStream, Status) {
	return PQTerm{
		P: big.NewInt(0),
		Q: big.NewInt(0),
	}, finitePQEOF, StatusEOF
}

func (s *eofPQStream) Range() Range {
	panic("Range() is undefined on EOF PQStream")
}

func isValidFinitePQTerm(term PQTerm, isFirst bool) bool {
	if term.P == nil || term.Q == nil {
		return false
	}

	if term.Q.Sign() == 0 {
		return false
	}

	if isFirst {
		return true
	}

	return term.P.Sign() > 0 && term.Q.Sign() > 0
}

func cloneFinitePQStep(step FinitePQStep) FinitePQStep {
	return FinitePQStep{
		Term:  clonePQTerm(step.Term),
		Range: cloneRange(step.Range),
	}
}

func clonePQTerm(term PQTerm) PQTerm {
	return PQTerm{
		P: cloneBigInt(term.P),
		Q: cloneBigInt(term.Q),
	}
}

func cloneRange(r Range) Range {
	return Range{
		Lo:     cloneEndpoint(r.Lo),
		Hi:     cloneEndpoint(r.Hi),
		Inside: r.Inside,
	}
}

func cloneEndpoint(e Endpoint) Endpoint {
	return Endpoint{
		Value: cloneRational(e.Value),
		Open:  e.Open,
	}
}

func cloneRational(r Rational) Rational {
	return Rational{
		num: cloneBigInt(r.num),
		den: cloneBigInt(r.den),
	}
}

func cloneBigInt(x *big.Int) *big.Int {
	if x == nil {
		return nil
	}
	return new(big.Int).Set(x)
}

// core/finitepqstream.go v1
