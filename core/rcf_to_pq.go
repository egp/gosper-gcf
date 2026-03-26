// core/rcf_to_pq.go v1
package core

import "math/big"

type rcfAsPQStream struct {
	src RCFStream
}

func PQStreamFromRCF(src RCFStream) PQStream {
	if src == nil {
		panic("PQStreamFromRCF: nil source")
	}
	return &rcfAsPQStream{src: src}
}

func (s *rcfAsPQStream) NextPQ() (PQTerm, PQStream, Status) {
	term, status := s.src.NextRCF()
	if status != StatusOK {
		return PQTerm{
			P: big.NewInt(0),
			Q: big.NewInt(0),
		}, s, status
	}

	return PQTerm{
		P: term.A(),
		Q: big.NewInt(1),
	}, s, StatusOK
}

func (s *rcfAsPQStream) Range() Range {
	return s.src.Range()
}

// core/rcf_to_pq.go v1
