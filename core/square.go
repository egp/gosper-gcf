// core/square.go v2
package core

import "math/big"

func Square(x PQStream) *GCF {
	return NewDLFT1(
		DLFTCoefficients{
			A: big.NewInt(1),
			B: big.NewInt(0),
			C: big.NewInt(0),
			D: big.NewInt(0),
			E: big.NewInt(0),
			F: big.NewInt(1),
		},
		x,
	)
}

// core/square.go v2
