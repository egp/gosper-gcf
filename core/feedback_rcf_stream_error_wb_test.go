// core/feedback_rcf_stream_error_wb_test.go v3
package core

import (
	"errors"
	"math/big"
	"testing"
)

func TestWB_FeedbackRCFStream_RangeAfterAllTermsConsumedReturnsError(t *testing.T) {
	s := newFeedbackRCFStream()
	if err := s.Append(NewRCFTerm(big.NewInt(2)), exactRangeFromRational(RationalFromInt64(2))); err != nil {
		t.Fatalf("Append error = %v", err)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close error = %v", err)
	}

	_, status, err := s.NextRCF()
	if err != nil {
		t.Fatalf("NextRCF error = %v", err)
	}
	if status != StatusOK {
		t.Fatalf("status = %v, want %v", status, StatusOK)
	}

	_, err = s.Range()
	if !errors.Is(err, ErrNoRemainingSuffixRange) {
		t.Fatalf("Range error = %v, want ErrNoRemainingSuffixRange", err)
	}
}

// core/feedback_rcf_stream_error_wb_test.go v3
