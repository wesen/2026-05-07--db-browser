#!/usr/bin/env bash
set -euo pipefail

DB_PATH="examples/playwright-smoke/data/app.db"
BIN="$(mktemp /tmp/db-browser-playwright-bin-XXXXXX)"
PORT="19090"
cleanup() {
  if [[ -n "${PID:-}" ]]; then kill "$PID" >/dev/null 2>&1 || true; fi
  rm -f "$BIN"
}
trap cleanup EXIT

go build -o "$BIN" ./cmd/db-browser
"$BIN" serve --db "$DB_PATH" --scripts-dir examples/playwright-smoke/scripts --addr ":$PORT" --dev > /tmp/db-browser-playwright-smoke.log 2>&1 &
PID=$!
for _ in $(seq 1 40); do
  if curl -fsS "http://127.0.0.1:$PORT/" >/tmp/db-browser-playwright-smoke.html 2>/dev/null; then
    exit 0
  fi
  sleep 0.2
done
cat /tmp/db-browser-playwright-smoke.log >&2 || true
exit 1
