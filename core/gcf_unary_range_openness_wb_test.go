// core/gcf_unary_range_openness_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

type staticRangePQStream struct {
	rng Range
}

func (s *staticRangePQStream) NextPQ() (PQTerm, PQStream, Status) {
	return PQTerm{
		P: big.NewInt(0),
		Q: big.NewInt(0),
	}, s, StatusEOF
}

func (s *staticRangePQStream) Range() Range {
	return s.rng
}

func TestWB_GCF_UnaryRange_IdentityPreservesEndpointOpenness(t *testing.T) {
	src := &staticRangePQStream{
		rng: Range{
			Lo: Endpoint{
				Value: NewRational(big.NewInt(5), big.NewInt(2)),
				Open:  false,
			},
			Hi: Endpoint{
				Value: RationalFromInt64(3),
				Open:  true,
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

	got := g.Range()

	if !got.Inside {
		t.Fatalf("Inside = false, want true")
	}
	if got.Lo.Open {
		t.Fatalf("Lo.Open = true, want false")
	}
	if !got.Hi.Open {
		t.Fatalf("Hi.Open = false, want true")
	}
	if got.Lo.Value.Cmp(NewRational(big.NewInt(5), big.NewInt(2))) != 0 {
		t.Fatalf("Lo = %v/%v, want 5/2", got.Lo.Value.Num(), got.Lo.Value.Den())
	}
	if got.Hi.Value.Cmp(RationalFromInt64(3)) != 0 {
		t.Fatalf("Hi = %v/%v, want 3/1", got.Hi.Value.Num(), got.Hi.Value.Den())
	}
}

// core/gcf_unary_range_openness_wb_test.go v1
