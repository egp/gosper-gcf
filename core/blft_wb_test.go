// core/blft_wb_test.go v3
package core

import (
	"math/big"
	"testing"
)

func TestWB_BLFT_IngestXMatchesGosperFormula(t *testing.T) {
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

	got := s.IngestX(PQTerm{
		P: big.NewInt(3),
		Q: big.NewInt(5),
	})

	want := blftState{
		A: big.NewInt(6),  // A*p + C = 1*3 + 3
		B: big.NewInt(10), // B*p + D = 2*3 + 4
		C: big.NewInt(5),  // A*q     = 1*5
		D: big.NewInt(10), // B*q     = 2*5
		E: big.NewInt(22), // E*p + G = 5*3 + 7
		F: big.NewInt(26), // F*p + H = 6*3 + 8
		G: big.NewInt(25), // E*q     = 5*5
		H: big.NewInt(30), // F*q     = 6*5
	}

	assertBLFTStateEqual(t, got, want)
}

func TestWB_BLFT_IngestYMatchesGosperFormula(t *testing.T) {
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

	got := s.IngestY(PQTerm{
		P: big.NewInt(3),
		Q: big.NewInt(5),
	})

	want := blftState{
		A: big.NewInt(5),  // A*p + B = 1*3 + 2
		B: big.NewInt(5),  // A*q     = 1*5
		C: big.NewInt(13), // C*p + D = 3*3 + 4
		D: big.NewInt(15), // C*q     = 3*5
		E: big.NewInt(21), // E*p + F = 5*3 + 6
		F: big.NewInt(25), // E*q     = 5*5
		G: big.NewInt(29), // G*p + H = 7*3 + 8
		H: big.NewInt(35), // G*q     = 7*5
	}

	assertBLFTStateEqual(t, got, want)
}

func assertBLFTStateEqual(t *testing.T, got, want blftState) {
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

func assertBigIntEqual(t *testing.T, name string, got, want *big.Int) {
	t.Helper()

	if got == nil && want == nil {
		return
	}
	if got == nil || want == nil {
		t.Fatalf("%s = %v, want %v", name, got, want)
	}
	if got.Cmp(want) != 0 {
		t.Fatalf("%s = %v, want %v", name, got, want)
	}
}

// core/blft_wb_test.go v3
