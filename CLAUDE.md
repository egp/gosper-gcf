# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Full check: format, vet, staticcheck, tests, coverage
./tools/check.sh

# Run all tests (no cache)
go test -count=1 ./core ./named ./trig

# Run a single package
go test -count=1 ./core
go test -count=1 ./named
go test -count=1 ./trig

# Run a specific test
go test -count=1 ./core -run TestBB_GCF_Add

# Run pending (normally-skipped) tests
RUN_PENDING_TESTS=1 go test -count=1 ./core ./named

# Static analysis
go vet ./core ./named ./trig
staticcheck ./core ./named ./trig
```

The "all green" bar includes `RUN_PENDING_TESTS=1` on `./core ./named`.

Never push while normal tests are red (OK for `PENDING` to be red).

## Requirements Specification

- `docs/gosper_cf_requirements_spec.md` - Primary specs
- `docs/newSpec.md` - clarification specs

These specs are derived from Gosper's work, especially HAKMEM 101A-101C.
There are a few deliberate deviations from Gosper that should be well marked in the specs.


## Architecture

GoGCF is a Go library for exact arithmetic on generalized continued fractions (GCF) in the style of Gosper. It uses `*big.Int` for all data, including GCF input, RCF output, Rational, and Range(Interval). Only things like counters and loop indexes are excluded. The early milestone (MVP) expression is `sqrt(3/pi^2 + e) / (tanh(sqrt(5)) - sin(69°))`.

### Packages

- **`core/`** — all evaluation machinery; zero external clients, so public API may be changed if necessary.
- **`named/`** — named procedural sources (`E`, `Pi`, `Sqrt2`, `SinDegrees`); consumes `core` and `trig`
- **`trig/`** — trigonometric/hyperbolic kernels (`Sin`, `Tanh`, Lambert-family); consumes `core`

### Two-tier exactness

The central invariant: clients never see a retracted term.

1. **Speculator tier** — BLFT and DLFT kernels emit "noisy" (speculative) terms to keep the pipeline flowing during uncertain intervals.
2. **Rectifier tier** — a homographic LFT at the output end absorbs speculative terms and only emits a term once the integer floor is proven stable across the full current interval.

### Key types

| Type | Role |
|---|---|
| `RCFTerm` | Single regular continued-fraction term (one `*big.Int`) |
| `PQTerm` | Generalized input term `p + q/x'` |
| `Rational` | Exact `num/den` pair |
| `Range` | Interval `[Lo, Hi]` with openness flags and `Inside` indicator |
| `RCFStream` | Interface: `NextRCF()`, `Range()` |
| `PQStream` | Interface: `Next()`, `CurrentInterval()` |
| `GCF` | Main evaluator; implements `RCFStream` |

GCF constructors: `NewGCF0` (constant), `NewGCF1` (unary), `NewGCF2` (binary), `NewExactTerminalGCF`, `newObservedRCFGCF`.

### Mathematical kernels

**BLFT** — Bi-Linear Fractional Transform. State: 8-tuple `(a,b,c,d,e,f,g,h)` representing `(axy+bx+cy+d)/(exy+fx+gy+h)`. Handles Add, Sub, Mul, Div, Reciprocal.

**DLFT** — Diagonal Quadratic LFT. State: 6-tuple `(a,b,c,d,e,f)` representing `(ax²+bx+c)/(dx²+ex+f)`. Handles Square and Sqrt. Not subject to monotonicity; feeds the Rectifier.

**Rectifier** — 4-tuple homographic LFT. Absorbs speculative BLFT/DLFT output; emits only proven terms.

### Sqrt feedback loop

`Sqrt(x)` uses Newton refinement via an ouroboros feedback stream: emitted approximation terms are recorded by `feedbackRCFStream`, converted back to a `PQStream` via `PQStreamFromRCF`, and fed into the DLFT kernel as the running approximation.

### Named sources

All procedural named sources expose an exact `Range()` at every step. BB tests for named sources verify `Range()` at each step and validate against published gold standards (OEIS) where available, with at least 50 terms when practical.

### Test naming conventions

- `*_bb_test.go` — black-box / integration tests
- `*_wb_test.go` — white-box / unit tests

### Pending tests

Tests guarded by `RUN_PENDING_TESTS=1` are intentionally red and represent the current TODO surface:

- `TestBB_GCF_SqrtOfTwoMatchesKnownPrefix` — parked sqrt(2) signoff; keep red until addressed
- `TestBB_SquareOfSqrt2IsExactlyTwo` — deferred; blocked by DLFT infinite algebraic-source limitation

### Feature-start workflow

1. Add public stubs
2. Add full BB signoff list as pending (guarded by `RUN_PENDING_TESTS=1`)
3. Activate tests as implementation turns them green
4. Add WB tests per R/G cycle
5. Remove skip logic from each test as it goes green

### Current priorities (as of 2026-04)

1. Finish remaining pending `SinDegrees` non-inside binary-range panic (non-inside BLFT path in `core`)
2. Keep parked `sqrt(2)` pending red explicit
3. Defer full general outside/outside BLFT range solver until after SinDegrees signoff
