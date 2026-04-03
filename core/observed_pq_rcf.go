// core/observed_pq_rcf.go v1
package core

import (
	"fmt"
	"math/big"
)

type observedRCFFromPQ struct {
	src PQStream
}

func newObservedRCFFromPQ(src PQStream) RCFStream {
	if src == nil {
		return newErrorRCFStream(fmt.Errorf("newObservedRCFFromPQ: %w", ErrNilObservedSource))
	}
	return &observedRCFFromPQ{src: src}
}

func (s *observedRCFFromPQ) NextRCF() (RCFTerm, Status, error) {
	if s == nil || s.src == nil {
		return NewRCFTerm(nil), StatusEOF, fmt.Errorf("observedRCFFromPQ.NextRCF: %w", ErrNilReceiver)
	}

	term, tail, status, err := s.src.NextPQ()
	if err != nil {
		return NewRCFTerm(nil), status, err
	}

	switch status {
	case StatusOK:
		if term.Q == nil || term.Q.Cmp(big.NewInt(1)) != 0 {
			return NewRCFTerm(nil), StatusEOF, fmt.Errorf("observedRCFFromPQ.NextRCF: non-RCF-compatible PQ term (%v,%v)", term.P, term.Q)
		}
		s.src = tail
		return NewRCFTerm(term.P), StatusOK, nil

	case StatusEOF:
		s.src = tail
		return NewRCFTerm(nil), StatusEOF, nil

	default:
		s.src = tail
		return NewRCFTerm(nil), status, fmt.Errorf("observedRCFFromPQ.NextRCF: invalid child status=%v", status)
	}
}

func (s *observedRCFFromPQ) CurrentInterval() (Interval, error) {
	if s == nil || s.src == nil {
		return Interval{}, fmt.Errorf("observedRCFFromPQ.CurrentInterval: %w", ErrNilReceiver)
	}
	return s.src.CurrentInterval()
}

func (s *observedRCFFromPQ) Range() (Range, error) {
	return s.CurrentInterval()
}

func isIdentityUnaryState(s blftState) bool {
	return coeffIsZero(s.A) &&
		!coeffIsZero(s.B) &&
		coeffIsZero(s.C) &&
		coeffIsZero(s.D) &&
		coeffIsZero(s.E) &&
		coeffIsZero(s.F) &&
		coeffIsZero(s.G) &&
		s.B != nil &&
		s.H != nil &&
		s.B.Cmp(s.H) == 0
}

// core/observed_pq_rcf.go v1
