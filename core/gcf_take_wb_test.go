// core/gcf_take_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

func TestWB_GCF_TakeRCFTermsUpTo_CollectsFinitePrefix(t *testing.T) {
	g := NewGCF1(
		takeIdentityCoeffs(),
		PQStreamFromRational(NewRational(big.NewInt(8), big.NewInt(3))),
	)

	got := g.takeRCFTermsUpTo(10)

	if len(got) != 3 {
		t.Fatalf("len(got) = %d, want 3", len(got))
	}

	want := []int64{2, 1, 2}
	for i, w := range want {
		if got[i].A().Cmp(big.NewInt(w)) != 0 {
			t.Fatalf("term %d = %v, want %d", i+1, got[i].A(), w)
		}
	}
}

func TestWB_FinitePQStepsFromRCFTerms_BuildsExactSuffixRanges(t *testing.T) {
	terms := []RCFTerm{
		NewRCFTerm(big.NewInt(2)),
		NewRCFTerm(big.NewInt(1)),
		NewRCFTerm(big.NewInt(2)),
	}

	got := finitePQStepsFromRCFTerms(terms)

	if len(got) != 3 {
		t.Fatalf("len(got) = %d, want 3", len(got))
	}

	assertTakeStep(t, got[0], 2, 1, NewRational(big.NewInt(8), big.NewInt(3)), 1)
	assertTakeStep(t, got[1], 1, 1, NewRational(big.NewInt(3), big.NewInt(2)), 2)
	assertTakeStep(t, got[2], 2, 1, RationalFromInt64(2), 3)
}

func TestWB_FinitePQStepsFromRCFTerms_EmptyInputYieldsEmptySteps(t *testing.T) {
	got := finitePQStepsFromRCFTerms(nil)
	if len(got) != 0 {
		t.Fatalf("len(got) = %d, want 0", len(got))
	}
}

func takeIdentityCoeffs() BLFTCoefficients {
	return BLFTCoefficients{
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

func assertTakeStep(t *testing.T, got FinitePQStep, wantP, wantQ int64, wantRange Rational, step int) {
	t.Helper()

	if got.Term.P.Cmp(big.NewInt(wantP)) != 0 {
		t.Fatalf("step %d P = %v, want %d", step, got.Term.P, wantP)
	}
	if got.Term.Q.Cmp(big.NewInt(wantQ)) != 0 {
		t.Fatalf("step %d Q = %v, want %d", step, got.Term.Q, wantQ)
	}
	if !got.Range.Inside {
		t.Fatalf("step %d Inside = false, want true", step)
	}
	if got.Range.Lo.Open {
		t.Fatalf("step %d Lo.Open = true, want false", step)
	}
	if got.Range.Hi.Open {
		t.Fatalf("step %d Hi.Open = true, want false", step)
	}
	if got.Range.Lo.Value.Cmp(wantRange) != 0 {
		t.Fatalf(
			"step %d Lo = %v/%v, want %v/%v",
			step,
			got.Range.Lo.Value.Num(), got.Range.Lo.Value.Den(),
			wantRange.Num(), wantRange.Den(),
		)
	}
	if got.Range.Hi.Value.Cmp(wantRange) != 0 {
		t.Fatalf(
			"step %d Hi = %v/%v, want %v/%v",
			step,
			got.Range.Hi.Value.Num(), got.Range.Hi.Value.Den(),
			wantRange.Num(), wantRange.Den(),
		)
	}
}

// core/gcf_take_wb_test.go v1
