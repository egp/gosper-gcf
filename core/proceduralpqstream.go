// core/proceduralpqstream.go v3
package core

import (
	"fmt"
	"math/big"
)

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

func (s *proceduralPQStream) NextPQ() (PQTerm, PQStream, Status, error) {
	if s == nil {
		return PQTerm{}, finitePQEOF, StatusEOF, fmt.Errorf("proceduralPQStream.NextPQ: %w", ErrNilReceiver)
	}

	if s.next == nil {
		return clonePQTerm(s.step.Term), finitePQEOF, StatusOK, nil
	}

	nextStep, nextNext, status := s.next()
	switch status {
	case StatusOK:
		if !isValidFinitePQTerm(nextStep.Term, false) {
			return clonePQTerm(s.step.Term), &invalidPQStream{status: StatusInvalidInput}, StatusOK, nil
		}
		return clonePQTerm(s.step.Term), &proceduralPQStream{
			step: cloneFinitePQStep(nextStep),
			next: nextNext,
		}, StatusOK, nil

	case StatusEOF:
		return clonePQTerm(s.step.Term), finitePQEOF, StatusOK, nil

	default:
		return clonePQTerm(s.step.Term), &invalidPQStream{status: status}, StatusOK, nil
	}
}

func (s *proceduralPQStream) CurrentInterval() (Interval, error) {
	if s == nil {
		return Interval{}, fmt.Errorf("proceduralPQStream.CurrentInterval: %w", ErrNilReceiver)
	}
	return cloneRange(s.step.Range), nil
}

func (s *proceduralPQStream) Range() (Range, error) {
	return s.CurrentInterval()
}

func (s *invalidPQStream) NextPQ() (PQTerm, PQStream, Status, error) {
	return PQTerm{
		P: big.NewInt(0),
		Q: big.NewInt(0),
	}, s, s.status, nil
}

func (s *invalidPQStream) CurrentInterval() (Interval, error) {
	return Interval{}, fmt.Errorf("invalidPQStream.CurrentInterval: %w", ErrUndefinedRangeOnBadStream)
}

func (s *invalidPQStream) Range() (Range, error) {
	return s.CurrentInterval()
}

// core/proceduralpqstream.go v3
