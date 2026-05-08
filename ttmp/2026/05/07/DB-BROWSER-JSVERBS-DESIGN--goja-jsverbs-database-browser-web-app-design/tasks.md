# Tasks

## Implementation sequence

### T01 — Ticket/task preparation

- [x] Expand this task list into implementation-sized tasks.
- [x] Add a ticket-local smoke-test script under `scripts/`.
- [x] Commit the updated ticket planning docs.

### T02 — Go project skeleton

- [x] Create `go.mod` for the DB browser app.
- [x] Add `cmd/db-browser/main.go` with root command, logging defaults, and a lazy `verbs` command placeholder.
- [x] Add a tiny embedded example verb repository for smoke tests.
- [x] Add `.gitignore` entries for build/test artifacts.
- [x] Run formatting and basic compile checks.
- [x] Commit the skeleton.

### T03 — Verb repository bootstrap

- [x] Implement repository source model: embedded, config, env, and CLI repositories.
- [x] Support `DB_BROWSER_VERB_REPOSITORIES` with `filepath.SplitList`.
- [x] Support leading `--repository` and `--verb-repository` flags before dynamic verb paths.
- [x] Support `.db-browser.yml` and `.db-browser.override.yml` with `verbs.repositories[]` entries.
- [x] Normalize `~`, relative paths, absolute paths, and duplicate repository identities.
- [x] Add unit tests for normalization and bootstrap ordering.
- [x] Commit repository bootstrap.

### T04 — jsverbs scanning and dynamic CLI mounting

- [x] Scan each repository with `jsverbs.ScanFS` / `jsverbs.ScanDir`.
- [x] Default to `IncludePublicFunctions=false` for explicit `__verb__`-only command export.
- [x] Register DB-browser shared sections if needed. (No shared sections needed yet.)
- [x] Detect duplicate `verb.FullPath()` values with provenance.
- [x] Build a lazy `verbs` command that rebuilds the real command tree after bootstrap flag parsing.
- [x] Add a `verbs list` command that prints discovered verbs.
- [x] Add tests for duplicate detection and list output.
- [x] Commit jsverbs dynamic CLI mounting.

### T05 — Runtime module profile for CLI verbs

- [x] Build a runtime factory for a scanned repository using `registry.RequireLoader()`.
- [x] Add repository root, repository `node_modules`, parent, and parent `node_modules` require roots.
- [x] Enable `fs`, `path`, `time`, `timer`, and `yaml` through default go-go-goja modules.
- [x] Add `database` and `db` aliases when `--db` is provided.
- [x] Add `--db`, `--readonly`, and `--allow-writes` runtime flags.
- [x] Execute a fixture verb that imports `yaml` and optionally `database`.
- [x] Commit runtime module profile.

### T06 — Minimal UI DSL package

- [x] Add `internal/uidsl` node/render primitives or import/refactor the goja-hosting-site version.
- [x] Expose `ui.dsl` / `ui` modules with `page`, tags, fragments, text, raw, and render.
- [x] Add `ui.table.fromRows` as the first high-level primitive.
- [x] Add renderer tests for escaping and documents via copied goja-hosting-site tests. (Basic table output remains with `ui.table.fromRows`.)
- [x] Commit minimal UI DSL.

### T07 — Minimal Express-style web host

- [x] Add `internal/web` host, route registry, request/response DTOs, and `express` registrar.
- [x] Add `serve` command with `--addr`, `--db`, `--scripts-dir`, and `--dev`.
- [x] Load scripts in deterministic order and let them register routes.
- [x] Wire `res.html` through the UI renderer.
- [x] Add `httptest` integration test for HTML, JSON, redirect, and POST body handling.
- [x] Commit web host.

### T08 — Rich table DSL v1

- [x] Implement `ui.table(id).data(...).columns(...).features(...).render(...)`.
- [x] Implement table context parsing for query, params, pagination, and sorting. (Filter builder remains future work.)
- [x] Implement column types: text, badge, money, date, tags.
- [x] Implement pagination, sorting headers, and column picker markers. (Interactive column picker UI remains future work.)
- [x] Add tests and an example SQLite browser script.
- [x] Commit rich table DSL v1.

### T09 — Examples and validation scripts

- [x] Add `examples/generic-browser` with a small SQLite fixture or fixture setup script.
- [x] Add `examples/yaml-dashboard` to exercise `require("yaml")`.
- [x] Add ticket-local scripts for build/test/smoke validation.
- [x] Commit examples and validation scripts.

### T10 — Documentation refresh and reMarkable update

- [x] Update design doc with implementation notes and any deviations.
- [x] Update diary with chronological commits and validation evidence.
- [x] Run `docmgr doctor --ticket DB-BROWSER-JSVERBS-DESIGN --stale-after 30`.
- [x] Upload refreshed bundle to reMarkable if requested.
- [x] Commit final documentation refresh.

## Completed research/design tasks

- [x] Create docmgr ticket workspace.
- [x] Add primary design document.
- [x] Add investigation diary.
- [x] Inspect go-go-goja runtime, modules, and jsverbs implementation.
- [x] Inspect go-minitrace repository-backed command scanning pattern.
- [x] Inspect goja-hosting-site Express and UI DSL prototype.
- [x] Validate design against css-visual-diff JavaScript playground patterns.
- [x] Add yaml to the proposed JavaScript host surface.
- [x] Write intern-facing design and implementation guide.
