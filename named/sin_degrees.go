// named/sin_degrees.go v5
package named

import (
	"math/big"

	"github.com/egp/gosper-gcf/core"
	"github.com/egp/gosper-gcf/trig"
)

func SinDegrees(x core.PQStream) *core.GCF {
	return trig.Sin(degreesToRadiansForSin(x))
}

func degreesToRadiansForSin(x core.PQStream) core.PQStream {
	if x == nil {
		panic("degreesToRadiansForSin: nil input")
	}

	if exact, ok := exactClosedCurrentValueSinDegrees(x); ok {
		if exact.Cmp(core.RationalFromInt64(0)) == 0 {
			return core.PQStreamFromRational(core.RationalFromInt64(0))
		}
		if exact.Cmp(core.RationalFromInt64(180)) == 0 {
			return Pi()
		}
	}

	radians := core.NewGCF2(
		core.BLFTCoefficients{
			A: big.NewInt(1),
			B: big.NewInt(0),
			C: big.NewInt(0),
			D: big.NewInt(0),
			E: big.NewInt(0),
			F: big.NewInt(0),
			G: big.NewInt(0),
			H: big.NewInt(180),
		},
		x,
		Pi(),
	)

	return core.PQStreamFromRCF(radians)
}

func exactClosedCurrentValueSinDegrees(x core.PQStream) (core.Rational, bool) {
	rng, err := x.Range()
	if err != nil {
		return core.Rational{}, false
	}
	if !rng.Inside {
		return core.Rational{}, false
	}
	if rng.Lo.Open || rng.Hi.Open {
		return core.Rational{}, false
	}
	if rng.Lo.Value.Cmp(rng.Hi.Value) != 0 {
		return core.Rational{}, false
	}
	return rng.Lo.Value, true
}

// named/sin_degrees.go v5
