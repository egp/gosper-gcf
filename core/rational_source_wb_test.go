// core/rational_source_wb_test.go v2
package core

import (
	"math/big"
	"testing"
)

func TestWB_RationalSource_RangeIsExactInput(t *testing.T) {
	src := PQStreamFromRational(NewRational(big.NewInt(7), big.NewInt(5)))

	r, err := src.Range()
	if err != nil {
		t.Fatalf("Range error = %v", err)
	}
	want := NewRational(big.NewInt(7), big.NewInt(5))
	if !r.Inside {
		t.Fatal("Inside = false, want true")
	}
	if r.Lo.Open || r.Hi.Open {
		t.Fatalf("openness = (%v,%v), want (false,false)", r.Lo.Open, r.Hi.Open)
	}
	if r.Lo.Value.Cmp(want) != 0 || r.Hi.Value.Cmp(want) != 0 {
		t.Fatalf("Range = [%v/%v,%v/%v], want exact 7/5",
			r.Lo.Value.Num(), r.Lo.Value.Den(),
			r.Hi.Value.Num(), r.Hi.Value.Den(),
		)
	}
}

func TestWB_RationalSource_EmitsFiniteRCFAsPQTerms(t *testing.T) {
	src := PQStreamFromRational(NewRational(big.NewInt(7), big.NewInt(5)))

	want := []int64{1, 2, 2}
	stream := src

	for i, w := range want {
		term, tail, status, err := stream.NextPQ()
		if err != nil {
			t.Fatalf("term %d NextPQ error = %v", i+1, err)
		}
		if status != StatusOK {
			t.Fatalf("term %d status = %v, want %v", i+1, status, StatusOK)
		}
		if term.P.Cmp(big.NewInt(w)) != 0 || term.Q.Cmp(big.NewInt(1)) != 0 {
			t.Fatalf("term %d = (%v,%v), want (%d,1)", i+1, term.P, term.Q, w)
		}
		stream = tail
	}

	term, _, status, err := stream.NextPQ()
	if err != nil {
		t.Fatalf("EOF NextPQ error = %v", err)
	}
	if status != StatusEOF {
		t.Fatalf("EOF status = %v, want %v", status, StatusEOF)
	}
	if term.P.Cmp(big.NewInt(0)) != 0 || term.Q.Cmp(big.NewInt(0)) != 0 {
		t.Fatalf("EOF term = (%v,%v), want (0,0)", term.P, term.Q)
	}
}

// new regression coverage for the error-channel migration.
func TestWB_RationalSource_RepeatedRangeDoesNotConsume(t *testing.T) {
	src := PQStreamFromRational(NewRational(big.NewInt(7), big.NewInt(5)))

	r1, err := src.Range()
	if err != nil {
		t.Fatalf("first Range error = %v", err)
	}
	r2, err := src.Range()
	if err != nil {
		t.Fatalf("second Range error = %v", err)
	}

	if r1.Lo.Value.Cmp(r2.Lo.Value) != 0 || r1.Hi.Value.Cmp(r2.Hi.Value) != 0 {
		t.Fatalf("Range changed across repeated calls")
	}
	if r1.Lo.Open != r2.Lo.Open || r1.Hi.Open != r2.Hi.Open || r1.Inside != r2.Inside {
		t.Fatalf("Range openness/inside changed across repeated calls")
	}
}

// core/rational_source_wb_test.go v2
