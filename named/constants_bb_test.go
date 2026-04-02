// named/constants_bb_test.go v6
package named_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/egp/gosper-gcf/core"
	"github.com/egp/gosper-gcf/named"
)

func TestBB_Named_E_MatchesKnownPrefix50(t *testing.T) {
	g := core.NewGCF1(identityUnaryCoeffsNamedConstants(), named.E())
	assertRCFPrefixNamedConstants(t, g, eTermsNamedConstants(50))
}

func identityUnaryCoeffsNamedConstants() core.BLFTCoefficients {
	return core.BLFTCoefficients{
		A: big.NewInt(0),
		B: big.NewInt(1),
		C: big.NewInt(0),
		D: big.NewInt(0),
		E: big.NewInt(0),
		F: big.NewInt(0),
		G: big.NewInt(0),
		H: big.NewInt(1),
	}
}

func eTermsNamedConstants(n int) []int64 {
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

func assertRCFPrefixNamedConstants(t *testing.T, g *core.GCF, want []int64) {
	t.Helper()

	for i, w := range want {
		term, status, err := nextRCFWithTimeoutNamedConstants(t, g, time.Second)
		if err != nil {
			t.Fatalf("term %d NextRCF error = %v", i+1, err)
		}
		if status != core.StatusOK {
			t.Fatalf("term %d status = %v, want %v", i+1, status, core.StatusOK)
		}
		if term.A().Cmp(big.NewInt(w)) != 0 {
			t.Fatalf("term %d = %v, want %d", i+1, term.A(), w)
		}
	}
}

func nextRCFWithTimeoutNamedConstants(t *testing.T, g *core.GCF, timeout time.Duration) (core.RCFTerm, core.Status, error) {
	t.Helper()

	type result struct {
		term   core.RCFTerm
		status core.Status
		err    error
	}

	ch := make(chan result, 1)
	go func() {
		term, status, err := g.NextRCF()
		ch <- result{term: term, status: status, err: err}
	}()

	select {
	case got := <-ch:
		return got.term, got.status, got.err
	case <-time.After(timeout):
		t.Fatalf("NextRCF timed out after %v", timeout)
		return core.NewRCFTerm(nil), core.StatusInvalidInput, nil
	}
}

// named/constants_bb_test.go v6
