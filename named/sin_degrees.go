// named/sin_degrees.go v6
package named

import (
	"math/big"

	"github.com/egp/gosper-gcf/core"
	"github.com/egp/gosper-gcf/trig"
)

// sinDegreesExactTable maps (integer degrees mod 360) → exact rational sin value
// for every multiple of 30° where sin is rational.
var sinDegreesExactTable = map[int64]core.Rational{
	0:   core.RationalFromInt64(0),
	30:  core.NewRational(big.NewInt(1), big.NewInt(2)),
	90:  core.RationalFromInt64(1),
	150: core.NewRational(big.NewInt(1), big.NewInt(2)),
	180: core.RationalFromInt64(0),
	210: core.NewRational(big.NewInt(-1), big.NewInt(2)),
	270: core.RationalFromInt64(-1),
	330: core.NewRational(big.NewInt(-1), big.NewInt(2)),
}

func SinDegrees(x core.PQStream) *core.GCF {
	if g, ok := exactSinDegreesShortcut(x); ok {
		return g
	}
	return trig.Sin(degreesToRadiansForSin(x))
}

// exactSinDegreesShortcut returns an exact terminal GCF when the input is an
// integer-degree angle whose sine is a rational number (multiples of 30°).
func exactSinDegreesShortcut(x core.PQStream) (*core.GCF, bool) {
	val, ok := exactClosedCurrentValueSinDegrees(x)
	if !ok {
		return nil, false
	}
	// Only integer angles (denominator 1) have entries in the table.
	if val.Den().Cmp(big.NewInt(1)) != 0 {
		return nil, false
	}
	// Reduce mod 360; big.Int.Mod always returns a non-negative result when
	// the divisor is positive, so deg360 is in [0, 359].
	deg360 := new(big.Int).Mod(val.Num(), big.NewInt(360))
	sinVal, found := sinDegreesExactTable[deg360.Int64()]
	if !found {
		return nil, false
	}
	return core.NewExactTerminalGCFFromRational(sinVal), true
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

// named/sin_degrees.go v6
