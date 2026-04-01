// core/finitepqstream.go v2
package core

import (
	"fmt"
	"math/big"
)

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

func (s *finitePQStream) NextPQ() (PQTerm, PQStream, Status, error) {
	if s == nil {
		return PQTerm{}, finitePQEOF, StatusEOF, fmt.Errorf("finitePQStream.NextPQ: %w", ErrNilReceiver)
	}
	return clonePQTerm(s.step.Term), s.tail, StatusOK, nil
}

func (s *finitePQStream) Range() (Range, error) {
	if s == nil {
		return Range{}, fmt.Errorf("finitePQStream.Range: %w", ErrNilReceiver)
	}
	return cloneRange(s.step.Range), nil
}

func (s *eofPQStream) NextPQ() (PQTerm, PQStream, Status, error) {
	return PQTerm{
		P: big.NewInt(0),
		Q: big.NewInt(0),
	}, finitePQEOF, StatusEOF, nil
}

func (s *eofPQStream) Range() (Range, error) {
	return Range{}, fmt.Errorf("eofPQStream.Range: %w", ErrUndefinedRangeOnEOFStream)
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
		Lo:     Endpoint{Value: NewRational(r.Lo.Value.Num(), r.Lo.Value.Den()), Open: r.Lo.Open},
		Hi:     Endpoint{Value: NewRational(r.Hi.Value.Num(), r.Hi.Value.Den()), Open: r.Hi.Open},
		Inside: r.Inside,
	}
}

// core/finitepqstream.go v2
