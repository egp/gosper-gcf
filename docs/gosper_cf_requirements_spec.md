# SPEC-1-Gosper-Continued-Fraction-Arithmetic-Library

**Version:** 6  
**Timestamp:** 2026-04-03T21:00:00-04:00

## Background

This document defines the requirements for a software library implementing Gosper-style continued-fraction arithmetic, grounded primarily in HAKMEM Item 101A (representation), Item 101B (continued-fraction arithmetic), and Item 101C (smallest rational in an interval).

The primary motivation is Gosper’s own: arithmetic should proceed exactly, incrementally, and on demand, producing continued-fraction output terms as soon as they are forced by the known input terms. The library therefore prioritizes:

- exact arithmetic over approximation;
- demand-driven production of output terms;
- exact interval tracking as part of the runtime model;
- composition of operators over continued-fraction streams rather than conversion through floating point.

HAKMEM 101B explicitly uses a composed expression involving `sqrt`, `pi`, `e`, `tanh`, and `sin` as its motivating example. This specification adopts a closely related MVP target expression:

`sqrt(3/pi^2 + e) / (tanh(sqrt(5)) - sin(69°))`

The degree symbol is a deliberate project choice. HAKMEM writes `sin(69)` without specifying units. This project interprets that example as degrees for the client-facing helper, while still defining the core operator in radians and requiring exact degree-to-radian conversion before dispatch.

## Source Fidelity and Conflict Policy

When the original HAKMEM text and modern implementation practice differ, this specification marks the choice explicitly with one of these dispositions:

1. **Faithful to HAKMEM**  
   Used when HAKMEM is clear and directly implementable.

2. **Project restriction / modern repair**  
   Used when the project intentionally adopts a stricter operational contract than HAKMEM, typically to avoid silent looping or retractable client-visible output.

3. **Project generalization**  
   Used when the project adds convenience or API structure not dictated by HAKMEM, while remaining consistent with Gosper’s model.

### Explicit conflicts and resolutions

- **Internal representation conflict** — **Faithful to HAKMEM**  
  HAKMEM 101B is formulated for generalized continued fractions whose terms are integer pairs `(p_i, q_i)`, with regular continued fractions as the special case `q_i = 1`.  
  **Resolution:** all internal sources, adapters, generators, transducers, and kernels use generalized `(p, q)` streams. Only the final client-facing emission is regular continued fraction (RCF).

- **Function-source conflict** — **Faithful to HAKMEM**  
  HAKMEM describes constants and functions such as `pi`, `e`, `sin`, and `tanh` as procedurally generated continued-fraction sources rather than values reconstructed from floating point.  
  **Resolution:** programmatic generators and transducers are mandatory.

- **Numeric-type conflict** — **Project restriction**  
  HAKMEM allows floating-point estimates to facilitate endpoint evaluation.  
  **Resolution:** this project uses exact big integers and exact rational bounds for all correctness-critical interval work.

- **Retractable-output conflict** — **Project restriction**  
  HAKMEM explicitly allows temporary wrong terms followed by recanting correction terms.  
  **Resolution:** this project forbids speculative client-visible output. A client-visible RCF term may be emitted only when exact interval analysis proves it.

- **Productivity conflict** — **Project restriction / modern repair**  
  HAKMEM acknowledges that a computation may get stuck trying to emit the last rational term when driven by irrational inputs.  
  **Resolution:** the engine is not required to guarantee next-term productivity, but it is required to keep refining exact interval knowledge until either:
  - the next term becomes provable, or
  - configured resource limits are exceeded, in which case it must fail explicitly.

- **Quadratic single-input form** — **Faithful to HAKMEM with stronger certification**  
  HAKMEM 101B explicitly introduces the diagonal quadratic form  
  `(aX^2 + bX + c) / (dX^2 + eX + f)`,  
  calls it preserved by continued-fraction term transactions, says it is more economical than `z(x,x)`, and says it is essential for Newton-feedback tricks. HAKMEM also notes that it is not guaranteed monotone and therefore weakens the interval-check justification.  
  **Resolution:** this project includes a diagonal quadratic kernel as part of the core model, but because client-visible retractions are forbidden, any client-visible output driven by this kernel must satisfy exact certification rules stronger than Gosper’s permissive self-correction stance.

- **Angle-unit choice for `sin(69)`** — **Project generalization**  
  HAKMEM does not specify units for `sin(69)`.  
  **Resolution:** the library defines:
  - `sinRadians(x)` as the core operator;
  - `sinDegrees(x)` as an exact wrapper that converts degrees to radians before dispatch.

