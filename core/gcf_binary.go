// core/gcf_binary.go V11
package core

import (
	"fmt"
	"math/big"
)

func (g *GCF) nextBinaryRCF() (RCFTerm, Status, error) {
	for {
		// Detect runtime independence before any emission or ingest attempt.
		// This catches both initial independence and post-emit independence
		// on the following iteration, so we never lose an already-emitted term.
		if g.binary.engine.IndependentOfX() {
			g.unary = &unaryEvaluatorState{
				engine:    g.binary.engine.CollapseBinaryXEOF(),
				rectifier: NewRectifier(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1)),
				x:         g.binary.y,
			}
			g.binary = nil
			return g.NextRCF()
		}
		if g.binary.engine.IndependentOfY() {
			g.unary = &unaryEvaluatorState{
				engine:    g.binary.engine.CollapseBinaryYEOF(),
				rectifier: NewRectifier(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1)),
				x:         g.binary.x,
			}
			g.binary = nil
			return g.NextRCF()
		}

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
			g.unary = &unaryEvaluatorState{
				engine:    g.binary.engine.CollapseBinaryXEOF(),
				rectifier: NewRectifier(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1)), // identity transform per newSpec.md §4
				x:         g.binary.y,
			}
			g.binary = nil
			return g.NextRCF()
		}

		if isEOFPQStream(g.binary.y) {
			g.unary = &unaryEvaluatorState{
				engine:    g.binary.engine.CollapseBinaryYEOF(),
				rectifier: NewRectifier(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1)), // identity transform per newSpec.md §4
				x:         g.binary.x,
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

			// TWO-TIER (this iteration): always follow internal emission so term sequence / collapse timing is unchanged
			pq := PQTerm{
				P: cloneBigIntOrZero(term.A()),
				Q: big.NewInt(1),
			}
			g.binary.rectifier.Absorb(pq)
			g.binary.engine = g.binary.engine.EmitBinary(term)
			g.binary.rectifier.Emit(term.A())

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

// core/gcf_binary.go V11
