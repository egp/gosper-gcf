// core/sqrt.go v1
package core

func Sqrt(x PQStream) *GCF {
	_ = x
	return NewExactTerminalGCF(
		nil,
		exactRangeFromRational(RationalFromInt64(0)),
	)
}

// core/sqrt.go v1
