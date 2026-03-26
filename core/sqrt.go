// core/sqrt.go v2
package core

func Sqrt(x PQStream) *GCF {
	controller := newSqrtController(x)
	return controller.buildRefinement(controller.seedApproximation())
}

// core/sqrt.go v2
