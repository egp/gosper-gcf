// core/rational_source.go v1
package core

import "math/big"

func PQStreamFromRational(r Rational) PQStream {
	terms := finiteRCFTermsFromRational(r)
	steps := make([]FinitePQStep, len(terms))
	suffixes := suffixRationalsFromRCFTerms(terms)

	for i, term := range terms {
		steps[i] = FinitePQStep{
			Term: PQTerm{
				P: big.NewInt(term),
				Q: big.NewInt(1),
			},
			Range: exactRangeFromRational(suffixes[i]),
		}
	}

	stream, status := NewFinitePQStream(steps)
	if status != StatusOK {
		panic("PQStreamFromRational: failed to construct finite PQ stream")
	}
	return stream
}

func finiteRCFTermsFromRational(r Rational) []int64 {
	n := r.Num()
	d := r.Den()

	if d.Sign() == 0 {
		panic("finiteRCFTermsFromRational: zero denominator")
	}

	out := make([]int64, 0, 8)

	for {
		q, rem := floorQuoRem(n, d)
		if !q.IsInt64() {
			panic("finiteRCFTermsFromRational: term does not fit int64")
		}
		out = append(out, q.Int64())

		if rem.Sign() == 0 {
			return out
		}

		n, d = d, rem
	}
}

func suffixRationalsFromRCFTerms(terms []int64) []Rational {
	n := len(terms)
	out := make([]Rational, n)

	current := NewRational(big.NewInt(terms[n-1]), big.NewInt(1))
	out[n-1] = current

	for i := n - 2; i >= 0; i-- {
		a := big.NewInt(terms[i])

		num := new(big.Int).Mul(a, current.Num())
		num.Add(num, current.Den())

		den := current.Num()
		current = NewRational(num, den)
		out[i] = current
	}

	return out
}

// core/rational_source.go v1
