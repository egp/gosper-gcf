// core/rcfstream.go v3
package core

type RCFStream interface {
	NextRCF() (RCFTerm, Status, error)
	CurrentInterval() (Interval, error)
	Range() (Range, error)
}

// core/rcfstream.go v3
