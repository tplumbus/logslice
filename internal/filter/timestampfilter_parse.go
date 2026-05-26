package filter

import (
	"fmt"
	"strings"
)

// Well-known shorthand aliases for common time layouts.
var timestampAliases = map[string]string{
	"rfc3339":  "2006-01-02T15:04:05Z07:00",
	"rfc3339ms": "2006-01-02T15:04:05.000Z07:00",
	"unix-date": "Mon Jan _2 15:04:05 MST 2006",
	"datetime":  "2006-01-02 15:04:05",
	"date":      "2006-01-02",
}

// resolveAlias returns the Go time layout for a shorthand alias or the
// original string if no alias is found.
func resolveAlias(s string) string {
	if v, ok := timestampAliases[strings.ToLower(s)]; ok {
		return v
	}
	return s
}

// ParseTimestampFlag parses a flag value of the form
// "field:in-format:out-format" and returns a configured TimestampFilter.
// Both in-format and out-format may use shorthand aliases defined above.
func ParseTimestampFlag(value string) (*TimestampFilter, error) {
	parts := strings.SplitN(value, ":", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("timestampfilter: expected field:in-format:out-format, got %q", value)
	}
	field := strings.TrimSpace(parts[0])
	inFmt := resolveAlias(strings.TrimSpace(parts[1]))
	outFmt := resolveAlias(strings.TrimSpace(parts[2]))
	return NewTimestampFilter(field, inFmt, outFmt)
}
