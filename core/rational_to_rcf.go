// core/rational_to_rcf.go v5
package core

import "math/big"

func rcfTermsFromRational(r Rational) []RCFTerm {
	terms, err := rcfTermsFromRationalChecked(r)
	if err != nil {
		return nil
	}
	return terms
}

func rcfTermsFromRationalChecked(r Rational) ([]RCFTerm, error) {
	n := r.Num()
	d := r.Den()

	if d.Sign() == 0 {
		return nil, ErrZeroRationalDenominator
	}

	terms := make([]RCFTerm, 0, 8)

	for {
		q, rem, err := floorQuoRemChecked(n, d)
		if err != nil {
			return nil, err
		}
		terms = append(terms, NewRCFTerm(q))

		if rem.Sign() == 0 {
			return terms, nil
		}

		n, d = d, rem
	}
}

func floorQuoRemChecked(n, d *big.Int) (*big.Int, *big.Int, error) {
	if d == nil || d.Sign() <= 0 {
		return nil, nil, ErrNonPositiveQuoDenominator
	}

	q := new(big.Int).Quo(n, d)
	rem := new(big.Int).Rem(n, d)

	if rem.Sign() < 0 {
		q.Sub(q, big.NewInt(1))
		rem.Add(rem, d)
	}

	return q, rem, nil
}

// core/rational_to_rcf.go v5
