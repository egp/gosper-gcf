// core/gcf_bb_test.go v2
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

	term1, status1 := g.NextRCF()
	if status1 != core.StatusOK {
		t.Fatalf("first status = %v, want %v", status1, core.StatusOK)
	}
	if term1.A().Cmp(big.NewInt(3)) != 0 {
		t.Fatalf("first term = %v, want 3", term1.A())
	}

	term2, status2 := g.NextRCF()
	if status2 != core.StatusOK {
		t.Fatalf("second status = %v, want %v", status2, core.StatusOK)
	}
	if term2.A().Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("second term = %v, want 1", term2.A())
	}

	term3, status3 := g.NextRCF()
	if status3 != core.StatusOK {
		t.Fatalf("third status = %v, want %v", status3, core.StatusOK)
	}
	if term3.A().Cmp(big.NewInt(4)) != 0 {
		t.Fatalf("third term = %v, want 4", term3.A())
	}

	_, status4 := g.NextRCF()
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

	r1 := g.Range()
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

	_, status := g.NextRCF()
	if status != core.StatusOK {
		t.Fatalf("NextRCF status = %v, want %v", status, core.StatusOK)
	}

	r2 := g.Range()
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
		term1, status1 := g1.NextRCF()
		term2, status2 := g2.NextRCF()

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

	_, status1 := g1.NextRCF()
	_, status2 := g2.NextRCF()
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

func identityCoefficients() core.TransformCoefficients {
	return core.TransformCoefficients{
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

// core/gcf_bb_test.go v2
