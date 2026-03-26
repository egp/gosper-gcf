// named/e_cycle1_wb_test.go v2
package named

import (
	"math/big"
	"testing"

	"github.com/egp/gosper-gcf/core"
)

func TestWB_ETermAt_MatchesKnownPeriodicPattern(t *testing.T) {
	want := []int64{2, 1, 2, 1, 1, 4, 1, 1, 6, 1, 1, 8}

	for i, w := range want {
		got := eTermAt(i + 1)
		if got != w {
			t.Fatalf("term %d = %d, want %d", i+1, got, w)
		}
	}
}

func TestWB_ELookaheadRange_UsesOpenOpenWindow_ForHead(t *testing.T) {
	got := eLookaheadRange(2, 1)

	want := core.Range{
		Lo: core.Endpoint{
			Value: core.NewRational(big.NewInt(5), big.NewInt(2)),
			Open:  true,
		},
		Hi: core.Endpoint{
			Value: core.NewRational(big.NewInt(3), big.NewInt(1)),
			Open:  true,
		},
		Inside: true,
	}

	assertEExactInterval(t, got, want)
}

func TestWB_ELookaheadRange_UsesOpenOpenWindow_ForLargerNext(t *testing.T) {
	got := eLookaheadRange(1, 8)

	want := core.Range{
		Lo: core.Endpoint{
			Value: core.NewRational(big.NewInt(10), big.NewInt(9)),
			Open:  true,
		},
		Hi: core.Endpoint{
			Value: core.NewRational(big.NewInt(9), big.NewInt(8)),
			Open:  true,
		},
		Inside: true,
	}

	assertEExactInterval(t, got, want)
}

func assertEExactInterval(t *testing.T, got core.Range, want core.Range) {
	t.Helper()

	if got.Inside != want.Inside {
		t.Fatalf("Inside = %v, want %v", got.Inside, want.Inside)
	}
	if got.Lo.Open != want.Lo.Open {
		t.Fatalf("Lo.Open = %v, want %v", got.Lo.Open, want.Lo.Open)
	}
	if got.Hi.Open != want.Hi.Open {
		t.Fatalf("Hi.Open = %v, want %v", got.Hi.Open, want.Hi.Open)
	}
	if got.Lo.Value.Cmp(want.Lo.Value) != 0 {
		t.Fatalf(
			"Lo = %v/%v, want %v/%v",
			got.Lo.Value.Num(), got.Lo.Value.Den(),
			want.Lo.Value.Num(), want.Lo.Value.Den(),
		)
	}
	if got.Hi.Value.Cmp(want.Hi.Value) != 0 {
		t.Fatalf(
			"Hi = %v/%v, want %v/%v",
			got.Hi.Value.Num(), got.Hi.Value.Den(),
			want.Hi.Value.Num(), want.Hi.Value.Den(),
		)
	}
}

// named/e_cycle1_wb_test.go v2
