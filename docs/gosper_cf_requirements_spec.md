# SPEC-1-Gosper-Continued-Fraction-Arithmetic-Library

**Version:** 4  
**Timestamp:** 2026-03-31T00:00:00-04:00

## Background

This document defines the requirements for a software library implementing Gosper-style continued fraction arithmetic, grounded first in HAKMEM Item 101A (continued fraction representation) and Item 101B (continued fraction arithmetic), and only secondarily in later expositions and repair work.

The primary motivation is the one Gosper states directly: arithmetic should proceed exactly, incrementally, and on demand, producing continued-fraction output terms as soon as they are forced by the known input terms. HAKMEM also explicitly uses an example expression of the same general form as the MVP target here: a composed expression involving `sqrt`, `pi`, `e`, `tanh`, and `sin`, evaluated as a stream of continued-fraction terms rather than by floating-point approximation.

### Source Fidelity and Conflict Policy

Where the original HAKMEM text and a modern retelling disagree, the requirements in this document mark the conflict explicitly and choose one of three dispositions:

1. **Faithful to HAKMEM** — used when the original formulation is clear and implementable.
2. **Generalized beyond HAKMEM** — used when the MVP needs capabilities not fully specified in HAKMEM, but still aligned with its model.
3. **Modern repair** — used when later work addresses an operational weakness in the original algorithms, such as non-productivity or silent infinite looping.

The currently identified conflicts are:

- **Internal representation conflict**: HAKMEM 101B is formulated for continued fractions whose terms are integer pairs `(p_i, q_i)`, with regular continued fractions as the special case `q_i = 1`. Many modern tutorials simplify the model to regular continued fractions only. For this library, the **internal model must use generalized `(p, q)` streams for all sources, transducers, and kernels**, while **only the final client-facing emission is regular continued fraction (RCF)**.
- **Function-source conflict**: HAKMEM states that numbers like `pi` and `e`, and functions like `sin` and `tanh`, can be produced by short programs as continued-fraction terms. Many modern retellings instead start from floating-point values or from precomputed regular terms. For this library, **programmatic generators and transducers are mandatory** and must not be reduced to float-to-CF conversion.
- **Numeric-type conflict**: HAKMEM suggests that interval endpoint evaluation may be accelerated with floating-point estimates. This specification rejects that optimization for MVP correctness. All transform coefficients, terms, interval numerators, and interval denominators must use **arbitrary-precision integers**, with exact rationals for bounds.
- **Productivity conflict**: HAKMEM 101B gives a demand-driven output rule based on interval narrowing, but also acknowledges cases where an implementation can get stuck trying to emit the last term of a rational result when driven by irrational inputs. The MVP therefore preserves Gosper semantics for final RCF emission, but adopts a **Modern repair** stance operationally: the engine is required to keep refining exact intervals when the next term is not yet provable, and to fail explicitly when client-specified resource thresholds are exceeded rather than looping silently.
- **Transcendental-on-irrational-input conflict**: HAKMEM explicitly says transcendental functions of irrational arguments are awkward and suggests symbolic term-producing subroutines. The MVP requires `tanh(sqrt(5))` and `sqrt(3/pi^2 + e)`, so the design must **generalize HAKMEM's procedural idea into first-class unary operator transducers**.
- **Quadratic-kernel extension**: In addition to homographic and bihomographic kernels, the library requires a diagonal quadratic kernel for single-input forms `X' = (aX^2 + bX + c)/(dX^2 + eX + f)`. This is **Generalized beyond HAKMEM** and is driven by the MVP operator set and the priority of making `sqrt` practical.

For an implementable MVP, the library is assumed to be:

- a reusable programming library rather than a CLI or service;
- language-agnostic at the specification level, with a Go-specific API mapping later in the document;
- centered on exact, lazy arithmetic over continued-fraction term sources;
- capable of producing terms for the concrete expression:

  `sqrt(3/pi^2 + e) / (tanh(sqrt(5)) - sin(69°))`

This expression forces the MVP beyond a pure rational-only core. The resulting requirements therefore include both Gosper arithmetic and a minimal special-function term-source framework sufficient to support the expression above.

## Requirements

### Must Have

