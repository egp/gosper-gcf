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
| `TestBB_Named_SinDegrees_ThirtyIsExactlyOneHalf` | `named` | wrong term 2 (gets 1, want 2) — diagnosis pending |
| `TestBB_Named_SinDegrees_NinetyIsExactlyOne` | `named` | Lambert series hits tan(π/4)=1 pole; needs exact-sin shortcut for 90° |
| `TestBB_Named_SinDegrees_MinusThirtyIsExactlyMinusOneHalf` | `named` | wrong term 2 (gets 3, want 2) — same root cause as 30° |
| `TestBB_Named_SinDegrees_SixtyNineMatchesKnownPrefix` | `named` | outside/outside BLFT range case (deferred) |
| `TestWB_SinDegrees_RationalValuedAngles` (subtests except 0°) | `named` | diagnostic scaffold — activatable case by case |
| `TestBB_GCF_SqrtOfTwoMatchesKnownPrefix` | `named` | sqrt(2) parked; keep red until Stage 2 |
| `TestBB_SquareOfSqrt2IsExactlyTwo` | `named` | DLFT infinite algebraic-source limitation |
| `TestBB_Trig_Sin_OneHalfMatchesKnownPrefix` | `trig` | binary BLFT non-inside range path (same family) |
| `TestWB_GCF_UnaryIdentity_OverPQStreamFromRCF_OfDegreesScaleOutsideCase_DoesNotPanic` | `core` | unary BLFT outside-range (projective wrap-around) case |

---

## MVP Stage Plan

Target expression: `sqrt(3/pi² + e) / (tanh(sqrt(5)) − sin(69°))`

| Stage | Focus | ~R/G cycles | Status |
|---|---|---|---|
| 1 | SinDegrees full signoff — binary non-inside BLFT path | 7 | **in progress** |
| 2 | Sqrt signoff — sqrt(2) prefix, 50-term OEIS validation | 7 | pending |
| 3 | Tanh signoff — tanh(sqrt(5)) prefix, full BB suite | 7 | pending |
| 4 | Named compositions: `sqrt(5)`, `3/pi²`, `3/pi² + e` | 7 | pending |
| 5 | Full MVP assembly: numerator, denominator, quotient, emit | 6 | pending |

### Stage 1 R/G cycle log

| Cycle | Target | Result |
|---|---|---|
| 1 | binary BLFT outside/inside CornerRange | ✓ GREEN — new `blft_range_outside.go`; updated stale guard tests |

**Remaining in Stage 1 (Cycle 2+):**
- 90°: add exact-sin special case for 90°/270° in `SinDegrees` (Lambert hits tan(π/4)=1 pole)
- 30°, −30°: diagnose wrong term 2 in binary BLFT inside-inside path
- 69°: outside/outside BLFT range (deferred until after 30°/90° signoff)

---

## Requirements Specifications

- `docs/gosper_cf_requirements_spec.md` — primary spec
- `docs/newSpec.md` — clarification spec

Derived from Gosper's work, especially HAKMEM 101A–101C. Deliberate deviations from Gosper are marked in the specs.
