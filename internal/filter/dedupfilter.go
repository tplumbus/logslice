package filter

import (
	"encoding/json"
	"fmt"
	"sync"
)

// DedupFilter eliminates duplicate log lines based on the value of a specific
// JSON field. Once a value has been seen, subsequent lines with the same value
// for that field are dropped. The filter is safe for concurrent use.
type DedupFilter struct {
	field string
	seen  map[string]struct{}
	mu    sync.Mutex
}

// NewDedupFilter creates a DedupFilter that deduplicates lines by the given
// field name. Returns an error if field is empty.
func NewDedupFilter(field string) (*DedupFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("dedupfilter: field name must not be empty")
	}
	return &DedupFilter{
		field: field,
		seen:  make(map[string]struct{}),
	}, nil
}

// MatchesLine returns true if the line has not been seen before (based on the
// configured field value). Lines with missing or unparseable fields are passed
// through once — their raw bytes are used as the dedup key.
func (d *DedupFilter) MatchesLine(line []byte) bool {
	key := d.extractKey(line)

	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.seen[key]; exists {
		return false
	}
	d.seen[key] = struct{}{}
	return true
}

// Reset clears all previously seen values, allowing the filter to start fresh.
func (d *DedupFilter) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]struct{})
}

// SeenCount returns the number of distinct values observed so far.
func (d *DedupFilter) SeenCount() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.seen)
}

// extractKey returns a string key for deduplication. It tries to parse the
// line as JSON and extract the target field; on failure it falls back to the
// raw line content.
func (d *DedupFilter) extractKey(line []byte) string {
	var record map[string]interface{}
	if err := json.Unmarshal(line, &record); err != nil {
		// Fallback: treat the whole raw line as the key.
		return string(line)
	}

	val, ok := record[d.field]
	if !ok {
		// Field absent — use the raw line so absent-field lines are not all
		// collapsed into a single bucket.
		return string(line)
	}

	return fmt.Sprintf("%v", val)
}
