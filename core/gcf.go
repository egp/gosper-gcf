// core/gcf.go v10
package core

import "math/big"

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
	coeffs BLFTCoefficients
	cfg    Config
	x      PQStream
	y      PQStream

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
		panic("newObservedRCFGCFWithResolvedConfig: nil source")
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
		final := exactRationalFromUnaryEngine(state, x)
		return newExactTerminalGCFWithResolvedConfig(
			rcfTermsFromRational(final),
			exactRangeFromRational(final),
			cfg,
		)

	case state.IndependentOfX():
		final := exactRationalFromUnaryEngine(state.CollapseBinaryXEOF(), y)
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

func (g *GCF) NextRCF() (RCFTerm, Status) {
	if g == nil {
		panic("GCF receiver is nil")
	}

	if g.stream != nil {
		return g.stream.NextRCF()
	}

	if g.terminal != nil {
		if g.terminal.next >= len(g.terminal.terms) {
			return NewRCFTerm(nil), StatusEOF
		}
		term := g.terminal.terms[g.terminal.next]
		g.terminal.next++
		return NewRCFTerm(term.A()), StatusOK
	}

	if g.unary != nil {
		return g.nextUnaryRCF()
	}

	if g.binary != nil {
		return g.nextBinaryRCF()
	}

	return NewRCFTerm(nil), StatusEOF
}

func (g *GCF) Range() Range {
	if g == nil {
		panic("GCF receiver is nil")
	}

	if g.stream != nil {
		return g.stream.Range()
	}

	if g.terminal != nil {
		return cloneRange(g.terminal.rng)
	}

	if g.unary != nil {
		if isEOFPQStream(g.unary.x) {
			return exactRangeFromRational(g.unary.engine.CollapseUnaryEOF())
		}
		return g.unary.engine.UnaryRange(g.unary.x.Range())
	}

	if g.binary != nil {
		return g.binaryRange()
	}

	return exactRangeFromRational(RationalFromInt64(0))
}

func (g *GCF) Config() Config {
	if g == nil {
		return DefaultConfig()
	}
	return g.cfg
}

func (g *GCF) nextUnaryRCF() (RCFTerm, Status) {
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

		currentRange := g.unary.engine.UnaryRange(g.unary.x.Range())
		if term, ok := g.unary.engine.CanEmitRCFTerm(currentRange); ok {
			if isExactIntegerRangeWithTerm(currentRange, term) {
				g.terminal = &exactTerminalState{
					terms: nil,
					rng:   exactRangeFromRational(RationalFromInt64(0)),
					next:  0,
				}
				g.unary = nil
				return term, StatusOK
			}

			g.unary.engine = g.unary.engine.EmitUnary(term)
			return term, StatusOK
		}

		term, tail, status := g.unary.x.NextPQ()
		switch status {
		case StatusOK:
			g.unary.engine = g.unary.engine.IngestUnaryX(term)
			g.unary.x = tail
		case StatusEOF:
			g.unary.x = tail
		default:
			panic("nextUnaryRCF: invalid input status")
		}
	}
}

func (g *GCF) nextBinaryRCF() (RCFTerm, Status) {
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
			final := exactRationalFromUnaryEngine(g.binary.engine.CollapseBinaryXEOF(), g.binary.y)
			g.terminal = &exactTerminalState{
				terms: cloneRCFTerms(rcfTermsFromRational(final)),
				rng:   exactRangeFromRational(final),
				next:  0,
			}
			g.binary = nil
			return g.NextRCF()
		}

		if isEOFPQStream(g.binary.y) {
			final := exactRationalFromUnaryEngine(g.binary.engine.CollapseBinaryYEOF(), g.binary.x)
			g.terminal = &exactTerminalState{
				terms: cloneRCFTerms(rcfTermsFromRational(final)),
				rng:   exactRangeFromRational(final),
				next:  0,
			}
			g.binary = nil
			return g.NextRCF()
		}

		currentRange := g.binary.engine.BinaryRange(g.binary.x.Range(), g.binary.y.Range())
		if term, ok := g.binary.engine.CanEmitRCFTerm(currentRange); ok {
			if isExactIntegerRangeWithTerm(currentRange, term) {
				g.terminal = &exactTerminalState{
					terms: nil,
					rng:   exactRangeFromRational(RationalFromInt64(0)),
					next:  0,
				}
				g.binary = nil
				return term, StatusOK
			}

			g.binary.engine = g.binary.engine.EmitBinary(term)
			return term, StatusOK
		}

		if chooseIngestX(g.binary.x.Range(), g.binary.y.Range()) {
			term, tail, status := g.binary.x.NextPQ()
			switch status {
			case StatusOK:
				g.binary.engine = g.binary.engine.IngestBinaryX(term)
				g.binary.x = tail
			case StatusEOF:
				g.binary.x = tail
			default:
				panic("nextBinaryRCF: invalid left-input status")
			}
		} else {
			term, tail, status := g.binary.y.NextPQ()
			switch status {
			case StatusOK:
				g.binary.engine = g.binary.engine.IngestBinaryY(term)
				g.binary.y = tail
			case StatusEOF:
				g.binary.y = tail
			default:
				panic("nextBinaryRCF: invalid right-input status")
			}
		}
	}
}

func (g *GCF) binaryRange() Range {
	if g.binary == nil {
		panic("binaryRange called without binary evaluator state")
	}

	switch {
	case isEOFPQStream(g.binary.x) && isEOFPQStream(g.binary.y):
		return exactRangeFromRational(g.binary.engine.CollapseBinaryBothEOF())

	case isEOFPQStream(g.binary.x):
		return g.binary.engine.CollapseBinaryXEOF().UnaryRange(g.binary.y.Range())

	case isEOFPQStream(g.binary.y):
		return g.binary.engine.CollapseBinaryYEOF().UnaryRange(g.binary.x.Range())

	default:
		return g.binary.engine.BinaryRange(g.binary.x.Range(), g.binary.y.Range())
	}
}

func exactRationalFromUnaryEngine(engine unaryEngine, stream PQStream) Rational {
	current := engine
	input := stream

	for {
		term, tail, status := input.NextPQ()
		switch status {
		case StatusOK:
			current = current.IngestUnaryX(term)
			input = tail

		case StatusEOF:
			return current.CollapseUnaryEOF()

		default:
			panic("exactRationalFromUnaryEngine: invalid input status")
		}
	}
}

func collapseUnaryXEOF(coeffs blftState) Rational {
	if coeffs.F != nil && coeffs.F.Sign() != 0 {
		return NewRational(coeffs.B, coeffs.F)
	}
	return NewRational(coeffs.D, coeffs.H)
}

func collapseXEOFToUnary(coeffs blftState) blftState {
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

func collapseYEOFToUnary(coeffs blftState) blftState {
	return blftState(coeffs.CollapseY()).Normalize()
}

func isIndependentOfY(coeffs blftState) bool {
	return coeffIsZero(coeffs.A) &&
		coeffIsZero(coeffs.C) &&
		coeffIsZero(coeffs.E) &&
		coeffIsZero(coeffs.G)
}

func isIndependentOfX(coeffs blftState) bool {
	return coeffIsZero(coeffs.A) &&
		coeffIsZero(coeffs.B) &&
		coeffIsZero(coeffs.E) &&
		coeffIsZero(coeffs.F)
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

// core/gcf.go v10
