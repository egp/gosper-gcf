// core/blft_collapse.go v2
package core

import "math/big"

func (s blftState) CollapseX() TransformCoefficients {
	return TransformCoefficients{
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

func (s blftState) CollapseY() TransformCoefficients {
	return TransformCoefficients{
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

// core/blft_collapse.go v2
