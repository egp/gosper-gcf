// core/rectifier.go V2

package core

import (
	"math/big"
)

// Rectifier implements the Version 7 "Rectification Tier".
// It is a Unary LFT [a, b; c, d] that certifies internal noisy terms.
type Rectifier struct {
	a, b, c, d *big.Int
}

// NewRectifier matches your existing signature: NewRectifier(a, b, c, d)
func NewRectifier(a, b, c, d *big.Int) *Rectifier {
	return &Rectifier{
		a: new(big.Int).Set(a),
		b: new(big.Int).Set(b),
		c: new(big.Int).Set(c),
		d: new(big.Int).Set(d),
	}
}

// Absorb matches your existing signature: Absorb(PQTerm)
func (r *Rectifier) Absorb(pq PQTerm) {
	p, q := pq.P, pq.Q

	// a' = a*p + b*q, b' = a
	// c' = c*p + d*q, d' = c
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

// CanEmit replaces the "Emit" check. It returns (true, term) if the floor is stable.
func (r *Rectifier) CanEmit() (bool, *big.Int) {
	// If denominators are zero, we are at infinity or undefined
	if r.c.Sign() == 0 && r.d.Sign() == 0 {
		return false, nil
	}

	// Bound 1: z(inf) = a/c
	if r.c.Sign() == 0 {
		return false, nil
	}
	fInf := new(big.Int).Div(r.a, r.c)

	// Bound 2: z(1) = (a+b)/(c+d)
	num1 := new(big.Int).Add(r.a, r.b)
	den1 := new(big.Int).Add(r.c, r.d)
	if den1.Sign() == 0 {
		return false, nil
	}
	fOne := new(big.Int).Div(num1, den1)

	if fInf.Cmp(fOne) == 0 {
		return true, fInf
	}
	return false, nil
}

// Emit is a convenience method that combines the state transition logic.
// In your loops, if CanEmit returns true, you call this to update the matrix.
func (r *Rectifier) Emit(t *big.Int) {
	// Production: z = t + 1/z'  => z' = 1/(z-t)
	// a'' = c, b'' = d
	// c'' = a - tc, d'' = b - td
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
	g := gcd4(r.a, r.b, r.c, r.d)
	if g.Cmp(big.NewInt(1)) > 1 {
		r.a.Div(r.a, g)
		r.b.Div(r.b, g)
		r.c.Div(r.c, g)
		r.d.Div(r.d, g)
	}
}

// gcd4 is a helper for normalization
func gcd4(a, b, c, d *big.Int) *big.Int {
	res := new(big.Int).GCD(nil, nil, a, b)
	res.GCD(nil, nil, res, c)
	res.GCD(nil, nil, res, d)
	return res
}

// core/rectifier.go V2
