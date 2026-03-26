// core/sqrt_controller.go v6
package core

import "math/big"

type sqrtController struct {
	x               PQStream
	ourorobos       *feedbackRCFStream
	ourorobosApprox PQStream
	half            PQStream
	seed            PQStream
}

func newSqrtController(x PQStream) *sqrtController {
	if x == nil {
		panic("newSqrtController: nil x")
	}

	validateSqrtNonNegativeRange(x.Range())

	ourorobos := newFeedbackRCFStream()

	return &sqrtController{
		x:               x,
		ourorobos:       ourorobos,
		ourorobosApprox: PQStreamFromRCF(ourorobos),
		half:            exactHalfPQStream(),
		seed:            selectSqrtSeedApproximation(x.Range()),
	}
}

func (c *sqrtController) ourorobosApproximation() PQStream {
	if c == nil {
		panic("sqrtController.ourorobosApproximation: nil receiver")
	}
	return c.ourorobosApprox
}

func (c *sqrtController) halfSource() PQStream {
	if c == nil {
		panic("sqrtController.halfSource: nil receiver")
	}
	return c.half
}

func (c *sqrtController) seedApproximation() PQStream {
	if c == nil {
		panic("sqrtController.seedApproximation: nil receiver")
	}
	return c.seed
}

func (c *sqrtController) buildRefinement(y PQStream) *GCF {
	if c == nil {
		panic("sqrtController.buildRefinement: nil receiver")
	}
	if y == nil {
		panic("sqrtController.buildRefinement: nil approximation stream")
	}

	divNode := Div(c.x, y)
	addNode := Add(y, PQStreamFromRCF(divNode))
	return Mul(PQStreamFromRCF(addNode), c.half)
}

func (c *sqrtController) feedCertifiedTerm(term RCFTerm, rng Range) {
	if c == nil {
		panic("sqrtController.feedCertifiedTerm: nil receiver")
	}

	c.ourorobos.Append(term, rng)
}

func exactHalfPQStream() PQStream {
	return PQStreamFromRational(NewRational(big.NewInt(1), big.NewInt(2)))
}

func selectSqrtSeedApproximation(r Range) PQStream {
	if root, ok := exactPositiveClosedSquareRoot(r); ok {
		return PQStreamFromRational(root)
	}

	one := RationalFromInt64(1)

	if !r.Inside {
		return PQStreamFromRational(one)
	}

	if !r.Lo.Open && r.Lo.Value.Num().Sign() > 0 {
		return PQStreamFromRational(r.Lo.Value)
	}

	return PQStreamFromRational(one)
}

func validateSqrtNonNegativeRange(r Range) {
	if rangeProvablyWhollyNegative(r) {
		panic("Sqrt: radicand range is wholly negative")
	}
}

func rangeProvablyWhollyNegative(r Range) bool {
	if !r.Inside {
		return false
	}

	hiSign := r.Hi.Value.Num().Sign()
	if hiSign < 0 {
		return true
	}
	if hiSign == 0 && r.Hi.Open {
		return true
	}

	return false
}

func exactPositiveClosedSquareRoot(r Range) (Rational, bool) {
	if !r.Inside || r.Lo.Open || r.Hi.Open {
		return Rational{}, false
	}
	if r.Lo.Value.Cmp(r.Hi.Value) != 0 {
		return Rational{}, false
	}

	value := r.Lo.Value
	if value.Num().Sign() <= 0 {
		return Rational{}, false
	}

	numRoot, numOK := exactIntegerSquareRoot(value.Num())
	denRoot, denOK := exactIntegerSquareRoot(value.Den())
	if !numOK || !denOK {
		return Rational{}, false
	}

	return NewRational(numRoot, denRoot), true
}

func exactIntegerSquareRoot(n *big.Int) (*big.Int, bool) {
	if n == nil || n.Sign() < 0 {
		return nil, false
	}

	root := new(big.Int).Sqrt(n)
	square := new(big.Int).Mul(new(big.Int).Set(root), new(big.Int).Set(root))
	if square.Cmp(n) != 0 {
		return nil, false
	}

	return root, true
}

// core/sqrt_controller.go v6
