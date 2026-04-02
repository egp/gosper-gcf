// core/gcf_unary_range_openness_wb_test.go v2
package core

import (
	"math/big"
	"testing"
)

type staticRangePQStream struct {
	rng Range
}

func (s *staticRangePQStream) NextPQ() (PQTerm, PQStream, Status, error) {
	return PQTerm{
		P: big.NewInt(0),
		Q: big.NewInt(0),
	}, s, StatusEOF, nil
}

func (s *staticRangePQStream) Range() (Range, error) {
	return s.rng, nil
}

func TestWB_GCF_UnaryRange_Identity_PreservesHalfOpenRange(t *testing.T) {
	src := &staticRangePQStream{
		rng: Range{
			Lo: Endpoint{
				Value: RationalFromInt64(2),
				Open:  false,
			},
			Hi: Endpoint{
				Value: RationalFromInt64(5),
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

	got, err := g.Range()
	if err != nil {
		t.Fatalf("Range error = %v", err)
	}

	if !got.Inside {
		t.Fatal("Inside = false, want true")
	}
	if got.Lo.Value.Cmp(RationalFromInt64(2)) != 0 {
		t.Fatalf("Lo = %v/%v, want 2/1", got.Lo.Value.Num(), got.Lo.Value.Den())
	}
	if got.Hi.Value.Cmp(RationalFromInt64(5)) != 0 {
		t.Fatalf("Hi = %v/%v, want 5/1", got.Hi.Value.Num(), got.Hi.Value.Den())
	}
	if got.Lo.Open {
		t.Fatal("Lo.Open = true, want false")
	}
	if !got.Hi.Open {
		t.Fatal("Hi.Open = false, want true")
	}
}

func TestWB_GCF_UnaryRange_ScaleHalf_PreservesOpennessAndScalesEndpoints(t *testing.T) {
	src := &staticRangePQStream{
		rng: Range{
			Lo: Endpoint{
				Value: RationalFromInt64(2),
				Open:  true,
			},
			Hi: Endpoint{
				Value: RationalFromInt64(6),
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
			H: big.NewInt(2),
		},
		src,
	)

	got, err := g.Range()
	if err != nil {
		t.Fatalf("Range error = %v", err)
	}

	if !got.Inside {
		t.Fatal("Inside = false, want true")
	}
	wantLo := RationalFromInt64(1)
	wantHi := RationalFromInt64(3)
	if got.Lo.Value.Cmp(wantLo) != 0 {
		t.Fatalf("Lo = %v/%v, want 1/1", got.Lo.Value.Num(), got.Lo.Value.Den())
	}
	if got.Hi.Value.Cmp(wantHi) != 0 {
		t.Fatalf("Hi = %v/%v, want 3/1", got.Hi.Value.Num(), got.Hi.Value.Den())
	}
	if !got.Lo.Open {
		t.Fatal("Lo.Open = false, want true")
	}
	if got.Hi.Open {
		t.Fatal("Hi.Open = true, want false")
	}
}

func TestWB_GCF_UnaryRange_Constant_ReturnsExactClosedRange(t *testing.T) {
	src := &staticRangePQStream{
		rng: Range{
			Lo: Endpoint{
				Value: RationalFromInt64(2),
				Open:  true,
			},
			Hi: Endpoint{
				Value: RationalFromInt64(6),
				Open:  true,
			},
			Inside: true,
		},
	}

	g := NewGCF1(
		BLFTCoefficients{
			A: big.NewInt(0),
			B: big.NewInt(0),
			C: big.NewInt(0),
			D: big.NewInt(7),
			E: big.NewInt(0),
			F: big.NewInt(0),
			G: big.NewInt(0),
			H: big.NewInt(1),
		},
		src,
	)

	got, err := g.Range()
	if err != nil {
		t.Fatalf("Range error = %v", err)
	}

	if !got.Inside {
		t.Fatal("Inside = false, want true")
	}
	if got.Lo.Value.Cmp(RationalFromInt64(7)) != 0 || got.Hi.Value.Cmp(RationalFromInt64(7)) != 0 {
		t.Fatalf("Range = [%v/%v,%v/%v], want exact 7/1",
			got.Lo.Value.Num(), got.Lo.Value.Den(),
			got.Hi.Value.Num(), got.Hi.Value.Den(),
		)
	}
	if got.Lo.Open || got.Hi.Open {
		t.Fatalf("Openness = (%v,%v), want (false,false)", got.Lo.Open, got.Hi.Open)
	}
}

func TestWB_GCF_UnaryRange_RangeDoesNotConsumeChildStream(t *testing.T) {
	src := &staticRangePQStream{
		rng: Range{
			Lo: Endpoint{
				Value: RationalFromInt64(10),
				Open:  false,
			},
			Hi: Endpoint{
				Value: RationalFromInt64(11),
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

	r1, err := g.Range()
	if err != nil {
		t.Fatalf("first Range error = %v", err)
	}
	r2, err := g.Range()
	if err != nil {
		t.Fatalf("second Range error = %v", err)
	}

	if r1.Lo.Value.Cmp(r2.Lo.Value) != 0 || r1.Hi.Value.Cmp(r2.Hi.Value) != 0 {
		t.Fatalf("Range changed across repeated calls: first=[%v/%v,%v/%v] second=[%v/%v,%v/%v]",
			r1.Lo.Value.Num(), r1.Lo.Value.Den(),
			r1.Hi.Value.Num(), r1.Hi.Value.Den(),
			r2.Lo.Value.Num(), r2.Lo.Value.Den(),
			r2.Hi.Value.Num(), r2.Hi.Value.Den(),
		)
	}
	if r1.Lo.Open != r2.Lo.Open || r1.Hi.Open != r2.Hi.Open || r1.Inside != r2.Inside {
		t.Fatalf("Range openness/inside changed across repeated calls")
	}
}

// core/gcf_unary_range_openness_wb_test.go v2
