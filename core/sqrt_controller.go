// core/sqrt_controller.go v8
package core

import (
	"fmt"
	"math/big"
)

type sqrtController struct {
	x                    PQStream
	ourorobos            *feedbackRCFStream
	ourorobosApprox      PQStream
	hasOurorobosFeedback bool
	half                 PQStream
	seed                 PQStream
	initErr              error
}

func newSqrtController(x PQStream) *sqrtController {
	c := &sqrtController{
		x:         x,
		ourorobos: newFeedbackRCFStream(),
		half:      exactHalfPQStream(),
	}

	if x == nil {
		c.initErr = fmt.Errorf("newSqrtController: %w", ErrNilReceiver)
		c.ourorobosApprox = newErrorPQStream(c.initErr)
		c.seed = newErrorPQStream(c.initErr)
		return c
	}

	xRange, err := x.Range()
	if err != nil {
		c.initErr = fmt.Errorf("newSqrtController: x.Range: %w", err)
		c.ourorobosApprox = newErrorPQStream(c.initErr)
		c.seed = newErrorPQStream(c.initErr)
		return c
	}

	if err := validateSqrtNonNegativeRange(xRange); err != nil {
		c.initErr = err
		c.ourorobosApprox = newErrorPQStream(c.initErr)
		c.seed = newErrorPQStream(c.initErr)
		return c
	}

	c.ourorobosApprox = PQStreamFromRCF(c.ourorobos)
	c.seed = selectSqrtSeedApproximation(xRange)
	return c
}

func (c *sqrtController) ourorobosApproximation() PQStream {
	if c == nil {
		return newErrorPQStream(fmt.Errorf("sqrtController.ourorobosApproximation: %w", ErrNilReceiver))
	}
	if c.initErr != nil {
		return newErrorPQStream(c.initErr)
	}
	return c.ourorobosApprox
}

func (c *sqrtController) halfSource() PQStream {
	if c == nil {
		return newErrorPQStream(fmt.Errorf("sqrtController.halfSource: %w", ErrNilReceiver))
	}
	return c.half
}

func (c *sqrtController) seedApproximation() PQStream {
	if c == nil {
		return newErrorPQStream(fmt.Errorf("sqrtController.seedApproximation: %w", ErrNilReceiver))
	}
	if c.initErr != nil {
		return newErrorPQStream(c.initErr)
	}
	return c.seed
}

func (c *sqrtController) buildRefinement(y PQStream) *GCF {
	if c == nil {
		return newObservedRCFGCF(newErrorRCFStream(fmt.Errorf("sqrtController.buildRefinement: %w", ErrNilReceiver)))
	}
	if c.initErr != nil {
		return newObservedRCFGCF(newErrorRCFStream(c.initErr))
	}
	if y == nil {
		return newObservedRCFGCF(newErrorRCFStream(fmt.Errorf("sqrtController.buildRefinement: nil approximation stream")))
	}

	divNode := Div(c.x, y)
	addNode := Add(y, PQStreamFromRCF(divNode))
	return Mul(PQStreamFromRCF(addNode), c.half)
}

func (c *sqrtController) feedCertifiedTerm(term RCFTerm, rng Range) error {
	if c == nil {
		return fmt.Errorf("sqrtController.feedCertifiedTerm: %w", ErrNilReceiver)
	}
	if c.ourorobos == nil {
		return fmt.Errorf("sqrtController.feedCertifiedTerm: %w", ErrNilReceiver)
	}
	if err := c.ourorobos.Append(term, rng); err != nil {
		return err
	}
	c.hasOurorobosFeedback = true
	return nil
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

func validateSqrtNonNegativeRange(r Range) error {
	if rangeProvablyWhollyNegative(r) {
		return fmt.Errorf("Sqrt: radicand range is wholly negative")
	}
	return nil
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

// core/sqrt_controller.go v8
