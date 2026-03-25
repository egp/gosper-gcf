// named/sqrt2.go v3
package named

import (
	"math/big"

	"github.com/egp/gosper-gcf/core"
)

func Sqrt2() core.PQStream {
	first := core.FinitePQStep{
		Term: core.PQTerm{
			P: big.NewInt(1),
			Q: big.NewInt(1),
		},
		Range: sqrt2HeadRange(),
	}

	stream, status := core.NewProceduralPQStream(first, sqrt2TailNext)
	if status != core.StatusOK {
		panic("named.Sqrt2: failed to construct procedural PQ stream")
	}

	return stream
}

func sqrt2TailNext() (core.FinitePQStep, core.ProceduralPQNext, core.Status) {
	return core.FinitePQStep{
		Term: core.PQTerm{
			P: big.NewInt(2),
			Q: big.NewInt(1),
		},
		Range: sqrt2TailRange(),
	}, sqrt2TailNext, core.StatusOK
}

func sqrt2HeadRange() core.Range {
	return core.Range{
		Lo: core.Endpoint{
			Value: core.RationalFromInt64(1),
			Open:  false,
		},
		Hi: core.Endpoint{
			Value: core.RationalFromInt64(2),
			Open:  false,
		},
		Inside: true,
	}
}

func sqrt2TailRange() core.Range {
	return core.Range{
		Lo: core.Endpoint{
			Value: core.RationalFromInt64(2),
			Open:  false,
		},
		Hi: core.Endpoint{
			Value: core.RationalFromInt64(3),
			Open:  false,
		},
		Inside: true,
	}
}

// named/sqrt2.go v3
