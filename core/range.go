// core/range.go v4
package core

type RangeKind int

const (
	RangeArc RangeKind = iota
	RangeFull
)

// Transitional compatibility aliases.
// Under the new projective-arc model, Kind() distinguishes Arc vs Full,
// while the Inside field distinguishes inside vs outside arcs.
const (
	InsideInterval  = RangeArc
	OutsideInterval = RangeArc
)

type Range struct {
	Lo     Endpoint
	Hi     Endpoint
	Inside bool
	Kind_  RangeKind
}

func (r Range) Kind() RangeKind {
	if r.Kind_ == RangeFull {
		return RangeFull
	}
	return RangeArc
}

// Cmp currently has only the final shape.
// Full uncertainty-order semantics remain to be implemented.
func (r Range) Cmp(_ Range) int {
	return 0
}

// core/range.go v4
