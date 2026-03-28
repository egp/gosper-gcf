// trig/tanh.go v1
package trig

import "github.com/egp/gosper-gcf/core"

func tanhRadians(x core.PQStream) *core.GCF {
	if x == nil {
		panic("tanhRadians: nil input")
	}
	return hyperbolicHalfAngleKernel(x)
}

// trig/tanh.go v1
