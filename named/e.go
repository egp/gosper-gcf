// named/e.go v8
package named

import (
	"fmt"
	"math/big"

	"github.com/egp/gosper-gcf/core"
)

type eProceduralStream struct {
	index int
}

func E() core.PQStream {
	return &eProceduralStream{index: 1}
}

func (s *eProceduralStream) NextPQ() (core.PQTerm, core.PQStream, core.Status, error) {
	if s == nil {
		return core.PQTerm{}, s, core.StatusEOF, fmt.Errorf("eProceduralStream.NextPQ: %w", core.ErrNilReceiver)
	}
	term, err := eTermAt(s.index)
	if err != nil {
		return core.PQTerm{}, s, core.StatusEOF, err
	}
	return core.PQTerm{
		P: big.NewInt(term),
		Q: big.NewInt(1),
	}, &eProceduralStream{index: s.index + 1}, core.StatusOK, nil
}

func (s *eProceduralStream) CurrentInterval() (core.Interval, error) {
	if s == nil {
		return core.Interval{}, fmt.Errorf("eProceduralStream.CurrentInterval: %w", core.ErrNilReceiver)
	}
	a, err := eTermAt(s.index)
	if err != nil {
		return core.Interval{}, err
	}
	next, err := eTermAt(s.index + 1)
	if err != nil {
		return core.Interval{}, err
	}
	return eLookaheadRange(a, next)
}

func (s *eProceduralStream) Range() (core.Range, error) {
	return s.CurrentInterval()
}

func eTermAt(index int) (int64, error) {
	if index <= 0 {
		return 0, fmt.Errorf("eTermAt: index must be >= 1, got %d", index)
	}
	if index == 1 {
		return 2, nil
	}
	offset := index - 2
	if offset%3 == 1 {
		return int64(2 * (offset/3 + 1)), nil
	}
	return 1, nil
}

func eLookaheadRange(a, next int64) (core.Range, error) {
	if next <= 0 {
		return core.Range{}, fmt.Errorf("eLookaheadRange: next must be > 0, got %d", next)
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
	}, nil
}

// named/e.go v8
