package filter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FormatFilter rewrites each JSON log line into a different output format.
// Supported formats: "json" (compact), "pretty" (indented JSON), "text" (key=value pairs).
type FormatFilter struct {
	format string
}

// NewFormatFilter creates a FormatFilter for the given format string.
func NewFormatFilter(format string) (*FormatFilter, error) {
	format = strings.ToLower(strings.TrimSpace(format))
	switch format {
	case "json", "pretty", "text":
		return &FormatFilter{format: format}, nil
	default:
		return nil, fmt.Errorf("formatfilter: unsupported format %q (want json|pretty|text)", format)
	}
}

// MatchesLine always returns true; FormatFilter is a transform-only filter.
func (f *FormatFilter) MatchesLine(line string) bool {
	return true
}

// TransformLine reformats the JSON line according to the chosen format.
func (f *FormatFilter) TransformLine(line string) (string, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, nil // pass through non-JSON lines unchanged
	}

	switch f.format {
	case "pretty":
		b, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			return line, nil
		}
		return string(b), nil
	case "text":
		parts := make([]string, 0, len(obj))
		for k, v := range obj {
			parts = append(parts, fmt.Sprintf("%s=%v", k, v))
		}
		return strings.Join(parts, " "), nil
	default: // "json" — compact
		b, err := json.Marshal(obj)
		if err != nil {
			return line, nil
		}
		return string(b), nil
	}
}
