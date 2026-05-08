---
Title: db-browser JavaScript API Reference
Slug: js-api-reference
Short: Complete reference for the JavaScript modules available in db-browser serve scripts and repository verbs.
Topics:
- javascript
- js-api
- goja
- sqlite
- express
- ui-dsl
- jsverbs
Commands:
- db-browser serve
- db-browser verbs
- db-browser inspect modules
Flags:
- --db
- --scripts-dir
- --addr
- --dev
- --readonly
- --allow-writes
- --repository
- --verb-repository
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: GeneralTopic
---

This page is the JavaScript runtime reference for `db-browser`. It documents the modules available to server-side app scripts loaded by `db-browser serve` and to repository-scanned JS verbs executed by `db-browser verbs`.

`db-browser` scripts run in Goja with CommonJS-style `require(...)`. They do not run in Node, do not have browser DOM APIs, and should not assume npm packages unless you explicitly arrange a compatible require path. The core authoring model is: query SQLite with `db`, register routes with `express`, build HTML with `ui.dsl`, and optionally read YAML/files from disk.

## Runtime availability matrix

| Module | `serve` scripts | `verbs` commands | Notes |
| --- | --- | --- | --- |
| `db` / `database` | Yes when `--db` is set | Yes when `--db` is set | SQLite database module from go-go-goja, wrapped by db-browser write guards. |
| `express` | Yes | No | Route registration for the embedded HTTP server. |
| `ui.dsl` / `ui` | Yes | Yes | Server-rendered HTML node DSL and rich table helpers. |
| `fs` | Yes | Yes | Goja NodeJS-style filesystem module. |
| `path` | Yes | Yes | Path helper module. |
| `yaml` | Yes | Yes | YAML parse/stringify module from go-go-goja. |
| `time` | Yes | Yes | Time utilities from go-go-goja default modules. |
| `timer` | Yes | Yes | Timer utilities from go-go-goja default modules. |

If a script calls `require("db")` without `--db`, the module is not registered. Pass a database path:

```bash
db-browser serve --db app.db --scripts-dir scripts
db-browser verbs --db app.db examples builtin tables
```

## Script lifecycle

### `db-browser serve`

`serve` loads every `.js` file in `--scripts-dir` in sorted filename order. Scripts run at startup and should register routes:

```js
const express = require("express");
const app = express.app();

app.get("/", (req, res) => {
  res.send("hello");
});
```

Route handlers run later for each HTTP request. Restart the server after editing scripts.

### `db-browser verbs`

`verbs` scans configured repositories for explicit `__verb__` definitions. By default it scans only the embedded `builtin` repository. Additional repositories come from `.db-browser.yml`, `.db-browser.override.yml`, `DB_BROWSER_VERB_REPOSITORIES`, `--repository`, or `--verb-repository`.

List discovered verbs with structured output:

```bash
db-browser verbs list --fields path,repository,source,location --output table
```

## `require("db")` / `require("database")`

The database module exposes SQLite access for the database configured with `--db`.

```js
const db = require("db");
```

The `db` and `database` module names are aliases for the same database handle.

### `db.query(sql, ...args)`

Runs a SQL query and returns rows as an array of JavaScript objects.

```js
const rows = db.query(
  "SELECT id, name FROM customers WHERE segment = ? ORDER BY name",
  "vip",
) || [];
```

- **sql**: SQL string.
- **args**: Optional positional parameters.
- **Returns**: `Array<Record<string, any>>` for result rows. Use `|| []` defensively because some zero-row paths can encode as `null`.
- **Throws**: SQL errors, driver errors, or module availability errors when `--db` is missing.

Use placeholders for values. Do not interpolate user input into SQL strings:

```js
// Good
const rows = db.query("SELECT * FROM customers WHERE id = ?", req.params.id) || [];

// Bad
const rows = db.query(`SELECT * FROM customers WHERE id = ${req.params.id}`) || [];
```

### `db.exec(sql, ...args)`

Runs a SQL statement that does not return rows.

```js
db.exec("UPDATE customers SET segment = ? WHERE id = ?", "vip", id);
```

`db-browser serve` and verb runtimes default to read-only write guards. To allow writes you must opt in explicitly:

```bash
db-browser serve --db app.db --scripts-dir scripts --readonly=false --allow-writes
db-browser verbs --db app.db --readonly=false --allow-writes path to verb
```

If either condition is missing, write statements are rejected. Prefer read-only apps unless the user explicitly asks for mutations.

### Dynamic identifiers

SQL placeholders bind values, not table or column names. If an app accepts a table name in a route, validate it before inserting it into SQL:

```js
function quoteIdent(name) {
  if (!/^[A-Za-z_][A-Za-z0-9_]*$/.test(name)) {
    throw new Error("unsafe identifier: " + name);
  }
  return '"' + name.replace(/"/g, '""') + '"';
}

const rows = db.query(`SELECT * FROM ${quoteIdent(req.params.name)} LIMIT 100`) || [];
```

## `require("express")`

The `express` module provides a small Express-style route registration API for `db-browser serve`.

```js
const express = require("express");
const app = express.app();
```

### `express.app()`

Returns the application route registry object. Call it once per script set and register routes on the returned object.

```js
const app = express.app();
```

### `app.get(pattern, handler)` and HTTP method variants

Registers a route handler.

```js
app.get("/customers/:id", (req, res) => {
  const customer = db.query("SELECT * FROM customers WHERE id = ?", req.params.id)[0];
  if (!customer) return res.status(404).send("not found");
  res.json(customer);
});
```

Available methods:

- `app.get(pattern, handler)`
- `app.post(pattern, handler)`
- `app.put(pattern, handler)`
- `app.patch(pattern, handler)`
- `app.delete(pattern, handler)`
- `app.all(pattern, handler)`

Patterns can include named params such as `/customers/:id`. The params appear in `req.params`.

### `app.static(prefix, directory)`

Serves files from a local directory under a URL prefix.

```js
app.static("/assets", "public");
```

- **prefix**: URL path prefix.
- **directory**: local filesystem directory.

Keep static directories narrow and avoid exposing project roots.

## Request object

Route handlers receive a plain request object.

| Property | Type | Description |
| --- | --- | --- |
| `req.method` | `string` | HTTP method. |
| `req.url` | `string` | Full request URL path and query. |
| `req.path` | `string` | URL path without query. |
| `req.query` | `Record<string, string | string[]>` | Query string values. Repeated keys become arrays. |
| `req.params` | `Record<string, string>` | Named route params. |
| `req.headers` | `Record<string, string>` | Request headers joined with `, `. |
| `req.cookies` | `Record<string, string>` | Parsed cookies. |
| `req.session` | `{ id, isNew, cookieName } | null` | Lightweight session identity. |
| `req.ip` | `string` | Remote IP address when available. |
| `req.body` | `any` | Parsed body for JSON/forms, raw string for other content types. |
| `req.rawBody` | `string` | Raw request body string. |

### Query values

```js
app.get("/", (req, res) => {
  const q = String(req.query.q || "");
  res.send("search: " + q);
});
```

### JSON bodies

For `Content-Type: application/json`, `req.body` is parsed JSON:

```js
app.post("/api/customers", (req, res) => {
  const name = String(req.body.name || "");
  res.json({ name });
});
```

### Form bodies

For `application/x-www-form-urlencoded` or `multipart/form-data`, `req.body` is an object whose values are strings or arrays of strings.

```js
app.post("/customers/:id/segment", (req, res) => {
  const segment = String(req.body.segment || "");
  db.exec("UPDATE customers SET segment = ? WHERE id = ?", segment, req.params.id);
  res.redirect("/customers/" + req.params.id);
});
```

## Response object

The response object controls status, headers, and body.

### `res.status(code)`

Sets the HTTP status and returns the response object for chaining.

```js
return res.status(404).send("not found");
```

### `res.set(name, value)`

Sets a response header and returns the response object.

```js
res.set("X-App", "db-browser").json({ ok: true });
```

### `res.type(value)`

Sets `Content-Type` and returns the response object.

```js
res.type("text/plain; charset=utf-8").send("ok");
```

### `res.json(value)`

Sends JSON with `Content-Type: application/json`.

```js
res.json({ ok: true, rows });
```

### `res.send(value)`

Sends strings as text or HTML depending on whether the trimmed string starts with `<`. Non-string values are sent as JSON.

