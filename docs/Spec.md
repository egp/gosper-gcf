# Spec.md

# gosper-gcf / core Specification

## Status

This document is the current design contract for the `core` package and the companion package `named`.

When a design choice conflicts with Gosper / HAKMEM 101B, Gosper wins.

For GitHub rendering, displayed formulas in this file use fenced `math` blocks.

---

## 1. Mission

Build a mathematically correct, testable Go library for arithmetic on generalized continued fractions in the style of Gosper / HAKMEM 101B.

Primary goals:

- mathematical correctness
- simplicity
- exact arithmetic only
- DRY implementation
- design for test

Performance matters only after correctness and clarity.

### Target formula

```math
\frac{\sqrt{\frac{3}{\pi^2} + e}}{\tanh(\sqrt{5}) - \sin(69^\circ)}
```

This is the future celebratory target.

---

## 2. Package layout

### 2.1 `core`

The `core` package contains:

- generalized source-stream abstractions
- the client-facing regular-output `GCF`
- BLFT-centered transform engine
- range logic
- normalization and invariant enforcement
- config and status codes

### 2.2 `named`

The companion package is `named`.

`named` contains:

- named streams such as `Pi()`, `E()`, `Sqrt2()`, etc.
- constructors from `Rational`
- constructors from finite or procedural source descriptions
- wrappers that help clients build source streams for testing and use

`named` does not weaken the `core` rule that engine ingestion is always in generalized `(p,q)` terms.

---

## 3. Core design principles

### 3.1 Exact arithmetic only

No floats anywhere in production code, tests, or correctness helpers.

### 3.2 Simplicity

Prefer one clear algebraic path over multiple convenience paths.

### 3.3 DRY

ULFT, BLFT, and any future DLFT support must share helpers wherever algebraically sound.

### 3.4 One internal engine shape

Internally, prefer one BLFT-sized transform state. Unary and diagonal cases are specialized initializations or constrained views of the same machinery whenever practical.

### 3.5 Certified output only

The engine emits a regular continued-fraction term only when Gosper’s emission criterion says it is forced by the current range.

### 3.6 EOF is valid

EOF is not an error. EOF triggers algebraic collapse and computation continues if possible.

### 3.7 Malformed client-visible input is reported

Malformed source terms supplied by client-visible streams return a failure status as soon as noticed, whether at construction time or pull time.

### 3.8 Internal invariant failure is fatal

Internal impossible states panic.

### 3.9 Invalid use should be unrepresentable where possible

APIs should prefer states and signatures that make misuse hard or impossible.

---

## 4. Mathematical conventions

### 4.1 Generalized source-term convention

Every generalized source term uses:

```math
X = p + \frac{q}{X'}
```

This is the exact convention everywhere in `core`.

### 4.2 Emission convention

Emission uses Gosper’s regular tail rewrite:

```math
Z = t + \frac{u}{Z'} \quad\Longleftrightarrow\quad Z' = \frac{u}{Z - t}
```

For regular emitted output, `u = 1`.

### 4.3 Regular emitted output convention

Emitted output is a regular continued-fraction term sequence.

- the first emitted regular term may be negative
- every later emitted regular term is a standard regular CF term

### 4.4 Input sign policy for generalized source terms

For v1 source validation:

- `p` must be nonzero or zero as algebra permits
- `q` must be nonzero on every term
- the first term may be negative in either component
- after the first term, both `p` and `q` must be positive
- any later negative `p` or `q` is malformed input

This is a practical v1 source-validity rule. The engine algebra is not assumed to require this stronger restriction once a term has been admitted.

### 4.5 No internal `q == 1` shortcut path

Regular CFs are the special case `q = 1`, but core logic does not fork into a separate arithmetic path for regular CFs.

### 4.6 Exact value vs integer value

If a range has `lo == hi`, then the value is known exactly.

That does **not** imply the value is an integer. It may be any exact rational.

---

