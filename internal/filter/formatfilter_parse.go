package filter

import (
	"fmt"
	"strings"
)

// ParseFormatFlag parses the --format flag value and returns a FormatFilter.
// Accepted values (case-insensitive): json, pretty, text.
func ParseFormatFlag(value string) (*FormatFilter, error) {
	v := strings.TrimSpace(value)
	if v == "" {
		return nil, fmt.Errorf("formatfilter: format flag must not be empty")
	}
	return NewFormatFilter(v)
}
