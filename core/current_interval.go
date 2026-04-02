// core/current_interval.go v1
package core

import "fmt"

type IntervalReader interface {
	CurrentInterval() (Interval, error)
}

func CurrentIntervalOfPQ(src PQStream) (Interval, error) {
	if src == nil {
		return Interval{}, fmt.Errorf("CurrentIntervalOfPQ: %w", ErrNilReceiver)
	}
	return src.Range()
}

func CurrentIntervalOfRCF(src RCFStream) (Interval, error) {
	if src == nil {
		return Interval{}, fmt.Errorf("CurrentIntervalOfRCF: %w", ErrNilReceiver)
	}
	return src.Range()
}

// core/current_interval.go v1
