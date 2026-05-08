---
Title: Goja jsverbs database browser design and implementation guide
Ticket: DB-BROWSER-JSVERBS-DESIGN
Status: active
Topics:
    - goja
    - jsverbs
    - sqlite
    - web-ui
    - docmgr
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ../../../../../../../2026-05-03--goja-hosting-site/pkg/uidsl/module.go
      Note: Existing low-level ui.dsl module
    - Path: ../../../../../../../2026-05-03--goja-hosting-site/pkg/web/express_module.go
      Note: Express-style JavaScript route registration
    - Path: ../../../../../../../2026-05-03--goja-hosting-site/sites/trail/scripts/app.js
      Note: Example target JS authoring style
    - Path: ../../../../../../../corporate-headquarters/css-visual-diff/README.md
      Note: JavaScript-first host/domain split used to validate DB-browser design
    - Path: ../../../../../../../corporate-headquarters/css-visual-diff/cmd/css-visual-diff/main.go
      Note: Root command mounts lazy verbs command
    - Path: ../../../../../../../corporate-headquarters/css-visual-diff/examples/verbs/review-sweep.js
      Note: YAML-backed repository verb example
    - Path: ../../../../../../../corporate-headquarters/css-visual-diff/internal/cssvisualdiff/dsl/host.go
      Note: Runtime factory composition for scanned JS workflows
    - Path: ../../../../../../../corporate-headquarters/css-visual-diff/internal/cssvisualdiff/dsl/registrar.go
      Note: Custom domain module registrar pattern
    - Path: ../../../../../../../corporate-headquarters/css-visual-diff/internal/cssvisualdiff/jsapi/module.go
      Note: Promise-first native JS API pattern
    - Path: ../../../../../../../corporate-headquarters/css-visual-diff/internal/cssvisualdiff/verbcli/bootstrap.go
      Note: Repository discovery
    - Path: ../../../../../../../corporate-headquarters/css-visual-diff/internal/cssvisualdiff/verbcli/command.go
      Note: Lazy dynamic jsverbs command bootstrap pattern
    - Path: ../../../../../../../corporate-headquarters/go-go-goja/README.md
      Note: Runtime lifecycle and module security flags used by the design
    - Path: ../../../../../../../corporate-headquarters/go-go-goja/modules/database/database.go
      Note: Database native module query/exec API
    - Path: ../../../../../../../corporate-headquarters/go-go-goja/modules/yaml/yaml.go
      Note: yaml parse/stringify/validate module to expose in DB-browser
    - Path: ../../../../../../../corporate-headquarters/go-go-goja/pkg/jsverbs/command.go
      Note: Glazed command generation from scanned verbs
    - Path: ../../../../../../../corporate-headquarters/go-go-goja/pkg/jsverbs/runtime.go
      Note: Caller-owned runtime invocation and scanned-source require loader
    - Path: ../../../../../../../corporate-headquarters/go-go-goja/pkg/jsverbs/scan.go
      Note: jsverbs scanning and default command parent derivation
    - Path: ../../../../../../../corporate-headquarters/go-minitrace/cmd/go-minitrace/cmds/query/commands.go
      Note: Nested repository-backed command mounting pattern
    - Path: ../../../../../../../corporate-headquarters/go-minitrace/cmd/go-minitrace/cmds/query/js_runtime.go
      Note: Goja jsverbs runtime integration with a custom module
    - Path: ../../../../../../../corporate-headquarters/go-minitrace/pkg/minitracecmd/repositories.go
      Note: Repository path config/env/flag loading pattern
ExternalSources: []
Summary: Design for a Go-hosted Goja web app that scans repository jsverbs, exposes them as CLI verbs, and serves SQLite exploratory UIs through database, fs, express, and ui.dsl modules.
LastUpdated: 2026-05-07T20:25:00-04:00
WhatFor: Use this to onboard an intern before implementing the db-browser host and rich UI DSL.
WhenToUse: Read before touching go-go-goja jsverbs, go-minitrace repository scanning, or the goja-hosting-site Express/UI runtime.
---



# Goja jsverbs database browser design and implementation guide

## Executive summary

Build a new Go CLI/web application in this repository that uses `../corporate-headquarters/go-go-goja/` as the JavaScript runtime substrate. The app should scan a configured set of repositories for JavaScript verbs, register those verbs as nested CLI commands, and also provide a web-hosting runtime for database-focused JavaScript applications. The JavaScript authoring surface should feel like a small Node/Express sandbox:

```javascript
const db = require("database");
const fs = require("fs");
const express = require("express");
const ui = require("ui.dsl");

const app = express.app();
app.get("/", (req, res) => {
  const rows = db.query("SELECT name FROM sqlite_schema WHERE type = 'table'");
  res.html(ui.page({ title: "DB Browser" }, ui.table.fromRows("tables", rows).render({ query: req.query })));
});
```

The key product decision is that JavaScript authors should describe *what* data view they want: filters, data query, columns, actions, layouts, metrics, and charts. The Go host and `ui.dsl` should own the repetitive mechanics: parsing/clamping query parameters, preserving URL state, building pagination links, sorting headers, empty states, safe HTML rendering, consistent styling, CSV export, and form/action wiring.

This project can reuse three existing codebases:

- `go-go-goja` already provides native modules, an explicit `engine.NewBuilder().Build().NewRuntime(ctx)` runtime lifecycle, and `pkg/jsverbs` scanning/command generation.
- `go-minitrace` shows how to load command repositories from config, environment variables, and CLI flags, then mount file-backed commands into a nested Cobra command tree.
- `2026-05-03--goja-hosting-site` already proves the Express-style Goja web host, database module wiring, response helpers, and a basic `ui.dsl` HTML node renderer.

