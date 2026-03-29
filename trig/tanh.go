// trig/tanh.go v4
package trig

import (
	"math/big"

	"github.com/egp/gosper-gcf/core"
)

func tanhRadians(x core.PQStream) *core.GCF {
	if x == nil {
		panic("tanhRadians: nil input")
	}

	return tanhFromTanhHalf(hyperbolicHalfAngleKernel(halfInputForTanh(x)))
}

func halfInputForTanh(x core.PQStream) core.PQStream {
	if x == nil {
		panic("halfInputForTanh: nil input")
	}

	half := core.NewGCF1(
		core.BLFTCoefficients{
			A: big.NewInt(0),
			B: big.NewInt(1),
			C: big.NewInt(0),
			D: big.NewInt(0),
			E: big.NewInt(0),
			F: big.NewInt(0),
			G: big.NewInt(0),
			H: big.NewInt(2),
		},
		x,
	)

	return newPQFromRCFReplay(half)
}

func tanhFromTanhHalf(t core.RCFStream) *core.GCF {
	if t == nil {
		panic("tanhFromTanhHalf: nil input")
	}

	return doubleAngleFromHalfQuotient(t)
}

// trig/tanh.go v4
