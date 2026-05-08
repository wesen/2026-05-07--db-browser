# Tasks

## DB-BROWSER-UIDSL-COMPONENTS implementation sequence

### T01 — Ticket planning and implementation guide

- [x] Create docmgr ticket workspace.
- [x] Add design/implementation guide for `ui.codeBlock`, `ui.badge`, and `ui.tabs`.
- [x] Add detailed implementation tasks.
- [x] Add initial diary entry and changelog.
- [x] Validate ticket with `docmgr doctor`.
- [x] Commit planning docs.

### T02 — `ui.codeBlock` and convenience aliases

- [ ] Add `internal/uidsl/components.go` or equivalent focused component file.
- [ ] Implement `ui.codeBlock(language, source, options?)`.
- [ ] Implement `ui.sql(source, options?)`.
- [ ] Implement `ui.js(source, options?)`.
- [ ] Implement `ui.jsonBlock(value, options?)` with pretty-print behavior.
- [ ] Normalize language tokens and default invalid/empty language to `text`.
- [ ] Support `title`, `lineNumbers`, `wrap`, `copy`, `maxHeight`, and `class` options.
- [ ] Add escaping and render-contract tests.
- [ ] Commit codeBlock implementation.

### T03 — `ui.badge`

- [ ] Implement `ui.badge(value, options?)`.
- [ ] Support `tone`, `title`, and `class` options.
- [ ] Normalize tone and value CSS tokens.
- [ ] Ensure unknown tone falls back to `default`.
- [ ] Add escaping, tone, and class tests.
- [ ] Commit badge implementation.

### T04 — `ui.tabs`

- [ ] Implement `ui.tabs(id, tabs, options?)`.
- [ ] Render CSS-only radio tab markup, or document and test a `<details>` fallback if needed.
- [ ] Normalize container and tab IDs.
- [ ] Suffix duplicate tab IDs.
- [ ] Resolve selected tab by id or index, falling back to first non-disabled tab.
- [ ] Render disabled tab labels safely and unselectably.
- [ ] Normalize tab content through existing UI node normalization.
- [ ] Add tests for selected/disabled/duplicate/escaping behavior.
- [ ] Commit tabs implementation.

### T05 — Example integration

- [ ] Update at least one example detail page to use `ui.badge`, `ui.sql`/`ui.codeBlock`, `ui.jsonBlock`, and `ui.tabs`.
- [ ] Add/update retro CSS classes for code blocks, badges, and tabs.
- [ ] Add a ticket-local serve/curl smoke script that checks rendered component classes.
- [ ] Validate with Go tests and smoke script.
- [ ] Commit example integration.

### T06 — Documentation updates

- [ ] Update `db-browser help js-api-reference` with all new APIs and safety notes.
- [ ] Update `db-browser help user-guide` with inspection/debug page examples.
- [ ] Mention components in README if useful.
- [ ] Validate embedded help rendering.
- [ ] Commit documentation updates.

### T07 — Final validation and handoff

- [ ] Run `go test ./...`.
- [ ] Run new component smoke scripts.
- [ ] Run `docmgr doctor --ticket DB-BROWSER-UIDSL-COMPONENTS --stale-after 30`.
- [ ] Update diary with final validation and review instructions.
- [ ] Commit final ticket docs if needed.

## Future follow-ups

- [ ] Add real client-side copy behavior for `ui.codeBlock` if a static JS asset path becomes available.
- [ ] Add CSS counter-based line numbers to the shared theme.
- [ ] Add ARIA-focused tab keyboard support if client-side JS becomes acceptable.
- [ ] Consider `ui.kv(...)` for summary/detail pages.
