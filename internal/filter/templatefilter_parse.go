package filter

import (
	"fmt"
	"strings"
)

// ParseTemplateFlag parses the --template flag value.
// It supports an optional "file:" prefix to load a template from disk,
// but for simplicity in this implementation the raw string is used directly.
// Returns a compiled *TemplateFilter or an error.
func ParseTemplateFlag(value string) (*TemplateFilter, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, fmt.Errorf("--template: value must not be empty")
	}
	// Unescape common escape sequences supplied on the command line.
	value = strings.ReplaceAll(value, `\n`, "\n")
	value = strings.ReplaceAll(value, `\t`, "\t")
	return NewTemplateFilter(value)
}
