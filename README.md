# db-browser (retired)

`db-browser` has been retired as an independent shell.

Its reusable runtime pieces and examples have moved to the canonical Goja web
shells:

- `github.com/go-go-golems/go-go-goja/pkg/gojahttp`
- `github.com/go-go-golems/go-go-goja/modules/express`
- `github.com/go-go-golems/go-go-goja/modules/uidsl`
- `github.com/go-go-golems/go-go-goja/pkg/jsverbrepos`
- `github.com/go-go-golems/go-go-goja/pkg/jsverbscli`
- `2026-05-03--goja-hosting-site/cmd/goja-site`

Use `goja-site` instead:

```bash
# Serve a generic SQLite browser-style app in read-only mode.
go run ./cmd/goja-site serve \
  --db /path/to/app.sqlite \
  --scripts examples/db-browser/generic-browser/scripts \
  --db-policy simple \
  --readonly \
  --dev

# Run JavaScript verbs.
go run ./cmd/goja-site verbs list
```

Migrated examples now live in `goja-hosting-site`:

- `examples/db-browser/generic-browser`
- `examples/db-browser/yaml-dashboard`
- `examples/db-browser/playwright-smoke`
- `examples/verbs/builtin`

This repository intentionally no longer contains the duplicated server/runtime
implementation. Historical documentation remains under `docs/` and `ttmp/` for
reference only.
