package filter

import (
	"encoding/json"
	"fmt"
)

// CountFilter counts how many lines pass through it, optionally grouped by a field.
type CountFilter struct {
	field  string
	counts map[string]int
	total  int
}

// NewCountFilter creates a CountFilter. If field is empty, only total count is tracked.
func NewCountFilter(field string) (*CountFilter, error) {
	return &CountFilter{
		field:  field,
		counts: make(map[string]int),
	}, nil
}

// MatchesLine always returns true (count filter is pass-through) but records the line.
func (c *CountFilter) MatchesLine(line string) (bool, error) {
	c.total++
	if c.field == "" {
		return true, nil
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		c.counts["<invalid>"]++
		return true, nil
	}
	val, ok := obj[c.field]
	if !ok {
		c.counts["<missing>"]++
		return true, nil
	}
	c.counts[fmt.Sprintf("%v", val)]++
	return true, nil
}

// Total returns the total number of lines observed.
func (c *CountFilter) Total() int {
	return c.total
}

// Counts returns the per-value counts for the tracked field.
// Returns nil if no field was set.
func (c *CountFilter) Counts() map[string]int {
	if c.field == "" {
		return nil
	}
	result := make(map[string]int, len(c.counts))
	for k, v := range c.counts {
		result[k] = v
	}
	return result
}

// Summary returns a human-readable summary string.
func (c *CountFilter) Summary() string {
	if c.field == "" {
		return fmt.Sprintf("total lines: %d", c.total)
	}
	return fmt.Sprintf("total lines: %d, unique values for %q: %d", c.total, c.field, len(c.counts))
}
