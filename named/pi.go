// named/pi.go v3
package named

import (
	"math/big"

	"github.com/egp/gosper-gcf/core"
)

func Pi() core.PQStream {
	return mustFiniteRCFPrefixStream([]int64{3, 7, 15, 1, 292, 1, 1, 1, 2})
}

func mustFiniteRCFPrefixStream(terms []int64) core.PQStream {
	if len(terms) == 0 {
		panic("named: finite RCF prefix stream requires at least one term")
	}

	steps := make([]core.FinitePQStep, len(terms))
	suffixes := suffixRationalsFromRCFTerms(terms)

	for i, term := range terms {
		steps[i] = core.FinitePQStep{
			Term: core.PQTerm{
				P: big.NewInt(term),
				Q: big.NewInt(1),
			},
			Range: exactRangeForNamedRational(suffixes[i]),
		}
	}

	stream, status := core.NewFinitePQStream(steps)
	if status != core.StatusOK {
		panic("named: failed to construct finite RCF prefix stream")
	}
	return stream
}

func suffixRationalsFromRCFTerms(terms []int64) []core.Rational {
	n := len(terms)
	out := make([]core.Rational, n)

	current := core.NewRational(big.NewInt(terms[n-1]), big.NewInt(1))
	out[n-1] = current

	for i := n - 2; i >= 0; i-- {
		a := big.NewInt(terms[i])

		num := new(big.Int).Mul(a, current.Num())
		num.Add(num, current.Den())

		den := current.Num()
		current = core.NewRational(num, den)
		out[i] = current
	}

	return out
}

func exactRangeForNamedRational(value core.Rational) core.Range {
	return core.Range{
		Lo: core.Endpoint{
			Value: value,
			Open:  false,
		},
		Hi: core.Endpoint{
			Value: value,
			Open:  false,
		},
		Inside: true,
	}
}

// named/pi.go v3
