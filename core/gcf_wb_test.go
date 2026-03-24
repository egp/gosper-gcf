// core/gcf_wb_test.go v1
package core

import (
	"math/big"
	"testing"
	"time"
)

func TestWB_GCF_TerminalExactStateEmitsCorrectly(t *testing.T) {
	g := NewExactTerminalGCF(
		[]RCFTerm{
			NewRCFTerm(big.NewInt(3)),
			NewRCFTerm(big.NewInt(4)),
		},
		testExactRange(7, 5),
	)

	term1, status1 := g.NextRCF()
	if status1 != StatusOK {
		t.Fatalf("first status = %v, want %v", status1, StatusOK)
	}
	if term1.A().Cmp(big.NewInt(3)) != 0 {
		t.Fatalf("first term = %v, want 3", term1.A())
	}
}

func TestWB_GCF_TerminalExactStateHasExactRange(t *testing.T) {
	g := NewExactTerminalGCF(
		[]RCFTerm{
			NewRCFTerm(big.NewInt(1)),
		},
		testExactRange(7, 5),
	)

	r := g.Range()

	want := NewRational(big.NewInt(7), big.NewInt(5))

	if !r.Inside {
		t.Fatal("Inside = false, want true")
	}
	if r.Lo.Value.Cmp(want) != 0 {
		t.Fatalf("Lo = %v/%v, want 7/5", r.Lo.Value.Num(), r.Lo.Value.Den())
	}
	if r.Hi.Value.Cmp(want) != 0 {
		t.Fatalf("Hi = %v/%v, want 7/5", r.Hi.Value.Num(), r.Hi.Value.Den())
	}
	if r.Lo.Open {
		t.Fatal("Lo.Open = true, want false")
	}
	if r.Hi.Open {
		t.Fatal("Hi.Open = true, want false")
	}
}

func TestWB_GCF_FirstEmittedRegularTermMayBeNegative(t *testing.T) {
	g := NewExactTerminalGCF(
		[]RCFTerm{
			NewRCFTerm(big.NewInt(-2)),
		},
		testExactRange(-2, 1),
	)

	term, status := g.NextRCF()
	if status != StatusOK {
		t.Fatalf("status = %v, want %v", status, StatusOK)
	}
	if term.A().Cmp(big.NewInt(-2)) != 0 {
		t.Fatalf("term = %v, want -2", term.A())
	}
}

func TestWB_GCF_EOFUsesStatusEOF(t *testing.T) {
	g := NewExactTerminalGCF(
		[]RCFTerm{
			NewRCFTerm(big.NewInt(9)),
		},
		testExactRange(9, 1),
	)

	_, status1 := g.NextRCF()
	if status1 != StatusOK {
		t.Fatalf("first status = %v, want %v", status1, StatusOK)
	}

	_, status2 := g.NextRCF()
	if status2 != StatusEOF {
		t.Fatalf("second status = %v, want %v", status2, StatusEOF)
	}
}

func TestWB_GCF_ConfigDefaultsApply(t *testing.T) {
	custom := Config{
		Timeout:             250 * time.Millisecond,
		BitLenLimit:         4096,
		BitLenWarnThreshold: 3072,
	}

	g1 := NewExactTerminalGCF(nil, testExactRange(0, 1))
	if g1.Config() != DefaultConfig() {
		t.Fatalf("default config = %#v, want %#v", g1.Config(), DefaultConfig())
	}

	g2 := NewExactTerminalGCFWithConfig(nil, testExactRange(0, 1), custom)
	if g2.Config() != custom {
		t.Fatalf("custom config = %#v, want %#v", g2.Config(), custom)
	}
}

func testExactRange(num, den int64) Range {
	value := NewRational(big.NewInt(num), big.NewInt(den))
	return Range{
		Lo: Endpoint{
			Value: value,
			Open:  false,
		},
		Hi: Endpoint{
			Value: value,
			Open:  false,
		},
		Inside: true,
	}
}

// core/gcf_wb_test.go v1
