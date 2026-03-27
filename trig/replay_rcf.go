// trig/replay_rcf.go v2
package trig

import (
	"math/big"

	"github.com/egp/gosper-gcf/core"
)

type replayRCF struct {
	src      core.RCFStream
	terms    []core.RCFTerm
	ranges   []core.Range
	eofKnown bool
}

type replayRCFFork struct {
	root  *replayRCF
	index int
}

func newReplayRCF(src core.RCFStream) *replayRCF {
	if src == nil {
		panic("newReplayRCF: nil source")
	}
	return &replayRCF{
		src:    src,
		ranges: []core.Range{cloneRangeReplay(src.Range())},
	}
}

func (r *replayRCF) Fork() *replayRCFFork {
	if r == nil {
		panic("(*replayRCF).Fork: nil receiver")
	}
	return &replayRCFFork{
		root:  r,
		index: 0,
	}
}

func (r *replayRCF) ensureCached(index int) {
	if r == nil {
		panic("(*replayRCF).ensureCached: nil receiver")
	}
	if index < 0 {
		panic("(*replayRCF).ensureCached: negative index")
	}

	for len(r.terms) <= index && !r.eofKnown {
		term, status := r.src.NextRCF()
		switch status {
		case core.StatusOK:
			current := r.ranges[len(r.terms)]
			next := suffixRangeAfterRCFTermReplay(current, term)
			r.terms = append(r.terms, core.NewRCFTerm(term.A()))
			r.ranges = append(r.ranges, next)
		case core.StatusEOF:
			r.eofKnown = true
		default:
			panic("(*replayRCF).ensureCached: invalid input status")
		}
	}
}

func (f *replayRCFFork) NextRCF() (core.RCFTerm, core.Status) {
	if f == nil {
		panic("(*replayRCFFork).NextRCF: nil receiver")
	}
	f.root.ensureCached(f.index)
	if f.index >= len(f.root.terms) {
		return core.NewRCFTerm(nil), core.StatusEOF
	}
	term := f.root.terms[f.index]
	f.index++
	return core.NewRCFTerm(term.A()), core.StatusOK
}

func (f *replayRCFFork) Range() core.Range {
	if f == nil {
		panic("(*replayRCFFork).Range: nil receiver")
	}
	return cloneRangeReplay(f.root.ranges[f.index])
}

func suffixRangeAfterRCFTermReplay(current core.Range, term core.RCFTerm) core.Range {
	if !current.Inside {
		panic("suffixRangeAfterRCFTermReplay: outside ranges not yet supported")
	}

	if isExactClosedRangeReplay(current) {
		exact := current.Lo.Value
		shifted := subtractIntegerFromRationalReplay(exact, term.A())
		if shifted.Num().Sign() == 0 {
			return exactRangeFromRationalReplay(core.RationalFromInt64(0))
		}
		return exactRangeFromRationalReplay(reciprocalRationalReplay(shifted))
	}

	loShift := subtractIntegerFromRationalReplay(current.Lo.Value, term.A())
	hiShift := subtractIntegerFromRationalReplay(current.Hi.Value, term.A())
	if loShift.Num().Sign() <= 0 || hiShift.Num().Sign() <= 0 {
		panic("suffixRangeAfterRCFTermReplay: nonpositive denominator interval")
	}

	return core.Range{
		Lo: core.Endpoint{
			Value: reciprocalRationalReplay(hiShift),
			Open:  current.Hi.Open,
		},
		Hi: core.Endpoint{
			Value: reciprocalRationalReplay(loShift),
			Open:  current.Lo.Open,
		},
		Inside: true,
	}
}

func isExactClosedRangeReplay(r core.Range) bool {
	return r.Inside &&
		!r.Lo.Open &&
		!r.Hi.Open &&
		r.Lo.Value.Cmp(r.Hi.Value) == 0
}

func subtractIntegerFromRationalReplay(r core.Rational, a *big.Int) core.Rational {
	num := r.Num()
	den := r.Den()
	scaled := new(big.Int).Mul(new(big.Int).Set(a), den)
	num.Sub(num, scaled)
	return core.NewRational(num, den)
}

func reciprocalRationalReplay(r core.Rational) core.Rational {
	if r.Num().Sign() == 0 {
		panic("reciprocalRationalReplay: zero numerator")
	}
	return core.NewRational(r.Den(), r.Num())
}

func exactRangeFromRationalReplay(value core.Rational) core.Range {
	return core.Range{
		Lo: core.Endpoint{
			Value: cloneRationalReplay(value),
			Open:  false,
		},
		Hi: core.Endpoint{
			Value: cloneRationalReplay(value),
			Open:  false,
		},
		Inside: true,
	}
}

func cloneRangeReplay(r core.Range) core.Range {
	return core.Range{
		Lo: core.Endpoint{
			Value: cloneRationalReplay(r.Lo.Value),
			Open:  r.Lo.Open,
		},
		Hi: core.Endpoint{
			Value: cloneRationalReplay(r.Hi.Value),
			Open:  r.Hi.Open,
		},
		Inside: r.Inside,
	}
}

func cloneRationalReplay(r core.Rational) core.Rational {
	return core.NewRational(r.Num(), r.Den())
}

// trig/replay_rcf.go v2
