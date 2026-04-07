// core/gcf_error_channel_wb_test.go v1
package core

import (
	"errors"
	"testing"
)

func TestWB_GCF_NilReceiverReturnsErrorsNotPanics(t *testing.T) {
	var g *GCF

	_, _, err := g.NextRCF()
	if !errors.Is(err, ErrNilReceiver) {
		t.Fatalf("NextRCF error = %v, want ErrNilReceiver", err)
	}

	_, err = g.Range()
	if !errors.Is(err, ErrNilReceiver) {
		t.Fatalf("Range error = %v, want ErrNilReceiver", err)
	}
}

func TestWB_NewObservedRCFGCFWithResolvedConfig_NilSourceReturnsErrorStream(t *testing.T) {
	g := newObservedRCFGCFWithResolvedConfig(nil, DefaultConfig())
	if g == nil {
		t.Fatal("newObservedRCFGCFWithResolvedConfig(nil, cfg) returned nil")
	}

	_, _, err := g.NextRCF()
	if !errors.Is(err, ErrNilObservedSource) {
		t.Fatalf("NextRCF error = %v, want ErrNilObservedSource", err)
	}
}

// core/gcf_error_channel_wb_test.go v1
