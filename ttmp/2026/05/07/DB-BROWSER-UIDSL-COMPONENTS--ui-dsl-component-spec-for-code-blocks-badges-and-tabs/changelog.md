# Changelog

## 2026-05-07

- Initial workspace created


## 2026-05-07

Created ticket, implementation guide, task list, and initial diary for ui.dsl codeBlock badge tabs components.

### Related Files

- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-UIDSL-COMPONENTS--ui-dsl-component-spec-for-code-blocks-badges-and-tabs/design-doc/01-ui-dsl-codeblock-badge-tabs-implementation-guide.md — Component implementation guide
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-UIDSL-COMPONENTS--ui-dsl-component-spec-for-code-blocks-badges-and-tabs/reference/01-implementation-diary.md — Initial implementation diary
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-UIDSL-COMPONENTS--ui-dsl-component-spec-for-code-blocks-badges-and-tabs/tasks.md — Implementation task list


## 2026-05-07

Implemented ui.codeBlock, ui.badge, ui.tabs, aliases, and component tests.

### Related Files

- /home/manuel/code/wesen/2026-05-07--db-browser/internal/uidsl/components.go — New UI DSL component implementations
- /home/manuel/code/wesen/2026-05-07--db-browser/internal/uidsl/components_test.go — Component render contract and escaping tests
- /home/manuel/code/wesen/2026-05-07--db-browser/internal/uidsl/module.go — Goja export wiring for new components
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-UIDSL-COMPONENTS--ui-dsl-component-spec-for-code-blocks-badges-and-tabs/reference/01-implementation-diary.md — Diary step 2
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-UIDSL-COMPONENTS--ui-dsl-component-spec-for-code-blocks-badges-and-tabs/tasks.md — Marked T02-T04 core implementation tasks complete


## 2026-05-07

Integrated new ui.dsl components into the generic browser example and added a smoke script.

### Related Files

- /home/manuel/code/wesen/2026-05-07--db-browser/examples/generic-browser/scripts/app.js — Example detail page uses badge tabs SQL and JSON code blocks
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-UIDSL-COMPONENTS--ui-dsl-component-spec-for-code-blocks-badges-and-tabs/reference/01-implementation-diary.md — Diary step 3
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-UIDSL-COMPONENTS--ui-dsl-component-spec-for-code-blocks-badges-and-tabs/scripts/001-uidsl-components-smoke.sh — Component integration smoke validation
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-UIDSL-COMPONENTS--ui-dsl-component-spec-for-code-blocks-badges-and-tabs/tasks.md — Marked T05 complete


## 2026-05-07

Documented new ui.dsl inspection components in embedded help and README.

### Related Files

- /home/manuel/code/wesen/2026-05-07--db-browser/README.md — UI DSL feature summary includes components
- /home/manuel/code/wesen/2026-05-07--db-browser/internal/doc/topics/js-api-reference.md — Full component API reference
- /home/manuel/code/wesen/2026-05-07--db-browser/internal/doc/topics/user-guide.md — Inspection component usage example
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-UIDSL-COMPONENTS--ui-dsl-component-spec-for-code-blocks-badges-and-tabs/reference/01-implementation-diary.md — Diary step 4
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-UIDSL-COMPONENTS--ui-dsl-component-spec-for-code-blocks-badges-and-tabs/tasks.md — Marked T06 complete


## 2026-05-07

Completed final validation and handoff for ui.dsl inspection components.

### Related Files

- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-UIDSL-COMPONENTS--ui-dsl-component-spec-for-code-blocks-badges-and-tabs/reference/01-implementation-diary.md — Final validation diary step
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-UIDSL-COMPONENTS--ui-dsl-component-spec-for-code-blocks-badges-and-tabs/scripts/001-uidsl-components-smoke.sh — Final smoke validation script
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-UIDSL-COMPONENTS--ui-dsl-component-spec-for-code-blocks-badges-and-tabs/tasks.md — Final task completion state


## 2026-05-07

Added safe code block syntax highlighting and updated generic browser CSS toward the supplied Macintosh trace-browser reference.

### Related Files

- /home/manuel/code/wesen/2026-05-07--db-browser/examples/generic-browser/scripts/app.js — Reference-inspired classic Mac CSS shell
- /home/manuel/code/wesen/2026-05-07--db-browser/internal/doc/topics/js-api-reference.md — Code block highlighting documentation
- /home/manuel/code/wesen/2026-05-07--db-browser/internal/uidsl/components.go — Safe SQL/JSON/JS token highlighting
- /home/manuel/code/wesen/2026-05-07--db-browser/internal/uidsl/components_test.go — Highlighting test expectations
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-UIDSL-COMPONENTS--ui-dsl-component-spec-for-code-blocks-badges-and-tabs/reference/01-implementation-diary.md — Diary step 6
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-UIDSL-COMPONENTS--ui-dsl-component-spec-for-code-blocks-badges-and-tabs/tasks.md — Added T08 highlighting/styling task


## 2026-05-07

Added Obsidian research article and ticket copy for server-interactive ui.dsl backend-dispatched events.

### Related Files

- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-UIDSL-COMPONENTS--ui-dsl-component-spec-for-code-blocks-badges-and-tabs/reference/01-implementation-diary.md — Diary step 7
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-UIDSL-COMPONENTS--ui-dsl-component-spec-for-code-blocks-badges-and-tabs/reference/02-server-interactive-ui-proposal.md — Ticket copy of Obsidian research article
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-UIDSL-COMPONENTS--ui-dsl-component-spec-for-code-blocks-badges-and-tabs/tasks.md — Added T09 research article task


## 2026-05-07

Fixed CSS-only ui.tabs switching by making radios and panels siblings and emitting per-instance tab CSS.

### Related Files

- /home/manuel/code/wesen/2026-05-07--db-browser/internal/doc/topics/js-api-reference.md — Documented per-instance tab style behavior
- /home/manuel/code/wesen/2026-05-07--db-browser/internal/uidsl/components.go — Tabs markup and generated CSS fix
- /home/manuel/code/wesen/2026-05-07--db-browser/internal/uidsl/components_test.go — Tab switching CSS test expectation
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-UIDSL-COMPONENTS--ui-dsl-component-spec-for-code-blocks-badges-and-tabs/reference/01-implementation-diary.md — Diary step 8
- /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-UIDSL-COMPONENTS--ui-dsl-component-spec-for-code-blocks-badges-and-tabs/tasks.md — Added T10 tab fix task

