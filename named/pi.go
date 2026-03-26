// named/pi.go v1
package named

import (
	"math/big"

	"github.com/egp/gosper-gcf/core"
)

func Pi() core.PQStream {
	stream, status := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term: core.PQTerm{
				P: big.NewInt(0),
				Q: big.NewInt(1),
			},
			Range: core.Range{
				Lo: core.Endpoint{
					Value: core.RationalFromInt64(0),
					Open:  false,
				},
				Hi: core.Endpoint{
					Value: core.RationalFromInt64(0),
					Open:  false,
				},
				Inside: true,
			},
		},
	})
	if status != core.StatusOK {
		panic("named.Pi: failed to construct stub source")
	}
	return stream
}

// named/pi.go v1
