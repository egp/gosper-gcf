// core/sub_div_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

func TestWB_GCF_BLFTCoefficients_XMinusY(t *testing.T) {
	x, y := wbSubDivStreams(t)

	g := NewGCF2(
		BLFTCoefficients{
			A: big.NewInt(0),
			B: big.NewInt(1),
			C: big.NewInt(-1),
			D: big.NewInt(0),
			E: big.NewInt(0),
			F: big.NewInt(0),
			G: big.NewInt(0),
			H: big.NewInt(1),
		},
		x,
		y,
	)

	assertWBRCFSequence(t, g, []int64{-1, 2})
}

func TestWB_GCF_BLFTCoefficients_YMinusX(t *testing.T) {
	x, y := wbSubDivStreams(t)

	g := NewGCF2(
		BLFTCoefficients{
			A: big.NewInt(0),
			B: big.NewInt(-1),
			C: big.NewInt(1),
			D: big.NewInt(0),
			E: big.NewInt(0),
			F: big.NewInt(0),
			G: big.NewInt(0),
			H: big.NewInt(1),
		},
		x,
		y,
	)

	assertWBRCFSequence(t, g, []int64{0, 2})
}

func TestWB_GCF_BLFTCoefficients_XDivY(t *testing.T) {
	x, y := wbSubDivStreams(t)

	g := NewGCF2(
		BLFTCoefficients{
			A: big.NewInt(0),
			B: big.NewInt(1),
			C: big.NewInt(0),
			D: big.NewInt(0),
			E: big.NewInt(0),
			F: big.NewInt(0),
			G: big.NewInt(1),
			H: big.NewInt(0),
		},
		x,
		y,
	)

	assertWBRCFSequence(t, g, []int64{0, 1, 3})
}

func TestWB_GCF_BLFTCoefficients_YDivX(t *testing.T) {
	x, y := wbSubDivStreams(t)

	g := NewGCF2(
		BLFTCoefficients{
			A: big.NewInt(0),
			B: big.NewInt(0),
			C: big.NewInt(1),
			D: big.NewInt(0),
			E: big.NewInt(0),
			F: big.NewInt(1),
			G: big.NewInt(0),
			H: big.NewInt(0),
		},
		x,
		y,
	)

	assertWBRCFSequence(t, g, []int64{1, 3})
}

func wbSubDivStreams(t *testing.T) (PQStream, PQStream) {
	t.Helper()

	x, xStatus := NewFinitePQStream([]FinitePQStep{
		{
			Term:  PQTerm{P: big.NewInt(1), Q: big.NewInt(1)},
			Range: wbSubDivExactRange(3, 2),
		},
		{
			Term:  PQTerm{P: big.NewInt(2), Q: big.NewInt(1)},
			Range: wbSubDivExactRange(2, 1),
		},
	})
	if xStatus != StatusOK {
		t.Fatalf("left stream status = %v, want %v", xStatus, StatusOK)
	}

	y, yStatus := NewFinitePQStream([]FinitePQStep{
		{
			Term:  PQTerm{P: big.NewInt(2), Q: big.NewInt(1)},
			Range: wbSubDivExactRange(2, 1),
		},
	})
	if yStatus != StatusOK {
		t.Fatalf("right stream status = %v, want %v", yStatus, StatusOK)
	}

	return x, y
}

func wbSubDivExactRange(num, den int64) Range {
	value := NewRational(big.NewInt(num), big.NewInt(den))
	return Range{
		Lo: Endpoint{
			Value: value,
			Open:  false,
		},
		Hi: Endpoint{
			Value: value,
			Open:  false,
		},
		Inside: true,
	}
}

func assertWBRCFSequence(t *testing.T, g *GCF, want []int64) {
	t.Helper()

	for i, w := range want {
		got, status := g.NextRCF()
		if status != StatusOK {
			t.Fatalf("term %d status = %v, want %v", i+1, status, StatusOK)
		}
		if got.A().Cmp(big.NewInt(w)) != 0 {
			t.Fatalf("term %d = %v, want %d", i+1, got.A(), w)
		}
	}

	_, eofStatus := g.NextRCF()
	if eofStatus != StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, StatusEOF)
	}
}

// core/sub_div_wb_test.go v1
