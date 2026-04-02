// core/gcf_binary.go v2
package core

import "fmt"

func (g *GCF) nextBinaryRCF() (RCFTerm, Status, error) {
	for {
		if isEOFPQStream(g.binary.x) && isEOFPQStream(g.binary.y) {
			collapsed := g.binary.engine.CollapseBinaryBothEOF()
			terms, err := rcfTermsFromRationalChecked(collapsed)
			if err != nil {
				return NewRCFTerm(nil), StatusEOF, fmt.Errorf("nextBinaryRCF: collapse both EOF terms: %w", err)
			}
			g.terminal = &exactTerminalState{
				terms: cloneRCFTerms(terms),
				rng:   exactRangeFromRational(collapsed),
				next:  0,
			}
			g.binary = nil
			return g.NextRCF()
		}

		if isEOFPQStream(g.binary.x) {
			final, err := exactRationalFromUnaryEngine(g.binary.engine.CollapseBinaryXEOF(), g.binary.y)
			if err != nil {
				return NewRCFTerm(nil), StatusEOF, fmt.Errorf("nextBinaryRCF: collapse X EOF: %w", err)
			}
			terms, err := rcfTermsFromRationalChecked(final)
			if err != nil {
				return NewRCFTerm(nil), StatusEOF, fmt.Errorf("nextBinaryRCF: collapse X EOF terms: %w", err)
			}
			g.terminal = &exactTerminalState{
				terms: cloneRCFTerms(terms),
				rng:   exactRangeFromRational(final),
				next:  0,
			}
			g.binary = nil
			return g.NextRCF()
		}

		if isEOFPQStream(g.binary.y) {
			final, err := exactRationalFromUnaryEngine(g.binary.engine.CollapseBinaryYEOF(), g.binary.x)
			if err != nil {
				return NewRCFTerm(nil), StatusEOF, fmt.Errorf("nextBinaryRCF: collapse Y EOF: %w", err)
			}
			terms, err := rcfTermsFromRationalChecked(final)
			if err != nil {
				return NewRCFTerm(nil), StatusEOF, fmt.Errorf("nextBinaryRCF: collapse Y EOF terms: %w", err)
			}
			g.terminal = &exactTerminalState{
				terms: cloneRCFTerms(terms),
				rng:   exactRangeFromRational(final),
				next:  0,
			}
			g.binary = nil
			return g.NextRCF()
		}

		xRange, err := g.binary.x.Range()
		if err != nil {
			return NewRCFTerm(nil), StatusEOF, fmt.Errorf("nextBinaryRCF: left range: %w", err)
		}
		yRange, err := g.binary.y.Range()
		if err != nil {
			return NewRCFTerm(nil), StatusEOF, fmt.Errorf("nextBinaryRCF: right range: %w", err)
		}

		currentRange, err := g.binary.engine.BinaryRange(xRange, yRange)
		if err != nil {
			return NewRCFTerm(nil), StatusEOF, fmt.Errorf("nextBinaryRCF: binary range: %w", err)
		}

		if term, ok := g.binary.engine.CanEmitRCFTerm(currentRange); ok {
			if isExactIntegerRangeWithTerm(currentRange, term) {
				g.terminal = &exactTerminalState{
					terms: nil,
					rng:   exactRangeFromRational(RationalFromInt64(0)),
					next:  0,
				}
				g.binary = nil
				return term, StatusOK, nil
			}
			g.binary.engine = g.binary.engine.EmitBinary(term)
			return term, StatusOK, nil
		}

		if chooseIngestX(xRange, yRange) {
			term, tail, status, err := g.binary.x.NextPQ()
			if err != nil {
				return NewRCFTerm(nil), status, fmt.Errorf("nextBinaryRCF: left NextPQ: %w", err)
			}
			switch status {
			case StatusOK:
				g.binary.engine = g.binary.engine.IngestBinaryX(term)
				g.binary.x = tail
			case StatusEOF:
				g.binary.x = tail
			default:
				return NewRCFTerm(nil), status, fmt.Errorf("nextBinaryRCF: %w: status=%v", ErrInvalidBinaryLeftStatus, status)
			}
		} else {
			term, tail, status, err := g.binary.y.NextPQ()
			if err != nil {
				return NewRCFTerm(nil), status, fmt.Errorf("nextBinaryRCF: right NextPQ: %w", err)
			}
			switch status {
			case StatusOK:
				g.binary.engine = g.binary.engine.IngestBinaryY(term)
				g.binary.y = tail
			case StatusEOF:
				g.binary.y = tail
			default:
				return NewRCFTerm(nil), status, fmt.Errorf("nextBinaryRCF: %w: status=%v", ErrInvalidBinaryRightStatus, status)
			}
		}
	}
}

func (g *GCF) binaryRange() (Range, error) {
	if g.binary == nil {
		return Range{}, fmt.Errorf("binaryRange: %w", ErrMissingBinaryState)
	}

	switch {
	case isEOFPQStream(g.binary.x) && isEOFPQStream(g.binary.y):
		return exactRangeFromRational(g.binary.engine.CollapseBinaryBothEOF()), nil

	case isEOFPQStream(g.binary.x):
		yRange, err := g.binary.y.Range()
		if err != nil {
			return Range{}, fmt.Errorf("binaryRange: right range: %w", err)
		}
		return g.binary.engine.CollapseBinaryXEOF().UnaryRange(yRange)

	case isEOFPQStream(g.binary.y):
		xRange, err := g.binary.x.Range()
		if err != nil {
			return Range{}, fmt.Errorf("binaryRange: left range: %w", err)
		}
		return g.binary.engine.CollapseBinaryYEOF().UnaryRange(xRange)

	default:
		xRange, err := g.binary.x.Range()
		if err != nil {
			return Range{}, fmt.Errorf("binaryRange: left range: %w", err)
		}
		yRange, err := g.binary.y.Range()
		if err != nil {
			return Range{}, fmt.Errorf("binaryRange: right range: %w", err)
		}
		return g.binary.engine.BinaryRange(xRange, yRange)
	}
}

// core/gcf_binary.go v2
