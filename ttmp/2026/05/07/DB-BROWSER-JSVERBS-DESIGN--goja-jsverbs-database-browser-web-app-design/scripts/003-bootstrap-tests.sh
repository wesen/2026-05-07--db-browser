#!/usr/bin/env bash
set -euo pipefail

gofmt -w internal/verbrepos/bootstrap.go internal/verbrepos/bootstrap_test.go
go test ./internal/verbrepos ./cmd/db-browser
