// core/blft_emit_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

func TestWB_BLFT_EmitFormulaMatchesGosper(t *testing.T) {
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

	got := s.Emit(NewRCFTerm(big.NewInt(2)))

	want := blftState{
		A: big.NewInt(5),
		B: big.NewInt(6),
		C: big.NewInt(7),
		D: big.NewInt(8),
		E: big.NewInt(-9),  // 1 - 2*5
		F: big.NewInt(-10), // 2 - 2*6
		G: big.NewInt(-11), // 3 - 2*7
		H: big.NewInt(-12), // 4 - 2*8
	}

	assertBLFTStateEqual(t, got, want)
}

func TestWB_BLFT_InsideIntervalCanEmitWhenFloorsAgree(t *testing.T) {
	s := blftState{
		A: big.NewInt(1),
		B: big.NewInt(0),
		C: big.NewInt(0),
		D: big.NewInt(0),
		E: big.NewInt(0),
		F: big.NewInt(0),
		G: big.NewInt(0),
		H: big.NewInt(1),
	}

	r := Range{
		Lo: Endpoint{
			Value: NewRational(big.NewInt(7), big.NewInt(3)), // 2.333...
			Open:  false,
		},
		Hi: Endpoint{
			Value: NewRational(big.NewInt(8), big.NewInt(3)), // 2.666...
			Open:  false,
		},
		Inside: true,
	}

	got, ok := s.CanEmitRCFTerm(r)
	if !ok {
		t.Fatal("CanEmitRCFTerm returned ok=false, want true")
	}
	if got.A().Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("emitted term = %v, want 2", got.A())
	}
}

func TestWB_BLFT_OutsideIntervalDoesNotDirectlyEmit(t *testing.T) {
	s := blftState{
		A: big.NewInt(1),
		B: big.NewInt(0),
		C: big.NewInt(0),
		D: big.NewInt(0),
		E: big.NewInt(0),
		F: big.NewInt(0),
		G: big.NewInt(0),
		H: big.NewInt(1),
	}

	r := Range{
		Lo: Endpoint{
			Value: RationalFromInt64(5),
			Open:  false,
		},
		Hi: Endpoint{
			Value: RationalFromInt64(2),
			Open:  false,
		},
		Inside: false,
	}

	got, ok := s.CanEmitRCFTerm(r)
	if ok {
		t.Fatalf("CanEmitRCFTerm returned ok=true with term %v, want false", got.A())
	}
}

// core/blft_emit_wb_test.go v1
