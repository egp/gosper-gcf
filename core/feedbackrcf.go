// core/feedbackrcf.go v3
package core

import "fmt"

// feedbackRCFBuffer is an internal growable replay buffer for emitted regular-CF
// terms. The owner supplies an ensure callback that appends more produced prefix
// material or closes the buffer at EOF.
//
// This is the substrate needed for Ouroboros-style live sqrt: the public output
// stream is the source of truth, and internal feedback consumers read forked
// cursors over the same emitted prefix.
type feedbackRCFBuffer struct {
	ensure func(target int)
	terms  []RCFTerm
	ranges []Range
	eof    bool
}

type feedbackRCFCursor struct {
	owner *feedbackRCFBuffer
	next  int
}

func newFeedbackRCFBuffer(ensure func(target int)) *feedbackRCFBuffer {
	if ensure == nil {
		return nil
	}
	return &feedbackRCFBuffer{
		ensure: ensure,
	}
}

func (b *feedbackRCFBuffer) Cursor() RCFStream {
	if b == nil {
		return newErrorRCFStream(fmt.Errorf("feedbackRCFBuffer.Cursor: %w", ErrNilReceiver))
	}
	return &feedbackRCFCursor{
		owner: b,
		next:  0,
	}
}

func (b *feedbackRCFBuffer) Append(term RCFTerm, rng Range) error {
	if b == nil {
		return fmt.Errorf("feedbackRCFBuffer.Append: %w", ErrNilReceiver)
	}
	if b.eof {
		return fmt.Errorf("feedbackRCFBuffer.Append: %w", ErrAppendAfterClose)
	}
	b.terms = append(b.terms, NewRCFTerm(term.A()))
	b.ranges = append(b.ranges, cloneRange(rng))
	return nil
}

func (b *feedbackRCFBuffer) Close() error {
	if b == nil {
		return fmt.Errorf("feedbackRCFBuffer.Close: %w", ErrNilReceiver)
	}
	b.eof = true
	return nil
}

func (b *feedbackRCFBuffer) ensureRange(index int) error {
	if b == nil {
		return fmt.Errorf("feedbackRCFBuffer.ensureRange: %w", ErrNilReceiver)
	}
	for len(b.ranges) <= index && !b.eof {
		beforeTerms := len(b.terms)
		beforeRanges := len(b.ranges)
		wasEOF := b.eof
		b.ensure(index)
		if len(b.terms) == beforeTerms && len(b.ranges) == beforeRanges && b.eof == wasEOF {
			return fmt.Errorf("feedbackRCFBuffer.ensureRange: %w", ErrUnsupportedRangeCase)
		}
	}
	return nil
}

func (b *feedbackRCFBuffer) ensureTerm(index int) error {
	if b == nil {
		return fmt.Errorf("feedbackRCFBuffer.ensureTerm: %w", ErrNilReceiver)
	}
	for len(b.terms) <= index && !b.eof {
		beforeTerms := len(b.terms)
		beforeRanges := len(b.ranges)
		wasEOF := b.eof
		b.ensure(index)
		if len(b.terms) == beforeTerms && len(b.ranges) == beforeRanges && b.eof == wasEOF {
			return fmt.Errorf("feedbackRCFBuffer.ensureTerm: %w", ErrUnsupportedRangeCase)
		}
	}
	return nil
}

func (c *feedbackRCFCursor) NextRCF() (RCFTerm, Status, error) {
	if c == nil || c.owner == nil {
		return NewRCFTerm(nil), StatusEOF, fmt.Errorf("feedbackRCFCursor.NextRCF: %w", ErrNilReceiver)
	}
	if err := c.owner.ensureTerm(c.next); err != nil {
		return NewRCFTerm(nil), StatusEOF, err
	}
	if c.next >= len(c.owner.terms) {
		return NewRCFTerm(nil), StatusEOF, nil
	}
	term := c.owner.terms[c.next]
	c.next++
	return NewRCFTerm(term.A()), StatusOK, nil
}

func (c *feedbackRCFCursor) CurrentInterval() (Interval, error) {
	if c == nil || c.owner == nil {
		return Interval{}, fmt.Errorf("feedbackRCFCursor.CurrentInterval: %w", ErrNilReceiver)
	}
	if err := c.owner.ensureRange(c.next); err != nil {
		return Interval{}, err
	}
	if c.next >= len(c.owner.ranges) {
		return Interval{}, fmt.Errorf("feedbackRCFCursor.CurrentInterval: %w", ErrUndefinedRangeOnEOFStream)
	}
	return cloneRange(c.owner.ranges[c.next]), nil
}

func (c *feedbackRCFCursor) Range() (Range, error) {
	return c.CurrentInterval()
}

// core/feedbackrcf.go v3
