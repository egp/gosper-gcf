// core/sqrt_cycle4.go v2
package core

func (c *sqrtController) activeApproximation() PQStream {
	if c == nil {
		panic("sqrtController.activeApproximation: nil receiver")
	}
	if c.hasOurorobosFeedback {
		return c.ourorobosApprox
	}
	return c.seed
}

func (c *sqrtController) currentRefinement() *GCF {
	if c == nil {
		panic("sqrtController.currentRefinement: nil receiver")
	}
	return c.buildRefinement(c.activeApproximation())
}

// core/sqrt_cycle4.go v2