- **Decimal emission scope** — **Faithful to HAKMEM, deferred in MVP**  
  HAKMEM describes decimal-digit emission by multiplying by `10` instead of reciprocating after outputting the current floor term.  
  **Resolution:** decimal emission is a planned capability but not an MVP requirement.

## Requirements

### Must Have

- **[Faithful to HAKMEM]** Accept generalized continued-fraction input streams whose terms are integer pairs `(p_i, q_i)`.
- **[Faithful to HAKMEM]** Use GCF `(p, q)` streams for all internal sources, adapters, generators, transducers, and kernels.
- **[Faithful to HAKMEM]** Emit only regular continued-fraction (RCF) output terms to clients.
- **[Project restriction]** Use arbitrary-precision integers for:
  - all GCF terms;
  - all BLFT coefficients;
  - all diagonal-kernel coefficients;
  - all intermediate coefficient updates.
- **[Project restriction]** Use exact rationals backed by arbitrary-precision integers for interval endpoints, exact bounds, and floor-stability checks.
- **[Project restriction]** Use no binary floating point and no fixed-width integer arithmetic in correctness-critical logic.
- **[Faithful to HAKMEM]** Implement comparison of continued fractions without conversion to floating point.
- **[Faithful to HAKMEM]** Implement a unified bihomographic kernel with bigint coefficients `(a, b, c, d, e, f, g, h)` for
  `z(x,y) = (axy + bx + cy + d) / (exy + fx + gy + h)`.
- **[Faithful to HAKMEM]** Support unary homographic behavior as a specialization of the same kernel rather than requiring a distinct runtime ULFT type.
- **[Faithful to HAKMEM]** Implement the diagonal quadratic kernel with bigint coefficients `(a, b, c, d, e, f)` for
  `z(x) = (aX^2 + bX + c) / (dX^2 + eX + f)`.
- **[Faithful to HAKMEM]** Implement Gosper-style binary arithmetic for `+`, `-`, `*`, and `/`, with output terms emitted incrementally and input terms requested only when needed.
- **[Faithful to HAKMEM]** Define deterministic source-exhaustion semantics for exact finite operands, with `EOF` as the public end-of-stream signal and HAKMEM-style tail interpretation applied internally by the kernels.
- **[Faithful to HAKMEM]** When one BLFT input stream is exhausted, evaluation continues in the same 8-coefficient representation, with the exhausted input folded into the coefficients using the exact HAKMEM limit rule.
- **[Faithful to HAKMEM]** When a computation yields a finite exact RCF output, the emitter returns the final term normally and then returns `io.EOF` on the next pull.
- **[Project restriction]** Define canonical normalization for finite RCF outputs, including deterministic treatment of the trailing-`1` equivalence.
- **[Faithful to HAKMEM]** Normalize coefficient tuples by dividing by a nontrivial `gcd` whenever doing so preserves the represented transform.
- **[Faithful to HAKMEM]** Provide conversion from emitted RCF prefixes into convergents `(P_n, Q_n)` and exact rational prefixes.
- **[Faithful to HAKMEM]** Expose the current exact interval as a sibling of next-term emission.
- **[Faithful to HAKMEM]** Provide a term-source abstraction for procedurally generated GCF streams, not only precomputed sequences.
- **[Faithful to HAKMEM]** Provide built-in constant generators for `pi` and `e`.
- **[Faithful to HAKMEM]** Provide adapters to create internal GCF streams from:
  - infinite RCF streams, by mapping each regular term `a_i` to `(a_i, 1)`;
  - exact rational values;
  - integer values by exact lifting.
- **[Faithful to HAKMEM / project generalization]** Provide unary operators for:
  - `sqrt`
  - `tanh`
  - `sinRadians`
  - `sinDegrees`
  - `square`
  - `reciprocal`
- **[Faithful to HAKMEM]** Treat `square(x)` as a distinct unary operator, not merely as binary multiplication with duplicated input consumption.
- **[Project generalization]** `sinDegrees(x)` must convert degrees to radians exactly before dispatching to `sinRadians(x)`.
- **[Project restriction]** Ensure the public API can evaluate composed expressions lazily rather than forcing eager expansion.
- **[Project restriction / modern repair]** Do not promise next-term productivity.
- **[Project restriction / modern repair]** Guarantee continued exact interval refinement until either:
  - a next term is provable, or
  - resource thresholds are exceeded.
