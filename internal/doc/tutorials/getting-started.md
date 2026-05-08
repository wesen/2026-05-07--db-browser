---
Title: Getting started with db-browser
Slug: getting-started
Short: Build and run your first Goja-backed SQLite browser app with db-browser.
Topics:
- getting-started
- sqlite
- goja
- web-apps
Commands:
- db-browser serve
- db-browser inspect modules
- db-browser verbs list
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

This tutorial gets you from an empty directory to a running browser UI backed by SQLite. It shows the smallest useful app shape: a seed database, one JavaScript script, one Express-style route, and one `ui.dsl` table.

The important idea is that `db-browser` does not ask you to build a frontend bundle. You write JavaScript that runs inside Goja, import host modules such as `db`, `express`, and `ui.dsl`, and return server-rendered HTML.

## What you will build

You will create a small customers app with:

- a SQLite database file;
- a `scripts/app.js` route file;
- a filterable, sortable HTML table;
- a read-only web server listening on `localhost`.

The same structure scales to dashboards, internal admin tools, CSV/database inspectors, and small analysis apps.

## 1. Check the runtime modules

Run this first so you know what JavaScript modules are available:

```bash
db-browser inspect modules
```

You should see modules such as:

```text
database
db
fs
yaml
express
ui.dsl
```

Use `db` for SQLite queries, `express` for routes, and `ui.dsl` for HTML nodes and higher-level widgets.

## 2. Create a tiny seed database

Create a working directory:

```bash
mkdir -p hello-db-browser/scripts hello-db-browser/data
cd hello-db-browser
```

Create a SQLite database with Python:

```bash
python3 - <<'PY'
import sqlite3
from pathlib import Path

path = Path('data/app.db')
con = sqlite3.connect(path)
con.execute('CREATE TABLE customers (id INTEGER PRIMARY KEY, name TEXT, segment TEXT, total_cents INTEGER)')
con.executemany(
    'INSERT INTO customers(name, segment, total_cents) VALUES (?, ?, ?)',
    [
        ('Alice Example', 'vip', 17998),
        ('Bob Browser', 'active', 2500),
        ('Carla Canvas', 'prospect', 0),
    ],
)
con.commit()
con.close()
print(path)
PY
```

## 3. Write the first app script

Create `scripts/app.js`:

```javascript
const db = require("db");
const express = require("express");
const ui = require("ui.dsl");

const app = express.app();

app.get("/favicon.ico", (req, res) => res.status(204).end());

function customers() {
  return db.query(`
    SELECT id, name, segment, total_cents
    FROM customers
    ORDER BY id
  `) || [];
}

app.get("/", (req, res) => {
  const rows = customers();
  res.html(ui.page(
    { title: "Customers" },
    ui.h1("Customers"),
    ui.p("A first db-browser app."),
    ui.table.fromRows("customers", rows)
      .columns(c => c
        .text("id").label("ID").sortable()
        .text("name").label("Customer").sortable().filterable()
        .badge("segment").label("Segment").filterable()
        .money("total_cents").label("Total").align("right").sortable()
      )
      .features(f => f.filters().pagination({ size: 25 }).sorting())
      .render({ query: req.query })
  ));
});
```

This script uses three host modules:

- `db.query(sql, ...args)` returns SQLite rows as JavaScript objects.
- `express.app()` registers routes for the embedded HTTP server.
- `ui.dsl` creates HTML nodes and renders tables.

## 4. Run the app

Start the server:

```bash
db-browser serve \
  --db data/app.db \
  --scripts-dir scripts \
  --addr :8080 \
  --dev
```

Open:

```text
http://127.0.0.1:8080/
```

Try query parameters that exercise the table features:

```text
http://127.0.0.1:8080/?q=alice
http://127.0.0.1:8080/?filter.segment=vip
http://127.0.0.1:8080/?sort=total_cents&dir=desc
```

## 5. Understand the execution model

`db-browser serve` loads all `.js` files from `--scripts-dir` in sorted order. Those scripts run once during startup and register routes. Route handlers run later when HTTP requests arrive.

The common route pattern is:

```javascript
app.get("/path", (req, res) => {
  const rows = db.query("SELECT ...") || [];
  res.html(ui.page({ title: "Page" }, ...nodes));
});
```

Request data is available through:

- `req.query` for query string parameters;
- `req.params` for route parameters such as `/customers/:id`;
- request body helpers in the Express-style host for POST routes.

## 6. Use read-only mode deliberately

By default, served scripts run with `--readonly=true`, which blocks write operations through `db.exec`. This is the safest mode for inspectors and dashboards.

If you are intentionally building a write-enabled internal tool, you must opt in explicitly:

```bash
db-browser serve \
  --db data/app.db \
  --scripts-dir scripts \
  --readonly=false \
  --allow-writes
```

Keep write-enabled apps small, auditable, and preferably local-only until you add authentication and deployment controls outside db-browser.

## 7. Next commands to try

List JavaScript verbs from configured repositories:

```bash
db-browser verbs list --fields path,repository,source,location --output table
```

With no config, environment variable, or CLI repository flag, `db-browser` scans only the embedded built-in repository. That repository is named `builtin`, has source `embedded`, and contains smoke/demo verbs such as:

```text
examples builtin hello
examples builtin yaml-keys
examples builtin tables
examples builtin render-sample-table
```

Add more verb repositories with one of these sources:

- `.db-browser.yml`;
- `.db-browser.override.yml`;
- `DB_BROWSER_VERB_REPOSITORIES`;
- leading CLI flags `--repository` or `--verb-repository`.

For example:

```bash
db-browser verbs --repository ./my-verbs list --output json
```

Run one of the included examples from the repository:

```bash
db-browser serve \
  --db examples/playwright-smoke/data/app.db \
  --scripts-dir examples/playwright-smoke/scripts \
  --addr :8080 \
  --dev
```

## Troubleshooting

| Problem | Cause | Solution |
| --- | --- | --- |
| `Cannot find module "db"` | The script is not running inside `db-browser serve`, or the runtime failed to initialize with a DB path. | Start the app with `db-browser serve --db path/to/app.db --scripts-dir scripts`. |
| Browser shows a JavaScript stack trace | `--dev` is enabled and a route handler threw an error. | Read the stack trace, fix the script, and restart the server. |
| Table filters do not affect rows | Dynamic `.data(ctx => ...)` tables must apply `ctx.filter` themselves. Static `table.fromRows(...)` tables filter automatically. | Use `table.fromRows` for in-memory rows, or read `ctx.filter` inside the SQL callback. |
| Writes fail even with valid SQL | Served apps default to read-only mode. | Use both `--readonly=false` and `--allow-writes`, and only do this for intentional write apps. |
| Route changes do not appear | Scripts are loaded at startup. | Stop and restart `db-browser serve`. |

## See also

- `db-browser help user-guide`
- `db-browser help app-playbook`
- `db-browser help`
