// core/feedback_rcf_stream.go v3
package core

import (
	"fmt"
	"math/big"
)

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

func (s *feedbackRCFStream) Append(term RCFTerm, rng Range) error {
	if s == nil {
		return fmt.Errorf("feedbackRCFStream.Append: %w", ErrNilReceiver)
	}
	if s.closed {
		return fmt.Errorf("feedbackRCFStream.Append: %w", ErrAppendAfterClose)
	}

	s.terms = append(s.terms, cloneFeedbackRCFTerm(term))
	s.ranges = append(s.ranges, cloneFeedbackRCFRange(rng))
	return nil
}

func (s *feedbackRCFStream) Close() error {
	if s == nil {
		return fmt.Errorf("feedbackRCFStream.Close: %w", ErrNilReceiver)
	}
	s.closed = true
	return nil
}

func (s *feedbackRCFStream) NextRCF() (RCFTerm, Status, error) {
	if s == nil {
		return NewRCFTerm(nil), StatusEOF, fmt.Errorf("feedbackRCFStream.NextRCF: %w", ErrNilReceiver)
	}
	if s.next < len(s.terms) {
		term := cloneFeedbackRCFTerm(s.terms[s.next])
		s.next++
		return term, StatusOK, nil
	}
	if s.closed {
		return NewRCFTerm(big.NewInt(0)), StatusEOF, nil
	}
	return NewRCFTerm(nil), StatusEOF, fmt.Errorf("feedbackRCFStream.NextRCF: %w", ErrNoTermAvailableBeforeClose)
}

func (s *feedbackRCFStream) Range() (Range, error) {
	if s == nil {
		return Range{}, fmt.Errorf("feedbackRCFStream.Range: %w", ErrNilReceiver)
	}
	if s.next >= len(s.ranges) {
		return Range{}, fmt.Errorf("feedbackRCFStream.Range: %w", ErrNoRemainingSuffixRange)
	}
	return cloneFeedbackRCFRange(s.ranges[s.next]), nil
}

func cloneFeedbackRCFTerm(term RCFTerm) RCFTerm {
	return NewRCFTerm(term.A())
}

func cloneFeedbackRCFRange(r Range) Range {
	return Range{
		Lo:     cloneFeedbackRCFEndpoint(r.Lo),
		Hi:     cloneFeedbackRCFEndpoint(r.Hi),
		Inside: r.Inside,
		Kind_:  r.Kind_,
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

// core/feedback_rcf_stream.go v3
