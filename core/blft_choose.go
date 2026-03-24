// core/blft_choose.go v2
package core

func chooseIngestX(xRange, yRange Range) bool {
	if !xRange.Inside || !yRange.Inside {
		panic("chooseIngestX currently supports only inside/inside ranges")
	}

	xWidth := rangeWidth(xRange)
	yWidth := rangeWidth(yRange)

	return xWidth.Cmp(yWidth) >= 0
}

// core/blft_choose.go v2
