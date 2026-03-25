// core/div.go v2
package core

import "math/big"

func Div(x, y PQStream) *GCF {
	return NewGCF2(
		BLFTCoefficients{
			A: big.NewInt(0),
			B: big.NewInt(1),
			C: big.NewInt(0),
			D: big.NewInt(0),
			E: big.NewInt(0),
			F: big.NewInt(0),
			G: big.NewInt(1),
			H: big.NewInt(0),
		},
		x,
		y,
	)
}

// core/div.go v2
