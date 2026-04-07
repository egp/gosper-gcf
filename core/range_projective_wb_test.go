// core/range_projective_wb_test.go v1
package core

import (
	"errors"
	"testing"
)

func TestWB_NormalizeRange_RejectsMalformedInsideDescendingArc(t *testing.T) {
	input := Range{
		Lo:     wbClosedInt(5),
		Hi:     wbClosedInt(2),
		Inside: true,
		Kind_:  RangeArc,
	}

	_, err := NormalizeRange(input)
	if !errors.Is(err, ErrNonCanonicalRange) {
		t.Fatalf("NormalizeRange error = %v, want ErrNonCanonicalRange", err)
	}
}

func TestWB_NormalizeRange_ExactInsideClosedEqualRemainsCanonical(t *testing.T) {
	input := Range{
		Lo:     wbClosedInt(3),
		Hi:     wbClosedInt(3),
		Inside: true,
		Kind_:  RangeArc,
	}

	got, err := NormalizeRange(input)
	if err != nil {
		t.Fatalf("NormalizeRange returned unexpected error: %v", err)
	}

	if got.Kind() != RangeArc {
		t.Fatalf("got.Kind() = %v, want RangeArc", got.Kind())
	}
	if !got.Inside {
		t.Fatal("got.Inside = false, want true")
	}
	if got.Lo.Open || got.Hi.Open {
		t.Fatalf("got exact point openness = (%t,%t), want (false,false)", got.Lo.Open, got.Hi.Open)
	}
	if got.Lo.Value.Cmp(bbRatInt(3)) != 0 || got.Hi.Value.Cmp(bbRatInt(3)) != 0 {
		t.Fatalf("got exact point values = (%s,%s), want (3/1,3/1)", bbFormatRational(got.Lo.Value), bbFormatRational(got.Hi.Value))
	}
}

func TestWB_NormalizeRange_FullAlwaysUsesCanonicalDummyPayload(t *testing.T) {
	input := Range{
		Lo:     wbOpenInt(99),
		Hi:     wbOpenInt(-17),
		Inside: false,
		Kind_:  RangeFull,
	}

	got, err := NormalizeRange(input)
	if err != nil {
		t.Fatalf("NormalizeRange returned unexpected error: %v", err)
	}

	if got.Kind() != RangeFull {
		t.Fatalf("got.Kind() = %v, want RangeFull", got.Kind())
	}
	if !got.Inside {
		t.Fatal("got.Inside = false, want true for canonical full dummy payload")
	}
	if got.Lo.Open || got.Hi.Open {
		t.Fatalf("canonical full openness = (%t,%t), want (false,false)", got.Lo.Open, got.Hi.Open)
	}
	if got.Lo.Value.Cmp(bbRatInt(0)) != 0 || got.Hi.Value.Cmp(bbRatInt(0)) != 0 {
		t.Fatalf("canonical full values = (%s,%s), want (0/1,0/1)", bbFormatRational(got.Lo.Value), bbFormatRational(got.Hi.Value))
	}
}
func TestWB_ComplementRange_TogglesEndpointOpennessExactly(t *testing.T) {
	input := Range{
		Lo:     wbClosedInt(2),
		Hi:     wbOpenInt(5),
		Inside: true,
		Kind_:  RangeArc,
	}

	got, err := ComplementRange(input)
	if err != nil {
		t.Fatalf("ComplementRange returned unexpected error: %v", err)
	}

	if got.Kind() != RangeArc {
		t.Fatalf("got.Kind() = %v, want RangeArc", got.Kind())
	}
	if got.Inside {
		t.Fatal("got.Inside = true, want false")
	}
	if got.Lo.Open {
		t.Fatal("got.Lo.Open = true, want false")
	}
	if !got.Hi.Open {
		t.Fatal("got.Hi.Open = false, want true")
	}
	if got.Lo.Value.Cmp(bbRatInt(5)) != 0 || got.Hi.Value.Cmp(bbRatInt(2)) != 0 {
		t.Fatalf("complement values = (%s,%s), want (5/1,2/1)", bbFormatRational(got.Lo.Value), bbFormatRational(got.Hi.Value))
	}
}

func wbClosedInt(n int64) Endpoint {
	return Endpoint{
		Value: RationalFromInt64(n),
		Open:  false,
	}
}

func wbOpenInt(n int64) Endpoint {
	return Endpoint{
		Value: RationalFromInt64(n),
		Open:  true,
	}
}

// core/range_projective_wb_test.go v1
