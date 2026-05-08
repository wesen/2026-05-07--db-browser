package uidsl

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/dop251/goja"
)

var allowedBadgeTones = map[string]bool{
	"default": true,
	"info":    true,
	"success": true,
	"warning": true,
	"danger":  true,
	"muted":   true,
}

type codeBlockOptions struct {
	Title       string
	LineNumbers bool
	Wrap        bool
	Copy        bool
	MaxHeight   string
	Class       string
}

func codeBlockNode(language string, source any, opts map[string]any) Node {
	lang := cssToken(language)
	if lang == "empty" {
		lang = "text"
	}
	options := parseCodeBlockOptions(opts)
	classes := []any{"ui-codeblock", "ui-codeblock--" + lang}
	if options.Wrap {
		classes = append(classes, "ui-codeblock--wrap")
	} else {
		classes = append(classes, "ui-codeblock--nowrap")
	}
	if options.LineNumbers {
		classes = append(classes, "ui-codeblock--line-numbers")
	}
	if options.Class != "" {
		classes = append(classes, options.Class)
	}
	preAttrs := map[string]any{}
	if options.MaxHeight != "" {
		preAttrs["style"] = map[string]any{"max-height": options.MaxHeight, "overflow": "auto"}
	}
	code := &Element{Tag: "code", Attrs: map[string]any{"class": "language-" + lang}, Children: []Node{&Text{Value: stringifySource(source)}}}
	if options.Title == "" && !options.Copy {
		preAttrs["class"] = classes
		return &Element{Tag: "pre", Attrs: preAttrs, Children: []Node{code}}
	}
	captionChildren := []Node{}
	if options.Title != "" {
		captionChildren = append(captionChildren, &Element{Tag: "span", Attrs: map[string]any{"class": "ui-codeblock__title"}, Children: []Node{&Text{Value: options.Title}}})
	}
	if options.Copy {
		captionChildren = append(captionChildren, &Element{Tag: "button", Attrs: map[string]any{"class": "ui-codeblock__copy", "type": "button"}, Children: []Node{&Text{Value: "Copy"}}})
	}
	preAttrs["class"] = "ui-codeblock__pre"
	return &Element{Tag: "figure", Attrs: map[string]any{"class": classes}, Children: []Node{
		&Element{Tag: "figcaption", Attrs: map[string]any{"class": "ui-codeblock__caption"}, Children: captionChildren},
		&Element{Tag: "pre", Attrs: preAttrs, Children: []Node{code}},
	}}
}

func parseCodeBlockOptions(opts map[string]any) codeBlockOptions {
	return codeBlockOptions{
		Title:       stringFromAny(opts["title"]),
		LineNumbers: boolFromAny(opts["lineNumbers"]),
		Wrap:        optionBool(opts, "wrap", true),
		Copy:        boolFromAny(opts["copy"]),
		MaxHeight:   stringFromAny(opts["maxHeight"]),
		Class:       stringFromAny(opts["class"]),
	}
}

func badgeNode(value any, opts map[string]any) Node {
	text := stringifySource(value)
	tone := cssToken(stringFromAny(opts["tone"]))
	if !allowedBadgeTones[tone] {
		tone = "default"
	}
	classes := []any{"ui-badge", "ui-badge--" + tone, "ui-badge--value-" + cssToken(text)}
	if className := stringFromAny(opts["class"]); className != "" {
		classes = append(classes, className)
	}
	attrs := map[string]any{"class": classes}
	if title := stringFromAny(opts["title"]); title != "" {
		attrs["title"] = title
	}
	return &Element{Tag: "span", Attrs: attrs, Children: []Node{&Text{Value: text}}}
}

type tabSpec struct {
	ID       string
	Label    string
	Content  Node
	Disabled bool
}

func tabsNode(id string, value goja.Value, opts map[string]any) (Node, error) {
	baseID := domToken(id, "tabs")
	tabs, err := parseTabs(value)
	if err != nil {
		return nil, err
	}
	if len(tabs) == 0 {
		return &Element{Tag: "div", Attrs: map[string]any{"class": []any{"ui-tabs", stringFromAny(opts["class"])}, "id": baseID}}, nil
	}
	assignTabIDs(tabs)
	selected := selectedTabIndex(tabs, opts["selected"])
	classes := []any{"ui-tabs"}
	if className := stringFromAny(opts["class"]); className != "" {
		classes = append(classes, className)
	}
	children := []Node{&Element{Tag: "div", Attrs: map[string]any{"class": "ui-tabs__tablist", "role": "tablist"}, Children: tabControlNodes(baseID, tabs, selected)}}
	panels := make([]Node, 0, len(tabs))
	for i, tab := range tabs {
		panelClasses := []any{"ui-tabs__panel"}
		if i == selected {
			panelClasses = append(panelClasses, "ui-tabs__panel--active")
		}
		panels = append(panels, &Element{Tag: "section", Attrs: map[string]any{"class": panelClasses, "data-tab": tab.ID}, Children: []Node{tab.Content}})
	}
	children = append(children, &Element{Tag: "div", Attrs: map[string]any{"class": "ui-tabs__panels"}, Children: panels})
	return &Element{Tag: "div", Attrs: map[string]any{"class": classes, "id": baseID}, Children: children}, nil
}

