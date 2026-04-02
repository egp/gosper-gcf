// core/rational_source.go v2
package core

import (
	"fmt"
	"math/big"
)

func PQStreamFromRational(r Rational) PQStream {
	stream, err := PQStreamFromRationalChecked(r)
	if err != nil {
		return newErrorPQStream(err)
	}
	return stream
}

func PQStreamFromRationalChecked(r Rational) (PQStream, error) {
	terms, err := finiteRCFTermsFromRationalChecked(r)
	if err != nil {
		return nil, fmt.Errorf("PQStreamFromRationalChecked: finite terms: %w", err)
	}

	suffixes, err := suffixRationalsFromRCFTermsChecked(terms)
	if err != nil {
		return nil, fmt.Errorf("PQStreamFromRationalChecked: suffix rationals: %w", err)
	}

	steps := make([]FinitePQStep, len(terms))
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
		return nil, fmt.Errorf("PQStreamFromRationalChecked: NewFinitePQStream status=%v", status)
	}
	return stream, nil
}

func finiteRCFTermsFromRational(r Rational) []int64 {
	terms, err := finiteRCFTermsFromRationalChecked(r)
	if err != nil {
		return nil
	}
	return terms
}

func finiteRCFTermsFromRationalChecked(r Rational) ([]int64, error) {
	n := r.Num()
	d := r.Den()

	if d.Sign() == 0 {
		return nil, ErrZeroRationalDenominator
	}

	out := make([]int64, 0, 8)

	for {
		q, rem, err := floorQuoRemChecked(n, d)
		if err != nil {
			return nil, err
		}
		if !q.IsInt64() {
			return nil, ErrRCFTermDoesNotFitInt64
		}
		out = append(out, q.Int64())

		if rem.Sign() == 0 {
			return out, nil
		}

		n, d = d, rem
	}
}

func suffixRationalsFromRCFTerms(terms []int64) []Rational {
	out, err := suffixRationalsFromRCFTermsChecked(terms)
	if err != nil {
		return nil
	}
	return out
}

func suffixRationalsFromRCFTermsChecked(terms []int64) ([]Rational, error) {
	n := len(terms)
	if n == 0 {
		return []Rational{}, nil
	}

	out := make([]Rational, n)

	current, err := NewRationalChecked(big.NewInt(terms[n-1]), big.NewInt(1))
	if err != nil {
		return nil, fmt.Errorf("suffixRationalsFromRCFTermsChecked: seed term: %w", err)
	}
	out[n-1] = current

	for i := n - 2; i >= 0; i-- {
		a := big.NewInt(terms[i])

		num := new(big.Int).Mul(a, current.Num())
		num.Add(num, current.Den())

		den := current.Num()
		current, err = NewRationalChecked(num, den)
		if err != nil {
			return nil, fmt.Errorf("suffixRationalsFromRCFTermsChecked: index %d: %w", i, err)
		}
		out[i] = current
	}

	return out, nil
}

// core/rational_source.go v2
