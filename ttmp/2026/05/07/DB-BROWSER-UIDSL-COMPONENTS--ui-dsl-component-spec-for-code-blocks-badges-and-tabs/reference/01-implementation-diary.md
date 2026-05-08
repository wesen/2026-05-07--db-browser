---
Title: Implementation diary
Ticket: DB-BROWSER-UIDSL-COMPONENTS
Status: active
Topics:
  - db-browser
  - ui-dsl
  - goja
  - server-rendered-ui
  - documentation
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
  - path: /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-UIDSL-COMPONENTS--ui-dsl-component-spec-for-code-blocks-badges-and-tabs/design-doc/01-ui-dsl-codeblock-badge-tabs-implementation-guide.md
    note: Primary implementation guide for the component work.
  - path: /home/manuel/code/wesen/2026-05-07--db-browser/ttmp/2026/05/07/DB-BROWSER-UIDSL-COMPONENTS--ui-dsl-component-spec-for-code-blocks-badges-and-tabs/tasks.md
    note: Task checklist for step-by-step implementation.
ExternalSources: []
Summary: "Chronological implementation diary for ui.dsl codeBlock, badge, and tabs components."
LastUpdated: 2026-05-07T22:30:00-04:00
WhatFor: "Use to understand what changed, why, what failed, and how to continue/review the component implementation."
WhenToUse: "Read before resuming or reviewing DB-BROWSER-UIDSL-COMPONENTS work."
---

# Implementation diary

## Goal

Track the implementation of three reusable server-rendered `ui.dsl` components for db-browser inspection/debug pages: `ui.codeBlock`, `ui.badge`, and `ui.tabs`. The diary records planning, implementation steps, validation, failures, commits, and follow-up questions.

## Step 1: Create ticket and implementation guide

I created a focused docmgr ticket for the colleague-requested component work and wrote a detailed implementation guide before touching code. The guide translates the requested JavaScript APIs and render contracts into the current db-browser architecture: Goja module exports, `Node` composition, safe text rendering, table DSL precedents, tests, examples, and embedded help documentation.

I also added a step-by-step task list so the work can proceed in small commits: planning, code blocks, badges, tabs, examples, documentation, and final validation.

### Prompt Context

