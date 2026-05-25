package filter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// MergeFilter merges fields from a static JSON object into each log line.
// Existing fields in the log line are preserved unless overwrite is true.
type MergeFilter struct {
	fields    map[string]interface{}
	overwrite bool
}

// NewMergeFilter creates a MergeFilter that merges the given key=value pairs
// into each log line. If overwrite is true, existing fields are replaced.
func NewMergeFilter(pairs []string, overwrite bool) (*MergeFilter, error) {
	if len(pairs) == 0 {
		return nil, fmt.Errorf("mergefilter: at least one field pair required")
	}
	fields := make(map[string]interface{}, len(pairs))
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			return nil, fmt.Errorf("mergefilter: invalid pair %q, expected key=value", p)
		}
		fields[parts[0]] = parts[1]
	}
	return &MergeFilter{fields: fields, overwrite: overwrite}, nil
}

// MatchesLine always returns true; MergeFilter is a transform-only filter.
func (f *MergeFilter) MatchesLine(line string) bool {
	return true
}

// TransformLine merges the configured fields into the JSON log line.
func (f *MergeFilter) TransformLine(line string) (string, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, nil
	}
	for k, v := range f.fields {
		if _, exists := obj[k]; !exists || f.overwrite {
			obj[k] = v
		}
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line, fmt.Errorf("mergefilter: marshal error: %w", err)
	}
	return string(out), nil
}
