package core

import "math/big"

type RCFTerm struct {
	a *big.Int
}

func NewRCFTerm(a *big.Int) RCFTerm {
	if a == nil {
		return RCFTerm{a: big.NewInt(0)}
	}
	return RCFTerm{a: new(big.Int).Set(a)}
}

func (t RCFTerm) A() *big.Int {
	if t.a == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Set(t.a)
}
