// trig/api.go v4
package trig

import "github.com/egp/gosper-gcf/core"

func Sin(x core.PQStream) *core.GCF {
	return sinRadians(x)
}

func Tanh(x core.PQStream) *core.GCF {
	return tanhRadians(x)
}

// trig/api.go v4
