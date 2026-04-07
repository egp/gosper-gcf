// core/rectifier.go V3
// Version: 3

package core

import (
	"math/big"
)

type Rectifier struct {
	a, b, c, d *big.Int
}

func NewRectifier(a, b, c, d *big.Int) *Rectifier {
	return &Rectifier{
		a: new(big.Int).Set(a),
		b: new(big.Int).Set(b),
		c: new(big.Int).Set(c),
		d: new(big.Int).Set(d),
	}
}

// Absorb folds a speculative term into the 4-tuple.
func (r *Rectifier) Absorb(pq PQTerm) {
	p, q := pq.P, pq.Q

	newA := new(big.Int).Mul(r.a, p)
	newA.Add(newA, new(big.Int).Mul(r.b, q))

	newC := new(big.Int).Mul(r.c, p)
	newC.Add(newC, new(big.Int).Mul(r.d, q))

	r.b.Set(r.a)
	r.a.Set(newA)
	r.d.Set(r.c)
	r.c.Set(newC)

	r.Normalize()
}

func (r *Rectifier) CanEmit() (bool, *big.Int) {
	if r.c.Sign() == 0 || new(big.Int).Add(r.c, r.d).Sign() == 0 {
		return false, nil
	}

	// z(inf) = a/c
	fInf := new(big.Int).Div(r.a, r.c)

	// z(1) = (a+b)/(c+d)
	num1 := new(big.Int).Add(r.a, r.b)
	den1 := new(big.Int).Add(r.c, r.d)
	fOne := new(big.Int).Div(num1, den1)

	if fInf.Cmp(fOne) == 0 {
		return true, fInf
	}
	return false, nil
}

// Emit performs the state transition after a term is successfully proven.
func (r *Rectifier) Emit(t *big.Int) {
	tc := new(big.Int).Mul(t, r.c)
	td := new(big.Int).Mul(t, r.d)

	newC := new(big.Int).Sub(r.a, tc)
	newD := new(big.Int).Sub(r.b, td)

	r.a.Set(r.c)
	r.b.Set(r.d)
	r.c.Set(newC)
	r.d.Set(newD)
	r.Normalize()
}

func (r *Rectifier) Normalize() {
	g := new(big.Int).GCD(nil, nil, r.a, r.b)
	g.GCD(nil, nil, g, r.c)
	g.GCD(nil, nil, g, r.d)

	if g.Cmp(big.NewInt(1)) > 0 {
		r.a.Div(r.a, g)
		r.b.Div(r.b, g)
		r.c.Div(r.c, g)
		r.d.Div(r.d, g)
	}
}

// core/rectifier.go V3
