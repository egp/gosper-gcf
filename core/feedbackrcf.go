// core/feedbackrcf.go v1
package core

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
		panic("newFeedbackRCFBuffer: nil ensure callback")
	}
	return &feedbackRCFBuffer{
		ensure: ensure,
	}
}

func (b *feedbackRCFBuffer) Cursor() RCFStream {
	if b == nil {
		panic("feedbackRCFBuffer.Cursor: nil receiver")
	}
	return &feedbackRCFCursor{
		owner: b,
		next:  0,
	}
}

func (b *feedbackRCFBuffer) Append(term RCFTerm, rng Range) {
	if b == nil {
		panic("feedbackRCFBuffer.Append: nil receiver")
	}
	if b.eof {
		panic("feedbackRCFBuffer.Append: buffer already closed")
	}

	b.terms = append(b.terms, NewRCFTerm(term.A()))
	b.ranges = append(b.ranges, cloneRange(rng))
}

func (b *feedbackRCFBuffer) Close() {
	if b == nil {
		panic("feedbackRCFBuffer.Close: nil receiver")
	}
	b.eof = true
}

func (b *feedbackRCFBuffer) ensureRange(index int) {
	for len(b.ranges) <= index && !b.eof {
		beforeTerms := len(b.terms)
		beforeRanges := len(b.ranges)
		wasEOF := b.eof

		b.ensure(index)

		if len(b.terms) == beforeTerms && len(b.ranges) == beforeRanges && b.eof == wasEOF {
			panic("feedbackRCFBuffer.ensureRange: ensure callback made no progress")
		}
	}
}

func (b *feedbackRCFBuffer) ensureTerm(index int) {
	for len(b.terms) <= index && !b.eof {
		beforeTerms := len(b.terms)
		beforeRanges := len(b.ranges)
		wasEOF := b.eof

		b.ensure(index)

		if len(b.terms) == beforeTerms && len(b.ranges) == beforeRanges && b.eof == wasEOF {
			panic("feedbackRCFBuffer.ensureTerm: ensure callback made no progress")
		}
	}
}

func (c *feedbackRCFCursor) NextRCF() (RCFTerm, Status) {
	if c == nil || c.owner == nil {
		panic("feedbackRCFCursor.NextRCF: nil cursor")
	}

	c.owner.ensureTerm(c.next)

	if c.next >= len(c.owner.terms) {
		return NewRCFTerm(nil), StatusEOF
	}

	term := c.owner.terms[c.next]
	c.next++
	return NewRCFTerm(term.A()), StatusOK
}

func (c *feedbackRCFCursor) Range() Range {
	if c == nil || c.owner == nil {
		panic("feedbackRCFCursor.Range: nil cursor")
	}

	c.owner.ensureRange(c.next)

	if c.next >= len(c.owner.ranges) {
		panic("Range() is undefined on EOF RCF stream")
	}

	return cloneRange(c.owner.ranges[c.next])
}

// core/feedbackrcf.go v1
