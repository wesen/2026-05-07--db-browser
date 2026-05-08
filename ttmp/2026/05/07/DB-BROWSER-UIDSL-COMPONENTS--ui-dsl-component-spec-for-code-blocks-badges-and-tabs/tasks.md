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

### T08 — Syntax highlighting and closer trace-browser styling

- [x] Add safe server-side token highlighting for SQL, JSON, and JavaScript code blocks.
- [x] Render highlighted tokens as escaped span-wrapped UI nodes, not raw HTML.
- [x] Add/adjust tests for keyword/string/key/number highlighting.
- [x] Update the generic browser CSS toward the supplied classic Macintosh trace-browser reference.
- [x] Validate with Go tests, smoke script, and a Playwright visual check.
- [x] Commit syntax highlighting and styling update.

### T09 — Server-interactive UI research article

- [x] Write a textbook-style Obsidian article exploring backend-dispatched ui.dsl events.
- [x] Include architecture diagrams, event time diagrams, implementation phases, examples, and failure modes.
- [x] Store the original article in the Obsidian vault under `Projects/2026/05/08/`.
- [x] Copy the article back into this ticket with `cp`.
- [x] Update diary/changelog and validate ticket hygiene.
- [x] Commit ticket copy and bookkeeping.

### T10 — Fix CSS-only tab switching

- [x] Diagnose why `ui.tabs` labels/radios did not switch SQL/Metadata panels in real apps.
- [x] Move radio inputs to be direct siblings of the panels so CSS `:checked ~ .ui-tabs__panels` can work.
- [x] Emit a per-instance component style block mapping each radio to its matching panel.
- [x] Keep labels in a `ui-tabs__tablist` so visible tab clicks still check hidden radios.
- [x] Add/adjust tests for emitted tab switching CSS.
- [x] Validate with Go tests and Playwright label click.
- [x] Commit the tab fix.

### T07 — Final validation and handoff

- [x] Run `go test ./...`.
- [x] Run new component smoke scripts.
- [x] Run `docmgr doctor --ticket DB-BROWSER-UIDSL-COMPONENTS --stale-after 30`.
- [x] Update diary with final validation and review instructions.
- [x] Commit final ticket docs if needed.

## Future follow-ups

- [ ] Add real client-side copy behavior for `ui.codeBlock` if a static JS asset path becomes available.
- [ ] Add CSS counter-based line numbers to the shared theme.
- [ ] Add ARIA-focused tab keyboard support if client-side JS becomes acceptable.
- [ ] Consider `ui.kv(...)` for summary/detail pages.
