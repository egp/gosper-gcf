// core/range_projective_bb_test.go v1
package core

import (
	"errors"
	"testing"
)

func TestBB_NormalizeRange_OutsideClosedEqualBecomesExactInside(t *testing.T) {
	input := Range{
		Lo:     bbClosedInt(7),
		Hi:     bbClosedInt(7),
		Inside: false,
		Kind_:  RangeArc,
	}

	got, err := NormalizeRange(input)
	if err != nil {
		t.Fatalf("NormalizeRange returned unexpected error: %v", err)
	}

	want := Range{
		Lo:     bbClosedInt(7),
		Hi:     bbClosedInt(7),
		Inside: true,
		Kind_:  RangeArc,
	}

	if !bbRangeEqual(got, want) {
		t.Fatalf("NormalizeRange outside closed-equal mismatch\ngot:  %s\nwant: %s", bbFormatRange(got), bbFormatRange(want))
	}
}

func TestBB_NormalizeRange_OutsideOpenEqualReturnsEmptyError(t *testing.T) {
	input := Range{
		Lo:     bbOpenInt(7),
		Hi:     bbOpenInt(7),
		Inside: false,
		Kind_:  RangeArc,
	}

	_, err := NormalizeRange(input)
	if !errors.Is(err, ErrEmptyRange) {
		t.Fatalf("NormalizeRange error = %v, want ErrEmptyRange", err)
	}
}

func TestBB_NormalizeRange_FullCanonicalizesDummyPayload(t *testing.T) {
	input := Range{
		Lo:     bbOpenInt(5),
		Hi:     bbClosedInt(2),
		Inside: false,
		Kind_:  RangeFull,
	}

	got, err := NormalizeRange(input)
	if err != nil {
		t.Fatalf("NormalizeRange returned unexpected error: %v", err)
	}

	want := Range{
		Lo:     bbClosedInt(0),
		Hi:     bbClosedInt(0),
		Inside: true,
		Kind_:  RangeFull,
	}

	if !bbRangeEqual(got, want) {
		t.Fatalf("NormalizeRange full canonicalization mismatch\ngot:  %s\nwant: %s", bbFormatRange(got), bbFormatRange(want))
	}
}

func TestBB_ComplementRange_InsideClosedBecomesOutsideOpen(t *testing.T) {
	input := Range{
		Lo:     bbClosedInt(2),
		Hi:     bbClosedInt(5),
		Inside: true,
		Kind_:  RangeArc,
	}

	got, err := ComplementRange(input)
	if err != nil {
		t.Fatalf("ComplementRange returned unexpected error: %v", err)
	}

	want := Range{
		Lo:     bbOpenInt(5),
		Hi:     bbOpenInt(2),
		Inside: false,
		Kind_:  RangeArc,
	}

	if !bbRangeEqual(got, want) {
		t.Fatalf("ComplementRange inside closed mismatch\ngot:  %s\nwant: %s", bbFormatRange(got), bbFormatRange(want))
	}
}

func TestBB_ComplementRange_InsideOpenBecomesOutsideClosed(t *testing.T) {
	input := Range{
		Lo:     bbOpenInt(2),
		Hi:     bbOpenInt(5),
		Inside: true,
		Kind_:  RangeArc,
	}

	got, err := ComplementRange(input)
	if err != nil {
		t.Fatalf("ComplementRange returned unexpected error: %v", err)
	}

	want := Range{
		Lo:     bbClosedInt(5),
		Hi:     bbClosedInt(2),
		Inside: false,
		Kind_:  RangeArc,
	}

	if !bbRangeEqual(got, want) {
		t.Fatalf("ComplementRange inside open mismatch\ngot:  %s\nwant: %s", bbFormatRange(got), bbFormatRange(want))
	}
}

func TestBB_ComplementRange_OutsideClosedBecomesInsideOpen(t *testing.T) {
	input := Range{
		Lo:     bbClosedInt(5),
		Hi:     bbClosedInt(2),
		Inside: false,
		Kind_:  RangeArc,
	}

	got, err := ComplementRange(input)
	if err != nil {
		t.Fatalf("ComplementRange returned unexpected error: %v", err)
	}

	want := Range{
		Lo:     bbOpenInt(2),
		Hi:     bbOpenInt(5),
		Inside: true,
		Kind_:  RangeArc,
	}

	if !bbRangeEqual(got, want) {
		t.Fatalf("ComplementRange outside closed mismatch\ngot:  %s\nwant: %s", bbFormatRange(got), bbFormatRange(want))
	}
}

func TestBB_ComplementRange_FullReturnsEmptyError(t *testing.T) {
	input := Range{
		Lo:     bbClosedInt(0),
		Hi:     bbClosedInt(0),
		Inside: true,
		Kind_:  RangeFull,
	}

	_, err := ComplementRange(input)
	if !errors.Is(err, ErrEmptyRange) {
		t.Fatalf("ComplementRange error = %v, want ErrEmptyRange", err)
	}
}

