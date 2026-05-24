package filter

import "strings"

// ParseAddFields splits a comma-separated string of key=value pairs into a
// slice suitable for NewAddFieldsFilter. Blank entries are skipped.
// Example input: "env=prod,region=us-east-1"
func ParseAddFields(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
