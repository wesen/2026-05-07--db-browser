---
Title: Goja jsverbs database browser web app design
Ticket: DB-BROWSER-JSVERBS-DESIGN
Status: active
Topics:
    - goja
    - jsverbs
    - sqlite
    - web-ui
    - docmgr
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Ticket for designing a Goja-hosted SQLite browser/web app that scans repository jsverbs, exposes them as CLI verbs, and provides database/fs/express/ui.dsl APIs."
LastUpdated: 2026-05-07T20:25:00-04:00
WhatFor: "Central entrypoint for the db-browser jsverbs design work."
WhenToUse: "Use when planning or implementing the app described by DB-BROWSER-JSVERBS-DESIGN."
---

# Goja jsverbs database browser web app design

## Overview

This ticket contains the design package for a new Go web/CLI app that uses `../corporate-headquarters/go-go-goja/` to run JavaScript, scans configured repositories for `jsverbs`, exposes those verbs as CLI commands, and serves SQLite exploratory UIs through `database`, `fs`, `yaml`, `express`, and a richer `ui.dsl` module. The design was validated against `go-minitrace`, `goja-hosting-site`, and the richer `css-visual-diff` JavaScript playground.

## Key Links

- [Design and implementation guide](./design-doc/01-goja-jsverbs-database-browser-design-and-implementation-guide.md)
- [Investigation diary](./reference/01-investigation-diary.md)
- [Tasks](./tasks.md)
- [Changelog](./changelog.md)

## Status

Current status: **active**. The initial research/design deliverable is complete; implementation tasks remain open.

## Topics

- goja
- jsverbs
- sqlite
- web-ui
- docmgr

## Deliverables in this ticket

- Architecture analysis with evidence from `go-go-goja`, `go-minitrace`, `goja-hosting-site`, and `css-visual-diff`.
- Proposed CLI, runtime, module, and `ui.dsl` architecture.
- JavaScript API reference for `database`, `fs`, `express`, and `ui.dsl`.
- Intern-facing phased implementation guide.
- Diary of the investigation and documentation work.

## Structure

- `design-doc/` - Architecture and design documents.
- `reference/` - Investigation diary and reusable context.
- `playbooks/` - Future command sequences and test procedures.
- `scripts/` - Future temporary tooling.
- `various/` - Future working notes.
- `archive/` - Deprecated or reference-only artifacts.
