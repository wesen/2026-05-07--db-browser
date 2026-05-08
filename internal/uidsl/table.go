package uidsl

import (
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/dop251/goja"
)

type TableBuilder struct {
	ID        string
	Rows      []map[string]any
	features  TableFeatures
	dataFn    goja.Callable
	columnsFn goja.Callable
}

type TableFeatures struct {
	Pagination   bool
	Sorting      bool
	ColumnPicker bool
	PageSize     int
}

type ColumnSpec struct {
	Name     string
	Label    string
	Kind     string
	Sortable bool
	Align    string
}

type RenderContext struct {
	Query  map[string]any    `json:"query"`
	Page   map[string]any    `json:"page"`
	Order  map[string]any    `json:"order"`
	State  map[string]any    `json:"state"`
	Filter map[string]any    `json:"filter"`
	Params map[string]string `json:"params"`
}

type FeatureBuilder struct{ features *TableFeatures }
type ColumnBuilder struct{ columns *[]ColumnSpec }

func newTableBuilder(id string) *TableBuilder {
	return &TableBuilder{ID: id, Rows: []map[string]any{}}
}

func tableFromRows(id string, value goja.Value) (*TableBuilder, error) {
	rows, err := rowsFromValue(value)
	if err != nil {
		return nil, err
	}
	return &TableBuilder{ID: id, Rows: rows}, nil
}

func tableBuilderObject(vm *goja.Runtime, builder *TableBuilder) goja.Value {
	obj := vm.NewObject()
	_ = obj.Set("features", func(callback goja.Value) goja.Value {
		if fn, ok := goja.AssertFunction(callback); ok {
			_, _ = fn(goja.Undefined(), featureBuilderObject(vm, &builder.features))
		}
		return obj
	})
	_ = obj.Set("data", func(callback goja.Value) goja.Value {
		if fn, ok := goja.AssertFunction(callback); ok {
			builder.dataFn = fn
		}
		return obj
	})
	_ = obj.Set("columns", func(callback goja.Value) goja.Value {
		if fn, ok := goja.AssertFunction(callback); ok {
			builder.columnsFn = fn
		}
		return obj
	})
	_ = obj.Set("render", func(input map[string]any) (Node, error) {
		return builder.render(vm, input)
	})
	return obj
}

func (t *TableBuilder) render(vm *goja.Runtime, input map[string]any) (Node, error) {
	ctx := t.context(input)
	rows := append([]map[string]any(nil), t.Rows...)
	total := len(rows)
	if t.dataFn != nil {
		value, err := t.dataFn(goja.Undefined(), vm.ToValue(ctx))
		if err != nil {
			return nil, err
		}
		dataRows, dataTotal, err := dataResultFromValue(value)
		if err != nil {
			return nil, err
		}
		rows = dataRows
		total = dataTotal
	}
	columns, err := t.resolveColumns(vm, ctx, rows)
	if err != nil {
		return nil, err
	}
	return t.node(ctx, rows, total, columns), nil
}

func (t *TableBuilder) context(input map[string]any) RenderContext {
	query := mapFromAny(input["query"])
	params := stringMapFromAny(input["params"])
	pageSize := intFromAny(query["limit"], t.features.PageSize)
	if pageSize <= 0 {
		pageSize = 25
	}
	if pageSize > 500 {
		pageSize = 500
	}
	pageIndex := intFromAny(query["page"], 1)
	if pageIndex < 1 {
		pageIndex = 1
	}
	sortKey := stringFromAny(query["sort"])
	dir := strings.ToLower(stringFromAny(query["dir"]))
	if dir != "desc" {
		dir = "asc"
	}
	return RenderContext{
		Query:  query,
		Params: params,
		State:  map[string]any{},
		Filter: query,
		Page: map[string]any{
			"index":  pageIndex,
			"size":   pageSize,
			"limit":  pageSize,
			"offset": (pageIndex - 1) * pageSize,
		},
		Order: map[string]any{"key": sortKey, "dir": dir},
	}
}

func (t *TableBuilder) resolveColumns(vm *goja.Runtime, ctx RenderContext, rows []map[string]any) ([]ColumnSpec, error) {
	if t.columnsFn == nil {
		cols := tableColumns(rows)
		ret := make([]ColumnSpec, 0, len(cols))
		for _, col := range cols {
			ret = append(ret, ColumnSpec{Name: col, Label: col, Kind: "text"})
		}
		return ret, nil
	}
	cols := []ColumnSpec{}
	cb := columnBuilderObject(vm, &cols)
	value, err := t.columnsFn(goja.Undefined(), cb, vm.ToValue(ctx))
	if err != nil {
		return nil, err
	}
	if exported := value.Export(); exported != nil {
		if parsed := columnsFromExport(exported); len(parsed) > 0 {
			return parsed, nil
		}
	}
	return cols, nil
}

