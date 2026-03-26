// named/e_cycle1.go v2
package named

import (
	"math/big"

	"github.com/egp/gosper-gcf/core"
)

func eTermAt(index int) int64 {
	if index <= 0 {
		panic("eTermAt: index must be >= 1")
	}
	if index == 1 {
		return 2
	}

	m := index - 2
	switch m % 3 {
	case 0:
		return 1
	case 1:
		return int64(2 * (m/3 + 1))
	default:
		return 1
	}
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
			Open:  false,
		},
		Hi: core.Endpoint{
			Value: core.NewRational(hiNum, hiDen),
			Open:  true,
		},
		Inside: true,
	}
}

// named/e_cycle1.go v2
