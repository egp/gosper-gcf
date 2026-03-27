// named/pi.go v7
package named

import (
	"math/big"

	"github.com/egp/gosper-gcf/core"
)

type piGaussStream struct {
	index int
}

func Pi() core.PQStream {
	return &piGaussStream{index: 0}
}

func piGaussTermAt(index int) core.PQTerm {
	if index < 0 {
		panic("piGaussTermAt: negative index")
	}

	if index == 0 {
		return core.PQTerm{
			P: big.NewInt(0),
			Q: big.NewInt(4),
		}
	}

	n := int64(index)
	return core.PQTerm{
		P: big.NewInt(2*n - 1),
		Q: big.NewInt(n * n),
	}
}

func piGaussLookaheadRange(index int) core.Range {
	if index < 0 {
		panic("piGaussLookaheadRange: negative index")
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
		}
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
	}
}

func (s *piGaussStream) NextPQ() (core.PQTerm, core.PQStream, core.Status) {
	if s == nil {
		panic("piGaussStream.NextPQ: nil receiver")
	}

	return piGaussTermAt(s.index), &piGaussStream{index: s.index + 1}, core.StatusOK
}

func (s *piGaussStream) Range() core.Range {
	if s == nil {
		panic("piGaussStream.Range: nil receiver")
	}

	return piGaussLookaheadRange(s.index)
}

// named/pi.go v7
