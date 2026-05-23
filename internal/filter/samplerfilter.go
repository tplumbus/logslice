package filter

import (
	"fmt"
	"sync/atomic"
)

// SamplerFilter passes every Nth log line, useful for reducing
// high-volume log streams to a representative sample.
type SamplerFilter struct {
	n       uint64
	counter atomic.Uint64
}

// NewSamplerFilter creates a SamplerFilter that passes every nth line.
// n must be >= 1; n=1 passes every line.
func NewSamplerFilter(n int) (*SamplerFilter, error) {
	if n < 1 {
		return nil, fmt.Errorf("sampler: n must be >= 1, got %d", n)
	}
	return &SamplerFilter{n: uint64(n)}, nil
}

// MatchesLine returns true for every nth line encountered.
func (s *SamplerFilter) MatchesLine(line []byte) (bool, error) {
	count := s.counter.Add(1)
	return count%s.n == 1, nil
}

// Reset resets the internal counter, allowing the filter to be reused.
func (s *SamplerFilter) Reset() {
	s.counter.Store(0)
}
