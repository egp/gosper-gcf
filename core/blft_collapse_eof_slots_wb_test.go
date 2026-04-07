// core/blft_collapse_eof_slots_wb_test.go v1
package core

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWB_BLFT_CollapseBinaryXEOF_RemapsABEFIntoUnarySlots(t *testing.T) {
	s := newBLFTState(BLFTCoefficients{
		A: big.NewInt(2),
		B: big.NewInt(3),
		C: big.NewInt(5),
		D: big.NewInt(7),
		E: big.NewInt(11),
		F: big.NewInt(13),
		G: big.NewInt(17),
		H: big.NewInt(19),
	})

	got, err := exactRationalFromUnaryEngine(
		s.CollapseBinaryXEOF(),
		PQStreamFromRational(RationalFromInt64(4)),
	)
	require.NoError(t, err)

	want := NewRational(big.NewInt(11), big.NewInt(57)) // (2*4 + 3) / (11*4 + 13)
	assert.Zero(t, got.Cmp(want))
}

func TestWB_BLFT_CollapseBinaryYEOF_RemapsACEGIntoUnarySlots(t *testing.T) {
	s := newBLFTState(BLFTCoefficients{
		A: big.NewInt(2),
		B: big.NewInt(3),
		C: big.NewInt(5),
		D: big.NewInt(7),
		E: big.NewInt(11),
		F: big.NewInt(13),
		G: big.NewInt(17),
		H: big.NewInt(19),
	})

	got, err := exactRationalFromUnaryEngine(
		s.CollapseBinaryYEOF(),
		PQStreamFromRational(RationalFromInt64(4)),
	)
	require.NoError(t, err)

	want := NewRational(big.NewInt(13), big.NewInt(61)) // (2*4 + 5) / (11*4 + 17)
	assert.Zero(t, got.Cmp(want))
}

// core/blft_collapse_eof_slots_wb_test.go v1
