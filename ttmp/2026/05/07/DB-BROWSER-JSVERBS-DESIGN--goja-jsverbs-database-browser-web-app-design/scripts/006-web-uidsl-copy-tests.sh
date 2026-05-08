#!/usr/bin/env bash
set -euo pipefail

gofmt -w internal/web/*.go internal/uidsl/*.go internal/verbcli/*.go
go test ./internal/web ./internal/uidsl ./internal/verbcli
