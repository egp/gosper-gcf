package core

import (
	"math/big"
	"testing"
)

func TestWB_Endpoint_OpenClosedFlagsAreIndependent(t *testing.T) {
	value := NewRational(big.NewInt(3), big.NewInt(7))

	openEndpoint := Endpoint{
		Value: value,
		Open:  true,
	}
	closedEndpoint := Endpoint{
		Value: value,
		Open:  false,
	}

	if !openEndpoint.Open {
		t.Fatal("openEndpoint.Open = false, want true")
	}
	if closedEndpoint.Open {
		t.Fatal("closedEndpoint.Open = true, want false")
	}
	if openEndpoint.Value.Cmp(closedEndpoint.Value) != 0 {
		t.Fatal("endpoint values differ, want equal values with independent flags")
	}
}
