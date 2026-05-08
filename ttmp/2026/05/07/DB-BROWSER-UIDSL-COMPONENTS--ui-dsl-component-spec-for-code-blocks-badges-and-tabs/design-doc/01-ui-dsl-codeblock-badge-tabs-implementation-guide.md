---
Title: ui.dsl codeBlock badge tabs implementation guide
Ticket: DB-BROWSER-UIDSL-COMPONENTS
Status: active
Topics:
  - db-browser
  - ui-dsl
  - goja
  - server-rendered-ui
  - documentation
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
  - path: /home/manuel/code/wesen/2026-05-07--db-browser/internal/uidsl/module.go
    note: ui.dsl module export registration point.
  - path: /home/manuel/code/wesen/2026-05-07--db-browser/internal/uidsl/node.go
    note: Server-rendered UI node model used by all new components.
  - path: /home/manuel/code/wesen/2026-05-07--db-browser/internal/uidsl/render.go
    note: HTML escaping and attribute rendering contracts.
  - path: /home/manuel/code/wesen/2026-05-07--db-browser/internal/uidsl/table.go
    note: Existing rich component builder patterns and CSS token helpers.
  - path: /home/manuel/code/wesen/2026-05-07--db-browser/internal/doc/topics/js-api-reference.md
    note: Embedded JavaScript API reference that must be updated after implementation.
ExternalSources: []
Summary: "Design and implementation guide for adding ui.codeBlock, ui.badge, and ui.tabs to db-browser's server-rendered ui.dsl."
LastUpdated: 2026-05-07T22:25:00-04:00
WhatFor: "Use this when implementing reusable server-rendered inspection/debug components in ui.dsl."
WhenToUse: "Read before changing internal/uidsl component APIs, render contracts, tests, examples, or JavaScript API docs."
---

# ui.dsl codeBlock badge tabs implementation guide

## Executive Summary

This ticket adds three reusable, server-rendered `ui.dsl` components for db-browser inspection and debug pages:

- `ui.codeBlock(language, source, options?)` for escaped formatted code and preformatted text.
- `ui.badge(value, options?)` for compact status/type labels outside tables.
- `ui.tabs(id, tabs, options?)` for compact multi-view detail pages that work without client-side JavaScript.

The implementation should fit the current db-browser UI model: Goja exports JavaScript functions from `internal/uidsl/module.go`; those functions return `Node` values from `internal/uidsl/node.go`; `internal/uidsl/render.go` escapes text and renders attributes; examples and help docs demonstrate usage. The components must be safe by default for database text, request text, JSON payloads, SQL schema strings, and generated app output.

## Problem Statement

The current `ui.dsl` has low-level tag helpers and a rich table component, but inspection pages still need to hand-author common detail-page primitives. Developers and LLM-generated apps repeatedly need to show SQL, JSON, scripts, status labels, record metadata, raw payloads, and alternate views of the same record.

Without first-class components, apps tend to drift toward unsafe patterns such as:

```js
ui.raw(`<pre>${row.sql}</pre>`)
```

or inconsistent patterns such as custom badge classes and ad hoc detail layouts. This is especially risky for db-browser because its most common content comes from databases and request parameters. The new components should make the safe path the convenient path.

## Goals

1. Add public JavaScript APIs:
   - `ui.codeBlock(language, source, options?)`
   - `ui.sql(source, options?)`
   - `ui.jsonBlock(value, options?)`
   - `ui.js(source, options?)`
   - `ui.badge(value, options?)`
   - `ui.tabs(id, tabs, options?)`
2. Preserve server-side rendering only; no client-side JavaScript required.
3. Escape untrusted source, labels, values, and content by default.
4. Normalize language, value, id, and tab tokens to CSS/DOM-safe strings.
5. Update `js-api-reference` and `user-guide` with signatures, examples, and safety notes.
6. Add unit tests that assert render contracts and escaping.
7. Add or update an example page that uses all three components together.

## Non-goals

- Full syntax highlighting is not part of this ticket. The `language-*` class should make future highlighting possible.
- Clipboard functionality does not need to work initially. A copy button/affordance may be inert.
- Client-side keyboard tab behavior is not required. CSS-only radio tabs or `<details>` fallback is acceptable.
- A general component styling system is not required. The components should expose stable classes for examples/themes to style.

## Existing Architecture

### Node model

`internal/uidsl/node.go` defines the renderable types:

```go
type Node interface{ isNode() }
type Document struct { Title string; Head []Node; Body []Node }
type Element struct { Tag string; Attrs map[string]any; Children []Node }
type Text struct{ Value string }
type RawHTML struct{ Value string }
type Fragment struct{ Children []Node }
```