func TestBB_RangeContains_InsideClosed(t *testing.T) {
	r := Range{
		Lo:     bbClosedInt(2),
		Hi:     bbClosedInt(5),
		Inside: true,
		Kind_:  RangeArc,
	}

	if !RangeContains(r, bbRatInt(2)) {
		t.Fatal("RangeContains([2,5], 2) = false, want true")
	}
	if !RangeContains(r, bbRatInt(4)) {
		t.Fatal("RangeContains([2,5], 4) = false, want true")
	}
	if !RangeContains(r, bbRatInt(5)) {
		t.Fatal("RangeContains([2,5], 5) = false, want true")
	}
	if RangeContains(r, bbRatInt(1)) {
		t.Fatal("RangeContains([2,5], 1) = true, want false")
	}
	if RangeContains(r, bbRatInt(6)) {
		t.Fatal("RangeContains([2,5], 6) = true, want false")
	}
}

func TestBB_RangeContains_OutsideClosed(t *testing.T) {
	r := Range{
		Lo:     bbClosedInt(5),
		Hi:     bbClosedInt(2),
		Inside: false,
		Kind_:  RangeArc,
	}

	if !RangeContains(r, bbRatInt(1)) {
		t.Fatal("RangeContains(outside[5,2], 1) = false, want true")
	}
	if !RangeContains(r, bbRatInt(2)) {
		t.Fatal("RangeContains(outside[5,2], 2) = false, want true")
	}
	if RangeContains(r, bbRatInt(3)) {
		t.Fatal("RangeContains(outside[5,2], 3) = true, want false")
	}
	if !RangeContains(r, bbRatInt(5)) {
		t.Fatal("RangeContains(outside[5,2], 5) = false, want true")
	}
	if !RangeContains(r, bbRatInt(9)) {
		t.Fatal("RangeContains(outside[5,2], 9) = false, want true")
	}
}

func TestBB_RangeContains_FullAlwaysTrue(t *testing.T) {
	r := Range{
		Lo:     bbClosedInt(0),
		Hi:     bbClosedInt(0),
		Inside: true,
		Kind_:  RangeFull,
	}

	for _, q := range []Rational{bbRatInt(-100), bbRatInt(0), bbRatInt(100)} {
		if !RangeContains(r, q) {
			t.Fatalf("RangeContains(full, %s) = false, want true", bbFormatRational(q))
		}
	}
}

func bbClosedInt(n int64) Endpoint {
	return Endpoint{
		Value: RationalFromInt64(n),
		Open:  false,
	}
}

func bbOpenInt(n int64) Endpoint {
	return Endpoint{
		Value: RationalFromInt64(n),
		Open:  true,
	}
}

func bbRatInt(n int64) Rational {
	return RationalFromInt64(n)
}

func bbRangeEqual(a, b Range) bool {
	return a.Kind() == b.Kind() &&
		a.Inside == b.Inside &&
		a.Lo.Open == b.Lo.Open &&
		a.Hi.Open == b.Hi.Open &&
		a.Lo.Value.Cmp(b.Lo.Value) == 0 &&
		a.Hi.Value.Cmp(b.Hi.Value) == 0
}

func bbFormatRange(r Range) string {
	return "Range{" +
		"Kind=" + bbFormatKind(r.Kind()) +
		", Lo=" + bbFormatEndpoint(r.Lo) +
		", Hi=" + bbFormatEndpoint(r.Hi) +
		", Inside=" + bbFormatBool(r.Inside) +
		"}"
}

func bbFormatEndpoint(ep Endpoint) string {
	return "Endpoint{Value=" + bbFormatRational(ep.Value) + ", Open=" + bbFormatBool(ep.Open) + "}"
}

func bbFormatRational(r Rational) string {
	return r.Num().String() + "/" + r.Den().String()
}

func bbFormatKind(k RangeKind) string {
	switch k {
	case RangeArc:
		return "RangeArc"
	case RangeFull:
		return "RangeFull"
	default:
		return "RangeKind(?)"
	}
}

func bbFormatBool(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

func TestBB_ComplementRange_HalfOpenInsideHasExpectedBoundaryMembership(t *testing.T) {
	input := Range{
		Lo:     bbClosedInt(2),
		Hi:     bbOpenInt(5),
		Inside: true,
		Kind_:  RangeArc,
	}

	got, err := ComplementRange(input)
	if err != nil {
		t.Fatalf("ComplementRange returned unexpected error: %v", err)
	}

	if RangeContains(got, bbRatInt(2)) {
		t.Fatal("RangeContains(complement([2,5)), 2) = true, want false")
	}
	if RangeContains(got, bbRatInt(3)) {
		t.Fatal("RangeContains(complement([2,5)), 3) = true, want false")
	}
	if !RangeContains(got, bbRatInt(5)) {
		t.Fatal("RangeContains(complement([2,5)), 5) = false, want true")
	}
	if !RangeContains(got, bbRatInt(9)) {
		t.Fatal("RangeContains(complement([2,5)), 9) = false, want true")
	}
}

// core/range_projective_bb_test.go v1
