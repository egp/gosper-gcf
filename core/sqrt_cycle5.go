// core/sqrt_cycle5.go v2
package core

type sqrtApproximationPQStream struct {
	controller *sqrtController
}

func newSqrtApproximationPQStream(controller *sqrtController) *sqrtApproximationPQStream {
	if controller == nil {
		panic("newSqrtApproximationPQStream: nil controller")
	}
	return &sqrtApproximationPQStream{
		controller: controller,
	}
}

func (s *sqrtApproximationPQStream) NextPQ() (PQTerm, PQStream, Status) {
	if s == nil || s.controller == nil {
		panic("sqrtApproximationPQStream.NextPQ: nil receiver")
	}

	active := s.controller.activeApproximation()
	term, tail, status := active.NextPQ()

	if s.controller.hasOurorobosFeedback {
		s.controller.ourorobosApprox = tail
	} else {
		s.controller.seed = tail
	}

	return term, s, status
}

func (s *sqrtApproximationPQStream) Range() Range {
	if s == nil || s.controller == nil {
		panic("sqrtApproximationPQStream.Range: nil receiver")
	}

	return s.controller.activeApproximation().Range()
}

type sqrtObservedRefinementStream struct {
	controller *sqrtController
	proxy      *sqrtApproximationPQStream
	inner      *GCF
}

func newSqrtObservedRefinementStream(controller *sqrtController) *sqrtObservedRefinementStream {
	if controller == nil {
		panic("newSqrtObservedRefinementStream: nil controller")
	}

	proxy := newSqrtApproximationPQStream(controller)

	return &sqrtObservedRefinementStream{
		controller: controller,
		proxy:      proxy,
		inner:      controller.buildRefinement(proxy),
	}
}

func (s *sqrtObservedRefinementStream) NextRCF() (RCFTerm, Status) {
	if s == nil || s.controller == nil || s.inner == nil {
		panic("sqrtObservedRefinementStream.NextRCF: nil receiver")
	}

	rng := s.inner.Range()
	term, status := s.inner.NextRCF()
	if status == StatusOK {
		s.controller.feedCertifiedTerm(term, rng)
	}

	return term, status
}

func (s *sqrtObservedRefinementStream) Range() Range {
	if s == nil || s.controller == nil || s.inner == nil {
		panic("sqrtObservedRefinementStream.Range: nil receiver")
	}

	return s.inner.Range()
}

// core/sqrt_cycle5.go v2
