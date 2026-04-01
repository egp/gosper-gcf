// core/runtime_errors.go v1
package core

import "errors"

var (
	ErrNilReceiver                = errors.New("nil receiver")
	ErrNilObservedSource          = errors.New("nil observed source")
	ErrAppendAfterClose           = errors.New("append after close")
	ErrNoTermAvailableBeforeClose = errors.New("no term available before close")
	ErrNoRemainingSuffixRange     = errors.New("no remaining suffix range")
	ErrUndefinedRangeOnEOFStream  = errors.New("range undefined on EOF stream")
	ErrUndefinedRangeOnBadStream  = errors.New("range undefined on invalid stream")
	ErrInvalidUnaryInputStatus    = errors.New("invalid unary input status")
	ErrInvalidBinaryLeftStatus    = errors.New("invalid binary left-input status")
	ErrInvalidBinaryRightStatus   = errors.New("invalid binary right-input status")
	ErrMissingBinaryState         = errors.New("missing binary evaluator state")
)

// core/runtime_errors.go v1
