package filter

import (
	"fmt"
	"strings"
)

// ParseSortOrder converts a string ("asc" or "desc") to a SortOrder constant.
func ParseSortOrder(s string) (SortOrder, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "asc", "":
		return SortAsc, nil
	case "desc":
		return SortDesc, nil
	default:
		return SortAsc, fmt.Errorf("sortfilter: unknown order %q, want \"asc\" or \"desc\"", s)
	}
}

// ParseSortFilter creates a SortFilter from a field name and order string.
// Example: ParseSortFilter("timestamp", "desc")
func ParseSortFilter(field, order string) (*SortFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("sortfilter: field must not be empty")
	}
	so, err := ParseSortOrder(order)
	if err != nil {
		return nil, err
	}
	return NewSortFilter(field, so)
}
