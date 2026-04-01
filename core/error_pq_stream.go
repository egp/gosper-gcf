// core/error_pq_stream.go v1
package core

import "math/big"

type errorPQStream struct {
	err error
}

func newErrorPQStream(err error) *errorPQStream {
	return &errorPQStream{err: err}
}

func (s *errorPQStream) NextPQ() (PQTerm, PQStream, Status, error) {
	return PQTerm{
		P: big.NewInt(0),
		Q: big.NewInt(0),
	}, s, StatusEOF, s.err
}

func (s *errorPQStream) Range() (Range, error) {
	return Range{}, s.err
}

// core/error_pq_stream.go v1
