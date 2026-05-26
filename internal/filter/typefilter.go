package filter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// TypeFilter filters log lines based on the JSON type of a field value.
// Supported types: string, number, bool, null, array, object
type TypeFilter struct {
	field    string
	wantType string
}

// NewTypeFilter creates a TypeFilter that passes lines where the named field
// has the specified JSON type.
func NewTypeFilter(field, typeName string) (*TypeFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("typefilter: field must not be empty")
	}
	valid := map[string]bool{
		"string": true, "number": true, "bool": true,
		"null": true, "array": true, "object": true,
	}
	t := strings.ToLower(strings.TrimSpace(typeName))
	if !valid[t] {
		return nil, fmt.Errorf("typefilter: unknown type %q (want string|number|bool|null|array|object)", typeName)
	}
	return &TypeFilter{field: field, wantType: t}, nil
}

func (f *TypeFilter) MatchesLine(line string) bool {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return false
	}
	raw, ok := obj[f.field]
	if !ok {
		return false
	}
	return jsonType(raw) == f.wantType
}

func (f *TypeFilter) TransformLine(line string) string { return line }

// jsonType returns the JSON type name of a raw JSON value.
func jsonType(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	switch raw[0] {
	case '"':
		return "string"
	case 't', 'f':
		return "bool"
	case 'n':
		return "null"
	case '[':
		return "array"
	case '{':
		return "object"
	default:
		return "number"
	}
}

// ParseTypeFlag parses a flag value of the form "field:type".
func ParseTypeFlag(s string) (*TypeFilter, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("typefilter: expected field:type, got %q", s)
	}
	return NewTypeFilter(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
}