- **[Project restriction / modern repair]** Provide resource guards for pathologically slow or non-productive computations, including caller-configurable limits on time, work, memory, and optionally emitted terms.
- **[Project restriction / modern repair]** Return an explicit error when configured resource thresholds are exceeded.
- **[Faithful to HAKMEM, strengthened]** A client-visible term may be emitted only when the exact admissible image lies wholly inside a single floor bucket.
- **[Project restriction]** No client-visible speculative or retractable terms are permitted.
- **[Faithful to HAKMEM, strengthened]** Any source of GCF terms must also provide the strongest exact interval it can currently justify for its unread tail.
- **[Project restriction]** A source may use read-ahead to tighten its interval only if:
  - the read-ahead is observationally pure;
  - it does not violate laziness beyond the amount needed to justify the tighter bound;
  - the resulting interval is exact with respect to the information consumed.
- **[Project restriction]** If a source can determine the unread tail exactly at the current frontier, its interval may collapse to a point.
- **[Faithful to HAKMEM]** Represent interval enclosures using inside/outside semantics compatible with Gosper’s toroidal/projective view.
- **[Faithful to HAKMEM]** Implement exact denominator-crossing logic when mapping intervals through transforms.
- **[Faithful to HAKMEM]** For BLFT range propagation, evaluate the four endpoint/corner-derived ranges and combine them by exact union semantics.
- **[Faithful to HAKMEM]** Use Gosper’s widest-range idea, or an exact equivalent, to decide which input stream to ingest when no output term is yet forced.
- **[Faithful to HAKMEM]** Support MVP expression evaluation sufficient to emit RCF terms for  
  `sqrt(3/pi^2 + e) / (tanh(sqrt(5)) - sin(69°))`.
- **[Faithful to HAKMEM 101C / project utility]** Provide an exact utility that, given a nonempty interval, returns a rational contained in that interval with the smallest numerator and denominator in the Gosper 101C sense.
- **[Project utility]** Make that 101C interval-rational utility available for testing, diagnostics, and canonical witness generation.
- **[Project quality requirement]** Include conformance tests with examples derived from HAKMEM where practical.
- **[Project quality requirement]** Include property-based tests for:
  - rational identities;
  - reciprocal and square round-trips where defined;
  - monotone interval refinement;
  - exact containment invariants;
  - agreement between certified output and interval floors.
- **[Project documentation requirement]** Every major requirement or algorithmic choice must be labeled as Faithful to HAKMEM, Project restriction / modern repair, or Project generalization.

### Should Have

- **[Faithful to HAKMEM, deferred]** Support decimal-digit emission as a consumer over the exact continued-fraction pipeline by the HAKMEM multiply-by-10 rule.
- **[Project utility]** Support periodic continued fractions as a dedicated representation for quadratic irrationals.
- **[Project utility]** Expose debugging state for:
  - transform coefficients;
  - exact interval bounds;
  - openness flags;
  - denominator-crossing decisions;
  - coefficient reductions;
  - source-consumption decisions.
- **[Project utility]** Support symbolic-sharing semantics so `square(x)` need not independently consume two copies of the same source.
- **[Project utility]** Support a configurable maximum number of emitted RCF terms as an additional resource guard.
- **[Project utility]** Provide visualization helpers for emitted terms, convergents, and narrowing intervals.

### Could Have

- **[Project generalization]** Provide unary operators beyond the MVP, such as `exp`, `log`, and inverse trigonometric functions, when exact generator strategies are available.
- **[Project generalization]** Provide expression graphs that memoize shared subexpressions during lazy evaluation.
- **[Project utility]** Provide optional serialization for GCF fixtures, RCF outputs, periodic forms, and debugging transcripts.

### Won’t Have (MVP)

- No reliance on binary floating point for core arithmetic decisions or exactness claims.
- No requirement to reconstruct transcendental continued fractions from sampled numeric values.
- No promise to support every generalized continued-fraction invariant form in the MVP.
- No guarantee that every arbitrary infinite generator is productive without caller-supplied limits.
- No distributed or parallel execution requirement in the MVP.

## Method

### 1. Core Architecture

The library is organized as a lazy exact-evaluation pipeline with a hard boundary:

- all internal sources, adapters, generators, transducers, and kernels consume and produce GCF `(p, q)` streams and/or exact interval state derived from them;
- only the final client-facing emitter regularizes to RCF terms `a_i`;
- decisions about when an RCF term may be emitted are made only from exact interval analysis over bigint-backed rationals.

#### Primary components

1. **GCF Source Interface**  
   Pull-based source of `(p, q)` terms, plus exact interval exposure for the unread tail.

2. **Unified BLFT Kernel**  
   Maintains bigint coefficients `(a, b, c, d, e, f, g, h)` for linear-fractional continued-fraction evaluation.

3. **Diagonal Quadratic Kernel**  
   Maintains bigint coefficients `(a, b, c, d, e, f)` for  
   `z(x) = (aX^2 + bX + c) / (dX^2 + eX + f)`.

