package filter

import (
	"encoding/json"
	"fmt"
	"sort"
)

// SortOrder defines ascending or descending sort direction.
type SortOrder int

const (
	SortAsc  SortOrder = iota
	SortDesc SortOrder = iota
)

// SortFilter buffers all matching lines and re-emits them sorted by a JSON field.
type SortFilter struct {
	field string
	order SortOrder
	buf   []scoredLine
}

type scoredLine struct {
	raw   string
	value string
}

// NewSortFilter returns a SortFilter that sorts buffered lines by field in the given order.
func NewSortFilter(field string, order SortOrder) (*SortFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("sortfilter: field must not be empty")
	}
	return &SortFilter{field: field, order: order}, nil
}

// Observe buffers a line for deferred sorting.
func (s *SortFilter) Observe(line string) {
	var m map[string]interface{}
	val := ""
	if err := json.Unmarshal([]byte(line), &m); err == nil {
		if v, ok := m[s.field]; ok {
			val = fmt.Sprintf("%v", v)
		}
	}
	s.buf = append(s.buf, scoredLine{raw: line, value: val})
}

// Sorted returns lines in the requested sort order.
func (s *SortFilter) Sorted() []string {
	sorted := make([]scoredLine, len(s.buf))
	copy(sorted, s.buf)
	sort.SliceStable(sorted, func(i, j int) bool {
		if s.order == SortAsc {
			return sorted[i].value < sorted[j].value
		}
		return sorted[i].value > sorted[j].value
	})
	out := make([]string, len(sorted))
	for i, sl := range sorted {
		out[i] = sl.raw
	}
	return out
}

// Reset clears the internal buffer.
func (s *SortFilter) Reset() {
	s.buf = s.buf[:0]
}
