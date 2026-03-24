// core/rational_to_rcf.go v2
package core

import "math/big"

func rcfTermsFromRational(r Rational) []RCFTerm {
	n := r.Num()
	d := r.Den()

	if d.Sign() == 0 {
		panic("rcfTermsFromRational: zero denominator")
	}

	terms := make([]RCFTerm, 0, 8)

	for {
		q, rem := floorQuoRem(n, d)
		terms = append(terms, NewRCFTerm(q))

		if rem.Sign() == 0 {
			return terms
		}

		n, d = d, rem
	}
}

func floorQuoRem(n, d *big.Int) (*big.Int, *big.Int) {
	if d.Sign() <= 0 {
		panic("floorQuoRem: denominator must be positive")
	}

	q := new(big.Int).Quo(n, d)
	rem := new(big.Int).Rem(n, d)

	if rem.Sign() < 0 {
		q.Sub(q, big.NewInt(1))
		rem.Add(rem, d)
	}

	return q, rem
}

// core/rational_to_rcf.go v2
