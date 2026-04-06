// core/config.go v1
package core

import (
	"fmt"
	"math/big"
	"time"
)

type Config struct {
	Timeout             time.Duration
	BitLenLimit         int
	BitLenWarnThreshold int
}

func DefaultConfig() Config {
	return Config{}
}

// core/config.go v1

// --- appended from core/current_interval.go ---
// core/current_interval.go v2

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

// --- appended from core/engine.go ---
// core/engine.go v2

type unaryEngine interface {
	UnaryRange(xRange Range) (Range, error)
	CanEmitRCFTerm(r Range) (RCFTerm, bool)
	EmitUnary(term RCFTerm) unaryEngine
	IngestUnaryX(term PQTerm) unaryEngine
	CollapseUnaryEOF() Rational
}

type binaryEngine interface {
	BinaryRange(xRange, yRange Range) (Range, error)
	CanEmitRCFTerm(r Range) (RCFTerm, bool)
	EmitBinary(term RCFTerm) binaryEngine
	IngestBinaryX(term PQTerm) binaryEngine
	IngestBinaryY(term PQTerm) binaryEngine
	CollapseBinaryXEOF() unaryEngine
	CollapseBinaryYEOF() unaryEngine
	CollapseBinaryBothEOF() Rational
	IndependentOfX() bool
	IndependentOfY() bool
}

// core/engine.go v2

// --- appended from core/gcf_dlft.go ---
// core/gcf_dlft.go v4

func NewDLFT1(coeffs DLFTCoefficients, x PQStream) *GCF {
	return newDLFT1WithResolvedConfig(coeffs, x, DefaultConfig())
}

func NewDLFT1WithConfig(coeffs DLFTCoefficients, x PQStream, cfg Config) *GCF {
	return newDLFT1WithResolvedConfig(coeffs, x, cfg)
}

func newDLFT1WithResolvedConfig(coeffs DLFTCoefficients, x PQStream, cfg Config) *GCF {
	g := &GCF{
		cfg: cfg,
		x:   x,
	}
	if x != nil {
		g.unary = &unaryEvaluatorState{
			engine:    newDLFTState(coeffs),
			rectifier: NewRectifier(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1)),
			x:         x,
		}
	}
	return g
}

// core/gcf_dlft.go v4
