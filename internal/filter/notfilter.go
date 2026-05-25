package filter

import (
	"fmt"
	"strings"
)

// NotFilter wraps multiple filters and passes a line only if NONE of them match.
// It is the logical NOR of all inner filters.
type NotFilter struct {
	filters []Filter
}

// NewNotFilter creates a NotFilter that rejects lines matching any of the given filters.
// At least one inner filter must be provided.
func NewNotFilter(filters ...Filter) (*NotFilter, error) {
	if len(filters) == 0 {
		return nil, fmt.Errorf("notfilter: at least one inner filter is required")
	}
	for i, f := range filters {
		if f == nil {
			return nil, fmt.Errorf("notfilter: filter at index %d is nil", i)
		}
	}
	return &NotFilter{filters: filters}, nil
}

// MatchesLine returns true only when none of the inner filters match the line.
func (n *NotFilter) MatchesLine(line string) bool {
	for _, f := range n.filters {
		if f.MatchesLine(line) {
			return false
		}
	}
	return true
}

// ParseNotFlag parses a comma-separated list of field=value pairs and builds
// a NotFilter that rejects lines matching any of those field queries.
// Example: "level=debug,env=staging"
func ParseNotFlag(raw string) (*NotFilter, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("notfilter: empty flag value")
	}
	parts := strings.Split(raw, ",")
	filters := make([]Filter, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		fq, err := ParseFieldQuery(part)
		if err != nil {
			return nil, fmt.Errorf("notfilter: invalid query %q: %w", part, err)
		}
		filters = append(filters, fq)
	}
	if len(filters) == 0 {
		return nil, fmt.Errorf("notfilter: no valid queries parsed from %q", raw)
	}
	return NewNotFilter(filters...)
}
