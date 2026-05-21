package filter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FieldQuery represents a key=value filter on a JSON log field.
type FieldQuery struct {
	Key   string
	Value string
}

// ParseFieldQuery parses a query string in the form "key=value".
func ParseFieldQuery(q string) (*FieldQuery, error) {
	parts := strings.SplitN(q, "=", 2)
	if len(parts) != 2 || parts[0] == "" {
		return nil, fmt.Errorf("invalid field query %q: expected key=value", q)
	}
	return &FieldQuery{Key: parts[0], Value: parts[1]}, nil
}

// MatchesLine returns true if the JSON log line contains the key with the expected value.
func (fq *FieldQuery) MatchesLine(line string) bool {
	var record map[string]interface{}
	if err := json.Unmarshal([]byte(line), &record); err != nil {
		return false
	}
	v, ok := record[fq.Key]
	if !ok {
		return false
	}
	return fmt.Sprintf("%v", v) == fq.Value
}
