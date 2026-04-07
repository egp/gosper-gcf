// core/sub_div_wb_test.go v2
package core

import (
	"math/big"
	"testing"
)

func TestWB_Sub_ExactFiniteInputs_ProducesExactDifference(t *testing.T) {
	g := Sub(
		PQStreamFromRational(RationalFromInt64(5)),
		PQStreamFromRational(RationalFromInt64(2)),
	)

	term, status, err := g.NextRCF()
	if err != nil {
		t.Fatalf("first NextRCF error = %v", err)
	}
	if status != StatusOK {
		t.Fatalf("first status = %v, want %v", status, StatusOK)
	}
	if term.A().Cmp(big.NewInt(3)) != 0 {
		t.Fatalf("first term = %v, want 3", term.A())
	}

	_, eofStatus, err := g.NextRCF()
	if err != nil {
		t.Fatalf("second NextRCF error = %v", err)
	}
	if eofStatus != StatusEOF {
		t.Fatalf("second status = %v, want %v", eofStatus, StatusEOF)
	}
}

func TestWB_Div_ExactFiniteInputs_ProducesExactQuotient(t *testing.T) {
	g := Div(
		PQStreamFromRational(NewRational(big.NewInt(22), big.NewInt(7))),
		PQStreamFromRational(NewRational(big.NewInt(11), big.NewInt(7))),
	)

	term, status, err := g.NextRCF()
	if err != nil {
		t.Fatalf("first NextRCF error = %v", err)
	}
	if status != StatusOK {
		t.Fatalf("first status = %v, want %v", status, StatusOK)
	}
	if term.A().Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("first term = %v, want 2", term.A())
	}

	_, eofStatus, err := g.NextRCF()
	if err != nil {
		t.Fatalf("second NextRCF error = %v", err)
	}
	if eofStatus != StatusEOF {
		t.Fatalf("second status = %v, want %v", eofStatus, StatusEOF)
	}
}

// new regression coverage for the error-channel migration.
func TestWB_Div_ExactFiniteInputs_RangeIsExactQuotient(t *testing.T) {
	g := Div(
		PQStreamFromRational(NewRational(big.NewInt(22), big.NewInt(7))),
		PQStreamFromRational(NewRational(big.NewInt(11), big.NewInt(7))),
	)

	r, err := g.Range()
	if err != nil {
		t.Fatalf("Range error = %v", err)
	}
	want := RationalFromInt64(2)
	if !r.Inside {
		t.Fatal("Inside = false, want true")
	}
	if r.Lo.Value.Cmp(want) != 0 || r.Hi.Value.Cmp(want) != 0 {
		t.Fatalf("Range = [%v/%v,%v/%v], want exact 2/1",
			r.Lo.Value.Num(), r.Lo.Value.Den(),
			r.Hi.Value.Num(), r.Hi.Value.Den(),
		)
	}
}

// core/sub_div_wb_test.go v2
