// core/dlft_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

func TestWB_DLFT_IngestXMatchesFormula(t *testing.T) {
	s := dlftState{
		A: big.NewInt(1),
		B: big.NewInt(2),
		C: big.NewInt(3),
		D: big.NewInt(4),
		E: big.NewInt(5),
		F: big.NewInt(6),
	}

	got := s.IngestX(PQTerm{
		P: big.NewInt(3),
		Q: big.NewInt(5),
	})

	want := dlftState{
		A: big.NewInt(18),  // a p^2 + b p + c = 1*9 + 2*3 + 3
		B: big.NewInt(40),  // 2 a p q + b q   = 2*1*3*5 + 2*5
		C: big.NewInt(25),  // a q^2           = 1*25
		D: big.NewInt(57),  // d p^2 + e p + f = 4*9 + 5*3 + 6
		E: big.NewInt(145), // 2 d p q + e q   = 2*4*3*5 + 5*5
		F: big.NewInt(100), // d q^2           = 4*25
	}

	assertDLFTStateEqual(t, got, want)
}

func TestWB_DLFT_EmitFormulaMatchesGosper(t *testing.T) {
	s := dlftState{
		A: big.NewInt(1),
		B: big.NewInt(2),
		C: big.NewInt(3),
		D: big.NewInt(4),
		E: big.NewInt(5),
		F: big.NewInt(6),
	}

	got := s.Emit(NewRCFTerm(big.NewInt(2)))

	want := dlftState{
		A: big.NewInt(4),
		B: big.NewInt(5),
		C: big.NewInt(6),
		D: big.NewInt(-7), // 1 - 2*4
		E: big.NewInt(-8), // 2 - 2*5
		F: big.NewInt(-9), // 3 - 2*6
	}

	assertDLFTStateEqual(t, got, want)
}

func TestWB_DLFT_CollapseChoosesHighestSurvivingDegree(t *testing.T) {
	t.Run("quadratic", func(t *testing.T) {
		s := dlftState{
			A: big.NewInt(6),
			B: big.NewInt(2),
			C: big.NewInt(3),
			D: big.NewInt(8),
			E: big.NewInt(5),
			F: big.NewInt(7),
		}

		got := s.CollapseToRational()
		want := NewRational(big.NewInt(6), big.NewInt(8))

		if got.Cmp(want) != 0 {
			t.Fatalf("CollapseToRational() = %v/%v, want %v/%v",
				got.Num(), got.Den(), want.Num(), want.Den())
		}
	})

	t.Run("linear fallback", func(t *testing.T) {
		s := dlftState{
			A: big.NewInt(0),
			B: big.NewInt(9),
			C: big.NewInt(3),
			D: big.NewInt(0),
			E: big.NewInt(12),
			F: big.NewInt(7),
		}

		got := s.CollapseToRational()
		want := NewRational(big.NewInt(9), big.NewInt(12))

		if got.Cmp(want) != 0 {
			t.Fatalf("CollapseToRational() = %v/%v, want %v/%v",
				got.Num(), got.Den(), want.Num(), want.Den())
		}
	})

	t.Run("constant fallback", func(t *testing.T) {
		s := dlftState{
			A: big.NewInt(0),
			B: big.NewInt(0),
			C: big.NewInt(10),
			D: big.NewInt(0),
			E: big.NewInt(0),
			F: big.NewInt(15),
		}

		got := s.CollapseToRational()
		want := NewRational(big.NewInt(10), big.NewInt(15))

		if got.Cmp(want) != 0 {
			t.Fatalf("CollapseToRational() = %v/%v, want %v/%v",
				got.Num(), got.Den(), want.Num(), want.Den())
		}
	})
}

func TestWB_DLFT_ConstantCandidateRangeIsExact(t *testing.T) {
	s := dlftState{
		A: big.NewInt(0),
		B: big.NewInt(0),
		C: big.NewInt(3),
		D: big.NewInt(0),
		E: big.NewInt(0),
		F: big.NewInt(2),
	}

	xRange := Range{
		Lo: Endpoint{
			Value: RationalFromInt64(2),
			Open:  false,
		},
		Hi: Endpoint{
			Value: RationalFromInt64(5),
			Open:  false,
		},
		Inside: true,
	}

	got := s.CandidateRange(xRange)
	want := NewRational(big.NewInt(3), big.NewInt(2))

	if !got.Inside {
		t.Fatal("CandidateRange().Inside = false, want true")
	}
	if got.Lo.Value.Cmp(want) != 0 {
		t.Fatalf("Lo = %v/%v, want %v/%v",
			got.Lo.Value.Num(), got.Lo.Value.Den(), want.Num(), want.Den())
	}
	if got.Hi.Value.Cmp(want) != 0 {
		t.Fatalf("Hi = %v/%v, want %v/%v",
			got.Hi.Value.Num(), got.Hi.Value.Den(), want.Num(), want.Den())
	}
}

func assertDLFTStateEqual(t *testing.T, got, want dlftState) {
	t.Helper()

	assertBigIntEqual(t, "A", got.A, want.A)
	assertBigIntEqual(t, "B", got.B, want.B)
	assertBigIntEqual(t, "C", got.C, want.C)
	assertBigIntEqual(t, "D", got.D, want.D)
	assertBigIntEqual(t, "E", got.E, want.E)
	assertBigIntEqual(t, "F", got.F, want.F)
}

// core/dlft_wb_test.go v1
