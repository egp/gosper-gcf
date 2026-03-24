// core/blft_emit.go v3
package core

import "math/big"

func (s blftState) CanEmitRCFTerm(r Range) (RCFTerm, bool) {
	if !r.Inside {
		return NewRCFTerm(nil), false
	}

	loFloor, _ := floorQuoRem(r.Lo.Value.Num(), r.Lo.Value.Den())
	hiFloor, _ := floorQuoRem(r.Hi.Value.Num(), r.Hi.Value.Den())

	if loFloor.Cmp(hiFloor) != 0 {
		return NewRCFTerm(nil), false
	}

	return NewRCFTerm(loFloor), true
}

func (s blftState) Emit(term RCFTerm) blftState {
	n := term.A()

	return blftState{
		A: cloneBigIntOrZero(s.E),
		B: cloneBigIntOrZero(s.F),
		C: cloneBigIntOrZero(s.G),
		D: cloneBigIntOrZero(s.H),
		E: subMul(s.A, n, s.E),
		F: subMul(s.B, n, s.F),
		G: subMul(s.C, n, s.G),
		H: subMul(s.D, n, s.H),
	}.Normalize()
}

func subMul(x, n, y *big.Int) *big.Int {
	left := cloneBigIntOrZero(x)
	right := mul(n, y)
	return left.Sub(left, right)
}

func cloneBigIntOrZero(x *big.Int) *big.Int {
	if x == nil {
		return big.NewInt(0)
	}
	return cloneBigInt(x)
}

// core/blft_emit.go v3
