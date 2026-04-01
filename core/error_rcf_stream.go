// core/error_rcf_stream.go v1
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

func (s *errorRCFStream) Range() (Range, error) {
	return Range{}, s.err
}

// core/error_rcf_stream.go v1
