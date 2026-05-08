# Tasks

## Scoped ui.table query state implementation sequence

### T01 — Planning and implementation guide

- [x] Create docmgr ticket workspace.
- [x] Write scoped table query state implementation guide.
- [x] Add later implementation task list.
- [x] Validate ticket hygiene.
- [x] Commit planning docs.

### T02 — Scoped query parsing tests

- [ ] Add `internal/uidsl/table_scoped_query_test.go`.
- [ ] Test global behavior remains unchanged without scope.
- [ ] Test `scope: true` uses table id.
- [ ] Test `scope: "custom"` uses custom namespace.
- [ ] Test strict scoped mode ignores global `page/sort/filter` keys.

### T03 — Scope support in render context

- [ ] Add `Scope` to `RenderContext`.
- [ ] Parse `render({ scope })` in `TableBuilder.context`.
- [ ] Normalize scope tokens safely.
- [ ] Keep current global behavior when scope is omitted or false.

### T04 — Scoped query helpers

- [ ] Add helpers for scoped query key lookup and updates.
- [ ] Add scoped filter parsing for `scope.q` and `scope.filter.<column>`.
- [ ] Add helper to remove one scope's keys for filter clear links.
- [ ] Unit-test pure helper behavior.

### T05 — Hrefs and filter form generation

- [ ] Update sortable header links to write scoped `sort/dir/page` keys.
- [ ] Update pagination links to write scoped `page` keys.
- [ ] Update filter input names to use scoped names.
- [ ] Update hidden sort fields in the filter form.
- [ ] Update clear link to clear only this table's scope.

### T06 — Examples and documentation

- [ ] Update `examples/generic-browser/scripts/app.js` to use scoped table state where multiple tables can appear.
- [ ] Update `js-api-reference` table docs.
- [ ] Update `user-guide` and `app-playbook` with multi-table guidance.
- [ ] Add trace-browser migration snippet if useful.

### T07 — Smoke validation

- [ ] Add ticket-local smoke script with two independent tables.
- [ ] Validate `left.page=2` does not affect right table.
- [ ] Validate generated links preserve unrelated scoped query state.
- [ ] Run `go test ./...` and `docmgr doctor`.
- [ ] Commit implementation and docs.

## Future follow-ups

- [ ] Add scoped tab selected-state if needed.
- [ ] Consider `t.<scope>.<key>` prefix if collisions appear.
- [ ] Consider server-interactive partial table refreshes after scoped GET state is implemented.
