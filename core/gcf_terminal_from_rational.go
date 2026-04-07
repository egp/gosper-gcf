// core/gcf_terminal_from_rational.go v3
package core

import "fmt"

// NewExactTerminalGCFFromRational builds an exact finite GCF whose RCF
// emission is exactly the canonical continued-fraction expansion of r.
func NewExactTerminalGCFFromRational(r Rational) *GCF {
	return newExactTerminalGCFFromRational(r)
}

func newExactTerminalGCFFromRational(r Rational) *GCF {
	terms, err := rcfTermsFromRationalChecked(r)
	if err != nil {
		return &GCF{
			cfg:    DefaultConfig(),
			stream: newErrorRCFStream(fmt.Errorf("newExactTerminalGCFFromRational: %w", err)),
		}
	}

	return NewExactTerminalGCF(
		terms,
		Range{
			Lo: Endpoint{
				Value: r,
				Open:  false,
			},
			Hi: Endpoint{
				Value: r,
				Open:  false,
			},
			Inside: true,
		},
	)
}

// core/gcf_terminal_from_rational.go v3