The new components should return these existing node types. They do not need a new renderer path.

### Escaping model

`internal/uidsl/render.go` escapes `Text` values with `html.EscapeString`. Therefore component implementations should use `Text` for database/request/source text. They should use `RawHTML` only for trusted component-owned markup, and ideally not at all.

### Module registration

`internal/uidsl/module.go` exports functions to JavaScript by setting fields on the CommonJS `exports` object. New component functions should be registered there, preferably delegating implementation to a new focused file such as `internal/uidsl/components.go`.

### Current component precedent

`internal/uidsl/table.go` already contains useful patterns:

- Goja callable wrappers for fluent builders.
- Map option parsing with `map[string]any`.
- CSS-safe token normalization via `cssToken`.
- Rendering by composing `Element`, `Text`, and `Fragment`.

The new components should reuse or centralize the token normalization helper rather than create divergent implementations.

## Proposed JavaScript API

## 1. `ui.codeBlock(language, source, options?)`

### Purpose

Render escaped code/text with optional title, line numbers, wrapping, copy affordance, max-height style, and language class for future syntax highlighting.

### Signature

```js
ui.codeBlock(language, source, options?)
```

### Examples

```js
ui.codeBlock("sql", row.sql, {
  title: "CREATE VIEW delivery_chain",
  lineNumbers: true,
  wrap: false,
  copy: true,
});

ui.codeBlock("json", JSON.stringify(obj, null, 2));

ui.codeBlock("javascript", scriptSource, {
  title: "trace_browser_app.js",
});
```

### Convenience aliases

```js
ui.sql(source, options?)       // ui.codeBlock("sql", source, options)
ui.jsonBlock(value, options?)  // pretty-print object or JSON string, then codeBlock("json", ...)
ui.js(source, options?)        // ui.codeBlock("javascript", source, options)
```

### Options

```ts
type CodeBlockOptions = {
  title?: string;
  lineNumbers?: boolean;
  wrap?: boolean;       // default true
  copy?: boolean;       // render copy button/affordance; may be inert initially
  maxHeight?: string;   // e.g. "480px"
  class?: string;
};
```

### Render contract

When title or copy is enabled:

```html
<figure class="ui-codeblock ui-codeblock--sql ui-codeblock--wrap ui-codeblock--line-numbers [custom-class]">
  <figcaption class="ui-codeblock__caption">
    <span class="ui-codeblock__title">CREATE VIEW delivery_chain</span>
    <button class="ui-codeblock__copy" type="button">Copy</button>
  </figcaption>
  <pre class="ui-codeblock__pre"><code class="language-sql">escaped source</code></pre>
</figure>
```

When no title/copy is requested:

```html
<pre class="ui-codeblock ui-codeblock--sql ui-codeblock--wrap"><code class="language-sql">escaped source</code></pre>
```

### Implementation notes

- Empty or invalid language normalizes to `text`.
- `source` converts with `fmt.Sprint` unless it is `nil`, in which case it should render as an empty string.
- `jsonBlock` should:
  - pretty-print Goja objects/arrays through JSON marshal indent;
  - parse valid JSON strings and re-indent them;
  - fall back to plain escaped text for invalid JSON strings.
- `maxHeight` should become inline style on the `<pre>`: `max-height:<value>;overflow:auto`.
- `lineNumbers` can initially be CSS-class only, or can render line wrappers if the implementation remains simple. The first implementation should not sacrifice escaping safety for line number markup.

## 2. `ui.badge(value, options?)`

### Purpose

Render compact status/type labels outside tables. Tables already render badge columns internally; this public helper makes badges available in page headers, summaries, tabs, and key-value sections.

### Signature

```js
ui.badge(value, options?)
```

### Examples

```js
ui.badge("view");

ui.badge(row.transport_fanout, {
  tone: row.transport_fanout === "yes" ? "success" : "danger",
});

ui.badge("provider_normalize_delta", {
  tone: "info",
  title: "Geppetto stage",
});
```

### Options

```ts
type BadgeOptions = {
  tone?: "default" | "info" | "success" | "warning" | "danger" | "muted";
  title?: string;
  class?: string;
};
```

### Render contract

```html
<span class="ui-badge ui-badge--success ui-badge--value-yes" title="...">yes</span>
```

### Implementation notes

- Value text must render through `Text` and be escaped by the renderer.
- Unknown tone falls back to `default`.
- Value token uses the same CSS-safe normalization as table badges.
- `nil` values render as an empty badge or possibly `unknown`; this should be decided in tests. Recommended initial behavior: empty string, class `ui-badge--value-empty`.