```js
res.send("plain text");
res.send("<h1>Hello</h1>");
res.send({ ok: true });
```

### `res.html(nodeOrDocument)`

Renders a `ui.dsl` node or document and sends HTML.

```js
res.html(ui.page({ title: "Home" }, ui.h1("Home")));
```

### `res.redirect(url)` and `res.redirect(status, url)`

Redirects the browser. One argument uses status `302`.

```js
res.redirect("/");
res.redirect(303, "/customers/1");
```

### `res.end()`

Sends headers/status with no body.

```js
app.get("/favicon.ico", (req, res) => res.status(204).end());
```

## `require("ui.dsl")` / `require("ui")`

The UI DSL builds HTML nodes. It escapes normal text automatically and renders a complete document with `ui.page`.

```js
const ui = require("ui.dsl");
```

### Tag functions

Most common tags are exported as functions:

```js
ui.div({ class: "panel" }, ui.h1("Title"), ui.p("Body"));
ui.a({ href: "/customers/1" }, "Open customer");
ui.input({ type: "search", name: "q", value: req.query.q || "" });
```

The first argument can be an attributes object. Remaining arguments are children. If the first argument is text, a node, an array, a number, or a boolean, it is treated as a child rather than attributes.

Current tag helpers include:

```text
html head body title meta link script style main img br hr time svg path rect line polyline circle div span h1 h2 h3 h4 p a form input button select option ul ol li table thead tbody tr th td section article header footer nav label textarea strong em small pre code
```

### `ui.page(options, ...children)`

Creates an HTML document.

```js
ui.page(
  { title: "Customers" },
  ui.style(ui.raw("body { font-family: system-ui; }")),
  ui.h1("Customers")
);
```

- `options.title`: document title.
- Head tags such as `style`, `meta`, `link`, and `title` are placed in `<head>`.
- Other nodes are placed in `<body>`.

### `ui.fragment(...children)`

Creates a fragment that renders children without a wrapper.

```js
ui.fragment(ui.strong("A"), " ", ui.em("B"));
```

### `ui.text(value)`

Creates an escaped text node.

```js
ui.text("<safe because escaped>");
```

### `ui.raw(html)`

Creates raw, unescaped HTML.

```js
ui.style(ui.raw("body { margin: 0; }"));
```

Only use this for trusted HTML or CSS that your script owns. Do not put database values, request values, or user input into `ui.raw`.

### `ui.render(value)`

Renders any UI node/document value to an HTML string.

```js
const html = ui.render(ui.p("hello"));
```

Most apps should use `res.html(...)` instead.

## Inspection/debug components

`ui.dsl` includes reusable server-rendered components for database inspection and debug pages. These components are safe for untrusted database/request text by default because they render values as escaped text nodes. Do not use `ui.raw` for SQL, JSON, scripts, request parameters, or other untrusted content.

### `ui.codeBlock(language, source, options?)`

Renders escaped code or preformatted text with stable classes and lightweight server-side token highlighting for SQL, JSON, and JavaScript.

```js
ui.codeBlock("sql", row.sql, {
  title: row.name,
  lineNumbers: true,
  wrap: false,
  copy: true,
  maxHeight: "480px",
});
```

Options:

| Option | Type | Default | Description |
| --- | --- | --- | --- |
| `title` | `string` | `""` | Optional caption title. |
| `lineNumbers` | `boolean` | `false` | Adds `ui-codeblock--line-numbers`. Initial implementation is class/affordance only. |
| `wrap` | `boolean` | `true` | Adds wrap/nowrap class. |
| `copy` | `boolean` | `false` | Renders an inert `Copy` button affordance. |
| `maxHeight` | `string` | `""` | Adds `max-height`/`overflow:auto` style to the `<pre>`. |
| `class` | `string` | `""` | Additional class on the outer element. |

Render shapes:

```html
<pre class="ui-codeblock ui-codeblock--sql ui-codeblock--wrap"><code class="language-sql">escaped source</code></pre>
```

or, when `title` or `copy` is set:

```html
<figure class="ui-codeblock ui-codeblock--sql ui-codeblock--nowrap ui-codeblock--line-numbers">
  <figcaption class="ui-codeblock__caption">
    <span class="ui-codeblock__title">customers</span>
    <button class="ui-codeblock__copy" type="button">Copy</button>
  </figcaption>
  <pre class="ui-codeblock__pre"><code class="language-sql">escaped source</code></pre>
</figure>
```