4. **Unary Operator Transducers**  
   Lazy exact transducers for `sqrt`, `tanh`, `sinRadians`, `sinDegrees`, `square`, and `reciprocal`.

5. **Interval Engine**  
   Computes and refines exact admissible ranges using inside/outside interval semantics.

6. **RCF Emitter**  
   Emits the next regular term only when floor stability is exactly certified.

7. **Expression Planner**  
   Builds a lazy graph for composed expressions so subexpressions are consumed only when demanded downstream.

### 2. Fundamental Numeric Rules

- Every coefficient in every transform is bigint.
- Every input term component `p` and `q` is bigint.
- Every interval endpoint is an exact rational.
- No `float32`, `float64`, machine `int`, or fixed-width arithmetic may participate in correctness-critical logic.
- A Möbius transform is synonymous here with a homographic or linear-fractional transform.
- Coefficient tuples may be normalized by dividing all coefficients by a nontrivial common divisor whenever that preserves the represented transform.
- Any optimization must be observationally equivalent to the exact bigint/rational semantics.

### 3. Tail Interpretation

When a finite GCF source is exhausted, the kernel does not treat exhaustion as an immediate numerical failure. Instead it applies HAKMEM’s limiting interpretation: the unread tail has become infinite.

For the bihomographic form

`z(x,y) = (axy + bx + cy + d) / (exy + fx + gy + h)`

the limit rules are:

- exhausting `x` yields  
  `z(∞, y) = (0·xy + 0·x + ay + b) / (0·xy + 0·x + ey + f)`

- exhausting `y` yields  
  `z(x, ∞) = (0·xy + ax + 0·y + c) / (0·xy + ex + 0·y + g)`

Operationally, the public source returns `io.EOF`; the kernel then applies the exact limit rule internally and continues evaluation in the same BLFT representation.

### 4. Interval Model

Interval semantics follow Gosper’s inside/outside model, interpreted on the projective line.

A range consists of:

- two exact rational endpoints `Lo` and `Hi`;
- endpoint openness/closedness flags;
- a polarity bit `Inside` versus `Outside`.

#### Inside intervals

An inside interval denotes the connected admissible arc from `Lo` to `Hi` in the ordinary affine picture.

#### Outside intervals

An outside interval denotes the complement of the excluded gap between `Lo` and `Hi`; in the ordinary affine picture it behaves like the union of two rays. This is the projective/toroidal viewpoint Gosper references when he notes that plus and minus infinity are represented by the same empty continued fraction.

#### Exactness requirements

- Every reported interval must be a true enclosure of the unread tail or transform image.
- A point value is represented by a degenerate exact interval.
- If a tighter exact interval is already provable from consumed information, a weaker interval is not acceptable.
- Read-ahead is permitted only to justify a tighter exact interval, never to speculate.

### 5. Range Propagation Rules

#### BLFT / homographic-style propagation

For fixed `y`, map the current `x` interval through `z(x,y)`.

- If the current `x` range is inside and the denominator does not change sign over that range, the image stays inside between the endpoint images.
- If the denominator changes sign across the range, the image becomes outside.
- Symmetric rules apply with `x` and `y` interchanged.

For the full BLFT, the certified range is derived from the union of the four endpoint/corner-derived ranges associated with the extremes of `x` and `y`.

The implementation must preserve exactness of this enclosure. A looser but exact enclosure is acceptable when unavoidable; an enclosure known to omit admissible values is not.

#### Widest-range input choice

When no output term is forced, the engine should ingest from the source associated with the widest current range, or use an exact strategy equivalent in effect.

The intended ordering follows Gosper’s intuition:

`Outside narrowness > Outside wideness > Inside wideness > Inside narrowness`

The implementation may use any mathematically equivalent exact metric that preserves the same decision quality.

#### Diagonal quadratic propagation

Because the diagonal quadratic form is not guaranteed monotone, the project may not rely on Gosper’s permissive “it will probably self-correct” argument for client-visible output.

Instead, the diagonal kernel must satisfy one of these stronger conditions before any client-visible emission depending on it:

1. it proves an exact enclosing interval by sound range analysis; or
2. it refines downstream through a certified interval-only path until a client-visible term is provable.

No uncertified diagonal-kernel output may reach the public RCF stream.

### 6. Output Rule

The final output rule is:

1. evaluate the current transform over the current exact bounds of its unread input tails;
2. derive the exact admissible image as an inside/outside interval enclosure;
3. emit `t = floor(z)` only when that floor is identical over the entire admissible image;
4. otherwise continue refining exact interval knowledge.

