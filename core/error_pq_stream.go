// core/error_pq_stream.go v2
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

func (s *errorPQStream) CurrentInterval() (Interval, error) {
	return Interval{}, s.err
}

func (s *errorPQStream) Range() (Range, error) {
	return s.CurrentInterval()
}

// core/error_pq_stream.go v2
