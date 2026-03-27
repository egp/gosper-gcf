// trig/lambert.go v12
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
	return k.previewRoot()
}

func (k *lambertKernel) previewDepth() int {
	if k == nil {
		panic("(*lambertKernel).previewDepth: nil receiver")
	}

	switch k.mode {
	case lambertModeCircular:
		return 5
	case lambertModeHyperbolic:
		return 7
	default:
		panic("(*lambertKernel).previewDepth: invalid mode")
	}
}

func (k *lambertKernel) previewRoot() *core.GCF {
	if k == nil {
		panic("(*lambertKernel).previewRoot: nil receiver")
	}
	return k.truncatedAt(0, k.previewDepth())
}

func (k *lambertKernel) stageAt(stage int, tail core.PQStream) *core.GCF {
	if k == nil {
		panic("(*lambertKernel).stageAt: nil receiver")
	}
	if stage < 0 {
		panic("(*lambertKernel).stageAt: negative stage")
	}
	if tail == nil {
		panic("(*lambertKernel).stageAt: nil tail")
	}
	return lambertStageGCF(k.mode, k.mode.stageOdd(stage), k.half, tail)
}

func (k *lambertKernel) truncatedAt(stage, depth int) *core.GCF {
	if k == nil {
		panic("(*lambertKernel).truncatedAt: nil receiver")
	}
	if stage < 0 {
		panic("(*lambertKernel).truncatedAt: negative stage")
	}
	if depth <= 0 {
		panic("(*lambertKernel).truncatedAt: nonpositive depth")
	}

	if depth == 1 {
		return k.stageAt(stage, core.PQStreamFromRational(core.RationalFromInt64(0)))
	}

	deeper := k.truncatedAt(stage+1, depth-1)
	return k.stageAt(stage, newPQFromRCFReplay(deeper))
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

type pqFromRCFReplay struct {
	fork *replayRCFFork
}

func newPQFromRCFReplay(src core.RCFStream) core.PQStream {
	if src == nil {
		panic("newPQFromRCFReplay: nil source")
	}
	root := newReplayRCF(src)
	return &pqFromRCFReplay{
		fork: root.Fork(),
	}
}

func (p *pqFromRCFReplay) NextPQ() (core.PQTerm, core.PQStream, core.Status) {
	if p == nil {
		panic("(*pqFromRCFReplay).NextPQ: nil receiver")
	}

	term, status := p.fork.NextRCF()
	switch status {
	case core.StatusOK:
		return core.PQTerm{
			P: new(big.Int).Set(term.A()),
			Q: big.NewInt(1),
		}, p, core.StatusOK
	case core.StatusEOF:
		return core.PQTerm{}, p, core.StatusEOF
	default:
		panic("(*pqFromRCFReplay).NextPQ: invalid input status")
	}
}

func (p *pqFromRCFReplay) Range() core.Range {
	if p == nil {
		panic("(*pqFromRCFReplay).Range: nil receiver")
	}
	return p.fork.Range()
}

// trig/lambert.go v12
