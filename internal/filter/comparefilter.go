package filter

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// CompareOp represents a comparison operator.
type CompareOp string

const (
	OpEq  CompareOp = "eq"
	OpNe  CompareOp = "ne"
	OpLt  CompareOp = "lt"
	OpLte CompareOp = "lte"
	OpGt  CompareOp = "gt"
	OpGte CompareOp = "gte"
)

// CompareFilter matches lines where a numeric or string field satisfies an operator/value comparison.
type CompareFilter struct {
	field string
	op    CompareOp
	value float64
}

// NewCompareFilter creates a CompareFilter for the given field, operator, and threshold.
func NewCompareFilter(field string, op CompareOp, value float64) (*CompareFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("comparefilter: field must not be empty")
	}
	switch op {
	case OpEq, OpNe, OpLt, OpLte, OpGt, OpGte:
		// valid
	default:
		return nil, fmt.Errorf("comparefilter: unknown operator %q", op)
	}
	return &CompareFilter{field: field, op: op, value: value}, nil
}

// MatchesLine returns true if the field in the JSON line satisfies the comparison.
func (f *CompareFilter) MatchesLine(line string) bool {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return false
	}
	v, ok := obj[f.field]
	if !ok {
		return false
	}
	var num float64
	switch val := v.(type) {
	case float64:
		num = val
	case string:
		parsed, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return false
		}
		num = parsed
	default:
		return false
	}
	switch f.op {
	case OpEq:
		return num == f.value
	case OpNe:
		return num != f.value
	case OpLt:
		return num < f.value
	case OpLte:
		return num <= f.value
	case OpGt:
		return num > f.value
	case OpGte:
		return num >= f.value
	}
	return false
}

// TransformLine returns the line unchanged.
func (f *CompareFilter) TransformLine(line string) string {
	return line
}
