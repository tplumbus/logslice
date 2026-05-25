package filter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// TrimFilter removes leading and/or trailing whitespace from string field values.
type TrimFilter struct {
	fields []string
	mode   string // "both", "left", "right"
}

// NewTrimFilter creates a TrimFilter that trims the given fields.
// mode must be one of: "both", "left", "right".
func NewTrimFilter(fields []string, mode string) (*TrimFilter, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("trimfilter: at least one field is required")
	}
	switch mode {
	case "both", "left", "right":
	default:
		return nil, fmt.Errorf("trimfilter: invalid mode %q, must be both/left/right", mode)
	}
	clean := make([]string, 0, len(fields))
	for _, f := range fields {
		if strings.TrimSpace(f) != "" {
			clean = append(clean, f)
		}
	}
	if len(clean) == 0 {
		return nil, fmt.Errorf("trimfilter: no valid fields provided")
	}
	return &TrimFilter{fields: clean, mode: mode}, nil
}

// MatchesLine always returns true; TrimFilter is a transform-only filter.
func (f *TrimFilter) MatchesLine(line string) bool {
	return true
}

// TransformLine trims whitespace from the configured string fields.
func (f *TrimFilter) TransformLine(line string) (string, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, nil
	}
	for _, field := range f.fields {
		val, ok := obj[field]
		if !ok {
			continue
		}
		s, ok := val.(string)
		if !ok {
			continue
		}
		switch f.mode {
		case "both":
			obj[field] = strings.TrimSpace(s)
		case "left":
			obj[field] = strings.TrimLeft(s, " \t\n\r")
		case "right":
			obj[field] = strings.TrimRight(s, " \t\n\r")
		}
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line, nil
	}
	return string(out), nil
}

// ParseTrimFlag parses a trim flag value of the form "field1,field2:mode".
func ParseTrimFlag(flag string) (*TrimFilter, error) {
	parts := strings.SplitN(flag, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("trimfilter: expected format 'field1,field2:mode', got %q", flag)
	}
	fields := strings.Split(parts[0], ",")
	mode := strings.TrimSpace(parts[1])
	return NewTrimFilter(fields, mode)
}
