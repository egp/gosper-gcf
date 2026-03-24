package core

type PQStream interface {
	NextPQ() (PQTerm, PQStream, Status)
	Range() Range
}