func parseTabs(value goja.Value) ([]tabSpec, error) {
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return nil, nil
	}
	exported := value.Export()
	items, ok := exported.([]any)
	if !ok {
		return nil, fmt.Errorf("ui.tabs expects an array of tab specs, got %T", exported)
	}
	tabs := make([]tabSpec, 0, len(items))
	for i, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("ui.tabs tab %d must be an object, got %T", i, item)
		}
		content, err := NormalizeExport(m["content"])
		if err != nil {
			return nil, err
		}
		label := stringFromAny(m["label"])
		if label == "" {
			label = fmt.Sprintf("Tab %d", i+1)
		}
		tabs = append(tabs, tabSpec{ID: stringFromAny(m["id"]), Label: label, Content: content, Disabled: boolFromAny(m["disabled"])})
	}
	return tabs, nil
}

func assignTabIDs(tabs []tabSpec) {
	seen := map[string]int{}
	for i := range tabs {
		base := domToken(tabs[i].ID, "")
		if base == "" {
			base = domToken(tabs[i].Label, fmt.Sprintf("tab-%d", i+1))
		}
		seen[base]++
		if seen[base] > 1 {
			tabs[i].ID = fmt.Sprintf("%s-%d", base, seen[base])
		} else {
			tabs[i].ID = base
		}
	}
}

func selectedTabIndex(tabs []tabSpec, selected any) int {
	fallback := -1
	for i, tab := range tabs {
		if !tab.Disabled && fallback == -1 {
			fallback = i
		}
	}
	if fallback == -1 {
		return 0
	}
	if idx, ok := intLike(selected); ok && idx >= 0 && idx < len(tabs) && !tabs[idx].Disabled {
		return idx
	}
	selectedID := domToken(stringFromAny(selected), "")
	if selectedID != "" {
		for i, tab := range tabs {
			if tab.ID == selectedID && !tab.Disabled {
				return i
			}
		}
	}
	return fallback
}

func tabControlNodes(baseID string, tabs []tabSpec, selected int) []Node {
	children := make([]Node, 0, len(tabs)*2)
	for i, tab := range tabs {
		inputID := baseID + "-" + tab.ID
		inputAttrs := map[string]any{"class": "ui-tabs__radio", "type": "radio", "name": baseID, "id": inputID}
		if i == selected && !tab.Disabled {
			inputAttrs["checked"] = true
		}
		if tab.Disabled {
			inputAttrs["disabled"] = true
		}
		labelClasses := []any{"ui-tabs__tab"}
		if tab.Disabled {
			labelClasses = append(labelClasses, "ui-tabs__tab--disabled")
		}
		children = append(children,
			&Element{Tag: "input", Attrs: inputAttrs},
			&Element{Tag: "label", Attrs: map[string]any{"class": labelClasses, "for": inputID}, Children: []Node{&Text{Value: tab.Label}}},
		)
	}
	return children
}

func jsonBlockSource(value goja.Value) string {
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return "null"
	}
	exported := value.Export()
	if s, ok := exported.(string); ok {
		var parsed any
		if err := json.Unmarshal([]byte(s), &parsed); err == nil {
			if b, err := json.MarshalIndent(parsed, "", "  "); err == nil {
				return string(b)
			}
		}
		return s
	}
	if b, err := json.MarshalIndent(exported, "", "  "); err == nil {
		return string(b)
	}
	return stringifySource(exported)
}

func stringifySource(value any) string {
	if value == nil {
		return ""
	}
	return fmt.Sprint(value)
}

func firstOptions(options []map[string]any) map[string]any {
	if len(options) == 0 || options[0] == nil {
		return map[string]any{}
	}
	return options[0]
}

func optionBool(opts map[string]any, key string, fallback bool) bool {
	if _, ok := opts[key]; !ok {
		return fallback
	}
	return boolFromAny(opts[key])
}

func domToken(value string, fallback string) string {
	token := cssToken(value)
	if token == "empty" {
		return fallback
	}
	return token
}

func intLike(value any) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	case string:
		n, err := strconv.Atoi(v)
		return n, err == nil
	default:
		return 0, false
	}
}

func classList(values ...any) []any {
	classes := []any{}
	for _, value := range values {
		if s := strings.TrimSpace(fmt.Sprint(value)); s != "" && s != "<nil>" {
			classes = append(classes, s)
		}
	}
	return classes
}
