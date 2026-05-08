#!/usr/bin/env bash
set -euo pipefail

go test ./internal/uidsl -run 'Test.*(CodeBlock|Badge|Tabs)' -count=1
go test ./...

BIN="$(mktemp /tmp/db-browser-uidsl-components-bin-XXXXXX)"
DB="$(mktemp /tmp/db-browser-uidsl-components-db-XXXXXX.sqlite)"
PORT="19101"
cleanup() {
  if [[ -n "${PID:-}" ]]; then kill "$PID" >/dev/null 2>&1 || true; fi
  rm -f "$BIN" "$DB"
}
trap cleanup EXIT

go build -o "$BIN" ./cmd/db-browser
DB="$DB" python3 - <<'PY'
import os, sqlite3
con = sqlite3.connect(os.environ['DB'])
con.execute('create table people(id integer primary key, name text, status text)')
con.executemany('insert into people(name,status) values (?,?)', [('Alice','active'), ('Bob','pending')])
con.commit(); con.close()
PY
"$BIN" serve --db "$DB" --scripts-dir examples/generic-browser/scripts --addr ":$PORT" --dev > /tmp/db-browser-uidsl-components.log 2>&1 &
PID=$!
for _ in $(seq 1 40); do
  if curl -fsS "http://127.0.0.1:$PORT/tables/people" > /tmp/db-browser-uidsl-components.html 2>/dev/null; then
    break
  fi
  sleep 0.2
done

grep -q 'ui-tabs' /tmp/db-browser-uidsl-components.html
grep -q 'ui-codeblock' /tmp/db-browser-uidsl-components.html
grep -q 'ui-badge' /tmp/db-browser-uidsl-components.html
grep -q 'Debug JSON' /tmp/db-browser-uidsl-components.html
grep -q 'CREATE TABLE people' /tmp/db-browser-uidsl-components.html
