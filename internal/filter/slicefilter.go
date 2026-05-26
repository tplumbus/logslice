package filter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SliceFilter emits only lines whose index (0-based) falls within [start, end).
// A negative end means "until EOF".
type SliceFilter struct {
	start int
	end   int
	cursor int
}

// NewSliceFilter creates a SliceFilter. start must be >= 0; end must be > start or -1.
func NewSliceFilter(start, end int) (*SliceFilter, error) {
	if start < 0 {
		return nil, fmt.Errorf("slicefilter: start must be >= 0, got %d", start)
	}
	if end != -1 && end <= start {
		return nil, fmt.Errorf("slicefilter: end (%d) must be greater than start (%d) or -1", end, start)
	}
	return &SliceFilter{start: start, end: end}, nil
}

// MatchesLine returns true when the current cursor falls within [start, end).
func (f *SliceFilter) MatchesLine(line string) bool {
	idx := f.cursor
	f.cursor++
	if idx < f.start {
		return false
	}
	if f.end != -1 && idx >= f.end {
		return false
	}
	return true
}

// TransformLine returns the line unchanged.
func (f *SliceFilter) TransformLine(line string) (string, error) {
	return line, nil
}

// ParseSliceFlag parses a flag value of the form "start:end" or "start:"
// where end is optional (defaults to -1 meaning open).
func ParseSliceFlag(s string) (int, int, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("slicefilter: expected format start:end, got %q", s)
	}
	var start, end int
	if _, err := fmt.Sscanf(parts[0], "%d", &start); err != nil {
		return 0, 0, fmt.Errorf("slicefilter: invalid start %q: %w", parts[0], err)
	}
	if parts[1] == "" {
		end = -1
	} else {
		if _, err := fmt.Sscanf(parts[1], "%d", &end); err != nil {
			return 0, 0, fmt.Errorf("slicefilter: invalid end %q: %w", parts[1], err)
		}
	}
	_ = json.RawMessage{} // satisfy import if needed
	return start, end, nil
}
