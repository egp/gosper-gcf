// core/gcf_no_observe_wb_test.go v1
package core

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWB_GCF_UnaryIdentity_GeneralizedInputDoesNotRequireObservedPQRCFShortcut(t *testing.T) {
	stream, status := NewFinitePQStream([]FinitePQStep{
		{
			Term: PQTerm{P: big.NewInt(1), Q: big.NewInt(2)},
			Range: Range{
				Lo:     Endpoint{Value: RationalFromInt64(1), Open: false},
				Hi:     Endpoint{Value: RationalFromInt64(2), Open: false},
				Inside: true,
			},
		},
		{
			Term: PQTerm{P: big.NewInt(3), Q: big.NewInt(1)},
			Range: Range{
				Lo:     Endpoint{Value: RationalFromInt64(3), Open: false},
				Hi:     Endpoint{Value: RationalFromInt64(4), Open: false},
				Inside: true,
			},
		},
	})
	require.Equal(t, StatusOK, status)

	g := NewGCF1(
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
		stream,
	)

	term1, status1, err1 := g.NextRCF()
	require.NoError(t, err1)
	require.Equal(t, StatusOK, status1)
	assert.Zero(t, term1.A().Cmp(big.NewInt(1)))

	term2, status2, err2 := g.NextRCF()
	require.NoError(t, err2)
	require.Equal(t, StatusOK, status2)
	assert.Zero(t, term2.A().Cmp(big.NewInt(1)))

	term3, status3, err3 := g.NextRCF()
	require.NoError(t, err3)
	require.Equal(t, StatusOK, status3)
	assert.Zero(t, term3.A().Cmp(big.NewInt(2)))
}

// core/gcf_no_observe_wb_test.go v1
