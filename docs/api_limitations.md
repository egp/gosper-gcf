# GoGCF Public API Limitations

_Living document — updated each stage. Last updated: 2026-04-07 (Stage 1 close)._

This file documents known limitations of the public API: behaviors that are broken,
unimplemented, or only partially working. All items are tracked as pending tests.

---

## `trig.Sin(x PQStream) *core.GCF`

**Status: broken for all irrational-result inputs.**

- Works correctly only when `sin(x)` is rational: currently `x = 0` (sin 0 = 0).
- For any input whose sine is irrational, the function produces wrong output terms.
  Root cause: the Lambert continued-fraction series does not handle large inputs (no
  range reduction mod 2π) and the binary BLFT range computation is missing the
  outside/outside case needed for the replay suffix ranges.
- Fix required before Stage 5 MVP assembly.

Pending tests:
- `TestBB_Trig_Sin_OneHalfMatchesKnownPrefix` — sin(1/2 radians)
- `TestBB_Trig_Sin_SixtyNineRadiansMatchesKnownPrefix` — sin(69 radians)

---

## `named.SinDegrees(x PQStream) *core.GCF`

**Status: partially working — rational-sine angles only.**

- Works correctly when `sin(x°)` is rational: multiples of 30° (0°, ±30°, ±90°, ±150°, 180°,
  ±210°, ±270°, ±330°) handled via exact shortcut.
- Broken for all other angles (e.g., 69°). Same root cause as `trig.Sin`.
- sin(69°) is required for the MVP expression; deferred until after other MVP stages are
  complete and the outside/outside BLFT range fix is implemented.

Pending tests:
- `TestBB_Named_SinDegrees_SixtyNineMatchesKnownPrefix` — sin(69°)

---

## `core.GCF` — unary BLFT with F=0 coefficient over outside-range input

**Status: returns error for this specific configuration.**

When a unary BLFT whose F coefficient is zero receives an input stream whose current
interval is an outside (projective-wrap-around) range, `NextRCF()` returns an error
rather than a term. The "observe/identity" BLFT pattern (B=1, H=1, all others zero)
has F=0 and is the most likely production use case.

Whether any current MVP computation path actually triggers this is unverified.

Pending test:
- `TestWB_GCF_UnaryIdentity_OverPQStreamFromRCF_OfDegreesScaleOutsideCase_DoesNotPanic`

---

## `named.Sqrt`, `named.Sqrt2`, `core` DLFT — square root operators

**Status: not yet fully validated (Stage 2 pending).**

`sqrt(2)` prefix test is parked. The DLFT handles finite algebraic sources; behavior
with infinite algebraic sources (e.g., `square(sqrt(x))`) is a known open problem.

Pending tests:
- `TestBB_GCF_SqrtOfTwoMatchesKnownPrefix`
- `TestBB_SquareOfSqrt2IsExactlyTwo`

---

## HAKMEM 101C — `SmallestRationalInInterval`

**Status: not yet implemented (Stage 1.5 pending).**

The utility that finds the rational with smallest numerator and denominator inside a
given interval is required by Rational Collapse and by Named composition tests.
No public function exists yet.

---

## Resource guards — timeout, bit-length limits

**Status: not yet implemented (Stage 4.5 pending).**

`GCF.NextRCF()` does not yet enforce caller-configured time or work limits. A stuck
computation will block indefinitely. Tests use goroutine timeouts as a workaround.

---

## Summary table

| API surface | Works for | Broken for | Stage fix |
|---|---|---|---|
| `trig.Sin` | x=0 (sin=0) | all irrational-result inputs | pre-Stage 5 |
| `named.SinDegrees` | multiples of 30° | all other angles | pre-Stage 5 |
| unary GCF, F=0 | inside-range inputs | outside-range inputs | TBD |
| sqrt/DLFT | finite algebraic sources | infinite algebraic sources | Stage 2 |
| 101C utility | — | not implemented | Stage 1.5 |
| resource guards | — | not implemented | Stage 4.5 |
