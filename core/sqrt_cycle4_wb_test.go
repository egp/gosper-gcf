// core/sqrt_cycle4_wb_test.go v2
package core

import (
	"math/big"
	"testing"
)

func TestWB_SqrtController_ActiveApproximation_UsesSeedBeforeFeedback(t *testing.T) {
	x := &sqrtControllerCountingPQStream{
		rng: exactRangeFromRational(RationalFromInt64(2)),
	}

	c := newSqrtController(x)
	if c.activeApproximation() != c.seedApproximation() {
		t.Fatalf("active approximation before feedback != seed approximation")
	}
}

func TestWB_SqrtController_ActiveApproximation_SwitchesToOurorobosAfterFeedback(t *testing.T) {
	x := &sqrtControllerCountingPQStream{
		rng: exactRangeFromRational(RationalFromInt64(2)),
	}

	c := newSqrtController(x)
	if err := c.feedCertifiedTerm(
		NewRCFTerm(big.NewInt(1)),
		exactRangeFromRational(NewRational(big.NewInt(3), big.NewInt(2))),
	); err != nil {
		t.Fatalf("feedCertifiedTerm error = %v", err)
	}
	if c.activeApproximation() != c.ourorobosApproximation() {
		t.Fatalf("active approximation after feedback != ourorobos approximation")
	}
}

func TestWB_SqrtController_CurrentRefinement_UsesOurorobosAfterFeedback(t *testing.T) {
	x := &sqrtControllerCountingPQStream{
		rng: exactRangeFromRational(RationalFromInt64(2)),
	}

	c := newSqrtController(x)
	if err := c.feedCertifiedTerm(
		NewRCFTerm(big.NewInt(1)),
		exactRangeFromRational(NewRational(big.NewInt(3), big.NewInt(2))),
	); err != nil {
		t.Fatalf("feedCertifiedTerm error = %v", err)
	}

	got := c.currentRefinement()
	if got == nil {
		t.Fatalf("current refinement = nil, want *GCF")
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
	if addNode.x != c.ourorobosApproximation() {
		t.Fatalf("add left operand != ourorobos approximation after feedback")
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
	if divNode.y != c.ourorobosApproximation() {
		t.Fatalf("div right operand != ourorobos approximation after feedback")
	}
}

// core/sqrt_cycle4_wb_test.go v2
