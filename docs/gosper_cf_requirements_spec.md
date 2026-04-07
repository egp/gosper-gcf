# GoGCF Specification

**Version:** 8 (merged v6 + v7 + addendum)  
**Timestamp:** 2026-04-07

---

## 1. Background

This document specifies a Go library implementing Gosper-style continued-fraction arithmetic, grounded in HAKMEM Items 101A (representation and best-approximation properties), 101B (GCF arithmetic), and 101C (simplest rational in an interval).

Gosper's motivating example is a composed expression involving `sqrt`, `pi`, `e`, `tanh`, and `sin`. This project targets:

> `sqrt(3/pi² + e) / (tanh(sqrt(5)) − sin(69°))`

The degree symbol is a deliberate project choice. HAKMEM writes `sin(69)` without specifying units; this project interprets it as degrees for the client-facing helper, while the core operator works in radians with an exact degree-to-radian conversion before dispatch.

### Priorities (faithful to Gosper)

- Exact arithmetic over approximation
- Demand-driven production of output terms
- Exact interval tracking at every stage
- Composition of operators over GCF streams rather than conversion through floating point

---

## 2. Source Fidelity and Conflict Policy

Every requirement is labeled with one of:

- **Faithful to HAKMEM** — directly implementable as Gosper describes
- **Project restriction / modern repair** — stricter than HAKMEM to avoid silent looping or retractable client output
- **Project generalization** — structure or API not dictated by HAKMEM, but consistent with his model

### 2.1. Internal representation — Faithful to HAKMEM

All internal sources, adapters, generators, transducers, and kernels use generalized `(p, q)` streams where a term represents `p + q/x'`. Regular continued fractions (`q = 1`) are the special case. Only the final client-facing emitter produces RCF terms.

### 2.2. Function sources — Faithful to HAKMEM

Constants and functions such as `pi`, `e`, `sin`, and `tanh` are procedurally generated GCF sources, not values reconstructed from floating point. Programmatic generators and transducers are mandatory.

### 2.3. Numeric type — Project restriction

All coefficients (BLFT, DLFT, Rectifier) and all interval endpoint values use arbitrary-precision integers (`*big.Int`) and exact rationals. Binary floating point and fixed-width arithmetic are forbidden in correctness-critical logic.

### 2.4. Retractable output — Project synthesis (Two-Tier Architecture)

HAKMEM allows temporary wrong terms followed by recanting correction terms. This project forbids retractable client-visible output while preserving Gosper's fluid pipeline via a two-tier model:

1. **Speculator tier** — BLFT and DLFT kernels emit speculative `(p, q)` terms internally to maintain pipeline flow.
2. **Rectifier tier** — a homographic LFT at the output boundary absorbs speculative terms into its coefficients and emits a client-visible RCF term only when the integer floor is proven stable across the full current interval.

### 2.5. Productivity and rational collapse — Project restriction / modern repair

A computation may get stuck approaching an exact rational result as bounds converge but never cross an integer floor.

**Resolution:** The engine implements Rational Collapse. When exact interval analysis isolates a single rational value, the engine invokes the HAKMEM 101C utility, emits the finite exact RCF sequence, and terminates. The engine need not guarantee productivity for every input, but must either produce the next term or fail explicitly when configured resource limits are exceeded.

### 2.6. Quadratic single-input form — Faithful to HAKMEM, stronger certification

HAKMEM 101B introduces the diagonal quadratic form `(ax² + bx + c) / (dx² + ex + f)` and notes it is not guaranteed monotone. This project includes the quadratic kernel but requires it to feed the Rectifier rather than emit client-visible terms directly.

### 2.7. Angle-unit choice — Project generalization

The library defines:
- `sinRadians(x)` as the core operator (in `trig/`)
- `SinDegrees(x)` as the exact wrapper (in `named/`) that converts degrees to radians before dispatch

### 2.8. Decimal emission — Faithful to HAKMEM, deferred from MVP

Decimal digits are obtained by multiplying the residual transform by 10 instead of reciprocating after outputting the floor term. This is planned but not required for MVP. The design must not preclude it.

---

## 3. Requirements

### 3.1. Must Have

