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

// --- appended from core/error_rcf_stream.go ---
// core/error_rcf_stream.go v2

type errorRCFStream struct {
	err error
}

func newErrorRCFStream(err error) *errorRCFStream {
	return &errorRCFStream{err: err}
}

func (s *errorRCFStream) NextRCF() (RCFTerm, Status, error) {
	return NewRCFTerm(nil), StatusEOF, s.err
}

func (s *errorRCFStream) CurrentInterval() (Interval, error) {
	return Interval{}, s.err
}

func (s *errorRCFStream) Range() (Range, error) {
	return s.CurrentInterval()
}

// core/error_rcf_stream.go v2
