# SPEC-1-Gosper-Continued-Fraction-Arithmetic-Library

**Version:** 7 (Two-Tier Exactness Revision)  
**Timestamp:** 2026-04-03T22:00:00-04:00  

## 1. Background

This document defines the requirements for a software library implementing Gosper-style continued-fraction arithmetic, grounded in HAKMEM Items 101A, 101B, and 101C.

The primary motivation is Gosper’s own: arithmetic should proceed exactly, incrementally, and on demand. This library prioritizes:
* Exact arithmetic over approximation;
* Demand-driven production of output terms;
* Exact interval tracking via state-matrix coefficients;
* Composition of operators over generalized continued-fraction (GCF) streams.

## 2. Source Fidelity and Conflict Policy

### 2.1. Internal Representation: Faithful to HAKMEM
All internal sources, adapters, and kernels use generalized $(p, q)$ streams where a term represents $p + q/x'$. 
* **Resolution:** Regular Continued Fractions (RCF) where $q=1$ are treated as a special case of the GCF pipeline.

### 2.2. Numeric-Type: Project Restriction
* **Resolution:** The project uses exact arbitrary-precision integers (BigInts) for all coefficients. Floating-point or hardware-double types are forbidden for interval certification.

### 2.3. Retractable-Output: Project Synthesis (Two-Tier Architecture)
HAKMEM allows temporary "wrong" terms followed by recanting terms. To prevent client-side confusion while maintaining Gosper's fluid math, the library adopts a **Two-Tier Architecture**:
1. **Internal Speculator Tier:** Kernels (BLFT, Quadratic) may emit "noisy" generalized terms (zeroes, negatives, or speculative guesses) to maintain pipeline flow.
2. **Rectification Tier:** The final node before the Public API is a strict Homographic LFT (the "Rectifier"). It mathematically absorbs speculative terms into its coefficients (acting as a state-driven buffer) and only emits a term to the client when the integer floor is proven stable across the entire interval.

### 2.4. Productivity & Rational Collapse: Modern Repair
Irrational inputs may cause the engine to "hang" when the result is an exact rational, as bounds infinitely approach but never cross an integer floor.
* **Resolution:** The engine implements **Rational Collapse**. If interval bounds isolate a single rational value, the engine invokes the HAKMEM 101C utility to identify the simplest rational, emits the finite sequence, and terminates.

## 3. Mathematical Kernels

### 3.1. Unified Bi-Linear Fractional Transform (BLFT)
State is the coefficient tuple $(a,b,c,d,e,f,g,h)$ representing:
$$z(x, y) = \frac{axy + bx + cy + d}{exy + fx + gy + h}$$

* **Ingest Left/Right:** Substitute $x = p + q/x'$ into the state.
* **Produce:** Substitute $z = t + 1/z'$ into the state.
* **Normalization:** Apply exact GCD reduction to the 8-tuple when coefficients exceed a defined bit-length.

### 3.2. Diagonal Quadratic Kernel
State is $(a,b,c,d,e,f)$ representing:
$$z(x) = \frac{ax^2 + bx + c}{dx^2 + ex + f}$$

* **Operational Intent:** Primary engine for `square` and `sqrt` (Newton-feedback).
* **Interval Rule:** This kernel is **exempt** from monotonicity requirements. It must feed its speculative $(p,q)$ output into a downstream Rectifier to ensure client-visible correctness.

## 4. The Rectifier (Final Emission Rule)

The Rectifier is a Homographic LFT $[a, b; c, d]$ that sits at the end of every operator chain.

1. **Absorb:** Ingest speculative $(p, q)$ terms from the internal tier. Update coefficients:
   * $a' = ap + cq$
   * $b' = a$
   * $c' = cp + dq$
   * $d' = c$
2. **Interval Check:** Evaluate the transform over the range $[1, \infty]$. 
   * The bounds are $z(1) = (a+b)/(c+d)$ and $z(\infty) = a/c$.
3. **Emit:** A client-visible RCF term $t$ is emitted **only if** $\lfloor z(1) \rfloor = \lfloor z(\infty) \rfloor$.
4. **Iterate:** If a term is emitted, apply the production substitution. If not, request the next speculative term from the internal tier.

## 5. HAKMEM 101C Utility

The library must provide a standalone utility to find the "Simplest Rational in an Interval":
1. Convert interval endpoints to continued fractions.
2. Find the first term where they differ.
3. Increment the smaller term and discard subsequent terms.
4. This utility is used for **Rational Collapse** and final simplification of finite results.

## 6. Conformance Expectations

The implementation is conformant if:
* It lazily evaluates the target: `sqrt(3/pi^2 + e) / (tanh(sqrt(5)) - sin(69°))`.
* It never "recants" or "backspaces" a term once emitted to the Public API.
* It maintains exactness via BigInt coefficients.
* The `sqrt` implementation uses a feedback loop into a Rectifier rather than an external floating-point approximation.