// core/rational_source_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

func TestWB_PQStreamFromRational_ProducesCanonicalFiniteGeneralizedTerms(t *testing.T) {
	src := PQStreamFromRational(NewRational(big.NewInt(3), big.NewInt(2)))

	term1, tail, status1 := src.NextPQ()
	if status1 != StatusOK {
		t.Fatalf("first status = %v, want %v", status1, StatusOK)
	}
	if term1.P.Cmp(big.NewInt(1)) != 0 || term1.Q.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("first term = (%v,%v), want (1,1)", term1.P, term1.Q)
	}

	term2, _, status2 := tail.NextPQ()
	if status2 != StatusOK {
		t.Fatalf("second status = %v, want %v", status2, StatusOK)
	}
	if term2.P.Cmp(big.NewInt(2)) != 0 || term2.Q.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("second term = (%v,%v), want (2,1)", term2.P, term2.Q)
	}
}

func TestWB_PQStreamFromRational_RangeTracksRemainingExactSuffix(t *testing.T) {
	src := PQStreamFromRational(NewRational(big.NewInt(3), big.NewInt(2)))

	r0 := src.Range()
	want0 := NewRational(big.NewInt(3), big.NewInt(2))
	assertExactRangeEqualsRationalSource(t, r0, want0, 0)

	_, tail, status := src.NextPQ()
	if status != StatusOK {
		t.Fatalf("status = %v, want %v", status, StatusOK)
	}

	r1 := tail.Range()
	want1 := NewRational(big.NewInt(2), big.NewInt(1))
	assertExactRangeEqualsRationalSource(t, r1, want1, 1)
}

func assertExactRangeEqualsRationalSource(t *testing.T, got Range, want Rational, step int) {
	t.Helper()

	if !got.Inside {
		t.Fatalf("range %d has Inside=false, want true", step)
	}
	if got.Lo.Open || got.Hi.Open {
		t.Fatalf("range %d openness wrong, want both closed", step)
	}
	if got.Lo.Value.Cmp(want) != 0 || got.Hi.Value.Cmp(want) != 0 {
		t.Fatalf("range %d = [%v/%v,%v/%v], want exact %v/%v",
			step,
			got.Lo.Value.Num(), got.Lo.Value.Den(),
			got.Hi.Value.Num(), got.Hi.Value.Den(),
			want.Num(), want.Den(),
		)
	}
}

// core/rational_source_wb_test.go v1
