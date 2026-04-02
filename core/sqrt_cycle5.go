// core/sqrt_cycle5.go v4
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

func (s *sqrtApproximationPQStream) CurrentInterval() (Interval, error) {
	if s == nil || s.controller == nil {
		return Interval{}, fmt.Errorf("sqrtApproximationPQStream.CurrentInterval: %w", ErrNilReceiver)
	}
	return s.controller.activeApproximation().CurrentInterval()
}

func (s *sqrtApproximationPQStream) Range() (Range, error) {
	return s.CurrentInterval()
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

	rng, err := s.inner.CurrentInterval()
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

func (s *sqrtObservedRefinementStream) CurrentInterval() (Interval, error) {
	if s == nil || s.controller == nil || s.inner == nil {
		return Interval{}, fmt.Errorf("sqrtObservedRefinementStream.CurrentInterval: %w", ErrNilReceiver)
	}
	return s.inner.CurrentInterval()
}

func (s *sqrtObservedRefinementStream) Range() (Range, error) {
	return s.CurrentInterval()
}

// core/sqrt_cycle5.go v4