## 3. `ui.tabs(id, tabs, options?)`

### Purpose

Render a no-JavaScript multi-view component for record/detail pages. Typical tabs are Summary, Pretty JSON, Raw JSON, SQL, Related rows, and Debug.

### Signature

```js
ui.tabs(id, tabs, options?)
```

### Example

```js
ui.tabs("record-tabs", [
  {
    id: "summary",
    label: "Summary",
    content: ui.kv(row),
  },
  {
    id: "json",
    label: "Raw JSON",
    content: ui.jsonBlock(row.raw_json, { lineNumbers: true }),
  },
  {
    id: "sql",
    label: "Schema SQL",
    content: ui.sql(row.sql),
  },
], {
  selected: "summary",
});
```

### Types

```ts
type TabSpec = {
  id?: string;
  label: string;
  content: UiNode | UiNode[] | string;
  disabled?: boolean;
};

type TabsOptions = {
  selected?: string | number; // tab id or index
  class?: string;
};
```

### Preferred render contract

CSS-only radio tabs:

```html
<div class="ui-tabs" id="record-tabs">
  <div class="ui-tabs__tablist" role="tablist">
    <input class="ui-tabs__radio" type="radio" name="record-tabs" id="record-tabs-summary" checked>
    <label class="ui-tabs__tab" for="record-tabs-summary">Summary</label>
    <input class="ui-tabs__radio" type="radio" name="record-tabs" id="record-tabs-json">
    <label class="ui-tabs__tab" for="record-tabs-json">Raw JSON</label>
  </div>
  <div class="ui-tabs__panels">
    <section class="ui-tabs__panel ui-tabs__panel--active" data-tab="summary">...</section>
    <section class="ui-tabs__panel" data-tab="json">...</section>
  </div>
</div>
```

### Acceptable fallback

If CSS-only radio tabs become too complex for the first implementation, render stacked `<details>` blocks. This preserves no-JS behavior, accessibility, and safety. The guide recommends starting with radio-tab markup because the user explicitly prefers it, but tests should focus on stable classes and safe content rather than a fragile exact full tree.

### Implementation notes

- Container `id` must normalize to a DOM-safe token. Empty id becomes `tabs` or a deterministic fallback.
- Tab ids normalize to DOM-safe tokens. Empty tab id can derive from label or index.
- Duplicate tab ids receive suffixes: `json`, `json-2`, `json-3`.
- Invalid `selected` chooses the first non-disabled tab.
- Disabled tabs render labels but should not be selected and should not render an enabled radio input. Use `disabled` attr and `ui-tabs__tab--disabled`.
- Labels render through `Text`.
- Content should pass through the existing `NormalizeExport` path so strings are escaped and `Node` values render normally.

## Proposed Go implementation plan

### Files

Add:

- `internal/uidsl/components.go` — component constructors and helpers.
- `internal/uidsl/components_test.go` — codeBlock, badge, tabs tests.

Modify:

- `internal/uidsl/module.go` — export new JS functions.
- `internal/uidsl/table.go` — if `cssToken` should be shared/renamed, keep compatibility or move helper carefully.
- `internal/doc/topics/js-api-reference.md` — document the new API.
- `internal/doc/topics/user-guide.md` — add a short inspection/debug page example.
- examples — use components in at least one detail page.

### Registration pattern

In `Loader`:

```go
_ = exports.Set("codeBlock", func(language string, source goja.Value, options ...map[string]any) goja.Value { ... })
_ = exports.Set("sql", func(source goja.Value, options ...map[string]any) goja.Value { ... })
_ = exports.Set("jsonBlock", func(value goja.Value, options ...map[string]any) goja.Value { ... })
_ = exports.Set("js", func(source goja.Value, options ...map[string]any) goja.Value { ... })
_ = exports.Set("badge", func(value goja.Value, options ...map[string]any) goja.Value { ... })
_ = exports.Set("tabs", func(id string, tabs goja.Value, options ...map[string]any) (goja.Value, error) { ... })
```

Return `vm.ToValue(node)` so existing normalization/rendering works.

### Option parsing

Use small helpers:

```go
func firstOptions(options []map[string]any) map[string]any
func optionString(opts map[string]any, key string) string
func optionBool(opts map[string]any, key string, fallback bool) bool
```

Avoid making each component parse Goja values differently.

### Token normalization

The existing `cssToken` helper lowercases, strips unsupported chars, and returns `empty` for empty output. Reuse that behavior for:

