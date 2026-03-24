package core_test

import (
	"math/big"
	"testing"

	"github.com/egp/gosper-gcf/core"
)

func TestBB_Range_ExactValueIsNotNecessarilyInteger(t *testing.T) {
	half := core.NewRational(big.NewInt(1), big.NewInt(2))

	r := core.Range{
		Lo:     half,
		Hi:     half,
		LoOpen: false,
		HiOpen: false,
		Inside: true,
	}

	if r.Kind() != core.InsideInterval {
		t.Fatalf("Kind() = %v, want %v", r.Kind(), core.InsideInterval)
	}
	if r.Lo.Num().Cmp(big.NewInt(1)) != 0 || r.Lo.Den().Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("Lo = %v/%v, want 1/2", r.Lo.Num(), r.Lo.Den())
	}
	if r.Hi.Num().Cmp(big.NewInt(1)) != 0 || r.Hi.Den().Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("Hi = %v/%v, want 1/2", r.Hi.Num(), r.Hi.Den())
	}
}

func TestBB_Range_InsideAndOutsideKindsConstructCleanly(t *testing.T) {
	one := core.FromInt64(1)
	two := core.FromInt64(2)

	inside := core.Range{
		Lo:     one,
		Hi:     two,
		LoOpen: false,
		HiOpen: false,
		Inside: true,
	}
	outside := core.Range{
		Lo:     one,
		Hi:     two,
		LoOpen: false,
		HiOpen: false,
		Inside: false,
	}

	if inside.Kind() != core.InsideInterval {
		t.Fatalf("inside.Kind() = %v, want %v", inside.Kind(), core.InsideInterval)
	}
	if outside.Kind() != core.OutsideInterval {
		t.Fatalf("outside.Kind() = %v, want %v", outside.Kind(), core.OutsideInterval)
	}
}