The implementation should therefore be an integration project more than a greenfield interpreter. The new work is primarily: multi-repository jsverbs discovery, a CLI command surface tailored to this app, a production-ready database browser `ui.dsl`, and safety boundaries around SQLite, filesystem access, static routes, and request-time JavaScript execution.

## Problem statement and scope

The original product goal is to create a small JavaScript interpreter/sandbox for building web UIs over SQLite databases, in the spirit of Datasette but scriptable through Goja. The host side is Go. The JavaScript side should have a minimal elegant API:

- `require("database")` or `require("db")` for `query(sql, ...params)` and controlled `exec(sql, ...params)`.
- `require("fs")` for file helpers when enabled.
- `require("express")` for `app.get`, `app.post`, `res.html`, `res.json`, and `res.redirect`.
- `require("ui.dsl")` for declarative UI primitives: page, forms, tables, filters, dashboards, metrics, charts, cards, grid, split, and nav.
- `jsverbs` discovery so scripts in configured repositories can become CLI verbs.

Out of scope for the first implementation:

- A browser-side SPA framework. Server-rendered HTML plus small progressive-enhancement scripts is enough for v1.
- Arbitrary untrusted multi-tenant execution. The existing runtime is a sandbox in the Goja sense, but native modules can access host resources. Treat scripts as trusted local app code unless a later ticket adds stronger isolation.
- A full charting library API. The DSL should expose simple line/bar specs and render them through a bundled lightweight client or static SVG.
- A complete Datasette clone. The goal is a small authoring runtime that makes custom SQLite explorers pleasant.

## Current-state architecture with evidence

### go-go-goja runtime and module system

`../corporate-headquarters/go-go-goja/README.md` describes the project as a place to wire Go native modules into Goja with Node-style `require()` and lists `modules/` as the place where native modules live. It also states the canonical runtime lifecycle: create an `engine.NewBuilder`, add modules/options, build a factory, create a runtime, and close it explicitly (`README.md:34-41`). That lifecycle should be copied for this app rather than introducing ad hoc runtime ownership.

Important evidence:

- `README.md:3-11` says the repo wires Go-implemented native modules into Goja and supports explicit runtime composition.
- `README.md:70-88` documents module security flags and module allow/deny behavior, which matters because this app must decide whether `fs`, `database`, and `express` are enabled by default.
- `README.md:141-163` documents `WithModuleRootsFromScript`, useful when loading scripts from nested repository folders.

Native modules follow a small interface. The database module in `modules/database/database.go` provides a model for the DB API:

- `DBModule` has options for `WithName`, `WithPreconfiguredDB`, `WithCloseFn`, and `WithConfigureEnabled` (`database.go:20-57`).
- The default module name is `database`, but the host can create an alias by setting the name (`database.go:70-89`).
- The module exports `configure`, `query`, `exec`, and `close` (`database.go:153-160`).
- `Query` returns `[]map[string]any`, which is the ideal JavaScript table-row shape for the DSL (`database.go:195-248`).

The fs module already exposes promise-based and synchronous file helpers and registers them under `fs` (`modules/fs/fs.go:20-24`, `fs.go:40-63`, `fs.go:82-180`). Use it as-is initially, but consider a later restricted filesystem wrapper if scripts should be limited to repository roots.

### jsverbs discovery and CLI command generation

`go-go-goja/pkg/jsverbs` is the core for discovering JavaScript functions and exposing them through Glazed/Cobra commands.

Key observed behavior:

- `ScanDir(root)` walks a real directory, skips `node_modules` and dot-directories, reads `.js` and `.cjs` files, and builds source inputs (`scan.go:17-74`, `scan.go:195-207`).
- `ScanFS(fsys, root)` performs the same scan over an `fs.FS`, which is valuable for repository catalog overlays or embedded examples (`scan.go:76-124`).
- Default scan options include public functions, `.js`/`.cjs`, and fail-on-error diagnostics (`model.go:57-69`).
- `finalizeVerbs` exposes explicit `__verb__` metadata first, then public non-underscore functions as verbs if `IncludePublicFunctions` is true (`scan.go:247-287`).
- Default CLI parents are inferred from package metadata, folders, and file stem (`scan.go:334-356`).
- `Registry.Commands()` produces Glazed commands; `CommandsWithInvoker` lets a host override execution (`command.go:41-59`).
- `InvokeInRuntime` lets a host use a caller-owned runtime rather than letting each command create and close a fresh runtime (`runtime.go:46-110`).
- `RequireLoader()` exposes the scanned source overlay as a Node require loader so relative requires can work (`runtime.go:40-44`, `runtime.go:169-175`).

The sample CLI in `cmd/jsverbs-example/main.go` demonstrates a minimal end-to-end command: choose a directory, call `jsverbs.ScanDir`, build commands, add them to a Cobra root through Glazed, and add a `list` command that emits discovered verbs (`main.go:51-121`). This is the fastest starting point for the new CLI surface.

### go-minitrace repository scanning pattern

`../corporate-headquarters/go-minitrace/` has the repository-loading pattern requested in the prompt. The important package is `pkg/minitracecmd/repositories.go`:

- Repositories can come from app config, environment, and repeated CLI flags (`repositories.go:15-24`, `repositories.go:111-127`).
- The env var is parsed with `filepath.SplitList`, so colon-separated Unix paths and platform-specific separators work (`repositories.go:129-136`).
- Paths are normalized, deduplicated, and relative config paths are resolved against the config file directory (`repositories.go:138-175`).
- `SourceRootsFromPaths` validates directories, wraps them in `os.DirFS`, marks them read-only, and appends embedded source roots (`repositories.go:185-201`).
- `ExtractRepositoryFlagValuesFromArgs` parses early `--query-repository` flags before command registration (`repositories.go:212-230`). This is important because dynamic commands must be known before Cobra parses the final command.

