// core/rectifier.go V4
package core

import "math/big"

type Rectifier struct {
	a, b, c, d *big.Int
}

func NewRectifier(a, b, c, d *big.Int) *Rectifier {
	return &Rectifier{
		a: cloneBigIntOrZero(a),
		b: cloneBigIntOrZero(b),
		c: cloneBigIntOrZero(c),
		d: cloneBigIntOrZero(d),
	}
}

// Absorb — exact rule from newSpec.md §4 + GCD normalization after each update
func (r *Rectifier) Absorb(term PQTerm) *Rectifier {
	if r == nil {
		return nil
	}
	p := term.P
	q := term.Q

	newA := new(big.Int).Add(new(big.Int).Mul(r.a, p), new(big.Int).Mul(r.c, q))
	newB := cloneBigIntOrZero(r.a)
	newC := new(big.Int).Add(new(big.Int).Mul(r.c, p), new(big.Int).Mul(r.d, q))
	newD := cloneBigIntOrZero(r.c)

	// GCD normalization (chained exactly like blft_normalize.go)
	g := new(big.Int).GCD(nil, nil, newA, newB)
	g = new(big.Int).GCD(g, nil, g, newC)
	g = new(big.Int).GCD(g, nil, g, newD)
	if g.Sign() > 0 {
		newA.Div(newA, g)
		newB.Div(newB, g)
		newC.Div(newC, g)
		newD.Div(newD, g)
	}

	r.a = newA
	r.b = newB
	r.c = newC
	r.d = newD
	return r
}

// CanEmit — exact interval check from newSpec.md §4
func (r *Rectifier) CanEmit() (RCFTerm, bool) {
	if r == nil || r.c.Sign() == 0 {
		return NewRCFTerm(nil), false
	}

	zInf := NewRational(cloneBigIntOrZero(r.a), cloneBigIntOrZero(r.c))
	z1Num := new(big.Int).Add(cloneBigIntOrZero(r.a), cloneBigIntOrZero(r.b))
	z1Den := new(big.Int).Add(cloneBigIntOrZero(r.c), cloneBigIntOrZero(r.d))
	z1 := NewRational(z1Num, z1Den)

	rng := Range{
		Lo:     Endpoint{Value: zInf, Open: false},
		Hi:     Endpoint{Value: z1, Open: false},
		Inside: true,
	}

	return canEmitRCFTermFromRange(rng)
}

// Emit — production substitution (z = t + 1/z') after emit
func (r *Rectifier) Emit(term RCFTerm) *Rectifier {
	if r == nil {
		return nil
	}
	n := term.A()

	newA := cloneBigIntOrZero(r.c)
	newB := cloneBigIntOrZero(r.d)
	newC := new(big.Int).Sub(cloneBigIntOrZero(r.a), new(big.Int).Mul(n, r.c))
	newD := new(big.Int).Sub(cloneBigIntOrZero(r.b), new(big.Int).Mul(n, r.d))

	r.a = newA
	r.b = newB
	r.c = newC
	r.d = newD
	return r
}

func (r *Rectifier) CanEmitRCFTerm(_ Range) (RCFTerm, bool) {
	return r.CanEmit()
}

// core/rectifier.go V4
