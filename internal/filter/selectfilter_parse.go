package filter

import (
	"fmt"
	"strings"
)

// ParseSelectFields splits a comma-separated list of field names and
// returns a SelectFilter. Whitespace around names is trimmed.
// Returns an error if the input string is empty or contains only blanks.
func ParseSelectFields(raw string) (*SelectFilter, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, fmt.Errorf("select: field list must not be empty")
	}
	parts := strings.Split(raw, ",")
	return NewSelectFilter(parts)
}
