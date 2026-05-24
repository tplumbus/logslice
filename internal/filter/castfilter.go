package filter

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// CastFilter converts a field's value to a target type (string, int, float, bool).
type CastFilter struct {
	field      string
	targetType string
}

// NewCastFilter creates a CastFilter that casts the given field to targetType.
// Valid targetType values: "string", "int", "float", "bool".
func NewCastFilter(field, targetType string) (*CastFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("castfilter: field must not be empty")
	}
	switch targetType {
	case "string", "int", "float", "bool":
	default:
		return nil, fmt.Errorf("castfilter: unsupported target type %q", targetType)
	}
	return &CastFilter{field: field, targetType: targetType}, nil
}

// MatchesLine always returns true; CastFilter is a transform-only filter.
func (f *CastFilter) MatchesLine(line string) bool { return true }

// TransformLine parses the JSON line, casts the field value, and re-encodes it.
func (f *CastFilter) TransformLine(line string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	raw, ok := obj[f.field]
	if !ok {
		return line
	}
	casted, err := castValue(raw, f.targetType)
	if err != nil {
		return line
	}
	obj[f.field] = casted
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}

func castValue(v interface{}, targetType string) (interface{}, error) {
	str := fmt.Sprintf("%v", v)
	switch targetType {
	case "string":
		return str, nil
	case "int":
		return strconv.ParseInt(str, 10, 64)
	case "float":
		return strconv.ParseFloat(str, 64)
	case "bool":
		return strconv.ParseBool(str)
	}
	return nil, fmt.Errorf("unsupported type: %s", targetType)
}
