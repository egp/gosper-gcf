// core/gcf_bb_test.go v3
package core_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/egp/gosper-gcf/core"
)

func TestBB_GCF_ReadFiniteRegularTermsThenEOF(t *testing.T) {
	g := core.NewExactTerminalGCF(
		[]core.RCFTerm{
			core.NewRCFTerm(big.NewInt(3)),
			core.NewRCFTerm(big.NewInt(1)),
			core.NewRCFTerm(big.NewInt(4)),
		},
		exactRange(19, 6),
	)

	term1, status1, err := g.NextRCF()
	if err != nil {
		t.Fatalf("first NextRCF error = %v", err)
	}
	if status1 != core.StatusOK {
		t.Fatalf("first status = %v, want %v", status1, core.StatusOK)
	}
	if term1.A().Cmp(big.NewInt(3)) != 0 {
		t.Fatalf("first term = %v, want 3", term1.A())
	}

	term2, status2, err := g.NextRCF()
	if err != nil {
		t.Fatalf("second NextRCF error = %v", err)
	}
	if status2 != core.StatusOK {
		t.Fatalf("second status = %v, want %v", status2, core.StatusOK)
	}
	if term2.A().Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("second term = %v, want 1", term2.A())
	}

	term3, status3, err := g.NextRCF()
	if err != nil {
		t.Fatalf("third NextRCF error = %v", err)
	}
	if status3 != core.StatusOK {
		t.Fatalf("third status = %v, want %v", status3, core.StatusOK)
	}
	if term3.A().Cmp(big.NewInt(4)) != 0 {
		t.Fatalf("third term = %v, want 4", term3.A())
	}

	_, status4, err := g.NextRCF()
	if err != nil {
		t.Fatalf("fourth NextRCF error = %v", err)
	}
	if status4 != core.StatusEOF {
		t.Fatalf("fourth status = %v, want %v", status4, core.StatusEOF)
	}
}

func TestBB_GCF_RangeAvailableWhileLive(t *testing.T) {
	g := core.NewExactTerminalGCF(
		[]core.RCFTerm{
			core.NewRCFTerm(big.NewInt(2)),
			core.NewRCFTerm(big.NewInt(7)),
		},
		exactRange(19, 8),
	)

	r1, err := g.Range()
	if err != nil {
		t.Fatalf("initial Range error = %v", err)
	}
	want := core.NewRational(big.NewInt(19), big.NewInt(8))
	if !r1.Inside {
		t.Fatal("initial range Inside = false, want true")
	}
	if r1.Lo.Value.Cmp(want) != 0 || r1.Hi.Value.Cmp(want) != 0 {
		t.Fatalf("initial range = [%v/%v, %v/%v], want exact 19/8",
			r1.Lo.Value.Num(), r1.Lo.Value.Den(),
			r1.Hi.Value.Num(), r1.Hi.Value.Den(),
		)
	}

	_, status, err := g.NextRCF()
	if err != nil {
		t.Fatalf("NextRCF error = %v", err)
	}
	if status != core.StatusOK {
		t.Fatalf("NextRCF status = %v, want %v", status, core.StatusOK)
	}

	r2, err := g.Range()
	if err != nil {
		t.Fatalf("live Range error = %v", err)
	}
	if !r2.Inside {
		t.Fatal("live range Inside = false, want true")
	}
	if r2.Lo.Value.Cmp(want) != 0 || r2.Hi.Value.Cmp(want) != 0 {
		t.Fatalf("live range = [%v/%v, %v/%v], want exact 19/8",
			r2.Lo.Value.Num(), r2.Lo.Value.Den(),
			r2.Hi.Value.Num(), r2.Hi.Value.Den(),
		)
	}
}

func TestBB_GCF_ExactStateProducesDeterministicTerms(t *testing.T) {
	terms := []core.RCFTerm{
		core.NewRCFTerm(big.NewInt(-1)),
		core.NewRCFTerm(big.NewInt(2)),
		core.NewRCFTerm(big.NewInt(3)),
	}
	rng := exactRange(5, 3)

	g1 := core.NewExactTerminalGCF(terms, rng)
	g2 := core.NewExactTerminalGCF(terms, rng)

	for i, want := range []*big.Int{big.NewInt(-1), big.NewInt(2), big.NewInt(3)} {
		term1, status1, err := g1.NextRCF()
		if err != nil {
			t.Fatalf("g1 term %d NextRCF error = %v", i+1, err)
		}
		term2, status2, err := g2.NextRCF()
		if err != nil {
			t.Fatalf("g2 term %d NextRCF error = %v", i+1, err)
		}

		if status1 != core.StatusOK {
			t.Fatalf("g1 term %d status = %v, want %v", i+1, status1, core.StatusOK)
		}
		if status2 != core.StatusOK {
			t.Fatalf("g2 term %d status = %v, want %v", i+1, status2, core.StatusOK)
		}
		if term1.A().Cmp(want) != 0 {
			t.Fatalf("g1 term %d = %v, want %v", i+1, term1.A(), want)
		}
		if term2.A().Cmp(want) != 0 {
			t.Fatalf("g2 term %d = %v, want %v", i+1, term2.A(), want)
		}
	}

	_, status1, err := g1.NextRCF()
	if err != nil {
		t.Fatalf("g1 EOF NextRCF error = %v", err)
	}
	_, status2, err := g2.NextRCF()
	if err != nil {
		t.Fatalf("g2 EOF NextRCF error = %v", err)
	}
	if status1 != core.StatusEOF {
		t.Fatalf("g1 EOF status = %v, want %v", status1, core.StatusEOF)
	}
	if status2 != core.StatusEOF {
		t.Fatalf("g2 EOF status = %v, want %v", status2, core.StatusEOF)
	}
}

