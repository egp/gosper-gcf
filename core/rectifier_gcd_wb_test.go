// core/rectifier_gcd_wb_test.go V4
package core

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWB_Rectifier_Absorb_GCD_Normalization_ReducesCoefficients(t *testing.T) {
	// Start with (2,4,6,8). Absorb {P=3,Q=6}:
	//   newA = 2*3 + 4*6 = 30, newC = 6*3 + 8*6 = 66
	//   shift: b=2, a=30, d=6, c=66 → (30,2,66,6)
	//   normalize GCD(30,2,66,6)=2 → (15,1,33,3)
	r := NewRectifier(big.NewInt(2), big.NewInt(4), big.NewInt(6), big.NewInt(8))
	pq := PQTerm{P: big.NewInt(3), Q: big.NewInt(6)}
	r.Absorb(pq)

	// The overall GCD of normalized coefficients must be 1.
	g := new(big.Int).GCD(nil, nil, r.a, r.b)
	g = new(big.Int).GCD(g, nil, g, r.c)
	g = new(big.Int).GCD(g, nil, g, r.d)
	assert.Equal(t, int64(1), g.Int64(), "final GCD is 1 after normalization")
	assert.Equal(t, int64(15), r.a.Int64(), "a reduced")
	assert.Equal(t, int64(1), r.b.Int64(), "b reduced")
	assert.Equal(t, int64(33), r.c.Int64(), "c reduced")
	assert.Equal(t, int64(3), r.d.Int64(), "d reduced")
}

// core/rectifier_gcd_wb_test.go V4
