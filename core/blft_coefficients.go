// core/blftcoefficients.go v1
package core

import "math/big"

type BLFTCoefficients struct {
	A *big.Int
	B *big.Int
	C *big.Int
	D *big.Int
	E *big.Int
	F *big.Int
	G *big.Int
	H *big.Int
}

func cloneBLFTCoefficients(tc BLFTCoefficients) BLFTCoefficients {
	return BLFTCoefficients{
		A: cloneBigInt(tc.A),
		B: cloneBigInt(tc.B),
		C: cloneBigInt(tc.C),
		D: cloneBigInt(tc.D),
		E: cloneBigInt(tc.E),
		F: cloneBigInt(tc.F),
		G: cloneBigInt(tc.G),
		H: cloneBigInt(tc.H),
	}
}

// core/blftcoefficients.go v1
