// core/feedback_rcf_stream_wb_test.go v3
package core

import (
	"errors"
	"math/big"
	"testing"
)

func TestWB_FeedbackRCFStream_NextRCFBeforeCloseReturnsError(t *testing.T) {
	s := newFeedbackRCFStream()

	_, _, err := s.NextRCF()
	if !errors.Is(err, ErrNoTermAvailableBeforeClose) {
		t.Fatalf("NextRCF error = %v, want ErrNoTermAvailableBeforeClose", err)
	}
}

func TestWB_FeedbackRCFStream_AppendAfterCloseReturnsError(t *testing.T) {
	s := newFeedbackRCFStream()
	if err := s.Close(); err != nil {
		t.Fatalf("Close error = %v", err)
	}

	err := s.Append(NewRCFTerm(big.NewInt(1)), exactRangeFromRational(RationalFromInt64(1)))
	if !errors.Is(err, ErrAppendAfterClose) {
		t.Fatalf("Append error = %v, want ErrAppendAfterClose", err)
	}
}

// core/feedback_rcf_stream_wb_test.go v3
