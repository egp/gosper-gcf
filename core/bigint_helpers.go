// core/bigint_helpers.go v1
package core

import "math/big"

func cloneBigInt(x *big.Int) *big.Int {
	if x == nil {
		return nil
	}
	return new(big.Int).Set(x)
}

// core/bigint_helpers.go v1
