// core/pqstream_error_channel_wb_test.go v1
package core

import (
	"errors"
	"testing"
)

func TestWB_EOFPQStream_RangeReturnsErrorInsteadOfPanicking(t *testing.T) {
	_, err := finitePQEOF.Range()
	if !errors.Is(err, ErrUndefinedRangeOnEOFStream) {
		t.Fatalf("Range error = %v, want ErrUndefinedRangeOnEOFStream", err)
	}
}

// core/pqstream_error_channel_wb_test.go v1
