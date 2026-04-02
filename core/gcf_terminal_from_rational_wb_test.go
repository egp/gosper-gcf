// core/gcf_terminal_from_rational_wb_test.go v2
package core

import (
	"math/big"
	"testing"
)

func TestWB_GCF_ExactTerminalFromRational_Positive(t *testing.T) {
	g := newExactTerminalGCFFromRational(NewRational(big.NewInt(7), big.NewInt(5)))

	r, err := g.Range()
	if err != nil {
		t.Fatalf("Range error = %v", err)
	}

	want := NewRational(big.NewInt(7), big.NewInt(5))
	if !r.Inside {
		t.Fatal("Inside = false, want true")
	}
	if r.Lo.Value.Cmp(want) != 0 {
		t.Fatalf("Lo = %v/%v, want 7/5", r.Lo.Value.Num(), r.Lo.Value.Den())
	}
	if r.Hi.Value.Cmp(want) != 0 {
		t.Fatalf("Hi = %v/%v, want 7/5", r.Hi.Value.Num(), r.Hi.Value.Den())
	}

	term1, status1, err := g.NextRCF()
	if err != nil {
		t.Fatalf("first NextRCF error = %v", err)
	}
	if status1 != StatusOK {
		t.Fatalf("first status = %v, want %v", status1, StatusOK)
	}
	if term1.A().Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("first term = %v, want 1", term1.A())
	}

	term2, status2, err := g.NextRCF()
	if err != nil {
		t.Fatalf("second NextRCF error = %v", err)
	}
	if status2 != StatusOK {
		t.Fatalf("second status = %v, want %v", status2, StatusOK)
	}
	if term2.A().Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("second term = %v, want 2", term2.A())
	}

	term3, status3, err := g.NextRCF()
	if err != nil {
		t.Fatalf("third NextRCF error = %v", err)
	}
	if status3 != StatusOK {
		t.Fatalf("third status = %v, want %v", status3, StatusOK)
	}
	if term3.A().Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("third term = %v, want 2", term3.A())
	}

	_, status4, err := g.NextRCF()
	if err != nil {
		t.Fatalf("fourth NextRCF error = %v", err)
	}
	if status4 != StatusEOF {
		t.Fatalf("fourth status = %v, want %v", status4, StatusEOF)
	}
}

func TestWB_GCF_ExactTerminalFromRational_Negative(t *testing.T) {
	g := newExactTerminalGCFFromRational(NewRational(big.NewInt(-7), big.NewInt(5)))

	r, err := g.Range()
	if err != nil {
		t.Fatalf("Range error = %v", err)
	}

	want := NewRational(big.NewInt(-7), big.NewInt(5))
	if !r.Inside {
		t.Fatal("Inside = false, want true")
	}
	if r.Lo.Value.Cmp(want) != 0 {
		t.Fatalf("Lo = %v/%v, want -7/5", r.Lo.Value.Num(), r.Lo.Value.Den())
	}
	if r.Hi.Value.Cmp(want) != 0 {
		t.Fatalf("Hi = %v/%v, want -7/5", r.Hi.Value.Num(), r.Hi.Value.Den())
	}

	wantTerms := []*big.Int{
		big.NewInt(-2),
		big.NewInt(1),
		big.NewInt(1),
		big.NewInt(2),
	}

	for i, wantTerm := range wantTerms {
		got, status, err := g.NextRCF()
		if err != nil {
			t.Fatalf("term %d NextRCF error = %v", i+1, err)
		}
		if status != StatusOK {
			t.Fatalf("term %d status = %v, want %v", i+1, status, StatusOK)
		}
		if got.A().Cmp(wantTerm) != 0 {
			t.Fatalf("term %d = %v, want %v", i+1, got.A(), wantTerm)
		}
	}

	_, eofStatus, err := g.NextRCF()
	if err != nil {
		t.Fatalf("EOF NextRCF error = %v", err)
	}
	if eofStatus != StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, StatusEOF)
	}
}

// core/gcf_terminal_from_rational_wb_test.go v2
