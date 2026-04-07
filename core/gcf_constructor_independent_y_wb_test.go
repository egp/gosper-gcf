// core/gcf_constructor_independent_y_wb_test.go v2
package core

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWB_NewGCF2_IndependentOfY_IdentityProjection_PreservesBinaryEvaluatorLazily(t *testing.T) {
	g := NewGCF2(
		BLFTCoefficients{
			A: big.NewInt(0),
			B: big.NewInt(1),
			C: big.NewInt(0),
			D: big.NewInt(0),
			E: big.NewInt(0),
			F: big.NewInt(0),
			G: big.NewInt(0),
			H: big.NewInt(1),
		},
		PQStreamFromRational(NewRational(big.NewInt(22), big.NewInt(7))),
		PQStreamFromRational(RationalFromInt64(5)),
	)

	require.NotNil(t, g)
	assert.Nil(t, g.stream)
	assert.Nil(t, g.unary)
	require.NotNil(t, g.binary)
}

func TestWB_NewGCF2_IndependentOfY_AffineProjection_PreservesBinaryEvaluatorLazily(t *testing.T) {
	g := NewGCF2(
		BLFTCoefficients{
			A: big.NewInt(0),
			B: big.NewInt(2),
			C: big.NewInt(0),
			D: big.NewInt(1),
			E: big.NewInt(0),
			F: big.NewInt(0),
			G: big.NewInt(0),
			H: big.NewInt(1),
		},
		PQStreamFromRational(RationalFromInt64(3)),
		PQStreamFromRational(RationalFromInt64(5)),
	)

	require.NotNil(t, g)
	assert.Nil(t, g.stream)
	assert.Nil(t, g.unary)
	require.NotNil(t, g.binary)
}

// core/gcf_constructor_independent_y_wb_test.go v2