The query command tree shows how loaded catalog commands become nested CLI groups. `NewCommandsCommand` loads a catalog from configured repositories, adds a persistent repository flag, creates group commands from folders, detects collisions, and mounts child commands (`cmds/query/commands.go:12-72`, `commands.go:74-135`). Copy this shape for `db-browser verbs` or make it the root command behavior.

`go-minitrace` also demonstrates Goja + jsverbs runtime composition for a domain-specific module. `RunJSCommandIntoProcessor` scans the relevant source root, finds the verb, builds a runtime with `registry.RequireLoader()`, adds default modules plus a custom `minitrace` module, invokes the verb, and emits rows (`cmds/query/js_runtime.go:24-86`). The new app should follow the same pattern but replace the custom `minitrace` module with `database`, `db`, `express`, and `ui.dsl` registrars depending on execution mode.

### goja-hosting-site Express and UI runtime

`../2026-05-03--goja-hosting-site/` is the closest prototype for the web side.

The trail script demonstrates the intended JS authoring style:

- Common modules are imported at the top: `database`, `express`, `ui.dsl`, and a domain DSL (`sites/trail/scripts/app.js:1-4`).
- `express.app()` returns an app object, static assets can be mounted, and routes are registered in JS (`app.js:6-8`, `app.js:295-326`).
- DB helper functions remain ordinary JavaScript functions that call `db.query` and `db.exec` (`app.js:29-49`, `app.js:85-143`).
- UI is composed through chainable DSL builders and rendered from routes (`app.js:209-268`, `app.js:270-292`).

The Express module itself is small and reusable. `pkg/web/express_module.go` registers native module `express`, exports `app()`, and the app object supports `get`, `post`, `put`, `patch`, `delete`, `all`, and `static` (`express_module.go:17-52`).

The request/response wrapper supports exactly the desired route surface:

- Request DTO includes method, URL, path, query, params, headers, cookies, session, IP, body, and raw body (`request_response.go:16-72`).
- Response JS object exposes `status`, `set`, `type`, `json`, `send`, `html`, `redirect`, and `end` (`request_response.go:89-105`).
- `res.html` renders a UI node through the configured renderer and sets `text/html` (`request_response.go:140-150`).

The existing `ui.dsl` is currently a low-level HTML node builder:

- It registers `ui.dsl` and `ui` modules (`pkg/uidsl/module.go:11-19`).
- It exposes basic tags, `fragment`, `text`, `raw`, `render`, and `page` (`module.go:21-37`).
- Rendering escapes text and attributes and supports document/head/body rendering (`render.go:13-23`, `render.go:74-125`, `render.go:128-199`).

The new app should reuse the basic node model but extend the module with higher-level builders for data applications.

### css-visual-diff JavaScript playground validation

`~/code/wesen/corporate-headquarters/css-visual-diff` is an important additional validation target because it is a mature JavaScript-first Goja playground that already combines repository-scanned verbs, custom native modules, YAML-driven workflows, asynchronous browser work, and a review web UI. It is closer to the desired “scripts own the workflow; Go owns reliable primitives” model than a plain command wrapper.

Observed patterns that should influence this project:

- The README explicitly frames the tool as **JavaScript-first**: Go owns browser/screenshot/CSS/artifact primitives while JavaScript owns project workflows such as page specs, selectors, policies, accepted differences, Storybook URLs, routes, and handoff formats (`README.md:11-12`). This validates the same split proposed here for DB browsing: Go should own stable host primitives; JS should own table/query/dashboard workflows.
- `css-visual-diff verbs` is already a lazy dynamic command tree. The root command adds `verbcli.NewLazyCommand()` (`cmd/css-visual-diff/main.go:46-50`), and the lazy command disables normal flag parsing until it has discovered repositories and rebuilt the real command tree (`verbcli/command.go:16-35`). This is stronger evidence than the earlier `go-minitrace` pattern for apps whose dynamic commands need early `--repository` parsing.
- Repository discovery supports an embedded built-in repository, config files, env var `CSS_VISUAL_DIFF_VERB_REPOSITORIES`, `--repository`, and `--verb-repository` flags (`verbcli/bootstrap.go:20-27`, `bootstrap.go:67-107`). The DB browser should adopt this richer shape rather than only the go-minitrace `queryRepositories` model.
- Config schema is nested under `verbs.repositories[]` with `name`, `path`, and optional `enabled` (`verbcli/bootstrap.go:53-65`, `bootstrap.go:172-198`). This is a good candidate schema for DB browser script repositories:
  ```yaml
  verbs:
    repositories:
      - name: local-db-tools
        path: ./scripts
        enabled: true
  ```
