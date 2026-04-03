// core/gcf_types.go v6
package core

import "fmt"

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
		return &GCF{
			cfg: cfg,
			stream: newErrorRCFStream(
				fmt.Errorf("newObservedRCFGCFWithResolvedConfig: %w", ErrNilObservedSource),
			),
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
	if x != nil && y != nil {
		g.binary = &binaryEvaluatorState{
			engine: newBLFTState(coeffs),
			x:      x,
			y:      y,
		}
	}
	return g
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

func (g *GCF) CurrentInterval() (Interval, error) {
	return g.Range()
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

// core/gcf_types.go v6
