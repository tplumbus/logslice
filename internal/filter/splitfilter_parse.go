package filter

import (
	"fmt"
	"strings"
)

// ParseSplitFlag parses a flag value of the form "field:delimiter" or
// "field:delimiter:outfield" and returns a SplitFilter.
// The delimiter may be "\t" or "\n" for tab/newline.
func ParseSplitFlag(value string) (*SplitFilter, error) {
	parts := strings.SplitN(value, ":", 3)
	if len(parts) < 2 {
		return nil, fmt.Errorf("splitfilter: expected 'field:delimiter' or 'field:delimiter:outfield', got %q", value)
	}

	field := strings.TrimSpace(parts[0])
	delimiter := parts[1]
	outField := ""
	if len(parts) == 3 {
		outField = strings.TrimSpace(parts[2])
	}

	// Unescape common escape sequences
	delimiter = strings.ReplaceAll(delimiter, `\t`, "\t")
	delimiter = strings.ReplaceAll(delimiter, `\n`, "\n")

	return NewSplitFilter(field, delimiter, outField)
}
