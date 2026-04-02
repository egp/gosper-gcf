// trig/replay_rcf.go v5
package trig

import (
	"fmt"
	"math/big"

	"github.com/egp/gosper-gcf/core"
)

type replayRCF struct {
	src      core.RCFStream
	terms    []core.RCFTerm
	ranges   []core.Range
	eofKnown bool
	err      error
}

type replayRCFFork struct {
	root  *replayRCF
	index int
}

type staticRangePQReplay struct {
	rng core.Range
}

type errorRCFReplay struct {
	err error
}

type errorPQReplay struct {
	err error
}

func newReplayRCF(src core.RCFStream) *replayRCF {
	if src == nil {
		return &replayRCF{
			err: fmt.Errorf("newReplayRCF: %w", core.ErrNilObservedSource),
		}
	}

	r := &replayRCF{src: src}

	initial, err := src.CurrentInterval()
	if err != nil {
		r.err = fmt.Errorf("newReplayRCF: source interval: %w", err)
		return r
	}
	r.ranges = []core.Range{cloneRangeReplay(initial)}
	return r
}

func (r *replayRCF) Fork() *replayRCFFork {
	return &replayRCFFork{
		root:  r,
		index: 0,
	}
}

func (r *replayRCF) ensureCached(index int) error {
	if r == nil {
		return fmt.Errorf("(*replayRCF).ensureCached: %w", core.ErrNilReceiver)
	}
	if index < 0 {
		return fmt.Errorf("(*replayRCF).ensureCached: negative index")
	}
	if r.err != nil {
		return r.err
	}

	for len(r.terms) <= index && !r.eofKnown && r.err == nil {
		term, status, err := r.src.NextRCF()
		if err != nil {
			r.err = fmt.Errorf("(*replayRCF).ensureCached: source NextRCF: %w", err)
			break
		}

		switch status {
		case core.StatusOK:
			current := r.ranges[len(r.terms)]
			next, nextErr := suffixRangeAfterRCFTermReplay(current, term)
			if nextErr != nil {
				r.err = fmt.Errorf("(*replayRCF).ensureCached: next suffix interval: %w", nextErr)
				break
			}
			r.terms = append(r.terms, core.NewRCFTerm(term.A()))
			r.ranges = append(r.ranges, next)

		case core.StatusEOF:
			r.eofKnown = true

		default:
			r.err = fmt.Errorf("(*replayRCF).ensureCached: invalid input status=%v", status)
		}
	}

	return r.err
}

func (f *replayRCFFork) NextRCF() (core.RCFTerm, core.Status, error) {
	if f == nil || f.root == nil {
		return core.NewRCFTerm(nil), core.StatusEOF, fmt.Errorf("(*replayRCFFork).NextRCF: %w", core.ErrNilReceiver)
	}
	if err := f.root.ensureCached(f.index); err != nil {
		return core.NewRCFTerm(nil), core.StatusEOF, err
	}
	if f.index >= len(f.root.terms) {
		return core.NewRCFTerm(nil), core.StatusEOF, nil
	}
	term := f.root.terms[f.index]
	f.index++
	return core.NewRCFTerm(term.A()), core.StatusOK, nil
}

func (f *replayRCFFork) CurrentInterval() (core.Interval, error) {
	if f == nil || f.root == nil {
		return core.Interval{}, fmt.Errorf("(*replayRCFFork).CurrentInterval: %w", core.ErrNilReceiver)
	}
	err := f.root.ensureCached(f.index)
	if f.index < len(f.root.ranges) {
		return cloneRangeReplay(f.root.ranges[f.index]), nil
	}
	if err != nil {
		return core.Interval{}, err
	}
	return core.Interval{}, fmt.Errorf("(*replayRCFFork).CurrentInterval: %w", core.ErrUndefinedRangeOnEOFStream)
}

func (f *replayRCFFork) Range() (core.Range, error) {
	return f.CurrentInterval()
}

func (s *staticRangePQReplay) NextPQ() (core.PQTerm, core.PQStream, core.Status, error) {
	if s == nil {
		return core.PQTerm{}, s, core.StatusEOF, fmt.Errorf("(*staticRangePQReplay).NextPQ: %w", core.ErrNilReceiver)
	}
	return core.PQTerm{}, s, core.StatusEOF, nil
}

