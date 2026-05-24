package filter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// AddFieldsFilter adds or overwrites static key=value fields in each log line.
type AddFieldsFilter struct {
	fields map[string]string
}

// NewAddFieldsFilter creates an AddFieldsFilter from a slice of "key=value" pairs.
// Returns an error if any pair is malformed or the key is empty.
func NewAddFieldsFilter(pairs []string) (*AddFieldsFilter, error) {
	if len(pairs) == 0 {
		return nil, fmt.Errorf("addfields: at least one key=value pair required")
	}
	fields := make(map[string]string, len(pairs))
	for _, p := range pairs {
		idx := strings.IndexByte(p, '=')
		if idx <= 0 {
			return nil, fmt.Errorf("addfields: invalid pair %q, expected key=value", p)
		}
		key := p[:idx]
		val := p[idx+1:]
		fields[key] = val
	}
	return &AddFieldsFilter{fields: fields}, nil
}

// MatchesLine always returns true — this filter transforms rather than filters.
func (f *AddFieldsFilter) MatchesLine(line string) bool {
	return true
}

// TransformLine adds or overwrites the configured fields in the JSON object.
// Returns the original line unchanged if it is not valid JSON.
func (f *AddFieldsFilter) TransformLine(line string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for k, v := range f.fields {
		obj[k] = v
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
