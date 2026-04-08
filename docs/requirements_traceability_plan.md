# Requirements Traceability Plan

_Status: **deferred — implement after Stage 5 / MVP complete**_  
_Written: 2026-04-08_

---

## Goal

Bidirectional traceability between spec requirements and test cases:

- Every leaf requirement has a unique numeric ID and a list of tests that cover it.
- Every test is annotated with the requirement IDs it exercises.
- A script can run all tests for a given requirement ID.
- A coverage report shows which requirements have no tests (nil leaves).

---

## Numeric ID Scheme

IDs are hierarchical dotted integers derived from the spec section structure,
up to four levels deep: `A.B.C.D`

| Level | Meaning | Example |
|---|---|---|
| A | Spec section (3 = Requirements) | `3` |
| B | Sub-section (1 = Must Have, 2 = Should Have) | `3.1` |
| C | Requirement clause within sub-section | `3.1.4` |
| D | Leaf sub-requirement within a clause | `3.1.4.2` |

Example IDs:

| ID | Requirement leaf |
|---|---|
| `3.1.1` | Accept generalized (p,q) input streams |
| `3.1.4.1` | BLFT: ingest x-stream term (p,q) |
| `3.1.4.2` | BLFT: ingest y-stream term (r,s) |
| `3.1.4.3` | BLFT: produce output term t |
| `3.1.4.4` | BLFT: normalize by GCD |
| `3.1.4.5` | BLFT: collapse x-EOF (x → ∞) |
| `3.1.4.6` | BLFT: collapse y-EOF (y → ∞) |
| `3.1.8` | Expose CurrentInterval() alongside NextRCF() |
| `3.1.9.1` | Inside interval semantics |
| `3.1.9.2` | Outside interval semantics |
| `3.1.9.3` | Denominator-crossing detection → outside |
| `3.1.10` | Widest-range ingest-ordering rule |
| `3.1.13` | SinDegrees exact shortcut for multiples of 30° |
| `3.1.15` | HAKMEM 101C SmallestRationalInInterval |
| `3.1.16` | Resource guards: time and bit-length limits |
| `3.2.1` | Decimal digit emission (multiply-by-10 rule) |

The full ID table is generated as part of implementation (Action 1 below).

---

## Implementation Plan

### Action 1 — Enumerate leaf requirements with IDs (~4–6 h)

Walk `docs/gosper_cf_requirements_spec.md` sections 3.1 (Must Have) and 3.2 (Should Have).
For each bullet, assign an ID and decompose into testable leaves where needed.
Write the result to `docs/requirements_ids.md` — one table row per leaf:

```
| ID | Fidelity | Leaf description | Coverage |
```

`Coverage` values: `full` | `partial` | `nil`. Start everything at `nil`.

### Action 2 — Annotate existing tests with `// req:` comments (~2–3 h)

Add a single comment line immediately before each test function:

```go
// req: 3.1.4.1, 3.1.4.4
func TestBB_GCF_Add(t *testing.T) {
```

No test renaming. The comment is machine-parseable by the tooling in Action 4.
Multiple IDs are comma-separated. A test with no clear requirement leaf gets `// req: TBD`.

### Action 3 — Cross-reference YAML (~2 h)

Create `docs/requirements_trace.yaml` keyed by requirement ID:

```yaml
"3.1.4.1":
  description: "BLFT: ingest x-stream term (p,q)"
  spec_ref: "3.1 / 4.2"
  coverage: nil
  tests: []

"3.1.4.4":
  description: "BLFT: normalize by GCD"
  spec_ref: "3.1 / 4.2"
  coverage: partial
  tests:
    - TestBB_GCF_Add
    - TestWB_BLFT_NormalizeGCD
```

The YAML is the source of truth for the req→test direction.
The `// req:` comments are authoritative for the test→req direction.
A validation script (Action 4) checks that the two are consistent.

### Action 4 — Tooling scripts (~2 h)

`tools/run_req.sh <req-id>`
: Looks up `req-id` in `requirements_trace.yaml`, extracts test names,
  builds a `-run` regex, and invokes `go test` on all three packages.

`tools/check_coverage.sh`
: Parses `requirements_trace.yaml`, reports all IDs with `coverage: nil`.
  Also parses `// req:` comments across all `*_test.go` files and flags
  any test whose IDs don't appear in the YAML (orphan annotations).

`tools/validate_trace.sh`
: Cross-checks YAML tests vs `// req:` comments; exits non-zero on mismatch.
  Intended to run as part of `check.sh` once the traceability system is live.

### Action 5 — Update `requirements_ids.md` coverage column (~ongoing)

After Actions 2–3, manually review each leaf and set `coverage` to
`full`, `partial`, or `nil`. Add missing tests for any nil leaf that
falls within the current stage scope.

---

## Level of Effort

| Action | Estimate |
|---|---|
| 1 — Enumerate leaf IDs | 4–6 h |
| 2 — Annotate ~60 existing tests | 2–3 h |
| 3 — Write cross-reference YAML | 2 h |
| 4 — Write three scripts | 2 h |
| 5 — Set coverage column | 1 h |
| **Total upfront** | **~11–14 h** |
| Ongoing per new test | ~5 min |

---

## Why Deferred

The requirements hierarchy will shift through Stage 5: requirements are
still being refined (trig.Sin irrational fix, resource guards, 101C),
and tagging against a moving spec creates busywork. Traceability pays off
when the spec is stable and the test suite is comprehensive.

**Trigger condition to start:** Stage 5 commit is green on main.

---

## Interim Discipline (zero cost, starts now)

Add `// req: TBD` to every new test written from Stage 1.5 onward.
This positions us to populate Action 3 from real annotation data
rather than reconstructing it from memory after MVP.