func (t *TableBuilder) node(ctx RenderContext, rows []map[string]any, total int, columns []ColumnSpec) Node {
	headCells := make([]Node, 0, len(columns))
	for _, column := range columns {
		label := column.Label
		if label == "" {
			label = column.Name
		}
		var child Node = &Text{Value: label}
		if t.features.Sorting && column.Sortable {
			dir := "asc"
			if stringFromAny(ctx.Order["key"]) == column.Name && stringFromAny(ctx.Order["dir"]) == "asc" {
				dir = "desc"
			}
			child = &Element{Tag: "a", Attrs: map[string]any{"href": queryHref(ctx.Query, map[string]any{"sort": column.Name, "dir": dir, "page": 1})}, Children: []Node{&Text{Value: label}}}
		}
		headCells = append(headCells, &Element{Tag: "th", Attrs: map[string]any{"data-column": column.Name}, Children: []Node{child}})
	}
	bodyRows := make([]Node, 0, len(rows))
	for _, row := range rows {
		cells := make([]Node, 0, len(columns))
		for _, column := range columns {
			attrs := map[string]any{"data-column": column.Name}
			if column.Align != "" {
				attrs["class"] = "align-" + column.Align
			}
			cells = append(cells, &Element{Tag: "td", Attrs: attrs, Children: []Node{&Text{Value: formatCell(row[column.Name], column)}}})
		}
		bodyRows = append(bodyRows, &Element{Tag: "tr", Children: cells})
	}
	attrs := map[string]any{"class": tableClass(t.features)}
	if t.ID != "" {
		attrs["id"] = t.ID
	}
	table := &Element{Tag: "table", Attrs: attrs, Children: []Node{
		&Element{Tag: "thead", Children: []Node{&Element{Tag: "tr", Children: headCells}}},
		&Element{Tag: "tbody", Children: bodyRows},
	}}
	if !t.features.Pagination {
		return table
	}
	return &Fragment{Children: []Node{table, paginationNode(ctx, total)}}
}

func featureBuilderObject(vm *goja.Runtime, features *TableFeatures) goja.Value {
	obj := vm.NewObject()
	_ = obj.Set("pagination", func(options ...map[string]any) goja.Value {
		features.Pagination = true
		if len(options) > 0 {
			features.PageSize = intFromAny(options[0]["size"], features.PageSize)
		}
		return obj
	})
	_ = obj.Set("sorting", func(_ ...any) goja.Value { features.Sorting = true; return obj })
	_ = obj.Set("columnPicker", func(_ ...any) goja.Value { features.ColumnPicker = true; return obj })
	return obj
}

func columnBuilderObject(vm *goja.Runtime, columns *[]ColumnSpec) goja.Value {
	obj := vm.NewObject()
	add := func(kind string, name string) goja.Value {
		col := ColumnSpec{Name: name, Label: name, Kind: kind}
		*columns = append(*columns, col)
		return columnObject(vm, &(*columns)[len(*columns)-1], columns)
	}
	for _, kind := range []string{"text", "badge", "money", "date", "tags"} {
		kind := kind
		_ = obj.Set(kind, func(name string, _ ...map[string]any) goja.Value { return add(kind, name) })
	}
	return obj
}

func columnObject(vm *goja.Runtime, col *ColumnSpec, columns *[]ColumnSpec) goja.Value {
	obj := vm.NewObject()
	_ = obj.Set("label", func(label string) goja.Value { col.Label = label; return obj })
	_ = obj.Set("sortable", func(_ ...any) goja.Value { col.Sortable = true; return obj })
	_ = obj.Set("align", func(align string) goja.Value { col.Align = align; return obj })
	_ = obj.Set("mono", func(_ ...any) goja.Value { return obj })
	_ = obj.Set("truncate", func(_ ...any) goja.Value { return obj })
	_ = obj.Set("link", func(_ ...any) goja.Value { return obj })
	for _, kind := range []string{"text", "badge", "money", "date", "tags"} {
		kind := kind
		_ = obj.Set(kind, func(name string, _ ...map[string]any) goja.Value {
			next := ColumnSpec{Name: name, Label: name, Kind: kind}
			*columns = append(*columns, next)
			return columnObject(vm, &(*columns)[len(*columns)-1], columns)
		})
	}
	return obj
}

func rowsFromValue(value goja.Value) ([]map[string]any, error) {
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return []map[string]any{}, nil
	}
	exported := value.Export()
	rows, ok := rowsFromExport(exported)
	if !ok {
		return nil, fmt.Errorf("ui.table.fromRows expects an array, got %T", exported)
	}
	return rows, nil
}

func rowsFromExport(exported any) ([]map[string]any, bool) {
	if exported == nil {
		return []map[string]any{}, true
	}
	if slice, ok := exported.([]any); ok {
		return rowsFromSlice(slice), true
	}
	rv := reflect.ValueOf(exported)
	if !rv.IsValid() || (rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array) {
		return nil, false
	}
	rows := make([]map[string]any, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		rows = append(rows, rowFromAny(rv.Index(i).Interface()))
	}
	return rows, true
}

