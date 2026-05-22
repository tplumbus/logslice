package filter

import (
	"strings"
	"testing"
)

func TestRegexFilter_WithPipeline_FiltersCorrectly(t *testing.T) {
	lines := []string{
		`{"ts":"2024-01-10T10:00:00Z","level":"error","msg":"disk full"}`,
		`{"ts":"2024-01-10T10:01:00Z","level":"info","msg":"startup complete"}`,
		`{"ts":"2024-01-10T10:02:00Z","level":"warn","msg":"memory pressure"}`,
		`{"ts":"2024-01-10T10:03:00Z","level":"error","msg":"connection refused"}`,
	}

	rf, err := NewRegexFilter("level", `^(error|warn)$`)
	if err != nil {
		t.Fatalf("failed to create regex filter: %v", err)
	}

	var matched []string
	for _, line := range lines {
		if rf.MatchesLine(line) {
			matched = append(matched, line)
		}
	}

	if len(matched) != 3 {
		t.Fatalf("expected 3 matched lines, got %d", len(matched))
	}

	for _, line := range matched {
		if strings.Contains(line, `"level":"info"`) {
			t.Error("info line should not be in results")
		}
	}
}

func TestRegexFilter_WithPipeline_MessageSubstring(t *testing.T) {
	lines := []string{
		`{"ts":"2024-01-10T10:00:00Z","level":"error","msg":"read timeout on socket"}`,
		`{"ts":"2024-01-10T10:01:00Z","level":"error","msg":"write failed"}`,
		`{"ts":"2024-01-10T10:02:00Z","level":"warn","msg":"dial timeout"}`,
	}

	rf, err := NewRegexFilter("msg", `timeout`)
	if err != nil {
		t.Fatalf("failed to create regex filter: %v", err)
	}

	var matched []string
	for _, line := range lines {
		if rf.MatchesLine(line) {
			matched = append(matched, line)
		}
	}

	if len(matched) != 2 {
		t.Fatalf("expected 2 timeout lines, got %d", len(matched))
	}

	for _, line := range matched {
		if !strings.Contains(line, "timeout") {
			t.Errorf("matched line does not contain 'timeout': %s", line)
		}
	}
}
