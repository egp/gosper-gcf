// core/current_interval.go v2
package core

import "fmt"

type IntervalReader interface {
	CurrentInterval() (Interval, error)
}

func CurrentIntervalOfPQ(src PQStream) (Interval, error) {
	if src == nil {
		return Interval{}, fmt.Errorf("CurrentIntervalOfPQ: %w", ErrNilReceiver)
	}
	return src.CurrentInterval()
}

func CurrentIntervalOfRCF(src RCFStream) (Interval, error) {
	if src == nil {
		return Interval{}, fmt.Errorf("CurrentIntervalOfRCF: %w", ErrNilReceiver)
	}
	return src.CurrentInterval()
}

// core/current_interval.go v2
