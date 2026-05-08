---
Title: Scoped ui.table query state implementation guide
Ticket: DB-BROWSER-SCOPED-TABLE-QUERY
Status: active
Topics:
  - db-browser
  - ui-dsl
  - web-ui
  - goja
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
  - path: /home/manuel/code/wesen/2026-05-07--db-browser/internal/uidsl/table.go
    note: ui.table render context, query parsing, filtering, sorting, pagination, and href generation live here.
  - path: /home/manuel/code/wesen/2026-05-07--db-browser/internal/uidsl/table_filters_test.go
    note: Existing tests for filter, sort, pagination behavior.
  - path: /home/manuel/code/wesen/2026-05-07--db-browser/internal/uidsl/table_rich_test.go
    note: Existing tests for rich table render context and pagination links.
  - path: /home/manuel/code/wesen/2026-05-07--db-browser/examples/generic-browser/scripts/app.js
    note: Example with multiple tables and detail pages that can exercise scoped state.
  - path: /home/manuel/code/wesen/2026-05-07--db-browser/internal/doc/topics/js-api-reference.md
    note: Embedded JavaScript API docs that must describe scoped table query state.
ExternalSources:
  - /home/manuel/workspaces/2026-05-02/use-sessionstream-coinvault/2026-03-16--gec-rag/ttmp/2026/05/07/SQLITE-TRACE-VERBS--design-sqlite-trace-inspection-verbs/scripts/serve/trace_browser_app.js
Summary: "Design for making ui.table pagination/sorting/filtering state table-scoped so multi-table pages can paginate independently while preserving existing global query behavior."
LastUpdated: 2026-05-07T23:35:00-04:00
WhatFor: "Use this when implementing table-scoped query parameters for ui.table in db-browser."
WhenToUse: "Read before changing ui.table query parsing, pagination href generation, filter forms, or examples with multiple tables per page."
---

# Scoped ui.table query state implementation guide

## Executive Summary

`ui.table` currently reads pagination, sorting, and filter state from a single global query namespace:

```text
?page=2&sort=name&dir=asc&q=alice&filter.status=open
```

That works for pages with one primary table. It becomes ambiguous on pages with several tables, such as the SQLite trace browser where one page can render schema objects, backend records, frontend records, nested record tabs, raw row tables, and JSON/code panels. If every table reads `page`, `sort`, and `filter.*`, then clicking pagination for one table can affect all tables on the page.

This ticket proposes a future `ui.table` enhancement: support **scoped table query state** while preserving the existing global behavior for backwards compatibility. A table with id `schema` should be able to read and emit parameters like:

```text
?schema.page=2&schema.sort=name&schema.dir=asc&schema.filter.type=table
```

or, if a different prefix convention is chosen:

```text
?t.schema.page=2&t.schema.sort=name&t.schema.filter.type=table
```

The implementation should be small and local to `internal/uidsl/table.go`: parse scoped values into `RenderContext`, update filter forms and pagination/sort links to emit scoped keys, add tests, update docs, and update example apps.

## Problem Statement

The immediate issue was reported from a real db-browser app:

```text
I get to ?page=2 but I'm still on page 1 it seems.
```

The direct cause in that app was script-level: the helper rendered every table with an empty query object:

```js
.render({ query: {} })
```

Changing that to pass `req.query` fixes basic pagination:

```js
.render({ query: req.query })
```

But the report exposes a deeper DSL limitation. The table DSL assumes one query namespace per page. Consider a page with two independent tables:

```js
ui.table.fromRows("backend", backendRows).features(f => f.pagination()).render({ query: req.query })
ui.table.fromRows("frontend", frontendRows).features(f => f.pagination()).render({ query: req.query })
```

When the user clicks “Next” on the backend table, the link is currently:

```text
?page=2
```

Now both tables see `page=2`. The frontend table may skip to page 2 even though the user did not interact with it. Likewise, `sort=name` applies to every table that has a sortable `name` column, and `filter.status=open` applies to every table with a `status` filter.

This is especially problematic for inspection/debug pages where multiple tables are normal, not exceptional.

## Goals

1. Let each table maintain independent pagination, sorting, and filtering state.
2. Preserve existing global query behavior for simple one-table pages.
3. Keep generated URLs readable and inspectable.
4. Avoid requiring app authors to manually rewrite query parameters.
5. Support both static `table.fromRows` tables and dynamic `.data(ctx => ...)` tables.
6. Update docs so LLM-generated apps use scoped state when rendering multiple tables.
7. Provide a migration path for existing scripts.

## Non-goals

