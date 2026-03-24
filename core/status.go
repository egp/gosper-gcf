package core

type Status int

const (
	StatusOK Status = iota
	StatusEOF
	StatusInvalidInput
)
