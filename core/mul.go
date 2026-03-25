// core/mul.go v2
package core

import "math/big"

func Mul(x, y PQStream) *GCF {
	return NewGCF2(
		BLFTCoefficients{
			A: big.NewInt(1),
			B: big.NewInt(0),
			C: big.NewInt(0),
			D: big.NewInt(0),
			E: big.NewInt(0),
			F: big.NewInt(0),
			G: big.NewInt(0),
			H: big.NewInt(1),
		},
		x,
		y,
	)
}

// core/mul.go v2
