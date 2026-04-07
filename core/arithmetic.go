// core/add.go v2 (unchanged — speculator tier only)
package core

import "math/big"

func Add(x, y PQStream) *GCF {
	return NewGCF2(
		BLFTCoefficients{
			A: big.NewInt(0),
			B: big.NewInt(1),
			C: big.NewInt(1),
			D: big.NewInt(0),
			E: big.NewInt(0),
			F: big.NewInt(0),
			G: big.NewInt(0),
			H: big.NewInt(1),
		},
		x, y,
	)
}

// --- appended from core/div.go ---
// core/div.go v2

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

// --- appended from core/mul.go ---
// core/mul.go v2 (unchanged — speculator tier only)

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
		x, y,
	)
}
