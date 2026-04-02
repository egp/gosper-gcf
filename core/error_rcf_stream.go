// core/error_rcf_stream.go v2
package core

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
