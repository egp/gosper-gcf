// core/gcf_constructor_independent_x_wb_test.go v3
package core

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWB_NewGCF2_IndependentOfX_ProjectY_FirstTermMatchesRightInput(t *testing.T) {
	g := NewGCF2(
		BLFTCoefficients{
			A: big.NewInt(0),
			B: big.NewInt(0),
			C: big.NewInt(1),
			D: big.NewInt(0),
			E: big.NewInt(0),
			F: big.NewInt(0),
			G: big.NewInt(0),
			H: big.NewInt(1),
		},
		PQStreamFromRational(RationalFromInt64(5)),
		PQStreamFromRational(NewRational(big.NewInt(22), big.NewInt(7))),
	)

	term, status, err := g.NextRCF()
	require.NoError(t, err)
	require.Equal(t, StatusOK, status)
	assert.Zero(t, term.A().Cmp(big.NewInt(3)))
}

func TestWB_NewGCF2_IndependentOfX_IdentityProjection_PreservesBinaryEvaluatorLazily(t *testing.T) {
	g := NewGCF2(
		BLFTCoefficients{
			A: big.NewInt(0),
			B: big.NewInt(0),
			C: big.NewInt(1),
			D: big.NewInt(0),
			E: big.NewInt(0),
			F: big.NewInt(0),
			G: big.NewInt(0),
			H: big.NewInt(1),
		},
		PQStreamFromRational(RationalFromInt64(5)),
		PQStreamFromRational(NewRational(big.NewInt(22), big.NewInt(7))),
	)

	require.NotNil(t, g)
	assert.Nil(t, g.stream)
	assert.Nil(t, g.unary)
	require.NotNil(t, g.binary)
}

func TestWB_NewGCF2_IndependentOfX_AffineProjection_PreservesBinaryEvaluatorLazily(t *testing.T) {
	g := NewGCF2(
		BLFTCoefficients{
			A: big.NewInt(0),
			B: big.NewInt(0),
			C: big.NewInt(2),
			D: big.NewInt(1),
			E: big.NewInt(0),
			F: big.NewInt(0),
			G: big.NewInt(0),
			H: big.NewInt(1),
		},
		PQStreamFromRational(RationalFromInt64(5)),
		PQStreamFromRational(RationalFromInt64(3)),
	)

	require.NotNil(t, g)
	assert.Nil(t, g.stream)
	assert.Nil(t, g.unary)
	require.NotNil(t, g.binary)
}

// core/gcf_constructor_independent_x_wb_test.go v3
