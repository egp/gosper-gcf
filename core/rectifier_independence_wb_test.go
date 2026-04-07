// core/rectifier_independence_wb_test.go V3
package core

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

// After absorbing the two-term CF [3;7] (= 22/7), the rectifier's interval
// [z(inf), z(1)] = [22/7, 25/8] = [3.14…, 3.125], both floor to 3, so
// CanEmit must succeed and return 3. After Emit(3) the state shifts to the
// next partial-quotient window.
func TestWB_Rectifier_AbsorbTwoTerms_CanEmit_ThenEmit_UpdatesState(t *testing.T) {
	r := NewRectifier(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1))

	r.Absorb(PQTerm{P: big.NewInt(3), Q: big.NewInt(1)}) // state becomes (3,1,1,0)
	r.Absorb(PQTerm{P: big.NewInt(7), Q: big.NewInt(1)}) // state becomes (22,3,7,1)

	ok, term := r.CanEmit()
	assert.True(t, ok, "expected CanEmit to succeed after absorbing [3;7]")
	assert.Equal(t, int64(3), term.Int64(), "floor of 22/7 is 3")

	r.Emit(term)

	// After emitting 3, state is (7,1,1,0): z(inf)=7, z(1)=8 → can't yet emit
	// without more input; verify CanEmit is false (not a panic).
	ok2, _ := r.CanEmit()
	assert.False(t, ok2, "expected CanEmit to be false without more input after emit")
}

// core/rectifier_independence_wb_test.go V3
