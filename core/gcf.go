// core/gcf.go v3
package core

import "math/big"

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
	g := &GCF{
		coeffs: cloneTransformCoefficients(coeffs),
		cfg:    cfg,
		x:      x,
	}

	if x != nil {
		value := exactRationalFromPQStream(x)
		transformed := applyUnaryXTransform(g.coeffs, value)
		return newExactTerminalGCFWithResolvedConfig(
			rcfTermsFromRational(transformed),
			exactRangeFromRational(transformed),
			cfg,
		)
	}

	return g
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

func exactRationalFromPQStream(stream PQStream) Rational {
	terms := make([]PQTerm, 0, 8)
	current := stream

	for {
		term, tail, status := current.NextPQ()

		switch status {
		case StatusOK:
			terms = append(terms, clonePQTerm(term))
			current = tail
		case StatusEOF:
			if len(terms) == 0 {
				panic("exactRationalFromPQStream: empty stream")
			}
			return exactRationalFromTerms(terms)
		default:
			panic("exactRationalFromPQStream: invalid input status")
		}
	}
}

func exactRationalFromTerms(terms []PQTerm) Rational {
	last := terms[len(terms)-1]
	value := RationalFromBigInt(last.P)

	for i := len(terms) - 2; i >= 0; i-- {
		value = generalizedStepToRational(terms[i], value)
	}

	return value
}

func generalizedStepToRational(term PQTerm, tail Rational) Rational {
	pn := cloneBigInt(term.P)
	qn := cloneBigInt(term.Q)
	tn := tail.Num()
	td := tail.Den()

	numLeft := new(big.Int).Mul(pn, tn)
	numRight := new(big.Int).Mul(qn, td)
	num := new(big.Int).Add(numLeft, numRight)

	return NewRational(num, tn)
}

func applyUnaryXTransform(coeffs TransformCoefficients, x Rational) Rational {
	xn := x.Num()
	xd := x.Den()

	numLeft := new(big.Int).Mul(cloneBigIntOrZero(coeffs.B), xn)
	numRight := new(big.Int).Mul(cloneBigIntOrZero(coeffs.D), xd)
	num := new(big.Int).Add(numLeft, numRight)

	denLeft := new(big.Int).Mul(cloneBigIntOrZero(coeffs.F), xn)
	denRight := new(big.Int).Mul(cloneBigIntOrZero(coeffs.H), xd)
	den := new(big.Int).Add(denLeft, denRight)

	return NewRational(num, den)
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

func RationalFromBigInt(n *big.Int) Rational {
	if n == nil {
		return RationalFromInt64(0)
	}
	return NewRational(n, big.NewInt(1))
}

func cloneBigIntOrZero(x *big.Int) *big.Int {
	if x == nil {
		return big.NewInt(0)
	}
	return cloneBigInt(x)
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

// core/gcf.go v3