- Accept generalized continued-fraction (GCF) input streams whose terms are integer pairs `(p_i, q_i)`, with infinite streams as the normal case and `EOF` used to signal exact finite exhaustion when applicable.
- Use GCF `(p, q)` streams for **all internal sources, adapters, generators, transducers, and kernels**.
- Emit only **regular continued-fraction (RCF) output terms** to clients.
- Use **arbitrary-precision integers** for all GCF terms, all homographic coefficients, all bihomographic coefficients, all diagonal-kernel coefficients, and all intermediate coefficient updates.
- Use **exact rationals backed by arbitrary-precision integers** for interval endpoints, error bounds, and floor-stability checks.
- Use **no binary floating point** and no fixed-width integer arithmetic in the core algorithm or in interval checking.
- Implement comparison of continued fractions without conversion to floating point.
- Implement **homographic transforms** `(ax + b) / (cx + d)` as a first-class primitive.
- Implement **bihomographic transforms** `(axy + bx + cy + d) / (exy + fx + gy + h)` as the core primitive for binary arithmetic.
- Implement a **diagonal quadratic kernel** for single-input transforms `X' = (aX^2 + bX + c)/(dX^2 + eX + f)`.
- Implement Gosper-style binary arithmetic for `+`, `-`, `*`, and `/`, with output terms emitted incrementally and input terms requested only when necessary.
- Define deterministic source-exhaustion semantics for exact finite operands, with `EOF` as the canonical end-of-stream signal and HAKMEM-style tail interpretation applied by the kernels.
- When a computation yields a finite exact RCF output, the emitter must return the final term normally and then return `io.EOF` immediately on the next pull.
- Define canonical normalization for finite RCF outputs, including deterministic treatment of the trailing-`1` equivalence.
- Normalize kernel coefficient tuples by dividing by a nontrivial `gcd` when doing so preserves the represented transform.
- Provide conversion from emitted RCF prefixes into convergents `(P_n, Q_n)` and exact rational prefixes.
- Expose the current exact interval as a sibling of next-term emission.
- Provide a **term-source interface** for procedurally generated GCF streams, not just precomputed sequences.
- Provide built-in constant generators for `pi` and `e`.
- Provide adapters to create internal GCF streams from:
  - infinite RCF streams, by mapping each regular term `a_i` to `(a_i, 1)`;
  - BigInt-based rational values;
  - `int64` values by first lifting them to BigInt and then through the rational/GCF adapter.
- Provide unary operators for `sqrt`, `tanh`, `sin`, `square`, and `reciprocal`.
- Treat `square` as a distinct unary operator, not merely as binary multiplication with duplicated input consumption.
- Support exact angle input for `sin`, including a degree-based path sufficient for `69°`.
- Convert degrees to radians exactly before dispatching into the radian-domain `sin(x)` operator path.
- Ensure the public API can evaluate composed expressions lazily, rather than forcing eager expansion of all subexpressions.
- Do not promise next-term productivity. Instead, guarantee interval refinement until either a next term is provable or client-specified resource thresholds are exceeded.
- Provide resource guards for non-productive or pathologically slow generators and computations, including client-specified time and memory thresholds.
- Return an explicit error when configured resource thresholds are exceeded.
- Support MVP expression evaluation sufficient to emit RCF terms for:

  `sqrt(3/pi^2 + e) / (tanh(sqrt(5)) - sin(69°))`

- Include a conformance test suite with examples taken from HAKMEM 101A/101B where possible.
- Include property-based tests for rational identities, reciprocal and square round-trips where defined, and monotone refinement of intervals and approximants.
- Document each requirement or algorithmic choice as one of: **Faithful to HAKMEM**, **Generalized beyond HAKMEM**, or **Modern repair**.

### Should Have

- Support periodic continued fractions as a dedicated representation for quadratic irrationals.
- Support decimal-digit emission as a consumer over the exact continued-fraction pipeline.
- Expose debugging state for transform coefficients, interval bounds, gcd reductions, and source-consumption decisions.
- Support symbolic-sharing semantics so `square(x)` need not duplicate and independently consume the same source twice.
- Allow callers to choose whether angle arguments are expressed in radians or degrees, while keeping an exact degree helper for the MVP use case.
- Support a configurable maximum number of emitted RCF terms as an additional resource guard.

