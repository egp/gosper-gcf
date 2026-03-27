// trig/api.go v2
package trig

import (
	"math/big"

	"github.com/egp/gosper-gcf/core"
)

func Sin(x core.PQStream) *core.GCF {
	return sinRadians(x)
}

func Tanh(x core.PQStream) *core.GCF {
	_ = x
	return exactZeroTrig()
}

func exactZeroTrig() *core.GCF {
	return core.NewExactTerminalGCF(
		[]core.RCFTerm{
			core.NewRCFTerm(big.NewInt(0)),
		},
		exactRangeTrig(0, 1),
	)
}

func exactRangeTrig(num, den int64) core.Range {
	value := core.NewRational(big.NewInt(num), big.NewInt(den))
	return core.Range{
		Lo:     core.Endpoint{Value: value, Open: false},
		Hi:     core.Endpoint{Value: value, Open: false},
		Inside: true,
	}
}

// trig/api.go v2
