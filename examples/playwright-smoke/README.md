# Playwright Smoke DB App

A tiny seeded SQLite app for browser smoke testing `db-browser serve` with Playwright.

Seed database:

- `data/app.db`
- tables: `customers`, `orders`

Run manually:

```bash
go run ./cmd/db-browser serve \
  --db examples/playwright-smoke/data/app.db \
  --scripts-dir examples/playwright-smoke/scripts \
  --addr :19090 \
  --dev
```

Then open <http://127.0.0.1:19090/>.
