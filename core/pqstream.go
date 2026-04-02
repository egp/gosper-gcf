// core/pqstream.go v3
package core

type PQStream interface {
	NextPQ() (PQTerm, PQStream, Status, error)
	CurrentInterval() (Interval, error)
	Range() (Range, error)
}

// core/pqstream.go v3
