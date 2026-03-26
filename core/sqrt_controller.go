// core/sqrt_controller.go v4
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

func exactHalfPQStream() PQStream {
	return PQStreamFromRational(NewRational(big.NewInt(1), big.NewInt(2)))
}

func selectSqrtSeedApproximation(r Range) PQStream {
	one := RationalFromInt64(1)

	if !r.Inside {
		return PQStreamFromRational(one)
	}

	if !r.Lo.Open && !r.Hi.Open && r.Lo.Value.Cmp(r.Hi.Value) == 0 {
		if r.Lo.Value.Num().Sign() > 0 {
			return PQStreamFromRational(r.Lo.Value)
		}
		return PQStreamFromRational(one)
	}

	if !r.Lo.Open && r.Lo.Value.Num().Sign() > 0 {
		return PQStreamFromRational(r.Lo.Value)
	}

	return PQStreamFromRational(one)
}

// core/sqrt_controller.go v4
