// named/e_cycle2_wb_test.go v4
package named

import (
	"math/big"
	"testing"

	"github.com/egp/gosper-gcf/core"
)

func TestWB_E_ReturnsProceduralSource(t *testing.T) {
	if _, ok := E().(*eProceduralStream); !ok {
		t.Fatalf("E() type = %T, want *eProceduralStream", E())
	}
}

func TestWB_EProceduralStream_NextPQ_AdvancesToNextTerm(t *testing.T) {
	src := &eProceduralStream{index: 1}

	first, tail, status, err := src.NextPQ()
	if err != nil {
		t.Fatalf("first NextPQ error = %v", err)
	}
	if status != core.StatusOK {
		t.Fatalf("first status = %v, want %v", status, core.StatusOK)
	}
	if first.P.Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("first P = %v, want 2", first.P)
	}
	if first.Q.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("first Q = %v, want 1", first.Q)
	}

	nextSrc, ok := tail.(*eProceduralStream)
	if !ok {
		t.Fatalf("tail type = %T, want *eProceduralStream", tail)
	}

	second, _, secondStatus, err := nextSrc.NextPQ()
	if err != nil {
		t.Fatalf("second NextPQ error = %v", err)
	}
	if secondStatus != core.StatusOK {
		t.Fatalf("second status = %v, want %v", secondStatus, core.StatusOK)
	}
	if second.P.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("second P = %v, want 1", second.P)
	}
	if second.Q.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("second Q = %v, want 1", second.Q)
	}
}

func TestWB_EProceduralStream_Range_UsesInfiniteLookaheadWindow(t *testing.T) {
	src := &eProceduralStream{index: 1}

	got, err := src.Range()
	if err != nil {
		t.Fatalf("Range error = %v", err)
	}
	want, err := eLookaheadRange(2, 1)
	if err != nil {
		t.Fatalf("eLookaheadRange error = %v", err)
	}
	assertEProceduralRange(t, got, want)
}

func TestWB_EProceduralStream_Range_BeforeThirdTerm_IsOpenOpenFiveHalvesToThree(t *testing.T) {
	src := &eProceduralStream{index: 1}

	_, tail1, status1, err := src.NextPQ()
	if err != nil {
		t.Fatalf("first NextPQ error = %v", err)
	}
	if status1 != core.StatusOK {
		t.Fatalf("first status = %v, want %v", status1, core.StatusOK)
	}

	_, tail2, status2, err := tail1.NextPQ()
	if err != nil {
		t.Fatalf("second NextPQ error = %v", err)
	}
	if status2 != core.StatusOK {
		t.Fatalf("second status = %v, want %v", status2, core.StatusOK)
	}

	got, err := tail2.Range()
	if err != nil {
		t.Fatalf("tail Range error = %v", err)
	}
	want, err := eLookaheadRange(2, 1)
	if err != nil {
		t.Fatalf("eLookaheadRange error = %v", err)
	}
	assertEProceduralRange(t, got, want)
}

func assertEProceduralRange(t *testing.T, got core.Range, want core.Range) {
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

// named/e_cycle2_wb_test.go v4