func TestBB_GCF_ConfigOptionalAtCreation(t *testing.T) {
	coeffs := identityCoefficients()

	g0 := core.NewGCF0(coeffs)
	if g0.Config() != core.DefaultConfig() {
		t.Fatalf("NewGCF0 default config = %#v, want %#v", g0.Config(), core.DefaultConfig())
	}

	stream, streamStatus := core.NewFinitePQStream([]core.FinitePQStep{
		{
			Term:  core.PQTerm{P: big.NewInt(1), Q: big.NewInt(1)},
			Range: insideRange(1, 2),
		},
	})
	if streamStatus != core.StatusOK {
		t.Fatalf("NewFinitePQStream status = %v, want %v", streamStatus, core.StatusOK)
	}

	g1 := core.NewGCF1(coeffs, stream)
	if g1.Config() != core.DefaultConfig() {
		t.Fatalf("NewGCF1 default config = %#v, want %#v", g1.Config(), core.DefaultConfig())
	}

	g2 := core.NewGCF2(coeffs, stream, stream)
	if g2.Config() != core.DefaultConfig() {
		t.Fatalf("NewGCF2 default config = %#v, want %#v", g2.Config(), core.DefaultConfig())
	}

	custom := core.Config{
		Timeout:             500 * time.Millisecond,
		BitLenLimit:         2048,
		BitLenWarnThreshold: 1536,
	}

	g0c := core.NewGCF0WithConfig(coeffs, custom)
	if g0c.Config() != custom {
		t.Fatalf("NewGCF0WithConfig config = %#v, want %#v", g0c.Config(), custom)
	}

	g1c := core.NewGCF1WithConfig(coeffs, stream, custom)
	if g1c.Config() != custom {
		t.Fatalf("NewGCF1WithConfig config = %#v, want %#v", g1c.Config(), custom)
	}

	g2c := core.NewGCF2WithConfig(coeffs, stream, stream, custom)
	if g2c.Config() != custom {
		t.Fatalf("NewGCF2WithConfig config = %#v, want %#v", g2c.Config(), custom)
	}
}

// new regression coverage for the error-channel migration.
func TestBB_GCF_ExactTerminalRangeStaysAvailableAfterEOF(t *testing.T) {
	g := core.NewExactTerminalGCF(
		[]core.RCFTerm{
			core.NewRCFTerm(big.NewInt(9)),
		},
		exactRange(9, 1),
	)

	_, status, err := g.NextRCF()
	if err != nil {
		t.Fatalf("first NextRCF error = %v", err)
	}
	if status != core.StatusOK {
		t.Fatalf("first status = %v, want %v", status, core.StatusOK)
	}

	_, eofStatus, err := g.NextRCF()
	if err != nil {
		t.Fatalf("EOF NextRCF error = %v", err)
	}
	if eofStatus != core.StatusEOF {
		t.Fatalf("EOF status = %v, want %v", eofStatus, core.StatusEOF)
	}

	r, err := g.Range()
	if err != nil {
		t.Fatalf("Range after EOF error = %v", err)
	}
	want := core.RationalFromInt64(9)
	if !r.Inside {
		t.Fatal("Inside = false, want true")
	}
	if r.Lo.Value.Cmp(want) != 0 || r.Hi.Value.Cmp(want) != 0 {
		t.Fatalf("Range after EOF = [%v/%v,%v/%v], want exact 9/1",
			r.Lo.Value.Num(), r.Lo.Value.Den(),
			r.Hi.Value.Num(), r.Hi.Value.Den(),
		)
	}
}

func exactRange(num, den int64) core.Range {
	value := core.NewRational(big.NewInt(num), big.NewInt(den))
	return core.Range{
		Lo: core.Endpoint{
			Value: value,
			Open:  false,
		},
		Hi: core.Endpoint{
			Value: value,
			Open:  false,
		},
		Inside: true,
	}
}

func insideRange(lo, hi int64) core.Range {
	return core.Range{
		Lo: core.Endpoint{
			Value: core.RationalFromInt64(lo),
			Open:  false,
		},
		Hi: core.Endpoint{
			Value: core.RationalFromInt64(hi),
			Open:  false,
		},
		Inside: true,
	}
}

func identityCoefficients() core.BLFTCoefficients {
	return core.BLFTCoefficients{
		A: big.NewInt(1),
		B: big.NewInt(0),
		C: big.NewInt(0),
		D: big.NewInt(0),
		E: big.NewInt(0),
		F: big.NewInt(0),
		G: big.NewInt(0),
		H: big.NewInt(1),
	}
}

// core/gcf_bb_test.go v3
