// core/gcf_unary_streaming_wb_test.go v3
package core

import (
	"math/big"
	"testing"
)

type countingPQStream struct {
	rng        Range
	nextCalls  int
	rangeCalls int
}

func (s *countingPQStream) NextPQ() (PQTerm, PQStream, Status, error) {
	s.nextCalls++
	return PQTerm{
		P: big.NewInt(0),
		Q: big.NewInt(0),
	}, s, StatusEOF, nil
}

func (s *countingPQStream) CurrentInterval() (Interval, error) {
	s.rangeCalls++
	return s.rng, nil
}

func (s *countingPQStream) Range() (Range, error) {
	return s.CurrentInterval()
}

func TestWB_GCF_UnaryStreaming_ExactIntegerRangeEmitsWithoutIngest(t *testing.T) {
	src := &countingPQStream{
		rng: exactRangeFromRational(RationalFromInt64(3)),
	}

	g := NewGCF1(
		BLFTCoefficients{
			A: big.NewInt(0),
			B: big.NewInt(1),
			C: big.NewInt(0),
			D: big.NewInt(0),
			E: big.NewInt(0),
			F: big.NewInt(0),
			G: big.NewInt(0),
			H: big.NewInt(1),
		},
		src,
	)

	term, status, err := g.NextRCF()
	if err != nil {
		t.Fatalf("NextRCF error = %v", err)
	}
	if status != StatusOK {
		t.Fatalf("status = %v, want %v", status, StatusOK)
	}
	if term.A().Cmp(big.NewInt(3)) != 0 {
		t.Fatalf("term = %v, want 3", term.A())
	}
	if src.nextCalls != 0 {
		t.Fatalf("nextCalls = %d, want 0", src.nextCalls)
	}
}

func TestWB_GCF_UnaryStreaming_RangeQueryDoesNotIngest(t *testing.T) {
	src := &countingPQStream{
		rng: Range{
			Lo: Endpoint{
				Value: RationalFromInt64(2),
				Open:  false,
			},
			Hi: Endpoint{
				Value: RationalFromInt64(5),
				Open:  false,
			},
			Inside: true,
		},
	}

	g := NewGCF1(
		BLFTCoefficients{
			A: big.NewInt(0),
			B: big.NewInt(1),
			C: big.NewInt(0),
			D: big.NewInt(0),
			E: big.NewInt(0),
			F: big.NewInt(0),
			G: big.NewInt(0),
			H: big.NewInt(1),
		},
		src,
	)

	_, err := g.Range()
	if err != nil {
		t.Fatalf("Range error = %v", err)
	}
	if src.nextCalls != 0 {
		t.Fatalf("nextCalls = %d, want 0", src.nextCalls)
	}
	if src.rangeCalls == 0 {
		t.Fatal("rangeCalls = 0, want > 0")
	}
}

// core/gcf_unary_streaming_wb_test.go v3
