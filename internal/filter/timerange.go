package filter

import (
	"fmt"
	"time"
)

// TimeRange represents an inclusive time window for log filtering.
type TimeRange struct {
	From time.Time
	To   time.Time
}

// ParseTimeRange parses two RFC3339 timestamp strings into a TimeRange.
// Either bound may be empty to indicate an open-ended range.
func ParseTimeRange(from, to string) (TimeRange, error) {
	var tr TimeRange

	if from != "" {
		t, err := time.Parse(time.RFC3339, from)
		if err != nil {
			return tr, fmt.Errorf("invalid --from timestamp %q: %w", from, err)
		}
		tr.From = t
	}

	if to != "" {
		t, err := time.Parse(time.RFC3339, to)
		if err != nil {
			return tr, fmt.Errorf("invalid --to timestamp %q: %w", to, err)
		}
		tr.To = t
	}

	if !tr.From.IsZero() && !tr.To.IsZero() && tr.To.Before(tr.From) {
		return tr, fmt.Errorf("--to (%s) must not be before --from (%s)", to, from)
	}

	return tr, nil
}

// Contains reports whether the given timestamp falls within the time range.
// A zero From or To bound is treated as unbounded.
func (tr TimeRange) Contains(ts time.Time) bool {
	if !tr.From.IsZero() && ts.Before(tr.From) {
		return false
	}
	if !tr.To.IsZero() && ts.After(tr.To) {
		return false
	}
	return true
}

// IsZero reports whether neither bound has been set.
func (tr TimeRange) IsZero() bool {
	return tr.From.IsZero() && tr.To.IsZero()
}