## 5. Numeric types

### 5.1 Big integers

Use `*big.Int` internally for:

- `PQTerm.p`
- `PQTerm.q`
- `RCFTerm`
- BLFT coefficients
- Rational numerator
- Rational denominator

### 5.2 Rational

`Rational` is exact and BigInt-based.

Expected operations include:

- normalization
- comparison
- sign
- arithmetic for endpoint evaluation
- exact conversion from integer

### 5.3 Exact conversion from integer

“Exact conversion from integer” means:

- integer `n` converts to rational `n / 1`
- no rounding
- no approximation
- no float intermediary

### 5.4 Int64 usage

Use `int64` only where overflow is impossible by construction.

No correctness-critical internal algorithm may rely on `int64` if overflow is possible.

---

## 6. Status codes

The project prefers status codes over plain `error` returns for operational APIs.

### 6.1 Status direction

Representative statuses include:

- `StatusOK`
- `StatusEOF`
- `StatusInvalidInput`

Exact names may change during cleanup.

### 6.2 Meaning

- `StatusOK`: a term was produced successfully
- `StatusEOF`: normal end of stream
- `StatusInvalidInput`: malformed client-visible source term or invalid constructor input

### 6.3 Limits and panic policy

Configured timeout and excessive bit-length are not ordinary statuses in v1.

They panic.

---

## 7. Public model

### 7.1 Input source stream

A generalized source stream provides:

- the next generalized term `(p,q)`
- the tail source `X'`
- the current range enclosure

Conceptual direction:

- `NextPQ() -> (PQTerm, tail, Status)`
- `Range() -> Range`

### 7.2 Client-facing `GCF`

In v1, `GCF` is the client-facing regular-output stream object, backed by one engine state plus zero, one, or two input source streams.

A client creates a `GCF` from:

- an initial transform plugboard
- zero, one, or two input source streams
- an optional config, otherwise defaults

Then calls:

- `NextRCF()` to obtain the next emitted regular term
- `Range()` to inspect the current enclosure
- `Cmp()` to compare values

### 7.3 No public `Emit()` method

There is no public `Emit()` method.

Emission is an internal event. To the client, emission simply appears as a successful `NextRCF()` result.

### 7.4 No decimal API in v1

Decimal output is not part of v1 and is intentionally omitted from the public API.

---

## 8. Client-facing operations

The client accesses arithmetic by selecting predefined plugboards.

Representative direction:

- `Add(X, Y)`
- `Sub(X, Y)`
- `Mul(X, Y)`
- `Div(X, Y)`
- `Reciprocal(X)`
- `Square(X)`
- later `Sin(X)`, `Tanh(X)`, `Sqrt(X)`

Operationally, these morph into choosing the proper transform initialization and then running the common engine.

There is no special core arithmetic path for division, addition, etc., beyond plugboard selection plus common engine logic.

---

## 9. Range model

### 9.1 Range kinds

A `Range` is one of:

- `InsideInterval`
- `OutsideInterval`

### 9.2 Endpoint type

Endpoints are exact `Rational`s.

Each endpoint has its own openness/closedness flag.

### 9.3 Inside interval meaning

An inside interval encloses values within bounded endpoints.

Invariant:

- `lo <= hi`

Interpretation:

- the value lies within the interval
- if `lo == hi`, the exact value is known

### 9.4 Outside interval meaning

An outside interval denotes the enclosure:

```math
(-\infty, lo] \cup [hi, +\infty)
```

with independently open/closed finite endpoints.

Canonical ordering details for outside intervals remain open.

### 9.5 Range comparison ordering

Use Gosper’s uncertainty ordering:

- inside narrow
- inside wide
- outside wide
- outside narrow

The intent is:

- earlier in the list means more certain / better / smaller uncertainty
- later in the list means less certain / worse / larger uncertainty

This ordering need not expose a numeric width to clients.

### 9.6 Range validity

`Range()` is defined only on live objects.

