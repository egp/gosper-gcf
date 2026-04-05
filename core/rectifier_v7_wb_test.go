// core/rectifier_v7_wb_test.go V1

package core
// Version: 1

package core

import (
	"math/big"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestWB_Rectifier_V7_SpeculativeIngestion(t *testing.T) {
	// Identity LFT [1, 0; 0, 1]
	r := NewRectifier(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1))

	// 1. Absorb a term (3, 1). Matrix becomes [3, 1; 1, 0]
	// z(inf) = 3/1 = 3. z(1) = (3+1)/(1+0) = 4. 
	// Floor is not stable (3 != 4).
	r.Absorb(PQTerm{P: big.NewInt(3), Q: big.NewInt(1)})
	can, _ := r.CanEmit()
	assert.False(t, can, "Should not emit when floor is unstable [3, 4]")

	// 2. Absorb next term (2, 1). Matrix updates.
	// New A = 3*2 + 1*1 = 7. New B = 3.
	// New C = 1*2 + 0*1 = 2. New D = 1.
	// z(inf) = 7/2 = 3. z(1) = (7+3)/(2+1) = 10/3 = 3.
	// Floor is stable (3 == 3).
	r.Absorb(PQTerm{P: big.NewInt(2), Q: big.NewInt(1)})
	can, term := r.CanEmit()
	assert.True(t, can, "Should emit when floor stabilizes")
	assert.Equal(t, int64(3), term.Int64())

	// 3. Produce (Emit) the term 3.
	r.Emit(term)
	// After Emit(3), verify matrix state for next iteration
	// a'' = c = 2, b'' = d = 1
	// c'' = a - 3c = 7 - 6 = 1, d'' = b - 3d = 3 - 3 = 0
	// Matrix is [2, 1; 1, 0]
	assert.Equal(t, int64(2), r.a.Int64())
	assert.Equal(t, int64(1), r.b.Int64())
}

// core/rectifier_v7_wb_test.go V1