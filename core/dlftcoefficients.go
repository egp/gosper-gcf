// core/dlftcoefficients.go v1
package core

import "math/big"

type DLFTCoefficients struct {
	A *big.Int
	B *big.Int
	C *big.Int
	D *big.Int
	E *big.Int
	F *big.Int
}

func cloneDLFTCoefficients(tc DLFTCoefficients) DLFTCoefficients {
	return DLFTCoefficients{
		A: cloneBigInt(tc.A),
		B: cloneBigInt(tc.B),
		C: cloneBigInt(tc.C),
		D: cloneBigInt(tc.D),
		E: cloneBigInt(tc.E),
		F: cloneBigInt(tc.F),
	}
}

// core/dlftcoefficients.go v1
