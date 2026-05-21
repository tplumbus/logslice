package filter

import (
	"bufio"
	"encoding/json"
	"io"
	"time"
)

// Pipeline applies a TimeRange and zero or more FieldQueries to a stream of
// JSON log lines, writing matching lines to dst.
type Pipeline struct {
	TimeRange    TimeRange
	FieldQueries []FieldQuery
	TimestampKey string // JSON key to read the timestamp from (default: "time")
}

// NewPipeline creates a Pipeline with sensible defaults.
func NewPipeline(tr TimeRange, queries []FieldQuery) *Pipeline {
	return &Pipeline{
		TimeRange:    tr,
		FieldQueries: queries,
		TimestampKey: "time",
	}
}

// Run reads lines from src and writes matching lines to dst.
// It returns the number of lines written and any read/write error.
func (p *Pipeline) Run(src io.Reader, dst io.Writer) (int, error) {
	scanner := bufio.NewScanner(src)
	written := 0

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		if !p.matchesTimeRange(line) {
			continue
		}

		if !p.matchesFields(line) {
			continue
		}

		if _, err := dst.Write(append(line, '\n')); err != nil {
			return written, err
		}
		written++
	}

	return written, scanner.Err()
}

func (p *Pipeline) matchesTimeRange(line []byte) bool {
	var record map[string]json.RawMessage
	if err := json.Unmarshal(line, &record); err != nil {
		return false
	}
	raw, ok := record[p.TimestampKey]
	if !ok {
		return false
	}
	var tsStr string
	if err := json.Unmarshal(raw, &tsStr); err != nil {
		return false
	}
	ts, err := time.Parse(time.RFC3339Nano, tsStr)
	if err != nil {
		return false
	}
	return p.TimeRange.Contains(ts)
}

func (p *Pipeline) matchesFields(line []byte) bool {
	for _, fq := range p.FieldQueries {
		ok, err := fq.MatchesLine(line)
		if err != nil || !ok {
			return false
		}
	}
	return true
}
