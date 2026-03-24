// core/blft_normalize_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

func TestWB_BLFT_NormalizesAfterMutation(t *testing.T) {
	t.Run("IngestX", func(t *testing.T) {
		s := blftState{
			A: big.NewInt(-2),
			B: big.NewInt(-4),
			C: big.NewInt(-6),
			D: big.NewInt(-8),
			E: big.NewInt(-10),
			F: big.NewInt(-12),
			G: big.NewInt(-14),
			H: big.NewInt(-16),
		}

		got := s.IngestX(PQTerm{
			P: big.NewInt(1),
			Q: big.NewInt(1),
		})

		want := blftState{
			A: big.NewInt(4),
			B: big.NewInt(6),
			C: big.NewInt(1),
			D: big.NewInt(2),
			E: big.NewInt(12),
			F: big.NewInt(14),
			G: big.NewInt(5),
			H: big.NewInt(6),
		}

		assertBLFTStateEqual(t, got, want)
	})

	t.Run("IngestY", func(t *testing.T) {
		s := blftState{
			A: big.NewInt(-2),
			B: big.NewInt(-4),
			C: big.NewInt(-6),
			D: big.NewInt(-8),
			E: big.NewInt(-10),
			F: big.NewInt(-12),
			G: big.NewInt(-14),
			H: big.NewInt(-16),
		}

		got := s.IngestY(PQTerm{
			P: big.NewInt(1),
			Q: big.NewInt(1),
		})

		want := blftState{
			A: big.NewInt(3),
			B: big.NewInt(1),
			C: big.NewInt(7),
			D: big.NewInt(3),
			E: big.NewInt(11),
			F: big.NewInt(5),
			G: big.NewInt(15),
			H: big.NewInt(7),
		}

		assertBLFTStateEqual(t, got, want)
	})
}

func TestWB_BLFT_NormalizeHelperCanonicalizes(t *testing.T) {
	s := blftState{
		A: big.NewInt(-8),
		B: big.NewInt(-12),
		C: big.NewInt(-2),
		D: big.NewInt(-4),
		E: big.NewInt(-24),
		F: big.NewInt(-28),
		G: big.NewInt(-10),
		H: big.NewInt(-12),
	}

	got := s.Normalize()

	want := blftState{
		A: big.NewInt(4),
		B: big.NewInt(6),
		C: big.NewInt(1),
		D: big.NewInt(2),
		E: big.NewInt(12),
		F: big.NewInt(14),
		G: big.NewInt(5),
		H: big.NewInt(6),
	}

	assertBLFTStateEqual(t, got, want)
}

// core/blft_normalize_wb_test.go v1
