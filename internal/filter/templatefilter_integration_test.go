package filter

import (
	"strings"
	"testing"
)

func TestTemplateFilter_WithPipeline_TransformsOutput(t *testing.T) {
	lines := []string{
		`{"level":"info","message":"started","ts":"2024-01-01T00:00:00Z"}`,
		`{"level":"error","message":"failed","ts":"2024-01-01T00:01:00Z"}`,
		`{"level":"warn","message":"retrying","ts":"2024-01-01T00:02:00Z"}`,
	}

	tf, err := NewTemplateFilter(`[{{.level}}] {{.message}}`)
	if err != nil {
		t.Fatalf("NewTemplateFilter: %v", err)
	}

	var results []string
	for _, line := range lines {
		if tf.MatchesLine(line) {
			results = append(results, tf.TransformLine(line))
		}
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[0] != "[info] started" {
		t.Errorf("line 0: got %q", results[0])
	}
	if results[1] != "[error] failed" {
		t.Errorf("line 1: got %q", results[1])
	}
}

func TestTemplateFilter_WithPipeline_CombinedWithFieldFilter(t *testing.T) {
	lines := []string{
		`{"level":"info","message":"ok"}`,
		`{"level":"error","message":"boom"}`,
		`{"level":"info","message":"also ok"}`,
	}

	ff, err := ParseFieldQuery("level=error")
	if err != nil {
		t.Fatalf("ParseFieldQuery: %v", err)
	}
	tf, err := NewTemplateFilter(`ALERT: {{.message}}`)
	if err != nil {
		t.Fatalf("NewTemplateFilter: %v", err)
	}

	var results []string
	for _, line := range lines {
		if ff.MatchesLine(line) {
			results = append(results, tf.TransformLine(line))
		}
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !strings.HasPrefix(results[0], "ALERT:") {
		t.Errorf("expected ALERT prefix, got %q", results[0])
	}
}
