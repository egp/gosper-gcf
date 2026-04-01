// core/pqstream.go v2
package core

type PQStream interface {
	NextPQ() (PQTerm, PQStream, Status, error)
	Range() (Range, error)
}

// core/pqstream.go v2
