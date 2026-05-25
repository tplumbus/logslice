package filter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// JSONPathFilter matches log lines where a dot-separated nested field equals a given value.
type JSONPathFilter struct {
	path  []string
	value string
}

// NewJSONPathFilter creates a filter that matches lines where the nested field
// at dotPath equals value. dotPath uses dot notation, e.g. "meta.user.id".
func NewJSONPathFilter(dotPath, value string) (*JSONPathFilter, error) {
	if strings.TrimSpace(dotPath) == "" {
		return nil, fmt.Errorf("jsonpath: dot path must not be empty")
	}
	parts := strings.Split(dotPath, ".")
	for _, p := range parts {
		if strings.TrimSpace(p) == "" {
			return nil, fmt.Errorf("jsonpath: path segment must not be empty in %q", dotPath)
		}
	}
	return &JSONPathFilter{path: parts, value: value}, nil
}

// MatchesLine returns true if the nested field resolved by the dot path equals the target value.
func (f *JSONPathFilter) MatchesLine(line string) bool {
	var root map[string]interface{}
	if err := json.Unmarshal([]byte(line), &root); err != nil {
		return false
	}
	resolved, ok := resolvePath(root, f.path)
	if !ok {
		return false
	}
	return fmt.Sprintf("%v", resolved) == f.value
}

// TransformLine returns the line unchanged; JSONPathFilter is a predicate only.
func (f *JSONPathFilter) TransformLine(line string) string {
	return line
}

// resolvePath walks a nested map following the given key segments.
func resolvePath(node map[string]interface{}, path []string) (interface{}, bool) {
	if len(path) == 0 {
		return nil, false
	}
	val, ok := node[path[0]]
	if !ok {
		return nil, false
	}
	if len(path) == 1 {
		return val, true
	}
	child, ok := val.(map[string]interface{})
	if !ok {
		return nil, false
	}
	return resolvePath(child, path[1:])
}

// ParseJSONPathFlag parses a flag value of the form "dot.path=value".
func ParseJSONPathFlag(s string) (*JSONPathFilter, error) {
	idx := strings.Index(s, "=")
	if idx < 1 {
		return nil, fmt.Errorf("jsonpath: expected format 'dot.path=value', got %q", s)
	}
	return NewJSONPathFilter(s[:idx], s[idx+1:])
}
