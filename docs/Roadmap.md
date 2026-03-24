# Roadmap.md

# gosper-gcf Initial Roadmap

## Purpose

This document is the execution roadmap for the early development of `gosper-gcf`.

Unlike `Spec.md`, this file is intentionally about sequencing, test strategy, and development slices.

---

## 1. Strategy

Start from the edges and work toward the center.

The high-level order is:

1. foundational domain types
2. input source-stream edge
3. client-facing regular-output edge
4. BLFT center
5. range-driven emission and comparison
6. richer operations and named sources

The goal is to get stable test scaffolding around the engine before implementing the engine itself.

---

## 2. Phase 0 — repository skeleton

## Goal

Get the repository into a clean, ready-to-develop state.

## Deliverables

- `Spec.md`
- `Roadmap.md`
- `README.md`
- root `.gitignore`
- `go.mod`
- `core/`
- `named/`
- `core/doc.go`
- `named/doc.go`

No arithmetic yet.

---

## 3. Phase 1 — foundational domain types

## Goal

Stabilize the core nouns before implementing streams or transforms.

## Deliverables in `core`

- `Status`
- `PQTerm`
- `RCFTerm`
- `Rational`
- `RangeKind`
- `Endpoint`
- `Range`

## Key design intent

These types should be exact, small, testable, and free of transform logic.

## White-box tests

- `TestWB_Rational_NormalizesSign`
- `TestWB_Rational_ReducesByGCD`
- `TestWB_Rational_FromIntIsExact`
- `TestWB_PQTerm_AllowsNegativeFirstTermConstruction`
- `TestWB_RCFTerm_HoldsBigIntExactly`
- `TestWB_Range_InsideLoEqHiMeansExact`
- `TestWB_Range_CompareOrdersUncertaintyCorrectly`
- `TestWB_Endpoint_OpenClosedFlagsAreIndependent`

## Black-box tests

- `TestBB_Rational_ExactIntegerConversion`
- `TestBB_Range_ExactValueIsNotNecessarilyInteger`
- `TestBB_Range_InsideAndOutsideKindsConstructCleanly`
- `TestBB_Status_PublicValuesAreStable`

## Notes

This phase should not introduce BLFT or source-stream logic.

---

## 4. Phase 2 — input source-stream edge

## Goal

Finish the generalized input-stream contract before the engine exists.

Temporary name in this roadmap: `GCFStream` / input stream.
The current spec direction calls its pull method `NextPQ()`.

## Deliverables in `core`

- input stream interface or concrete contract
- validating finite source for tests
- pull-time and construction-time validation behavior
- `Range()` alongside `NextPQ()`
- EOF handling

## Validation policy

Raise a failure as soon as it is noticed:

- eagerly at construction time, when possible
- incrementally at pull time, for procedural or infinite sources

## White-box tests

- `TestWB_GCFStream_FirstTermMayBeNegative`
- `TestWB_GCFStream_LaterNegativePRejected`
- `TestWB_GCFStream_LaterNegativeQRejected`
- `TestWB_GCFStream_ZeroQRejected`
- `TestWB_GCFStream_EOFReturnsStatusEOF`
- `TestWB_GCFStream_TailAdvancesCorrectly`
- `TestWB_GCFStream_RangeAvailableOnLiveStream`
- `TestWB_GCFStream_IncrementalValidationTriggersAtBadTerm`

## Black-box tests

- `TestBB_GCFStream_ReadFiniteSequenceThenEOF`
- `TestBB_GCFStream_FirstMalformedTermFails`
- `TestBB_GCFStream_LaterMalformedTermFails`
- `TestBB_GCFStream_RangeTracksReturnedTail`
- `TestBB_GCFStream_ProceduralStreamCanFailLate`

## Property tests to co-develop

- valid finite sequences never produce invalid-input status
- any later negative `p` or `q` fails at or before that term
- any `q == 0` fails at that term

---

## 5. Phase 3 — client-facing regular-output edge

## Goal

Finish the client-facing regular-output shell before the center engine is implemented.

In v1, the client sees `GCF` as the regular-output stream object.

## Deliverables in `core`

- `GCF` shell type
- `NextRCF()` shape
- `Range()` shape
- exact terminal-state representation
- initial config plumbing
- placeholder constructors for plugboard-based creation