- code block language class suffix;
- badge value class suffix;
- tab DOM ids.

If tab IDs need a slightly different helper, build on top of `cssToken` and ensure the first character is safe enough for practical HTML IDs.

## Testing plan

### Unit tests

Add tests for:

1. `ui.codeBlock("sql", "SELECT '<x>'")` escapes source.
2. `ui.codeBlock` with title/copy/wrap/lineNumbers renders the expected classes and caption/button.
3. `ui.sql`, `ui.js`, and `ui.jsonBlock` aliases work.
4. `ui.jsonBlock` pretty-prints objects and valid JSON strings.
5. `ui.jsonBlock` falls back to escaped text for invalid JSON strings.
6. `ui.badge("yes", { tone: "success", title: "..." })` renders escaped text and normalized classes.
7. Unknown badge tone falls back to default.
8. `ui.tabs` renders labels escaped, content escaped, selected tab checked/active, and disabled tabs disabled.
9. Duplicate tab IDs are suffixed.
10. Invalid selected tab selects the first non-disabled tab.

### Integration/example smoke

Add/update a ticket script such as:

```bash
ttmp/.../scripts/001-uidsl-components-tests.sh
```

It should run:

```bash
go test ./internal/uidsl -run 'Test.*(CodeBlock|Badge|Tabs)' -count=1
go test ./...
```

If examples are updated, add a serve/curl check for the rendered component classes.

## Documentation plan

Update `db-browser help js-api-reference` with:

- API signatures;
- render contracts;
- examples for SQL schema rendering;
- examples for JSON debug records;
- safety note: code block content is escaped; do not use `ui.raw` for database/request text.

Add a concise `user-guide` section showing:

```js
res.html(ui.page({ title: "Schema" },
  ui.h1(row.name),
  ui.badge(row.type),
  ui.codeBlock("sql", row.sql, {
    title: row.name,
    lineNumbers: true,
    copy: true,
  })
));
```

## Implementation sequence

1. Planning docs and task list.
2. `ui.codeBlock`, aliases, and tests.
3. `ui.badge` and tests.
4. `ui.tabs` and tests.
5. Example app update using all three components.
6. Embedded help docs update.
7. Full validation and final diary/changelog.

Commit after each coherent step so regressions are easy to bisect.

## Design Decisions

### Decision: server-rendered components only

All components return static HTML nodes. This keeps db-browser useful in restricted/internal environments, avoids asset bundling, and makes generated apps simple.

### Decision: escape with `Text`, not `ui.raw`

The public API should accept raw database/request strings and safely escape them. This is the central safety guarantee for `codeBlock`, badge labels, tab labels, and string tab content.

### Decision: CSS classes are part of the component contract

The components should be minimally styled by examples and future themes. Stable classes are therefore part of the contract even if the core package does not ship a global stylesheet yet.

### Decision: copy affordance may be inert initially

A rendered `Copy` button communicates affordance and provides a future hook. It should not require client-side JavaScript for the component to be valid.

### Decision: `tabs` should prefer CSS-only radio markup

Radio tabs preserve the requested compact UI without JavaScript. A `<details>` fallback is acceptable only if radio markup blocks progress.

## Alternatives Considered

### Alternative: keep using low-level tag helpers

Rejected because it repeats unsafe patterns and makes LLM-generated apps less consistent.

### Alternative: client-side tab JavaScript

Rejected for initial implementation because the stated requirement is no client-side JavaScript by default.

### Alternative: syntax highlighting now

Rejected because the requirement is language classes for future highlighting, not full highlighting in this ticket.

### Alternative: expose badge only through table columns

Rejected because inspection pages need badges in headings, summaries, tabs, and detail panels outside tables.

## Open Questions

1. Should `ui.codeBlock(..., { lineNumbers: true })` render actual line-number markup now, or only stable classes for CSS counters?
2. Should `ui.badge(null)` render empty text, `unknown`, or `null`?
3. Should `ui.tabs` include ARIA roles in the initial implementation, and if so how far should no-JS keyboard semantics go?
4. Should a shared retro theme stylesheet eventually style these classes globally?

## References

- Existing UI DSL registration: `/home/manuel/code/wesen/2026-05-07--db-browser/internal/uidsl/module.go`
- Existing rich table DSL: `/home/manuel/code/wesen/2026-05-07--db-browser/internal/uidsl/table.go`
- Existing JS API docs: `/home/manuel/code/wesen/2026-05-07--db-browser/internal/doc/topics/js-api-reference.md`
- Ticket source request: colleague component spec in user prompt for `DB-BROWSER-UIDSL-COMPONENTS`.
