// core/rcfstream.go v1
package core

type RCFStream interface {
	NextRCF() (RCFTerm, Status)
	Range() Range
}

// core/rcfstream.go v1