- Do not implement client-side state or local storage.
- Do not introduce JavaScript event dispatch or partial refreshes in this ticket.
- Do not solve independent tab state unless it naturally shares the same scoping helper.
- Do not change SQL-backed `.data(ctx => ...)` semantics beyond passing scoped context.

## Current `ui.table` behavior

The relevant code lives in `internal/uidsl/table.go`.

Today `TableBuilder.context(input)` does roughly this:

```go
query := mapFromAny(input["query"])
pageSize := intFromAny(query["limit"], t.features.PageSize)
pageIndex := intFromAny(query["page"], 1)
sortKey := stringFromAny(query["sort"])
dir := stringFromAny(query["dir"])
filter := filterMapFromQuery(query)
```

Links are generated with `queryHref(ctx.Query, updates)`:

```go
queryHref(ctx.Query, map[string]any{"page": page + 1})
queryHref(ctx.Query, map[string]any{"sort": column.Name, "dir": dir, "page": 1})
```

Filter forms render names such as:

```html
<input name="q" ...>
<input name="filter.status" ...>
<input type="hidden" name="sort" ...>
```

This design is simple and useful. The scoped design should keep it as the default for existing apps unless authors opt into scoped state.

## Proposed API

### Recommended public API

Add a render option named `scope`:

```js
ui.table.fromRows("backend", rows)
  .features(f => f.filters().pagination({ size: 50 }).sorting())
  .render({ query: req.query, scope: "backend" })
```

When `scope` is present, the table reads and writes scoped parameters:

```text
?backend.page=2
?backend.limit=50
?backend.sort=ordinal
?backend.dir=desc
?backend.q=error
?backend.filter.kind=transport
```

The table id can be used as the default scope if requested explicitly:

```js
.render({ query: req.query, scope: true })
```

This lets a helper be concise:

```js
function table(id, rows, configure) {
  const builder = ui.table.fromRows(id, rows || []);
  if (configure) builder.columns(configure);
  return builder
    .features(f => f.filters().pagination({ size: 50 }).sorting().columnPicker())
    .render({ query: req.query, scope: true });
}
```

### Optional fluent API

A later ergonomic addition could be:

```js
ui.table.fromRows("backend", rows)
  .scope("backend")
  .render({ query: req.query })
```

or:

```js
ui.table.fromRows("backend", rows)
  .scoped()
  .render({ query: req.query })
```

This guide recommends implementing `render({ scope })` first because it is local to context parsing and avoids expanding the fluent builder API.

## Query key convention

There are two reasonable key conventions.

### Option A: `<scope>.<key>`

```text
?backend.page=2&backend.sort=ordinal&backend.filter.kind=transport
```

Pros:

- Short and readable.
- Mirrors existing `filter.<column>` syntax.
- Easy to type manually.

Cons:

- Can collide with app-owned query parameters that happen to start with `backend.`.
- Nested keys like `backend.filter.kind` require careful parsing.

### Option B: `t.<scope>.<key>`

```text
?t.backend.page=2&t.backend.sort=ordinal&t.backend.filter.kind=transport
```

Pros:

- Clear ownership: `t.` means table state.
- Lower risk of collision with app query parameters.
- Easier to scan when many app-level filters exist.

Cons:

- More verbose.
- Slightly less friendly when typing URLs by hand.

### Recommendation

Use Option A initially, but keep the implementation helper configurable enough that Option B can be adopted if collisions become a problem.

In code, centralize key construction:

```go
func scopedKey(scope, key string) string {
    if scope == "" { return key }
    return scope + "." + key
}
```

If we later choose `t.` prefix, only this helper changes:

```go
return "t." + scope + "." + key
```

## Backwards compatibility strategy

The safest compatibility rule is:

1. If no scope is configured, preserve current global behavior exactly.
2. If a scope is configured, prefer scoped keys.
3. Optionally allow scoped tables to fall back to global keys only when the scoped key is absent.

The fallback choice is subtle. Suppose the URL contains:

```text
?page=2&backend.page=1
```

A scoped backend table should use `backend.page=1`, not global `page=2`.

If the URL contains only:

```text
?page=2
```

Should scoped backend use page 2? There are two possible rules:

- **Strict scoped mode:** no; scoped tables ignore global keys.
- **Fallback mode:** yes; scoped tables use global keys when scoped keys are absent.

Recommendation: implement **strict scoped mode**. It is easier to reason about and prevents surprising cross-table effects. Existing apps get global behavior by omitting `scope`.

## Proposed internal data model

Extend `RenderContext`:

