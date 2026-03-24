// core/range_wb_test.go v2
package core

import (
	"math/big"
	"testing"
)

func TestWB_Range_InsideLoEqHiMeansExact(t *testing.T) {
	half := NewRational(big.NewInt(1), big.NewInt(2))

	r := Range{
		Lo: Endpoint{
			Value: half,
			Open:  false,
		},
		Hi: Endpoint{
			Value: half,
			Open:  false,
		},
		Inside: true,
	}

	if !r.Inside {
		t.Fatal("Inside = false, want true")
	}
	if r.Kind() != InsideInterval {
		t.Fatalf("Kind() = %v, want %v", r.Kind(), InsideInterval)
	}
	if r.Lo.Value.Cmp(r.Hi.Value) != 0 {
		t.Fatalf("Lo and Hi differ: Lo=%v/%v Hi=%v/%v",
			r.Lo.Value.Num(), r.Lo.Value.Den(), r.Hi.Value.Num(), r.Hi.Value.Den())
	}
	if r.Lo.Open {
		t.Fatal("Lo.Open = true, want false")
	}
	if r.Hi.Open {
		t.Fatal("Hi.Open = true, want false")
	}
}

// core/range_wb_test.go v2
