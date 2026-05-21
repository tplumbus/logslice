package filter

import (
	"encoding/json"
	"time"
)

// TimeRangeStats tracks statistics about log lines evaluated against a TimeRange.
type TimeRangeStats struct {
	Total    int
	Matched  int
	Skipped  int
	Earliest *time.Time
	Latest   *time.Time
}

// EvaluateLine parses the "time" field from a raw JSON log line, checks it
// against the TimeRange, and updates the stats. It returns true if the line
// falls within the range (or if the range is unbounded).
func (s *TimeRangeStats) EvaluateLine(tr TimeRange, line []byte) bool {
	s.Total++

	var entry struct {
		Time time.Time `json:"time"`
	}

	if err := json.Unmarshal(line, &entry); err != nil || entry.Time.IsZero() {
		// Lines without a parseable time field are always passed through.
		s.Matched++
		return true
	}

	t := entry.Time

	if s.Earliest == nil || t.Before(*s.Earliest) {
		copy := t
		s.Earliest = &copy
	}
	if s.Latest == nil || t.After(*s.Latest) {
		copy := t
		s.Latest = &copy
	}

	if !tr.Contains(t) {
		s.Skipped++
		return false
	}

	s.Matched++
	return true
}

// MatchRatio returns the fraction of total lines that matched, or 1.0 if no
// lines were evaluated.
func (s *TimeRangeStats) MatchRatio() float64 {
	if s.Total == 0 {
		return 1.0
	}
	return float64(s.Matched) / float64(s.Total)
}
