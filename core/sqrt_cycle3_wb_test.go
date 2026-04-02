// core/sqrt_cycle3_wb_test.go v2
package core

import (
	"math/big"
	"testing"
)

func TestWB_SqrtController_OurorobosApproximation_UsesFeedbackWhenAvailable(t *testing.T) {
	c := newSqrtController(PQStreamFromRational(RationalFromInt64(2)))

	feedbackRange := exactRangeFromRational(NewRational(big.NewInt(3), big.NewInt(2)))
	if err := c.feedCertifiedTerm(NewRCFTerm(big.NewInt(1)), feedbackRange); err != nil {
		t.Fatalf("feedCertifiedTerm error = %v", err)
	}

	term, _, status, err := c.ourorobosApproximation().NextPQ()
	if err != nil {
		t.Fatalf("NextPQ error = %v", err)
	}
	if status != StatusOK {
		t.Fatalf("status = %v, want %v", status, StatusOK)
	}
	if term.P.Cmp(big.NewInt(1)) != 0 || term.Q.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("term = (%v,%v), want (1,1)", term.P, term.Q)
	}
}

func TestWB_SqrtController_OurorobosApproximation_RangeMatchesFedCertifiedRange(t *testing.T) {
	c := newSqrtController(PQStreamFromRational(RationalFromInt64(2)))

	feedbackRange := exactRangeFromRational(NewRational(big.NewInt(3), big.NewInt(2)))
	if err := c.feedCertifiedTerm(NewRCFTerm(big.NewInt(1)), feedbackRange); err != nil {
		t.Fatalf("feedCertifiedTerm error = %v", err)
	}

	got, err := c.ourorobosApproximation().Range()
	if err != nil {
		t.Fatalf("Range error = %v", err)
	}
	if !got.Inside {
		t.Fatal("Inside = false, want true")
	}
	if got.Lo.Value.Cmp(feedbackRange.Lo.Value) != 0 || got.Hi.Value.Cmp(feedbackRange.Hi.Value) != 0 {
		t.Fatalf("Range = [%v/%v,%v/%v], want exact 3/2",
			got.Lo.Value.Num(), got.Lo.Value.Den(),
			got.Hi.Value.Num(), got.Hi.Value.Den(),
		)
	}
	if got.Lo.Open || got.Hi.Open {
		t.Fatalf("openness = (%v,%v), want (false,false)", got.Lo.Open, got.Hi.Open)
	}
}

// core/sqrt_cycle3_wb_test.go v2
