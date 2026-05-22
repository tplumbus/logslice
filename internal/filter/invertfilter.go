package filter

import "fmt"

// InvertFilter wraps another Filter and negates its MatchesLine result.
// Lines that would be excluded by the inner filter are included, and vice versa.
type InvertFilter struct {
	inner Filter
}

// Filter is the interface expected by InvertFilter.
// It matches the MatchesLine signature used by FieldQuery, RegexFilter, etc.
type Filter interface {
	MatchesLine(line string) bool
}

// NewInvertFilter creates an InvertFilter wrapping the provided Filter.
// Returns an error if inner is nil.
func NewInvertFilter(inner Filter) (*InvertFilter, error) {
	if inner == nil {
		return nil, fmt.Errorf("invertfilter: inner filter must not be nil")
	}
	return &InvertFilter{inner: inner}, nil
}

// MatchesLine returns true when the inner filter does NOT match the line.
func (f *InvertFilter) MatchesLine(line string) bool {
	return !f.inner.MatchesLine(line)
}