```go
type RenderContext struct {
    Query  map[string]any    `json:"query"`
    Page   map[string]any    `json:"page"`
    Order  map[string]any    `json:"order"`
    State  map[string]any    `json:"state"`
    Filter map[string]any    `json:"filter"`
    Params map[string]string `json:"params"`
    Scope  string            `json:"scope"`
}
```

Add helpers:

```go
func tableScope(input map[string]any, tableID string) string
func scopedQueryValue(query map[string]any, scope, key string) any
func scopedQueryUpdate(scope string, updates map[string]any) map[string]any
func filterMapFromScopedQuery(query map[string]any, scope string) map[string]any
func scopedInputName(scope, name string) string
```

The context parser should become:

```go
scope := tableScope(input, t.ID)
pageSize := intFromAny(scopedQueryValue(query, scope, "limit"), t.features.PageSize)
pageIndex := intFromAny(scopedQueryValue(query, scope, "page"), 1)
sortKey := stringFromAny(scopedQueryValue(query, scope, "sort"))
dir := stringFromAny(scopedQueryValue(query, scope, "dir"))
filter := filterMapFromScopedQuery(query, scope)
```

When `scope == ""`, helpers should return current behavior.

## Filter form behavior

Current filter input names:

```text
q
filter.name
filter.status
sort
dir
```

Scoped filter input names:

```text
backend.q
backend.filter.name
backend.filter.status
backend.sort
backend.dir
```

Filter forms should preserve current table's sort state using scoped hidden inputs:

```html
<input type="hidden" name="backend.sort" value="ordinal">
<input type="hidden" name="backend.dir" value="desc">
```

The `Clear` link needs a decision. Today `Clear` points to `?`, clearing the entire query string. In scoped mode, `Clear` should remove only this table's scoped keys while preserving other query parameters.

Example:

```text
Before: ?backend.q=err&frontend.page=2&tab=schema
Clear backend: ?frontend.page=2&tab=schema
```

Add helper:

```go
func clearScopeHref(query map[string]any, scope string) string
```

If `scope == ""`, preserve the old behavior: `?`.

## Pagination and sorting links

Scoped pagination should preserve unrelated query parameters and update only the scoped keys.

Before clicking backend next:

```text
?frontend.page=3&tab=rows
```

Backend next link:

```text
?frontend.page=3&tab=rows&backend.page=2
```

Backend sort link:

```text
?frontend.page=3&tab=rows&backend.sort=ordinal&backend.dir=asc&backend.page=1
```

Implementation:

```go
queryHref(ctx.Query, scopedUpdates(ctx.Scope, map[string]any{
    "sort": column.Name,
    "dir": dir,
    "page": 1,
}))
```

## Multiple tables in tabs

The SQLite trace browser often renders tables inside tabs. Scoped table state should work inside tabs without knowing about the tab component.

Example:

```js
ui.tabs("raw-tabs-" + name, [
  {
    id: "table",
    label: "Table",
    content: table("raw-" + name, rows, ...),
  },
  {
    id: "json",
    label: "JSON",
    content: jsonBlock(rows, "first 500 rows"),
  },
])
```

If the table helper uses `scope: true`, table id `raw-people` creates keys:

```text
?raw-people.page=2
```

Open question: tab state itself currently uses no query state unless the app passes `selected: req.query.tab`. A future ticket could add scoped tabs state, but this ticket should not mix those concerns.

## Example migration for the trace browser script

The immediate fix was script-only:

```js
let currentQuery = {};
function useQuery(req) { currentQuery = (req && req.query) || {}; }

function table(id, rows, configure) {
  const builder = ui.table.fromRows(id, rows || []);
  if (configure) builder.columns(configure);
  else builder.columns(c => c.text("name").text("n"));

  return builder
    .features(f => f.filters().pagination({ size: 50 }).sorting().columnPicker())
    .render({ query: currentQuery });
}
```

After scoped table support, the helper should become:

```js
let currentQuery = {};
function useQuery(req) { currentQuery = (req && req.query) || {}; }

function table(id, rows, configure) {
  const builder = ui.table.fromRows(id, rows || []);
  if (configure) builder.columns(configure);
  else builder.columns(c => c.text("name").text("n"));

  return builder
    .features(f => f.filters().pagination({ size: 50 }).sorting().columnPicker())
    .render({ query: currentQuery, scope: true });
}
```

If a page needs a custom scope:

```js
.render({ query: currentQuery, scope: "schema-main" })
```

## Implementation Plan

### T01 — Tests first for scoped parsing

Add focused tests in `internal/uidsl/table_scoped_query_test.go`:

