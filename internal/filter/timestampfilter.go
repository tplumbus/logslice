package filter

import (
	"encoding/json"
	"fmt"
	"time"
)

// TimestampFilter parses a timestamp field and reformats it into a target format.
type TimestampFilter struct {
	field     string
	inFmt     string
	outFmt    string
}

// NewTimestampFilter creates a filter that parses the given field using inFmt
// and rewrites it using outFmt. Both formats use Go reference-time layout.
func NewTimestampFilter(field, inFmt, outFmt string) (*TimestampFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("timestampfilter: field must not be empty")
	}
	if inFmt == "" {
		return nil, fmt.Errorf("timestampfilter: input format must not be empty")
	}
	if outFmt == "" {
		return nil, fmt.Errorf("timestampfilter: output format must not be empty")
	}
	return &TimestampFilter{field: field, inFmt: inFmt, outFmt: outFmt}, nil
}

// MatchesLine always returns true; this filter only transforms.
func (f *TimestampFilter) MatchesLine(line string) bool {
	return true
}

// TransformLine parses the timestamp field and rewrites it in the output format.
// If the field is missing or unparseable the line is returned unchanged.
func (f *TimestampFilter) TransformLine(line string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	raw, ok := obj[f.field]
	if !ok {
		return line
	}
	str, ok := raw.(string)
	if !ok {
		return line
	}
	t, err := time.Parse(f.inFmt, str)
	if err != nil {
		return line
	}
	obj[f.field] = t.Format(f.outFmt)
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
