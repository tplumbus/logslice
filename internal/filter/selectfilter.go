package filter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SelectFilter emits only the specified fields from each JSON log line.
type SelectFilter struct {
	fields []string
}

// NewSelectFilter creates a SelectFilter that retains only the given fields.
// Returns an error if no fields are provided.
func NewSelectFilter(fields []string) (*SelectFilter, error) {
	var kept []string
	for _, f := range fields {
		f = strings.TrimSpace(f)
		if f != "" {
			kept = append(kept, f)
		}
	}
	if len(kept) == 0 {
		return nil, fmt.Errorf("selectfilter: at least one field is required")
	}
	return &SelectFilter{fields: kept}, nil
}

// MatchesLine always returns true; SelectFilter only transforms output.
func (s *SelectFilter) MatchesLine(line string) bool {
	return true
}

// TransformLine returns a new JSON object containing only the selected fields.
// If a field is missing from the source object it is omitted from the output.
// Returns the original line if it cannot be parsed.
func (s *SelectFilter) TransformLine(line string) string {
	var src map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &src); err != nil {
		return line
	}

	out := make(map[string]json.RawMessage, len(s.fields))
	for _, f := range s.fields {
		if v, ok := src[f]; ok {
			out[f] = v
		}
	}

	b, err := json.Marshal(out)
	if err != nil {
		return line
	}
	return string(b)
}
