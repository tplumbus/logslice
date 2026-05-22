package filter

import (
	"encoding/json"
	"fmt"
)

// LineMatcher is satisfied by any filter that can evaluate a single JSON log line.
type LineMatcher interface {
	MatchesLine(line string) bool
}

// Pipeline combines a TimeRange gate with an optional LineMatcher to filter
// structured JSON log lines.
type Pipeline struct {
	tr      *TimeRange
	filter  LineMatcher
}

// NewPipeline creates a Pipeline. filter may be nil (time-range only).
func NewPipeline(tr *TimeRange, filter LineMatcher) *Pipeline {
	return &Pipeline{tr: tr, filter: filter}
}

// Run iterates over lines, applies the time-range and field filter, and calls
// emit for every line that passes both gates. Processing stops on the first
// error returned by emit.
func (p *Pipeline) Run(lines []string, emit func(string) error) error {
	for _, line := range lines {
		if line == "" {
			continue
		}

		var raw map[string]interface{}
		if err := json.Unmarshal([]byte(line), &raw); err != nil {
			// skip non-JSON lines silently
			continue
		}

		if p.tr != nil {
			ts, ok := extractTimestamp(raw)
			if !ok {
				continue
			}
			if !p.tr.Contains(ts) {
				continue
			}
		}

		if p.filter != nil && !p.filter.MatchesLine(line) {
			continue
		}

		if err := emit(line); err != nil {
			return fmt.Errorf("pipeline emit: %w", err)
		}
	}
	return nil
}
