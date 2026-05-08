#!/usr/bin/env bash
set -euo pipefail

gofmt -w internal/uidsl/*.go internal/app/*.go cmd/db-browser/main.go
go test ./internal/uidsl ./internal/app ./cmd/db-browser

DB_PATH="$(mktemp /tmp/db-browser-rich-table-XXXX.sqlite)"
BIN="$(mktemp /tmp/db-browser-bin-XXXXXX)"
PORT=18987
cleanup() {
  if [[ -n "${PID:-}" ]]; then kill "$PID" >/dev/null 2>&1 || true; fi
  rm -f "$DB_PATH" "$BIN"
}
trap cleanup EXIT
DB="$DB_PATH" python3 - <<'PY'
import os, sqlite3
path = os.environ["DB"]
con = sqlite3.connect(path)
con.execute("create table people(id integer primary key, name text)")
con.commit()
con.close()
PY
go build -o "$BIN" ./cmd/db-browser
"$BIN" serve --db "$DB_PATH" --scripts-dir examples/generic-browser/scripts --addr ":$PORT" --dev > /tmp/db-browser-rich-table.log 2>&1 &
PID=$!
for _ in $(seq 1 30); do
  if curl -fsS "http://127.0.0.1:$PORT/" >/tmp/db-browser-rich-table.html 2>/dev/null; then
    break
  fi
  sleep 0.2
done
grep -q 'Generic SQLite Browser' /tmp/db-browser-rich-table.html
grep -q '<table class="ui-table ui-table--pagination ui-table--sorting ui-table--column-picker" id="tables">' /tmp/db-browser-rich-table.html
grep -q 'people' /tmp/db-browser-rich-table.html
