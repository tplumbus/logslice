package filter

import "fmt"

// LimitFilter stops the pipeline after emitting at most N matching lines.
type LimitFilter struct {
	max     int
	emitted int
}

// NewLimitFilter creates a LimitFilter that passes at most max lines.
// max must be greater than zero.
func NewLimitFilter(max int) (*LimitFilter, error) {
	if max <= 0 {
		return nil, fmt.Errorf("limit must be greater than zero, got %d", max)
	}
	return &LimitFilter{max: max}, nil
}

// MatchesLine returns true while the number of emitted lines is below the limit.
// Once the limit is reached it always returns false, causing the pipeline to
// stop forwarding lines.
func (f *LimitFilter) MatchesLine(line string) bool {
	if f.emitted >= f.max {
		return false
	}
	f.emitted++
	return true
}

// Reset resets the internal counter so the filter can be reused.
func (f *LimitFilter) Reset() {
	f.emitted = 0
}
