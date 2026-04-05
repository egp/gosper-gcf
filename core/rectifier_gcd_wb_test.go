// core/rectifier_gcd_wb_test.go V2
package core

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWB_Rectifier_Absorb_GCD_Normalization_ReducesCoefficients(t *testing.T) {
	r := NewRectifier(big.NewInt(2), big.NewInt(4), big.NewInt(6), big.NewInt(8))
	pq := PQTerm{P: big.NewInt(3), Q: big.NewInt(6)}
	r = r.Absorb(pq)

	// chained GCD (exactly like blft_normalize.go)
	g := new(big.Int).GCD(nil, nil, r.a, r.b)
	g = new(big.Int).GCD(g, nil, g, r.c)
	g = new(big.Int).GCD(g, nil, g, r.d)
	assert.Equal(t, int64(1), g.Int64(), "final GCD is 1 after normalization")
	assert.Equal(t, int64(1), r.a.Int64(), "a reduced")
	assert.Equal(t, int64(2), r.b.Int64(), "b reduced")
	assert.Equal(t, int64(3), r.c.Int64(), "c reduced")
	assert.Equal(t, int64(4), r.d.Int64(), "d reduced")
}

// core/rectifier_gcd_wb_test.go V2
