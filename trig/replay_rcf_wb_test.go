// trig/replay_rcf_wb_test.go v2
package trig

import (
	"math/big"
	"testing"

	"github.com/egp/gosper-gcf/core"
)

func TestWB_SuffixRangeAfterRCFTermReplay_MatchesCoreUnaryRange_WhenIntervalTouchesDenominatorZero(t *testing.T) {
	current := core.Range{
		Lo: core.Endpoint{
			Value: core.RationalFromInt64(3),
			Open:  false,
		},
		Hi: core.Endpoint{
			Value: core.RationalFromInt64(4),
			Open:  true,
		},
		Inside: true,
	}
	term := core.NewRCFTerm(big.NewInt(3))

	want, err := core.NewGCF1(
		core.BLFTCoefficients{
			A: big.NewInt(0),
			B: big.NewInt(0),
			C: big.NewInt(0),
			D: big.NewInt(1),
			E: big.NewInt(0),
			F: big.NewInt(1),
			G: big.NewInt(0),
			H: big.NewInt(-3),
		},
		&staticRangePQReplay{rng: current},
	).Range()
	if err != nil {
		t.Fatalf("want Range error = %v", err)
	}

	got, err := suffixRangeAfterRCFTermReplay(current, term)
	if err != nil {
		t.Fatalf("suffixRangeAfterRCFTermReplay error = %v", err)
	}

	assertSameRangeReplayWB(t, got, want)
}

func assertSameRangeReplayWB(t *testing.T, got, want core.Range) {
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

// trig/replay_rcf_wb_test.go v2