It does not return a status.

Invalid internal use panics.

---

## 10. Exact openness rules at integer boundaries

These rules matter for emission.

### 10.1 Inside intervals

For an inside interval, the next regular term is forced only if **every** value in the interval has the same floor.

Define:

- `lowerTerm = floor(lo)`
- `upperTerm = floor(hi)`, except:
  - if `hi` is open and `hi` is exactly an integer `n`, then `upperTerm = n - 1`

Then an inside interval certifies a unique next regular term iff:

- `lowerTerm == upperTerm`

Examples:

- `[1.2, 1.9]` emits `1`
- `[1.2, 2)` emits `1`
- `[1, 2)` emits `1`
- `[1, 2]` does **not** emit
- `(2, 2.5]` emits `2`

### 10.2 Outside intervals

For v1, `OutsideInterval` does not directly certify a unique regular output term.

Corner analysis and infinity analysis may still be used to construct or refine the current range, but an outside interval itself is not a direct emission certificate.

---

## 11. BLFT core engine

### 11.1 Primary internal state

The primary internal transform state is BLFT-shaped:

```math
Z(X,Y) = \frac{aXY + bX + cY + d}{eXY + fX + gY + h}
```

with BigInt coefficients.

ULFT and any future diagonal support begin as constrained or specialized cases of the same internal machinery where practical.

### 11.2 Core engine actions

At each step the engine does exactly one of:

1. ingest one `(p,q)` term from one input
2. emit one regular CF term
3. collapse on EOF from an input

The difficult decision is between ingest and emit.

---

## 12. Gosper ingest formulas

If `X` produces `(p,q)` with

```math
X = p + \frac{q}{X'}
```

then ingesting `X` updates BLFT coefficients to:

```math
(a,b,c,d;e,f,g,h)
\mapsto
(pa+c,\ pb+d,\ qa,\ qb;\ pe+g,\ pf+h,\ qe,\ qf)
```

If `Y` produces `(r,s)` with

```math
Y = r + \frac{s}{Y'}
```

then ingesting `Y` updates BLFT coefficients to:

```math
(a,b,c,d;e,f,g,h)
\mapsto
(ra+b,\ sa,\ rc+d,\ sc;\ re+f,\ se,\ rg+h,\ sg)
```

---

## 13. Gosper emission formula

To emit `(t,u)` so that

```math
Z = t + \frac{u}{Z'}
```

the BLFT coefficients update to:

```math
(a,b,c,d;e,f,g,h)
\mapsto
(ue,\ uf,\ ug,\ uh;\ a-te,\ b-tf,\ c-tg,\ d-th)
```

For regular emitted output, `u = 1`.

---

## 14. Collapse rules

### 14.1 EOF is not an error

EOF is represented by `StatusEOF` and causes algebraic collapse.

### 14.2 BLFT collapse when `X` hits EOF

BLFT collapses to the induced unary transform in `Y`.

### 14.3 BLFT collapse when `Y` hits EOF

BLFT collapses to the induced unary transform in `X`.

### 14.4 BLFT collapse when both hit EOF

The state becomes a terminal exact state and emits the rest directly.

### 14.5 ULFT collapse on EOF

When a ULFT loses its input due to EOF, it collapses directly and emits the rest.

### 14.6 Terminal exact state

A terminal exact state is a live regular-output state whose value is exactly known as a rational and which emits the remainder directly without further source ingestion.

---

## 15. Emit-vs-ingest decision rule

### 15.1 Primary rule

Follow Gosper’s current widest-corner-range ingest rule exactly.

Do not replace it with speculative “try ingest X and Y, compare future ranges” logic.

### 15.2 Emit only when forced

Emit only when the current range proves the next regular term is uniquely determined.

### 15.3 Otherwise ingest

If emission is not forced, ingest from the variable associated with the widest current corner range.

### 15.4 Tie-break

If a tie-break is needed, tie goes to `X`.