**User prompt (verbatim):** "Create a new ticket and implementation guide for (from our colleague): 


 ui.dsl component spec: codeBlock, badge, tabs                                                                            
                                                                                                                          
 Goal                                                                                                                     
                                                                                                                          
 Add three reusable server-rendered ui.dsl components for inspection/debug pages:                                         
                                                                                                                          
 - ui.codeBlock(...) for safe formatted code/pre blocks.                                                                  
 - ui.badge(...) for status/type labels.                                                                                  
 - ui.tabs(...) for compact multi-view record/detail pages.                                                               
                                                                                                                          
 These should work without client-side JavaScript by default and be safe for untrusted database/request text.             
                                                                                                                          
 ────────────────────────────────────────────────────────────────────────────────                                         
                                                                                                                          
 1. ui.codeBlock(language, source, options?)                                                                              
                                                                                                                          
 ### Purpose                                                                                                              
                                                                                                                          
 Render escaped code/text with optional title, line numbers, wrapping, copy affordance, and language class for future     
 syntax highlighting.                                                                                                     
                                                                                                                          
 ### JS API                                                                                                               
                                                                                                                          
 ```js                                                                                                                    
   ui.codeBlock(language, source, options?)                                                                               
 ```                                                                                                                      
                                                                                                                          
 Examples:                                                                                                                
                                                                                                                          
 ```js                                                                                                                    
   ui.codeBlock("sql", row.sql, {                                                                                         
     title: "CREATE VIEW delivery_chain",                                                                                 
     lineNumbers: true,                                                                                                   
     wrap: false,                                                                                                         
     copy: true,                                                                                                          
   });                                                                                                                    
                                                                                                                          
   ui.codeBlock("json", JSON.stringify(obj, null, 2));                                                                    
                                                                                                                          
   ui.codeBlock("javascript", scriptSource, {                                                                             
     title: "trace_browser_app.js",                                                                                       
   });                                                                                                                    
 ```                                                                                                                      
                                                                                                                          
 Convenience aliases, optional but desirable:                                                                             
                                                                                                                          
 ```js                                                                                                                    
   ui.sql(source, options?)       // ui.codeBlock("sql", source, options)                                                 
   ui.jsonBlock(value, options?)  // pretty-print object or JSON string, then codeBlock("json", ...)                      
   ui.js(source, options?)        // ui.codeBlock("javascript", source, options)                                          
 ```                                                                                                                      
                                                                                                                          
 ### Options                                                                                                              
                                                                                                                          
 ```ts                                                                                                                    
   type CodeBlockOptions = {                                                                                              
     title?: string;                                                                                                      
     lineNumbers?: boolean;                                                                                               
     wrap?: boolean;       // default true                                                                                
     copy?: boolean;       // render copy button/affordance; may be inert initially                                        
     maxHeight?: string;   // e.g. "480px"                                                                                
     class?: string;                                                                                                      
   };                                                                                                                     
 ```                                                                                                                      
                                                                                                                          
 ### Render contract                                                                                                      
                                                                                                                          
 ```html                                                                                                                  
   <figure class="ui-codeblock ui-codeblock--sql [custom-class]">                                                         
     <figcaption class="ui-codeblock__caption">                                                                           
       <span class="ui-codeblock__title">CREATE VIEW delivery_chain</span>                                                
       <button class="ui-codeblock__copy" type="button">Copy</button>                                                     
     </figcaption>                                                                                                        
     <pre class="ui-codeblock__pre"><code class="language-sql">escaped source</code></pre>                                
   </figure>                                                                                                              
 ```                                                                                                                      
                                                                                                                          
 If no title/copy:                                                                                                        
                                                                                                                          
 ```html                                                                                                                  
   <pre class="ui-codeblock ui-codeblock--sql"><code class="language-sql">escaped source</code></pre>                     
 ```                                                                                                                      
                                                                                                                          
 ### Safety requirements                                                                                                  
                                                                                                                          
 - Escape source by default.                                                                                              
 - Never require callers to pass source through ui.raw.                                                                   
 - language must be normalized to a CSS-safe token.                                                                        
 - Invalid/empty language should become "text".                                                                           
 - jsonBlock should:                                                                                                      
     - pretty-print objects;                                                                                              
     - parse and pretty-print valid JSON strings;                                                                         
     - fall back to escaped plain text for invalid JSON.                                                                  
                                                                                                                          
 ### Styling classes                                                                                                      
                                                                                                                          
 Minimum classes:                                                                                                         
                                                                                                                          
 ```text                                                                                                                  
   ui-codeblock                                                                                                           
   ui-codeblock--<language>                                                                                               
   ui-codeblock--wrap                                                                                                     
   ui-codeblock--nowrap                                                                                                   
   ui-codeblock--line-numbers                                                                                             
   ui-codeblock__caption                                                                                                  
   ui-codeblock__title                                                                                                    
   ui-codeblock__copy                                                                                                     
   ui-codeblock__pre                                                                                                      
 ```                                                                                                                      
                                                                                                                          
 ────────────────────────────────────────────────────────────────────────────────                                         
                                                                                                                          
 2. ui.badge(value, options?)                                                                                             
                                                                                                                          
 ### Purpose                                                                                                              
                                                                                                                          
 Render compact status/type labels outside tables. Useful for yes/no, ok/error, schema type, event stage, entity kind.    
                                                                                                                          
 ### JS API                                                                                                               
                                                                                                                          
 ```js                                                                                                                    
   ui.badge(value, options?)                                                                                              
 ```                                                                                                                      
                                                                                                                          
 Examples:                                                                                                                
                                                                                                                          
 ```js                                                                                                                    
   ui.badge("view");                                                                                                      
                                                                                                                          
   ui.badge(row.transport_fanout, {                                                                                       
     tone: row.transport_fanout === "yes" ? "success" : "danger",                                                         
   });                                                                                                                    
                                                                                                                          
   ui.badge("provider_normalize_delta", {                                                                                 
     tone: "info",                                                                                                        
     title: "Geppetto stage",                                                                                             
   });                                                                                                                    
 ```                                                                                                                      
                                                                                                                          
 ### Options                                                                                                              
                                                                                                                          
 ```ts                                                                                                                    
   type BadgeOptions = {                                                                                                  
     tone?: "default" | "info" | "success" | "warning" | "danger" | "muted";                                              
     title?: string;                                                                                                      
     class?: string;                                                                                                      
   };                                                                                                                     
 ```                                                                                                                      
                                                                                                                          
 ### Render contract                                                                                                      
                                                                                                                          
 ```html                                                                                                                  
   <span class="ui-badge ui-badge--success ui-badge--value-yes" title="...">yes</span>                                    
 ```                                                                                                                      
                                                                                                                          
 ### Safety requirements                                                                                                  
                                                                                                                          
 - Escape text.                                                                                                           
 - CSS token for value class must be normalized.                                                                          
 - Unknown tone falls back to default.                                                                                    
                                                                                                                          
 ### Styling classes                                                                                                      
                                                                                                                          
 ```text                                                                                                                  
   ui-badge                                                                                                               
   ui-badge--default                                                                                                      
   ui-badge--info                                                                                                         
   ui-badge--success                                                                                                      
   ui-badge--warning                                                                                                      
   ui-badge--danger                                                                                                       
   ui-badge--muted                                                                                                        
   ui-badge--value-<token>                                                                                                
 ```                                                                                                                      
                                                                                                                          
 ────────────────────────────────────────────────────────────────────────────────                                         
                                                                                                                          
 3. ui.tabs(id, tabs, options?)                                                                                           
                                                                                                                          
 ### Purpose                                                                                                              
                                                                                                                          
 Render a server-side tab/disclosure component for detail pages with multiple views, such as:                             
                                                                                                                          
 - Summary                                                                                                                
 - Pretty JSON                                                                                                            
 - Raw JSON                                                                                                               
 - SQL                                                                                                                    
 - Related rows                                                                                                           
                                                                                                                          
 Should work without custom JavaScript.                                                                                   
                                                                                                                          
 ### JS API                                                                                                               
                                                                                                                          
 ```js                                                                                                                    
   ui.tabs(id, tabs, options?)                                                                                            
 ```                                                                                                                      
                                                                                                                          
 Example:                                                                                                                 
                                                                                                                          
 ```js                                                                                                                    
   ui.tabs("record-tabs", [                                                                                               
     {                                                                                                                    
       id: "summary",                                                                                                     
       label: "Summary",                                                                                                  
       content: ui.kv(row),                                                                                               
     },                                                                                                                   
     {                                                                                                                    
       id: "json",                                                                                                        
       label: "Raw JSON",                                                                                                 
       content: ui.jsonBlock(row.raw_json, { lineNumbers: true }),                                                        
     },                                                                                                                   
     {                                                                                                                    
       id: "sql",                                                                                                         
       label: "Schema SQL",                                                                                                
       content: ui.sql(row.sql),                                                                                          
     },                                                                                                                   
   ], {                                                                                                                   
     selected: "summary",                                                                                                 
   });                                                                                                                    
 ```                                                                                                                      
                                                                                                                          
 ### Types                                                                                                                
                                                                                                                          
 ```ts                                                                                                                    
   type TabSpec = {                                                                                                       
     id?: string;                                                                                                         
     label: string;                                                                                                       
     content: UiNode | UiNode[] | string;                                                                                 
     disabled?: boolean;                                                                                                  
   };                                                                                                                     
                                                                                                                          
   type TabsOptions = {                                                                                                   
     selected?: string | number; // tab id or index                                                                       
     class?: string;                                                                                                      
   };                                                                                                                     
 ```                                                                                                                      
                                                                                                                          
 ### Render contract                                                                                                      
                                                                                                                          
 Prefer CSS-only radio tabs:                                                                                              
                                                                                                                          
 ```html                                                                                                                  
   <div class="ui-tabs" id="record-tabs">                                                                                 
     <input class="ui-tabs__radio" type="radio" name="record-tabs" id="record-tabs-summary" checked>                      
     <label class="ui-tabs__tab" for="record-tabs-summary">Summary</label>                                                
                                                                                                                          
     <input class="ui-tabs__radio" type="radio" name="record-tabs" id="record-tabs-json">                                 
     <label class="ui-tabs__tab" for="record-tabs-json">Raw JSON</label>                                                  
                                                                                                                          
     <div class="ui-tabs__panels">                                                                                        
       <section class="ui-tabs__panel" data-tab="summary">...</section>                                                   
       <section class="ui-tabs__panel" data-tab="json">...</section>                                                      
     </div>                                                                                                                
   </div>                                                                                                                 
 ```                                                                                                                      
                                                                                                                          
 Alternative acceptable initial implementation: render as stacked <details> blocks if CSS-only tabs are too much.         
                                                                                                                          
 ### Safety requirements                                                                                                  
                                                                                                                          
 - Escape labels.                                                                                                         
 - Normalize id and tab IDs to safe DOM tokens.                                                                            
 - If duplicate tab IDs are provided, suffix them.                                                                        
 - If selected tab is invalid, select first non-disabled tab.                                                             
 - Disabled tabs render label but cannot be selected.                                                                     
                                                                                                                          
 ### Styling classes                                                                                                      
                                                                                                                          
 ```text                                                                                                                  
   ui-tabs                                                                                                                
   ui-tabs__radio                                                                                                         
   ui-tabs__tablist                                                                                                       
   ui-tabs__tab                                                                                                           
   ui-tabs__tab--disabled                                                                                                 
   ui-tabs__panels                                                                                                        
   ui-tabs__panel                                                                                                         
   ui-tabs__panel--active                                                                                                 
 ```                                                                                                                      
                                                                                                                          
 ────────────────────────────────────────────────────────────────────────────────                                         
                                                                                                                          
 Documentation updates                                                                                                    
                                                                                                                          
 Update db-browser help js-api-reference and user-guide with:                                                             
                                                                                                                          
 - API signatures.                                                                                                        
 - Examples for SQL schema rendering.                                                                                     
 - Examples for JSON debug records.                                                                                       
 - Safety note: code block content is escaped; do not use ui.raw for database/request text.                               
                                                                                                                          
 Example doc snippet:                                                                                                     
                                                                                                                          
 ```js                                                                                                                    
   res.html(ui.page({ title: "Schema" },                                                                                  
     ui.h1(row.name),                                                                                                     
     ui.badge(row.type),                                                                                                  
     ui.codeBlock("sql", row.sql, {                                                                                       
       title: row.name,                                                                                                   
       lineNumbers: true,                                                                                                 
       copy: true,                                                                                                        
     })                                                                                                                   
   ));                                                                                                                    
 ```"

