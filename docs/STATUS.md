# GoGCF Project Status

_Last updated: 2026-04-07_

---

## Green Bar

Normal tests (`go test -count=1 ./core ./named ./trig`) — **all pass**.

Full bar (`RUN_PENDING_TESTS=1 go test -count=1 ./core ./named`) — **red on PENDING tests below** (intentional).

---

## PENDING Test Inventory

Tests guarded by `RUN_PENDING_TESTS=1`; intentionally red. Use `grep skipIfPending` to list all.

| Test | Package | Reason |
|---|---|---|
| `TestBB_Named_SinDegrees_SixtyNineMatchesKnownPrefix` | `named` | outside/outside BLFT range case (deferred until Stage 1 complete) |
| `TestBB_GCF_SqrtOfTwoMatchesKnownPrefix` | `named` | sqrt(2) parked; keep red until Stage 2 |
| `TestBB_SquareOfSqrt2IsExactlyTwo` | `named` | DLFT infinite algebraic-source limitation |
| `TestBB_Trig_Sin_OneHalfMatchesKnownPrefix` | `trig` | binary BLFT non-inside range path (same family as 69°) |
| `TestWB_GCF_UnaryIdentity_OverPQStreamFromRCF_OfDegreesScaleOutsideCase_DoesNotPanic` | `core` | unary BLFT outside-range (projective wrap-around) case |

---

## MVP Stage Plan

Target expression: `sqrt(3/pi² + e) / (tanh(sqrt(5)) − sin(69°))`

| Stage | Focus | ~R/G cycles | Status |
|---|---|---|---|
| 1 | SinDegrees full signoff — binary non-inside BLFT path | 7 | **in progress** |
| 1.5 | HAKMEM 101C utility (`SmallestRationalInInterval`) | 2 | pending (prerequisite for Stage 4) |
| 2 | Sqrt signoff — sqrt(2) prefix, 50-term OEIS validation | 7 | pending |
| 3 | Tanh signoff — tanh(sqrt(5)) prefix, full BB suite | 7 | pending |
| 4 | Named compositions: `sqrt(5)`, `3/pi²`, `3/pi² + e` | 7 | pending (requires 101C) |
| 4.5 | Resource guards: Timeout + BitLen limits, explicit error return | 2 | pending (prerequisite for Stage 5) |
| 5 | Full MVP assembly: numerator, denominator, quotient, emit | 6 | pending (requires resource guards) |

### Stage 1 R/G cycle log

| Cycle | Target | Result |
|---|---|---|
| 1 | binary BLFT outside/inside CornerRange | ✓ GREEN — new `blft_range_outside.go`; updated stale guard tests |
| 2 | exact-sin shortcut for rational-valued angles (30°, 90°, −30°, and all multiples of 30°) | ✓ GREEN — `exactSinDegreesShortcut` in `named/sin_degrees.go`; exported `NewExactTerminalGCFFromRational`; activated 3 BB tests + all WB subtests |

**Remaining in Stage 1 (Cycle 3+):**
- 69°: outside/outside BLFT range (deferred; needs full projective outside/outside case in `core`)
- Unary BLFT outside-range (projective wrap-around) — also required for full Stage 1 signoff

---

## Requirements Specifications

- `docs/gosper_cf_requirements_spec.md` — authoritative specification (merged v6+v7+addendum, 2026-04-07)

Derived from Gosper's work, especially HAKMEM 101A–101C. Deliberate deviations from Gosper are marked in the spec.
