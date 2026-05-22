package filter

import (
	"strings"
	"testing"
	"time"
)

func TestInvertFilter_WithPipeline_ExcludesMatchingLines(t *testing.T) {
	lines := []string{
		`{"ts":"2024-01-01T10:00:00Z","level":"error","msg":"disk full"}`,
		`{"ts":"2024-01-01T10:01:00Z","level":"info","msg":"startup"}`,
		`{"ts":"2024-01-01T10:02:00Z","level":"error","msg":"timeout"}`,
		`{"ts":"2024-01-01T10:03:00Z","level":"warn","msg":"slow query"}`,
	}

	q, err := ParseFieldQuery("level=error")
	if err != nil {
		t.Fatalf("ParseFieldQuery: %v", err)
	}
	inv, err := NewInvertFilter(q)
	if err != nil {
		t.Fatalf("NewInvertFilter: %v", err)
	}

	tr, _ := ParseTimeRange("", "")
	p := NewPipeline(tr, inv)

	var matched []string
	err = p.Run(lines, func(line string) error {
		matched = append(matched, line)
		return nil
	})
	if err != nil {
		t.Fatalf("pipeline Run: %v", err)
	}

	if len(matched) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(matched), matched)
	}
	for _, m := range matched {
		if strings.Contains(m, `"level":"error"`) {
			t.Errorf("unexpected error-level line in output: %s", m)
		}
	}
}

func TestInvertFilter_WithPipeline_DoubleInvert(t *testing.T) {
	lines := []string{
		`{"ts":"2024-01-01T10:00:00Z","level":"info","msg":"boot"}`,
		`{"ts":"2024-01-01T10:01:00Z","level":"debug","msg":"verbose"}`,
	}

	q, _ := ParseFieldQuery("level=info")
	inner, _ := NewInvertFilter(q)  // NOT info  → debug passes
	outer, _ := NewInvertFilter(inner) // NOT (NOT info) → info passes

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	tr := &TimeRange{From: &from, To: &to}
	p := NewPipeline(tr, outer)

	var matched []string
	_ = p.Run(lines, func(line string) error {
		matched = append(matched, line)
		return nil
	})

	if len(matched) != 1 {
		t.Fatalf("expected 1 line, got %d", len(matched))
	}
	if !strings.Contains(matched[0], `"level":"info"`) {
		t.Errorf("expected info line, got: %s", matched[0])
	}
}
