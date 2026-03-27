// named/constants_range_bb_test.go v5
package named_test

import (
	"math/big"
	"os"
	"testing"

	"github.com/egp/gosper-gcf/core"
	"github.com/egp/gosper-gcf/named"
)

const pendingTestNamedERange = true

func TestBB_Named_E_RangeTracksBestLookaheadInterval50(t *testing.T) {

	terms := eTermsNamedRange(51)
	var src core.PQStream = named.E()

	for i := 0; i < 50; i++ {
		wantRange := namedLookaheadRange(terms[i], terms[i+1])
		assertIntervalRangeEqualsNamed(t, src.Range(), wantRange, i+1)

		got, tail, status := src.NextPQ()
		if status != core.StatusOK {
			t.Fatalf("term %d status = %v, want %v", i+1, status, core.StatusOK)
		}
		if got.P.Cmp(big.NewInt(terms[i])) != 0 {
			t.Fatalf("term %d P = %v, want %d", i+1, got.P, terms[i])
		}
		if got.Q.Cmp(big.NewInt(1)) != 0 {
			t.Fatalf("term %d Q = %v, want 1", i+1, got.Q)
		}

		src = tail
	}
}

func shouldSkipPendingNamedERange() bool {
	return pendingTestNamedERange && os.Getenv("RUN_PENDING_TESTS") == ""
}

func eTermsNamedRange(n int) []int64 {
	if n <= 0 {
		return nil
	}

	out := make([]int64, 0, n)
	out = append(out, 2)

	k := int64(1)
	for len(out) < n {
		out = append(out, 1)
		if len(out) >= n {
			break
		}
		out = append(out, 2*k)
		if len(out) >= n {
			break
		}
		out = append(out, 1)
		k++
	}

	return out
}

func namedLookaheadRange(a, next int64) core.Range {
	loNum := big.NewInt(a*(next+1) + 1)
	loDen := big.NewInt(next + 1)

	hiNum := big.NewInt(a*next + 1)
	hiDen := big.NewInt(next)

	return core.Range{
		Lo: core.Endpoint{
			Value: core.NewRational(loNum, loDen),
			Open:  true,
		},
		Hi: core.Endpoint{
			Value: core.NewRational(hiNum, hiDen),
			Open:  true,
		},
		Inside: true,
	}
}

func assertIntervalRangeEqualsNamed(t *testing.T, got core.Range, want core.Range, step int) {
	t.Helper()

	if got.Inside != want.Inside {
		t.Fatalf("range before term %d Inside = %v, want %v", step, got.Inside, want.Inside)
	}
	if got.Lo.Open != want.Lo.Open {
		t.Fatalf("range before term %d Lo.Open = %v, want %v", step, got.Lo.Open, want.Lo.Open)
	}
	if got.Hi.Open != want.Hi.Open {
		t.Fatalf("range before term %d Hi.Open = %v, want %v", step, got.Hi.Open, want.Hi.Open)
	}
	if got.Lo.Value.Cmp(want.Lo.Value) != 0 {
		t.Fatalf(
			"range before term %d Lo = %v/%v, want %v/%v",
			step, got.Lo.Value.Num(), got.Lo.Value.Den(),
			want.Lo.Value.Num(), want.Lo.Value.Den(),
		)
	}
	if got.Hi.Value.Cmp(want.Hi.Value) != 0 {
		t.Fatalf(
			"range before term %d Hi = %v/%v, want %v/%v",
			step, got.Hi.Value.Num(), got.Hi.Value.Den(),
			want.Hi.Value.Num(), want.Hi.Value.Den(),
		)
	}
}

// named/constants_range_bb_test.go v5
