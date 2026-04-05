// core/rectifier_wb_test.go V1
package core

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWB_Rectifier_IdentityAbsorb_CanEmit_ReturnsCorrectFloor(t *testing.T) {
	r := NewRectifier(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1))
	pq := PQTerm{P: big.NewInt(5), Q: big.NewInt(1)}
	r = r.Absorb(pq)
	term, ok := r.CanEmit()
	assert.True(t, ok, "expected CanEmit to succeed after absorb")
	assert.Equal(t, int64(5), term.A().Int64(), "floor should be the absorbed integer term")
	assert.Equal(t, big.NewInt(1), term.A(), "RCFTerm value exact")
}

func TestWB_Rectifier_Emit_ProductionSubstitution_UpdatesState(t *testing.T) {
	r := NewRectifier(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1))
	pq := PQTerm{P: big.NewInt(3), Q: big.NewInt(1)}
	r = r.Absorb(pq)
	term, _ := r.CanEmit()
	r = r.Emit(term)
	// After emit the state should be ready for next term (z' = 1/(z - 3))
	nextTerm, ok := r.CanEmit()
	assert.True(t, ok)
	assert.Equal(t, int64(0), nextTerm.A().Int64(), "post-emit state starts at 0 for next floor")
}

// core/rectifier_wb_test.go V1
