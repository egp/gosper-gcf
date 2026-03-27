// trig/lambert.go v4
package trig

import (
	"math/big"

	"github.com/egp/gosper-gcf/core"
)

type lambertMode int

const (
	lambertModeCircular lambertMode = iota
	lambertModeHyperbolic
)

type lambertKernel struct {
	mode lambertMode
	half core.PQStream
}

func newLambertKernel(mode lambertMode, half core.PQStream) *lambertKernel {
	if half == nil {
		panic("newLambertKernel: nil half stream")
	}
	return &lambertKernel{
		mode: mode,
		half: half,
	}
}

func circularHalfAngleKernel(half core.PQStream) *core.GCF {
	return newLambertKernel(lambertModeCircular, half).Root()
}

func hyperbolicHalfAngleKernel(half core.PQStream) *core.GCF {
	return newLambertKernel(lambertModeHyperbolic, half).Root()
}

func (k *lambertKernel) Root() *core.GCF {
	if k == nil {
		panic("(*lambertKernel).Root: nil receiver")
	}
	return exactZeroTrig()
}

func (m lambertMode) stageOdd(stage int) *big.Int {
	_ = m
	if stage < 0 {
		panic("lambertMode.stageOdd: negative stage")
	}
	return big.NewInt(int64(2*stage + 1))
}

func (m lambertMode) xySign() int {
	switch m {
	case lambertModeCircular:
		return -1
	case lambertModeHyperbolic:
		return 1
	default:
		panic("lambertMode.xySign: invalid mode")
	}
}

func (m lambertMode) stageCoefficients(odd *big.Int) core.BLFTCoefficients {
	if odd == nil {
		panic("lambertMode.stageCoefficients: nil odd")
	}
	if odd.Sign() <= 0 {
		panic("lambertMode.stageCoefficients: odd must be positive")
	}

	return core.BLFTCoefficients{
		A: big.NewInt(0),
		B: big.NewInt(1),
		C: big.NewInt(0),
		D: big.NewInt(0),
		E: big.NewInt(int64(m.xySign())),
		F: big.NewInt(0),
		G: big.NewInt(0),
		H: new(big.Int).Set(odd),
	}
}

func lambertStageGCF(mode lambertMode, odd *big.Int, x, y core.PQStream) *core.GCF {
	if x == nil {
		panic("lambertStageGCF: nil x")
	}
	if y == nil {
		panic("lambertStageGCF: nil y")
	}
	return core.NewGCF2(mode.stageCoefficients(odd), x, y)
}

// trig/lambert.go v4
