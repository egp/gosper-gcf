# GoGCF Project Status

_Last updated: 2026-04-07_

---

## Green Bar

Normal tests (`go test -count=1 ./core ./named ./trig`) — **all pass**.

Full bar (`RUN_PENDING_TESTS=1 go test -count=1 ./core ./named`) — **red on PENDING tests below** (intentional).

---

## PENDING Test Inventory

Tests guarded by `RUN_PENDING_TESTS=1`; intentionally red. Use `grep skipIfPending` to list all.
See `docs/api_limitations.md` for the corresponding public API limitations.

| Test | Package | Reason |
|---|---|---|
| `TestBB_Named_SinDegrees_SixtyNineMatchesKnownPrefix` | `named` | trig.Sin broken for irrational results; deferred pre-Stage 5 |
| `TestBB_Trig_Sin_OneHalfMatchesKnownPrefix` | `trig` | trig.Sin broken for irrational results; deferred pre-Stage 5 |
| `TestBB_Trig_Sin_SixtyNineRadiansMatchesKnownPrefix` | `trig` | trig.Sin broken for irrational results (wrong terms for large integer input); deferred pre-Stage 5 |
| `TestBB_GCF_SqrtOfTwoMatchesKnownPrefix` | `named` | sqrt(2) parked; keep red until Stage 2 |
| `TestBB_SquareOfSqrt2IsExactlyTwo` | `named` | DLFT infinite algebraic-source limitation; Stage 2 |
| `TestWB_GCF_UnaryIdentity_OverPQStreamFromRCF_OfDegreesScaleOutsideCase_DoesNotPanic` | `core` | unary BLFT F=0 + outside-range; different root cause from sin tests; TBD |

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
| 3 | Stage 1 close: document trig.Sin limitations; add sin(69 radians) pending test | ✓ GREEN — `docs/api_limitations.md` created; `TestBB_Trig_Sin_SixtyNineRadiansMatchesKnownPrefix` confirmed failing then guarded; `CLAUDE.md` priorities updated. trig.Sin irrational-result fix deferred pre-Stage 5 |

**Stage 1 closed.** trig.Sin broken for irrational results (wrong terms, not just hang); outside/outside BLFT fix required before Stage 5 assembly.

---

## Requirements Specifications

- `docs/gosper_cf_requirements_spec.md` — authoritative specification (merged v6+v7+addendum, 2026-04-07)

Derived from Gosper's work, especially HAKMEM 101A–101C. Deliberate deviations from Gosper are marked in the spec.
