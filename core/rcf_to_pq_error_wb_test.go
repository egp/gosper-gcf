// core/rcf_to_pq_error_wb_test.go v1
package core

import (
	"errors"
	"testing"
)

func TestWB_PQStreamFromRCF_NilSourceReturnsErrorPQStream(t *testing.T) {
	pq := PQStreamFromRCF(nil)

	_, _, _, err := pq.NextPQ()
	if !errors.Is(err, ErrNilObservedSource) {
		t.Fatalf("NextPQ error = %v, want ErrNilObservedSource", err)
	}

	_, err = pq.Range()
	if !errors.Is(err, ErrNilObservedSource) {
		t.Fatalf("Range error = %v, want ErrNilObservedSource", err)
	}
}

// core/rcf_to_pq_error_wb_test.go v1
