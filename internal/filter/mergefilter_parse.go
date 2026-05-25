package filter

import "strings"

// ParseMergeFields splits a comma-separated list of key=value pairs
// and returns them as individual pair strings suitable for NewMergeFilter.
// Blank entries are skipped.
func ParseMergeFields(raw string) []string {
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
