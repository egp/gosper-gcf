package core

import "math/big"

type Rational struct {
	num *big.Int
	den *big.Int
}

func FromInt64(n int64) Rational {
	return Rational{
		num: big.NewInt(n),
		den: big.NewInt(1),
	}
}

func NewRational(num, den *big.Int) Rational {
	if den == nil {
		panic("Rational denominator is nil")
	}
	if den.Sign() == 0 {
		panic("Rational denominator is zero")
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
		}
	}

	g := new(big.Int).GCD(nil, nil, new(big.Int).Abs(new(big.Int).Set(n)), d)
	n.Quo(n, g)
	d.Quo(d, g)

	return Rational{
		num: n,
		den: d,
	}
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
