package filter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SplitFilter splits a string field into an array by a delimiter.
type SplitFilter struct {
	field string
	delimiter string
	outField string
}

// NewSplitFilter creates a SplitFilter that splits field by delimiter,
// writing the result array to outField (defaults to field if empty).
func NewSplitFilter(field, delimiter, outField string) (*SplitFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("splitfilter: field must not be empty")
	}
	if delimiter == "" {
		return nil, fmt.Errorf("splitfilter: delimiter must not be empty")
	}
	if outField == "" {
		outField = field
	}
	return &SplitFilter{field: field, delimiter: delimiter, outField: outField}, nil
}

// MatchesLine always returns true; SplitFilter is a transform-only filter.
func (f *SplitFilter) MatchesLine(line string) bool {
	return true
}

// TransformLine splits the target field value and writes an array to outField.
func (f *SplitFilter) TransformLine(line string) (string, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, nil
	}

	raw, ok := obj[f.field]
	if !ok {
		return line, nil
	}

	str, ok := raw.(string)
	if !ok {
		return line, nil
	}

	parts := strings.Split(str, f.delimiter)
	obj[f.outField] = parts

	out, err := json.Marshal(obj)
	if err != nil {
		return line, nil
	}
	return string(out), nil
}
