// core/gcf.go v14
package core

import (
	"fmt"
	"math/big"
)

type exactTerminalState struct {
	terms []RCFTerm
	rng   Range
	next  int
}

type unaryEvaluatorState struct {
	engine unaryEngine
	x      PQStream
}

type binaryEvaluatorState struct {
	engine binaryEngine
	x      PQStream
	y      PQStream
}

type GCF struct {
	coeffs   BLFTCoefficients
	cfg      Config
	x        PQStream
	y        PQStream
	stream   RCFStream
	terminal *exactTerminalState
	unary    *unaryEvaluatorState
	binary   *binaryEvaluatorState
}

func NewExactTerminalGCF(terms []RCFTerm, rng Range) *GCF {
	return newExactTerminalGCFWithResolvedConfig(terms, rng, DefaultConfig())
}

func NewExactTerminalGCFWithConfig(terms []RCFTerm, rng Range, cfg Config) *GCF {
	return newExactTerminalGCFWithResolvedConfig(terms, rng, cfg)
}

func NewGCF0(coeffs BLFTCoefficients) *GCF {
	return newGCF0WithResolvedConfig(coeffs, DefaultConfig())
}

func NewGCF0WithConfig(coeffs BLFTCoefficients, cfg Config) *GCF {
	return newGCF0WithResolvedConfig(coeffs, cfg)
}

func NewGCF1(coeffs BLFTCoefficients, x PQStream) *GCF {
	return newGCF1WithResolvedConfig(coeffs, x, DefaultConfig())
}

func NewGCF1WithConfig(coeffs BLFTCoefficients, x PQStream, cfg Config) *GCF {
	return newGCF1WithResolvedConfig(coeffs, x, cfg)
}

func NewGCF2(coeffs BLFTCoefficients, x, y PQStream) *GCF {
	return newGCF2WithResolvedConfig(coeffs, x, y, DefaultConfig())
}

func NewGCF2WithConfig(coeffs BLFTCoefficients, x, y PQStream, cfg Config) *GCF {
	return newGCF2WithResolvedConfig(coeffs, x, y, cfg)
}

func newObservedRCFGCF(src RCFStream) *GCF {
	return newObservedRCFGCFWithResolvedConfig(src, DefaultConfig())
}

func newObservedRCFGCFWithResolvedConfig(src RCFStream, cfg Config) *GCF {
	if src == nil {
		return &GCF{
			cfg:    cfg,
			stream: newErrorRCFStream(fmt.Errorf("newObservedRCFGCFWithResolvedConfig: %w", ErrNilObservedSource)),
		}
	}

	return &GCF{
		cfg:    cfg,
		stream: src,
	}
}

func newExactTerminalGCFWithResolvedConfig(terms []RCFTerm, rng Range, cfg Config) *GCF {
	return &GCF{
		cfg: cfg,
		terminal: &exactTerminalState{
			terms: cloneRCFTerms(terms),
			rng:   cloneRange(rng),
			next:  0,
		},
	}
}

func newGCF0WithResolvedConfig(coeffs BLFTCoefficients, cfg Config) *GCF {
	return &GCF{
		coeffs: cloneBLFTCoefficients(coeffs),
		cfg:    cfg,
	}
}

func newGCF1WithResolvedConfig(coeffs BLFTCoefficients, x PQStream, cfg Config) *GCF {
	g := &GCF{
		coeffs: cloneBLFTCoefficients(coeffs),
		cfg:    cfg,
		x:      x,
	}
	if x != nil {
		g.unary = &unaryEvaluatorState{
			engine: newBLFTState(coeffs),
			x:      x,
		}
	}
	return g
}

func newGCF2WithResolvedConfig(coeffs BLFTCoefficients, x, y PQStream, cfg Config) *GCF {
	g := &GCF{
		coeffs: cloneBLFTCoefficients(coeffs),
		cfg:    cfg,
		x:      x,
		y:      y,
	}
	if x == nil || y == nil {
		return g
	}

	state := newBLFTState(coeffs)

	switch {
	case state.IndependentOfY():
		final, err := exactRationalFromUnaryEngine(collapseIndependentOfYToUnary(state), x)
		if err != nil {
			return &GCF{
				cfg:    cfg,
				stream: newErrorRCFStream(fmt.Errorf("newGCF2WithResolvedConfig: collapse independent of Y: %w", err)),
			}
		}
		return newExactTerminalGCFWithResolvedConfig(
			rcfTermsFromRational(final),
			exactRangeFromRational(final),
			cfg,
		)

	case state.IndependentOfX():
		final, err := exactRationalFromUnaryEngine(collapseIndependentOfXToUnary(state), y)
		if err != nil {
			return &GCF{
				cfg:    cfg,
				stream: newErrorRCFStream(fmt.Errorf("newGCF2WithResolvedConfig: collapse independent of X: %w", err)),
			}
		}
		return newExactTerminalGCFWithResolvedConfig(
			rcfTermsFromRational(final),
			exactRangeFromRational(final),
			cfg,
		)

	default:
		g.binary = &binaryEvaluatorState{
			engine: state,
			x:      x,
			y:      y,
		}
		return g
	}
}

func (g *GCF) NextRCF() (RCFTerm, Status, error) {
	if g == nil {
		return NewRCFTerm(nil), StatusEOF, fmt.Errorf("GCF.NextRCF: %w", ErrNilReceiver)
	}

	if g.stream != nil {
		return g.stream.NextRCF()
	}

	if g.terminal != nil {
		if g.terminal.next >= len(g.terminal.terms) {
			return NewRCFTerm(nil), StatusEOF, nil
		}
		term := g.terminal.terms[g.terminal.next]
		g.terminal.next++
		return NewRCFTerm(term.A()), StatusOK, nil
	}

	if g.unary != nil {
		return g.nextUnaryRCF()
	}

	if g.binary != nil {
		return g.nextBinaryRCF()
	}

	return NewRCFTerm(nil), StatusEOF, nil
}

