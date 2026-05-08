---
Title: db-browser user guide
Slug: user-guide
Short: A practical reference for running db-browser, writing scripts, using modules, and building tables.
Topics:
- user-guide
- sqlite
- goja
- express
- ui-dsl
Commands:
- db-browser serve
- db-browser inspect modules
- db-browser verbs
- db-browser verbs list
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

This guide explains how to use `db-browser` day to day. It covers the command model, JavaScript runtime, SQLite module, Express-style web host, `ui.dsl` rendering layer, table DSL, examples, and safety defaults.

Use this as the reference after you finish `getting-started`. When something is unclear in an app script, check the relevant section here and copy one of the patterns.

## Command overview

`db-browser` has three main surfaces:

| Command | Purpose |
| --- | --- |
| `db-browser serve` | Run a JavaScript web app against a SQLite database. |
| `db-browser inspect modules` | List the JavaScript modules exposed by the runtime. |
| `db-browser verbs ...` | Discover and run JavaScript verbs from configured repositories. |

The most common command is:

```bash
db-browser serve --db app.db --scripts-dir scripts --addr :8080 --dev
```

Scripts are loaded from `--scripts-dir` in sorted filename order. Keep route registration code deterministic and avoid hidden startup side effects.

## Runtime modules

The hosted JavaScript runtime includes modules for database access, files, YAML, routing, and HTML generation.

| Module | Typical import | Use |
| --- | --- | --- |
| SQLite database | `const db = require("db")` | Query the configured database. |
| Database alias | `const database = require("database")` | Same database module under a longer name. |
| Filesystem | `const fs = require("fs")` | Read local config/templates/data files. |
| Paths | `const path = require("path")` | Build portable filesystem paths. |
| YAML | `const yaml = require("yaml")` | Parse dashboard specs and app config files. |
| Express-style host | `const express = require("express")` | Register HTTP routes. |
| UI DSL | `const ui = require("ui.dsl")` | Build HTML pages, forms, tables, and fragments. |

Check the live list with:

```bash
db-browser inspect modules
```

## Database access

Use `db.query` for reads:

```javascript
const rows = db.query(
  "SELECT id, name FROM customers WHERE segment = ? ORDER BY name",
  "vip",
) || [];
```

Use placeholders instead of string concatenation for values. This keeps app scripts safe and predictable.

Use `db.exec` only when you intentionally run with writes enabled:

```javascript
db.exec("UPDATE customers SET segment = ? WHERE id = ?", "vip", id);
```

Served apps default to read-only mode. To permit writes you must pass both flags:

```bash
db-browser serve --db app.db --scripts-dir scripts --readonly=false --allow-writes
```

## Express-style routes

Create one app per script set:

```javascript
const express = require("express");
const app = express.app();
```

Register routes with handlers:

```javascript
app.get("/", (req, res) => {
  res.html(ui.page({ title: "Home" }, ui.h1("Home")));
});

app.get("/customers/:id", (req, res) => {
  const customer = db.query("SELECT * FROM customers WHERE id = ?", req.params.id)[0];
  if (!customer) return res.status(404).send("not found");
  res.json(customer);
});
```

The route handler receives:

- `req.query`: query string values;
- `req.params`: route path parameters;
- request body helpers for routes that accept POST data;
- `res.html(nodeOrDocument)`: render a `ui.dsl` node to HTML;
- `res.json(value)`: return JSON;
- `res.status(code)`: set the HTTP status;
- `res.redirect(url)`: redirect the browser.

Add a favicon route in browser-tested examples to avoid noisy console errors:

```javascript
app.get("/favicon.ico", (req, res) => res.status(204).end());
```

## UI DSL basics

`ui.dsl` creates structured HTML nodes rather than strings. This keeps escaping automatic for normal text.

```javascript
const ui = require("ui.dsl");

res.html(ui.page(
  { title: "Customers" },
  ui.h1("Customers"),
  ui.p("Server-rendered HTML from Goja."),
  ui.a({ href: "/customers/1" }, "Open Alice")
));
```

Most common HTML tags are functions: `ui.main`, `ui.section`, `ui.div`, `ui.form`, `ui.input`, `ui.table`, `ui.tr`, `ui.td`, and so on. The first argument can be an attributes object. Remaining arguments become children.

Use `ui.raw(html)` only for trusted HTML such as inline CSS that you wrote yourself:

```javascript
ui.style(ui.raw("body { font-family: system-ui; }"))
```

Do not put user-provided content into `ui.raw`.

## Table DSL

The table DSL is the highest-level UI primitive in db-browser.

Use `table.fromRows` when you already have rows:

```javascript
ui.table.fromRows("customers", rows)
  .columns(c => c
    .text("id").label("ID").sortable()
    .text("name").label("Customer").sortable().filterable()
    .badge("segment").label("Segment").filterable()
    .money("total_cents").label("Total").align("right").sortable()
  )
  .features(f => f.filters().pagination({ size: 25 }).sorting().columnPicker())
  .render({ query: req.query })
```