The language is normalized to a CSS-safe token. Empty/invalid language becomes `text`. Highlighted tokens render as escaped `<span>` nodes with classes such as `ui-codeblock__token--keyword`, `ui-codeblock__token--string`, `ui-codeblock__token--number`, `ui-codeblock__token--comment`, `ui-codeblock__token--key`, and `ui-codeblock__token--literal`.

### Convenience code block aliases

```js
ui.sql(source, options?)       // ui.codeBlock("sql", source, options)
ui.js(source, options?)        // ui.codeBlock("javascript", source, options)
ui.jsonBlock(value, options?)  // pretty JSON code block
```

`ui.jsonBlock` pretty-prints objects and arrays. If passed a valid JSON string, it parses and re-indents it. If the string is not valid JSON, it falls back to escaped plain text.

```js
ui.jsonBlock({ schema, columns, sampleRows }, {
  title: "debug payload",
  lineNumbers: true,
});
```

### `ui.badge(value, options?)`

Renders compact status/type labels outside tables.

```js
ui.badge(row.type, { tone: "info", title: "SQLite schema type" })
ui.badge(row.ok ? "ok" : "error", { tone: row.ok ? "success" : "danger" })
```

Options:

| Option | Type | Description |
| --- | --- | --- |
| `tone` | `default | info | success | warning | danger | muted` | Unknown tones fall back to `default`. |
| `title` | `string` | Optional title attribute. |
| `class` | `string` | Additional CSS class. |

Render shape:

```html
<span class="ui-badge ui-badge--success ui-badge--value-yes" title="...">yes</span>
```

The text is escaped and the value class is normalized.

### `ui.tabs(id, tabs, options?)`

Renders a no-JavaScript multi-view detail component. The initial implementation uses CSS-friendly radio tab markup and also marks the selected panel with `ui-tabs__panel--active` for server-rendered themes.

```js
ui.tabs("record-tabs", [
  { id: "summary", label: "Summary", content: ui.table.fromRows("summary", rows).render({ query: req.query }) },
  { id: "json", label: "Raw JSON", content: ui.jsonBlock(row.raw_json, { lineNumbers: true }) },
  { id: "sql", label: "Schema SQL", content: ui.sql(row.sql) },
], { selected: req.query.tab || "summary" })
```

Tab spec:

| Field | Type | Description |
| --- | --- | --- |
| `id` | `string` | Optional tab id. Normalized and suffixed if duplicated. |
| `label` | `string` | Escaped tab label. |
| `content` | `UiNode | UiNode[] | string` | Content normalized through the normal UI node path. |
| `disabled` | `boolean` | Disabled tabs render disabled labels and cannot be selected. |

Options:

| Option | Type | Description |
| --- | --- | --- |
| `selected` | `string | number` | Selected tab id or zero-based index. Invalid selection falls back to the first non-disabled tab. |
| `class` | `string` | Additional class on the container. |

Minimum classes:

```text
ui-tabs
ui-tabs__radio
ui-tabs__tablist
ui-tabs__tab
ui-tabs__tab--disabled
ui-tabs__panels
ui-tabs__panel
ui-tabs__panel--active
```

Example schema/debug detail page:

```js
res.html(ui.page({ title: "Schema" },
  ui.h1(row.name),
  ui.badge(row.type),
  ui.tabs("schema-tabs", [
    { id: "sql", label: "SQL", content: ui.sql(row.sql, { title: row.name, lineNumbers: true, copy: true }) },
    { id: "json", label: "Debug JSON", content: ui.jsonBlock(row, { lineNumbers: true }) },
  ])
));
```

## `ui.table` rich table API

The table API has two entry points.

### `ui.table.fromRows(id, rows)`

Builds a table from an in-memory array of objects. Static row tables apply filters, sorting, and pagination automatically.

```js
ui.table.fromRows("customers", rows)
  .columns(c => c
    .text("id").label("ID").sortable()
    .text("name").label("Customer").sortable().filterable().link(row => "/customers/" + row.id)
    .badge("segment").label("Segment").filterable()
    .money("total_cents").label("Total").align("right").sortable()
  )
  .features(f => f.filters().pagination({ size: 25 }).sorting().columnPicker())
  .render({ query: req.query, params: req.params })
```