- Repository scanning sets `IncludePublicFunctions = false` and therefore requires explicit `__verb__` metadata (`verbcli/bootstrap.go:288-307`). That is safer for rich application repositories than auto-exporting every public function. The DB browser should default to explicit verbs for configured repositories and reserve public-function scanning for examples or a development flag.
- Duplicate verb paths are rejected with provenance from both repositories (`verbcli/bootstrap.go:314-342`). The DB browser should preserve this behavior because multi-repository CLI trees otherwise become unpredictable.
- Runtime module roots include the repository root, repository `node_modules`, the parent directory, and parent `node_modules` (`verbcli/bootstrap.go:344-355`). This is practical for real project-local scripts and complements `jsverbs.RequireLoader()`.
- The runtime factory composes `registry.RequireLoader()`, `engine.DefaultRegistryModules()`, and a custom runtime registrar (`dsl/host.go:44-50`). That means `fs`, `path`, `yaml`, and other default `go-go-goja` modules are available unless the app applies stricter middleware.
- The custom runtime registrar installs domain modules such as `diff`, `report`, and `css-visual-diff` (`dsl/registrar.go:25-82`, `jsapi/module.go:16-42`). This supports the proposed DB-browser design where `database`, `express`, and `ui.dsl` are domain modules installed beside default modules.
- The JavaScript API is Promise-first for browser/page work, and uses runtime-owner posting to resolve promises back on the Goja runtime (`jsapi/module.go:78-98`). DB-browser v1 can keep `db.query` synchronous, but any future filesystem, background import, or streaming query feature should copy this owner-posting pattern.
- The example `review-sweep.js` reads YAML specs through `require("yaml")`, reads files through `fs`, invokes a domain module through `require("diff")`, writes JSON artifacts, and exposes explicit verbs through `__verb__` (`examples/verbs/review-sweep.js:277-285`, `review-sweep.js:340-365`, `review-sweep.js:432-438`). This confirms that `yaml` is a first-class practical module for repository workflows and should be included in the DB browser host API.

Design update from this validation: add `yaml` to the intended JS surface. The minimal import block for real repository scripts should be allowed to look like this:

```javascript
const db = require("database");
const fs = require("fs");
const yaml = require("yaml");
const express = require("express");
const ui = require("ui.dsl");
```

The `yaml` module already exists in `go-go-goja`: it exports `parse`, `stringify`, and `validate` and is registered into the default module registry (`modules/yaml/yaml.go:22-23`, `yaml.go:61-71`, `yaml.go:74-147`). For this project, YAML is useful for saved view definitions, table/browser app manifests, dashboard packs, repository app configs, and fixture specifications.

## Proposed architecture

### High-level component diagram

```text
+----------------------------+        +-----------------------------+
| CLI root (Cobra + Glazed)  |        | HTTP server                 |
|                            |        |                             |
| db-browser list-verbs      |        | net/http -> web.Host        |
| db-browser <js verb path>  |        | routes registered by JS     |
| db-browser serve           |        | res.html/json/redirect      |
+-------------+--------------+        +--------------+--------------+
              |                                      |
              v                                      v
+---------------------------------------------------------------+
| App runtime builder                                            |
| - go-go-goja engine.NewBuilder                                 |
| - registry.RequireLoader() for scanned JS sources              |
| - database/db native modules                                   |
| - fs/path/time/timer modules through middleware                |
| - express registrar                                            |
| - ui.dsl registrar                                             |
+--------------------------+------------------------------------+
                           |
                           v
+---------------------------------------------------------------+
| Repository discovery + jsverbs registry                        |
| - config/env/--verb-repository paths                           |
| - per-root jsverbs.ScanDir or combined ScanFS/ScanSources      |
| - collision checks and provenance                              |
| - Registry.CommandsWithInvoker for CLI                         |
+---------------------------------------------------------------+
```

### CLI shape

Recommended command tree:

```text
db-browser
  serve --db ./app.db --script ./scripts/app.js --verb-repository ../apps
  verbs
    list --verb-repository ../repo --output table
    <dynamic nested path> [flags from jsverbs metadata]
  inspect
    modules
    routes
    schema
```

The `css-visual-diff` validation changes the CLI recommendation slightly: implement `verbs` as a lazy command with `DisableFlagParsing: true`, parse only leading repository flags first, build the dynamic command tree, then execute the resolved command with the remaining args. This lets commands loaded from repositories define their own flags without conflicting with the bootstrap flags.

Alternative: mount dynamic verbs directly under the root. This is shorter, but riskier because repository scripts can collide with built-in commands. Prefer `db-browser verbs ...` for v1.

Repository configuration should mirror go-minitrace with app-specific names:

```yaml
# .db-browser.yml or app config
authoringRepositories:
  - ./scripts
  - ../team-db-tools/jsverbs
```

Environment and flags:

```bash
DB_BROWSER_VERB_REPOSITORIES=./scripts:../shared db-browser verbs list
db-browser verbs --repository ./scripts --verb-repository ../shared list
```

Use early flag extraction before building dynamic commands. Follow `css-visual-diff` more closely than `go-minitrace` here: allow both `--repository` and `--verb-repository` aliases, support an embedded built-in repository, and keep repository bootstrap flags before the dynamic verb path.

### Runtime modes

There are two distinct runtime modes. Keeping them separate avoids subtle lifecycle bugs.

#### CLI verb invocation mode

One command invocation creates one runtime, invokes one JS verb, emits Glazed rows or text, then closes the runtime.

Pseudocode:

```go
func invokeVerbCommand(ctx context.Context, registry *jsverbs.Registry, verb *jsverbs.VerbSpec, vals *values.Values, cfg RuntimeConfig, gp middlewares.Processor) error {
    db, closeDB := openSQLite(cfg.DBPath)
    defer closeDB()

    databaseModule := databasemod.New(
        databasemod.WithPreconfiguredDB(db),
        databasemod.WithConfigureEnabled(false),
    )
    dbAlias := databasemod.New(
        databasemod.WithName("db"),
        databasemod.WithPreconfiguredDB(db),
        databasemod.WithConfigureEnabled(false),
    )

    factory, err := engine.NewBuilder().
        WithRequireOptions(require.WithLoader(registry.RequireLoader())).
        WithModules(
            engine.NativeModuleSpec{ModuleID: "database", ModuleName: "database", Loader: databaseModule.Loader},
            engine.NativeModuleSpec{ModuleID: "db", ModuleName: "db", Loader: dbAlias.Loader},
        ).
        UseModuleMiddleware(engine.MiddlewareOnly("fs", "path", "time", "timer")).
        WithRuntimeModuleRegistrars(uidsl.NewRegistrar()).
        Build()
    if err != nil { return err }

    rt, err := factory.NewRuntime(ctx)
    if err != nil { return err }
    defer rt.Close(context.Background())

    result, err := registry.InvokeInRuntime(ctx, rt, verb, vals)
    if err != nil { return err }
    return emitJSResult(ctx, gp, result)
}
```

