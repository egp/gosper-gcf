// core/rectifier_independence_wb_test.go V1
package core

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWB_Rectifier_IndependenceTransition_AfterEmit_CorrectlySwitchesToUnary(t *testing.T) {
	// Simulate a binary GCF that becomes independent of X after emit
	r := NewRectifier(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1))
	pq := PQTerm{P: big.NewInt(5), Q: big.NewInt(1)}
	r = r.Absorb(pq)
	term, ok := r.CanEmit()
	assert.True(t, ok)
	r = r.Emit(term)

	// After independence transition the rectifier state should be ready for unary continuation
	nextTerm, nextOk := r.CanEmit()
	assert.True(t, nextOk)
	assert.Equal(t, int64(0), nextTerm.A().Int64(), "post-transition floor is 0 as expected for continued unary")
}

// core/rectifier_independence_wb_test.go V1
