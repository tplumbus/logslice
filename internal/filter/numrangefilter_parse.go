package filter

import (
	"fmt"
	"strconv"
	"strings"
)

// splitN splits s by sep into exactly n parts, returning an error-safe slice.
func splitN(s, sep string, n int) []string {
	return strings.SplitN(s, sep, n)
}

// parseFloat wraps strconv.ParseFloat with a friendlier error.
func parseFloat(s string) (float64, error) {
	v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0, fmt.Errorf("not a valid number: %q", s)
	}
	return v, nil
}
