// core/gcf.go v2
package core

type exactTerminalState struct {
	terms []RCFTerm
	rng   Range
	next  int
}

type GCF struct {
	coeffs TransformCoefficients
	cfg    Config
	x      PQStream
	y      PQStream

	terminal *exactTerminalState
}

func NewExactTerminalGCF(terms []RCFTerm, rng Range) *GCF {
	return newExactTerminalGCFWithResolvedConfig(terms, rng, DefaultConfig())
}

func NewExactTerminalGCFWithConfig(terms []RCFTerm, rng Range, cfg Config) *GCF {
	return newExactTerminalGCFWithResolvedConfig(terms, rng, cfg)
}

func NewGCF0(coeffs TransformCoefficients) *GCF {
	return newGCF0WithResolvedConfig(coeffs, DefaultConfig())
}

func NewGCF0WithConfig(coeffs TransformCoefficients, cfg Config) *GCF {
	return newGCF0WithResolvedConfig(coeffs, cfg)
}

func NewGCF1(coeffs TransformCoefficients, x PQStream) *GCF {
	return newGCF1WithResolvedConfig(coeffs, x, DefaultConfig())
}

func NewGCF1WithConfig(coeffs TransformCoefficients, x PQStream, cfg Config) *GCF {
	return newGCF1WithResolvedConfig(coeffs, x, cfg)
}

func NewGCF2(coeffs TransformCoefficients, x, y PQStream) *GCF {
	return newGCF2WithResolvedConfig(coeffs, x, y, DefaultConfig())
}

func NewGCF2WithConfig(coeffs TransformCoefficients, x, y PQStream, cfg Config) *GCF {
	return newGCF2WithResolvedConfig(coeffs, x, y, cfg)
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

func newGCF0WithResolvedConfig(coeffs TransformCoefficients, cfg Config) *GCF {
	return &GCF{
		coeffs: cloneTransformCoefficients(coeffs),
		cfg:    cfg,
	}
}

func newGCF1WithResolvedConfig(coeffs TransformCoefficients, x PQStream, cfg Config) *GCF {
	return &GCF{
		coeffs: cloneTransformCoefficients(coeffs),
		cfg:    cfg,
		x:      x,
	}
}

func newGCF2WithResolvedConfig(coeffs TransformCoefficients, x, y PQStream, cfg Config) *GCF {
	return &GCF{
		coeffs: cloneTransformCoefficients(coeffs),
		cfg:    cfg,
		x:      x,
		y:      y,
	}
}

func (g *GCF) NextRCF() (RCFTerm, Status) {
	if g == nil {
		panic("GCF receiver is nil")
	}

	if g.terminal != nil {
		if g.terminal.next >= len(g.terminal.terms) {
			return NewRCFTerm(nil), StatusEOF
		}

		term := g.terminal.terms[g.terminal.next]
		g.terminal.next++
		return NewRCFTerm(term.A()), StatusOK
	}

	return NewRCFTerm(nil), StatusEOF
}

func (g *GCF) Range() Range {
	if g == nil {
		panic("GCF receiver is nil")
	}

	if g.terminal != nil {
		return cloneRange(g.terminal.rng)
	}

	return Range{
		Lo: Endpoint{
			Value: RationalFromInt64(0),
			Open:  false,
		},
		Hi: Endpoint{
			Value: RationalFromInt64(0),
			Open:  false,
		},
		Inside: true,
	}
}

func (g *GCF) Config() Config {
	if g == nil {
		return DefaultConfig()
	}
	return g.cfg
}

func cloneTransformCoefficients(tc TransformCoefficients) TransformCoefficients {
	return TransformCoefficients{
		A: cloneBigInt(tc.A),
		B: cloneBigInt(tc.B),
		C: cloneBigInt(tc.C),
		D: cloneBigInt(tc.D),
		E: cloneBigInt(tc.E),
		F: cloneBigInt(tc.F),
		G: cloneBigInt(tc.G),
		H: cloneBigInt(tc.H),
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

// core/gcf.go v2
