// core/rcf_to_pq.go v3
package core

import (
	"fmt"
	"math/big"
)

type rcfAsPQStream struct {
	src RCFStream
}

func PQStreamFromRCF(src RCFStream) PQStream {
	if src == nil {
		return newErrorPQStream(fmt.Errorf("PQStreamFromRCF: %w", ErrNilObservedSource))
	}
	return &rcfAsPQStream{src: src}
}

func (s *rcfAsPQStream) NextPQ() (PQTerm, PQStream, Status, error) {
	term, status, err := s.src.NextRCF()
	if err != nil {
		return PQTerm{
			P: big.NewInt(0),
			Q: big.NewInt(0),
		}, s, status, err
	}
	if status != StatusOK {
		return PQTerm{
			P: big.NewInt(0),
			Q: big.NewInt(0),
		}, s, status, nil
	}
	return PQTerm{
		P: term.A(),
		Q: big.NewInt(1),
	}, s, StatusOK, nil
}

func (s *rcfAsPQStream) CurrentInterval() (Interval, error) {
	if s == nil || s.src == nil {
		return Interval{}, fmt.Errorf("rcfAsPQStream.CurrentInterval: %w", ErrNilReceiver)
	}
	return s.src.CurrentInterval()
}

func (s *rcfAsPQStream) Range() (Range, error) {
	return s.CurrentInterval()
}

// core/rcf_to_pq.go v3