**Assistant interpretation:** Create a new ticket and implementation guide from the colleague's component spec, then proceed step-by-step with tasks, diary, and focused commits.

**Inferred user intent:** Preserve the design in a ticket before implementing, then build the components carefully with reviewable history and documentation updates.

**Commit (code):** pending — planning docs.

### What I did

- Created docmgr ticket `DB-BROWSER-UIDSL-COMPONENTS`.
- Added design document:
  - `design-doc/01-ui-dsl-codeblock-badge-tabs-implementation-guide.md`.
- Added implementation diary:
  - `reference/01-implementation-diary.md`.
- Rewrote `tasks.md` with T01–T07 implementation tasks.
- The implementation guide maps the colleague spec to concrete repo files, helper functions, tests, examples, and embedded docs.

### Why

- The component spec touches public JavaScript APIs, render contracts, safety guarantees, tests, examples, and help docs. A ticket and implementation guide reduce the chance of building only the happy path or missing the safety requirements.

### What worked

- Ticket and documents were created successfully.
- The guide now contains enough file-level and test-level direction for implementation to proceed.

### What didn't work

- N/A in this step.

### What I learned

- The existing table DSL already has several reusable patterns for these components: CSS token normalization, explicit Goja wrapper functions, and node composition.

### What was tricky to build

