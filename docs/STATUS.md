# GoGCF Project Status

_Last updated: 2026-04-06_

---

## Green Bar

Normal tests (`go test -count=1 ./core ./named ./trig`) — **all pass**.

Full bar (`RUN_PENDING_TESTS=1 go test -count=1 ./core ./named`) — **red on PENDING tests below** (intentional).

---

## PENDING Test Inventory

Tests guarded by `RUN_PENDING_TESTS=1`; intentionally red. Use `grep skipIfPending` to list all.

| Test | Package | Reason |
|---|---|---|
| `TestBB_Named_SinDegrees_ThirtyIsExactlyOneHalf` | `named` | binary BLFT non-inside range path unimplemented |
| `TestBB_Named_SinDegrees_NinetyIsExactlyOne` | `named` | same |
| `TestBB_Named_SinDegrees_MinusThirtyIsExactlyMinusOneHalf` | `named` | same |
| `TestBB_Named_SinDegrees_SixtyNineMatchesKnownPrefix` | `named` | same |
| `TestWB_SinDegrees_RationalValuedAngles` (subtests except 0°) | `named` | same — diagnostic scaffold |
| `TestBB_GCF_SqrtOfTwoMatchesKnownPrefix` | `named` | sqrt(2) parked; keep red until Stage 2 |
| `TestBB_SquareOfSqrt2IsExactlyTwo` | `named` | DLFT infinite algebraic-source limitation |
| `TestBB_Trig_Sin_OneHalfMatchesKnownPrefix` | `trig` | binary BLFT non-inside range path unimplemented |
| `TestWB_GCF_UnaryIdentity_OverPQStreamFromRCF_OfDegreesScaleOutsideCase_DoesNotPanic` | `core` | unary BLFT outside-range (projective wrap-around) case |

---

## MVP Stage Plan

Target expression: `sqrt(3/pi² + e) / (tanh(sqrt(5)) − sin(69°))`

| Stage | Focus | ~R/G cycles | Status |
|---|---|---|---|
| 1 | SinDegrees full signoff — binary non-inside BLFT path | 7 | **current** |
| 2 | Sqrt signoff — sqrt(2) prefix, 50-term OEIS validation | 7 | pending |
| 3 | Tanh signoff — tanh(sqrt(5)) prefix, full BB suite | 7 | pending |
| 4 | Named compositions: `sqrt(5)`, `3/pi²`, `3/pi² + e` | 7 | pending |
| 5 | Full MVP assembly: numerator, denominator, quotient, emit | 6 | pending |

---

## Requirements Specifications

- `docs/gosper_cf_requirements_spec.md` — primary spec
- `docs/newSpec.md` — clarification spec

Derived from Gosper's work, especially HAKMEM 101A–101C. Deliberate deviations from Gosper are marked in the specs.