- **[Faithful]** Accept generalized `(p, q)` GCF input streams; use them for all internal sources, adapters, kernels, and transducers.
- **[Faithful]** Emit only RCF terms to clients.
- **[Restriction]** Use arbitrary-precision integers for all coefficients, terms, and interval endpoints. No `float32`, `float64`, or fixed-width arithmetic in correctness-critical logic.
- **[Faithful]** Implement a unified bihomographic kernel (BLFT) with bigint coefficients `(a,b,c,d,e,f,g,h)` for `z(x,y) = (axy+bx+cy+d)/(exy+fx+gy+h)`. Support binary arithmetic (`+`, `−`, `×`, `÷`) and unary specializations.
- **[Faithful]** Implement the diagonal quadratic kernel (DLFT) with bigint coefficients `(a,b,c,d,e,f)` for `z(x) = (ax²+bx+c)/(dx²+ex+f)`. Required for `square` and `sqrt`.
- **[Faithful]** Implement Gosper-style demand-driven ingest/emit: request input terms only when needed; emit output terms as soon as proven.
- **[Faithful]** Define deterministic source-exhaustion semantics: when a stream is exhausted, apply HAKMEM's tail limit rule (fold `x → ∞`) into the existing kernel coefficients and continue.
- **[Restriction]** Normalize coefficient tuples by dividing by GCD whenever that preserves the represented transform.
- **[Faithful]** Expose the current exact interval as a sibling of next-term emission (`CurrentInterval()` alongside `NextRCF()`).
- **[Faithful]** Implement inside/outside interval semantics on the projective line, including exact denominator-crossing logic and union of corner-derived ranges for BLFT propagation.
- **[Faithful]** Use Gosper's widest-range idea (or an exact equivalent) to choose which input stream to ingest when no output term is yet forced. Ordering: `Outside narrowness > Outside wideness > Inside wideness > Inside narrowness`.
- **[Faithful]** Provide built-in GCF generators for `pi` and `e`.
- **[Faithful/Generalization]** Provide unary operators: `sqrt`, `tanh`, `sinRadians`, `SinDegrees`, `square`, `reciprocal`. Trig and hyperbolic operators live in `trig/`; `SinDegrees` lives in `named/`. The core kernels (BLFT, DLFT) remain pure.
- **[Generalization]** `SinDegrees(x)` converts degrees to radians exactly (via BLFT) before dispatching to `sinRadians`. Exact shortcut: when the input angle is an integer multiple of 30° with a rational sine value, return the exact rational GCF directly without invoking the Lambert series.
- **[Faithful]** Treat `square(x)` as a distinct diagonal-quadratic operator, not binary multiplication with a duplicated source.
- **[Faithful]** Implement the HAKMEM 101C utility: given a nonempty interval, find the rational in that interval with the smallest numerator and denominator by the Gosper method (find first CF term where endpoints differ; increment the lesser; discard subsequent terms). Use for Rational Collapse and canonical witness generation in tests.
- **[Restriction]** Provide resource guards: caller-configurable limits on time, work (bit-length growth), and optionally emitted terms. Return an explicit error when limits are exceeded.
- **[Faithful]** Provide conversion from RCF prefixes to convergents `(Pₙ, Qₙ)` and exact rational values.
- **[Restriction]** A client-visible RCF term may be emitted only when exact interval analysis proves the integer floor is stable across the entire admissible image.
- **[Restriction]** No client-visible speculative or retractable terms.
- **[Restriction]** Ensure public API supports lazy evaluation of composed expressions.
- **[Generalization]** Support evaluation of the MVP target: `sqrt(3/pi² + e) / (tanh(sqrt(5)) − sin(69°))`.

### 3.2. Should Have

- **[Faithful]** Decimal-digit emission by the HAKMEM multiply-by-10 rule (deferred from MVP; design must not preclude it).
- **[Faithful]** Hurwitz number detection: many constants and functions (e, certain hyperbolic functions, ratios of Bessel functions) are Hurwitz numbers whose BLFT/DLFT coefficient sequences are naturally periodic or bounded. The kernels should preserve Hurwitz form when all active input streams are Hurwitz; optionally detect and report emerging periodicity for diagnostics and performance.
- **[Generalization]** Expose debugging state (transform coefficients, interval bounds, openness flags, denominator-crossing decisions, GCD reductions, source-consumption decisions).
- **[Generalization]** Support symbolic-sharing semantics so `square(x)` need not independently consume two copies of the same source.
- **[Generalization]** Configurable maximum emitted RCF term count as an additional resource guard.
- **[Faithful]** Tests and documentation should demonstrate the best-approximation and error-bound properties of HAKMEM 101A: that truncated CF prefixes give best rational approximations with error bounded by `1/(Q²·a_{n+1})`.

### 3.3. Could Have

- **[Generalization]** Additional unary operators (`exp`, `log`, inverse trig) when exact generator strategies are available.
- **[Generalization]** Expression graphs that memoize shared subexpressions during lazy evaluation.
- **[Generalization]** Support for approximate GCF sources (inputs known only approximately, with widening intervals) propagated exactly as Gosper describes.
- **[Generalization]** Optional serialization for GCF fixtures, RCF outputs, periodic forms, and debugging transcripts.
- **[Faithful]** Support piecewise transmission of large terms: a single large `(p, q)` value may be split into equivalent piecewise GCF segments ("thinly disguised multiprecision" per Gosper). All kernels must handle such streams transparently.

### 3.4. Won't Have (MVP)

- No reliance on binary floating point for correctness decisions.
- No reconstruction of transcendental CFs from sampled numeric values.
- No promise to support every generalized CF invariant form.
- No guarantee of productivity without caller-supplied resource limits.
- No distributed or parallel execution.

