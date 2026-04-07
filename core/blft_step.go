// core/blft_choose.go v5
package core

import "math/big"

func chooseIngestX(xRange, yRange Range) bool {
	if xRange.Inside && yRange.Inside {
		xWidth := rangeWidth(xRange)
		yWidth := rangeWidth(yRange)
		return xWidth.Cmp(yWidth) >= 0
	}

	if isExactClosedRangeChoose(xRange) {
		return true
	}
	if isExactClosedRangeChoose(yRange) {
		return false
	}

	if xRange.Inside != yRange.Inside {
		return xRange.Inside
	}

	xWidth := rangeWidth(xRange)
	yWidth := rangeWidth(yRange)
	return xWidth.Cmp(yWidth) >= 0
}

func isExactClosedRangeChoose(r Range) bool {
	return r.Inside &&
		!r.Lo.Open &&
		!r.Hi.Open &&
		r.Lo.Value.Cmp(r.Hi.Value) == 0
}

// core/blft_choose.go v5

// --- appended from core/blft_collapse.go ---
// core/blft_collapse.go v5

func (s blftState) CollapseX() BLFTCoefficients {
	return BLFTCoefficients{
		A: big.NewInt(0),
		B: big.NewInt(0),
		C: cloneBigInt(s.C),
		D: cloneBigInt(s.D),
		E: big.NewInt(0),
		F: big.NewInt(0),
		G: cloneBigInt(s.G),
		H: cloneBigInt(s.H),
	}
}

func (s blftState) CollapseY() BLFTCoefficients {
	return BLFTCoefficients{
		A: big.NewInt(0),
		B: cloneBigInt(s.B),
		C: big.NewInt(0),
		D: cloneBigInt(s.D),
		E: big.NewInt(0),
		F: cloneBigInt(s.F),
		G: big.NewInt(0),
		H: cloneBigInt(s.H),
	}
}

func (s blftState) CollapseToRational() Rational {
	return NewRational(s.D, s.H)
}

// core/blft_collapse.go v5

// --- appended from core/blft_emit.go ---
// core/blft_emit.go v4

func (s blftState) CanEmitRCFTerm(r Range) (RCFTerm, bool) {
	return canEmitRCFTermFromRange(r)
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

// core/blft_emit.go v4
