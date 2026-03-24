package named

import (
	"math/big"

	"github.com/egp/gosper-gcf/core"
)

type sqrt2Head struct{}
type sqrt2Tail struct{}

var (
	sqrt2HeadSingleton = &sqrt2Head{}
	sqrt2TailSingleton = &sqrt2Tail{}
)

func Sqrt2() core.PQStream {
	return sqrt2HeadSingleton
}

func (s *sqrt2Head) NextPQ() (core.PQTerm, core.PQStream, core.Status) {
	return core.PQTerm{
		P: big.NewInt(1),
		Q: big.NewInt(1),
	}, sqrt2TailSingleton, core.StatusOK
}

func (s *sqrt2Head) Range() core.Range {
	return core.Range{
		Lo:     core.FromInt64(1),
		Hi:     core.FromInt64(2),
		LoOpen: false,
		HiOpen: false,
		Inside: true,
	}
}

func (s *sqrt2Tail) NextPQ() (core.PQTerm, core.PQStream, core.Status) {
	return core.PQTerm{
		P: big.NewInt(2),
		Q: big.NewInt(1),
	}, sqrt2TailSingleton, core.StatusOK
}

func (s *sqrt2Tail) Range() core.Range {
	return core.Range{
		Lo:     core.FromInt64(1),
		Hi:     core.FromInt64(2),
		LoOpen: false,
		HiOpen: false,
		Inside: true,
	}
}