### `ui.table(id)`

Builds a dynamic table. Use `.data(ctx => ...)` to provide rows. Dynamic tables receive context but must apply SQL filters, order, and pagination themselves.

```js
ui.table("orders")
  .data(ctx => {
    const rows = db.query("SELECT * FROM orders LIMIT ? OFFSET ?", ctx.page.limit, ctx.page.offset) || [];
    const total = db.query("SELECT COUNT(*) AS n FROM orders")[0].n;
    return { rows, total };
  })
  .columns(c => c.text("id").sortable().text("status").filterable())
  .features(f => f.filters().pagination({ size: 50 }).sorting())
  .render({ query: req.query, params: req.params })
```

### Builder methods

| Method | Purpose |
| --- | --- |
| `.data(fn)` | Set dynamic data callback. Callback receives render context. |
| `.columns(fn)` | Declare columns with a column builder. |
| `.features(fn)` | Enable filters, pagination, sorting, and column-picker markers. |
| `.render(ctx)` | Render to a UI node. Pass `{ query: req.query, params: req.params }`. |

### Data callback return values

A data callback may return an array:

```js
.data(ctx => db.query("SELECT * FROM orders") || [])
```

or an object with rows and total:

```js
.data(ctx => ({ rows, total }))
```

Use `{ rows, total }` when pagination is enabled and the full row count differs from the returned page size.

### Render context

Tables parse the render input into this context:

| Property | Description |
| --- | --- |
| `ctx.query` | Original query object. |
| `ctx.params` | Route params passed from `req.params`. |
| `ctx.page.index` | Current page number, minimum 1. |
| `ctx.page.size` | Page size. Defaults to feature size or 25. Maximum 500. |
| `ctx.page.limit` | Same as page size, for SQL `LIMIT`. |
| `ctx.page.offset` | SQL-style row offset. |
| `ctx.order.key` | Requested sort column from `?sort=...`. |
| `ctx.order.dir` | `asc` or `desc`. |
| `ctx.filter.q` | Global search string from `?q=...`. |
| `ctx.filter.<name>` | Per-column filters from `?filter.name=...` or `?filter_name=...`. |

### Feature builder

```js
.features(f => f.filters().pagination({ size: 25 }).sorting().columnPicker())
```

| Feature | Behavior |
| --- | --- |
| `f.filters()` | Renders a GET filter form. Static row tables filter automatically. |
| `f.pagination({ size })` | Renders pagination state and previous/next links. |
| `f.sorting()` | Renders sortable header links for sortable columns. |
| `f.columnPicker()` | Adds table marker class for future column-picker UI. |

### Column builder

Columns can be declared fluently:

```js
.columns(c => c
  .text("id").label("ID").sortable()
  .text("name").label("Customer").filterable().link(row => "/customers/" + row.id)
  .badge("status").label("Status").filterable()
  .money("total_cents").label("Total").align("right")
  .date("created_at").label("Created")
  .tags("labels").label("Labels")
)
```

Column kinds:

| Kind | Builder | Rendering |
| --- | --- | --- |
| Text | `c.text("field")` | Escaped text. |
| Badge | `c.badge("field")` | `<span class="ui-badge ui-badge--value">`. |
| Money | `c.money("field")` | Integer cents as dollars, e.g. `17998` -> `$179.98`. |
| Date | `c.date("field")` | Currently escaped text; reserved for richer formatting. |
| Tags | `c.tags("field")` | Comma/semicolon-separated values or arrays as tag spans. |

Column modifiers:

| Modifier | Behavior |
| --- | --- |
| `.label(text)` | Header label. |
| `.sortable()` | Enables sort header link. |
| `.filterable()` | Renders a per-column filter input when filters are enabled. |
| `.align("right")` | Adds `class="align-right"` to cells. |
| `.link(row => href)` | Wraps cell content in a dynamic anchor. Callback receives `(row, value)`. |
| `.link("/path/{field}")` | Wraps cell content in a template anchor with row-field substitution. |
| `.mono()` | Placeholder/no-op marker for future monospace styling. |
| `.truncate()` | Placeholder/no-op marker for future truncation styling. |

