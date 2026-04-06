// core/pqstream.go v3
package core

import "math/big"

type PQStream interface {
	NextPQ() (PQTerm, PQStream, Status, error)
	CurrentInterval() (Interval, error)
	Range() (Range, error)
}

// core/pqstream.go v3

// --- appended from core/pqterm.go ---

type PQTerm struct {
	P *big.Int
	Q *big.Int
}
