// core/rectifier_wb_test.go V3
package core

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

// After absorbing one integer term p into the identity rectifier, the output
// interval is [p, p+1), so CanEmit requires two terms to narrow the floor.
// Absorb [5;2] (= 11/2): z(inf)=11/2=5.5 → floor 5, z(1)=16/3=5.33 → floor 5.
func TestWB_Rectifier_IdentityAbsorb_CanEmit_ReturnsCorrectFloor(t *testing.T) {
	r := NewRectifier(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1))
	r.Absorb(PQTerm{P: big.NewInt(5), Q: big.NewInt(1)})
	r.Absorb(PQTerm{P: big.NewInt(2), Q: big.NewInt(1)})
	ok, term := r.CanEmit()
	assert.True(t, ok, "expected CanEmit to succeed after absorbing [5;2]")
	assert.Equal(t, int64(5), term.Int64(), "floor of 11/2 is 5")
}

// Absorb [3;4] (= 13/4 = 3.25) into identity; both z(inf)=3.25 and z(1)=3.2
// floor to 3. Emit(3) performs the production substitution; the post-emit state
// (4,1,1,0) represents the remaining partial-quotient window.
func TestWB_Rectifier_Emit_ProductionSubstitution_UpdatesState(t *testing.T) {
	r := NewRectifier(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1))
	r.Absorb(PQTerm{P: big.NewInt(3), Q: big.NewInt(1)}) // → (3,1,1,0)
	r.Absorb(PQTerm{P: big.NewInt(4), Q: big.NewInt(1)}) // → (13,3,4,1)

	ok, term := r.CanEmit()
	assert.True(t, ok)
	assert.Equal(t, int64(3), term.Int64())

	r.Emit(term) // state → (4,1,1,0)

	// Verify post-emit state by absorbing one more term and checking CanEmit.
	// Absorb {P=10,Q=1}: (41,4,10,1); z(inf)=4.1, z(1)=45/11=4.09 → floor 4.
	r.Absorb(PQTerm{P: big.NewInt(10), Q: big.NewInt(1)})
	ok2, term2 := r.CanEmit()
	assert.True(t, ok2, "expected CanEmit after emit + one more absorb")
	assert.Equal(t, int64(4), term2.Int64(), "next floor is 4")
}

// core/rectifier_wb_test.go V3