Notes:

- CLI verbs do not need `express` by default unless a verb wants to generate app code or inspect routes. Keep `express` disabled in CLI invocation unless requested.
- `ui.dsl` may still be useful for verbs that emit HTML snippets or static pages.
- For read-only query verbs, use a guarded DB wrapper. For mutating admin verbs, require an explicit `--allow-writes` flag.

#### Web serve mode

A long-lived runtime owns route registration and serves many HTTP requests. This follows `goja-hosting-site/pkg/app/server.go`.

Pseudocode:

```go
func NewServer(cfg Config) (*Server, error) {
    db := openSQLite(cfg.DBPath)
    guard := dbguard.New(db, cfg.DBPath)
    metered := dbguard.NewMeteredDB(db, guard)

    host := web.NewHost(web.HostOptions{Dev: cfg.Dev, Renderer: uidsl.RenderAny})

    databaseModule := databasemod.New(databasemod.WithPreconfiguredDB(metered), databasemod.WithConfigureEnabled(false))
    dbAlias := databasemod.New(databasemod.WithName("db"), databasemod.WithPreconfiguredDB(metered), databasemod.WithConfigureEnabled(false))

    registry := scanConfiguredRepositories(cfg.VerbRepositories)

    factory := engine.NewBuilder().
        WithRequireOptions(require.WithLoader(registry.RequireLoader())).
        WithModules(databaseSpec(databaseModule), databaseSpec(dbAlias)).
        UseModuleMiddleware(engine.MiddlewareOnly("fs", "path", "time", "timer")).
        WithRuntimeModuleRegistrars(web.NewExpressRegistrar(host), uidsl.NewRegistrar(), dbguard.NewRegistrar(guard)).
        Build()

    rt := factory.NewRuntime(context.Background())
    host.SetRuntime(rt.Owner)

    // Load selected app scripts. These scripts call express.app().get/post and register routes.
    for _, script := range cfg.AppScripts { rt.RunScript(script.Path, script.Source) }

    return &Server{db: db, runtime: rt, host: host}, nil
}
```

The important distinction: route registration happens during script load; request handling runs inside the same runtime owner through `web.Host.ServeHTTP`, which serializes calls through the runtime owner (`host.go:79-94`).

## JavaScript API reference

### `database` / `db`

Minimal API:

```typescript
type Row = Record<string, unknown>;

interface DatabaseModule {
  query(sql: string, ...params: unknown[]): Row[];
  exec(sql: string, ...params: unknown[]): { rowsAffected?: number; lastInsertId?: number };
}
```

Rules:

- In hosted apps, the DB module is preconfigured by Go; JS must not call `configure`.
- Prefer `db.query` for SELECT and PRAGMA reads.
- `db.exec` should be write-gated by host config. If the app is in read-only mode, writes should fail clearly.
- Always parameterize values. Identifier interpolation must go through a JS helper like `ident(name)` that validates table and column names.

### `fs`

Use the existing `go-go-goja/modules/fs` surface:

```typescript
fs.readFileSync(path: string, encoding?: string | object): string | Buffer;
fs.writeFileSync(path: string, data: string | Buffer | Uint8Array, encoding?: string | object): void;
fs.existsSync(path: string): boolean;
fs.readdirSync(path: string): string[];
fs.statSync(path: string): FileStats;
```

For v1, enable `fs` only for trusted local apps. For later hardening, wrap it with an allowlisted root and deny path traversal outside configured repositories/assets.

### `yaml`

Expose the existing `go-go-goja/modules/yaml` module for practical repository workflows. `css-visual-diff` validates this need: its review sweep verb reads YAML specs through `require("yaml")` and turns them into browser comparison runs.

```typescript
interface YAMLModule {
  parse(input: string): unknown;
  stringify(value: unknown, options?: { indent?: number }): string;
  validate(input: string): { valid: boolean; errors?: string[] };
}
```

Recommended DB-browser uses:

- Saved table/browser view files.
- Dashboard packs checked into project repositories.
- App manifests that declare database path, scripts, static mounts, and startup routes.
- Fixture and smoke-test specs.

### `express`

Use the `goja-hosting-site` API:

```typescript
interface ExpressModule {
  app(): App;
}

interface App {
  get(pattern: string, handler: Handler): void;
  post(pattern: string, handler: Handler): void;
  put(pattern: string, handler: Handler): void;
  patch(pattern: string, handler: Handler): void;
  delete(pattern: string, handler: Handler): void;
  all(pattern: string, handler: Handler): void;
  static(prefix: string, dir: string): void;
}

type Handler = (req: Request, res: Response) => unknown;
```

Request object:

```typescript
interface Request {
  method: string;
  url: string;
  path: string;
  query: Record<string, string | string[]>;
  params: Record<string, string>;
  headers: Record<string, string>;
  cookies: Record<string, string>;
  session: Record<string, unknown>;
  ip: string;
  body: unknown;
  rawBody: string;
}
```

Response object:

```typescript
interface Response {
  status(code: number): Response;
  set(name: string, value: string): Response;
  type(value: string): Response;
  json(value: unknown): void;
  send(value: unknown): void;
  html(value: unknown): void;
  redirect(url: string): void;
  redirect(status: number, url: string): void;
  end(): void;
}
```

### `ui.dsl`: low-level primitives

Keep the existing low-level node API:

```typescript
ui.page(opts, ...children)
ui.fragment(...children)
ui.raw(html)
ui.text(value)
ui.div(attrs?, ...children)
ui.form(attrs?, ...children)
ui.input(attrs?)
ui.select(attrs?, ...children)
ui.option(attrs?, ...children)
ui.button(attrsOrText?, text?)
```

These primitives are useful escape hatches and implementation building blocks.

### `ui.dsl`: high-level data-app primitives

Add a data table builder:

```typescript
ui.table(id: string)
  .title(textOrFn)
  .state(ctx => object)
  .filters(filters => filters
    .search(name, opts)
    .select(name, opts)
    .dateRange(name, opts))
  .toolbar(ctx => Node)
  .data(ctx => ({ rows, total }))
  .columns(columns => columns
    .text(name)
    .badge(name, labels?)
    .money(name, opts?)
    .date(name, opts?)
    .tags(name, opts?))
  .rowActions(actions => actions
    .link(label, hrefFn, opts?)
    .post(label, actionFn, opts?))
  .features(features => features
    .pagination(opts?)
    .sorting(opts?)
    .filters(opts?)
    .columnPicker(opts?)
    .savedViews(opts?)
    .csvExport(opts?))
  .render(ctx)
```

Table render context:

```typescript
interface TableContext {
  query: Record<string, unknown>;
  params: Record<string, string>;
  state: Record<string, unknown>;
  filter: Record<string, unknown>;
  page: { limit: number; offset: number; index: number; size: number };
  order: { key: string; dir: "asc" | "desc"; sql(defaultOrder: string): string };
}
```

The DSL should own these mechanics:

- Parse `?page`, `?limit`, `?sort`, and `?dir`.
- Clamp page sizes to allowed values.
- Preserve filters across pagination and sorting links.
- Render selected filter values.
- Validate sort keys against declared sortable columns.
- Provide an empty state when `rows.length === 0`.
- Render error panels in dev mode.

Add dashboard/layout primitives:

```typescript
ui.dashboard(opts, ...widgets)
ui.metric(label, value, opts?)
ui.chart.line({ data, x, y, series?, opts? })
ui.chart.bar({ data, x, y, opts? })
ui.card(title, ...children)
ui.grid(opts, ...children)
ui.split({ sidebar, main })
ui.nav(title, links)
ui.breadcrumbs(links)
ui.header(...children)
```

Add table shortcuts:

```typescript
ui.table.fromRows(id, rows)
  .features(f => f.pagination().sorting().columnPicker())
  .render({ query: req.query })
```

## Example authoring patterns

### Generic SQLite browser

```javascript
const browser = ui.table("table-browser")
  .state(ctx => {
    const table = safeTable(ctx.query.table || "");
    const cols = columns(table).map(c => c.name);
    return { table, cols, q: String(ctx.query.q || ""), sort: safeColumn(table, ctx.query.sort || cols[0]), dir: ctx.query.dir === "desc" ? "desc" : "asc" };
  })
  .toolbar(ctx => ui.form({ method: "get", class: "toolbar" }, tableSelect(ctx), ui.input({ name: "q", value: ctx.state.q }), ui.button("Search")))
  .data(ctx => selectRows(ctx.state, ctx.page))
  .columns(ctx => ctx.state.cols.map(name => ui.column(name).label(name).sortable().mono().truncate(80)))
  .features(f => f.pagination({ sizes: [25, 50, 100] }).sorting().columnPicker());
```

The host should make this concise by providing table paging/order helpers, but SQL identifier safety remains the script author's responsibility unless a later `db.schema` helper is added.

### Curated table explorer

```javascript
const orders = ui.table("orders")
  .filters(f => f.search("q", { label: "Search" }).select("status", { options: [["", "Any"], ["paid", "Paid"]] }))
  .data(ctx => queryOrders(ctx.filter, ctx.order, ctx.page))
  .columns(c => c.text("id").link(row => `/orders/${row.id}`).text("customer").badge("status").money("total_cents", { cents: true }))
  .rowActions(a => a.post("Mark shipped", row => `/orders/${row.id}/ship`, { visible: row.status === "paid", confirm: "Mark shipped?" }))
  .features(f => f.filters().pagination().sorting());
```

### Detail page with related rows and forms

```javascript
app.get("/customers/:id", (req, res) => {
  const c = customer(req.params.id);
  if (!c) return res.status(404).send("not found");
  res.html(ui.page({ title: c.name },
    ui.breadcrumbs([["Customers", "/customers"], [c.name, `/customers/${c.id}`]]),
    ui.dashboard({ columns: 3 }, ui.metric("Orders", stats.order_count), ui.metric("Lifetime", ui.money(stats.lifetime_cents, { cents: true }))),
    ui.grid({ columns: [2, 1] },
      ui.card("Orders", customerOrders.render({ query: req.query, params: { customerId: c.id } })),
      ui.card("Edit customer", editCustomerForm(c))
    )
  ));
});
```

### SQL workbench

A named SQL workbench should be a simple user-land pattern rather than a special host feature. The DSL only needs forms and `ui.table.fromRows`.

```javascript
const queries = {
  recent_signups: { title: "Recent signups", params: { days: { type: "number", default: 7 } }, sql: "SELECT * FROM customers WHERE created_at >= date('now', '-' || ? || ' days')", args: p => [Number(p.days || 7)] }
};

app.get("/workbench", (req, res) => {
  const q = queries[String(req.query.query || "recent_signups")] || queries.recent_signups;
  const rows = db.query(q.sql, ...q.args(req.query));
  res.html(ui.page({ title: "SQL Workbench" }, queryForm(queries, q, req.query), ui.table.fromRows("result", rows).features(f => f.pagination().sorting()).render({ query: req.query })));
});
```

