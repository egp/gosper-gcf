// core/rational_wb_test.go v2
package core

import (
	"math/big"
	"testing"
)

func TestWB_Rational_NormalizesSign(t *testing.T) {
	r := NewRational(big.NewInt(2), big.NewInt(-4))

	if r.num.Cmp(big.NewInt(-1)) != 0 {
		t.Fatalf("normalized numerator = %v, want -1", r.num)
	}
	if r.den.Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("normalized denominator = %v, want 2", r.den)
	}
}

func TestWB_Rational_ReducesByGCD(t *testing.T) {
	r := NewRational(big.NewInt(6), big.NewInt(8))

	if r.num.Cmp(big.NewInt(3)) != 0 {
		t.Fatalf("reduced numerator = %v, want 3", r.num)
	}
	if r.den.Cmp(big.NewInt(4)) != 0 {
		t.Fatalf("reduced denominator = %v, want 4", r.den)
	}
}

func TestWB_Rational_FromIntIsExact(t *testing.T) {
	r := RationalFromInt64(-7)

	if r.num.Cmp(big.NewInt(-7)) != 0 {
		t.Fatalf("numerator = %v, want -7", r.num)
	}
	if r.den.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("denominator = %v, want 1", r.den)
	}
}

// core/rational_wb_test.go v2
