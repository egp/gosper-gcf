// core/gcf_unary.go V12
package core

import (
	"errors"
	"fmt"
	"math/big"
)

func (g *GCF) nextUnaryRCF() (RCFTerm, Status, error) {
	for {
		if isEOFPQStream(g.unary.x) {
			collapsed := g.unary.engine.CollapseUnaryEOF()
			terms, err := rcfTermsFromRationalChecked(collapsed)
			if err != nil {
				return NewRCFTerm(nil), StatusEOF, fmt.Errorf("nextUnaryRCF: collapse EOF terms: %w", err)
			}
			g.terminal = &exactTerminalState{
				terms: cloneRCFTerms(terms),
				rng:   exactRangeFromRational(collapsed),
				next:  0,
			}
			g.unary = nil
			return g.NextRCF()
		}

		xRange, err := g.unary.x.Range()
		if err != nil {
			return NewRCFTerm(nil), StatusEOF, fmt.Errorf("nextUnaryRCF: child range: %w", err)
		}

		currentRange, err := g.unary.engine.UnaryRange(xRange)
		if err != nil {
			if errors.Is(err, ErrUnsupportedRangeCase) && canAdvancePastUnsupportedUnaryRange(g.unary.x) {
				term, tail, status, nextErr := g.unary.x.NextPQ()
				if nextErr != nil {
					return NewRCFTerm(nil), status, fmt.Errorf("nextUnaryRCF: advance past unsupported range NextPQ: %w", nextErr)
				}
				switch status {
				case StatusOK:
					g.unary.engine = g.unary.engine.IngestUnaryX(term)
					g.unary.x = tail
					continue
				case StatusEOF:
					g.unary.x = tail
				default:
					return NewRCFTerm(nil), status, fmt.Errorf("nextUnaryRCF: %w: status=%v", ErrInvalidUnaryInputStatus, status)
				}
			}
			return NewRCFTerm(nil), StatusEOF, fmt.Errorf("nextUnaryRCF: unary range: %w", err)
		}

		if term, ok := g.unary.engine.CanEmitRCFTerm(currentRange); ok {
			if isExactIntegerRangeWithTerm(currentRange, term) {
				g.terminal = &exactTerminalState{
					terms: nil,
					rng:   exactRangeFromRational(RationalFromInt64(0)),
					next:  0,
				}
				g.unary = nil
				return term, StatusOK, nil
			}

			// TWO-TIER (this iteration): always follow internal emission so term sequence / collapse timing is unchanged
			pq := PQTerm{
				P: cloneBigIntOrZero(term.A()),
				Q: big.NewInt(1),
			}
			g.unary.rectifier.Absorb(pq)
			g.unary.engine = g.unary.engine.EmitUnary(term)
			g.unary.rectifier.Emit(term.A())
			return term, StatusOK, nil
		}

		term, tail, status, err := g.unary.x.NextPQ()
		if err != nil {
			return NewRCFTerm(nil), status, fmt.Errorf("nextUnaryRCF: NextPQ: %w", err)
		}

		switch status {
		case StatusOK:
			g.unary.engine = g.unary.engine.IngestUnaryX(term)
			g.unary.x = tail
		case StatusEOF:
			g.unary.x = tail
		default:
			return NewRCFTerm(nil), status, fmt.Errorf("nextUnaryRCF: %w: status=%v", ErrInvalidUnaryInputStatus, status)
		}
	}
}

func exactRationalFromUnaryEngine(engine unaryEngine, stream PQStream) (Rational, error) {
	current := engine
	input := stream

	for {
		term, tail, status, err := input.NextPQ()
		if err != nil {
			return Rational{}, fmt.Errorf("exactRationalFromUnaryEngine: NextPQ: %w", err)
		}

		switch status {
		case StatusOK:
			current = current.IngestUnaryX(term)
			input = tail
		case StatusEOF:
			return current.CollapseUnaryEOF(), nil
		default:
			return Rational{}, fmt.Errorf("exactRationalFromUnaryEngine: %w: status=%v", ErrInvalidUnaryInputStatus, status)
		}
	}
}

func canAdvancePastUnsupportedUnaryRange(x PQStream) bool {
	_, ok := x.(*rcfAsPQStream)
	return ok
}

// core/gcf_unary.go V12
