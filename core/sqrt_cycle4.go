// core/sqrt_cycle4.go v3
package core

import "fmt"

func (c *sqrtController) activeApproximation() PQStream {
	if c == nil {
		return newErrorPQStream(fmt.Errorf("sqrtController.activeApproximation: %w", ErrNilReceiver))
	}
	if c.initErr != nil {
		return newErrorPQStream(c.initErr)
	}
	if c.hasOurorobosFeedback {
		return c.ourorobosApprox
	}
	return c.seed
}

func (c *sqrtController) currentRefinement() *GCF {
	if c == nil {
		return newObservedRCFGCF(newErrorRCFStream(fmt.Errorf("sqrtController.currentRefinement: %w", ErrNilReceiver)))
	}
	return c.buildRefinement(c.activeApproximation())
}

// core/sqrt_cycle4.go v3
