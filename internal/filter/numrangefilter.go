package filter

import (
	"encoding/json"
	"fmt"
	"math"
)

// NumRangeFilter keeps lines where a numeric field falls within [Min, Max].
type NumRangeFilter struct {
	Field string
	Min   float64
	Max   float64
}

// NewNumRangeFilter creates a NumRangeFilter for the given field and bounds.
// Use math.Inf(-1) / math.Inf(1) for open bounds.
func NewNumRangeFilter(field string, min, max float64) (*NumRangeFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("numrangefilter: field must not be empty")
	}
	if min > max {
		return nil, fmt.Errorf("numrangefilter: min %.4g > max %.4g", min, max)
	}
	return &NumRangeFilter{Field: field, Min: min, Max: max}, nil
}

// MatchesLine returns true when the field exists and its numeric value is
// within [Min, Max]. Lines with missing or non-numeric fields are skipped.
func (f *NumRangeFilter) MatchesLine(line string) bool {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return false
	}
	v, ok := obj[f.Field]
	if !ok {
		return false
	}
	num, ok := toFloat64(v)
	if !ok {
		return false
	}
	return num >= f.Min && num <= f.Max
}

// TransformLine returns the line unchanged.
func (f *NumRangeFilter) TransformLine(line string) string {
	return line
}

// toFloat64 converts a JSON-decoded numeric value to float64.
func toFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case json.Number:
		f, err := n.Float64()
		if err != nil {
			return 0, false
		}
		return f, true
	}
	return 0, false
}

// ParseNumRangeFlag parses a "field:min:max" string into a NumRangeFilter.
// Either bound may be "*" to indicate open (−∞ or +∞).
func ParseNumRangeFlag(s string) (*NumRangeFilter, error) {
	parts := splitN(s, ":", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("numrangefilter: expected field:min:max, got %q", s)
	}
	field, rawMin, rawMax := parts[0], parts[1], parts[2]
	min := math.Inf(-1)
	max := math.Inf(1)
	var err error
	if rawMin != "*" {
		if min, err = parseFloat(rawMin); err != nil {
			return nil, fmt.Errorf("numrangefilter: invalid min %q: %w", rawMin, err)
		}
	}
	if rawMax != "*" {
		if max, err = parseFloat(rawMax); err != nil {
			return nil, fmt.Errorf("numrangefilter: invalid max %q: %w", rawMax, err)
		}
	}
	return NewNumRangeFilter(field, min, max)
}
