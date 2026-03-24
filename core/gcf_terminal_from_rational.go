// core/gcf_terminal_from_rational.go v2
package core

func newExactTerminalGCFFromRational(r Rational) *GCF {
	return NewExactTerminalGCF(
		rcfTermsFromRational(r),
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

// core/gcf_terminal_from_rational.go v2
