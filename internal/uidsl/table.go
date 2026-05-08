package uidsl

import (
	"fmt"
	"sort"

	"github.com/dop251/goja"
)

type TableBuilder struct {
	ID       string
	Rows     []map[string]any
	features TableFeatures
}

type TableFeatures struct {
	Pagination   bool
	Sorting      bool
	ColumnPicker bool
}

type FeatureBuilder struct {
	features *TableFeatures
}

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
	_ = obj.Set("render", func(_ map[string]any) Node {
		return builder.node()
	})
	return obj
}

func (t *TableBuilder) node() Node {
	columns := tableColumns(t.Rows)
	headCells := make([]Node, 0, len(columns))
	for _, column := range columns {
		headCells = append(headCells, &Element{Tag: "th", Children: []Node{&Text{Value: column}}})
	}
	bodyRows := make([]Node, 0, len(t.Rows))
	for _, row := range t.Rows {
		cells := make([]Node, 0, len(columns))
		for _, column := range columns {
			cells = append(cells, &Element{Tag: "td", Children: []Node{&Text{Value: fmt.Sprint(row[column])}}})
		}
		bodyRows = append(bodyRows, &Element{Tag: "tr", Children: cells})
	}
	attrs := map[string]any{"class": tableClass(t.features)}
	if t.ID != "" {
		attrs["id"] = t.ID
	}
	return &Element{Tag: "table", Attrs: attrs, Children: []Node{
		&Element{Tag: "thead", Children: []Node{&Element{Tag: "tr", Children: headCells}}},
		&Element{Tag: "tbody", Children: bodyRows},
	}}
}

func featureBuilderObject(vm *goja.Runtime, features *TableFeatures) goja.Value {
	obj := vm.NewObject()
	_ = obj.Set("pagination", func(_ ...any) goja.Value {
		features.Pagination = true
		return obj
	})
	_ = obj.Set("sorting", func(_ ...any) goja.Value {
		features.Sorting = true
		return obj
	})
	_ = obj.Set("columnPicker", func(_ ...any) goja.Value {
		features.ColumnPicker = true
		return obj
	})
	return obj
}

func rowsFromValue(value goja.Value) ([]map[string]any, error) {
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return []map[string]any{}, nil
	}
	exported := value.Export()
	slice, ok := exported.([]any)
	if !ok {
		return nil, fmt.Errorf("ui.table.fromRows expects an array, got %T", exported)
	}
	rows := make([]map[string]any, 0, len(slice))
	for _, item := range slice {
		switch row := item.(type) {
		case map[string]any:
			rows = append(rows, row)
		default:
			rows = append(rows, map[string]any{"value": row})
		}
	}
	return rows, nil
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