1. Global behavior remains unchanged when no scope is passed.
2. Scoped page reads `schema.page`, not global `page`.
3. Scoped sort reads `schema.sort` and `schema.dir`.
4. Scoped filters read `schema.q`, `schema.filter.name`, and `schema.filter_name` if underscore compatibility is desired.
5. Strict scoped mode ignores global `page` when scope is set.

Example JS test shape:

```js
ui.table.fromRows("schema", rows)
  .features(f => f.pagination({ size: 2 }).sorting().filters())
  .render({ query: { "schema.page": "2" }, scope: "schema" })
```

Assert that page 2 rows render and page 1 rows do not.

### T02 — Add scope to `RenderContext`

Modify `RenderContext` and `TableBuilder.context` to parse scope from render input.

Scope parsing rules:

```go
switch input["scope"] {
case string:
    use normalized string if non-empty
case bool true:
    use table ID
case nil / false:
    no scope
}
```

Normalize with a query-safe token helper. It should preserve enough readability for table ids such as `raw-people`.

### T03 — Add scoped query helper functions

Implement helpers in `table.go` or a small `table_query.go`:

```go
func queryKey(scope, key string) string
func queryValue(query map[string]any, scope, key string) any
func queryUpdates(scope string, updates map[string]any) map[string]any
func removeScopedKeys(query map[string]any, scope string) map[string]any
func scopedFilterMap(query map[string]any, scope string) map[string]any
```

Prefer small pure functions with unit tests.

### T04 — Update href generation

Update:

- sortable header hrefs;
- previous/next pagination hrefs;
- filter form hidden inputs;
- filter input names;
- clear link.

### T05 — Update examples and docs

Update:

- `examples/generic-browser/scripts/app.js` helper to use `scope: true` when there are multiple tables on a page.
- `internal/doc/topics/js-api-reference.md` table context docs.
- `internal/doc/topics/user-guide.md` multi-table page guidance.
- Possibly `app-playbook` to instruct LLMs to use scoped tables for multi-table pages.

### T06 — Validate against real trace-browser pattern

Create or update a ticket-local smoke script:

```bash
ttmp/.../scripts/001-scoped-table-query-smoke.sh
```

It should:

1. Build db-browser.
2. Create a SQLite DB with enough rows for multiple pages.
3. Serve an example page with at least two tables.
4. Fetch `?left.page=2&right.page=1`.
5. Assert left table shows page 2 rows and right table shows page 1 rows.
6. Assert generated pagination links preserve the other table's scoped query state.

## Pseudocode

### Scope extraction

```go
func (t *TableBuilder) context(input map[string]any) RenderContext {
    query := mapFromAny(input["query"])
    scope := tableScope(input["scope"], t.ID)

    pageSize := intFromAny(queryValue(query, scope, "limit"), t.features.PageSize)
    pageIndex := intFromAny(queryValue(query, scope, "page"), 1)
    sortKey := stringFromAny(queryValue(query, scope, "sort"))
    dir := normalizeDir(queryValue(query, scope, "dir"))

    return RenderContext{
        Query: query,
        Scope: scope,
        Filter: filterMapFromQuery(query, scope),
        Page: map[string]any{
            "index": pageIndex,
            "limit": pageSize,
            "offset": (pageIndex-1)*pageSize,
        },
        Order: map[string]any{"key": sortKey, "dir": dir},
    }
}
```

### Scoped query value

```go
func queryValue(query map[string]any, scope, key string) any {
    if scope == "" {
        return query[key]
    }
    return query[scope+"."+key]
}
```

### Scoped updates

```go
func scopedUpdates(scope string, updates map[string]any) map[string]any {
    if scope == "" { return updates }
    ret := map[string]any{}
    for k, v := range updates {
        ret[scope+"."+k] = v
    }
    return ret
}
```

### Scoped filters

```go
func filterMapFromQuery(query map[string]any, scope string) map[string]any {
    ret := map[string]any{}
    prefix := ""
    if scope != "" { prefix = scope + "." }

    qKey := prefix + "q"
    if query[qKey] != "" { ret["q"] = query[qKey] }

    filterPrefix := prefix + "filter."
    for key, value := range query {
        if strings.HasPrefix(key, filterPrefix) {
            ret[strings.TrimPrefix(key, filterPrefix)] = value
        }
    }
    return ret
}
```

## Test cases

### Scoped pagination

Input rows:

```js
[
  { id: 1 }, { id: 2 }, { id: 3 }, { id: 4 }
]
```

Render:

```js
ui.table.fromRows("orders", rows)
  .features(f => f.pagination({ size: 2 }))
  .render({ query: { "orders.page": "2" }, scope: true })
```

Expected:

