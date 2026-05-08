#!/usr/bin/env bash
set -euo pipefail

go test ./...
ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/scripts/005-runtime-smoke.sh
ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/scripts/009-rich-table-tests.sh
ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/scripts/010-examples-smoke.sh
docmgr doctor --ticket DB-BROWSER-JSVERBS-DESIGN --stale-after 30
