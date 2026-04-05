// core/rectifier.go V5

package core

import "math/big"

// Rectifier implements the Version 7 "Rectification Tier".
// It is a Unary LFT [a, b; c, d] that certifies internal noisy terms.
type Rectifier struct {
	a, b, c, d *big.Int
}

func NewRectifier() *Rectifier {
	return &Rectifier{
		a: big.NewInt(1), b: big.NewInt(0),
		c: big.NewInt(0), d: big.NewInt(1),
	}
}

// Absorb takes a generalized term (p, q) and folds it into the state.
// This is the "State-driven buffer" from the spec.
func (r *Rectifier) Absorb(p, q *big.Int) {
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

// CanEmit checks if floor(z(1)) == floor(z(inf))
func (r *Rectifier) CanEmit() (bool, *big.Int) {
	if r.c.Sign() == 0 && r.d.Sign() == 0 {
		return false, nil
	}

	// Case 1: z(inf) = a/c
	if r.c.Sign() == 0 {
		return false, nil
	} // Infinite
	fInf := new(big.Int).Div(r.a, r.c)

	// Case 2: z(1) = (a+b)/(c+d)
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

// Produce updates state after an RCF term t is emitted: z = t + 1/z' => z' = 1/(z-t)
func (r *Rectifier) Produce(t *big.Int) {
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
	if g.Cmp(big.NewInt(1)) > 0 {
		r.a.Div(r.a, g)
		r.b.Div(r.b, g)
		r.c.Div(r.c, g)
		r.d.Div(r.d, g)
	}
}

// core/rectifier.go V
