// named/pi_gauss_wb_test.go v2
package named

import (
	"math/big"
	"testing"

	"github.com/egp/gosper-gcf/core"
)

func TestWB_PiGaussTermAt_HeadAndEarlyTerms(t *testing.T) {
	got0, err := piGaussTermAt(0)
	if err != nil {
		t.Fatalf("piGaussTermAt(0) error = %v", err)
	}
	assertPiGaussTerm(t, got0, 0, 4, 0)

	got1, err := piGaussTermAt(1)
	if err != nil {
		t.Fatalf("piGaussTermAt(1) error = %v", err)
	}
	assertPiGaussTerm(t, got1, 1, 1, 1)

	got2, err := piGaussTermAt(2)
	if err != nil {
		t.Fatalf("piGaussTermAt(2) error = %v", err)
	}
	assertPiGaussTerm(t, got2, 3, 4, 2)

	got3, err := piGaussTermAt(3)
	if err != nil {
		t.Fatalf("piGaussTermAt(3) error = %v", err)
	}
	assertPiGaussTerm(t, got3, 5, 9, 3)
}

func TestWB_PiGaussLookaheadRange_HeadIsOpenOpenZeroToFour(t *testing.T) {
	got, err := piGaussLookaheadRange(0)
	if err != nil {
		t.Fatalf("piGaussLookaheadRange(0) error = %v", err)
	}

	want := core.Range{
		Lo:     core.Endpoint{Value: core.RationalFromInt64(0), Open: true},
		Hi:     core.Endpoint{Value: core.RationalFromInt64(4), Open: true},
		Inside: true,
	}
	assertPiGaussRangeWB(t, got, want)
}

func TestWB_PiGaussLookaheadRange_ThirdStageUsesNextOddBound(t *testing.T) {
	got, err := piGaussLookaheadRange(2)
	if err != nil {
		t.Fatalf("piGaussLookaheadRange(2) error = %v", err)
	}

	want := core.Range{
		Lo: core.Endpoint{Value: core.RationalFromInt64(3), Open: true},
		Hi: core.Endpoint{
			Value: core.NewRational(big.NewInt(19), big.NewInt(5)),
			Open:  true,
		},
		Inside: true,
	}
	assertPiGaussRangeWB(t, got, want)
}

func assertPiGaussTerm(t *testing.T, got core.PQTerm, wantP, wantQ int64, index int) {
	t.Helper()

	if got.P.Cmp(big.NewInt(wantP)) != 0 {
		t.Fatalf("index %d P = %v, want %d", index, got.P, wantP)
	}
	if got.Q.Cmp(big.NewInt(wantQ)) != 0 {
		t.Fatalf("index %d Q = %v, want %d", index, got.Q, wantQ)
	}
}

func assertPiGaussRangeWB(t *testing.T, got core.Range, want core.Range) {
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

// named/pi_gauss_wb_test.go v2
