// core/blft_collapse_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

func TestWB_BLFT_CollapseOnXEOF(t *testing.T) {
	s := blftState{
		A: big.NewInt(1),
		B: big.NewInt(2),
		C: big.NewInt(3),
		D: big.NewInt(4),
		E: big.NewInt(5),
		F: big.NewInt(6),
		G: big.NewInt(7),
		H: big.NewInt(8),
	}

	got := s.CollapseX()

	want := TransformCoefficients{
		A: big.NewInt(0),
		B: big.NewInt(0),
		C: big.NewInt(3),
		D: big.NewInt(4),
		E: big.NewInt(0),
		F: big.NewInt(0),
		G: big.NewInt(7),
		H: big.NewInt(8),
	}

	assertTransformCoefficientsEqual(t, got, want)
}

func TestWB_BLFT_CollapseOnYEOF(t *testing.T) {
	s := blftState{
		A: big.NewInt(1),
		B: big.NewInt(2),
		C: big.NewInt(3),
		D: big.NewInt(4),
		E: big.NewInt(5),
		F: big.NewInt(6),
		G: big.NewInt(7),
		H: big.NewInt(8),
	}

	got := s.CollapseY()

	want := TransformCoefficients{
		A: big.NewInt(0),
		B: big.NewInt(2),
		C: big.NewInt(0),
		D: big.NewInt(4),
		E: big.NewInt(0),
		F: big.NewInt(6),
		G: big.NewInt(0),
		H: big.NewInt(8),
	}

	assertTransformCoefficientsEqual(t, got, want)
}

func TestWB_BLFT_CollapseOnBothEOFProducesExactState(t *testing.T) {
	s := blftState{
		A: big.NewInt(1),
		B: big.NewInt(2),
		C: big.NewInt(3),
		D: big.NewInt(4),
		E: big.NewInt(5),
		F: big.NewInt(6),
		G: big.NewInt(7),
		H: big.NewInt(8),
	}

	got := s.CollapseToRational()
	want := NewRational(big.NewInt(4), big.NewInt(8))

	if got.Cmp(want) != 0 {
		t.Fatalf("CollapseToRational() = %v/%v, want %v/%v",
			got.Num(), got.Den(), want.Num(), want.Den())
	}
}

func assertTransformCoefficientsEqual(t *testing.T, got, want TransformCoefficients) {
	t.Helper()

	assertBigIntEqual(t, "A", got.A, want.A)
	assertBigIntEqual(t, "B", got.B, want.B)
	assertBigIntEqual(t, "C", got.C, want.C)
	assertBigIntEqual(t, "D", got.D, want.D)
	assertBigIntEqual(t, "E", got.E, want.E)
	assertBigIntEqual(t, "F", got.F, want.F)
	assertBigIntEqual(t, "G", got.G, want.G)
	assertBigIntEqual(t, "H", got.H, want.H)
}

// core/blft_collapse_wb_test.go v1
