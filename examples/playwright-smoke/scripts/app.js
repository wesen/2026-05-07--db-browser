const db = require("db");
const express = require("express");
const ui = require("ui.dsl");

const app = express.app();

app.get("/favicon.ico", (req, res) => res.status(204).end());

function customerRows() {
  return db.query(`
    SELECT
      customers.id,
      customers.name,
      customers.segment,
      customers.email,
      COUNT(orders.id) AS order_count,
      COALESCE(SUM(orders.total_cents), 0) AS total_cents
    FROM customers
    LEFT JOIN orders ON orders.customer_id = customers.id
    GROUP BY customers.id
    ORDER BY customers.id
  `) || [];
}

function orderRows(customerId) {
  return db.query(`
    SELECT id, status, total_cents, created_at
    FROM orders
    WHERE customer_id = ?
    ORDER BY created_at DESC
  `, customerId) || [];
}

app.get("/", (req, res) => {
  const rows = customerRows();
  res.html(ui.page(
    { title: "Playwright Smoke DB" },
    ui.main({ class: "page" },
      ui.h1("Playwright Smoke DB"),
      ui.p("A tiny seeded SQLite app rendered through Goja, Express, db, and ui.dsl."),
      ui.table("customers")
        .data(ctx => ({ rows, total: rows.length }))
        .columns(c => c
          .text("id").label("ID").sortable()
          .text("name").label("Customer").sortable()
          .badge("segment").label("Segment")
          .text("email").label("Email")
          .text("order_count").label("Orders").align("right")
          .money("total_cents").label("Total cents").align("right")
        )
        .features(f => f.pagination({ size: 10 }).sorting().columnPicker())
        .render({ query: req.query })
    )
  ));
});

app.get("/customers/:id", (req, res) => {
  const customer = db.query("SELECT * FROM customers WHERE id = ?", req.params.id)[0];
  if (!customer) return res.status(404).send("not found");
  const rows = orderRows(customer.id);
  res.html(ui.page(
    { title: customer.name },
    ui.h1(customer.name),
    ui.p("Segment: " + customer.segment),
    ui.table.fromRows("orders", rows)
      .features(f => f.sorting())
      .render({ query: req.query })
  ));
});
