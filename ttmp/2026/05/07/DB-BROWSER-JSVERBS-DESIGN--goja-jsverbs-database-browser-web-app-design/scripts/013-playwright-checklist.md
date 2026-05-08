---
Title: Playwright smoke checklist
Ticket: DB-BROWSER-JSVERBS-DESIGN
Status: active
Topics:
  - goja
  - jsverbs
  - sqlite
  - web-ui
DocType: script
Intent: short-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Manual Playwright validation checklist for the seeded smoke app."
LastUpdated: 2026-05-07T21:10:00-04:00
WhatFor: "Use to reproduce the browser checks for examples/playwright-smoke."
WhenToUse: "Run after starting the Playwright smoke app server."
---

# Playwright smoke checklist

This checklist records the manual Playwright validation performed for `examples/playwright-smoke`.

## Server

```bash
go build -o /tmp/db-browser-playwright ./cmd/db-browser
/tmp/db-browser-playwright serve \
  --db examples/playwright-smoke/data/app.db \
  --scripts-dir examples/playwright-smoke/scripts \
  --addr :19090 \
  --dev
```

## Browser checks performed

- Navigated to `http://127.0.0.1:19090/`.
- Confirmed page title: `Playwright Smoke DB`.
- Confirmed main table shows seeded customers:
  - `Alice Example`
  - `Bob Browser`
  - `Carla Canvas`
- Clicked `Customer` sort header and confirmed URL changed to `/?dir=asc&page=1&sort=name`.
- Navigated to `http://127.0.0.1:19090/customers/1`.
- Confirmed page title: `Alice Example`.
- Confirmed detail table shows orders with statuses `shipped` and `paid`.
- Confirmed current console messages had zero errors/warnings after adding `/favicon.ico` route.

## Notes

The first navigation produced a 404 favicon console error. The app was updated with:

```js
app.get("/favicon.ico", (req, res) => res.status(204).end());
```

After server restart, the current Playwright console message check returned zero errors and warnings.
