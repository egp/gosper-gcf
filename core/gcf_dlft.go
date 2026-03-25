// core/gcf_dlft.go v1
package core

func NewDLFT1(coeffs DLFTCoefficients, x PQStream) *GCF {
	return newDLFT1WithResolvedConfig(coeffs, x, DefaultConfig())
}

func NewDLFT1WithConfig(coeffs DLFTCoefficients, x PQStream, cfg Config) *GCF {
	return newDLFT1WithResolvedConfig(coeffs, x, cfg)
}

// Shape-only stub.
// Later this will create a live unary DLFT-backed GCF evaluator.
func newDLFT1WithResolvedConfig(coeffs DLFTCoefficients, x PQStream, cfg Config) *GCF {
	_ = coeffs
	_ = x

	return NewExactTerminalGCFWithConfig(
		nil,
		exactRangeFromRational(RationalFromInt64(0)),
		cfg,
	)
}

// core/gcf_dlft.go v1
