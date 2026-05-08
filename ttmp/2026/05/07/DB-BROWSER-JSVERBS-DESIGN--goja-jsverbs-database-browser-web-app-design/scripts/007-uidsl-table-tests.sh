#!/usr/bin/env bash
set -euo pipefail

gofmt -w internal/uidsl/*.go internal/verbcli/*.go internal/verbrepos/*.go
go test ./internal/uidsl ./internal/verbcli ./internal/verbrepos
OUT="$(go run ./cmd/db-browser verbs examples builtin render-sample-table)"
printf '%s\n' "$OUT"
grep -q '<table class="ui-table ui-table--pagination ui-table--sorting" id="sample">' <<<"$OUT"
grep -q '<td data-column="name">Alice</td>' <<<"$OUT"
