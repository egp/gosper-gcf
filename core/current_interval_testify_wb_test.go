// core/current_interval_testify_wb_test.go v1
package core

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWB_CurrentIntervalOfPQ_UsesTestify(t *testing.T) {
	src, status := NewFinitePQStream([]FinitePQStep{
		{
			Term: PQTerm{
				P: big.NewInt(1),
				Q: big.NewInt(1),
			},
			Range: exactRangeFromRational(RationalFromInt64(2)),
		},
	})
	require.Equal(t, StatusOK, status)

	got, err := CurrentIntervalOfPQ(src)
	require.NoError(t, err)
	assert.True(t, got.Inside)
	assert.False(t, got.Lo.Open)
	assert.False(t, got.Hi.Open)
	assert.Zero(t, got.Lo.Value.Cmp(RationalFromInt64(2)))
	assert.Zero(t, got.Hi.Value.Cmp(RationalFromInt64(2)))
}

// core/current_interval_testify_wb_test.go v1
