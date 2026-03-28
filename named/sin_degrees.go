// named/sin_degrees.go v4
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

// named/sin_degrees.go v4
