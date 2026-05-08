#!/usr/bin/env bash
set -euo pipefail

go test ./...

run_server_check() {
  local db_path="$1"
  local scripts_dir="$2"
  local port="$3"
  local expect="$4"
  local bin="$5"
  "$bin" serve --db "$db_path" --scripts-dir "$scripts_dir" --addr ":$port" --dev > "/tmp/db-browser-example-$port.log" 2>&1 &
  local pid=$!
  for _ in $(seq 1 40); do
    if curl -fsS "http://127.0.0.1:$port/" > "/tmp/db-browser-example-$port.html" 2>/dev/null; then
      break
    fi
    sleep 0.2
  done
  kill "$pid" >/dev/null 2>&1 || true
  grep -q "$expect" "/tmp/db-browser-example-$port.html"
}

BIN="$(mktemp /tmp/db-browser-examples-bin-XXXXXX)"
DB1="$(mktemp /tmp/db-browser-example-generic-XXXX.sqlite)"
DB2="$(mktemp /tmp/db-browser-example-yaml-XXXX.sqlite)"
trap 'rm -f "$BIN" "$DB1" "$DB2"' EXIT

go build -o "$BIN" ./cmd/db-browser
DB="$DB1" python3 - <<'PY'
import os, sqlite3
con = sqlite3.connect(os.environ['DB'])
con.execute('create table people(id integer primary key, name text)')
con.commit(); con.close()
PY
run_server_check "$DB1" examples/generic-browser/scripts 18988 people "$BIN"

DB="$DB2" python3 - <<'PY'
import os, sqlite3
con = sqlite3.connect(os.environ['DB'])
con.execute('create table people(id integer primary key, name text)')
con.executemany('insert into people(name) values (?)', [('Alice',), ('Bob',)])
con.commit(); con.close()
PY
run_server_check "$DB2" examples/yaml-dashboard/scripts 18989 'YAML Dashboard' "$BIN"
grep -q 'People' /tmp/db-browser-example-18989.html
grep -q '>2<' /tmp/db-browser-example-18989.html
