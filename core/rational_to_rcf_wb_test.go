// core/rational_to_rcf_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

func TestWB_RCFTermsFromRational_FinitePositive(t *testing.T) {
	r := NewRational(big.NewInt(7), big.NewInt(5))

	got := rcfTermsFromRational(r)

	want := []RCFTerm{
		NewRCFTerm(big.NewInt(1)),
		NewRCFTerm(big.NewInt(2)),
		NewRCFTerm(big.NewInt(2)),
	}

	assertRCFTermsEqual(t, got, want)
}

func TestWB_RCFTermsFromRational_ExactInteger(t *testing.T) {
	r := RationalFromInt64(9)

	got := rcfTermsFromRational(r)

	want := []RCFTerm{
		NewRCFTerm(big.NewInt(9)),
	}

	assertRCFTermsEqual(t, got, want)
}

func TestWB_RCFTermsFromRational_FiniteNegative(t *testing.T) {
	r := NewRational(big.NewInt(-7), big.NewInt(5))

	got := rcfTermsFromRational(r)

	want := []RCFTerm{
		NewRCFTerm(big.NewInt(-2)),
		NewRCFTerm(big.NewInt(1)),
		NewRCFTerm(big.NewInt(1)),
		NewRCFTerm(big.NewInt(2)),
	}

	assertRCFTermsEqual(t, got, want)
}

func assertRCFTermsEqual(t *testing.T, got, want []RCFTerm) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("len(got) = %d, want %d", len(got), len(want))
	}

	for i := range want {
		if got[i].A().Cmp(want[i].A()) != 0 {
			t.Fatalf("term[%d] = %v, want %v", i, got[i].A(), want[i].A())
		}
	}
}

// core/rational_to_rcf_wb_test.go v1
