package filter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// RenameFilter renames one or more fields in each JSON log line.
// Pairs are specified as "oldName=newName".
type RenameFilter struct {
	mappings map[string]string
}

// NewRenameFilter constructs a RenameFilter from a slice of "old=new" pairs.
func NewRenameFilter(pairs []string) (*RenameFilter, error) {
	if len(pairs) == 0 {
		return nil, fmt.Errorf("renamefilter: at least one field mapping is required")
	}
	mappings := make(map[string]string, len(pairs))
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("renamefilter: invalid mapping %q, expected old=new", p)
		}
		mappings[parts[0]] = parts[1]
	}
	return &RenameFilter{mappings: mappings}, nil
}

// MatchesLine always returns true; RenameFilter only transforms.
func (r *RenameFilter) MatchesLine(line string) bool {
	return true
}

// TransformLine renames fields in the JSON object according to the configured mappings.
// Fields not present in the object are silently skipped.
func (r *RenameFilter) TransformLine(line string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for oldKey, newKey := range r.mappings {
		if val, ok := obj[oldKey]; ok {
			obj[newKey] = val
			delete(obj, oldKey)
		}
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