## Implementation guide for a new intern

### Phase 1: Project skeleton and CLI

1. Create the Go module/command if it does not exist yet.
2. Add dependencies on `go-go-goja`, `glazed`, `cobra`, `sqlite3`, and any local replace directives needed for development.
3. Implement `cmd/db-browser/main.go` with a root Cobra command.
4. Add persistent flags:
   - `--db path`
   - `--script path` or `--scripts-dir path`
   - `--verb-repository path` repeatable
   - `--readonly` default true for browser mode
   - `--addr :8080`
   - `--dev`
5. Add a built-in `verbs list` command that scans configured repositories and prints: full path, source root, source file, function, output mode, short description.

Validation:

```bash
go test ./...
go run ./cmd/db-browser verbs list --verb-repository ./examples
```

### Phase 2: Repository discovery

Implement an app-local package, for example `pkg/verbrepos`, modeled on `go-minitrace/pkg/minitracecmd/repositories.go`.

Suggested types:

```go
type AppConfig struct {
    VerbRepositories []string `yaml:"verbRepositories"`
}

type SourceRoot struct {
    Name string
    FS fs.FS
    RootDir string
    Readonly bool
}
```

Required behavior:

- Load from app config, `.db-browser.yml`, `.db-browser.override.yml`, env var, and flags.
- Normalize and dedupe paths.
- Resolve relative paths in config relative to the config file.
- Ignore missing paths with a warning for interactive use, but fail in `--strict` mode.
- Preserve source root name for provenance and collision diagnostics.

### Phase 3: Combined jsverbs registry

Start simple: scan each repository with `jsverbs.ScanDir`, then merge command metadata at the CLI layer. If `jsverbs.Registry` does not expose safe merge operations, do not mutate private internals. Keep a slice of registries:

```go
type LoadedRegistry struct {
    Root SourceRoot
    Registry *jsverbs.Registry
}
```

For CLI registration, iterate all registries and all verbs. Detect duplicate `verb.FullPath()` before adding commands. For invocation, capture the `(registry, verb)` pair in the command wrapper.

Later, consider adding `jsverbs.ScanRepositories` upstream to `go-go-goja` if many apps need this.

### Phase 4: Runtime builder

Create `pkg/runtime/build.go` with functions for CLI and web modes. Centralize module decisions here.

Checklist:

- Register `database` and `db` aliases with the same preconfigured DB.
- Register `ui.dsl` high-level module.
- Register `yaml` unless the operator explicitly selects a strict module profile.
- Register `express` only in web mode.
- Use `registry.RequireLoader()` so scanned scripts can require each other.
- Add repository and parent `node_modules` folders like `css-visual-diff` does for external script roots.
- Use middleware allowlists for standard modules.
- Make DB writes configurable.
- Close runtime and DB handles reliably.

### Phase 5: Web host

Either vendor/copy the minimal `web` and `uidsl` packages from `goja-hosting-site` or depend on that module if it is intended to become reusable. Prefer copying only if the original is a prototype and not stable.

Implement:

- `Server` with `db`, `runtime`, `host`, and `http.Server` fields.
- `NewServer(cfg)` that opens DB, builds runtime, loads scripts, and returns an HTTP handler.
- `LoadScripts(ctx)` that loads configured app scripts in deterministic order.
- `/favicon.ico` default 204 if no route is registered.
- Optional `/__debug/routes` in dev mode.

### Phase 6: High-level `ui.dsl`

Build the high-level DSL in layers:

1. Keep existing `Node`, `Element`, `Fragment`, `Document`, and renderer.
2. Add Go-backed builder objects for table/columns/features, or implement builders in JavaScript and expose only low-level rendering helpers.
3. Prefer Go-backed builders if you want stronger validation and easier tests.
4. Prefer JS builders if you want faster iteration and a smaller Go surface.

Recommended v1 split:

- Go owns rendering primitives and escaping.
- Go owns table context parsing and safe URL generation.
- JavaScript builder methods capture callbacks (`data`, `columns`, `toolbar`) as Goja callables.

Core table render algorithm:

```go
func (t *Table) Render(vm *goja.Runtime, input RenderInput) (uidsl.Node, error) {
    ctx := t.BuildContext(input.Query, input.Params)
    state := callOptional(t.stateFn, ctx)
    ctx.State = state
    filters := parseFilters(t.filterSpecs, input.Query)
    ctx.Filter = filters
    columns := callColumns(t.columnsFn, ctx)
    ctx.Order = parseOrder(input.Query, columns)
    data := callData(t.dataFn, ctx) // { rows, total }
    return renderTable(t, ctx, columns, data)
}
```

Do not implement every feature at once. Suggested order:

1. `ui.table.fromRows`.
2. Manual `.data` + `.columns`.
3. Pagination.
4. Sorting.
5. Filters.
6. Column picker.
7. Row actions.
8. CSV export and saved views.
9. Charts.

### Phase 7: Example apps

Add examples as executable documentation:

- `examples/generic-browser/scripts/app.js`
- `examples/orders/scripts/app.js`
- `examples/documents/scripts/app.js`
- `examples/workbench/scripts/app.js`

Each example should have a small SQLite fixture and a README with commands:

```bash
db-browser serve --db examples/orders/orders.db --scripts-dir examples/orders/scripts --addr :8080
```

### Phase 8: Tests

Test at four levels:

- Repository config unit tests: env/flags/config normalization, dedupe, relative path resolution.
- jsverbs integration tests: scan fixture repos, detect collisions, run a fixture verb.
- UI rendering tests: render table HTML snapshots for pagination, sorting, filters, empty rows, escaping.
- HTTP integration tests: load a script, hit routes with `httptest`, verify HTML/JSON/redirect behavior.

Security tests:

