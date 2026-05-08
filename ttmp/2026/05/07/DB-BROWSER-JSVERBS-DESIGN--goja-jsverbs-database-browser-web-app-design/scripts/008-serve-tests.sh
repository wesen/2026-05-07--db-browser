#!/usr/bin/env bash
set -euo pipefail

gofmt -w cmd/db-browser/main.go internal/app/*.go internal/web/*.go internal/uidsl/*.go internal/verbcli/*.go
go test ./internal/app ./internal/web ./internal/uidsl ./internal/verbcli ./cmd/db-browser
