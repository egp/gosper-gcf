// core/rcfstream.go v2
package core

type RCFStream interface {
	NextRCF() (RCFTerm, Status, error)
	Range() (Range, error)
}

// core/rcfstream.go v2
