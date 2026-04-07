// core/bigint_helpers_wb_test.go v1
package core

import (
	"math/big"
	"testing"
)

func TestWB_CloneBigInt_NilAndDeepCopy(t *testing.T) {
	if cloneBigInt(nil) != nil {
		t.Fatal("cloneBigInt(nil) != nil")
	}

	src := big.NewInt(123)
	got := cloneBigInt(src)
	if got == nil {
		t.Fatal("cloneBigInt(src) returned nil")
	}
	if got == src {
		t.Fatal("cloneBigInt(src) returned same pointer, want deep copy")
	}
	if got.Cmp(src) != 0 {
		t.Fatalf("cloneBigInt(src) = %s, want %s", got.String(), src.String())
	}

	src.SetInt64(999)
	if got.Int64() != 123 {
		t.Fatalf("clone changed after mutating source: got %d, want 123", got.Int64())
	}
}

// core/bigint_helpers_wb_test.go v1
