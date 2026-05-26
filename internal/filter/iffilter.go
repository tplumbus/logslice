package filter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// IfFilter applies a transform only when a condition filter matches,
// passing the original line through unchanged if the condition does not match.
type IfFilter struct {
	condition Filter
	transform  Filter
}

// Filter is the shared interface implemented by all filters in this package.
// It is defined here for reference; the canonical definition lives in pipeline.go.
type ifFilterIface interface {
	MatchesLine(line string) bool
	TransformLine(line string) string
}

// NewIfFilter creates an IfFilter that applies transform only when condition matches.
func NewIfFilter(condition, transform Filter) (*IfFilter, error) {
	if condition == nil {
		return nil, fmt.Errorf("iffilter: condition must not be nil")
	}
	if transform == nil {
		return nil, fmt.Errorf("iffilter: transform must not be nil")
	}
	return &IfFilter{condition: condition, transform: transform}, nil
}

// MatchesLine always returns true; IfFilter never drops lines.
func (f *IfFilter) MatchesLine(_ string) bool {
	return true
}

// TransformLine applies the transform only when the condition matches.
func (f *IfFilter) TransformLine(line string) string {
	if f.condition.MatchesLine(line) {
		return f.transform.TransformLine(line)
	}
	return line
}

// ParseIfFlag parses a flag of the form "field=value:set_field=new_value".
// The part before ":" is the condition (field query), the part after is the
// add-fields spec applied when the condition matches.
func ParseIfFlag(flag string) (*IfFilter, error) {
	parts := strings.SplitN(flag, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("iffilter: expected 'condition:transform', got %q", flag)
	}
	condStr := strings.TrimSpace(parts[0])
	transStr := strings.TrimSpace(parts[1])
	if condStr == "" {
		return nil, fmt.Errorf("iffilter: condition part is empty")
	}
	if transStr == "" {
		return nil, fmt.Errorf("iffilter: transform part is empty")
	}

	cond, err := ParseFieldQuery(condStr)
	if err != nil {
		return nil, fmt.Errorf("iffilter: invalid condition: %w", err)
	}

	pairs, err := ParseAddFields([]string{transStr})
	if err != nil {
		return nil, fmt.Errorf("iffilter: invalid transform: %w", err)
	}
	trans, err := NewAddFieldsFilter(pairs)
	if err != nil {
		return nil, fmt.Errorf("iffilter: invalid transform: %w", err)
	}

	return NewIfFilter(cond, trans)
}

// jsonRoundTrip is a helper used in tests; kept here for package-level reuse.
func jsonRoundTrip(line string) (map[string]interface{}, bool) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return nil, false
	}
	return m, true
}