Use `table(id).data(...)` when rows come from a callback:

```javascript
ui.table("orders")
  .data(ctx => {
    const rows = db.query("SELECT * FROM orders LIMIT ? OFFSET ?", ctx.page.limit, ctx.page.offset) || [];
    const total = db.query("SELECT COUNT(*) AS n FROM orders")[0].n;
    return { rows, total };
  })
  .columns(c => c.text("id").sortable().text("status").filterable())
  .features(f => f.pagination({ size: 50 }).sorting().filters())
  .render({ query: req.query, params: req.params })
```

Static `fromRows` tables apply filters, sorting, and pagination automatically. Dynamic `data(ctx => ...)` tables receive `ctx` and must decide how to apply filters/sorting in SQL.

### Column kinds

| Kind | Builder | Typical data |
| --- | --- | --- |
| Text | `c.text("name")` | strings, ids, counts |
| Badge | `c.badge("status")` | enum/status values |
| Money | `c.money("total_cents")` | integer cents rendered as dollars |
| Date | `c.date("created_at")` | ISO date/time strings |
| Tags | `c.tags("tags")` | comma-separated or array-like tags |

### Column modifiers

| Modifier | Purpose |
| --- | --- |
| `.label("Customer")` | Human-readable column heading. |
| `.sortable()` | Enable sort links for the column. |
| `.filterable()` | Render a per-column filter input when filters are enabled. |
| `.align("right")` | Add alignment class, usually for numbers. |
| `.link(row => "/items/" + row.id)` | Wrap the cell in a dynamic link. |
| `.link("/tables/{name}")` | Wrap the cell in a template link. |

## Filters and query state

When filters are enabled, the table renders a GET form. It understands:

```text
?q=global search
?filter.segment=vip
?filter_name=alice
?sort=name&dir=asc&page=2
```

The render context passed to table callbacks contains:

```javascript
{
  query:  req.query,
  params: req.params,
  page:   { index, size, limit, offset },
  order:  { key, dir },
  filter: { q, segment, name }
}
```

Use this context when writing SQL-backed tables:

```javascript
.data(ctx => {
  const where = [];
  const args = [];
  if (ctx.filter.segment) {
    where.push("segment LIKE ?");
    args.push("%" + ctx.filter.segment + "%");
  }
  const sql = "SELECT * FROM customers" + (where.length ? " WHERE " + where.join(" AND ") : "") + " LIMIT ? OFFSET ?";
  return { rows: db.query(sql, ...args, ctx.page.limit, ctx.page.offset) || [], total: 0 };
})
```

For production-quality generated apps, compute an accurate `total` using the same `WHERE` conditions.

## JavaScript verb repositories

The `verbs` command scans configured JavaScript repositories and mounts explicit `__verb__` definitions as CLI commands.

List discovered verbs:

```bash
db-browser verbs list
```

Repository sources include:

- embedded built-in verbs;
- `.db-browser.yml`;
- `.db-browser.override.yml`;
- `DB_BROWSER_VERB_REPOSITORIES`;
- leading `--repository` or `--verb-repository` flags.

The runtime for CLI verbs includes the same useful host modules: filesystem, YAML, time/timer modules, `ui.dsl`, and optional database modules when `--db` is supplied.

## Safety and deployment notes

`db-browser` is currently best suited for local tools, internal exploration, generated prototypes, and reviewable experiments. Treat it as a server-side app host with filesystem and database access.

Before exposing an app beyond localhost, decide how you will handle:

- authentication;
- authorization;
- network binding;
- write access;
- input validation;
- database backups;
- secrets and filesystem paths.

Prefer read-only apps unless the use case truly requires writes.

## Troubleshooting

| Problem | Cause | Solution |
| --- | --- | --- |
| A route returns 500 | The JavaScript handler threw an exception. | Run with `--dev`, read the stack trace, fix the script, and restart. |
| Empty SQL results render as `null` in JSON | Some database paths can preserve nil Go slices. | In app scripts, use `db.query(...) || []` before rendering. |
| Filter form submits but rows do not change | The table uses `.data(ctx => ...)`, so app SQL must apply `ctx.filter`. | Either use `table.fromRows` or implement parameterized SQL filtering in the callback. |
| Sort links appear but order does not change | Dynamic tables must apply `ctx.order` themselves. | Map known sortable column names to safe SQL `ORDER BY` clauses. |
| Favicon 404 appears in Playwright console | The browser requests `/favicon.ico`. | Add `app.get("/favicon.ico", (req, res) => res.status(204).end());`. |
| `db.exec` is blocked | The server is in read-only mode. | Pass `--readonly=false --allow-writes` only for intentional write apps. |

## See also

- `db-browser help getting-started`
- `db-browser help app-playbook`
- `db-browser inspect modules`
- `db-browser verbs list`
