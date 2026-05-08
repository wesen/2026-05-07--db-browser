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

- [x] Add `internal/uidsl/components.go` or equivalent focused component file.
- [x] Implement `ui.codeBlock(language, source, options?)`.
- [x] Implement `ui.sql(source, options?)`.
- [x] Implement `ui.js(source, options?)`.
- [x] Implement `ui.jsonBlock(value, options?)` with pretty-print behavior.
- [x] Normalize language tokens and default invalid/empty language to `text`.
- [x] Support `title`, `lineNumbers`, `wrap`, `copy`, `maxHeight`, and `class` options.
- [x] Add escaping and render-contract tests.
- [x] Commit codeBlock implementation.

### T03 — `ui.badge`

- [x] Implement `ui.badge(value, options?)`.
- [x] Support `tone`, `title`, and `class` options.
- [x] Normalize tone and value CSS tokens.
- [x] Ensure unknown tone falls back to `default`.
- [x] Add escaping, tone, and class tests.
- [x] Commit badge implementation.

### T04 — `ui.tabs`

- [x] Implement `ui.tabs(id, tabs, options?)`.
- [x] Render CSS-only radio tab markup, or document and test a `<details>` fallback if needed.
- [x] Normalize container and tab IDs.
- [x] Suffix duplicate tab IDs.
- [x] Resolve selected tab by id or index, falling back to first non-disabled tab.
- [x] Render disabled tab labels safely and unselectably.
- [x] Normalize tab content through existing UI node normalization.
- [x] Add tests for selected/disabled/duplicate/escaping behavior.
- [x] Commit tabs implementation.

### T05 — Example integration

- [x] Update at least one example detail page to use `ui.badge`, `ui.sql`/`ui.codeBlock`, `ui.jsonBlock`, and `ui.tabs`.
- [x] Add/update retro CSS classes for code blocks, badges, and tabs.
- [x] Add a ticket-local serve/curl smoke script that checks rendered component classes.
- [x] Validate with Go tests and smoke script.
- [x] Commit example integration.

### T06 — Documentation updates

- [x] Update `db-browser help js-api-reference` with all new APIs and safety notes.
- [x] Update `db-browser help user-guide` with inspection/debug page examples.
- [x] Mention components in README if useful.
- [x] Validate embedded help rendering.
- [x] Commit documentation updates.

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
