#!/usr/bin/env bash
set -euo pipefail

gofmt -w cmd/db-browser/main.go internal/verbcli/*.go internal/verbrepos/*.go
go test ./internal/verbrepos ./internal/verbcli ./cmd/db-browser

YAML_OUT="$(go run ./cmd/db-browser verbs examples builtin yaml-keys --text $'alpha: 1\nbeta: 2' --output json)"
printf '%s\n' "$YAML_OUT"
grep -q '"key": "alpha"' <<<"$YAML_OUT"
grep -q '"key": "beta"' <<<"$YAML_OUT"

DB_PATH="$(mktemp /tmp/db-browser-runtime-XXXX.sqlite)"
trap 'rm -f "$DB_PATH"' EXIT
DB="$DB_PATH" python3 - <<'PY'
import os, sqlite3
path = os.environ["DB"]
con = sqlite3.connect(path)
con.execute("create table users(id integer primary key, name text)")
con.commit()
con.close()
PY
DB_OUT="$(go run ./cmd/db-browser verbs --db "$DB_PATH" examples builtin tables --output json)"
printf '%s\n' "$DB_OUT"
grep -q '"name": "users"' <<<"$DB_OUT"
