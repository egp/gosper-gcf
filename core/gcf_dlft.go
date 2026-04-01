// core/gcf_dlft.go v4
package core

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
			engine: newDLFTState(coeffs),
			x:      x,
		}
	}
	return g
}

// core/gcf_dlft.go v4
