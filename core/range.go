// core/range.go v2
package core

type RangeKind int

const (
	InsideInterval RangeKind = iota
	OutsideInterval
)

type Range struct {
	Lo     Endpoint
	Hi     Endpoint
	Inside bool
}

func (r Range) Kind() RangeKind {
	if r.Inside {
		return InsideInterval
	}
	return OutsideInterval
}

// Cmp currently has only the final shape.
// Full uncertainty-order semantics remain to be implemented.
func (r Range) Cmp(_ Range) int {
	return 0
}

// core/range.go v2
