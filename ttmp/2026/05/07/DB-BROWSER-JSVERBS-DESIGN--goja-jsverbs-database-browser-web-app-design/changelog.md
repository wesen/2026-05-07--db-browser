# Changelog

## 2026-05-07

- Initial workspace created.
- Added the primary design/implementation guide for the Goja jsverbs SQLite browser web app.
- Added the investigation diary with commands, evidence, and continuation notes.
- Updated ticket index and task checklist.

## 2026-05-07

Completed initial research/design package and diary for Goja jsverbs database browser.

### Related Files

- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/design-doc/01-goja-jsverbs-database-browser-design-and-implementation-guide.md — Primary design deliverable
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/reference/01-investigation-diary.md — Investigation diary


## 2026-05-07

Recorded validation and reMarkable upload evidence in the diary.

### Related Files

- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/reference/01-investigation-diary.md — Diary updated with validation and upload evidence


## 2026-05-07

Validated design against css-visual-diff and added yaml/lazy-verb bootstrap recommendations.

### Related Files

- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/design-doc/01-goja-jsverbs-database-browser-design-and-implementation-guide.md — Design guide updated
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/reference/01-investigation-diary.md — Diary updated with css-visual-diff validation


## 2026-05-07

Expanded implementation tasks and added ticket doctor script.

### Related Files

- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/reference/01-investigation-diary.md — Diary step 6
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/scripts/001-ticket-doctor.sh — Ticket validation script
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/tasks.md — Detailed implementation checklist


## 2026-05-07

Implemented initial Go project skeleton and build smoke script.

### Related Files

- /home/manuel/code/wesen/2026-05-07--db-browser/cmd/db-browser/main.go — Root CLI skeleton
- /home/manuel/code/wesen/2026-05-07--db-browser/examples/builtin-verbs/hello.js — Built-in smoke verb fixture
- /home/manuel/code/wesen/2026-05-07--db-browser/go.mod — Go module definition
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/reference/01-investigation-diary.md — Diary step 7
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/scripts/002-build-skeleton.sh — Skeleton validation script


## 2026-05-07

Implemented verb repository bootstrap with embedded/config/env/CLI sources and tests.

### Related Files

- /home/manuel/code/wesen/2026-05-07--db-browser/internal/verbrepos/bootstrap.go — Repository bootstrap implementation
- /home/manuel/code/wesen/2026-05-07--db-browser/internal/verbrepos/bootstrap_test.go — Bootstrap unit tests
- /home/manuel/code/wesen/2026-05-07--db-browser/internal/verbrepos/builtin/hello.js — Embedded built-in verb fixture
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/reference/01-investigation-diary.md — Diary step 8
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/scripts/003-bootstrap-tests.sh — Bootstrap validation script


## 2026-05-07

Mounted scanned jsverbs as a lazy dynamic CLI with list and duplicate tests.

### Related Files

- /home/manuel/code/wesen/2026-05-07--db-browser/cmd/db-browser/main.go — Root command now uses verbcli
- /home/manuel/code/wesen/2026-05-07--db-browser/internal/verbcli/command.go — Lazy dynamic jsverbs command implementation
- /home/manuel/code/wesen/2026-05-07--db-browser/internal/verbcli/command_test.go — jsverbs CLI tests
- /home/manuel/code/wesen/2026-05-07--db-browser/internal/verbcli/list.go — verbs list implementation
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/reference/01-investigation-diary.md — Diary step 9
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/scripts/004-jsverbs-cli-tests.sh — jsverbs CLI validation script


## 2026-05-07

Added CLI verb runtime profile with yaml and configured SQLite db aliases.

### Related Files

- /home/manuel/code/wesen/2026-05-07--db-browser/internal/verbcli/command.go — Runtime flags and invoker wiring
- /home/manuel/code/wesen/2026-05-07--db-browser/internal/verbcli/runtime.go — Goja runtime factory and database/yaml module profile
- /home/manuel/code/wesen/2026-05-07--db-browser/internal/verbrepos/builtin/db.js — SQLite runtime smoke verb
- /home/manuel/code/wesen/2026-05-07--db-browser/internal/verbrepos/builtin/yaml.js — YAML runtime smoke verb
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/reference/01-investigation-diary.md — Diary step 10
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/scripts/005-runtime-smoke.sh — Runtime validation script


## 2026-05-07

Recorded commit checkpoints for completed T01-T05 implementation tasks.

### Related Files

- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/reference/01-investigation-diary.md — Commit checkpoint diary step


## 2026-05-07

