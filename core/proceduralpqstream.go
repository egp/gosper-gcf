// core/proceduralpqstream.go v1
package core

import "math/big"

type ProceduralPQNext func() (FinitePQStep, ProceduralPQNext, Status)

type proceduralPQStream struct {
	step FinitePQStep
	next ProceduralPQNext
}

type invalidPQStream struct {
	status Status
}

func NewProceduralPQStream(first FinitePQStep, next ProceduralPQNext) (PQStream, Status) {
	if !isValidFinitePQTerm(first.Term, true) {
		return nil, StatusInvalidInput
	}

	return &proceduralPQStream{
		step: cloneFinitePQStep(first),
		next: next,
	}, StatusOK
}

func (s *proceduralPQStream) NextPQ() (PQTerm, PQStream, Status) {
	if s == nil {
		panic("proceduralPQStream receiver is nil")
	}

	if s.next == nil {
		return clonePQTerm(s.step.Term), finitePQEOF, StatusOK
	}

	nextStep, nextNext, status := s.next()

	switch status {
	case StatusOK:
		if !isValidFinitePQTerm(nextStep.Term, false) {
			return clonePQTerm(s.step.Term), &invalidPQStream{status: StatusInvalidInput}, StatusOK
		}

		return clonePQTerm(s.step.Term), &proceduralPQStream{
			step: cloneFinitePQStep(nextStep),
			next: nextNext,
		}, StatusOK

	case StatusEOF:
		return clonePQTerm(s.step.Term), finitePQEOF, StatusOK

	default:
		return clonePQTerm(s.step.Term), &invalidPQStream{status: status}, StatusOK
	}
}

func (s *proceduralPQStream) Range() Range {
	if s == nil {
		panic("proceduralPQStream receiver is nil")
	}

	return cloneRange(s.step.Range)
}

func (s *invalidPQStream) NextPQ() (PQTerm, PQStream, Status) {
	return PQTerm{
		P: big.NewInt(0),
		Q: big.NewInt(0),
	}, s, s.status
}

func (s *invalidPQStream) Range() Range {
	panic("Range() is undefined on invalid PQStream")
}

// core/proceduralpqstream.go v1