## White-box tests

- `TestWB_GCF_TerminalExactStateEmitsCorrectly`
- `TestWB_GCF_TerminalExactStateHasExactRange`
- `TestWB_GCF_FirstEmittedRegularTermMayBeNegative`
- `TestWB_GCF_EOFUsesStatusEOF`
- `TestWB_GCF_ConfigDefaultsApply`

## Black-box tests

- `TestBB_GCF_ReadFiniteRegularTermsThenEOF`
- `TestBB_GCF_RangeAvailableWhileLive`
- `TestBB_GCF_ExactStateProducesDeterministicTerms`
- `TestBB_GCF_ConfigOptionalAtCreation`

## Notes

This phase is still edge work. It should avoid full BLFT logic except for exact terminal-state behavior.

---

## 6. Phase 4 — BLFT center, part 1

## Goal

Bring in the algebraic heart without yet finishing all range-driven emission decisions.

## Deliverables in `core`

- BLFT state struct
- normalization helpers
- sign-canonicalization helpers
- BitLen checks
- ingest-X formula
- ingest-Y formula
- EOF collapse formulas
- terminal exact-state collapse

## White-box tests

- `TestWB_BLFT_IngestXMatchesGosperFormula`
- `TestWB_BLFT_IngestYMatchesGosperFormula`
- `TestWB_BLFT_CollapseOnXEOF`
- `TestWB_BLFT_CollapseOnYEOF`
- `TestWB_BLFT_CollapseOnBothEOFProducesExactState`
- `TestWB_BLFT_NormalizesAfterMutation`
- `TestWB_BLFT_BitLenCheckRunsAtMutationBoundary`

## Black-box tests

- `TestBB_BLFT_UnaryIdentityPassesThroughRegularInput`
- `TestBB_BLFT_GeneralizedInputHandlesQNotOne`
- `TestBB_BLFT_CollapseContinuesAfterEOF`

---

## 7. Phase 5 — BLFT center, part 2

## Goal

Add range-driven emission and comparison.

## Deliverables in `core`

- corner-range evaluation
- widest-corner-range choice logic
- `CanEmitRCFTerm`
- emission formula
- `Cmp()` based on subtraction

## White-box tests

- `TestWB_BLFT_EmitFormulaMatchesGosper`
- `TestWB_BLFT_InsideIntervalCanEmitWhenFloorsAgree`
- `TestWB_BLFT_OutsideIntervalDoesNotDirectlyEmit`
- `TestWB_BLFT_TieBreakGoesToX`
- `TestWB_GCF_CmpNegative`
- `TestWB_GCF_CmpZero`
- `TestWB_GCF_CmpPositive`

## Black-box tests

- `TestBB_GCF_AddProducesExpectedPrefix`
- `TestBB_GCF_SubProducesExpectedPrefix`
- `TestBB_GCF_MulProducesExpectedPrefix`
- `TestBB_GCF_DivProducesExpectedPrefix`
- `TestBB_GCF_CmpMatchesKnownCases`

---

## 8. Phase 6 — richer operations and `named`

## Goal

Expose a useful client toolkit once the core is stable.

## Deliverables

In `core`:

- `Add`
- `Sub`
- `Mul`
- `Div`
- `Reciprocal`
- `Square`

In `named`:

- `FromRational`
- finite/source wrappers for testing
- `Pi()`
- `E()`
- `Sqrt2()`

Later:

- `Sin`
- `Tanh`
- `Sqrt`
- Ouroboros / Newton feedback

---

## 9. Invariants and property testing

Co-develop white-box property tests in parallel with invariant enforcement.

Every important struct should have:

- an explicit invariant list
- explicit invariant-check helpers where appropriate
- matching white-box tests

Priority structs:

- `Rational`
- `Range`
- input stream implementation
- `GCF`
- BLFT state

---

## 10. Immediate next slice

Start with Phase 1.

## First files to create in `core`

- `status.go`
- `pqterm.go`
- `rcfterm.go`
- `rational.go`
- `endpoint.go`
- `range.go`

## First red tests to write

- `rational_wb_test.go`
- `range_wb_test.go`
- `rational_bb_test.go`
- `range_bb_test.go`

## First green target

Have exact, normalized, well-tested foundational types committed before introducing any stream or engine code.