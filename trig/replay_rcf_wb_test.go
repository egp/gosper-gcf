// trig/replay_rcf_wb_test.go v2
package trig

import (
	"math/big"
	"testing"

	"github.com/egp/gosper-gcf/core"
)

func TestWB_ReplayRCF_FreshForkStartsAtOriginalRange(t *testing.T) {
	root := newReplayRCF(replayExactTerminalSource())
	fork := root.Fork()

	assertExactRangeReplay(t, fork.Range(), 19, 5)
}

func TestWB_ReplayRCF_ForksReplaySameFirstTermIndependently(t *testing.T) {
	root := newReplayRCF(replayExactTerminalSource())
	left := root.Fork()
	right := root.Fork()

	leftTerm, leftStatus := left.NextRCF()
	if leftStatus != core.StatusOK {
		t.Fatalf("left first status = %v, want %v", leftStatus, core.StatusOK)
	}
	if leftTerm.A().Cmp(big.NewInt(3)) != 0 {
		t.Fatalf("left first term = %v, want 3", leftTerm.A())
	}

	rightTerm, rightStatus := right.NextRCF()
	if rightStatus != core.StatusOK {
		t.Fatalf("right first status = %v, want %v", rightStatus, core.StatusOK)
	}
	if rightTerm.A().Cmp(big.NewInt(3)) != 0 {
		t.Fatalf("right first term = %v, want 3", rightTerm.A())
	}
}

func TestWB_ReplayRCF_ForkRangeAdvancesToSuffixAfterRead(t *testing.T) {
	root := newReplayRCF(replayExactTerminalSource())
	fork := root.Fork()

	term1, status1 := fork.NextRCF()
	if status1 != core.StatusOK {
		t.Fatalf("first status = %v, want %v", status1, core.StatusOK)
	}
	if term1.A().Cmp(big.NewInt(3)) != 0 {
		t.Fatalf("first term = %v, want 3", term1.A())
	}
	assertExactRangeReplay(t, fork.Range(), 5, 4)

	term2, status2 := fork.NextRCF()
	if status2 != core.StatusOK {
		t.Fatalf("second status = %v, want %v", status2, core.StatusOK)
	}
	if term2.A().Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("second term = %v, want 1", term2.A())
	}
	assertExactRangeReplay(t, fork.Range(), 4, 1)
}

func TestWB_ReplayRCF_FinalSuffixIsZeroThenEOF(t *testing.T) {
	root := newReplayRCF(replayExactTerminalSource())
	fork := root.Fork()

	for i, want := range []int64{3, 1, 4} {
		term, status := fork.NextRCF()
		if status != core.StatusOK {
			t.Fatalf("term %d status = %v, want %v", i+1, status, core.StatusOK)
		}
		if term.A().Cmp(big.NewInt(want)) != 0 {
			t.Fatalf("term %d = %v, want %d", i+1, term.A(), want)
		}
	}

	assertExactRangeReplay(t, fork.Range(), 0, 1)

	_, eofStatus := fork.NextRCF()
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}
}

func replayExactTerminalSource() *core.GCF {
	return core.NewExactTerminalGCF(
		[]core.RCFTerm{
			core.NewRCFTerm(big.NewInt(3)),
			core.NewRCFTerm(big.NewInt(1)),
			core.NewRCFTerm(big.NewInt(4)),
		},
		exactRangeReplay(19, 5),
	)
}

func exactRangeReplay(num, den int64) core.Range {
	value := core.NewRational(big.NewInt(num), big.NewInt(den))
	return core.Range{
		Lo:     core.Endpoint{Value: value, Open: false},
		Hi:     core.Endpoint{Value: value, Open: false},
		Inside: true,
	}
}

func assertExactRangeReplay(t *testing.T, got core.Range, wantNum, wantDen int64) {
	t.Helper()

	want := core.NewRational(big.NewInt(wantNum), big.NewInt(wantDen))
	if !got.Inside {
		t.Fatal("Inside = false, want true")
	}
	if got.Lo.Open {
		t.Fatal("Lo.Open = true, want false")
	}
	if got.Hi.Open {
		t.Fatal("Hi.Open = true, want false")
	}
	if got.Lo.Value.Cmp(want) != 0 {
		t.Fatalf(
			"Lo = %v/%v, want %v/%v",
			got.Lo.Value.Num(), got.Lo.Value.Den(),
			want.Num(), want.Den(),
		)
	}
	if got.Hi.Value.Cmp(want) != 0 {
		t.Fatalf(
			"Hi = %v/%v, want %v/%v",
			got.Hi.Value.Num(), got.Hi.Value.Den(),
			want.Num(), want.Den(),
		)
	}
}

// trig/replay_rcf_wb_test.go v2
