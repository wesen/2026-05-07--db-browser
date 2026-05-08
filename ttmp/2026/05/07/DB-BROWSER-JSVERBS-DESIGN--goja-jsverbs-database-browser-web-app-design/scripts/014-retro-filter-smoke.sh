#!/usr/bin/env bash
set -euo pipefail

BIN="$(mktemp /tmp/db-browser-retro-bin-XXXXXX)"
PORT="19092"
cleanup() {
  if [[ -n "${PID:-}" ]]; then kill "$PID" >/dev/null 2>&1 || true; fi
  rm -f "$BIN"
}
trap cleanup EXIT

go build -o "$BIN" ./cmd/db-browser
"$BIN" serve --db examples/playwright-smoke/data/app.db --scripts-dir examples/playwright-smoke/scripts --addr ":$PORT" --dev > /tmp/db-browser-retro-filter.log 2>&1 &
PID=$!
for _ in $(seq 1 40); do
  if curl -fsS "http://127.0.0.1:$PORT/?filter.segment=vip" > /tmp/db-browser-retro-filter-vip.html 2>/dev/null; then
    break
  fi
  sleep 0.2
done

grep -q 'macos-desktop' /tmp/db-browser-retro-filter-vip.html
grep -q 'Alice Example' /tmp/db-browser-retro-filter-vip.html
grep -q '/customers/1' /tmp/db-browser-retro-filter-vip.html
! grep -q 'Bob Browser' /tmp/db-browser-retro-filter-vip.html
! grep -q 'Carla Canvas' /tmp/db-browser-retro-filter-vip.html
grep -q 'Page 1 of 1 (1 rows)' /tmp/db-browser-retro-filter-vip.html
curl -fsS "http://127.0.0.1:$PORT/?q=bob&filter.segment=vip" > /tmp/db-browser-retro-filter-empty.html
grep -q 'No rows match the current filters' /tmp/db-browser-retro-filter-empty.html
grep -q 'Page 1 of 1 (0 rows)' /tmp/db-browser-retro-filter-empty.html
! grep -q '&lt;nil&gt;' /tmp/db-browser-retro-filter-empty.html