---

## 16. `CanEmitRCFTerm`

Model `CanEmitRCFTerm` as a predicate/helper, not the primary public API.

Conceptual direction:

- `CanEmitRCFTerm(r Range) -> (RCFTerm, bool)`

For v1:

- only `InsideInterval` can certify emission
- use the integer-boundary openness rules from Section 10
- `OutsideInterval` returns `false`

A state-transition diagram may be added later as explanatory documentation, but the core rule is the predicate.

---

## 17. Normalization and invariants

### 17.1 Transform normalization

After every state mutation:

- ingest
- emit
- collapse

normalize transform coefficients by gcd where possible.

### 17.2 Sign normalization

Apply a conservative canonical sign placement whenever algebraically sound.

### 17.3 Exactness

Normalization must preserve exact semantics.

### 17.4 Preconditions

Validate client-visible input early when practical.

### 17.5 Postconditions

Use invariant checks to catch engine bugs early.

### 17.6 Bit-length checks

If enabled by config, inspect all eight BLFT coefficients using `BitLen()` around each operation, just before or just after gcd normalization.

### 17.7 Panic policy

Panic on:

- internal invariant failure
- impossible internal algebraic states
- configured timeout breach
- configured excessive bit-length breach

Do not panic for normal EOF.

---

## 18. Config

The client can supply config either:

- at creation time as an optional parameter, or
- via a set/get config object

Defaults must exist either way.

Representative config controls include:

- bit-length limit, default off
- timeout limit, default off

---

## 19. `Cmp`

### 19.1 Intent

`GCF.Cmp()` is important and should follow Gosper-style range reasoning.

### 19.2 Direction

Recommended semantics:

- compare `X` and `Y` by comparing the sign of `X - Y`
- operationally, construct the subtraction transform and refine until:
  - range proves negative
  - range proves zero exactly
  - range proves positive

### 19.3 Proposed result style

Direction:

- return a comparison code such as `-1`, `0`, `+1`
- plus a `Status`

Representative behavior:

- `(-1, StatusOK)` if `X < Y`
- `(0, StatusOK)` if equality is proven
- `(+1, StatusOK)` if `X > Y`
- non-OK statuses only for malformed input reaching comparison logic

### 19.4 Exact zero

If subtraction yields a range whose exact value is proven to be zero, `Cmp` returns equality.

---

## 20. Sources

HAKMEM / Gosper:
https://w3.pppl.gov/~hammett/work/2009/AIM-239-ocr.pdf

GitHub math rendering docs:
https://docs.github.com/en/get-started/writing-on-github/working-with-advanced-formatting/writing-mathematical-expressions

Go `math/big` docs:
https://pkg.go.dev/math/big

---

## 21. Open questions

The following questions remain open.

### 21.1 Final public type names

Whether temporary names such as `PQTerm`, `RCFTerm`, and `GCF` remain the final exported names.

### 21.2 Exact Go signatures

Especially:

- the exact signature of `NextPQ()`
- the exact signature of `GCF.NextRCF()`
- constructor signatures
- the exact config injection pattern

### 21.3 Exact `Range` struct layout

Including:

- endpoint structs
- inside/outside encoding
- canonical representation for outside intervals

### 21.4 Outside-interval canonicalization

We have settled the meaning, but not the exact internal ordering and invariants.

### 21.5 Full corner-range algorithm details

The high-level Gosper rule is fixed, but the precise implementation recipe for computing and comparing the four corner ranges still needs to be written down carefully.

### 21.6 Infinity analysis details

How `+\infty` and `-\infty` are used in range refinement and collapse-related reasoning still needs to be specified precisely.

### 21.7 `Cmp` final signature

The semantics are outlined, but the final exported signature is still open.

### 21.8 Config object shape

Creation-time optional config versus set/get config object still needs one concrete choice.

### 21.9 Future diagonal support

“BLFT but diagonal” remains deferred as a final guarantee question.