- The user prompt contains a detailed spec, so the main challenge was preserving its contracts while translating them into the current implementation architecture and task sequence.

### What warrants a second pair of eyes

- Review the tabs render contract before implementation. CSS-only radio tabs are preferred, but the doc still allows a `<details>` fallback if needed.
- Review the open question about `lineNumbers`: class-only vs actual line markup.

### What should be done in the future

- Implement T02 next: `ui.codeBlock`, `ui.sql`, `ui.js`, and `ui.jsonBlock` with tests.

### Code review instructions

- Start with the design doc and `tasks.md` in the ticket workspace.
- Validate ticket hygiene with:

```bash
docmgr doctor --ticket DB-BROWSER-UIDSL-COMPONENTS --stale-after 30
```

### Technical details

Ticket path:

```text
ttmp/2026/05/07/DB-BROWSER-UIDSL-COMPONENTS--ui-dsl-component-spec-for-code-blocks-badges-and-tabs
```

## Step 2: Implement codeBlock, badge, and tabs core components

I implemented the three requested components in the UI DSL core and registered them with the Goja module loader. The implementation keeps the components server-rendered and safe by default: source code, badge values, labels, and string tab content are represented as `Text` nodes and therefore go through the existing HTML escaping renderer.

I implemented the three components in one focused pass because they share helper code for options, token normalization, and node composition. The tests cover render contracts, escaping, aliases, duplicate tab IDs, disabled tabs, and selected-tab fallback behavior.

