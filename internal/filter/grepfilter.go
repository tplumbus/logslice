package filter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// GrepFilter matches log lines where a specific field's string value
// contains the given substring (case-sensitive or case-insensitive).
type GrepFilter struct {
	field     string
	substring string
	ignoreCase bool
}

// NewGrepFilter creates a GrepFilter for the given field and substring.
// If ignoreCase is true, matching is case-insensitive.
func NewGrepFilter(field, substring string, ignoreCase bool) (*GrepFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("grepfilter: field must not be empty")
	}
	if substring == "" {
		return nil, fmt.Errorf("grepfilter: substring must not be empty")
	}
	sub := substring
	if ignoreCase {
		sub = strings.ToLower(sub)
	}
	return &GrepFilter{field: field, substring: sub, ignoreCase: ignoreCase}, nil
}

// MatchesLine returns true if the field value contains the substring.
func (g *GrepFilter) MatchesLine(line string) bool {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return false
	}
	val, ok := obj[g.field]
	if !ok {
		return false
	}
	s := fmt.Sprintf("%v", val)
	if g.ignoreCase {
		s = strings.ToLower(s)
	}
	return strings.Contains(s, g.substring)
}

// TransformLine returns the line unchanged.
func (g *GrepFilter) TransformLine(line string) string {
	return line
}

// ParseGrepFlag parses a grep flag of the form "field:substring" or
// "field:substring:i" (the trailing ":i" enables case-insensitive matching).
func ParseGrepFlag(flag string) (*GrepFilter, error) {
	parts := strings.SplitN(flag, ":", 3)
	if len(parts) < 2 {
		return nil, fmt.Errorf("grepfilter: expected 'field:substring' or 'field:substring:i', got %q", flag)
	}
	ignoreCase := len(parts) == 3 && strings.ToLower(parts[2]) == "i"
	return NewGrepFilter(parts[0], parts[1], ignoreCase)
}