### Filter URL contract

The generated filter form uses GET parameters:

```text
?q=alice
?filter.segment=vip
?filter_name=alice
?sort=name&dir=asc&page=2
```

Static `table.fromRows` tables apply these in memory. Dynamic `.data(ctx => ...)` tables should use `ctx.filter`, `ctx.order`, and `ctx.page` in their own SQL.

### Empty states

When filtering produces no visible rows, the table renders one empty-state row:

```text
No rows match the current filters.
```

## `require("yaml")`

The YAML module parses and emits YAML. It is useful for dashboard specs and app manifests.

```js
const fs = require("fs");
const yaml = require("yaml");

const spec = yaml.parse(fs.readFileSync("dashboard.yaml", "utf8"));
```

Common operations:

```js
yaml.parse(text);
yaml.stringify(value);
```

Use YAML files for trusted local configuration. Do not allow untrusted users to submit arbitrary SQL-bearing YAML specs unless you have a separate authorization layer.

## `require("fs")`

The filesystem module follows the Goja NodeJS-style API exposed by go-go-goja.

Common operations:

```js
const fs = require("fs");

const text = fs.readFileSync("config.json", "utf8");
fs.writeFileSync("out.json", JSON.stringify({ ok: true }, null, 2));
const stat = fs.statSync("out.json");
```

The module accesses the local filesystem of the `db-browser` process. Treat paths as sensitive and avoid exposing arbitrary path reads through HTTP routes.

## `require("path")`

The path module provides path helpers.

```js
const path = require("path");
const configPath = path.join("examples", "yaml-dashboard", "dashboard.yaml");
```

Use it when building paths that should not rely on string concatenation.

## `require("time")` and `require("timer")`

These modules come from the default go-go-goja module set.

`timer` provides scheduling helpers such as `setTimeout`, `setInterval`, and `clearInterval`:

```js
const { setTimeout } = require("timer");

function later() {
  return new Promise(resolve => setTimeout(() => resolve("done"), 100));
}
```

`time` provides time utilities from go-go-goja. Prefer simple app-level date strings unless you need runtime time helpers.

## JS verb declarations

Repository-scanned verbs use go-go-goja `jsverbs` declarations. A typical file looks like:

```js
__package__({
  name: "reports",
  parents: ["examples"],
  short: "Report verbs",
});

function listCustomers() {
  const db = require("db");
  return db.query("SELECT id, name FROM customers ORDER BY id") || [];
}

__verb__("list-customers", {
  short: "List customers from the configured database",
});
```

Run it with a configured repository and database:

```bash
db-browser verbs --repository ./verbs --db app.db examples reports list-customers --output json
```

Important distinctions:

- `--scripts-dir` is for web app route scripts.
- `--repository` / config/env repository sources are for JS verbs.
- The two are intentionally separate inputs.
- Do not point `db-browser serve --scripts-dir` at a directory that also contains verb-only files using `__package__` or `__verb__`; the serve runtime does not define those symbols. Keep mixed projects in separate subdirectories such as `scripts/serve/` and `scripts/verbs/`, then run `db-browser serve --scripts-dir scripts/serve` and `db-browser verbs --repository scripts/verbs ...`.

## Error handling patterns

### Missing database module

If a script uses `require("db")`, always tell the user to pass `--db`.

```bash
db-browser serve --db data/app.db --scripts-dir scripts
db-browser verbs --db data/app.db examples builtin tables
```

### Route errors

Run with `--dev` while developing web apps. It shows detailed JavaScript errors in HTTP responses.

```bash
db-browser serve --db app.db --scripts-dir scripts --dev
```

### Safe generated apps

When prompting LLMs to generate apps, require these rules:

- Use only documented modules.
- Use `db.query(...) || []` for reads.
- Use SQL placeholders for values.
- Validate dynamic identifiers.
- Use `ui.dsl` nodes instead of concatenated HTML.
- Add `/favicon.ico` with status 204 for browser tests.
- Keep writes disabled unless explicitly requested.

## See also

- `db-browser help getting-started`
- `db-browser help user-guide`
- `db-browser help app-playbook`
- `db-browser inspect modules`
- `db-browser verbs list --output json`
