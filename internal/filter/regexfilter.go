package filter

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// RegexFilter matches log lines where a specific JSON field value matches a regular expression.
type RegexFilter struct {
	field string
	re    *regexp.Regexp
}

// NewRegexFilter creates a RegexFilter for the given field and regex pattern.
// Returns an error if the pattern is not a valid regular expression.
func NewRegexFilter(field, pattern string) (*RegexFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("field name must not be empty")
	}
	if pattern == "" {
		return nil, fmt.Errorf("pattern must not be empty")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern %q: %w", pattern, err)
	}
	return &RegexFilter{field: field, re: re}, nil
}

// MatchesLine returns true if the given JSON log line contains the target field
// and its string representation matches the compiled regular expression.
func (f *RegexFilter) MatchesLine(line string) bool {
	var record map[string]interface{}
	if err := json.Unmarshal([]byte(line), &record); err != nil {
		return false
	}
	val, ok := record[f.field]
	if !ok {
		return false
	}
	strVal := fmt.Sprintf("%v", val)
	return f.re.MatchString(strVal)
}

// Field returns the field name this filter operates on.
func (f *RegexFilter) Field() string {
	return f.field
}

// Pattern returns the original regex pattern string.
func (f *RegexFilter) Pattern() string {
	return f.re.String()
}
