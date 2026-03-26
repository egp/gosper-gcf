// core/sqrt_cycle5_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

func TestWB_SqrtApproximationPQStream_SwitchesToOurorobosAfterFeedback(t *testing.T) {
	x := &sqrtControllerCountingPQStream{
		rng: exactRangeFromRational(RationalFromInt64(2)),
	}

	c := newSqrtController(x)
	proxy := newSqrtApproximationPQStream(c)

	feedbackRange := exactRangeFromRational(NewRational(big.NewInt(3), big.NewInt(2)))
	c.feedCertifiedTerm(NewRCFTerm(big.NewInt(1)), feedbackRange)

	got := proxy.Range()
	assertSqrtControllerExactRange(t, got, feedbackRange)
}

func TestWB_SqrtObservedRefinementStream_BuildsInnerFromApproximationProxy(t *testing.T) {
	x := &sqrtControllerCountingPQStream{
		rng: exactRangeFromRational(RationalFromInt64(2)),
	}

	c := newSqrtController(x)
	observed := newSqrtObservedRefinementStream(c)

	if observed.inner == nil {
		t.Fatalf("inner = nil, want *GCF")
	}
	if observed.inner.y != c.half {
		t.Fatalf("top-level right operand = %T, want controller half source", observed.inner.y)
	}

	leftAdapter, ok := observed.inner.x.(*rcfAsPQStream)
	if !ok {
		t.Fatalf("top-level left operand type = %T, want *rcfAsPQStream", observed.inner.x)
	}

	addNode, ok := leftAdapter.src.(*GCF)
	if !ok {
		t.Fatalf("adapter src type = %T, want *GCF add node", leftAdapter.src)
	}

	if _, ok := addNode.x.(*sqrtApproximationPQStream); !ok {
		t.Fatalf("add left operand type = %T, want *sqrtApproximationPQStream", addNode.x)
	}

	rightAdapter, ok := addNode.y.(*rcfAsPQStream)
	if !ok {
		t.Fatalf("add right operand type = %T, want *rcfAsPQStream", addNode.y)
	}

	divNode, ok := rightAdapter.src.(*GCF)
	if !ok {
		t.Fatalf("right adapter src type = %T, want *GCF div node", rightAdapter.src)
	}

	if divNode.x != x {
		t.Fatalf("div left operand != x radicand")
	}
	if _, ok := divNode.y.(*sqrtApproximationPQStream); !ok {
		t.Fatalf("div right operand type = %T, want *sqrtApproximationPQStream", divNode.y)
	}
}

func TestWB_SqrtObservedRefinementStream_NextRCF_FeedsBackEmittedTerm(t *testing.T) {
	x := &sqrtControllerCountingPQStream{
		rng: exactRangeFromRational(RationalFromInt64(2)),
	}

	c := newSqrtController(x)
	observed := newSqrtObservedRefinementStream(c)

	_, status := observed.NextRCF()
	if status != StatusOK {
		t.Fatalf("first status = %v, want %v", status, StatusOK)
	}

	if !c.hasOurorobosFeedback {
		t.Fatalf("hasOurorobosFeedback = false, want true after emitted term")
	}
}

// core/sqrt_cycle5_wb_test.go v1
