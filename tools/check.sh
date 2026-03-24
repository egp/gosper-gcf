#!/usr/bin/env zsh
# tools/check.zsh v2

set -euo pipefail

# Run from anywhere inside the repo.
REPO_ROOT="$(cd "$(dirname "${0:A}")/.." && pwd)"
cd "$REPO_ROOT"

print -P "%F{cyan}==> gofmt (and goimports if available)%f"

# gofmt is always available
gofmt -w .

# goimports is optional; nice to have
if command -v goimports >/dev/null 2>&1; then
  goimports -w .
else
  print -P "%F{yellow}note:%f goimports not found; skipping import normalization"
fi

print -P "%F{cyan}==> go vet%f"
go vet ./core
go vet ./named

print -P "%F{cyan}==> staticcheck%f"
staticcheck ./core
staticcheck ./named


print -P "%F{cyan}==> tests (no cache)%f"
go test -count=9 ./core
go test -count=9 ./named


print -P "%F{cyan}==> coverage summary%f"
mkdir -p ./tmp
go test -count=9 -coverprofile=./tmp/cover_cf.out ./core
go test -count=9 -coverprofile=./tmp/cover_cfsource.out ./named

print -P "%F{green}OK%f"
# tools/check.zsh v2