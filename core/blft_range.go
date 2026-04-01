// core/blft_range.go v9
package core

import "fmt"

func (s blftState) CornerRange(xr, yr Range) (Range, error) {
	if special, ok := s.specialCaseRange(xr, yr); ok {
		return special, nil
	}

	if !xr.Inside || !yr.Inside {
		return Range{}, fmt.Errorf(
			"CornerRange: %w\nBLFT=%s\nxRange=%s\nyRange=%s",
			ErrUnsupportedRangeCase,
			formatBLFTDebug(s),
			formatRangeDebug(xr),
			formatRangeDebug(yr),
		)
	}

	type cornerValue struct {
		value Rational
		open  bool
	}

	corners := []struct {
		x    Rational
		y    Rational
		open bool
	}{
		{xr.Lo.Value, yr.Lo.Value, xr.Lo.Open || yr.Lo.Open},
		{xr.Lo.Value, yr.Hi.Value, xr.Lo.Open || yr.Hi.Open},
		{xr.Hi.Value, yr.Lo.Value, xr.Hi.Open || yr.Lo.Open},
		{xr.Hi.Value, yr.Hi.Value, xr.Hi.Open || yr.Hi.Open},
	}

	values := make([]cornerValue, 0, 4)
	sawZeroDen := false
	sawPosDen := false
	sawNegDen := false

	for _, corner := range corners {
		num, den := evalBLFTNumDenAt(s, corner.x, corner.y)
		switch den.Sign() {
		case 0:
			sawZeroDen = true
			continue
		case 1:
			sawPosDen = true
		case -1:
			sawNegDen = true
		}

		values = append(values, cornerValue{
			value: NewRational(num, den),
			open:  corner.open,
		})
	}

	if sawZeroDen || (sawPosDen && sawNegDen) {
		raw := make([]Rational, 0, len(values))
		for _, v := range values {
			raw = append(raw, v.value)
		}
		return outsideRangeFromValues(raw), nil
	}

	if len(values) == 0 {
		return outsideRangeFromValues(nil), nil
	}

	lo := values[0].value
	hi := values[0].value
	loOpen := values[0].open
	hiOpen := values[0].open

	for _, v := range values[1:] {
		switch cmp := v.value.Cmp(lo); {
		case cmp < 0:
			lo = v.value
			loOpen = v.open
		case cmp == 0:
			loOpen = loOpen && v.open
		}

		switch cmp := v.value.Cmp(hi); {
		case cmp > 0:
			hi = v.value
			hiOpen = v.open
		case cmp == 0:
			hiOpen = hiOpen && v.open
		}
	}

	return Range{
		Lo: Endpoint{
			Value: lo,
			Open:  loOpen,
		},
		Hi: Endpoint{
			Value: hi,
			Open:  hiOpen,
		},
		Inside: true,
	}, nil
}

func (s blftState) specialCaseRange(xr, yr Range) (Range, bool) {
	if r, ok := s.constantRange(); ok {
		return r, true
	}

	if r, ok := s.rangeWithExactX(xr, yr); ok {
		return r, true
	}

	if r, ok := s.rangeWithExactY(xr, yr); ok {
		return r, true
	}

	if r, ok := s.rangeWithDegenerateClosedX(xr, yr); ok {
		return r, true
	}

	if r, ok := s.rangeWithDegenerateClosedY(xr, yr); ok {
		return r, true
	}

	if r, ok := s.affineXRange(xr); ok {
		return r, true
	}

	if r, ok := s.affineYRange(yr); ok {
		return r, true
	}

	return Range{}, false
}

// core/blft_range.go v9
