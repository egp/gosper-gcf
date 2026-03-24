// core/gcf_terminal_from_rational_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

func TestWB_GCF_ExactTerminalFromRational_Positive(t *testing.T) {
	g := newExactTerminalGCFFromRational(NewRational(big.NewInt(7), big.NewInt(5)))

	r := g.Range()
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

	term1, status1 := g.NextRCF()
	if status1 != StatusOK {
		t.Fatalf("first status = %v, want %v", status1, StatusOK)
	}
	if term1.A().Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("first term = %v, want 1", term1.A())
	}

	term2, status2 := g.NextRCF()
	if status2 != StatusOK {
		t.Fatalf("second status = %v, want %v", status2, StatusOK)
	}
	if term2.A().Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("second term = %v, want 2", term2.A())
	}

	term3, status3 := g.NextRCF()
	if status3 != StatusOK {
		t.Fatalf("third status = %v, want %v", status3, StatusOK)
	}
	if term3.A().Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("third term = %v, want 2", term3.A())
	}

	_, status4 := g.NextRCF()
	if status4 != StatusEOF {
		t.Fatalf("fourth status = %v, want %v", status4, StatusEOF)
	}
}

func TestWB_GCF_ExactTerminalFromRational_Negative(t *testing.T) {
	g := newExactTerminalGCFFromRational(NewRational(big.NewInt(-7), big.NewInt(5)))

	r := g.Range()
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
		got, status := g.NextRCF()
		if status != StatusOK {
			t.Fatalf("term %d status = %v, want %v", i+1, status, StatusOK)
		}
		if got.A().Cmp(wantTerm) != 0 {
			t.Fatalf("term %d = %v, want %v", i+1, got.A(), wantTerm)
		}
	}

	_, eofStatus := g.NextRCF()
	if eofStatus != StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, StatusEOF)
	}
}

// core/gcf_terminal_from_rational_wb_test.go v1
