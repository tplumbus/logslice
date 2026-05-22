package filter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/filter"
)

// TestFieldStats_WithPipeline verifies that FieldStats can be used alongside
// an existing Pipeline to collect statistics only on matching lines.
func TestFieldStats_WithPipeline_MatchingLines(t *testing.T) {
	lines := []string{
		buildLine("2024-01-15T10:00:00Z", "info", "started"),
		buildLine("2024-01-15T11:00:00Z", "warn", "slow query"),
		buildLine("2024-01-15T12:00:00Z", "info", "checkpoint"),
		buildLine("2024-01-15T13:00:00Z", "error", "crash"),
		buildLine("2024-01-15T14:00:00Z", "info", "recovered"),
	}

	tr, err := filter.ParseTimeRange("2024-01-15T10:30:00Z", "2024-01-15T13:30:00Z")
	if err != nil {
		t.Fatalf("ParseTimeRange: %v", err)
	}

	pipe := filter.NewPipeline(tr, nil)
	var out bytes.Buffer
	matched, err := pipe.Run(strings.NewReader(strings.Join(lines, "\n")), &out)
	if err != nil {
		t.Fatalf("pipeline run: %v", err)
	}

	fs := filter.NewFieldStats("level")
	for _, line := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		if line != "" {
			fs.Observe(line)
		}
	}

	if matched != 3 {
		t.Errorf("expected 3 matched lines, got %d", matched)
	}
	if fs.Counts["warn"] != 1 {
		t.Errorf("expected warn=1, got %d", fs.Counts["warn"])
	}
	if fs.Counts["info"] != 1 {
		t.Errorf("expected info=1, got %d", fs.Counts["info"])
	}
	if fs.Counts["error"] != 1 {
		t.Errorf("expected error=1, got %d", fs.Counts["error"])
	}
}

func TestFieldStats_NumericFieldValue(t *testing.T) {
	fs := filter.NewFieldStats("code")
	fs.Observe(`{"code":200,"path":"/health"}`)
	fs.Observe(`{"code":200,"path":"/api"}`)
	fs.Observe(`{"code":500,"path":"/api"}`)

	top := fs.TopN(1)
	if len(top) == 0 {
		t.Fatal("expected at least one result")
	}
	if top[0].Value != "200" {
		t.Errorf("expected top value '200', got '%s'", top[0].Value)
	}
	if top[0].Count != 2 {
		t.Errorf("expected count 2, got %d", top[0].Count)
	}
}
