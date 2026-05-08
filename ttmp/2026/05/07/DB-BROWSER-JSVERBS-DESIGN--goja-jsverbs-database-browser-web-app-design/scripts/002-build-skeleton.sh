#!/usr/bin/env bash
set -euo pipefail

gofmt -w cmd/db-browser/main.go
go test ./...
go run ./cmd/db-browser inspect modules
