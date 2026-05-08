const db = require("db");
const express = require("express");
const ui = require("ui.dsl");

const app = express.app();
app.get("/favicon.ico", (req, res) => res.status(204).end());

const retroCSS = `
:root{--ink:#111;--paper:#f2efe5;--panel:#fffdf4;--line:#111;--shadow:#969288;--accent:#6f8f8a;--warn:#9b8378}*{box-sizing:border-box}body{margin:0;background:var(--paper);color:var(--ink);font-family:"Chicago","Geneva",system-ui,sans-serif;font-size:14px;background-image:linear-gradient(90deg,rgba(0,0,0,.04) 1px,transparent 1px),linear-gradient(rgba(0,0,0,.04) 1px,transparent 1px);background-size:10px 10px}.desktop{max-width:1180px;margin:22px auto;padding:0 16px}.window{background:var(--panel);border:2px solid var(--line);box-shadow:6px 6px 0 var(--shadow);margin-bottom:18px}.titlebar{border-bottom:2px solid var(--line);padding:5px 8px;font-weight:900;background:repeating-linear-gradient(0deg,#fffdf4 0 2px,#dcd7cb 2px 4px)}.content{padding:16px}.toolbar,.ui-table-filters{display:flex;flex-wrap:wrap;gap:8px;align-items:center;border:2px solid var(--line);background:#e7e2d7;padding:9px;margin:10px 0}.card-grid{display:grid;grid-template-columns:repeat(3,1fr);gap:10px}.card{border:2px solid var(--line);background:#fffdf7;padding:10px;border-left:10px solid var(--accent);box-shadow:3px 3px 0 #c8c2b6}.card strong{font-size:24px;display:block}input,button{font:inherit;border:2px solid var(--line);background:#fffdf7;padding:4px 6px;box-shadow:2px 2px 0 var(--shadow)}button{background:#dcd7cb;font-weight:900}table.ui-table{width:100%;border-collapse:separate;border-spacing:0;background:#fffdf7;border:2px solid var(--line)}th,td{padding:7px 8px;border-right:1px solid var(--line);border-bottom:1px solid var(--line);text-align:left}th{background:#dcd7cb}.align-right{text-align:right}.ui-badge,.ui-tag{display:inline-block;border:1px solid var(--line);padding:1px 5px;background:#d7e3df;font-size:12px;font-weight:800}.ui-table-pagination{margin-top:8px;display:inline-block;border:2px solid var(--line);background:#fffdf7;padding:5px 8px;box-shadow:3px 3px 0 var(--shadow)}a{color:#172f44;font-weight:800}.muted{color:#4d4a44}.danger{border-left-color:var(--warn)}
`;

function page(title, ...children) {
  return ui.page({ title }, ui.style(ui.raw(retroCSS)), ui.main({ class: "desktop" }, ui.section({ class: "window" }, ui.div({ class: "titlebar" }, title), ui.div({ class: "content" }, children))));
}

function quoteIdent(name) {
  if (!/^[A-Za-z_][A-Za-z0-9_]*$/.test(name)) throw new Error("unsafe table name: " + name);
  return '"' + name.replace(/"/g, '""') + '"';
}

function tableList() {
  const rows = db.query(`
    SELECT name
    FROM sqlite_schema
    WHERE type = 'table'
      AND name NOT LIKE 'sqlite_%'
    ORDER BY name
  `) || [];
  return rows.map(row => ({ name: row.name, column_count: columnsFor(row.name).length }));
}

function columnsFor(name) {
  return db.query(`PRAGMA table_info(${quoteIdent(name)})`) || [];
}

function rowsFor(name) {
  return db.query(`SELECT * FROM ${quoteIdent(name)} LIMIT 100`) || [];
}

app.get("/", (req, res) => {
  const tables = tableList();
  res.html(page("Generic SQLite Browser",
    ui.h1("Generic SQLite Browser"),
    ui.p({ class: "muted" }, "A monochrome, System-1-ish explorer for any SQLite file."),
    ui.div({ class: "card-grid" },
      ui.div({ class: "card" }, ui.small("Tables"), ui.strong(tables.length)),
      ui.div({ class: "card" }, ui.small("Database"), ui.strong("SQLite")),
      ui.div({ class: "card danger" }, ui.small("Mode"), ui.strong("Read-only UI"))
    ),
    ui.table.fromRows("tables", tables)
      .columns(c => c
        .text("name").label("Table").sortable().filterable()
        .text("column_count").label("Columns").align("right").sortable()
      )
      .features(f => f.filters().pagination({ size: 25 }).sorting().columnPicker())
      .render({ query: req.query })
  ));
});

app.get("/tables/:name", (req, res) => {
  const name = req.params.name;
  const cols = columnsFor(name);
  const rows = rowsFor(name);
  res.html(page("Table " + name,
    ui.p(ui.a({ href: "/" }, "← Tables")),
    ui.h1("Table: " + name),
    ui.h2("Columns"),
    ui.table.fromRows("columns", cols)
      .columns(c => c.text("cid").label("#").text("name").label("Name").filterable().text("type").label("Type").filterable().text("notnull").label("Required"))
      .features(f => f.filters().sorting())
      .render({ query: req.query }),
    ui.h2("Rows"),
    ui.table.fromRows("rows", rows)
      .features(f => f.filters().pagination({ size: 25 }).sorting())
      .render({ query: req.query })
  ));
});
