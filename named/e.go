// named/e.go v6
package named

import (
	"math/big"

	"github.com/egp/gosper-gcf/core"
)

type eProceduralStream struct {
	index int
}

func E() core.PQStream {
	return &eProceduralStream{index: 1}
}

func (s *eProceduralStream) NextPQ() (core.PQTerm, core.PQStream, core.Status) {
	if s == nil {
		panic("eProceduralStream.NextPQ: nil receiver")
	}

	return core.PQTerm{
			P: big.NewInt(eTermAt(s.index)),
			Q: big.NewInt(1),
		},
		&eProceduralStream{index: s.index + 1},
		core.StatusOK
}

func (s *eProceduralStream) Range() core.Range {
	if s == nil {
		panic("eProceduralStream.Range: nil receiver")
	}

	return eLookaheadRange(
		eTermAt(s.index),
		eTermAt(s.index+1),
	)
}

func eTermAt(index int) int64 {
	if index <= 0 {
		panic("eTermAt: index must be >= 1")
	}
	if index == 1 {
		return 2
	}

	offset := index - 2
	if offset%3 == 1 {
		return int64(2 * (offset/3 + 1))
	}
	return 1
}

func eLookaheadRange(a, next int64) core.Range {
	if next <= 0 {
		panic("eLookaheadRange: next must be > 0")
	}

	loNum := big.NewInt(a*(next+1) + 1)
	loDen := big.NewInt(next + 1)

	hiNum := big.NewInt(a*next + 1)
	hiDen := big.NewInt(next)

	return core.Range{
		Lo: core.Endpoint{
			Value: core.NewRational(loNum, loDen),
			Open:  true,
		},
		Hi: core.Endpoint{
			Value: core.NewRational(hiNum, hiDen),
			Open:  true,
		},
		Inside: true,
	}
}

// named/e.go v6
