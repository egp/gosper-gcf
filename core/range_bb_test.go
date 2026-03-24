// core/range_bb_test.go v2
package core_test

import (
	"math/big"
	"testing"

	"github.com/egp/gosper-gcf/core"
)

func TestBB_Range_ExactValueIsNotNecessarilyInteger(t *testing.T) {
	half := core.NewRational(big.NewInt(1), big.NewInt(2))

	r := core.Range{
		Lo: core.Endpoint{
			Value: half,
			Open:  false,
		},
		Hi: core.Endpoint{
			Value: half,
			Open:  false,
		},
		Inside: true,
	}

	if r.Kind() != core.InsideInterval {
		t.Fatalf("Kind() = %v, want %v", r.Kind(), core.InsideInterval)
	}
	if r.Lo.Value.Num().Cmp(big.NewInt(1)) != 0 || r.Lo.Value.Den().Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("Lo = %v/%v, want 1/2", r.Lo.Value.Num(), r.Lo.Value.Den())
	}
	if r.Hi.Value.Num().Cmp(big.NewInt(1)) != 0 || r.Hi.Value.Den().Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("Hi = %v/%v, want 1/2", r.Hi.Value.Num(), r.Hi.Value.Den())
	}
}

func TestBB_Range_InsideAndOutsideKindsConstructCleanly(t *testing.T) {
	one := core.RationalFromInt64(1)
	two := core.RationalFromInt64(2)

	inside := core.Range{
		Lo: core.Endpoint{
			Value: one,
			Open:  false,
		},
		Hi: core.Endpoint{
			Value: two,
			Open:  false,
		},
		Inside: true,
	}
	outside := core.Range{
		Lo: core.Endpoint{
			Value: one,
			Open:  false,
		},
		Hi: core.Endpoint{
			Value: two,
			Open:  false,
		},
		Inside: false,
	}

	if inside.Kind() != core.InsideInterval {
		t.Fatalf("inside.Kind() = %v, want %v", inside.Kind(), core.InsideInterval)
	}
	if outside.Kind() != core.OutsideInterval {
		t.Fatalf("outside.Kind() = %v, want %v", outside.Kind(), core.OutsideInterval)
	}
}

// core/range_bb_test.go v2
