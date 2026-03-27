// named/sin_degrees.go v1
package named

import (
	"github.com/egp/gosper-gcf/core"
	"github.com/egp/gosper-gcf/trig"
)

func SinDegrees(x core.PQStream) *core.GCF {
	return trig.Sin(x)
}

// named/sin_degrees.go v1
