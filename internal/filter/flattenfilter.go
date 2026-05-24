package filter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FlattenFilter flattens nested JSON objects into dot-notation keys.
// For example, {"a":{"b":1}} becomes {"a.b":1}.
type FlattenFilter struct {
	prefix string
	sep    string
}

// NewFlattenFilter creates a FlattenFilter. sep is the separator used for
// joining nested keys (e.g. "."). Returns an error if sep is empty.
func NewFlattenFilter(sep string) (*FlattenFilter, error) {
	if sep == "" {
		return nil, fmt.Errorf("flattenfilter: separator must not be empty")
	}
	return &FlattenFilter{sep: sep}, nil
}

// MatchesLine always returns true; FlattenFilter is a transform-only filter.
func (f *FlattenFilter) MatchesLine(line string) bool {
	return true
}

// TransformLine flattens the JSON object in line. If the line cannot be
// parsed as a JSON object, it is returned unchanged.
func (f *FlattenFilter) TransformLine(line string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	flat := make(map[string]interface{})
	flattenMap("", obj, f.sep, flat)
	out, err := json.Marshal(flat)
	if err != nil {
		return line
	}
	return string(out)
}

// flattenMap recursively walks src, building dot-notation keys in dst.
func flattenMap(prefix string, src map[string]interface{}, sep string, dst map[string]interface{}) {
	for k, v := range src {
		key := k
		if prefix != "" {
			key = strings.Join([]string{prefix, k}, sep)
		}
		switch child := v.(type) {
		case map[string]interface{}:
			flattenMap(key, child, sep, dst)
		default:
			dst[key] = v
		}
	}
}