- rows 3 and 4 render;
- rows 1 and 2 do not;
- pagination says `Page 2 of 2 (4 rows)`;
- previous link contains `orders.page=1`.

### Scoped sorting does not affect another table

Render two tables with one shared query:

```js
const query = { "left.sort": "name", "left.dir": "desc", "right.sort": "id" };
```

Expected:

- left sorts by name desc;
- right sorts by id asc;
- no global `sort` is read.

### Clear only one scope

Input query:

```js
{
  "left.q": "error",
  "left.page": "2",
  "right.page": "3",
  "tab": "schema"
}
```

Clear href for left should preserve:

```text
?right.page=3&tab=schema
```

## Documentation updates

Add to `js-api-reference` under the table section:

```js
ui.table.fromRows("backend", rows)
  .features(f => f.filters().pagination({ size: 50 }).sorting())
  .render({ query: req.query, scope: true })
```

Explain:

- `scope: true` uses table id as the query namespace.
- `scope: "custom"` uses a custom namespace.
- no scope preserves legacy global parameters.

Add a warning to `app-playbook`:

> If a generated page renders more than one independently pageable/filterable table, pass `scope: true` to each table render call.

## Alternatives Considered

### Alternative A: keep only global query state

Rejected for multi-table inspection pages. It is simple but causes cross-table interactions.

### Alternative B: require app authors to manually pass filtered query objects

Example:

```js
.render({ query: pickTableQuery(req.query, "backend") })
```

Rejected because every app would invent its own convention, and pagination links would still need to know how to write scoped keys.

### Alternative C: use hash fragments

Example:

```text
/page#backend.page=2
```

Rejected because the server does not receive hash fragments in HTTP requests.

### Alternative D: POST state for table navigation

Rejected for now. GET query parameters are inspectable, bookmarkable, and work without JavaScript.

### Alternative E: namespace all tables by default

Tempting, but rejected for compatibility. Existing one-table apps and docs use `?page=2`; changing the default would break URLs and surprise users.

## Design Decisions

### Decision: opt-in scoping first

Scope is explicit through `render({ scope })`. This preserves existing behavior and makes migration straightforward.

### Decision: strict scoped mode

When scope is set, the table ignores global `page`, `sort`, and `filter.*`. This prevents accidental cross-table coupling.

### Decision: scope should be query-visible

Use readable query keys rather than opaque encoded state. Debugging internal tools should remain easy.

### Decision: page helpers can hide the boilerplate

Apps with many tables can define:

```js
function table(id, rows, configure) {
  return ui.table.fromRows(id, rows)
    .features(...)
    .render({ query: currentQuery, scope: true })
}
```

so authors do not repeat the option at every call site.

## Open Questions

1. Should the canonical key convention be `<scope>.<key>` or `t.<scope>.<key>`?
2. Should scoped filter input names support both `scope.filter.name` and `scope.filter_name`?
3. Should table `limit` be scoped too? Recommendation: yes.
4. Should the table id be normalized for query scope, or preserved exactly? Recommendation: normalize to a safe readable token.
5. Should `ui.tabs` eventually use the same scope convention for selected tab state?
6. Should `scope: true` be recommended in all docs, or only multi-table pages?

## Risks

### Query URLs become long

Multi-table pages can produce URLs like:

```text
?backend.page=2&backend.sort=ordinal&frontend.page=4&schema.q=record
```

That is acceptable for internal tools. It is still more transparent than hidden state.

### Table IDs may be unstable

If table ids are generated from row names and those names change, URLs break. This is already true for element IDs and links. Docs should encourage stable table IDs.

### Multiple tables with the same ID

Scoped query state assumes unique table ids. The DSL cannot fully prevent duplicate ids, but examples and docs should treat duplicate table ids as a bug.

## Recommended future implementation order

1. Add tests for scoped pagination/sorting/filtering.
2. Add `Scope` to `RenderContext` and parse `render({ scope })`.
3. Add scoped query helper functions.
4. Update href and form generation.
5. Update docs and examples.
6. Add smoke script with two independent tables.
7. Validate against the trace browser script pattern.

## References

- Current table DSL: `/home/manuel/code/wesen/2026-05-07--db-browser/internal/uidsl/table.go`
- Current JS API reference: `/home/manuel/code/wesen/2026-05-07--db-browser/internal/doc/topics/js-api-reference.md`
- Real script that exposed the problem: `/home/manuel/workspaces/2026-05-02/use-sessionstream-coinvault/2026-03-16--gec-rag/ttmp/2026/05/07/SQLITE-TRACE-VERBS--design-sqlite-trace-inspection-verbs/scripts/serve/trace_browser_app.js`
