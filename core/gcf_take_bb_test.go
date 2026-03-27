// core/gcf_take_bb_test.go v2
package core_test

import (
	"math/big"
	"os"
	"testing"

	"github.com/egp/gosper-gcf/core"
	"github.com/egp/gosper-gcf/named"
)

const pendingTestGCFTake = true
const pendingTestGCFRational = true

func TestBB_GCF_Take_PrefixOfInfiniteSourceIsFiniteExactPQStream(t *testing.T) {

	g := core.NewGCF1(identityUnaryCoeffsTake(), named.E())
	taken := g.Take(3)

	_, tail2 := assertFinitePQStepTake(
		t,
		taken,
		1,
		2,
		1,
		core.NewRational(big.NewInt(8), big.NewInt(3)),
	)

	_, tail3 := assertFinitePQStepTake(
		t,
		tail2,
		2,
		1,
		1,
		core.NewRational(big.NewInt(3), big.NewInt(2)),
	)

	_, tail4 := assertFinitePQStepTake(
		t,
		tail3,
		3,
		2,
		1,
		core.RationalFromInt64(2),
	)

	assertFinitePQEOF(t, tail4, 4)
}

func TestBB_GCF_Rational_ReturnsConvergentOfTakenPrefix(t *testing.T) {

	g := core.NewGCF1(identityUnaryCoeffsTake(), named.E())
	got := g.Rational(3)
	want := core.NewRational(big.NewInt(8), big.NewInt(3))

	if got.Cmp(want) != 0 {
		t.Fatalf(
			"Rational(3) = %v/%v, want %v/%v",
			got.Num(), got.Den(),
			want.Num(), want.Den(),
		)
	}
}

func TestBB_GCF_Take_StopsAtEOF_WhenSourceHasFewerTerms(t *testing.T) {

	g := core.NewGCF1(
		identityUnaryCoeffsTake(),
		core.PQStreamFromRational(core.NewRational(big.NewInt(8), big.NewInt(3))),
	)

	taken := g.Take(10)

	_, tail2 := assertFinitePQStepTake(
		t,
		taken,
		1,
		2,
		1,
		core.NewRational(big.NewInt(8), big.NewInt(3)),
	)

	_, tail3 := assertFinitePQStepTake(
		t,
		tail2,
		2,
		1,
		1,
		core.NewRational(big.NewInt(3), big.NewInt(2)),
	)

	_, tail4 := assertFinitePQStepTake(
		t,
		tail3,
		3,
		2,
		1,
		core.RationalFromInt64(2),
	)

	assertFinitePQEOF(t, tail4, 4)
}

func shouldSkipPendingGCFTake() bool {
	return pendingTestGCFTake && os.Getenv("RUN_PENDING_TESTS") == ""
}

func shouldSkipPendingGCFRational() bool {
	return pendingTestGCFRational && os.Getenv("RUN_PENDING_TESTS") == ""
}

func identityUnaryCoeffsTake() core.BLFTCoefficients {
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

func assertFinitePQStepTake(
	t *testing.T,
	src core.PQStream,
	step int,
	wantP int64,
	wantQ int64,
	wantRange core.Rational,
) (core.PQTerm, core.PQStream) {
	t.Helper()

	gotRange := src.Range()
	assertExactTakeRange(t, gotRange, wantRange, step)

	term, tail, status := src.NextPQ()
	if status != core.StatusOK {
		t.Fatalf("step %d status = %v, want %v", step, status, core.StatusOK)
	}

	if term.P.Cmp(big.NewInt(wantP)) != 0 {
		t.Fatalf("step %d P = %v, want %d", step, term.P, wantP)
	}
	if term.Q.Cmp(big.NewInt(wantQ)) != 0 {
		t.Fatalf("step %d Q = %v, want %d", step, term.Q, wantQ)
	}

	return term, tail
}

func assertFinitePQEOF(t *testing.T, src core.PQStream, step int) {
	t.Helper()

	_, _, status := src.NextPQ()
	if status != core.StatusEOF {
		t.Fatalf("step %d status = %v, want %v", step, status, core.StatusEOF)
	}
}

func assertExactTakeRange(t *testing.T, got core.Range, want core.Rational, step int) {
	t.Helper()

	if !got.Inside {
		t.Fatalf("step %d Inside = false, want true", step)
	}
	if got.Lo.Open {
		t.Fatalf("step %d Lo.Open = true, want false", step)
	}
	if got.Hi.Open {
		t.Fatalf("step %d Hi.Open = true, want false", step)
	}
	if got.Lo.Value.Cmp(want) != 0 {
		t.Fatalf(
			"step %d Lo = %v/%v, want %v/%v",
			step,
			got.Lo.Value.Num(), got.Lo.Value.Den(),
			want.Num(), want.Den(),
		)
	}
	if got.Hi.Value.Cmp(want) != 0 {
		t.Fatalf(
			"step %d Hi = %v/%v, want %v/%v",
			step,
			got.Hi.Value.Num(), got.Hi.Value.Den(),
			want.Num(), want.Den(),
		)
	}
}

// core/gcf_take_bb_test.go v2