---

## 4. Architecture

### 4.1. Two-Tier Pipeline

```
PQStream sources → [BLFT / DLFT speculator tier] → [Rectifier] → client RCFStream
```

- **Speculator tier**: BLFT and DLFT kernels operate on GCF `(p,q)` streams, emit speculative terms, and track exact interval state.
- **Rectifier**: a homographic LFT `[a, b; c, d]` that sits at the end of every operator chain. It absorbs speculative terms (updating `a' = ap + cq`, `b' = a`, `c' = cp + dq`, `d' = c`), evaluates the output range as `[z(1), z(∞)] = [(a+b)/(c+d), a/c]`, and emits a client-visible RCF term `t` only when `⌊z(1)⌋ = ⌊z(∞)⌋`.

### 4.2. BLFT Kernel

State: `(a,b,c,d,e,f,g,h)` representing `z(x,y) = (axy+bx+cy+d)/(exy+fx+gy+h)`.

Operations:
- **Ingest x term `(p,q)`**: substitute `x = p + q/x'`.
- **Ingest y term `(r,s)`**: substitute `y = r + s/y'`.
- **Produce output term `t`**: substitute `z = t + 1/z'`.
- **Normalize**: divide all 8 coefficients by their GCD.
- **Collapse x-EOF**: apply `x → ∞`, zeroing the `a,b` dependence on x.
- **Collapse y-EOF**: apply `y → ∞`, zeroing the `a,c` dependence on y.

Initial states: Add `(0,1,1,0; 0,0,0,1)`, Sub `(0,1,−1,0; 0,0,0,1)`, Mul `(1,0,0,0; 0,0,0,1)`, Div `(0,1,0,0; 0,0,1,0)`, Reciprocal `(0,0,1,0; 0,1,0,0)`.

### 4.3. Diagonal Quadratic Kernel (DLFT)

State: `(a,b,c,d,e,f)` representing `z(x) = (ax²+bx+c)/(dx²+ex+f)`.

Not guaranteed monotone; must feed the Rectifier. Primary uses: `square`, and the Newton-feedback loop for `sqrt`.

### 4.4. Interval Model

A `Range` has two exact rational endpoints with openness flags and an `Inside`/`Outside` polarity bit.

- **Inside**: the admissible arc from `Lo` to `Hi` in the ordinary affine picture.
- **Outside**: the complement of the excluded gap — behaves like the union of two rays (the projective/toroidal view Gosper describes for `+∞ = −∞`).

BLFT range propagation: evaluate the four corner-derived ranges and combine by exact union. When the denominator changes sign over an inside range, the image becomes outside.

### 4.5. Tail Interpretation

When a finite source is exhausted (returns `io.EOF`), the kernel applies HAKMEM's limit rule: the unread tail has become infinite. For BLFT:
- `x → ∞`: retain only the `a,c` coefficients in numerator and `e,g` in denominator.
- `y → ∞`: retain only the `a,b` and `e,f` coefficients.

Evaluation continues in the same kernel representation.

### 4.6. HAKMEM 101C Utility

Given a nonempty interval, find the rational inside it with the smallest numerator and denominator:

1. Express both endpoints as continued fractions.
2. Find the first term where they differ.
3. Increment the smaller term (unless it is the last); discard subsequent terms.
4. If one endpoint terminates while matching the other, append `∞` and continue.

Roles: Rational Collapse (when bounds isolate an exact rational), canonical witness generation, interval-debugging support, test verification.

### 4.7. Sqrt Feedback Loop

`sqrt(x)` uses Newton refinement via an ouroboros feedback stream: emitted approximation terms are recorded by `feedbackRCFStream`, converted back to a `PQStream` via `PQStreamFromRCF`, and fed into the DLFT kernel as the running approximation. All steps satisfy the no-retractable-term invariant through the Rectifier.

### 4.8. Output Rule

A client-visible RCF term `t = ⌊z⌋` is emitted only when:
1. The current exact interval of the remaining input tails is computed.
2. The exact admissible image is derived as an inside/outside enclosure.
3. The floor is identical over the entire admissible image.

Otherwise: ingest the next term from the widest-range input stream and repeat.

---

## 5. Conformance

The implementation is conformant when it:

- Performs exact GCF arithmetic using generalized internal `(p,q)` streams.
- Preserves demand-driven ingest/emit behavior.
- Implements exact interval reasoning using inside/outside semantics.
- Includes the diagonal quadratic kernel for `square` and `sqrt`.
- Forbids speculative client-visible output.
- Provides the 101C rational-in-interval utility.
- Provides explicit failure on configured resource exhaustion.
- Can lazily evaluate the MVP expression and emit certified RCF terms.
- Includes conformance tests with HAKMEM examples and property-based tests for rational identities, monotone interval refinement, and agreement between certified output and interval floors.
