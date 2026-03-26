// core/sqrt.go v3
package core

func Sqrt(x PQStream) *GCF {
	controller := newSqrtController(x)
	return newObservedRCFGCF(newSqrtObservedRefinementStream(controller))
}

// core/sqrt.go v3
