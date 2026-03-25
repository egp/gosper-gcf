// core/gcf.go v8
package core

import "math/big"

type exactTerminalState struct {
	terms []RCFTerm
	rng   Range
	next  int
}

type unaryEvaluatorState struct {
	coeffs blftState
	x      PQStream
}

type binaryEvaluatorState struct {
	coeffs blftState
	x      PQStream
	y      PQStream
}

type GCF struct {
	coeffs TransformCoefficients
	cfg    Config
	x      PQStream
	y      PQStream

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
		g.unary = &unaryEvaluatorState{
			coeffs: blftState(cloneTransformCoefficients(coeffs)),
			x:      x,
		}
	}

	return g
}

func newGCF2WithResolvedConfig(coeffs TransformCoefficients, x, y PQStream, cfg Config) *GCF {
	g := &GCF{
		coeffs: cloneTransformCoefficients(coeffs),
		cfg:    cfg,
		x:      x,
		y:      y,
	}

	if x == nil || y == nil {
		return g
	}

	state := blftState(cloneTransformCoefficients(coeffs))

	switch {
	case isIndependentOfY(state):
		final := applyUnaryXTransform(TransformCoefficients(state), exactRationalFromPQStream(x))
		return newExactTerminalGCFWithResolvedConfig(
			rcfTermsFromRational(final),
			exactRangeFromRational(final),
			cfg,
		)

	case isIndependentOfX(state):
		unaryCoeffs := collapseXEOFToUnary(state)
		final := applyUnaryXTransform(TransformCoefficients(unaryCoeffs), exactRationalFromPQStream(y))
		return newExactTerminalGCFWithResolvedConfig(
			rcfTermsFromRational(final),
			exactRangeFromRational(final),
			cfg,
		)

	default:
		g.binary = &binaryEvaluatorState{
			coeffs: state,
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

	if g.terminal != nil {
		return cloneRange(g.terminal.rng)
	}

	if g.unary != nil {
		if isEOFPQStream(g.unary.x) {
			collapsed := collapseUnaryXEOF(g.unary.coeffs)
			return exactRangeFromRational(collapsed)
		}
		return g.unaryRange()
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
			collapsed := collapseUnaryXEOF(g.unary.coeffs)
			g.terminal = &exactTerminalState{
				terms: cloneRCFTerms(rcfTermsFromRational(collapsed)),
				rng:   exactRangeFromRational(collapsed),
				next:  0,
			}
			g.unary = nil
			return g.NextRCF()
		}

		currentRange := g.unaryRange()

		if term, ok := g.unary.coeffs.CanEmitRCFTerm(currentRange); ok {
			if isExactIntegerRangeWithTerm(currentRange, term) {
				g.terminal = &exactTerminalState{
					terms: nil,
					rng:   exactRangeFromRational(RationalFromInt64(0)),
					next:  0,
				}
				g.unary = nil
				return term, StatusOK
			}

			g.unary.coeffs = g.unary.coeffs.Emit(term)
			return term, StatusOK
		}

		term, tail, status := g.unary.x.NextPQ()

		switch status {
		case StatusOK:
			g.unary.coeffs = g.unary.coeffs.IngestX(term)
			g.unary.x = tail

		case StatusEOF:
			g.unary.x = tail

		default:
			panic("nextUnaryRCF: invalid input status")
		}
	}
}

func (g *GCF) unaryRange() Range {
	if g.unary == nil {
		panic("unaryRange called without unary evaluator state")
	}

	return g.unary.coeffs.CornerRange(
		g.unary.x.Range(),
		exactRangeFromRational(RationalFromInt64(0)),
	)
}

func (g *GCF) nextBinaryRCF() (RCFTerm, Status) {
	for {
		if isEOFPQStream(g.binary.x) && isEOFPQStream(g.binary.y) {
			collapsed := g.binary.coeffs.CollapseToRational()
			g.terminal = &exactTerminalState{
				terms: cloneRCFTerms(rcfTermsFromRational(collapsed)),
				rng:   exactRangeFromRational(collapsed),
				next:  0,
			}
			g.binary = nil
			return g.NextRCF()
		}

		if isEOFPQStream(g.binary.x) {
			unaryCoeffs := collapseXEOFToUnary(g.binary.coeffs)
			final := applyUnaryXTransform(TransformCoefficients(unaryCoeffs), exactRationalFromPQStream(g.binary.y))
			g.terminal = &exactTerminalState{
				terms: cloneRCFTerms(rcfTermsFromRational(final)),
				rng:   exactRangeFromRational(final),
				next:  0,
			}
			g.binary = nil
			return g.NextRCF()
		}

		if isEOFPQStream(g.binary.y) {
			unaryCoeffs := collapseYEOFToUnary(g.binary.coeffs)
			final := applyUnaryXTransform(TransformCoefficients(unaryCoeffs), exactRationalFromPQStream(g.binary.x))
			g.terminal = &exactTerminalState{
				terms: cloneRCFTerms(rcfTermsFromRational(final)),
				rng:   exactRangeFromRational(final),
				next:  0,
			}
			g.binary = nil
			return g.NextRCF()
		}

		currentRange := g.binaryRange()
		if term, ok := g.binary.coeffs.CanEmitRCFTerm(currentRange); ok {
			if isExactIntegerRangeWithTerm(currentRange, term) {
				g.terminal = &exactTerminalState{
					terms: nil,
					rng:   exactRangeFromRational(RationalFromInt64(0)),
					next:  0,
				}
				g.binary = nil
				return term, StatusOK
			}

			g.binary.coeffs = g.binary.coeffs.Emit(term)
			return term, StatusOK
		}

		if chooseIngestX(g.binary.x.Range(), g.binary.y.Range()) {
			term, tail, status := g.binary.x.NextPQ()
			switch status {
			case StatusOK:
				g.binary.coeffs = g.binary.coeffs.IngestX(term)
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
				g.binary.coeffs = g.binary.coeffs.IngestY(term)
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
		return exactRangeFromRational(g.binary.coeffs.CollapseToRational())

	case isEOFPQStream(g.binary.x):
		unaryCoeffs := collapseXEOFToUnary(g.binary.coeffs)
		return unaryCoeffs.CornerRange(
			g.binary.y.Range(),
			exactRangeFromRational(RationalFromInt64(0)),
		)

	case isEOFPQStream(g.binary.y):
		unaryCoeffs := collapseYEOFToUnary(g.binary.coeffs)
		return unaryCoeffs.CornerRange(
			g.binary.x.Range(),
			exactRangeFromRational(RationalFromInt64(0)),
		)

	default:
		return g.binary.coeffs.CornerRange(g.binary.x.Range(), g.binary.y.Range())
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
	value := NewRational(last.P, big.NewInt(1))

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

	numLeft := new(big.Int).Mul(cloneBigIntOrZeroLocal(coeffs.B), xn)
	numRight := new(big.Int).Mul(cloneBigIntOrZeroLocal(coeffs.D), xd)
	num := new(big.Int).Add(numLeft, numRight)

	denLeft := new(big.Int).Mul(cloneBigIntOrZeroLocal(coeffs.F), xn)
	denRight := new(big.Int).Mul(cloneBigIntOrZeroLocal(coeffs.H), xd)
	den := new(big.Int).Add(denLeft, denRight)

	return NewRational(num, den)
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

func cloneBigIntOrZeroLocal(x *big.Int) *big.Int {
	if x == nil {
		return big.NewInt(0)
	}
	return cloneBigInt(x)
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

// core/gcf.go v8
