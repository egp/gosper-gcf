// trig/sin.go v2
package trig

import (
	"math/big"

	"github.com/egp/gosper-gcf/core"
)

func sinRadians(x core.PQStream) *core.GCF {
	return sinFromTanHalf(circularHalfAngleKernel(halfInputForSin(x)))
}

func halfInputForSin(x core.PQStream) core.PQStream {
	if x == nil {
		panic("halfInputForSin: nil input")
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

func sinFromTanHalf(t core.RCFStream) *core.GCF {
	if t == nil {
		panic("sinFromTanHalf: nil input")
	}

	left, right := newPQPairFromRCFReplay(t)

	return core.NewGCF2(
		core.BLFTCoefficients{
			A: big.NewInt(0),
			B: big.NewInt(2),
			C: big.NewInt(0),
			D: big.NewInt(0),
			E: big.NewInt(1),
			F: big.NewInt(0),
			G: big.NewInt(0),
			H: big.NewInt(1),
		},
		left,
		right,
	)
}

func newPQPairFromRCFReplay(src core.RCFStream) (core.PQStream, core.PQStream) {
	if src == nil {
		panic("newPQPairFromRCFReplay: nil source")
	}

	root := newReplayRCF(src)
	return &pqFromRCFReplay{fork: root.Fork()}, &pqFromRCFReplay{fork: root.Fork()}
}

// trig/sin.go v2