- SQL read-only mode rejects `INSERT`, `UPDATE`, `DELETE`, `DROP`, and multi-statement writes.
- Sort keys are constrained to declared sortable columns.
- HTML text and attributes are escaped.
- Static mounts cannot escape configured roots.

## Design decisions

1. **Use go-go-goja engine lifecycle directly.** It is already documented and tested, and it avoids inventing runtime ownership.
2. **Reuse jsverbs instead of writing a new scanner.** It already handles top-level functions, metadata, command schemas, binding plans, relative require overlays, and invocation.
3. **Mirror go-minitrace repository configuration.** Dynamic command repositories are a solved problem there, including config layering, env vars, flags, normalization, and nested command mounting.
4. **Keep `express` server-rendered.** It matches the prototype and keeps the first version small. The DSL can add progressive client behavior later.
5. **Make the high-level DSL declarative.** Authors describe filters/data/columns/features; the host manages repeated mechanics.
6. **Prefer safety by construction.** The database module should be preconfigured by Go, JS should not configure arbitrary DSNs, and write access should be explicit.

## Alternatives considered

- **Build a custom JavaScript interpreter instead of Goja.** Rejected because Goja and the module system already exist and support the desired API.
- **Expose only low-level HTML tags.** Rejected because the prompt explicitly wants rich tables, dashboards, filtering, ordering, and forms without hand-rolling state each time.
- **Make everything a SPA.** Rejected for v1 because server-rendered Express routes are already proven and easier to inspect, test, and use from SQLite scripts.
- **Merge all repositories into one synthetic jsverbs registry immediately.** Deferred because `Registry` internals are private; a slice of registries with host-side collision detection is less invasive.
- **Mount dynamic verbs at CLI root.** Deferred because it increases collision risk with built-ins. A `verbs` namespace is clearer for v1.

## Risks and open questions

- **Runtime concurrency:** `web.Host` routes calls through a runtime owner, which likely serializes access. Long DB queries can block all requests. Consider per-request runtimes or a worker pool later.
- **Filesystem trust:** Existing `fs` is broad. If scripts are not fully trusted, implement a restricted fs module.
- **SQL writes:** Datasette-like browsing is mostly read-only, but forms need writes. Decide whether writes are app-wide, route-scoped, or command-scoped.
- **Chart rendering:** Decide between static SVG server rendering and a bundled browser chart library.
- **Registry caching:** Scanning repositories per command is simple but may be slow. Add cache invalidation only after measuring.
- **Type declarations:** `go-go-goja` has TypeScript declaration generation. Add declarations for `express` and high-level `ui.dsl` once APIs stabilize.

## References

- `../corporate-headquarters/go-go-goja/README.md:34-41` — canonical runtime lifecycle.
- `../corporate-headquarters/go-go-goja/modules/database/database.go:153-248` — database module exports and query rows.
- `../corporate-headquarters/go-go-goja/modules/fs/fs.go:20-63` — fs module API surface.
- `../corporate-headquarters/go-go-goja/modules/yaml/yaml.go:22-147` — yaml module API surface and default registration.
- `../corporate-headquarters/go-go-goja/pkg/jsverbs/scan.go:17-74` — directory scanning.
- `../corporate-headquarters/go-go-goja/pkg/jsverbs/scan.go:247-356` — verb finalization and default command parents.
- `../corporate-headquarters/go-go-goja/pkg/jsverbs/command.go:41-100` — Glazed command generation.
- `../corporate-headquarters/go-go-goja/pkg/jsverbs/runtime.go:40-110` — require loader and caller-owned runtime invocation.
- `../corporate-headquarters/go-go-goja/cmd/jsverbs-example/main.go:51-121` — minimal jsverbs CLI.
- `../corporate-headquarters/go-minitrace/pkg/minitracecmd/repositories.go:111-201` — repository collection and source roots.
- `../corporate-headquarters/go-minitrace/cmd/go-minitrace/cmds/query/commands.go:12-72` — nested dynamic command mounting.
- `../corporate-headquarters/go-minitrace/cmd/go-minitrace/cmds/query/js_runtime.go:24-86` — jsverbs runtime integration with a custom module.
- `../2026-05-03--goja-hosting-site/sites/trail/scripts/app.js:1-8` — target authoring style.
- `../2026-05-03--goja-hosting-site/pkg/web/express_module.go:17-52` — Express module registration.
- `../2026-05-03--goja-hosting-site/pkg/web/request_response.go:89-105` — response API.
- `../2026-05-03--goja-hosting-site/pkg/uidsl/module.go:11-37` — current low-level UI DSL.
- `../corporate-headquarters/css-visual-diff/README.md:11-12` — JavaScript-first host/domain split.
- `../corporate-headquarters/css-visual-diff/internal/cssvisualdiff/verbcli/command.go:16-35` — lazy dynamic verbs command with disabled initial flag parsing.
- `../corporate-headquarters/css-visual-diff/internal/cssvisualdiff/verbcli/bootstrap.go:20-107` — built-in/config/env/CLI verb repository bootstrap.
- `../corporate-headquarters/css-visual-diff/internal/cssvisualdiff/verbcli/bootstrap.go:288-355` — explicit-verb scanning, duplicate detection, and module root options.
- `../corporate-headquarters/css-visual-diff/internal/cssvisualdiff/dsl/host.go:44-50` — runtime factory with scanned-source loader, default modules, and custom registrar.
- `../corporate-headquarters/css-visual-diff/internal/cssvisualdiff/jsapi/module.go:16-98` — Promise-first native API and runtime-owner promise resolution.
- `../corporate-headquarters/css-visual-diff/examples/verbs/review-sweep.js:277-285` — YAML + fs usage in a repository-scanned verb.
