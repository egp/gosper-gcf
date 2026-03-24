// core/gcf_unary_streaming_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

type countingPQStream struct {
	steps []FinitePQStep
	calls *int
}

func (s *countingPQStream) NextPQ() (PQTerm, PQStream, Status) {
	if s.calls == nil {
		panic("countingPQStream calls counter is nil")
	}
	*s.calls++

	if len(s.steps) == 0 {
		return PQTerm{
			P: big.NewInt(0),
			Q: big.NewInt(0),
		}, s, StatusEOF
	}

	head := cloneFinitePQStep(s.steps[0])
	tail := &countingPQStream{
		steps: append([]FinitePQStep(nil), s.steps[1:]...),
		calls: s.calls,
	}

	return head.Term, tail, StatusOK
}

func (s *countingPQStream) Range() Range {
	if len(s.steps) == 0 {
		panic("Range() is undefined on EOF PQStream")
	}
	return cloneRange(s.steps[0].Range)
}

func TestWB_GCF_UnaryConstructorDoesNotConsumeSource(t *testing.T) {
	calls := 0

	stream := &countingPQStream{
		steps: []FinitePQStep{
			{
				Term: PQTerm{P: big.NewInt(3), Q: big.NewInt(1)},
				Range: Range{
					Lo: Endpoint{
						Value: RationalFromInt64(3),
						Open:  false,
					},
					Hi: Endpoint{
						Value: RationalFromInt64(4),
						Open:  false,
					},
					Inside: true,
				},
			},
			{
				Term: PQTerm{P: big.NewInt(1), Q: big.NewInt(1)},
				Range: Range{
					Lo: Endpoint{
						Value: RationalFromInt64(1),
						Open:  false,
					},
					Hi: Endpoint{
						Value: RationalFromInt64(2),
						Open:  false,
					},
					Inside: true,
				},
			},
		},
		calls: &calls,
	}

	g := NewGCF1(identityUnaryXCoeffs(), stream)
	if g == nil {
		t.Fatal("NewGCF1 returned nil")
	}

	if calls != 0 {
		t.Fatalf("constructor consumed source %d times, want 0", calls)
	}
}

func TestWB_GCF_UnaryLiveRangeComesFromCurrentEvaluatorState(t *testing.T) {
	calls := 0

	stream := &countingPQStream{
		steps: []FinitePQStep{
			{
				Term: PQTerm{P: big.NewInt(3), Q: big.NewInt(1)},
				Range: Range{
					Lo: Endpoint{
						Value: RationalFromInt64(3),
						Open:  false,
					},
					Hi: Endpoint{
						Value: RationalFromInt64(4),
						Open:  false,
					},
					Inside: true,
				},
			},
			{
				Term: PQTerm{P: big.NewInt(1), Q: big.NewInt(1)},
				Range: Range{
					Lo: Endpoint{
						Value: RationalFromInt64(1),
						Open:  false,
					},
					Hi: Endpoint{
						Value: RationalFromInt64(2),
						Open:  false,
					},
					Inside: true,
				},
			},
		},
		calls: &calls,
	}

	g := NewGCF1(identityUnaryXCoeffs(), stream)

	r := g.Range()

	if !r.Inside {
		t.Fatal("Range.Inside = false, want true")
	}
	if r.Lo.Value.Cmp(RationalFromInt64(3)) != 0 {
		t.Fatalf("Lo = %v/%v, want 3/1", r.Lo.Value.Num(), r.Lo.Value.Den())
	}
	if r.Hi.Value.Cmp(RationalFromInt64(4)) != 0 {
		t.Fatalf("Hi = %v/%v, want 4/1", r.Hi.Value.Num(), r.Hi.Value.Den())
	}
}

func identityUnaryXCoeffs() TransformCoefficients {
	return TransformCoefficients{
		A: big.NewInt(0),
		B: big.NewInt(1),
		C: big.NewInt(0),
		D: big.NewInt(0),
		E: big.NewInt(0),
		F: big.NewInt(0),
		G: big.NewInt(0),
		H: big.NewInt(1),
	}
}

// core/gcf_unary_streaming_wb_test.go v1
