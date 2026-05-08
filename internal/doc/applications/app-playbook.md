---
Title: Writing full db-browser apps with an LLM
Slug: app-playbook
Short: A playbook for prompting an LLM to generate complete db-browser apps with SQLite, routes, tables, filters, and styling.
Topics:
- playbook
- llm
- app-generation
- sqlite
- ui-dsl
Commands:
- db-browser serve
- db-browser inspect modules
Flags:
- --db
- --scripts-dir
- --addr
- --dev
- --readonly
- --allow-writes
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: Tutorial
---

This playbook is written for humans who want to hand a clear app-building recipe to an LLM. It explains the app shape, the safe patterns, the expected files, and the checks that keep generated apps usable.

Use it when you want a complete db-browser app, not a one-off snippet. The output should include a seed database or seed script, JavaScript route files, styled UI pages, filters, detail views, and validation commands.

## The app contract to give the LLM

Ask the LLM to produce a directory with this shape:

```text
my-app/
  README.md
  data/
    seed.sql or seed.py
    app.db              # optional; include if binary fixtures are acceptable
  scripts/
    app.js
    00-theme.js         # optional if the app splits helpers
    10-routes.js        # optional if the app has many routes
```

The script files run inside `db-browser serve`, not in Node. The LLM must use CommonJS `require(...)` for host modules and must not assume a browser bundler, React, npm packages, or DOM APIs.

Available host modules:

```javascript
const db = require("db");
const express = require("express");
const fs = require("fs");
const path = require("path");
const yaml = require("yaml");
const ui = require("ui.dsl");
```

## Verb repository discovery contract

If the generated app also includes JavaScript CLI verbs, or if the LLM tells the user to inspect verbs, it must explain repository discovery explicitly.

By default, `db-browser verbs ...` scans only the embedded built-in repository:

```bash
db-browser verbs list --fields path,repository,source,location --output table
```

The default rows come from repository `builtin`, source `embedded`, location `builtin:builtin`. They are smoke/demo verbs, not the user's app scripts.

Additional verb repositories are discovered from:

- `.db-browser.yml`;
- `.db-browser.override.yml`;
- `DB_BROWSER_VERB_REPOSITORIES`;
- leading CLI flags `--repository` or `--verb-repository`.

A generated app that wants to ship its own verbs should include either a config file:

```yaml
verbs:
  repositories:
    - name: app-verbs
      path: ./verbs
```

or run instructions that pass the repository explicitly:

```bash
db-browser verbs --repository ./verbs list --output json
```

Do not imply that `scripts/` used by `db-browser serve` is scanned as a verb repository. Serve scripts and verb repositories are separate inputs.

## Master prompt template

Copy this prompt and adapt the domain details.

```text
You are writing a complete db-browser app.

Runtime constraints:
- The app runs with `db-browser serve --db data/app.db --scripts-dir scripts --dev`.
- JavaScript runs server-side in Goja, not Node and not a browser bundle.
- Use only these modules: db, express, fs, path, yaml, ui.dsl.
- Register routes with `const app = express.app()`.
- Return HTML with `res.html(ui.page(...))`.
- Use `db.query(sql, ...args) || []` for reads.
- Do not use string interpolation for SQL values; use placeholders.
- Do not require npm packages.
- Add `app.get("/favicon.ico", (req, res) => res.status(204).end())`.

Files to produce:
- README.md with run instructions.
- data/seed.py or data/seed.sql that creates and seeds the SQLite database.
- scripts/app.js containing the app.

Application domain:
<describe domain: customers/orders, tickets/comments, inventory/suppliers, experiments/results, etc.>

Pages required:
- `/`: overview dashboard with 3-5 metric cards and a primary table.
- A detail page for the primary entity, such as `/customers/:id`.
- At least one secondary table on the detail page.

UI requirements:
- Use ui.dsl only; do not concatenate HTML strings except trusted CSS inside `ui.raw`.
- Use ui.table.fromRows for in-memory rows.
- Use `.features(f => f.filters().pagination({ size: 25 }).sorting().columnPicker())` on primary tables.
- Mark useful columns `.filterable()` and `.sortable()`.
- Use `.link(row => "/entity/" + row.id)` for drill-down cells.
- Use badge columns for status/segment fields.
- Use money columns for integer cent values.

Style requirements:
- Include a self-contained monochrome retro Macintosh/System-1-inspired stylesheet.
- Use hard black borders, title bars, window panels, simple shadows, and muted accent colors.
- Keep the app readable without JavaScript in the browser.

Validation requirements:
- Provide commands to create the seed DB.
- Provide the `db-browser serve` command.
- Provide curl checks for the overview page, one filtered URL, and one detail page.
```

## Recommended generated app structure

A good generated `scripts/app.js` usually has five sections.

### 1. Imports and app setup

Start with explicit imports and a favicon route:

```javascript
const db = require("db");
const express = require("express");
const ui = require("ui.dsl");

const app = express.app();
app.get("/favicon.ico", (req, res) => res.status(204).end());
```