### Could Have

- Provide unary operator support beyond the MVP, such as `exp`, `log`, and inverse trigonometric functions, when exact generator strategies are available.
- Provide expression graphs that memoize shared subexpressions during lazy evaluation.
- Provide visualization helpers for emitted terms, convergents, and narrowing intervals.
- Provide adapters for arbitrary-precision integer ecosystems in multiple languages.
- Provide optional serialization for GCF fixtures, RCF outputs, and periodic forms for reproducible tests, debugging transcripts, and saved examples.

### Won’t Have (MVP)

- No reliance on binary floating point for core arithmetic decisions or exactness claims.
- No requirement to reconstruct transcendental continued fractions from sampled numeric values.
- No promise to support every generalized continued-fraction form in the MVP.
- No guarantee that every arbitrary infinite generator is productive without caller-supplied limits.
- No distributed or parallel execution requirement in the MVP.

## Method

### 1. Core Architecture

The library is organized as a lazy exact-evaluation pipeline with a hard boundary:

- **All internal sources, adapters, generators, transducers, and kernels consume and produce GCF `(p, q)` streams or exact interval state derived from them.**
- **Only the final client-facing emitter regularizes to RCF terms `a_i`.**
- **Decisions about when an RCF term may be emitted** are made only from exact interval analysis over bigint-backed rationals.

#### Primary components

1. **GCF Source Interface**  
   Pull-based source of `(p, q)` terms. Sources include literal generators, constant generators, operator transducers, adapters over finite exact values, and adapters over infinite RCF streams.

2. **Homographic Kernel**  
   Maintains bigint coefficients `(a, b, c, d)` for unary linear-fractional transforms such as reciprocal and affine-rational transforms.

3. **Bihomographic Kernel**  
   Maintains bigint coefficients `(a, b, c, d, e, f, g, h)` for binary operations `+`, `-`, `*`, and `/`.

4. **Diagonal Quadratic Kernel**  
   Maintains bigint coefficients `(a, b, c, d, e, f)` for single-input transforms `X' = (aX^2 + bX + c)/(dX^2 + eX + f)`. This kernel is internal and is free to operate as a GCF-producing transducer, an interval-refining kernel, or a hybrid, whichever best supports the `sqrt` implementation and other unary quadratic transforms.

5. **Unary Operator Transducers**  
   Lazy exact transducers for `sqrt`, `tanh`, `sin`, `square`, and `reciprocal`. These consume one GCF input source and expose a new GCF source or interval-refining internal state suitable for later RCF emission.

6. **Interval Engine**  
   Represents bounds as exact rationals `(num, den)` with bigint numerator and denominator, plus inclusion polarity where needed. It determines whether the current image of a transform lies wholly inside a single floor bucket.

7. **RCF Emitter**  
   Emits the next regular term only when the exact floor is proven stable across the full current interval image.

8. **Expression Planner**  
   Builds a lazy graph for composed expressions so subexpressions are consumed only when demanded by downstream kernels.

### 2. Fundamental Numeric Rules

- Every coefficient in every transform is bigint.
- Every input term component `p` and `q` is bigint.
- Every interval endpoint is a rational with bigint numerator and denominator.
- No `float32`, `float64`, machine `int`, or `int64` may participate in correctness-critical logic.
- A **Möbius transform** is another name for a homographic or linear-fractional transform of the form `(ax + b)/(cx + d)`.
- Coefficient tuples may be normalized by dividing all coefficients by their nontrivial `gcd` whenever that preserves the represented transform.
- Any optional optimization must be observationally equivalent to the exact bigint/rational semantics.

### 3. Tail Interpretation, Intervals, and Output Rule

**HAKMEM-style tail interpretation applied by the kernels** means the following.

When a finite GCF source is exhausted, the kernel does not treat exhaustion as an immediate numerical failure. Instead it applies the exact limiting substitution described by HAKMEM: an exhausted continued fraction is interpreted as having become infinite at its unread tail.

For the bihomographic form

`z(x,y) = (axy + bx + cy + d)/(exy + fx + gy + h)`

exhausting `x` maps the state to the limiting form

