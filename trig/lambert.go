// trig/lambert.go v13
package trig

import (
	"fmt"
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
		half = &errorPQReplay{err: fmt.Errorf("newLambertKernel: nil half stream")}
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
		return errorGCF(fmt.Errorf("(*lambertKernel).Root: %w", core.ErrNilReceiver))
	}
	return k.previewRoot()
}

func (k *lambertKernel) previewDepth() int {
	if k == nil {
		return 0
	}
	switch k.mode {
	case lambertModeCircular:
		return 5
	case lambertModeHyperbolic:
		return 7
	default:
		return 0
	}
}

func (k *lambertKernel) previewRoot() *core.GCF {
	if k == nil {
		return errorGCF(fmt.Errorf("(*lambertKernel).previewRoot: %w", core.ErrNilReceiver))
	}
	depth := k.previewDepth()
	if depth <= 0 {
		return errorGCF(fmt.Errorf("(*lambertKernel).previewRoot: invalid preview depth"))
	}
	return k.truncatedAt(0, depth)
}

func (k *lambertKernel) stageAt(stage int, tail core.PQStream) *core.GCF {
	if k == nil {
		return errorGCF(fmt.Errorf("(*lambertKernel).stageAt: %w", core.ErrNilReceiver))
	}
	if stage < 0 {
		return errorGCF(fmt.Errorf("(*lambertKernel).stageAt: negative stage"))
	}
	if tail == nil {
		tail = &errorPQReplay{err: fmt.Errorf("(*lambertKernel).stageAt: nil tail")}
	}
	odd, err := k.mode.stageOdd(stage)
	if err != nil {
		return errorGCF(err)
	}
	return lambertStageGCF(k.mode, odd, k.half, tail)
}

func (k *lambertKernel) truncatedAt(stage, depth int) *core.GCF {
	if k == nil {
		return errorGCF(fmt.Errorf("(*lambertKernel).truncatedAt: %w", core.ErrNilReceiver))
	}
	if stage < 0 {
		return errorGCF(fmt.Errorf("(*lambertKernel).truncatedAt: negative stage"))
	}
	if depth <= 0 {
		return errorGCF(fmt.Errorf("(*lambertKernel).truncatedAt: nonpositive depth"))
	}
	if depth == 1 {
		return k.stageAt(stage, core.PQStreamFromRational(core.RationalFromInt64(0)))
	}
	deeper := k.truncatedAt(stage+1, depth-1)
	return k.stageAt(stage, newPQFromRCFReplay(deeper))
}

func (m lambertMode) stageOdd(stage int) (*big.Int, error) {
	if stage < 0 {
		return nil, fmt.Errorf("lambertMode.stageOdd: negative stage")
	}
	return big.NewInt(int64(2*stage + 1)), nil
}

func (m lambertMode) xySign() (int, error) {
	switch m {
	case lambertModeCircular:
		return -1, nil
	case lambertModeHyperbolic:
		return 1, nil
	default:
		return 0, fmt.Errorf("lambertMode.xySign: invalid mode")
	}
}

func (m lambertMode) stageCoefficients(odd *big.Int) (core.BLFTCoefficients, error) {
	if odd == nil {
		return core.BLFTCoefficients{}, fmt.Errorf("lambertMode.stageCoefficients: nil odd")
	}
	if odd.Sign() <= 0 {
		return core.BLFTCoefficients{}, fmt.Errorf("lambertMode.stageCoefficients: odd must be positive")
	}
	sign, err := m.xySign()
	if err != nil {
		return core.BLFTCoefficients{}, err
	}
	return core.BLFTCoefficients{
		A: big.NewInt(0),
		B: big.NewInt(1),
		C: big.NewInt(0),
		D: big.NewInt(0),
		E: big.NewInt(int64(sign)),
		F: big.NewInt(0),
		G: big.NewInt(0),
		H: new(big.Int).Set(odd),
	}, nil
}

func lambertStageGCF(mode lambertMode, odd *big.Int, x, y core.PQStream) *core.GCF {
	if x == nil {
		x = &errorPQReplay{err: fmt.Errorf("lambertStageGCF: nil x")}
	}
	if y == nil {
		y = &errorPQReplay{err: fmt.Errorf("lambertStageGCF: nil y")}
	}
	coeffs, err := mode.stageCoefficients(odd)
	if err != nil {
		return errorGCF(err)
	}
	return core.NewGCF2(coeffs, x, y)
}

type pqFromRCFReplay struct {
	fork *replayRCFFork
}

func newPQFromRCFReplay(src core.RCFStream) core.PQStream {
	if src == nil {
		return &errorPQReplay{err: fmt.Errorf("newPQFromRCFReplay: nil source")}
	}
	root := newReplayRCF(src)
	return &pqFromRCFReplay{
		fork: root.Fork(),
	}
}

func (p *pqFromRCFReplay) NextPQ() (core.PQTerm, core.PQStream, core.Status, error) {
	if p == nil || p.fork == nil {
		return core.PQTerm{}, p, core.StatusEOF, fmt.Errorf("(*pqFromRCFReplay).NextPQ: %w", core.ErrNilReceiver)
	}

	term, status, err := p.fork.NextRCF()
	if err != nil {
		return core.PQTerm{}, p, core.StatusEOF, err
	}

	switch status {
	case core.StatusOK:
		return core.PQTerm{
			P: new(big.Int).Set(term.A()),
			Q: big.NewInt(1),
		}, p, core.StatusOK, nil
	case core.StatusEOF:
		return core.PQTerm{}, p, core.StatusEOF, nil
	default:
		return core.PQTerm{}, p, status, fmt.Errorf("(*pqFromRCFReplay).NextPQ: invalid input status=%v", status)
	}
}

func (p *pqFromRCFReplay) Range() (core.Range, error) {
	if p == nil || p.fork == nil {
		return core.Range{}, fmt.Errorf("(*pqFromRCFReplay).Range: %w", core.ErrNilReceiver)
	}
	return p.fork.Range()
}

func errorGCF(err error) *core.GCF {
	return core.NewGCF1(
		core.BLFTCoefficients{
			A: big.NewInt(0),
			B: big.NewInt(1),
			C: big.NewInt(0),
			D: big.NewInt(0),
			E: big.NewInt(0),
			F: big.NewInt(0),
			G: big.NewInt(0),
			H: big.NewInt(1),
		},
		&errorPQReplay{err: err},
	)
}

// trig/lambert.go v13
