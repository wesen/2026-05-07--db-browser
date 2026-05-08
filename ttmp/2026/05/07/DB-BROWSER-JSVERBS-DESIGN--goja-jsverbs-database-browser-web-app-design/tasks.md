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

- [ ] Scan each repository with `jsverbs.ScanFS` / `jsverbs.ScanDir`.
- [ ] Default to `IncludePublicFunctions=false` for explicit `__verb__`-only command export.
- [ ] Register DB-browser shared sections if needed.
- [ ] Detect duplicate `verb.FullPath()` values with provenance.
- [ ] Build a lazy `verbs` command that rebuilds the real command tree after bootstrap flag parsing.
- [ ] Add a `verbs list` command that prints discovered verbs.
- [ ] Add tests for duplicate detection and list output.
- [ ] Commit jsverbs dynamic CLI mounting.

### T05 — Runtime module profile for CLI verbs

- [ ] Build a runtime factory for a scanned repository using `registry.RequireLoader()`.
- [ ] Add repository root, repository `node_modules`, parent, and parent `node_modules` require roots.
- [ ] Enable `fs`, `path`, `time`, `timer`, and `yaml` through default go-go-goja modules.
- [ ] Add `database` and `db` aliases when `--db` is provided.
- [ ] Add `--db`, `--readonly`, and `--allow-writes` runtime flags.
- [ ] Execute a fixture verb that imports `yaml` and optionally `database`.
- [ ] Commit runtime module profile.

### T06 — Minimal UI DSL package

- [ ] Add `internal/uidsl` node/render primitives or import/refactor the goja-hosting-site version.
- [ ] Expose `ui.dsl` / `ui` modules with `page`, tags, fragments, text, raw, and render.
- [ ] Add `ui.table.fromRows` as the first high-level primitive.
- [ ] Add renderer tests for escaping, documents, and basic table output.
- [ ] Commit minimal UI DSL.

### T07 — Minimal Express-style web host

- [ ] Add `internal/web` host, route registry, request/response DTOs, and `express` registrar.
- [ ] Add `serve` command with `--addr`, `--db`, `--scripts-dir`, and `--dev`.
- [ ] Load scripts in deterministic order and let them register routes.
- [ ] Wire `res.html` through the UI renderer.
- [ ] Add `httptest` integration test for HTML, JSON, redirect, and POST body handling.
- [ ] Commit web host.

### T08 — Rich table DSL v1

- [ ] Implement `ui.table(id).data(...).columns(...).features(...).render(...)`.
- [ ] Implement table context parsing for query, params, pagination, sorting, and filters.
- [ ] Implement column types: text, badge, money, date, tags.
- [ ] Implement pagination, sorting headers, empty states, and column picker.
- [ ] Add tests and an example SQLite browser script.
- [ ] Commit rich table DSL v1.

### T09 — Examples and validation scripts

- [ ] Add `examples/generic-browser` with a small SQLite fixture or fixture setup script.
- [ ] Add `examples/yaml-dashboard` to exercise `require("yaml")`.
- [ ] Add ticket-local scripts for build/test/smoke validation.
- [ ] Commit examples and validation scripts.

### T10 — Documentation refresh and reMarkable update

- [ ] Update design doc with implementation notes and any deviations.
- [ ] Update diary with chronological commits and validation evidence.
- [ ] Run `docmgr doctor --ticket DB-BROWSER-JSVERBS-DESIGN --stale-after 30`.
- [ ] Upload refreshed bundle to reMarkable if requested.
- [ ] Commit final documentation refresh.

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