`z(∞, y) = (0·xy + 0·x + ay + b)/(0·xy + 0·x + ey + f)`

and exhausting `y` maps it to

`z(x, ∞) = (0·xy + ax + 0·y + c)/(0·xy + ex + 0·y + g)`.

Operationally, the public source still returns `io.EOF`; the kernel then applies this exact limit rule internally before deciding whether more output can be produced.

The public API must expose:

- `NextRCF()` => next final client-visible RCF term, or finite exhaustion;
- `CurrentInterval()` => current exact rational enclosure of the remaining value state.

The output rule is:

- evaluate the transform over the current exact bounds of its unread input tails;
- derive the exact image interval or inside/outside interval set;
- emit `t = floor(z)` only when that floor is identical for the whole admissible image;
- otherwise continue refining exact interval knowledge.

This is **Faithful to HAKMEM** in its final-emission semantics and **Modern repair** in its operational contract: the engine does not promise to always produce the next RCF term, but it must continue producing tighter exact enclosures until either a term becomes provable or configured resource limits are exceeded.

### 4. Ingest / Produce State Machines

#### Bihomographic kernel

State is the coefficient tuple `(a,b,c,d,e,f,g,h)`.

Operations:

- **Ingest left term `(p,q)`**: update the eight coefficients by the HAKMEM substitution for `x = p + q/x'`.
- **Ingest right term `(r,s)`**: update the eight coefficients by the HAKMEM substitution for `y = r + s/y'`.
- **Produce output term `t`**: update the eight coefficients by the HAKMEM output substitution for `z = t + 1/z'`.
- **Normalize**: reduce coefficients by a shared nontrivial `gcd` when valid.

Initial states:

- `x + y`  => `(0,1,1,0, 0,0,0,1)`
- `x - y`  => `(0,1,-1,0, 0,0,0,1)`
- `y - x`  => `(0,-1,1,0, 0,0,0,1)`
- `x * y`  => `(1,0,0,0, 0,0,0,1)`
- `x / y`  => `(0,1,0,0, 0,0,1,0)`
- `y / x`  => `(0,0,1,0, 0,1,0,0)`

#### Homographic kernel

State is `(a,b,c,d)`.

Operations:

- ingest input term `(p,q)`;
- emit an internal GCF term or contribute refined interval knowledge when the state makes that possible;
- support final RCF emission only through the client-facing emitter;
- normalize coefficients by `gcd` when valid.

Primary uses in MVP:

- reciprocal;
- rational scaling and shifting;
- substeps inside unary transducers.

#### Diagonal quadratic kernel

State is `(a,b,c,d,e,f)` and represents `X' = (aX^2 + bX + c)/(dX^2 + eX + f)`.

Operations:

- ingest input term `(p,q)`;
- refine exact interval knowledge under the quadratic map;
- optionally emit internal GCF terms if that best supports the active operator strategy;
- normalize coefficients by `gcd` when valid.

Primary uses in MVP:

- single-input quadratic forms;
- `square` and related transforms;
- support machinery for `sqrt`, including Newton-style internal strategies when exactness is preserved.

### 5. Unary Operator Strategy

#### Reciprocal

Implemented as a homographic transform. This is **Faithful to HAKMEM**.

#### Square

Implemented as a dedicated single-source operator with source sharing, not by blindly wiring the same input stream into both sides of a binary multiply.

This operator is semantically distinct because a single-variable quadratic form is operationally different from a generic bihomographic call with duplicated inputs.

The diagonal quadratic kernel is the preferred primitive for this family of transforms.

#### Sqrt

Implemented as a dedicated unary transducer using exact rational bounds and exact floor tests. This is **Generalized beyond HAKMEM**.

The diagonal quadratic kernel exists primarily to make `sqrt` practical. The implementation may use Newton-style internal refinement, diagonal-kernel interval refinement, or a hybrid strategy, provided that:

- all correctness-critical arithmetic remains exact;
- all internal observable numeric state remains GCF- or exact-interval-based;
- final client-visible output is still RCF only.

#### Sin and Tanh

Implemented as symbolic unary transducers that consume an exact input source lazily and expose a continued-fraction-producing or interval-refining interface downstream.

