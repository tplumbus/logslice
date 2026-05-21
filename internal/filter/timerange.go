package filter

import (
	"fmt"
	"time"
)

// TimeRange represents an optional start/end window for log filtering.
type TimeRange struct {
	From *time.Time
	To   *time.Time
}

// ParseTimeRange parses optional from/to RFC3339 timestamp strings into a TimeRange.
// Either or both may be empty strings to indicate open bounds.
func ParseTimeRange(from, to string) (TimeRange, error) {
	var tr TimeRange

	if from != "" {
		t, err := time.Parse(time.RFC3339, from)
		if err != nil {
			return tr, fmt.Errorf("invalid --from timestamp %q: %w", from, err)
		}
		tr.From = &t
	}

	if to != "" {
		t, err := time.Parse(time.RFC3339, to)
		if err != nil {
			return tr, fmt.Errorf("invalid --to timestamp %q: %w", to, err)
		}
		tr.To = &t
	}

	if tr.From != nil && tr.To != nil && tr.To.Before(*tr.From) {
		return tr, fmt.Errorf("--to (%s) must not be before --from (%s)", to, from)
	}

	return tr, nil
}

// Contains reports whether the given time falls within the TimeRange.
// Open bounds are treated as unbounded.
func (tr TimeRange) Contains(t time.Time) bool {
	if tr.From != nil && t.Before(*tr.From) {
		return false
	}
	if tr.To != nil && t.After(*tr.To) {
		return false
	}
	return true
}

// IsZero reports whether the TimeRange has no bounds set.
func (tr TimeRange) IsZero() bool {
	return tr.From == nil && tr.To == nil
}
