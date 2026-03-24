package core

import (
	"math/big"
	"testing"
)

func TestWB_Range_InsideLoEqHiMeansExact(t *testing.T) {
	half := NewRational(big.NewInt(1), big.NewInt(2))

	r := Range{
		Lo:     half,
		Hi:     half,
		LoOpen: false,
		HiOpen: false,
		Inside: true,
	}

	if !r.Inside {
		t.Fatal("Inside = false, want true")
	}
	if r.Kind() != InsideInterval {
		t.Fatalf("Kind() = %v, want %v", r.Kind(), InsideInterval)
	}
	if r.Lo.Cmp(r.Hi) != 0 {
		t.Fatalf("Lo and Hi differ: Lo=%v/%v Hi=%v/%v",
			r.Lo.Num(), r.Lo.Den(), r.Hi.Num(), r.Hi.Den())
	}
	if r.LoOpen {
		t.Fatal("LoOpen = true, want false")
	}
	if r.HiOpen {
		t.Fatal("HiOpen = true, want false")
	}
}
