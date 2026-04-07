// core/sqrt_controller_wb_test.go v4
package core

import (
	"math/big"
	"testing"
)

type sqrtControllerCountingPQStream struct {
	nextCalls  int
	rangeCalls int
	rng        Range
}

func (s *sqrtControllerCountingPQStream) NextPQ() (PQTerm, PQStream, Status, error) {
	s.nextCalls++
	return PQTerm{
		P: big.NewInt(1),
		Q: big.NewInt(1),
	}, s, StatusOK, nil
}

func (s *sqrtControllerCountingPQStream) CurrentInterval() (Interval, error) {
	s.rangeCalls++
	return s.rng, nil
}

func (s *sqrtControllerCountingPQStream) Range() (Range, error) {
	return s.CurrentInterval()
}

func TestWB_SqrtController_BuildRefinementGraph_UsesNewtonStructure(t *testing.T) {
	x := &sqrtControllerCountingPQStream{
		rng: exactRangeFromRational(RationalFromInt64(2)),
	}

	c := newSqrtController(x)
	y := c.ourorobosApproximation()
	got := c.buildRefinement(y)

	if got == nil {
		t.Fatalf("refinement = nil, want *GCF")
	}
	if got.y != c.half {
		t.Fatalf("top-level right operand = %T, want controller half source", got.y)
	}

	leftAdapter, ok := got.x.(*rcfAsPQStream)
	if !ok {
		t.Fatalf("top-level left operand type = %T, want *rcfAsPQStream", got.x)
	}
	addNode, ok := leftAdapter.src.(*GCF)
	if !ok {
		t.Fatalf("adapter src type = %T, want *GCF add node", leftAdapter.src)
	}
	if addNode.x != y {
		t.Fatalf("add left operand != y ourorobos approximation")
	}

	rightAdapter, ok := addNode.y.(*rcfAsPQStream)
	if !ok {
		t.Fatalf("add right operand type = %T, want *rcfAsPQStream", addNode.y)
	}
	divNode, ok := rightAdapter.src.(*GCF)
	if !ok {
		t.Fatalf("right adapter src type = %T, want *GCF div node", rightAdapter.src)
	}
	if divNode.x != PQStream(x) {
		t.Fatalf("div left operand != x radicand")
	}
	if divNode.y != y {
		t.Fatalf("div right operand != y ourorobos approximation")
	}
}

func TestWB_SqrtController_HalfSource_IsExactOneHalf(t *testing.T) {
	x := &sqrtControllerCountingPQStream{
		rng: exactRangeFromRational(RationalFromInt64(2)),
	}

	c := newSqrtController(x)
	got, err := c.halfSource().Range()
	if err != nil {
		t.Fatalf("halfSource Range error = %v", err)
	}
	want := exactRangeFromRational(NewRational(big.NewInt(1), big.NewInt(2)))
	assertSqrtControllerExactRange(t, got, want)
}

func TestWB_SqrtController_SeedApproximation_DoesNotPrereadX(t *testing.T) {
	x := &sqrtControllerCountingPQStream{
		rng: exactRangeFromRational(RationalFromInt64(2)),
	}

	c := newSqrtController(x)
	seed := c.seedApproximation()

	_, _, _, _ = seed.NextPQ()

	if x.nextCalls != 0 {
		t.Fatalf("x.NextPQ calls = %d, want 0", x.nextCalls)
	}
}

func assertSqrtControllerExactRange(t *testing.T, got Range, want Range) {
	t.Helper()

	if got.Inside != want.Inside {
		t.Fatalf("Inside = %v, want %v", got.Inside, want.Inside)
	}
	if got.Lo.Open != want.Lo.Open || got.Hi.Open != want.Hi.Open {
		t.Fatalf(
			"openness = (%v,%v), want (%v,%v)",
			got.Lo.Open, got.Hi.Open, want.Lo.Open, want.Hi.Open,
		)
	}
	if got.Lo.Value.Cmp(want.Lo.Value) != 0 || got.Hi.Value.Cmp(want.Hi.Value) != 0 {
		t.Fatalf(
			"range = [%v/%v,%v/%v], want [%v/%v,%v/%v]",
			got.Lo.Value.Num(), got.Lo.Value.Den(),
			got.Hi.Value.Num(), got.Hi.Value.Den(),
			want.Lo.Value.Num(), want.Lo.Value.Den(),
			want.Hi.Value.Num(), want.Hi.Value.Den(),
		)
	}
}

// core/sqrt_controller_wb_test.go v4