func (g *GCF) Range() (Range, error) {
	if g == nil {
		return Range{}, fmt.Errorf("GCF.Range: %w", ErrNilReceiver)
	}

	if g.stream != nil {
		return g.stream.Range()
	}

	if g.terminal != nil {
		return cloneRange(g.terminal.rng), nil
	}

	if g.unary != nil {
		if isEOFPQStream(g.unary.x) {
			return exactRangeFromRational(g.unary.engine.CollapseUnaryEOF()), nil
		}
		xRange, err := g.unary.x.Range()
		if err != nil {
			return Range{}, fmt.Errorf("GCF.Range: unary child range: %w", err)
		}
		return g.unary.engine.UnaryRange(xRange)
	}

	if g.binary != nil {
		return g.binaryRange()
	}

	return exactRangeFromRational(RationalFromInt64(0)), nil
}

func (g *GCF) Config() Config {
	if g == nil {
		return DefaultConfig()
	}
	return g.cfg
}

func (g *GCF) nextUnaryRCF() (RCFTerm, Status, error) {
	for {
		if isEOFPQStream(g.unary.x) {
			collapsed := g.unary.engine.CollapseUnaryEOF()
			g.terminal = &exactTerminalState{
				terms: cloneRCFTerms(rcfTermsFromRational(collapsed)),
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
			g.unary.engine = g.unary.engine.EmitUnary(term)
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

func (g *GCF) nextBinaryRCF() (RCFTerm, Status, error) {
	for {
		if isEOFPQStream(g.binary.x) && isEOFPQStream(g.binary.y) {
			collapsed := g.binary.engine.CollapseBinaryBothEOF()
			g.terminal = &exactTerminalState{
				terms: cloneRCFTerms(rcfTermsFromRational(collapsed)),
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
			g.terminal = &exactTerminalState{
				terms: cloneRCFTerms(rcfTermsFromRational(final)),
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
			g.terminal = &exactTerminalState{
				terms: cloneRCFTerms(rcfTermsFromRational(final)),
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

func collapseUnaryXEOF(coeffs blftState) Rational {
	if coeffs.F != nil && coeffs.F.Sign() != 0 {
		return NewRational(coeffs.B, coeffs.F)
	}
	return NewRational(coeffs.D, coeffs.H)
}

func collapseIndependentOfXToUnary(coeffs blftState) blftState {
	return blftState{
		A: big.NewInt(0),
		B: cloneBigInt(coeffs.C),
		C: big.NewInt(0),
		D: cloneBigInt(coeffs.D),
		E: big.NewInt(0),
		F: cloneBigInt(coeffs.G),
		G: big.NewInt(0),
		H: cloneBigInt(coeffs.H),
	}.Normalize()
}

func collapseIndependentOfYToUnary(coeffs blftState) blftState {
	return blftState(coeffs.CollapseY()).Normalize()
}

func collapseXEOFToUnary(coeffs blftState) blftState {
	if isIndependentOfX(coeffs) {
		return collapseIndependentOfXToUnary(coeffs)
	}
	return blftState{
		A: big.NewInt(0),
		B: cloneBigInt(coeffs.A),
		C: big.NewInt(0),
		D: cloneBigInt(coeffs.B),
		E: big.NewInt(0),
		F: cloneBigInt(coeffs.E),
		G: big.NewInt(0),
		H: cloneBigInt(coeffs.F),
	}.Normalize()
}

func collapseYEOFToUnary(coeffs blftState) blftState {
	if isIndependentOfY(coeffs) {
		return collapseIndependentOfYToUnary(coeffs)
	}
	return blftState{
		A: big.NewInt(0),
		B: cloneBigInt(coeffs.A),
		C: big.NewInt(0),
		D: cloneBigInt(coeffs.C),
		E: big.NewInt(0),
		F: cloneBigInt(coeffs.E),
		G: big.NewInt(0),
		H: cloneBigInt(coeffs.G),
	}.Normalize()
}

func isIndependentOfY(coeffs blftState) bool {
	return coeffIsZero(coeffs.A) && coeffIsZero(coeffs.C) && coeffIsZero(coeffs.E) && coeffIsZero(coeffs.G)
}

func isIndependentOfX(coeffs blftState) bool {
	return coeffIsZero(coeffs.A) && coeffIsZero(coeffs.B) && coeffIsZero(coeffs.E) && coeffIsZero(coeffs.F)
}

func coeffIsZero(x *big.Int) bool {
	return x == nil || x.Sign() == 0
}

func isExactIntegerRangeWithTerm(r Range, term RCFTerm) bool {
	if !r.Inside {
		return false
	}
	if r.Lo.Value.Cmp(r.Hi.Value) != 0 {
		return false
	}
	if r.Lo.Value.Den().Cmp(big.NewInt(1)) != 0 {
		return false
	}
	return r.Lo.Value.Num().Cmp(term.A()) == 0
}

func isEOFPQStream(x PQStream) bool {
	_, ok := x.(*eofPQStream)
	return ok
}

func exactRangeFromRational(r Rational) Range {
	return Range{
		Lo: Endpoint{
			Value: r,
			Open:  false,
		},
		Hi: Endpoint{
			Value: r,
			Open:  false,
		},
		Inside: true,
	}
}

func cloneRCFTerms(terms []RCFTerm) []RCFTerm {
	if terms == nil {
		return nil
	}
	out := make([]RCFTerm, len(terms))
	for i, term := range terms {
		out[i] = NewRCFTerm(term.A())
	}
	return out
}

// core/gcf.go v14
