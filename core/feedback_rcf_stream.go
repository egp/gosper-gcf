// core/feedback_rcf_stream.go v2
package core

import "math/big"

type feedbackRCFStream struct {
	terms  []RCFTerm
	ranges []Range
	next   int
	closed bool
}

func newFeedbackRCFStream() *feedbackRCFStream {
	return &feedbackRCFStream{
		terms:  make([]RCFTerm, 0, 8),
		ranges: make([]Range, 0, 8),
	}
}

func (s *feedbackRCFStream) Append(term RCFTerm, rng Range) {
	if s == nil {
		panic("feedbackRCFStream.Append: nil receiver")
	}
	if s.closed {
		panic("feedbackRCFStream.Append: append after Close")
	}

	s.terms = append(s.terms, cloneFeedbackRCFTerm(term))
	s.ranges = append(s.ranges, cloneFeedbackRCFRange(rng))
}

func (s *feedbackRCFStream) Close() {
	if s == nil {
		panic("feedbackRCFStream.Close: nil receiver")
	}
	s.closed = true
}

func (s *feedbackRCFStream) NextRCF() (RCFTerm, Status) {
	if s == nil {
		panic("feedbackRCFStream.NextRCF: nil receiver")
	}

	if s.next < len(s.terms) {
		term := cloneFeedbackRCFTerm(s.terms[s.next])
		s.next++
		return term, StatusOK
	}

	if s.closed {
		return NewRCFTerm(big.NewInt(0)), StatusEOF
	}

	panic("feedbackRCFStream.NextRCF: no term available before Close")
}

func (s *feedbackRCFStream) Range() Range {
	if s == nil {
		panic("feedbackRCFStream.Range: nil receiver")
	}
	if s.next >= len(s.ranges) {
		panic("feedbackRCFStream.Range: no remaining suffix range")
	}

	return cloneFeedbackRCFRange(s.ranges[s.next])
}

func cloneFeedbackRCFTerm(term RCFTerm) RCFTerm {
	return NewRCFTerm(term.A())
}

func cloneFeedbackRCFRange(r Range) Range {
	return Range{
		Lo:     cloneFeedbackRCFEndpoint(r.Lo),
		Hi:     cloneFeedbackRCFEndpoint(r.Hi),
		Inside: r.Inside,
	}
}

func cloneFeedbackRCFEndpoint(e Endpoint) Endpoint {
	return Endpoint{
		Value: cloneFeedbackRCFRational(e.Value),
		Open:  e.Open,
	}
}

func cloneFeedbackRCFRational(r Rational) Rational {
	return NewRational(r.Num(), r.Den())
}

// core/feedback_rcf_stream.go v2
