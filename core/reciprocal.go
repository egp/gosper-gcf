// core/reciprocal.go v2
package core

import "math/big"

func Reciprocal(x PQStream) *GCF {
	return NewGCF1(
		BLFTCoefficients{
			A: big.NewInt(0),
			B: big.NewInt(0),
			C: big.NewInt(0),
			D: big.NewInt(1),
			E: big.NewInt(0),
			F: big.NewInt(1),
			G: big.NewInt(0),
			H: big.NewInt(0),
		},
		x,
	)
}

// core/reciprocal.go v2
