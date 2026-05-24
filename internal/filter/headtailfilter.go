package filter

import "fmt"

// HeadTailFilter emits only the first N lines (head) or last N lines (tail).
type HeadTailFilter struct {
	n    int
	mode string // "head" or "tail"
	buf  []string
	seen int
}

// NewHeadTailFilter creates a filter that passes the first or last n lines.
// mode must be "head" or "tail"; n must be positive.
func NewHeadTailFilter(mode string, n int) (*HeadTailFilter, error) {
	if mode != "head" && mode != "tail" {
		return nil, fmt.Errorf("headtailfilter: mode must be \"head\" or \"tail\", got %q", mode)
	}
	if n <= 0 {
		return nil, fmt.Errorf("headtailfilter: n must be positive, got %d", n)
	}
	return &HeadTailFilter{n: n, mode: mode}, nil
}

// MatchesLine returns true for head mode when fewer than n lines have been seen.
// For tail mode it buffers all lines and always returns false during streaming;
// Flush must be called to retrieve the final lines.
func (f *HeadTailFilter) MatchesLine(line string) bool {
	if f.mode == "head" {
		if f.seen < f.n {
			f.seen++
			return true
		}
		return false
	}
	// tail: accumulate in a ring buffer of size n
	if len(f.buf) < f.n {
		f.buf = append(f.buf, line)
	} else {
		f.buf[f.seen%f.n] = line
	}
	f.seen++
	return false
}

// Flush returns buffered lines for tail mode in order, or nil for head mode.
func (f *HeadTailFilter) Flush() []string {
	if f.mode != "tail" || len(f.buf) == 0 {
		return nil
	}
	size := len(f.buf)
	out := make([]string, size)
	if f.seen <= f.n {
		copy(out, f.buf)
	} else {
		start := f.seen % f.n
		for i := 0; i < size; i++ {
			out[i] = f.buf[(start+i)%f.n]
		}
	}
	return out
}