### 2. Theme and page shell

Keep the first app self-contained. Inline CSS is acceptable because `db-browser` does not yet require a static asset pipeline.

```javascript
const retroCSS = `
body { margin: 0; background: #f3efe4; color: #111; font-family: "Chicago", "Geneva", system-ui, sans-serif; }
.window { background: #fffdf4; border: 2px solid #111; box-shadow: 6px 6px 0 #989287; margin: 24px auto; max-width: 1100px; }
.titlebar { border-bottom: 2px solid #111; padding: 5px 8px; font-weight: 900; background: repeating-linear-gradient(0deg, #fffdf4 0 2px, #ddd8cc 2px 4px); }
.content { padding: 16px; }
`;

function page(title, ...children) {
  return ui.page(
    { title },
    ui.style(ui.raw(retroCSS)),
    ui.main({ class: "window" },
      ui.div({ class: "titlebar" }, title),
      ui.div({ class: "content" }, children)
    )
  );
}
```

Tell the LLM: use `ui.raw` only for trusted CSS, never for user or database values.

### 3. Data access helpers

Keep SQL in small functions. Use placeholders for values.

```javascript
function customers() {
  return db.query(`
    SELECT id, name, segment, email, total_cents
    FROM customers
    ORDER BY id
  `) || [];
}

function customerById(id) {
  return db.query("SELECT * FROM customers WHERE id = ?", id)[0];
}

function ordersForCustomer(id) {
  return db.query(`
    SELECT id, status, total_cents, created_at
    FROM orders
    WHERE customer_id = ?
    ORDER BY created_at DESC
  `, id) || [];
}
```

For dynamic SQL identifiers, such as table names in a generic browser, whitelist or validate them first. Never insert arbitrary query strings into SQL identifiers.

### 4. Overview route

The overview route should combine metrics and a primary table.

```javascript
app.get("/", (req, res) => {
  const rows = customers();
  const total = rows.reduce((sum, row) => sum + Number(row.total_cents || 0), 0);

  res.html(page("Customer Console",
    ui.h1("Customer Console"),
    ui.div({ class: "metric-grid" },
      ui.div({ class: "metric" }, ui.small("Customers"), ui.strong(rows.length)),
      ui.div({ class: "metric" }, ui.small("Revenue"), ui.strong(formatMoney(total)))
    ),
    ui.table.fromRows("customers", rows)
      .columns(c => c
        .text("id").label("ID").sortable()
        .text("name").label("Customer").sortable().filterable().link(row => "/customers/" + row.id)
        .badge("segment").label("Segment").filterable()
        .text("email").label("Email").filterable()
        .money("total_cents").label("Total").align("right").sortable()
      )
      .features(f => f.filters().pagination({ size: 25 }).sorting().columnPicker())
      .render({ query: req.query })
  ));
});
```

### 5. Detail route

The detail route should show one record plus related rows.

```javascript
app.get("/customers/:id", (req, res) => {
  const customer = customerById(req.params.id);
  if (!customer) return res.status(404).send("not found");

  const orders = ordersForCustomer(customer.id);
  res.html(page(customer.name,
    ui.p(ui.a({ href: "/" }, "← Customers")),
    ui.h1(customer.name),
    ui.p("Segment: " + customer.segment),
    ui.table.fromRows("orders", orders)
      .columns(c => c
        .date("created_at").label("Date").sortable()
        .badge("status").label("Status").filterable()
        .money("total_cents").label("Total").align("right").sortable()
      )
      .features(f => f.filters().sorting())
      .render({ query: req.query, params: req.params })
  ));
});
```

## App pattern catalog

Use these patterns as building blocks for LLM-generated apps.

### Pattern A: Read-only operations dashboard

Best for: support dashboards, incident views, data-quality boards, CI status boards.

Required pages:

- `/`: metric cards, recent items, status distribution table;
- `/items/:id`: item details, event log, related records.

Data model prompt:

```text
Create tables `items`, `events`, and `owners`. Include statuses open, investigating, resolved, and ignored. Use integer cents only if money is needed. Seed at least 20 items and 60 events.
```

Important UI choices:

- `badge("status")` for status;
- `date("created_at")` for timestamps;
- `tags("labels")` for labels;
- filters on status, owner, and severity.

### Pattern B: Generic SQLite inspector

Best for: exploring arbitrary SQLite files.

Required pages:

- `/`: table list with column counts;
- `/tables/:name`: schema table and first 100 rows.

Important safety rule:

```javascript
function quoteIdent(name) {
  if (!/^[A-Za-z_][A-Za-z0-9_]*$/.test(name)) throw new Error("unsafe identifier: " + name);
  return '"' + name.replace(/"/g, '""') + '"';
}
```

The LLM must validate table names before inserting them into SQL because SQLite placeholders cannot bind identifiers.

### Pattern C: YAML-configured metrics dashboard

Best for: handing non-programmers a small YAML file that defines cards and queries.

Files:

```text
scripts/app.js
dashboard.yaml
```

YAML shape:

```yaml
title: Support Dashboard
metrics:
  - label: Open Tickets
    group: Tickets
    tone: warning
    sql: |
      SELECT COUNT(*) AS value FROM tickets WHERE status != 'closed'
```

Script pattern:

```javascript
const fs = require("fs");
const yaml = require("yaml");
const spec = yaml.parse(fs.readFileSync("dashboard.yaml", "utf8"));
```

Keep YAML SQL read-only unless you fully trust the dashboard author.

### Pattern D: Small write-enabled admin tool

Best for: local-only maintenance tools where the operator understands the database.

Extra safety requirements:

- run only with `--readonly=false --allow-writes`;
- use POST routes for mutations;
- validate all required fields;
- redirect after writes;
- show a clear banner that writes are enabled.

Mutation route shape:

```javascript
app.post("/customers/:id/segment", (req, res) => {
  const segment = String(req.body.segment || "");
  if (!/^[a-z_ -]{1,32}$/i.test(segment)) return res.status(400).send("invalid segment");
  db.exec("UPDATE customers SET segment = ? WHERE id = ?", segment, req.params.id);
  res.redirect("/customers/" + req.params.id);
});
```

Do not generate write-enabled apps unless the user explicitly asks for writes.

## Validation checklist for generated apps

Ask the LLM to include and run these checks.

### Build or locate the binary

```bash
go build -o /tmp/db-browser ./cmd/db-browser
```

or, if installed:

```bash
db-browser inspect modules
```

### Create the database

```bash
python3 data/seed.py
```

or:

```bash
sqlite3 data/app.db < data/seed.sql
```

### Start the app

```bash
db-browser serve --db data/app.db --scripts-dir scripts --addr :18080 --dev
```

### Curl smoke tests

```bash
curl -fsS http://127.0.0.1:18080/ | grep -q '<title>'
curl -fsS 'http://127.0.0.1:18080/?q=alice' | grep -q 'Alice'
curl -fsS http://127.0.0.1:18080/customers/1 | grep -q 'Customer'
```

### Browser checks

When using Playwright or another browser driver, check:

- the page title;
- primary table row text;
- a filter interaction;
- a detail-page link;
- current console messages for errors and warnings.

## Common LLM mistakes to reject

| Mistake | Why it is wrong | Correction |
| --- | --- | --- |
| Uses `import` syntax or npm packages | The runtime is Goja with CommonJS-style host modules. | Use `const db = require("db")`; do not require external packages. |
| Generates React/Vite/browser code | db-browser serves server-rendered `ui.dsl` HTML. | Build pages with `ui.page`, tag helpers, and `ui.table`. |
| Concatenates database values into HTML strings | It risks broken markup and escaping bugs. | Use `ui.text`, normal tag children, and table cells. |
| Builds SQL values with template strings | It risks SQL injection and quoting bugs. | Use `?` placeholders and pass arguments to `db.query`. |
| Assumes filters automatically affect SQL-backed `.data(ctx => ...)` tables | Dynamic tables receive context but own their SQL. | Apply `ctx.filter`, `ctx.order`, and `ctx.page` inside the callback. |
| Omits run instructions | Reviewers cannot reproduce the app. | Include seed, serve, curl, and browser check commands. |
| Omits `/favicon.ico` | Browser tests get a noisy 404 console error. | Add a 204 favicon route. |

## Final handoff format

Ask the LLM to finish with:

1. File tree.
2. Full contents of each generated file.
3. Seed command.
4. Serve command.
5. Curl smoke checks.
6. Browser validation checklist.
7. Known limitations and safety notes.

This makes the generated app reviewable by both humans and coding agents.

## Troubleshooting

| Problem | Cause | Solution |
| --- | --- | --- |
| Generated app fails at startup with module errors | The LLM used unsupported modules or syntax. | Restrict the prompt to `db`, `express`, `fs`, `path`, `yaml`, and `ui.dsl`; use CommonJS `require`. |
| Detail links render but 404 | Route path and generated `link(...)` target do not match. | Compare `app.get("/entities/:id")` with `.link(row => "/entities/" + row.id)`. |
| Filters show inputs but SQL rows do not change | The app used `.data(ctx => ...)` without applying `ctx.filter`. | Apply filters in SQL or switch to `table.fromRows` for smaller data sets. |
| The app looks unstyled | The page omitted `ui.style(ui.raw(css))` or the CSS selector does not match generated markup. | Put a small theme in the page shell and inspect the rendered class names. |
| Browser tests show favicon errors | No favicon route. | Add the 204 route. |
| The app writes unexpectedly | The LLM generated mutation routes. | Keep read-only by default; require explicit user approval and serve flags for writes. |

## See also

- `db-browser help getting-started`
- `db-browser help user-guide`
- `db-browser help js-api-reference`
- `db-browser serve --help`
- `db-browser inspect modules`