func rowsFromSlice(slice []any) []map[string]any {
	rows := make([]map[string]any, 0, len(slice))
	for _, item := range slice {
		rows = append(rows, rowFromAny(item))
	}
	return rows
}

func rowFromAny(item any) map[string]any {
	if row, ok := item.(map[string]any); ok {
		return row
	}
	rv := reflect.ValueOf(item)
	if rv.IsValid() && rv.Kind() == reflect.Map && rv.Type().Key().Kind() == reflect.String {
		row := map[string]any{}
		iter := rv.MapRange()
		for iter.Next() {
			row[iter.Key().String()] = iter.Value().Interface()
		}
		return row
	}
	return map[string]any{"value": item}
}

func dataResultFromValue(value goja.Value) ([]map[string]any, int, error) {
	exported := value.Export()
	if rows, ok := rowsFromExport(exported); ok {
		return rows, len(rows), nil
	}
	m, ok := exported.(map[string]any)
	if !ok {
		return nil, 0, fmt.Errorf("table data callback must return rows array or {rows,total}, got %T", exported)
	}
	rows, _ := rowsFromExport(m["rows"])
	total := intFromAny(m["total"], len(rows))
	return rows, total, nil
}

func tableColumns(rows []map[string]any) []string {
	seen := map[string]struct{}{}
	columns := []string{}
	for _, row := range rows {
		keys := make([]string, 0, len(row))
		for key := range row {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			columns = append(columns, key)
		}
	}
	return columns
}

func columnsFromExport(exported any) []ColumnSpec {
	items, ok := exported.([]any)
	if !ok {
		return nil
	}
	ret := []ColumnSpec{}
	for _, item := range items {
		if m, ok := item.(map[string]any); ok {
			name := stringFromAny(m["name"])
			if name == "" {
				continue
			}
			ret = append(ret, ColumnSpec{Name: name, Label: stringFromAny(m["label"]), Kind: stringFromAny(m["kind"]), Sortable: boolFromAny(m["sortable"]), Align: stringFromAny(m["align"])})
		}
	}
	return ret
}

func tableClass(features TableFeatures) []any {
	classes := []any{"ui-table"}
	if features.Pagination {
		classes = append(classes, "ui-table--pagination")
	}
	if features.Sorting {
		classes = append(classes, "ui-table--sorting")
	}
	if features.ColumnPicker {
		classes = append(classes, "ui-table--column-picker")
	}
	return classes
}

func paginationNode(ctx RenderContext, total int) Node {
	page := intFromAny(ctx.Page["index"], 1)
	limit := intFromAny(ctx.Page["limit"], 25)
	pages := 1
	if limit > 0 && total > 0 {
		pages = (total + limit - 1) / limit
	}
	children := []Node{&Text{Value: fmt.Sprintf("Page %d of %d (%d rows)", page, pages, total)}}
	if page > 1 {
		children = append(children, &Text{Value: " "}, &Element{Tag: "a", Attrs: map[string]any{"href": queryHref(ctx.Query, map[string]any{"page": page - 1})}, Children: []Node{&Text{Value: "Previous"}}})
	}
	if page < pages {
		children = append(children, &Text{Value: " "}, &Element{Tag: "a", Attrs: map[string]any{"href": queryHref(ctx.Query, map[string]any{"page": page + 1})}, Children: []Node{&Text{Value: "Next"}}})
	}
	return &Element{Tag: "nav", Attrs: map[string]any{"class": "ui-table-pagination"}, Children: children}
}

func formatCell(value any, col ColumnSpec) string {
	if value == nil {
		return ""
	}
	return fmt.Sprint(value)
}

func queryHref(query map[string]any, updates map[string]any) string {
	values := url.Values{}
	for key, value := range query {
		if value == nil {
			continue
		}
		values.Set(key, fmt.Sprint(value))
	}
	for key, value := range updates {
		values.Set(key, fmt.Sprint(value))
	}
	encoded := values.Encode()
	if encoded == "" {
		return "?"
	}
	return "?" + encoded
}

func mapFromAny(value any) map[string]any {
	ret := map[string]any{}
	if m, ok := value.(map[string]any); ok {
		for k, v := range m {
			ret[k] = v
		}
	}
	return ret
}

func stringMapFromAny(value any) map[string]string {
	ret := map[string]string{}
	if m, ok := value.(map[string]any); ok {
		for k, v := range m {
			ret[k] = fmt.Sprint(v)
		}
	}
	return ret
}

func intFromAny(value any, fallback int) int {
	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func stringFromAny(value any) string {
	if value == nil {
		return ""
	}
	return fmt.Sprint(value)
}

func boolFromAny(value any) bool {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		return v == "true"
	default:
		return false
	}
}
