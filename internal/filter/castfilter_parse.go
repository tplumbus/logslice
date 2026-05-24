package filter

import (
	"fmt"
	"strings"
)

// ParseCastFlag parses a cast flag value in the form "field:type".
// Example: "status:int", "score:float", "active:bool", "code:string".
func ParseCastFlag(s string) (*CastFilter, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fmt.Errorf("castfilter: empty cast expression")
	}
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("castfilter: expected format field:type, got %q", s)
	}
	field := strings.TrimSpace(parts[0])
	typeName := strings.TrimSpace(parts[1])
	if field == "" {
		return nil, fmt.Errorf("castfilter: field name must not be empty")
	}
	if typeName == "" {
		return nil, fmt.Errorf("castfilter: type name must not be empty")
	}
	return NewCastFilter(field, typeName)
}
