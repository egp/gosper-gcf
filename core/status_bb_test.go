package core_test

import (
	"testing"

	"github.com/egp/gosper-gcf/core"
)

func TestBB_Status_PublicValuesAreStable(t *testing.T) {
	if core.StatusOK != 0 {
		t.Fatalf("StatusOK = %d, want 0", core.StatusOK)
	}
	if core.StatusEOF != 1 {
		t.Fatalf("StatusEOF = %d, want 1", core.StatusEOF)
	}
	if core.StatusInvalidInput != 2 {
		t.Fatalf("StatusInvalidInput = %d, want 2", core.StatusInvalidInput)
	}
}
