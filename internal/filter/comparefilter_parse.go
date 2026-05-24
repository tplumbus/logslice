package filter

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseCompareFlag parses a flag value in the form "field:op:value".
// Example: "status_code:gte:400"
func ParseCompareFlag(s string) (*CompareFilter, error) {
	parts := strings.SplitN(s, ":", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("comparefilter: invalid format %q, expected field:op:value", s)
	}
	field := strings.TrimSpace(parts[0])
	op := CompareOp(strings.TrimSpace(parts[1]))
	valStr := strings.TrimSpace(parts[2])

	if field == "" {
		return nil, fmt.Errorf("comparefilter: field must not be empty")
	}
	value, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return nil, fmt.Errorf("comparefilter: invalid value %q: %w", valStr, err)
	}
	return NewCompareFilter(field, op, value)
}
