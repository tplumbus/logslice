package filter

import (
	"encoding/json"
	"io"
	"sort"
)

// FieldStats tracks value frequency for a specific JSON field across log lines.
type FieldStats struct {
	Field  string
	Counts map[string]int
	Total  int
	Missed int
}

// NewFieldStats creates a FieldStats collector for the given field name.
func NewFieldStats(field string) *FieldStats {
	return &FieldStats{
		Field:  field,
		Counts: make(map[string]int),
	}
}

// Observe parses a JSON log line and records the value of the tracked field.
func (fs *FieldStats) Observe(line string) {
	fs.Total++
	var record map[string]interface{}
	if err := json.Unmarshal([]byte(line), &record); err != nil {
		fs.Missed++
		return
	}
	v, ok := record[fs.Field]
	if !ok {
		fs.Missed++
		return
	}
	key := formatValue(v)
	fs.Counts[key]++
}

// TopN returns the top n field values sorted by descending frequency.
// If n <= 0 all entries are returned.
func (fs *FieldStats) TopN(n int) []FieldCount {
	result := make([]FieldCount, 0, len(fs.Counts))
	for k, v := range fs.Counts {
		result = append(result, FieldCount{Value: k, Count: v})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count != result[j].Count {
			return result[i].Count > result[j].Count
		}
		return result[i].Value < result[j].Value
	})
	if n > 0 && n < len(result) {
		return result[:n]
	}
	return result
}

// WriteSummary writes a human-readable summary to the provided writer.
func (fs *FieldStats) WriteSummary(w io.Writer) {
	io.WriteString(w, "field: "+fs.Field+"\n")
	for _, fc := range fs.TopN(0) {
		io.WriteString(w, "  "+fc.Value+": "+itoa(fc.Count)+"\n")
	}
	io.WriteString(w, "total lines: "+itoa(fs.Total)+", missing field: "+itoa(fs.Missed)+"\n")
}

// FieldCount holds a field value and its occurrence count.
type FieldCount struct {
	Value string
	Count int
}

func formatValue(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}
