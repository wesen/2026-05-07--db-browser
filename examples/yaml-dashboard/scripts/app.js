const db = require("db");
const express = require("express");
const fs = require("fs");
const path = require("path");
const yaml = require("yaml");
const ui = require("ui.dsl");

const app = express.app();
const specPath = path.join("examples", "yaml-dashboard", "dashboard.yaml");

function loadSpec() {
  return yaml.parse(fs.readFileSync(specPath, "utf8"));
}

function metricRows(spec) {
  return (spec.metrics || []).map(metric => {
    const row = db.query(metric.sql)[0] || { value: null };
    return { metric: metric.label, value: row.value };
  });
}

app.get("/", (req, res) => {
  const spec = loadSpec();
  const rows = metricRows(spec);
  res.html(ui.page(
    { title: spec.title || "YAML Dashboard" },
    ui.h1(spec.title || "YAML Dashboard"),
    ui.table.fromRows("metrics", rows)
      .features(f => f.sorting())
      .render({ query: req.query })
  ));
});
