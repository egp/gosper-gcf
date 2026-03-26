// core/gcf_take.go v3
package core

import "math/big"

func (g *GCF) Take(n int) PQStream {
	if g == nil {
		panic("GCF.Take: nil receiver")
	}

	terms := g.takeRCFTermsUpTo(n)
	steps := finitePQStepsFromRCFTerms(terms)

	stream, status := NewFinitePQStream(steps)
	if status != StatusOK {
		panic("GCF.Take: internal invalid finite PQ steps")
	}

	return stream
}

func (g *GCF) Rational(n int) Rational {
	if g == nil {
		panic("GCF.Rational: nil receiver")
	}

	terms := g.takeRCFTermsUpTo(n)
	if len(terms) == 0 {
		return RationalFromInt64(0)
	}

	steps := finitePQStepsFromRCFTerms(terms)
	if len(steps) == 0 {
		return RationalFromInt64(0)
	}

	return steps[0].Range.Lo.Value
}

func (g *GCF) takeRCFTermsUpTo(n int) []RCFTerm {
	if g == nil {
		panic("GCF.takeRCFTermsUpTo: nil receiver")
	}
	if n <= 0 {
		return nil
	}

	out := make([]RCFTerm, 0, n)
	for len(out) < n {
		term, status := g.NextRCF()
		if status == StatusEOF {
			break
		}
		if status != StatusOK {
			panic("GCF.takeRCFTermsUpTo: invalid RCF status")
		}
		out = append(out, NewRCFTerm(term.A()))
	}

	return out
}

func finitePQStepsFromRCFTerms(terms []RCFTerm) []FinitePQStep {
	if len(terms) == 0 {
		return nil
	}

	suffixes := make([]Rational, len(terms))

	current := NewRational(terms[len(terms)-1].A(), big.NewInt(1))
	suffixes[len(terms)-1] = current

	for i := len(terms) - 2; i >= 0; i-- {
		a := terms[i].A()

		num := new(big.Int).Mul(a, current.Num())
		num.Add(num, current.Den())

		den := current.Num()

		current = NewRational(num, den)
		suffixes[i] = current
	}

	steps := make([]FinitePQStep, len(terms))
	for i, term := range terms {
		steps[i] = FinitePQStep{
			Term: PQTerm{
				P: term.A(),
				Q: big.NewInt(1),
			},
			Range: exactRangeFromRational(suffixes[i]),
		}
	}

	return steps
}

// core/gcf_take.go v3
