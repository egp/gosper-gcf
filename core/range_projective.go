// core/range_projective.go v3
package core

import "errors"

var (
	ErrEmptyRange                = errors.New("empty projective range")
	ErrNonCanonicalRange         = errors.New("noncanonical projective range")
	ErrUnsupportedRangeCase      = errors.New("unsupported projective range case")
	ErrUnresolvedProjectiveImage = errors.New("unresolved projective image")
)

func NormalizeRange(r Range) (Range, error) {
	if r.Kind() == RangeFull {
		return canonicalFullRange(), nil
	}

	cmp := r.Lo.Value.Cmp(r.Hi.Value)

	if r.Inside {
		if cmp > 0 {
			return Range{}, ErrNonCanonicalRange
		}
		if cmp == 0 {
			if r.Lo.Open || r.Hi.Open {
				return Range{}, ErrEmptyRange
			}
			return exactArcRange(r.Lo.Value), nil
		}
		return Range{
			Lo:     r.Lo,
			Hi:     r.Hi,
			Inside: true,
			Kind_:  RangeArc,
		}, nil
	}

	if cmp < 0 {
		return Range{}, ErrNonCanonicalRange
	}
	if cmp == 0 {
		switch {
		case !r.Lo.Open && !r.Hi.Open:
			return exactArcRange(r.Lo.Value), nil
		case r.Lo.Open && r.Hi.Open:
			return Range{}, ErrEmptyRange
		default:
			return canonicalFullRange(), nil
		}
	}

	return Range{
		Lo:     r.Lo,
		Hi:     r.Hi,
		Inside: false,
		Kind_:  RangeArc,
	}, nil
}

func ComplementRange(r Range) (Range, error) {
	norm, err := NormalizeRange(r)
	if err != nil {
		return Range{}, err
	}
	if norm.Kind() == RangeFull {
		return Range{}, ErrEmptyRange
	}

	comp := Range{
		Lo: Endpoint{
			Value: norm.Hi.Value,
			Open:  !norm.Hi.Open,
		},
		Hi: Endpoint{
			Value: norm.Lo.Value,
			Open:  !norm.Lo.Open,
		},
		Inside: !norm.Inside,
		Kind_:  RangeArc,
	}

	return NormalizeRange(comp)
}

func RangeContains(r Range, q Rational) bool {
	norm, err := NormalizeRange(r)
	if err != nil {
		return false
	}
	if norm.Kind() == RangeFull {
		return true
	}

	if norm.Inside {
		return atOrAboveLower(norm.Lo, q) && atOrBelowUpper(norm.Hi, q)
	}

	return atOrAboveLower(norm.Lo, q) || atOrBelowUpper(norm.Hi, q)
}

func atOrAboveLower(lo Endpoint, q Rational) bool {
	cmp := q.Cmp(lo.Value)
	if lo.Open {
		return cmp > 0
	}
	return cmp >= 0
}

func atOrBelowUpper(hi Endpoint, q Rational) bool {
	cmp := q.Cmp(hi.Value)
	if hi.Open {
		return cmp < 0
	}
	return cmp <= 0
}

func exactArcRange(v Rational) Range {
	return Range{
		Lo: Endpoint{
			Value: v,
			Open:  false,
		},
		Hi: Endpoint{
			Value: v,
			Open:  false,
		},
		Inside: true,
		Kind_:  RangeArc,
	}
}

func canonicalFullRange() Range {
	zero := RationalFromInt64(0)
	return Range{
		Lo: Endpoint{
			Value: zero,
			Open:  false,
		},
		Hi: Endpoint{
			Value: zero,
			Open:  false,
		},
		Inside: true,
		Kind_:  RangeFull,
	}
}

// core/range_projective.go v3