### Prompt Context

**User prompt (verbatim):** "continue"

**Assistant interpretation:** Continue after planning by implementing the component tasks step by step, with validation, diary updates, and commits.

**Inferred user intent:** Move from the implementation guide into working code while preserving the ticket workflow.

**Commit (code):** pending — core component implementation.

### What I did

- Added `internal/uidsl/components.go`.
- Registered new exports in `internal/uidsl/module.go`:
  - `ui.codeBlock`;
  - `ui.sql`;
  - `ui.js`;
  - `ui.jsonBlock`;
  - `ui.badge`;
  - `ui.tabs`.
- Added `internal/uidsl/components_test.go`.
- Implemented code block options:
  - `title`;
  - `lineNumbers`;
  - `wrap` defaulting to true;
  - `copy`;
  - `maxHeight`;
  - `class`.
- Implemented badge options:
  - `tone` with fallback to `default`;
  - `title`;
  - `class`.
- Implemented CSS-only radio tab markup with:
  - normalized container/tab IDs;
  - duplicate tab suffixes;
  - selected tab resolution by id/index;
  - disabled tabs;
  - content normalization through `NormalizeExport`.

### Why

- These components are required for inspection/debug pages and reduce unsafe ad hoc use of `ui.raw` for SQL, JSON, and request/database text.

### What worked

- `go test ./internal/uidsl -count=1` passed.
- `go test ./...` passed.

### What didn't work

- The first test run failed in three places:

```text
--- FAIL: TestCodeBlockAliasesAndJSON
missing "&quot;a&quot;: &quot;&lt;x&gt;&quot;" ... actual JSON used &#34; and \u003c escaping

--- FAIL: TestTabsRenderSelectedDuplicateDisabledAndEscaped
missing checked attribute substring because renderer sorts boolean attributes first

--- FAIL: TestTabsInvalidSelectionFallsBack
same checked-attribute ordering issue
```

- Fixes:
  - Updated JSON expectations to match Go's safe JSON escaping (`\u003c`) plus HTML numeric quote escaping (`&#34;`).
  - Updated checked-radio assertions to match the renderer's sorted boolean attribute output.

### What I learned

- JSON pretty printing through Go's encoder escapes `<` as `\u003c`, which is desirable for HTML safety but should be documented/expected in tests.
- The renderer sorts attributes and renders boolean attributes before string attributes when their key sorts first, so exact render-contract tests need to account for that deterministic ordering.

### What was tricky to build

- The main tricky part was tabs content normalization. Goja exports UI nodes as Go values, strings as strings, and arrays as `[]any`; the implementation uses the existing `NormalizeExport` path so it stays consistent with the rest of the DSL.

### What warrants a second pair of eyes

- Review `jsonBlock` formatting expectations. It currently uses Go JSON marshal indentation and therefore escapes `<`, `>`, and `&` as unicode sequences.
- Review whether `ui.badge(nil)` should render an empty badge or a textual fallback such as `unknown`.
- Review whether `lineNumbers` should eventually produce per-line markup instead of only a class.

### What should be done in the future

- Update examples to use the new components.
- Update `js-api-reference` and `user-guide`.
- Add component smoke script.

### Code review instructions

- Start with `internal/uidsl/components.go`.
- Review `internal/uidsl/module.go` for export wiring.
- Review `internal/uidsl/components_test.go` for the public JS API contract.
- Validate with:

```bash
go test ./internal/uidsl -count=1
go test ./...
```

### Technical details

Example API now implemented:

```js
ui.codeBlock("sql", row.sql, { title: row.name, lineNumbers: true, copy: true })
ui.badge(row.type, { tone: "info" })
ui.tabs("record-tabs", [
  { id: "sql", label: "SQL", content: ui.sql(row.sql) },
  { id: "json", label: "JSON", content: ui.jsonBlock(row) },
])
```
