package filter

import (
	"encoding/json"
	"fmt"
)

// ExistsFilter passes lines where the specified field either exists or does not
// exist in the JSON object, depending on the negate flag.
type ExistsFilter struct {
	field  string
	negate bool // if true, passes lines where field is ABSENT
}

// NewExistsFilter creates a filter that passes lines containing the given field.
// If negate is true, it passes lines where the field is absent instead.
func NewExistsFilter(field string, negate bool) (*ExistsFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("existsfilter: field name must not be empty")
	}
	return &ExistsFilter{field: field, negate: negate}, nil
}

// MatchesLine returns true when the line satisfies the existence condition.
func (f *ExistsFilter) MatchesLine(line string) bool {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return false
	}
	_, exists := obj[f.field]
	if f.negate {
		return !exists
	}
	return exists
}

// TransformLine returns the line unchanged.
func (f *ExistsFilter) TransformLine(line string) string {
	return line
}

// ParseExistsFlag parses a field name from the --exists or --not-exists flag.
// Returns the field name and a negate bool.
func ParseExistsFlag(value string, negate bool) (*ExistsFilter, error) {
	if value == "" {
		return nil, fmt.Errorf("existsfilter: field name must not be empty")
	}
	return NewExistsFilter(value, negate)
}
