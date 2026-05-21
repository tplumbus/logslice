package filter

import (
	"fmt"
	"strings"
)

// MultiFieldFilter holds multiple FieldQuery filters and matches lines
// that satisfy ALL of the specified field conditions (AND semantics).
type MultiFieldFilter struct {
	queries []*FieldQuery
}

// NewMultiFieldFilter parses a slice of "key=value" strings into a MultiFieldFilter.
// Returns an error if any query string is invalid.
func NewMultiFieldFilter(rawQueries []string) (*MultiFieldFilter, error) {
	queries := make([]*FieldQuery, 0, len(rawQueries))
	for _, raw := range rawQueries {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		fq, err := ParseFieldQuery(raw)
		if err != nil {
			return nil, fmt.Errorf("multifield: %w", err)
		}
		queries = append(queries, fq)
	}
	return &MultiFieldFilter{queries: queries}, nil
}

// MatchesLine returns true only when the line satisfies every FieldQuery.
// An empty filter (no queries) always returns true.
func (mf *MultiFieldFilter) MatchesLine(line string) bool {
	for _, q := range mf.queries {
		if !q.MatchesLine(line) {
			return false
		}
	}
	return true
}

// Len returns the number of field queries in the filter.
func (mf *MultiFieldFilter) Len() int {
	return len(mf.queries)
}