This preserves Gosper’s final-emission idea while rejecting speculative client-visible terms.

The public API must expose at least:

- `NextRCF()` → next final client-visible RCF term, or finite exhaustion, or explicit resource-limit error;
- `CurrentInterval()` → current exact enclosure of the remaining value state.

### 7. Ingest / Produce State Machines

#### Unified BLFT kernel

State is the coefficient tuple `(a,b,c,d,e,f,g,h)`.

Operations:

- **Ingest left term `(p,q)`** by the HAKMEM substitution for `x = p + q/x'`.
- **Ingest right term `(r,s)`** by the HAKMEM substitution for `y = r + s/y'`.
- **Produce output term `t`** by the HAKMEM substitution for `z = t + 1/z'`.
- **Normalize** by valid common-divisor reduction.
- **Collapse left exhaustion** by the HAKMEM `x -> ∞` rule.
- **Collapse right exhaustion** by the HAKMEM `y -> ∞` rule.

Initial states include:

- `x + y` => `(0,1,1,0, 0,0,0,1)`
- `x - y` => `(0,1,-1,0, 0,0,0,1)`
- `y - x` => `(0,-1,1,0, 0,0,0,1)`
- `x * y` => `(1,0,0,0, 0,0,0,1)`
- `x / y` => `(0,1,0,0, 0,0,1,0)`
- `y / x` => `(0,0,1,0, 0,1,0,0)`

Unary linear-fractional transforms are represented by the same BLFT state, either as initial degenerate configurations or as the result of one-sided exhaustion.

#### Diagonal quadratic kernel

State is `(a,b,c,d,e,f)` and represents

`z(x) = (aX^2 + bX + c) / (dX^2 + eX + f)`

Operations:

- ingest input term `(p,q)`;
- refine exact interval knowledge under the quadratic map;
- optionally transform internally in whatever certified way best supports the operator strategy;
- normalize coefficients by valid common-divisor reduction.

Primary MVP uses:

- `square`;
- quadratic single-input transforms;
- internal machinery supporting `sqrt`, including Newton-style feedback strategies when exactness is preserved.

### 8. Exact Interval Utility from HAKMEM 101C

The library must include an exact interval utility implementing the 101C rule:

Given a nonempty interval, find a rational contained in it with the smallest numerator and denominator in Gosper’s sense.

Method:

1. express the interval endpoints as continued fractions;
2. find the first term where they differ;
3. increment the lesser term unless it is the last;
4. discard all terms to the right;
5. if one endpoint terminates while matching the other so far, append infinity and continue.

This utility is not a substitute for the streaming arithmetic engine. Its roles are:

- canonical witness selection for a certified interval;
- interval-debugging support;
- tests for exact interval machinery;
- optional simplification when a finite representative inside a certified interval is needed.

### 9. Transcendental and Symbolic-Term Support

HAKMEM notes that transcendental functions of irrational inputs are awkward because their continued-fraction terms may themselves be symbolic functions of the input.

For MVP purposes, this specification requires that such operators be expressible as first-class transducers or equivalent lazy exact sources. The implementation does not need to mimic Gosper’s exact subroutine sketch literally, but it must preserve the same demand-driven semantics:

- no float-to-CF reconstruction;
- no eager full expansion;
- no client-visible speculative terms.

### 10. Decimal Emission

Decimal emission is not required for MVP. However, the design must not preclude it.

The intended future rule is the HAKMEM one: after outputting the current floor term `t`, decimal digits are obtained by multiplying the residual transform by `10` instead of reciprocating.

### 11. Conformance Expectations

The implementation should be judged conformant to this specification when it satisfies all of the following:

- it performs exact continued-fraction arithmetic using generalized internal `(p,q)` streams;
- it preserves Gosper-style demand-driven ingest/emit behavior;
- it implements exact interval reasoning using inside/outside semantics;
- it supports the diagonal quadratic form needed for `square` and `sqrt`-oriented work;
- it forbids speculative client-visible output;
- it provides explicit failure on configured resource exhaustion;
- it includes the 101C rational-in-interval utility;
- it can lazily evaluate the MVP composed expression and emit certified RCF terms.

## Notes for Implementation

- Identity, project-X, and project-Y cases are not permission to bypass the GCF kernel by reinterpreting `(p,q)` terms as already-emitted RCF terms.
- Any source that already knows the next exact tail value should surface that fact through its interval interface.
- Interval exactness is part of the semantic contract, not merely a debugging aid.
- The toroidal/projective interpretation is central to outside-interval reasoning and denominator-crossing behavior.
- The project deliberately chooses stronger public-output correctness guarantees than Gosper’s permissive self-correction model.