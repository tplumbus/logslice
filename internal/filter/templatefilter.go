package filter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
)

// TemplateFilter rewrites each matched log line using a Go text/template.
// The template receives the parsed JSON object as a map[string]any.
type TemplateFilter struct {
	tmpl *template.Template
}

// NewTemplateFilter compiles tmplStr as a Go text/template.
// Returns an error if tmplStr is empty or fails to parse.
func NewTemplateFilter(tmplStr string) (*TemplateFilter, error) {
	if tmplStr == "" {
		return nil, fmt.Errorf("templatefilter: template string must not be empty")
	}
	t, err := template.New("logslice").Option("missingkey=zero").Parse(tmplStr)
	if err != nil {
		return nil, fmt.Errorf("templatefilter: invalid template: %w", err)
	}
	return &TemplateFilter{tmpl: t}, nil
}

// MatchesLine always returns true; TemplateFilter is a transformer, not a gate.
func (f *TemplateFilter) MatchesLine(line string) bool {
	return true
}

// TransformLine applies the template to the JSON line.
// If the line is not valid JSON the original line is returned unchanged.
func (f *TemplateFilter) TransformLine(line string) string {
	var data map[string]any
	if err := json.Unmarshal([]byte(line), &data); err != nil {
		return line
	}
	var buf bytes.Buffer
	if err := f.tmpl.Execute(&buf, data); err != nil {
		return line
	}
	return buf.String()
}
