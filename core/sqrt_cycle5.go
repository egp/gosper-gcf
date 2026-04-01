// core/sqrt_cycle5.go v3
package core

import "fmt"

type sqrtApproximationPQStream struct {
	controller *sqrtController
}

func newSqrtApproximationPQStream(controller *sqrtController) *sqrtApproximationPQStream {
	return &sqrtApproximationPQStream{
		controller: controller,
	}
}

func (s *sqrtApproximationPQStream) NextPQ() (PQTerm, PQStream, Status, error) {
	if s == nil || s.controller == nil {
		return PQTerm{}, s, StatusEOF, fmt.Errorf("sqrtApproximationPQStream.NextPQ: %w", ErrNilReceiver)
	}

	active := s.controller.activeApproximation()
	term, tail, status, err := active.NextPQ()
	if err != nil {
		return term, s, status, err
	}

	if s.controller.hasOurorobosFeedback {
		s.controller.ourorobosApprox = tail
	} else {
		s.controller.seed = tail
	}
	return term, s, status, nil
}

func (s *sqrtApproximationPQStream) Range() (Range, error) {
	if s == nil || s.controller == nil {
		return Range{}, fmt.Errorf("sqrtApproximationPQStream.Range: %w", ErrNilReceiver)
	}
	return s.controller.activeApproximation().Range()
}

type sqrtObservedRefinementStream struct {
	controller *sqrtController
	proxy      *sqrtApproximationPQStream
	inner      *GCF
}

func newSqrtObservedRefinementStream(controller *sqrtController) *sqrtObservedRefinementStream {
	proxy := newSqrtApproximationPQStream(controller)

	var inner *GCF
	if controller == nil {
		inner = newObservedRCFGCF(newErrorRCFStream(fmt.Errorf("newSqrtObservedRefinementStream: %w", ErrNilReceiver)))
	} else {
		inner = controller.buildRefinement(proxy)
	}

	return &sqrtObservedRefinementStream{
		controller: controller,
		proxy:      proxy,
		inner:      inner,
	}
}

func (s *sqrtObservedRefinementStream) NextRCF() (RCFTerm, Status, error) {
	if s == nil || s.controller == nil || s.inner == nil {
		return NewRCFTerm(nil), StatusEOF, fmt.Errorf("sqrtObservedRefinementStream.NextRCF: %w", ErrNilReceiver)
	}

	rng, err := s.inner.Range()
	if err != nil {
		return NewRCFTerm(nil), StatusEOF, err
	}

	term, status, err := s.inner.NextRCF()
	if err != nil {
		return term, status, err
	}
	if status == StatusOK {
		if err := s.controller.feedCertifiedTerm(term, rng); err != nil {
			return term, status, err
		}
	}
	return term, status, nil
}

func (s *sqrtObservedRefinementStream) Range() (Range, error) {
	if s == nil || s.controller == nil || s.inner == nil {
		return Range{}, fmt.Errorf("sqrtObservedRefinementStream.Range: %w", ErrNilReceiver)
	}
	return s.inner.Range()
}

// core/sqrt_cycle5.go v3
