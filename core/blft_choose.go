// core/blft_choose.go v5
package core

func chooseIngestX(xRange, yRange Range) bool {
	if xRange.Inside && yRange.Inside {
		xWidth := rangeWidth(xRange)
		yWidth := rangeWidth(yRange)
		return xWidth.Cmp(yWidth) >= 0
	}

	if isExactClosedRangeChoose(xRange) {
		return true
	}
	if isExactClosedRangeChoose(yRange) {
		return false
	}

	if xRange.Inside != yRange.Inside {
		return xRange.Inside
	}

	xWidth := rangeWidth(xRange)
	yWidth := rangeWidth(yRange)
	return xWidth.Cmp(yWidth) >= 0
}

func isExactClosedRangeChoose(r Range) bool {
	return r.Inside &&
		!r.Lo.Open &&
		!r.Hi.Open &&
		r.Lo.Value.Cmp(r.Hi.Value) == 0
}

// core/blft_choose.go v5
