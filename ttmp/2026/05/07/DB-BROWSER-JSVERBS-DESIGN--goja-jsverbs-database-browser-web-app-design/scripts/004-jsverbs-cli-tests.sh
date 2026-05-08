#!/usr/bin/env bash
set -euo pipefail

gofmt -w cmd/db-browser/main.go internal/verbcli/*.go internal/verbrepos/*.go
go test ./internal/verbrepos ./internal/verbcli ./cmd/db-browser
OUT="$(go run ./cmd/db-browser verbs list)"
printf '%s\n' "$OUT"
grep -q 'examples builtin hello' <<<"$OUT"
