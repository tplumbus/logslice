package filter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FieldQuery represents a single field=value filter condition.
type FieldQuery struct {
	Field string
	Value string
}

// ParseFieldQuery parses a query string of the form "field=value".
func ParseFieldQuery(q string) (FieldQuery, error) {
	parts := strings.SplitN(q, "=", 2)
	if len(parts) != 2 {
		return FieldQuery{}, fmt.Errorf("invalid field query %q: expected field=value format", q)
	}
	field := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	if field == "" {
		return FieldQuery{}, fmt.Errorf("invalid field query %q: field name must not be empty", q)
	}
	return FieldQuery{Field: field, Value: value}, nil
}

// MatchesLine reports whether the JSON log line satisfies the field query.
// It performs a string comparison against the field's JSON value.
func (fq FieldQuery) MatchesLine(line []byte) (bool, error) {
	var record map[string]json.RawMessage
	if err := json.Unmarshal(line, &record); err != nil {
		return false, fmt.Errorf("failed to parse JSON: %w", err)
	}

	raw, ok := record[fq.Field]
	if !ok {
		return false, nil
	}

	// Compare as unquoted string or raw value.
	var strVal string
	if err := json.Unmarshal(raw, &strVal); err == nil {
		return strVal == fq.Value, nil
	}

	// Fall back to raw JSON comparison (e.g. numbers, booleans).
	return strings.TrimSpace(string(raw)) == fq.Value, nil
}
