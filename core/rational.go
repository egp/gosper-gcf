// core/rational.go v3
package core

import (
	"errors"
	"math/big"
)

var (
	ErrNilRationalDenominator    = errors.New("rational denominator is nil")
	ErrZeroRationalDenominator   = errors.New("rational denominator is zero")
	ErrNonPositiveQuoDenominator = errors.New("floor quotient denominator must be positive")
	ErrRCFTermDoesNotFitInt64    = errors.New("regular-CF term does not fit int64")
)

type Rational struct {
	num *big.Int
	den *big.Int
}

func RationalFromInt64(n int64) Rational {
	return Rational{
		num: big.NewInt(n),
		den: big.NewInt(1),
	}
}

// NewRational is kept for compatibility during the phase-1 cleanup.
// New core runtime paths should prefer NewRationalChecked so they can
// propagate an explicit error instead of silently normalizing invalid input.
func NewRational(num, den *big.Int) Rational {
	r, err := NewRationalChecked(num, den)
	if err != nil {
		return RationalFromInt64(0)
	}
	return r
}

func NewRationalChecked(num, den *big.Int) (Rational, error) {
	if den == nil {
		return Rational{}, ErrNilRationalDenominator
	}
	if den.Sign() == 0 {
		return Rational{}, ErrZeroRationalDenominator
	}

	n := new(big.Int)
	if num != nil {
		n.Set(num)
	}

	d := new(big.Int).Set(den)

	if d.Sign() < 0 {
		n.Neg(n)
		d.Neg(d)
	}

	if n.Sign() == 0 {
		return Rational{
			num: big.NewInt(0),
			den: big.NewInt(1),
		}, nil
	}

	g := new(big.Int).GCD(nil, nil, new(big.Int).Abs(new(big.Int).Set(n)), d)
	n.Quo(n, g)
	d.Quo(d, g)

	return Rational{
		num: n,
		den: d,
	}, nil
}

func (r Rational) Num() *big.Int {
	if r.num == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Set(r.num)
}

func (r Rational) Den() *big.Int {
	if r.den == nil || r.den.Sign() == 0 {
		return big.NewInt(1)
	}
	return new(big.Int).Set(r.den)
}

func (r Rational) Cmp(other Rational) int {
	left := new(big.Int).Mul(r.Num(), other.Den())
	right := new(big.Int).Mul(other.Num(), r.Den())
	return left.Cmp(right)
}

// core/rational.go v3