This follows HAKMEM's procedural model because HAKMEM explicitly recommends subroutines that produce symbolic terms on demand. It goes beyond the original text because the MVP requires those operators to work concretely on nontrivial irrational inputs such as `tanh(sqrt(5))` and `sqrt(3/pi^2 + e)`.

### 6. Termination and Productivity

The library must distinguish three behaviors:

1. **Productive** — a next RCF term is proven and emitted.
2. **Undecided but refining** — more source terms are required before emission is possible, but the exact interval is still being tightened.
3. **Resource-exhausted** — the engine has not proven the next term before a client-specified time, memory, or other configured budget is exceeded.

The MVP adopts the stronger rule that the engine does **not** promise next-term productivity.

Instead, it promises:

- exact interval refinement while computation remains within configured budgets;
- final RCF emission whenever the next term becomes provable;
- explicit failure when configured resource limits are exceeded.

The modern repair path for ambiguous near-integer cases is therefore:

- if the current exact enclosure crosses an integer boundary, the engine must not guess the next RCF term;
- it must continue refining the enclosure using additional source terms or operator-specific exact refinement;
- if the enclosure still does not isolate a unique floor before budgets are exhausted, the API must surface the current exact enclosure and return a resource-limit error rather than looping silently.

### 7. Go-specific API Mapping

The language-agnostic model should map cleanly to Go as follows:

- bigint => `math/big.Int`
- rational => `math/big.Rat`
- term => struct `{ P *big.Int; Q *big.Int }`
- RCF term => `*big.Int`
- source => pull iterator interface, context-aware, using `io.EOF` as the exact graceful exhaustion signal for finite sources
- kernels => mutable structs owning reusable scratch bigints to reduce allocation churn
- time/resource budgets => constructor or evaluator options, with context deadline/cancellation used for wall-clock control

Suggested Go interfaces:

```go
type GCFPair struct {
    P *big.Int
    Q *big.Int
}

type Interval struct {
    Lo *big.Rat
    Hi *big.Rat
}

type GCFSource interface {
    Next(ctx context.Context) (GCFPair, error) // returns io.EOF unwrapped on graceful finite exhaustion
}

type RCFSource interface {
    NextRCF(ctx context.Context) (*big.Int, error) // may also return io.EOF for finite exact outputs
    CurrentInterval(ctx context.Context) (Interval, error)
}
```

### 8. Component Diagram

```plantuml
@startuml
component "GCF Source
(bigint (p,q) stream)" as GCF
component "Unary Transducer
(sqrt/sin/tanh/square/reciprocal)" as U
component "Bihomographic Kernel
(a..h : bigint)" as B
component "Diagonal Kernel
(a..f : bigint)" as D
component "Homographic Kernel
(a..d : bigint)" as H
component "Interval Engine
(big.Rat bounds)" as I
component "RCF Emitter
(bigint a_i stream)" as R

GCF --> U
GCF --> B
U --> H
U --> D
U --> B
B --> I
D --> I
H --> I
I --> R
R --> H
R --> D
R --> B
@enduml
```

## Implementation

### Phase 1 — Numeric substrate

- Implement bigint coefficient containers for homographic, bihomographic, and diagonal-kernel states.
- Implement exact rational interval utilities using normalized bigint numerator/denominator pairs.
- Implement exact floor/ceil and interval-image helpers for Möbius, bilinear, and diagonal quadratic transforms.
- Implement coefficient-tuple `gcd` normalization.
- Define error taxonomy: invalid term, zero denominator, division by exact zero, resource limit exceeded, context cancellation, and graceful `io.EOF`.

### Phase 2 — Streaming contracts

- Implement the GCF source contract and finite-source exhaustion semantics.
- Implement the RCF emitter contract, including the rule that a finite exact result returns its last term and then `io.EOF` on the next call.
- Implement adapters for infinite RCF streams, BigInt rationals, `int64` values, finite GCF fixtures, and test harness sources.
- Add transcript logging hooks for consumed input terms, emitted RCF terms, coefficient transitions, gcd reductions, and interval snapshots.

### Phase 3 — Core kernels

