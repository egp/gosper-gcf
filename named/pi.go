// named/pi.go v8
package named

import (
	"fmt"
	"math/big"

	"github.com/egp/gosper-gcf/core"
)

type piGaussStream struct {
	index int
}

func Pi() core.PQStream {
	return &piGaussStream{index: 0}
}

func piGaussTermAt(index int) (core.PQTerm, error) {
	if index < 0 {
		return core.PQTerm{}, fmt.Errorf("piGaussTermAt: negative index %d", index)
	}
	if index == 0 {
		return core.PQTerm{
			P: big.NewInt(0),
			Q: big.NewInt(4),
		}, nil
	}

	n := int64(index)
	return core.PQTerm{
		P: big.NewInt(2*n - 1),
		Q: big.NewInt(n * n),
	}, nil
}

func piGaussLookaheadRange(index int) (core.Range, error) {
	if index < 0 {
		return core.Range{}, fmt.Errorf("piGaussLookaheadRange: negative index %d", index)
	}
	if index == 0 {
		return core.Range{
			Lo: core.Endpoint{
				Value: core.RationalFromInt64(0),
				Open:  true,
			},
			Hi: core.Endpoint{
				Value: core.RationalFromInt64(4),
				Open:  true,
			},
			Inside: true,
		}, nil
	}

	n := int64(index)
	p := 2*n - 1
	q := n * n
	nextOdd := 2*n + 1

	return core.Range{
		Lo: core.Endpoint{
			Value: core.RationalFromInt64(p),
			Open:  true,
		},
		Hi: core.Endpoint{
			Value: core.NewRational(
				big.NewInt(p*nextOdd+q),
				big.NewInt(nextOdd),
			),
			Open: true,
		},
		Inside: true,
	}, nil
}

func (s *piGaussStream) NextPQ() (core.PQTerm, core.PQStream, core.Status, error) {
	if s == nil {
		return core.PQTerm{}, s, core.StatusEOF, fmt.Errorf("piGaussStream.NextPQ: %w", core.ErrNilReceiver)
	}

	term, err := piGaussTermAt(s.index)
	if err != nil {
		return core.PQTerm{}, s, core.StatusEOF, err
	}

	return term, &piGaussStream{index: s.index + 1}, core.StatusOK, nil
}

func (s *piGaussStream) Range() (core.Range, error) {
	if s == nil {
		return core.Range{}, fmt.Errorf("piGaussStream.Range: %w", core.ErrNilReceiver)
	}
	return piGaussLookaheadRange(s.index)
}

// named/pi.go v8
