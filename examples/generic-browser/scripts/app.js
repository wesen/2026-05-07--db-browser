const db = require("db");
const express = require("express");
const ui = require("ui.dsl");

const app = express.app();

function tables() {
  return db.query(`
    SELECT name
    FROM sqlite_schema
    WHERE type = 'table'
      AND name NOT LIKE 'sqlite_%'
    ORDER BY name
    LIMIT ? OFFSET ?
  `, 25, 0) || [];
}

app.get("/", (req, res) => {
  res.html(ui.page(
    { title: "Generic SQLite Browser" },
    ui.h1("Generic SQLite Browser"),
    ui.table("tables")
      .data(ctx => ({ rows: tables(), total: tables().length }))
      .columns(c => c
        .text("name").label("Table").sortable()
      )
      .features(f => f.pagination({ size: 25 }).sorting().columnPicker())
      .render({ query: req.query })
  ));
});