- Implement homographic ingest and update transitions over bigint coefficients.
- Implement bihomographic ingest-left, ingest-right, and output transitions over bigint coefficients.
- Implement diagonal-kernel ingest and interval-update transitions over bigint coefficients.
- Implement exact floor-stability checks on transformed tail intervals.
- Implement exact exhausted-tail substitution rules after `io.EOF` from finite sources.
- Implement resource accounting hooks for time, memory, and optional emitted-term limits.

### Phase 4 — Constants and unary operators

- Implement built-in generators for `pi` and `e` as GCF-producing sources.
- Implement `reciprocal` as a homographic operator.
- Implement `square` as a dedicated single-source operator with source sharing, preferably on top of the diagonal kernel.
- Implement `sqrt` as a dedicated unary transducer using exact rational bounds and exact refinement, with the diagonal kernel available for Newton-style or hybrid internal strategies.
- Implement `sin` and `tanh` as exact unary transducers over continued-fraction sources, with no float-based fallback.
- Implement exact degree-to-radian conversion for `sin(69°)` before dispatch into the radian-domain `sin` path.

### Phase 5 — Expression engine

- Build a lazy expression planner that wires sources, unary transducers, binary kernels, and interval tracking into a demand-driven graph.
- Ensure subexpressions are evaluated only when downstream consumers request more information.
- Ensure shared unary inputs can be memoized where single-source semantics matter, especially for `square`.

### Phase 6 — Verification

- Add HAKMEM-derived golden tests for elementary operations, canonical trailing-`1` handling, and early term emission.
- Add tests for finite exhaustion behavior, including last-term-then-`io.EOF`.
- Add property-based tests for rational identities, reciprocal involution where defined, square monotonicity on nonnegative intervals, and convergent/interval refinement.
- Add targeted tests for the MVP expression `sqrt(3/pi^2 + e) / (tanh(sqrt(5)) - sin(69°))` to verify that terms begin streaming without eager whole-expression evaluation.
- Add regression tests for historically non-productive edge cases and resource-limit behavior.

### Phase 7 — Go packaging

- Publish a minimal Go package surface centered on pull iterators, `math/big.Int`, `math/big.Rat`, `context.Context`, `io.EOF`, `NextRCF`, and `CurrentInterval`.
- Reuse receiver-owned scratch objects to reduce allocation churn in hot coefficient-update paths.
- Keep interfaces small and place resource budgets behind optional configuration structs.

## Milestones

1. **M1 — Exact numeric core**  
   Bigint coefficient types, exact rational intervals, gcd normalization, floor-stability checks, and canonical error model complete.

2. **M2 — Streaming kernels**  
   Homographic, bihomographic, and diagonal kernels pass golden tests for `+`, `-`, `*`, `/`, reciprocal, finite exhaustion semantics, and interval refinement.

3. **M3 — Unary MVP operators**  
   `square`, `sqrt`, `sin`, and `tanh` integrated as lazy exact transducers; `pi` and `e` generators available.

4. **M4 — MVP expression success**  
   The library emits RCF terms for `sqrt(3/pi^2 + e) / (tanh(sqrt(5)) - sin(69°))` through the public API and exposes exact intervals throughout evaluation.

5. **M5 — Hardening**  
   Property-based tests, regression suite for non-productive cases, resource-guard enforcement, tracing, and Go packaging polish complete.

## Gathering Results

The implementation will be considered successful if it demonstrates all of the following:

- elementary arithmetic over GCF inputs emits correct RCF outputs matching HAKMEM-derived examples;
- all correctness-critical paths use bigint and exact rational arithmetic only;
- finite exact outputs return the final term and then `io.EOF` immediately on the next pull;
- unary MVP operators compose lazily and correctly with the binary and diagonal kernels;
- the MVP expression begins producing terms incrementally without eager global evaluation;
- interval bounds narrow monotonically and remain sound for every emitted term;
- known non-productive scenarios are surfaced as interval refinement followed by either provable emission or explicit resource-limit failure.

Primary evaluation artifacts:

- golden-term transcripts;
- convergent and bound traces;
- allocation and latency benchmarks per emitted term;
- regression corpus for ambiguous or historically looping cases.

## Unanswered Questions

- Whether a maximum emitted-term count should be a required MVP guard or remain optional configuration.