Copied goja-hosting-site web/uidsl packages and switched CLI runtime to go-go-goja database module.

### Related Files

- /home/manuel/code/wesen/2026-05-07--db-browser/internal/uidsl — Copied and adapted UI renderer package
- /home/manuel/code/wesen/2026-05-07--db-browser/internal/verbcli/runtime.go — Now wires go-go-goja database module
- /home/manuel/code/wesen/2026-05-07--db-browser/internal/web — Copied and adapted Express/web host package
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/reference/01-investigation-diary.md — Diary step 12
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/scripts/006-web-uidsl-copy-tests.sh — Web/uidsl validation script


## 2026-05-07

Added ui.table.fromRows and CLI ui.dsl runtime smoke coverage.

### Related Files

- /home/manuel/code/wesen/2026-05-07--db-browser/internal/uidsl/table.go — ui.table.fromRows implementation
- /home/manuel/code/wesen/2026-05-07--db-browser/internal/uidsl/table_test.go — UI table renderer tests
- /home/manuel/code/wesen/2026-05-07--db-browser/internal/verbcli/runtime.go — Registers ui.dsl in CLI verb runtime
- /home/manuel/code/wesen/2026-05-07--db-browser/internal/verbrepos/builtin/ui.js — CLI smoke verb for ui.dsl table rendering
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/reference/01-investigation-diary.md — Diary step 13
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/scripts/007-uidsl-table-tests.sh — UI DSL validation script


## 2026-05-07

Wired copied Express host into db-browser serve command with SQLite and ui.dsl modules.

### Related Files

- /home/manuel/code/wesen/2026-05-07--db-browser/cmd/db-browser/main.go — serve command now runs app server
- /home/manuel/code/wesen/2026-05-07--db-browser/internal/app/server.go — Goja web server runtime and script loader
- /home/manuel/code/wesen/2026-05-07--db-browser/internal/app/server_test.go — Serve integration tests
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/reference/01-investigation-diary.md — Diary step 14
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/scripts/008-serve-tests.sh — Serve validation script


## 2026-05-07

Expanded ui.table into rich DSL v1 and added generic browser smoke example.

### Related Files

- /home/manuel/code/wesen/2026-05-07--db-browser/examples/generic-browser/README.md — Example run instructions
- /home/manuel/code/wesen/2026-05-07--db-browser/examples/generic-browser/scripts/app.js — Serve-mode generic SQLite browser example
- /home/manuel/code/wesen/2026-05-07--db-browser/internal/uidsl/table.go — Rich table DSL v1 implementation
- /home/manuel/code/wesen/2026-05-07--db-browser/internal/uidsl/table_rich_test.go — Rich table behavior tests
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/reference/01-investigation-diary.md — Diary step 15
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/scripts/009-rich-table-tests.sh — Rich table smoke validation script


## 2026-05-07

Added generic browser and YAML dashboard examples with smoke validation.

### Related Files

- /home/manuel/code/wesen/2026-05-07--db-browser/examples/generic-browser/scripts/app.js — Generic browser example
- /home/manuel/code/wesen/2026-05-07--db-browser/examples/yaml-dashboard/dashboard.yaml — YAML dashboard spec
- /home/manuel/code/wesen/2026-05-07--db-browser/examples/yaml-dashboard/scripts/app.js — YAML dashboard serve example
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/reference/01-investigation-diary.md — Diary step 16
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/scripts/010-examples-smoke.sh — Examples smoke validation script


## 2026-05-07

Refreshed design documentation with implementation status and added final validation script.

### Related Files

- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/design-doc/01-goja-jsverbs-database-browser-design-and-implementation-guide.md — Implementation status update
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/reference/01-investigation-diary.md — Diary step 17
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/scripts/011-final-validation.sh — Final validation script


## 2026-05-07

Added seeded Playwright smoke DB/app and validated it in browser.

### Related Files

- /home/manuel/code/wesen/2026-05-07--db-browser/examples/playwright-smoke/README.md — Run instructions for Playwright smoke app
- /home/manuel/code/wesen/2026-05-07--db-browser/examples/playwright-smoke/data/app.db — Seed SQLite fixture for browser smoke app
- /home/manuel/code/wesen/2026-05-07--db-browser/examples/playwright-smoke/scripts/app.js — Playwright smoke Express/ui.dsl app
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/reference/01-investigation-diary.md — Diary step 18
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/scripts/012-playwright-smoke.sh — Server smoke script for Playwright app
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-JSVERBS-DESIGN--goja-jsverbs-database-browser-web-app-design/scripts/013-playwright-checklist.md — Manual Playwright validation checklist

