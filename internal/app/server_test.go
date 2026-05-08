package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestServerLoadsScriptsAndServesHTMLAndJSON(t *testing.T) {
	dir := t.TempDir()
	scriptsDir := filepath.Join(dir, "scripts")
	if err := os.MkdirAll(scriptsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeScript(t, filepath.Join(scriptsDir, "app.js"), `
const express = require("express");
const ui = require("ui.dsl");
const db = require("db");
const app = express.app();
app.get("/", (req, res) => res.html(ui.page({ title: "Demo" }, ui.h1("Demo"))));
app.get("/tables", (req, res) => res.json(db.query("SELECT name FROM sqlite_schema WHERE type = 'table' AND name NOT LIKE 'sqlite_%' ORDER BY name")));
`)
	dbPath := filepath.Join(dir, "app.db")
	server, err := NewServer(context.Background(), Config{DBPath: dbPath, ScriptsDir: scriptsDir, Dev: true})
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close(context.Background())

	rr := httptest.NewRecorder()
	server.Handler().ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "<h1>Demo</h1>") {
		t.Fatalf("html body=%s", rr.Body.String())
	}

	rr = httptest.NewRecorder()
	server.Handler().ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/tables", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	if body := strings.TrimSpace(rr.Body.String()); body != "[]" && body != "null" {
		t.Fatalf("tables body=%s", rr.Body.String())
	}
}

func TestServerWriteGate(t *testing.T) {
	dir := t.TempDir()
	scriptsDir := filepath.Join(dir, "scripts")
	if err := os.MkdirAll(scriptsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeScript(t, filepath.Join(scriptsDir, "app.js"), `
const express = require("express");
const db = require("db");
const app = express.app();
app.post("/migrate", (req, res) => {
  db.exec("CREATE TABLE notes(id INTEGER PRIMARY KEY, body TEXT)");
  res.json({ ok: true });
});
`)
	server, err := NewServer(context.Background(), Config{DBPath: filepath.Join(dir, "app.db"), ScriptsDir: scriptsDir, Dev: true})
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close(context.Background())

	rr := httptest.NewRecorder()
	server.Handler().ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/migrate", nil))
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "database writes are disabled") {
		t.Fatalf("expected write gate error, got %s", rr.Body.String())
	}
}

func writeScript(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write script: %v", err)
	}
}