func (s *staticRangePQReplay) CurrentInterval() (core.Interval, error) {
	if s == nil {
		return core.Interval{}, fmt.Errorf("(*staticRangePQReplay).CurrentInterval: %w", core.ErrNilReceiver)
	}
	return cloneRangeReplay(s.rng), nil
}

func (s *staticRangePQReplay) Range() (core.Range, error) {
	return s.CurrentInterval()
}

func (s *errorRCFReplay) NextRCF() (core.RCFTerm, core.Status, error) {
	if s == nil {
		return core.NewRCFTerm(nil), core.StatusEOF, fmt.Errorf("(*errorRCFReplay).NextRCF: %w", core.ErrNilReceiver)
	}
	return core.NewRCFTerm(nil), core.StatusEOF, s.err
}

func (s *errorRCFReplay) CurrentInterval() (core.Interval, error) {
	if s == nil {
		return core.Interval{}, fmt.Errorf("(*errorRCFReplay).CurrentInterval: %w", core.ErrNilReceiver)
	}
	return core.Interval{}, s.err
}

func (s *errorRCFReplay) Range() (core.Range, error) {
	return s.CurrentInterval()
}

func (s *errorPQReplay) NextPQ() (core.PQTerm, core.PQStream, core.Status, error) {
	if s == nil {
		return core.PQTerm{}, s, core.StatusEOF, fmt.Errorf("(*errorPQReplay).NextPQ: %w", core.ErrNilReceiver)
	}
	return core.PQTerm{}, s, core.StatusEOF, s.err
}

func (s *errorPQReplay) CurrentInterval() (core.Interval, error) {
	if s == nil {
		return core.Interval{}, fmt.Errorf("(*errorPQReplay).CurrentInterval: %w", core.ErrNilReceiver)
	}
	return core.Interval{}, s.err
}

func (s *errorPQReplay) Range() (core.Range, error) {
	return s.CurrentInterval()
}

func suffixRangeAfterRCFTermReplay(current core.Range, term core.RCFTerm) (core.Range, error) {
	if isExactClosedRangeReplay(current) {
		exact := current.Lo.Value
		shifted := subtractIntegerFromRationalReplay(exact, term.A())
		if shifted.Num().Sign() == 0 {
			return exactRangeFromRationalReplay(core.RationalFromInt64(0)), nil
		}
		recip, err := reciprocalRationalReplay(shifted)
		if err != nil {
			return core.Range{}, err
		}
		return exactRangeFromRationalReplay(recip), nil
	}
	return suffixRangeViaUnaryRangeReplay(current, term.A())
}

func suffixRangeViaUnaryRangeReplay(current core.Range, a *big.Int) (core.Range, error) {
	g := core.NewGCF1(
		core.BLFTCoefficients{
			A: big.NewInt(0),
			B: big.NewInt(0),
			C: big.NewInt(0),
			D: big.NewInt(1),
			E: big.NewInt(0),
			F: big.NewInt(1),
			G: big.NewInt(0),
			H: new(big.Int).Neg(new(big.Int).Set(a)),
		},
		&staticRangePQReplay{rng: current},
	)
	rng, err := g.CurrentInterval()
	if err != nil {
		return core.Range{}, err
	}
	return cloneRangeReplay(rng), nil
}

func isExactClosedRangeReplay(r core.Range) bool {
	return r.Inside && !r.Lo.Open && !r.Hi.Open && r.Lo.Value.Cmp(r.Hi.Value) == 0
}

func subtractIntegerFromRationalReplay(r core.Rational, a *big.Int) core.Rational {
	num := r.Num()
	den := r.Den()
	scaled := new(big.Int).Mul(new(big.Int).Set(a), den)
	num.Sub(num, scaled)
	return core.NewRational(num, den)
}

func reciprocalRationalReplay(r core.Rational) (core.Rational, error) {
	if r.Num().Sign() == 0 {
		return core.Rational{}, fmt.Errorf("reciprocalRationalReplay: zero numerator")
	}
	return core.NewRational(r.Den(), r.Num()), nil
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
		Kind_:  r.Kind_,
	}
}

func cloneRationalReplay(r core.Rational) core.Rational {
	return core.NewRational(r.Num(), r.Den())
}

// trig/replay_rcf.go v5
